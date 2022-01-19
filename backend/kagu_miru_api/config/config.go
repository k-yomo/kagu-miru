package config

import "github.com/kelseyhightower/envconfig"

type Env string

const (
	EnvLocal Env = "local"
	EnvTest  Env = "test"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

func (e Env) IsDeployed() bool {
	return e != EnvLocal && e != EnvTest
}

type Config struct {
	Env            Env      `default:"local" envconfig:"APP_ENV"`
	Port           int      `default:"8000" envconfig:"PORT"`
	AllowedOrigins []string `default:"http://localhost:3000,http://localhost:3333" envconfig:"ALLOWED_ORIGINS"`

	GCPProjectID       string `default:"local" envconfig:"GCP_PROJECT_ID"`
	PubSubEventTopicID string `envconfig:"PUBSUB_EVENT_TOPIC_ID"`

	SpannerInstanceID string `envconfig:"SPANNER_INSTANCE_ID"`
	SpannerDatabaseID string `envconfig:"SPANNER_DATABASE_ID"`

	VertexAICategoryClassificationEndpointID string `envconfig:"VERTEX_AI_CATEGORY_CLASSIFICATION_ENDPOINT_ID"`

	SanityProjectID string `envconfig:"SANITY_PROJECT_ID"`
	SanityDatasetID string `default:"development" envconfig:"SANITY_DATASET_ID"`

	ElasticSearchUsername string `envconfig:"ELASTICSEARCH_USERNAME"`
	ElasticSearchPassword string `envconfig:"ELASTICSEARCH_PASSWORD"`
	ElasticSearchURL      string `default:"http://localhost:9200" envconfig:"ELASTICSEARCH_URL"`
	// ItemsIndexName                 string `default:"items" envconfig:"ITEMS_INDEX_NAME"`
	ItemsIndexName                 string `default:"items" envconfig:"ITEMS_INDEX_NAME"`
	ItemsQuerySuggestionsIndexName string `default:"items.query_suggestions" envconfig:"ITEMS_QUERY_SUGGESTIONS_INDEX_NAME"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
