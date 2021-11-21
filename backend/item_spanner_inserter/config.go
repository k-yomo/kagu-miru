package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	Port int `default:"8222" envconfig:"PORT"`

	GCPProjectID                   string `envconfig:"GCP_PROJECT_ID"`
	PubsubItemUpdateSubscriptionID string `default:"item-update.item-spanner-inserter" envconfig:"PUBSUB_ITEM_UPDATE_SUBSCRIPTION_ID"`

	SpannerInstanceID string `envconfig:"SPANNER_INSTANCE_ID"`
	SpannerDatabaseID string `envconfig:"SPANNER_DATABASE_ID"`
}

func newConfig() (*config, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
