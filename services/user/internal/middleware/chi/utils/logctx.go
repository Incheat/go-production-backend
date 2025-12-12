// Package chimiddlewareutils defines the context for the chi middleware utils.
package chimiddlewareutils

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

// WithLogger adds a logger to the context.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger gets the logger from the context.
func GetLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		// fallback to global logger
		return zap.L()
	}
	return logger
}
