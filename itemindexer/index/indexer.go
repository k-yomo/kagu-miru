package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/k-yomo/kagu-miru/internal/es"
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

type indexParams struct {
	Index *indexParamsIndex `json:"index"`
}

type indexParamsIndex struct {
	Index string `json:"_index"`
	ID    string `json:"_id"`
}

func (i *ItemIndexer) BulkIndex(ctx context.Context, items []*es.Item) error {
	if len(items) == 0 {
		return nil
	}
	var bulkIndexParamsByte []byte
	for _, item := range items {
		params := &indexParams{Index: &indexParamsIndex{Index: i.indexName, ID: item.ID}}
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("json.Marshal failed, param: %v,  err: %w", params, err)
		}
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return err
		}
		bulkIndexParamsByte = append(bulkIndexParamsByte, paramsJSON...)
		bulkIndexParamsByte = append(bulkIndexParamsByte, []byte("\n")...)
		bulkIndexParamsByte = append(bulkIndexParamsByte, itemJSON...)
		bulkIndexParamsByte = append(bulkIndexParamsByte, []byte("\n")...)
	}

	response, err := i.esClient.Bulk(bytes.NewReader(bulkIndexParamsByte), i.esClient.Bulk.WithContext(ctx))
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
