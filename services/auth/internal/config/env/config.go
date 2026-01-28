// Package envconfig defines the configuration for the auth service.
package envconfig

import "time"

// EnvName is the name of the environment.
type EnvName string

const (
	// EnvDev is the development environment.
	EnvDev EnvName = "dev"
	// EnvStaging is the staging environment.
	EnvStaging EnvName = "staging"
	// EnvProd is the production environment.
	EnvProd EnvName = "prod"
)

// Config is the configuration for the application.
type Config struct {
	Env         EnvName
	Version     string
	Server      Server
	Redis       Redis
	JWT         JWT
	Refresh     Refresh
	UserGateway UserGateway
	Obs         Obs
}

// Server is the configuration for the server.
type Server struct {
	HTTPPort Port
}

// UserGateway is the configuration for the user gateway.
type UserGateway struct {
	InternalAddress string
}

// Port is the port for the server.
type Port int

// Redis is the configuration for the Redis.
type Redis struct {
	Host     string
	Password string
	DB       int
}

// JWT is the configuration for the JWT.
type JWT struct {
	PrivateKeyPEM string
	KeyID         string
	Issuer        string
	Audience      string
	Expire        time.Duration
	JWKSPath      string
}

// Refresh is the configuration for the refresh.
type Refresh struct {
	NumBytes int
	EndPoint string
	MaxAge   int
}

// Obs is the configuration for the observability.
type Obs struct {
	Logging Logging
	Metrics Metrics
	Tracing Tracing
	OTLP    OTLP
}

// Logging is the configuration for the logging.
type Logging struct {
	Level string
}

// Metrics is the configuration for the metrics.
type Metrics struct {
	Port Port
}

// Tracing is the configuration for the tracing.
type Tracing struct {
	SamplingRatio float64
}

// OTLP is the configuration for the OpenTelemetry.
type OTLP struct {
	Endpoint string
}
