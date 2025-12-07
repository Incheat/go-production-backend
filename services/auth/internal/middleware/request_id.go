// Package middleware defines the middleware for the auth service.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HeaderRequestID is the header name for the request ID.
const HeaderRequestID = "X-Request-ID"

// RequestID is a middleware that generates a request ID and stores it in the context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(HeaderRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		// store in context so ZapLogger can log it
		c.Set("request_id", reqID)

		// write to response for client tracing
		c.Writer.Header().Set(HeaderRequestID, reqID)

		c.Next()
	}
}
