package middleware

import (
	"context"
	"fmt"
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
	templateService          interfaces.TemplateService
	layoutService            interfaces.LayoutService
	errorService             interfaces.ErrorService
	parameterExtractor       ParameterExtractor
	templateRegistry         interfaces.TemplateRegistry
	componentMetadataService interfaces.ComponentMetadataService
	i18nService              interfaces.I18nService
	configService            interfaces.ConfigService
	logger                   *zap.Logger
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
		templateService:          templateService,
		layoutService:            layoutService,
		errorService:             errorService,
		parameterExtractor:       parameterExtractor,
		templateRegistry:         templateRegistry,
		componentMetadataService: componentMetadataService,
		i18nService:              i18nService,
		configService:            configService,
		logger:                   logger,
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

		// Get the full template key from route mapping
		templateKey := routeMapping[route.Path]
		if templateKey == "" {
			tm.logger.Debug("No template key found for route path",
				zap.String("route_path", route.Path),
				zap.String("template_file", route.TemplateFile))
			templateKey = route.TemplateFile // fallback to template file
		}

		// Load metadata and i18n for ALL templates using unified approach
		// Everything is a component - pages, layouts, and components use the same loading mechanism
		ctx = tm.loadTemplateMetadataAndI18n(ctx, templateKey, route.Path)

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

// loadTemplateMetadataAndI18n loads metadata and i18n for ALL templates using unified approach
// Everything is a component - pages, layouts, and components use the same loading mechanism
func (tm *templateMiddleware) loadTemplateMetadataAndI18n(ctx context.Context, templateKey, routePath string) context.Context {
	// Get template metadata from registry
	metadata, err := tm.templateRegistry.GetTemplateMetadata(templateKey)
	if err != nil {
		tm.logger.Debug("No metadata found in registry for template",
			zap.String("template_key", templateKey),
			zap.String("route_path", routePath))
		return ctx
	}

	tm.logger.Debug("Loading template metadata and i18n",
		zap.String("template_key", templateKey),
		zap.String("route_path", routePath),
		zap.String("template_type", string(metadata.Type)),
		zap.String("component_name", metadata.ComponentName))

	// Initialize i18n context for ALL templates
	ctx = tm.initializeI18nContext(ctx, metadata.TemplatePath)

	// Load hierarchical metadata and i18n: Layout -> Page -> Component
	// Components have highest priority and override parent settings
	mergedConfig, mergedTranslations, componentsConfigs, err := tm.loadHierarchicalMetadataAndI18n(ctx, metadata)
	if err != nil {
		tm.logger.Error("Failed loading metadata and i18n config for template",
			zap.String("template_key", templateKey),
			zap.String("route_path", routePath),
			zap.Error(err),
		)
		return ctx
	}

	// Add merged template config to context for metadata.M() access
	ctx = context.WithValue(ctx, shared.TemplateConfigKey, mergedConfig)

	// Add components metadata to context for component-specific metadata resolution
	ctx = context.WithValue(ctx, shared.ComponentsMetadataKey, componentsConfigs)

	// Update i18n context with merged translations
	if i18nData, ok := ctx.Value(shared.I18nDataKey).(*i18n.I18nData); ok {
		i18nData.Translations = mergedTranslations
		ctx = context.WithValue(ctx, shared.I18nDataKey, i18nData)
	}

	tm.logger.Debug("Processed template metadata and i18n",
		zap.String("component_name", metadata.ComponentName),
		zap.String("template_type", string(metadata.Type)),
		zap.Bool("yaml_exists", metadata.YAMLExists),
		zap.Bool("has_i18n", metadata.HasI18n))

	return ctx
}

// initializeI18nContext creates i18n context for any template
func (tm *templateMiddleware) initializeI18nContext(ctx context.Context, templatePath string) context.Context {
	// Use i18n service to create context for any template type
	return tm.i18nService.CreateContext(ctx, templatePath)
}

// buildCorrectYAMLPath builds the correct YAML file path from registry metadata
func (tm *templateMiddleware) buildCorrectYAMLPath(registryYAMLPath string) string {
	// The registry gives us paths like "app/components/footer.templ.yaml"
	// layoutRoot = "demo/app" means layout root directory is "demo/app"
	// working directory = project root = "/app"
	// registry path = "app/components/footer.templ.yaml"
	// actual file location = "demo/app/components/footer.templ.yaml"
	// We need to replace "app/" with the full layoutRoot path
	layoutRoot := tm.configService.GetLayoutRootDirectory()

	// If registry path starts with "app/", replace it with the full layout root
	if strings.HasPrefix(registryYAMLPath, "app/") {
		// Replace "app/" with "demo/app/" (or whatever layoutRoot is)
		return strings.Replace(registryYAMLPath, "app/", layoutRoot+"/", 1)
	}

	// Fallback: prepend layout root if no "app/" prefix
	return layoutRoot + "/" + registryYAMLPath
}

// loadHierarchicalMetadataAndI18n loads and merges metadata and i18n data in hierarchical order
// Layout -> Page -> Component (higher priority overrides lower priority)
func (tm *templateMiddleware) loadHierarchicalMetadataAndI18n(ctx context.Context, currentMetadata *interfaces.TemplateMetadata) (*shared.ConfigFile, map[string]string, map[string]*shared.ConfigFile, error) {
	// Step 1: Load layout metadata (lowest priority)
	layoutConfig := tm.loadLayoutMetadata()

	// Step 2: Load current template metadata (page or component)
	currentConfig := tm.loadCurrentTemplateMetadata(currentMetadata)

	// Step 3: Load embedded components metadata (highest priority)
	componentsConfigs := tm.loadEmbeddedComponentsMetadata()

	// Step 4: Merge hierarchically: Layout + Current + Components
	mergedConfig, err := tm.mergeConfigsHierarchically(layoutConfig, currentConfig, componentsConfigs)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get current locale from context
	locale, _ := ctx.Value(shared.LocaleKey).(string)
	if locale == "" {
		locale = tm.configService.GetDefaultLocale() // Use configured default locale
	}

	// Step 5: Merge i18n translations hierarchically for current locale only
	mergedTranslations := tm.mergeTranslationsHierarchically(layoutConfig, currentConfig, componentsConfigs, locale)

	tm.logger.Debug("Hierarchical metadata merge completed",
		zap.String("template_type", string(currentMetadata.Type)),
		zap.String("template_path", currentMetadata.TemplatePath),
		zap.String("locale", locale),
		zap.Int("layout_metadata_count", len(tm.getRouteMetadata(layoutConfig))),
		zap.Int("current_metadata_count", len(tm.getRouteMetadata(currentConfig))),
		zap.Int("components_count", len(componentsConfigs)),
		zap.Int("merged_metadata_count", len(tm.getRouteMetadata(mergedConfig))),
		zap.Int("merged_translation_count", len(mergedTranslations)))

	return mergedConfig, mergedTranslations, componentsConfigs, nil
}

// loadLayoutMetadata loads layout metadata if available
func (tm *templateMiddleware) loadLayoutMetadata() *shared.ConfigFile {
	// Find layout template in registry
	layoutTemplateKey := shared.GenerateTemplateKey("app/layout.templ#Layout")
	layoutMetadata, err := tm.templateRegistry.GetTemplateMetadata(layoutTemplateKey)
	if err != nil || !layoutMetadata.YAMLExists {
		tm.logger.Debug("No layout metadata found or layout YAML doesn't exist")
		return tm.createEmptyConfig()
	}

	// Load layout YAML
	yamlPath := tm.buildCorrectYAMLPath(layoutMetadata.YAMLFile)
	_, config, err := shared.ParseYAMLMetadata(yamlPath)
	if err != nil {
		tm.logger.Debug("Failed to load layout metadata",
			zap.String("yaml_path", yamlPath),
			zap.Error(err))
		return tm.createEmptyConfig()
	}

	metadataCount := len(tm.getRouteMetadata(config))
	tm.logger.Debug("Successfully loaded layout metadata",
		zap.String("yaml_path", yamlPath),
		zap.Int("metadata_count", metadataCount))

	return config
}

// loadCurrentTemplateMetadata loads the current template's metadata
func (tm *templateMiddleware) loadCurrentTemplateMetadata(metadata *interfaces.TemplateMetadata) *shared.ConfigFile {
	if !metadata.YAMLExists || metadata.YAMLFile == "" {
		tm.logger.Debug("No YAML file for current template",
			zap.String("template_path", metadata.TemplatePath))
		return tm.createEmptyConfig()
	}

	// Load current template YAML
	yamlPath := tm.buildCorrectYAMLPath(metadata.YAMLFile)
	_, config, err := shared.ParseYAMLMetadata(yamlPath)
	if err != nil {
		tm.logger.Debug("Failed to load current template metadata",
			zap.String("yaml_path", yamlPath),
			zap.String("template_path", metadata.TemplatePath),
			zap.Error(err))
		return tm.createEmptyConfig()
	}

	tm.logger.Debug("Successfully loaded current template metadata",
		zap.String("template_path", metadata.TemplatePath),
		zap.String("yaml_path", yamlPath),
		zap.String("template_type", string(metadata.Type)),
		zap.Any("metadata", config.GetRouteMetadata()))

	return config
}

// loadEmbeddedComponentsMetadata finds and loads all embedded components metadata
func (tm *templateMiddleware) loadEmbeddedComponentsMetadata() map[string]*shared.ConfigFile {
	components := make(map[string]*shared.ConfigFile)

	// For now, load all available components as a simple approach
	// In a more sophisticated implementation, we could analyze template files
	// to find exactly which components are used
	routeMapping := tm.templateRegistry.GetRouteToTemplateMapping()
	for _, templateKey := range routeMapping {

		// Get component metadata
		componentMetadata, err := tm.templateRegistry.GetTemplateMetadata(templateKey)
		if err != nil {
			continue
		}

		// Only process component routes
		if componentMetadata.Type != interfaces.TemplateTypeComponent {
			continue
		}

		// Only load if YAML exists
		if !componentMetadata.YAMLExists || componentMetadata.YAMLFile == "" {
			continue
		}

		// Load component YAML
		yamlPath := tm.buildCorrectYAMLPath(componentMetadata.YAMLFile)
		_, config, err := shared.ParseYAMLMetadata(yamlPath)
		if err != nil {
			tm.logger.Debug("Failed to load component metadata",
				zap.String("component_name", componentMetadata.ComponentName),
				zap.String("yaml_path", yamlPath),
				zap.Error(err))
			continue
		}

		components[componentMetadata.ComponentName] = config
	}

	tm.logger.Debug("Loaded embedded components metadata",
		zap.Int("components_count", len(components)))

	return components
}

// createEmptyConfig creates an empty config structure
func (tm *templateMiddleware) createEmptyConfig() *shared.ConfigFile {
	config := &shared.ConfigFile{}

	// Initialize all the type-safe structures with empty data
	config.Metadata = &shared.MetadataConfig{
		Custom: make(map[string]interface{}),
	}
	config.I18n = &shared.I18nConfig{
		FlatMappings: make(map[string]string),
		Translations: make(map[string]*shared.LocaleTranslations),
	}
	config.Auth = &shared.AuthConfig{
		Settings: make(map[string]interface{}),
	}
	config.Layout = &shared.LayoutConfig{
		Settings: make(map[string]interface{}),
	}
	config.Error = &shared.ErrorConfig{
		Settings: make(map[string]interface{}),
	}
	config.Dynamic = &shared.DynamicConfig{
		Rules:    make(map[string]*shared.ValidationRule),
		Settings: make(map[string]interface{}),
	}

	return config
}

// mergeConfigsHierarchically merges configs in hierarchical order: Layout + Current + Components
// Components (highest priority) override Current, which overrides Layout (lowest priority)
func (tm *templateMiddleware) mergeConfigsHierarchically(layoutConfig, currentConfig *shared.ConfigFile, componentsConfigs map[string]*shared.ConfigFile) (*shared.ConfigFile, error) {
	// Start with layout as base
	merged := tm.cloneConfig(layoutConfig)

	// Merge current template (overrides layout)
	merged, err := tm.mergeTwoConfigs(merged, currentConfig, "current_template")
	if err != nil {
		return nil, fmt.Errorf("failed to merge config [layout]: %w", err)
	}

	// Merge all components (overrides current + layout)
	for componentName, componentConfig := range componentsConfigs {
		merged, err = tm.mergeTwoConfigs(merged, componentConfig, componentName)
		if err != nil {
			return nil, fmt.Errorf("failed to merge config [component]: %w", err)
		}
	}

	return merged, nil
}

// mergeTranslationsHierarchically merges i18n translations hierarchically: Layout + Current + Components
// Only merges translations for the specified locale
func (tm *templateMiddleware) mergeTranslationsHierarchically(layoutConfig, currentConfig *shared.ConfigFile, componentsConfigs map[string]*shared.ConfigFile, locale string) map[string]string {
	merged := make(map[string]string)

	// Step 1: Start with layout translations for CURRENT LOCALE only (lowest priority)
	layoutI18n := layoutConfig.GetMultiLocaleI18n()
	if layoutTranslations, exists := layoutI18n[locale]; exists {
		for key, value := range layoutTranslations {
			merged[key] = value // This will be overridden by higher priority items
		}
	}

	// Step 2: Merge current template translations for CURRENT LOCALE only (overrides layout)
	currentI18n := currentConfig.GetMultiLocaleI18n()
	if currentTranslations, exists := currentI18n[locale]; exists {
		for key, value := range currentTranslations {
			merged[key] = value // Overrides layout
		}
	}

	// Step 3: Merge all components translations for CURRENT LOCALE only (highest priority)
	for _, componentConfig := range componentsConfigs {
		componentI18n := componentConfig.GetMultiLocaleI18n()
		if componentTranslations, exists := componentI18n[locale]; exists {
			for key, value := range componentTranslations {
				merged[key] = value // Components have highest priority
			}
		}
	}

	return merged
}

// mergeTwoConfigs merges two configs with the second overriding the first
func (tm *templateMiddleware) mergeTwoConfigs(base, override *shared.ConfigFile, sourceName string) (*shared.ConfigFile, error) {
	merged := tm.cloneConfig(base)

	// Merge route metadata
	overrideMetadata := override.GetRouteMetadata()
	if overrideMetadata != nil {
		if overrideMetadataMap, ok := overrideMetadata.(map[string]interface{}); ok {
			if merged.Metadata == nil {
				merged.Metadata = &shared.MetadataConfig{
					Custom: make(map[string]interface{}),
				}
			}
			// Use the MergeMetadata method which handles all the merging logic
			overrideConfig := &shared.ConfigFile{Metadata: &shared.MetadataConfig{Custom: overrideMetadataMap}}
			merged.MergeMetadata(overrideConfig)
		}
	}

	// Merge i18n data
	overrideI18n := override.GetMultiLocaleI18n()
	if overrideI18n != nil {
		// For i18n, we need to work with the I18nConfig directly
		if merged.I18n == nil {
			merged.I18n = &shared.I18nConfig{
				FlatMappings: make(map[string]string),
				Translations: make(map[string]*shared.LocaleTranslations),
			}
		}
		for locale, translations := range overrideI18n {
			if merged.I18n.Translations[locale] == nil {
				merged.I18n.Translations[locale] = &shared.LocaleTranslations{
					Locale:       locale,
					Translations: make(map[string]interface{}),
				}
			}
			for key, value := range translations {
				merged.I18n.Translations[locale].Translations[key] = value
			}
		}
	}

	// Merge auth settings
	overrideAuth := override.GetAuthSettings()
	if overrideAuth != nil {
		if overrideAuthMap, ok := overrideAuth.(map[string]interface{}); ok {
			if merged.Auth == nil {
				merged.Auth = &shared.AuthConfig{
					Settings: make(map[string]interface{}),
				}
			}
			// Set auth fields from map
			if authType, ok := overrideAuthMap["type"].(string); ok {
				at, err := shared.ParseAuthType(authType)
				if err != nil {
					return nil, fmt.Errorf("invalid auth type %s in %s", authType, sourceName)
				}
				merged.Auth.Type = at
			}
			if redirectURL, ok := overrideAuthMap["redirect_url"].(string); ok {
				merged.Auth.RedirectURL = redirectURL
			}
			if roles, ok := overrideAuthMap["roles"].([]string); ok {
				merged.Auth.Roles = roles
			}
			// Merge other settings
			for key, value := range overrideAuthMap {
				if key != "type" && key != "redirect_url" && key != "roles" {
					merged.Auth.Settings[key] = value
				}
			}
		}
	}

	// Merge i18n mappings (flat mappings)
	overrideI18nMappings := override.GetI18nMappings()
	if overrideI18nMappings != nil {
		if merged.I18n == nil {
			merged.I18n = &shared.I18nConfig{
				FlatMappings: make(map[string]string),
				Translations: make(map[string]*shared.LocaleTranslations),
			}
		}
		for key, value := range overrideI18nMappings {
			merged.I18n.FlatMappings[key] = value
		}
	}

	// Merge error settings
	overrideError := override.GetErrorSettings()
	if overrideError != nil {
		if overrideErrorMap, ok := overrideError.(map[string]interface{}); ok {
			if merged.Error == nil {
				merged.Error = &shared.ErrorConfig{
					Settings: make(map[string]interface{}),
				}
			}
			if template, ok := overrideErrorMap["template"].(string); ok {
				merged.Error.Template = template
			}
			// Merge other settings
			for key, value := range overrideErrorMap {
				if key != "template" {
					merged.Error.Settings[key] = value
				}
			}
		}
	}

	// Merge dynamic settings
	overrideDynamic := override.GetDynamicSettings()
	if overrideDynamic != nil {
		if overrideDynamicMap, ok := overrideDynamic.(map[string]interface{}); ok {
			if merged.Dynamic == nil {
				merged.Dynamic = &shared.DynamicConfig{
					Rules:    make(map[string]*shared.ValidationRule),
					Settings: make(map[string]interface{}),
				}
			}
			// Handle rules separately if they exist
			if rulesData, ok := overrideDynamicMap["rules"].(map[string]interface{}); ok {
				for ruleName, ruleData := range rulesData {
					if ruleMap, ok := ruleData.(map[string]interface{}); ok {
						rule := &shared.ValidationRule{
							Name:     ruleName,
							Settings: make(map[string]interface{}),
						}
						if ruleType, ok := ruleMap["type"].(string); ok {
							rule.Type = ruleType
						}
						if required, ok := ruleMap["required"].(bool); ok {
							rule.Required = required
						}
						if pattern, ok := ruleMap["pattern"].(string); ok {
							rule.Pattern = pattern
						}
						if defaultValue, ok := ruleMap["default"]; ok {
							rule.Default = defaultValue
						}
						// Merge other rule settings
						for key, value := range ruleMap {
							if key != "name" && key != "type" && key != "required" && key != "pattern" && key != "default" {
								rule.Settings[key] = value
							}
						}
						merged.Dynamic.Rules[ruleName] = rule
					}
				}
			}
			// Merge other settings
			for key, value := range overrideDynamicMap {
				if key != "rules" {
					merged.Dynamic.Settings[key] = value
				}
			}
		}
	}

	tm.logger.Debug("Merged config with override",
		zap.String("source", sourceName),
		zap.Int("base_metadata_count", len(tm.getRouteMetadata(base))),
		zap.Int("override_metadata_count", len(tm.getRouteMetadata(override))),
		zap.Int("merged_metadata_count", len(tm.getRouteMetadata(merged))))

	return merged, nil
}

// cloneConfig creates a deep copy of a config
func (tm *templateMiddleware) cloneConfig(original *shared.ConfigFile) *shared.ConfigFile {
	clone := &shared.ConfigFile{
		FilePath:         original.FilePath,
		TemplateFilePath: original.TemplateFilePath,
	}

	// Deep copy Metadata
	if original.Metadata != nil {
		clone.Metadata = &shared.MetadataConfig{
			Title:       original.Metadata.Title,
			Description: original.Metadata.Description,
			Author:      original.Metadata.Author,
			Version:     original.Metadata.Version,
		}
		if len(original.Metadata.Keywords) > 0 {
			clone.Metadata.Keywords = make([]string, len(original.Metadata.Keywords))
			copy(clone.Metadata.Keywords, original.Metadata.Keywords)
		}
		if original.Metadata.Custom != nil {
			clone.Metadata.Custom = make(map[string]interface{})
			for key, value := range original.Metadata.Custom {
				clone.Metadata.Custom[key] = value
			}
		}
	}

	// Deep copy I18n
	if original.I18n != nil {
		clone.I18n = &shared.I18nConfig{
			FlatMappings: make(map[string]string),
			Translations: make(map[string]*shared.LocaleTranslations),
		}
		for key, value := range original.I18n.FlatMappings {
			clone.I18n.FlatMappings[key] = value
		}
		for locale, translations := range original.I18n.Translations {
			clone.I18n.Translations[locale] = &shared.LocaleTranslations{
				Locale:       translations.Locale,
				Translations: make(map[string]interface{}),
			}
			for key, value := range translations.Translations {
				clone.I18n.Translations[locale].Translations[key] = value
			}
		}
	}

	// Deep copy Auth
	if original.Auth != nil {
		clone.Auth = &shared.AuthConfig{
			Type:        original.Auth.Type,
			RedirectURL: original.Auth.RedirectURL,
			Settings:    make(map[string]interface{}),
		}
		if len(original.Auth.Roles) > 0 {
			clone.Auth.Roles = make([]string, len(original.Auth.Roles))
			copy(clone.Auth.Roles, original.Auth.Roles)
		}
		for key, value := range original.Auth.Settings {
			clone.Auth.Settings[key] = value
		}
	}

	// Deep copy Layout
	if original.Layout != nil {
		clone.Layout = &shared.LayoutConfig{
			Template: original.Layout.Template,
			Settings: make(map[string]interface{}),
		}
		for key, value := range original.Layout.Settings {
			clone.Layout.Settings[key] = value
		}
	}

	// Deep copy Error
	if original.Error != nil {
		clone.Error = &shared.ErrorConfig{
			Template: original.Error.Template,
			Settings: make(map[string]interface{}),
		}
		for key, value := range original.Error.Settings {
			clone.Error.Settings[key] = value
		}
	}

	// Deep copy Dynamic
	if original.Dynamic != nil {
		clone.Dynamic = &shared.DynamicConfig{
			Rules:    make(map[string]*shared.ValidationRule),
			Settings: make(map[string]interface{}),
		}
		for ruleName, rule := range original.Dynamic.Rules {
			clone.Dynamic.Rules[ruleName] = &shared.ValidationRule{
				Name:     rule.Name,
				Type:     rule.Type,
				Required: rule.Required,
				Pattern:  rule.Pattern,
				Default:  rule.Default,
				Settings: make(map[string]interface{}),
			}
			for key, value := range rule.Settings {
				clone.Dynamic.Rules[ruleName].Settings[key] = value
			}
		}
		for key, value := range original.Dynamic.Settings {
			clone.Dynamic.Settings[key] = value
		}
	}

	return clone
}

// Helper functions for safe interface{} access

// getRouteMetadata safely extracts route metadata as map[string]interface{}
func (tm *templateMiddleware) getRouteMetadata(config *shared.ConfigFile) map[string]interface{} {
	if config == nil || config.Metadata == nil {
		return make(map[string]interface{})
	}

	result := make(map[string]interface{})

	// Add standard fields
	if config.Metadata.Title != "" {
		result["title"] = config.Metadata.Title
	}
	if config.Metadata.Description != "" {
		result["description"] = config.Metadata.Description
	}
	if len(config.Metadata.Keywords) > 0 {
		result["keywords"] = config.Metadata.Keywords
	}
	if config.Metadata.Author != "" {
		result["author"] = config.Metadata.Author
	}
	if config.Metadata.Version != "" {
		result["version"] = config.Metadata.Version
	}

	// Add custom fields
	for key, value := range config.Metadata.Custom {
		result[key] = value
	}

	return result
}

// isComponentRoute determines if a route is a component route using registry
// This is generic - components can be in any route, not just /components/
func (tm *templateMiddleware) isComponentRoute(routePath string) bool {
	// Get the route mapping to check if this path corresponds to a component template
	routeMapping := tm.templateRegistry.GetRouteToTemplateMapping()

	templateKey, exists := routeMapping[routePath]
	if !exists {
		return false
	}

	// Check if the template is a component type
	if metadata, err := tm.templateRegistry.GetTemplateMetadata(templateKey); err == nil {
		return metadata.Type == "component"
	}

	return false
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
