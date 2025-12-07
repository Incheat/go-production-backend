// Package config defines the configuration for the helloworld service.
package config

// Config is the configuration for the application.
type Config struct {
	Server struct {
		Port int `koanf:"port"`
	} `koanf:"server"`

	CORS struct {
		Rules []CORSRule `koanf:"rules"`
	} `koanf:"cors"`

	Database struct {
		Host     string `koanf:"host"`
		Port     int    `koanf:"port"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
		Name     string `koanf:"name"`
	} `koanf:"database"`

	Redis struct {
		Addr     string `koanf:"addr"`
		Password string `koanf:"password"`
		DB       int    `koanf:"db"`
	} `koanf:"redis"`

	JWT struct {
		Secret string `koanf:"secret"`
		Expire int    `koanf:"expire"` // seconds
	} `koanf:"jwt"`
}

// CORSRule is a rule that defines the CORS configuration for a specific path.
type CORSRule struct {
	Path           string   `koanf:"path"`
	AllowedOrigins []string `koanf:"allowed_origins"`
}
