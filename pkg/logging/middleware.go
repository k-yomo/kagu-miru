package logging

import (
	"net/http"
	"strings"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

// NewMiddleware returns middleware to set logger to context
func NewMiddleware(gcpProjectID string, logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := newTraceFromTraceContext(gcpProjectID, r.Header.Get("X-Cloud-Trace-Context"))
			zapFields := append(
				zapdriver.TraceContext(t.TraceID, t.SpanID, true, t.ProjectID),
				zap.String("ip", r.RemoteAddr),
				zap.String("requestId", middleware.GetReqID(r.Context())),
			)
			ctx := ctxzap.ToContext(r.Context(), logger.With(zapFields...))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type traceInfo struct {
	ProjectID string
	TraceID   string
	SpanID    string
}

func newTraceFromTraceContext(projectID, traceContext string) traceInfo {
	t := traceInfo{ProjectID: projectID}
	if traceContext != "" {
		params := strings.Split(traceContext, "/")
		if len(params) >= 2 {
			t.TraceID = params[0]
			t.SpanID = params[1]
		}
	}
	return t
}
