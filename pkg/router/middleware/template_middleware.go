package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// templateMiddleware handles template rendering concerns (private implementation)
type templateMiddleware struct {
	templateService    interfaces.TemplateService
	layoutService      interfaces.LayoutService
	errorService       interfaces.ErrorService
	parameterExtractor ParameterExtractor
	logger             *zap.Logger
}

// ParameterExtractor interface for extracting parameters from URLs (library-agnostic)
type ParameterExtractor interface {
	ExtractParameters(urlPath string, route interfaces.Route) map[string]string
	ExtractParametersFromRequest(r *http.Request, route interfaces.Route) map[string]string
}

// Import central types to eliminate redundancy
// Route and LayoutTemplate are now imported from interfaces package

// NewTemplateMiddleware creates a new template middleware for DI
func NewTemplateMiddleware(i do.Injector) (TemplateMiddlewareInterface, error) {
	templateService := do.MustInvoke[interfaces.TemplateService](i)
	layoutService := do.MustInvoke[interfaces.LayoutService](i)
	errorService := do.MustInvoke[interfaces.ErrorService](i)
	parameterExtractor := do.MustInvoke[ParameterExtractor](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &templateMiddleware{
		templateService:    templateService,
		layoutService:      layoutService,
		errorService:       errorService,
		parameterExtractor: parameterExtractor,
		logger:             logger,
	}, nil
}

// Handle processes template rendering for a request
func (tm *templateMiddleware) Handle(route interfaces.Route, params map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract parameters from the HTTP request using Chi's URL parameters (modern approach)
		extractedParams := tm.parameterExtractor.ExtractParametersFromRequest(r, route)

		// Merge with existing params (extracted params take precedence)
		for k, v := range extractedParams {
			params[k] = v
		}
		ctx := r.Context()

		// Load template config and add to context for router.M() access
		ctx = tm.addTemplateConfigToContext(ctx, route.TemplateFile)

		tm.logger.Debug("Rendering template",
			zap.String("template", route.TemplateFile),
			zap.String("path", route.Path),
			zap.Any("params", params),
			zap.Bool("requires_data_service", route.RequiresDataService),
			zap.String("data_service_interface", route.DataServiceInterface))

		// Render the page component (TemplateService now handles DataService templates directly)
		var component templ.Component
		var err error
		
		component, err = tm.templateService.RenderComponent(route, ctx, params)
		if err != nil {
			tm.logger.Error("Template rendering failed",
				zap.String("route", route.Path),
				zap.String("template", route.TemplateFile),
				zap.Error(err))

			// Render error component
			component = tm.errorService.CreateErrorComponent(err.Error(), route.Path)
		}

		// Wrap in layout if available
		if layout := tm.layoutService.FindLayoutForTemplate(route.TemplateFile); layout != nil {
			tm.logger.Debug("Wrapping component in layout",
				zap.String("layout", layout.FilePath),
				zap.Int("layout_level", layout.LayoutLevel))

			component = tm.layoutService.WrapInLayout(component, layout, ctx)
		}

		// Render the final component
		if component != nil {
			w.Header().Set("Content-Type", "text/html")
			if err := component.Render(ctx, w); err != nil {
				tm.logger.Error("Component rendering failed",
					zap.String("route", route.Path),
					zap.Error(err))
				http.Error(w, "Template rendering error", http.StatusInternalServerError)
			}
		} else {
			tm.renderFallback(w, route)
		}
	})
}

// renderFallback renders a fallback response when template is not found
func (tm *templateMiddleware) renderFallback(w http.ResponseWriter, route interfaces.Route) {
	tm.logger.Warn("Rendering fallback for missing template",
		zap.String("template", route.TemplateFile),
		zap.String("path", route.Path))

	w.Header().Set("Content-Type", "text/html")
	response := "<html><head><title>Template Not Found</title></head><body>"
	response += "<h1>Template not found: " + route.TemplateFile + "</h1>"
	response += "<p>Route: " + route.Path + "</p>"
	response += "<p>Please implement the templ component for this route.</p>"
	response += "</body></html>"
	w.Write([]byte(response))
}

// addTemplateConfigToContext loads template config and adds it to context for router.M() access
func (tm *templateMiddleware) addTemplateConfigToContext(ctx context.Context, templateFile string) context.Context {
	// Build YAML metadata path from template file
	yamlPath := tm.buildYamlPath(templateFile)

	tm.logger.Debug("Loading template config for context",
		zap.String("template_file", templateFile),
		zap.String("yaml_path", yamlPath))

	// Load shared config
	sharedConfig, err := shared.ParseYAMLMetadata(yamlPath)
	if err != nil {
		tm.logger.Debug("No template config found or failed to load",
			zap.String("yaml_path", yamlPath),
			zap.Error(err))
		return ctx // Return original context if no config
	}

	// Add shared config to context for router.M() access
	ctx = context.WithValue(ctx, "template_config", sharedConfig)
	tm.logger.Info("Added template metadata to context",
		zap.String("yaml_path", yamlPath),
		zap.Any("metadata", sharedConfig.RouteMetadata))

	return ctx
}

// buildYamlPath builds the YAML metadata path from template file path
func (tm *templateMiddleware) buildYamlPath(templateFile string) string {
	// Remove .templ extension and add .templ.yaml
	// e.g., "app/page.templ" -> "app/page.templ.yaml"
	if strings.HasSuffix(templateFile, ".templ") {
		return templateFile + ".yaml"
	}

	// Fallback: add .yaml to whatever we have
	return templateFile + ".yaml"
}

// REMOVED: extractParametersFromURL - replaced with pluggable ParameterExtractor interface
// This eliminates hardcoded "user" and "product" route assumptions, making the middleware library-agnostic
