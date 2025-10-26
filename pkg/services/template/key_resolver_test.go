package template

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

// Mock ConfigService for template tests
type mockTemplateConfigService struct {
	layoutRootDir string
}

func (m *mockTemplateConfigService) GetLayoutRootDirectory() string { 
	if m.layoutRootDir != "" {
		return m.layoutRootDir
	}
	return "app" 
}

// Implement other required methods with defaults
func (m *mockTemplateConfigService) GetServerHost() string                     { return "localhost" }
func (m *mockTemplateConfigService) GetServerPort() int                        { return 8080 }
func (m *mockTemplateConfigService) GetServerBaseURL() string                  { return "http://localhost:8080" }
func (m *mockTemplateConfigService) GetServerReadTimeout() time.Duration       { return 30 * time.Second }
func (m *mockTemplateConfigService) GetServerWriteTimeout() time.Duration      { return 30 * time.Second }
func (m *mockTemplateConfigService) GetServerIdleTimeout() time.Duration       { return 60 * time.Second }
func (m *mockTemplateConfigService) GetServerShutdownTimeout() time.Duration   { return 10 * time.Second }
func (m *mockTemplateConfigService) GetSupportedLocales() []string             { return []string{"en"} }
func (m *mockTemplateConfigService) GetDefaultLocale() string                  { return "en" }
func (m *mockTemplateConfigService) GetFallbackLocale() string                 { return "en" }
func (m *mockTemplateConfigService) GetLayoutFileName() string                 { return "layout.templ" }
func (m *mockTemplateConfigService) GetLayoutAssetsDirectory() string          { return "assets" }
func (m *mockTemplateConfigService) GetLayoutAssetsRouteName() string          { return "/assets/" }
func (m *mockTemplateConfigService) GetTemplateExtension() string              { return ".templ" }
func (m *mockTemplateConfigService) GetMetadataExtension() string              { return ".yaml" }
func (m *mockTemplateConfigService) IsLayoutInheritanceEnabled() bool          { return true }
func (m *mockTemplateConfigService) GetTemplateOutputDir() string              { return "generated" }
func (m *mockTemplateConfigService) GetTemplatePackageName() string            { return "templates" }
func (m *mockTemplateConfigService) GetDatabaseHost() string                   { return "localhost" }
func (m *mockTemplateConfigService) GetDatabasePort() int                      { return 5432 }
func (m *mockTemplateConfigService) GetDatabaseUser() string                   { return "user" }
func (m *mockTemplateConfigService) GetDatabasePassword() string               { return "password" }
func (m *mockTemplateConfigService) GetDatabaseName() string                   { return "testdb" }
func (m *mockTemplateConfigService) GetDatabaseSSLMode() string                { return "disable" }
func (m *mockTemplateConfigService) IsEmailVerificationRequired() bool         { return false }
func (m *mockTemplateConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *mockTemplateConfigService) GetSessionCookieName() string              { return "session" }
func (m *mockTemplateConfigService) GetSessionExpiry() time.Duration           { return 24 * time.Hour }
func (m *mockTemplateConfigService) IsSessionSecure() bool                     { return false }
func (m *mockTemplateConfigService) IsSessionHttpOnly() bool                   { return true }
func (m *mockTemplateConfigService) GetSessionSameSite() string                { return "Lax" }
func (m *mockTemplateConfigService) GetMinPasswordLength() int                 { return 8 }
func (m *mockTemplateConfigService) IsStrongPasswordRequired() bool            { return false }
func (m *mockTemplateConfigService) ShouldCreateDefaultAdmin() bool            { return false }
func (m *mockTemplateConfigService) GetDefaultAdminEmail() string              { return "" }
func (m *mockTemplateConfigService) GetDefaultAdminPassword() string           { return "" }
func (m *mockTemplateConfigService) GetDefaultAdminFirstName() string          { return "" }
func (m *mockTemplateConfigService) GetDefaultAdminLastName() string           { return "" }
func (m *mockTemplateConfigService) GetCSRFSecret() string                     { return "secret" }
func (m *mockTemplateConfigService) IsCSRFSecure() bool                        { return false }
func (m *mockTemplateConfigService) IsCSRFHttpOnly() bool                      { return true }
func (m *mockTemplateConfigService) GetCSRFSameSite() string                   { return "Lax" }
func (m *mockTemplateConfigService) IsRateLimitEnabled() bool                  { return false }
func (m *mockTemplateConfigService) GetRateLimitRequests() int                 { return 100 }
func (m *mockTemplateConfigService) AreSecurityHeadersEnabled() bool           { return false }
func (m *mockTemplateConfigService) IsHSTSEnabled() bool                       { return false }
func (m *mockTemplateConfigService) GetHSTSMaxAge() int                        { return 31536000 }
func (m *mockTemplateConfigService) GetLogLevel() string                       { return "info" }
func (m *mockTemplateConfigService) GetLogFormat() string                      { return "json" }
func (m *mockTemplateConfigService) GetLogOutput() string                      { return "stdout" }
func (m *mockTemplateConfigService) IsFileLoggingEnabled() bool                { return false }
func (m *mockTemplateConfigService) GetLogFilePath() string                    { return "" }
func (m *mockTemplateConfigService) GetSMTPHost() string                       { return "" }
func (m *mockTemplateConfigService) GetSMTPPort() int                          { return 587 }
func (m *mockTemplateConfigService) GetSMTPUsername() string                   { return "" }
func (m *mockTemplateConfigService) GetSMTPPassword() string                   { return "" }
func (m *mockTemplateConfigService) IsSMTPTLSEnabled() bool                    { return true }
func (m *mockTemplateConfigService) GetFromEmail() string                      { return "" }
func (m *mockTemplateConfigService) GetFromName() string                       { return "" }
func (m *mockTemplateConfigService) GetReplyToEmail() string                   { return "" }
func (m *mockTemplateConfigService) IsEmailDummyModeEnabled() bool             { return true }
func (m *mockTemplateConfigService) IsDevelopment() bool                       { return true }
func (m *mockTemplateConfigService) IsProduction() bool                        { return false }

func TestKeyResolver_CreateTemplateKeyFromPath(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		functionName string
		rootDir      string
		expected     string
	}{
		{
			name:         "root level template",
			filePath:     "app/page.templ",
			functionName: "HomePage",
			rootDir:      "app",
			expected:     "HomePage",
		},
		{
			name:         "nested template",
			filePath:     "app/dashboard/page.templ",
			functionName: "DashboardPage",
			rootDir:      "app",
			expected:     "dashboard.DashboardPage",
		},
		{
			name:         "locale template",
			filePath:     "app/locale_/dashboard/page.templ",
			functionName: "DashboardPage",
			rootDir:      "app",
			expected:     "locale.dashboard.DashboardPage",
		},
		{
			name:         "deep nested template",
			filePath:     "app/admin/users/edit/page.templ",
			functionName: "EditUserPage",
			rootDir:      "app",
			expected:     "admin.users.edit.EditUserPage",
		},
		{
			name:         "template with underscore directory",
			filePath:     "app/user_/profile/page.templ",
			functionName: "ProfilePage",
			rootDir:      "app",
			expected:     "user.profile.ProfilePage",
		},
		{
			name:         "empty root directory",
			filePath:     "page.templ",
			functionName: "HomePage",
			rootDir:      "",
			expected:     "HomePage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockTemplateConfigService{layoutRootDir: tt.rootDir}
			logger := zap.NewNop()
			resolver := NewKeyResolver(logger, config, "test-module")

			result := resolver.CreateTemplateKeyFromPath(tt.filePath, tt.functionName)

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestKeyResolver_ResolveTemplateKey(t *testing.T) {
	tests := []struct {
		name         string
		templateKey  string
		expectedPath string
		expectedFunc string
	}{
		{
			name:         "simple key",
			templateKey:  "HomePage",
			expectedPath: "",
			expectedFunc: "HomePage",
		},
		{
			name:         "nested key",
			templateKey:  "dashboard.DashboardPage",
			expectedPath: "dashboard",
			expectedFunc: "DashboardPage",
		},
		{
			name:         "deep nested key",
			templateKey:  "admin.users.edit.EditUserPage",
			expectedPath: "admin/users/edit",
			expectedFunc: "EditUserPage",
		},
		{
			name:         "locale key",
			templateKey:  "locale.dashboard.DashboardPage",
			expectedPath: "locale/dashboard",
			expectedFunc: "DashboardPage",
		},
	}

	config := &mockTemplateConfigService{}
	logger := zap.NewNop()
	resolver := NewKeyResolver(logger, config, "test-module")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, functionName := resolver.ResolveTemplateKey(tt.templateKey)

			if path != tt.expectedPath {
				t.Errorf("Expected path %q, got %q", tt.expectedPath, path)
			}

			if functionName != tt.expectedFunc {
				t.Errorf("Expected function name %q, got %q", tt.expectedFunc, functionName)
			}
		})
	}
}

func TestKeyResolver_NormalizeTemplateKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already normalized",
			input:    "dashboard.page",
			expected: "dashboard.page",
		},
		{
			name:     "with uppercase",
			input:    "Dashboard.Page",
			expected: "dashboard.page",
		},
		{
			name:     "with slashes",
			input:    "dashboard/page",
			expected: "dashboard.page",
		},
		{
			name:     "with underscores",
			input:    "dashboard_page",
			expected: "dashboard.page",
		},
		{
			name:     "mixed separators",
			input:    "Admin/Users_Edit.Page",
			expected: "admin.users.edit.page",
		},
		{
			name:     "complex case",
			input:    "Locale_/Dashboard/User_Profile.EditPage",
			expected: "locale..dashboard.user.profile.editpage",
		},
	}

	config := &mockTemplateConfigService{}
	logger := zap.NewNop()
	resolver := NewKeyResolver(logger, config, "test-module")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.NormalizeTemplateKey(tt.input)

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNewKeyResolver(t *testing.T) {
	config := &mockTemplateConfigService{}
	logger := zap.NewNop()
	moduleName := "test-module"

	resolver := NewKeyResolver(logger, config, moduleName)

	if resolver == nil {
		t.Fatal("NewKeyResolver returned nil")
	}

	if resolver.logger != logger {
		t.Error("Logger not set correctly")
	}

	if resolver.config != config {
		t.Error("Config not set correctly")
	}

	if resolver.moduleName != moduleName {
		t.Error("Module name not set correctly")
	}
}
// Auth redirect routes (only for success cases)
func (m *mockTemplateConfigService) GetSignInSuccessRoute() string  { return "/dashboard" }
func (m *mockTemplateConfigService) GetSignUpSuccessRoute() string  { return "/welcome" }
func (m *mockTemplateConfigService) GetSignOutSuccessRoute() string { return "/" }

// Auth routes
func (m *mockTemplateConfigService) GetSignInRoute() string { return "/login" }

// Router configuration methods
func (m *mockTemplateConfigService) GetRouterEnableTrailingSlash() bool     { return true }
func (m *mockTemplateConfigService) GetRouterEnableSlashRedirect() bool     { return true }
func (m *mockTemplateConfigService) GetRouterEnableMethodNotAllowed() bool  { return true }
