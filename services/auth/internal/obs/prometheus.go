// Package obs provides Prometheus metrics for the auth service.
package obs

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTPRequestsTotal is the total number of HTTP requests.
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	// HTTPRequestDuration is the duration of the HTTP request.
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)
)

// Init initializes the Prometheus metrics.
func Init() {
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration)
}
