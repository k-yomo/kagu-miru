package rakutenichiba

import (
	"context"
	"errors"

	"github.com/cenkalti/backoff/v4"
)

var Done = errors.New("DONE")

type GenreItemCursor struct {
	ichibaClient *Client
	genreID      int
	curPage      int
	curMinPrice  int

	isDone bool
}

func (c *Client) NewGenreItemCursor(genreID int) *GenreItemCursor {
	return &GenreItemCursor{
		ichibaClient: c,
		genreID:      genreID,
		curPage:      1,
		curMinPrice:  0,
	}
}

func (g *GenreItemCursor) CurPage() int {
	return g.curPage
}

func (g *GenreItemCursor) CurMinPrice() int {
	return g.curMinPrice
}

func (g *GenreItemCursor) Next(ctx context.Context) (*SearchItemResponse, error) {
	if g.isDone {
		return nil, Done
	}
	searchItemParams := &SearchItemParams{
		GenreID:  g.genreID,
		MinPrice: g.curMinPrice,
		Page:     g.curPage,
		SortType: SearchItemSortTypeItemPriceAsc,
	}
	var searchItemRes *SearchItemResponse
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
	err := backoff.Retry(func() error {
		var err error
		searchItemRes, err = g.ichibaClient.SearchItem(ctx, searchItemParams)
		return err
	}, b)
	if err != nil {
		return nil, err
	}

	if g.curPage == searchItemRes.PageCount {
		if searchItemRes.PageCount < 100 {
			g.isDone = true
		} else {
			nextPrice := searchItemRes.Items[len(searchItemRes.Items)-1].Item.ItemPrice
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
