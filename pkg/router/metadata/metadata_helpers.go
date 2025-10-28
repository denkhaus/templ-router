package metadata

import (
	"context"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
)

// M retrieves a metadata value from the current template's YAML configuration
func M(ctx context.Context, key string) string {
	// Extract config from context
	configValue := ctx.Value(shared.TemplateConfigKey)
	if configValue == nil {
		return "[MISSING_METADATA_CONTEXT: " + key + "]" // No config available
	}

	// Try router.ConfigFile first
	if config, ok := configValue.(*interfaces.ConfigFile); ok {
		return extractMetadataFromConfig(config.RouteMetadata, key)
	}

	// Try shared.ConfigFile as fallback
	if sharedConfig, ok := configValue.(*shared.ConfigFile); ok {
		return extractMetadataFromConfig(sharedConfig.RouteMetadata, key)
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
