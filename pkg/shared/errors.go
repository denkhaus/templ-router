package shared

import (
	"errors"
	"fmt"
)

// Common error types for consistent error handling across the application
var (
	// Template-related errors
	ErrTemplateNotFound     = errors.New("template not found")
	ErrTemplateInvalid      = errors.New("template is invalid")
	ErrTemplateRenderFailed = errors.New("template rendering failed")

	// Configuration errors
	ErrConfigInvalid    = errors.New("configuration is invalid")
	ErrConfigMissing    = errors.New("required configuration is missing")
	ErrConfigValidation = errors.New("configuration validation failed")

	// Service errors
	ErrServiceNotFound      = errors.New("service not found")
	ErrServiceInitFailed    = errors.New("service initialization failed")
	ErrServiceUnavailable   = errors.New("service is unavailable")
	ErrDependencyInjection  = errors.New("dependency injection failed")

	// Route errors
	ErrRouteNotFound    = errors.New("route not found")
	ErrRouteInvalid     = errors.New("route is invalid")
	ErrRouteConflict    = errors.New("route conflict detected")

	// Parameter errors
	ErrParameterMissing = errors.New("required parameter is missing")
	ErrParameterInvalid = errors.New("parameter is invalid")

	// Data service errors
	ErrDataServiceNotFound = errors.New("data service not found")
	ErrDataServiceFailed   = errors.New("data service operation failed")

	// Authentication errors
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrAuthorizationFailed  = errors.New("authorization failed")
	ErrSessionInvalid       = errors.New("session is invalid")

	// I18n errors
	ErrLocaleNotSupported = errors.New("locale is not supported")
	ErrTranslationMissing = errors.New("translation is missing")
)

// ErrorType represents different categories of errors for consistent handling
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "validation"
	ErrorTypeConfiguration ErrorType = "configuration"
	ErrorTypeService        ErrorType = "service"
	ErrorTypeTemplate       ErrorType = "template"
	ErrorTypeRoute          ErrorType = "route"
	ErrorTypeAuth           ErrorType = "auth"
	ErrorTypeI18n           ErrorType = "i18n"
	ErrorTypeInternal       ErrorType = "internal"
)

// AppError represents a structured application error with context
type AppError struct {
	Type        ErrorType `json:"type"`
	Code        string    `json:"code"`
	Message     string    `json:"message"`
	Details     string    `json:"details,omitempty"`
	Cause       error     `json:"-"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause for error wrapping
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new structured application error
func NewAppError(errorType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// WithDetails adds details to an AppError
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithCause adds a cause to an AppError
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// WithContext adds context information to an AppError
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// Common error constructors for consistent error creation

// NewValidationError creates a validation error
func NewValidationError(message string, details ...string) *AppError {
	err := NewAppError(ErrorTypeValidation, "VALIDATION_FAILED", message)
	if len(details) > 0 {
		err.WithDetails(details[0])
	}
	return err
}

// NewConfigurationError creates a configuration error
func NewConfigurationError(message string, details ...string) *AppError {
	err := NewAppError(ErrorTypeConfiguration, "CONFIG_ERROR", message)
	if len(details) > 0 {
		err.WithDetails(details[0])
	}
	return err
}

// NewServiceError creates a service error
func NewServiceError(message string, details ...string) *AppError {
	err := NewAppError(ErrorTypeService, "SERVICE_ERROR", message)
	if len(details) > 0 {
		err.WithDetails(details[0])
	}
	return err
}

// NewTemplateError creates a template error
func NewTemplateError(message string, details ...string) *AppError {
	err := NewAppError(ErrorTypeTemplate, "TEMPLATE_ERROR", message)
	if len(details) > 0 {
		err.WithDetails(details[0])
	}
	return err
}

// NewRouteError creates a route error
func NewRouteError(message string, details ...string) *AppError {
	err := NewAppError(ErrorTypeRoute, "ROUTE_ERROR", message)
	if len(details) > 0 {
		err.WithDetails(details[0])
	}
	return err
}

// NewDependencyInjectionError creates a dependency injection error
func NewDependencyInjectionError(message string, details ...string) *AppError {
	err := NewAppError(ErrorTypeService, "DI_ERROR", message)
	if len(details) > 0 {
		err.WithDetails(details[0])
	}
	return err
}

// WrapError wraps an existing error with additional context
func WrapError(err error, errorType ErrorType, code, message string) *AppError {
	return NewAppError(errorType, code, message).WithCause(err)
}

// IsErrorType checks if an error is of a specific type
func IsErrorType(err error, errorType ErrorType) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == errorType
	}
	return false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return "UNKNOWN_ERROR"
}