package pipeline

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"

	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// HandlerPipeline creates clean, composable HTTP handlers using middleware pattern
type HandlerPipeline struct {
	authMiddleware     interfaces.AuthMiddlewareInterface
	i18nMiddleware     interfaces.I18nMiddlewareInterface
	templateMiddleware interfaces.TemplateMiddlewareInterface
	templateRegistry   interfaces.TemplateRegistry
	logger             *zap.Logger
}

// PipelineConfig contains configuration for building a handler pipeline
type PipelineConfig struct {
	Route      interfaces.Route
	AuthConfig *shared.AuthConfig
	ConfigFile *shared.ConfigFile
	Params     map[string]string
}

func NewHandlerPipeline(i do.Injector) (*HandlerPipeline, error) {
	authMiddleware := do.MustInvoke[interfaces.AuthMiddlewareInterface](i)
	i18nMiddleware := do.MustInvoke[interfaces.I18nMiddlewareInterface](i)
	templateMiddleware := do.MustInvoke[interfaces.TemplateMiddlewareInterface](i)
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
func (hp *HandlerPipeline) resolveAuthSettings(config PipelineConfig) *shared.AuthConfig {
	// Template-level auth settings take precedence
	if config.ConfigFile != nil && config.ConfigFile.Auth != nil {
		hp.logger.Debug("Using template-level auth settings",
			zap.String("route", config.Route.Path),
			zap.String("auth_type", string(config.ConfigFile.Auth.Type)))
		return config.ConfigFile.Auth
	}

	// Route-level auth settings
	if config.AuthConfig != nil {
		hp.logger.Debug("Using route-level auth settings",
			zap.String("route", config.Route.Path),
			zap.String("auth_type", string(config.AuthConfig.Type)))
		return config.AuthConfig
	}

	// Default to public
	hp.logger.Debug("Using default public auth settings",
		zap.String("route", config.Route.Path))
	return &shared.AuthConfig{Type: shared.AuthTypePublic}
}

// BuildHandlerFunc creates an http.HandlerFunc using the pipeline
func (hp *HandlerPipeline) BuildHandlerFunc(config PipelineConfig) http.HandlerFunc {
	handler := hp.BuildHandler(config)
	return handler.ServeHTTP
}
