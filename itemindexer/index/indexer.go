package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"io/ioutil"
)

type ItemIndexer struct {
	indexName string
	esClient  *elasticsearch.Client
}

func NewItemIndexer(indexName string, esClient *elasticsearch.Client) *ItemIndexer {
	return &ItemIndexer{
		indexName: indexName,
		esClient:  esClient,
	}
}

func (i *ItemIndexer) Index(ctx context.Context, item *Item) error {
	params := indexParams{Index: &index{Index: i.indexName, ID: item.ID}}
	postIndexByte, err := json.Marshal(&params)
	if err != nil {
		return fmt.Errorf("json.Marshal failed, param: %v,  err: %w", params, err)
	}

	response, err := i.esClient.Index(i.indexName, bytes.NewReader(postIndexByte), i.esClient.Index.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("esClient.Index failed, err: %w", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read body failed, err: %w", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("index failed, body: %v, err: %w", body, err)
	}

	return nil
}

func (i *ItemIndexer) BulkIndex(ctx context.Context, items []*Item) error {
	var bulkIndexParamsByte []byte
	for _, item := range items {
		params := indexParams{Index: &index{Index: i.indexName, ID: item.ID}}
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("json.Marshal failed, param: %v,  err: %w", params, err)
		}
		postByte, err := json.Marshal(item)
		if err != nil {
			return err
		}
		bulkIndexParamsByte = append(bulkIndexParamsByte, paramsJSON...)
		bulkIndexParamsByte = append(bulkIndexParamsByte, []byte("\n")...)
		bulkIndexParamsByte = append(bulkIndexParamsByte, postByte...)
		bulkIndexParamsByte = append(bulkIndexParamsByte, []byte("\n")...)
	}

	response, err := i.esClient.Bulk(bytes.NewReader(bulkIndexParamsByte), i.esClient.Bulk.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("esClient.Bulk failed, err: %w", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read body failed, err: %w", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("bulk index failed, body: %s, err: %w", body, err)
	}

	return nil
}
