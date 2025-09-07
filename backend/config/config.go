package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	JWTSecret       string
	JWTExpiryPeriod time.Duration
	Port            string
	Environment     string
}

var cfg *Config

// Load initializes configuration from environment variables
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	// Load .env file if it exists
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env") // Try loading from backend root if running from cmd/app

	cfg = &Config{
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		JWTExpiryPeriod: 24 * time.Hour, // 24 hours
		Port:            getEnv("PORT", "8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}

	// Validate required configuration
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		panic("Configuration not loaded")
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
