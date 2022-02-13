package yahoo_shopping

import (
	"context"
	"errors"

	"github.com/cenkalti/backoff/v4"
)

var Done = errors.New("DONE")

type CategoryItemCursor struct {
	shoppingClient *Client

	categoryID  int
	curPage     int
	curMinPrice int
	maxPrice    int

	isDone bool
}

func (c *Client) NewCategoryItemCursor(categoryID int, minPrice int, maxPrice int) *CategoryItemCursor {
	return &CategoryItemCursor{
		shoppingClient: c,
		categoryID:     categoryID,
		curPage:        1,
		curMinPrice:    minPrice,
		maxPrice:       maxPrice,
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
		CategoryID: g.categoryID,
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

	for i, hit := range searchItemRes.Hits {
		if hit.Price > g.maxPrice {
			searchItemRes.Hits = searchItemRes.Hits[:i]
			g.isDone = true
			return searchItemRes, nil
		}
	}

	if noMoreItems := searchItemRes.TotalResultsReturned == searchItemRes.TotalResultsAvailable; noMoreItems {
		g.isDone = true
	} else if g.curPage*maxResultCount >= maxResultTotalCount {
		nextPrice := searchItemRes.Hits[len(searchItemRes.Hits)-1].Price
		if g.curMinPrice == nextPrice {
			nextPrice++
		}
		g.curPage = 1
		g.curMinPrice = nextPrice
	} else {
		g.curPage++
	}

	return searchItemRes, nil
}
