package template

import (
	"github.com/a-h/templ"
)

// TemplateService defines the interface for template loading and rendering
type TemplateService interface {
	// LoadTemplate loads a template by key and renders it with the given parameters
	LoadTemplate(templateKey string, params map[string]string) templ.Component
	
	// IsTemplateAvailable checks if a template is available for the given key
	IsTemplateAvailable(templateKey string) bool
	
	// GetAvailableTemplates returns a list of all available template keys
	GetAvailableTemplates() []string
	
	// RefreshTemplates reloads the template registry (useful for development)
	RefreshTemplates() error
}

