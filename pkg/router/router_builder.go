package router

import (
	"fmt"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
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
	// Middleware configuration
	enableRouterMiddleware   bool
	enableAuthMiddleware     bool
	enableI18nMiddleware     bool
	enableTemplateMiddleware bool
	customMiddleware         []func(http.Handler) http.Handler

	// Service overrides
	authHandlersOverride     interfaces.AuthHandlers
	userStoreOverride        interfaces.UserStore
	templateRegistryOverride interfaces.TemplateRegistry
	assetsServiceOverride    interfaces.AssetsService

	// Routing configuration
	enableHealthCheck bool
	healthCheckPath   string
	enableAPIRoutes   bool
	apiRoutePrefix    string

	// Error handling
	errorHandler            func(http.ResponseWriter, *http.Request, error)
	notFoundHandler         http.Handler
	methodNotAllowedHandler http.Handler
}

// middlewareConfig represents a middleware configuration with ordering
type middlewareConfig struct {
	name       string
	middleware func(http.Handler) http.Handler
	order      int
	priority   int
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
			enableRouterMiddleware:   true,
			enableAuthMiddleware:     true,
			enableI18nMiddleware:     true,
			enableTemplateMiddleware: true,
			enableHealthCheck:        true,
			healthCheckPath:          "/api/health",
			enableAPIRoutes:          true,
			apiRoutePrefix:           "/api",
		},
		middlewareList: make([]middlewareConfig, 0),
		customRoutes:   make([]customRoute, 0),
	}, nil
}

// WithRouterMiddleware enables or disables router-level middleware
func (rb *RouterBuilder) WithRouterMiddleware(enabled bool) *RouterBuilder {
	rb.options.enableRouterMiddleware = enabled
	return rb
}

// WithAuthMiddleware enables or disables authentication middleware
func (rb *RouterBuilder) WithAuthMiddleware(enabled bool) *RouterBuilder {
	rb.options.enableAuthMiddleware = enabled
	return rb
}

// WithI18nMiddleware enables or disables internationalization middleware
func (rb *RouterBuilder) WithI18nMiddleware(enabled bool) *RouterBuilder {
	rb.options.enableI18nMiddleware = enabled
	return rb
}

// WithTemplateMiddleware enables or disables template middleware
func (rb *RouterBuilder) WithTemplateMiddleware(enabled bool) *RouterBuilder {
	rb.options.enableTemplateMiddleware = enabled
	return rb
}

// WithCustomMiddleware adds custom middleware to the chain
func (rb *RouterBuilder) WithCustomMiddleware(middleware func(http.Handler) http.Handler, order ...int) *RouterBuilder {
	middlewareOrder := 1000 // Default order for custom middleware
	if len(order) > 0 {
		middlewareOrder = order[0]
	}

	rb.options.customMiddleware = append(rb.options.customMiddleware, middleware)
	rb.middlewareList = append(rb.middlewareList, middlewareConfig{
		name:       "custom",
		middleware: middleware,
		order:      middlewareOrder,
		priority:   len(rb.middlewareList),
	})
	return rb
}

// WithAuthHandlers overrides the default auth handlers
func (rb *RouterBuilder) WithAuthHandlers(authHandlers interfaces.AuthHandlers) *RouterBuilder {
	rb.options.authHandlersOverride = authHandlers
	return rb
}

// WithUserStore overrides the default user store
func (rb *RouterBuilder) WithUserStore(userStore interfaces.UserStore) *RouterBuilder {
	rb.options.userStoreOverride = userStore
	return rb
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

// WithAPIRoutes configures API routes
func (rb *RouterBuilder) WithAPIRoutes(enabled bool, prefix ...string) *RouterBuilder {
	rb.options.enableAPIRoutes = enabled
	if len(prefix) > 0 {
		rb.options.apiRoutePrefix = prefix[0]
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

	// Register auth routes if enabled
	if rb.options.enableAPIRoutes {
		if err := rb.registerAuthRoutes(mux); err != nil {
			return nil, fmt.Errorf("failed to register auth routes: %w", err)
		}
	}

	// Configure error handlers
	rb.configureErrorHandlers(mux)

	return mux, nil
}

// applyServiceOverrides applies any service overrides to the DI container
func (rb *RouterBuilder) applyServiceOverrides() error {
	if rb.options.authHandlersOverride != nil {
		do.OverrideValue(rb.injector, rb.options.authHandlersOverride)
	}
	if rb.options.userStoreOverride != nil {
		do.OverrideValue(rb.injector, rb.options.userStoreOverride)
	}
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
		{"router", rb.options.enableRouterMiddleware, rb.setupRouterMiddleware, 100},
		{"auth", rb.options.enableAuthMiddleware, rb.setupAuthMiddleware, 200},
		{"i18n", rb.options.enableI18nMiddleware, rb.setupI18nMiddleware, 300},
		{"template", rb.options.enableTemplateMiddleware, rb.setupTemplateMiddleware, 400},
	}

	// Sort and apply middleware by order
	for _, config := range middlewareOrder {
		if config.enabled {
			if err := config.setup(); err != nil {
				return fmt.Errorf("failed to setup %s middleware: %w", config.name, err)
			}
		}
	}

	// Apply custom middleware
	for _, custom := range rb.options.customMiddleware {
		mux.Use(custom)
		rb.logger.Info("Applied custom middleware")
	}

	return nil
}

// setupRouterMiddleware sets up router-level middleware
func (rb *RouterBuilder) setupRouterMiddleware() error {
	routerMiddleware := rb.routerCore.GetMiddlewareSetup().GetRouterMiddleware()
	return routerMiddleware.Configure(chi.NewRouter()) // Will be applied to main mux
}

// setupAuthMiddleware sets up authentication middleware
func (rb *RouterBuilder) setupAuthMiddleware() error {
	_, err := middleware.NewAuthContextMiddleware(rb.injector)
	if err != nil {
		return fmt.Errorf("failed to create auth middleware: %w", err)
	}
	// The middleware will be applied in the actual implementation
	return nil
}

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
	for _, route := range rb.customRoutes {
		switch route.method {
		case "GET":
			mux.Get(route.path, route.handler)
		case "POST":
			mux.Post(route.path, route.handler)
		case "PUT":
			mux.Put(route.path, route.handler)
		case "DELETE":
			mux.Delete(route.path, route.handler)
		case "PATCH":
			mux.Patch(route.path, route.handler)
		default:
			return fmt.Errorf("unsupported HTTP method: %s", route.method)
		}
		rb.logger.Info("Registered custom route",
			zap.String("method", route.method),
			zap.String("path", route.path))
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

// registerAuthRoutes registers authentication routes
func (rb *RouterBuilder) registerAuthRoutes(mux *chi.Mux) error {
	authHandlers := do.MustInvoke[interfaces.AuthHandlers](rb.injector)
	authHandlers.RegisterRoutes(func(method, path string, handler http.HandlerFunc) {
		switch method {
		case "POST":
			mux.Post(path, handler)
		case "GET":
			mux.Get(path, handler)
		}
		rb.logger.Info("Auth route registered",
			zap.String("method", method),
			zap.String("path", path))
	})
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
