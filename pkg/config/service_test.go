package config

import (
	"os"
	"strings"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

		// Auth config
		"TR_AUTH_SIGNIN_ROUTE":          "/auth/login",
		"TR_AUTH_SIGNIN_SUCCESS_ROUTE":  "/dashboard",
		"TR_AUTH_SIGNUP_SUCCESS_ROUTE":  "/welcome",
		"TR_AUTH_SIGNOUT_SUCCESS_ROUTE": "/goodbye",

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

	// Test Auth route accessors (only route configuration remains in current implementation)
	assert.Equal(t, "/auth/login", service.GetSignInRoute())
	assert.Equal(t, "/dashboard", service.GetSignInSuccessRoute())
	assert.Equal(t, "/welcome", service.GetSignUpSuccessRoute())
	assert.Equal(t, "/goodbye", service.GetSignOutSuccessRoute())

	// Test Security accessors
	assert.Equal(t, "test-csrf-secret", service.GetCSRFSecret())
	assert.True(t, service.IsCSRFSecure())
	assert.False(t, service.IsCSRFHttpOnly())
	assert.Equal(t, "none", service.GetCSRFSameSite())
	assert.False(t, service.AreSecurityHeadersEnabled())
	assert.True(t, service.IsHSTSEnabled())
	assert.Equal(t, 63072000, service.GetHSTSMaxAge())

	// Test Logging accessors
	assert.Equal(t, "debug", service.GetLogLevel())
	assert.Equal(t, "text", service.GetLogFormat())
	assert.Equal(t, "stderr", service.GetLogOutput())
	assert.True(t, service.IsFileLoggingEnabled())
	assert.Equal(t, "/tmp/test.log", service.GetLogFilePath())

	// Test I18n accessors
	assert.Equal(t, []string{"en", "de", "fr"}, service.GetSupportedLocales())
	assert.Equal(t, "de", service.GetDefaultLocale())
	assert.Equal(t, "en", service.GetFallbackLocale())

	// Test Layout accessors
	// GetLayoutRootDirectory returns absolute path, so check base name
	assert.True(t, strings.HasSuffix(service.GetLayoutRootDirectory(), "testapp"))
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

	assert.Equal(t, "change-me-in-production", service.GetCSRFSecret())
	assert.False(t, service.IsCSRFSecure())
	assert.True(t, service.IsCSRFHttpOnly())
	assert.Equal(t, "strict", service.GetCSRFSameSite())

	assert.Equal(t, "info", service.GetLogLevel())
	assert.Equal(t, "json", service.GetLogFormat())
	assert.Equal(t, "stdout", service.GetLogOutput())
	assert.False(t, service.IsFileLoggingEnabled())
	assert.Equal(t, "logs/router.log", service.GetLogFilePath())

	assert.Equal(t, []string{"en", "de"}, service.GetSupportedLocales())
	assert.Equal(t, "en", service.GetDefaultLocale())
	assert.Equal(t, "en", service.GetFallbackLocale())

	// GetLayoutRootDirectory returns absolute path, so check base name
	assert.True(t, strings.HasSuffix(service.GetLayoutRootDirectory(), "app"))
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
		"TR_AUTH_SIGNIN_ROUTE", "TR_AUTH_SIGNIN_SUCCESS_ROUTE", "TR_AUTH_SIGNUP_SUCCESS_ROUTE", "TR_AUTH_SIGNOUT_SUCCESS_ROUTE",
		"TR_SECURITY_CSRF_SECRET", "TR_SECURITY_CSRF_SECURE", "TR_SECURITY_CSRF_HTTP_ONLY", "TR_SECURITY_CSRF_SAME_SITE",
		"TR_SECURITY_ENABLE_RATE_LIMIT", "TR_SECURITY_RATE_LIMIT_REQUESTS", "TR_SECURITY_ENABLE_SECURITY_HEADERS",
		"TR_SECURITY_ENABLE_HSTS", "TR_SECURITY_HSTS_MAX_AGE",
		"TR_LOGGING_LEVEL", "TR_LOGGING_FORMAT", "TR_LOGGING_OUTPUT", "TR_LOGGING_ENABLE_FILE", "TR_LOGGING_FILE_PATH",
		"TR_I18N_SUPPORTED_LOCALES", "TR_I18N_DEFAULT_LOCALE", "TR_I18N_FALLBACK_LOCALE",
		"TR_LAYOUT_ROOT_DIRECTORY", "TR_LAYOUT_ASSETS_DIRECTORY", "TR_LAYOUT_ASSETS_ROUTE_NAME",
		"TR_LAYOUT_LAYOUT_FILE_NAME", "TR_LAYOUT_TEMPLATE_EXTENSION", "TR_LAYOUT_METADATA_EXTENSION", "TR_LAYOUT_ENABLE_INHERITANCE",
		"TR_TEMPLATE_GENERATOR_OUTPUT_DIR", "TR_TEMPLATE_GENERATOR_PACKAGE_NAME",
		"TR_ENVIRONMENT_KIND", "TR_CONFIG_PRINT_SUMMARY",
		// Also clear system environment variables that might interfere with defaults
		"USER", "NAME",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
