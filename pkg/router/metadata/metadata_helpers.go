package metadata

import (
	"context"
	"strings"

	"github.com/denkhaus/templ-router/pkg/shared"
)

// M retrieves a metadata value from the current template's YAML configuration
func M(ctx context.Context, key string) string {
	// Extract config from context
	configValue := ctx.Value(shared.TemplateConfigKey)
	if configValue == nil {
		return "[MISSING_METADATA_CONTEXT: " + key + "]" // No config available
	}

	// Use shared.ConfigFile
	if sharedConfig, ok := configValue.(*shared.ConfigFile); ok {
		result := extractMetadataFromConfig(sharedConfig.GetRouteMetadata(), key)
		if result != "[MISSING_METADATA: "+key+"]" {
			return result // Found in main config
		}

		// Try to load metadata from the component's own YAML file
		// This handles embedded components that have their own metadata
		if componentResult := tryLoadComponentMetadata(ctx, key); componentResult != "" {
			return componentResult
		}
	}

	return "[INVALID_METADATA_CONFIG: " + key + "]" // Invalid config type
}

// extractMetadataFromConfig extracts metadata from RouteMetadata (works with both router and shared configs)
func extractMetadataFromConfig(routeMetadata interface{}, key string) string {
	if routeMetadata == nil {
		return ""
	}

	// Try interface{}-keyed map
	if routeMap, ok := routeMetadata.(map[interface{}]interface{}); ok {
		if value, exists := routeMap[key]; exists {
			if strValue, ok := value.(string); ok {
				return strValue
			}
		}
	}

	// Try string-keyed map
	if routeMap, ok := routeMetadata.(map[string]interface{}); ok {
		if value, exists := routeMap[key]; exists {
			if strValue, ok := value.(string); ok {
				return strValue
			}
		}
	}

	return "[MISSING_METADATA: " + key + "]" // Key not found
}

// tryLoadComponentMetadata attempts to load metadata from the component's own YAML file
// This handles embedded components that have their own metadata separate from the page/layout
func tryLoadComponentMetadata(ctx context.Context, key string) string {
	// Get the current template path from context
	templatePath, ok := ctx.Value(shared.I18nTemplateKey).(string)
	if !ok {
		return "" // No template path available
	}

	// Extract component name from template path
	// e.g., "app/components/footer.templ" -> "footer"
	componentName := extractComponentNameFromPath(templatePath)
	if componentName == "" {
		return "" // Not a component
	}

	// Build the YAML path for this component
	yamlPath := buildComponentYAMLPath(templatePath)
	if yamlPath == "" {
		return "" // No YAML path
	}

	// Load the component's metadata directly
	_, config, err := shared.ParseYAMLMetadata(yamlPath)
	if err != nil {
		return "" // Failed to load component metadata
	}

	// Extract the metadata from the component's config
	return extractMetadataFromConfig(config.GetRouteMetadata(), key)
}

// extractComponentNameFromPath extracts component name from template path
func extractComponentNameFromPath(templatePath string) string {
	// Check if this is a component template
	if !strings.Contains(templatePath, "/components/") {
		return "" // Not a component
	}

	// Extract filename without extension
	parts := strings.Split(templatePath, "/")
	if len(parts) == 0 {
		return ""
	}

	filename := parts[len(parts)-1]
	filename = strings.TrimSuffix(filename, ".templ")

	return filename
}

// buildComponentYAMLPath builds the YAML path for a component
func buildComponentYAMLPath(templatePath string) string {
	// Replace .templ with .templ.yaml
	if strings.HasSuffix(templatePath, ".templ") {
		return templatePath + ".yaml"
	}
	return templatePath + ".templ.yaml"
}
