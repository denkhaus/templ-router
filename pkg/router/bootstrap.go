package router

import (
	"fmt"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
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
	// Middleware configuration
	enableRouterMiddleware    bool
	enableAuthMiddleware      bool
	enableI18nMiddleware      bool
	enableTemplateMiddleware  bool
	middlewareOrder          []string

	// Route configuration
	enableHealthCheck   bool
	healthCheckPath     string
	enableAPIRoutes     bool
	apiRoutePrefix      string
	customRoutes        []map[string]interface{}

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
		enableRouterMiddleware:   true, // Default to enabled
		enableAuthMiddleware:     true,
		enableI18nMiddleware:     true,
		enableTemplateMiddleware: true,
		enableHealthCheck:        true,
		healthCheckPath:          "/api/health",
		enableAPIRoutes:          true,
		apiRoutePrefix:           "/api",
		middlewareOrder:          []string{"router", "auth", "i18n", "template"},
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

// configureMiddleware configures all middleware in the correct order
func (rb *RouterBootstrap) configureMiddleware(mux *chi.Mux) error {
	// Always configure router middleware first if enabled
	if rb.options.enableRouterMiddleware {
		if err := rb.setupRouterMiddleware(mux); err != nil {
			return fmt.Errorf("failed to setup router middleware: %w", err)
		}
	}

	// Configure auth middleware if enabled
	if rb.options.enableAuthMiddleware {
		if err := rb.setupAuthMiddleware(mux); err != nil {
			return fmt.Errorf("failed to setup auth middleware: %w", err)
		}
	}

	// Configure other middleware as needed
	// Note: i18n and template middleware are handled internally by the router

	return nil
}

// setupRouterMiddleware sets up router-level middleware
func (rb *RouterBootstrap) setupRouterMiddleware(mux *chi.Mux) error {
	// Get the middleware setup interface
	middlewareSetup := rb.routerCore.GetMiddlewareSetup()

	// Use type assertion to access the concrete implementation
	if setup, ok := middlewareSetup.(interface{ GetRouterMiddleware() interface{} }); ok {
		routerMiddleware := setup.GetRouterMiddleware()

		// Try to configure the router middleware
		if configurator, ok := routerMiddleware.(interface{ Configure(*chi.Mux) error }); ok {
			return configurator.Configure(mux)
		}
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

// setupAuthMiddleware sets up authentication middleware
func (rb *RouterBootstrap) setupAuthMiddleware(mux *chi.Mux) error {
	authMiddleware, err := middleware.NewAuthContextMiddleware(rb.injector)
	if err != nil {
		return fmt.Errorf("failed to create auth middleware: %w", err)
	}

	mux.Use(authMiddleware.Middleware)
	rb.logger.Info("Applied authentication middleware")
	return nil
}

// registerCustomRoutes registers all custom routes
func (rb *RouterBootstrap) registerCustomRoutes(mux *chi.Mux) error {
	for _, route := range rb.options.customRoutes {
		method, _ := route["method"].(string)
		path, _ := route["path"].(string)
		handler, _ := route["handler"].(http.HandlerFunc)

		if handler == nil {
			rb.logger.Warn("Invalid custom route handler", zap.String("method", method), zap.String("path", path))
			continue
		}

		switch method {
		case "GET":
			mux.Get(path, handler)
		case "POST":
			mux.Post(path, handler)
		case "PUT":
			mux.Put(path, handler)
		case "DELETE":
			mux.Delete(path, handler)
		case "PATCH":
			mux.Patch(path, handler)
		default:
			return fmt.Errorf("unsupported HTTP method: %s", method)
		}
		rb.logger.Info("Registered custom route",
			zap.String("method", method),
			zap.String("path", path))
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

// registerAuthRoutes registers authentication routes
func (rb *RouterBootstrap) registerAuthRoutes(mux *chi.Mux) error {
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