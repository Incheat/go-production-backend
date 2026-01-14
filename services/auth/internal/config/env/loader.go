package envconfig

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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

	authRedisHost := getString("AUTH_REDIS_HOST")
	authRedisPassword := getString("AUTH_REDIS_PASSWORD")
	authRedisDB, err := getIntRequired("AUTH_REDIS_DB")
	if err != nil {
		return nil, err
	}

	authJWTPrivateKeyPEM := getString("AUTH_JWT_PRIVATE_KEY_PEM")
	authJWTKeyID := getString("AUTH_JWT_KEY_ID")
	authJWTIssuer := getString("AUTH_JWT_ISSUER")
	authJWTAudience := getString("AUTH_JWT_AUDIENCE")
	authJWKSPath := getString("AUTH_JWT_JWKS_PATH")
	authJWTExpireRaw, err := getIntRequired("AUTH_JWT_EXPIRE")
	if err != nil {
		return nil, err
	}
	authJWTExpire := time.Duration(authJWTExpireRaw) * time.Minute

	authRefreshNumBytes, err := getIntRequired("AUTH_REFRESH_NUM_BYTES")
	if err != nil {
		return nil, err
	}
	authRefreshEndPoint := getString("AUTH_REFRESH_END_POINT")
	authRefreshMaxAge, err := getIntRequired("AUTH_REFRESH_MAX_AGE")
	if err != nil {
		return nil, err
	}

	authUserGatewayInternalAddress := getString("USER_GRPC_ADDR")

	cfg := &Config{
		Env: EnvName(env),
		Server: Server{
			PublicPort: Port(authPublicPort),
		},
		UserGateway: UserGateway{
			InternalAddress: authUserGatewayInternalAddress,
		},
		Redis: Redis{
			Host:     authRedisHost,
			Password: authRedisPassword,
			DB:       authRedisDB,
		},
		JWT: JWT{
			PrivateKeyPEM: authJWTPrivateKeyPEM,
			KeyID:         authJWTKeyID,
			Issuer:        authJWTIssuer,
			Audience:      authJWTAudience,
			Expire:        authJWTExpire,
			JWKSPath:      authJWKSPath,
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

func validate(cfg *Config) error {
	if cfg.Server.PublicPort <= 0 || cfg.Server.PublicPort > 65535 {
		return fmt.Errorf("AUTH_HTTP_PORT: must be between 1 and 65535")
	}

	if cfg.JWT.PrivateKeyPEM == "" {
		return fmt.Errorf("AUTH_JWT_PRIVATE_KEY_PEM is empty")
	}
	if cfg.JWT.KeyID == "" {
		return fmt.Errorf("AUTH_JWT_KEY_ID is empty")
	}
	if cfg.JWT.Issuer == "" {
		return fmt.Errorf("AUTH_JWT_ISSUER is empty")
	}
	if cfg.JWT.Audience == "" {
		return fmt.Errorf("AUTH_JWT_AUDIENCE is empty")
	}
	return nil
}
