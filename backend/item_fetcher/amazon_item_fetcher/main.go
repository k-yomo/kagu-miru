package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/spanner"
	"github.com/blendle/zapdriver"
	"github.com/k-yomo/kagu-miru/backend/pkg/amazon"
	"github.com/k-yomo/kagu-miru/backend/pkg/spannerutil"
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
	pubsubItemUpdateTopic.EnableMessageOrdering = true

	spannerClient, err := spanner.NewClient(
		context.Background(),
		spannerutil.BuildSpannerDBPath(cfg.GCPProjectID, cfg.SpannerInstanceID, cfg.SpannerDatabaseID),
	)
	if err != nil {
		logger.Fatal("failed to initialize spanner client", zap.Error(err))
	}

	amazonIchibaClient, err := amazon.NewClient(cfg.AmazonPartnerTag, cfg.AmazonAccessKey, cfg.AmazonSecretKey)
	if err != nil {
		logger.Fatal("failed to initialize amazon api client", zap.Error(err))
	}
	amazonItemWorker := newWorker(pubsubItemUpdateTopic, spannerClient, amazonIchibaClient, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneCh := make(chan struct{}, 1)
	go func() {
		logger.Info("amazonItemWorker started running")
		if err := amazonItemWorker.run(ctx, &amazonWorkerOption{
			StartBrowseNodeID: cfg.AmazonStartGenreID,
		}); err != nil {
			logger.Error("amazonItemWorker failed", zap.Error(err))
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
