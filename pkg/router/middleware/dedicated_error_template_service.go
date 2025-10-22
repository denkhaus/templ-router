package middleware

import (
	"fmt"
	"reflect"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// DedicatedErrorTemplateService interface for error template rendering (PUBLIC)
type DedicatedErrorTemplateService interface {
	RenderErrorTemplate(errorTemplate *interfaces.ErrorTemplate, errorContext *ErrorContext) (templ.Component, error)
	IsErrorTemplateAvailable(errorTemplate *interfaces.ErrorTemplate) bool
}

// dedicatedErrorTemplateServiceImpl handles error template rendering separately from regular templates (PRIVATE)
// This solves the conflict issue where OptimizedTemplateService would resolve wrong templates for error paths
type dedicatedErrorTemplateServiceImpl struct {
	templateRegistry interfaces.TemplateRegistry
	logger           *zap.Logger
}

// NewDedicatedErrorTemplateService creates a dedicated error template service (RETURNS INTERFACE)
func NewDedicatedErrorTemplateService(i do.Injector) (DedicatedErrorTemplateService, error) {
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &dedicatedErrorTemplateServiceImpl{
		templateRegistry: templateRegistry,
		logger:           logger,
	}, nil
}

// RenderErrorTemplate renders an error template by direct template key lookup
// This bypasses route-based resolution to avoid conflicts
func (dets *dedicatedErrorTemplateServiceImpl) RenderErrorTemplate(errorTemplate *interfaces.ErrorTemplate, errorContext *ErrorContext) (templ.Component, error) {
	dets.logger.Debug("Rendering error template directly",
		zap.String("template_path", errorTemplate.FilePath),
		zap.Int("error_code", errorTemplate.ErrorCode))

	// Generate template key directly from error template path
	templateKey := dets.generateErrorTemplateKey(errorTemplate)
	
	dets.logger.Debug("Looking up error template by key",
		zap.String("template_key", templateKey))

	// Get template function directly by key (no route resolution)
	templateFunc, exists := dets.templateRegistry.GetTemplateFunction(templateKey)
	if !exists {
		dets.logger.Debug("Error template not found in registry",
			zap.String("template_key", templateKey))
		return nil, fmt.Errorf("error template not found: %s", templateKey)
	}

	// Execute template function with error context
	return dets.executeErrorTemplateFunction(templateFunc, errorContext)
}

// generateErrorTemplateKey generates a template key for error templates
// This follows the same pattern as regular templates but for error.templ files
func (dets *dedicatedErrorTemplateServiceImpl) generateErrorTemplateKey(errorTemplate *interfaces.ErrorTemplate) string {
	// Convert file path to template key
	// Example: "app/admin/error.templ" -> "app/admin/error.templ#Error"
	return fmt.Sprintf("%s#Error", errorTemplate.FilePath)
}

// executeErrorTemplateFunction executes error template function with proper error context
func (dets *dedicatedErrorTemplateServiceImpl) executeErrorTemplateFunction(templateFunc func() interface{}, errorContext *ErrorContext) (templ.Component, error) {
	result := templateFunc()

	// Handle parameterless error template (most common)
	if fn, ok := result.(func() templ.Component); ok {
		component := fn()
		dets.logger.Debug("Parameterless error template executed successfully")
		return component, nil
	}

	// Handle error template with error context parameter
	if fn, ok := result.(func(*ErrorContext) templ.Component); ok {
		component := fn(errorContext)
		dets.logger.Debug("Error template with context executed successfully")
		return component, nil
	}

	// Handle generic interface{} parameter (fallback)
	if fn, ok := result.(func(interface{}) templ.Component); ok {
		component := fn(errorContext)
		dets.logger.Debug("Generic error template executed successfully")
		return component, nil
	}

	// Use reflection for unknown signatures (last resort)
	return dets.executeErrorTemplateWithReflection(result, errorContext)
}

// executeErrorTemplateWithReflection handles unknown template function signatures using reflection
func (dets *dedicatedErrorTemplateServiceImpl) executeErrorTemplateWithReflection(templateFunc interface{}, errorContext *ErrorContext) (templ.Component, error) {
	funcValue := reflect.ValueOf(templateFunc)
	funcType := funcValue.Type()

	if funcType.Kind() != reflect.Func {
		return nil, fmt.Errorf("template is not a function")
	}

	// Prepare arguments based on function signature
	var args []reflect.Value

	if funcType.NumIn() == 1 {
		// Single parameter - pass error context
		args = []reflect.Value{reflect.ValueOf(errorContext)}
	} else if funcType.NumIn() == 0 {
		// No parameters
		args = []reflect.Value{}
	} else {
		return nil, fmt.Errorf("unsupported error template function signature: %d parameters", funcType.NumIn())
	}

	// Call the function
	results := funcValue.Call(args)
	if len(results) != 1 {
		return nil, fmt.Errorf("error template function should return one value")
	}

	// Convert result to templ.Component
	component, ok := results[0].Interface().(templ.Component)
	if !ok {
		return nil, fmt.Errorf("error template function did not return templ.Component")
	}

	dets.logger.Debug("Error template executed via reflection")
	return component, nil
}

// IsErrorTemplateAvailable checks if an error template is available in the registry
func (dets *dedicatedErrorTemplateServiceImpl) IsErrorTemplateAvailable(errorTemplate *interfaces.ErrorTemplate) bool {
	templateKey := dets.generateErrorTemplateKey(errorTemplate)
	return dets.templateRegistry.IsAvailable(templateKey)
}