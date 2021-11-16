package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.uber.org/zap"
)

type GraphQLResponseInterceptor struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
} = GraphQLResponseInterceptor{}

func (g GraphQLResponseInterceptor) ExtensionName() string {
	return "Logging"
}

func (g GraphQLResponseInterceptor) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (g GraphQLResponseInterceptor) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	requestStarted := time.Now()
	resp := next(ctx)

	status := statusCode(len(resp.Errors) > 0)
	oc := graphql.GetOperationContext(ctx)
	Logger(ctx).Info(
		fmt.Sprintf("%s - %s", oc.OperationName, status),
		zap.String("type", "request"),
		zap.String("operation", oc.OperationName),
		zap.String("status", status),
		zap.Int64("latencyMs", time.Since(requestStarted).Milliseconds()),
		zap.String("query", oc.RawQuery),
		zap.Any("variables", oc.Variables),
	)
	return resp
}

func statusCode(isErr bool) string {
	if isErr {
		return "error"
	}
	return "ok"
}
