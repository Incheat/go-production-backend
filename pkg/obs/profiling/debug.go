// Package profiling defines the profiling server for the observability.
package profiling

import (
	"context"
	"net/http"
	"time"

	//nolint:gosec // pprof is used for debugging
	_ "net/http/pprof"

	"go.uber.org/zap"
)

// StartServer starts the profiling server.
func StartServer(_ context.Context, addr string, logger *zap.Logger) (shutdown func(context.Context) error) {

	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start profiling server", zap.Error(err))
		}
	}()

	return srv.Shutdown
}
