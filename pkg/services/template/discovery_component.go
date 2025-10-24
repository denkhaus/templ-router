package template

import (
	"context"
	"embed"
	"html/template"
	"io"

	"github.com/a-h/templ"
)

//go:embed templates/*.html
var discoveryTemplates embed.FS

// DiscoveryComponent handles creation of discovery-related template components
type DiscoveryComponent struct{}

// DiscoveredTemplateData holds data for the discovered template page
type DiscoveredTemplateData struct {
	TemplateKey  string
	FunctionName string
	Params       map[string]string
}

// TemplateNotFoundData holds data for the template not found page
type TemplateNotFoundData struct {
	TemplateKey   string
	AvailableKeys []string
}

// NewDiscoveryComponent creates a new discovery component generator
func NewDiscoveryComponent() *DiscoveryComponent {
	return &DiscoveryComponent{}
}

// CreateDiscoveredTemplateComponent creates a component showing discovered template info
func (dc *DiscoveryComponent) CreateDiscoveredTemplateComponent(templateKey, functionName string, params map[string]string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Parse embedded template
		tmpl, err := template.ParseFS(discoveryTemplates, "templates/discovered_template.html")
		if err != nil {
			// Fallback to simple message
			_, writeErr := w.Write([]byte("Template discovered: " + templateKey))
			return writeErr
		}

		// Prepare template data
		data := DiscoveredTemplateData{
			TemplateKey:  templateKey,
			FunctionName: functionName,
			Params:       params,
		}

		// Execute template
		return tmpl.Execute(w, data)
	})
}

// CreateNotFoundComponent creates a component for when a template is not found
func (dc *DiscoveryComponent) CreateNotFoundComponent(templateKey string, availableKeys []string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Parse embedded template
		tmpl, err := template.ParseFS(discoveryTemplates, "templates/template_not_found.html")
		if err != nil {
			// Fallback to simple message
			_, writeErr := w.Write([]byte("Template not found: " + templateKey))
			return writeErr
		}

		// Prepare template data
		data := TemplateNotFoundData{
			TemplateKey:   templateKey,
			AvailableKeys: availableKeys,
		}

		// Execute template
		return tmpl.Execute(w, data)
	})
}