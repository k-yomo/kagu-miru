package search

import (
	"context"
	"sort"
	"strings"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/pkg/interfaceconv"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

type FacetType int

const (
	FacetTypeCategoryIDs FacetType = iota + 1
	FacetTypeBrandNames
	FacetTypeColors
	FacetTypeMetadata
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

func applyAggregationsAndPostFiltersForFacets(search *elastic.SearchService, searchFilter *gqlmodel.SearchFilter) (*elastic.SearchService, map[string]elastic.Query, map[string]elastic.Query) {
	// filters for facetable fields
	postFilterMap := make(map[string]elastic.Query)
	var postFilters []elastic.Query
	if len(searchFilter.CategoryIds) > 0 {
		var categoryIDs []interface{}
		for _, id := range searchFilter.CategoryIds {
			categoryIDs = append(categoryIDs, id)
		}
		// Since we aggregate for category_id
		filter := elastic.NewTermsQuery(es.ItemFieldCategoryIDs, categoryIDs...)
		postFilterMap[es.ItemFieldCategoryID] = filter
		postFilters = append(postFilters, filter)
	}
	if len(searchFilter.BrandNames) > 0 {
		var brandNames []interface{}
		for _, brandName := range searchFilter.BrandNames {
			brandNames = append(brandNames, brandName)
		}
		filter := elastic.NewTermsQuery(es.ItemFieldBrandName, brandNames...)
		postFilterMap[es.ItemFieldBrandName] = filter
		postFilters = append(postFilters, filter)
	}
	if len(searchFilter.Colors) > 0 {
		var colors []interface{}
		for _, gqlColor := range searchFilter.Colors {
			if color := mapGraphqlItemColorToSearchItemColor(gqlColor); color != "" {
				colors = append(colors, color)
			}
		}
		filter := elastic.NewTermsQuery(es.ItemFieldColors, colors...)
		postFilterMap[es.ItemFieldColors] = filter
		postFilters = append(postFilters, filter)
	}

	var metadataFilters []elastic.Query
	postMetadataFilterMap := make(map[string]elastic.Query)
	if len(searchFilter.Metadata) > 0 {
		for _, metadata := range searchFilter.Metadata {
			if len(metadata.Values) == 0 {
				continue
			}
			boolQuery := elastic.NewBoolQuery().Filter(
				elastic.NewTermQuery(es.MetadataNameFullPath, metadata.Name),
				elastic.NewTermsQuery(es.MetadataValueFullPath, interfaceconv.StringArrayToInterfaceArray(metadata.Values)...),
			)
			filter := elastic.NewNestedQuery(es.ItemFieldMetadata, boolQuery)
			postMetadataFilterMap[metadata.Name] = filter
			metadataFilters = append(metadataFilters, filter)
		}
	}

	postFilterBoolQuery := elastic.NewBoolQuery().Must()
	for _, filter := range append(postFilters, metadataFilters...) {
		postFilterBoolQuery.Filter(filter)
	}
	search.PostFilter(postFilterBoolQuery)

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

	allFilters := append(postFilters, metadataFilters...)
	if len(allFilters) > 0 {
		search.Aggregation(es.ItemFieldMetadata, elastic.NewFiltersAggregation().
			Filters(allFilters...).
			SubAggregation(
				es.ItemFieldMetadata,
				elastic.NewNestedAggregation().
					Path(es.ItemFieldMetadata).
					SubAggregation(
						es.MetadataNameFullPath,
						elastic.NewTermsAggregation().
							Field(es.MetadataNameFullPath).
							SubAggregation(
								es.MetadataValueFullPath,
								elastic.NewTermsAggregation().Field(es.MetadataValueFullPath),
							),
					),
			))
	} else {
		search.Aggregation(es.ItemFieldMetadata, elastic.NewNestedAggregation().
			Path(es.ItemFieldMetadata).
			SubAggregation(
				es.MetadataNameFullPath,
				elastic.NewTermsAggregation().
					Field(es.MetadataNameFullPath).
					SubAggregation(
						es.MetadataValueFullPath,
						elastic.NewTermsAggregation().Field(es.MetadataValueFullPath),
					),
			),
		)
	}

	for _, metadata := range searchFilter.Metadata {
		if len(metadata.Values) == 0 {
			continue
		}
		metadataFilters := extractFiltersExceptForField(metadata.Name, postMetadataFilterMap)
		search.Aggregation(metadata.Name, newFilterMetadataAggregationForFacet(
			metadata.Name,
			append(metadataFilters, postFilters...),
		),
		)
	}

	return search, postFilterMap, postMetadataFilterMap
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

func newFilterMetadataAggregationForFacet(metadataName string, filters []elastic.Query) elastic.Aggregation {
	filters = append(filters, elastic.NewNestedQuery(es.ItemFieldMetadata, elastic.NewTermQuery(es.MetadataNameFullPath, metadataName)))
	return elastic.NewFiltersAggregation().
		Filters(filters...).
		SubAggregation(metadataName, elastic.NewNestedAggregation().
			Path(es.ItemFieldMetadata).
			SubAggregation(metadataName, elastic.NewTermsAggregation().Field(es.MetadataValueFullPath).Size(30)))
}

func (s *searchClient) mapAggregationToFacets(ctx context.Context, agg elastic.Aggregations, postFilterMap map[string]elastic.Query, postMetadataFilterMap map[string]elastic.Query) []*Facet {
	var facets []*Facet

	if result, ok := agg.Terms(es.ItemFieldCategoryID); ok {
		isFilterAggregation := len(extractFiltersExceptForField(es.ItemFieldCategoryID, postFilterMap)) > 0 || len(postMetadataFilterMap) > 0
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
		isFilterAggregation := len(extractFiltersExceptForField(es.ItemFieldBrandName, postFilterMap)) > 0 || len(postMetadataFilterMap) > 0
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
		// facet must have multiple variations if not selected
		if isFilterAggregation || len(facet.Values) > 1 {
			facets = append(facets, facet)
		}
	}

	if result, ok := agg.Terms(es.ItemFieldColors); ok {
		isFilterAggregation := len(extractFiltersExceptForField(es.ItemFieldColors, postFilterMap)) > 0 || len(postMetadataFilterMap) > 0
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
		// facet must have multiple variations if not selected
		if isFilterAggregation || len(facet.Values) > 1 {
			facets = append(facets, facet)
		}
	}

	var metadataFacets []*Facet
	if result, ok := agg.Terms(es.ItemFieldMetadata); ok {
		isFilterAggregation := len(postFilterMap) > 0 || len(postMetadataFilterMap) > 0
		if isFilterAggregation {
			result, _ = result.Buckets[0].Terms(es.ItemFieldMetadata)
		}
		result, _ = result.Terms(es.MetadataNameFullPath)

		for _, bucket := range result.Buckets {
			metadataName, ok := bucket.Key.(string)
			if !ok {
				continue
			}
			// facet for filtered metadata is covered by their individual aggregation
			// NOTE: it might be more efficient to exclude from result with aggregation filter
			if _, ok := postMetadataFilterMap[metadataName]; ok {
				continue
			}

			bucketKeyItems, _ := bucket.Terms(es.MetadataValueFullPath)
			facetValues := make([]*FacetValue, 0, len(bucketKeyItems.Buckets))
			totalCount := bucketKeyItems.SumOfOtherDocCount
			for _, bucket := range bucketKeyItems.Buckets {
				keyStr, ok := bucket.Key.(string)
				if !ok {
					continue
				}
				facetValues = append(facetValues, &FacetValue{
					ID:    keyStr,
					Name:  keyStr,
					Count: int(bucket.DocCount),
				})
				totalCount += bucket.DocCount
			}
			// facet must have multiple variations
			if len(facetValues) > 1 {
				sortMetadataValues(metadataName, facetValues)
				metadataFacets = append(metadataFacets, &Facet{
					Title:      metadataName,
					FacetType:  FacetTypeMetadata,
					Values:     facetValues,
					TotalCount: int(totalCount),
				})
			}
		}
	}

	for metadataName := range postMetadataFilterMap {
		if result, ok := agg.Terms(metadataName); ok {
			result, _ := result.Buckets[0].Terms(metadataName)
			result, _ = result.Terms(metadataName)

			facetValues := make([]*FacetValue, 0, len(result.Buckets))
			totalCount := result.SumOfOtherDocCount
			for _, bucket := range result.Buckets {
				keyStr, ok := bucket.Key.(string)
				if !ok {
					continue
				}
				facetValues = append(facetValues, &FacetValue{
					ID:    keyStr,
					Name:  keyStr,
					Count: int(bucket.DocCount),
				})
				totalCount += bucket.DocCount
			}

			if len(facetValues) > 0 {
				sortMetadataValues(metadataName, facetValues)
				metadataFacets = append(metadataFacets, &Facet{
					Title:      metadataName,
					FacetType:  FacetTypeMetadata,
					Values:     facetValues,
					TotalCount: int(totalCount),
				})
			}
		}
	}

	sort.Slice(metadataFacets, func(i, j int) bool {
		iOrder := es.MetadataNameSortOrderMap[metadataFacets[i].Title]
		jOrder := es.MetadataNameSortOrderMap[metadataFacets[j].Title]
		return iOrder < jOrder
	})

	return append(facets, metadataFacets...)
}

func sortMetadataValues(metadataName string, facetValues []*FacetValue) {
	switch metadataName {
	case es.MetadataNameWidthRange, es.MetadataNameDepthRange, es.MetadataNameHeightRange:
		sort.Slice(facetValues, func(i, j int) bool {
			iOrder := es.MetadataValueLengthSortOrderMap[facetValues[i].ID]
			jOrder := es.MetadataValueLengthSortOrderMap[facetValues[j].ID]
			return iOrder < jOrder
		})
	}
}
