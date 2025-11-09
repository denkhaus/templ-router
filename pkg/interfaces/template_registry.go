package interfaces

import (
	"github.com/a-h/templ"
)

// TemplateType represents the type of template based on naming conventions
type TemplateType string

const (
	TemplateTypeLayout    TemplateType = "layout"    // layout.templ
	TemplateTypePage      TemplateType = "page"      // page.templ
	TemplateTypeError     TemplateType = "error"     // error.templ
	TemplateTypeComponent TemplateType = "component" // any other .templ file
)

// TemplateMetadata contains metadata about a template for generic discovery
type TemplateMetadata struct {
	// Original template file path (e.g., "app/components/footer.templ")
	TemplatePath string `json:"template_path"`

	// Template type based on naming convention
	Type TemplateType `json:"type"`

	// Component name extracted from filename (e.g., "footer" from "footer.templ")
	ComponentName string `json:"component_name,omitempty"`

	// Full route pattern for this template
	Route string `json:"route"`

	// YAML metadata file existence and path
	YAMLFile    string `json:"yaml_file,omitempty"` // Full path to YAML file if it exists
	YAMLExists  bool   `json:"yaml_exists"`         // Whether YAML file exists
	HasI18n     bool   `json:"has_i18n"`            // Whether YAML contains i18n data
	HasMetadata bool   `json:"has_metadata"`        // Whether YAML contains metadata
	HasAuth     bool   `json:"has_auth"`            // Whether YAML contains auth settings
}

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

	// Enhanced metadata access for generic discovery
	GetTemplateMetadata(key string) (*TemplateMetadata, error)
	GetTemplateMetadataByRoute(route string) (*TemplateMetadata, error)
	GetAllTemplateMetadata() map[string]*TemplateMetadata

	// Component discovery methods
	FindComponentTemplates() map[string]*TemplateMetadata // componentName -> metadata
	GetTemplateKeyByComponentName(componentName string) (string, bool)

	// Data Service Integration
	RequiresDataService(key string) bool
	GetDataServiceInfo(key string) (DataServiceInfo, bool)
}
