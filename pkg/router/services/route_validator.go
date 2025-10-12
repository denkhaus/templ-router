package services

import (
	"fmt"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

type RouteValidator interface {
	ValidateTemplateFileExists(route *interfaces.Route, result *ValidationResult)
	ValidateRouteConfig(route *interfaces.Route, config *interfaces.ConfigFile, result *ValidationResult)
	ValidateRouteConflicts(routes []interfaces.Route, result *ValidationResult)
}

// routeValidator handles route-specific validation logic
type routeValidator struct {
	logger     *zap.Logger
	config     interfaces.ConfigService
	fileSystem middleware.FileSystemChecker
}

// NewRouteValidator creates a new route validator for DI
func NewRouteValidator(i do.Injector) (RouteValidator, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	config := do.MustInvoke[interfaces.ConfigService](i)
	fileSystem := do.MustInvoke[middleware.FileSystemChecker](i)

	return &routeValidator{
		logger:     logger,
		config:     config,
		fileSystem: fileSystem,
	}, nil
}

// ValidateTemplateFileExists checks if template file exists in filesystem
func (rv *routeValidator) ValidateTemplateFileExists(route *interfaces.Route, result *ValidationResult) {
	if route.TemplateFile == "" {
		result.Errors = append(result.Errors, ValidationError{
			Type:      "MISSING_TEMPLATE_FILE",
			Message:   "Route has no template file specified",
			RoutePath: route.Path,
			FilePath:  "",
		})
		return
	}

	// Check if file exists
	if !rv.fileSystem.FileExists(route.TemplateFile) {
		result.Errors = append(result.Errors, ValidationError{
			Type:      "TEMPLATE_FILE_NOT_FOUND",
			Message:   fmt.Sprintf("Template file does not exist: %s", route.TemplateFile),
			RoutePath: route.Path,
			FilePath:  route.TemplateFile,
		})
		return
	}

	rv.logger.Debug("Template file exists",
		zap.String("route", route.Path),
		zap.String("template", route.TemplateFile))
}

// ValidateRouteConfig validates route configuration settings
func (rv *routeValidator) ValidateRouteConfig(route *interfaces.Route, config *interfaces.ConfigFile, result *ValidationResult) {
	if config == nil {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:      "MISSING_CONFIG",
			Message:   "Route has no configuration file",
			RoutePath: route.Path,
			FilePath:  "",
		})
		return
	}

	// Validate route-specific settings
	if config.RouteMetadata != nil {
		rv.validateRouteSettings(route, config, result)
	}

	rv.logger.Debug("Route config validated",
		zap.String("route", route.Path),
		zap.Bool("has_config", config != nil))
}

// ValidateRouteConflicts checks for conflicting routes
func (rv *routeValidator) ValidateRouteConflicts(routes []interfaces.Route, result *ValidationResult) {
	routeMap := make(map[string]*interfaces.Route)

	for i := range routes {
		route := &routes[i]
		normalizedPath := rv.normalizeRoutePath(route.Path)

		if existingRoute, exists := routeMap[normalizedPath]; exists {
			result.Errors = append(result.Errors, ValidationError{
				Type:      "ROUTE_CONFLICT",
				Message:   fmt.Sprintf("Route path conflicts with existing route: %s", existingRoute.Path),
				RoutePath: route.Path,
				FilePath:  route.TemplateFile,
				Suggestions: []string{
					"Use different route paths",
					"Consider using dynamic parameters",
				},
			})
		} else {
			routeMap[normalizedPath] = route
		}
	}

	// Check for dynamic route conflicts
	rv.validateDynamicRouteConflicts(routes, result)
}

// validateRouteSettings validates specific route configuration settings
func (rv *routeValidator) validateRouteSettings(route *interfaces.Route, config *interfaces.ConfigFile, result *ValidationResult) {
	// Basic route metadata validation
	if config.RouteMetadata == nil {
		rv.logger.Debug("No route metadata found", zap.String("route", route.Path))
		return
	}

	// Additional route-specific validations can be added here
	rv.logger.Debug("Route settings validated",
		zap.String("route", route.Path),
		zap.Bool("has_metadata", config.RouteMetadata != nil))
}

// normalizeRoutePath normalizes a route path for comparison
func (rv *routeValidator) normalizeRoutePath(path string) string {
	// Remove leading/trailing slashes and normalize
	normalized := strings.Trim(path, "/")
	if normalized == "" {
		return "/"
	}
	return "/" + normalized
}

// validateDynamicRouteConflicts checks for ambiguous dynamic routes
func (rv *routeValidator) validateDynamicRouteConflicts(routes []interfaces.Route, result *ValidationResult) {
	for i := 0; i < len(routes); i++ {
		for j := i + 1; j < len(routes); j++ {
			if rv.routesAreAmbiguous(routes[i].Path, routes[j].Path) {
				result.Errors = append(result.Errors, ValidationError{
					Type:      "AMBIGUOUS_ROUTES",
					Message:   fmt.Sprintf("Routes are ambiguous: %s and %s", routes[i].Path, routes[j].Path),
					RoutePath: routes[i].Path,
					FilePath:  routes[i].TemplateFile,
					Suggestions: []string{
						"Make route patterns more specific",
						"Use different parameter names",
						"Add static path segments",
					},
				})
			}
		}
	}
}

// routesAreAmbiguous checks if two routes could conflict
func (rv *routeValidator) routesAreAmbiguous(path1, path2 string) bool {
	parts1 := strings.Split(strings.Trim(path1, "/"), "/")
	parts2 := strings.Split(strings.Trim(path2, "/"), "/")

	if len(parts1) != len(parts2) {
		return false
	}

	for i := 0; i < len(parts1); i++ {
		p1, p2 := parts1[i], parts2[i]

		// If both are static and different, no conflict
		if !strings.HasPrefix(p1, "$") && !strings.HasPrefix(p2, "$") && p1 != p2 {
			return false
		}
	}

	return true
}
