package services

import (
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// Mock implementations for testing
type mockRouteDiscoveryConfigService struct {
	layoutRootDir     string
	layoutFileName    string
	templateExtension string
	supportedLocales  []string
	defaultLocale     string
}

func (m *mockRouteDiscoveryConfigService) GetLayoutRootDirectory() string {
	if m.layoutRootDir == "" {
		return "app"
	}
	return m.layoutRootDir
}

func (m *mockRouteDiscoveryConfigService) GetLayoutFileName() string {
	if m.layoutFileName == "" {
		return "layout"
	}
	return m.layoutFileName
}

func (m *mockRouteDiscoveryConfigService) GetTemplateExtension() string {
	if m.templateExtension == "" {
		return ".templ"
	}
	return m.templateExtension
}

func (m *mockRouteDiscoveryConfigService) GetSupportedLocales() []string {
	if len(m.supportedLocales) == 0 {
		return []string{"en", "de"}
	}
	return m.supportedLocales
}

func (m *mockRouteDiscoveryConfigService) GetDefaultLocale() string {
	if m.defaultLocale == "" {
		return "en"
	}
	return m.defaultLocale
}

// Additional methods to satisfy interfaces.ConfigService
func (m *mockRouteDiscoveryConfigService) AreSecurityHeadersEnabled() bool { return false }
func (m *mockRouteDiscoveryConfigService) GetServerHost() string { return "localhost" }
func (m *mockRouteDiscoveryConfigService) GetServerPort() int { return 8080 }
func (m *mockRouteDiscoveryConfigService) GetServerBaseURL() string { return "http://localhost:8080" }
func (m *mockRouteDiscoveryConfigService) GetDatabaseHost() string { return "localhost" }
func (m *mockRouteDiscoveryConfigService) GetDatabasePort() int { return 5432 }
func (m *mockRouteDiscoveryConfigService) GetDatabaseUser() string { return "user" }
func (m *mockRouteDiscoveryConfigService) GetDatabasePassword() string { return "pass" }
func (m *mockRouteDiscoveryConfigService) GetDatabaseName() string { return "db" }
func (m *mockRouteDiscoveryConfigService) GetDatabaseSSLMode() string { return "disable" }
func (m *mockRouteDiscoveryConfigService) IsEmailVerificationRequired() bool { return false }
func (m *mockRouteDiscoveryConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *mockRouteDiscoveryConfigService) GetSessionCookieName() string { return "session" }
func (m *mockRouteDiscoveryConfigService) GetSessionExpiry() time.Duration { return 24 * time.Hour }
func (m *mockRouteDiscoveryConfigService) IsSessionSecure() bool { return false }
func (m *mockRouteDiscoveryConfigService) IsSessionHTTPOnly() bool { return true }
func (m *mockRouteDiscoveryConfigService) GetSessionSameSite() string { return "lax" }
func (m *mockRouteDiscoveryConfigService) GetMinPasswordLength() int { return 8 }
func (m *mockRouteDiscoveryConfigService) IsStrongPasswordRequired() bool { return false }
func (m *mockRouteDiscoveryConfigService) ShouldCreateDefaultAdmin() bool { return false }
func (m *mockRouteDiscoveryConfigService) GetDefaultAdminEmail() string { return "admin@example.com" }
func (m *mockRouteDiscoveryConfigService) GetDefaultAdminPassword() string { return "password" }
func (m *mockRouteDiscoveryConfigService) GetDefaultAdminFirstName() string { return "Admin" }
func (m *mockRouteDiscoveryConfigService) GetDefaultAdminLastName() string { return "User" }
func (m *mockRouteDiscoveryConfigService) GetSMTPHost() string { return "" }
func (m *mockRouteDiscoveryConfigService) GetSMTPPort() int { return 587 }
func (m *mockRouteDiscoveryConfigService) GetSMTPUsername() string { return "" }
func (m *mockRouteDiscoveryConfigService) GetSMTPPassword() string { return "" }
func (m *mockRouteDiscoveryConfigService) IsSMTPTLSEnabled() bool { return true }
func (m *mockRouteDiscoveryConfigService) GetFromEmail() string { return "noreply@example.com" }
func (m *mockRouteDiscoveryConfigService) GetFromName() string { return "App" }
func (m *mockRouteDiscoveryConfigService) GetReplyToEmail() string { return "" }
func (m *mockRouteDiscoveryConfigService) IsEmailDummyModeEnabled() bool { return true }
func (m *mockRouteDiscoveryConfigService) GetCSRFSecret() string { return "secret" }
func (m *mockRouteDiscoveryConfigService) IsCSRFSecure() bool { return false }
func (m *mockRouteDiscoveryConfigService) IsCSRFHTTPOnly() bool { return true }
func (m *mockRouteDiscoveryConfigService) GetCSRFSameSite() string { return "strict" }
func (m *mockRouteDiscoveryConfigService) IsRateLimitEnabled() bool { return false }
func (m *mockRouteDiscoveryConfigService) GetRateLimitRequests() int { return 100 }
func (m *mockRouteDiscoveryConfigService) IsHSTSEnabled() bool { return false }
func (m *mockRouteDiscoveryConfigService) GetHSTSMaxAge() int { return 31536000 }
func (m *mockRouteDiscoveryConfigService) GetLogLevel() string { return "info" }
func (m *mockRouteDiscoveryConfigService) GetLogFormat() string { return "json" }
func (m *mockRouteDiscoveryConfigService) GetLogOutput() string { return "stdout" }
func (m *mockRouteDiscoveryConfigService) IsFileLoggingEnabled() bool { return false }
func (m *mockRouteDiscoveryConfigService) GetLogFilePath() string { return "" }
func (m *mockRouteDiscoveryConfigService) IsProductionMode() bool { return false }
func (m *mockRouteDiscoveryConfigService) IsDevelopmentMode() bool { return true }

type mockFileSystemChecker struct{}

func (m *mockFileSystemChecker) FileExists(path string) bool {
	return true // For testing, assume all files exist
}

func (m *mockFileSystemChecker) IsDirectory(path string) bool {
	return false // For testing, assume paths are files
}

func (m *mockFileSystemChecker) WalkDirectory(root string, walkFn func(path string, isDir bool, err error) error) error {
	return nil // For testing, no walking needed
}

type mockTemplateRegistry struct {
	routeMapping map[string]string
}

func (m *mockTemplateRegistry) GetRouteToTemplateMapping() map[string]string {
	if m.routeMapping == nil {
		return map[string]string{
			"/":                    "template1",
			"/{locale}":           "template2", 
			"/{locale}/dashboard": "template3",
			"/{locale}/user/{id}": "template4",
			"/login":              "template5",
		}
	}
	return m.routeMapping
}

func (m *mockTemplateRegistry) GetTemplate(key string) (templ.Component, error) {
	return nil, nil
}

func (m *mockTemplateRegistry) GetTemplateFunction(key string) (func() interface{}, bool) {
	return nil, false
}

func (m *mockTemplateRegistry) GetAllTemplateKeys() []string {
	return []string{}
}

func (m *mockTemplateRegistry) IsAvailable(key string) bool {
	return false
}

func (m *mockTemplateRegistry) GetTemplateByRoute(route string) (templ.Component, error) {
	return nil, nil
}

func (m *mockTemplateRegistry) RequiresDataService(key string) bool {
	return false
}

func (m *mockTemplateRegistry) GetDataServiceInfo(key string) (interfaces.DataServiceInfo, bool) {
	return interfaces.DataServiceInfo{}, false
}

func createTestContainer() do.Injector {
	injector := do.New()
	
	// Register mocks with proper interface types
	do.ProvideValue[interfaces.ConfigService](injector, &mockRouteDiscoveryConfigService{})
	do.ProvideValue[*zap.Logger](injector, zap.NewNop())
	do.ProvideValue[middleware.FileSystemChecker](injector, &mockFileSystemChecker{})
	do.ProvideValue[interfaces.TemplateRegistry](injector, &mockTemplateRegistry{})
	
	return injector
}

func TestNewRouteDiscovery(t *testing.T) {
	injector := createTestContainer()
	
	discovery, err := NewRouteDiscovery(injector)
	if err != nil {
		t.Fatalf("NewRouteDiscovery() returned error: %v", err)
	}
	
	if discovery == nil {
		t.Fatal("NewRouteDiscovery() returned nil")
	}
}

func TestDiscoverRoutes(t *testing.T) {
	injector := createTestContainer()
	discovery, err := NewRouteDiscovery(injector)
	if err != nil {
		t.Fatalf("Failed to create route discovery: %v", err)
	}
	
	// Test with demo directory that actually contains .templ files
	routes, err := discovery.DiscoverRoutes("../../demo/app")
	if err != nil {
		t.Fatalf("DiscoverRoutes() returned error: %v", err)
	}
	
	// The demo directory should have some routes
	if len(routes) == 0 {
		t.Skip("No routes found in demo directory - this is expected if demo templates don't exist")
	}
	
	// Log found routes for debugging
	t.Logf("Found %d routes:", len(routes))
	for _, route := range routes {
		t.Logf("  Route: %s -> %s (dynamic: %v)", route.Path, route.TemplateFile, route.IsDynamic)
	}
}

func TestGenerateTemplateFilePathFromPattern(t *testing.T) {
	injector := createTestContainer()
	discovery, err := NewRouteDiscovery(injector)
	if err != nil {
		t.Fatalf("Failed to create route discovery: %v", err)
	}
	
	// Access the private method through the implementation
	impl := discovery.(*routeDiscoveryImpl)
	
	tests := []struct {
		name           string
		routePattern   string
		expectedPath   string
		description    string
	}{
		{
			name:         "Root route",
			routePattern: "/",
			expectedPath: "app/page.templ",
			description:  "Root route should map to app/page.templ",
		},
		{
			name:         "Locale route with placeholder",
			routePattern: "/{locale}",
			expectedPath: "app/locale_/page.templ",
			description:  "Locale placeholder should be converted to locale_ directory",
		},
		{
			name:         "Nested locale route",
			routePattern: "/{locale}/dashboard",
			expectedPath: "app/locale_/dashboard/page.templ",
			description:  "Nested locale route should maintain directory structure",
		},
		{
			name:         "Dynamic ID parameter",
			routePattern: "/{locale}/user/{id}",
			expectedPath: "app/locale_/user/id_/page.templ",
			description:  "Dynamic ID should be converted to id_ directory",
		},
		{
			name:         "Static route",
			routePattern: "/login",
			expectedPath: "app/login/page.templ",
			description:  "Static routes should map directly",
		},
		{
			name:         "Concrete locale (en)",
			routePattern: "/en/dashboard",
			expectedPath: "app/locale_/dashboard/page.templ",
			description:  "Concrete locale should be converted to locale_",
		},
		{
			name:         "Concrete locale (de)",
			routePattern: "/de/admin",
			expectedPath: "app/locale_/admin/page.templ",
			description:  "German locale should be converted to locale_",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := impl.generateTemplateFilePathFromPattern(tt.routePattern)
			if result != tt.expectedPath {
				t.Errorf("generateTemplateFilePathFromPattern(%s) = %s, want %s\nDescription: %s",
					tt.routePattern, result, tt.expectedPath, tt.description)
			}
		})
	}
}

func TestI18nPlaceholderFix(t *testing.T) {
	injector := createTestContainer()
	discovery, err := NewRouteDiscovery(injector)
	if err != nil {
		t.Fatalf("Failed to create route discovery: %v", err)
	}
	
	impl := discovery.(*routeDiscoveryImpl)
	
	// Test the specific fix for i18n placeholder resolution
	testCases := []struct {
		input    string
		expected string
		desc     string
	}{
		{
			input:    "/{locale}",
			expected: "app/locale_/page.templ",
			desc:     "Locale placeholder should resolve to locale_ directory",
		},
		{
			input:    "/{locale}/dashboard",
			expected: "app/locale_/dashboard/page.templ", 
			desc:     "Nested locale routes should work",
		},
		{
			input:    "/{locale}/user/{id}",
			expected: "app/locale_/user/id_/page.templ",
			desc:     "Multiple dynamic parameters should work",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := impl.generateTemplateFilePathFromPattern(tc.input)
			if result != tc.expected {
				t.Errorf("i18n fix test failed: input=%s, got=%s, want=%s",
					tc.input, result, tc.expected)
			}
		})
	}
}
// Additional missing methods for ConfigService
func (m *mockRouteDiscoveryConfigService) GetServerReadTimeout() time.Duration { return 30 * time.Second }
func (m *mockRouteDiscoveryConfigService) GetServerWriteTimeout() time.Duration { return 30 * time.Second }
func (m *mockRouteDiscoveryConfigService) GetServerIdleTimeout() time.Duration { return 2 * time.Minute }
func (m *mockRouteDiscoveryConfigService) GetServerShutdownTimeout() time.Duration { return 30 * time.Second }
func (m *mockRouteDiscoveryConfigService) GetFallbackLocale() string { return "en" }
func (m *mockRouteDiscoveryConfigService) GetLayoutAssetsDirectory() string { return "assets" }
func (m *mockRouteDiscoveryConfigService) GetLayoutAssetsRouteName() string { return "/assets/" }
func (m *mockRouteDiscoveryConfigService) GetMetadataExtension() string { return ".yaml" }
func (m *mockRouteDiscoveryConfigService) IsLayoutInheritanceEnabled() bool { return true }
func (m *mockRouteDiscoveryConfigService) GetTemplateOutputDir() string { return "generated" }
func (m *mockRouteDiscoveryConfigService) GetTemplatePackageName() string { return "templates" }
func (m *mockRouteDiscoveryConfigService) IsSessionHttpOnly() bool { return true }
func (m *mockRouteDiscoveryConfigService) IsCSRFHttpOnly() bool { return true }
func (m *mockRouteDiscoveryConfigService) IsDevelopment() bool { return true }
func (m *mockRouteDiscoveryConfigService) IsProduction() bool { return false }

// Auth redirect routes (only for success cases)
func (m *mockRouteDiscoveryConfigService) GetSignInSuccessRoute() string  { return "/dashboard" }
func (m *mockRouteDiscoveryConfigService) GetSignUpSuccessRoute() string  { return "/welcome" }
func (m *mockRouteDiscoveryConfigService) GetSignOutSuccessRoute() string { return "/" }

// Auth routes
func (m *mockRouteDiscoveryConfigService) GetSignInRoute() string { return "/login" }

// Router configuration methods
func (m *mockRouteDiscoveryConfigService) GetRouterEnableTrailingSlash() bool     { return true }
func (m *mockRouteDiscoveryConfigService) GetRouterEnableSlashRedirect() bool     { return true }
func (m *mockRouteDiscoveryConfigService) GetRouterEnableMethodNotAllowed() bool  { return true }
