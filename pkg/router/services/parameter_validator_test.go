package services

import (
	"testing"
	"time"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock config service for tests
type mockConfigService struct{}

// Router configuration methods
func (m *mockConfigService) GetRouterEnableTrailingSlash() bool     { return true }
func (m *mockConfigService) GetRouterEnableSlashRedirect() bool     { return true }
func (m *mockConfigService) GetRouterEnableMethodNotAllowed() bool  { return true }

// Implement all required ConfigService methods
func (m *mockConfigService) GetLayoutRootDirectory() string            { return "app" }
func (m *mockConfigService) GetSupportedLocales() []string             { return []string{"en", "de"} }
func (m *mockConfigService) GetDefaultLocale() string                  { return "en" }
func (m *mockConfigService) GetFallbackLocale() string                 { return "en" }
func (m *mockConfigService) GetLayoutFileName() string                 { return "layout" }
func (m *mockConfigService) GetTemplateExtension() string              { return ".templ" }
func (m *mockConfigService) GetMetadataExtension() string              { return ".yaml" }
func (m *mockConfigService) IsLayoutInheritanceEnabled() bool          { return true }
func (m *mockConfigService) GetTemplateOutputDir() string              { return "generated" }
func (m *mockConfigService) GetTemplatePackageName() string            { return "templates" }
func (m *mockConfigService) GetLayoutAssetsDirectory() string          { return "assets" }
func (m *mockConfigService) GetLayoutAssetsRouteName() string          { return "/assets/" }
func (m *mockConfigService) IsDevelopment() bool                       { return true }
func (m *mockConfigService) IsProduction() bool                        { return false }
func (m *mockConfigService) GetServerHost() string                     { return "localhost" }
func (m *mockConfigService) GetServerPort() int                        { return 8080 }
func (m *mockConfigService) GetServerBaseURL() string                  { return "http://localhost:8080" }
func (m *mockConfigService) GetServerReadTimeout() time.Duration       { return 30 * time.Second }
func (m *mockConfigService) GetServerWriteTimeout() time.Duration      { return 30 * time.Second }
func (m *mockConfigService) GetServerIdleTimeout() time.Duration       { return 60 * time.Second }
func (m *mockConfigService) GetServerShutdownTimeout() time.Duration   { return 10 * time.Second }
func (m *mockConfigService) GetDatabaseHost() string                   { return "localhost" }
func (m *mockConfigService) GetDatabasePort() int                      { return 5432 }
func (m *mockConfigService) GetDatabaseUser() string                   { return "user" }
func (m *mockConfigService) GetDatabasePassword() string               { return "password" }
func (m *mockConfigService) GetDatabaseName() string                   { return "testdb" }
func (m *mockConfigService) GetDatabaseSSLMode() string                { return "disable" }
func (m *mockConfigService) IsEmailVerificationRequired() bool         { return false }
func (m *mockConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *mockConfigService) GetSessionCookieName() string              { return "session" }
func (m *mockConfigService) GetSessionExpiry() time.Duration           { return 24 * time.Hour }
func (m *mockConfigService) IsSessionSecure() bool                     { return false }
func (m *mockConfigService) IsSessionHttpOnly() bool                   { return true }
func (m *mockConfigService) GetSessionSameSite() string                { return "Lax" }
func (m *mockConfigService) GetMinPasswordLength() int                 { return 8 }
func (m *mockConfigService) IsStrongPasswordRequired() bool            { return false }
func (m *mockConfigService) ShouldCreateDefaultAdmin() bool            { return false }
func (m *mockConfigService) GetDefaultAdminEmail() string              { return "" }
func (m *mockConfigService) GetDefaultAdminPassword() string           { return "" }
func (m *mockConfigService) GetDefaultAdminFirstName() string          { return "" }
func (m *mockConfigService) GetDefaultAdminLastName() string           { return "" }
func (m *mockConfigService) GetSignInRoute() string                    { return "/login" }
func (m *mockConfigService) GetSignInSuccessRoute() string             { return "/dashboard" }
func (m *mockConfigService) GetSignUpSuccessRoute() string             { return "/welcome" }
func (m *mockConfigService) GetSignOutSuccessRoute() string            { return "/" }
func (m *mockConfigService) GetCSRFSecret() string                     { return "secret" }
func (m *mockConfigService) IsCSRFSecure() bool                        { return false }
func (m *mockConfigService) IsCSRFHttpOnly() bool                      { return true }
func (m *mockConfigService) GetCSRFSameSite() string                   { return "Lax" }
func (m *mockConfigService) IsRateLimitEnabled() bool                  { return false }
func (m *mockConfigService) GetRateLimitRequests() int                 { return 100 }
func (m *mockConfigService) AreSecurityHeadersEnabled() bool           { return false }
func (m *mockConfigService) IsHSTSEnabled() bool                       { return false }
func (m *mockConfigService) GetHSTSMaxAge() int                        { return 31536000 }
func (m *mockConfigService) GetLogLevel() string                       { return "info" }
func (m *mockConfigService) GetLogFormat() string                      { return "json" }
func (m *mockConfigService) GetLogOutput() string                      { return "stdout" }
func (m *mockConfigService) IsFileLoggingEnabled() bool                { return false }
func (m *mockConfigService) GetLogFilePath() string                    { return "" }
func (m *mockConfigService) GetSMTPHost() string                       { return "" }
func (m *mockConfigService) GetSMTPPort() int                          { return 587 }
func (m *mockConfigService) GetSMTPUsername() string                   { return "" }
func (m *mockConfigService) GetSMTPPassword() string                   { return "" }
func (m *mockConfigService) IsSMTPTLSEnabled() bool                    { return true }
func (m *mockConfigService) GetFromEmail() string                      { return "" }
func (m *mockConfigService) GetFromName() string                       { return "" }
func (m *mockConfigService) GetReplyToEmail() string                   { return "" }
func (m *mockConfigService) IsEmailDummyModeEnabled() bool             { return true }

func TestParameterValidator_ValidateParameterInheritance(t *testing.T) {
	// Setup
	injector := do.New()
	defer injector.Shutdown()

	// Provide logger and config
	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockConfigService{}, nil
	})

	validator, err := NewParameterValidator(injector)
	require.NoError(t, err)

	tests := []struct {
		name           string
		route          *interfaces.Route
		config         *interfaces.ConfigFile
		hierarchyMap   map[string][]string
		configs        map[string]*interfaces.ConfigFile
		expectedErrors int
		expectedWarnings int
		expectedWarningTypes []string
	}{
		{
			name: "no inheritance - no conflicts",
			route: &interfaces.Route{
				Path:         "/user/profile",
				TemplateFile: "user/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[0-9]+$",
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{},
			configs:      map[string]*interfaces.ConfigFile{},
			expectedErrors: 0,
			expectedWarnings: 0,
		},
		{
			name: "compatible parameter inheritance",
			route: &interfaces.Route{
				Path:         "/user/$id/profile",
				TemplateFile: "user/id_/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[0-9]+$",
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/profile": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								Validation: "^[0-9]+$",
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 0,
		},
		{
			name: "conflicting validation regex",
			route: &interfaces.Route{
				Path:         "/user/$id/profile",
				TemplateFile: "user/id_/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[a-z]+$", // Different from parent
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/profile": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								Validation: "^[0-9]+$", // Different from child
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"PARAMETER_VALIDATION_CONFLICT"},
		},
		{
			name: "conflicting supported values",
			route: &interfaces.Route{
				Path:         "/user/$role/dashboard",
				TemplateFile: "user/role_/dashboard/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"role": {
							SupportedValues: []string{"admin", "user", "guest", "superuser"}, // Contains values not in parent
							Description: "User role",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$role/dashboard": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"role": {
								SupportedValues: []string{"admin", "user", "guest"}, // Subset of child
								Description: "User role",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"PARAMETER_VALUES_CONFLICT"},
		},
		{
			name: "child adds validation (good practice)",
			route: &interfaces.Route{
				Path:         "/user/$id/settings",
				TemplateFile: "user/id_/settings/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[0-9]+$", // Child adds validation
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/settings": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								// No validation in parent
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"PARAMETER_INHERITANCE_INFO"},
		},
		{
			name: "missing inherited parameter config",
			route: &interfaces.Route{
				Path:         "/user/$id/profile",
				TemplateFile: "user/id_/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						// Missing "id" parameter config
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/profile": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								Validation: "^[0-9]+$",
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"INHERITED_PARAMETER_MISSING_CONFIG"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &ValidationResult{
				Errors:   make([]ValidationError, 0),
				Warnings: make([]ValidationWarning, 0),
			}

			validator.ValidateParameterInheritance(tt.route, tt.config, tt.hierarchyMap, tt.configs, result)

			assert.Equal(t, tt.expectedErrors, len(result.Errors), "Expected %d errors, got %d", tt.expectedErrors, len(result.Errors))
			assert.Equal(t, tt.expectedWarnings, len(result.Warnings), "Expected %d warnings, got %d", tt.expectedWarnings, len(result.Warnings))

			// Check warning types
			if len(tt.expectedWarningTypes) > 0 {
				actualWarningTypes := make([]string, len(result.Warnings))
				for i, warning := range result.Warnings {
					actualWarningTypes[i] = warning.Type
				}
				assert.ElementsMatch(t, tt.expectedWarningTypes, actualWarningTypes, "Warning types don't match")
			}
		})
	}
}

func TestParameterValidator_IsSubsetOfValues(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockConfigService{}, nil
	})

	validator, err := NewParameterValidator(injector)
	require.NoError(t, err)

	tests := []struct {
		name         string
		childValues  []string
		parentValues []string
		expected     bool
	}{
		{
			name:         "empty child values",
			childValues:  []string{},
			parentValues: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:         "child is subset of parent",
			childValues:  []string{"a", "b"},
			parentValues: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:         "child equals parent",
			childValues:  []string{"a", "b", "c"},
			parentValues: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:         "child has extra values",
			childValues:  []string{"a", "b", "c", "d"},
			parentValues: []string{"a", "b", "c"},
			expected:     false,
		},
		{
			name:         "completely different values",
			childValues:  []string{"x", "y", "z"},
			parentValues: []string{"a", "b", "c"},
			expected:     false,
		},
		{
			name:         "partial overlap",
			childValues:  []string{"a", "x"},
			parentValues: []string{"a", "b", "c"},
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isSubsetOfValues(tt.childValues, tt.parentValues)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParameterValidator_IsConfigForRoute(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockConfigService{}, nil
	})

	validator, err := NewParameterValidator(injector)
	require.NoError(t, err)

	tests := []struct {
		name         string
		templateFile string
		routePath    string
		expected     bool
	}{
		{
			name:         "exact match",
			templateFile: "user/page.templ",
			routePath:    "/user",
			expected:     true,
		},
		{
			name:         "nested route",
			templateFile: "user/profile/page.templ",
			routePath:    "/user/profile",
			expected:     true,
		},
		{
			name:         "dynamic parameter",
			templateFile: "user/id_/page.templ",
			routePath:    "/user/$id",
			expected:     true, // "user" segment matches
		},
		{
			name:         "no match",
			templateFile: "admin/page.templ",
			routePath:    "/user",
			expected:     false,
		},
		{
			name:         "partial match",
			templateFile: "user/settings/page.templ",
			routePath:    "/user/profile",
			expected:     false, // "profile" not in template file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isConfigForRoute(tt.templateFile, tt.routePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

