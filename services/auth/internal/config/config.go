// Package config defines the configuration for the auth service.
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

	MySQL struct {
		DSN string `koanf:"dsn"` // e.g. user:pass@tcp(mysql:3306)/auth?parseTime=true
	} `koanf:"mysql"`

	Redis struct {
		Addr     string `koanf:"addr"`
		Password string `koanf:"password"`
		DB       int    `koanf:"db"`
	} `koanf:"redis"`

	JWT struct {
		Secret string `koanf:"secret"`
		Expire int    `koanf:"expire"` // minutes
	} `koanf:"jwt"`

	Refresh struct {
		NumBytes int    `koanf:"num_bytes"`
		EndPoint string `koanf:"end_point"`
		MaxAge   int    `koanf:"max_age"` // seconds
	} `koanf:"refresh"`
}

// CORSRule is a rule that defines the CORS configuration for a specific path.
type CORSRule struct {
	Path           string   `koanf:"path"`
	AllowedOrigins []string `koanf:"allowed_origins"`
}
