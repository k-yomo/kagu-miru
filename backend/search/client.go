package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strings"
)

type Client struct {
	itemsIndexName string
	esClient       *elasticsearch.Client
}

func NewSearchClient(itemsIndexName string, esClient *elasticsearch.Client) *Client {
	return &Client{
		itemsIndexName: itemsIndexName,
		esClient:       esClient,
	}
}


const DefaultPage uint64 = 1
const defaultPageSize uint64 = 100

type Request struct {
	Query    string
	Page     uint64
	PageSize *int
}

type Response struct {
	Result    *Result
	Page      uint64
	TotalPage uint64
}

func (p *Client) SearchItems(ctx context.Context, req *Request) (*Response, error) {
	var page uint64
	if req.Page > 1 {
		page = req.Page - 1 // page starts from 0 in elasticsearch
	}
	pageSize := defaultPageSize
	if req.PageSize != nil {
		pageSize = uint64(*req.PageSize)
	}
	response, err := p.esClient.Search(
		p.esClient.Search.WithContext(ctx),
		p.esClient.Search.WithIndex(p.itemsIndexName),
		p.esClient.Search.WithBody(buildSearchQuery(req.Query, page, pageSize)),
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
	if response.StatusCode >= 400 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "read error response failed, response: %s", response.String())
		}
		return nil, errors.Errorf("search request to elasticsearch failed with status %s, body: %s", response.Status(), body)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body failed")
	}
	sr := &Result{}
	if err := json.Unmarshal(body, sr); err != nil {
		return nil, errors.Wrapf(err, "unmarshal response body json failed, body: %s", body)
	}
	return sr, nil
}

func buildSearchQuery(query string, page, pageSize uint64) io.Reader {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`{
	"_source": [
		"id", 
		"name", 
		"description", 
		"status", 
		"sellingPageUrl", 
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
	"from": %d,
	"size": %d,
	"sort": [{ "_score" : "desc" }, { "_doc" : "asc" }]
}`, query, page*pageSize, pageSize))

	return strings.NewReader(b.String())
}

func calcTotalPage(totalItems, pageSize uint64) uint64 {
	total := totalItems + pageSize - 1
	if total == 0 {
		return 0
	}
	return total / pageSize
}
