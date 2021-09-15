package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RakutenApplicationID string `envconfig:"RAKUTEN_APPLICATION_ID" required:"true"`
	RakutenStartGenreID  int    `envconfig:"RAKUTEN_START_GENRE_ID" default:"0"`
	RakutenMinPrice      int    `envconfig:"RAKUTEN_MIN_PRICE" default:"0"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
