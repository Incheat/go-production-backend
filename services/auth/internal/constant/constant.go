// Package constant defines the constants for the auth service.
package constant

const (
	// APIResponseVersionV1 is the version of the API response.
	APIResponseVersionV1 = "v1"
	// RedisRefreshTokenPrefix is the prefix for the refresh token in Redis.
	RedisRefreshTokenPrefix = "refresh_token:"
	// JWKSPath is the path for the JWKS endpoint.
	JWKSPath = "/.well-known/jwks.json"
	// ServiceName is the name of the service for the auth.
	ServiceName = "auth"
	// SpanNameAuthHTTP is the name of the span for the auth HTTP server.
	SpanNameAuthHTTP = "auth.http"
	// DefaultOTLPEndpoint is the default endpoint for the OpenTelemetry.
	DefaultOTLPEndpoint = "otel-collector:4317"
)
