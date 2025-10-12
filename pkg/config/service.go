package config

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)

// configService is the private implementation of ConfigService interface
type configService struct {
	config *Config
}

// NewConfigService creates a new config service for DI
func NewConfigService(i do.Injector) (interfaces.ConfigService, error) {
	config, err := NewService(i)
	if err != nil {
		return nil, err
	}

	return &configService{
		config: config,
	}, nil
}
