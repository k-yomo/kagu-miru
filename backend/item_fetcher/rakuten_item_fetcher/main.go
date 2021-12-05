package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/k-yomo/kagu-miru/backend/pkg/spannerutil"

	"cloud.google.com/go/spanner"

	"cloud.google.com/go/pubsub"
	"github.com/blendle/zapdriver"
	"github.com/k-yomo/kagu-miru/backend/pkg/rakutenichiba"
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

	rakutenIchibaClient := rakutenichiba.NewClient(cfg.RakutenApplicationIDs, cfg.RakutenAffiliateID)
	rakutenItemWorker := newWorker(pubsubItemUpdateTopic, spannerClient, rakutenIchibaClient, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneCh := make(chan struct{}, 1)
	go func() {
		logger.Info("rakutenItemWorker started running")
		if err := rakutenItemWorker.run(ctx, &rakutenWorkerOption{
			StartGenreID: cfg.RakutenStartGenreID,
		}); err != nil {
			logger.Error("rakutenItemWorker failed", zap.Error(err))
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
