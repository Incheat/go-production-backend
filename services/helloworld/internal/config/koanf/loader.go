package koanfconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	envKey        = "APP_ENV" // ex: "test" / "staging" / "prod"
	envPrefix     = "APP_"    // APP_SERVER_PORT => server.port
	envConfigDir  = "config"
	baseConfig    = "config.yaml"
	envConfigTmpl = "config.%s.yaml"
)

// Load loads the configuration from the environment variables and the configuration file.
func Load() (*Config, error) {
	k := koanf.New(".")

	// 1. Base config.yaml
	if err := loadYAMLIfExists(k, filepath.Join(envConfigDir, baseConfig)); err != nil {
		return nil, fmt.Errorf("load base config: %w", err)
	}

	// 2. Environment-specific config.<env>.yaml
	envName := strings.TrimSpace(os.Getenv(envKey)) // ex: "test" / "staging" / "prod"
	if envName != "" {
		envFile := filepath.Join(envConfigDir, fmt.Sprintf(envConfigTmpl, envName))
		if err := loadYAMLIfExists(k, envFile); err != nil {
			return nil, fmt.Errorf("load %s config: %w", envName, err)
		}
	}

	// 3. Environment variable overrides: APP_SERVER_PORT -> server.port
	if err := k.Load(env.Provider(envPrefix, ".", normalizeEnvKey), nil); err != nil {
		return nil, fmt.Errorf("load env overrides: %w", err)
	}

	// 4. Unmarshal into struct
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
	s = strings.TrimPrefix(s, envPrefix)
	s = strings.ToLower(s)
	return strings.ReplaceAll(s, "_", ".")
}
