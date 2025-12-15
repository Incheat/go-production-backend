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

	authPublicPort, err := getIntRequired("AUTH_PUBLIC_PORT")
	if err != nil {
		return nil, err
	}

	authCorsPublicAllowedOrigins := parseOrigins(getString("AUTH_CORS_PUBLIC_ALLOWED_ORIGINS"))
	authCorsPublicAllowCredentials, err := getBoolRequired("AUTH_CORS_PUBLIC_ALLOW_CREDENTIALS")
	if err != nil {
		return nil, err
	}

	authCorsInternalAllowedOrigins := parseOrigins(getString("AUTH_CORS_INTERNAL_ALLOWED_ORIGINS"))
	authCorsInternalAllowCredentials, err := getBoolRequired("AUTH_CORS_INTERNAL_ALLOW_CREDENTIALS")
	if err != nil {
		return nil, err
	}

	authRedisHost := getString("AUTH_REDIS_HOST")
	authRedisPassword := getString("AUTH_REDIS_PASSWORD")
	authRedisDB, err := getIntRequired("AUTH_REDIS_DB")
	if err != nil {
		return nil, err
	}

	authJWTSecret := getString("AUTH_JWT_SECRET")
	authJWTExpire, err := getIntRequired("AUTH_JWT_EXPIRE")
	if err != nil {
		return nil, err
	}

	authRefreshNumBytes, err := getIntRequired("AUTH_REFRESH_NUM_BYTES")
	if err != nil {
		return nil, err
	}
	authRefreshEndPoint := getString("AUTH_REFRESH_END_POINT")
	authRefreshMaxAge, err := getIntRequired("AUTH_REFRESH_MAX_AGE")
	if err != nil {
		return nil, err
	}

	authUserGatewayInternalPort, err := getIntRequired("USER_INTERNAL_PORT")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Env: EnvName(env),
		Server: Server{
			PublicPort: Port(authPublicPort),
		},
		UserGateway: UserGateway{
			InternalPort: Port(authUserGatewayInternalPort),
		},
		CORS: CORS{
			Internal: CORSRule{
				AllowedOrigins:   authCorsInternalAllowedOrigins,
				AllowCredentials: authCorsInternalAllowCredentials,
			},
			Public: CORSRule{
				AllowedOrigins:   authCorsPublicAllowedOrigins,
				AllowCredentials: authCorsPublicAllowCredentials,
			},
		},
		Redis: Redis{
			Host:     authRedisHost,
			Password: authRedisPassword,
			DB:       authRedisDB,
		},
		JWT: JWT{
			Secret: authJWTSecret,
			Expire: authJWTExpire,
		},
		Refresh: Refresh{
			NumBytes: authRefreshNumBytes,
			EndPoint: authRefreshEndPoint,
			MaxAge:   authRefreshMaxAge,
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

func getBoolRequired(name string) (bool, error) {
	raw := getString(name)
	if raw == "" {
		return false, fmt.Errorf("%s: %w", name, errMissingEnv)
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s: %w", name, err)
	}
	return v, nil
}

// parseOrigins supports:
// - "" => nil
// - "*" => []string{"*"}  (so CORS layer can treat it as allow-all)
// - "a,b,c" => []string{"a","b","c"} (trimmed, empties removed)
func parseOrigins(env string) []string {
	env = strings.TrimSpace(env)
	if env == "" {
		return nil
	}
	if env == "*" {
		return []string{"*"}
	}

	parts := strings.Split(env, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func validate(cfg *Config) error {
	if cfg.Server.PublicPort <= 0 || cfg.Server.PublicPort > 65535 {
		return fmt.Errorf("AUTH_HTTP_PORT: must be between 1 and 65535")
	}

	if cfg.JWT.Secret == "" {
		return fmt.Errorf("AUTH_JWT_SECRET: %w", errMissingEnv)
	}
	return nil
}
