package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
)

// Common errors
var (
	ErrTemplateNotFound = errors.New("template not found")
)

// AuthService handles authentication and authorization
type AuthService interface {
	Authenticate(req *http.Request, requirements *interfaces.AuthSettings) (*interfaces.AuthResult, error)
	HasRequiredPermissions(req *http.Request, settings *interfaces.AuthSettings) bool
}

// I18nService handles internationalization
type I18nService interface {
	ExtractLocale(req *http.Request) string
	CreateContext(ctx context.Context, locale, templatePath string) context.Context
}

// TemplateService handles template rendering
type TemplateService interface {
	RenderComponent(route interfaces.Route, ctx context.Context, params map[string]string) (templ.Component, error)
	RenderLayoutComponent(layoutPath string, content templ.Component, ctx context.Context) (templ.Component, error)
}

// LayoutService handles layout resolution and wrapping
type LayoutService interface {
	FindLayoutForTemplate(templatePath string) *interfaces.LayoutTemplate
	WrapInLayout(component templ.Component, layout *interfaces.LayoutTemplate, ctx context.Context) templ.Component
}

// ErrorService handles error template resolution
type ErrorService interface {
	CreateErrorComponent(message, path string) templ.Component
}

// AuthMiddlewareInterface handles authentication middleware
type AuthMiddlewareInterface interface {
	Handle(next http.Handler, requirements *interfaces.AuthSettings) http.Handler
}

// I18nMiddlewareInterface handles internationalization middleware
type I18nMiddlewareInterface interface {
	Handle(next http.Handler, templatePath string) http.Handler
}

// TemplateMiddlewareInterface handles template rendering middleware
type TemplateMiddlewareInterface interface {
	Handle(route interfaces.Route, params map[string]string) http.Handler
}

// FileSystemChecker provides filesystem operations for library-agnostic file access
type FileSystemChecker interface {
	FileExists(path string) bool
	IsDirectory(path string) bool
	WalkDirectory(root string, walkFn func(path string, isDir bool, err error) error) error
}
