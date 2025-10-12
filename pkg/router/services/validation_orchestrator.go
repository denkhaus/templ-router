package services

import (
	"fmt"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// validationOrchestrator coordinates all validation services (replaces UnifiedValidationService)
type validationOrchestrator struct {
	logger         *zap.Logger
	config         interfaces.ConfigService
	routeValidator RouteValidator
	paramValidator *ParameterValidator
	authValidator  *AuthValidator
}

// NewValidationOrchestrator creates a new validation orchestrator for DI
func NewValidationOrchestrator(i do.Injector) (interfaces.ValidationService, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	config := do.MustInvoke[interfaces.ConfigService](i)

	// Create specialized validators
	routeValidator, err := NewRouteValidator(i)
	if err != nil {
		return nil, fmt.Errorf("failed to create route validator: %w", err)
	}

	paramValidator, err := NewParameterValidator(i)
	if err != nil {
		return nil, fmt.Errorf("failed to create parameter validator: %w", err)
	}

	authValidator, err := NewAuthValidator(i)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth validator: %w", err)
	}

	return &validationOrchestrator{
		logger:         logger,
		config:         config,
		routeValidator: routeValidator,
		paramValidator: paramValidator,
		authValidator:  authValidator,
	}, nil
}

// ValidateConfiguration validates the complete configuration (implements interfaces.ValidationService)
func (vo *validationOrchestrator) ValidateConfiguration(routes []interfaces.Route, configs map[string]*interfaces.ConfigFile) error {
	vo.logger.Info("Starting configuration validation",
		zap.Int("routes", len(routes)),
		zap.Int("configs", len(configs)))

	result := vo.validateAll(routes, configs)

	// Log validation results
	vo.logValidationResults(result)

	// Return error if there are validation errors
	if result.HasErrors() {
		return fmt.Errorf("validation failed with %d errors and %d warnings",
			result.GetErrorCount(), result.GetWarningCount())
	}

	vo.logger.Info("Configuration validation completed successfully",
		zap.Int("warnings", result.GetWarningCount()))

	return nil
}

// validateAll performs comprehensive validation using specialized validators
func (vo *validationOrchestrator) validateAll(routes []interfaces.Route, configs map[string]*interfaces.ConfigFile) *ValidationResult {
	result := &ValidationResult{
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationWarning, 0),
	}

	// Build hierarchy map for parameter inheritance validation
	hierarchyMap := vo.buildHierarchyMap(routes)

	// Validate each route individually
	for i := range routes {
		route := &routes[i]
		config := configs[route.TemplateFile]
		vo.validateSingleRoute(route, config, hierarchyMap, result)
	}

	// Validate route conflicts (cross-route validation)
	vo.routeValidator.ValidateRouteConflicts(routes, result)

	return result
}

// validateSingleRoute validates a single route using all specialized validators
func (vo *validationOrchestrator) validateSingleRoute(route *interfaces.Route, config *interfaces.ConfigFile, hierarchyMap map[string][]string, result *ValidationResult) {
	vo.logger.Debug("Validating route",
		zap.String("path", route.Path),
		zap.String("template", route.TemplateFile))

	// Route-level validation
	vo.routeValidator.ValidateTemplateFileExists(route, result)
	vo.routeValidator.ValidateRouteConfig(route, config, result)

	// Parameter validation
	vo.paramValidator.ValidateParameters(route, config, hierarchyMap, result)

	// Auth validation
	vo.authValidator.ValidateAuthSettings(route, config, result)
}

// buildHierarchyMap builds a map of route hierarchies for inheritance validation
func (vo *validationOrchestrator) buildHierarchyMap(routes []interfaces.Route) map[string][]string {
	hierarchyMap := make(map[string][]string)

	for _, route := range routes {
		parents := vo.findParentRoutes(route.Path, routes)
		if len(parents) > 0 {
			hierarchyMap[route.Path] = parents
		}
	}

	return hierarchyMap
}

// findParentRoutes finds potential parent routes for a given route path
func (vo *validationOrchestrator) findParentRoutes(routePath string, routes []interfaces.Route) []string {
	var parents []string

	// Split path into segments
	segments := strings.Split(strings.Trim(routePath, "/"), "/")

	// Check for parent paths (shorter paths that could be parents)
	for i := len(segments) - 1; i > 0; i-- {
		parentPath := "/" + strings.Join(segments[:i], "/")

		// Check if this parent path exists in routes
		for _, route := range routes {
			if route.Path == parentPath {
				parents = append(parents, parentPath)
				break
			}
		}
	}

	return parents
}

// logValidationResults logs all validation results
func (vo *validationOrchestrator) logValidationResults(result *ValidationResult) {
	for _, err := range result.Errors {
		vo.logger.Error("Validation error",
			zap.String("type", err.Type),
			zap.String("message", err.Message),
			zap.String("file", err.FilePath),
			zap.String("route", err.RoutePath),
			zap.Strings("suggestions", err.Suggestions))
	}

	for _, warn := range result.Warnings {
		vo.logger.Warn("Validation warning",
			zap.String("type", warn.Type),
			zap.String("message", warn.Message),
			zap.String("file", warn.FilePath),
			zap.String("route", warn.RoutePath),
			zap.Strings("suggestions", warn.Suggestions))
	}
}
