package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port           int      `default:"8000" envconfig:"PORT"`
	GCPProjectID   string   `default:"local" envconfig:"GCP_PROJECT_ID"`
	AllowedOrigins []string `default:"http://localhost:3000" envconfig:"ALLOWED_ORIGINS"`

	ItemsIndexName string `default:"items" envconfig:"ITEMS_INDEX_NAME"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
