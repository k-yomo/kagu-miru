package logging

import (
	"context"

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
	oc := graphql.GetOperationContext(ctx)
	Logger(ctx).Info(oc.OperationName, zap.String("query", oc.RawQuery), zap.Any("variables", oc.Variables))
	return next(ctx)
}
