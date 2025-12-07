package config

import (
	"strings"
	"testing"

	"github.com/incheat/go-playground/services/user/internal/constant"
)

// TestLoad_InvalidEnv verifies that Load returns an error
// when the environment value cannot be treated as a valid EnvName.
func TestLoad_InvalidEnv(t *testing.T) {
	t.Helper()

	// Set an invalid environment name.
	//
	// NOTE: this assumes your normalizeEnvKey + EnvPrefix logic maps
	// something like APP_ENV -> env (or whatever constant.EnvKey is).
	t.Setenv(constant.EnvPrefix+"ENV", "totally-invalid-env")

	cfg, err := Load()
	if err == nil {
		t.Fatalf("expected error, got nil (cfg=%#v)", cfg)
	}

	if !strings.Contains(err.Error(), "invalid environment") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestLoad_ValidEnv_NoFile verifies that Load succeeds when a valid
// environment is set and the env-specific config file is either absent
// or loadable without error.
//
// This relies on loadYAMLIfExists *not* failing when the file is missing
// (which its name strongly implies).
func TestLoad_ValidEnv_NoFile(t *testing.T) {
	tests := []struct {
		name    string
		env     EnvName
		wantErr bool
	}{
		{"dev", EnvDev, false},
		{"staging", EnvStaging, false},
		{"prod", EnvProd, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv(constant.EnvPrefix+"ENV", string(test.env))
			cfg, err := Load()
			if (err != nil) != test.wantErr {
				t.Fatalf("expected error: %v, got: %v", test.wantErr, err)
			}
			if cfg == nil {
				t.Fatalf("expected non-nil config")
			}
			if cfg.Env != test.env {
				t.Fatalf("expected Env=%q, got %q", test.env, cfg.Env)
			}
		})
	}
}

// TestMustLoad_PanicsOnError ensures MustLoad panics when Load fails.
func TestMustLoad_PanicsOnError(t *testing.T) {
	t.Helper()

	// Force an error by setting an invalid env again.
	t.Setenv(constant.EnvPrefix+"ENV", "bad-env")

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected MustLoad to panic, but it did not")
		}
	}()

	_ = MustLoad() // should panic
}
