package router

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/pipeline"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// routerCore implements clean architecture principles with proper separation of concerns (private implementation)
type routerCore struct {
	// Core configuration
	scanPath      string
	config        interfaces.ConfigService
	assetsService interfaces.AssetsService
	logger        *zap.Logger
	injector      do.Injector

	routeRegistrar  interfaces.RouteRegistrar
	handlerBuilder  interfaces.HandlerBuilder
	middlewareSetup interfaces.MiddlewareSetup

	// Handler pipeline
	handlerPipeline *pipeline.HandlerPipeline

	// Route discovery and processing
	routeDiscovery interfaces.RouteDiscovery
	configLoader   interfaces.ConfigLoader

	// Data storage (clean, no business logic)
	routes          []interfaces.Route
	layoutTemplates []interfaces.LayoutTemplate
	errorTemplates  []interfaces.ErrorTemplate
}

// NewRouterCore creates a new clean router with separated concerns for DI
func NewRouterCore(i do.Injector) (interfaces.RouterCore, error) {
	// Inject core dependencies
	config := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	handlerPipeline := do.MustInvoke[*pipeline.HandlerPipeline](i)
	routeDiscovery := do.MustInvoke[interfaces.RouteDiscovery](i)
	assetsService := do.MustInvoke[interfaces.AssetsService](i)
	configLoader := do.MustInvoke[interfaces.ConfigLoader](i)

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

	return &routerCore{
		scanPath:        config.GetLayoutRootDirectory(),
		config:          config,
		assetsService:   assetsService,
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
func (crc *routerCore) Initialize() error {
	crc.logger.Info("Initializing clean router core", zap.String("scan_path", crc.scanPath))

	// Discover routes
	routes, err := crc.routeDiscovery.DiscoverRoutes(crc.scanPath)
	if err != nil {
		return fmt.Errorf("failed to discover routes: %w", err)
	}
	crc.routes = routes

	// Load translations for all discovered templates
	if err := crc.loadTranslationsForDiscoveredRoutes(); err != nil {
		crc.logger.Warn("Failed to load some translations during initialization", zap.Error(err))
		// Don't fail initialization for translation loading errors
	}

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
func (crc *routerCore) RegisterRoutes(chiRouter *chi.Mux) error {
	crc.logger.Info("Registering routes with Chi router")

	// Note: Router middleware should be configured BEFORE calling RegisterRoutes
	// This is now handled in the application layer to ensure proper middleware order

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

	// Authentication handlers are now registered by client applications

	// Register error handlers
	crc.routeRegistrar.Register404Handler()
	crc.routeRegistrar.RegisterMethodNotAllowedHandler()

	crc.logger.Info("All routes registered successfully",
		zap.Int("total_routes", len(crc.routes)))

	return nil
}

// convertToInterfaceRoutes converts router.Route to interfaces.Route
func (crc *routerCore) convertToInterfaceRoutes(routes []interfaces.Route) []interfaces.Route {
	interfaceRoutes := make([]interfaces.Route, len(routes))

	for i, route := range routes {
		interfaceRoutes[i] = interfaces.Route{
			Path:                 route.Path,
			TemplateFile:         route.TemplateFile,
			IsDynamic:            route.IsDynamic,
			RequiresDataService:  route.RequiresDataService,
			DataServiceInterface: route.DataServiceInterface,
		}
	}

	return interfaceRoutes
}

// GetRoutes returns all discovered routes
func (crc *routerCore) GetRoutes() []interfaces.Route {
	return crc.routes
}

// GetLayoutTemplates returns all discovered layout templates
func (crc *routerCore) GetLayoutTemplates() []interfaces.LayoutTemplate {
	return crc.layoutTemplates
}

// GetErrorTemplates returns all discovered error templates
func (crc *routerCore) GetErrorTemplates() []interfaces.ErrorTemplate {
	return crc.errorTemplates
}

// GetMiddlewareSetup returns the middleware setup for external access
func (crc *routerCore) GetMiddlewareSetup() interfaces.MiddlewareSetup {
	return crc.middlewareSetup
}

// GetHandlerBuilder returns the handler builder for external access
func (crc *routerCore) GetHandlerBuilder() interfaces.HandlerBuilder {
	return crc.handlerBuilder
}

// GetRouteRegistrar returns the route registrar for external access
func (crc *routerCore) GetRouteRegistrar() interfaces.RouteRegistrar {
	return crc.routeRegistrar
}

// loadTranslationsForDiscoveredRoutes loads translations for all discovered routes
func (crc *routerCore) loadTranslationsForDiscoveredRoutes() error {
	crc.logger.Info("Loading translations for discovered routes", zap.Int("route_count", len(crc.routes)))

	// Extract template paths from discovered routes
	templatePaths := make([]string, 0, len(crc.routes))
	for _, route := range crc.routes {
		if route.TemplateFile != "" {
			templatePaths = append(templatePaths, route.TemplateFile)
		}
	}

	// Get translation store from DI container
	translationStore := do.MustInvoke[interfaces.I18nService](crc.injector)

	// Load all translations in bulk
	if err := translationStore.LoadAllTranslations(templatePaths); err != nil {
		return fmt.Errorf("failed to bulk load translations: %w", err)
	}

	crc.logger.Info("Successfully loaded translations for all discovered routes",
		zap.Int("template_count", len(templatePaths)))

	return nil
}
