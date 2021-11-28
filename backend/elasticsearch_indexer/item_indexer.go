package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/k-yomo/kagu-miru/backend/internal/es"
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

type indexingParams struct {
	Index *documentMeta `json:"index"`
}

type deleteParams struct {
	Delete *documentMeta `json:"delete"`
}

type documentMeta struct {
	Index string `json:"_index"`
	ID    string `json:"_id"`
}

func (i *ItemIndexer) BulkIndex(ctx context.Context, items []*es.Item) error {
	if len(items) == 0 {
		return nil
	}
	var bulkParamsBytes []byte
	for _, item := range items {
		if item.IsActive() {
			params := &indexingParams{Index: &documentMeta{Index: i.indexName, ID: item.ID}}
			paramsJSON, err := json.Marshal(params)
			if err != nil {
				return fmt.Errorf("json.Marshal failed, params: %v,  err: %w", params, err)
			}
			itemJSON, err := json.Marshal(item)
			if err != nil {
				return err
			}
			bulkParamsBytes = append(bulkParamsBytes, paramsJSON...)
			bulkParamsBytes = append(bulkParamsBytes, []byte("\n")...)
			bulkParamsBytes = append(bulkParamsBytes, itemJSON...)
			bulkParamsBytes = append(bulkParamsBytes, []byte("\n")...)
		} else {
			params := &deleteParams{Delete: &documentMeta{Index: i.indexName, ID: item.ID}}
			paramsJSON, err := json.Marshal(params)
			if err != nil {
				return fmt.Errorf("json.Marshal failed, params: %v,  err: %w", params, err)
			}
			bulkParamsBytes = append(bulkParamsBytes, paramsJSON...)
			bulkParamsBytes = append(bulkParamsBytes, []byte("\n")...)
		}
	}

	response, err := i.esClient.Bulk(bytes.NewReader(bulkParamsBytes), i.esClient.Bulk.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("esClient.Bulk failed: %w", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("bulk index failed, body: %s: %w", body, err)
	}

	return nil
}
