// Package interceptor defines the interceptors for the user service.
package interceptor

import (
	"context"
	"time"

	"github.com/incheat/go-production-backend/services/user/internal/obs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// PromMetrics adds Prometheus metrics to the gRPC request.
func PromMetrics() grpc.UnaryServerInterceptor {
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
		obs.GrpcRequestsTotal.WithLabelValues(method, code).Inc()
		obs.GrpcRequestDuration.WithLabelValues(method, code).Observe(time.Since(start).Seconds())
		return resp, err
	}
}
