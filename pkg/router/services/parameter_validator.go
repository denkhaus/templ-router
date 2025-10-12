package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// ParameterValidator handles parameter-specific validation logic
type ParameterValidator struct {
	logger *zap.Logger
	config interfaces.ConfigService
}

// NewParameterValidator creates a new parameter validator for DI
func NewParameterValidator(i do.Injector) (*ParameterValidator, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	config := do.MustInvoke[interfaces.ConfigService](i)

	return &ParameterValidator{
		logger: logger,
		config: config,
	}, nil
}

// ValidateParameters validates all parameters for a route
func (pv *ParameterValidator) ValidateParameters(route *interfaces.Route, config *interfaces.ConfigFile, hierarchyMap map[string][]string, result *ValidationResult) {
	if config == nil || config.DynamicSettings == nil {
		return
	}

	dirParams := pv.extractDirectoryParameters(route.Path)

	// Validate each configured parameter
	for paramName, paramConfig := range config.DynamicSettings.Parameters {
		pv.ValidateSingleParameter(paramName, paramConfig, route, config, dirParams, result)
	}

	// Check for missing parameters
	pv.ValidateMissingParameters(route, config, dirParams, result)

	// Check parameter inheritance
	pv.ValidateParameterInheritance(route, config, hierarchyMap, result)
}

// ValidateSingleParameter validates a single parameter configuration
func (pv *ParameterValidator) ValidateSingleParameter(
	paramName string,
	paramConfig *interfaces.DynamicParameterConfig,
	route *interfaces.Route,
	config *interfaces.ConfigFile,
	dirParams []string,
	result *ValidationResult,
) {
	// Check if parameter exists in route path
	if !pv.parameterExistsInRoute(paramName, route.Path, dirParams) {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:      "UNUSED_PARAMETER",
			Message:   fmt.Sprintf("Parameter '%s' is configured but not used in route path", paramName),
			RoutePath: route.Path,
			FilePath:  route.TemplateFile,
			Suggestions: []string{
				fmt.Sprintf("Add $%s to route path", paramName),
				"Remove parameter from configuration",
			},
		})
	}

	// Validate parameter regex if provided
	if paramConfig.Validation != "" {
		if err := pv.validateParameterRegex(paramConfig.Validation); err != nil {
			result.Errors = append(result.Errors, ValidationError{
				Type:      "INVALID_PARAMETER_REGEX",
				Message:   fmt.Sprintf("Invalid regex for parameter '%s': %v", paramName, err),
				RoutePath: route.Path,
				FilePath:  route.TemplateFile,
			})
		}
	}

	// Log parameter validation completion
	pv.logger.Debug("Parameter validation completed",
		zap.String("parameter", paramName),
		zap.String("route", route.Path))

}

// ValidateMissingParameters checks for parameters in route path that lack configuration
func (pv *ParameterValidator) ValidateMissingParameters(route *interfaces.Route, config *interfaces.ConfigFile, dirParams []string, result *ValidationResult) {
	configuredParams := make(map[string]bool)
	if config != nil && config.DynamicSettings != nil {
		for paramName := range config.DynamicSettings.Parameters {
			configuredParams[paramName] = true
		}
	}

	// Check directory parameters
	for _, param := range dirParams {
		if !configuredParams[param] {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:      "MISSING_PARAMETER_CONFIG",
				Message:   fmt.Sprintf("Parameter '$%s' in route path lacks configuration", param),
				RoutePath: route.Path,
				FilePath:  route.TemplateFile,
				Suggestions: []string{
					fmt.Sprintf("Add configuration for parameter '%s'", param),
					"Define validation rules",
					"Set parameter type",
				},
			})
		}
	}
}

// ValidateParameterInheritance checks parameter inheritance consistency
func (pv *ParameterValidator) ValidateParameterInheritance(route *interfaces.Route, config *interfaces.ConfigFile, hierarchyMap map[string][]string, result *ValidationResult) {
	if config == nil || config.DynamicSettings == nil {
		return
	}

	// Get parent routes
	parentRoutes := hierarchyMap[route.Path]

	for _, parentPath := range parentRoutes {
		// Check if parent has conflicting parameter definitions
		pv.checkParameterInheritanceConflicts(route, parentPath, config, result)
	}
}

// extractDirectoryParameters extracts parameter names from route path
func (pv *ParameterValidator) extractDirectoryParameters(routePath string) []string {
	var params []string
	parts := strings.Split(routePath, "/")

	for _, part := range parts {
		if strings.HasPrefix(part, "$") {
			paramName := strings.TrimPrefix(part, "$")
			if paramName != "" {
				params = append(params, paramName)
			}
		}
	}

	return params
}

// parameterExistsInRoute checks if a parameter exists in the route path
func (pv *ParameterValidator) parameterExistsInRoute(paramName, routePath string, dirParams []string) bool {
	// Check in directory parameters
	for _, param := range dirParams {
		if param == paramName {
			return true
		}
	}

	// Check for $paramName in path
	return strings.Contains(routePath, "$"+paramName)
}

// validateParameterRegex validates a parameter regex pattern
func (pv *ParameterValidator) validateParameterRegex(pattern string) error {
	_, err := regexp.Compile(pattern)
	return err
}

// checkParameterInheritanceConflicts checks for parameter inheritance conflicts
func (pv *ParameterValidator) checkParameterInheritanceConflicts(route *interfaces.Route, parentPath string, config *interfaces.ConfigFile, result *ValidationResult) {
	// TODO: Implementat checking parameter inheritance conflicts
	// This would compare parameter definitions between parent and child routes
	pv.logger.Debug("Checking parameter inheritance",
		zap.String("route", route.Path),
		zap.String("parent", parentPath))
}
