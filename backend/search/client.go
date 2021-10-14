package search

import (
	"context"

	"github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
)

type Client interface {
	SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error)
	GetQuerySuggestions(ctx context.Context, query string) ([]string, error)
}

const (
	defaultPage     uint64 = 0
	defaultPageSize uint64 = 100
	maxPageSize     uint64 = 1000
)
