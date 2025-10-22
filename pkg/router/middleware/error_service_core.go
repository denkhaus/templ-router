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
	templateResolver       ErrorTemplateResolver // Interface (proper DI)
	renderer              *ErrorRenderer
	templateService       interfaces.TemplateService // Integration with OptimizedTemplateService
	dedicatedErrorService DedicatedErrorTemplateService // NEW: Dedicated error template service
	logger                *zap.Logger
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

	// NEW: Invoke dedicated error template service from DI
	dedicatedErrorService := do.MustInvoke[DedicatedErrorTemplateService](i)

	return &ErrorServiceCore{
		templateResolver:       templateResolver,
		renderer:              renderer,
		templateService:       templateService,
		dedicatedErrorService: dedicatedErrorService,
		logger:                logger,
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

// tryRenderErrorTemplate attempts to render an error template using the dedicated error service
func (esc *ErrorServiceCore) tryRenderErrorTemplate(errorTemplate *interfaces.ErrorTemplate, message, path string) templ.Component {
	// FIXED: Use dedicated error template service instead of OptimizedTemplateService
	// This solves the conflict where OptimizedTemplateService would resolve wrong templates for error paths

	esc.logger.Debug("Attempting error template resolution through DedicatedErrorTemplateService",
		zap.String("template", errorTemplate.FilePath),
		zap.String("message", message),
		zap.String("path", path))

	// Check if error template is available in registry
	if !esc.dedicatedErrorService.IsErrorTemplateAvailable(errorTemplate) {
		esc.logger.Debug("Error template not available in registry",
			zap.String("template", errorTemplate.FilePath))
		return nil
	}

	// Create error context for template
	errorContext := &ErrorContext{
		StatusCode:  errorTemplate.ErrorCode,
		Message:     message,
		RequestPath: path,
	}

	// Render error template using dedicated service
	component, err := esc.dedicatedErrorService.RenderErrorTemplate(errorTemplate, errorContext)
	if err != nil {
		esc.logger.Error("Failed to render error template",
			zap.String("template", errorTemplate.FilePath),
			zap.Error(err))
		return nil
	}

	esc.logger.Info("Successfully rendered error template",
		zap.String("template", errorTemplate.FilePath),
		zap.String("path", path))

	return component
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
