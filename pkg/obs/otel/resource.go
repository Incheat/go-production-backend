package otel

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// ResourceConfig is the configuration for the resource.
type ResourceConfig struct {
	ServiceName string
	// add more later (version/env/instance) if needed.
}

// NewResource creates a new resource.
func NewResource(ctx context.Context, cfg ResourceConfig) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
}
