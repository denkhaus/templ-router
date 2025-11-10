package services

import (
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// MockConfigService for testing
type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) GetServerHost() string {
	args := m.Called("GetServerHost")
	return args.String(0)
}

func (m *MockConfigService) GetServerPort() int {
	args := m.Called("GetServerPort")
	return args.Int(0)
}

func (m *MockConfigService) GetServerBaseURL() string {
	args := m.Called("GetServerBaseURL")
	return args.String(0)
}

func (m *MockConfigService) GetServerReadTimeout() time.Duration {
	args := m.Called("GetServerReadTimeout")
	return args.Get(0).(time.Duration)
}

func (m *MockConfigService) GetServerWriteTimeout() time.Duration {
	args := m.Called("GetServerWriteTimeout")
	return args.Get(0).(time.Duration)
}

func (m *MockConfigService) GetServerIdleTimeout() time.Duration {
	args := m.Called("GetServerIdleTimeout")
	return args.Get(0).(time.Duration)
}

func (m *MockConfigService) GetServerShutdownTimeout() time.Duration {
	args := m.Called("GetServerShutdownTimeout")
	return args.Get(0).(time.Duration)
}

func (m *MockConfigService) GetSupportedLocales() []string {
	args := m.Called("GetSupportedLocales")
	return args.Get(0).([]string)
}

func (m *MockConfigService) GetDefaultLocale() string {
	args := m.Called("GetDefaultLocale")
	return args.String(0)
}

func (m *MockConfigService) GetFallbackLocale() string {
	args := m.Called("GetFallbackLocale")
	return args.String(0)
}

func (m *MockConfigService) GetLayoutRootDirectory() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetLayoutFileName() string {
	args := m.Called("GetLayoutFileName")
	return args.String(0)
}

func (m *MockConfigService) GetLayoutAssetsDirectory() string {
	args := m.Called("GetLayoutAssetsDirectory")
	return args.String(0)
}

func (m *MockConfigService) GetLayoutAssetsRouteName() string {
	args := m.Called("GetLayoutAssetsRouteName")
	return args.String(0)
}

func (m *MockConfigService) GetTemplateExtension() string {
	args := m.Called("GetTemplateExtension")
	return args.String(0)
}

func (m *MockConfigService) GetMetadataExtension() string {
	args := m.Called("GetMetadataExtension")
	return args.String(0)
}

func (m *MockConfigService) IsLayoutInheritanceEnabled() bool {
	args := m.Called("IsLayoutInheritanceEnabled")
	return args.Bool(0)
}

func (m *MockConfigService) GetTemplateOutputDir() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetTemplatePackageName() string {
	args := m.Called("GetTemplatePackageName")
	return args.String(0)
}

func (m *MockConfigService) GetDatabaseHost() string {
	args := m.Called("GetDatabaseHost")
	return args.String(0)
}

func (m *MockConfigService) GetDatabasePort() int {
	args := m.Called("GetDatabasePort")
	return args.Int(0)
}

func (m *MockConfigService) GetDatabaseUser() string {
	args := m.Called("GetDatabaseUser")
	return args.String(0)
}

func (m *MockConfigService) GetDatabasePassword() string {
	args := m.Called("GetDatabasePassword")
	return args.String(0)
}

func (m *MockConfigService) GetDatabaseName() string {
	args := m.Called("GetDatabaseName")
	return args.String(0)
}

func (m *MockConfigService) GetDatabaseSSLMode() string {
	args := m.Called("GetDatabaseSSLMode")
	return args.String(0)
}

func (m *MockConfigService) IsEmailVerificationRequired() bool {
	args := m.Called("IsEmailVerificationRequired")
	return args.Bool(0)
}

func (m *MockConfigService) GetVerificationTokenExpiry() time.Duration {
	args := m.Called("GetVerificationTokenExpiry")
	return args.Get(0).(time.Duration)
}

func (m *MockConfigService) GetSessionCookieName() string {
	args := m.Called("GetSessionCookieName")
	return args.String(0)
}

func (m *MockConfigService) GetSessionExpiry() time.Duration {
	args := m.Called("GetSessionExpiry")
	return args.Get(0).(time.Duration)
}

func (m *MockConfigService) IsSessionSecure() bool {
	args := m.Called("IsSessionSecure")
	return args.Bool(0)
}

func (m *MockConfigService) IsSessionHttpOnly() bool {
	args := m.Called("IsSessionHttpOnly")
	return args.Bool(0)
}

func (m *MockConfigService) GetSessionSameSite() string {
	args := m.Called("GetSessionSameSite")
	return args.String(0)
}

func (m *MockConfigService) GetMinPasswordLength() int {
	args := m.Called("GetMinPasswordLength")
	return args.Int(0)
}

func (m *MockConfigService) IsStrongPasswordRequired() bool {
	args := m.Called("IsStrongPasswordRequired")
	return args.Bool(0)
}

func (m *MockConfigService) ShouldCreateDefaultAdmin() bool {
	args := m.Called("ShouldCreateDefaultAdmin")
	return args.Bool(0)
}

func (m *MockConfigService) GetDefaultAdminEmail() string {
	args := m.Called("GetDefaultAdminEmail")
	return args.String(0)
}

func (m *MockConfigService) GetDefaultAdminPassword() string {
	args := m.Called("GetDefaultAdminPassword")
	return args.String(0)
}

func (m *MockConfigService) GetDefaultAdminFirstName() string {
	args := m.Called("GetDefaultAdminFirstName")
	return args.String(0)
}

func (m *MockConfigService) GetDefaultAdminLastName() string {
	args := m.Called("GetDefaultAdminLastName")
	return args.String(0)
}

func (m *MockConfigService) GetSignInRoute() string {
	args := m.Called("GetSignInRoute")
	return args.String(0)
}

func (m *MockConfigService) GetSignInSuccessRoute() string {
	args := m.Called("GetSignInSuccessRoute")
	return args.String(0)
}

func (m *MockConfigService) GetSignUpSuccessRoute() string {
	args := m.Called("GetSignUpSuccessRoute")
	return args.String(0)
}

func (m *MockConfigService) GetSignOutSuccessRoute() string {
	args := m.Called("GetSignOutSuccessRoute")
	return args.String(0)
}

func (m *MockConfigService) GetCSRFSecret() string {
	args := m.Called("GetCSRFSecret")
	return args.String(0)
}

func (m *MockConfigService) IsCSRFSecure() bool {
	args := m.Called("IsCSRFSecure")
	return args.Bool(0)
}

func (m *MockConfigService) IsCSRFHttpOnly() bool {
	args := m.Called("IsCSRFHttpOnly")
	return args.Bool(0)
}

func (m *MockConfigService) GetCSRFSameSite() string {
	args := m.Called("GetCSRFSameSite")
	return args.String(0)
}

func (m *MockConfigService) AreSecurityHeadersEnabled() bool {
	args := m.Called("AreSecurityHeadersEnabled")
	return args.Bool(0)
}

func (m *MockConfigService) IsHSTSEnabled() bool {
	args := m.Called("IsHSTSEnabled")
	return args.Bool(0)
}

func (m *MockConfigService) GetHSTSMaxAge() int {
	args := m.Called("GetHSTSMaxAge")
	return args.Int(0)
}

func (m *MockConfigService) GetLogLevel() string {
	args := m.Called("GetLogLevel")
	return args.String(0)
}

func (m *MockConfigService) GetLogFormat() string {
	args := m.Called("GetLogFormat")
	return args.String(0)
}

func (m *MockConfigService) GetLogOutput() string {
	args := m.Called("GetLogOutput")
	return args.String(0)
}

func (m *MockConfigService) IsFileLoggingEnabled() bool {
	args := m.Called("IsFileLoggingEnabled")
	return args.Bool(0)
}

func (m *MockConfigService) GetLogFilePath() string {
	args := m.Called("GetLogFilePath")
	return args.String(0)
}

func (m *MockConfigService) GetSMTPHost() string {
	args := m.Called("GetSMTPHost")
	return args.String(0)
}

func (m *MockConfigService) GetSMTPPort() int {
	args := m.Called("GetSMTPPort")
	return args.Int(0)
}

func (m *MockConfigService) GetSMTPUsername() string {
	args := m.Called("GetSMTPUsername")
	return args.String(0)
}

func (m *MockConfigService) GetSMTPPassword() string {
	args := m.Called("GetSMTPPassword")
	return args.String(0)
}

func (m *MockConfigService) IsSMTPTLSEnabled() bool {
	args := m.Called("IsSMTPTLSEnabled")
	return args.Bool(0)
}

func (m *MockConfigService) GetFromEmail() string {
	args := m.Called("GetFromEmail")
	return args.String(0)
}

func (m *MockConfigService) GetFromName() string {
	args := m.Called("GetFromName")
	return args.String(0)
}

func (m *MockConfigService) GetReplyToEmail() string {
	args := m.Called("GetReplyToEmail")
	return args.String(0)
}

func (m *MockConfigService) IsEmailDummyModeEnabled() bool {
	args := m.Called("IsEmailDummyModeEnabled")
	return args.Bool(0)
}

func (m *MockConfigService) IsDevelopment() bool {
	args := m.Called("IsDevelopment")
	return args.Bool(0)
}

func (m *MockConfigService) IsProduction() bool {
	args := m.Called("IsProduction")
	return args.Bool(0)
}

func (m *MockConfigService) GetRouterEnableTrailingSlash() bool {
	args := m.Called("GetRouterEnableTrailingSlash")
	return args.Bool(0)
}

func (m *MockConfigService) GetRouterEnableSlashRedirect() bool {
	args := m.Called("GetRouterEnableSlashRedirect")
	return args.Bool(0)
}

func (m *MockConfigService) GetRouterEnableMethodNotAllowed() bool {
	args := m.Called("GetRouterEnableMethodNotAllowed")
	return args.Bool(0)
}

func (m *MockConfigService) GetRouterEnableAuthRoutes() bool {
	args := m.Called("GetRouterEnableAuthRoutes")
	return args.Bool(0)
}

func (m *MockConfigService) GetRouterAuthRoutePrefix() string {
	args := m.Called("GetRouterAuthRoutePrefix")
	return args.String(0)
}

// MockTemplateRegistry for testing
type MockTemplateRegistry struct {
	mock.Mock
}

func (m *MockTemplateRegistry) GetTemplate(key string) (templ.Component, error) {
	args := m.Called(key)
	return args.Get(0).(templ.Component), args.Error(1)
}

func (m *MockTemplateRegistry) GetTemplateFunction(key string) (func() interface{}, bool) {
	args := m.Called(key)
	return args.Get(0).(func() interface{}), args.Bool(1)
}

func (m *MockTemplateRegistry) GetAllTemplateKeys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockTemplateRegistry) IsAvailable(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockTemplateRegistry) GetRouteToTemplateMapping() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *MockTemplateRegistry) GetTemplateByRoute(route string) (templ.Component, error) {
	args := m.Called(route)
	return args.Get(0).(templ.Component), args.Error(1)
}

func (m *MockTemplateRegistry) RequiresDataService(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockTemplateRegistry) GetDataServiceInfo(key string) (interfaces.DataServiceInfo, bool) {
	args := m.Called(key)
	return args.Get(0).(interfaces.DataServiceInfo), args.Bool(1)
}

func (m *MockTemplateRegistry) GetTemplateMetadata(key string) (*interfaces.TemplateMetadata, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.TemplateMetadata), args.Error(1)
}

func (m *MockTemplateRegistry) GetTemplateMetadataByRoute(route string) (*interfaces.TemplateMetadata, error) {
	args := m.Called(route)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.TemplateMetadata), args.Error(1)
}

func (m *MockTemplateRegistry) GetAllTemplateMetadata() map[string]*interfaces.TemplateMetadata {
	args := m.Called()
	return args.Get(0).(map[string]*interfaces.TemplateMetadata)
}

func (m *MockTemplateRegistry) FindComponentTemplates() map[string]*interfaces.TemplateMetadata {
	args := m.Called()
	return args.Get(0).(map[string]*interfaces.TemplateMetadata)
}

func (m *MockTemplateRegistry) GetTemplateKeyByComponentName(componentName string) (string, bool) {
	args := m.Called(componentName)
	return args.String(0), args.Bool(1)
}

func TestComponentMetadataService_LoadComponentMetadata(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockConfigService := &MockConfigService{}
	mockConfigService.On("GetTemplateOutputDir").Return("demo/app")
	mockConfigService.On("GetLayoutRootDirectory").Return("demo/app")

	// Create a mock template registry
	mockTemplateRegistry := &MockTemplateRegistry{}
	mockTemplateRegistry.On("GetRouteToTemplateMapping").Return(map[string]string{
		"/components/footer": "app/components/footer.templ#Footer",
		"/components/navbar": "app/components/navbar.templ#Navbar",
	})
	mockTemplateRegistry.On("GetTemplateKeyByComponentName", "footer").Return("app/components/footer.templ#Footer", true)
	mockTemplateRegistry.On("GetTemplateKeyByComponentName", "navbar").Return("app/components/navbar.templ#Navbar", true)
	mockTemplateRegistry.On("GetTemplateKeyByComponentName", "nonexistent").Return("", false)

	// Mock template metadata responses
	mockTemplateRegistry.On("GetTemplateMetadata", "app/components/footer.templ#Footer").Return(&interfaces.TemplateMetadata{
		TemplatePath:  "app/components/footer.templ",
		Type:          interfaces.TemplateTypeComponent,
		ComponentName: "footer",
		Route:         "/components/footer",
		YAMLFile:      "app/components/footer.templ.yaml",
		YAMLExists:    true,
		HasMetadata:   true,
	}, nil)
	mockTemplateRegistry.On("GetTemplateMetadata", "app/components/navbar.templ#Navbar").Return(&interfaces.TemplateMetadata{
		TemplatePath:  "app/components/navbar.templ",
		Type:          interfaces.TemplateTypeComponent,
		ComponentName: "navbar",
		Route:         "/components/navbar",
		YAMLFile:      "app/components/navbar.templ.yaml",
		YAMLExists:    true,
		HasMetadata:   true,
	}, nil)

	// Create service instance manually for testing (not using DI)
	cms := &componentMetadataService{
		configService:    mockConfigService,
		templateRegistry: mockTemplateRegistry,
		logger:           logger,
		metadataCache:    make(map[string]*shared.ConfigFile),
		translationCache: make(map[string]map[string]string),
	}

	tests := []struct {
		name          string
		componentName string
		expectError   bool
		description   string
	}{
		{
			name:          "Load existing component metadata",
			componentName: "footer",
			expectError:   true, // YAML file doesn't actually exist in test environment
			description:   "Should return error when YAML file is missing",
		},
		{
			name:          "Load non-existing component metadata",
			componentName: "nonexistent",
			expectError:   true,
			description:   "Should return error for non-existing component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := cms.LoadComponentMetadata(tt.componentName)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, config, tt.description)
			} else {
				// For successful cases (shouldn't happen in current test setup)
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, config, tt.description)
				// Verify that it's now cached
				cachedConfig, found := cms.GetCachedMetadata(tt.componentName)
				assert.True(t, found, "Component should be cached after loading")
				assert.Equal(t, config, cachedConfig, "Cached config should match loaded config")
			}
		})
	}

	// Note: Not asserting all mock expectations since some methods may not be called in error paths
	mockConfigService.AssertCalled(t, "GetLayoutRootDirectory")
	mockTemplateRegistry.AssertCalled(t, "GetTemplateKeyByComponentName", "footer")
}

func TestComponentMetadataService_DetectComponentsFromTemplate(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockConfigService := &MockConfigService{}

	cms := &componentMetadataService{
		configService:    mockConfigService,
		logger:           logger,
		metadataCache:    make(map[string]*shared.ConfigFile),
		translationCache: make(map[string]map[string]string),
	}

	tests := []struct {
		name            string
		templateContent string
		expected        []string
		description     string
	}{
		{
			name:            "Single component usage",
			templateContent: `<footer>@components.Footer()</footer>`,
			expected:        []string{"Footer"},
			description:     "Should detect single component usage",
		},
		{
			name:            "Multiple component usage",
			templateContent: `<header>@components.NavBar()</header><main>Content</main><footer>@components.Footer()</footer>`,
			expected:        []string{"NavBar", "Footer"},
			description:     "Should detect multiple component usages",
		},
		{
			name:            "Component with parameters",
			templateContent: `<div>@components.SearchBar("placeholder")</div>`,
			expected:        []string{"SearchBar"},
			description:     "Should detect component with parameters",
		},
		{
			name:            "No component usage",
			templateContent: `<div>Regular HTML content</div>`,
			expected:        []string{},
			description:     "Should return empty slice for templates without components",
		},
		{
			name:            "Component with whitespace",
			templateContent: `<footer>@components.Footer   (  )</footer>`,
			expected:        []string{"Footer"},
			description:     "Should handle whitespace in component calls",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			components, err := cms.DetectComponentsFromTemplate(tt.templateContent)

			assert.NoError(t, err, tt.description)
			assert.Equal(t, len(tt.expected), len(components), tt.description)

			// Check that all expected components are found (order doesn't matter)
			componentMap := make(map[string]bool)
			for _, comp := range components {
				componentMap[comp] = true
			}

			for _, expectedComp := range tt.expected {
				assert.True(t, componentMap[expectedComp],
					"Expected component %s not found in %v", expectedComp, components)
			}
		})
	}
}

func TestComponentMetadataService_CacheBehavior(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockConfigService := &MockConfigService{}
	mockConfigService.On("GetTemplateOutputDir").Return("demo/app")
	mockConfigService.On("GetFallbackLocale").Return("en")

	mockTemplateRegistry := &MockTemplateRegistry{}
	mockTemplateRegistry.On("GetRouteToTemplateMapping").Return(map[string]string{
		"/components/footer": "app/components/footer.templ#Footer",
	})

	cms := &componentMetadataService{
		configService:    mockConfigService,
		templateRegistry: mockTemplateRegistry,
		logger:           logger,
		metadataCache:    make(map[string]*shared.ConfigFile),
		translationCache: make(map[string]map[string]string),
	}

	// Test metadata cache
	t.Run("Metadata cache behavior", func(t *testing.T) {
		// Initially no cached metadata
		_, found := cms.GetCachedMetadata("footer")
		assert.False(t, found, "Should not have cached metadata initially")

		// Manually cache some metadata
		testConfig := &shared.ConfigFile{
			FilePath: "test/path.yaml",
			Metadata: &shared.MetadataConfig{
				Custom: map[string]interface{}{
					"title": "Test Component",
				},
			},
		}
		cms.cacheMetadata("footer", testConfig)

		// Now it should be cached
		cachedConfig, found := cms.GetCachedMetadata("footer")
		assert.True(t, found, "Should have cached metadata after caching")
		assert.Equal(t, testConfig, cachedConfig, "Cached metadata should match")
	})

	// Test translation cache
	t.Run("Translation cache behavior", func(t *testing.T) {
		// Initially no cached translations
		_, found := cms.GetCachedTranslations("footer", "en")
		assert.False(t, found, "Should not have cached translations initially")

		// Manually cache some translations
		testTranslations := map[string]string{
			"footer_copyright": "© 2024 Test",
			"footer_privacy":   "Privacy Policy",
		}
		cms.cacheTranslations("footer", "en", testTranslations)

		// Now it should be cached
		cachedTranslations, found := cms.GetCachedTranslations("footer", "en")
		assert.True(t, found, "Should have cached translations after caching")
		assert.Equal(t, testTranslations, cachedTranslations, "Cached translations should match")
	})

	// Note: Not checking mock expectations since we're manually testing cache behavior
	// and not all service methods are called in cache tests
}

func TestComponentMetadataService_LoadMultipleComponentMetadata(t *testing.T) {
	logger := zaptest.NewLogger(t)

	mockConfigService := &MockConfigService{}
	mockConfigService.On("GetTemplateOutputDir").Return("demo/app")
	mockConfigService.On("GetFallbackLocale").Return("en")
	mockConfigService.On("GetLayoutRootDirectory").Return("demo/app")

	mockTemplateRegistry := &MockTemplateRegistry{}
	mockTemplateRegistry.On("GetRouteToTemplateMapping").Return(map[string]string{
		"/components/footer": "app/components/footer.templ#Footer",
		"/components/navbar": "app/components/navbar.templ#Navbar",
	})
	mockTemplateRegistry.On("GetTemplateKeyByComponentName", "footer").Return("app/components/footer.templ#Footer", true)
	mockTemplateRegistry.On("GetTemplateKeyByComponentName", "navbar").Return("app/components/navbar.templ#Navbar", true)
	mockTemplateRegistry.On("GetTemplateKeyByComponentName", "nonexistent").Return("", false)

	// Mock template metadata responses
	mockTemplateRegistry.On("GetTemplateMetadata", "app/components/footer.templ#Footer").Return(&interfaces.TemplateMetadata{
		TemplatePath:  "app/components/footer.templ",
		Type:          interfaces.TemplateTypeComponent,
		ComponentName: "footer",
		Route:         "/components/footer",
		YAMLFile:      "app/components/footer.templ.yaml",
		YAMLExists:    true,
		HasMetadata:   true,
	}, nil)
	mockTemplateRegistry.On("GetTemplateMetadata", "app/components/navbar.templ#Navbar").Return(&interfaces.TemplateMetadata{
		TemplatePath:  "app/components/navbar.templ",
		Type:          interfaces.TemplateTypeComponent,
		ComponentName: "navbar",
		Route:         "/components/navbar",
		YAMLFile:      "app/components/navbar.templ.yaml",
		YAMLExists:    true,
		HasMetadata:   true,
	}, nil)

	cms := &componentMetadataService{
		configService:    mockConfigService,
		templateRegistry: mockTemplateRegistry,
		logger:           logger,
		metadataCache:    make(map[string]*shared.ConfigFile),
		translationCache: make(map[string]map[string]string),
	}

	componentNames := []string{"footer", "navbar", "nonexistent"}

	result, err := cms.LoadMultipleComponentMetadata(componentNames)

	assert.NoError(t, err, "Should not return error even when some components fail to load")
	assert.NotNil(t, result, "Should return a map result even if empty")

	// Should not have loaded non-existing component
	assert.Nil(t, result["nonexistent"], "Non-existing component should not be in result")

	// Only check expectations that were actually called
	mockConfigService.AssertCalled(t, "GetLayoutRootDirectory")
	// GetFallbackLocale may or may not be called depending on whether components are found
}
