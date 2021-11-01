package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/kagu-miru/backend/internal/es"
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

	genreMap map[string]*rakutenichiba.Genre

	logger *zap.Logger
}

type worker struct {
	rakutenIchibaAPIClient *rakutenichiba.Client
	pool                   <-chan *genreItemsFetcher
	workers                []*genreItemsFetcher
	logger                 *zap.Logger

	wg *sync.WaitGroup
}

func newWorker(pubsubItemUpdateTopic *pubsub.Topic, rakutenIchibaAPIClient *rakutenichiba.Client, logger *zap.Logger) *worker {
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
		rakutenIchibaAPIClient: rakutenIchibaAPIClient,
		wg:                     wg,
		pool:                   pool,
		workers:                workers,
		logger:                 logger,
	}
}

type rakutenWorkerOption struct {
	StartGenreID int
	MinPrice     int
}

func (r *worker) run(ctx context.Context, option *rakutenWorkerOption) error {
	furnitureGenres, err := r.getFurnitureGenres(ctx)
	if err != nil {
		return fmt.Errorf("getFurnitureGenre: %w", err)
	}

	genreMap := buildGenreMapFromGenres(furnitureGenres)
	for _, w := range r.workers {
		w.genreMap = genreMap
		w.start(ctx)
	}

	fetchGenreIDs := make([]int, 0, len(furnitureGenres))
	for _, genre := range furnitureGenres {
		fetchGenreIDs = append(fetchGenreIDs, genre.ID)
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

// buildGenreMapFromGenre builds genreID => *Genre map by tracking
func buildGenreMapFromGenres(genres []*rakutenichiba.Genre) map[string]*rakutenichiba.Genre {
	queue := genres
	genreMap := map[string]*rakutenichiba.Genre{}
	for len(queue) != 0 {
		g := queue[0]
		queue = queue[1:]

		genreMap[strconv.Itoa(g.ID)] = g

		queue = append(queue, g.Children...)
	}
	return genreMap
}

func (r *worker) getFurnitureGenres(ctx context.Context) ([]*rakutenichiba.Genre, error) {
	furnitureGenre, err := r.rakutenIchibaAPIClient.SearchGenre(ctx, rakutenichiba.GenreFurnitureID)
	if err != nil {
		return nil, fmt.Errorf("rakutenIchibaAPIClient.SearchGenre: %w", err)
	}
	var genres []*rakutenichiba.Genre
	for _, child := range furnitureGenre.Children {
		genre, err := r.rakutenIchibaAPIClient.GetGenreWithAllChildren(ctx, strconv.Itoa(child.Child.ID))
		if err != nil {
			return nil, fmt.Errorf("rakutenIchibaAPIClient.GetGenreWithAllChildren: %w", err)
		}
		genres = append(genres, genre)
	}
	return genres, nil
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
				cursor := w.rakutenIchibaAPIClient.NewGenreItemCursor(genreID)
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
					items, err := mapRakutenItemsToIndexItems(rakutenItems, w.genreMap)
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

func mapRakutenItemsToIndexItems(rakutenItems []*rakutenichiba.Item, genreMap map[string]*rakutenichiba.Genre) ([]*es.Item, error) {
	items := make([]*es.Item, 0, len(rakutenItems))
	var errors []error
	for _, rakutenItem := range rakutenItems {
		genre, ok := genreMap[rakutenItem.GenreID]
		if !ok {
			errors = append(errors, fmt.Errorf("failed to get genre, item id: %s", rakutenItem.ID()))
			continue
		}
		item, err := mapRakutenItemToIndexItem(rakutenItem, genre)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		items = append(items, item)
	}
	return items, multierr.Combine(errors...)
}

func mapRakutenItemToIndexItem(rakutenItem *rakutenichiba.Item, genre *rakutenichiba.Genre) (*es.Item, error) {
	var status es.Status
	switch rakutenItem.Availability {
	case 0:
		status = es.StatusInactive
	case 1:
		status = es.StatusActive
	default:
		return nil, fmt.Errorf("unknown status %d, item id: %v", rakutenItem.Availability, rakutenItem.ID())
	}

	imageURLs := make([]string, 0, len(rakutenItem.MediumImageURLs))
	for _, mediumImage := range rakutenItem.MediumImageURLs {
		imageURLs = append(imageURLs, mediumImage.ImageURL)
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
		CategoryID:    rakutenItem.GenreID,
		CategoryIDs:   genre.GenreIDs(),
		CategoryNames: genre.GenreNames(),
		TagIDs:        rakutenItem.TagIDs,
		JANCode:       janCode,
		Platform:      es.PlatformRakuten,
		// TODO fix since this is actually not indexedAt but fetchedAt
		IndexedAt: time.Now().UnixMilli(),
	}, nil
}