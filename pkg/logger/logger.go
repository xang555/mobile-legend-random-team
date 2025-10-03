package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config captures runtime logging settings.
type Config struct {
	Level    string
	Encoding string
}

// New creates a zap logger based on configuration.
func New(cfg Config) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         cfg.Encoding,
		EncoderConfig:    encoderConfig(cfg.Encoding),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}

func parseLevel(lvl string) (zapcore.Level, error) {
	switch lvl {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("logger: unsupported level %s", lvl)
	}
}

func encoderConfig(encoding string) zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if encoding != "json" {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return cfg
}
