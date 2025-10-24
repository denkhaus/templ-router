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
func (pv *ParameterValidator) ValidateParameters(route *interfaces.Route, config *interfaces.ConfigFile, hierarchyMap map[string][]string, configs map[string]*interfaces.ConfigFile, result *ValidationResult) {
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
	pv.ValidateParameterInheritance(route, config, hierarchyMap, configs, result)
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
func (pv *ParameterValidator) ValidateParameterInheritance(route *interfaces.Route, config *interfaces.ConfigFile, hierarchyMap map[string][]string, configs map[string]*interfaces.ConfigFile, result *ValidationResult) {
	if config == nil || config.DynamicSettings == nil {
		return
	}

	// Get parent routes
	parentRoutes := hierarchyMap[route.Path]

	for _, parentPath := range parentRoutes {
		// Check if parent has conflicting parameter definitions
		pv.checkParameterInheritanceConflicts(route, parentPath, config, configs, result)
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
func (pv *ParameterValidator) checkParameterInheritanceConflicts(route *interfaces.Route, parentPath string, config *interfaces.ConfigFile, configs map[string]*interfaces.ConfigFile, result *ValidationResult) {
	// Find parent route configuration
	var parentConfig *interfaces.ConfigFile
	for templateFile, cfg := range configs {
		// Find the route that matches the parent path
		if pv.isConfigForRoute(templateFile, parentPath) {
			parentConfig = cfg
			break
		}
	}

	if parentConfig == nil || parentConfig.DynamicSettings == nil || parentConfig.DynamicSettings.Parameters == nil {
		// No parent configuration or no parameters to check
		return
	}

	if config.DynamicSettings == nil || config.DynamicSettings.Parameters == nil {
		// Child has no parameters, no conflicts possible
		return
	}

	// Check for parameter conflicts between parent and child
	for paramName, childParamConfig := range config.DynamicSettings.Parameters {
		if parentParamConfig, exists := parentConfig.DynamicSettings.Parameters[paramName]; exists {
			// Parameter exists in both parent and child - check for conflicts
			pv.validateParameterCompatibility(route, parentPath, paramName, childParamConfig, parentParamConfig, result)
		}
	}

	// Check for inherited parameters that should be available in child
	pv.validateParameterInheritanceAvailability(route, parentPath, parentConfig, config, result)

	pv.logger.Debug("Checked parameter inheritance",
		zap.String("route", route.Path),
		zap.String("parent", parentPath),
		zap.Int("parent_params", len(parentConfig.DynamicSettings.Parameters)),
		zap.Int("child_params", len(config.DynamicSettings.Parameters)))
}

// isConfigForRoute checks if a template file corresponds to a route path
func (pv *ParameterValidator) isConfigForRoute(templateFile, routePath string) bool {
	// Simple heuristic: check if the template file path contains the route path structure
	// This is a simplified approach - in a real implementation, you might need a more sophisticated mapping
	normalizedRoute := strings.Trim(routePath, "/")
	normalizedTemplate := strings.ReplaceAll(templateFile, "\\", "/")
	
	// Check if the route path segments appear in the template file path
	routeSegments := strings.Split(normalizedRoute, "/")
	for _, segment := range routeSegments {
		if segment != "" {
			// Handle dynamic parameters: $id becomes id_
			if strings.HasPrefix(segment, "$") {
				paramName := strings.TrimPrefix(segment, "$")
				expectedDir := paramName + "_"
				if !strings.Contains(normalizedTemplate, expectedDir) {
					return false
				}
			} else {
				// Regular segment
				if !strings.Contains(normalizedTemplate, segment) {
					return false
				}
			}
		}
	}
	return true
}

// validateParameterCompatibility checks if child parameter is compatible with parent parameter
func (pv *ParameterValidator) validateParameterCompatibility(
	route *interfaces.Route,
	parentPath string,
	paramName string,
	childConfig *interfaces.DynamicParameterConfig,
	parentConfig *interfaces.DynamicParameterConfig,
	result *ValidationResult,
) {
	// Check validation regex compatibility
	if childConfig.Validation != "" && parentConfig.Validation != "" {
		if childConfig.Validation != parentConfig.Validation {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:      "PARAMETER_VALIDATION_CONFLICT",
				Message:   fmt.Sprintf("Parameter '%s' has different validation regex in child route ('%s') than parent route ('%s')", paramName, childConfig.Validation, parentConfig.Validation),
				RoutePath: route.Path,
				FilePath:  route.TemplateFile,
				Suggestions: []string{
					"Use the same validation regex as parent route",
					"Make child validation more restrictive than parent",
					"Document why different validation is needed",
				},
			})
		}
	}

	// Check supported values compatibility
	if len(childConfig.SupportedValues) > 0 && len(parentConfig.SupportedValues) > 0 {
		if !pv.isSubsetOfValues(childConfig.SupportedValues, parentConfig.SupportedValues) {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:      "PARAMETER_VALUES_CONFLICT",
				Message:   fmt.Sprintf("Parameter '%s' has supported values in child route that are not supported by parent route", paramName),
				RoutePath: route.Path,
				FilePath:  route.TemplateFile,
				Suggestions: []string{
					"Ensure child supported values are a subset of parent supported values",
					"Update parent route to support additional values",
					"Remove conflicting values from child route",
				},
			})
		}
	}

	// Check if child parameter is more restrictive (which is generally good)
	if childConfig.Validation != "" && parentConfig.Validation == "" {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:      "PARAMETER_INHERITANCE_INFO",
			Message:   fmt.Sprintf("Parameter '%s' adds validation in child route that doesn't exist in parent - this is generally good practice", paramName),
			RoutePath: route.Path,
			FilePath:  route.TemplateFile,
			Suggestions: []string{
				"Consider adding the same validation to parent route for consistency",
			},
		})
	}
}

// validateParameterInheritanceAvailability checks if inherited parameters are properly available
func (pv *ParameterValidator) validateParameterInheritanceAvailability(
	route *interfaces.Route,
	parentPath string,
	parentConfig *interfaces.ConfigFile,
	childConfig *interfaces.ConfigFile,
	result *ValidationResult,
) {
	// Extract parameters from route paths
	childParams := pv.extractDirectoryParameters(route.Path)

	// Check if child route uses parameters that are defined in parent config but not in child config
	// This covers the case where parent defines parameter configs that child routes should inherit
	for paramName := range parentConfig.DynamicSettings.Parameters {
		// If parameter is used in child route path but not configured in child
		if pv.parameterExistsInRoute(paramName, route.Path, childParams) {
			if _, exists := childConfig.DynamicSettings.Parameters[paramName]; !exists {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:      "INHERITED_PARAMETER_MISSING_CONFIG",
					Message:   fmt.Sprintf("Parameter '%s' is inherited from parent route '%s' but lacks configuration in child route", paramName, parentPath),
					RoutePath: route.Path,
					FilePath:  route.TemplateFile,
					Suggestions: []string{
						fmt.Sprintf("Add configuration for inherited parameter '%s'", paramName),
						"Copy parameter configuration from parent route",
						"Ensure parameter validation is consistent with parent",
					},
				})
			}
		}
	}
}

// isSubsetOfValues checks if child values are a subset of parent values
func (pv *ParameterValidator) isSubsetOfValues(childValues, parentValues []string) bool {
	parentSet := make(map[string]bool)
	for _, value := range parentValues {
		parentSet[value] = true
	}

	for _, childValue := range childValues {
		if !parentSet[childValue] {
			return false
		}
	}
	return true
}
