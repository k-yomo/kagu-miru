package search

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/db"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"github.com/k-yomo/kagu-miru/backend/pkg/xesquery"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	defaultPage     int = 0
	defaultPageSize int = 100
	maxPageSize     int = 1000
)

var NotFoundErr = errors.New("item not found")

type Client interface {
	GetItem(ctx context.Context, id string) (*es.Item, error)
	SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error)
	GetSimilarItems(ctx context.Context, input *gqlmodel.GetSimilarItemsInput, itemCategoryID string) (*Response, error)
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

func (s *searchClient) GetItem(ctx context.Context, id string) (*es.Item, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_GetItem")
	defer span.End()

	resp, err := s.esClient.Get().Index(s.itemsIndexName).Id(id).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("esClient.Get: %w", err)
	}

	if !resp.Found {
		return nil, fmt.Errorf("get '%s': %w", id, NotFoundErr)
	}

	var item es.Item
	if err := json.Unmarshal(resp.Source, &item); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &item, nil
}

func (s *searchClient) SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_SearchItems")
	defer span.End()

	go func() {
		if err := s.insertQuerySuggestion(context.Background(), input.Query); err != nil {
			logging.Logger(ctx).Error("insertQuerySuggestion failed", zap.Error(err))
		}
	}()

	searchQuery, postFilterMap, err := buildSearchQuery(input)
	if err != nil {
		return nil, fmt.Errorf("buildSearchQuery: %w", err)
	}
	postFilterBoolQuery := elastic.NewBoolQuery().Must()
	for _, filter := range postFilterMap {
		postFilterBoolQuery.Filter(filter)
	}
	pageSize := defaultPageSize
	if input.PageSize != nil {
		pageSize = int(math.Min(float64(*input.PageSize), float64(maxPageSize)))
	}

	search := s.esClient.Search().
		Index(s.itemsIndexName).
		Query(searchQuery)
	search = addAggregationsForFacets(search, postFilterMap)
	resp, err := search.
		PostFilter(postFilterBoolQuery).
		SortBy(getSorters(input.SortType)...).
		From(calcElasticSearchPage(input.Page) * pageSize).
		Size(pageSize).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}

	return &Response{
		Items:      mapElasticsearchHitsToItems(ctx, resp.Hits.Hits),
		Facets:     s.mapAggregationToFacets(ctx, resp.Aggregations, postFilterMap),
		Page:       calcElasticSearchPage(input.Page) + 1,
		TotalPage:  calcTotalPage(int(resp.Hits.TotalHits.Value), 100),
		TotalCount: int(resp.Hits.TotalHits.Value),
	}, nil
}

func buildSearchQuery(input *gqlmodel.SearchInput) (query elastic.Query, postFilterMap map[string]elastic.Query, err error) {
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
				return nil, nil, fmt.Errorf("platform comversion faild: %w", err)
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

	// filters for facetable fields
	postFilterQueryMap := make(map[string]elastic.Query)
	if len(input.Filter.CategoryIds) > 0 {
		var categoryIDs []interface{}
		for _, id := range input.Filter.CategoryIds {
			categoryIDs = append(categoryIDs, id)
		}
		// Since we aggregate for category_id
		postFilterQueryMap[es.ItemFieldCategoryID] = elastic.NewTermsQuery(es.ItemFieldCategoryIDs, categoryIDs...)
	}
	if len(input.Filter.BrandNames) > 0 {
		var brandNames []interface{}
		for _, brandName := range input.Filter.BrandNames {
			brandNames = append(brandNames, brandName)
		}
		postFilterQueryMap[es.ItemFieldBrandName] = elastic.NewTermsQuery(es.ItemFieldBrandName, brandNames...)
	}
	if len(input.Filter.Colors) > 0 {
		var colors []interface{}
		for _, gqlColor := range input.Filter.Colors {
			if color := mapGraphqlItemColorToSearchItemColor(gqlColor); color != "" {
				colors = append(colors, color)
			}
		}
		postFilterQueryMap[es.ItemFieldColors] = elastic.NewTermsQuery(es.ItemFieldColors, colors...)
	}

	return searchQuery, postFilterQueryMap, nil
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

func (s *searchClient) GetSimilarItems(ctx context.Context, input *gqlmodel.GetSimilarItemsInput, itemCategoryID string) (*Response, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_GetSimilarIts")
	defer span.End()

	boolQuery := elastic.NewBoolQuery().Must(
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
	).Filter(elastic.NewTermsQuery(es.ItemFieldCategoryIDs, itemCategoryID))

	functionScoreQuery := elastic.NewFunctionScoreQuery().
		Query(boolQuery).
		AddScoreFunc(elastic.NewGaussDecayFunction().FieldName(es.ItemFieldAverageRating).Origin(5).Offset(1).Scale(1).Decay(0.4)).
		AddScoreFunc(elastic.NewFieldValueFactorFunction().Field(es.ItemFieldReviewCount)).
		MaxBoost(1)

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
		Items:      mapElasticsearchHitsToItems(ctx, resp.Hits.Hits),
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

func mapGraphqlPlatformToPlatform(platform gqlmodel.ItemSellingPlatform) (xitem.Platform, error) {
	switch platform {
	case gqlmodel.ItemSellingPlatformRakuten:
		return xitem.PlatformRakuten, nil
	case gqlmodel.ItemSellingPlatformYahooShopping:
		return xitem.PlatformYahooShopping, nil
	case gqlmodel.ItemSellingPlatformPaypayMall:
		return xitem.PlatformPayPayMall, nil
	default:
		return "", fmt.Errorf("unknown platform %s", platform.String())
	}
}

func mapGraphqlItemColorToSearchItemColor(color gqlmodel.ItemColor) string {
	switch color {
	case gqlmodel.ItemColorWhite:
		return "ホワイト"
	case gqlmodel.ItemColorYellow:
		return "イエロー"
	case gqlmodel.ItemColorOrange:
		return "オレンジ"
	case gqlmodel.ItemColorPink:
		return "ピンク"
	case gqlmodel.ItemColorRed:
		return "レッド"
	case gqlmodel.ItemColorBeige:
		return "ベージュ"
	case gqlmodel.ItemColorSilver:
		return "シルバー"
	case gqlmodel.ItemColorGold:
		return "ゴールド"
	case gqlmodel.ItemColorGray:
		return "グレー"
	case gqlmodel.ItemColorPurple:
		return "パープル"
	case gqlmodel.ItemColorBrown:
		return "ブラウン"
	case gqlmodel.ItemColorGreen:
		return "グリーン"
	case gqlmodel.ItemColorBlue:
		return "ブルー"
	case gqlmodel.ItemColorBlack:
		return "ブラック"
	case gqlmodel.ItemColorNavy:
		return "ネイビー"
	case gqlmodel.ItemColorKhaki:
		return "カーキ"
	case gqlmodel.ItemColorWineRed:
		return "ワインレッド"
	case gqlmodel.ItemColorTransparent:
		return "透明"
	default:
		return ""
	}
}

func mapElasticsearchHitsToItems(ctx context.Context, hits []*elastic.SearchHit) []*es.Item {
	items := make([]*es.Item, 0, len(hits))
	for _, hit := range hits {
		var item es.Item
		if err := json.Unmarshal(hit.Source, &item); err != nil {
			logging.Logger(ctx).Error("Failed to unmarshal hit.Source into es.Item", zap.String("source", string(hit.Source)))
			continue
		}

		items = append(items, &item)
	}

	return items
}
