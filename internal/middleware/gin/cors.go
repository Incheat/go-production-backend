// Package ginmiddleware defines the middleware for the application.
package ginmiddleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSRule is a rule that defines the CORS configuration for a specific path.
type CORSRule struct {
	Path           string
	AllowedOrigins []string
}

// PathBasedCORS is a middleware that applies CORS rules based on the request path.
// It uses a simple path match algorithm to determine if the request path matches a rule.
// The middleware can be used to apply CORS rules to specific paths.
func PathBasedCORS(rules []CORSRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPath := c.FullPath()

		for _, rule := range rules {
			if matchPath(requestPath, rule.Path) {
				cors.New(cors.Config{
					AllowOrigins:     rule.AllowedOrigins,
					AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
					ExposeHeaders:    []string{"Content-Length"},
					AllowCredentials: true,
					MaxAge:           12 * time.Hour,
				})(c)
				break
			}
		}

		c.Next()
	}
}

// Simple Path Match (can be replaced with more sophisticated glob/pattern match).
func matchPath(requestPath, rulePath string) bool {
	if rulePath == requestPath {
		return true
	}
	// Support simple wildcard.
	if len(rulePath) > 0 && rulePath[len(rulePath)-1] == '*' {
		prefix := rulePath[:len(rulePath)-1]
		return len(requestPath) >= len(prefix) && requestPath[:len(prefix)] == prefix
	}
	return false
}
