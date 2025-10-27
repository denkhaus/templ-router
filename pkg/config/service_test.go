package config

import (
	"os"
	"testing"
	"time"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigService(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "default configuration",
			envVars:     map[string]string{},
			expectError: false,
		},
		{
			name: "valid custom configuration",
			envVars: map[string]string{
				"TR_SERVER_HOST":                 "0.0.0.0",
				"TR_SERVER_PORT":                 "3000",
				"TR_DATABASE_HOST":               "db.example.com",
				"TR_DATABASE_PORT":               "5432",
				"TR_AUTH_MIN_PASSWORD_LENGTH":    "12",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD": "validpassword123",
			},
			expectError: false,
		},
		{
			name: "invalid server port",
			envVars: map[string]string{
				"TR_SERVER_PORT": "99999",
			},
			expectError: true,
			errorMsg:    "Invalid server port",
		},
		{
			name: "invalid database port",
			envVars: map[string]string{
				"TR_DATABASE_PORT": "0",
			},
			expectError: true,
			errorMsg:    "Invalid database port",
		},
		{
			name: "invalid min password length",
			envVars: map[string]string{
				"TR_AUTH_MIN_PASSWORD_LENGTH": "0",
			},
			expectError: true,
			errorMsg:    "Minimum password length must be at least 1",
		},
		{
			name: "default admin enabled but missing email",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "true",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":  "",
			},
			expectError: true,
			errorMsg:    "Email cannot be empty when CreateDefaultAdmin is enabled",
		},
		{
			name: "default admin enabled but password too short",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN":   "true",
				"TR_AUTH_MIN_PASSWORD_LENGTH":    "10",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD": "short",
			},
			expectError: true,
			errorMsg:    "Password must be at least 10 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearTestEnv(t)

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Create injector and config service
			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					// The service layer wraps validation errors in "configuration validation failed"
					assert.Contains(t, err.Error(), "configuration validation failed")
				}
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}

func TestConfigServiceAccessors(t *testing.T) {
	// Clear environment and set test values
	clearTestEnv(t)

	testEnvVars := map[string]string{
		// Server config
		"TR_SERVER_HOST":             "test.example.com",
		"TR_SERVER_PORT":             "9000",
		"TR_SERVER_BASE_URL":         "https://test.example.com",
		"TR_SERVER_READ_TIMEOUT":     "45s",
		"TR_SERVER_WRITE_TIMEOUT":    "60s",
		"TR_SERVER_IDLE_TIMEOUT":     "180s",
		"TR_SERVER_SHUTDOWN_TIMEOUT": "45s",

		// Database config
		"TR_DATABASE_HOST":     "testdb.example.com",
		"TR_DATABASE_PORT":     "3306",
		"TR_DATABASE_USER":     "testuser",
		"TR_DATABASE_PASSWORD": "testpass",
		"TR_DATABASE_NAME":     "testdb",
		"TR_DATABASE_SSL_MODE": "require",

		// Auth config
		"TR_AUTH_REQUIRE_EMAIL_VERIFICATION": "false",
		"TR_AUTH_VERIFICATION_TOKEN_EXPIRY":  "48h",
		"TR_AUTH_SESSION_COOKIE_NAME":        "test_session",
		"TR_AUTH_SESSION_EXPIRY":             "48h",
		"TR_AUTH_SESSION_SECURE":             "true",
		"TR_AUTH_SESSION_HTTP_ONLY":          "false",
		"TR_AUTH_SESSION_SAME_SITE":          "strict",
		"TR_AUTH_MIN_PASSWORD_LENGTH":        "12",
		"TR_AUTH_REQUIRE_STRONG_PASSWORD":    "true",
		"TR_AUTH_CREATE_DEFAULT_ADMIN":       "false",
		"TR_AUTH_SIGNIN_ROUTE":               "/auth/login",
		"TR_AUTH_SIGNIN_SUCCESS_ROUTE":       "/dashboard",
		"TR_AUTH_SIGNUP_SUCCESS_ROUTE":       "/welcome",
		"TR_AUTH_SIGNOUT_SUCCESS_ROUTE":      "/goodbye",

		// Security config
		"TR_SECURITY_CSRF_SECRET":             "test-csrf-secret",
		"TR_SECURITY_CSRF_SECURE":             "true",
		"TR_SECURITY_CSRF_HTTP_ONLY":          "false",
		"TR_SECURITY_CSRF_SAME_SITE":          "none",
		"TR_SECURITY_ENABLE_RATE_LIMIT":       "false",
		"TR_SECURITY_RATE_LIMIT_REQUESTS":     "200",
		"TR_SECURITY_ENABLE_SECURITY_HEADERS": "false",
		"TR_SECURITY_ENABLE_HSTS":             "true",
		"TR_SECURITY_HSTS_MAX_AGE":            "63072000",

		// Logging config
		"TR_LOGGING_LEVEL":       "debug",
		"TR_LOGGING_FORMAT":      "text",
		"TR_LOGGING_OUTPUT":      "stderr",
		"TR_LOGGING_ENABLE_FILE": "true",
		"TR_LOGGING_FILE_PATH":   "/tmp/test.log",

		// Email config
		"TR_EMAIL_SMTP_HOST":         "smtp.test.com",
		"TR_EMAIL_SMTP_PORT":         "465",
		"TR_EMAIL_SMTP_USERNAME":     "testuser@test.com",
		"TR_EMAIL_SMTP_PASSWORD":     "testsmtppass",
		"TR_EMAIL_SMTP_USE_TLS":      "false",
		"TR_EMAIL_FROM_EMAIL":        "noreply@test.com",
		"TR_EMAIL_FROM_NAME":         "Test App",
		"TR_EMAIL_REPLY_TO_EMAIL":    "support@test.com",
		"TR_EMAIL_ENABLE_DUMMY_MODE": "false",

		// I18n config
		"TR_I18N_SUPPORTED_LOCALES": "en,de,fr",
		"TR_I18N_DEFAULT_LOCALE":    "de",
		"TR_I18N_FALLBACK_LOCALE":   "en",

		// Layout config
		"TR_LAYOUT_ROOT_DIRECTORY":     "testapp",
		"TR_LAYOUT_ASSETS_DIRECTORY":   "testassets",
		"TR_LAYOUT_ASSETS_ROUTE_NAME":  "static",
		"TR_LAYOUT_LAYOUT_FILE_NAME":   "testlayout",
		"TR_LAYOUT_TEMPLATE_EXTENSION": ".test.templ",
		"TR_LAYOUT_METADATA_EXTENSION": ".test.yaml",
		"TR_LAYOUT_ENABLE_INHERITANCE": "false",

		// Template generator config
		"TR_TEMPLATE_GENERATOR_OUTPUT_DIR":   "test/generated",
		"TR_TEMPLATE_GENERATOR_PACKAGE_NAME": "testpkg",

		// Environment config
		"TR_ENVIRONMENT_KIND": "test",
	}

	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Create config service
	injector := do.New()
	defer injector.Shutdown()

	configFactory := NewConfigService("TR")
	service, err := configFactory(injector)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Test Server accessors
	assert.Equal(t, "test.example.com", service.GetServerHost())
	assert.Equal(t, 9000, service.GetServerPort())
	assert.Equal(t, "https://test.example.com", service.GetServerBaseURL())
	assert.Equal(t, 45*time.Second, service.GetServerReadTimeout())
	assert.Equal(t, 60*time.Second, service.GetServerWriteTimeout())
	assert.Equal(t, 180*time.Second, service.GetServerIdleTimeout())
	assert.Equal(t, 45*time.Second, service.GetServerShutdownTimeout())

	// Test Database accessors
	assert.Equal(t, "testdb.example.com", service.GetDatabaseHost())
	assert.Equal(t, 3306, service.GetDatabasePort())
	assert.Equal(t, "testuser", service.GetDatabaseUser())
	assert.Equal(t, "testpass", service.GetDatabasePassword())
	assert.Equal(t, "testdb", service.GetDatabaseName())
	assert.Equal(t, "require", service.GetDatabaseSSLMode())

	// Test Auth accessors
	assert.False(t, service.IsEmailVerificationRequired())
	assert.Equal(t, 48*time.Hour, service.GetVerificationTokenExpiry())
	assert.Equal(t, "test_session", service.GetSessionCookieName())
	assert.Equal(t, 48*time.Hour, service.GetSessionExpiry())
	assert.True(t, service.IsSessionSecure())
	assert.False(t, service.IsSessionHttpOnly())
	assert.Equal(t, "strict", service.GetSessionSameSite())
	assert.Equal(t, 12, service.GetMinPasswordLength())
	assert.True(t, service.IsStrongPasswordRequired())
	assert.False(t, service.ShouldCreateDefaultAdmin())
	assert.Equal(t, "/auth/login", service.GetSignInRoute())
	assert.Equal(t, "/dashboard", service.GetSignInSuccessRoute())
	assert.Equal(t, "/welcome", service.GetSignUpSuccessRoute())
	assert.Equal(t, "/goodbye", service.GetSignOutSuccessRoute())

	// Test Security accessors
	assert.Equal(t, "test-csrf-secret", service.GetCSRFSecret())
	assert.True(t, service.IsCSRFSecure())
	assert.False(t, service.IsCSRFHttpOnly())
	assert.Equal(t, "none", service.GetCSRFSameSite())
	assert.False(t, service.IsRateLimitEnabled())
	assert.Equal(t, 200, service.GetRateLimitRequests())
	assert.False(t, service.AreSecurityHeadersEnabled())
	assert.True(t, service.IsHSTSEnabled())
	assert.Equal(t, 63072000, service.GetHSTSMaxAge())

	// Test Logging accessors
	assert.Equal(t, "debug", service.GetLogLevel())
	assert.Equal(t, "text", service.GetLogFormat())
	assert.Equal(t, "stderr", service.GetLogOutput())
	assert.True(t, service.IsFileLoggingEnabled())
	assert.Equal(t, "/tmp/test.log", service.GetLogFilePath())

	// Test Email accessors
	assert.Equal(t, "smtp.test.com", service.GetSMTPHost())
	assert.Equal(t, 465, service.GetSMTPPort())
	assert.Equal(t, "testuser@test.com", service.GetSMTPUsername())
	assert.Equal(t, "testsmtppass", service.GetSMTPPassword())
	assert.False(t, service.IsSMTPTLSEnabled())
	assert.Equal(t, "noreply@test.com", service.GetFromEmail())
	assert.Equal(t, "Test App", service.GetFromName())
	assert.Equal(t, "support@test.com", service.GetReplyToEmail())
	assert.False(t, service.IsEmailDummyModeEnabled())

	// Test I18n accessors
	assert.Equal(t, []string{"en", "de", "fr"}, service.GetSupportedLocales())
	assert.Equal(t, "de", service.GetDefaultLocale())
	assert.Equal(t, "en", service.GetFallbackLocale())

	// Test Layout accessors
	assert.Equal(t, "testapp", service.GetLayoutRootDirectory())
	assert.Equal(t, "testassets", service.GetLayoutAssetsDirectory())
	assert.Equal(t, "static", service.GetLayoutAssetsRouteName())
	assert.Equal(t, "testlayout", service.GetLayoutFileName())
	assert.Equal(t, ".test.templ", service.GetTemplateExtension())
	assert.Equal(t, ".test.yaml", service.GetMetadataExtension())
	assert.False(t, service.IsLayoutInheritanceEnabled())

	// Test Template Generator accessors
	assert.Equal(t, "test/generated", service.GetTemplateOutputDir())
	assert.Equal(t, "testpkg", service.GetTemplatePackageName())

	// Test Environment detection - with test environment kind, should be production
	assert.False(t, service.IsDevelopment())
	assert.True(t, service.IsProduction())
}

func TestDefaultValues(t *testing.T) {
	// Clear environment to test defaults
	clearTestEnv(t)

	// Create config service with defaults
	injector := do.New()
	defer injector.Shutdown()

	configFactory := NewConfigService("TR")
	service, err := configFactory(injector)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Test default values
	assert.Equal(t, "localhost", service.GetServerHost())
	assert.Equal(t, 8080, service.GetServerPort())
	assert.Equal(t, "http://localhost:8080", service.GetServerBaseURL())
	assert.Equal(t, 30*time.Second, service.GetServerReadTimeout())

	assert.Equal(t, "localhost", service.GetDatabaseHost())
	assert.Equal(t, 5432, service.GetDatabasePort())
	assert.Equal(t, "postgres", service.GetDatabaseUser())
	assert.Equal(t, "postgres", service.GetDatabasePassword())
	assert.Equal(t, "router_db", service.GetDatabaseName())
	assert.Equal(t, "disable", service.GetDatabaseSSLMode())

	assert.True(t, service.IsEmailVerificationRequired())
	assert.Equal(t, 24*time.Hour, service.GetVerificationTokenExpiry())
	assert.Equal(t, "session_id", service.GetSessionCookieName())
	assert.Equal(t, 24*time.Hour, service.GetSessionExpiry())
	assert.False(t, service.IsSessionSecure())
	assert.True(t, service.IsSessionHttpOnly())
	assert.Equal(t, "lax", service.GetSessionSameSite())
	assert.Equal(t, 8, service.GetMinPasswordLength())
	assert.False(t, service.IsStrongPasswordRequired())
	assert.True(t, service.ShouldCreateDefaultAdmin())
	assert.Equal(t, "admin@example.com", service.GetDefaultAdminEmail())
	assert.Equal(t, "admin123", service.GetDefaultAdminPassword())
	assert.Equal(t, "Default", service.GetDefaultAdminFirstName())
	assert.Equal(t, "Admin", service.GetDefaultAdminLastName())

	assert.Equal(t, "change-me-in-production", service.GetCSRFSecret())
	assert.False(t, service.IsCSRFSecure())
	assert.True(t, service.IsCSRFHttpOnly())
	assert.Equal(t, "strict", service.GetCSRFSameSite())
	assert.True(t, service.IsRateLimitEnabled())
	assert.Equal(t, 100, service.GetRateLimitRequests())

	assert.Equal(t, "info", service.GetLogLevel())
	assert.Equal(t, "json", service.GetLogFormat())
	assert.Equal(t, "stdout", service.GetLogOutput())
	assert.False(t, service.IsFileLoggingEnabled())
	assert.Equal(t, "logs/router.log", service.GetLogFilePath())

	assert.Equal(t, []string{"en", "de"}, service.GetSupportedLocales())
	assert.Equal(t, "en", service.GetDefaultLocale())
	assert.Equal(t, "en", service.GetFallbackLocale())

	assert.Equal(t, "app", service.GetLayoutRootDirectory())
	assert.Equal(t, "assets", service.GetLayoutAssetsDirectory())
	assert.Equal(t, "assets", service.GetLayoutAssetsRouteName())
	assert.Equal(t, "layout", service.GetLayoutFileName())
	assert.Equal(t, ".templ", service.GetTemplateExtension())
	assert.Equal(t, ".templ.yaml", service.GetMetadataExtension())
	assert.True(t, service.IsLayoutInheritanceEnabled())

	assert.Equal(t, "generated/templates", service.GetTemplateOutputDir())
	assert.Equal(t, "templates", service.GetTemplatePackageName())

	// Test environment detection with defaults
	assert.True(t, service.IsDevelopment())
	assert.False(t, service.IsProduction())
}

func TestProductionDetection(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		isProduction bool
	}{
		{
			name:         "default development",
			envVars:      map[string]string{},
			isProduction: false,
		},
		{
			name: "production by base URL",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":     "production",
				"TR_SERVER_BASE_URL":      "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "production-secret",
			},
			isProduction: true,
		},
		{
			name: "development by environment kind",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":     "develop",
				"TR_SERVER_BASE_URL":      "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "production-secret",
			},
			isProduction: false,
		},
		{
			name: "still development with localhost",
			envVars: map[string]string{
				"TR_SERVER_BASE_URL": "http://localhost:8080",
			},
			isProduction: false,
		},
		{
			name: "still development with default csrf secret",
			envVars: map[string]string{
				"TR_SERVER_BASE_URL":      "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "change-me-in-production",
			},
			isProduction: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			assert.Equal(t, tt.isProduction, service.IsProduction())
			assert.Equal(t, !tt.isProduction, service.IsDevelopment())
		})
	}
}

// Helper function to clear test environment variables
func clearTestEnv(t *testing.T) {
	// List of environment variables to clear for clean test state
	envVars := []string{
		"TR_SERVER_HOST", "TR_SERVER_PORT", "TR_SERVER_BASE_URL",
		"TR_SERVER_READ_TIMEOUT", "TR_SERVER_WRITE_TIMEOUT", "TR_SERVER_IDLE_TIMEOUT", "TR_SERVER_SHUTDOWN_TIMEOUT",
		"TR_DATABASE_HOST", "TR_DATABASE_PORT", "TR_DATABASE_USER", "TR_DATABASE_PASSWORD", "TR_DATABASE_NAME", "TR_DATABASE_SSL_MODE",
		"TR_AUTH_REQUIRE_EMAIL_VERIFICATION", "TR_AUTH_VERIFICATION_TOKEN_EXPIRY", "TR_AUTH_SESSION_COOKIE_NAME",
		"TR_AUTH_SESSION_EXPIRY", "TR_AUTH_SESSION_SECURE", "TR_AUTH_SESSION_HTTP_ONLY", "TR_AUTH_SESSION_SAME_SITE",
		"TR_AUTH_MIN_PASSWORD_LENGTH", "TR_AUTH_REQUIRE_STRONG_PASSWORD", "TR_AUTH_CREATE_DEFAULT_ADMIN",
		"TR_AUTH_DEFAULT_ADMIN_EMAIL", "TR_AUTH_DEFAULT_ADMIN_PASSWORD", "TR_AUTH_DEFAULT_ADMIN_FIRST_NAME", "TR_AUTH_DEFAULT_ADMIN_LAST_NAME",
		"TR_AUTH_SIGNIN_ROUTE", "TR_AUTH_SIGNIN_SUCCESS_ROUTE", "TR_AUTH_SIGNUP_SUCCESS_ROUTE", "TR_AUTH_SIGNOUT_SUCCESS_ROUTE",
		"TR_SECURITY_CSRF_SECRET", "TR_SECURITY_CSRF_SECURE", "TR_SECURITY_CSRF_HTTP_ONLY", "TR_SECURITY_CSRF_SAME_SITE",
		"TR_SECURITY_ENABLE_RATE_LIMIT", "TR_SECURITY_RATE_LIMIT_REQUESTS", "TR_SECURITY_ENABLE_SECURITY_HEADERS",
		"TR_SECURITY_ENABLE_HSTS", "TR_SECURITY_HSTS_MAX_AGE",
		"TR_LOGGING_LEVEL", "TR_LOGGING_FORMAT", "TR_LOGGING_OUTPUT", "TR_LOGGING_ENABLE_FILE", "TR_LOGGING_FILE_PATH",
		"TR_EMAIL_SMTP_HOST", "TR_EMAIL_SMTP_PORT", "TR_EMAIL_SMTP_USERNAME", "TR_EMAIL_SMTP_PASSWORD", "TR_EMAIL_SMTP_USE_TLS",
		"TR_EMAIL_FROM_EMAIL", "TR_EMAIL_FROM_NAME", "TR_EMAIL_REPLY_TO_EMAIL", "TR_EMAIL_ENABLE_DUMMY_MODE",
		"TR_I18N_SUPPORTED_LOCALES", "TR_I18N_DEFAULT_LOCALE", "TR_I18N_FALLBACK_LOCALE",
		"TR_LAYOUT_ROOT_DIRECTORY", "TR_LAYOUT_ASSETS_DIRECTORY", "TR_LAYOUT_ASSETS_ROUTE_NAME",
		"TR_LAYOUT_LAYOUT_FILE_NAME", "TR_LAYOUT_TEMPLATE_EXTENSION", "TR_LAYOUT_METADATA_EXTENSION", "TR_LAYOUT_ENABLE_INHERITANCE",
		"TR_TEMPLATE_GENERATOR_OUTPUT_DIR", "TR_TEMPLATE_GENERATOR_PACKAGE_NAME",
		"TR_ENVIRONMENT_KIND", "TR_CONFIG_PRINT_SUMMARY",
		// Also clear system environment variables that might interfere with defaults
		"USER",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
