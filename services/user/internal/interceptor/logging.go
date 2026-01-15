package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var healthMethods = map[string]struct{}{
	"/grpc.health.v1.Health/Check": {},
	"/grpc.health.v1.Health/Watch": {},
}

// Logging logs the request and response of the gRPC method.
func Logging(logger *zap.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = zap.NewNop()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		if _, ok := healthMethods[info.FullMethod]; ok {
			start := time.Now()
			resp, err := handler(ctx, req)
			st := status.Convert(err)
			logger.Debug("grpc healthcheck",
				zap.String("grpc.method", info.FullMethod),
				zap.String("grpc.code", st.Code().String()),
				zap.Duration("duration", time.Since(start)),
			)
			return resp, err
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		st := status.Convert(err)

		fields := []zap.Field{
			zap.String("grpc.method", info.FullMethod),
			zap.String("grpc.code", st.Code().String()),
			zap.Duration("duration", time.Since(start)),
		}

		switch st.Code() {
		case codes.OK:
			logger.Info("grpc request", fields...)
		default:
			fields = append(fields, zap.Error(err))
			logger.Warn("grpc request", fields...)
		}

		return resp, err
	}
}
