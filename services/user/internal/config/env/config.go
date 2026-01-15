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
	Env    EnvName
	Server Server
	MySQL  MySQL
}

// Server is the configuration for the server.
type Server struct {
	InternalPort Port
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
