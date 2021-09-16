package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/blendle/zapdriver"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/k-yomo/kagu-miru/itemindexer/config"
	"github.com/k-yomo/kagu-miru/itemindexer/index"
	"github.com/k-yomo/kagu-miru/itemindexer/indexworker"
	"github.com/k-yomo/kagu-miru/pkg/rakuten"
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
	})
	if err != nil {
		logger.Fatal("failed to initialize elasticsearch client", zap.Error(err))
	}

	indexer := index.NewItemIndexer(cfg.ItemsIndexName, esClient)
	rakutenIchibaClient := rakuten.NewIchibaClient(cfg.RakutenApplicationID)
	rakutenItemWorker := indexworker.NewRakutenItemIndexWorker(indexer, rakutenIchibaClient, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneCh := make(chan struct{}, 1)
	go func() {
		logger.Info("rakutenItemWorker started running")
		if err := rakutenItemWorker.Run(ctx, &indexworker.RakutenWorkerOption{
			StartGenreID: cfg.RakutenStartGenreID,
			MinPrice:     cfg.RakutenMinPrice,
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
