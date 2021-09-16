package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/blendle/zapdriver"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k-yomo/kagu-miru/backend/config"
	"github.com/k-yomo/kagu-miru/backend/graph"
	"github.com/k-yomo/kagu-miru/backend/graph/gqlgen"
	"github.com/k-yomo/kagu-miru/backend/search"
	"github.com/k-yomo/kagu-miru/pkg/csrf"
	"github.com/k-yomo/kagu-miru/pkg/logging"
	"github.com/k-yomo/kagu-miru/pkg/tracing"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	logger, err := zapdriver.NewProduction()
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("failed to initialize config", zap.Error(err))
	}

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticSearchURL},
		Username:  cfg.ElasticSearchUsername,
		Password:  cfg.ElasticSearchPassword,

		Transport: otelhttp.NewTransport(http.DefaultTransport),
	})
	if err != nil {
		logger.Fatal("failed to initialize elasticsearch client", zap.Error(err))
	}

	gqlConfig := gqlgen.Config{
		Resolvers: graph.NewResolver(search.NewSearchClient(cfg.ItemsIndexName, esClient)),
	}
	gqlServer := handler.NewDefaultServer(gqlgen.NewExecutableSchema(gqlConfig))
	gqlServer.Use(tracing.GraphqlExtension{})
	gqlServer.Use(logging.GraphQLResponseInterceptor{})

	r := newBaseRouter(cfg, logger)
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", gqlServer)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: r}

	go func() {
		logger.Info(fmt.Sprintf("server listening on port: %d", cfg.Port))
		logger.Fatal(httpServer.ListenAndServe().Error())
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	logger.Info("Signal received, shutting down gracefully...", zap.Any("signal", <-sigCh))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
}

func newBaseRouter(cfg *config.Config, logger *zap.Logger) *chi.Mux {
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
		middleware.RealIP,
		middleware.Recoverer,
		logging.NewMiddleware(cfg.GCPProjectID, logger),
		csrf.NewCSRFValidationMiddleware(cfg.Env.IsDeployed()),
	)
	return r
}
