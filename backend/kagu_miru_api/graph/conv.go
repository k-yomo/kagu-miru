package graph

import (
	"fmt"

	"github.com/k-yomo/kagu-miru/backend/internal/xitem"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/search"
)

func mapSearchItemToGraphqlItem(item *es.Item) (*gqlmodel.Item, error) {
	var status gqlmodel.ItemStatus
	switch item.Status {
	case xitem.StatusActive:
		status = gqlmodel.ItemStatusActive
	case xitem.StatusInactive:
		status = gqlmodel.ItemStatusInactive
	default:
		return nil, fmt.Errorf("unknown status %d, item: %v", item.Status, item)
	}

	var platform gqlmodel.ItemSellingPlatform
	switch item.Platform {
	case xitem.PlatformRakuten:
		platform = gqlmodel.ItemSellingPlatformRakuten
	case xitem.PlatformYahooShopping:
		platform = gqlmodel.ItemSellingPlatformYahooShopping
	default:
		return nil, fmt.Errorf("unknown platform %s, item: %v", item.Platform, item)
	}

	return &gqlmodel.Item{
		ID:            item.ID,
		Name:          item.Name,
		Description:   item.Description,
		Status:        status,
		URL:           item.URL,
		AffiliateURL:  item.AffiliateURL,
		Price:         item.Price,
		ImageUrls:     item.ImageURLs,
		AverageRating: item.AverageRating,
		ReviewCount:   item.ReviewCount,
		CategoryIds:   item.CategoryIDs,
		Platform:      platform,
	}, nil
}

func mapSearchToGraphqlItems(items []*es.Item) ([]*gqlmodel.Item, error) {
	gqlItems := make([]*gqlmodel.Item, 0, len(items))
	for _, item := range items {
		gqlItem, err := mapSearchItemToGraphqlItem(item)
		if err != nil {
			return nil, err
		}
		gqlItems = append(gqlItems, gqlItem)
	}
	return gqlItems, nil
}

func mapSearchResponseToGraphqlSearchResponse(res *search.Response, searchID string) (*gqlmodel.SearchResponse, error) {
	graphqlItems, err := mapSearchToGraphqlItems(res.Items)
	if err != nil {
		return nil, err
	}
	return &gqlmodel.SearchResponse{
		SearchID: searchID,
		ItemConnection: &gqlmodel.ItemConnection{
			PageInfo: &gqlmodel.PageInfo{
				Page:       int(res.Page),
				TotalPage:  int(res.TotalPage),
				TotalCount: int(res.TotalCount),
			},
			Nodes: graphqlItems,
		},
	}, nil
}
