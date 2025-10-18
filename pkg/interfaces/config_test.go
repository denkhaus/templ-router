package interfaces

import (
	"testing"
	"time"
)

// MockConfigService implements ConfigService for testing
type MockConfigService struct {
	config map[string]interface{}
}

func NewMockConfigService() *MockConfigService {
	return &MockConfigService{
		config: map[string]interface{}{
			"server.host":                    "localhost",
			"server.port":                    8080,
			"server.base_url":                "http://localhost:8080",
			"server.read_timeout":            30 * time.Second,
			"server.write_timeout":           30 * time.Second,
			"server.idle_timeout":            60 * time.Second,
			"server.shutdown_timeout":        10 * time.Second,
			"i18n.supported_locales":         []string{"en", "de", "fr"},
			"i18n.default_locale":            "en",
			"i18n.fallback_locale":           "en",
			"layout.root_directory":          "app",
			"layout.file_name":               "layout.templ",
			"layout.assets_directory":        "assets",
			"layout.assets_route_name":       "/assets/",
			"layout.template_extension":      ".templ",
			"layout.metadata_extension":      ".yaml",
			"layout.inheritance_enabled":     true,
			"template.output_dir":            "generated/templates",
			"template.package_name":          "templates",
			"database.host":                  "localhost",
			"database.port":                  5432,
			"database.user":                  "testuser",
			"database.password":              "testpass",
			"database.name":                  "testdb",
			"database.ssl_mode":              "disable",
			"auth.email_verification_required": true,
			"auth.verification_token_expiry":   24 * time.Hour,
			"auth.session_cookie_name":        "session_id",
			"auth.session_expiry":             7 * 24 * time.Hour,
			"auth.session_secure":             true,
			"auth.session_http_only":          true,
			"auth.session_same_site":          "Strict",
			"auth.min_password_length":        8,
			"auth.strong_password_required":   true,
			"auth.create_default_admin":       false,
			"auth.default_admin_email":        "admin@example.com",
			"auth.default_admin_password":     "admin123",
			"auth.default_admin_first_name":   "Admin",
			"auth.default_admin_last_name":    "User",
			"security.csrf_secret":            "csrf-secret-key",
			"security.csrf_secure":            true,
			"security.csrf_http_only":         true,
			"security.csrf_same_site":         "Strict",
			"security.rate_limit_enabled":     true,
			"security.rate_limit_requests":    100,
			"security.headers_enabled":        true,
			"security.hsts_enabled":           true,
			"security.hsts_max_age":           31536000,
			"logging.level":                   "info",
			"logging.format":                  "json",
			"logging.output":                  "stdout",
			"logging.file_enabled":            false,
			"logging.file_path":               "/var/log/app.log",
			"email.smtp_host":                 "smtp.example.com",
			"email.smtp_port":                 587,
			"email.smtp_username":             "user@example.com",
			"email.smtp_password":             "password",
			"email.smtp_tls_enabled":          true,
			"email.from_email":                "noreply@example.com",
			"email.from_name":                 "Test App",
			"email.reply_to_email":            "support@example.com",
			"email.dummy_mode_enabled":        true,
			"environment.development":         true,
			"environment.production":          false,
		},
	}
}

func (m *MockConfigService) GetServerHost() string {
	return m.config["server.host"].(string)
}

func (m *MockConfigService) GetServerPort() int {
	return m.config["server.port"].(int)
}

func (m *MockConfigService) GetServerBaseURL() string {
	return m.config["server.base_url"].(string)
}

func (m *MockConfigService) GetServerReadTimeout() time.Duration {
	return m.config["server.read_timeout"].(time.Duration)
}

func (m *MockConfigService) GetServerWriteTimeout() time.Duration {
	return m.config["server.write_timeout"].(time.Duration)
}

func (m *MockConfigService) GetServerIdleTimeout() time.Duration {
	return m.config["server.idle_timeout"].(time.Duration)
}

func (m *MockConfigService) GetServerShutdownTimeout() time.Duration {
	return m.config["server.shutdown_timeout"].(time.Duration)
}

func (m *MockConfigService) GetSupportedLocales() []string {
	return m.config["i18n.supported_locales"].([]string)
}

func (m *MockConfigService) GetDefaultLocale() string {
	return m.config["i18n.default_locale"].(string)
}

func (m *MockConfigService) GetFallbackLocale() string {
	return m.config["i18n.fallback_locale"].(string)
}

func (m *MockConfigService) GetLayoutRootDirectory() string {
	return m.config["layout.root_directory"].(string)
}

func (m *MockConfigService) GetLayoutFileName() string {
	return m.config["layout.file_name"].(string)
}

func (m *MockConfigService) GetLayoutAssetsDirectory() string {
	return m.config["layout.assets_directory"].(string)
}

func (m *MockConfigService) GetLayoutAssetsRouteName() string {
	return m.config["layout.assets_route_name"].(string)
}

func (m *MockConfigService) GetTemplateExtension() string {
	return m.config["layout.template_extension"].(string)
}

func (m *MockConfigService) GetMetadataExtension() string {
	return m.config["layout.metadata_extension"].(string)
}

func (m *MockConfigService) IsLayoutInheritanceEnabled() bool {
	return m.config["layout.inheritance_enabled"].(bool)
}

func (m *MockConfigService) GetTemplateOutputDir() string {
	return m.config["template.output_dir"].(string)
}

func (m *MockConfigService) GetTemplatePackageName() string {
	return m.config["template.package_name"].(string)
}

func (m *MockConfigService) GetDatabaseHost() string {
	return m.config["database.host"].(string)
}

func (m *MockConfigService) GetDatabasePort() int {
	return m.config["database.port"].(int)
}

func (m *MockConfigService) GetDatabaseUser() string {
	return m.config["database.user"].(string)
}

func (m *MockConfigService) GetDatabasePassword() string {
	return m.config["database.password"].(string)
}

func (m *MockConfigService) GetDatabaseName() string {
	return m.config["database.name"].(string)
}

func (m *MockConfigService) GetDatabaseSSLMode() string {
	return m.config["database.ssl_mode"].(string)
}

func (m *MockConfigService) IsEmailVerificationRequired() bool {
	return m.config["auth.email_verification_required"].(bool)
}

func (m *MockConfigService) GetVerificationTokenExpiry() time.Duration {
	return m.config["auth.verification_token_expiry"].(time.Duration)
}

func (m *MockConfigService) GetSessionCookieName() string {
	return m.config["auth.session_cookie_name"].(string)
}

func (m *MockConfigService) GetSessionExpiry() time.Duration {
	return m.config["auth.session_expiry"].(time.Duration)
}

func (m *MockConfigService) IsSessionSecure() bool {
	return m.config["auth.session_secure"].(bool)
}

func (m *MockConfigService) IsSessionHttpOnly() bool {
	return m.config["auth.session_http_only"].(bool)
}

func (m *MockConfigService) GetSessionSameSite() string {
	return m.config["auth.session_same_site"].(string)
}

func (m *MockConfigService) GetMinPasswordLength() int {
	return m.config["auth.min_password_length"].(int)
}

func (m *MockConfigService) IsStrongPasswordRequired() bool {
	return m.config["auth.strong_password_required"].(bool)
}

func (m *MockConfigService) ShouldCreateDefaultAdmin() bool {
	return m.config["auth.create_default_admin"].(bool)
}

func (m *MockConfigService) GetDefaultAdminEmail() string {
	return m.config["auth.default_admin_email"].(string)
}

func (m *MockConfigService) GetDefaultAdminPassword() string {
	return m.config["auth.default_admin_password"].(string)
}

func (m *MockConfigService) GetDefaultAdminFirstName() string {
	return m.config["auth.default_admin_first_name"].(string)
}

func (m *MockConfigService) GetDefaultAdminLastName() string {
	return m.config["auth.default_admin_last_name"].(string)
}

func (m *MockConfigService) GetCSRFSecret() string {
	return m.config["security.csrf_secret"].(string)
}

func (m *MockConfigService) IsCSRFSecure() bool {
	return m.config["security.csrf_secure"].(bool)
}

func (m *MockConfigService) IsCSRFHttpOnly() bool {
	return m.config["security.csrf_http_only"].(bool)
}

func (m *MockConfigService) GetCSRFSameSite() string {
	return m.config["security.csrf_same_site"].(string)
}

func (m *MockConfigService) IsRateLimitEnabled() bool {
	return m.config["security.rate_limit_enabled"].(bool)
}

func (m *MockConfigService) GetRateLimitRequests() int {
	return m.config["security.rate_limit_requests"].(int)
}

func (m *MockConfigService) AreSecurityHeadersEnabled() bool {
	return m.config["security.headers_enabled"].(bool)
}

func (m *MockConfigService) IsHSTSEnabled() bool {
	return m.config["security.hsts_enabled"].(bool)
}

func (m *MockConfigService) GetHSTSMaxAge() int {
	return m.config["security.hsts_max_age"].(int)
}

func (m *MockConfigService) GetLogLevel() string {
	return m.config["logging.level"].(string)
}

func (m *MockConfigService) GetLogFormat() string {
	return m.config["logging.format"].(string)
}

func (m *MockConfigService) GetLogOutput() string {
	return m.config["logging.output"].(string)
}

func (m *MockConfigService) IsFileLoggingEnabled() bool {
	return m.config["logging.file_enabled"].(bool)
}

func (m *MockConfigService) GetLogFilePath() string {
	return m.config["logging.file_path"].(string)
}

func (m *MockConfigService) GetSMTPHost() string {
	return m.config["email.smtp_host"].(string)
}

func (m *MockConfigService) GetSMTPPort() int {
	return m.config["email.smtp_port"].(int)
}

func (m *MockConfigService) GetSMTPUsername() string {
	return m.config["email.smtp_username"].(string)
}

func (m *MockConfigService) GetSMTPPassword() string {
	return m.config["email.smtp_password"].(string)
}

func (m *MockConfigService) IsSMTPTLSEnabled() bool {
	return m.config["email.smtp_tls_enabled"].(bool)
}

func (m *MockConfigService) GetFromEmail() string {
	return m.config["email.from_email"].(string)
}

func (m *MockConfigService) GetFromName() string {
	return m.config["email.from_name"].(string)
}

func (m *MockConfigService) GetReplyToEmail() string {
	return m.config["email.reply_to_email"].(string)
}

func (m *MockConfigService) IsEmailDummyModeEnabled() bool {
	return m.config["email.dummy_mode_enabled"].(bool)
}

func (m *MockConfigService) IsDevelopment() bool {
	return m.config["environment.development"].(bool)
}

func (m *MockConfigService) IsProduction() bool {
	return m.config["environment.production"].(bool)
}

// Tests for ConfigService interface
func TestConfigService_ServerConfiguration(t *testing.T) {
	config := NewMockConfigService()

	tests := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"GetServerHost", func() interface{} { return config.GetServerHost() }, "localhost"},
		{"GetServerPort", func() interface{} { return config.GetServerPort() }, 8080},
		{"GetServerBaseURL", func() interface{} { return config.GetServerBaseURL() }, "http://localhost:8080"},
		{"GetServerReadTimeout", func() interface{} { return config.GetServerReadTimeout() }, 30 * time.Second},
		{"GetServerWriteTimeout", func() interface{} { return config.GetServerWriteTimeout() }, 30 * time.Second},
		{"GetServerIdleTimeout", func() interface{} { return config.GetServerIdleTimeout() }, 60 * time.Second},
		{"GetServerShutdownTimeout", func() interface{} { return config.GetServerShutdownTimeout() }, 10 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestConfigService_I18nConfiguration(t *testing.T) {
	config := NewMockConfigService()

	locales := config.GetSupportedLocales()
	expectedLocales := []string{"en", "de", "fr"}
	if len(locales) != len(expectedLocales) {
		t.Errorf("GetSupportedLocales() length = %d, want %d", len(locales), len(expectedLocales))
	}

	defaultLocale := config.GetDefaultLocale()
	if defaultLocale != "en" {
		t.Errorf("GetDefaultLocale() = %s, want en", defaultLocale)
	}

	fallbackLocale := config.GetFallbackLocale()
	if fallbackLocale != "en" {
		t.Errorf("GetFallbackLocale() = %s, want en", fallbackLocale)
	}
}

func TestConfigService_LayoutConfiguration(t *testing.T) {
	config := NewMockConfigService()

	tests := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"GetLayoutRootDirectory", func() interface{} { return config.GetLayoutRootDirectory() }, "app"},
		{"GetLayoutFileName", func() interface{} { return config.GetLayoutFileName() }, "layout.templ"},
		{"GetLayoutAssetsDirectory", func() interface{} { return config.GetLayoutAssetsDirectory() }, "assets"},
		{"GetLayoutAssetsRouteName", func() interface{} { return config.GetLayoutAssetsRouteName() }, "/assets/"},
		{"GetTemplateExtension", func() interface{} { return config.GetTemplateExtension() }, ".templ"},
		{"GetMetadataExtension", func() interface{} { return config.GetMetadataExtension() }, ".yaml"},
		{"IsLayoutInheritanceEnabled", func() interface{} { return config.IsLayoutInheritanceEnabled() }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestConfigService_DatabaseConfiguration(t *testing.T) {
	config := NewMockConfigService()

	tests := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"GetDatabaseHost", func() interface{} { return config.GetDatabaseHost() }, "localhost"},
		{"GetDatabasePort", func() interface{} { return config.GetDatabasePort() }, 5432},
		{"GetDatabaseUser", func() interface{} { return config.GetDatabaseUser() }, "testuser"},
		{"GetDatabasePassword", func() interface{} { return config.GetDatabasePassword() }, "testpass"},
		{"GetDatabaseName", func() interface{} { return config.GetDatabaseName() }, "testdb"},
		{"GetDatabaseSSLMode", func() interface{} { return config.GetDatabaseSSLMode() }, "disable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestConfigService_AuthConfiguration(t *testing.T) {
	config := NewMockConfigService()

	tests := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"IsEmailVerificationRequired", func() interface{} { return config.IsEmailVerificationRequired() }, true},
		{"GetVerificationTokenExpiry", func() interface{} { return config.GetVerificationTokenExpiry() }, 24 * time.Hour},
		{"GetSessionCookieName", func() interface{} { return config.GetSessionCookieName() }, "session_id"},
		{"GetSessionExpiry", func() interface{} { return config.GetSessionExpiry() }, 7 * 24 * time.Hour},
		{"IsSessionSecure", func() interface{} { return config.IsSessionSecure() }, true},
		{"IsSessionHttpOnly", func() interface{} { return config.IsSessionHttpOnly() }, true},
		{"GetSessionSameSite", func() interface{} { return config.GetSessionSameSite() }, "Strict"},
		{"GetMinPasswordLength", func() interface{} { return config.GetMinPasswordLength() }, 8},
		{"IsStrongPasswordRequired", func() interface{} { return config.IsStrongPasswordRequired() }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestConfigService_SecurityConfiguration(t *testing.T) {
	config := NewMockConfigService()

	tests := []struct {
		name     string
		testFunc func() interface{}
		expected interface{}
	}{
		{"GetCSRFSecret", func() interface{} { return config.GetCSRFSecret() }, "csrf-secret-key"},
		{"IsCSRFSecure", func() interface{} { return config.IsCSRFSecure() }, true},
		{"IsRateLimitEnabled", func() interface{} { return config.IsRateLimitEnabled() }, true},
		{"GetRateLimitRequests", func() interface{} { return config.GetRateLimitRequests() }, 100},
		{"AreSecurityHeadersEnabled", func() interface{} { return config.AreSecurityHeadersEnabled() }, true},
		{"IsHSTSEnabled", func() interface{} { return config.IsHSTSEnabled() }, true},
		{"GetHSTSMaxAge", func() interface{} { return config.GetHSTSMaxAge() }, 31536000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestConfigService_EnvironmentConfiguration(t *testing.T) {
	config := NewMockConfigService()

	if !config.IsDevelopment() {
		t.Error("Expected development environment to be true")
	}

	if config.IsProduction() {
		t.Error("Expected production environment to be false")
	}
}

// Test interface compliance
func TestConfigServiceInterfaceCompliance(t *testing.T) {
	var _ ConfigService = (*MockConfigService)(nil)
	
	// This test ensures our mock implements the interface correctly
	config := NewMockConfigService()
	
	// Test that all methods are callable without panicking
	_ = config.GetServerHost()
	_ = config.GetSupportedLocales()
	_ = config.GetLayoutRootDirectory()
	_ = config.GetDatabaseHost()
	_ = config.IsEmailVerificationRequired()
	_ = config.GetCSRFSecret()
	_ = config.GetLogLevel()
	_ = config.GetSMTPHost()
	_ = config.IsDevelopment()
}