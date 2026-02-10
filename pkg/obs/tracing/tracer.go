// Package tracing defines the tracer for the observability.
package tracing

import (
	"context"
	"time"

	"github.com/incheat/go-production-backend/pkg/obs/config"
	obsotel "github.com/incheat/go-production-backend/pkg/obs/otel"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Shutdown is the function to shutdown the tracer.
type Shutdown func(context.Context) error

// InitTracer initializes the tracer.
func InitTracer(ctx context.Context, cfg config.TelemetryConfig) (Shutdown, error) {
	// endpoint default
	endpoint := cfg.OTLP.Endpoint
	if endpoint == "" {
		// Prefer shared default:
		endpoint = config.DefaultOTLPEndpoint
	}

	exp, err := obsotel.NewTraceExporter(ctx, obsotel.OTLPConfig{
		Endpoint: endpoint,
		Insecure: cfg.OTLP.Insecure,
	})
	if err != nil {
		return nil, err
	}

	res, err := obsotel.NewResource(ctx, obsotel.ResourceConfig{
		ServiceName: cfg.Resource.ServiceName,
	})
	if err != nil {
		return nil, err
	}

	// sampler
	sampler := sdktrace.ParentBased(sdktrace.AlwaysSample())
	// if cfg.Tracing.SamplingRatio < 1.0 && cfg.Tracing.SamplingRatio >= 0.0{
	// 	sampler = sdktrace.TraceIDRatioBased(cfg.Tracing.SamplingRatio)
	// }

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	// Propagate trace context (W3C + B3 + baggage)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
			b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader|b3.B3SingleHeader)),
		),
	)

	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}, nil
}
