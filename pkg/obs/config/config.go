// Package config defines the configuration for the observability.
package config

// DefaultOTLPEndpoint is the default OpenTelemetry Collector endpoint.
const DefaultOTLPEndpoint = "otel-collector:4317"

// TelemetryConfig defines configuration for logging, tracing,
// and metrics via OpenTelemetry.
type TelemetryConfig struct {
	// Resource describes the service emitting telemetry.
	Resource ResourceConfig

	// Logging controls application logging behavior.
	Logging LoggingConfig

	// OTLP configures the OpenTelemetry Protocol exporter.
	OTLP OTLPConfig

	// Tracing controls distributed tracing behavior.
	Tracing TracingConfig

	// Metrics controls metrics collection.
	Metrics MetricsConfig
}

// ResourceConfig describes service-level resource attributes.
type ResourceConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// LoggingConfig configures structured logging.
type LoggingConfig struct {
	Level string // e.g. debug, info, warn, error
	// JSON  bool   // enable JSON output
}

// OTLPConfig configures the OTLP exporter.
type OTLPConfig struct {
	Endpoint string // Collector endpoint (host:port)
	Insecure bool   // Disable TLS
}

// TracingConfig configures tracing behavior.
type TracingConfig struct {
	// SamplingRatio is a value between 0.0 and 1.0.
	SamplingRatio float64
}

// MetricsConfig configures metrics collection.
type MetricsConfig struct {
	// Enabled bool
}
