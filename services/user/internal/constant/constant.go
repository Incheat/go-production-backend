// Package constant defines the constants for the user service.
package constant

const (
	// EnvKey is the key for the environment. APP_ENV => env => ex: "test" / "staging" / "prod"
	EnvKey = "env"
	// EnvPrefix is the prefix for the environment. APP_SERVER_PORT => server.port
	EnvPrefix = "APP_"
	// EnvConfigDir is the directory for the configuration files.
	EnvConfigDir = "config"
	// EnvConfigTmpl is the template for the configuration files.
	EnvConfigTmpl = "config.%s.yaml"
)
