package services

import (
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
		} else if part == "{locale}" || (len(part) == 2 && (part == "en" || part == "de" || part == "fr" || part == "es")) {
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
// Uses the routeMapping to find the original template path by regenerating keys
// Template key format: "app/components/footer.templ#Footer" or hash format
// Returns: "app/components/footer.templ"
func (rd *routeDiscoveryImpl) extractTemplateFileFromRegistryKey(templateKey, routePattern string) string {
	// First try to split by # to separate the path from the template name
	parts := strings.Split(templateKey, "#")
	if len(parts) >= 2 && parts[0] != "" {
		templatePath := parts[0]
		rd.logger.Debug("Extracted template file path from registry key",
			zap.String("template_key", templateKey),
			zap.String("template_file", templatePath))
		return templatePath
	}

	// If it's a hash format, use the routeMapping to find the original template path
	// We know the routePattern maps to this templateKey, so we can try to reconstruct
	// the original template path by testing common patterns
	routeMapping := rd.templateRegistry.GetRouteToTemplateMapping()

	// Look for the current route in the mapping to confirm our templateKey
	if mappedKey, exists := routeMapping[routePattern]; exists && mappedKey == templateKey {
		// Try different template name patterns for this route
		candidateTemplateNames := []string{"Page", "Footer", "Navbar", "Error", "Layout"}

		for _, templateName := range candidateTemplateNames {
			// Generate candidate template path based on route pattern
			candidatePath := rd.buildCandidatePathFromRoute(routePattern, templateName)
			if candidatePath != "" {
				// Generate the key this path would create
				candidateKey := shared.GenerateTemplateKey(candidatePath + "#" + templateName)

				rd.logger.Debug("Testing candidate template path",
					zap.String("route_pattern", routePattern),
					zap.String("candidate_path", candidatePath),
					zap.String("template_name", templateName),
					zap.String("expected_key", candidateKey),
					zap.String("target_key", templateKey))

				if candidateKey == templateKey {
					rd.logger.Debug("Found matching template path",
						zap.String("template_key", templateKey),
						zap.String("route_pattern", routePattern),
						zap.String("matched_path", candidatePath))
					return candidatePath
				}
			}
		}
	}

	rd.logger.Warn("Could not extract template file path from registry key",
		zap.String("template_key", templateKey),
		zap.String("route_pattern", routePattern))
	return ""
}

// buildCandidatePathFromRoute builds a candidate template path from route pattern and template name
func (rd *routeDiscoveryImpl) buildCandidatePathFromRoute(routePattern, templateName string) string {
	// Remove leading slash and split into parts
	parts := strings.Split(strings.Trim(routePattern, "/"), "/")

	if len(parts) == 0 {
		return ""
	}

	// For component routes like "/components/footer"
	if strings.HasPrefix(routePattern, "/components/") && len(parts) >= 2 {
		// Component template path: "app/components/footer.templ"
		return "app/components/" + parts[1] + ".templ"
	}

	// For other routes, try standard patterns
	if len(parts) == 1 && parts[0] == "" {
		// Root route: "app/page.templ"
		return "app/page.templ"
	}

	if len(parts) == 1 {
		// Single part route like "login": "app/login/page.templ"
		return "app/" + parts[0] + "/page.templ"
	}

	if len(parts) >= 2 {
		// Multi-part route like "locale/dashboard": "app/locale_/dashboard/page.templ"
		return "app/" + strings.Join(parts[:len(parts)-1], "_") + "/" + parts[len(parts)-1] + "/page.templ"
	}

	return ""
}

// generateCandidateTemplatePath generates a candidate template file path from a route pattern
// This is a generic approach that learns from existing template registry patterns
func (rd *routeDiscoveryImpl) generateCandidateTemplatePath(routePattern string) string {
	// Analyze existing route mappings to understand the pattern
	routeMapping := rd.templateRegistry.GetRouteToTemplateMapping()

	// Find similar routes to learn the pattern
	for existingRoute, templateKey := range routeMapping {
		// Try to extract template path from non-hash keys
		if parts := strings.Split(templateKey, "#"); len(parts) >= 2 && parts[0] != "" {
			existingTemplatePath := parts[0]

			// Extract the pattern by comparing routes
			inferredPath := rd.inferTemplatePathFromExample(routePattern, existingRoute, existingTemplatePath)
			if inferredPath != "" {
				return inferredPath
			}
		}
	}

	// Fallback: use the existing generic path generation
	return rd.generateTemplateFilePathFromPattern(routePattern)
}

// inferTemplatePathFromExample infers a template path by comparing two routes
func (rd *routeDiscoveryImpl) inferTemplatePathFromExample(targetRoute, exampleRoute, exampleTemplatePath string) string {
	// Split both routes into parts
	targetParts := strings.Split(strings.Trim(targetRoute, "/"), "/")
	exampleParts := strings.Split(strings.Trim(exampleRoute, "/"), "/")

	// Split template path into parts
	templateParts := strings.Split(exampleTemplatePath, "/")

	// Find the mapping between route parts and template parts
	// This is a simplified pattern matching - replace corresponding parts
	if len(targetParts) == len(exampleParts) && len(templateParts) > 0 {
		// Replace the last route part with the target's last part
		if len(templateParts) >= 1 {
			templateParts[len(templateParts)-1] = targetParts[len(targetParts)-1] + rd.config.GetTemplateExtension()
			return strings.Join(templateParts, "/")
		}
	}

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
		} else if len(part) == 2 && (part == "en" || part == "de" || part == "fr" || part == "es") {
			// Handle locale
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
	rd.logger.Debug("Discovering layouts", zap.String("scan_path", scanPath))

	var layouts []interfaces.LayoutTemplate

	err := rd.fileSystem.WalkDirectory(scanPath, func(path string, isDir bool, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if isDir {
			return nil
		}

		// Only process layout templates
		if !rd.isLayoutTemplate(path) {
			return nil
		}

		layout, err := rd.createLayoutFromTemplate(path, scanPath)
		if err != nil {
			rd.logger.Warn("Failed to create layout from template",
				zap.String("template", path),
				zap.Error(err))
			return nil // Continue processing other files
		}

		layouts = append(layouts, layout)
		return nil
	})

	if err != nil {
		return nil, shared.NewRouteError("Failed to walk directory during layout discovery").
			WithCause(err).
			WithContext("scan_path", scanPath).
			WithContext("operation", "layout_discovery")
	}

	rd.logger.Info("Layout discovery completed",
		zap.String("scan_path", scanPath),
		zap.Int("layouts_found", len(layouts)))

	return layouts, nil
}

// DiscoverErrorTemplates implements router.RouteDiscovery
func (rd *routeDiscoveryImpl) DiscoverErrorTemplates(scanPath string) ([]interfaces.ErrorTemplate, error) {
	rd.logger.Debug("Discovering error templates", zap.String("scan_path", scanPath))

	var errorTemplates []interfaces.ErrorTemplate

	err := rd.fileSystem.WalkDirectory(scanPath, func(path string, isDir bool, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if isDir {
			return nil
		}

		// Only process error templates
		if !rd.isErrorTemplate(path) {
			return nil
		}

		errorTemplate, err := rd.createErrorTemplateFromTemplate(path, scanPath)
		if err != nil {
			rd.logger.Warn("Failed to create error template from template",
				zap.String("template", path),
				zap.Error(err))
			return nil // Continue processing other files
		}

		errorTemplates = append(errorTemplates, errorTemplate)
		return nil
	})

	if err != nil {
		return nil, shared.NewRouteError("Failed to walk directory during error template discovery").
			WithCause(err).
			WithContext("scan_path", scanPath).
			WithContext("operation", "error_template_discovery")
	}

	rd.logger.Info("Error template discovery completed",
		zap.String("scan_path", scanPath),
		zap.Int("error_templates_found", len(errorTemplates)))

	return errorTemplates, nil
}

// isLayoutTemplate checks if a template file is a layout template
func (rd *routeDiscoveryImpl) isLayoutTemplate(path string) bool {
	return strings.Contains(path, "layout.templ")
}

// isErrorTemplate checks if a template file is an error template
func (rd *routeDiscoveryImpl) isErrorTemplate(path string) bool {
	return strings.Contains(path, "error.templ")
}

// createLayoutFromTemplate creates a layout template from a template file
func (rd *routeDiscoveryImpl) createLayoutFromTemplate(templatePath, scanPath string) (interfaces.LayoutTemplate, error) {
	relativePath, err := filepath.Rel(scanPath, templatePath)
	if err != nil {
		return interfaces.LayoutTemplate{}, shared.NewRouteError("Failed to get relative path for layout template").
			WithCause(err).
			WithContext("template_path", templatePath).
			WithContext("scan_path", scanPath).
			WithContext("operation", "layout_template_creation")
	}

	// Calculate layout level based on directory depth
	layoutLevel := strings.Count(relativePath, string(filepath.Separator))

	layout := interfaces.LayoutTemplate{
		FilePath:      templatePath,
		DirectoryPath: filepath.Dir(templatePath),
		LayoutLevel:   layoutLevel,
	}

	return layout, nil
}

// createErrorTemplateFromTemplate creates an error template from a template file
func (rd *routeDiscoveryImpl) createErrorTemplateFromTemplate(templatePath, scanPath string) (interfaces.ErrorTemplate, error) {
	relativePath, err := filepath.Rel(scanPath, templatePath)
	if err != nil {
		return interfaces.ErrorTemplate{}, shared.NewRouteError("Failed to get relative path for error template").
			WithCause(err).
			WithContext("template_path", templatePath).
			WithContext("scan_path", scanPath).
			WithContext("operation", "error_template_creation")
	}

	// Extract error type from path (e.g., 404, 500, etc.)
	errorType := rd.extractErrorType(relativePath)

	errorTemplate := interfaces.ErrorTemplate{
		FilePath:        templatePath,
		DirectoryPath:   filepath.Dir(templatePath),
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
