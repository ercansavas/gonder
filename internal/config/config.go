package config

import (
	"os"
)

// Config represents application configuration
type Config struct {
	Port     string
	Host     string
	LogLevel string
}

// Load loads configuration from environment variables or default values
func Load() *Config {
	cfg := &Config{
		Port:     getEnv("PORT", "8080"),
		Host:     getEnv("HOST", "localhost"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
	return cfg
}

// getEnv gets environment variable, returns default value if not found
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
