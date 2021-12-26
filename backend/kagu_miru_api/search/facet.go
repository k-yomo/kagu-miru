package search

import (
	"context"
	"strings"

	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/pkg/interfaceconv"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

type FacetType int

const (
	FacetTypeCategoryIDs FacetType = iota + 1
	FacetTypeBrandNames
	FacetTypeColors
)

var facetTypeTitleMap = map[FacetType]string{
	FacetTypeCategoryIDs: "カテゴリー",
	FacetTypeBrandNames:  "ブランド",
	FacetTypeColors:      "カラー",
}

type Facet struct {
	Title      string
	FacetType  FacetType
	Values     []*FacetValue
	TotalCount int
}

type FacetValue struct {
	ID    string
	Name  string
	Count int
}

func newFacetFromBucketKeyItems(
	bucketKeyItems *elastic.AggregationBucketKeyItems,
	facetType FacetType,
	idNameMap map[string]string,
) *Facet {
	facetValues := make([]*FacetValue, 0, len(bucketKeyItems.Buckets))
	totalCount := bucketKeyItems.SumOfOtherDocCount
	for _, bucket := range bucketKeyItems.Buckets {
		keyStr, ok := bucket.Key.(string)
		if !ok {
			continue
		}
		facetValues = append(facetValues, &FacetValue{
			ID:    keyStr,
			Name:  idNameMap[keyStr],
			Count: int(bucket.DocCount),
		})
		totalCount += bucket.DocCount
	}
	return &Facet{
		Title:      facetTypeTitleMap[facetType],
		FacetType:  facetType,
		Values:     facetValues,
		TotalCount: int(totalCount),
	}
}

func addAggregationsAndPostFiltersForFacets(search *elastic.SearchService, searchFilter *gqlmodel.SearchFilter) (*elastic.SearchService, map[string]elastic.Query) {
	// filters for facetable fields
	postFilterMap := make(map[string]elastic.Query)
	if len(searchFilter.CategoryIds) > 0 {
		var categoryIDs []interface{}
		for _, id := range searchFilter.CategoryIds {
			categoryIDs = append(categoryIDs, id)
		}
		// Since we aggregate for category_id
		postFilterMap[es.ItemFieldCategoryID] = elastic.NewTermsQuery(es.ItemFieldCategoryIDs, categoryIDs...)
	}
	if len(searchFilter.BrandNames) > 0 {
		var brandNames []interface{}
		for _, brandName := range searchFilter.BrandNames {
			brandNames = append(brandNames, brandName)
		}
		postFilterMap[es.ItemFieldBrandName] = elastic.NewTermsQuery(es.ItemFieldBrandName, brandNames...)
	}
	if len(searchFilter.Colors) > 0 {
		var colors []interface{}
		for _, gqlColor := range searchFilter.Colors {
			if color := mapGraphqlItemColorToSearchItemColor(gqlColor); color != "" {
				colors = append(colors, color)
			}
		}
		postFilterMap[es.ItemFieldColors] = elastic.NewTermsQuery(es.ItemFieldColors, colors...)
	}

	var metadataFilters []elastic.Query
	metadataFilterMap := make(map[string]elastic.Query)
	if len(searchFilter.Metadata) > 0 {
		for _, metadata := range searchFilter.Metadata {
			boolQuery := elastic.NewBoolQuery().Filter(
				elastic.NewTermQuery(es.MetadataNameFullPath, metadata.Name),
				elastic.NewTermsQuery(es.MetadataValueFullPath, interfaceconv.StringArrayToInterfaceArray(metadata.Values)),
			)
			filter := elastic.NewNestedQuery(es.ItemFieldMetadata, boolQuery)
			metadataFilterMap[metadata.Name] = filter
			metadataFilters = append(metadataFilters, filter)
		}
	}

	search.
		Aggregation(
			es.ItemFieldCategoryID,
			newFilterAggregationForFacet(
				es.ItemFieldCategoryID,
				append(extractFiltersExceptForField(es.ItemFieldCategoryID, postFilterMap), metadataFilters...),
			),
		).
		Aggregation(
			es.ItemFieldBrandName,
			newFilterAggregationForFacet(
				es.ItemFieldBrandName,
				append(extractFiltersExceptForField(es.ItemFieldBrandName, postFilterMap), metadataFilters...),
			),
		).
		Aggregation(
			es.ItemFieldColors,
			newFilterAggregationForFacet(
				es.ItemFieldColors,
				append(extractFiltersExceptForField(es.ItemFieldColors, postFilterMap), metadataFilters...),
			),
		)

	// for _, metadata := range searchFilter.Metadata {
	// 	metadataFilters := extractFiltersExceptForField(metadata.Name, metadataFilterMap)
	// 	search.Aggregation(metadata.Name)
	// }
	//
	return search, postFilterMap
}

// newFilterAggregationForFacet initializes the aggregation with filters applied to the other fields
//  to get the correct facet count for the given field
func newFilterAggregationForFacet(field string, filters []elastic.Query) elastic.Aggregation {
	// We can't use filters aggregation when 0 filters
	if len(filters) == 0 {
		return elastic.NewTermsAggregation().Field(field).Size(30)
	}
	return elastic.NewFiltersAggregation().
		Filters(filters...).
		SubAggregation(field, elastic.NewTermsAggregation().Field(field).Size(30))
}

// func newFilterMetadataAggregationForFacet(field string, filters []elastic.Query) elastic.Aggregation {
// 	// We can't use filters aggregation when 0 filters
// 	if len(filters) == 0 {
// 		return elastic.NewNestedAggregation().Path(es.MetadataNameFullPath).SubAggregation(field)
// 	}
// 	return elastic.NewFiltersAggregation().
// 		Filters(filters...).
// 		SubAggregation(field, elastic.NewTermsAggregation().Field(field).Size(30))
// }

func (s *searchClient) mapAggregationToFacets(ctx context.Context, agg elastic.Aggregations, postFilterMap map[string]elastic.Query) []*Facet {
	var facets []*Facet

	if result, ok := agg.Terms(es.ItemFieldCategoryID); ok {
		isFilterAggregation := len(extractFiltersExceptForField(es.ItemFieldCategoryID, postFilterMap)) > 0
		if isFilterAggregation {
			result, _ = result.Buckets[0].Terms(es.ItemFieldCategoryID)
		}
		itemCategories, err := s.dbClient.GetAllItemCategoriesWithParent(ctx)
		if err != nil {
			logging.Logger(ctx).Error("failed to get item categories", zap.Error(err))
		}
		idNameMap := make(map[string]string)
		for _, itemCategory := range itemCategories {
			idNameMap[itemCategory.ID] = strings.Join(itemCategory.CategoryNames(), " > ")
		}
		facet := newFacetFromBucketKeyItems(result, FacetTypeCategoryIDs, idNameMap)
		if facet.TotalCount > 0 {
			facets = append(facets, facet)
		}
	}

	if result, ok := agg.Terms(es.ItemFieldBrandName); ok {
		isFilterAggregation := len(extractFiltersExceptForField(es.ItemFieldBrandName, postFilterMap)) > 0
		if isFilterAggregation {
			result, _ = result.Buckets[0].Terms(es.ItemFieldBrandName)
		}
		idNameMap := make(map[string]string)
		for _, bucket := range result.Buckets {
			keyStr, ok := bucket.Key.(string)
			if !ok {
				continue
			}
			idNameMap[keyStr] = keyStr
		}
		facet := newFacetFromBucketKeyItems(result, FacetTypeBrandNames, idNameMap)
		if facet.TotalCount > 0 {
			facets = append(facets, facet)
		}
	}

	if result, ok := agg.Terms(es.ItemFieldColors); ok {
		isFilterAggregation := len(extractFiltersExceptForField(es.ItemFieldColors, postFilterMap)) > 0
		if isFilterAggregation {
			result, _ = result.Buckets[0].Terms(es.ItemFieldColors)
		}
		idNameMap := make(map[string]string)
		for _, bucket := range result.Buckets {
			keyStr, ok := bucket.Key.(string)
			if !ok {
				continue
			}
			idNameMap[keyStr] = keyStr
		}
		facet := newFacetFromBucketKeyItems(result, FacetTypeColors, idNameMap)
		if facet.TotalCount > 0 {
			facets = append(facets, facet)
		}
	}

	return facets
}
