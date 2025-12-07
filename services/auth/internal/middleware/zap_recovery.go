// Package middleware defines the middleware for the auth service.
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapRecovery is a middleware that recovers from any panics and logs the error.
func ZapRecovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("panic recovered",
			zap.Any("error", recovered),
			zap.String("path", c.Request.URL.Path),
		)
		c.AbortWithStatus(500)
	})
}
