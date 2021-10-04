package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/k-yomo/kagu-miru/pkg/xesquery"

	"github.com/aquasecurity/esquery"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/k-yomo/kagu-miru/internal/es"
	"github.com/k-yomo/kagu-miru/pkg/logging"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type Client interface {
	SearchItems(ctx context.Context, req *Request) (*Response, error)
	GetQuerySuggestions(ctx context.Context, query string) ([]string, error)
}

type client struct {
	itemsIndexName                 string
	itemsQuerySuggestionsIndexName string
	esClient                       *elasticsearch.Client
}

func NewSearchClient(itemsIndexName string, itemsQuerySuggestionsIndexName string, esClient *elasticsearch.Client) Client {
	return &client{
		itemsIndexName:                 itemsIndexName,
		itemsQuerySuggestionsIndexName: itemsQuerySuggestionsIndexName,
		esClient:                       esClient,
	}
}

const DefaultPage uint64 = 1
const defaultPageSize uint64 = 100

type SortType int

const (
	SortTypeBestMatch SortType = iota
	SortTypePriceAsc
	SortTypePriceDesc
	SortTypeReviewCount
	SortTypeRating
)

type Request struct {
	Query    string
	SortType SortType
	Page     uint64
	PageSize *int
}

type Response struct {
	Items     []*es.Item
	Page      uint64
	TotalPage uint64
}

func (c *client) SearchItems(ctx context.Context, req *Request) (*Response, error) {
	ctx, span := otel.Tracer("search").Start(ctx, "search.Client_SearchItems")
	defer span.End()

	go func() {
		if err := c.insertQuerySuggestion(context.Background(), req.Query); err != nil {
			logging.Logger(ctx).Error("insertQuerySuggestion failed", zap.Error(err))
		}
	}()

	var page uint64
	if req.Page > 1 {
		page = req.Page - 1 // page starts from 0 in elasticsearch
	}
	pageSize := defaultPageSize
	if req.PageSize != nil {
		pageSize = uint64(*req.PageSize)
	}

	esQuery, err := buildSearchQuery(req.Query, req.SortType, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("buildSearchQuery: %w", err)
	}
	searchResult, err := c.search(ctx, c.itemsIndexName, esQuery)
	if err != nil {
		return nil, fmt.Errorf("client.search: %w", err)
	}

	var items []*es.Item
	for _, hit := range searchResult.Hits.Hits {
		var item es.Item
		if err := json.Unmarshal(hit.Source, &item); err != nil {
			logging.Logger(ctx).Error("Failed to unmarshal hit.Source into es.Item", zap.String("source", string(hit.Source)))
			continue
		}

		items = append(items, &item)
	}

	return &Response{
		Items:     items,
		Page:      req.Page,
		TotalPage: calcTotalPage(uint64(searchResult.Hits.TotalHits.Value), 100),
	}, nil
}

func buildSearchQuery(query string, sortType SortType, page, pageSize uint64) (io.Reader, error) {
	// queryTokens := strings.Fields(query)
	// if len(queryTokens) > 1 {
	// 	queryTokens[0] = xesquery.Boost(queryTokens[0], 10)
	// }
	// query = strings.Join(queryTokens, " ")
	esQuery := esquery.Search().
		Query(
			esquery.Bool().
				Must(
					esquery.MultiMatch(query).
						Type(esquery.MatchTypeMostFields).
						Fields(
							xesquery.Boost(es.ItemFieldName, 5),
							es.ItemFieldDescription,
							xesquery.Boost(es.ItemFieldCategoryNames, 10),
						),
				),
		).
		SourceIncludes(es.AllItemFields...).
		From(page * pageSize).
		Size(pageSize)

	switch sortType {
	case SortTypePriceAsc:
		esQuery.Sort(es.ItemFieldPrice, esquery.OrderAsc)
	case SortTypePriceDesc:
		esQuery.Sort(es.ItemFieldPrice, esquery.OrderDesc)
	case SortTypeReviewCount:
		esQuery.Sort(es.ItemFieldReviewCount, esquery.OrderDesc).Sort(es.ItemFieldAverageRating, esquery.OrderDesc)
	case SortTypeRating:
		esQuery.Sort(es.ItemFieldAverageRating, esquery.OrderDesc).Sort(es.ItemFieldReviewCount, esquery.OrderDesc)
	default:
		esQuery.Sort("_score", esquery.OrderDesc)
	}

	esQueryJSON, err := esQuery.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("esQuery.MarshalJSON(): %w", err)
	}

	return bytes.NewReader(esQueryJSON), nil
}

func (c *client) GetQuerySuggestions(ctx context.Context, query string) ([]string, error) {
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

func (c *client) search(ctx context.Context, indexName string, esQuery io.Reader) (*elastic.SearchResult, error) {
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

func (c *client) insertQuerySuggestion(ctx context.Context, query string) error {
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
		return fmt.Errorf("esClient.Index failed: %w", err)
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
