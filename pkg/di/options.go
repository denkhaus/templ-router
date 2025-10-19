package di

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)

// ApplicationOption defines an option for configuring application services
type ApplicationOption func(c *Container)

// WithTemplateRegistry sets the template registry
func WithTemplateRegistry(registry interfaces.TemplateRegistry) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, registry)
	}
}

// WithAssetsService sets the assets service
func WithAssetsService(assetsService interfaces.AssetsService) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, assetsService)
	}
}

// WithUserStore sets a custom user store implementation
func WithUserStore(userStore interfaces.UserStore) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, userStore)
	}
}

// WithAuthHandlers sets custom authentication handlers
func WithAuthHandlers(authHandlers interfaces.AuthHandlers) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, authHandlers)
	}
}
