package router

import (
	"fmt"
)

// MetadataRouteConfigurator handles route configuration operations
// Extracted from metadata.go for better separation of concerns
type MetadataRouteConfigurator struct{}

// NewMetadataRouteConfigurator creates a new route configurator
func NewMetadataRouteConfigurator() *MetadataRouteConfigurator {
	return &MetadataRouteConfigurator{}
}

// ApplyYAMLConfigToRoute applies YAML configuration to a single route
func (mrc *MetadataRouteConfigurator) ApplyYAMLConfigToRoute(route *Route, config *ConfigFile) error {
	if route == nil {
		return fmt.Errorf("route cannot be nil")
	}
	
	if config == nil {
		return nil
	}
	
	// Apply route metadata if specified
	if config.RouteMetadata != nil {
		// Note: RouteMetadata is interface{} - needs proper type assertion
		// For now, skip complex metadata parsing to avoid interface issues
		// TODO: Implement proper RouteMetadata type assertion if needed
	}
	
	// Apply i18n mappings if specified
	// Note: Route interface doesn't have I18nMappings field
	// I18n is handled separately by I18nService
	// Skip i18n mapping application to avoid interface issues
	
	// Apply other AuthSettings if specified directly in the config
	if config.AuthSettings != nil {
		route.AuthSettings = config.AuthSettings
	}
	
	return nil
}

// ApplyYAMLConfigsToRoutes applies YAML configurations to multiple routes
func (mrc *MetadataRouteConfigurator) ApplyYAMLConfigsToRoutes(routes []Route, configs map[string]*ConfigFile) ([]Route, error) {
	if routes == nil {
		return nil, fmt.Errorf("routes cannot be nil")
	}
	
	if configs == nil {
		return routes, nil
	}

	updatedRoutes := make([]Route, len(routes))
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
func ApplyYAMLConfigsToRoutes(routes []Route, configs map[string]*ConfigFile) ([]Route, error) {
	configurator := NewMetadataRouteConfigurator()
	return configurator.ApplyYAMLConfigsToRoutes(routes, configs)
}