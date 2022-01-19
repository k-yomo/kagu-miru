package graph

import (
	"fmt"

	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/cms"

	"github.com/k-yomo/kagu-miru/backend/internal/xerror"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/search"
	"github.com/k-yomo/kagu-miru/backend/pkg/pointerconv"
)

func mapSpannerItemCategoriesToGraphqlItemCategories(itemCategories []*xspanner.ItemCategory) []*gqlmodel.ItemCategory {
	gqlItemCategories := make([]*gqlmodel.ItemCategory, 0, len(itemCategories))
	for _, itemCategory := range itemCategories {
		gqlItemCategory := &gqlmodel.ItemCategory{
			ID:       itemCategory.ID,
			Name:     itemCategory.Name,
			Level:    int(itemCategory.Level),
			ParentID: pointerconv.StringToPointer(itemCategory.ParentID.StringVal),
			ImageURL: pointerconv.StringToPointer(itemCategory.ImageURL.StringVal),
		}
		gqlItemCategories = append(gqlItemCategories, gqlItemCategory)
	}
	return gqlItemCategories
}

func mapSearchItemToGraphqlItem(item *es.Item) (*gqlmodel.Item, error) {
	var status gqlmodel.ItemStatus
	switch item.Status {
	case xitem.StatusActive:
		status = gqlmodel.ItemStatusActive
	case xitem.StatusInactive:
		status = gqlmodel.ItemStatusInactive
	default:
		return nil, fmt.Errorf("unknown status %d, item: %v", item.Status, item)
	}

	var platform gqlmodel.ItemSellingPlatform
	switch item.Platform {
	case xitem.PlatformRakuten:
		platform = gqlmodel.ItemSellingPlatformRakuten
	case xitem.PlatformYahooShopping:
		platform = gqlmodel.ItemSellingPlatformYahooShopping
	case xitem.PlatformPayPayMall:
		platform = gqlmodel.ItemSellingPlatformPaypayMall
	default:
		return nil, fmt.Errorf("unknown platform %s, item: %v", item.Platform, item)
	}

	colors := make([]gqlmodel.ItemColor, 0, len(item.Colors))
	for _, color := range item.Colors {
		gqlColor := mapSearchItemColorToGraphqlItemColor(color)
		if gqlColor.IsValid() {
			colors = append(colors, gqlColor)
		}
	}

	return &gqlmodel.Item{
		ID:            item.ID,
		Name:          item.Name,
		Description:   item.Description,
		Status:        status,
		URL:           item.URL,
		AffiliateURL:  item.AffiliateURL,
		Price:         item.Price,
		ImageUrls:     item.ImageURLs,
		AverageRating: item.AverageRating,
		ReviewCount:   item.ReviewCount,
		CategoryID:    item.CategoryID,
		Colors:        colors,
		Platform:      platform,
	}, nil
}

func mapSearchToGraphqlItems(items []*es.Item) ([]*gqlmodel.Item, error) {
	gqlItems := make([]*gqlmodel.Item, 0, len(items))
	for _, item := range items {
		gqlItem, err := mapSearchItemToGraphqlItem(item)
		if err != nil {
			return nil, err
		}
		gqlItems = append(gqlItems, gqlItem)
	}
	return gqlItems, nil
}

func mapSearchItemColorToGraphqlItemColor(color string) gqlmodel.ItemColor {
	switch color {
	case "ホワイト":
		return gqlmodel.ItemColorWhite
	case "イエロー":
		return gqlmodel.ItemColorYellow
	case "オレンジ":
		return gqlmodel.ItemColorOrange
	case "ピンク":
		return gqlmodel.ItemColorPink
	case "レッド":
		return gqlmodel.ItemColorRed
	case "ベージュ":
		return gqlmodel.ItemColorBeige
	case "シルバー":
		return gqlmodel.ItemColorSilver
	case "ゴールド":
		return gqlmodel.ItemColorGold
	case "グレー":
		return gqlmodel.ItemColorGray
	case "パープル":
		return gqlmodel.ItemColorPurple
	case "ブラウン":
		return gqlmodel.ItemColorBrown
	case "グリーン":
		return gqlmodel.ItemColorGreen
	case "ブルー":
		return gqlmodel.ItemColorBlue
	case "ブラック":
		return gqlmodel.ItemColorBlack
	case "ネイビー":
		return gqlmodel.ItemColorNavy
	case "カーキ":
		return gqlmodel.ItemColorKhaki
	case "ワインレッド":
		return gqlmodel.ItemColorWineRed
	case "透明":
		return gqlmodel.ItemColorTransparent
	default:
		return ""
	}
}

func mapSearchResponseToGraphqlSearchResponse(resp *search.Response, searchID string) (*gqlmodel.SearchResponse, error) {
	graphqlItems, err := mapSearchToGraphqlItems(resp.Items)
	if err != nil {
		return nil, err
	}
	return &gqlmodel.SearchResponse{
		SearchID: searchID,
		ItemConnection: &gqlmodel.ItemConnection{
			PageInfo: &gqlmodel.PageInfo{
				Page:       resp.Page,
				TotalPage:  resp.TotalPage,
				TotalCount: resp.TotalCount,
			},
			Nodes: graphqlItems,
		},
		Facets: mapSearchFacetsToGraphqlFacets(resp.Facets),
	}, nil
}

func mapSearchFacetsToGraphqlFacets(facets []*search.Facet) []*gqlmodel.Facet {
	graphqlFacets := make([]*gqlmodel.Facet, 0, len(facets))
	for _, facet := range facets {
		var facetValues []*gqlmodel.FacetValue
		switch facet.FacetType {
		case search.FacetTypeColors:
			facetValues = mapSearchFacetColorValuesToGraphqlFacetValues(facet.Values)
		default:
			facetValues = mapSearchFacetValuesToGraphqlFacetValues(facet.Values)
		}
		graphqlFacets = append(graphqlFacets, &gqlmodel.Facet{
			Title:      facet.Title,
			FacetType:  mapSearchFacetTypeToGraphqlFacetType(facet.FacetType),
			Values:     facetValues,
			TotalCount: facet.TotalCount,
		})
	}

	return graphqlFacets
}

func mapSearchFacetTypeToGraphqlFacetType(facetType search.FacetType) gqlmodel.FacetType {
	switch facetType {
	case search.FacetTypeCategoryIDs:
		return gqlmodel.FacetTypeCategoryIDS
	case search.FacetTypeBrandNames:
		return gqlmodel.FacetTypeBrandNames
	case search.FacetTypeColors:
		return gqlmodel.FacetTypeColors
	case search.FacetTypeMetadata:
		return gqlmodel.FacetTypeMetadata
	default:
		return ""
	}
}

func mapSearchFacetValuesToGraphqlFacetValues(facetValues []*search.FacetValue) []*gqlmodel.FacetValue {
	graphqlFacetValues := make([]*gqlmodel.FacetValue, 0, len(facetValues))
	for _, facetValue := range facetValues {
		graphqlFacetValues = append(graphqlFacetValues, &gqlmodel.FacetValue{
			ID:    facetValue.ID,
			Name:  facetValue.Name,
			Count: facetValue.Count,
		})
	}
	return graphqlFacetValues
}

func mapSearchFacetColorValuesToGraphqlFacetValues(facetValues []*search.FacetValue) []*gqlmodel.FacetValue {
	graphqlFacetValues := make([]*gqlmodel.FacetValue, 0, len(facetValues))
	for _, facetValue := range facetValues {
		graphqlFacetValues = append(graphqlFacetValues, &gqlmodel.FacetValue{
			ID:    mapSearchItemColorToGraphqlItemColor(facetValue.ID).String(),
			Name:  facetValue.Name,
			Count: facetValue.Count,
		})
	}
	return graphqlFacetValues
}

func mapSearchResponseToGraphqlGetSimilarItemsResponse(res *search.Response, searchID string) (*gqlmodel.GetSimilarItemsResponse, error) {
	graphqlItems, err := mapSearchToGraphqlItems(res.Items)
	if err != nil {
		return nil, err
	}
	return &gqlmodel.GetSimilarItemsResponse{
		SearchID: searchID,
		ItemConnection: &gqlmodel.ItemConnection{
			PageInfo: &gqlmodel.PageInfo{
				Page:       int(res.Page),
				TotalPage:  int(res.TotalPage),
				TotalCount: int(res.TotalCount),
			},
			Nodes: graphqlItems,
		},
	}, nil
}

func mapSpannerItemToGraphqlItem(item *xspanner.Item) (*gqlmodel.Item, error) {
	var status gqlmodel.ItemStatus
	switch xitem.Status(item.Status) {
	case xitem.StatusActive:
		status = gqlmodel.ItemStatusActive
	case xitem.StatusInactive:
		status = gqlmodel.ItemStatusInactive
	default:
		return nil, fmt.Errorf("unknown status %d, item: %v", item.Status, item)
	}

	var platform gqlmodel.ItemSellingPlatform
	switch item.Platform {
	case xitem.PlatformRakuten:
		platform = gqlmodel.ItemSellingPlatformRakuten
	case xitem.PlatformYahooShopping:
		platform = gqlmodel.ItemSellingPlatformYahooShopping
	case xitem.PlatformPayPayMall:
		platform = gqlmodel.ItemSellingPlatformPaypayMall
	default:
		return nil, fmt.Errorf("unknown platform %s, item: %v", item.Platform, item)
	}

	return &gqlmodel.Item{
		ID:            item.ID,
		Name:          item.Name,
		Description:   item.Description,
		Status:        status,
		URL:           item.URL,
		AffiliateURL:  item.AffiliateURL,
		Price:         int(item.Price),
		ImageUrls:     item.ImageURLs,
		AverageRating: item.AverageRating,
		ReviewCount:   int(item.ReviewCount),
		CategoryID:    item.CategoryID,
		Platform:      platform,
	}, nil
}

func mapFromXErrorType(errType xerror.Type) gqlmodel.ErrorCode {
	switch errType {
	case xerror.TypeNotFound:
		return gqlmodel.ErrorCodeNotFound
	default:
		return gqlmodel.ErrorCodeInternal
	}
}

func mapPostsToGraphqlPosts(posts []*cms.Post) []*gqlmodel.MediaPost {
	gqlPosts := make([]*gqlmodel.MediaPost, 0, len(posts))
	for _, post := range posts {
		gqlPosts = append(gqlPosts, &gqlmodel.MediaPost{
			ID:           post.ID,
			Slug:         post.Slug,
			Title:        post.Title,
			Description:  post.Description,
			MainImageURL: post.MainImage.Asset.Ref,
			PublishedAt:  post.PublishedAt,
			Categories:   mapPostCategoriesToGraphqlPostCategories(post.Categories),
		})
	}

	return gqlPosts
}

func mapPostCategoriesToGraphqlPostCategories(postCategories []*cms.Category) []*gqlmodel.MediaPostCategory {
	gqlPostCategories := make([]*gqlmodel.MediaPostCategory, 0, len(postCategories))
	for _, category := range postCategories {
		gqlPostCategories = append(gqlPostCategories, &gqlmodel.MediaPostCategory{
			ID:    category.ID,
			Names: category.Names,
		})
	}
	return gqlPostCategories
}
