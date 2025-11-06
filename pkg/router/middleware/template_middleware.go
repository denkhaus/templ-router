package middleware

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	templateService           interfaces.TemplateService
	layoutService             interfaces.LayoutService
	errorService              interfaces.ErrorService
	parameterExtractor        ParameterExtractor
	templateRegistry          interfaces.TemplateRegistry
	componentMetadataService  interfaces.ComponentMetadataService
	i18nService               interfaces.I18nService
	configService             interfaces.ConfigService
	logger                    *zap.Logger
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
	componentMetadataService := do.MustInvoke[interfaces.ComponentMetadataService](i)
	i18nService := do.MustInvoke[interfaces.I18nService](i)
	configService := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &templateMiddleware{
		templateService:           templateService,
		layoutService:             layoutService,
		errorService:              errorService,
		parameterExtractor:        parameterExtractor,
		templateRegistry:          templateRegistry,
		componentMetadataService:  componentMetadataService,
		i18nService:               i18nService,
		configService:             configService,
		logger:                    logger,
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

		// Load component metadata based on route type
		if tm.isComponentTemplate(route.TemplateFile) {
			tm.logger.Debug("Component template detected, loading metadata",
				zap.String("template_file", route.TemplateFile))
			// Load metadata for components
			ctx = tm.addComponentMetadataForComponentRoute(ctx, route.TemplateFile)
		} else {
			tm.logger.Debug("Page template detected, embedded component loading temporarily disabled",
				zap.String("template_file", route.TemplateFile))
			// Temporarily disable embedded component loading for pages
			// ctx = tm.addEmbeddedComponentsMetadataToContext(ctx, route.TemplateFile)
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

		// Check if we should skip layout wrapping
		shouldSkipLayout := tm.shouldSkipLayout(route, r)

		// Wrap in layout if available and not skipping
		if !shouldSkipLayout {
			if layout := tm.layoutService.FindLayoutForTemplate(route.TemplateFile); layout != nil {
				tm.logger.Debug("Wrapping component in layout",
					zap.String("layout", layout.FilePath),
					zap.Int("layout_level", layout.LayoutLevel))

				component = tm.layoutService.WrapInLayout(component, layout, ctx)
			}
		} else {
			tm.logger.Debug("Skipping layout wrapping for partial rendering",
				zap.String("route", route.Path),
				zap.Bool("is_component_route", tm.isComponentRoute(route.Path)),
				zap.Bool("is_htmx_request", tm.isHTMXRequest(r)))
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

// REMOVED: loadComponentMetadata and buildComponentYamlPath
// Component metadata loading is now handled by ComponentMetadataService
// This removes tight coupling and improves separation of concerns

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
	// Extract component name from template file path using a generic approach
	componentName := tm.extractComponentNameFromTemplatePath(templateFile)
	if componentName == "" {
		tm.logger.Debug("Could not extract component name from template path",
			zap.String("template_file", templateFile))
		return ctx
	}

	// Load component metadata using the service
	componentConfig, err := tm.componentMetadataService.LoadComponentMetadata(componentName)
	if err != nil {
		// No component metadata available, return unchanged context
		tm.logger.Debug("No component metadata found, using page context only",
			zap.String("component", componentName),
			zap.String("template_file", templateFile),
			zap.Error(err))
		return ctx
	}

	tm.logger.Debug("Successfully loaded component metadata via service",
		zap.String("component", componentName),
		zap.String("template_file", templateFile))

	// Load component translations into i18n context
	ctx = tm.loadComponentTranslations(ctx, componentName)

	// Merge component metadata with existing page context
	return tm.mergeComponentMetadata(ctx, componentConfig)
}

// addEmbeddedComponentsMetadataToContext parses template content to find embedded components and loads their metadata
func (tm *templateMiddleware) addEmbeddedComponentsMetadataToContext(ctx context.Context, templateFile string) context.Context {
	// Read template file content
	templateContent, err := tm.readTemplateFile(templateFile)
	if err != nil {
		tm.logger.Debug("Failed to read template file for component detection",
			zap.String("template_file", templateFile),
			zap.Error(err))
		return ctx
	}

	// Extract component names from template content
	componentNames := tm.extractComponentsFromTemplateContent(templateContent)
	if len(componentNames) == 0 {
		tm.logger.Debug("No components found in template content",
			zap.String("template_file", templateFile))
		return ctx
	}

	tm.logger.Debug("Found embedded components in template",
		zap.String("template_file", templateFile),
		zap.Strings("components", componentNames))

	// Load metadata for each discovered component
	for _, componentName := range componentNames {
		ctx = tm.loadSingleComponentMetadata(ctx, componentName)
	}

	return ctx
}

// readTemplateFile reads the content of a template file
func (tm *templateMiddleware) readTemplateFile(templateFile string) (string, error) {
	// TemplateFile should already be the full path from route discovery
	// If it's a relative path, make it absolute using the layout root
	if !filepath.IsAbs(templateFile) {
		layoutRoot := tm.configService.GetLayoutRootDirectory()
		templateFile = filepath.Join(layoutRoot, templateFile)
	}

	content, err := os.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// extractComponentsFromTemplateContent parses template content to find component usages
// Looks for patterns like @components.Footer(), @navbar(), @components.UserSection()
func (tm *templateMiddleware) extractComponentsFromTemplateContent(content string) []string {
	var components []string

	// Pattern to match component usage in templates
	// Examples: @components.Footer(), @navbar(), @components.UserSection()
	pattern := regexp.MustCompile(`@(\w+(?:\.\w+)*)\s*\(`)

	matches := pattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			componentRef := match[1]

			// Convert the reference to component name
			// - "components.Footer" -> "Footer"
			// - "navbar" -> "navbar"
			// - "components.UserSection" -> "UserSection"
			parts := strings.Split(componentRef, ".")
			var componentName string
			if len(parts) > 1 {
				componentName = parts[len(parts)-1]
			} else {
				componentName = parts[0]
			}

			// Avoid duplicates
			found := false
			for _, existing := range components {
				if existing == componentName {
					found = true
					break
				}
			}
			if !found {
				components = append(components, componentName)
			}
		}
	}

	return components
}

// loadSingleComponentMetadata loads metadata for a single component
func (tm *templateMiddleware) loadSingleComponentMetadata(ctx context.Context, componentName string) context.Context {
	// Load component metadata using the service
	componentConfig, err := tm.componentMetadataService.LoadComponentMetadata(componentName)
	if err != nil {
		// No component metadata available, return unchanged context
		tm.logger.Debug("No component metadata found for embedded component",
			zap.String("component", componentName),
			zap.Error(err))
		return ctx
	}

	tm.logger.Debug("Successfully loaded embedded component metadata",
		zap.String("component", componentName))

	// Load component translations into i18n context
	ctx = tm.loadComponentTranslations(ctx, componentName)

	// Merge component metadata with existing page context
	return tm.mergeComponentMetadata(ctx, componentConfig)
}

// extractComponentNameFromTemplatePath extracts component name from template file path in a generic way
// This method works with any project structure by looking for the last directory name before the template file
// Examples:
// - "components/footer/page.templ" -> "footer"
// - "app/components/navbar.templ" -> "navbar"
// - "templates/ui/header.templ" -> "header"
// - "web/partials/sidebar.templ" -> "sidebar"
func (tm *templateMiddleware) extractComponentNameFromTemplatePath(templatePath string) string {
	// Remove file extension if present
	basePath := templatePath
	if strings.HasSuffix(basePath, ".templ") {
		basePath = basePath[:len(basePath)-len(".templ")]
	}

	// Split by directory separators
	parts := strings.Split(basePath, "/")
	if len(parts) < 2 {
		return ""
	}

	// For component templates, the component name is typically the filename
	// For page templates, it might be the directory name
	// We'll use the last part (filename) for component templates
	// and the second-to-last part for page templates
	var componentName string
	lastPart := parts[len(parts)-1]

	// If the last part is "page", use the directory name (second-to-last)
	// Otherwise, use the filename itself
	if lastPart == "page" && len(parts) >= 2 {
		componentName = parts[len(parts)-2]
	} else {
		componentName = lastPart
	}

	// This works for structures like:
	// - components/footer/page -> footer (uses directory)
	// - components/footer -> footer (uses filename)
	// - app/components/navbar -> navbar (uses filename)

	// Clean up the component name
	componentName = strings.TrimSpace(componentName)
	if componentName == "" {
		return ""
	}

	return componentName
}

// shouldSkipLayout determines if layout wrapping should be skipped for partial rendering
func (tm *templateMiddleware) shouldSkipLayout(route interfaces.Route, r *http.Request) bool {
	// Skip layout for component routes (e.g., /components/footer)
	if tm.isComponentRoute(route.Path) {
		tm.logger.Debug("Skipping layout for component route",
			zap.String("route", route.Path))
		return true
	}

	// Skip layout for HTMX requests (AJAX partial loading)
	if tm.isHTMXRequest(r) {
		tm.logger.Debug("Skipping layout for HTMX request",
			zap.String("route", route.Path))
		return true
	}

	return false
}

// isHTMXRequest detects if the request is an HTMX AJAX request
func (tm *templateMiddleware) isHTMXRequest(r *http.Request) bool {
	// HTMX sends this header for AJAX requests
	return r.Header.Get("HX-Request") != ""
}

// isPageTemplate checks if a template file is a page template
// Page templates are typically named page.templ
func (tm *templateMiddleware) isPageTemplate(templateFile string) bool {
	filename := filepath.Base(templateFile)
	filename = strings.TrimSuffix(filename, ".templ")

	// Check if this is a page template
	isPage := filename == "page"

	tm.logger.Debug("Template type check",
		zap.String("template_file", templateFile),
		zap.String("filename", filename),
		zap.Bool("is_page", isPage))

	return isPage
}


// addComponentMetadataForComponentRoute loads component metadata for component routes
func (tm *templateMiddleware) addComponentMetadataForComponentRoute(ctx context.Context, templateFile string) context.Context {
	// Extract component name from template file path using a generic approach
	componentName := tm.extractComponentNameFromTemplatePath(templateFile)
	if componentName == "" {
		tm.logger.Debug("Could not extract component name from template path for component route",
			zap.String("template_file", templateFile))
		return ctx
	}

	// Load component metadata using the service
	componentConfig, err := tm.componentMetadataService.LoadComponentMetadata(componentName)
	if err != nil {
		// No component metadata available, return unchanged context
		tm.logger.Debug("No component metadata found for component route",
			zap.String("component", componentName),
			zap.String("template_file", templateFile),
			zap.Error(err))
		return ctx
	}

	tm.logger.Debug("Successfully loaded component metadata for component route",
		zap.String("component", componentName),
		zap.String("template_file", templateFile))

	// Load component translations into i18n context
	ctx = tm.loadComponentTranslations(ctx, componentName)

	// Merge component metadata with existing context
	return tm.mergeComponentMetadata(ctx, componentConfig)
}

// templateContainsComponents checks if a template file contains embedded components
// Uses route mapping to check if any component routes exist in the system
func (tm *templateMiddleware) templateContainsComponents(templateFile string) bool {
	// Get the route mapping to check for component routes
	routeMapping := tm.templateRegistry.GetRouteToTemplateMapping()

	// If there are any routes that look like component routes,
	// then the system has component templates that could be embedded
	for route := range routeMapping {
		// Check if this route is for a component (using existing logic)
		if tm.isComponentRoute(route) {
			tm.logger.Debug("Found component route in system",
				zap.String("template_file", templateFile),
				zap.String("component_route", route))
			return true
		}
	}

	// No component routes found
	tm.logger.Debug("No component routes found in system",
		zap.String("template_file", templateFile))
	return false
}

// isComponentTemplate determines if a template is a component template
// Generic approach: components are anything that's not page.templ, layout.templ, or error.templ
func (tm *templateMiddleware) isComponentTemplate(templatePath string) bool {
	filename := filepath.Base(templatePath)
	filename = strings.TrimSuffix(filename, ".templ")

	// These are the standard template types that are NOT components
	standardTemplates := map[string]bool{
		"page":    true,  // page.templ
		"layout":  true,  // layout.templ
		"error":   true,  // error.templ
	}

	// If it's not a standard template type, it's a component
	isComponent := !standardTemplates[filename]

	tm.logger.Debug("Template type classification",
		zap.String("template_path", templatePath),
		zap.String("filename", filename),
		zap.Bool("is_component", isComponent))

	return isComponent
}

// extractTemplatePathFromKey extracts template path from template registry key
// Helper method to avoid duplication
func (tm *templateMiddleware) extractTemplatePathFromKey(templateKey string) string {
	// First try to split by # to separate the path from the template name
	parts := strings.Split(templateKey, "#")
	if len(parts) >= 2 && parts[0] != "" {
		return parts[0]
	}

	// For hash keys, this would need reverse lookup, but for now return empty
	// The component detection can work with just the presence of non-standard templates
	return ""
}

// loadComponentTranslations loads component translations into the i18n context
func (tm *templateMiddleware) loadComponentTranslations(ctx context.Context, componentName string) context.Context {
	// Use I18nService to load component translations into context
	if tm.i18nService == nil {
		tm.logger.Debug("I18nService not available, skipping component translation loading",
			zap.String("component", componentName))
		return ctx
	}

	// Call the I18nService method to load component translations
	newCtx := tm.i18nService.LoadComponentTranslationsIntoContext(ctx, componentName)

	tm.logger.Debug("Loaded component translations using I18nService",
		zap.String("component", componentName))

	return newCtx
}

// REMOVED: extractParametersFromURL - replaced with pluggable ParameterExtractor interface
// This eliminates hardcoded "user" and "product" route assumptions, making the middleware library-agnostic
