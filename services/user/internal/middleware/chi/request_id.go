// Package chimiddleware defines the middleware for the user service.
package chimiddleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	// HeaderRequestID is the header name for the request ID.
	HeaderRequestID = "X-Request-ID"

	// ContextRequestIDKey is the context key for storing the request ID.
	ContextRequestIDKey contextKey = "request_id"
)

// RequestID is a middleware that generates or propagates a request ID
// and stores it in the request context.
func RequestID() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Read existing request ID or generate new one
			reqID := r.Header.Get(HeaderRequestID)
			if reqID == "" {
				reqID = uuid.NewString()
			}

			// Store in context for retrieval by handlers/loggers
			ctx := context.WithValue(r.Context(), ContextRequestIDKey, reqID)
			r = r.WithContext(ctx)

			// Write header so client can trace it
			w.Header().Set(HeaderRequestID, reqID)

			next.ServeHTTP(w, r)
		})
	}
}
