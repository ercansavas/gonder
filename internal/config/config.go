package config

import (
	"os"
)

// Config uygulama konfigürasyonunu temsil eder
type Config struct {
	Port     string
	Host     string
	LogLevel string
}

// Load çevre değişkenlerinden veya varsayılan değerlerden konfigürasyon yükler
func Load() *Config {
	cfg := &Config{
		Port:     getEnv("PORT", "8080"),
		Host:     getEnv("HOST", "localhost"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
	return cfg
}

// getEnv çevre değişkenini alır, yoksa varsayılan değeri döner
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
