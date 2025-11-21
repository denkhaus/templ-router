package di

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// ApplicationOption defines an option for configuring application services
type ApplicationOption func(c *Container)

// WithTemplateRegistry sets the template registry
func WithTemplateRegistry(registry interfaces.TemplateRegistry) ApplicationOption {
	return func(c *Container) {
		do.OverrideValue(c.injector, registry)
	}
}

// WithAssetsServiceFactory sets the assets service using a factory function
func WithAssetsServiceFactory(factory func(do.Injector) (interfaces.AssetsService, error)) ApplicationOption {
	return func(c *Container) {
		do.Override(c.injector, factory)
	}
}

// WithAuthValidatorFactory sets a custom authentication validator using a factory function
func WithAuthValidatorFactory(factory func(do.Injector) (interfaces.AuthValidator, error)) ApplicationOption {
	return func(c *Container) {
		do.Override(c.injector, factory)
	}
}

// WithCustomMiddleware adds custom middleware to the chain in definition order
func WithCustomMiddleware(middlewareName string, middlewareFunc func(http.Handler) http.Handler) ApplicationOption {
	return func(c *Container) {
		// Get existing middleware definitions or create new slice
		var middlewareDefs []interfaces.CustomMiddlewareDefinition
		if existing, err := do.Invoke[[]interfaces.CustomMiddlewareDefinition](c.injector); err == nil {
			middlewareDefs = existing
		}

		// Calculate next order (append to end)
		nextOrder := len(middlewareDefs)

		// Log middleware registration
		logger := do.MustInvoke[*zap.Logger](c.injector)
		logger.Info("Registering custom middleware",
			zap.String("name", middlewareName),
			zap.Int("order", nextOrder))

		// Add new middleware with definition order
		middlewareDefs = append(middlewareDefs, interfaces.CustomMiddlewareDefinition{
			Name:  middlewareName,
			Func:  middlewareFunc,
			Order: nextOrder,
		})

		// Override the middleware definitions slice
		do.OverrideValue(c.injector, middlewareDefs)

		logger.Info("Custom middleware registered successfully",
			zap.Int("total_middleware", len(middlewareDefs)))
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

// WithCustomRoute adds a custom route to the router
func WithCustomRoute(method, path string, handler http.HandlerFunc) ApplicationOption {
	return func(c *Container) {
		// Get existing custom routes or create new slice
		var customRoutes []interfaces.CustomRouteDefinition
		if existing, err := do.Invoke[[]interfaces.CustomRouteDefinition](c.injector); err == nil {
			customRoutes = existing
		}

		// Calculate next order (append to end)
		nextOrder := len(customRoutes)

		// Add new route with definition order
		customRoutes = append(customRoutes, interfaces.CustomRouteDefinition{
			Method:  method,
			Path:    path,
			Handler: handler,
			Order:   nextOrder,
		})

		// Override the custom routes slice
		do.OverrideValue(c.injector, customRoutes)
	}
}

// WithErrorHandling configures custom error handlers
func WithErrorHandling(notFoundHandler http.Handler, methodNotAllowedHandler http.Handler) ApplicationOption {
	return func(c *Container) {
		errorHandlers := map[string]interface{}{
			"not_found":          notFoundHandler,
			"method_not_allowed": methodNotAllowedHandler,
		}
		do.OverrideValue(c.injector, errorHandlers)
	}
}
