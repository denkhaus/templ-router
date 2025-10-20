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
type mockConfigService struct {
	layoutRootDir     string
	layoutFileName    string
	templateExtension string
	supportedLocales  []string
	defaultLocale     string
}

func (m *mockConfigService) GetLayoutRootDirectory() string {
	if m.layoutRootDir == "" {
		return "app"
	}
	return m.layoutRootDir
}

func (m *mockConfigService) GetLayoutFileName() string {
	if m.layoutFileName == "" {
		return "layout"
	}
	return m.layoutFileName
}

func (m *mockConfigService) GetTemplateExtension() string {
	if m.templateExtension == "" {
		return ".templ"
	}
	return m.templateExtension
}

func (m *mockConfigService) GetSupportedLocales() []string {
	if len(m.supportedLocales) == 0 {
		return []string{"en", "de"}
	}
	return m.supportedLocales
}

func (m *mockConfigService) GetDefaultLocale() string {
	if m.defaultLocale == "" {
		return "en"
	}
	return m.defaultLocale
}

// Additional methods to satisfy interfaces.ConfigService
func (m *mockConfigService) AreSecurityHeadersEnabled() bool { return false }
func (m *mockConfigService) GetServerHost() string { return "localhost" }
func (m *mockConfigService) GetServerPort() int { return 8080 }
func (m *mockConfigService) GetServerBaseURL() string { return "http://localhost:8080" }
func (m *mockConfigService) GetDatabaseHost() string { return "localhost" }
func (m *mockConfigService) GetDatabasePort() int { return 5432 }
func (m *mockConfigService) GetDatabaseUser() string { return "user" }
func (m *mockConfigService) GetDatabasePassword() string { return "pass" }
func (m *mockConfigService) GetDatabaseName() string { return "db" }
func (m *mockConfigService) GetDatabaseSSLMode() string { return "disable" }
func (m *mockConfigService) IsEmailVerificationRequired() bool { return false }
func (m *mockConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *mockConfigService) GetSessionCookieName() string { return "session" }
func (m *mockConfigService) GetSessionExpiry() time.Duration { return 24 * time.Hour }
func (m *mockConfigService) IsSessionSecure() bool { return false }
func (m *mockConfigService) IsSessionHTTPOnly() bool { return true }
func (m *mockConfigService) GetSessionSameSite() string { return "lax" }
func (m *mockConfigService) GetMinPasswordLength() int { return 8 }
func (m *mockConfigService) IsStrongPasswordRequired() bool { return false }
func (m *mockConfigService) ShouldCreateDefaultAdmin() bool { return false }
func (m *mockConfigService) GetDefaultAdminEmail() string { return "admin@example.com" }
func (m *mockConfigService) GetDefaultAdminPassword() string { return "password" }
func (m *mockConfigService) GetDefaultAdminFirstName() string { return "Admin" }
func (m *mockConfigService) GetDefaultAdminLastName() string { return "User" }
func (m *mockConfigService) GetSMTPHost() string { return "" }
func (m *mockConfigService) GetSMTPPort() int { return 587 }
func (m *mockConfigService) GetSMTPUsername() string { return "" }
func (m *mockConfigService) GetSMTPPassword() string { return "" }
func (m *mockConfigService) IsSMTPTLSEnabled() bool { return true }
func (m *mockConfigService) GetFromEmail() string { return "noreply@example.com" }
func (m *mockConfigService) GetFromName() string { return "App" }
func (m *mockConfigService) GetReplyToEmail() string { return "" }
func (m *mockConfigService) IsEmailDummyModeEnabled() bool { return true }
func (m *mockConfigService) GetCSRFSecret() string { return "secret" }
func (m *mockConfigService) IsCSRFSecure() bool { return false }
func (m *mockConfigService) IsCSRFHTTPOnly() bool { return true }
func (m *mockConfigService) GetCSRFSameSite() string { return "strict" }
func (m *mockConfigService) IsRateLimitEnabled() bool { return false }
func (m *mockConfigService) GetRateLimitRequests() int { return 100 }
func (m *mockConfigService) IsHSTSEnabled() bool { return false }
func (m *mockConfigService) GetHSTSMaxAge() int { return 31536000 }
func (m *mockConfigService) GetLogLevel() string { return "info" }
func (m *mockConfigService) GetLogFormat() string { return "json" }
func (m *mockConfigService) GetLogOutput() string { return "stdout" }
func (m *mockConfigService) IsFileLoggingEnabled() bool { return false }
func (m *mockConfigService) GetLogFilePath() string { return "" }
func (m *mockConfigService) IsProductionMode() bool { return false }
func (m *mockConfigService) IsDevelopmentMode() bool { return true }

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

func createTestContainer() do.Injector {
	injector := do.New()
	
	// Register mocks with proper interface types
	do.ProvideValue[interfaces.ConfigService](injector, &mockConfigService{})
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
func (m *mockConfigService) GetServerReadTimeout() time.Duration { return 30 * time.Second }
func (m *mockConfigService) GetServerWriteTimeout() time.Duration { return 30 * time.Second }
func (m *mockConfigService) GetServerIdleTimeout() time.Duration { return 2 * time.Minute }
func (m *mockConfigService) GetServerShutdownTimeout() time.Duration { return 30 * time.Second }
func (m *mockConfigService) GetFallbackLocale() string { return "en" }
func (m *mockConfigService) GetLayoutAssetsDirectory() string { return "assets" }
func (m *mockConfigService) GetLayoutAssetsRouteName() string { return "/assets/" }
func (m *mockConfigService) GetMetadataExtension() string { return ".yaml" }
func (m *mockConfigService) IsLayoutInheritanceEnabled() bool { return true }
func (m *mockConfigService) GetTemplateOutputDir() string { return "generated" }
func (m *mockConfigService) GetTemplatePackageName() string { return "templates" }
func (m *mockConfigService) IsSessionHttpOnly() bool { return true }
func (m *mockConfigService) IsCSRFHttpOnly() bool { return true }
func (m *mockConfigService) IsDevelopment() bool { return true }
func (m *mockConfigService) IsProduction() bool { return false }

// Auth redirect routes (only for success cases)
func (m *mockConfigService) GetSignInSuccessRoute() string  { return "/dashboard" }
func (m *mockConfigService) GetSignUpSuccessRoute() string  { return "/welcome" }
func (m *mockConfigService) GetSignOutSuccessRoute() string { return "/" }
