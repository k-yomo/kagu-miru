package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/kagu-miru/backend/pkg/yahoo_shopping"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type categoryItemsFetcher struct {
	pubsubItemUpdateTopic  *pubsub.Topic
	yahooShoppingAPIClient *yahoo_shopping.Client

	wg         *sync.WaitGroup
	pool       chan<- *categoryItemsFetcher
	categoryID chan int

	logger *zap.Logger
}

type worker struct {
	yahooShoppingAPIClient *yahoo_shopping.Client
	pool                   <-chan *categoryItemsFetcher
	workers                []*categoryItemsFetcher
	logger                 *zap.Logger

	wg *sync.WaitGroup
}

func newWorker(pubsubItemUpdateTopic *pubsub.Topic, yahooShoppingAPIClient *yahoo_shopping.Client, logger *zap.Logger) *worker {
	wg := &sync.WaitGroup{}
	pool := make(chan *categoryItemsFetcher, yahooShoppingAPIClient.ApplicationIDNum())
	workers := make([]*categoryItemsFetcher, 0, cap(pool))
	for i := 0; i < cap(pool); i++ {
		workers = append(workers, &categoryItemsFetcher{
			pubsubItemUpdateTopic:  pubsubItemUpdateTopic,
			yahooShoppingAPIClient: yahooShoppingAPIClient,
			wg:                     wg,
			pool:                   pool,
			categoryID:             make(chan int),
			logger:                 logger,
		})
	}
	return &worker{
		yahooShoppingAPIClient: yahooShoppingAPIClient,
		wg:                     wg,
		pool:                   pool,
		workers:                workers,
		logger:                 logger,
	}
}

type yahooShoppingWorkerOption struct {
	StartCategoryID int
	MinPrice        int
}

func (r *worker) run(ctx context.Context, option *yahooShoppingWorkerOption) error {
	furnitureCategories, err := r.getFurnitureCategories(ctx)
	if err != nil {
		return fmt.Errorf("getFurnitureGenre: %w", err)
	}

	for _, w := range r.workers {
		w.start(ctx)
	}

	fetchCategoryIDs := make([]int, 0, len(furnitureCategories))
	for _, category := range furnitureCategories {
		fetchCategoryIDs = append(fetchCategoryIDs, category.ID)
	}
	sort.Slice(fetchCategoryIDs, func(i, j int) bool {
		return fetchCategoryIDs[i] < fetchCategoryIDs[j]
	})

	startCategoryIdx := 0
	if option.StartCategoryID != 0 {
		for i, categoryID := range fetchCategoryIDs {
			if categoryID == option.StartCategoryID {
				startCategoryIdx = i
				break
			}
		}
	}

	r.logger.Info(fmt.Sprintf("[start] fetching %d category", len(fetchCategoryIDs[startCategoryIdx:])))

	for _, categoryID := range fetchCategoryIDs[startCategoryIdx:] {
		r.wg.Add(1)
		(<-r.pool).categoryID <- categoryID
	}
	r.wg.Wait()

	r.logger.Info(fmt.Sprintf("[end] fetching %d category", len(fetchCategoryIDs[startCategoryIdx:])))
	return nil
}

func (r *worker) getFurnitureCategories(ctx context.Context) ([]*yahoo_shopping.Category, error) {
	category, err := r.yahooShoppingAPIClient.GetCategoryWithAllChildren(ctx, yahoo_shopping.CategoryFurnitureID)
	if err != nil {
		return nil, fmt.Errorf("yahooShoppingAPIClient.GetCategoryWithAllChildren: %w", err)
	}
	return category.Children, nil
}

// We traverse all items in the given category with following way
// due to Yahoo Shopping API limitation(max 30 items at once, 100 page for the given condition)
// 1. get items in price ascending order
// 2. when we reach 100th page, set the last item's price to `minPrice` and fetch more 100 pages
// 3. when we get 0 items, it means we reached the end.
// Ideally, we want to refetch only updated items since we need to do full-reindex with the current approach.
// But currently we don't have a way to get item's updated time(API doesn't return it) and set `from` parameter for search
func (w *categoryItemsFetcher) start(ctx context.Context) {
	rateLimiter := rate.NewLimiter(1, 1)

	go func() {
		for {
			w.pool <- w

			select {
			case <-ctx.Done():
				return
			case categoryID := <-w.categoryID:

				totalPublishedCount := 0
				cursor := w.yahooShoppingAPIClient.NewCategoryItemCursor(categoryID)
				for {
					if err := rateLimiter.Wait(ctx); err != nil {
						w.logger.Error("rateLimiter.Wait failed", zap.Error(err))
					}

					res, err := cursor.Next(ctx)
					if err == yahoo_shopping.Done {
						w.logger.Info(fmt.Sprintf(
							"fetched all items in category %d", categoryID),
							zap.Int("categoryID", categoryID),
							zap.Int("total", totalPublishedCount),
						)
						break
					}
					if err != nil {
						w.logger.Error("cursor.Next failed",
							zap.Error(err),
							zap.Int("categoryID", categoryID),
							zap.Int("minPrice", cursor.CurMinPrice()),
							zap.Int("page", cursor.CurPage()),
						)
						break
					}

					items, err := mapYahooShoppingItemsToIndexItems(res.Hits)
					if err != nil {
						w.logger.Error(
							"mapYahooShoppingItemsToIndexItems failed for some items",
							zap.Error(err),
							zap.Int("totalCount", len(res.Hits)),
							zap.Int("failedCount", len(res.Hits)-len(items)),
						)
					}

					wg := sync.WaitGroup{}
					var publishedCount int64
					for _, item := range items {
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
								return
							}
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
							zap.Int("categoryID", categoryID),
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

func mapYahooShoppingItemsToIndexItems(yahooShoppingItems []*yahoo_shopping.Item) ([]*xitem.Item, error) {
	items := make([]*xitem.Item, 0, len(yahooShoppingItems))
	var errors []error
	for _, yahooShoppingItem := range yahooShoppingItems {
		item, err := mapYahooShoppingItemToIndexItem(yahooShoppingItem)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		items = append(items, item)
	}
	return items, multierr.Combine(errors...)
}

func mapYahooShoppingItemToIndexItem(yahooShoppingItem *yahoo_shopping.Item) (*xitem.Item, error) {
	var status xitem.Status
	if yahooShoppingItem.InStock {
		status = xitem.StatusActive
	} else {
		status = xitem.StatusInactive
	}

	// TODO: convert to kagu-miru category id
	// categoryID := strconv.Itoa(yahooShoppingItem.GenreCategory.Id)

	return &xitem.Item{
		ID:            xitem.ItemUniqueID(xitem.PlatformYahooShopping, yahooShoppingItem.Code),
		Name:          yahooShoppingItem.Name,
		Description:   yahooShoppingItem.Description,
		Status:        status,
		URL:           yahooShoppingItem.Url,
		AffiliateURL:  yahooShoppingItem.Url,
		Price:         yahooShoppingItem.Price,
		ImageURLs:     []string{yahooShoppingItem.Image.Medium, yahooShoppingItem.Image.Small},
		AverageRating: yahooShoppingItem.Review.Rate,
		ReviewCount:   yahooShoppingItem.Review.Count,
		CategoryID:    "101859", // その他, TODO: Replace with appropriate category id
		// CategoryID:    categoryID,
		// CategoryIDs:
		// CategoryNames:
		// TagIDs:
		JANCode:  yahooShoppingItem.JanCode,
		Platform: xitem.PlatformYahooShopping,
	}, nil
}
