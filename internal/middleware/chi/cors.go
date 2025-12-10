// Package chimiddleware defines the middleware for the application.
package chimiddleware

import (
	"net/http"
	"time"

	chicors "github.com/go-chi/cors"
)

// CORSRule defines the CORS configuration for a specific path.
type CORSRule struct {
	Path           string
	AllowedOrigins []string
}

// PathBasedCORS is a Chi middleware that applies CORS rules based on the request path.
// It supports exact match and simple wildcard suffixes like "/api/*".
func PathBasedCORS(rules []CORSRule) func(next http.Handler) http.Handler {
	// Prebuild CORS handlers for performance
	corsHandlers := make([]struct {
		Path    string
		Handler func(http.Handler) http.Handler
	}, len(rules))

	for i, rule := range rules {
		corsHandlers[i] = struct {
			Path    string
			Handler func(http.Handler) http.Handler
		}{
			Path: rule.Path,
			Handler: chicors.Handler(chicors.Options{
				AllowedOrigins:   rule.AllowedOrigins,
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
				ExposedHeaders:   []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           int((12 * time.Hour).Seconds()),
			}),
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestPath := r.URL.Path

			for _, ch := range corsHandlers {
				if matchPath(requestPath, ch.Path) {
					ch.Handler(next).ServeHTTP(w, r)
					return
				}
			}

			// No rule matched → just continue
			next.ServeHTTP(w, r)
		})
	}
}

// matchPath matches exact path or supports "*" wildcard at the end.
func matchPath(requestPath, rulePath string) bool {
	if rulePath == requestPath {
		return true
	}

	// Wildcard handling like "/api/*"
	if len(rulePath) > 0 && rulePath[len(rulePath)-1] == '*' {
		prefix := rulePath[:len(rulePath)-1]
		return len(requestPath) >= len(prefix) && requestPath[:len(prefix)] == prefix
	}

	return false
}
