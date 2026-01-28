// Package logging defines the HTTP middleware for the observability.
package logging

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/incheat/go-production-backend/pkg/obs/correlation"
)

const (
	// HeaderRequestID is the header name for the request ID. Envoy generated/inherited
	HeaderRequestID = "X-Request-ID"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// HTTPRequestLogging injects request-scoped logger into context.
// It uses chi middleware.RequestID() (or middleware.GetReqID) for request_id.
// Trace/span are pulled from OTel context via correlation.TraceFields.
func HTTPRequestLogging(base *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			reqID := r.Header.Get(HeaderRequestID) // Envoy generated/inherited
			if reqID == "" {
				reqID = "-"
			}

			ctx, err := correlation.SetBaggage(r.Context(), correlation.BaggageRequestID, reqID)
			if err != nil {
				base.Warn("failed to set request_id baggage", zap.Error(err))
			}

			// request-scoped logger (request_id + trace/span + optional baggage)
			l := base.With(
				zap.String("request_id", reqID),
			).With(correlation.TraceFields(ctx)...)

			l = l.With(correlation.BaggageFields(ctx,
				correlation.BaggageTenantID,
				correlation.BaggageRequestID,
			)...)

			// Put logger into context so handlers can use it
			ctxWithLogger := correlation.ContextWithLogger(ctx, l)
			r = r.WithContext(ctxWithLogger)

			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(ww, r)
		})
	}
}
