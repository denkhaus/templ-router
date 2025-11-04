package i18n

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/denkhaus/templ-router/pkg/shared"
	"go.uber.org/zap"
)

// I18nData holds translation data for the current request
type I18nData struct {
	Locale          string
	CurrentTemplate string
	Translations    map[string]string // key -> translation
	FallbackLocale  string
	Logger          *zap.Logger
	mu              sync.RWMutex
}

// T translates a key using the current context
func T(ctx context.Context, key string) string {
	data, ok := ctx.Value(shared.I18nDataKey).(*I18nData)
	if !ok {
		// Graceful fallback when i18n context is missing
		return fmt.Sprintf("[MISSING_I18N_CONTEXT: %s]", key)
	}

	data.mu.RLock()
	defer data.mu.RUnlock()

	if translation, exists := data.Translations[key]; exists {
		data.Logger.Debug("Translation found",
			zap.String("key", key),
			zap.String("locale", data.Locale),
			zap.String("template", data.CurrentTemplate),
			zap.String("translation", translation))
		return translation
	}

	// Graceful fallback for missing translations
	data.Logger.Warn("Translation key not found - using fallback",
		zap.String("key", key),
		zap.String("locale", data.Locale),
		zap.String("template", data.CurrentTemplate),
		zap.String("fallback", fmt.Sprintf("[MISSING: %s]", key)))

	// Return a visible but non-breaking fallback
	return fmt.Sprintf("[MISSING_I18N: %s]", key)
}

// GetCurrentLocale returns the current locale from context
func GetCurrentLocale(ctx context.Context) string {
	locale, ok := ctx.Value(shared.LocaleKey).(string)
	if !ok {
		// Graceful fallback to default locale
		return "en"
	}
	return locale
}

// GetCurrentTemplate returns the current template from context
func GetCurrentTemplate(ctx context.Context) string {
	template, ok := ctx.Value(shared.I18nTemplateKey).(string)
	if !ok {
		// Graceful fallback
		return "unknown"
	}
	return template
}

// GetAvailableKeys returns all available translation keys for current template
func GetAvailableKeys(ctx context.Context) []string {
	data, ok := ctx.Value(shared.I18nDataKey).(*I18nData)
	if !ok {
		return nil
	}

	data.mu.RLock()
	defer data.mu.RUnlock()

	keys := make([]string, 0, len(data.Translations))
	for key := range data.Translations {
		keys = append(keys, key)
	}

	return keys
}

// GetCurrentRoute returns the current route path from context
// Returns empty string if request path is not stored in context
func GetCurrentRoute(ctx context.Context) string {
	path, ok := ctx.Value(shared.RequestPathKey).(string)
	if !ok {
		return ""
	}
	return path
}

// IsLocalizedRoute checks if the given path matches any route pattern that contains {locale}
// Made public for testing
func IsLocalizedRoute(path string, routeMapping map[string]string) bool {
	// Check each route pattern to see if it matches our path and contains {locale}
	for routePattern := range routeMapping {
		if !strings.Contains(routePattern, "{locale}") {
			continue
		}

		// Simple pattern matching - replace all path parameters with wildcards and check if path matches
		// This is a basic implementation; for more complex routing, you might need a more sophisticated matcher
		pattern := strings.ReplaceAll(routePattern, "{locale}", "*")
		// Replace all other path parameters like {id}, {userId}, etc. with *
		for i := 0; i < len(pattern); i++ {
			if pattern[i] == '{' {
				// Find the closing brace
				end := strings.Index(pattern[i:], "}")
				if end != -1 {
					// Replace the entire parameter with *
					pattern = pattern[:i] + "*" + pattern[i+end+1:]
				}
			}
		}
			if PathMatchesPattern(path, pattern) {
			return true
		}
	}
	return false
}

// PathMatchesPattern performs simple wildcard matching for route patterns
// Supports * as wildcard for path segments
// Made public for testing
func PathMatchesPattern(path, pattern string) bool {
	// Split both paths into segments
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")
	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")

	// If different number of segments, no match
	if len(pathSegments) != len(patternSegments) {
		return false
	}

	// Check each segment
	for i, patternSegment := range patternSegments {
		if patternSegment != "*" && patternSegment != pathSegments[i] {
			return false
		}
	}

	return true
}

// GetCurrentRouteWithoutLocale returns the current route path with locale stripped
// Uses route mapping to accurately determine if a route is localized
// Returns the full path if route is not localized or route mapping is not available
func GetCurrentRouteWithoutLocale(ctx context.Context) string {
	fullPath := GetCurrentRoute(ctx)
	if fullPath == "" {
		return ""
	}

	// Get route mapping from context
	routeMapping, ok := ctx.Value(shared.RouteMappingKey).(map[string]string)
	if !ok {
		// Fallback to old method if route mapping is not available
		return getCurrentRouteWithoutLocaleFallback(ctx, fullPath)
	}

	// Check if this path matches a localized route pattern
	if !IsLocalizedRoute(fullPath, routeMapping) {
		// This is not a localized route, return full path
		return fullPath
	}

	// This is a localized route, strip the locale segment
	parts := strings.Split(fullPath, "/")

	// parts[0] will be empty due to leading slash, so check parts[1] for locale
	if len(parts) >= 2 && len(parts[1]) == 2 {
		// Strip the locale segment (parts[1]) and reconstruct path
		pathWithoutLocale := strings.Join(parts[2:], "/")
		if pathWithoutLocale == "" {
			return "/"
		}
		return "/" + pathWithoutLocale
	}

	// No locale segment found despite being a localized route, return full path
	return fullPath
}

// getCurrentRouteWithoutLocaleFallback is the fallback method when route mapping is not available
func getCurrentRouteWithoutLocaleFallback(ctx context.Context, fullPath string) string {
	// Check if localization is active by checking if we have a valid current template
	// If GetCurrentTemplate returns "unknown", localization is likely not active for this route
	currentTemplate := GetCurrentTemplate(ctx)
	if currentTemplate == "unknown" {
		// No localization active, return full path
		return fullPath
	}

	// Localization is active, try to strip locale from path
	parts := strings.Split(fullPath, "/")

	// parts[0] will be empty due to leading slash, so check parts[1] for locale
	if len(parts) >= 2 && len(parts[1]) == 2 {
		// Strip the locale segment (parts[1]) and reconstruct path
		pathWithoutLocale := strings.Join(parts[2:], "/")
		if pathWithoutLocale == "" {
			return "/"
		}
		return "/" + pathWithoutLocale
	}

	// No locale segment found, return full path
	return fullPath
}
