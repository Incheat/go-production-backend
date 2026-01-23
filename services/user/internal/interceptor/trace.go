package interceptor

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type loggerKey struct{}

// ZapTraceUnaryInterceptor adds Zap logging to the gRPC request.
func ZapTraceUnaryInterceptor(base *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		sc := trace.SpanFromContext(ctx).SpanContext()
		l := base
		if sc.IsValid() {
			l = l.With(
				zap.String("trace_id", sc.TraceID().String()),
				zap.String("span_id", sc.SpanID().String()),
			)
		}

		ctx = context.WithValue(ctx, loggerKey{}, l)

		return handler(ctx, req)
	}
}
