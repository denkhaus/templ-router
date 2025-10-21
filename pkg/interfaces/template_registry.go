package interfaces

import (
	"context"
	"github.com/a-h/templ"
)

// TemplateRegistry provides access to application templates
// This interface allows the router to work with any template registry implementation
type TemplateRegistry interface {
	// GetTemplate retrieves a template component by key
	GetTemplate(key string) (templ.Component, error)
	
	// GetTemplateFunction retrieves a template function by key
	GetTemplateFunction(key string) (func() interface{}, bool)
	
	// GetAllTemplateKeys returns all available template keys
	GetAllTemplateKeys() []string
	
	// IsAvailable checks if a template exists
	IsAvailable(key string) bool
	
	// Route-to-Template mapping
	GetRouteToTemplateMapping() map[string]string
	GetTemplateByRoute(route string) (templ.Component, error)
	
	// Data Service Integration
	RequiresDataService(key string) bool
	GetDataServiceInfo(key string) (DataServiceInfo, bool)
}

// TemplateDataService handles data resolution for templates
type TemplateDataService interface {
	// ResolveTemplateWithData resolves template data and renders the template
	ResolveTemplateWithData(ctx context.Context, templateKey string, routeParams map[string]string) (templ.Component, error)
	
	// RequiresData checks if a template requires data service
	RequiresData(templateKey string) bool
}

