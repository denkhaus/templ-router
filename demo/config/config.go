package config

import (
	"os"
)

type Config struct {
	Environment string
}

func NewConfig() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}