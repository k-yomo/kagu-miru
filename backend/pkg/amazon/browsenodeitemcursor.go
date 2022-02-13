package amazon

import (
	"context"
	"errors"

	"github.com/utekaravinash/gopaapi5/api"

	"github.com/cenkalti/backoff/v4"
)

var Done = errors.New("DONE")

type BrowseNodeItemCursor struct {
	amazonClient *Client

	browseNodeID string
	curPage      int
	curMinPrice  int
	maxPrice     int

	isDone bool
}

func (c *Client) NewBrowseNodeItemCursor(browseNodeID string, minPrice int, maxPrice int) *BrowseNodeItemCursor {
	return &BrowseNodeItemCursor{
		amazonClient: c,
		browseNodeID: browseNodeID,
		curPage:      1,
		curMinPrice:  minPrice,
		maxPrice:     maxPrice,
	}
}

func (g *BrowseNodeItemCursor) CurPage() int {
	return g.curPage
}

func (g *BrowseNodeItemCursor) CurMinPrice() int {
	return g.curMinPrice
}

func (g *BrowseNodeItemCursor) Next(ctx context.Context) ([]api.Item, error) {
	if g.isDone {
		return nil, Done
	}
	searchItemParams := &api.SearchItemsParams{
		BrowseNodeId: g.browseNodeID,
		MinPrice:     g.curMinPrice,
		MaxPrice:     g.maxPrice,
		Condition:    api.New,
		// There is an api call limit for Amazon API, so fetch only high rating items
		MinReviewsRating:      3,
		LanguagesOfPreference: []api.Language{api.JapaneseJapan},
		ItemPage:              g.curPage,
		SortBy:                api.PriceLowToHigh,
		Resources: []api.Resource{
			api.BrowseNodeInfoBrowseNodes,
			api.ImagesPrimaryLarge,
			api.ImagesVariantsLarge,
			api.ItemInfoByLineInfo,
			api.ItemInfoProductInfo,
			api.ItemInfoTitle,
			api.OffersListingsAvailabilityType,
			api.OffersListingsCondition,
			api.OffersListingsConditionSubCondition,
			api.OffersListingsDeliveryInfoIsAmazonFulfilled,
			api.OffersListingsDeliveryInfoIsFreeShippingEligible,
			api.OffersListingsDeliveryInfoIsPrimeEligible,
			api.OffersListingsDeliveryInfoShippingCharges,
			api.OffersListingsIsBuyBoxWinner,
			api.OffersListingsLoyaltyPointsPoints,
			api.OffersListingsMerchantInfo,
			api.OffersListingsPrice,
			api.OffersListingsProgramEligibilityIsPrimeExclusive,
			api.OffersListingsProgramEligibilityIsPrimePantry,
			api.OffersListingsPromotions,
			api.OffersListingsSavingBasis,
			api.OffersSummariesHighestPrice,
			api.OffersSummariesLowestPrice,
			api.OffersSummariesOfferCount,
			api.ParentASIN,
		},
	}
	var searchItemRes *api.SearchResult
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
	err := backoff.Retry(func() error {
		var err error
		searchItemRes, err = g.amazonClient.SearchItems(ctx, searchItemParams)
		return err
	}, b)
	if err != nil {
		return nil, err
	}

	for i, item := range searchItemRes.Items {
		if int(item.Offers.Listings[0].Price.Amount) > g.maxPrice {
			searchItemRes.Items = searchItemRes.Items[:i]
			g.isDone = true
			return searchItemRes.Items, nil
		}
	}

	if noMoreItems := len(searchItemRes.Items) == searchItemRes.TotalResultCount; noMoreItems {
		g.isDone = true
	} else if g.curPage == maxPage {
		nextPrice := int(searchItemRes.Items[len(searchItemRes.Items)-1].Offers.Listings[0].Price.Amount)
		if g.curMinPrice == nextPrice {
			nextPrice++
		}
		g.curPage = 1
		g.curMinPrice = nextPrice
	} else {
		g.curPage++
	}

	return searchItemRes.Items, nil
}
