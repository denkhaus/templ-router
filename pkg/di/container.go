package di

import (
	"github.com/denkhaus/templ-router/pkg/config"
	"github.com/denkhaus/templ-router/pkg/router"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/denkhaus/templ-router/pkg/router/pipeline"
	"github.com/denkhaus/templ-router/pkg/router/services"
	"github.com/denkhaus/templ-router/pkg/services/auth"
	"github.com/denkhaus/templ-router/pkg/services/logger"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// Container is the new DI container for library usage without generated templates
type Container struct {
	injector do.Injector
}

// NewContainer creates a new library-compatible DI container
func NewContainer() *Container {
	injector := do.New()
	return &Container{
		injector: injector,
	}
}

// GetInjector returns the underlying injector
func (c *Container) GetInjector() do.Injector {
	return c.injector
}

// RegisterApplicationServices registers all services the application provides to override internal behaviour
// Uses options pattern for flexible configuration
func (c *Container) RegisterApplicationServices(options ...ApplicationOption) {
	// Apply all options
	for _, option := range options {
		option(c)
	}
}

// RegisterRouterServices registers all router services (without template dependencies)
func (c *Container) RegisterRouterServices() {
	// Register logger
	do.Provide(c.injector, logger.NewService)

	// Register configuration
	// this should not be exposed. use config.NewConfigService instead
	do.Provide(c.injector, config.NewConfigService)

	// Register stores - constructors already return interfaces! (pluggable)
	// Option 1: In-Memory Stores (simple, for development)
	do.Provide(c.injector, services.NewInMemorySessionStore)
	do.Provide(c.injector, services.NewInMemoryUserStore)

	// Option 2: Productive Stores (advanced, for production) - commented out
	// do.Provide(c.injector, auth.NewProductiveSessionStore)
	// do.Provide(c.injector, auth.NewProductiveUserStore)

	// Internal services (these can remain concrete for now)
	do.Provide(c.injector, services.NewInMemoryTranslationStore)

	// Register clean services - constructors already return interfaces!
	do.Provide(c.injector, services.NewAuthService)
	do.Provide(c.injector, services.NewI18nService)

	// UNIFIED TEMPLATE ARCHITECTURE - Performance Optimized
	// Note: This will use the externally registered TemplateRegistry
	do.Provide(c.injector, services.NewOptimizedTemplateService)

	// UNIFIED VALIDATION ARCHITECTURE - Orchestrated Validation Logic
	do.Provide(c.injector, services.NewValidationOrchestrator)

	// Register middleware services - constructors already return interfaces!
	do.Provide(c.injector, middleware.NewProductiveFileSystemChecker)
	do.Provide(c.injector, middleware.NewLayoutService)
	do.Provide(c.injector, middleware.NewErrorServiceCore)
	do.Provide(c.injector, middleware.NewConfigurableParameterExtractor)

	do.Provide(c.injector, middleware.NewAuthMiddleware)
	do.Provide(c.injector, middleware.NewI18nMiddleware)
	do.Provide(c.injector, middleware.NewTemplateMiddleware)

	do.Provide(c.injector, pipeline.NewHandlerPipeline)
	do.Provide(c.injector, services.NewRouteDiscovery)
	do.Provide(c.injector, services.NewConfigLoader)

	// Register auth services (external, pluggable)
	do.Provide(c.injector, auth.NewAuthProvider)
	do.Provide(c.injector, auth.NewAuthHandlers)

	// Register clean router (refactored with separation of concerns)
	do.Provide(c.injector, router.NewCleanRouterCore)

}

// GetRouter returns the clean router from the container
func (c *Container) GetRouter() router.RouterCore {
	return do.MustInvoke[router.RouterCore](c.injector)
}

// GetLogger returns the logger from the container
func (c *Container) GetLogger() *zap.Logger {
	return do.MustInvoke[*zap.Logger](c.injector)
}

// Shutdown gracefully shuts down all services
func (c *Container) Shutdown() error {
	return c.injector.Shutdown()
}
