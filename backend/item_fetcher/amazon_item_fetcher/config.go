package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	GCPProjectID            string `envconfig:"GCP_PROJECT_ID"`
	PubsubItemUpdateTopicID string `default:"item-update" envconfig:"PUBSUB_ITEM_UPDATE_TOPIC_ID"`

	SpannerInstanceID string `envconfig:"SPANNER_INSTANCE_ID"`
	SpannerDatabaseID string `envconfig:"SPANNER_DATABASE_ID"`

	AmazonPartnerTag string `required:"true" envconfig:"AMAZON_PARTNER_TAG"`
	AmazonAccessKey  string `required:"true" envconfig:"AMAZON_ACCESS_KEY"`
	AmazonSecretKey  string `required:"true" envconfig:"AMAZON_SECRET_KEY"`

	AmazonStartGenreID string `envconfig:"AMAZON_START_GENRE_ID"`
}

func newConfig() (*config, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
