package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/denkhaus/templ-router/pkg/router/pipeline"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// Shared mock implementations for router tests

type MockConfigService struct {
	LayoutRootDir string
}

func (m *MockConfigService) GetLayoutRootDirectory() string {
	if m.LayoutRootDir != "" {
		return m.LayoutRootDir
	}
	return "/app"
}
func (m *MockConfigService) GetServerHost() string                     { return "localhost" }
func (m *MockConfigService) GetServerPort() int                        { return 8080 }
func (m *MockConfigService) GetServerBaseURL() string                  { return "http://localhost:8080" }
func (m *MockConfigService) GetSupportedLocales() []string             { return []string{"en", "de"} }
func (m *MockConfigService) GetDefaultLocale() string                  { return "en" }
func (m *MockConfigService) GetFallbackLocale() string                 { return "en" }
func (m *MockConfigService) GetLayoutFileName() string                 { return "layout.templ" }
func (m *MockConfigService) GetLayoutAssetsDirectory() string          { return "assets" }
func (m *MockConfigService) GetLayoutAssetsRouteName() string          { return "/assets/" }
func (m *MockConfigService) GetTemplateExtension() string              { return ".templ" }
func (m *MockConfigService) GetMetadataExtension() string              { return ".yaml" }
func (m *MockConfigService) IsLayoutInheritanceEnabled() bool          { return true }
func (m *MockConfigService) GetTemplateOutputDir() string              { return "generated" }
func (m *MockConfigService) GetTemplatePackageName() string            { return "templates" }
func (m *MockConfigService) IsDevelopment() bool                       { return true }
func (m *MockConfigService) IsProduction() bool                        { return false }

// Router configuration methods
func (m *MockConfigService) GetRouterEnableTrailingSlash() bool     { return true }
func (m *MockConfigService) GetRouterEnableSlashRedirect() bool     { return true }
func (m *MockConfigService) GetRouterEnableMethodNotAllowed() bool  { return true }
func (m *MockConfigService) GetServerReadTimeout() time.Duration       { return 30 * time.Second }
func (m *MockConfigService) GetServerWriteTimeout() time.Duration      { return 30 * time.Second }
func (m *MockConfigService) GetServerIdleTimeout() time.Duration       { return 60 * time.Second }
func (m *MockConfigService) GetServerShutdownTimeout() time.Duration   { return 10 * time.Second }
func (m *MockConfigService) GetDatabaseHost() string                   { return "localhost" }
func (m *MockConfigService) GetDatabasePort() int                      { return 5432 }
func (m *MockConfigService) GetDatabaseUser() string                   { return "user" }
func (m *MockConfigService) GetDatabasePassword() string               { return "password" }
func (m *MockConfigService) GetDatabaseName() string                   { return "testdb" }
func (m *MockConfigService) GetDatabaseSSLMode() string                { return "disable" }
func (m *MockConfigService) IsEmailVerificationRequired() bool         { return false }
func (m *MockConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *MockConfigService) GetSessionCookieName() string              { return "session" }
func (m *MockConfigService) GetSessionExpiry() time.Duration           { return 24 * time.Hour }
func (m *MockConfigService) IsSessionSecure() bool                     { return false }
func (m *MockConfigService) IsSessionHttpOnly() bool                   { return true }
func (m *MockConfigService) GetSessionSameSite() string                { return "Lax" }
func (m *MockConfigService) GetMinPasswordLength() int                 { return 8 }
func (m *MockConfigService) IsStrongPasswordRequired() bool            { return false }
func (m *MockConfigService) ShouldCreateDefaultAdmin() bool            { return false }
func (m *MockConfigService) GetDefaultAdminEmail() string              { return "" }
func (m *MockConfigService) GetDefaultAdminPassword() string           { return "" }
func (m *MockConfigService) GetDefaultAdminFirstName() string          { return "" }
func (m *MockConfigService) GetDefaultAdminLastName() string           { return "" }
func (m *MockConfigService) GetCSRFSecret() string                     { return "secret" }
func (m *MockConfigService) IsCSRFSecure() bool                        { return false }
func (m *MockConfigService) IsCSRFHttpOnly() bool                      { return true }
func (m *MockConfigService) GetCSRFSameSite() string                   { return "Lax" }
func (m *MockConfigService) IsRateLimitEnabled() bool                  { return false }
func (m *MockConfigService) GetRateLimitRequests() int                 { return 100 }
func (m *MockConfigService) AreSecurityHeadersEnabled() bool           { return false }
func (m *MockConfigService) IsHSTSEnabled() bool                       { return false }
func (m *MockConfigService) GetHSTSMaxAge() int                        { return 31536000 }
func (m *MockConfigService) GetLogLevel() string                       { return "info" }
func (m *MockConfigService) GetLogFormat() string                      { return "json" }
func (m *MockConfigService) GetLogOutput() string                      { return "stdout" }
func (m *MockConfigService) IsFileLoggingEnabled() bool                { return false }
func (m *MockConfigService) GetLogFilePath() string                    { return "" }
func (m *MockConfigService) GetSMTPHost() string                       { return "" }
func (m *MockConfigService) GetSMTPPort() int                          { return 587 }
func (m *MockConfigService) GetSMTPUsername() string                   { return "" }
func (m *MockConfigService) GetSMTPPassword() string                   { return "" }
func (m *MockConfigService) IsSMTPTLSEnabled() bool                    { return true }
func (m *MockConfigService) GetFromEmail() string                      { return "" }
func (m *MockConfigService) GetFromName() string                       { return "" }
func (m *MockConfigService) GetReplyToEmail() string                   { return "" }
func (m *MockConfigService) IsEmailDummyModeEnabled() bool             { return true }

type MockAssetsService struct{}

func (m *MockAssetsService) SetupRoutesWithRouter(mux chi.Router) {}
func (m *MockAssetsService) SetupRoutes(mux *chi.Mux)             {}

type MockRouteDiscovery struct {
	Routes         []interfaces.Route
	Layouts        []LayoutTemplate
	ErrorTemplates []ErrorTemplate
	ShouldError    bool
}

func (m *MockRouteDiscovery) DiscoverRoutes(scanPath string) ([]interfaces.Route, error) {
	if m.ShouldError {
		return nil, errors.New("mock discovery error")
	}
	return m.Routes, nil
}

func (m *MockRouteDiscovery) DiscoverLayouts(scanPath string) ([]LayoutTemplate, error) {
	if m.ShouldError {
		return nil, errors.New("mock layout discovery error")
	}
	return m.Layouts, nil
}

func (m *MockRouteDiscovery) DiscoverErrorTemplates(scanPath string) ([]ErrorTemplate, error) {
	if m.ShouldError {
		return nil, errors.New("mock error template discovery error")
	}
	return m.ErrorTemplates, nil
}

type MockConfigLoader struct {
	ShouldError  bool
	Config       *interfaces.ConfigFile
	AuthSettings *interfaces.AuthSettings
}

func (m *MockConfigLoader) LoadRouteConfig(templateFile string) (*interfaces.ConfigFile, error) {
	if m.ShouldError {
		return nil, errors.New("mock config load error")
	}
	if m.Config != nil {
		return m.Config, nil
	}
	return &interfaces.ConfigFile{FilePath: templateFile}, nil
}

func (m *MockConfigLoader) LoadConfig(templatePath string) (*interfaces.ConfigFile, error) {
	if m.ShouldError {
		return nil, errors.New("mock config load error")
	}
	if m.Config != nil {
		return m.Config, nil
	}
	return &interfaces.ConfigFile{FilePath: templatePath}, nil
}

func (m *MockConfigLoader) LoadAuthSettings(templatePath string) (*interfaces.AuthSettings, error) {
	if m.ShouldError {
		return nil, errors.New("mock auth settings load error")
	}
	if m.AuthSettings != nil {
		return m.AuthSettings, nil
	}
	return &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}, nil
}

type MockHandlerPipeline struct {
	ShouldError bool
}

func (m *MockHandlerPipeline) BuildHandler(config pipeline.PipelineConfig) http.Handler {
	if m.ShouldError {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Pipeline error", http.StatusInternalServerError)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Handler response"))
	})
}

// Service mocks
type MockAuthService struct{}

func (m *MockAuthService) Authenticate(req *http.Request, requirements *interfaces.AuthSettings) (*interfaces.AuthResult, error) {
	return &interfaces.AuthResult{IsAuthenticated: true}, nil
}

func (m *MockAuthService) HasRequiredPermissions(req *http.Request, settings *interfaces.AuthSettings) bool {
	return true
}

type MockI18nService struct{}

func (m *MockI18nService) ExtractLocale(req *http.Request) string {
	return "en"
}

func (m *MockI18nService) CreateContext(ctx context.Context, locale, templatePath string) context.Context {
	return ctx
}

func (m *MockI18nService) GetSupportedLocales() []string {
	return []string{"en", "de"}
}

func (m *MockI18nService) LoadAllTranslations(templatePaths []string) error {
	return nil
}

type MockTemplateService struct{}

func (m *MockTemplateService) RenderComponent(route interfaces.Route, ctx context.Context, params map[string]string) (templ.Component, error) {
	return templ.Raw("test content"), nil
}

func (m *MockTemplateService) RenderLayoutComponent(layoutPath string, content templ.Component, ctx context.Context) (templ.Component, error) {
	return templ.Raw("layout content"), nil
}

type MockLayoutService struct{}

func (m *MockLayoutService) FindLayoutForTemplate(templatePath string) *interfaces.LayoutTemplate {
	return &interfaces.LayoutTemplate{FilePath: "/app/layout.templ"}
}

func (m *MockLayoutService) WrapInLayout(component templ.Component, layout *interfaces.LayoutTemplate, ctx context.Context) templ.Component {
	return templ.Raw("wrapped content")
}

type MockErrorService struct{}

func (m *MockErrorService) FindErrorTemplateForPath(path string) *interfaces.ErrorTemplate {
	return &interfaces.ErrorTemplate{FilePath: "/app/error.templ"}
}

func (m *MockErrorService) CreateErrorComponent(message, path string) templ.Component {
	return templ.Raw("error content")
}

// Middleware mocks
type MockAuthMiddleware struct{}

func (m *MockAuthMiddleware) Handle(next http.Handler, requirements *interfaces.AuthSettings) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

type MockI18nMiddleware struct{}

func (m *MockI18nMiddleware) Handle(next http.Handler, templatePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

type MockTemplateMiddleware struct{}

func (m *MockTemplateMiddleware) Handle(route interfaces.Route, params map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock template response"))
	})
}

type MockRouterMiddleware struct{}

func (m *MockRouterMiddleware) ConfigureRouterMiddleware(chiRouter *chi.Mux) error {
	// Mock implementation - do nothing
	return nil
}

// Helper function to create a complete test DI container
func CreateTestContainer() do.Injector {
	injector := do.New()

	// Register core services
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &MockConfigService{}, nil
	})

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	do.Provide(injector, func(i do.Injector) (*pipeline.HandlerPipeline, error) {
		// Return a real HandlerPipeline since we can't mock it easily
		return &pipeline.HandlerPipeline{}, nil
	})

	do.Provide(injector, func(i do.Injector) (RouteDiscovery, error) {
		return &MockRouteDiscovery{
			Routes: []interfaces.Route{
				{Path: "/", TemplateFile: "index.templ", IsDynamic: false},
				{Path: "/about", TemplateFile: "about.templ", IsDynamic: false},
				{Path: "/user/{id}", TemplateFile: "user.templ", IsDynamic: true},
			},
			Layouts: []LayoutTemplate{
				{FilePath: "/app/layout.templ", DirectoryPath: "/app"},
			},
			ErrorTemplates: []ErrorTemplate{
				{FilePath: "/app/error.templ", DirectoryPath: "/app", ErrorTypes: []string{"404"}},
			},
		}, nil
	})

	do.Provide(injector, func(i do.Injector) (interfaces.AssetsService, error) {
		return &MockAssetsService{}, nil
	})

	do.Provide(injector, func(i do.Injector) (ConfigLoader, error) {
		return &MockConfigLoader{}, nil
	})

	// Register services for middleware
	do.Provide(injector, func(i do.Injector) (interfaces.AuthService, error) {
		return &MockAuthService{}, nil
	})

	do.Provide(injector, func(i do.Injector) (interfaces.I18nService, error) {
		return &MockI18nService{}, nil
	})

	do.Provide(injector, func(i do.Injector) (interfaces.TemplateService, error) {
		return &MockTemplateService{}, nil
	})

	do.Provide(injector, func(i do.Injector) (interfaces.LayoutService, error) {
		return &MockLayoutService{}, nil
	})

	do.Provide(injector, func(i do.Injector) (interfaces.ErrorService, error) {
		return &MockErrorService{}, nil
	})

	// Register all middleware interfaces for complete testing
	do.Provide(injector, func(i do.Injector) (middleware.AuthMiddlewareInterface, error) {
		return &MockAuthMiddleware{}, nil
	})

	do.Provide(injector, func(i do.Injector) (middleware.I18nMiddlewareInterface, error) {
		return &MockI18nMiddleware{}, nil
	})

	do.Provide(injector, func(i do.Injector) (middleware.TemplateMiddlewareInterface, error) {
		return &MockTemplateMiddleware{}, nil
	})

	do.Provide(injector, func(i do.Injector) (middleware.RouterMiddlewareInterface, error) {
		return &MockRouterMiddleware{}, nil
	})

	return injector
}

// Add missing methods to MockConfigService for auth redirect routes
func (m *MockConfigService) GetSignInSuccessRoute() string  { return "/dashboard" }
func (m *MockConfigService) GetSignUpSuccessRoute() string  { return "/welcome" }
func (m *MockConfigService) GetSignOutSuccessRoute() string { return "/" }

// Auth routes
func (m *MockConfigService) GetSignInRoute() string { return "/login" }
