package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"time"

	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
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

var NotFoundErr = errors.New("item not found")

type elasticsearchClient struct {
	itemsIndexName                 string
	itemsQuerySuggestionsIndexName string
	esClient                       *elasticsearch.Client
}

func NewElasticsearchClient(itemsIndexName string, itemsQuerySuggestionsIndexName string, esClient *elasticsearch.Client) Client {
	return &elasticsearchClient{
		itemsIndexName:                 itemsIndexName,
		itemsQuerySuggestionsIndexName: itemsQuerySuggestionsIndexName,
		esClient:                       esClient,
	}
}

type Response struct {
	Items      []*es.Item
	Page       uint64
	TotalPage  uint64
	TotalCount uint64
}

func (c *elasticsearchClient) GetItem(ctx context.Context, id string) (*es.Item, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_GetItem")
	defer span.End()

	response, err := c.esClient.Get(c.itemsIndexName, id, c.esClient.Get.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("esClient.Get: %w", err)
	}

	result := elastic.GetResult{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("json.Decoder.Decode: %w", err)
	}

	if !result.Found {
		return nil, fmt.Errorf("get '%s': %w", id, NotFoundErr)
	}

	var item es.Item
	if err := json.Unmarshal(result.Source, &item); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &item, nil
}

func (c *elasticsearchClient) SearchItems(ctx context.Context, input *gqlmodel.SearchInput) (*Response, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_SearchItems")
	defer span.End()

	go func() {
		if err := c.insertQuerySuggestion(context.Background(), input.Query); err != nil {
			logging.Logger(ctx).Error("insertQuerySuggestion failed", zap.Error(err))
		}
	}()

	esQuery, err := buildSearchQuery(input)
	if err != nil {
		return nil, fmt.Errorf("buildSearchQuery: %w", err)
	}
	searchResult, err := c.search(ctx, c.itemsIndexName, esQuery)
	if err != nil {
		return nil, fmt.Errorf("elasticsearchClient.search: %w", err)
	}

	return &Response{
		Items:      mapElasticsearchHitsToItems(ctx, searchResult.Hits.Hits),
		Page:       calcElasticSearchPage(input.Page) + 1,
		TotalPage:  calcTotalPage(uint64(searchResult.Hits.TotalHits.Value), 100),
		TotalCount: uint64(searchResult.Hits.TotalHits.Value),
	}, nil
}

func buildSearchQuery(input *gqlmodel.SearchInput) (io.Reader, error) {
	var mustQueries []esquery.Mappable
	if input.Query != "" {
		mustQueries = append(mustQueries, esquery.MultiMatch(input.Query).
			Type(esquery.MatchTypeMostFields).
			Fields(
				xesquery.Boost(es.ItemFieldName, 20),
				xesquery.Boost(es.ItemFieldBrandName, 5),
				xesquery.Boost(es.ItemFieldCategoryNames, 5),
				xesquery.Boost(es.ItemFieldColors, 5),
				es.ItemFieldDescription,
			))
	} else {
		mustQueries = append(mustQueries, esquery.MatchAll())
	}

	boolQuery := esquery.Bool().Must(mustQueries...)

	if len(input.Filter.CategoryIds) > 0 {
		var categoryIDs []interface{}
		for _, id := range input.Filter.CategoryIds {
			categoryIDs = append(categoryIDs, id)
		}
		boolQuery.Filter(esquery.Terms(es.ItemFieldCategoryIDs, categoryIDs...))
	}

	if len(input.Filter.Platforms) > 0 {
		var platforms []interface{}
		for _, filterPlatform := range input.Filter.Platforms {
			platform, err := mapGraphqlPlatformToPlatform(filterPlatform)
			if err != nil {
				return nil, fmt.Errorf("platform comversion faild: %w", err)
			}
			platforms = append(platforms, platform)
		}
		boolQuery.Filter(esquery.Terms(es.ItemFieldPlatform, platforms...))
	}

	if len(input.Filter.Colors) > 0 {
		var colors []interface{}
		for _, gqlColor := range input.Filter.Colors {
			if color := mapGraphqlItemColorToSearchItemColor(gqlColor); color != "" {
				colors = append(colors, color)
			}
		}
		boolQuery.Filter(esquery.Terms(es.ItemFieldColors, colors...))
	}

	if input.Filter.MinPrice != nil && input.Filter.MaxPrice != nil {
		boolQuery.Filter(
			esquery.Range(es.ItemFieldPrice).
				Gte(*input.Filter.MinPrice).
				Lte(*input.Filter.MaxPrice),
		)
	} else if input.Filter.MinPrice != nil {
		boolQuery.Filter(esquery.Range(es.ItemFieldPrice).Gte(*input.Filter.MinPrice))
	} else if input.Filter.MaxPrice != nil {
		boolQuery.Filter(esquery.Range(es.ItemFieldPrice).Lte(*input.Filter.MaxPrice))
	}

	if input.Filter.MinRating != nil {
		boolQuery.Filter(esquery.Range(es.ItemFieldAverageRating).Gte(*input.Filter.MinRating))
	}

	esQuery := esquery.Search().Query(esquery.CustomQuery(map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": boolQuery.Map(),
			"functions": []map[string]interface{}{
				{
					"gauss": map[string]interface{}{
						es.ItemFieldAverageRating: map[string]interface{}{
							"origin": 5,
							"offset": 1,
							"scale":  1,
							"decay":  0.4,
						},
					},
				},
				{
					"field_value_factor": map[string]string{
						"field": es.ItemFieldReviewCount,
					},
				},
			},
			"max_boost": 3,
			// "min_score": 0,
		},
	}))

	switch input.SortType {
	case gqlmodel.SearchSortTypePriceAsc:
		esQuery.Sort(es.ItemFieldPrice, esquery.OrderAsc)
	case gqlmodel.SearchSortTypePriceDesc:
		esQuery.Sort(es.ItemFieldPrice, esquery.OrderDesc)
	case gqlmodel.SearchSortTypeReviewCount:
		esQuery.Sort(es.ItemFieldReviewCount, esquery.OrderDesc).Sort(es.ItemFieldAverageRating, esquery.OrderDesc)
	case gqlmodel.SearchSortTypeRating:
		esQuery.Sort(es.ItemFieldAverageRating, esquery.OrderDesc).Sort(es.ItemFieldReviewCount, esquery.OrderDesc)
	default:
		esQuery.Sort("_score", esquery.OrderDesc)
	}

	pageSize := defaultPageSize
	if input.PageSize != nil {
		pageSize = uint64(math.Min(float64(*input.PageSize), float64(maxPageSize)))
	}
	esQuery.
		SourceIncludes(es.AllItemFields...).
		From(calcElasticSearchPage(input.Page) * pageSize).
		Size(pageSize)

	esQueryJSON, err := esQuery.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("esQuery.MarshalJSON(): %w", err)
	}

	return bytes.NewReader(esQueryJSON), nil
}

func (c *elasticsearchClient) GetSimilarItems(ctx context.Context, input *gqlmodel.GetSimilarItemsInput, itemCategoryID string) (*Response, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_GetSimilarItems")
	defer span.End()

	mustQueries := []esquery.Mappable{esquery.CustomQuery(map[string]interface{}{
		"more_like_this": map[string]interface{}{
			"fields": []string{
				xesquery.Boost(es.ItemFieldName, 30),
				xesquery.Boost(es.ItemFieldBrandName, 10),
				xesquery.Boost(es.ItemFieldCategoryNames, 5),
				xesquery.Boost(es.ItemFieldColors, 5),
				es.ItemFieldDescription,
			},
			"like": map[string]interface{}{
				"_index": c.itemsIndexName,
				"_id":    input.ItemID,
			},
		},
	})}
	boolQuery := esquery.Bool().Must(mustQueries...)
	boolQuery.Filter(esquery.Terms(es.ItemFieldCategoryIDs, itemCategoryID))
	esQuery := esquery.Search().Query(esquery.CustomQuery(map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": boolQuery.Map(),
			"functions": []map[string]interface{}{
				{
					"gauss": map[string]interface{}{
						es.ItemFieldAverageRating: map[string]interface{}{
							"origin": 5,
							"offset": 1,
							"scale":  1,
							"decay":  0.4,
						},
					},
				},
				{
					"field_value_factor": map[string]string{
						"field": es.ItemFieldReviewCount,
					},
				},
			},
			"max_boost": 1,
		},
	}))

	esQueryJSON, err := esQuery.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("esQuery.MarshalJSON(): %w", err)
	}

	searchResult, err := c.search(ctx, c.itemsIndexName, bytes.NewReader(esQueryJSON))
	if err != nil {
		return nil, fmt.Errorf("elasticsearchClient.search: %w", err)
	}

	return &Response{
		Items:      mapElasticsearchHitsToItems(ctx, searchResult.Hits.Hits),
		Page:       calcElasticSearchPage(input.Page) + 1,
		TotalPage:  calcTotalPage(uint64(searchResult.Hits.TotalHits.Value), 100),
		TotalCount: uint64(searchResult.Hits.TotalHits.Value),
	}, nil
}

func (c *elasticsearchClient) GetQuerySuggestions(ctx context.Context, query string) ([]string, error) {
	ctx, span := otel.Tracer("").Start(ctx, "search.elasticsearchClient_GetQuerySuggestions")
	defer span.End()

	const aggregationTerm = "queries"
	esQuery, err := esquery.Search().
		Query(
			esquery.Bool().Should(
				esquery.Match("query.autocomplete", query),
				esquery.Match("query.readingform", query).Fuzziness("AUTO").Operator(esquery.OperatorAnd),
			),
		).
		Aggs(
			esquery.TermsAgg(aggregationTerm, "query").
				Order(map[string]string{"_count": string(esquery.OrderDesc)}).
				Size(10),
		).
		Size(0).
		MarshalJSON()

	if err != nil {
		return nil, fmt.Errorf("esquery.MarshalJSON(): %w", err)
	}

	searchResult, err := c.search(ctx, c.itemsQuerySuggestionsIndexName, bytes.NewReader(esQuery))
	if err != nil {
		return nil, fmt.Errorf("c.search: %w", err)
	}

	bucketKeyItems, ok := searchResult.Aggregations.Terms(aggregationTerm)
	if !ok {
		return nil, fmt.Errorf("aggregation term '%s' not found in the search result, aggs: %+v", aggregationTerm, searchResult.Aggregations)
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

func (c *elasticsearchClient) search(ctx context.Context, indexName string, esQuery io.Reader) (*elastic.SearchResult, error) {
	response, err := c.esClient.Search(
		c.esClient.Search.WithContext(ctx),
		c.esClient.Search.WithIndex(indexName),
		c.esClient.Search.WithBody(esQuery),
		c.esClient.Search.WithRequestCache(true),
	)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read response body failed, response: %s", response.String())
	}
	if response.StatusCode >= 400 {
		return nil, errors.Errorf("search request to elasticsearch failed with status %s, body: %s", response.Status(), body)
	}

	searchResult := elastic.SearchResult{}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, errors.Wrapf(err, "unmarshal response body json failed, body: %s", body)
	}
	return &searchResult, nil
}

func (c *elasticsearchClient) insertQuerySuggestion(ctx context.Context, query string) error {
	querySuggestion := es.QuerySuggestion{Query: query, CreatedAt: time.Now()}
	querySuggestionJSON, err := json.Marshal(querySuggestion)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	response, err := c.esClient.Index(
		c.itemsQuerySuggestionsIndexName,
		bytes.NewBuffer(querySuggestionJSON),
		c.esClient.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("esClient.Delete failed: %w", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("index failed, body: %s: %w", body, err)
	}

	return nil
}

func calcTotalPage(totalItems, pageSize uint64) uint64 {
	if totalItems == 0 {
		return 1
	}
	total := totalItems + pageSize - 1
	return total / pageSize
}

func calcElasticSearchPage(inputPage *int) uint64 {
	page := defaultPage
	if inputPage != nil && *inputPage > 1 {
		page = uint64(*inputPage) - 1
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
