package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sort"

	"github.com/k-yomo/kagu-miru/backend/internal/xerror"

	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlgen"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/tracking"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"go.uber.org/zap"
)

func (r *mutationResolver) TrackEvent(ctx context.Context, event gqlmodel.Event) (bool, error) {
	r.EventLoader.Load(ctx, tracking.NewEvent(ctx, event))
	return true, nil
}

func (r *queryResolver) Search(ctx context.Context, input gqlmodel.SearchInput) (*gqlmodel.SearchResponse, error) {
	// Categorizing Query is slow (~ 1sec), so disabling for the time being
	// if input.Query != "" && len(input.Filter.CategoryIds) == 0 {
	// 	categoryIDs, err := r.QueryClassifierClient.CategorizeQuery(ctx, input.Query)
	// 	if err != nil {
	// 		logging.Logger(ctx).Error("failed to predict query's category", zap.Error(err))
	// 	}
	// 	input.Filter.CategoryIds = categoryIDs
	// }

	resp, err := r.SearchClient.SearchItems(ctx, &input)
	if err != nil {
		logging.Logger(ctx).Error(fmt.Sprintf("SearchClient.SearchItems: %v", err), zap.Error(err))
		return nil, fmt.Errorf("SearchClient.SearchItems: %w", err)
	}
	return mapSearchResponseToGraphqlSearchResponse(resp, r.SearchIDManager.GetSearchID(ctx))
}

func (r *queryResolver) GetSimilarItems(ctx context.Context, input gqlmodel.GetSimilarItemsInput) (*gqlmodel.GetSimilarItemsResponse, error) {
	item, err := r.DBClient.GetItem(ctx, input.ItemID)
	if err != nil {
		return nil, fmt.Errorf("DBClient.GetItem: %w", err)
	}
	resp, err := r.SearchClient.GetSimilarItems(ctx, &input, item)
	if err != nil {
		return nil, fmt.Errorf("SearchClient.GetSimilarItems: %w", err)
	}

	return mapSearchResponseToGraphqlGetSimilarItemsResponse(resp, r.SearchIDManager.GetSearchID(ctx))
}

func (r *queryResolver) GetQuerySuggestions(ctx context.Context, query string) (*gqlmodel.QuerySuggestionsResponse, error) {
	suggestedQueries, err := r.SearchClient.GetQuerySuggestions(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("SearchClient.GetQuerySuggestions: %w", err)
	}

	return &gqlmodel.QuerySuggestionsResponse{
		Query:            query,
		SuggestedQueries: suggestedQueries,
	}, nil
}

func (r *queryResolver) GetItem(ctx context.Context, id string) (*gqlmodel.Item, error) {
	// This is temp implementation while migration since eventually item must have group id
	item, err := r.DBClient.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("DBClient.GetItem: %w", err)
	}
	if !item.GroupID.Valid {
		return mapSpannerItemToGraphqlItem(item)
	}

	items, err := r.DBClient.GetSameGroupItemsByItemID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("DBClient.GetSameGroupItemsByItemID: %w", err)
	}
	if len(items) == 0 {
		return nil, xerror.NewNotFound(fmt.Errorf("item '%s' is not found", id))
	}

	var targetItem *gqlmodel.Item
	var sameGroupItems []*gqlmodel.Item
	for _, item := range items {
		gqlItem, err := mapSpannerItemToGraphqlItem(item)
		if err != nil {
			return nil, err
		}
		if item.ID == id {
			targetItem = gqlItem
		} else {
			sameGroupItems = append(sameGroupItems, gqlItem)
		}
	}
	sort.Slice(sameGroupItems, func(i, j int) bool {
		if sameGroupItems[i].AverageRating == sameGroupItems[j].AverageRating {
			return sameGroupItems[i].ReviewCount > sameGroupItems[j].ReviewCount
		}
		return sameGroupItems[i].AverageRating > sameGroupItems[j].AverageRating
	})
	targetItem.SameGroupItems = sameGroupItems

	return targetItem, nil
}

func (r *queryResolver) GetAllItemCategories(ctx context.Context) ([]*gqlmodel.ItemCategory, error) {
	allItemCategories, err := r.DBClient.GetAllItemCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("DBClient.GetAllItemCategories: %w", err)
	}

	gqlItemCategories := mapSpannerItemCategoriesToGraphqlItemCategories(allItemCategories)
	sort.Slice(gqlItemCategories, func(i, j int) bool {
		return gqlItemCategories[i].Level < gqlItemCategories[j].Level
	})

	itemCategoryMap := make(map[string]*gqlmodel.ItemCategory)
	for _, itemCategory := range gqlItemCategories {
		if itemCategory.Level == 0 {
			itemCategoryMap[itemCategory.ID] = itemCategory
		} else {
			itemCategoryMap[*itemCategory.ParentID].Children = append(itemCategoryMap[*itemCategory.ParentID].Children, itemCategory)
			itemCategoryMap[itemCategory.ID] = itemCategory
		}
	}

	var topLevelItemCategories []*gqlmodel.ItemCategory
	for _, itemCategory := range gqlItemCategories {
		if itemCategory.Level == 0 {
			topLevelItemCategories = append(topLevelItemCategories, itemCategory)
		}
	}
	return topLevelItemCategories, nil
}

// Mutation returns gqlgen.MutationResolver implementation.
func (r *Resolver) Mutation() gqlgen.MutationResolver { return &mutationResolver{r} }

// Query returns gqlgen.QueryResolver implementation.
func (r *Resolver) Query() gqlgen.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
