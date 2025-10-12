package services

import (
	"strings"

	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// RouteConverter defines the contract for route path conversions
type RouteConverter interface {
	GenerateTemplateKey(templateFile string) string
	ConvertLayoutPathToRoute(layoutPath string) string
	GenerateRouteVariations(routePath string) []string
}

// routeConverter handles route path conversions only (private implementation)
type routeConverter struct {
	logger *zap.Logger
}

// NewRouteConverter creates a new route converter for DI
func NewRouteConverter(i do.Injector) (RouteConverter, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	return &routeConverter{
		logger: logger,
	}, nil
}

// GenerateTemplateKey converts a template file path to a registry key
func (rc *routeConverter) GenerateTemplateKey(templateFile string) string {
	// Convert app/login/page.templ -> /login
	// Convert app/page.templ -> /
	// Convert app/locale_/admin/page.templ -> /$locale/admin
	// Convert app/locale_/page.templ -> /$locale

	// Remove root directory prefix and ".templ" suffix (library-agnostic)
	// Extract root directory from config or use first path segment
	parts := strings.Split(templateFile, "/")
	var key string
	if len(parts) > 1 {
		// Remove first directory (root) and rebuild path
		key = "/" + strings.Join(parts[1:], "/")
		key = strings.TrimSuffix(key, ".templ")
	} else {
		key = strings.TrimSuffix(templateFile, ".templ")
	}

	// Remove "/page" suffix for route mapping
	key = strings.TrimSuffix(key, "/page")

	// Convert locale_ to $locale for route matching
	key = strings.ReplaceAll(key, "/locale_", "/$locale")

	// Convert id_ to $id for dynamic route matching
	key = strings.ReplaceAll(key, "/id_", "/$id")

	// Special case: if empty, it's the root route
	if key == "" {
		key = "/"
	} else if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}

	rc.logger.Debug("Generated template key",
		zap.String("template_file", templateFile),
		zap.String("key", key))

	return key
}

// ConvertLayoutPathToRoute converts a layout file path to a route pattern (library-agnostic)
func (rc *routeConverter) ConvertLayoutPathToRoute(layoutPath string) string {
	// Extract the directory part and convert to route pattern
	// Examples (library-agnostic):
	// templates/layout.templ -> /layout
	// src/dashboard/layout.templ -> /dashboard/layout
	// components/layout.templ -> /layout

	// Remove file extension
	pathWithoutExt := strings.TrimSuffix(layoutPath, ".templ")

	// Find the root directory (first directory in path)
	parts := strings.Split(pathWithoutExt, "/")
	if len(parts) < 2 {
		// Single file in root
		return "/layout"
	}

	// Remove root directory and build route
	routeParts := parts[1:] // Skip first part (root directory)
	route := "/" + strings.Join(routeParts, "/")

	rc.logger.Debug("Converted layout path to route",
		zap.String("layout_path", layoutPath),
		zap.String("route", route),
		zap.Strings("parts", parts))

	return route
}

// GenerateRouteVariations creates alternative route patterns for dynamic routes
func (rc *routeConverter) GenerateRouteVariations(routePath string) []string {
	alternativeRoutes := []string{routePath}

	// Convert specific locale patterns to $locale pattern based on detected patterns
	// Extract locale from path if it matches common 2-letter patterns
	pathParts := strings.Split(strings.TrimPrefix(routePath, "/"), "/")
	if len(pathParts) > 0 && len(pathParts[0]) == 2 {
		// Assume first 2-letter segment is a locale
		localePattern := "/" + pathParts[0]
		alternativeRoutes = append(alternativeRoutes,
			strings.ReplaceAll(routePath, localePattern, "/$locale"))
	}

	// Convert numeric patterns to $id (configuration-agnostic)
	pathSegments := strings.Split(routePath, "/")
	for i, segment := range pathSegments {
		if segment != "" && rc.isNumeric(segment) {
			// Create alternative with $id
			altSegments := make([]string, len(pathSegments))
			copy(altSegments, pathSegments)
			altSegments[i] = "$id"
			alternativeRoutes = append(alternativeRoutes, strings.Join(altSegments, "/"))
		}
	}

	// Convert string patterns to $id for non-numeric dynamic routes
	for i, segment := range pathSegments {
		if segment != "" && !rc.isNumeric(segment) && len(pathSegments) > 2 {
			// Create alternative with $id for potential dynamic segments
			altSegments := make([]string, len(pathSegments))
			copy(altSegments, pathSegments)
			altSegments[i] = "$id"
			alternativeRoutes = append(alternativeRoutes, strings.Join(altSegments, "/"))
		}
	}

	rc.logger.Debug("Generated route variations",
		zap.String("original_route", routePath),
		zap.Strings("variations", alternativeRoutes))

	return alternativeRoutes
}

// isNumeric checks if a string contains only digits
func (rc *routeConverter) isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
