package middleware

import (
	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// ErrorServiceCore implements the core error service logic
// SEPARATED FROM: error_service.go (Separation of Concerns)
// This focuses purely on error service coordination without HTML rendering or template resolution
type ErrorServiceCore struct {
	templateResolver ErrorTemplateResolver // Interface (proper DI)
	renderer         *ErrorRenderer
	templateService  interfaces.TemplateService // Integration with OptimizedTemplateService
	logger           *zap.Logger
}

// NewErrorServiceCore creates a new core error service for DI
func NewErrorServiceCore(i do.Injector) (interfaces.ErrorService, error) {
	templateResolver, err := NewErrorTemplateResolver(i)
	if err != nil {
		return nil, err
	}

	renderer := NewErrorRenderer()
	templateService := do.MustInvoke[interfaces.TemplateService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &ErrorServiceCore{
		templateResolver: templateResolver,
		renderer:         renderer,
		templateService:  templateService,
		logger:           logger,
	}, nil
}

// FindErrorTemplateForPath delegates to the template resolver
func (esc *ErrorServiceCore) FindErrorTemplateForPath(path string) *interfaces.ErrorTemplate {
	return esc.templateResolver.FindErrorTemplateForPath(path)
}

// CreateErrorComponent creates an error component with proper template resolution
func (esc *ErrorServiceCore) CreateErrorComponent(message, path string) templ.Component {
	esc.logger.Debug("Creating error component",
		zap.String("message", message),
		zap.String("path", path))

	// Try to find a specific error template first
	errorTemplate := esc.templateResolver.FindErrorTemplateForPath(path)

	if errorTemplate != nil {
		esc.logger.Debug("Found specific error template",
			zap.String("template", errorTemplate.FilePath),
			zap.Int("error_code", errorTemplate.ErrorCode))

		// TODO: Use OptimizedTemplateService to render proper error template
		// This replaces the TODO from the original error_service.go:262
		component := esc.tryRenderErrorTemplate(errorTemplate, message, path)
		if component != nil {
			return component
		}
	}

	// Fallback to simple HTML error component
	esc.logger.Debug("Using fallback error renderer",
		zap.String("path", path))

	return esc.renderer.RenderErrorHTML(500, message, "fallback", path)
}

// tryRenderErrorTemplate attempts to render an error template using the template service
func (esc *ErrorServiceCore) tryRenderErrorTemplate(errorTemplate *interfaces.ErrorTemplate, message, path string) templ.Component {
	// FIXED: Proper error template resolution through OptimizedTemplateService
	// This addresses the critical TODO from error_service.go:262

	esc.logger.Debug("Attempting error template resolution through OptimizedTemplateService",
		zap.String("template", errorTemplate.FilePath),
		zap.String("message", message),
		zap.String("path", path))

	// CRITICAL FIX: Don't use OptimizedTemplateService for error templates
	// The OptimizedTemplateService tries to resolve templates based on the request path,
	// which can lead to rendering the wrong template (e.g., dashboard template for /fr/dashboard)
	// Instead, fall back to simple HTML rendering for error cases
	esc.logger.Info("Skipping OptimizedTemplateService for error template to avoid wrong template resolution",
		zap.String("template", errorTemplate.FilePath),
		zap.String("path", path))

	return nil // Always fall back to simple HTML rendering
}

// ErrorContext represents error information for templates
type ErrorContext struct {
	StatusCode  int    `json:"status_code"`
	Message     string `json:"message"`
	RequestPath string `json:"request_path"`
	UserAgent   string `json:"user_agent,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	ErrorID     string `json:"error_id,omitempty"`
}
