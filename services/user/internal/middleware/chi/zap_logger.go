// Package chimiddleware defines the middleware for the user service.
package chimiddleware

import (
	"net/http"
	"time"

	chimiddlewareutils "github.com/incheat/go-production-backend/services/user/internal/middleware/chi/utils"
	"go.uber.org/zap"
)

// ZapLogger logs the request and response using Uber Zap.
func ZapLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Create a logger with the request ID
			reqID, _ := r.Context().Value(ContextRequestIDKey).(string)
			if reqID == "" {
				reqID = "-"
			}

			// Create a new logger with the request ID
			requestLogger := logger.With(
				zap.String(string(ContextRequestIDKey), reqID),
			)
			ctx := chimiddlewareutils.WithLogger(r.Context(), requestLogger)
			r = r.WithContext(ctx)

			start := time.Now()

			// Wrap ResponseWriter to capture status code
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Process request
			next.ServeHTTP(ww, r)

			latency := time.Since(start)

			clientIP := r.RemoteAddr // similar to Gin's ClientIP()

			logger.Info("request handled",
				zap.String("request_id", reqID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.status),
				zap.Duration("latency", latency),
				zap.String("client_ip", clientIP),
			)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
