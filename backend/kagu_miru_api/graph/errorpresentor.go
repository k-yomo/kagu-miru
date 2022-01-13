package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/k-yomo/kagu-miru/backend/internal/xerror"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func NewErrorPresenter() graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		code := mapFromXErrorType(xerror.ErrorType(gqlErr.Unwrap()))
		if code == gqlmodel.ErrorCodeInternal {
			gqlErr.Message = "internal server error"
		}
		gqlErr.Extensions = map[string]interface{}{
			"code": code,
		}
		return gqlErr
	}
}
