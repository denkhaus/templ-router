package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// RouteRegistrar defines the contract for route registration
type RouteRegistrar interface {
	RegisterRoutes(routes []interfaces.Route) error
	RegisterStaticRoutes()
	Register404Handler()
	RegisterMethodNotAllowedHandler()
}

// routeRegistrar handles route registration logic (private implementation)
// SEPARATED FROM: clean_router.go (Separation of Concerns)
type routeRegistrar struct {
	router          *chi.Mux
	handlerBuilder  HandlerBuilder
	middlewareSetup MiddlewareSetup
	configService   interfaces.ConfigService
	logger          *zap.Logger
}

// NewRouteRegistrar creates a new route registrar
func NewRouteRegistrar(i do.Injector, router *chi.Mux) (RouteRegistrar, error) {
	handlerBuilder, err := NewHandlerBuilder(i)
	if err != nil {
		return nil, err
	}

	middlewareSetup, err := NewMiddlewareSetup(i)
	if err != nil {
		return nil, err
	}

	configService := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &routeRegistrar{
		router:          router,
		handlerBuilder:  handlerBuilder,
		middlewareSetup: middlewareSetup,
		configService:   configService,
		logger:          logger,
	}, nil
}

// RegisterRoutes registers all discovered routes with the router
func (rr *routeRegistrar) RegisterRoutes(routes []interfaces.Route) error {
	rr.logger.Info("Registering routes", zap.Int("count", len(routes)))

	for _, route := range routes {
		if err := rr.registerSingleRoute(route); err != nil {
			rr.logger.Error("Failed to register route",
				zap.String("path", route.Path),
				zap.Error(err))
			return err
		}
	}

	return nil
}

// registerSingleRoute registers a single route with proper handler and middleware
func (rr *routeRegistrar) registerSingleRoute(route interfaces.Route) error {
	// LOCALE EXPANSION: Handle $locale routes specially
	if strings.Contains(route.Path, "$locale") {
		return rr.registerLocaleSpecificRoutes(route)
	}

	// Convert route pattern for Chi router (replace $id with Chi syntax)
	chiPattern := rr.convertRoutePattern(route.Path)

	// Build handler with middleware pipeline
	handler := rr.handlerBuilder.BuildHandler(route)

	// Register with Chi router
	rr.router.Get(chiPattern, handler.ServeHTTP)

	rr.logger.Debug("Route registered",
		zap.String("original_pattern", route.Path),
		zap.String("chi_pattern", chiPattern),
		zap.String("template", route.TemplateFile),
		zap.Bool("dynamic", route.IsDynamic))

	return nil
}

// convertRoutePattern converts router patterns to Chi router syntax
func (rr *routeRegistrar) convertRoutePattern(pattern string) string {
	// Convert $locale to {locale} for Chi
	chiPattern := strings.ReplaceAll(pattern, "$locale", "{locale}")

	// Convert $id to {id} for Chi
	chiPattern = strings.ReplaceAll(chiPattern, "$id", "{id}")

	// Handle other dynamic parameters if needed
	// Add more conversions as the router grows

	return chiPattern
}

// RegisterStaticRoutes registers static file serving routes
func (rr *routeRegistrar) RegisterStaticRoutes() {
	// Serve static assets
	rr.router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	rr.logger.Debug("Static routes registered", zap.String("path", "/assets/*"))
}

// Register404Handler registers the 404 not found handler
func (rr *routeRegistrar) Register404Handler() {
	rr.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rr.logger.Info("404 handler triggered",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))

		// Use error service to render 404 page
		errorService := rr.middlewareSetup.GetErrorService()
		if errorService != nil {
			component := errorService.CreateErrorComponent("The requested page could not be found.", r.URL.Path)
			if component != nil {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusNotFound)
				if err := component.Render(r.Context(), w); err != nil {
					rr.logger.Error("Failed to render 404 page", zap.Error(err))
					http.Error(w, "Page not found", http.StatusNotFound)
				}
				return
			}
		}

		// Fallback to simple 404
		http.Error(w, "Page not found", http.StatusNotFound)
	})
}

// RegisterMethodNotAllowedHandler registers the 405 method not allowed handler
func (rr *routeRegistrar) RegisterMethodNotAllowedHandler() {
	rr.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rr.logger.Info("405 handler triggered",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))

		errorService := rr.middlewareSetup.GetErrorService()
		if errorService != nil {
			component := errorService.CreateErrorComponent(
				fmt.Sprintf("Method %s not allowed for this path.", r.Method),
				r.URL.Path)
			if component != nil {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusMethodNotAllowed)
				if err := component.Render(r.Context(), w); err != nil {
					rr.logger.Error("Failed to render 405 page", zap.Error(err))
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
		}

		// Fallback to simple 405
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}

// TODO: Why is this validation not implemented.
// isValidLocale checks if a locale is supported using config
// func (rr *routeRegistrar) isValidLocale(locale string) bool {

// 	// Get supported locales from config service
// 	supportedLocales := rr.configService.GetSupportedLocales()
// 	for _, valid := range supportedLocales {
// 		if locale == valid {
// 			return true
// 		}
// 	}
// 	return false
// }

// registerLocaleSpecificRoutes registers specific routes for each valid locale
func (rr *routeRegistrar) registerLocaleSpecificRoutes(route interfaces.Route) error {
	// ConfigService must be properly injected - no fallbacks
	if rr.configService == nil {
		return fmt.Errorf("ConfigService is nil - this is a DI configuration error")
	}

	// Get supported locales from config service
	validLocales := rr.configService.GetSupportedLocales()

	if len(validLocales) == 0 {
		return fmt.Errorf("no supported locales configured")
	}

	for _, locale := range validLocales {
		// Create locale-specific route
		localeRoute := interfaces.Route{
			Path:         strings.ReplaceAll(route.Path, "$locale", locale),
			TemplateFile: route.TemplateFile,
			IsDynamic:    route.IsDynamic,
		}

		// Convert pattern for Chi router
		chiPattern := rr.convertRoutePattern(localeRoute.Path)

		// Build handler with middleware pipeline
		handler := rr.handlerBuilder.BuildHandler(localeRoute)

		// Register with Chi router
		rr.router.Get(chiPattern, handler.ServeHTTP)

		rr.logger.Debug("Locale-specific route registered",
			zap.String("locale", locale),
			zap.String("original_pattern", route.Path),
			zap.String("chi_pattern", chiPattern),
			zap.String("template", route.TemplateFile))
	}

	return nil
}
