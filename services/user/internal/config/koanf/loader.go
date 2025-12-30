// Package koanfconfig defines the loader for the user service.
package koanfconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/incheat/go-production-backend/services/user/internal/constant"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Load loads the configuration from the environment variables and the configuration file.
func Load() (*Config, error) {

	k := koanf.New(".")

	if err := k.Load(env.Provider(constant.EnvPrefix, ".", normalizeEnvKey), nil); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	envRaw := k.String(constant.EnvKey)
	envName := EnvName(envRaw)

	switch envName {
	case EnvDev, EnvStaging, EnvProd:
		// ok
	default:
		return nil, fmt.Errorf("invalid environment: %s", envRaw)
	}

	envFile := filepath.Join(constant.EnvConfigDir, fmt.Sprintf(constant.EnvConfigTmpl, envName))
	if err := loadYAMLIfExists(k, envFile); err != nil {
		return nil, fmt.Errorf("load %s config: %w", envName, err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// MustLoad is convenient for main(); it panics on error.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err) // or log.Fatalf in main
	}
	return cfg
}

// loadYAMLIfExists loads a YAML file into koanf if it exists.
// If the file doesn't exist, it silently returns nil.
func loadYAMLIfExists(k *koanf.Koanf, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("File %s does not exist\n", path)
			return nil
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		// probably a misconfiguration; treat as error
		return fmt.Errorf("%s is a directory, expected file", path)
	}

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return fmt.Errorf("load %s: %w", path, err)
	}
	return nil
}

// normalizeEnvKey turns APP_SERVER_PORT -> server.port
func normalizeEnvKey(s string) string {
	s = strings.TrimPrefix(s, constant.EnvPrefix)
	s = strings.ToLower(s)
	return strings.ReplaceAll(s, "_", ".")
}
