package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/profiler"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/spanner"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/config"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/db"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlgen"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/queryclassifier"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/request"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/search"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/tracking"
	"github.com/k-yomo/kagu-miru/backend/pkg/csrf"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"github.com/k-yomo/kagu-miru/backend/pkg/spannerutil"
	"github.com/k-yomo/kagu-miru/backend/pkg/tracing"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	logger, err := logging.NewLogger(!cfg.Env.IsDeployed())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if cfg.Env.IsDeployed() {
		err, shutdown := tracing.InitTracer(cfg.GCPProjectID)
		if err != nil {
			logger.Error("set trace provider failed", zap.Error(err))
		} else {
			defer func() { _ = shutdown(ctx) }()
		}

		if err := profiler.Start(profiler.Config{}); err != nil {
			logger.Error("start profiler failed", zap.Error(err))
		}
	}

	spannerClient, err := spanner.NewClient(
		context.Background(),
		spannerutil.BuildSpannerDBPath(cfg.GCPProjectID, cfg.SpannerInstanceID, cfg.SpannerDatabaseID),
	)
	if err != nil {
		logger.Fatal("failed to initialize spanner client", zap.Error(err))
	}

	dbClient := db.NewSpannerDBClient(spannerClient)

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticSearchURL},
		Username:  cfg.ElasticSearchUsername,
		Password:  cfg.ElasticSearchPassword,

		Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("Elasticsearch%s", r.URL.EscapedPath())
		})),
	})
	if err != nil {
		logger.Fatal("failed to initialize elasticsearch client", zap.Error(err))
	}
	searchClient := search.NewElasticsearchClient(cfg.ItemsIndexName, cfg.ItemsQuerySuggestionsIndexName, esClient)

	predictionClient, err := aiplatform.NewPredictionClient(
		context.Background(),
		option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"),
		option.WithGRPCDialOption(grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor())),
	)
	if err != nil {
		logger.Fatal("failed to initialize elasticsearch client", zap.Error(err))
	}
	queryClassifierClient := queryclassifier.NewQueryClassifierClient(predictionClient, cfg.GCPProjectID, cfg.VertexAICategoryClassificationEndpointID)

	var eventLoader tracking.EventLoader
	if cfg.Env.IsDeployed() {
		pubsubClient, err := pubsub.NewClient(ctx, cfg.GCPProjectID)
		if err != nil {
			logger.Fatal("failed to initialize pubsub client", zap.Error(err))
		}
		eventLoader = tracking.NewEventLoader(pubsubClient.Topic(cfg.PubSubEventTopicID))
	} else {
		eventLoader = &tracking.NoopEventLoader{}
	}

	searchIDManager := tracking.NewSearchIDManager(cfg.Env.IsDeployed())

	gqlConfig := gqlgen.Config{
		Resolvers: graph.NewResolver(dbClient, searchClient, queryClassifierClient, searchIDManager, eventLoader),
	}
	gqlServer := handler.NewDefaultServer(gqlgen.NewExecutableSchema(gqlConfig))
	gqlServer.Use(tracing.GraphqlExtension{})
	gqlServer.Use(logging.GraphQLResponseInterceptor{})

	r := newBaseRouter(cfg, logger, searchIDManager)
	r.Route("/api", func(r chi.Router) {
		r.Handle("/graphql/playground", playground.Handler("GraphQL playground", "/api/graphql"))
		r.Handle("/graphql", gqlServer)
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		logger.Info(fmt.Sprintf("server listening on port: %d", cfg.Port))
		logger.Fatal(httpServer.ListenAndServe().Error())
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	logger.Info("Signal received, shutting down gracefully...", zap.Any("signal", <-sigCh))

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
}

func newBaseRouter(cfg *config.Config, logger *zap.Logger, searchIDManager *tracking.SearchIDManager) *chi.Mux {
	r := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedHeaders:   []string{"Origin", "Authorization", "Accept", "Content-Type", csrf.HeaderKey},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})

	r.Use(c.Handler)
	r.Use(
		middleware.RequestID,
		middleware.Recoverer,
		logging.NewMiddleware(cfg.GCPProjectID, logger),
		csrf.NewCSRFValidationMiddleware(cfg.Env.IsDeployed()),
		request.NewMiddleware(),
		searchIDManager.Middleware(),
	)
	return r
}
