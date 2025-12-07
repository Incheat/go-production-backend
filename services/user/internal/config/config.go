// Package config defines the configuration for the user service.
package config

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
	Env EnvName `koanf:"env"`

	Server struct {
		Port int `koanf:"port"`
	} `koanf:"server"`

	CORS struct {
		Rules []CORSRule `koanf:"rules"`
	} `koanf:"cors"`
}

// CORSRule is a rule that defines the CORS configuration for a specific path.
type CORSRule struct {
	Path           string   `koanf:"path"`
	AllowedOrigins []string `koanf:"allowed_origins"`
}
