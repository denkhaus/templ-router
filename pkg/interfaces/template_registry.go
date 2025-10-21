package interfaces

import (
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
