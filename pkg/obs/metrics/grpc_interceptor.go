// Package metrics defines the gRPC interceptor for the observability.
package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// PromGRPCMetrics adds Prometheus metrics to the gRPC request.
func PromGRPCMetrics() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err).String()
		method := info.FullMethod
		GrpcRequestsTotal.WithLabelValues(method, code).Inc()
		GrpcRequestDuration.WithLabelValues(method, code).Observe(time.Since(start).Seconds())
		return resp, err
	}
}
