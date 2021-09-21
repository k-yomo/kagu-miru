package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	gqlgend "github.com/k-yomo/kagu-miru/backend/graph/gqlgen"
	gqlmodell "github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/search"
)

func (r *queryResolver) SearchItems(ctx context.Context, input *gqlmodell.SearchItemsInput) (*gqlmodell.ItemConnection, error) {
	sortType, err := mapGraphqlSortTypeToSearchSortType(input.SortType)
	if err != nil {
		return nil, fmt.Errorf("mapGraphqlSortTypeToSearchSortType: %w", err)
	}
	page := search.DefaultPage
	if input.Page != nil {
		page = uint64(*input.Page)
	}
	searchResponse, err := r.SearchClient.SearchItems(ctx, &search.Request{
		Query:    input.Query,
		SortType: sortType,
		Page:     page,
		PageSize: input.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("SearchClient.SearchItems: %w", err)
	}

	return mapSearchResponseToGraphqlItemConnection(searchResponse)
}

func (r *queryResolver) GetQuerySuggestions(ctx context.Context, query string) ([]string, error) {
	querySuggestions, err := r.SearchClient.GetQuerySuggestions(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("SearchClient.GetQuerySuggestions: %w", err)
	}

	suggestedQueries := make([]string, 0, len(querySuggestions.Aggregations.Queries.Buckets))
	for _, bucket := range querySuggestions.Aggregations.Queries.Buckets {
		suggestedQueries = append(suggestedQueries, bucket.Key)
	}
	fmt.Println(suggestedQueries)

	return suggestedQueries, nil
}

// Query returns gqlgend.QueryResolver implementation.
func (r *Resolver) Query() gqlgend.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
