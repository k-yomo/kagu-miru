package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sort"

	gqlgend "github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlgen"
	gqlmodell "github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/tracking"
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

func (r *queryResolver) GetItem(ctx context.Context, id string) (*gqlmodell.Item, error) {
	item, err := r.SearchClient.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("SearchClient.GetItem: %w", err)
	}
	return mapSearchItemToGraphqlItem(item)
}

func (r *queryResolver) GetAllItemCategories(ctx context.Context) ([]*gqlmodell.ItemCategory, error) {
	allItemCategories, err := r.DBClient.GetAllItemCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("DBClient.GetAllItemCategories: %w", err)
	}

	gqlItemCategories := mapSpannerItemCategoriesToGraphqlItemCategories(allItemCategories)
	sort.Slice(gqlItemCategories, func(i, j int) bool {
		return gqlItemCategories[i].Level < gqlItemCategories[j].Level
	})

	itemCategoryMap := make(map[string]*gqlmodell.ItemCategory)
	for _, itemCategory := range gqlItemCategories {
		if itemCategory.Level == 0 {
			itemCategoryMap[itemCategory.ID] = itemCategory
		} else {
			itemCategoryMap[*itemCategory.ParentID].Children = append(itemCategoryMap[*itemCategory.ParentID].Children, itemCategory)
			itemCategoryMap[itemCategory.ID] = itemCategory
		}
	}

	var topLevelItemCategories []*gqlmodell.ItemCategory
	for _, itemCategory := range gqlItemCategories {
		if itemCategory.Level == 0 {
			topLevelItemCategories = append(topLevelItemCategories, itemCategory)
		}
	}
	return topLevelItemCategories, nil
}

// Mutation returns gqlgend.MutationResolver implementation.
func (r *Resolver) Mutation() gqlgend.MutationResolver { return &mutationResolver{r} }

// Query returns gqlgend.QueryResolver implementation.
func (r *Resolver) Query() gqlgend.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
