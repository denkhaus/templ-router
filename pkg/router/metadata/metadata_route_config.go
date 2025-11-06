package metadata

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
)

// MetadataRouteConfigurator handles route configuration operations
// Extracted from metadata.go for better separation of concerns
type MetadataRouteConfigurator struct{}

// NewMetadataRouteConfigurator creates a new route configurator
func NewMetadataRouteConfigurator() *MetadataRouteConfigurator {
	return &MetadataRouteConfigurator{}
}

// ApplyYAMLConfigToRoute applies YAML configuration to a single route
func (mrc *MetadataRouteConfigurator) ApplyYAMLConfigToRoute(route *interfaces.Route, config *shared.ConfigFile) error {
	if route == nil {
		return fmt.Errorf("route cannot be nil")
	}

	if config == nil {
		return nil
	}

	// Apply route metadata if specified
	if config.GetRouteMetadata() != nil {
		// Note: RouteMetadata is interface{} - needs proper type assertion
		// For now, skip complex metadata parsing to avoid interface issues
		// TODO: Implement proper RouteMetadata type assertion if needed
	}

	// Apply i18n mappings if specified
	// Note: Route interface doesn't have I18nMappings field
	// I18n is handled separately by I18nService
	// Skip i18n mapping application to avoid interface issues

	// Apply other AuthSettings if specified directly in the config
	if config.GetAuthSettings() != nil {
		// Convert the auth settings map to AuthSettings struct
		if authSettingsMap, ok := config.GetAuthSettings().(map[string]interface{}); ok {
			route.AuthSettings = convertMapToAuthSettings(authSettingsMap)
		}
	}

	return nil
}

// ApplyYAMLConfigsToRoutes applies YAML configurations to multiple routes
func (mrc *MetadataRouteConfigurator) ApplyYAMLConfigsToRoutes(routes []interfaces.Route, configs map[string]*shared.ConfigFile) ([]interfaces.Route, error) {
	if routes == nil {
		return nil, fmt.Errorf("routes cannot be nil")
	}

	if configs == nil {
		return routes, nil
	}

	updatedRoutes := make([]interfaces.Route, len(routes))
	copy(updatedRoutes, routes)

	for i := range updatedRoutes {
		route := &updatedRoutes[i]

		// Find the corresponding config for this route's template
		config, exists := configs[route.TemplateFile]
		if exists && config != nil {
			err := mrc.ApplyYAMLConfigToRoute(route, config)
			if err != nil {
				return nil, fmt.Errorf("failed to apply config to route %s: %w", route.Path, err)
			}
		}
	}

	return updatedRoutes, nil
}

// ApplyYAMLConfigsToRoutes is the legacy global function (DEPRECATED)
// Use MetadataRouteConfigurator.ApplyYAMLConfigsToRoutes instead
func ApplyYAMLConfigsToRoutes(routes []interfaces.Route, configs map[string]*shared.ConfigFile) ([]interfaces.Route, error) {
	configurator := NewMetadataRouteConfigurator()
	return configurator.ApplyYAMLConfigsToRoutes(routes, configs)
}

// convertMapToAuthSettings converts a map[string]interface{} to AuthSettings
// This helper function bridges the gap between the map-based auth settings and the typed AuthSettings
func convertMapToAuthSettings(authMap map[string]interface{}) *interfaces.AuthSettings {
	if authMap == nil {
		return nil
	}

	settings := &interfaces.AuthSettings{
		Type:     "Public", // Default to Public
		Settings: make(map[string]interface{}),
	}

	// Extract auth type
	if authType, exists := authMap["type"]; exists {
		if authTypeStr, ok := authType.(string); ok {
			settings.Type = authTypeStr
		}
	}

	// Extract redirect URL
	if redirectURL, exists := authMap["redirect_url"]; exists {
		if redirectURLStr, ok := redirectURL.(string); ok {
			settings.RedirectURL = redirectURLStr
		}
	}

	// Extract roles
	if roles, exists := authMap["roles"]; exists {
		if rolesSlice, ok := roles.([]string); ok {
			settings.Roles = rolesSlice
		} else if rolesInterface, ok := roles.([]interface{}); ok {
			settings.Roles = make([]string, 0, len(rolesInterface))
			for _, role := range rolesInterface {
				if roleStr, ok := role.(string); ok {
					settings.Roles = append(settings.Roles, roleStr)
				}
			}
		}
	}

	// Copy other settings
	for key, value := range authMap {
		if key != "type" && key != "redirect_url" && key != "roles" {
			settings.Settings[key] = value
		}
	}

	return settings
}
