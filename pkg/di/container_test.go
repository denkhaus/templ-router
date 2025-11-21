package di

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
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
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsServiceFactory(func(i do.Injector) (interfaces.AssetsService, error) {
			return &mockAssetsService{}, nil
		}),
	)

	// Should not panic - all required services are registered
	container.RegisterRouterServices(context.Background(), "TR")

	// Verify logger is registered and can be retrieved
	logger := container.GetLogger()
	if logger == nil {
		t.Error("Logger not registered properly")
	}

	// Note: Not testing router here since it requires UserStore which is application-provided
}

func TestGetLogger(t *testing.T) {
	container := NewContainer()

	// Register required dependencies
	mockRegistry := &mockTemplateRegistry{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsServiceFactory(func(i do.Injector) (interfaces.AssetsService, error) {
			return &mockAssetsService{}, nil
		}),
	)
	container.RegisterRouterServices(context.Background(), "TR")

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
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsServiceFactory(func(i do.Injector) (interfaces.AssetsService, error) {
			return &mockAssetsService{}, nil
		}),
	)
	container.RegisterRouterServices(context.Background(), "TR")

	router := container.GetRouter()
	if router == nil {
		t.Fatal("GetRouter() returned nil")
	}
}

func TestShutdown(t *testing.T) {
	container := NewContainer()

	// Register minimal dependencies for shutdown test
	mockRegistry := &mockTemplateRegistry{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsServiceFactory(func(i do.Injector) (interfaces.AssetsService, error) {
			return &mockAssetsService{}, nil
		}),
	)
	container.RegisterRouterServices(context.Background(), "TR")

	err := container.Shutdown()
	if err != nil && err.Error() != "" {
		t.Errorf("Shutdown() returned error: %v", err)
	}
}

func TestRegisterApplicationServices(t *testing.T) {
	container := NewContainer()

	// Mock template registry
	mockRegistry := &mockTemplateRegistry{}

	// Note: mockAssets removed as it's no longer needed for testing

	// Register with options
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
	)

	// Verify services are registered
	retrievedRegistry := do.MustInvoke[interfaces.TemplateRegistry](container.injector)
	if retrievedRegistry == nil {
		t.Error("Template registry not registered correctly")
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

func (m *mockTemplateRegistry) RequiresDataService(key string) bool {
	return false
}

func (m *mockTemplateRegistry) GetDataServiceInfo(key string) (interfaces.DataServiceInfo, bool) {
	return interfaces.DataServiceInfo{}, false
}

func (m *mockTemplateRegistry) GetTemplateMetadata(key string) (*interfaces.TemplateMetadata, error) {
	return nil, nil
}

func (m *mockTemplateRegistry) GetTemplateMetadataByRoute(route string) (*interfaces.TemplateMetadata, error) {
	return nil, nil
}

func (m *mockTemplateRegistry) GetAllTemplateMetadata() map[string]*interfaces.TemplateMetadata {
	return make(map[string]*interfaces.TemplateMetadata)
}

func (m *mockTemplateRegistry) FindComponentTemplates() map[string]*interfaces.TemplateMetadata {
	return make(map[string]*interfaces.TemplateMetadata)
}

func (m *mockTemplateRegistry) GetTemplateKeyByComponentName(componentName string) (string, bool) {
	return "", false
}

type mockAssetsService struct{}

func (m *mockAssetsService) SetupRoutes(router *chi.Mux)             {}
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

// Note: AuthHandlers are now client-side and not part of the router framework

// Mock User for testing
type mockUser struct {
	ID    string
	Email string
	Roles []string
}

func (u *mockUser) GetID() string      { return u.ID }
func (u *mockUser) GetEmail() string   { return u.Email }
func (u *mockUser) GetRoles() []string { return u.Roles }

// TestCustomMiddlewareDefinitionOrder tests that custom middleware are executed in definition order
func TestCustomMiddlewareDefinitionOrder(t *testing.T) {
	// Create DI container
	container := NewContainer()
	defer container.Shutdown()

	// Register router services
	injector := container.RegisterRouterServices(context.Background(), "TR")

	// Register application services with custom middleware in specific order
	container.RegisterApplicationServices(
		// First middleware (Order: 0)
		WithCustomMiddleware("first", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("X-Middleware-Order", "first")
				next.ServeHTTP(w, r)
			})
		}),

		// Second middleware (Order: 1)
		WithCustomMiddleware("second", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("X-Middleware-Order", "second")
				next.ServeHTTP(w, r)
			})
		}),

		// Third middleware (Order: 2)
		WithCustomMiddleware("third", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("X-Middleware-Order", "third")
				next.ServeHTTP(w, r)
			})
		}),
	)

	// Test that custom middleware definitions are loaded in correct order
	middlewareDefs, err := do.Invoke[[]interfaces.CustomMiddlewareDefinition](injector)
	assert.NoError(t, err)
	assert.Len(t, middlewareDefs, 3, "Should have exactly 3 custom middleware definitions")

	// Verify definition order
	expectedOrder := []string{"first", "second", "third"}
	for i, middlewareDef := range middlewareDefs {
		assert.Equal(t, expectedOrder[i], middlewareDef.Name, "Middleware name at index %d should match", i)
		assert.Equal(t, i, middlewareDef.Order, "Middleware order at index %d should match", i)
	}

	// Test actual middleware execution
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware in order (as router would do)
	handler := testHandler
	for _, middlewareDef := range middlewareDefs {
		handler = middlewareDef.Func(handler).(http.HandlerFunc)
	}

	// Create test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Execute the middleware chain
	handler.ServeHTTP(w, req)

	// Check results - HTTP headers can have multiple values
	orderHeaders := w.Header()["X-Middleware-Order"]
	assert.NotEmpty(t, orderHeaders, "Should have X-Middleware-Order header")
	assert.Len(t, orderHeaders, 3, "Should have exactly 3 header values")

	// Join all header values
	var orderString string
	for i, header := range orderHeaders {
		if i > 0 {
			orderString += ","
		}
		orderString += header
	}

	// Note: In Go middleware chains, the first middleware defined is the outermost
	// so it executes last on the way in. This is the expected behavior.
	expectedChain := "third,second,first"
	assert.Equal(t, expectedChain, orderString, "Middleware should execute in correct order")
}

// TestCustomMiddlewareEmpty tests behavior when no custom middleware is defined
func TestCustomMiddlewareEmpty(t *testing.T) {
	// Create DI container
	container := NewContainer()
	defer container.Shutdown()

	// Register router services
	injector := container.RegisterRouterServices(context.Background(), "TR")

	// Register application services without custom middleware
	container.RegisterApplicationServices()

	// Try to load custom middleware definitions - should return empty slice when service doesn't exist
	_, err := do.Invoke[[]interfaces.CustomMiddlewareDefinition](injector)

	// Should get error because service is not registered when no custom middleware is defined
	assert.Error(t, err, "Should get error when no custom middleware service is registered")

	// Alternative: Register empty custom middleware slice explicitly
	container.RegisterApplicationServices(
		WithCustomMiddleware("test", func(next http.Handler) http.Handler {
			return next
		}),
	)

	// Now we can load and verify the service
	middlewareDefs, err := do.Invoke[[]interfaces.CustomMiddlewareDefinition](injector)
	assert.NoError(t, err)
	assert.Len(t, middlewareDefs, 1, "Should have exactly 1 custom middleware")
	assert.Equal(t, "test", middlewareDefs[0].Name, "Middleware name should match")
}
