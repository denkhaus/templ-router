package router

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/pipeline"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// HandlerBuilder defines the contract for handler building
type HandlerBuilder interface {
	BuildHandler(route interfaces.Route) http.Handler
	BuildStaticHandler(path string) http.Handler
	BuildErrorHandler(statusCode int, message string) http.HandlerFunc
}

// handlerBuilder handles HTTP handler building logic (private implementation)
// SEPARATED FROM: clean_router.go (Separation of Concerns)
type handlerBuilder struct {
	handlerPipeline *pipeline.HandlerPipeline
	configLoader    ConfigLoader
	logger          *zap.Logger
}

// NewHandlerBuilder creates a new handler builder
func NewHandlerBuilder(i do.Injector) (HandlerBuilder, error) {
	handlerPipeline := do.MustInvoke[*pipeline.HandlerPipeline](i)
	configLoader := do.MustInvoke[ConfigLoader](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &handlerBuilder{
		handlerPipeline: handlerPipeline,
		configLoader:    configLoader,
		logger:          logger,
	}, nil
}

// BuildHandler creates an HTTP handler for a specific route with middleware pipeline
func (hb *handlerBuilder) BuildHandler(route interfaces.Route) http.Handler {
	hb.logger.Debug("Building handler for route",
		zap.String("path", route.Path),
		zap.String("template", route.TemplateFile),
		zap.Bool("dynamic", route.IsDynamic))

	// Load configuration for this route
	config, err := hb.configLoader.LoadConfig(route.TemplateFile)
	if err != nil {
		hb.logger.Warn("Failed to load config for route",
			zap.String("route", route.Path),
			zap.String("template", route.TemplateFile),
			zap.Error(err))
		// Continue with nil config
		config = nil
	}

	// Load auth settings
	authSettings, err := hb.configLoader.LoadAuthSettings(route.TemplateFile)
	if err != nil {
		hb.logger.Warn("Failed to load auth settings for route",
			zap.String("route", route.Path),
			zap.String("template", route.TemplateFile),
			zap.Error(err))
		// Continue with nil auth settings
		authSettings = nil
	}

	// Extract parameters for dynamic routes
	params := hb.extractParametersForRoute(route)

	// Build handler using the pipeline
	pipelineConfig := pipeline.PipelineConfig{
		Route:        route,
		AuthSettings: authSettings,
		Params:       params,
	}

	// Add config file if available
	if config != nil {
		pipelineConfig.ConfigFile = &pipeline.ConfigFile{
			AuthSettings: authSettings,
		}
	}

	handler := hb.handlerPipeline.BuildHandler(pipelineConfig)

	hb.logger.Debug("Handler built successfully",
		zap.String("route", route.Path),
		zap.String("template", route.TemplateFile))

	return handler
}

// extractParametersForRoute extracts parameters from a route pattern
func (hb *handlerBuilder) extractParametersForRoute(route interfaces.Route) map[string]string {
	params := make(map[string]string)

	// For dynamic routes, we need to extract parameters from the URL
	// This is a simplified version - the actual parameter extraction
	// happens in the middleware pipeline during request processing

	if route.IsDynamic {
		hb.logger.Debug("Route has dynamic parameters",
			zap.String("route", route.Path))

		// The actual parameter values will be extracted from the HTTP request
		// in the parameter extraction middleware
		// For now, we just mark that this route expects parameters
		params["_dynamic"] = "true"
	}

	return params
}

// convertToMiddlewareRoute is no longer needed as we use interfaces.Route directly
// This function is kept for backward compatibility but should be removed in future versions

// BuildStaticHandler creates a handler for static file serving
func (hb *handlerBuilder) BuildStaticHandler(path string) http.Handler {
	hb.logger.Debug("Building static file handler", zap.String("path", path))

	// Create file server for static assets
	fileServer := http.FileServer(http.Dir("assets/"))
	handler := http.StripPrefix("/assets/", fileServer)

	hb.logger.Debug("Static handler built successfully", zap.String("path", path))

	return handler
}

// BuildErrorHandler creates a handler for error pages (404, 500, etc.)
func (hb *handlerBuilder) BuildErrorHandler(statusCode int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hb.logger.Debug("Error handler triggered",
			zap.Int("status_code", statusCode),
			zap.String("path", r.URL.Path),
			zap.String("message", message))

		// Set appropriate status code
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "text/html")

		// For now, return simple error message
		// TODO: Integrate with error service for proper error templates
		w.Write([]byte(message))
	}
}
