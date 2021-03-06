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

func (s *SpannerDBClient) GetSameGroupItemsByItemID(ctx context.Context, itemID string) ([]*xspanner.Item, error) {
	return xspanner.GetSameGroupItemsByItemID(ctx, s.spannerClient, itemID)
}

func (s *SpannerDBClient) GetAllActiveItemCategories(ctx context.Context) ([]*xspanner.ItemCategory, error) {
	return xspanner.GetAllActiveItemCategories(ctx, s.spannerClient)
}

func (s *SpannerDBClient) GetAllItemCategoriesWithParent(ctx context.Context) ([]*xspanner.ItemCategoryWithParent, error) {
	return xspanner.GetAllActiveItemCategoriesWithParent(ctx, s.spannerClient)
}

func (s *SpannerDBClient) GetTopLevelItemCategories(ctx context.Context) ([]*xspanner.ItemCategory, error) {
	return xspanner.GetTopLevelItemCategories(ctx, s.spannerClient)

}
