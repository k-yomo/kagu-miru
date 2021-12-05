package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	GCPProjectID            string `envconfig:"GCP_PROJECT_ID"`
	PubsubItemUpdateTopicID string `default:"item-update" envconfig:"PUBSUB_ITEM_UPDATE_TOPIC_ID"`

	SpannerInstanceID string `envconfig:"SPANNER_INSTANCE_ID"`
	SpannerDatabaseID string `envconfig:"SPANNER_DATABASE_ID"`

	// To avoid late limit, we use multiple ids
	YahooShoppingApplicationIDs  []string `required:"true" envconfig:"YAHOO_SHOPPING_APPLICATION_IDS"`
	YahooShoppingStartCategoryID int      `default:"0" envconfig:"YAHOO_SHOPPING_START_CATEGORY_ID"`
	YahooShoppingMinPrice        int      `default:"0" envconfig:"YAHOO_SHOPPING_MIN_PRICE"`
}

func newConfig() (*config, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
