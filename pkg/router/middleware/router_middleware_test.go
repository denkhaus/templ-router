package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock config service for router middleware tests
type mockRouterConfigService struct {
	enableTrailingSlash    bool
	enableSlashRedirect    bool
	enableMethodNotAllowed bool
}

func (m *mockRouterConfigService) GetRouterEnableTrailingSlash() bool     { return m.enableTrailingSlash }
func (m *mockRouterConfigService) GetRouterEnableSlashRedirect() bool     { return m.enableSlashRedirect }
func (m *mockRouterConfigService) GetRouterEnableMethodNotAllowed() bool  { return m.enableMethodNotAllowed }

// Implement all required ConfigService methods (minimal implementation for tests)
func (m *mockRouterConfigService) GetLayoutRootDirectory() string            { return "app" }
func (m *mockRouterConfigService) GetSupportedLocales() []string             { return []string{"en", "de"} }
func (m *mockRouterConfigService) GetDefaultLocale() string                  { return "en" }
func (m *mockRouterConfigService) GetFallbackLocale() string                 { return "en" }
func (m *mockRouterConfigService) GetLayoutFileName() string                 { return "layout" }
func (m *mockRouterConfigService) GetTemplateExtension() string              { return ".templ" }
func (m *mockRouterConfigService) GetMetadataExtension() string              { return ".yaml" }
func (m *mockRouterConfigService) IsLayoutInheritanceEnabled() bool          { return true }
func (m *mockRouterConfigService) GetTemplateOutputDir() string              { return "generated" }
func (m *mockRouterConfigService) GetTemplatePackageName() string            { return "templates" }
func (m *mockRouterConfigService) GetLayoutAssetsDirectory() string          { return "assets" }
func (m *mockRouterConfigService) GetLayoutAssetsRouteName() string          { return "/assets/" }
func (m *mockRouterConfigService) IsDevelopment() bool                       { return true }
func (m *mockRouterConfigService) IsProduction() bool                        { return false }
func (m *mockRouterConfigService) GetServerHost() string                     { return "localhost" }
func (m *mockRouterConfigService) GetServerPort() int                        { return 8080 }
func (m *mockRouterConfigService) GetServerBaseURL() string                  { return "http://localhost:8080" }
func (m *mockRouterConfigService) GetServerReadTimeout() time.Duration       { return 30 * time.Second }
func (m *mockRouterConfigService) GetServerWriteTimeout() time.Duration      { return 30 * time.Second }
func (m *mockRouterConfigService) GetServerIdleTimeout() time.Duration       { return 60 * time.Second }
func (m *mockRouterConfigService) GetServerShutdownTimeout() time.Duration   { return 10 * time.Second }
func (m *mockRouterConfigService) GetDatabaseHost() string                   { return "localhost" }
func (m *mockRouterConfigService) GetDatabasePort() int                      { return 5432 }
func (m *mockRouterConfigService) GetDatabaseUser() string                   { return "user" }
func (m *mockRouterConfigService) GetDatabasePassword() string               { return "password" }
func (m *mockRouterConfigService) GetDatabaseName() string                   { return "testdb" }
func (m *mockRouterConfigService) GetDatabaseSSLMode() string                { return "disable" }
func (m *mockRouterConfigService) IsEmailVerificationRequired() bool         { return false }
func (m *mockRouterConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *mockRouterConfigService) GetSessionCookieName() string              { return "session" }
func (m *mockRouterConfigService) GetSessionExpiry() time.Duration           { return 24 * time.Hour }
func (m *mockRouterConfigService) IsSessionSecure() bool                     { return false }
func (m *mockRouterConfigService) IsSessionHttpOnly() bool                   { return true }
func (m *mockRouterConfigService) GetSessionSameSite() string                { return "Lax" }
func (m *mockRouterConfigService) GetMinPasswordLength() int                 { return 8 }
func (m *mockRouterConfigService) IsStrongPasswordRequired() bool            { return false }
func (m *mockRouterConfigService) ShouldCreateDefaultAdmin() bool            { return false }
func (m *mockRouterConfigService) GetDefaultAdminEmail() string              { return "" }
func (m *mockRouterConfigService) GetDefaultAdminPassword() string           { return "" }
func (m *mockRouterConfigService) GetDefaultAdminFirstName() string          { return "" }
func (m *mockRouterConfigService) GetDefaultAdminLastName() string           { return "" }
func (m *mockRouterConfigService) GetSignInRoute() string                    { return "/login" }
func (m *mockRouterConfigService) GetSignInSuccessRoute() string             { return "/dashboard" }
func (m *mockRouterConfigService) GetSignUpSuccessRoute() string             { return "/welcome" }
func (m *mockRouterConfigService) GetSignOutSuccessRoute() string            { return "/" }
func (m *mockRouterConfigService) GetCSRFSecret() string                     { return "secret" }
func (m *mockRouterConfigService) IsCSRFSecure() bool                        { return false }
func (m *mockRouterConfigService) IsCSRFHttpOnly() bool                      { return true }
func (m *mockRouterConfigService) GetCSRFSameSite() string                   { return "Lax" }
func (m *mockRouterConfigService) IsRateLimitEnabled() bool                  { return false }
func (m *mockRouterConfigService) GetRateLimitRequests() int                 { return 100 }
func (m *mockRouterConfigService) AreSecurityHeadersEnabled() bool           { return false }
func (m *mockRouterConfigService) IsHSTSEnabled() bool                       { return false }
func (m *mockRouterConfigService) GetHSTSMaxAge() int                        { return 31536000 }
func (m *mockRouterConfigService) GetLogLevel() string                       { return "info" }
func (m *mockRouterConfigService) GetLogFormat() string                      { return "json" }
func (m *mockRouterConfigService) GetLogOutput() string                      { return "stdout" }
func (m *mockRouterConfigService) IsFileLoggingEnabled() bool                { return false }
func (m *mockRouterConfigService) GetLogFilePath() string                    { return "" }
func (m *mockRouterConfigService) GetSMTPHost() string                       { return "" }
func (m *mockRouterConfigService) GetSMTPPort() int                          { return 587 }
func (m *mockRouterConfigService) GetSMTPUsername() string                   { return "" }
func (m *mockRouterConfigService) GetSMTPPassword() string                   { return "" }
func (m *mockRouterConfigService) IsSMTPTLSEnabled() bool                    { return true }
func (m *mockRouterConfigService) GetFromEmail() string                      { return "" }
func (m *mockRouterConfigService) GetFromName() string                       { return "" }
func (m *mockRouterConfigService) GetReplyToEmail() string                   { return "" }
func (m *mockRouterConfigService) IsEmailDummyModeEnabled() bool             { return true }

func TestNewRouterMiddleware(t *testing.T) {
	// Setup DI container
	injector := do.New()
	defer injector.Shutdown()

	// Register dependencies
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockRouterConfigService{
			enableTrailingSlash:    true,
			enableSlashRedirect:    true,
			enableMethodNotAllowed: true,
		}, nil
	})

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	// Test middleware creation
	middleware, err := NewRouterMiddleware(injector)
	require.NoError(t, err)
	assert.NotNil(t, middleware)
}

func TestRouterMiddleware_ConfigureRouterMiddleware_TrailingSlashEnabled(t *testing.T) {
	// Setup DI container
	injector := do.New()
	defer injector.Shutdown()

	// Register dependencies with trailing slash enabled
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockRouterConfigService{
			enableTrailingSlash:    true,
			enableSlashRedirect:    false,
			enableMethodNotAllowed: false,
		}, nil
	})

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	// Create middleware
	middleware, err := NewRouterMiddleware(injector)
	require.NoError(t, err)

	// Create Chi router
	router := chi.NewRouter()

	// Configure router middleware BEFORE adding routes
	err = middleware.ConfigureRouterMiddleware(router)
	require.NoError(t, err)

	// Add a test route AFTER middleware configuration
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Test trailing slash redirection
	// Request to /test/ should redirect to /test
	req := httptest.NewRequest("GET", "/test/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should redirect (301 or 302)
	assert.True(t, w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound,
		"Expected redirect status, got %d", w.Code)

	// Check location header
	location := w.Header().Get("Location")
	assert.Equal(t, "/test", location, "Expected redirect to /test, got %s", location)
}

func TestRouterMiddleware_ConfigureRouterMiddleware_SlashRedirectEnabled(t *testing.T) {
	// Setup DI container
	injector := do.New()
	defer injector.Shutdown()

	// Register dependencies with slash redirect enabled
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockRouterConfigService{
			enableTrailingSlash:    false,
			enableSlashRedirect:    true,
			enableMethodNotAllowed: false,
		}, nil
	})

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	// Create middleware
	middleware, err := NewRouterMiddleware(injector)
	require.NoError(t, err)

	// Create Chi router
	router := chi.NewRouter()

	// Configure router middleware BEFORE adding routes
	err = middleware.ConfigureRouterMiddleware(router)
	require.NoError(t, err)

	// Add a test route AFTER middleware configuration
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})


	// Test double slash cleanup with a path that has double slashes
	// Add a route that can be cleaned
	router.Get("/test/path", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("cleaned path response"))
	})

	// Request to /test//path should be cleaned to /test/path
	req := httptest.NewRequest("GET", "/test//path", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// CleanPath middleware should clean the path and find the route
	// If CleanPath is working, it should either redirect or serve the cleaned route
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound,
		"Expected success or redirect status, got %d", w.Code)
}

func TestRouterMiddleware_ConfigureRouterMiddleware_BothDisabled(t *testing.T) {
	// Setup DI container
	injector := do.New()
	defer injector.Shutdown()

	// Register dependencies with both features disabled
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockRouterConfigService{
			enableTrailingSlash:    false,
			enableSlashRedirect:    false,
			enableMethodNotAllowed: false,
		}, nil
	})

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	// Create middleware
	middleware, err := NewRouterMiddleware(injector)
	require.NoError(t, err)

	// Create Chi router
	router := chi.NewRouter()

	// Configure router middleware BEFORE adding routes
	err = middleware.ConfigureRouterMiddleware(router)
	require.NoError(t, err)

	// Add a test route AFTER middleware configuration
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})


	// Test that no redirection happens
	req := httptest.NewRequest("GET", "/test/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 404 since /test/ is not registered and no redirection middleware is active
	assert.Equal(t, http.StatusNotFound, w.Code,
		"Expected 404 when middleware is disabled, got %d", w.Code)
}

func TestRouterMiddleware_ConfigureRouterMiddleware_BothEnabled(t *testing.T) {
	// Setup DI container
	injector := do.New()
	defer injector.Shutdown()

	// Register dependencies with both features enabled
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockRouterConfigService{
			enableTrailingSlash:    true,
			enableSlashRedirect:    true,
			enableMethodNotAllowed: true,
		}, nil
	})

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	// Create middleware
	middleware, err := NewRouterMiddleware(injector)
	require.NoError(t, err)

	// Create Chi router
	router := chi.NewRouter()

	// Configure router middleware BEFORE adding routes
	err = middleware.ConfigureRouterMiddleware(router)
	require.NoError(t, err)

	// Add a test route AFTER middleware configuration
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})


	// Test that both middleware work together
	req := httptest.NewRequest("GET", "/test//", nil) // Double slash + trailing slash
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should redirect to cleaned path
	assert.True(t, w.Code == http.StatusMovedPermanently || w.Code == http.StatusFound,
		"Expected redirect status when both middleware are enabled, got %d", w.Code)
}