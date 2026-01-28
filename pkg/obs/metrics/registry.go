package metrics

import "github.com/prometheus/client_golang/prometheus"

// NewRegistry gives each service its own registry (cleaner than global default).
func NewRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}
