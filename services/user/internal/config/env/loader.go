package envconfig

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// errMissingEnv is the error returned when a required environment variable is missing.
var errMissingEnv = errors.New("missing env var")

// Load loads the configuration from the environment variables.
func Load() (*Config, error) {
	env := getString("ENV") // if missing, it will be ""

	userVersion := getString("USER_VERSION")
	if userVersion == "" {
		return nil, fmt.Errorf("USER_VERSION: %w", errMissingEnv)
	}

	userGrpcPort, err := getIntRequired("USER_GRPC_PORT")
	if err != nil {
		return nil, err
	}

	userMySQLHost := getString("USER_MYSQL_HOST")
	userMySQLUser := getString("USER_MYSQL_USER")
	userMySQLPassword := getString("USER_MYSQL_PASSWORD")
	userMySQLDBName := getString("USER_MYSQL_DB_NAME")
	userMySQLMaxOpenConns, err := getIntRequired("USER_MYSQL_MAX_OPEN_CONNS")
	if err != nil {
		return nil, err
	}
	userMySQLMaxIdleConns, err := getIntRequired("USER_MYSQL_MAX_IDLE_CONNS")
	if err != nil {
		return nil, err
	}
	userMySQLConnMaxLifetime, err := getIntRequired("USER_MYSQL_CONN_MAX_LIFETIME")
	if err != nil {
		return nil, err
	}

	userMetricsPort, err := getIntRequired("PROM_METRICS_PORT")
	if err != nil {
		return nil, err
	}
	userLoggingLevel := getString("USER_LOGGING_LEVEL")
	userTracingSamplingRatio, err := getFloat64Required("USER_TRACING_SAMPLING_RATIO")
	if err != nil {
		return nil, err
	}
	otelEndpoint := getString("OTEL_ENDPOINT")

	cfg := &Config{
		Env: EnvName(env),
		Server: Server{
			GrpcPort: Port(userGrpcPort),
		},
		MySQL: MySQL{
			Host:            userMySQLHost,
			User:            userMySQLUser,
			Password:        userMySQLPassword,
			DBName:          userMySQLDBName,
			MaxOpenConns:    userMySQLMaxOpenConns,
			MaxIdleConns:    userMySQLMaxIdleConns,
			ConnMaxLifetime: userMySQLConnMaxLifetime,
		},
		Obs: Obs{
			Logging: Logging{
				Level: userLoggingLevel,
			},
			Metrics: Metrics{
				Port: Port(userMetricsPort),
			},
			Tracing: Tracing{
				SamplingRatio: userTracingSamplingRatio,
			},
			OTLP: OTLP{
				Endpoint: otelEndpoint,
			},
		},
	}

	// Optional sanity checks (keep or remove as you like)
	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getString(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

func getIntRequired(name string) (int, error) {
	raw := getString(name)
	if raw == "" {
		return 0, fmt.Errorf("%s: %w", name, errMissingEnv)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return v, nil
}

func getFloat64Required(name string) (float64, error) {
	raw := getString(name)
	if raw == "" {
		return 0, fmt.Errorf("%s: %w", name, errMissingEnv)
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return v, nil
}

func validate(cfg *Config) error {
	if cfg.Server.GrpcPort <= 0 || cfg.Server.GrpcPort > 65535 {
		return fmt.Errorf("USER_PUBLIC_PORT: must be between 1 and 65535")
	}

	return nil
}
