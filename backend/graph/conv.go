package graph

import (
	"fmt"

	"github.com/k-yomo/kagu-miru/backend/search"

	"github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/internal/es"
)

func mapGraphqlSortTypeToSearchSortType(st gqlmodel.SearchItemsSortType) (search.SortType, error) {
	switch st {
	case gqlmodel.SearchItemsSortTypeBestMatch:
		return search.SortTypeBestMatch, nil
	case gqlmodel.SearchItemsSortTypePriceAsc:
		return search.SortTypePriceAsc, nil
	case gqlmodel.SearchItemsSortTypePriceDesc:
		return search.SortTypePriceDesc, nil
	case gqlmodel.SearchItemsSortTypeReviewCount:
		return search.SortTypeReviewCount, nil
	case gqlmodel.SearchItemsSortTypeRating:
		return search.SortTypeRating, nil
	default:
		return 0, fmt.Errorf("unknown sort type '%s' is given", st)
	}
}

func mapSearchItemToGraphqlItem(item *es.Item) (*gqlmodel.Item, error) {
	var status gqlmodel.ItemStatus
	switch item.Status {
	case es.StatusActive:
		status = gqlmodel.ItemStatusActive
	case es.StatusInactive:
		status = gqlmodel.ItemStatusInactive
	default:
		return nil, fmt.Errorf("unknown status %d, item: %v", item.Status, item)
	}

	var platform gqlmodel.ItemSellingPlatform
	switch item.Platform {
	case es.PlatformRakuten:
		platform = gqlmodel.ItemSellingPlatformRakuten
	default:
		return nil, fmt.Errorf("unknown platform %d, item: %v", item.Status, item)
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
		Platform:      platform,
	}, nil
}

func mapSearchItemsToGraphqlItems(items []*es.Item) ([]*gqlmodel.Item, error) {
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

func mapSearchResponseToGraphqlItemConnection(res *search.Response) (*gqlmodel.ItemConnection, error) {
	items := make([]*es.Item, 0, len(res.Result.Hits.Hits))
	for _, hit := range res.Result.Hits.Hits {
		items = append(items, hit.Source)
	}

	graphqlItems, err := mapSearchItemsToGraphqlItems(items)
	if err != nil {
		return nil, err
	}
	return &gqlmodel.ItemConnection{
		PageInfo: &gqlmodel.PageInfo{
			Page:      int(res.Page),
			TotalPage: int(res.TotalPage),
		},
		Nodes: graphqlItems,
	}, nil
}
