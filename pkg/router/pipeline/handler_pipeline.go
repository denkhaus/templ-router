package pipeline

import (
	"context"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"

	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// HandlerPipeline creates clean, composable HTTP handlers using middleware pattern
type HandlerPipeline struct {
	authMiddleware     middleware.AuthMiddlewareInterface
	i18nMiddleware     middleware.I18nMiddlewareInterface
	templateMiddleware middleware.TemplateMiddlewareInterface
	templateRegistry   interfaces.TemplateRegistry
	logger             *zap.Logger
}

// PipelineConfig contains configuration for building a handler pipeline
type PipelineConfig struct {
	Route        interfaces.Route
	AuthSettings *interfaces.AuthSettings
	ConfigFile   *ConfigFile
	Params       map[string]string
}

// ConfigFile represents template configuration (simplified)
type ConfigFile struct {
	AuthSettings *interfaces.AuthSettings
	// Add other config fields as needed
}

func NewHandlerPipeline(i do.Injector) (*HandlerPipeline, error) {
	authMiddleware := do.MustInvoke[middleware.AuthMiddlewareInterface](i)
	i18nMiddleware := do.MustInvoke[middleware.I18nMiddlewareInterface](i)
	templateMiddleware := do.MustInvoke[middleware.TemplateMiddlewareInterface](i)
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &HandlerPipeline{
		authMiddleware:     authMiddleware,
		i18nMiddleware:     i18nMiddleware,
		templateMiddleware: templateMiddleware,
		templateRegistry:   templateRegistry,
		logger:             logger,
	}, nil

}

// BuildHandler creates a complete HTTP handler using the middleware pipeline
func (hp *HandlerPipeline) BuildHandler(config PipelineConfig) http.Handler {
	hp.logger.Debug("Building handler pipeline",
		zap.String("route", config.Route.Path),
		zap.String("template", config.Route.TemplateFile),
		zap.Bool("requires_data_service", config.Route.RequiresDataService))

	// Start with the innermost handler (TemplateService now handles DataService templates directly)
	var handler http.Handler
	
	if config.Route.RequiresDataService {
		hp.logger.Debug("Route requires DataService - will be handled by TemplateService",
			zap.String("route", config.Route.Path),
			zap.String("data_service_interface", config.Route.DataServiceInterface))
	}
	
	// All routes use template middleware (which now handles DataService templates internally)
	handler = hp.templateMiddleware.Handle(config.Route, config.Params)

	// Wrap with i18n middleware
	handler = hp.i18nMiddleware.Handle(handler, config.Route.TemplateFile)

	// Wrap with auth middleware (outermost)
	authSettings := hp.resolveAuthSettings(config)
	handler = hp.authMiddleware.Handle(handler, authSettings)

	return handler
}

// resolveAuthSettings determines the final auth settings for a route
func (hp *HandlerPipeline) resolveAuthSettings(config PipelineConfig) *interfaces.AuthSettings {
	// Template-level auth settings take precedence
	if config.ConfigFile != nil && config.ConfigFile.AuthSettings != nil {
		hp.logger.Debug("Using template-level auth settings",
			zap.String("route", config.Route.Path),
			zap.String("auth_type", config.ConfigFile.AuthSettings.Type.String()))
		return config.ConfigFile.AuthSettings
	}

	// Route-level auth settings
	if config.AuthSettings != nil {
		hp.logger.Debug("Using route-level auth settings",
			zap.String("route", config.Route.Path),
			zap.String("auth_type", config.AuthSettings.Type.String()))
		return config.AuthSettings
	}

	// Default to public
	hp.logger.Debug("Using default public auth settings",
		zap.String("route", config.Route.Path))
	return &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}
}

// wrapWithTemplateKeyContext wraps a handler to set template_key in context
func (hp *HandlerPipeline) wrapWithTemplateKeyContext(next http.Handler, route interfaces.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		// Get route-to-template mapping from template registry
		routeMapping := hp.templateRegistry.GetRouteToTemplateMapping()
		
		// Try to find template key for this route
		templateKey, found := routeMapping[route.Path]
		if found {
			hp.logger.Debug("Setting template_key in context",
				zap.String("route", route.Path),
				zap.String("template_key", templateKey))
			
			// Add template_key to context
			ctx = context.WithValue(ctx, "template_key", templateKey)
			r = r.WithContext(ctx)
		} else {
			hp.logger.Warn("Could not resolve template_key for route",
				zap.String("route", route.Path),
				zap.String("template_file", route.TemplateFile))
		}
		
		next.ServeHTTP(w, r)
	})
}

// BuildHandlerFunc creates an http.HandlerFunc using the pipeline
func (hp *HandlerPipeline) BuildHandlerFunc(config PipelineConfig) http.HandlerFunc {
	handler := hp.BuildHandler(config)
	return handler.ServeHTTP
}
