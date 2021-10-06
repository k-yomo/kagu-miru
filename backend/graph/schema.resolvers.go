package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	gqlgend "github.com/k-yomo/kagu-miru/backend/graph/gqlgen"
	gqlmodell "github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/tracking"
)

func (r *mutationResolver) TrackEvent(ctx context.Context, event gqlmodell.Event) (bool, error) {
	r.EventLoader.Load(ctx, tracking.NewEvent(ctx, event))
	return true, nil
}

func (r *queryResolver) Search(ctx context.Context, input *gqlmodell.SearchInput) (*gqlmodell.SearchResponse, error) {
	searchResponse, err := r.SearchClient.SearchItems(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("SearchClient.SearchItems: %w", err)
	}

	return mapSearchResponseToGraphqlSearchResponse(searchResponse, r.SearchIDManager.GetSearchID(ctx))
}

func (r *queryResolver) GetQuerySuggestions(ctx context.Context, query string) (*gqlmodell.QuerySuggestionsResponse, error) {
	suggestedQueries, err := r.SearchClient.GetQuerySuggestions(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("SearchClient.GetQuerySuggestions: %w", err)
	}

	return &gqlmodell.QuerySuggestionsResponse{
		Query:            query,
		SuggestedQueries: suggestedQueries,
	}, nil
}

func (r *queryResolver) GetAllCategories(ctx context.Context) ([]*gqlmodell.ItemCategory, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns gqlgend.MutationResolver implementation.
func (r *Resolver) Mutation() gqlgend.MutationResolver { return &mutationResolver{r} }

// Query returns gqlgend.QueryResolver implementation.
func (r *Resolver) Query() gqlgend.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
