// Package correlation defines the logger context for the observability.
package correlation

import (
	"context"

	"go.uber.org/zap"
)

type ctxKeyLogger struct{}

// ContextWithLogger adds the logger to the context.
func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger{}, l)
}

// LoggerFromContext gets the logger from the context.
func LoggerFromContext(ctx context.Context) *zap.Logger {
	l, ok := FromContext(ctx)
	if !ok {
		fallbackLogger := zap.L()
		fallbackLogger.Warn("logger not found in context, falling back to global logger")
		return fallbackLogger
	}
	return l
}

// FromContext gets the logger from the context.
func FromContext(ctx context.Context) (*zap.Logger, bool) {
	l, ok := ctx.Value(ctxKeyLogger{}).(*zap.Logger)
	return l, ok
}
