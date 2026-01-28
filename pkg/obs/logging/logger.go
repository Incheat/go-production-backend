// Package logging defines the logger for the observability.
package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config is the configuration for the logger.
type Config struct {
	Service string // "auth-service"
	Env     string // "dev" / "prod"
	Level   string // "debug" / "info" / "warn" / "error"
}

// New creates a new logger.
func New(cfg Config) (*zap.Logger, error) {
	var zcfg zap.Config

	if cfg.Env == "prod" {
		zcfg = zap.NewProductionConfig()
	} else {
		zcfg = zap.NewDevelopmentConfig()
	}

	lvl := zapcore.InfoLevel
	_ = lvl.Set(cfg.Level) // if invalid, stays Info

	zcfg.Level = zap.NewAtomicLevelAt(lvl)
	zcfg.EncoderConfig.TimeKey = "ts"

	logger, err := zcfg.Build()
	if err != nil {
		return nil, err
	}

	// Add stable service fields once
	return logger.With(
		zap.String("service", cfg.Service),
		zap.String("env", cfg.Env),
	), nil
}
