package services

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// routeDiscoveryImpl implements clean route discovery
type routeDiscoveryImpl struct {
	config           interfaces.ConfigService
	logger           *zap.Logger
	injector         do.Injector
	fileSystem       middleware.FileSystemChecker
	templateRegistry interfaces.TemplateRegistry
}

// NewRouteDiscovery creates a new route discovery implementation for DI
func NewRouteDiscovery(i do.Injector) (router.RouteDiscovery, error) {
	config := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	fileSystem := do.MustInvoke[middleware.FileSystemChecker](i)
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

		// Create route object
		route := interfaces.Route{
			Path:                 routePattern,
			TemplateFile:         rd.generateTemplateFilePathFromPattern(routePattern),
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
			handlerParts = append(handlerParts, strings.Title(paramName))
		} else if len(part) == 2 && (part == "en" || part == "de" || part == "fr" || part == "es") {
			// Handle locale
			handlerParts = append(handlerParts, "Locale")
		} else {
			handlerParts = append(handlerParts, strings.Title(part))
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
func (rd *routeDiscoveryImpl) DiscoverLayouts(scanPath string) ([]router.LayoutTemplate, error) {
	rd.logger.Debug("Discovering layouts", zap.String("scan_path", scanPath))

	var layouts []router.LayoutTemplate

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
func (rd *routeDiscoveryImpl) DiscoverErrorTemplates(scanPath string) ([]router.ErrorTemplate, error) {
	rd.logger.Debug("Discovering error templates", zap.String("scan_path", scanPath))

	var errorTemplates []router.ErrorTemplate

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
func (rd *routeDiscoveryImpl) createLayoutFromTemplate(templatePath, scanPath string) (router.LayoutTemplate, error) {
	relativePath, err := filepath.Rel(scanPath, templatePath)
	if err != nil {
		return router.LayoutTemplate{}, shared.NewRouteError("Failed to get relative path for layout template").
			WithCause(err).
			WithContext("template_path", templatePath).
			WithContext("scan_path", scanPath).
			WithContext("operation", "layout_template_creation")
	}

	// Calculate layout level based on directory depth
	layoutLevel := strings.Count(relativePath, string(filepath.Separator))

	layout := router.LayoutTemplate{
		FilePath:      templatePath,
		DirectoryPath: filepath.Dir(templatePath),
		LayoutLevel:   layoutLevel,
	}

	return layout, nil
}

// createErrorTemplateFromTemplate creates an error template from a template file
func (rd *routeDiscoveryImpl) createErrorTemplateFromTemplate(templatePath, scanPath string) (router.ErrorTemplate, error) {
	relativePath, err := filepath.Rel(scanPath, templatePath)
	if err != nil {
		return router.ErrorTemplate{}, shared.NewRouteError("Failed to get relative path for error template").
			WithCause(err).
			WithContext("template_path", templatePath).
			WithContext("scan_path", scanPath).
			WithContext("operation", "error_template_creation")
	}

	// Extract error type from path (e.g., 404, 500, etc.)
	errorType := rd.extractErrorType(relativePath)

	errorTemplate := router.ErrorTemplate{
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
