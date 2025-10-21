package middleware

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)

// DataServiceMiddleware handles data service resolution for templates
type DataServiceMiddleware struct {
	templateDataService interfaces.TemplateDataService
}

// Ensure DataServiceMiddleware implements DataServiceMiddlewareInterface
var _ DataServiceMiddlewareInterface = (*DataServiceMiddleware)(nil)

// NewDataServiceMiddleware creates a new data service middleware for DI
func NewDataServiceMiddleware(i do.Injector) (DataServiceMiddlewareInterface, error) {
	templateDataService := do.MustInvoke[interfaces.TemplateDataService](i)
	return &DataServiceMiddleware{
		templateDataService: templateDataService,
	}, nil
}

// DataServiceContextKey is the context key for storing resolved template data
type DataServiceContextKey string

const (
	// TemplateDataKey is the context key for template data
	TemplateDataKey DataServiceContextKey = "template_data"
	// TemplateComponentKey is the context key for resolved template component
	TemplateComponentKey DataServiceContextKey = "template_component"
)

// ResolveTemplateData middleware resolves template data before rendering
func (m *DataServiceMiddleware) ResolveTemplateData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		// Get template key from context (should be set by routing middleware)
		templateKey, ok := ctx.Value("template_key").(string)
		if !ok {
			// No template key found, continue without data resolution
			next.ServeHTTP(w, r)
			return
		}

		// Check if template requires data
		if !m.templateDataService.RequiresData(templateKey) {
			// No data required, continue normally
			next.ServeHTTP(w, r)
			return
		}

		// Extract route parameters from context
		routeParams, ok := ctx.Value("route_params").(map[string]string)
		if !ok {
			routeParams = make(map[string]string)
		}

		// Resolve template with data
		component, err := m.templateDataService.ResolveTemplateWithData(ctx, templateKey, routeParams)
		if err != nil {
			// Handle data resolution error
			http.Error(w, "Failed to resolve template data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Store resolved component in context
		ctx = context.WithValue(ctx, TemplateComponentKey, component)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetResolvedComponent retrieves the resolved template component from context
func GetResolvedComponent(ctx context.Context) (templ.Component, bool) {
	component, ok := ctx.Value(TemplateComponentKey).(templ.Component)
	return component, ok
}

// HasResolvedComponent checks if a template component has been resolved
func HasResolvedComponent(ctx context.Context) bool {
	_, ok := ctx.Value(TemplateComponentKey).(templ.Component)
	return ok
}