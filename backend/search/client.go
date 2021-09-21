package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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
)

type Request struct {
	Query    string
	SortType SortType
	Page     uint64
	PageSize *int
}

type Response struct {
	Result    *Result
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
	response, err := c.esClient.Search(
		c.esClient.Search.WithContext(ctx),
		c.esClient.Search.WithIndex(c.itemsIndexName),
		c.esClient.Search.WithBody(buildSearchQuery(req.Query, req.SortType, page, pageSize)),
	)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}
	defer response.Body.Close()

	searchResult, err := mapResponseToSearchResult(response)
	if err != nil {
		return nil, err
	}

	return &Response{
		Result:    searchResult,
		Page:      req.Page,
		TotalPage: calcTotalPage(uint64(searchResult.Hits.Total.Value), 100),
	}, nil
}

func mapResponseToSearchResult(response *esapi.Response) (*Result, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read response body failed, response: %s", response.String())
	}
	if response.StatusCode >= 400 {
		return nil, errors.Errorf("search request to elasticsearch failed with status %s, body: %s", response.Status(), body)
	}

	sr := &Result{}
	if err := json.Unmarshal(body, sr); err != nil {
		return nil, errors.Wrapf(err, "unmarshal response body json failed, body: %s", body)
	}
	return sr, nil
}

func buildSearchQuery(query string, sortType SortType, page, pageSize uint64) io.Reader {
	var sort string
	switch sortType {
	case SortTypePriceAsc:
		sort = `[{ "price" : "asc" }]`
	case SortTypePriceDesc:
		sort = `[{ "price" : "desc" }]`
	default:
		sort = `[{ "_score" : "desc" }]`
	}

	// TODO: refactor query construction
	esQuery := fmt.Sprintf(`{
	"_source": [
		"id", 
		"name", 
		"description", 
		"status", 
		"url", 
		"affiliateUrl", 
		"price",
		"imageUrls",
		"averageRating",
		"reviewCount",
		"platform"
	],
	"query": {
		"multi_match" : {
			"query": %q,
			"type": "most_fields",
			"fields": ["name^100", "description"]
		}
	},
	"sort": %s,
	"from": %d,
	"size": %d
}`, query, sort, page*pageSize, pageSize)

	return strings.NewReader(esQuery)
}

func (c *client) GetQuerySuggestions(ctx context.Context, query string) ([]string, error) {
	esQuery := fmt.Sprintf(`{
  "size": 0,
  "query": {
    "bool": {
      "should": [{
        "match": {
          "query.autocomplete": {
            "query": %q
          }
        }
      }, {
        "match": {
          "query.readingform": {
            "query":  %q,
            "fuzziness": "AUTO",
            "operator": "and"
          }
        }
      }]
    }
  },
  "aggs": {
    "queries": {
      "terms": {
        "field": "query",
        "order": {
          "_count": "desc"
        },
        "size": "10"
      }
    }
  }
}`, query, query)
	esResponse, err := c.esClient.Search(
		c.esClient.Search.WithContext(ctx),
		c.esClient.Search.WithIndex(c.itemsQuerySuggestionsIndexName),
		c.esClient.Search.WithBody(strings.NewReader(esQuery)),
	)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}
	defer esResponse.Body.Close()

	body, err := ioutil.ReadAll(esResponse.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read response body failed, response: %s", esResponse.String())
	}
	if esResponse.StatusCode >= 400 {
		return nil, fmt.Errorf("search request to elasticsearch failed with status %s, body: %s", esResponse.Status(), body)
	}

	var response struct {
		Aggregations struct {
			Queries struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				} `json:"buckets"`
			} `json:"queries"`
		} `json:"aggregations"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("search request to elasticsearch failed with body: %s: %w", body, err)
	}

	suggestedQueries := make([]string, 0, len(response.Aggregations.Queries.Buckets))
	for _, bucket := range response.Aggregations.Queries.Buckets {
		if bucket.Key == query {
			continue
		}
		suggestedQueries = append(suggestedQueries, bucket.Key)
	}

	return suggestedQueries, nil
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
