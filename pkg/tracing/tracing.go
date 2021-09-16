package tracing

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(gcpProjectID string) (error, func(ctx context.Context) error) {
	exporter, err := texporter.New(texporter.WithProjectID(gcpProjectID))
	if err != nil {
		return err, nil
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)
	return nil, exporter.Shutdown
}
