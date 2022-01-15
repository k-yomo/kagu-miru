package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/pm"
	"go.uber.org/zap"
)

func newItemUpdateHandler(itemIndexer *ItemIndexer, logger *zap.Logger) pm.MessageBatchHandler {
	return func(messages []*pubsub.Message) error {
		items := make([]*xitem.Item, 0, len(messages))
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
			items = append(items, &item)
		}

		if err := itemIndexer.BulkIndex(context.Background(), items); err != nil {
			logger.Error(
				"itemIndexer.BulkIndex failed",
				zap.Error(err),
				zap.Any("items", items),
			)
			return fmt.Errorf("itemIndexer.BulkIndex: %w", err)
		}
		logger.Info(fmt.Sprintf("bluk indexed %d items", len(items)))
		return nil
	}
}
