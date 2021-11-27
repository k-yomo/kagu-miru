package db

import (
	"context"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
)

type SpannerDBClient struct {
	spannerClient *spanner.Client
}

func NewSpannerDBClient(spannerClient *spanner.Client) *SpannerDBClient {
	return &SpannerDBClient{spannerClient: spannerClient}
}

func (s *SpannerDBClient) GetAllItemCategories(ctx context.Context) ([]*xspanner.ItemCategory, error) {
	stmt := spanner.NewStatement(`SELECT * FROM item_categories`)
	iter := s.spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var itemCategories []*xspanner.ItemCategory
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var itemCategory xspanner.ItemCategory
		if err := row.ToStruct(&itemCategory); err != nil {
			return nil, err
		}
		itemCategories = append(itemCategories, &itemCategory)
	}
	return itemCategories, nil
}
