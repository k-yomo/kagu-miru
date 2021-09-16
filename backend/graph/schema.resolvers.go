package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	gqlgend "github.com/k-yomo/kagu-miru/backend/graph/gqlgen"
	gqlmodell "github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/search"
	"github.com/k-yomo/kagu-miru/internal/es"
)

func (r *queryResolver) SearchItems(ctx context.Context, input *gqlmodell.SearchItemsInput) ([]*gqlmodell.Item, error) {
	var page uint64
	if input.Page != nil {
		page = 0
	}
	searchResponse, err := r.SearchClient.SearchItems(ctx, &search.Request{
		Query:    input.Query,
		Page:     page,
		PageSize: input.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("SearchClient.SearchItems: %w", err)
	}

	items := make([]*es.Item, 0, len(searchResponse.Result.Hits.Hits))
	for _, hit := range searchResponse.Result.Hits.Hits {
		items = append(items, hit.Source)
	}

	return mapSearchItemsToGraphqlItems(items)
}

// Query returns gqlgend.QueryResolver implementation.
func (r *Resolver) Query() gqlgend.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
