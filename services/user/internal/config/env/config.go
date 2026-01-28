// Package envconfig defines the configuration for the auth service.
package envconfig

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
	Env     EnvName
	Version string
	Server  Server
	MySQL   MySQL
	Obs     Obs
}

// Server is the configuration for the server.
type Server struct {
	GrpcPort Port
}

// Port is the port for the server.
type Port int

// MySQL is the configuration for the MySQL.
type MySQL struct {
	User            string
	Password        string
	Host            string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // seconds
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
