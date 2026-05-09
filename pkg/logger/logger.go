// Package logger provides structured logging using zap.
// Layer -1: Infrastructure package — can be imported by any layer.
package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger for structured logging.
type Logger struct {
	*zap.Logger
}

// Field represents a zap field.
type Field = zap.Field

// NewProduction creates a production logger.
func NewProduction() (*Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

// NewDevelopment creates a development logger.
func NewDevelopment() (*Logger, error) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &Logger{Logger: zapLogger}, nil
}

// String creates a string field.
func String(key, val string) Field {
	return zap.String(key, val)
}

// Int creates an int field.
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// Error creates an error field.
func Error(err error) Field {
	return zap.Error(err)
}

// Any creates an any field.
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

// Duration creates a duration field.
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}
