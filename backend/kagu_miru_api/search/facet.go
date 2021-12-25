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
		facetValues = append(facetValues, &FacetValue{
			ID:    bucket.Key.(string),
			Name:  idNameMap[bucket.Key.(string)],
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

func addAggregationsForFacets(search *elastic.SearchService) *elastic.SearchService {
	return search.
		Aggregation(es.ItemFieldCategoryID, elastic.NewTermsAggregation().Field(es.ItemFieldCategoryID).Size(30)).
		Aggregation(es.ItemFieldBrandName, elastic.NewTermsAggregation().Field(es.ItemFieldBrandName).Size(30)).
		Aggregation(es.ItemFieldColors, elastic.NewTermsAggregation().Field(es.ItemFieldColors).Size(30))
}

func (s *searchClient) mapAggregationToFacets(ctx context.Context, agg elastic.Aggregations) []*Facet {
	var facets []*Facet

	if result, ok := agg.Terms(es.ItemFieldCategoryID); ok {
		itemCategories, err := s.dbClient.GetAllItemCategoriesWithParent(ctx)
		if err != nil {
			logging.Logger(ctx).Error("failed to get item categories", zap.Error(err))
		}
		idNameMap := make(map[string]string)
		for _, itemCategory := range itemCategories {
			idNameMap[itemCategory.ID] = strings.Join(itemCategory.CategoryNames(), " > ")
		}
		facets = append(facets, newFacetFromBucketKeyItems(result, FacetTypeCategoryIDs, idNameMap))
	}

	if result, ok := agg.Terms(es.ItemFieldBrandName); ok {
		idNameMap := make(map[string]string)
		for _, bucket := range result.Buckets {
			idNameMap[bucket.Key.(string)] = bucket.Key.(string)
		}
		facets = append(facets, newFacetFromBucketKeyItems(result, FacetTypeBrandNames, idNameMap))
	}

	if result, ok := agg.Terms(es.ItemFieldColors); ok {
		idNameMap := make(map[string]string)
		for _, bucket := range result.Buckets {
			idNameMap[bucket.Key.(string)] = bucket.Key.(string)
		}
		facets = append(facets, newFacetFromBucketKeyItems(result, FacetTypeColors, idNameMap))
	}

	return facets
}
