// Package correlation defines the baggage for the observability.
package correlation

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
	"go.uber.org/zap"
)

// BaggageKey defines allowed baggage keys.
// Keep this small and stable.
type BaggageKey string

const (
	// BaggageRequestID is the request ID.
	BaggageRequestID BaggageKey = "request.id"
	// BaggageTenantID is the tenant ID.
	BaggageTenantID BaggageKey = "tenant.id"
)

// SetBaggage adds or overwrites a baggage key/value on ctx.
func SetBaggage(ctx context.Context, key BaggageKey, value string) (context.Context, error) {
	if value == "" {
		return ctx, nil
	}

	member, err := baggage.NewMember(string(key), value)
	if err != nil {
		return ctx, err
	}

	bg := baggage.FromContext(ctx)
	bg, err = bg.SetMember(member)
	if err != nil {
		return ctx, err
	}

	return baggage.ContextWithBaggage(ctx, bg), nil
}

// BaggageFields extracts selected baggage keys from ctx as zap fields.
func BaggageFields(ctx context.Context, keys ...BaggageKey) []zap.Field {
	bg := baggage.FromContext(ctx)
	if bg.Len() == 0 {
		return nil
	}

	fields := make([]zap.Field, 0, len(keys))
	for _, k := range keys {
		v := bg.Member(string(k)).Value()
		if v == "" {
			continue
		}
		fields = append(fields, zap.String(string(k), v))
	}

	if len(fields) == 0 {
		return nil
	}
	return fields
}
