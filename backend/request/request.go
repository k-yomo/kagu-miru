package request

import (
	"context"
	"net/http"
	"strings"
)

type ctxRequestKey struct {
}

// NewMiddleware creates middleware to set request to given context
// We are using this since we can't access http.Request directly from graphql resolver
func NewMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ctxRequestKey{}, r)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestFromCtx gets request from given context
func GetRequestFromCtx(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(ctxRequestKey{}).(*http.Request)
	return req, ok
}

func RealClientIP(req *http.Request) string {
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 2 {
			// Since we have a GCLB in front, and it sets client-ip and lb-ip the second last IP is the client IP.
			// X-Forwarded-For: <supplied-value>,<client-ip>,<load-balancer-ip>
			return ips[len(ips)-2]
		}
		return ips[len(ips)-1]
	}
	return req.RemoteAddr
}
