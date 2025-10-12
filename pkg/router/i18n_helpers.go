package router

import (
	"context"

	"github.com/denkhaus/templ-router/pkg/shared"
)

// M retrieves a metadata value from the current template's YAML configuration
func M(ctx context.Context, key string) string {
	// Extract config from context
	configValue := ctx.Value("template_config")
	if configValue == nil {
		return "[MISSING_METADATA_CONTEXT: " + key + "]" // No config available
	}

	// Try router.ConfigFile first
	if config, ok := configValue.(*ConfigFile); ok {
		return extractMetadataFromConfig(config.RouteMetadata, key)
	}

	// Try shared.ConfigFile as fallback
	if sharedConfig, ok := configValue.(*shared.ConfigFile); ok {
		return extractMetadataFromConfig(sharedConfig.RouteMetadata, key)
	}

	return "[INVALID_METADATA_CONFIG: " + key + "]" // Invalid config type
}

// extractLocaleFromRequest extracts locale from URL path or headers
// func extractLocaleFromRequest(r *http.Request) string {
// 	// First, try to extract from URL path (e.g., /en/dashboard, /de/admin)
// 	path := strings.TrimPrefix(r.URL.Path, "/")
// 	pathParts := strings.Split(path, "/")

// 	if len(pathParts) > 0 {
// 		firstPart := pathParts[0]
// 		if isValidLocaleCode(firstPart) {
// 			return firstPart
// 		}
// 	}

// 	// Fallback to Accept-Language header
// 	acceptLang := r.Header.Get("Accept-Language")
// 	if acceptLang != "" {
// 		// Parse Accept-Language header (simplified)
// 		langs := strings.Split(acceptLang, ",")
// 		for _, lang := range langs {
// 			// Remove quality values (e.g., "en-US;q=0.9" -> "en-US")
// 			lang = strings.Split(strings.TrimSpace(lang), ";")[0]
// 			// Extract primary language (e.g., "en-US" -> "en")
// 			primaryLang := strings.Split(lang, "-")[0]
// 			if isValidLocaleCode(primaryLang) {
// 				return primaryLang
// 			}
// 		}
// 	}

// 	// Default fallback
// 	return "en"
// }

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
