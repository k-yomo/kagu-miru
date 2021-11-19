package search

import (
	"context"

	"github.com/k-yomo/kagu-miru/backend/internal/es"

	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
)

type Client interface {
	GetItem(ctx context.Context, id string) (*es.Item, error)
	SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error)
	GetQuerySuggestions(ctx context.Context, query string) ([]string, error)
}

const (
	defaultPage     uint64 = 0
	defaultPageSize uint64 = 100
	maxPageSize     uint64 = 1000
)
