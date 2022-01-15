package main

import (
	"time"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
)

func mapItemFetcherItemToElasticsearchItem(item *xitem.Item) *es.Item {
	return &es.Item{
		ID:            item.ID,
		Name:          item.Name,
		Description:   item.Description,
		Status:        item.Status,
		URL:           item.URL,
		AffiliateURL:  item.AffiliateURL,
		Price:         item.Price,
		ImageURLs:     item.ImageURLs,
		AverageRating: item.AverageRating,
		ReviewCount:   item.ReviewCount,
		CategoryID:    item.CategoryID,
		CategoryIDs:   item.CategoryIDs,
		CategoryNames: item.CategoryNames,
		BrandName:     item.BrandName,
		Colors:        item.Colors,
		Metadata:      extractMetadata(item),
		JANCode:       item.JANCode,
		Platform:      item.Platform,
		IndexedAt:     time.Now().UnixMilli(),
	}
}

func extractMetadata(item *xitem.Item) []es.Metadata {
	var facets []es.Metadata
	if item.WidthRange != nil {
		facets = append(facets, es.Metadata{
			Name:  es.MetadataNameWidthRange,
			Value: es.NewMetadataValueLengthRange(item.WidthRange.Gte, item.WidthRange.Lte),
		})
	}
	if item.DepthRange != nil {
		facets = append(facets, es.Metadata{
			Name:  es.MetadataNameDepthRange,
			Value: es.NewMetadataValueLengthRange(item.DepthRange.Gte, item.DepthRange.Lte),
		})
	}
	if item.HeightRange != nil {
		facets = append(facets, es.Metadata{
			Name:  es.MetadataNameHeightRange,
			Value: es.NewMetadataValueLengthRange(item.HeightRange.Gte, item.HeightRange.Lte),
		})
	}

	return facets
}

func mapItemToSpannerItem(item *xitem.Item, groupID string) *xspanner.Item {
	return &xspanner.Item{
		ID:            item.ID,
		GroupID:       spanner.NullString{StringVal: groupID, Valid: true},
		Name:          item.Name,
		Description:   item.Description,
		Status:        int64(item.Status),
		URL:           item.URL,
		AffiliateURL:  item.AffiliateURL,
		Price:         int64(item.Price),
		ImageURLs:     item.ImageURLs,
		AverageRating: item.AverageRating,
		ReviewCount:   int64(item.ReviewCount),
		CategoryID:    item.CategoryID,
		BrandName:     spanner.NullString{StringVal: item.BrandName, Valid: item.BrandName != ""},
		Colors:        item.Colors,
		WidthRange:    mapIntRangeToSpannerRange(item.WidthRange),
		DepthRange:    mapIntRangeToSpannerRange(item.DepthRange),
		HeightRange:   mapIntRangeToSpannerRange(item.HeightRange),
		JANCode:       spanner.NullString{StringVal: item.JANCode, Valid: item.JANCode != ""},
		Platform:      item.Platform,
		UpdatedAt:     time.Now(),
	}
}

func mapIntRangeToSpannerRange(r *xitem.IntRange) []int64 {
	if r == nil {
		return nil
	}
	return []int64{int64(r.Gte), int64(r.Lte)}
}
