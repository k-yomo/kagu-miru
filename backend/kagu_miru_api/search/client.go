package search

import (
	"context"

	"github.com/k-yomo/kagu-miru/backend/internal/es"

	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
)

type Client interface {
	GetItem(ctx context.Context, id string) (*es.Item, error)
	SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error)
	GetSimilarItems(ctx context.Context, input *gqlmodel.GetSimilarItemsInput, itemCategoryID string) (*Response, error)
	GetQuerySuggestions(ctx context.Context, query string) ([]string, error)
}

const (
	defaultPage     int = 0
	defaultPageSize int = 100
	maxPageSize     int = 1000
)
