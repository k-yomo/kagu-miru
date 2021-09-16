package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	ElasticSearchUsername string `envconfig:"ELASTICSEARCH_USERNAME"`
	ElasticSearchPassword string `envconfig:"ELASTICSEARCH_PASSWORD"`
	ElasticSearchURL string `default:"http://localhost:9200" envconfig:"ELASTICSEARCH_URL"`
	ItemsIndexName   string `default:"items" envconfig:"ITEMS_INDEX_NAME"`

	RakutenApplicationID string `required:"true" envconfig:"RAKUTEN_APPLICATION_ID"`
	RakutenStartGenreID  int    `default:"0" envconfig:"RAKUTEN_START_GENRE_ID"`
	RakutenMinPrice      int    `default:"0" envconfig:"RAKUTEN_MIN_PRICE"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
