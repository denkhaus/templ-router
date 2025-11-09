package metadata

import (
	"context"

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
	return extractLocaleSpecificMetadata(routeMetadata, key, "")
}

// extractLocaleSpecificMetadata extracts metadata from RouteMetadata with locale support
// If locale is empty, returns regular metadata. Otherwise tries locale-specific first.
func extractLocaleSpecificMetadata(routeMetadata interface{}, key string, locale string) string {
	if routeMetadata == nil {
		return ""
	}

	// Try interface{}-keyed map
	if routeMap, ok := routeMetadata.(map[interface{}]interface{}); ok {
		// Try locale-specific metadata first if locale is provided
		if locale != "" {
			if localeData, exists := routeMap[locale]; exists {
				if localeMap, ok := localeData.(map[interface{}]interface{}); ok {
					if value, exists := localeMap[key]; exists {
						if strValue, ok := value.(string); ok {
							return strValue
						}
					}
				}
			}
		}

		// Fallback to regular metadata
		if value, exists := routeMap[key]; exists {
			if strValue, ok := value.(string); ok {
				return strValue
			}
		}
	}

	// Try string-keyed map
	if routeMap, ok := routeMetadata.(map[string]interface{}); ok {
		// Try locale-specific metadata first if locale is provided
		if locale != "" {
			if localeData, exists := routeMap[locale]; exists {
				if localeMap, ok := localeData.(map[string]interface{}); ok {
					if value, exists := localeMap[key]; exists {
						if strValue, ok := value.(string); ok {
							return strValue
						}
					}
				}
			}
		}

		// Fallback to regular metadata
		if value, exists := routeMap[key]; exists {
			if strValue, ok := value.(string); ok {
				return strValue
			}
		}
	}

	return "[MISSING_METADATA: " + key + "]" // Key not found
}

// tryLoadComponentMetadata attempts to load metadata from component cache
// This handles embedded components that have their own metadata separate from the page/layout
// Uses the pre-loaded component metadata cache from middleware for performance
// Now supports locale-specific metadata resolution
func tryLoadComponentMetadata(ctx context.Context, key string) string {
	// Get current locale from context
	locale, _ := ctx.Value(shared.LocaleKey).(string)

	// Get component metadata from cache (only source of truth - no fallbacks)
	if componentsConfigs, ok := ctx.Value(shared.ComponentsMetadataKey).(map[string]*shared.ConfigFile); ok {
		// Try all component configs to find one with the key
		// Components are pre-loaded by middleware with proper locale awareness
		for _, componentConfig := range componentsConfigs {
			// Try locale-specific metadata first if locale is available
			result := extractLocaleSpecificMetadata(componentConfig.GetRouteMetadata(), key, locale)
			if result != "[MISSING_METADATA: "+key+"]" {
				return result
			}
		}
	}

	// No fallbacks - if it's not in the pre-loaded cache, it's not available
	return ""
}



