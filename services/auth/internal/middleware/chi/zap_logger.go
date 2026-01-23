// Package chimiddleware defines the middleware for the auth service.
package chimiddleware

import (
	"net/http"
	"time"

	chimiddlewareutils "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi/utils"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	// ContextTraceIDKey is the context key for storing the trace ID.
	ContextTraceIDKey = "trace_id"
	// ContextSpanIDKey is the context key for storing the span ID.
	ContextSpanIDKey = "span_id"
	// ContextRequestIDKey is the context key for storing the request ID.
	ContextRequestIDKey = "request_id"
)

// ZapLogger logs the request and response using Uber Zap.
func ZapLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// ----------------------------
			// Request ID
			// ----------------------------
			var reqID string
			requestMeta, ok := chimiddlewareutils.GetRequestMeta(r.Context())
			if !ok || requestMeta.RequestID == "" {
				reqID = "-"
			} else {
				reqID = requestMeta.RequestID
			}

			// ----------------------------
			// ⭐ Trace / Span
			// ----------------------------
			var traceID, spanID string
			sc := trace.SpanFromContext(r.Context()).SpanContext()
			if sc.IsValid() {
				traceID = sc.TraceID().String()
				spanID = sc.SpanID().String()
			} else {
				traceID = "-"
				spanID = "-"
			}

			// ----------------------------
			// Request-scoped logger
			// ----------------------------
			requestLogger := logger.With(
				zap.String(string(ContextRequestIDKey), reqID),
				zap.String(string(ContextTraceIDKey), traceID), // ⭐
				zap.String(string(ContextSpanIDKey), spanID),   // ⭐
			)

			ctx := chimiddlewareutils.WithLogger(r.Context(), requestLogger)
			r = r.WithContext(ctx)

			start := time.Now()

			// Wrap ResponseWriter to capture status code
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Process request
			next.ServeHTTP(ww, r)

			latency := time.Since(start)

			clientIP := "-"
			if ok {
				clientIP = requestMeta.IPAddress
			}

			// ----------------------------
			// Access log (also with trace/span)
			// ----------------------------
			requestLogger.Info("request handled",
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
