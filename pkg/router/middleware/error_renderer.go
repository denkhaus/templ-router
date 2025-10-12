package middleware

import (
	"context"
	"embed"
	"html/template"
	"io"
)

//go:embed templates/*.html
var errorTemplates embed.FS

// ErrorRenderer handles HTML rendering for error pages
// SEPARATED FROM: error_service.go (Separation of Concerns)
type ErrorRenderer struct{
	template *template.Template
}

// NewErrorRenderer creates a new error renderer
func NewErrorRenderer() *ErrorRenderer {
	tmpl, err := template.ParseFS(errorTemplates, "templates/error_fallback.html")
	if err != nil {
		// Fallback to nil template - will use simple string rendering
		return &ErrorRenderer{template: nil}
	}
	return &ErrorRenderer{template: tmpl}
}

// SimpleErrorComponent renders a simple error page using embedded HTML template
// FIXED: Now uses professional embed.FS template system
type SimpleErrorComponent struct {
	StatusCode   int
	Message      string
	TemplatePath string
	RequestPath  string
	renderer     *ErrorRenderer
}

// Render renders the error component using embedded HTML template
// FIXED: Replaced hardcoded HTML with professional embed.FS template system
func (sec *SimpleErrorComponent) Render(ctx context.Context, w io.Writer) error {
	if sec.renderer == nil || sec.renderer.template == nil {
		// Fallback to simple text if template is not available
		_, err := w.Write([]byte("Error: Template not available"))
		return err
	}
	
	// Execute the embedded template with error data
	return sec.renderer.template.Execute(w, sec)
}

// RenderErrorHTML renders error HTML using embedded template
func (er *ErrorRenderer) RenderErrorHTML(statusCode int, message, templatePath, requestPath string) *SimpleErrorComponent {
	component := &SimpleErrorComponent{
		StatusCode:   statusCode,
		Message:      message,
		TemplatePath: templatePath,
		RequestPath:  requestPath,
	}
	
	// Set the renderer reference for template execution
	component.renderer = er
	return component
}