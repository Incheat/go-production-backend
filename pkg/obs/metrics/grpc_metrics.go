// Package metrics defines the gRPC metrics for the observability.
package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// GrpcRequestsTotal is the total number of gRPC requests.
	GrpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "code"},
	)

	// GrpcRequestDuration is the duration of the gRPC request.
	GrpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)
)

// RegisterGRPC registers the gRPC metrics into the provided registry.
func RegisterGRPC(reg prometheus.Registerer) {
	reg.MustRegister(GrpcRequestsTotal, GrpcRequestDuration)
}
