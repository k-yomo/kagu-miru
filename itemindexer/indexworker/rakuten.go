package indexworker

import (
	"context"
	"fmt"
	"sort"

	"github.com/k-yomo/kagu-miru/internal/es"
	"github.com/k-yomo/kagu-miru/itemindexer/index"
	"github.com/k-yomo/kagu-miru/pkg/rakuten"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type RakutenWorker struct {
	itemIndexer            *index.ItemIndexer
	rakutenIchibaAPIClient *rakuten.IchibaClient
	logger                 *zap.Logger
}

func NewRakutenItemIndexWorker(indexer *index.ItemIndexer, rakutenIchibaAPIClient *rakuten.IchibaClient, logger *zap.Logger) *RakutenWorker {
	return &RakutenWorker{
		itemIndexer:            indexer,
		rakutenIchibaAPIClient: rakutenIchibaAPIClient,
		logger:                 logger,
	}
}

type RakutenWorkerOption struct {
	StartGenreID int
	MinPrice     int
}

// Run starts rakuten item indexing worker
func (r *RakutenWorker) Run(ctx context.Context, option *RakutenWorkerOption) error {
	furnitureGenreIDs, err := r.getFurnitureGenreIDs(ctx)
	if err != nil {
		return fmt.Errorf("getFurnitureGenreIDs: %w", err)
	}
	sort.Slice(furnitureGenreIDs, func(i, j int) bool {
		return furnitureGenreIDs[i] < furnitureGenreIDs[j]
	})

	curGenreIdx := 0
	curPage := 1
	curMinPrice := option.MinPrice

	if option.StartGenreID != 0 {
		for i, genreID := range furnitureGenreIDs {
			if genreID == option.StartGenreID {
				curGenreIdx = i
				break
			}
		}
	}

	rateLimiter := rate.NewLimiter(1, 1)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := rateLimiter.Wait(ctx); err != nil {
				return fmt.Errorf("rateLimiter.Wait: %w", err)
			}

			genreID := furnitureGenreIDs[curGenreIdx]
			searchItemRes, err := r.rakutenIchibaAPIClient.SearchItem(ctx, &rakuten.SearchItemParams{
				GenreID:  genreID,
				MinPrice: curMinPrice,
				Page:     curPage,
				SortType: rakuten.SearchItemSortTypeItemPriceAsc,
			})
			if err != nil {
				return fmt.Errorf("rakutenIchibaAPIClient.SearchItem: %w", err)
			}

			// when finished to get all items in the genre
			if searchItemRes.Count == 0 {
				curPage = 1
				curMinPrice = 0
				curGenreIdx++
				if curGenreIdx == len(furnitureGenreIDs)-1 {
					curGenreIdx = 0
				}
			}

			if len(searchItemRes.Items) > 0 {
				rakutenItems := make([]*rakuten.Item, 0, len(searchItemRes.Items))
				for _, item := range searchItemRes.Items {
					rakutenItems = append(rakutenItems, item.Item)
				}
				items, err := mapRakutenItemsToIndexItems(rakutenItems)
				if err != nil {
					return fmt.Errorf("mapRakutenItemsToIndexItems: %w", err)
				}

				if err := r.itemIndexer.BulkIndex(ctx, items); err != nil {
					return fmt.Errorf("itemIndexer.BulkIndex: %w", err)
				}

				r.logger.Info(fmt.Sprintf(
					"indexed %d items", len(items)),
					zap.Int("genreID", genreID),
					zap.Int("minPrice", curMinPrice),
					zap.Int("page", curPage),
				)
			}

			if curPage == searchItemRes.PageCount {
				curPage = 1
				// TODO: this implementation might miss items because the next item might be the same price
				curMinPrice = searchItemRes.Items[len(searchItemRes.Items)-1].Item.ItemPrice + 1
			} else {
				curPage++
			}
		}
	}
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

	return &es.Item{
		ID:             rakutenItem.ID(),
		Name:           rakutenItem.ItemName,
		Description:    rakutenItem.ItemCaption,
		Status:         status,
		SellingPageURL: rakutenItem.ItemURL,
		Price:          rakutenItem.ItemPrice,
		ImageURLs:      imageURLs,
		AverageRating:  rakutenItem.ReviewAverage,
		ReviewCount:    rakutenItem.ReviewCount,
		Platform:       es.PlatformRakuten,
	}, nil
}
