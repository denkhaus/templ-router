package config

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
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
		if err := envconfig.Process(envVarPraefix, &cfg); err != nil {
			return nil, shared.NewConfigurationError("failed to load configuration from environment variables").
				WithDetails("Environment variable processing failed").
				WithCause(err).
				WithContext("prefix", envVarPraefix)
		}

		if err := cfg.Validate(); err != nil {
			return nil, shared.NewConfigurationError("configuration validation failed").
				WithDetails("Configuration values do not meet validation requirements").
				WithCause(err)
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
