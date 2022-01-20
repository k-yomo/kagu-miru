package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sort"

	"github.com/k-yomo/kagu-miru/backend/internal/xerror"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/cms"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlgen"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/tracking"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"golang.org/x/sync/errgroup"
)

func (r *mutationResolver) TrackEvent(ctx context.Context, event gqlmodel.Event) (bool, error) {
	r.EventLoader.Load(ctx, tracking.NewEvent(ctx, event))
	return true, nil
}

func (r *queryResolver) Home(ctx context.Context) (*gqlmodel.HomeResponse, error) {
	eg := errgroup.Group{}

	var rakutenItemsRes *gqlmodel.SearchResponse
	eg.Go(func() error {
		var err error
		rakutenItemsRes, err = r.Search(ctx, gqlmodel.SearchInput{
			Filter: &gqlmodel.SearchFilter{
				Platforms: []gqlmodel.ItemSellingPlatform{gqlmodel.ItemSellingPlatformRakuten},
			},
			SortType: gqlmodel.SearchSortTypeReviewCount,
			PageSize: func() *int { i := 10; return &i }(),
		})
		return err
	})

	var yahooShoppingItemsRes *gqlmodel.SearchResponse
	eg.Go(func() error {
		var err error
		yahooShoppingItemsRes, err = r.Search(ctx, gqlmodel.SearchInput{
			Filter: &gqlmodel.SearchFilter{
				Platforms: []gqlmodel.ItemSellingPlatform{gqlmodel.ItemSellingPlatformYahooShopping},
			},
			SortType: gqlmodel.SearchSortTypeReviewCount,
			PageSize: func() *int { i := 10; return &i }(),
		})
		return err
	})

	var paypayMallItemsRes *gqlmodel.SearchResponse
	eg.Go(func() error {
		var err error
		paypayMallItemsRes, err = r.Search(ctx, gqlmodel.SearchInput{
			Filter: &gqlmodel.SearchFilter{
				Platforms: []gqlmodel.ItemSellingPlatform{gqlmodel.ItemSellingPlatformPaypayMall},
			},
			SortType: gqlmodel.SearchSortTypeReviewCount,
			PageSize: func() *int { i := 10; return &i }(),
		})
		return err
	})

	var topLevelItemCategories []*gqlmodel.ItemCategory
	eg.Go(func() error {
		categories, err := r.DBClient.GetTopLevelItemCategories(ctx)
		if err != nil {
			return err
		}
		topLevelItemCategories = mapSpannerItemCategoriesToGraphqlItemCategories(categories)
		sort.Slice(topLevelItemCategories, func(i, j int) bool {
			return topLevelItemCategories[i].ID < topLevelItemCategories[j].ID
		})

		return nil
	})

	var featuredPostsResp *cms.GetFeaturedPostsResponse
	eg.Go(func() error {
		var err error
		featuredPostsResp, err = r.CMSClient.GetFeaturedPosts(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	platformPopularItemGroupsComponent := &gqlmodel.HomeComponent{
		ID: "platformPopularItemGroups",
		Payload: gqlmodel.HomeComponentPayloadItemGroups{
			Title: "EC人気アイテム",
			Payload: []*gqlmodel.HomeComponentPayloadItems{
				{
					Title: "楽天",
					Items: rakutenItemsRes.ItemConnection.Nodes,
				},
				{
					Title: "Yahooショッピング",
					Items: yahooShoppingItemsRes.ItemConnection.Nodes,
				},
				{
					Title: "PayPayモール",
					Items: paypayMallItemsRes.ItemConnection.Nodes,
				},
			},
		},
	}

	categoriesComponent := &gqlmodel.HomeComponent{
		ID: "categories",
		Payload: &gqlmodel.HomeComponentPayloadCategories{
			Title:      "カテゴリーから探す",
			Categories: topLevelItemCategories,
		},
	}

	featuredPostsComponent := &gqlmodel.HomeComponent{
		ID: "featuredPosts",
		Payload: &gqlmodel.HomeComponentPayloadMediaPosts{
			Title: featuredPostsResp.Title,
			Posts: mapPostsToGraphqlPosts(featuredPostsResp.Posts),
		},
	}

	return &gqlmodel.HomeResponse{Components: []*gqlmodel.HomeComponent{
		platformPopularItemGroupsComponent,
		categoriesComponent,
		featuredPostsComponent,
	}}, nil
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
		return nil, fmt.Errorf("SearchClient.SearchItems: %w", err)
	}

	gqlRes, err := mapSearchResponseToGraphqlSearchResponse(resp, r.SearchIDManager.GetSearchID(ctx))
	if err != nil {
		return nil, logging.Error(ctx, fmt.Errorf("mapSearchResponseToGraphqlGetSimilarItemsResponse: %w", err))
	}
	return gqlRes, nil
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

	gqlRes, err := mapSearchResponseToGraphqlGetSimilarItemsResponse(resp, r.SearchIDManager.GetSearchID(ctx))
	if err != nil {
		return nil, logging.Error(ctx, fmt.Errorf("mapSearchResponseToGraphqlGetSimilarItemsResponse: %w", err))
	}
	return gqlRes, nil
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
			return nil, logging.Error(ctx, fmt.Errorf("mapSpannerItemToGraphqlItem :%w", err))
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
