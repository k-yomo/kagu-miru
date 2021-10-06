package event_bq_inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
)

type EventForBQInsert struct {
	PubSubMessageID string
	Value           map[string]interface{}
	bigquery.ValueSaver
}

func (e *EventForBQInsert) Save() (row map[string]bigquery.Value, insertID string, err error) {
	bqValueMap := make(map[string]bigquery.Value)
	for k, v := range e.Value {
		bqValueMap[k] = v
	}
	return bqValueMap, e.PubSubMessageID, nil
}

func InsertEventToBigquery(ctx context.Context, m *pubsub.Message) error {
	datasetID := mustEnv("BQ_DATASET_ID")
	tableID := mustEnv("BQ_TABLE_ID")

	bqClient, err := bigquery.NewClient(ctx, bigquery.DetectProjectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}

	var event EventForBQInsert
	if err := json.Unmarshal(m.Data, &event.Value); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	inserter := bqClient.Dataset(datasetID).Table(tableID).Inserter()
	if err := inserter.Put(ctx, &event); err != nil {
		return fmt.Errorf("inserter.Put: %w", err)
	}

	return nil
}

func mustEnv(key string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		panic(fmt.Sprintf("env variable '%s' is not set", key))
	}
	return envVar
}
