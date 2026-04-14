package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all values the app needs to start.
type Config struct {
	DBDsn   string // PostgreSQL connection string
	BaseURL string // public-facing server URL, e.g. "http://localhost:3000"
	Port    string // port the server listens on, default "3000"
	AppEnv  string // "development" or "production"
}

// Load reads config from a .env file (dev) or real environment variables (production).
func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error — production won't have a .env file

	cfg := &Config{
		DBDsn:   getEnv("DB_DSN", ""),
		BaseURL: getEnv("BASE_URL", "http://localhost:3000"),
		Port:    getEnv("PORT", "3000"),
		AppEnv:  getEnv("APP_ENV", "development"),
	}

	if cfg.DBDsn == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// IsDevelopment returns true when running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}
