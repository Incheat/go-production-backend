// Package interceptor defines the interceptors for the user service.
package interceptor

import (
	"github.com/incheat/go-production-backend/pkg/obs/logging"
	"github.com/incheat/go-production-backend/pkg/obs/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// DefaultChain returns a default chain of interceptors for the user service.
func DefaultChain(
	logger *zap.Logger,
) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		Recovery(),
		logging.GRPCRequestLogging(logger),
		metrics.PromGRPCMetrics(),
		// ZapTraceUnaryInterceptor(logger),
		// Logging(logger),
		// PromMetrics(),
	}
}
