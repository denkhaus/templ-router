package router

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/denkhaus/templ-router/pkg/router/pipeline"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// Mock implementations for testing
type mockRouterConfigService struct{}

func (m *mockRouterConfigService) GetLayoutRootDirectory() string            { return "app" }
func (m *mockRouterConfigService) GetLayoutFileName() string                 { return "layout" }
func (m *mockRouterConfigService) GetTemplateExtension() string              { return ".templ" }
func (m *mockRouterConfigService) GetSupportedLocales() []string             { return []string{"en", "de"} }
func (m *mockRouterConfigService) GetDefaultLocale() string                  { return "en" }
func (m *mockRouterConfigService) AreSecurityHeadersEnabled() bool           { return false }
func (m *mockRouterConfigService) GetServerHost() string                     { return "localhost" }
func (m *mockRouterConfigService) GetServerPort() int                        { return 8080 }
func (m *mockRouterConfigService) GetServerBaseURL() string                  { return "http://localhost:8080" }
func (m *mockRouterConfigService) GetServerReadTimeout() time.Duration       { return 30 * time.Second }
func (m *mockRouterConfigService) GetServerWriteTimeout() time.Duration      { return 30 * time.Second }
func (m *mockRouterConfigService) GetServerIdleTimeout() time.Duration       { return 2 * time.Minute }
func (m *mockRouterConfigService) GetServerShutdownTimeout() time.Duration   { return 30 * time.Second }
func (m *mockRouterConfigService) GetFallbackLocale() string                 { return "en" }
func (m *mockRouterConfigService) GetLayoutAssetsDirectory() string          { return "assets" }
func (m *mockRouterConfigService) GetLayoutAssetsRouteName() string          { return "/assets/" }
func (m *mockRouterConfigService) GetMetadataExtension() string              { return ".yaml" }
func (m *mockRouterConfigService) IsLayoutInheritanceEnabled() bool          { return true }
func (m *mockRouterConfigService) GetTemplateOutputDir() string              { return "generated" }
func (m *mockRouterConfigService) GetTemplatePackageName() string            { return "templates" }
func (m *mockRouterConfigService) GetDatabaseHost() string                   { return "localhost" }
func (m *mockRouterConfigService) GetDatabasePort() int                      { return 5432 }
func (m *mockRouterConfigService) GetDatabaseUser() string                   { return "user" }
func (m *mockRouterConfigService) GetDatabasePassword() string               { return "pass" }
func (m *mockRouterConfigService) GetDatabaseName() string                   { return "db" }
func (m *mockRouterConfigService) GetDatabaseSSLMode() string                { return "disable" }
func (m *mockRouterConfigService) IsEmailVerificationRequired() bool         { return false }
func (m *mockRouterConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *mockRouterConfigService) GetSessionCookieName() string              { return "session" }
func (m *mockRouterConfigService) GetSessionExpiry() time.Duration           { return 24 * time.Hour }
func (m *mockRouterConfigService) IsSessionSecure() bool                     { return false }
func (m *mockRouterConfigService) IsSessionHttpOnly() bool                   { return true }
func (m *mockRouterConfigService) GetSessionSameSite() string                { return "lax" }
func (m *mockRouterConfigService) GetMinPasswordLength() int                 { return 8 }
func (m *mockRouterConfigService) IsStrongPasswordRequired() bool            { return false }
func (m *mockRouterConfigService) ShouldCreateDefaultAdmin() bool            { return false }
func (m *mockRouterConfigService) GetDefaultAdminEmail() string              { return "admin@example.com" }
func (m *mockRouterConfigService) GetDefaultAdminPassword() string           { return "password" }
func (m *mockRouterConfigService) GetDefaultAdminFirstName() string          { return "Admin" }
func (m *mockRouterConfigService) GetDefaultAdminLastName() string           { return "User" }

// Auth routes
func (m *mockRouterConfigService) GetSignInRoute() string { return "/login" }

// Auth redirect routes (only for success cases)
func (m *mockRouterConfigService) GetSignInSuccessRoute() string  { return "/dashboard" }
func (m *mockRouterConfigService) GetSignUpSuccessRoute() string  { return "/welcome" }
func (m *mockRouterConfigService) GetSignOutSuccessRoute() string { return "/" }
func (m *mockRouterConfigService) GetSMTPHost() string            { return "" }
func (m *mockRouterConfigService) GetSMTPPort() int               { return 587 }
func (m *mockRouterConfigService) GetSMTPUsername() string        { return "" }
func (m *mockRouterConfigService) GetSMTPPassword() string        { return "" }
func (m *mockRouterConfigService) IsSMTPTLSEnabled() bool         { return true }
func (m *mockRouterConfigService) GetFromEmail() string           { return "noreply@example.com" }
func (m *mockRouterConfigService) GetFromName() string            { return "App" }
func (m *mockRouterConfigService) GetReplyToEmail() string        { return "" }
func (m *mockRouterConfigService) IsEmailDummyModeEnabled() bool  { return true }
func (m *mockRouterConfigService) GetCSRFSecret() string          { return "secret" }
func (m *mockRouterConfigService) IsCSRFSecure() bool             { return false }
func (m *mockRouterConfigService) IsCSRFHttpOnly() bool           { return true }
func (m *mockRouterConfigService) GetCSRFSameSite() string        { return "strict" }
func (m *mockRouterConfigService) IsRateLimitEnabled() bool       { return false }
func (m *mockRouterConfigService) GetRateLimitRequests() int      { return 100 }
func (m *mockRouterConfigService) IsHSTSEnabled() bool            { return false }
func (m *mockRouterConfigService) GetHSTSMaxAge() int             { return 31536000 }
func (m *mockRouterConfigService) GetLogLevel() string            { return "info" }
func (m *mockRouterConfigService) GetLogFormat() string           { return "json" }
func (m *mockRouterConfigService) GetLogOutput() string           { return "stdout" }
func (m *mockRouterConfigService) IsFileLoggingEnabled() bool     { return false }
func (m *mockRouterConfigService) GetLogFilePath() string         { return "" }
func (m *mockRouterConfigService) IsDevelopment() bool            { return true }
func (m *mockRouterConfigService) IsProduction() bool             { return false }

type mockRouterAssetsService struct{}

func (m *mockRouterAssetsService) SetupRoutes(router *chi.Mux)             {}
func (m *mockRouterAssetsService) SetupRoutesWithRouter(router chi.Router) {}

type mockRouterTemplateRegistry struct{}

func (m *mockRouterTemplateRegistry) GetTemplate(key string) (templ.Component, error) {
	return nil, nil
}
func (m *mockRouterTemplateRegistry) GetTemplateFunction(key string) (func() interface{}, bool) {
	return nil, false
}
func (m *mockRouterTemplateRegistry) GetAllTemplateKeys() []string { return []string{} }
func (m *mockRouterTemplateRegistry) IsAvailable(key string) bool  { return true }
func (m *mockRouterTemplateRegistry) GetRouteToTemplateMapping() map[string]string {
	return map[string]string{
		"/":         "template1",
		"/{locale}": "template2",
	}
}
func (m *mockRouterTemplateRegistry) GetTemplateByRoute(route string) (templ.Component, error) {
	return nil, nil
}

func (m *mockRouterTemplateRegistry) RequiresDataService(key string) bool {
	return false
}

func (m *mockRouterTemplateRegistry) GetDataServiceInfo(key string) (interfaces.DataServiceInfo, bool) {
	return interfaces.DataServiceInfo{}, false
}

type mockRouteDiscovery struct{}

func (m *mockRouteDiscovery) DiscoverRoutes(scanPath string) ([]interfaces.Route, error) {
	return []interfaces.Route{
		{Path: "/", TemplateFile: "app/page.templ", IsDynamic: false},
		{Path: "/{locale}", TemplateFile: "app/locale_/page.templ", IsDynamic: true},
	}, nil
}

func (m *mockRouteDiscovery) DiscoverLayouts(scanPath string) ([]LayoutTemplate, error) {
	return []LayoutTemplate{
		{FilePath: "app/layout.templ", DirectoryPath: "app", LayoutLevel: 0},
	}, nil
}

func (m *mockRouteDiscovery) DiscoverErrorTemplates(scanPath string) ([]ErrorTemplate, error) {
	return []ErrorTemplate{
		{FilePath: "app/error.templ", DirectoryPath: "app", ErrorTypes: []string{"404"}},
	}, nil
}

type mockConfigLoader struct{}

func (m *mockConfigLoader) LoadConfig(templatePath string) (*interfaces.ConfigFile, error) {
	return &interfaces.ConfigFile{}, nil
}

func (m *mockConfigLoader) LoadAuthSettings(templatePath string) (*interfaces.AuthSettings, error) {
	return &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}, nil
}

func (m *mockConfigLoader) LoadRouteConfig(templatePath string) (*interfaces.ConfigFile, error) {
	return &interfaces.ConfigFile{}, nil
}

type mockAuthService struct{}

func (m *mockAuthService) Authenticate(req *http.Request, requirements *interfaces.AuthSettings) (*interfaces.AuthResult, error) {
	return &interfaces.AuthResult{
		IsAuthenticated: true,
		User:            &testUser{ID: "test123", Email: "test@example.com", Roles: []string{"user"}},
	}, nil
}

func (m *mockAuthService) HasRequiredPermissions(req *http.Request, settings *interfaces.AuthSettings) bool {
	return true
}

// testUser implements UserEntity interface for router tests
type testUser struct {
	ID    string
	Email string
	Roles []string
}

func (u *testUser) GetID() string      { return u.ID }
func (u *testUser) GetEmail() string   { return u.Email }
func (u *testUser) GetRoles() []string { return u.Roles }

func (m *mockAuthService) RefreshToken(token string) (string, error) {
	return "new-mock-token", nil
}

func (m *mockAuthService) RevokeToken(token string) error {
	return nil
}

func (m *mockAuthService) ChangePassword(userID, oldPassword, newPassword string) error {
	return nil
}

func (m *mockAuthService) ResetPassword(email string) error {
	return nil
}

func (m *mockAuthService) VerifyEmail(token string) error {
	return nil
}

func (m *mockAuthService) ResendVerificationEmail(email string) error {
	return nil
}

type mockI18nService struct{}

func (m *mockI18nService) GetTranslation(key, locale string) string {
	return key // Return key as fallback
}

func (m *mockI18nService) LoadTranslations(templatePath string) error {
	return nil
}

func (m *mockI18nService) GetAvailableKeys(templatePath string) []string {
	return []string{}
}

func (m *mockI18nService) GetCurrentLocale(r *http.Request) string {
	return "en"
}

func (m *mockI18nService) CreateContext(ctx context.Context, templatePath string, locale string) context.Context {
	i18nData := &i18n.I18nData{
		Locale:          locale,
		CurrentTemplate: templatePath,
		Translations:    make(map[string]string),
		FallbackLocale:  "en",
		Logger:          zap.NewNop(),
	}
	return context.WithValue(ctx, i18n.I18nDataKey, i18nData)
}

func (m *mockI18nService) ExtractLocale(r *http.Request) string {
	return "en"
}

func (m *mockI18nService) GetSupportedLocales() []string {
	return []string{"en", "de"}
}

type mockTemplateService struct{}

func (m *mockTemplateService) GetTemplate(templatePath string) (templ.Component, error) {
	return nil, nil
}

func (m *mockTemplateService) IsTemplateAvailable(templatePath string) bool {
	return true
}

func (m *mockTemplateService) GetTemplateKey(templatePath string) (string, error) {
	return "mock-key", nil
}

func (m *mockTemplateService) RenderComponent(route interfaces.Route, ctx context.Context, params map[string]string) (templ.Component, error) {
	return nil, nil
}

func (m *mockTemplateService) RenderLayoutComponent(layoutPath string, content templ.Component, ctx context.Context) (templ.Component, error) {
	return content, nil
}

type mockLayoutService struct{}

func (m *mockLayoutService) GetLayoutForRoute(route interfaces.Route) (interfaces.LayoutTemplate, error) {
	return interfaces.LayoutTemplate{
		FilePath:    "app/layout.templ",
		LayoutLevel: 0,
	}, nil
}

func (m *mockLayoutService) RenderWithLayout(route interfaces.Route, content templ.Component, ctx context.Context) (templ.Component, error) {
	return content, nil
}

func (m *mockLayoutService) FindLayoutForTemplate(templatePath string) *interfaces.LayoutTemplate {
	return &interfaces.LayoutTemplate{
		FilePath:    "app/layout.templ",
		LayoutLevel: 0,
	}
}

func (m *mockLayoutService) WrapInLayout(component templ.Component, layout *interfaces.LayoutTemplate, ctx context.Context) templ.Component {
	return component
}

// Mock ErrorService
type mockErrorService struct{}

func (m *mockErrorService) FindErrorTemplateForPath(path string) *interfaces.ErrorTemplate {
	return &interfaces.ErrorTemplate{
		FilePath:  "app/error.templ",
		ErrorCode: 404,
	}
}

func (m *mockErrorService) CreateErrorComponent(message, path string) templ.Component {
	return nil
}

// Mock Middleware Interfaces
type mockAuthMiddleware struct{}

func (m *mockAuthMiddleware) Handle(next http.Handler, requirements *interfaces.AuthSettings) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (m *mockAuthMiddleware) AuthenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (m *mockAuthMiddleware) RequireAuth(authSettings *interfaces.AuthSettings) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

type mockI18nMiddleware struct{}

func (m *mockI18nMiddleware) Handle(next http.Handler, templatePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (m *mockI18nMiddleware) ExtractLocale(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (m *mockI18nMiddleware) SetupI18nContext(templatePath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

type mockTemplateMiddleware struct{}

func (m *mockTemplateMiddleware) Handle(route interfaces.Route, params map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock template response"))
	})
}

func (m *mockTemplateMiddleware) RenderTemplate(route interfaces.Route) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func (m *mockTemplateMiddleware) HandleTemplateError(err error, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Template error", http.StatusInternalServerError)
}

// Mock AuthHandlers for router tests
type mockRouterAuthHandlers struct{}

func (m *mockRouterAuthHandlers) RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc)) {
	// Register mock auth routes
	registerFunc("GET", "/signin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login page"))
	})
	registerFunc("POST", "/signin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login processed"))
	})
	registerFunc("GET", "/signup", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Signup page"))
	})
	registerFunc("POST", "/signout", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged out"))
	})
}

func (m *mockRouterAuthHandlers) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login handled"))
}

func (m *mockRouterAuthHandlers) HandleSignOut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout handled"))
}

func (m *mockRouterAuthHandlers) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signup handled"))
}

func createRouterTestContainer() do.Injector {
	injector := do.New()

	// Register all required dependencies
	do.ProvideValue[interfaces.ConfigService](injector, &mockRouterConfigService{})
	do.ProvideValue[*zap.Logger](injector, zap.NewNop())
	do.ProvideValue[interfaces.AssetsService](injector, &mockRouterAssetsService{})
	do.ProvideValue[interfaces.TemplateRegistry](injector, &mockRouterTemplateRegistry{})
	do.ProvideValue[RouteDiscovery](injector, &mockRouteDiscovery{})
	do.ProvideValue[ConfigLoader](injector, &mockConfigLoader{})
	do.ProvideValue[interfaces.AuthService](injector, &mockAuthService{})
	do.ProvideValue[interfaces.I18nService](injector, &mockI18nService{})
	do.ProvideValue[interfaces.TemplateService](injector, &mockTemplateService{})
	do.ProvideValue[interfaces.LayoutService](injector, &mockLayoutService{})

	// Register ErrorService (required by middleware setup)
	do.ProvideValue[interfaces.ErrorService](injector, &mockErrorService{})

	// Register middleware interfaces (required by middleware setup)
	do.ProvideValue[middleware.AuthMiddlewareInterface](injector, &mockAuthMiddleware{})
	do.ProvideValue[middleware.I18nMiddlewareInterface](injector, &mockI18nMiddleware{})
	do.ProvideValue[middleware.TemplateMiddlewareInterface](injector, &mockTemplateMiddleware{})

	// Register AuthHandlers (required by RegisterRoutes)
	do.ProvideValue[interfaces.AuthHandlers](injector, &mockRouterAuthHandlers{})

	// Create and register handler pipeline
	do.Provide(injector, pipeline.NewHandlerPipeline)

	return injector
}

func TestNewCleanRouterCore(t *testing.T) {
	injector := createRouterTestContainer()

	router, err := NewCleanRouterCore(injector)
	if err != nil {
		t.Fatalf("NewCleanRouterCore() returned error: %v", err)
	}

	if router == nil {
		t.Fatal("NewCleanRouterCore() returned nil")
	}
}

func TestCleanRouterCoreInitialize(t *testing.T) {
	injector := createRouterTestContainer()
	router, err := NewCleanRouterCore(injector)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	err = router.Initialize()
	if err != nil {
		t.Fatalf("Initialize() returned error: %v", err)
	}

	// Verify routes were discovered
	routes := router.GetRoutes()
	if len(routes) == 0 {
		t.Error("Initialize() did not discover any routes")
	}

	// Verify layouts were discovered
	layouts := router.GetLayoutTemplates()
	if len(layouts) == 0 {
		t.Error("Initialize() did not discover any layouts")
	}

	// Verify error templates were discovered
	errorTemplates := router.GetErrorTemplates()
	if len(errorTemplates) == 0 {
		t.Error("Initialize() did not discover any error templates")
	}
}

func TestCleanRouterCoreRegisterRoutes(t *testing.T) {
	injector := createRouterTestContainer()
	router, err := NewCleanRouterCore(injector)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	// Initialize first
	err = router.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize router: %v", err)
	}

	// Create Chi router
	chiRouter := chi.NewRouter()

	// Register routes
	err = router.RegisterRoutes(chiRouter)
	if err != nil {
		t.Fatalf("RegisterRoutes() returned error: %v", err)
	}

	// Verify route registrar was created
	registrar := router.GetRouteRegistrar()
	if registrar == nil {
		t.Error("RegisterRoutes() did not create route registrar")
	}
}

func TestCleanRouterCoreGetters(t *testing.T) {
	injector := createRouterTestContainer()
	router, err := NewCleanRouterCore(injector)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	// Initialize to populate data
	err = router.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize router: %v", err)
	}

	// Test GetRoutes
	routes := router.GetRoutes()
	if routes == nil {
		t.Error("GetRoutes() returned nil")
	}

	// Test GetLayoutTemplates
	layouts := router.GetLayoutTemplates()
	if layouts == nil {
		t.Error("GetLayoutTemplates() returned nil")
	}

	// Test GetErrorTemplates
	errorTemplates := router.GetErrorTemplates()
	if errorTemplates == nil {
		t.Error("GetErrorTemplates() returned nil")
	}

	// Test GetMiddlewareSetup
	middlewareSetup := router.GetMiddlewareSetup()
	if middlewareSetup == nil {
		t.Error("GetMiddlewareSetup() returned nil")
	}

	// Test GetHandlerBuilder
	handlerBuilder := router.GetHandlerBuilder()
	if handlerBuilder == nil {
		t.Error("GetHandlerBuilder() returned nil")
	}
}

func TestConvertToInterfaceRoutes(t *testing.T) {
	injector := createRouterTestContainer()
	router, err := NewCleanRouterCore(injector)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	// Create test routes
	testRoutes := []interfaces.Route{
		{Path: "/test", TemplateFile: "test.templ", IsDynamic: false},
		{Path: "/dynamic/{id}", TemplateFile: "dynamic.templ", IsDynamic: true},
	}

	// Access the private method through type assertion
	crc := router.(*cleanRouterCore)
	interfaceRoutes := crc.convertToInterfaceRoutes(testRoutes)

	if len(interfaceRoutes) != len(testRoutes) {
		t.Errorf("convertToInterfaceRoutes() returned %d routes, expected %d",
			len(interfaceRoutes), len(testRoutes))
	}

	// Verify conversion
	for i, route := range interfaceRoutes {
		if route.Path != testRoutes[i].Path {
			t.Errorf("Route %d path mismatch: got %s, want %s",
				i, route.Path, testRoutes[i].Path)
		}
		if route.TemplateFile != testRoutes[i].TemplateFile {
			t.Errorf("Route %d template file mismatch: got %s, want %s",
				i, route.TemplateFile, testRoutes[i].TemplateFile)
		}
		if route.IsDynamic != testRoutes[i].IsDynamic {
			t.Errorf("Route %d dynamic flag mismatch: got %v, want %v",
				i, route.IsDynamic, testRoutes[i].IsDynamic)
		}
	}
}
