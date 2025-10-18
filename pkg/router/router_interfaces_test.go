package router

import (
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
)

// Test implementations to verify interface compliance

type testRouterCore struct {
	routes          []Route
	layouts         []LayoutTemplate
	errorTemplates  []ErrorTemplate
	middlewareSetup MiddlewareSetup
	handlerBuilder  HandlerBuilder
	routeRegistrar  RouteRegistrar
	initialized     bool
}

func (t *testRouterCore) Initialize() error {
	t.initialized = true
	return nil
}

func (t *testRouterCore) RegisterRoutes(chiRouter *chi.Mux) error {
	return nil
}

func (t *testRouterCore) GetRoutes() []Route {
	return t.routes
}

func (t *testRouterCore) GetLayoutTemplates() []LayoutTemplate {
	return t.layouts
}

func (t *testRouterCore) GetErrorTemplates() []ErrorTemplate {
	return t.errorTemplates
}

func (t *testRouterCore) GetMiddlewareSetup() MiddlewareSetup {
	return t.middlewareSetup
}

func (t *testRouterCore) GetHandlerBuilder() HandlerBuilder {
	return t.handlerBuilder
}

func (t *testRouterCore) GetRouteRegistrar() RouteRegistrar {
	return t.routeRegistrar
}

type testRouteDiscovery struct {
	routes         []Route
	layouts        []LayoutTemplate
	errorTemplates []ErrorTemplate
}

func (t *testRouteDiscovery) DiscoverRoutes(scanPath string) ([]Route, error) {
	return t.routes, nil
}

func (t *testRouteDiscovery) DiscoverLayouts(scanPath string) ([]LayoutTemplate, error) {
	return t.layouts, nil
}

func (t *testRouteDiscovery) DiscoverErrorTemplates(scanPath string) ([]ErrorTemplate, error) {
	return t.errorTemplates, nil
}

type testConfigLoader struct {
	config       *interfaces.ConfigFile
	authSettings *interfaces.AuthSettings
}

func (t *testConfigLoader) LoadRouteConfig(templateFile string) (*interfaces.ConfigFile, error) {
	return t.config, nil
}

func (t *testConfigLoader) LoadConfig(templatePath string) (*interfaces.ConfigFile, error) {
	return t.config, nil
}

func (t *testConfigLoader) LoadAuthSettings(templatePath string) (*interfaces.AuthSettings, error) {
	return t.authSettings, nil
}

func TestRouterCore_InterfaceCompliance(t *testing.T) {
	// Verify that our test implementation satisfies the RouterCore interface
	var _ RouterCore = (*testRouterCore)(nil)

	router := &testRouterCore{
		routes: []Route{
			{Path: "/", TemplateFile: "index.templ"},
			{Path: "/about", TemplateFile: "about.templ"},
		},
		layouts: []LayoutTemplate{
			{FilePath: "/app/layout.templ", DirectoryPath: "/app"},
		},
		errorTemplates: []ErrorTemplate{
			{FilePath: "/app/error.templ", DirectoryPath: "/app", ErrorTypes: []string{"404"}},
		},
	}

	// Test Initialize
	err := router.Initialize()
	if err != nil {
		t.Errorf("Initialize() error = %v", err)
	}
	if !router.initialized {
		t.Error("Initialize() should set initialized flag")
	}

	// Test RegisterRoutes
	mux := chi.NewMux()
	err = router.RegisterRoutes(mux)
	if err != nil {
		t.Errorf("RegisterRoutes() error = %v", err)
	}

	// Test getter methods
	routes := router.GetRoutes()
	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}

	layouts := router.GetLayoutTemplates()
	if len(layouts) != 1 {
		t.Errorf("Expected 1 layout, got %d", len(layouts))
	}

	errorTemplates := router.GetErrorTemplates()
	if len(errorTemplates) != 1 {
		t.Errorf("Expected 1 error template, got %d", len(errorTemplates))
	}

	// Test that getters return expected values
	if routes[0].Path != "/" {
		t.Errorf("Expected first route path '/', got %s", routes[0].Path)
	}

	if layouts[0].DirectoryPath != "/app" {
		t.Errorf("Expected layout directory '/app', got %s", layouts[0].DirectoryPath)
	}

	if errorTemplates[0].DirectoryPath != "/app" {
		t.Errorf("Expected error template directory '/app', got %s", errorTemplates[0].DirectoryPath)
	}
}

func TestRouteDiscovery_InterfaceCompliance(t *testing.T) {
	// Verify that our test implementation satisfies the RouteDiscovery interface
	var _ RouteDiscovery = (*testRouteDiscovery)(nil)

	discovery := &testRouteDiscovery{
		routes: []Route{
			{Path: "/test", TemplateFile: "test.templ"},
		},
		layouts: []LayoutTemplate{
			{FilePath: "/app/test-layout.templ"},
		},
		errorTemplates: []ErrorTemplate{
			{FilePath: "/app/test-error.templ"},
		},
	}

	// Test DiscoverRoutes
	routes, err := discovery.DiscoverRoutes("/app")
	if err != nil {
		t.Errorf("DiscoverRoutes() error = %v", err)
	}
	if len(routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(routes))
	}
	if routes[0].Path != "/test" {
		t.Errorf("Expected route path '/test', got %s", routes[0].Path)
	}

	// Test DiscoverLayouts
	layouts, err := discovery.DiscoverLayouts("/app")
	if err != nil {
		t.Errorf("DiscoverLayouts() error = %v", err)
	}
	if len(layouts) != 1 {
		t.Errorf("Expected 1 layout, got %d", len(layouts))
	}
	if layouts[0].FilePath != "/app/test-layout.templ" {
		t.Errorf("Expected layout path '/app/test-layout.templ', got %s", layouts[0].FilePath)
	}

	// Test DiscoverErrorTemplates
	errorTemplates, err := discovery.DiscoverErrorTemplates("/app")
	if err != nil {
		t.Errorf("DiscoverErrorTemplates() error = %v", err)
	}
	if len(errorTemplates) != 1 {
		t.Errorf("Expected 1 error template, got %d", len(errorTemplates))
	}
	if errorTemplates[0].FilePath != "/app/test-error.templ" {
		t.Errorf("Expected error template path '/app/test-error.templ', got %s", errorTemplates[0].FilePath)
	}
}

func TestConfigLoader_InterfaceCompliance(t *testing.T) {
	// Verify that our test implementation satisfies the ConfigLoader interface
	var _ ConfigLoader = (*testConfigLoader)(nil)

	loader := &testConfigLoader{
		config: &interfaces.ConfigFile{
			FilePath: "/app/test.yaml",
			I18nMappings: map[string]string{
				"en": "English",
				"de": "German",
			},
		},
		authSettings: &interfaces.AuthSettings{
			Type:        interfaces.AuthTypeUser,
			RedirectURL: "/login",
			Roles:       []string{"user"},
		},
	}

	// Test LoadRouteConfig
	config, err := loader.LoadRouteConfig("test.templ")
	if err != nil {
		t.Errorf("LoadRouteConfig() error = %v", err)
	}
	if config == nil {
		t.Error("LoadRouteConfig() returned nil config")
	}
	if config.FilePath != "/app/test.yaml" {
		t.Errorf("Expected config path '/app/test.yaml', got %s", config.FilePath)
	}

	// Test LoadConfig
	config2, err := loader.LoadConfig("/app/test.templ")
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
	}
	if config2 == nil {
		t.Error("LoadConfig() returned nil config")
	}

	// Test LoadAuthSettings
	authSettings, err := loader.LoadAuthSettings("/app/test.templ")
	if err != nil {
		t.Errorf("LoadAuthSettings() error = %v", err)
	}
	if authSettings == nil {
		t.Error("LoadAuthSettings() returned nil auth settings")
	}
	if authSettings.Type != interfaces.AuthTypeUser {
		t.Errorf("Expected auth type User, got %v", authSettings.Type)
	}
	if authSettings.RedirectURL != "/login" {
		t.Errorf("Expected redirect URL '/login', got %s", authSettings.RedirectURL)
	}
	if len(authSettings.Roles) != 1 || authSettings.Roles[0] != "user" {
		t.Errorf("Expected roles ['user'], got %v", authSettings.Roles)
	}
}

func TestRouterInterfaces_MethodSignatures(t *testing.T) {
	// This test ensures that interface method signatures are correct
	// by attempting to call them with the expected parameters

	// RouterCore interface methods
	var router RouterCore = &testRouterCore{}
	
	_ = router.Initialize()
	_ = router.RegisterRoutes(chi.NewMux())
	_ = router.GetRoutes()
	_ = router.GetLayoutTemplates()
	_ = router.GetErrorTemplates()
	_ = router.GetMiddlewareSetup()
	_ = router.GetHandlerBuilder()
	_ = router.GetRouteRegistrar()

	// RouteDiscovery interface methods
	var discovery RouteDiscovery = &testRouteDiscovery{}
	
	_, _ = discovery.DiscoverRoutes("/app")
	_, _ = discovery.DiscoverLayouts("/app")
	_, _ = discovery.DiscoverErrorTemplates("/app")

	// ConfigLoader interface methods
	var loader ConfigLoader = &testConfigLoader{}
	
	_, _ = loader.LoadRouteConfig("test.templ")
	_, _ = loader.LoadConfig("/app/test.templ")
	_, _ = loader.LoadAuthSettings("/app/test.templ")
}

func TestRouterInterfaces_ReturnTypes(t *testing.T) {
	// Test that interface methods return the expected types

	router := &testRouterCore{
		routes: []Route{{Path: "/test"}},
		layouts: []LayoutTemplate{{FilePath: "/test"}},
		errorTemplates: []ErrorTemplate{{FilePath: "/test"}},
	}

	// Test return types
	routes := router.GetRoutes()
	if routes == nil {
		t.Error("GetRoutes() should not return nil")
	}

	layouts := router.GetLayoutTemplates()
	if layouts == nil {
		t.Error("GetLayoutTemplates() should not return nil")
	}

	errorTemplates := router.GetErrorTemplates()
	if errorTemplates == nil {
		t.Error("GetErrorTemplates() should not return nil")
	}

	// Test that slices have expected length
	if len(routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(routes))
	}

	if len(layouts) != 1 {
		t.Errorf("Expected 1 layout, got %d", len(layouts))
	}

	if len(errorTemplates) != 1 {
		t.Errorf("Expected 1 error template, got %d", len(errorTemplates))
	}
}