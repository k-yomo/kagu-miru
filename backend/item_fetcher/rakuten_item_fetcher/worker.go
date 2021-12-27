package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/k-yomo/kagu-miru/backend/item_fetcher"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/kagu-miru/backend/pkg/jancode"
	"github.com/k-yomo/kagu-miru/backend/pkg/rakutenichiba"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type genreItemsFetcher struct {
	pubsubItemUpdateTopic  *pubsub.Topic
	rakutenIchibaAPIClient *rakutenichiba.Client

	wg      *sync.WaitGroup
	pool    chan<- *genreItemsFetcher
	genreID chan int

	genreIDItemCategoryMap map[int]*xspanner.ItemCategoryWithParent
	tagMap                 map[int]*xspanner.RakutenTag

	logger *zap.Logger
}

type worker struct {
	rakutenIchibaAPIClient *rakutenichiba.Client
	spannerClient          *spanner.Client
	pool                   <-chan *genreItemsFetcher
	workers                []*genreItemsFetcher
	logger                 *zap.Logger

	wg *sync.WaitGroup
}

func newWorker(pubsubItemUpdateTopic *pubsub.Topic, spannerClient *spanner.Client, rakutenIchibaAPIClient *rakutenichiba.Client, logger *zap.Logger) *worker {
	wg := &sync.WaitGroup{}
	pool := make(chan *genreItemsFetcher, rakutenIchibaAPIClient.ApplicationIDNum())
	workers := make([]*genreItemsFetcher, 0, cap(pool))
	for i := 0; i < cap(pool); i++ {
		workers = append(workers, &genreItemsFetcher{
			pubsubItemUpdateTopic:  pubsubItemUpdateTopic,
			rakutenIchibaAPIClient: rakutenIchibaAPIClient,
			wg:                     wg,
			pool:                   pool,
			genreID:                make(chan int),
			logger:                 logger,
		})
	}
	return &worker{
		spannerClient:          spannerClient,
		rakutenIchibaAPIClient: rakutenIchibaAPIClient,
		wg:                     wg,
		pool:                   pool,
		workers:                workers,
		logger:                 logger,
	}
}

type rakutenWorkerOption struct {
	StartGenreID int
}

func (r *worker) run(ctx context.Context, option *rakutenWorkerOption) error {
	genreIDItemCategoryMap, err := r.getGenreIDItemCategoryMap(ctx)
	if err != nil {
		return fmt.Errorf("getGenreIDItemCategoryMap: %w", err)
	}
	tagMap, err := r.getTagMap(ctx)
	if err != nil {
		return fmt.Errorf("getTagMap: %w", err)
	}
	for _, w := range r.workers {
		w.genreIDItemCategoryMap = genreIDItemCategoryMap
		w.tagMap = tagMap
		w.start(ctx)
	}

	rakutenItemGenres, err := xspanner.GetAllRakutenItemGenres(ctx, r.spannerClient)
	if err != nil {
		return fmt.Errorf("xspanner.GetAllRakutenItemGenres: %w", err)
	}
	var fetchGenreIDs []int
	for _, genre := range rakutenItemGenres {
		if genre.Level == 0 {
			fetchGenreIDs = append(fetchGenreIDs, int(genre.ID))
		}
	}
	sort.Slice(fetchGenreIDs, func(i, j int) bool {
		return fetchGenreIDs[i] < fetchGenreIDs[j]
	})

	startGenreIdx := 0
	if option.StartGenreID != 0 {
		for i, genreID := range fetchGenreIDs {
			if genreID == option.StartGenreID {
				startGenreIdx = i
				break
			}
		}
	}

	r.logger.Info(fmt.Sprintf("[start] fetching %d genre", len(fetchGenreIDs[startGenreIdx:])))

	for _, genreID := range fetchGenreIDs[startGenreIdx:] {
		r.wg.Add(1)
		(<-r.pool).genreID <- genreID
	}
	r.wg.Wait()

	r.logger.Info(fmt.Sprintf("[end] fetching %d genre", len(fetchGenreIDs[startGenreIdx:])))
	return nil
}

func (r *worker) getTagMap(ctx context.Context) (map[int]*xspanner.RakutenTag, error) {
	tags, err := xspanner.GetAllRakutenTags(ctx, r.spannerClient)
	if err != nil {
		return nil, fmt.Errorf("xspanner.GetAllRakutenTags: %w", err)
	}

	tagMap := make(map[int]*xspanner.RakutenTag, len(tags))
	for _, tag := range tags {
		tagMap[int(tag.ID)] = tag
	}

	return tagMap, nil
}

func (r *worker) getGenreIDItemCategoryMap(ctx context.Context) (map[int]*xspanner.ItemCategoryWithParent, error) {
	rakutenItemGenres, err := xspanner.GetAllRakutenItemGenres(ctx, r.spannerClient)
	if err != nil {
		return nil, fmt.Errorf("xspanner.GetAllRakutenItemGenres: %w", err)
	}
	genreIDItemCategoryIDMap := make(map[int]string)
	for _, genre := range rakutenItemGenres {
		genreIDItemCategoryIDMap[int(genre.ID)] = genre.ItemCategoryID
	}

	itemCategoriesWithParent, err := xspanner.GetAllItemCategoriesWithParent(ctx, r.spannerClient)
	if err != nil {
		return nil, fmt.Errorf("xspanner.GetAllItemCategoriesWithParent: %w", err)
	}
	itemCategoryMap := make(map[string]*xspanner.ItemCategoryWithParent)
	for _, itemCategory := range itemCategoriesWithParent {
		itemCategoryMap[itemCategory.ID] = itemCategory
	}

	genreIDItemCategoryMap := make(map[int]*xspanner.ItemCategoryWithParent)
	for genreID, itemCategoryID := range genreIDItemCategoryIDMap {
		genreIDItemCategoryMap[genreID] = itemCategoryMap[itemCategoryID]
	}

	return genreIDItemCategoryMap, nil
}

// We traverse all items in the given genre with following way
// due to Rakuten Ichiba API limitation(max 30 items at once, 100 page for the given condition)
// 1. get items in price ascending order
// 2. when we reach 100th page, set the last item's price to `minPrice` and fetch more 100 pages
// 3. when we get 0 items, it means we reached the end.
// Ideally, we want to refetch only updated items since we need to do full-reindex with the current approach.
// But currently we don't have a way to get item's updated time(API doesn't return it) and set `from` parameter for search
func (w *genreItemsFetcher) start(ctx context.Context) {
	rateLimiter := rate.NewLimiter(1, 1)

	go func() {
		for {
			w.pool <- w

			select {
			case <-ctx.Done():
				return
			case genreID := <-w.genreID:

				totalPublishedCount := 0
				cursor := w.rakutenIchibaAPIClient.NewGenreItemCursor(genreID, item_fetcher.MinFetchItemPrice, item_fetcher.MaxFetchItemPrice)
				for {
					if err := rateLimiter.Wait(ctx); err != nil {
						w.logger.Error("rateLimiter.Wait failed", zap.Error(err))
					}

					res, err := cursor.Next(ctx)
					if err == rakutenichiba.Done {
						w.logger.Info(fmt.Sprintf(
							"fetched all items in genre %d", genreID),
							zap.Int("genreID", genreID),
							zap.Int("total", totalPublishedCount),
						)
						break
					}
					if err != nil {
						w.logger.Error("cursor.Next failed",
							zap.Error(err),
							zap.Int("genreID", genreID),
							zap.Int("minPrice", cursor.CurMinPrice()),
							zap.Int("page", cursor.CurPage()),
						)
						break
					}

					rakutenItems := make([]*rakutenichiba.Item, 0, len(res.Items))
					for _, item := range res.Items {
						rakutenItems = append(rakutenItems, item.Item)
					}
					items, err := mapRakutenItemsToIndexItems(rakutenItems, w.genreIDItemCategoryMap, w.tagMap)
					if err != nil {
						w.logger.Error(
							"mapRakutenItemsToIndexItems failed for some items",
							zap.Error(err),
							zap.Int("totalCount", len(rakutenItems)),
							zap.Int("failedCount", len(rakutenItems)-len(items)),
						)
					}

					wg := sync.WaitGroup{}
					var publishedCount int64
					for _, item := range items {
						item := item
						itemJSON, err := json.Marshal(item)
						if err != nil {
							w.logger.Error(
								"json.Marshal item failed",
								zap.Error(err),
								zap.Any("item", item),
							)
							continue
						}

						wg.Add(1)
						go func() {
							defer wg.Done()

							res := w.pubsubItemUpdateTopic.Publish(ctx, &pubsub.Message{
								Data: itemJSON,
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
							zap.Int("genreID", genreID),
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

func mapRakutenItemsToIndexItems(
	rakutenItems []*rakutenichiba.Item,
	genreIDItemCategoryMap map[int]*xspanner.ItemCategoryWithParent,
	tagMap map[int]*xspanner.RakutenTag,
) ([]*xitem.Item, error) {
	items := make([]*xitem.Item, 0, len(rakutenItems))
	var errors []error
	for _, rakutenItem := range rakutenItems {
		genreID, err := strconv.Atoi(rakutenItem.GenreID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to convert genre id '%s': %w", rakutenItem.GenreID, err))
			continue
		}
		itemCategory, ok := genreIDItemCategoryMap[genreID]
		if !ok {
			errors = append(errors, fmt.Errorf("failed to get itemCategory, item id: %s", rakutenItem.ItemCode))
			continue
		}
		item, err := mapRakutenItemToIndexItem(rakutenItem, itemCategory, tagMap)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		items = append(items, item)
	}
	return items, multierr.Combine(errors...)
}

func mapRakutenItemToIndexItem(
	rakutenItem *rakutenichiba.Item,
	itemCategory *xspanner.ItemCategoryWithParent,
	tagMap map[int]*xspanner.RakutenTag,
) (*xitem.Item, error) {
	var status xitem.Status
	switch rakutenItem.Availability {
	case 0:
		status = xitem.StatusInactive
	case 1:
		status = xitem.StatusActive
	default:
		return nil, fmt.Errorf("unknown status %d, item id: %v", rakutenItem.Availability, rakutenItem.ItemCode)
	}

	imageURLs := make([]string, 0, len(rakutenItem.MediumImageURLs))
	for _, mediumImage := range rakutenItem.MediumImageURLs {
		imageURLs = append(imageURLs, mediumImage.ImageURL)
	}

	janCode := jancode.ExtractJANCode(rakutenItem.ItemCaption)

	metadata := extractMetadataFromTags(rakutenItem.TagIDs, tagMap)

	return &xitem.Item{
		ID:            xitem.ItemUniqueID(xitem.PlatformRakuten, rakutenItem.ItemCode),
		Name:          rakutenItem.ItemName,
		Description:   rakutenItem.ItemCaption,
		Status:        status,
		URL:           rakutenItem.ItemURL,
		AffiliateURL:  rakutenItem.AffiliateUrl,
		Price:         rakutenItem.ItemPrice,
		ImageURLs:     imageURLs,
		AverageRating: rakutenItem.ReviewAverage,
		ReviewCount:   rakutenItem.ReviewCount,
		CategoryID:    itemCategory.ID,
		CategoryIDs:   itemCategory.CategoryIDs(),
		CategoryNames: itemCategory.CategoryNames(),
		BrandName:     metadata.brandName,
		Colors:        metadata.colors,
		WidthRange:    mapIntRangeToItemIntRange(metadata.widthRange),
		DepthRange:    mapIntRangeToItemIntRange(metadata.depthRange),
		HeightRange:   mapIntRangeToItemIntRange(metadata.heightRange),
		TagIDs:        rakutenItem.TagIDs,
		JANCode:       janCode,
		Platform:      xitem.PlatformRakuten,
	}, nil
}

type itemMetadata struct {
	brandName   string
	colors      []string
	widthRange  *intRange
	depthRange  *intRange
	heightRange *intRange
}

type intRange struct {
	Gte int
	Lte *int
}

func mapIntRangeToItemIntRange(r *intRange) *xitem.IntRange {
	if r == nil {
		return nil
	}
	return xitem.NewIntRange(r.Gte, r.Lte)
}

func extractMetadataFromTags(tagIDs []int, tagMap map[int]*xspanner.RakutenTag) *itemMetadata {
	metadata := itemMetadata{}
	for _, tagID := range tagIDs {
		tag, ok := tagMap[tagID]
		if !ok {
			continue
		}

		switch tag.TagGroupID {
		case xspanner.TagGroupIDBrand:
			metadata.brandName = tag.Name
		case xspanner.TagGroupIDColor:
			metadata.colors = append(metadata.colors, tag.Name)
		case xspanner.TagGroupIDWidth:
			const width0To19ID = 1000483
			metadata.widthRange = getDimensionRangeByTagID(tag.ID, width0To19ID)
		case xspanner.TagGroupIDDepth:
			const depth0To19ID = 1000503
			metadata.depthRange = getDimensionRangeByTagID(tag.ID, depth0To19ID)
		case xspanner.TagGroupIDHeight:
			const height0To19ID = 1000523
			metadata.heightRange = getDimensionRangeByTagID(tag.ID, height0To19ID)
		}
	}

	return &metadata
}

// getDimensionRangeByTagID gets int range from dimension tag id (width, depth or height)
// dimension tag id must be in within 20 consecutive id starting from 1: ~19cm up to 20: 200cm ~
func getDimensionRangeByTagID(tagID int64, tag0To19ID int64) *intRange {
	const dimensionGteLimit = 200
	curTagID := tag0To19ID
	gte, lte := 0, 19
	for {
		// not found
		if gte > dimensionGteLimit {
			return nil
		}
		if tagID == curTagID {
			if gte == dimensionGteLimit {
				return &intRange{Gte: gte, Lte: nil}
			}
			return &intRange{Gte: gte, Lte: &lte}
		}
		// since the number is increase as 0 ~ 19, 20 ~ 29, 30 ~ 39
		if curTagID == tag0To19ID {
			gte += 20
		} else {
			gte += 10
		}
		curTagID++
		lte += 10
	}
}
