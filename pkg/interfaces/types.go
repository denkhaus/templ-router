package interfaces

import (
	"context"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

type ConfigLoader interface {
	LoadRouteConfig(templateFile string) (*shared.ConfigFile, error)
	LoadConfig(templatePath string) (*shared.ConfigFile, error)
	LoadAuthSettings(templatePath string) (*shared.AuthConfig, error)
}

type Route struct {
	// Core routing information
	Path         string `json:"path"`
	Handler      string `json:"handler,omitempty"`
	TemplateFile string `json:"template_file"`
	MetadataFile string `json:"metadata_file,omitempty"`

	// Dynamic routing
	IsDynamic  bool `json:"is_dynamic"`
	Precedence int  `json:"precedence,omitempty"`

	// Internationalization
	Locale string `json:"locale,omitempty"`

	// Security
	AuthConfig *shared.AuthConfig `json:"auth_settings,omitempty"`

	// Data Service Integration
	RequiresDataService  bool   `json:"requires_data_service,omitempty"`
	DataServiceInterface string `json:"data_service_interface,omitempty"`
	DataParameterType    string `json:"data_parameter_type,omitempty"`
}

// LayoutTemplate represents a layout template
type LayoutTemplate struct {
	FilePath      string `json:"file_path"`
	YamlPath      string `json:"yaml_path,omitempty"`
	ComponentName string `json:"component_name,omitempty"`
	Content       string `json:"content,omitempty"`
	LayoutLevel   int    `json:"layout_level,omitempty"`
	// DirectoryPath is the directory containing this layout
	DirectoryPath string
}

// ErrorTemplate represents an error template

type ErrorTemplate struct {
	// FilePath is the full path to the error.templ file
	FilePath      string `json:"file_path"`
	ComponentName string `json:"component_name,omitempty"`
	Content       string `json:"content,omitempty"`
	ErrorCode     int    `json:"error_code,omitempty"`
	// ErrorTypes is a list of error types handled by this template
	ErrorTypes []string
	// DirectoryPath is the directory containing this error template
	DirectoryPath string
	// PrecedenceLevel is the level of precedence (closer templates override further ones)
	PrecedenceLevel int
	// ErrorMessages contains mapping of error codes to specific messages
	ErrorMessages map[int]string
}

// AuthType represents different authentication types
type AuthType int

const (
	AuthTypePublic AuthType = iota
	AuthTypeUser
	AuthTypeAdmin
)

// String returns the string representation of AuthType
func (at AuthType) String() string {
	switch at {
	case AuthTypePublic:
		return "public"
	case AuthTypeUser:
		return "user"
	case AuthTypeAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// AuthResult contains authentication result (generic)
type AuthResult struct {
	IsAuthenticated bool       `json:"is_authenticated"`
	User            UserEntity `json:"user,omitempty"`
	RedirectURL     string     `json:"redirect_url,omitempty"`
	ErrorMessage    string     `json:"error_message,omitempty"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Valid     bool      `json:"valid"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Template struct {
	// File information
	FilePath      string `json:"file_path"`
	FileName      string `json:"file_name"`
	DirectoryPath string `json:"directory_path"`

	// Template metadata
	Type          string                 `json:"type,omitempty"`
	ComponentName string                 `json:"component_name,omitempty"`
	Content       string                 `json:"content,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
}

// InternationalizationIdentifier represents a structured key for translations
type InternationalizationIdentifier struct {
	// Key is the identifier key (e.g., "admin.dashboard.create.title")
	Key string

	// Source is the source of the key ("opinionated-schema" or "yaml-metadata")
	Source string

	// TemplatePath is the path to the template that uses this identifier
	TemplatePath string

	// DefaultValue is the default value if translation is missing
	DefaultValue string

	// Locales contains translations for different locales
	Locales map[string]string
}

type RouteDiscovery interface {
	DiscoverRoutes(scanPath string) ([]Route, error)
	DiscoverLayouts(scanPath string) ([]LayoutTemplate, error)
	DiscoverErrorTemplates(scanPath string) ([]ErrorTemplate, error)
}

// RouteRegistrar defines the contract for route registration
type RouteRegistrar interface {
	RegisterRoutes(routes []Route) error
	RegisterStaticRoutes()
	Register404Handler()
	RegisterMethodNotAllowedHandler()
}

// AuthService handles authentication and authorization
type AuthService interface {
	Authenticate(req *http.Request, requirements *shared.AuthConfig) (*AuthResult, error)
	HasRequiredPermissions(req *http.Request, settings *shared.AuthConfig) bool
}

// I18nService handles internationalization
// I18nService handles internationalization
type I18nService interface {
	ExtractLocale(req *http.Request) string
	CreateContext(ctx context.Context, templatePath string) context.Context
	GetSupportedLocales() []string
	LoadAllTranslations(templatePaths []string) error
	LoadComponentTranslationsIntoContext(ctx context.Context, componentName string) context.Context
}

// TemplateService handles template rendering
type TemplateService interface {
	RenderComponent(route Route, routerCtx RouterContext, ctx context.Context) (templ.Component, error)
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

// ComponentMetadataService handles component metadata loading and caching
type ComponentMetadataService interface {
	// LoadComponentMetadata loads metadata for a specific component by name
	// componentName is the name without extension (e.g., "footer", "navbar")
	LoadComponentMetadata(componentName string) (*shared.ConfigFile, error)

	// GetCachedMetadata returns cached metadata if available
	GetCachedMetadata(componentName string) (*shared.ConfigFile, bool)

	// LoadComponentTranslations loads i18n translations for a specific component and locale
	LoadComponentTranslations(componentName, locale string) (map[string]string, error)

	// GetCachedTranslations returns cached translations if available
	GetCachedTranslations(componentName, locale string) (map[string]string, bool)

	// DetectComponentsFromTemplate parses a template to find component usage
	DetectComponentsFromTemplate(templateContent string) ([]string, error)

	// LoadMultipleComponentMetadata loads metadata for multiple components efficiently
	LoadMultipleComponentMetadata(componentNames []string) (map[string]*shared.ConfigFile, error)
}

// AuthMiddlewareInterface handles authentication middleware
type AuthMiddlewareInterface interface {
	Handle(next http.Handler, requirements *shared.AuthConfig) http.Handler
}

// I18nMiddlewareInterface handles internationalization middleware
type I18nMiddlewareInterface interface {
	Handle(next http.Handler, templatePath string) http.Handler
}

// TemplateMiddlewareInterface handles template rendering middleware
type TemplateMiddlewareInterface interface {
	Handle(route Route, params map[string]string) http.Handler
}

// RouterMiddlewareInterface handles router-level middleware configuration
type RouterMiddlewareInterface interface {
	Configure(chiRouter *chi.Mux) error
}

// FileSystemChecker provides filesystem operations for library-agnostic file access
type FileSystemChecker interface {
	FileExists(path string) bool
	IsDirectory(path string) bool
	WalkDirectory(root string, walkFn func(path string, isDir bool, err error) error) error
}

// MiddlewareSetup defines the contract for middleware configuration
type MiddlewareSetup interface {
	GetAuthService() AuthService
	GetI18nService() I18nService
	GetTemplateService() TemplateService
	GetLayoutService() LayoutService
	GetErrorService() ErrorService
	GetAuthMiddleware() AuthMiddlewareInterface
	GetI18nMiddleware() I18nMiddlewareInterface
	GetTemplateMiddleware() TemplateMiddlewareInterface
	GetRouterMiddleware() RouterMiddlewareInterface
	ConfigureMiddlewareChain(route Route, authSettings interface{}) []interface{}
	ValidateMiddlewareSetup() error
}

// CustomMiddlewareDefinition holds a custom middleware with its definition order
type CustomMiddlewareDefinition struct {
	Name  string
	Func  func(http.Handler) http.Handler
	Order int // Definition order
}

// HandlerBuilder defines the contract for handler building
type HandlerBuilder interface {
	BuildHandler(route Route) http.Handler
	BuildStaticHandler(path string) http.Handler
	BuildErrorHandler(statusCode int, message string) http.HandlerFunc
}
type RouterCore interface {
	Initialize() error
	RegisterRoutes(chiRouter *chi.Mux) error
	GetRoutes() []Route
	GetLayoutTemplates() []LayoutTemplate
	GetErrorTemplates() []ErrorTemplate
	GetMiddlewareSetup() MiddlewareSetup
	GetHandlerBuilder() HandlerBuilder
	GetRouteRegistrar() RouteRegistrar
}

// ApplicationOption defines an option for configuring application services
type ApplicationOption func(c Container)

// Container defines the interface for DI container operations
type Container interface {
	GetInjector() do.Injector
	GetRouter() RouterCore
	GetRouterBootstrap() interface{}
	GetLogger() *zap.Logger
	GetConfigService() ConfigService
	RegisterApplicationServices(options ...ApplicationOption)
	Shutdown() error
}
