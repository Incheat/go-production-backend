package otel

import (
	"context"
	"errors"
)

// ShutdownFunc is the function to shutdown the observability.
type ShutdownFunc func(context.Context) error

// MultiShutdown shuts down multiple observability components.
func MultiShutdown(funcs ...ShutdownFunc) ShutdownFunc {
	return func(ctx context.Context) error {
		var errs []error
		for i := len(funcs) - 1; i >= 0; i-- {
			if funcs[i] == nil {
				continue
			}
			if err := funcs[i](ctx); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}
}
