package search

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"go.uber.org/zap"

	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"

	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/db"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/pkg/xesquery"
	"github.com/olivere/elastic/v7"
	"go.opentelemetry.io/otel"
)

const (
	defaultPage     int = 0
	defaultPageSize int = 100
	maxPageSize     int = 1000

	minRequiredHitsForQuerySuggestion = 100
)

type Client interface {
	SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error)
	GetSimilarItems(ctx context.Context, input *gqlmodel.GetSimilarItemsInput, item *xspanner.Item) (*Response, error)
	GetQuerySuggestions(ctx context.Context, query string) ([]string, error)
}

type searchClient struct {
	itemsIndexName                 string
	itemsQuerySuggestionsIndexName string
	esClient                       *elastic.Client
	dbClient                       db.Client
}

func NewSearchClient(
	itemsIndexName string,
	itemsQuerySuggestionsIndexName string,
	esClient *elastic.Client,
	dbClient db.Client,
) Client {
	return &searchClient{
		itemsIndexName:                 itemsIndexName,
		itemsQuerySuggestionsIndexName: itemsQuerySuggestionsIndexName,
		esClient:                       esClient,
		dbClient:                       dbClient,
	}
}

type Response struct {
	Items      []*es.Item
	Facets     []*Facet
	Page       int
	TotalPage  int
	TotalCount int
}

func (s *searchClient) SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_SearchItems")
	defer span.End()

	searchQuery, err := buildSearchQuery(input)
	if err != nil {
		return nil, fmt.Errorf("buildSearchQuery: %w", err)
	}
	pageSize := defaultPageSize
	if input.PageSize != nil {
		pageSize = int(math.Min(float64(*input.PageSize), float64(maxPageSize)))
	}

	search := s.esClient.Search().
		Index(s.itemsIndexName).
		Query(searchQuery)
	search, postFilterMap, postMetadataFilterMap := applyAggregationsAndPostFiltersForFacets(search, input.Filter)
	resp, err := search.
		SortBy(getSorters(input.SortType)...).
		From(calcElasticSearchPage(input.Page) * pageSize).
		Size(pageSize).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}

	if resp.Hits.TotalHits.Value >= minRequiredHitsForQuerySuggestion {
		go func() {
			if err := s.insertQuerySuggestion(context.Background(), input.Query); err != nil {
				logging.Logger(ctx).Error("insertQuerySuggestion failed", zap.Error(err))
			}
		}()
	}

	return &Response{
		Items:      dedupItems(mapElasticsearchHitsToItems(ctx, resp.Hits.Hits)),
		Facets:     s.mapAggregationToFacets(ctx, resp.Aggregations, postFilterMap, postMetadataFilterMap),
		Page:       calcElasticSearchPage(input.Page) + 1,
		TotalPage:  calcTotalPage(int(resp.Hits.TotalHits.Value), 100),
		TotalCount: int(resp.Hits.TotalHits.Value),
	}, nil
}

func buildSearchQuery(input *gqlmodel.SearchInput) (query elastic.Query, err error) {
	var mustQueries []elastic.Query
	if input.Query != "" {
		mustQueries = append(mustQueries, elastic.NewMultiMatchQuery(
			input.Query,
			xesquery.Boost(es.ItemFieldName, 20),
			xesquery.Boost(es.ItemFieldBrandName, 5),
			xesquery.Boost(es.ItemFieldCategoryNames, 5),
			xesquery.Boost(es.ItemFieldColors, 5),
			es.ItemFieldDescription,
		).Type("cross_fields").Operator("AND"))
	} else {
		mustQueries = append(mustQueries, elastic.NewMatchAllQuery())
	}

	boolQuery := elastic.NewBoolQuery().Must(mustQueries...)

	if len(input.Filter.Platforms) > 0 {
		var platforms []interface{}
		for _, filterPlatform := range input.Filter.Platforms {
			platform, err := mapGraphqlPlatformToPlatform(filterPlatform)
			if err != nil {
				return nil, fmt.Errorf("platform comversion faild: %w", err)
			}
			platforms = append(platforms, platform)
		}
		boolQuery.Filter(elastic.NewTermsQuery(es.ItemFieldPlatform, platforms...))
	}

	if input.Filter.MinPrice != nil && input.Filter.MaxPrice != nil {
		boolQuery.Filter(
			elastic.NewRangeQuery(es.ItemFieldPrice).
				Gte(*input.Filter.MinPrice).
				Lte(*input.Filter.MaxPrice),
		)
	} else if input.Filter.MinPrice != nil {
		boolQuery.Filter(elastic.NewRangeQuery(es.ItemFieldPrice).Gte(*input.Filter.MinPrice))
	} else if input.Filter.MaxPrice != nil {
		boolQuery.Filter(elastic.NewRangeQuery(es.ItemFieldPrice).Lte(*input.Filter.MaxPrice))
	}

	if input.Filter.MinRating != nil {
		boolQuery.Filter(elastic.NewRangeQuery(es.ItemFieldAverageRating).Gte(*input.Filter.MinRating))
	}

	searchQuery := elastic.NewFunctionScoreQuery().Query(boolQuery).
		AddScoreFunc(elastic.NewGaussDecayFunction().FieldName(es.ItemFieldAverageRating).Origin(5).Offset(1).Scale(1).Decay(0.4)).
		AddScoreFunc(elastic.NewFieldValueFactorFunction().Field(es.ItemFieldReviewCount)).
		MaxBoost(3)

	return searchQuery, nil
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

func getSorters(sortType gqlmodel.SearchSortType) []elastic.Sorter {
	var sorters []elastic.Sorter
	switch sortType {
	case gqlmodel.SearchSortTypePriceAsc:
		sorters = []elastic.Sorter{elastic.NewFieldSort(es.ItemFieldPrice).Asc()}
	case gqlmodel.SearchSortTypePriceDesc:
		sorters = []elastic.Sorter{elastic.NewFieldSort(es.ItemFieldPrice).Desc()}
	case gqlmodel.SearchSortTypeReviewCount:
		sorters = []elastic.Sorter{
			elastic.NewFieldSort(es.ItemFieldReviewCount).Desc(),
			elastic.NewFieldSort(es.ItemFieldAverageRating).Desc(),
		}
	case gqlmodel.SearchSortTypeRating:
		sorters = []elastic.Sorter{
			elastic.NewFieldSort(es.ItemFieldAverageRating).Desc(),
			elastic.NewFieldSort(es.ItemFieldReviewCount).Desc(),
		}
	default:
		sorters = []elastic.Sorter{elastic.NewScoreSort().Desc()}
	}

	return sorters
}

func (s *searchClient) GetSimilarItems(ctx context.Context, input *gqlmodel.GetSimilarItemsInput, item *xspanner.Item) (*Response, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_GetSimilarIts")
	defer span.End()

	boolQuery := elastic.NewBoolQuery().
		Must(
			elastic.NewMoreLikeThisQuery().Field(
				xesquery.Boost(es.ItemFieldName, 30),
				xesquery.Boost(es.ItemFieldBrandName, 10),
				xesquery.Boost(es.ItemFieldCategoryNames, 5),
				xesquery.Boost(es.ItemFieldColors, 5),
				es.ItemFieldDescription,
			).LikeItems(
				elastic.NewMoreLikeThisQueryItem().
					Index(s.itemsIndexName).
					Id(input.ItemID),
			),
		).
		// MustNot(elastic.NewTermQuery(es.ItemFieldGroupID, item.GroupID)).
		Filter(elastic.NewTermQuery(es.ItemFieldCategoryIDs, item.CategoryID))
	// This is temp implementation while migration since eventually item must have group id
	if item.GroupID.Valid {
		boolQuery.MustNot(elastic.NewTermQuery(es.ItemFieldGroupID, item.GroupID))
	}

	functionScoreQuery := elastic.NewFunctionScoreQuery().
		Query(boolQuery).
		AddScoreFunc(
			elastic.NewGaussDecayFunction().
				FieldName(es.ItemFieldPrice).Origin(item.Price).
				Offset(int(math.Round(float64(item.Price) * 0.2))).
				Scale(int(math.Round(float64(item.Price) * 0.2))).
				Decay(0.5),
		).
		AddScoreFunc(
			elastic.NewGaussDecayFunction().
				FieldName(es.ItemFieldAverageRating).
				Origin(5).
				Offset(1).
				Scale(1).
				Decay(0.5),
		).
		AddScoreFunc(elastic.NewFieldValueFactorFunction().Field(es.ItemFieldReviewCount)).
		MaxBoost(3)

	pageSize := defaultPageSize
	if input.PageSize != nil {
		pageSize = int(math.Min(float64(*input.PageSize), float64(maxPageSize)))
	}
	resp, err := s.esClient.Search().
		Index(s.itemsIndexName).
		Query(functionScoreQuery).
		SortBy(elastic.NewScoreSort()).
		From(calcElasticSearchPage(input.Page) * pageSize).
		Size(pageSize).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}

	return &Response{
		Items:      dedupItems(mapElasticsearchHitsToItems(ctx, resp.Hits.Hits)),
		Page:       calcElasticSearchPage(input.Page) + 1,
		TotalPage:  calcTotalPage(int(resp.Hits.TotalHits.Value), 100),
		TotalCount: int(resp.Hits.TotalHits.Value),
	}, nil
}

func (s *searchClient) GetQuerySuggestions(ctx context.Context, query string) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_GetQuerySuggestions")
	defer span.End()

	const aggregationTerm = "queries"
	boolQuery := elastic.NewBoolQuery().Should(
		elastic.NewMatchQuery("query.autocomplete", query),
		elastic.NewMatchQuery("query.readingform", query).Fuzziness("AUTO").Operator("AND"),
	)

	aggregation := elastic.NewTermsAggregation().
		Field("query").
		Order("_count", false).
		Size(10)

	resp, err := s.esClient.Search().
		Index(s.itemsQuerySuggestionsIndexName).
		Query(boolQuery).
		Aggregation(aggregationTerm, aggregation).
		Size(0).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}

	bucketKeyItems, ok := resp.Aggregations.Terms(aggregationTerm)
	if !ok {
		return nil, fmt.Errorf("aggregation term '%s' not found in the search result, aggs: %+v", aggregationTerm, resp.Aggregations)
	}

	suggestedQueries := make([]string, 0, len(bucketKeyItems.Buckets))
	for _, bucket := range bucketKeyItems.Buckets {
		if bucket.Key == query {
			continue
		}
		suggestedQueries = append(suggestedQueries, bucket.Key.(string))
	}

	return suggestedQueries, nil
}

func (s *searchClient) insertQuerySuggestion(ctx context.Context, query string) error {
	resp, err := s.esClient.Index().
		Index(s.itemsQuerySuggestionsIndexName).
		BodyJson(es.QuerySuggestion{Query: query, CreatedAt: time.Now()}).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("esClient.Index failed: %w", err)
	}
	if resp.Status >= 400 {
		return fmt.Errorf("index failed, body: %s: %w", resp.Result, err)
	}

	return nil
}

func calcTotalPage(totalItems, pageSize int) int {
	if totalItems == 0 {
		return 1
	}
	total := totalItems + pageSize - 1
	return total / pageSize
}

func calcElasticSearchPage(inputPage *int) int {
	page := defaultPage
	if inputPage != nil && *inputPage > 1 {
		page = int(*inputPage) - 1
	}
	return page
}

func dedupItems(items []*es.Item) []*es.Item {
	groupIDMap := make(map[string]bool)
	dedupedItems := make([]*es.Item, 0, len(items))
	for _, item := range items {
		if item.GroupID != "" && groupIDMap[item.GroupID] {
			continue
		}
		dedupedItems = append(dedupedItems, item)
		groupIDMap[item.GroupID] = true
	}
	return dedupedItems
}
