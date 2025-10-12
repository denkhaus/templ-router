package services

import (
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"go.uber.org/zap"
)

// TemplateKeyResolver handles template key resolution and mapping
type TemplateKeyResolver struct {
	logger           *zap.Logger
	templateRegistry interfaces.TemplateRegistry
}

// NewTemplateKeyResolver creates a new template key resolver
func NewTemplateKeyResolver(logger *zap.Logger, templateRegistry interfaces.TemplateRegistry) *TemplateKeyResolver {
	return &TemplateKeyResolver{
		logger:           logger,
		templateRegistry: templateRegistry,
	}
}

// ResolveTemplateKey resolves the correct template key for a given route and template file
func (tkr *TemplateKeyResolver) ResolveTemplateKey(routePath, templateFile string) (string, bool) {
	// Strategy 1: Direct route lookup in RouteToTemplate
	// Get route-to-template mapping from template registry
	routeMapping := tkr.templateRegistry.GetRouteToTemplateMapping()

	// Try direct route mapping first
	if templateUUID, exists := routeMapping[routePath]; exists {
		tkr.logger.Debug("Template key found via direct route mapping",
			zap.String("route", routePath),
			zap.String("template_uuid", templateUUID))
		return templateUUID, true
	}

	// Strategy 2: Try alternative route patterns for dynamic routes
	alternativeRoutes := tkr.generateAlternativeRoutes(routePath)
	for _, altRoute := range alternativeRoutes {
		// Check alternative route in template registry
		if templateUUID, exists := routeMapping[altRoute]; exists {
			tkr.logger.Debug("Template key found via alternative route mapping",
				zap.String("original_route", routePath),
				zap.String("matched_route", altRoute),
				zap.String("template_uuid", templateUUID))
			return templateUUID, true
		}
	}

	// Strategy 3: Try to construct key from template file path
	constructedKey := tkr.constructKeyFromTemplatePath(templateFile)
	// Constructed key resolution via template registry interface
	if constructedKey != "" {
		tkr.logger.Debug("Template key found via constructed path",
			zap.String("template_file", templateFile),
			zap.String("constructed_key", constructedKey))
		return constructedKey, true
	}

	tkr.logger.Debug("Template key not found",
		zap.String("route", routePath),
		zap.String("template_file", templateFile),
		zap.Strings("tried_alternatives", alternativeRoutes))

	return "", false
}

// generateAlternativeRoutes creates alternative route patterns for dynamic routes
func (tkr *TemplateKeyResolver) generateAlternativeRoutes(routePath string) []string {
	alternativeRoutes := []string{routePath}

	// Convert detected locale patterns to $locale pattern
	pathParts := strings.Split(strings.TrimPrefix(routePath, "/"), "/")
	if len(pathParts) > 0 && len(pathParts[0]) == 2 {
		// Assume first 2-letter segment is a locale
		localePattern := "/" + pathParts[0]
		alternativeRoutes = append(alternativeRoutes,
			strings.ReplaceAll(routePath, localePattern, "/$locale"))
	}

	// Convert numeric patterns to $id (configuration-agnostic)
	for i, segment := range pathParts {
		if segment != "" && tkr.isNumeric(segment) {
			// Create alternative with $id
			altSegments := make([]string, len(pathParts))
			copy(altSegments, pathParts)
			altSegments[i] = "$id"
			alternativeRoutes = append(alternativeRoutes, "/"+strings.Join(altSegments, "/"))
		}
	}

	return alternativeRoutes
}

// constructKeyFromTemplatePath constructs a route key from template file path
func (tkr *TemplateKeyResolver) constructKeyFromTemplatePath(templateFile string) string {
	// FAIL FAST: This function MUST have config access
	// No fallbacks, no hardcoded values - proper DI required
	panic("template_key_resolver.go: Config injection required - no hardcoded fallbacks allowed")
}

// isNumeric checks if a string contains only digits
func (tkr *TemplateKeyResolver) isNumeric(s string) bool {
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
