// Package ginmiddleware defines the validator for the OpenAPI specification.
package ginmiddleware

import (
	"log"

	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
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
func NewValidatorOptions(cfg ValidatorConfig) *ginmiddleware.Options {
	if cfg.Logger == nil {
		cfg.Logger = log.Printf
	}
	if cfg.ProdError == "" {
		cfg.ProdError = "invalid request"
	}

	return &ginmiddleware.Options{
		ErrorHandler: func(c *gin.Context, message string, statusCode int) {
			cfg.Logger("validation error (%d): %s", statusCode, message)

			if cfg.ProdMode {
				c.AbortWithStatusJSON(statusCode, gin.H{"error": cfg.ProdError})
				return
			}

			c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
		},
	}
}
