package di

import (
	"net/http"

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

// WithRouterMiddleware enables or disables router-level middleware
func WithRouterMiddleware(enabled bool) ApplicationOption {
	return func(c *Container) {
		// This will be used during router configuration
		do.OverrideValue(c.injector, enabled)
	}
}

// WithAuthMiddleware enables or disables authentication middleware
func WithAuthMiddleware(enabled bool) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, enabled)
	}
}

// WithI18nMiddleware enables or disables internationalization middleware
func WithI18nMiddleware(enabled bool) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, enabled)
	}
}

// WithTemplateMiddleware enables or disables template middleware
func WithTemplateMiddleware(enabled bool) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, enabled)
	}
}

// WithCustomMiddleware adds custom middleware to the chain
func WithCustomMiddleware(middlewareName string, middlewareFunc func(http.Handler) http.Handler) ApplicationOption {
	return func(c *Container) {
		// Register named custom middleware
		do.ProvideNamedValue(c.injector, "CustomMiddleware_"+middlewareName, middlewareFunc)
	}
}

// WithHealthCheck configures the health check endpoint
func WithHealthCheck(enabled bool, path ...string) ApplicationOption {
	return func(c *Container) {
		config := map[string]interface{}{
			"enabled": enabled,
		}
		if len(path) > 0 {
			config["path"] = path[0]
		} else {
			config["path"] = "/api/health"
		}
		do.OverrideValue(c.injector, config)
	}
}

// WithAPIRoutes configures API routes
func WithAPIRoutes(enabled bool, prefix ...string) ApplicationOption {
	return func(c *Container) {
		config := map[string]interface{}{
			"enabled": enabled,
		}
		if len(prefix) > 0 {
			config["prefix"] = prefix[0]
		} else {
			config["prefix"] = "/api"
		}
		do.OverrideValue(c.injector, config)
	}
}

// WithCustomRoute adds a custom route to the router
func WithCustomRoute(method, path string, handler http.HandlerFunc) ApplicationOption {
	return func(c *Container) {
		route := map[string]interface{}{
			"method":  method,
			"path":    path,
			"handler": handler,
		}
		do.ProvideNamedValue(c.injector, "CustomRoute", route)
	}
}

// WithMiddlewareOrder sets the order of middleware execution
func WithMiddlewareOrder(middlewareOrder ...string) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, middlewareOrder)
	}
}

// WithErrorHandling configures custom error handlers
func WithErrorHandling(notFoundHandler http.Handler, methodNotAllowedHandler http.Handler) ApplicationOption {
	return func(c *Container) {
		errorHandlers := map[string]interface{}{
			"not_found":           notFoundHandler,
			"method_not_allowed":  methodNotAllowedHandler,
		}
		do.OverrideValue(c.injector, errorHandlers)
	}
}
