// Package middleware defines the middleware for the auth service.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger is a middleware that logs the request and response.
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // process request

		latency := time.Since(start)
		status := c.Writer.Status()

		reqID := c.GetString("request_id")
		if reqID == "" {
			reqID = "-"
		}

		logger.Info("request handled",
			zap.String("request_id", reqID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.Any("errors", c.Errors),
		)
	}
}
