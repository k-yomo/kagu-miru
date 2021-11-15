package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/pm"
	"go.uber.org/zap"
)

func newItemUpdateHandler(itemIndexer *ItemIndexer, logger *zap.Logger) pm.MessageBatchHandler {
	return func(messages []*pubsub.Message) error {
		items := make([]*es.Item, 0, len(messages))
		for _, m := range messages {
			var item es.Item
			if err := json.Unmarshal(m.Data, &item); err != nil {
				logger.Error(
					"json.Unmarshal failed",
					zap.Error(err),
					zap.Any("messageId", m.ID),
					zap.String("data", string(m.Data)),
				)
				continue
			}
			items = append(items, &item)
		}

		if err := itemIndexer.BulkIndex(context.Background(), items); err != nil {
			return fmt.Errorf("itemIndexer.BulkIndex: %w", err)
		}
		logger.Info(fmt.Sprintf("bluk indexed %d items", len(items)))
		return nil
	}
}
