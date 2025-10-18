package interfaces

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type AssetsService interface {
	SetupRoutesWithRouter(mux chi.Router)
	SetupRoutes(mux *chi.Mux)
}

// Import central type definitions to eliminate struct redundancy
// All Route, Template, Auth, User, Session types are now centralized in types.go

// AuthService handles authentication and authorization
type AuthService interface {
	Authenticate(req *http.Request, requirements *AuthSettings) (*AuthResult, error)
	HasRequiredPermissions(req *http.Request, settings *AuthSettings) bool
}

// I18nService handles internationalization
type I18nService interface {
	ExtractLocale(req *http.Request) string
	CreateContext(ctx context.Context, locale, templatePath string) context.Context
	GetSupportedLocales() []string
}

// TemplateService handles template rendering
type TemplateService interface {
	RenderComponent(route Route, ctx context.Context, params map[string]string) (templ.Component, error)
	RenderLayoutComponent(layoutPath string, content templ.Component, ctx context.Context) (templ.Component, error)
}

// LayoutService handles layout resolution and wrapping
type LayoutService interface {
	FindLayoutForTemplate(templatePath string) *LayoutTemplate
	WrapInLayout(component templ.Component, layout *LayoutTemplate, ctx context.Context) templ.Component
}

// ErrorService handles error template resolution
type ErrorService interface {
	FindErrorTemplateForPath(path string) *ErrorTemplate
	CreateErrorComponent(message, path string) templ.Component
}

// ValidationService handles unified validation of routes and configurations
type ValidationService interface {
	ValidateConfiguration(routes []Route, configs map[string]*ConfigFile) error
}

// SessionStore interface for session management (pluggable)
type SessionStore interface {
	GetSession(req *http.Request) (*Session, error)
	CreateSession(userID string) (*Session, error)
	DeleteSession(sessionID string) error
}

// UserStore interface for user management (pluggable)
type UserStore interface {
	GetUserByID(userID string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ValidateCredentials(email, password string) (*User, error)
	CreateUser(username, email, password string) (*User, error)
	UserExists(username, email string) (bool, error)
}
