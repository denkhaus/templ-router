package router

import (
	"fmt"
	"regexp"
	"strings"
)

// DynamicRouteSegment represents a route parameter defined using dollar sign convention
type DynamicRouteSegment struct {
	// Pattern is the pattern name (e.g., "$id", "$slug", "$locale")
	Pattern string

	// RoutePath is the full route path containing this segment
	RoutePath string

	// TemplateName is the template that handles this dynamic route
	TemplateName string

	// ParameterName is the name of the parameter (e.g., "id", "slug", "locale")
	ParameterName string

	// ValidationRegex is an optional regex to validate the parameter value
	ValidationRegex string

	// Description is a human-readable description of the parameter
	Description string

	// SupportedValues is a list of explicitly allowed values (optional)
	SupportedValues []string

	// Config contains the YAML configuration for this parameter (if available)
	Config *DynamicParameterConfig
}

// RecognizeDynamicRoutes identifies dynamic route patterns using dollar sign convention
// Only $locale is reserved for localization, all other parameters use the $ prefix (e.g., $id, $slug)
func RecognizeDynamicRoutes(routePath string, templateName string) []DynamicRouteSegment {
	return RecognizeDynamicRoutesWithConfig(routePath, templateName, nil)
}

// RecognizeDynamicRoutesWithConfig identifies dynamic route patterns with YAML configuration
func RecognizeDynamicRoutesWithConfig(routePath string, templateName string, dynamicSettings *DynamicSettings) []DynamicRouteSegment {
	segments := []DynamicRouteSegment{}

	// Split the route path into parts
	parts := strings.Split(routePath, "/")

	for _, part := range parts {
		if strings.HasPrefix(part, "$") && len(part) > 1 {
			parameterName := part[1:] // Remove the $ prefix
			pattern := part

			segment := DynamicRouteSegment{
				Pattern:       pattern,
				RoutePath:     routePath,
				TemplateName:  templateName,
				ParameterName: parameterName,
			}

			// Check if we have YAML configuration for this parameter
			if dynamicSettings != nil && dynamicSettings.Parameters != nil {
				if config, exists := dynamicSettings.Parameters[parameterName]; exists {
					segment.Config = config
					segment.ValidationRegex = config.Validation
					segment.Description = config.Description
					segment.SupportedValues = config.SupportedValues
				}
			}

			// No fallbacks - config is required for validation
			if segment.ValidationRegex == "" {
				panic(fmt.Sprintf("dynamic.go: validation regex not configured for parameter '%s' - config.Validation required", parameterName))
			}

			segments = append(segments, segment)
		}
	}

	return segments
}

// IsDynamicRoute checks if a route path contains dynamic segments
func IsDynamicRoute(routePath string) bool {
	return strings.Contains(routePath, "/$")
}

// ValidateDynamicSegmentValue validates a value against the segment's validation regex
func (d *DynamicRouteSegment) ValidateDynamicSegmentValue(value string) bool {
	// First check supported values if they are defined
	if len(d.SupportedValues) > 0 {
		for _, supportedValue := range d.SupportedValues {
			if value == supportedValue {
				return true
			}
		}
		return false // Value not in supported values list
	}

	// Fall back to regex validation
	if d.ValidationRegex == "" {
		return true // No validation regex means all values are valid
	}

	matched, err := regexp.MatchString(d.ValidationRegex, value)
	return err == nil && matched
}

// ValidateParameterValue validates a parameter value using YAML configuration
func ValidateParameterValue(paramName string, value string, config *DynamicParameterConfig) (bool, string) {
	if config == nil {
		return true, ""
	}

	// Check supported values first (more specific)
	if len(config.SupportedValues) > 0 {
		for _, supportedValue := range config.SupportedValues {
			if value == supportedValue {
				return true, ""
			}
		}
		return false, fmt.Sprintf("parameter '%s' value '%s' is not in supported values: %v", paramName, value, config.SupportedValues)
	}

	// Check regex validation
	if config.Validation != "" {
		matched, err := regexp.MatchString(config.Validation, value)
		if err != nil {
			return false, fmt.Sprintf("parameter '%s' validation regex error: %v", paramName, err)
		}
		if !matched {
			return false, fmt.Sprintf("parameter '%s' value '%s' does not match validation pattern '%s'", paramName, value, config.Validation)
		}
	}

	return true, ""
}
