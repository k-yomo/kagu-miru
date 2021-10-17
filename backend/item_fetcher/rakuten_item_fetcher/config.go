package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	GCPProjectID            string `envconfig:"GCP_PROJECT_ID"`
	PubsubItemUpdateTopicID string `default:"item-update" envconfig:"PUBSUB_ITEM_UPDATE_TOPIC_ID"`

	// To avoid late limit, we use multiple ids
	RakutenApplicationIDs []string `required:"true" envconfig:"RAKUTEN_APPLICATION_IDS"`
	RakutenAffiliateID    string   `required:"true" envconfig:"RAKUTEN_AFFILIATE_ID"`
	RakutenStartGenreID   int      `default:"0" envconfig:"RAKUTEN_START_GENRE_ID"`
	RakutenMinPrice       int      `default:"0" envconfig:"RAKUTEN_MIN_PRICE"`
}

func newConfig() (*config, error) {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
