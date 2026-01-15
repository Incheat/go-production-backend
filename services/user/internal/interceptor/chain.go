// Package interceptor defines the interceptors for the user service.
package interceptor

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// DefaultChain returns a default chain of interceptors for the user service.
func DefaultChain(
	logger *zap.Logger,
) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		Recovery(),
		Logging(logger),
	}
}
