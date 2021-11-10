package yahoo_shopping

import (
	"context"
	"errors"

	"github.com/cenkalti/backoff/v4"
)

var Done = errors.New("DONE")

type CategoryItemCursor struct {
	shoppingClient *Client
	genreID        int
	curPage        int
	curMinPrice    int

	isDone bool
}

func (c *Client) NewCategoryItemCursor(genreID int) *CategoryItemCursor {
	return &CategoryItemCursor{
		shoppingClient: c,
		genreID:        genreID,
		curPage:        1,
		curMinPrice:    0,
	}
}

func (g *CategoryItemCursor) CurPage() int {
	return g.curPage
}

func (g *CategoryItemCursor) CurMinPrice() int {
	return g.curMinPrice
}

func (g *CategoryItemCursor) Next(ctx context.Context) (*SearchItemResponse, error) {
	if g.isDone {
		return nil, Done
	}
	searchItemParams := &SearchItemParams{
		CategoryID: g.genreID,
		PriceFrom:  g.curMinPrice,
		Page:       g.curPage,
		SortType:   SearchItemSortTypePriceAsc,
	}
	var searchItemRes *SearchItemResponse
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
	err := backoff.Retry(func() error {
		var err error
		searchItemRes, err = g.shoppingClient.SearchItem(ctx, searchItemParams)
		return err
	}, b)
	if err != nil {
		return nil, err
	}

	lastResultPosition := (searchItemRes.FirstResultPosition - 1) + searchItemRes.TotalResultsReturned
	if lastResultPosition == searchItemRes.TotalResultsAvailable {
		if searchItemRes.TotalResultsAvailable < maxResultTotalCount {
			g.isDone = true
		} else {
			nextPrice := searchItemRes.Hits[len(searchItemRes.Hits)-1].Price
			if g.curMinPrice == nextPrice {
				nextPrice++
			}
			g.curPage = 1
			g.curMinPrice = nextPrice
		}
	} else {
		g.curPage++
	}

	return searchItemRes, nil
}
