package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/k-yomo/kagu-miru/backend/internal/xerror"

	"github.com/k-yomo/kagu-miru/backend/pkg/uuid"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
	"github.com/k-yomo/kagu-miru/backend/pkg/imageutil"
	"golang.org/x/sync/errgroup"

	"github.com/k-yomo/kagu-miru/backend/internal/xitem"

	"cloud.google.com/go/spanner"

	"github.com/olivere/elastic/v7"
	"go.uber.org/multierr"
)

type ItemIndexer struct {
	spannerClient *spanner.Client
	esClient      *elastic.Client
	indexName     string
}

func NewItemIndexer(spannerClient *spanner.Client, esClient *elastic.Client, indexName string) *ItemIndexer {
	return &ItemIndexer{
		spannerClient: spannerClient,
		esClient:      esClient,
		indexName:     indexName,
	}
}

func (i *ItemIndexer) BulkIndex(ctx context.Context, items []*xitem.Item) error {
	if len(items) == 0 {
		return nil
	}

	itemIDGroupIDMap, err := i.getGroupIDItemIDMap(ctx, items)
	if err != nil {
		return err
	}

	var esItems []*es.Item
	for _, item := range items {
		esItem := mapItemFetcherItemToElasticsearchItem(item)
		groupID, ok := itemIDGroupIDMap[esItem.ID]
		if !ok {
			continue
		}
		esItem.GroupID = groupID
	}

	spannerItems := make([]*xspanner.Item, 0, len(items))
	for _, item := range items {
		groupID, ok := itemIDGroupIDMap[item.ID]
		if !ok {
			continue
		}
		spannerItem := mapItemToSpannerItem(item, groupID)
		spannerItems = append(spannerItems, spannerItem)
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		return i.insertOrUpdateItemsToSpanner(ctx, spannerItems)
	})
	eg.Go(func() error {
		return i.bulkIndexItemsToElasticsearch(ctx, esItems)
	})

	return eg.Wait()
}

func (i *ItemIndexer) insertOrUpdateItemsToSpanner(ctx context.Context, items []*xspanner.Item) error {
	var mutations []*spanner.Mutation
	for _, item := range items {
		m, err := spanner.InsertOrUpdateStruct(xspanner.ItemsTableName, item)
		if err != nil {
			// logging
			continue
		}
		mutations = append(mutations, m)
	}

	if _, err := i.spannerClient.Apply(ctx, mutations); err != nil {
		return err
	}
	return nil
}

func (i *ItemIndexer) bulkIndexItemsToElasticsearch(ctx context.Context, items []*es.Item) error {
	if len(items) == 0 {
		return nil
	}

	bulk := i.esClient.Bulk().Index(i.indexName)
	for _, item := range items {
		if item.IsActive() {
			bulk.Add(elastic.NewBulkIndexRequest().Index(i.indexName).Id(item.ID).Doc(item))
		} else {
			bulk.Add(elastic.NewBulkDeleteRequest().Index(i.indexName).Id(item.ID))
		}
	}

	resp, err := bulk.Do(ctx)
	if err != nil {
		return fmt.Errorf("esClient.Bulk failed: %w", err)
	}
	if resp.Errors {
		var errs []error
		for _, failed := range resp.Failed() {
			errs = append(errs, fmt.Errorf("id: %s, status: %d, err: %s", failed.Id, failed.Status, failed.Result))
		}
		return fmt.Errorf("read body failed: %w", multierr.Combine(errs...))
	}

	return nil
}

func (i *ItemIndexer) getGroupIDItemIDMap(ctx context.Context, items []*xitem.Item) (map[string]string, error) {
	eg := errgroup.Group{}
	itemIDGroupIDChan := make(chan map[string]string)

	for _, item := range items {
		item := item
		eg.Go(func() error {
			itemGroupID, err := i.findOrInitializeItemGroupID(ctx, item)
			if err != nil {
				return err
			}
			itemIDGroupIDChan <- map[string]string{item.ID: itemGroupID}
			return nil
		})
	}

	itemIDGroupIDMap := make(map[string]string)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			m, ok := <-itemIDGroupIDChan
			if !ok {
				return
			}
			for itemID, groupID := range m {
				itemIDGroupIDMap[itemID] = groupID
			}
		}
	}()

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(itemIDGroupIDChan)
	wg.Wait()

	return itemIDGroupIDMap, nil
}

func (i *ItemIndexer) findOrInitializeItemGroupID(ctx context.Context, item *xitem.Item) (string, error) {
	// 1. find group id if already indexed
	dbItem, err := xspanner.GetItem(ctx, i.spannerClient, item.ID)
	if err != nil && !xerror.IsErrorType(err, xerror.TypeNotFound) {
		return "", err
	}
	if dbItem != nil && dbItem.GroupID.Valid {
		return dbItem.GroupID.StringVal, nil
	}
	// 2. find similar items
	similarItems, err := i.getSimilarItems(ctx, item)
	if err != nil {
		return "", err
	}
	for _, similarItem := range similarItems {
		esItem := mapItemFetcherItemToElasticsearchItem(item)
		isSimilar, err := isSimilarItem(ctx, esItem, similarItem)
		if err != nil {
			// logging
			continue
		}
		if isSimilar {
			return similarItem.GroupID, nil
		}
	}

	// 3. initialize new group id
	return uuid.UUID(), nil
}

func isSimilarItem(ctx context.Context, a *es.Item, b *es.Item) (bool, error) {
	if a.JANCode == b.JANCode {
		return true, nil
	}
	if a.Name == b.Name {
		return true, nil
	}

	if len(a.ImageURLs) == 0 || len(b.ImageURLs) == 0 {
		return false, nil
	}

	return imageutil.IsSimilarImageByURLs(ctx, a.ImageURLs[0], b.ImageURLs[0])
}

func (i *ItemIndexer) getSimilarItems(ctx context.Context, item *xitem.Item) ([]*es.Item, error) {
	boolQuery := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery(es.ItemFieldName, item.Name),
	).Filter(
		elastic.NewTermsQuery(es.ItemFieldCategoryID, item.CategoryID),
	)

	resp, err := i.esClient.Search().
		Index(i.indexName).
		Query(boolQuery).
		SortBy(elastic.NewScoreSort()).
		From(0).
		Size(10).
		RequestCache(true).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("esClient.Search: %w", err)
	}

	return mapElasticsearchHitsToItems(resp.Hits.Hits), nil
}

func mapElasticsearchHitsToItems(hits []*elastic.SearchHit) []*es.Item {
	items := make([]*es.Item, 0, len(hits))
	for _, hit := range hits {
		var item es.Item
		if err := json.Unmarshal(hit.Source, &item); err != nil {
			continue
		}

		items = append(items, &item)
	}

	return items
}
