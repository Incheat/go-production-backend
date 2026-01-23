// Package chimiddleware defines the middleware for the auth service.
package chimiddleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/incheat/go-production-backend/services/auth/internal/obs"
)

// PromMetrics adds Prometheus metrics to the request.
func PromMetrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			if routePattern == "" {
				routePattern = "unknown"
			}

			status := strconv.Itoa(ww.Status())
			method := r.Method
			duration := time.Since(start).Seconds()

			obs.HTTPRequestsTotal.WithLabelValues(method, routePattern, status).Inc()
			obs.HTTPRequestDuration.WithLabelValues(method, routePattern, status).Observe(duration)
		})
	}
}
