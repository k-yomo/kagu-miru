package indexworker

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/k-yomo/kagu-miru/pkg/jancode"

	"github.com/cenkalti/backoff/v4"
	"github.com/k-yomo/kagu-miru/internal/es"
	"github.com/k-yomo/kagu-miru/itemindexer/index"
	"github.com/k-yomo/kagu-miru/pkg/rakuten"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type genreIndexWorker struct {
	itemIndexer            *index.ItemIndexer
	rakutenIchibaAPIClient *rakuten.IchibaClient

	wg      *sync.WaitGroup
	pool    chan<- *genreIndexWorker
	genreID chan int

	logger *zap.Logger
}

type RakutenWorker struct {
	rakutenIchibaAPIClient *rakuten.IchibaClient
	pool                   <-chan *genreIndexWorker
	workers                []*genreIndexWorker
	logger                 *zap.Logger

	wg *sync.WaitGroup
}

func NewRakutenItemWorker(indexer *index.ItemIndexer, rakutenIchibaAPIClient *rakuten.IchibaClient, logger *zap.Logger) *RakutenWorker {
	wg := &sync.WaitGroup{}
	pool := make(chan *genreIndexWorker, rakutenIchibaAPIClient.ApplicationIDNum())
	workers := make([]*genreIndexWorker, 0, cap(pool))
	for i := 0; i < cap(pool); i++ {
		workers = append(workers, &genreIndexWorker{
			itemIndexer:            indexer,
			rakutenIchibaAPIClient: rakutenIchibaAPIClient,
			wg:                     wg,
			pool:                   pool,
			genreID:                make(chan int),
			logger:                 logger,
		})
	}
	return &RakutenWorker{
		rakutenIchibaAPIClient: rakutenIchibaAPIClient,
		wg:                     wg,
		pool:                   pool,
		workers:                workers,
		logger:                 logger,
	}
}

type RakutenWorkerOption struct {
	StartGenreID int
	MinPrice     int
}

// Run starts rakuten item indexing genreIndexWorker
func (r *RakutenWorker) Run(ctx context.Context, option *RakutenWorkerOption) error {
	for _, w := range r.workers {
		w.start(ctx)
	}

	furnitureGenreIDs, err := r.getFurnitureGenreIDs(ctx)
	if err != nil {
		return fmt.Errorf("getFurnitureGenreIDs: %w", err)
	}
	sort.Slice(furnitureGenreIDs, func(i, j int) bool {
		return furnitureGenreIDs[i] < furnitureGenreIDs[j]
	})

	startGenreIdx := 0
	if option.StartGenreID != 0 {
		for i, genreID := range furnitureGenreIDs {
			if genreID == option.StartGenreID {
				startGenreIdx = i
				break
			}
		}
	}

	r.logger.Info(fmt.Sprintf("[start] indexing %d genre", len(furnitureGenreIDs[startGenreIdx:])))

	for _, genreID := range furnitureGenreIDs[startGenreIdx:] {
		r.wg.Add(1)
		(<-r.pool).genreID <- genreID
	}
	r.wg.Wait()

	r.logger.Info(fmt.Sprintf("[end] indexing %d genre", len(furnitureGenreIDs[startGenreIdx:])))
	return nil
}

func (r *RakutenWorker) getFurnitureGenreIDs(ctx context.Context) ([]int, error) {
	furnitureGenre, err := r.rakutenIchibaAPIClient.SearchGenre(ctx, rakuten.GenreFurnitureID)
	if err != nil {
		return nil, fmt.Errorf("rakutenIchibaAPIClient.SearchGenre: %w", err)
	}
	genreIDs := make([]int, 0, len(furnitureGenre.Children))
	for _, genre := range furnitureGenre.Children {
		genreIDs = append(genreIDs, genre.Child.ID)
	}
	return genreIDs, nil
}

// We traverse all items in the given genre with following way
// due to Rakuten Ichiba API limitation(max 30 items at once, 100 page for the given condition)
// 1. get items in price ascending order
// 2. when we reach 100th page, set the last item's price to `minPrice` and fetch more 100 pages
// 3. when we get 0 items, it means we reached the end.
// Ideally, we want to reindex only updated items since we need to do full-reindex with the current approach.
// But currently we don't have a way to get item's updated time(API doesn't return it) and set `from` parameter for search
func (w *genreIndexWorker) start(ctx context.Context) {
	rateLimiter := rate.NewLimiter(1, 1)

	go func() {
		for {
			w.pool <- w

			select {
			case <-ctx.Done():
				return
			case genreID := <-w.genreID:
				totalIndexCount := 0
				curPage := 1
				curMinPrice := 0

				for {
					if err := rateLimiter.Wait(ctx); err != nil {
						w.logger.Error("rateLimiter.Wait failed", zap.Error(err))
					}

					searchItemParams := &rakuten.SearchItemParams{
						GenreID:  genreID,
						MinPrice: curMinPrice,
						Page:     curPage,
						SortType: rakuten.SearchItemSortTypeItemPriceAsc,
					}
					var searchItemRes *rakuten.SearchItemResponse
					b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
					err := backoff.Retry(func() error {
						var err error
						searchItemRes, err = w.rakutenIchibaAPIClient.SearchItem(ctx, searchItemParams)
						return err
					}, b)
					if err != nil {
						w.logger.Error("rakutenIchibaAPIClient.SearchItem failed", zap.Error(err), zap.Any("params", searchItemParams))
						break
					}

					// when finished to get all items in the genre
					rakutenItems := make([]*rakuten.Item, 0, len(searchItemRes.Items))
					for _, item := range searchItemRes.Items {
						rakutenItems = append(rakutenItems, item.Item)
					}
					items, err := mapRakutenItemsToIndexItems(rakutenItems)
					if err != nil {
						w.logger.Error("mapRakutenItemsToIndexItems failed", zap.Error(err))
					}

					if err := w.itemIndexer.BulkIndex(ctx, items); err != nil {
						// skip if indexing failed
						w.logger.Error("itemIndexer.BulkIndex failed", zap.Error(err))
					} else {
						totalIndexCount += len(items)
						if totalIndexCount%300 == 0 {
							w.logger.Info(fmt.Sprintf(
								"indexed %d items", totalIndexCount),
								zap.Int("genreID", genreID),
								zap.Int("minPrice", curMinPrice),
								zap.Int("page", curPage),
							)
						}
					}

					if curPage == searchItemRes.PageCount {
						if searchItemRes.PageCount < 100 {
							w.logger.Info(fmt.Sprintf(
								"indexed all items in genre %d", genreID),
								zap.Int("genreID", genreID),
								zap.Int("total", totalIndexCount),
							)
							break
						}
						nextPrice := searchItemRes.Items[len(searchItemRes.Items)-1].Item.ItemPrice
						if curMinPrice == nextPrice {
							nextPrice++
						}
						curPage = 1
						curMinPrice = nextPrice
					} else {
						curPage++
					}
				}

				w.wg.Done()
			}
		}
	}()
}

func mapRakutenItemsToIndexItems(rakutenItems []*rakuten.Item) ([]*es.Item, error) {
	items := make([]*es.Item, 0, len(rakutenItems))
	for _, rakutenItem := range rakutenItems {
		item, err := mapRakutenItemToIndexItem(rakutenItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func mapRakutenItemToIndexItem(rakutenItem *rakuten.Item) (*es.Item, error) {
	var status es.Status
	switch rakutenItem.Availability {
	case 0:
		status = es.StatusInactive
	case 1:
		status = es.StatusActive
	default:
		return nil, fmt.Errorf("unknown status %d, item: %v", rakutenItem.Availability, rakutenItem)
	}

	imageURLs := make([]string, 0, len(rakutenItem.MediumImageURLs))
	for _, mediumImage := range rakutenItem.MediumImageURLs {
		imageURLs = append(imageURLs, mediumImage.ImageURL)
	}

	genreID, err := strconv.Atoi(rakutenItem.GenreID)
	if err != nil {
		return nil, fmt.Errorf("invalid genreId '%d': %w", genreID, err)
	}

	janCode := jancode.ExtractJANCode(rakutenItem.ItemCaption)

	return &es.Item{
		ID:            es.ItemUniqueID(es.PlatformRakuten, rakutenItem.ID()),
		Name:          rakutenItem.ItemName,
		Description:   rakutenItem.ItemCaption,
		Status:        status,
		URL:           rakutenItem.ItemURL,
		AffiliateURL:  rakutenItem.AffiliateUrl,
		Price:         rakutenItem.ItemPrice,
		ImageURLs:     imageURLs,
		AverageRating: rakutenItem.ReviewAverage,
		ReviewCount:   rakutenItem.ReviewCount,
		GenreID:       genreID,
		TagIDs:        rakutenItem.TagIDs,
		JANCode:       janCode,
		Platform:      es.PlatformRakuten,
		IndexedAt:     time.Now().UnixMilli(),
	}, nil
}
