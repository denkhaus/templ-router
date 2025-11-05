package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
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
	templateRegistry   interfaces.TemplateRegistry
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
func NewTemplateMiddleware(i do.Injector) (interfaces.TemplateMiddlewareInterface, error) {
	templateService := do.MustInvoke[interfaces.TemplateService](i)
	layoutService := do.MustInvoke[interfaces.LayoutService](i)
	errorService := do.MustInvoke[interfaces.ErrorService](i)
	parameterExtractor := do.MustInvoke[ParameterExtractor](i)
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &templateMiddleware{
		templateService:    templateService,
		layoutService:      layoutService,
		errorService:       errorService,
		parameterExtractor: parameterExtractor,
		templateRegistry:   templateRegistry,
		logger:             logger,
	}, nil
}

// Handle processes template rendering for a request
func (tm *templateMiddleware) Handle(route interfaces.Route, params map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Store request path in context for i18n.GetCurrentRoute() access
		ctx = context.WithValue(ctx, shared.RequestPathKey, r.URL.Path)

		// Store route mapping in context for i18n.GetCurrentRouteWithoutLocale() to determine localized routes
		routeMapping := tm.templateRegistry.GetRouteToTemplateMapping()
		ctx = context.WithValue(ctx, shared.RouteMappingKey, routeMapping)

		// Create RouterContext with unified access to all request data
		routerCtx := NewRouterContext(ctx, r)

		// Load template config and add to context for router.M() access
		ctx = tm.addTemplateConfigToContext(ctx, route.TemplateFile)

		// Load component metadata if this is a component route and merge with page metadata
		if tm.isComponentRoute(route.Path) {
			ctx = tm.addComponentMetadataToContext(ctx, route.TemplateFile)
		}

		tm.logger.Debug("Rendering template",
			zap.String("template", route.TemplateFile),
			zap.String("path", route.Path),
			zap.Any("url_params", routerCtx.GetAllURLParams()),
			zap.Any("query_params", routerCtx.GetAllQueryParams()),
			zap.Bool("requires_data_service", route.RequiresDataService),
			zap.String("data_service_interface", route.DataServiceInterface))

		// Render the page component (TemplateService now handles DataService templates directly)
		var component templ.Component
		var err error

		component, err = tm.templateService.RenderComponent(route, routerCtx, ctx)
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
	configFileFound, sharedConfig, err := shared.ParseYAMLMetadata(yamlPath)
	if err != nil {
		if configFileFound {
			tm.logger.Debug("Failed to load template config",
				zap.String("yaml_path", yamlPath),
				zap.Error(err),
			)
		}

		return ctx // Return original context if no config
	}

	// Add shared config to context for router.M() access
	ctx = context.WithValue(ctx, shared.TemplateConfigKey, sharedConfig)
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

// loadComponentMetadata loads component YAML metadata for a given template file
// Returns the loaded config or nil if no metadata file exists
func (tm *templateMiddleware) loadComponentMetadata(templateFile string) *shared.ConfigFile {
	// Build component-specific YAML metadata path
	yamlPath := tm.buildComponentYamlPath(templateFile)

	tm.logger.Debug("Loading component metadata",
		zap.String("template_file", templateFile),
		zap.String("component_yaml_path", yamlPath))

	// Load component config
	configFileFound, config, err := shared.ParseYAMLMetadata(yamlPath)
	if err != nil {
		if configFileFound {
			tm.logger.Debug("Failed to load component metadata",
				zap.String("yaml_path", yamlPath),
				zap.Error(err))
		} else {
			tm.logger.Debug("No component metadata file found",
				zap.String("yaml_path", yamlPath))
		}
		return nil // No metadata available
	}

	tm.logger.Debug("Successfully loaded component metadata",
		zap.String("yaml_path", yamlPath),
		zap.Any("metadata", config.RouteMetadata))

	return config
}

// buildComponentYamlPath builds the correct YAML metadata path for component templates
// Components have different path structure than pages
// e.g., "app/components/footer/page.templ" -> "app/components/footer.templ.yaml"
func (tm *templateMiddleware) buildComponentYamlPath(templateFile string) string {
	// Check if this is a component template path (contains /components/ and ends with /page.templ)
	if strings.Contains(templateFile, "/components/") && strings.HasSuffix(templateFile, "/page.templ") {
		// Convert component path: app/components/footer/page.templ -> app/components/footer.templ
		componentPath := strings.TrimSuffix(templateFile, "/page.templ")
		return componentPath + ".templ.yaml"
	}

	// Fallback to standard YAML path generation
	return tm.buildYamlPath(templateFile)
}

// mergeComponentMetadata merges component metadata with existing page context
// Component metadata takes precedence over page metadata
func (tm *templateMiddleware) mergeComponentMetadata(ctx context.Context, componentConfig *shared.ConfigFile) context.Context {
	// Get existing config from context
	existingConfig := ctx.Value(shared.TemplateConfigKey)
	if existingConfig == nil {
		// No existing config, use component config directly
		ctx = context.WithValue(ctx, shared.TemplateConfigKey, componentConfig)
		return tm.addComponentI18nToContext(ctx, componentConfig)
	}

	// Get page/template config
	pageConfig, ok := existingConfig.(*shared.ConfigFile)
	if !ok {
		// Invalid config type, use component config
		ctx = context.WithValue(ctx, shared.TemplateConfigKey, componentConfig)
		return tm.addComponentI18nToContext(ctx, componentConfig)
	}

	// Merge configs: component takes precedence over page
	merged := tm.mergeConfigs(pageConfig, componentConfig)

	tm.logger.Info("Merged component metadata with page metadata (component takes precedence)",
		zap.Any("page_metadata", pageConfig.RouteMetadata),
		zap.Any("component_metadata", componentConfig.RouteMetadata),
		zap.Any("merged_metadata", merged.RouteMetadata))

	ctx = context.WithValue(ctx, shared.TemplateConfigKey, merged)
	return tm.addComponentI18nToContext(ctx, componentConfig)
}

// addComponentI18nToContext adds component i18n data to the I18n context
func (tm *templateMiddleware) addComponentI18nToContext(ctx context.Context, componentConfig *shared.ConfigFile) context.Context {
	if componentConfig.MultiLocaleI18n == nil {
		tm.logger.Debug("No component i18n data to add to context")
		return ctx
	}

	// Get current locale from context
	locale, _ := ctx.Value(shared.LocaleKey).(string)
	if locale == "" {
		locale = "en" // Default fallback
	}

	// Check if component has translations for current locale
	componentTranslations, hasLocale := componentConfig.MultiLocaleI18n[locale]
	if !hasLocale {
		tm.logger.Debug("No component translations found for locale",
			zap.String("locale", locale),
			zap.Strings("available_locales", tm.getAvailableLocales(componentConfig.MultiLocaleI18n)))
		return ctx
	}

	// Get existing i18n data from context
	i18nData, ok := ctx.Value(shared.I18nDataKey).(*i18n.I18nData)
	if !ok {
		tm.logger.Debug("No existing i18n data found in context")
		return ctx
	}

	// Add component translations to existing i18n data
	tm.logger.Info("Adding component i18n translations to context",
		zap.String("locale", locale),
		zap.Int("translation_count", len(componentTranslations)))

	// Create a copy of the translations map to avoid race conditions
	newTranslations := make(map[string]string)
	for key, value := range i18nData.Translations {
		newTranslations[key] = value
	}

	// Component translations take precedence over existing translations
	for key, value := range componentTranslations {
		newTranslations[key] = value
	}

	// Update i18n data with merged translations
	updatedI18nData := &i18n.I18nData{
		Locale:           i18nData.Locale,
		Translations:     newTranslations,
		CurrentTemplate:  i18nData.CurrentTemplate,
		Logger:           i18nData.Logger,
	}

	return context.WithValue(ctx, shared.I18nDataKey, updatedI18nData)
}

// getAvailableLocales returns a slice of available locale strings
func (tm *templateMiddleware) getAvailableLocales(i18nData map[string]map[string]string) []string {
	locales := make([]string, 0, len(i18nData))
	for locale := range i18nData {
		locales = append(locales, locale)
	}
	return locales
}

// mergeConfigs creates a merged config where higherPriorityConfig overrides baseConfig
func (tm *templateMiddleware) mergeConfigs(baseConfig, higherPriorityConfig *shared.ConfigFile) *shared.ConfigFile {
	// Start with base config as base
	merged := &shared.ConfigFile{
		RouteMetadata:   baseConfig.RouteMetadata,
		MultiLocaleI18n: make(map[string]map[string]string),
		AuthSettings:    baseConfig.AuthSettings,
		FilePath:        baseConfig.FilePath,
		TemplateFilePath: baseConfig.TemplateFilePath,
		I18nMappings:    make(map[string]string),
		ErrorSettings:   baseConfig.ErrorSettings,
		DynamicSettings: baseConfig.DynamicSettings,
	}

	// Copy base i18n data first
	if baseConfig.MultiLocaleI18n != nil {
		for locale, translations := range baseConfig.MultiLocaleI18n {
			merged.MultiLocaleI18n[locale] = make(map[string]string)
			for key, value := range translations {
				merged.MultiLocaleI18n[locale][key] = value
			}
		}
	}

	// Copy base i18n mappings
	if baseConfig.I18nMappings != nil {
		for key, value := range baseConfig.I18nMappings {
			merged.I18nMappings[key] = value
		}
	}

	// Override with higher priority metadata
	if higherPriorityConfig.RouteMetadata != nil {
		merged.RouteMetadata = higherPriorityConfig.RouteMetadata
	}

	// Override with higher priority i18n data
	if higherPriorityConfig.MultiLocaleI18n != nil {
		for locale, translations := range higherPriorityConfig.MultiLocaleI18n {
			if merged.MultiLocaleI18n[locale] == nil {
				merged.MultiLocaleI18n[locale] = make(map[string]string)
			}
			for key, value := range translations {
				merged.MultiLocaleI18n[locale][key] = value
			}
		}
	}

	// Override with higher priority i18n mappings
	if higherPriorityConfig.I18nMappings != nil {
		for key, value := range higherPriorityConfig.I18nMappings {
			merged.I18nMappings[key] = value
		}
	}

	// Use higher priority settings if available
	if higherPriorityConfig.AuthSettings != nil {
		merged.AuthSettings = higherPriorityConfig.AuthSettings
	}
	if higherPriorityConfig.ErrorSettings != nil {
		merged.ErrorSettings = higherPriorityConfig.ErrorSettings
	}
	if higherPriorityConfig.DynamicSettings != nil {
		merged.DynamicSettings = higherPriorityConfig.DynamicSettings
	}

	return merged
}

// isComponentRoute checks if the given route path represents a component route
// Component routes follow the pattern: /components/{component-name}
func (tm *templateMiddleware) isComponentRoute(routePath string) bool {
	return strings.HasPrefix(routePath, "/components/") && routePath != "/components/"
}

// addComponentMetadataToContext loads component metadata and merges it with existing context
func (tm *templateMiddleware) addComponentMetadataToContext(ctx context.Context, templateFile string) context.Context {
	// Load component metadata
	componentConfig := tm.loadComponentMetadata(templateFile)
	if componentConfig == nil {
		// No component metadata available, return unchanged context
		tm.logger.Debug("No component metadata found, using page context only",
			zap.String("template_file", templateFile))
		return ctx
	}

	// Merge component metadata with existing page context
	return tm.mergeComponentMetadata(ctx, componentConfig)
}

// REMOVED: extractParametersFromURL - replaced with pluggable ParameterExtractor interface
// This eliminates hardcoded "user" and "product" route assumptions, making the middleware library-agnostic
