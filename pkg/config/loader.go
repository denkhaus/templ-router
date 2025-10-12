package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/samber/do/v2"
)

// NewService creates the configuration service for DI
func NewService(injector do.Injector) (*Config, error) {
	var cfg Config
	if err := envconfig.Process("FB", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Log configuration summary
	cfg.LogSummary()

	return &cfg, nil
}
