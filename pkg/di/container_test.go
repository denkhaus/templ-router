package di

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	
	if container == nil {
		t.Fatal("NewContainer() returned nil")
	}
	
	if container.injector == nil {
		t.Fatal("Container injector is nil")
	}
}

func TestGetInjector(t *testing.T) {
	container := NewContainer()
	injector := container.GetInjector()
	
	if injector == nil {
		t.Fatal("GetInjector() returned nil")
	}
	
	// Verify it's the same injector
	if injector != container.injector {
		t.Error("GetInjector() returned different injector")
	}
}

func TestRegisterRouterServices(t *testing.T) {
	container := NewContainer()
	
	// Register required dependencies
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
	)
	
	// Should not panic
	container.RegisterRouterServices()
	
	// Verify logger is registered and can be retrieved
	logger := container.GetLogger()
	if logger == nil {
		t.Error("Logger not registered properly")
	}
	
	// Verify router is registered and can be retrieved
	router := container.GetRouter()
	if router == nil {
		t.Error("Router not registered properly")
	}
}

func TestGetLogger(t *testing.T) {
	container := NewContainer()
	
	// Register required dependencies
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
	)
	container.RegisterRouterServices()
	
	logger := container.GetLogger()
	if logger == nil {
		t.Fatal("GetLogger() returned nil")
	}
	
	// Verify it's a zap logger
	if _, ok := interface{}(logger).(*zap.Logger); !ok {
		t.Error("GetLogger() did not return *zap.Logger")
	}
}

func TestGetRouter(t *testing.T) {
	container := NewContainer()
	
	// Register required dependencies
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
	)
	container.RegisterRouterServices()
	
	router := container.GetRouter()
	if router == nil {
		t.Fatal("GetRouter() returned nil")
	}
}

func TestShutdown(t *testing.T) {
	container := NewContainer()
	container.RegisterRouterServices()
	
	err := container.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() returned error: %v", err)
	}
}

func TestRegisterApplicationServices(t *testing.T) {
	container := NewContainer()
	
	// Mock template registry
	mockRegistry := &mockTemplateRegistry{}
	
	// Mock assets service  
	mockAssets := &mockAssetsService{}
	
	// Register with options
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
	)
	
	// Verify services are registered
	retrievedRegistry := do.MustInvoke[interfaces.TemplateRegistry](container.injector)
	if retrievedRegistry == nil {
		t.Error("Template registry not registered correctly")
	}
	
	retrievedAssets := do.MustInvoke[interfaces.AssetsService](container.injector)
	if retrievedAssets == nil {
		t.Error("Assets service not registered correctly")
	}
}

// Mock implementations for testing
type mockTemplateRegistry struct{}

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

func (m *mockTemplateRegistry) GetRouteToTemplateMapping() map[string]string {
	return map[string]string{}
}

func (m *mockTemplateRegistry) GetTemplateByRoute(route string) (templ.Component, error) {
	return nil, nil
}

type mockAssetsService struct{}

func (m *mockAssetsService) SetupRoutes(router *chi.Mux) {}
func (m *mockAssetsService) SetupRoutesWithRouter(router chi.Router) {}