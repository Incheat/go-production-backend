// Package logging defines the gRPC interceptor for the observability.
package logging

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/incheat/go-production-backend/pkg/obs/correlation"
)

// GRPCRequestLogging logs gRPC unary requests and injects a request-scoped logger into ctx.
func GRPCRequestLogging(base *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		// start := time.Now()

		// Build request-scoped logger
		l := base.With(
			zap.String("grpc.method", info.FullMethod),
		).With(correlation.TraceFields(ctx)...)

		l = l.With(correlation.BaggageFields(ctx,
			correlation.BaggageRequestID,
			correlation.BaggageTenantID,
		)...)

		// Optional: remote peer address
		if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
			l = l.With(zap.String("client_addr", p.Addr.String()))
		}

		// Put logger into ctx for handlers
		ctx = correlation.ContextWithLogger(ctx, l)

		resp, err = handler(ctx, req)

		// latency := time.Since(start)
		// st := status.Convert(err)

		// // Access log
		// if err != nil {
		// 	l.Warn("grpc request handled",
		// 		zap.String("grpc.code", st.Code().String()),
		// 		zap.Duration("latency", latency),
		// 		zap.Error(err),
		// 	)
		// } else {
		// 	l.Info("grpc request handled",
		// 		zap.String("grpc.code", st.Code().String()),
		// 		zap.Duration("latency", latency),
		// 	)
		// }

		return resp, err
	}
}
