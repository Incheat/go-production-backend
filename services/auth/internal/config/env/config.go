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
	Env         EnvName
	Server      Server
	CORS        CORS
	Redis       Redis
	JWT         JWT
	Refresh     Refresh
	UserGateway UserGateway
}

// Server is the configuration for the server.
type Server struct {
	PublicPort Port
}

// UserGateway is the configuration for the user gateway.
type UserGateway struct {
	InternalPort Port
}

// Port is the port for the server.
type Port int

// CORS is the configuration for the CORS.
type CORS struct {
	Internal CORSRule
	Public   CORSRule
}

// CORSRule is the configuration for the CORS rule.
type CORSRule struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

// Redis is the configuration for the Redis.
type Redis struct {
	Host     string
	Password string
	DB       int
}

// JWT is the configuration for the JWT.
type JWT struct {
	Secret string
	Expire int
}

// Refresh is the configuration for the refresh.
type Refresh struct {
	NumBytes int
	EndPoint string
	MaxAge   int
}
