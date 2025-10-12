package router

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/pipeline"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// cleanRouterCore implements clean architecture principles with proper separation of concerns (private implementation)
// REFACTORED FROM: clean_router.go (301 lines → 4 separate files)
type cleanRouterCore struct {
	// Core configuration
	scanPath string
	config   interfaces.ConfigService
	logger   *zap.Logger
	injector do.Injector // Store injector for proper DI

	// Separated components (Separation of Concerns)
	routeRegistrar  RouteRegistrar
	handlerBuilder  HandlerBuilder
	middlewareSetup MiddlewareSetup

	// Handler pipeline
	handlerPipeline *pipeline.HandlerPipeline

	// Route discovery and processing
	routeDiscovery RouteDiscovery
	configLoader   ConfigLoader

	// Data storage (clean, no business logic)
	routes          []Route
	layoutTemplates []LayoutTemplate
	errorTemplates  []ErrorTemplate
}

// NewCleanRouterCore creates a new clean router with separated concerns for DI
func NewCleanRouterCore(i do.Injector) (RouterCore, error) {
	// Inject core dependencies
	config := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	handlerPipeline := do.MustInvoke[*pipeline.HandlerPipeline](i)
	routeDiscovery := do.MustInvoke[RouteDiscovery](i)
	configLoader := do.MustInvoke[ConfigLoader](i)

	// Create separated components
	handlerBuilder, err := NewHandlerBuilder(i)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler builder: %w", err)
	}

	middlewareSetup, err := NewMiddlewareSetup(i)
	if err != nil {
		return nil, fmt.Errorf("failed to create middleware setup: %w", err)
	}

	// Validate middleware setup
	if err := middlewareSetup.ValidateMiddlewareSetup(); err != nil {
		return nil, fmt.Errorf("middleware setup validation failed: %w", err)
	}

	return &cleanRouterCore{
		scanPath:        config.GetTemplateScanPath(), // Use config interface method
		config:          config,
		logger:          logger,
		injector:        i, // Store injector for RouteRegistrar creation
		handlerBuilder:  handlerBuilder,
		middlewareSetup: middlewareSetup,
		handlerPipeline: handlerPipeline,
		routeDiscovery:  routeDiscovery,
		configLoader:    configLoader,
	}, nil
}

// Initialize discovers and processes all routes, layouts, and error templates
func (crc *cleanRouterCore) Initialize() error {
	crc.logger.Info("Initializing clean router core", zap.String("scan_path", crc.scanPath))

	// Discover routes
	routes, err := crc.routeDiscovery.DiscoverRoutes(crc.scanPath)
	if err != nil {
		return fmt.Errorf("failed to discover routes: %w", err)
	}
	crc.routes = routes

	// Discover layouts
	layouts, err := crc.routeDiscovery.DiscoverLayouts(crc.scanPath)
	if err != nil {
		return fmt.Errorf("failed to discover layouts: %w", err)
	}
	crc.layoutTemplates = layouts

	// Discover error templates
	errorTemplates, err := crc.routeDiscovery.DiscoverErrorTemplates(crc.scanPath)
	if err != nil {
		return fmt.Errorf("failed to discover error templates: %w", err)
	}
	crc.errorTemplates = errorTemplates

	crc.logger.Info("Clean router core initialized successfully",
		zap.Int("routes", len(crc.routes)),
		zap.Int("layouts", len(crc.layoutTemplates)),
		zap.Int("error_templates", len(crc.errorTemplates)))

	return nil
}

// RegisterRoutes registers all discovered routes with a Chi router
func (crc *cleanRouterCore) RegisterRoutes(chiRouter *chi.Mux) error {
	crc.logger.Info("Registering routes with Chi router")

	// Create route registrar through DI to ensure proper ConfigService injection
	routeRegistrar, err := NewRouteRegistrar(crc.injector, chiRouter)
	if err != nil {
		return fmt.Errorf("failed to create route registrar: %w", err)
	}
	crc.routeRegistrar = routeRegistrar

	// Convert routes to interfaces.Route format
	interfaceRoutes := crc.convertToInterfaceRoutes(crc.routes)

	// Register all routes
	if err := crc.routeRegistrar.RegisterRoutes(interfaceRoutes); err != nil {
		return fmt.Errorf("failed to register routes: %w", err)
	}

	// Register static routes
	crc.routeRegistrar.RegisterStaticRoutes()

	// Register error handlers
	crc.routeRegistrar.Register404Handler()
	crc.routeRegistrar.RegisterMethodNotAllowedHandler()

	crc.logger.Info("All routes registered successfully",
		zap.Int("total_routes", len(crc.routes)))

	return nil
}

// convertToInterfaceRoutes converts router.Route to interfaces.Route
func (crc *cleanRouterCore) convertToInterfaceRoutes(routes []Route) []interfaces.Route {
	interfaceRoutes := make([]interfaces.Route, len(routes))

	for i, route := range routes {
		interfaceRoutes[i] = interfaces.Route{
			Path:         route.Path,
			TemplateFile: route.TemplateFile,
			IsDynamic:    route.IsDynamic,
		}
	}

	return interfaceRoutes
}

// GetRoutes returns all discovered routes
func (crc *cleanRouterCore) GetRoutes() []Route {
	return crc.routes
}

// GetLayoutTemplates returns all discovered layout templates
func (crc *cleanRouterCore) GetLayoutTemplates() []LayoutTemplate {
	return crc.layoutTemplates
}

// GetErrorTemplates returns all discovered error templates
func (crc *cleanRouterCore) GetErrorTemplates() []ErrorTemplate {
	return crc.errorTemplates
}

// GetMiddlewareSetup returns the middleware setup for external access
func (crc *cleanRouterCore) GetMiddlewareSetup() MiddlewareSetup {
	return crc.middlewareSetup
}

// GetHandlerBuilder returns the handler builder for external access
func (crc *cleanRouterCore) GetHandlerBuilder() HandlerBuilder {
	return crc.handlerBuilder
}

// GetRouteRegistrar returns the route registrar for external access
func (crc *cleanRouterCore) GetRouteRegistrar() RouteRegistrar {
	return crc.routeRegistrar
}
