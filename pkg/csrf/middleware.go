package csrf

import "net/http"

const HeaderKey = "X-Requested-By"

// NewCSRFValidationMiddleware initializes middleware that checks custom header to prevent csrf attack
func NewCSRFValidationMiddleware(enable bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enable {
				if r.Header.Get(HeaderKey) == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					_, _ = w.Write([]byte(`{"error": "Can't verify CSRF header"}`))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
