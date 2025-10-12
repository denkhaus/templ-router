package router

import (
	"strings"
)

// UpdateRoutePathWithLocale updates a route path with the specified locale
func UpdateRoutePathWithLocale(routePath string, locale string) string {
	// If locale is empty, return the original path
	if locale == "" {
		return routePath
	}

	// Check if the path already contains a locale parameter
	// If it does, we might need to handle it differently
	if strings.Contains(routePath, "/$locale") {
		// Replace $locale placeholder with actual locale value
		return strings.Replace(routePath, "/$locale", "/"+locale, 1)
	}

	// Otherwise, prepend the locale to the path
	if routePath == "/" {
		return "/" + locale
	}

	return "/" + locale + routePath
}

// CreateLocalizedRoutePath creates a path for a route with locale support
func CreateLocalizedRoutePath(basePath string, locale string) string {
	// This function would create a route path that incorporates locale
	// It might check if the route supports locale, then create the appropriate path

	if locale == "" {
		return basePath
	}

	// Check if this is a locale-aware route (contains $locale)
	if strings.Contains(basePath, "/$locale") {
		// Replace $locale with the actual locale parameter
		return strings.Replace(basePath, "/$locale", "/"+locale, 1)
	}

	// If the route doesn't support locales, return as is
	return basePath
}
