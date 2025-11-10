package router

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// RouterBuilder provides a fluent API for configuring and building a router
type RouterBuilder struct {
	container      interfaces.Container
	injector       do.Injector
	config         interfaces.ConfigService
	logger         *zap.Logger
	routerCore     interfaces.RouterCore
	options        *RouterOptions
	middlewareList []middlewareConfig
	customRoutes   []customRoute
}

// RouterOptions holds configuration for router behavior
type RouterOptions struct {
	// Middleware configuration (all middleware enabled/disabled via env vars)
	// Custom middleware is loaded from DI container in definition order

	// Service overrides
	templateRegistryOverride interfaces.TemplateRegistry
	assetsServiceOverride    interfaces.AssetsService

	// Routing configuration (auth routes controlled by env vars: TR_ROUTER_ENABLE_AUTH_ROUTES, TR_ROUTER_AUTH_ROUTE_PREFIX)
	enableHealthCheck bool
	healthCheckPath   string

	// Error handling
	errorHandler            func(http.ResponseWriter, *http.Request, error)
	notFoundHandler         http.Handler
	methodNotAllowedHandler http.Handler
}

// middlewareConfig represents a middleware configuration with ordering
// TODO: Use this for advanced middleware configuration in future versions
type middlewareConfig struct {
	name       string
	middleware func(http.Handler) http.Handler
	order      int       // Future: Priority order for middleware execution
	priority   int       // Future: Priority for middleware registration
}

// customRoute represents a custom route to register
type customRoute struct {
	method  string
	path    string
	handler http.HandlerFunc
}

// NewRouterBuilder creates a new RouterBuilder with default options
func NewRouterBuilder(container interfaces.Container) (*RouterBuilder, error) {
	injector := container.GetInjector()
	config := do.MustInvoke[interfaces.ConfigService](injector)
	logger := do.MustInvoke[*zap.Logger](injector)

	// Get router core from container
	routerCore := do.MustInvoke[interfaces.RouterCore](injector)

	return &RouterBuilder{
		container:  container,
		injector:   injector,
		config:     config,
		logger:     logger,
		routerCore: routerCore,
		options: &RouterOptions{
			enableHealthCheck: true,
			healthCheckPath:   "/api/health",
		},
		middlewareList: make([]middlewareConfig, 0),
		customRoutes:   make([]customRoute, 0),
	}, nil
}





// WithTemplateRegistry overrides the default template registry
func (rb *RouterBuilder) WithTemplateRegistry(registry interfaces.TemplateRegistry) *RouterBuilder {
	rb.options.templateRegistryOverride = registry
	return rb
}

// WithAssetsService overrides the default assets service
func (rb *RouterBuilder) WithAssetsService(assetsService interfaces.AssetsService) *RouterBuilder {
	rb.options.assetsServiceOverride = assetsService
	return rb
}

// WithHealthCheck configures the health check endpoint
func (rb *RouterBuilder) WithHealthCheck(enabled bool, path ...string) *RouterBuilder {
	rb.options.enableHealthCheck = enabled
	if len(path) > 0 {
		rb.options.healthCheckPath = path[0]
	}
	return rb
}


// WithCustomRoute adds a custom route to the router
func (rb *RouterBuilder) WithCustomRoute(method, path string, handler http.HandlerFunc) *RouterBuilder {
	rb.customRoutes = append(rb.customRoutes, customRoute{
		method:  method,
		path:    path,
		handler: handler,
	})
	return rb
}

// WithErrorHandler sets a custom error handler
func (rb *RouterBuilder) WithErrorHandler(handler func(http.ResponseWriter, *http.Request, error)) *RouterBuilder {
	rb.options.errorHandler = handler
	return rb
}

// WithNotFoundHandler sets a custom 404 handler
func (rb *RouterBuilder) WithNotFoundHandler(handler http.Handler) *RouterBuilder {
	rb.options.notFoundHandler = handler
	return rb
}

// WithMethodNotAllowedHandler sets a custom 405 handler
func (rb *RouterBuilder) WithMethodNotAllowedHandler(handler http.Handler) *RouterBuilder {
	rb.options.methodNotAllowedHandler = handler
	return rb
}

// Build builds and configures the Chi router with all middleware and routes
func (rb *RouterBuilder) Build() (*chi.Mux, error) {
	// Initialize router core first
	if err := rb.routerCore.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize router core: %w", err)
	}

	// Apply service overrides
	if err := rb.applyServiceOverrides(); err != nil {
		return nil, fmt.Errorf("failed to apply service overrides: %w", err)
	}

	// Create Chi router
	mux := chi.NewRouter()

	// Configure middleware in correct order
	if err := rb.configureMiddleware(mux); err != nil {
		return nil, fmt.Errorf("failed to configure middleware: %w", err)
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

	// Apply custom middleware from DI container in definition order
	if err := rb.applyCustomMiddleware(mux); err != nil {
		return nil, fmt.Errorf("failed to apply custom middleware: %w", err)
	}

	// Configure error handlers
	rb.configureErrorHandlers(mux)

	return mux, nil
}

// applyServiceOverrides applies any service overrides to the DI container
func (rb *RouterBuilder) applyServiceOverrides() error {
	if rb.options.templateRegistryOverride != nil {
		do.OverrideValue(rb.injector, rb.options.templateRegistryOverride)
	}
	if rb.options.assetsServiceOverride != nil {
		do.OverrideValue(rb.injector, rb.options.assetsServiceOverride)
	}
	return nil
}

// configureMiddleware configures all middleware in the correct order
func (rb *RouterBuilder) configureMiddleware(mux *chi.Mux) error {
	// Define middleware order (lower number = higher priority)
	middlewareOrder := []struct {
		name    string
		enabled bool
		setup   func() error
		order   int
	}{
		// All middleware always enabled, controlled by environment variables
		{"router", true, rb.setupRouterMiddleware, 100},
		{"i18n", true, rb.setupI18nMiddleware, 300},
		{"template", true, rb.setupTemplateMiddleware, 400},
	}

	// Sort and apply middleware by order
	for _, config := range middlewareOrder {
		if config.enabled {
			if err := config.setup(); err != nil {
				return fmt.Errorf("failed to setup %s middleware: %w", config.name, err)
			}
		}
	}

	return nil
}

// setupRouterMiddleware sets up router-level middleware
func (rb *RouterBuilder) setupRouterMiddleware() error {
	routerMiddleware := rb.routerCore.GetMiddlewareSetup().GetRouterMiddleware()
	return routerMiddleware.Configure(chi.NewRouter()) // Will be applied to main mux
}

// Auth middleware is now handled by the new hook-based AuthValidator interface

// setupI18nMiddleware sets up internationalization middleware
func (rb *RouterBuilder) setupI18nMiddleware() error {
	_ = rb.routerCore.GetMiddlewareSetup()
	// The middleware will be applied in the actual implementation
	return nil
}

// setupTemplateMiddleware sets up template middleware
func (rb *RouterBuilder) setupTemplateMiddleware() error {
	_ = rb.routerCore.GetMiddlewareSetup()
	// The middleware will be applied in the actual implementation
	return nil
}

// registerCustomRoutes registers all custom routes
func (rb *RouterBuilder) registerCustomRoutes(mux *chi.Mux) error {
	// Try to load custom routes from DI container
	customRoutes, err := do.Invoke[[]interfaces.CustomRouteDefinition](rb.injector)
	if err != nil {
		// No custom routes found, use the legacy approach for backward compatibility
		customRoutes = make([]interfaces.CustomRouteDefinition, len(rb.customRoutes))
		for i, route := range rb.customRoutes {
			customRoutes[i] = interfaces.CustomRouteDefinition{
				Method:  route.method,
				Path:    route.path,
				Handler: route.handler,
				Order:   i,
			}
		}
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
func (rb *RouterBuilder) registerHealthCheck(mux *chi.Mux) {
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

// Auth routes are now registered by client applications using AuthHandlers interface

// applyCustomMiddleware loads and applies custom middleware from DI container in definition order
func (rb *RouterBuilder) applyCustomMiddleware(mux *chi.Mux) error {
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
func (rb *RouterBuilder) configureErrorHandlers(mux *chi.Mux) {
	if rb.options.notFoundHandler != nil {
		mux.NotFound(rb.options.notFoundHandler.ServeHTTP)
	}
	if rb.options.methodNotAllowedHandler != nil {
		mux.MethodNotAllowed(rb.options.methodNotAllowedHandler.ServeHTTP)
	}
}

// GetRouterCore returns the underlying router core for advanced usage
func (rb *RouterBuilder) GetRouterCore() interfaces.RouterCore {
	return rb.routerCore
}

// GetLogger returns the logger for debugging
func (rb *RouterBuilder) GetLogger() *zap.Logger {
	return rb.logger
}
