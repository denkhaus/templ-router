package template

// TemplateService defines the interface for template loading and rendering
type TemplateService interface {
	// IsTemplateAvailable checks if a template is available for the given key
	IsTemplateAvailable(templateKey string) bool
	
	// GetAvailableTemplates returns a list of all available template keys
	GetAvailableTemplates() []string
	
	// RefreshTemplates reloads the template registry (useful for development)
	RefreshTemplates() error
}

