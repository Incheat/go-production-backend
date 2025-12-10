// Package chimiddleware defines the middleware for the auth service.
package chimiddleware

import (
	"net/http"

	"go.uber.org/zap"
)

// ZapRecovery is a middleware that recovers from panics and logs them using Zap.
func ZapRecovery(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Panic recovery wrapper
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						zap.Any("error", rec),
						zap.String("path", r.URL.Path),
					)

					// Always return 500 on panic
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
