package router

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// RouterBootstrap provides a streamlined way to bootstrap a router with all middleware and routes
type RouterBootstrap struct {
	injector   do.Injector
	config     interfaces.ConfigService
	logger     *zap.Logger
	routerCore interfaces.RouterCore
	options    *BootstrapConfig
}

// BootstrapConfig holds configuration for router bootstrap
type BootstrapConfig struct {
	// Middleware configuration (all middleware enabled/disabled via env vars)
	middlewareOrder []string

	// Route configuration (auth routes controlled by env vars: TR_ROUTER_ENABLE_AUTH_ROUTES, TR_ROUTER_AUTH_ROUTE_PREFIX)
	enableHealthCheck bool
	healthCheckPath   string
	customRoutes      []interfaces.CustomRouteDefinition

	// Error handling
	errorHandlers map[string]interface{}
}

// NewRouterBootstrap creates a new router bootstrap service
func NewRouterBootstrap(injector do.Injector) (*RouterBootstrap, error) {
	config := do.MustInvoke[interfaces.ConfigService](injector)
	logger := do.MustInvoke[*zap.Logger](injector)
	routerCore := do.MustInvoke[interfaces.RouterCore](injector)

	// Load configuration from DI container or use defaults
	bootstrapConfig := &BootstrapConfig{
		enableHealthCheck: true,
		healthCheckPath:   "/api/health",
		middlewareOrder:   []string{"router", "auth", "i18n", "template"},
	}

	// Try to load overridden configuration from DI container
	if err := bootstrapConfig.loadFromDI(injector); err != nil {
		logger.Info("No bootstrap configuration found in DI container, using defaults")
	}

	return &RouterBootstrap{
		injector:   injector,
		config:     config,
		logger:     logger,
		routerCore: routerCore,
		options:    bootstrapConfig,
	}, nil
}

// loadFromDI attempts to load configuration from the DI container
func (bc *BootstrapConfig) loadFromDI(injector do.Injector) error {
	// This would look for configuration values in the DI container
	// For now, we'll use the defaults
	return nil
}

// Bootstrap initializes and configures the router with all middleware and routes
func (rb *RouterBootstrap) Bootstrap() (*chi.Mux, error) {
	// Initialize router core first
	if err := rb.routerCore.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize router core: %w", err)
	}

	// Create Chi router
	mux := chi.NewRouter()

	// Configure middleware in correct order
	if err := rb.configureMiddleware(mux); err != nil {
		return nil, fmt.Errorf("failed to configure middleware: %w", err)
	}

	// Apply custom middleware from DI container in definition order (BEFORE routes)
	if err := rb.applyCustomMiddleware(mux); err != nil {
		return nil, fmt.Errorf("failed to apply custom middleware: %w", err)
	}

	// Register custom routes first (highest priority)
	if err := rb.registerCustomRoutes(mux); err != nil {
		return nil, fmt.Errorf("failed to register custom routes: %w", err)
	}

	// Register health check if enabled
	if rb.options.enableHealthCheck {
		rb.registerHealthCheck(mux)
	}

	// Register file-based routes
	if err := rb.routerCore.RegisterRoutes(mux); err != nil {
		return nil, fmt.Errorf("failed to register routes: %w", err)
	}

	// Auth routes are now registered by client applications

	// Configure error handlers
	rb.configureErrorHandlers(mux)

	return mux, nil
}

// configureMiddleware configures all middleware in the correct order
func (rb *RouterBootstrap) configureMiddleware(mux *chi.Mux) error {
	// Always configure router middleware (controlled by env vars: TR_ROUTER_ENABLE_*)
	if err := rb.setupRouterMiddleware(mux); err != nil {
		return fmt.Errorf("failed to setup router middleware: %w", err)
	}

	// Auth middleware is now handled by the new hook-based AuthValidator interface

	// Note: i18n and template middleware are handled internally by the router core
	// based on environment variables: TR_I18N_ENABLE_*, TR_TEMPLATE_ENABLE_*

	return nil
}

// setupRouterMiddleware sets up router-level middleware
func (rb *RouterBootstrap) setupRouterMiddleware(mux *chi.Mux) error {
	// Get the middleware setup interface
	middlewareSetup := rb.routerCore.GetMiddlewareSetup()

	// Use type assertion to access the concrete implementation
	if setup, ok := middlewareSetup.(interface {
		GetRouterMiddleware() interfaces.RouterMiddlewareInterface
	}); ok {
		routerMiddleware := setup.GetRouterMiddleware()

		// Configure the router middleware
		return routerMiddleware.Configure(mux)
	}

	// Fallback: apply basic router middleware directly
	if rb.config.GetRouterEnableTrailingSlash() {
		mux.Use(chimiddleware.RedirectSlashes)
		rb.logger.Info("Enabled trailing slash redirection")
	}

	if rb.config.GetRouterEnableSlashRedirect() {
		mux.Use(chimiddleware.CleanPath)
		rb.logger.Info("Enabled slash redirection")
	}

	if rb.config.GetRouterEnableMethodNotAllowed() {
		// Method not allowed is handled by chi by default
		rb.logger.Info("Method not allowed handling enabled")
	}

	return nil
}

// Auth middleware is now handled by the new hook-based AuthValidator interface

// registerCustomRoutes registers all custom routes
func (rb *RouterBootstrap) registerCustomRoutes(mux *chi.Mux) error {
	// Try to load custom routes from DI container
	customRoutes, err := do.Invoke[[]interfaces.CustomRouteDefinition](rb.injector)
	if err != nil {
		// No custom routes found in DI container, use bootstrap config
		customRoutes = rb.options.customRoutes
	}

	if len(customRoutes) == 0 {
		return nil
	}

	// Sort routes by definition order to ensure consistent registration
	sort.Slice(customRoutes, func(i, j int) bool {
		return customRoutes[i].Order < customRoutes[j].Order
	})

	// Register routes in definition order
	for _, route := range customRoutes {
		if route.Handler == nil {
			rb.logger.Warn("Invalid custom route handler",
				zap.String("method", route.Method),
				zap.String("path", route.Path))
			continue
		}

		switch route.Method {
		case "GET":
			mux.Get(route.Path, route.Handler)
		case "POST":
			mux.Post(route.Path, route.Handler)
		case "PUT":
			mux.Put(route.Path, route.Handler)
		case "DELETE":
			mux.Delete(route.Path, route.Handler)
		case "PATCH":
			mux.Patch(route.Path, route.Handler)
		default:
			return fmt.Errorf("unsupported HTTP method: %s", route.Method)
		}
		rb.logger.Info("Registered custom route",
			zap.String("method", route.Method),
			zap.String("path", route.Path),
			zap.Int("order", route.Order))
	}
	return nil
}

// registerHealthCheck registers the health check endpoint
func (rb *RouterBootstrap) registerHealthCheck(mux *chi.Mux) {
	mux.HandleFunc(rb.options.healthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"status": "healthy",
			"architecture": "clean",
			"dependency_injection": "samber/do",
			"router": "multi-language file-based",
			"i18n": "decentralized",
			"languages": ["en", "de"]
		}`))
	})
	rb.logger.Info("Registered health check", zap.String("path", rb.options.healthCheckPath))
}

// applyCustomMiddleware loads and applies custom middleware from DI container in definition order
func (rb *RouterBootstrap) applyCustomMiddleware(mux *chi.Mux) error {
	// Try to load custom middleware definitions from DI container
	middlewareDefs, err := do.Invoke[[]interfaces.CustomMiddlewareDefinition](rb.injector)
	if err != nil {
		// No custom middleware found, which is fine
		rb.logger.Info("No custom middleware found in DI container")
		return nil
	}

	if len(middlewareDefs) == 0 {
		rb.logger.Info("No custom middleware to apply")
		return nil
	}

	// Sort middleware by definition order to ensure correct execution sequence
	sort.Slice(middlewareDefs, func(i, j int) bool {
		return middlewareDefs[i].Order < middlewareDefs[j].Order
	})

	// Apply middleware in definition order
	for _, middlewareDef := range middlewareDefs {
		mux.Use(middlewareDef.Func)
		rb.logger.Info("Applied custom middleware",
			zap.String("name", middlewareDef.Name),
			zap.Int("order", middlewareDef.Order))
	}

	return nil
}

// configureErrorHandlers configures custom error handlers if provided
func (rb *RouterBootstrap) configureErrorHandlers(mux *chi.Mux) {
	if rb.options.errorHandlers != nil {
		if notFoundHandler, ok := rb.options.errorHandlers["not_found"].(http.Handler); ok {
			mux.NotFound(notFoundHandler.ServeHTTP)
		}
		if methodNotAllowedHandler, ok := rb.options.errorHandlers["method_not_allowed"].(http.Handler); ok {
			mux.MethodNotAllowed(methodNotAllowedHandler.ServeHTTP)
		}
	}
}

// GetRouterCore returns the underlying router core for advanced usage
func (rb *RouterBootstrap) GetRouterCore() interfaces.RouterCore {
	return rb.routerCore
}

// GetLogger returns the logger for debugging
func (rb *RouterBootstrap) GetLogger() *zap.Logger {
	return rb.logger
}
