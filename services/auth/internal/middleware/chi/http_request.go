// Package chimiddleware defines the middleware for the auth service.
package chimiddleware

import (
	"net/http"

	chimiddlewareutils "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi/utils"
)

// HTTPRequest is a middleware that stores the http request in the context.
func HTTPRequest() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Store in context for retrieval by handlers/loggers
			ctx := chimiddlewareutils.WithHTTPRequest(r.Context(), r)
			r = r.WithContext(ctx)

			// Process request
			next.ServeHTTP(w, r)
		})
	}
}
