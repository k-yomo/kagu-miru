package tracking

import (
	"context"
	"encoding/json"

	"github.com/cenkalti/backoff/v4"

	"go.uber.org/zap"

	"github.com/k-yomo/kagu-miru/pkg/logging"

	"cloud.google.com/go/pubsub"
)

type EventLoader interface {
	Load(ctx context.Context, event *Event)
}

type eventLoader struct {
	pubsubEventTopic *pubsub.Topic
}

func NewEventLoader(pubsubEventTopic *pubsub.Topic) *eventLoader {
	return &eventLoader{
		pubsubEventTopic: pubsubEventTopic,
	}
}

func (l *eventLoader) Load(ctx context.Context, event *Event) {
	logger := logging.Logger(ctx)
	eventData, err := json.Marshal(event)
	if err != nil {
		logger.Error("json marshal event data failed", zap.Any("event", event))
		return
	}

	go func() {
		msg := pubsub.Message{
			Data: eventData,
			Attributes: map[string]string{
				"eventId": event.ID.String(),
				"action":  event.Action.String(),
			},
		}

		b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3)
		err := backoff.Retry(func() error {
			_, err := l.pubsubEventTopic.Publish(context.Background(), &msg).Get(context.Background())
			return err
		}, b)
		if err != nil {
			logger.Warn("publish event failed", zap.Any("event", event))
		}
	}()
}

// NoopEventLoader ca be used for local env / test
type NoopEventLoader struct {
}

func (l *NoopEventLoader) Load(ctx context.Context, event *Event) {
	logging.Logger(ctx).Info("tracking event", zap.Any("event", event))
}
