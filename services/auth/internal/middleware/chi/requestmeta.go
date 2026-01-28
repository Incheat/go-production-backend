// Package chimiddleware defines the middleware for the auth service.
package chimiddleware

import (
	"net"
	"net/http"

	chimiddlewareutils "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi/utils"
)

// RequestMeta adds the request metadata to the context.
func RequestMeta() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			meta := chimiddlewareutils.RequestMeta{
				UserAgent: r.UserAgent(),
				IPAddress: getClientIP(r),
			}
			ctx := chimiddlewareutils.WithRequestMeta(r.Context(), meta)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getClientIP(httpRequest *http.Request) string {
	if httpRequest == nil {
		return ""
	}
	// Try real IP from common proxy headers
	ipAddress := httpRequest.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = httpRequest.Header.Get("X-Real-IP")
	}
	if ipAddress == "" {
		// fallback: use connection remote address
		// but strip port (ip:port)
		host, _, err := net.SplitHostPort(httpRequest.RemoteAddr)
		if err == nil {
			ipAddress = host
		} else {
			ipAddress = httpRequest.RemoteAddr
		}
	}
	return ipAddress
}
