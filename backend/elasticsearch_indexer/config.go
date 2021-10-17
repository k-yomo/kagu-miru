package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	Port int `default:"8111" envconfig:"PORT"`

	GCPProjectID                   string `envconfig:"GCP_PROJECT_ID"`
	PubsubItemUpdateSubscriptionID string `default:"item-update.elasticsearch-indexer" envconfig:"PUBSUB_ITEM_UPDATE_SUBSCRIPTION_ID"`

	ElasticSearchUsername string `envconfig:"ELASTICSEARCH_USERNAME"`
	ElasticSearchPassword string `envconfig:"ELASTICSEARCH_PASSWORD"`
	ElasticSearchURL      string `default:"http://localhost:9200" envconfig:"ELASTICSEARCH_URL"`
	ItemsIndexName        string `default:"items" envconfig:"ITEMS_INDEX_NAME"`
}

func newConfig() (*config, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
