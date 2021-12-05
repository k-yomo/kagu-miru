package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/pkg/spannerutil"

	"cloud.google.com/go/pubsub"
	"github.com/blendle/zapdriver"
	"github.com/k-yomo/kagu-miru/backend/pkg/yahoo_shopping"
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

	pubsubClient, err := pubsub.NewClient(context.Background(), cfg.GCPProjectID)
	if err != nil {
		logger.Fatal("failed to initialize pubsub client", zap.Error(err))
	}
	pubsubItemUpdateTopic := pubsubClient.Topic(cfg.PubsubItemUpdateTopicID)

	spannerClient, err := spanner.NewClient(
		context.Background(),
		spannerutil.BuildSpannerDBPath(cfg.GCPProjectID, cfg.SpannerInstanceID, cfg.SpannerDatabaseID),
	)
	if err != nil {
		logger.Fatal("failed to initialize spanner client", zap.Error(err))
	}

	yahooShoppingClient := yahoo_shopping.NewClient(cfg.YahooShoppingApplicationIDs)
	yahooShoppingItemWorker := newWorker(pubsubItemUpdateTopic, spannerClient, yahooShoppingClient, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneCh := make(chan struct{}, 1)
	go func() {
		logger.Info("yahooShoppingItemWorker started running")
		if err := yahooShoppingItemWorker.run(ctx, &yahooShoppingWorkerOption{
			StartCategoryID: cfg.YahooShoppingStartCategoryID,
			MinPrice:        cfg.YahooShoppingMinPrice,
		}); err != nil {
			logger.Error("yahooShoppingItemWorker failed", zap.Error(err))
		}
		doneCh <- struct{}{}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-doneCh:
	case sig := <-sigCh:
		logger.Info("Signal received, shutting down gracefully...", zap.Any("signal", sig))
	}

	cancel()
	logger.Info("stop indexer")
}
