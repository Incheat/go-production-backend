// Package chimiddleware defines the validator for the OpenAPI specification.
package chimiddleware

import (
	"encoding/json"
	"log"
	"net/http"

	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

// ValidatorConfig is the configuration for the validator.
type ValidatorConfig struct {
	ProdMode  bool
	ProdError string
	Logger    func(format string, args ...any)
}

// NewValidatorOptions creates a new validator options.
// If ProdMode is true, the validator will return a production error message.
// If ProdMode is false, the validator will return a development error message.
func NewValidatorOptions(cfg ValidatorConfig) *nethttpmiddleware.Options {
	if cfg.Logger == nil {
		cfg.Logger = log.Printf
	}
	if cfg.ProdError == "" {
		cfg.ProdError = "invalid request"
	}

	return &nethttpmiddleware.Options{
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			cfg.Logger("validation error (%d): %s", statusCode, message)

			errMsg := message
			if cfg.ProdMode {
				errMsg = cfg.ProdError
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)

			// Best-effort JSON response; fall back to plain text on error.
			if err := json.NewEncoder(w).Encode(map[string]string{"error": errMsg}); err != nil {
				http.Error(w, errMsg, statusCode)
			}
		},
	}
}
