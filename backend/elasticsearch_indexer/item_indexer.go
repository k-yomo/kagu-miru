package main

import (
	"context"
	"fmt"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/olivere/elastic/v7"
	"go.uber.org/multierr"
)

type ItemIndexer struct {
	indexName string
	esClient  *elastic.Client
}

func NewItemIndexer(indexName string, esClient *elastic.Client) *ItemIndexer {
	return &ItemIndexer{
		indexName: indexName,
		esClient:  esClient,
	}
}

func (i *ItemIndexer) BulkIndex(ctx context.Context, items []*es.Item) error {
	if len(items) == 0 {
		return nil
	}

	bulk := i.esClient.Bulk().Index(i.indexName)
	for _, item := range items {
		if item.IsActive() {
			bulk.Add(elastic.NewBulkIndexRequest().Index(i.indexName).Id(item.ID).Doc(item))
		} else {
			bulk.Add(elastic.NewBulkDeleteRequest().Index(i.indexName).Id(item.ID))
		}
	}

	resp, err := bulk.Do(ctx)
	if err != nil {
		return fmt.Errorf("esClient.Bulk failed: %w", err)
	}
	if resp.Errors {
		var errs []error
		for _, failed := range resp.Failed() {
			errs = append(errs, fmt.Errorf("id: %s, status: %d, err: %s", failed.Id, failed.Status, failed.Result))
		}
		return fmt.Errorf("read body failed: %w", multierr.Combine(errs...))
	}

	return nil
}
