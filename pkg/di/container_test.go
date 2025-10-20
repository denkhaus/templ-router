package di

import (
	"net/http"
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
	
	// Register required dependencies (AuthHandlers are provided by default)
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	mockUserStore := &mockUserStore{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
		WithUserStore(mockUserStore),
		// Note: No AuthHandlers needed - provided by default
	)
	
	// Should not panic - default AuthHandlers are automatically available
	container.RegisterRouterServices("TR")
	
	// Verify logger is registered and can be retrieved
	logger := container.GetLogger()
	if logger == nil {
		t.Error("Logger not registered properly")
	}
	
	// Note: Not testing router here since it requires UserStore which is application-provided
}

func TestGetLogger(t *testing.T) {
	container := NewContainer()
	
	// Register required dependencies (without custom AuthHandlers - use default)
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	mockUserStore := &mockUserStore{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
		WithUserStore(mockUserStore),
		// Note: Using default AuthHandlers from container
	)
	container.RegisterRouterServices("TR")
	
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
	
	// Register required dependencies including mock UserStore
	// This simulates how an application would provide its own UserStore implementation
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	mockUserStore := &mockUserStore{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
		WithUserStore(mockUserStore),
		// Note: Using default AuthHandlers from container
	)
	container.RegisterRouterServices("TR")
	
	router := container.GetRouter()
	if router == nil {
		t.Fatal("GetRouter() returned nil")
	}
}

func TestShutdown(t *testing.T) {
	container := NewContainer()
	
	// Register minimal dependencies for shutdown test
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	mockUserStore := &mockUserStore{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
		WithUserStore(mockUserStore),
		// Note: Using default AuthHandlers from container
	)
	container.RegisterRouterServices("TR")
	
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

// Mock UserStore for testing
type mockUserStore struct{}

func (m *mockUserStore) GetUserByID(userID string) (interfaces.UserEntity, error) {
	return &mockUser{ID: userID, Email: "test@example.com", Roles: []string{"user"}}, nil
}

func (m *mockUserStore) GetUserByEmail(email string) (interfaces.UserEntity, error) {
	return &mockUser{ID: "test123", Email: email, Roles: []string{"user"}}, nil
}

func (m *mockUserStore) ValidateCredentials(email, password string) (interfaces.UserEntity, error) {
	return &mockUser{ID: "test123", Email: email, Roles: []string{"user"}}, nil
}

func (m *mockUserStore) CreateUser(username, email, password string) (interfaces.UserEntity, error) {
	return &mockUser{ID: "new123", Email: email, Roles: []string{"user"}}, nil
}

func (m *mockUserStore) UserExists(username, email string) (bool, error) {
	return false, nil
}

func (m *mockUserStore) ValidateCredentialsFromRequest(req *http.Request) (interfaces.UserEntity, error) {
	return &mockUser{ID: "test123", Email: "test@example.com", Roles: []string{"user"}}, nil
}

func (m *mockUserStore) CreateUserFromRequest(req *http.Request) (interfaces.UserEntity, error) {
	return &mockUser{ID: "new123", Email: "new@example.com", Roles: []string{"user"}}, nil
}

// Mock AuthHandlers for testing
type mockAuthHandlers struct{}

func (m *mockAuthHandlers) RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc)) {}
func (m *mockAuthHandlers) HandleLogin(w http.ResponseWriter, r *http.Request) {}
func (m *mockAuthHandlers) HandleLogout(w http.ResponseWriter, r *http.Request) {}
func (m *mockAuthHandlers) HandleSignup(w http.ResponseWriter, r *http.Request) {}

// Mock User for testing
type mockUser struct {
	ID    string
	Email string
	Roles []string
}

func (u *mockUser) GetID() string    { return u.ID }
func (u *mockUser) GetEmail() string { return u.Email }
func (u *mockUser) GetRoles() []string { return u.Roles }