package indexworker

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/k-yomo/kagu-miru/pkg/jancode"

	"github.com/k-yomo/kagu-miru/internal/es"
	"github.com/k-yomo/kagu-miru/itemindexer/index"
	"github.com/k-yomo/kagu-miru/pkg/rakutenichiba"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type genreIndexWorker struct {
	itemIndexer            *index.ItemIndexer
	rakutenIchibaAPIClient *rakutenichiba.Client

	wg      *sync.WaitGroup
	pool    chan<- *genreIndexWorker
	genreID chan int

	genreMap map[string]*Genre

	logger *zap.Logger
}

type RakutenWorker struct {
	rakutenIchibaAPIClient *rakutenichiba.Client
	pool                   <-chan *genreIndexWorker
	workers                []*genreIndexWorker
	logger                 *zap.Logger

	wg *sync.WaitGroup
}

func NewRakutenItemWorker(indexer *index.ItemIndexer, rakutenIchibaAPIClient *rakutenichiba.Client, logger *zap.Logger) *RakutenWorker {
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
	furnitureGenres, err := r.getFurnitureGenres(ctx)
	if err != nil {
		return fmt.Errorf("getFurnitureGenre: %w", err)
	}

	genreMap := buildGenreMapFromGenres(furnitureGenres)
	for _, w := range r.workers {
		w.genreMap = genreMap
		w.start(ctx)
	}

	indexGenreIDs := make([]int, 0, len(furnitureGenres))
	for _, genre := range furnitureGenres {
		indexGenreIDs = append(indexGenreIDs, genre.ID)
	}
	sort.Slice(indexGenreIDs, func(i, j int) bool {
		return indexGenreIDs[i] < indexGenreIDs[j]
	})

	startGenreIdx := 0
	if option.StartGenreID != 0 {
		for i, genreID := range indexGenreIDs {
			if genreID == option.StartGenreID {
				startGenreIdx = i
				break
			}
		}
	}

	r.logger.Info(fmt.Sprintf("[start] indexing %d genre", len(indexGenreIDs[startGenreIdx:])))

	for _, genreID := range indexGenreIDs[startGenreIdx:] {
		r.wg.Add(1)
		(<-r.pool).genreID <- genreID
	}
	r.wg.Wait()

	r.logger.Info(fmt.Sprintf("[end] indexing %d genre", len(indexGenreIDs[startGenreIdx:])))
	return nil
}

// buildGenreMapFromGenre builds genreID => *Genre map by tracking
func buildGenreMapFromGenres(genres []*Genre) map[string]*Genre {
	queue := genres
	genreMap := map[string]*Genre{}
	for len(queue) != 0 {
		g := queue[0]
		queue = queue[1:]

		genreMap[strconv.Itoa(g.ID)] = g

		for _, child := range g.Children {
			queue = append(queue, child)
		}
	}
	return genreMap
}

func (r *RakutenWorker) getFurnitureGenres(ctx context.Context) ([]*Genre, error) {
	furnitureGenre, err := r.rakutenIchibaAPIClient.SearchGenre(ctx, rakutenichiba.GenreFurnitureID)
	if err != nil {
		return nil, fmt.Errorf("rakutenIchibaAPIClient.SearchGenre: %w", err)
	}
	genres := make([]*Genre, 0, len(furnitureGenre.Children))
	for _, genre := range furnitureGenre.Children {
		genre := &Genre{Genre: genre.Child}
		if err := r.setChildGenres(ctx, genre); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (r *RakutenWorker) setChildGenres(ctx context.Context, genre *Genre) error {
	time.Sleep((time.Duration(1000.0 / float64(r.rakutenIchibaAPIClient.ApplicationIDNum()))) * time.Millisecond)

	rakutenGenre, err := r.rakutenIchibaAPIClient.SearchGenre(ctx, strconv.Itoa(genre.ID))
	if err != nil {
		return fmt.Errorf("rakutenIchibaAPIClient.SearchGenre: %w", err)
	}
	genres := make([]*Genre, 0, len(genre.Children))
	for _, childGenre := range rakutenGenre.Children {
		cd := &Genre{
			Parent: genre,
			Genre:  childGenre.Child,
		}
		if err := r.setChildGenres(ctx, cd); err != nil {
			return err
		}
		genres = append(genres, cd)
	}
	genre.Children = genres
	return nil
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
				cursor := w.rakutenIchibaAPIClient.NewGenreItemCursor(genreID)
				for {
					if err := rateLimiter.Wait(ctx); err != nil {
						w.logger.Error("rateLimiter.Wait failed", zap.Error(err))
					}

					res, err := cursor.Next(ctx)
					if err == rakutenichiba.Done {
						w.logger.Info(fmt.Sprintf(
							"indexed all items in genre %d", genreID),
							zap.Int("genreID", genreID),
							zap.Int("total", totalIndexCount),
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

					if err := w.itemIndexer.BulkIndex(ctx, items); err != nil {
						// skip if indexing failed
						w.logger.Error("itemIndexer.BulkIndex failed", zap.Error(err))
					} else {
						totalIndexCount += len(items)
						if totalIndexCount%300 == 0 {
							w.logger.Info(fmt.Sprintf(
								"indexed %d items", totalIndexCount),
								zap.Int("genreID", genreID),
								zap.Int("minPrice", cursor.CurMinPrice()),
								zap.Int("page", cursor.CurPage()),
							)
						}
					}
				}

				w.wg.Done()
			}
		}
	}()
}

func mapRakutenItemsToIndexItems(rakutenItems []*rakutenichiba.Item, genreMap map[string]*Genre) ([]*es.Item, error) {
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

func mapRakutenItemToIndexItem(rakutenItem *rakutenichiba.Item, genre *Genre) (*es.Item, error) {
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
		IndexedAt:     time.Now().UnixMilli(),
	}, nil
}
