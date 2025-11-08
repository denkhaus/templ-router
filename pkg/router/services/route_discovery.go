package services

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// routeDiscoveryImpl implements clean route discovery
type routeDiscoveryImpl struct {
	config           interfaces.ConfigService
	logger           *zap.Logger
	injector         do.Injector
	fileSystem       interfaces.FileSystemChecker
	templateRegistry interfaces.TemplateRegistry
}

// NewRouteDiscovery creates a new route discovery implementation for DI
func NewRouteDiscovery(i do.Injector) (interfaces.RouteDiscovery, error) {
	config := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	fileSystem := do.MustInvoke[interfaces.FileSystemChecker](i)
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	return &routeDiscoveryImpl{
		config:           config,
		logger:           logger,
		injector:         i,
		fileSystem:       fileSystem,
		templateRegistry: templateRegistry,
	}, nil
}

// DiscoverRoutes implements router.RouteDiscovery using generated template registry
func (rd *routeDiscoveryImpl) DiscoverRoutes(scanPath string) ([]interfaces.Route, error) {
	rd.logger.Debug("Discovering routes using generated templates", zap.String("scan_path", scanPath))

	var routes []interfaces.Route

	// Get route-to-template mapping from template registry
	routeMapping := rd.templateRegistry.GetRouteToTemplateMapping()

	rd.logger.Debug("Discovering routes from template registry",
		zap.String("scan_path", scanPath),
		zap.Int("route_mappings", len(routeMapping)))

	// Convert template registry mappings to router.Route objects
	for routePattern, templateKey := range routeMapping {
		// Verify template exists
		if !rd.templateRegistry.IsAvailable(templateKey) {
			rd.logger.Warn("Template not available for route",
				zap.String("route", routePattern),
				zap.String("template", templateKey))
			continue
		}

		// Check if template requires data service
		requiresDataService := rd.templateRegistry.RequiresDataService(templateKey)
		var dataServiceInterface string
		if requiresDataService {
			if dataServiceInfo, exists := rd.templateRegistry.GetDataServiceInfo(templateKey); exists {
				dataServiceInterface = dataServiceInfo.InterfaceType
			}
		}

		// Debug logging for DataService detection
		rd.logger.Info("Route discovery DataService check",
			zap.String("route", routePattern),
			zap.String("template_key", templateKey),
			zap.Bool("requires_data_service", requiresDataService),
			zap.String("data_service_interface", dataServiceInterface))

		// First try the original working approach
		templateFile := rd.generateTemplateFilePathFromPattern(routePattern)

		rd.logger.Debug("Template path generation",
			zap.String("route_pattern", routePattern),
			zap.String("generated_template_file", templateFile),
			zap.Bool("file_exists", rd.fileSystem.FileExists(templateFile)))

		// Only if the generated template file doesn't exist, try registry-based extraction
		if !rd.fileSystem.FileExists(templateFile) {
			rd.logger.Debug("Generated template file doesn't exist, trying registry extraction",
				zap.String("route_pattern", routePattern),
				zap.String("template_key", templateKey),
				zap.String("missing_file", templateFile))

			extractedFile := rd.extractTemplateFileFromRegistryKey(templateKey, routePattern)
			if extractedFile != "" {
				templateFile = extractedFile
				rd.logger.Info("Template path extraction successful - using extracted path",
					zap.String("route_pattern", routePattern),
					zap.String("template_key", templateKey),
					zap.String("original_generated", templateFile),
					zap.String("final_extracted", extractedFile))
			} else {
				rd.logger.Error("Both generated and extracted template files don't exist",
					zap.String("route_pattern", routePattern),
					zap.String("template_key", templateKey),
					zap.String("generated_file", templateFile))
			}
		} else {
			rd.logger.Debug("Using generated template file (file exists)",
				zap.String("route_pattern", routePattern),
				zap.String("template_file", templateFile))
		}

		// Create route object
		route := interfaces.Route{
			Path:                 routePattern,
			TemplateFile:         templateFile,
			IsDynamic:            strings.Contains(routePattern, "$"),
			Handler:              rd.generateHandlerName(routePattern),
			Precedence:           rd.calculateRoutePrecedence(routePattern),
			RequiresDataService:  requiresDataService,
			DataServiceInterface: dataServiceInterface,
		}

		routes = append(routes, route)

		rd.logger.Info("Route discovered from template registry",
			zap.String("pattern", routePattern),
			zap.String("template", templateKey),
			zap.String("file", route.TemplateFile),
			zap.Bool("dynamic", route.IsDynamic),
			zap.Bool("requires_data_service", route.RequiresDataService),
			zap.String("data_service_interface", route.DataServiceInterface))
	}

	rd.logger.Info("Route discovery completed using template registry",
		zap.String("scan_path", scanPath),
		zap.Int("routes_found", len(routes)))

	return routes, nil
}

// generateTemplateFilePathFromPattern generates a template file path from a route pattern
func (rd *routeDiscoveryImpl) generateTemplateFilePathFromPattern(routePattern string) string {
	// Get configurable template root directory
	templateRoot := rd.config.GetLayoutRootDirectory()
	templateExtension := rd.config.GetTemplateExtension()

	// Convert route pattern to file path
	// Example: "/en/dashboard" -> "app/locale_/dashboard/page.templ"
	// Example: "/en/user/$id" -> "app/locale_/user/id_/page.templ"

	// Remove leading slash and split into parts
	parts := strings.Split(strings.Trim(routePattern, "/"), "/")

	var pathParts []string
	for _, part := range parts {
		if part == "" {
			continue
		}

		// Handle dynamic parameters
		if strings.HasPrefix(part, "$") {
			// Convert $id to id_
			paramName := strings.TrimPrefix(part, "$")
			pathParts = append(pathParts, paramName+"_")
		} else if part == "{locale}" || rd.isLocaleCode(part) {
			// Handle locale parameters - both placeholder {locale} and actual locale codes
			pathParts = append(pathParts, "locale_")
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			// Handle other dynamic parameters like {id}
			paramName := strings.Trim(part, "{}")
			pathParts = append(pathParts, paramName+"_")
		} else {
			pathParts = append(pathParts, part)
		}
	}

	// Add page template file
	pathParts = append(pathParts, "page"+templateExtension)

	return filepath.Join(templateRoot, filepath.Join(pathParts...))
}

// extractTemplateFileFromRegistryKey extracts the template file path from a template registry key
// Uses the template metadata directly to get the correct template path
// Template key format: MD5 hash (e.g., "7210696de26402b63095cea9005a4b7c")
// Returns: "app/components/footer.templ"
func (rd *routeDiscoveryImpl) extractTemplateFileFromRegistryKey(templateKey, routePattern string) string {
	// Use the template metadata to get the correct template path
	if metadata, err := rd.templateRegistry.GetTemplateMetadata(templateKey); err == nil {
		if metadata.TemplatePath != "" {
			rd.logger.Debug("Extracted template file path from registry metadata",
				zap.String("template_key", templateKey),
				zap.String("route_pattern", routePattern),
				zap.String("template_file", metadata.TemplatePath))
			return metadata.TemplatePath
		}
	}

	rd.logger.Warn("Could not extract template file path from registry key",
		zap.String("template_key", templateKey),
		zap.String("route_pattern", routePattern))
	return ""
}

// generateHandlerName generates a handler name from a route pattern
func (rd *routeDiscoveryImpl) generateHandlerName(routePattern string) string {
	// Convert route pattern to handler name
	// Example: "/en/dashboard" -> "LocaleDashboardHandler"
	// Example: "/en/user/$id" -> "LocaleUserIdHandler"

	parts := strings.Split(strings.Trim(routePattern, "/"), "/")
	var handlerParts []string

	for _, part := range parts {
		if part == "" {
			continue
		}

		if strings.HasPrefix(part, "$") {
			// Convert $id to Id
			paramName := strings.TrimPrefix(part, "$")
			if len(paramName) > 0 {
				handlerParts = append(handlerParts, strings.ToUpper(paramName[:1])+paramName[1:])
			} else {
				handlerParts = append(handlerParts, "Id")
			}
		} else if rd.isLocaleCode(part) {
			// Handle locale using supported locales from config service
			handlerParts = append(handlerParts, "Locale")
		} else {
			// Regular part - capitalize
			if len(part) > 0 {
				handlerParts = append(handlerParts, strings.ToUpper(part[:1])+part[1:])
			} else {
				handlerParts = append(handlerParts, "Index")
			}
		}
	}

	return strings.Join(handlerParts, "") + "Handler"
}

// isLocaleCode checks if a given part matches any of the supported locales from config service
func (rd *routeDiscoveryImpl) isLocaleCode(part string) bool {
	supportedLocales := rd.config.GetSupportedLocales()
	for _, locale := range supportedLocales {
		if part == locale {
			return true
		}
	}
	return false
}

// calculateRoutePrecedence calculates route precedence for ordering
func (rd *routeDiscoveryImpl) calculateRoutePrecedence(routePattern string) int {
	// Static routes have higher precedence than dynamic routes
	// More specific routes have higher precedence

	precedence := 100
	parts := strings.Split(strings.Trim(routePattern, "/"), "/")

	for _, part := range parts {
		if strings.HasPrefix(part, "$") {
			// Dynamic parameter reduces precedence
			precedence -= 10
		} else {
			// Static part increases precedence
			precedence += 5
		}
	}

	return precedence
}

// DiscoverLayouts implements router.RouteDiscovery
func (rd *routeDiscoveryImpl) DiscoverLayouts(scanPath string) ([]interfaces.LayoutTemplate, error) {
	rd.logger.Debug("Discovering layouts using registry metadata", zap.String("scan_path", scanPath))

	var layouts []interfaces.LayoutTemplate

	// Get all template metadata from registry
	allMetadata := rd.templateRegistry.GetAllTemplateMetadata()

	for templateKey, metadata := range allMetadata {
		// Only process layout templates
		if metadata.Type != "layout" {
			continue
		}

		rd.logger.Debug("Found layout template from registry",
			zap.String("template_key", templateKey),
			zap.String("template_path", metadata.TemplatePath),
			zap.String("component_name", metadata.ComponentName))

		// Create layout from registry metadata
		layout, err := rd.createLayoutFromMetadata(metadata, scanPath)
		if err != nil {
			rd.logger.Warn("Failed to create layout from registry metadata",
				zap.String("template_key", templateKey),
				zap.String("template_path", metadata.TemplatePath),
				zap.Error(err))
			continue // Continue processing other templates
		}

		layouts = append(layouts, layout)
	}

	rd.logger.Info("Layout discovery completed using registry",
		zap.String("scan_path", scanPath),
		zap.Int("layouts_found", len(layouts)))

	return layouts, nil
}

// DiscoverErrorTemplates implements router.RouteDiscovery
func (rd *routeDiscoveryImpl) DiscoverErrorTemplates(scanPath string) ([]interfaces.ErrorTemplate, error) {
	rd.logger.Debug("Discovering error templates using registry metadata", zap.String("scan_path", scanPath))

	var errorTemplates []interfaces.ErrorTemplate

	// Get all template metadata from registry
	allMetadata := rd.templateRegistry.GetAllTemplateMetadata()

	for templateKey, metadata := range allMetadata {
		// Only process error templates
		if metadata.Type != interfaces.TemplateTypeError {
			continue
		}

		rd.logger.Debug("Found error template from registry",
			zap.String("template_key", templateKey),
			zap.String("template_path", metadata.TemplatePath),
			zap.String("component_name", metadata.ComponentName))

		// Create error template from registry metadata
		errorTemplate, err := rd.createErrorTemplateFromMetadata(metadata, scanPath)
		if err != nil {
			rd.logger.Warn("Failed to create error template from registry metadata",
				zap.String("template_key", templateKey),
				zap.String("template_path", metadata.TemplatePath),
				zap.Error(err))
			continue // Continue processing other templates
		}

		errorTemplates = append(errorTemplates, errorTemplate)
	}

	rd.logger.Info("Error template discovery completed using registry",
		zap.String("scan_path", scanPath),
		zap.Int("error_templates_found", len(errorTemplates)))

	return errorTemplates, nil
}

// createLayoutFromMetadata creates a layout template from registry metadata
func (rd *routeDiscoveryImpl) createLayoutFromMetadata(metadata *interfaces.TemplateMetadata, scanPath string) (interfaces.LayoutTemplate, error) {
	// Normalize template path - convert relative to absolute if needed
	templatePath := metadata.TemplatePath
	if !filepath.IsAbs(templatePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return interfaces.LayoutTemplate{}, shared.NewRouteError("Failed to get working directory for path normalization").
				WithCause(err).
				WithContext("template_path", metadata.TemplatePath).
				WithContext("operation", "layout_template_creation_from_metadata")
		}
		templatePath = filepath.Clean(filepath.Join(cwd, templatePath))
	}

	relativePath, err := filepath.Rel(scanPath, templatePath)
	if err != nil {
		return interfaces.LayoutTemplate{}, shared.NewRouteError("Failed to get relative path for layout template").
			WithCause(err).
			WithContext("template_path", templatePath).
			WithContext("scan_path", scanPath).
			WithContext("operation", "layout_template_creation_from_metadata")
	}

	// Calculate layout level based on directory depth
	layoutLevel := strings.Count(relativePath, string(filepath.Separator))

	layout := interfaces.LayoutTemplate{
		FilePath:      metadata.TemplatePath,
		DirectoryPath: filepath.Dir(metadata.TemplatePath),
		LayoutLevel:   layoutLevel,
	}

	return layout, nil
}

// createErrorTemplateFromMetadata creates an error template from registry metadata
func (rd *routeDiscoveryImpl) createErrorTemplateFromMetadata(metadata *interfaces.TemplateMetadata, scanPath string) (interfaces.ErrorTemplate, error) {
	// Normalize template path - convert relative to absolute if needed
	templatePath := metadata.TemplatePath
	if !filepath.IsAbs(templatePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return interfaces.ErrorTemplate{}, shared.NewRouteError("Failed to get working directory for path normalization").
				WithCause(err).
				WithContext("template_path", metadata.TemplatePath).
				WithContext("operation", "error_template_creation_from_metadata")
		}
		templatePath = filepath.Clean(filepath.Join(cwd, templatePath))
	}

	relativePath, err := filepath.Rel(scanPath, templatePath)
	if err != nil {
		return interfaces.ErrorTemplate{}, shared.NewRouteError("Failed to get relative path for error template").
			WithCause(err).
			WithContext("template_path", templatePath).
			WithContext("scan_path", scanPath).
			WithContext("operation", "error_template_creation_from_metadata")
	}

	// Extract error type from path (e.g., 404, 500, etc.)
	errorType := rd.extractErrorType(relativePath)

	errorTemplate := interfaces.ErrorTemplate{
		FilePath:        metadata.TemplatePath,
		DirectoryPath:   filepath.Dir(metadata.TemplatePath),
		ErrorTypes:      []string{errorType},
		PrecedenceLevel: strings.Count(relativePath, string(filepath.Separator)),
		ErrorMessages:   make(map[int]string),
	}

	return errorTemplate, nil
}

// extractErrorType extracts the error type from an error template path
func (rd *routeDiscoveryImpl) extractErrorType(templatePath string) string {
	// Look for numeric error codes in the path
	parts := strings.Split(templatePath, "/")
	for _, part := range parts {
		if strings.Contains(part, "error") {
			// Try to extract error code
			if strings.Contains(part, "404") {
				return "404"
			}
			if strings.Contains(part, "500") {
				return "500"
			}
			if strings.Contains(part, "403") {
				return "403"
			}
		}
	}
	return "generic"
}
