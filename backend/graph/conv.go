package graph

import (
	"fmt"
	"github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/internal/es"
)

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
		ID:             item.ID,
		Name:           item.Name,
		Description:    item.Description,
		Status:         status,
		SellingPageURL: item.SellingPageURL,
		Price:          item.Price,
		ImageUrls:      item.ImageURLs,
		AverageRating:  item.AverageRating,
		ReviewCount:    item.ReviewCount,
		Platform:       platform,
	}, nil
}
