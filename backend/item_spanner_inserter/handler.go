package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
	"github.com/k-yomo/pm"
	"go.uber.org/zap"
)

func newItemUpdateHandler(spannerClient *spanner.Client, logger *zap.Logger) pm.MessageBatchHandler {
	return func(messages []*pubsub.Message) error {
		mutations := make([]*spanner.Mutation, 0, len(messages))
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

			mutation, err := spanner.InsertOrUpdateStruct(xspanner.ItemsTableName, mapItemToSpannerItem(&item))
			if err != nil {
				logger.Error(
					"spanner.InsertOrUpdateStruct failed",
					zap.Error(err),
					zap.Any("messageId", m.ID),
					zap.String("data", string(m.Data)),
				)
				continue
			}
			mutations = append(mutations, mutation)
		}

		_, err := spannerClient.Apply(context.Background(), mutations)
		if err != nil {
			return fmt.Errorf("spannerClient.Apply: %w", err)
		}
		logger.Info(fmt.Sprintf("bluk inserted / updated %d items", len(mutations)))
		return nil
	}
}

func mapItemToSpannerItem(item *xitem.Item) *xspanner.Item {
	return &xspanner.Item{
		ID:            item.ID,
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
		JANCode:       spanner.NullString{StringVal: item.JANCode, Valid: item.JANCode != ""},
		Platform:      item.Platform,
		UpdatedAt:     time.Now(),
	}
}
