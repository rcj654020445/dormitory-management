// Package config provides configuration loading.
// Layer -1: Infrastructure package — can be imported by any layer.
package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds application configuration.
type Config struct {
	Port          string
	Env           string
	LogLevel      string
	DatabaseURL   string
	JWTSecret     string
	CORSOrigins   []string
}

// Load loads configuration from environment variables.
func Load(configDir string) (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		CORSOrigins: parseOrigins(getEnv("CORS_ORIGINS", "http://localhost:3000")),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseOrigins(origins string) []string {
	if origins == "" {
		return []string{"http://localhost:3000"}
	}
	return strings.Split(origins, ",")
}

// Error represents a configuration error.
type Error struct {
	Field   string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("config error: %s - %s", e.Field, e.Message)
}
