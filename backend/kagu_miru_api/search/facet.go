package search

import (
	"context"
	"strings"

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

func addAggregationsForFacets(search *elastic.SearchService, postFilterMap map[string]elastic.Query) *elastic.SearchService {
	return search.
		Aggregation(es.ItemFieldCategoryID, newFilterAggregationForFacet(es.ItemFieldCategoryID, postFilterMap)).
		Aggregation(es.ItemFieldBrandName, newFilterAggregationForFacet(es.ItemFieldBrandName, postFilterMap)).
		Aggregation(es.ItemFieldColors, newFilterAggregationForFacet(es.ItemFieldColors, postFilterMap))
}

func newFilterAggregationForFacet(field string, postFilterMap map[string]elastic.Query) elastic.Aggregation {
	filters := extractFiltersExceptForField(field, postFilterMap)
	// We can't use filters aggregation when 0 filters
	if len(filters) == 0 {
		return elastic.NewTermsAggregation().Field(field).Size(30)
	}
	return elastic.NewFiltersAggregation().
		Filters(filters...).
		SubAggregation(field, elastic.NewTermsAggregation().Field(field).Size(30))
}

func extractFiltersExceptForField(field string, filterMap map[string]elastic.Query) []elastic.Query {
	var filters []elastic.Query
	for filteredField, filter := range filterMap {
		if filteredField != field {
			filters = append(filters, filter)
		}
	}
	return filters
}

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
