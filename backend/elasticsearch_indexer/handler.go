package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/pm"
	"go.uber.org/zap"
)

func newItemUpdateHandler(itemIndexer *ItemIndexer, logger *zap.Logger) pm.MessageBatchHandler {
	return func(messages []*pubsub.Message) error {
		items := make([]*es.Item, 0, len(messages))
		for _, m := range messages {
			var item xitem.Item
			if err := json.Unmarshal(m.Data, &item); err != nil {
				logger.Error(
					"json.Unmarshal failed",
					zap.Error(err),
					zap.Any("messageId", m.ID),
					zap.String("data", string(m.Data)),
				)
				continue
			}
			// Remove top level category not to show irrelevant items.
			if len(item.CategoryNames) > 0 {
				item.CategoryNames = item.CategoryNames[1:]
			}
			items = append(items, mapItemFetcherItemToElasticsearchItem(&item))
		}

		if err := itemIndexer.BulkIndex(context.Background(), items); err != nil {
			return fmt.Errorf("itemIndexer.BulkIndex: %w", err)
		}
		logger.Info(fmt.Sprintf("bluk indexed %d items", len(items)))
		return nil
	}
}

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
		TagIDs:        item.TagIDs,
		JANCode:       item.JANCode,
		Platform:      item.Platform,
		IndexedAt:     time.Now().UnixMilli(),
	}
}
