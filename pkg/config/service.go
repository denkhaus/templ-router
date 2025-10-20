package config

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/kelseyhightower/envconfig"
	"github.com/samber/do/v2"
)

// configService is the private implementation of ConfigService interface
type configService struct {
	config *configImpl
}

// NewConfigService creates a new config service for DI
func NewConfigService(envVarPraefix string) func(i do.Injector) (interfaces.ConfigService, error) {
	return func(i do.Injector) (interfaces.ConfigService, error) {
		var cfg configImpl
		if err := envconfig.Process("TR", &cfg); err != nil {
			return nil, fmt.Errorf("failed to load configuration: %w", err)
		}

		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("configuration validation failed: %w", err)
		}

		if cfg.Config.PrintSummary {
			// Log configuration summary
			cfg.LogSummary()
		}

		return &configService{
			config: &cfg,
		}, nil
	}
}
