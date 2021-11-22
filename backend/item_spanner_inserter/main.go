package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/spanner"
	"github.com/blendle/zapdriver"
	"github.com/k-yomo/pm"
	"github.com/k-yomo/pm/middleware/logging/pm_zap"
	"github.com/k-yomo/pm/middleware/pm_autoack"
	"github.com/k-yomo/pm/middleware/pm_recovery"
	"go.uber.org/zap"
)

func main() {
	logger, err := zapdriver.NewProduction()
	if err != nil {
		panic(err)
	}

	cfg, err := newConfig()
	if err != nil {
		logger.Fatal("failed to initialize config", zap.Error(err))
	}

	spannerClient, err := spanner.NewClient(
		context.Background(),
		fmt.Sprintf("projects/%s/instances/%s/databases/%s", cfg.GCPProjectID, cfg.SpannerInstanceID, cfg.SpannerDatabaseID),
	)
	if err != nil {
		logger.Fatal("failed to initialize spanner client", zap.Error(err))
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), cfg.GCPProjectID)
	if err != nil {
		logger.Fatal("failed to initialize pubsub client", zap.Error(err))
	}
	pubsubSubscriber := pm.NewSubscriber(
		pubsubClient,
		pm.WithSubscriptionInterceptor(
			pm_zap.SubscriptionInterceptor(logger),
			pm_autoack.SubscriptionInterceptor(),
			pm_recovery.SubscriptionInterceptor(),
		),
	)

	defer pubsubSubscriber.Close()

	err = pubsubSubscriber.HandleSubscriptionFunc(
		pubsubClient.Subscription(cfg.PubsubItemUpdateSubscriptionID),
		pm.NewBatchMessageHandler(newItemUpdateHandler(spannerClient, logger), pm.BatchMessageHandlerConfig{
			DelayThreshold: 1 * time.Second,
			CountThreshold: 1000,
			ByteThreshold:  1e7, // 10MB
			NumGoroutines:  10,
		}),
	)
	if err != nil {
		logger.Fatal("failed to register item update subscription", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	pubsubSubscriber.Run(ctx)
	logger.Info("pubsub subscriber started running")

	go func() {
		// this is a dummy http server just to meet Cloud Run requirements
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
			logger.Error("failed to register item update subscription", zap.Error(err))
			cancel()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info("Signal received, shutting down gracefully...", zap.Error(ctx.Err()))
	case sig := <-sigCh:
		logger.Info("Signal received, shutting down gracefully...", zap.Any("signal", sig))
	}
}