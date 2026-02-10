// Package otel defines the OTLP exporter for the observability.
package otel

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

// OTLPConfig is the configuration for the OTLP exporter.
type OTLPConfig struct {
	Endpoint string
	Insecure bool
}

// NewTraceExporter creates a new trace exporter.
func NewTraceExporter(ctx context.Context, cfg OTLPConfig) (*otlptrace.Exporter, error) {

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.New(ctx, opts...)
}
