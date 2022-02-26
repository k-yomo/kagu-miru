package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
	"github.com/k-yomo/kagu-miru/backend/item_fetcher"
	"github.com/k-yomo/kagu-miru/backend/pkg/amazon"
	"github.com/utekaravinash/gopaapi5/api"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type browseNodeItemsFetcher struct {
	pubsubItemUpdateTopic *pubsub.Topic
	amazonAPIClient       *amazon.Client

	wg           *sync.WaitGroup
	pool         chan<- *browseNodeItemsFetcher
	browseNodeID chan string

	browseNodeIDItemCategoryMap map[string]*xspanner.ItemCategoryWithParent

	logger *zap.Logger
}

type worker struct {
	amazonAPIClient *amazon.Client
	spannerClient   *spanner.Client
	pool            <-chan *browseNodeItemsFetcher
	workers         []*browseNodeItemsFetcher
	logger          *zap.Logger

	wg *sync.WaitGroup
}

func newWorker(pubsubItemUpdateTopic *pubsub.Topic, spannerClient *spanner.Client, amazonAPIClient *amazon.Client, logger *zap.Logger) *worker {
	wg := &sync.WaitGroup{}
	pool := make(chan *browseNodeItemsFetcher, 1)
	workers := make([]*browseNodeItemsFetcher, 0, cap(pool))
	for i := 0; i < cap(pool); i++ {
		workers = append(workers, &browseNodeItemsFetcher{
			pubsubItemUpdateTopic: pubsubItemUpdateTopic,
			amazonAPIClient:       amazonAPIClient,
			wg:                    wg,
			pool:                  pool,
			browseNodeID:          make(chan string),
			logger:                logger,
		})
	}
	return &worker{
		spannerClient:   spannerClient,
		amazonAPIClient: amazonAPIClient,
		wg:              wg,
		pool:            pool,
		workers:         workers,
		logger:          logger,
	}
}

type amazonWorkerOption struct {
	StartBrowseNodeID string
}

func (r *worker) run(ctx context.Context, option *amazonWorkerOption) error {
	browseNodeIDItemCategoryMap, err := r.getBrowseNodeIDItemCategoryMap(ctx)
	if err != nil {
		return fmt.Errorf("getBrowseNodeIDItemCategoryMap: %w", err)
	}
	for _, w := range r.workers {
		w.browseNodeIDItemCategoryMap = browseNodeIDItemCategoryMap
		w.start(ctx)
	}

	amazonBrowseNodes, err := xspanner.GetAllAmazonBrowseNodes(ctx, r.spannerClient)
	if err != nil {
		return fmt.Errorf("xspanner.GetAllAmazonItemBrowseNodes: %w", err)
	}

	// TODO: use bottom level browse node to narrow down the search result for each search
	var fetchBrowseNodes []string
	for _, browseNode := range amazonBrowseNodes {
		if browseNode.Level == 0 {
			fetchBrowseNodes = append(fetchBrowseNodes, browseNode.ID)
		}
	}
	sort.Slice(fetchBrowseNodes, func(i, j int) bool {
		return fetchBrowseNodes[i] < fetchBrowseNodes[j]
	})

	startBrowseNodeIdx := 0
	if option.StartBrowseNodeID != "" {
		for i, browseNodeID := range fetchBrowseNodes {
			if browseNodeID == option.StartBrowseNodeID {
				startBrowseNodeIdx = i
				break
			}
		}
	}

	r.logger.Info(fmt.Sprintf("[start] fetching %d browseNode", len(fetchBrowseNodes[startBrowseNodeIdx:])))

	for _, browseNodeID := range fetchBrowseNodes[startBrowseNodeIdx:] {
		r.wg.Add(1)
		(<-r.pool).browseNodeID <- browseNodeID
	}
	r.wg.Wait()

	r.logger.Info(fmt.Sprintf("[end] fetching %d browseNode", len(fetchBrowseNodes[startBrowseNodeIdx:])))
	return nil
}

func (r *worker) getBrowseNodeIDItemCategoryMap(ctx context.Context) (map[string]*xspanner.ItemCategoryWithParent, error) {
	amazonItemBrowseNodes, err := xspanner.GetAllAmazonBrowseNodes(ctx, r.spannerClient)
	if err != nil {
		return nil, fmt.Errorf("xspanner.GetAllAmazonItemBrowseNodes: %w", err)
	}
	browseNodeIDItemCategoryIDMap := make(map[string]string)
	for _, browseNode := range amazonItemBrowseNodes {
		browseNodeIDItemCategoryIDMap[browseNode.ID] = browseNode.ItemCategoryID
	}

	itemCategoriesWithParent, err := xspanner.GetAllActiveItemCategoriesWithParent(ctx, r.spannerClient)
	if err != nil {
		return nil, fmt.Errorf("xspanner.GetAllActiveItemCategoriesWithParent: %w", err)
	}
	itemCategoryMap := make(map[string]*xspanner.ItemCategoryWithParent)
	for _, itemCategory := range itemCategoriesWithParent {
		itemCategoryMap[itemCategory.ID] = itemCategory
	}

	browseNodeIDItemCategoryMap := make(map[string]*xspanner.ItemCategoryWithParent)
	for browseNodeID, itemCategoryID := range browseNodeIDItemCategoryIDMap {
		browseNodeIDItemCategoryMap[browseNodeID] = itemCategoryMap[itemCategoryID]
	}

	return browseNodeIDItemCategoryMap, nil
}

// We traverse all items in the given browseNode with following way
// due to Amazon Ichiba API limitation(max 30 items at once, 100 page for the given condition)
// 1. get items in price ascending order
// 2. when we reach 100th page, set the last item's price to `minPrice` and fetch more 100 pages
// 3. when we get 0 items, it means we reached the end.
// Ideally, we want to refetch only updated items since we need to do full-reindex with the current approach.
// But currently we don't have a way to get item's updated time(API doesn't return it) and set `from` parameter for search
func (w *browseNodeItemsFetcher) start(ctx context.Context) {
	rateLimiter := rate.NewLimiter(1, 1)

	go func() {
		for {
			w.pool <- w

			select {
			case <-ctx.Done():
				return
			case browseNodeID := <-w.browseNodeID:

				totalPublishedCount := 0
				cursor := w.amazonAPIClient.NewBrowseNodeItemCursor(browseNodeID, item_fetcher.MinFetchItemPrice, item_fetcher.MaxFetchItemPrice)
				for {
					if err := rateLimiter.Wait(ctx); err != nil {
						w.logger.Error("rateLimiter.Wait failed", zap.Error(err))
					}

					amazonItems, err := cursor.Next(ctx)
					if err == amazon.Done {
						w.logger.Info(fmt.Sprintf(
							"fetched all amazonItems in browseNode %s", browseNodeID),
							zap.String("browseNodeID", browseNodeID),
							zap.Int("total", totalPublishedCount),
						)
						break
					}
					if err != nil {
						w.logger.Error("cursor.Next failed",
							zap.Error(err),
							zap.String("browseNodeID", browseNodeID),
							zap.Int("minPrice", cursor.CurMinPrice()),
							zap.Int("page", cursor.CurPage()),
						)
						break
					}

					items, err := mapAmazonItemsToIndexItems(amazonItems, w.browseNodeIDItemCategoryMap)
					if err != nil {
						w.logger.Error(
							"mapAmazonItemsToIndexItems failed for some amazonItems",
							zap.Error(err),
							zap.Int("totalCount", len(items)),
						)
					}

					wg := sync.WaitGroup{}
					var publishedCount int64
					for _, item := range items {
						if !item.IsIndexable() {
							continue
						}

						item := item
						wg.Add(1)
						go func() {
							defer wg.Done()

							itemJSON, err := json.Marshal(item)
							if err != nil {
								w.logger.Error(
									"json.Marshal item failed",
									zap.Error(err),
									zap.Any("item", item),
								)
							}
							res := w.pubsubItemUpdateTopic.Publish(ctx, &pubsub.Message{
								Data:        itemJSON,
								OrderingKey: item_fetcher.ItemOrderingKey(item),
							})
							if _, err := res.Get(ctx); err != nil {
								w.logger.Error("publish item update failed",
									zap.Error(err),
									zap.String("itemId", item.ID),
								)
								return
							}
							atomic.AddInt64(&publishedCount, 1)
						}()
					}
					wg.Wait()

					totalPublishedCount += int(publishedCount)
					if totalPublishedCount%300 == 0 {
						w.logger.Info(fmt.Sprintf(
							"published %d items", totalPublishedCount),
							zap.String("browseNodeID", browseNodeID),
							zap.Int("minPrice", cursor.CurMinPrice()),
							zap.Int("page", cursor.CurPage()),
						)
					}
				}

				w.wg.Done()
			}
		}
	}()
}

func mapAmazonItemsToIndexItems(
	amazonItems []api.Item,
	browseNodeIDItemCategoryMap map[string]*xspanner.ItemCategoryWithParent,
) ([]*xitem.Item, error) {
	items := make([]*xitem.Item, 0, len(amazonItems))
	var errors []error
	for _, amazonItem := range amazonItems {
		if amazonItem.ItemInfo.ProductInfo.IsAdultProduct.DisplayValue {
			continue
		}
		browseNodeID := amazonItem.BrowseNodeInfo.BrowseNodes[0].Id
		itemCategory, ok := browseNodeIDItemCategoryMap[browseNodeID]
		if !ok {
			errors = append(errors, fmt.Errorf("failed to get itemCategory, item ASIN: %s", amazonItem.ASIN))
			continue
		}
		item, err := mapAmazonItemToIndexItem(amazonItem, itemCategory)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		items = append(items, item)
	}
	return items, multierr.Combine(errors...)
}

func mapAmazonItemToIndexItem(
	amazonItem api.Item,
	itemCategory *xspanner.ItemCategoryWithParent,
) (*xitem.Item, error) {
	listing := amazonItem.Offers.Listings[0]
	var status xitem.Status
	switch listing.Availability.Type {
	case "Now", "Available":
		status = xitem.StatusActive
	case "IncludeOutOfStock":
		status = xitem.StatusInactive
	default:
		return nil, fmt.Errorf("unknown type %s, item ASIN: %v", listing.Availability.Type, amazonItem.ASIN)
	}

	imageURLs := make([]string, 0, 1+len(amazonItem.Images.Variants)) // primary + variants
	imageURLs = append(imageURLs, amazonItem.Images.Primary.Large.URL)
	for _, variantImage := range amazonItem.Images.Variants {
		imageURLs = append(imageURLs, variantImage.Large.URL)
	}

	// janCode := jancode.ExtractJANCode(amazonItem.ItemInfo)
	// if janCode == "" {
	// 	janCode = jancode.ExtractJANCode(amazonItem.ItemCaption)
	// }

	return &xitem.Item{
		ID:   xitem.ItemUniqueID(xitem.PlatformAmazon, amazonItem.ASIN),
		Name: amazonItem.ItemInfo.Title.DisplayValue,
		// Description:   amazonItem.ItemInfo.ProductInfo.,
		Status:       status,
		URL:          amazonItem.DetailPageURL,
		AffiliateURL: amazonItem.DetailPageURL,
		Price:        int(listing.Price.Amount),
		ImageURLs:    imageURLs,
		// rating and review count are not available in PA-API5
		// AverageRating: amazonItem,
		// ReviewCount:   amazonItem,
		CategoryID:    itemCategory.ID,
		CategoryIDs:   itemCategory.CategoryIDs(),
		CategoryNames: itemCategory.CategoryNames(),
		BrandName:     amazonItem.ItemInfo.ByLineInfo.Brand.DisplayValue,
		Colors:        []string{amazonItem.ItemInfo.ProductInfo.Color.DisplayValue},
		WidthRange:    mapLengthToItemIntRange(amazonItem.ItemInfo.ProductInfo.ItemDimensions.Width),
		DepthRange:    mapLengthToItemIntRange(amazonItem.ItemInfo.ProductInfo.ItemDimensions.Length),
		HeightRange:   mapLengthToItemIntRange(amazonItem.ItemInfo.ProductInfo.ItemDimensions.Height),
		// JANCode:       janCode,
		Platform: xitem.PlatformAmazon,
	}, nil
}

func mapLengthToItemIntRange(uba api.UnitBasedAttribute) *xitem.IntRange {
	length := 0
	switch uba.Unit {
	case "インチ":
		length = int(math.Round(float64(uba.DisplayValue) * 2.54))
	case "センチメートル":
		length = int(math.Round(float64(uba.DisplayValue)))
	case "メートル":
		length = int(math.Round(float64(uba.DisplayValue) * 100))
	default:
		return nil
	}
	if length == 0 {
		return nil
	}
	return xitem.NewIntRange(length, &length)
}
