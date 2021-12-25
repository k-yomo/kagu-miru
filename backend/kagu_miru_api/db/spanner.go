package db

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
)

type SpannerDBClient struct {
	spannerClient *spanner.Client
}

func NewSpannerDBClient(spannerClient *spanner.Client) *SpannerDBClient {
	return &SpannerDBClient{spannerClient: spannerClient}
}

func (s *SpannerDBClient) GetItem(ctx context.Context, itemID string) (*xspanner.Item, error) {
	return xspanner.GetItem(ctx, s.spannerClient, itemID)
}

func (s *SpannerDBClient) GetAllItemCategories(ctx context.Context) ([]*xspanner.ItemCategory, error) {
	return xspanner.GetAllItemCategories(ctx, s.spannerClient)
}

func (s *SpannerDBClient) GetAllItemCategoriesWithParent(ctx context.Context) ([]*xspanner.ItemCategoryWithParent, error) {
	return xspanner.GetAllItemCategoriesWithParent(ctx, s.spannerClient)
}
