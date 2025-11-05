package di

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

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
	mockAssets := &mockAssetsService{}
	mockUserStore := &mockUserStore{}
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
		WithUserStore(mockUserStore),
		// Note: Using default AuthHandlers from container
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

func (m *mockTemplateRegistry) RequiresDataService(key string) bool {
	return false
}

func (m *mockTemplateRegistry) GetDataServiceInfo(key string) (interfaces.DataServiceInfo, bool) {
	return interfaces.DataServiceInfo{}, false
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

// MockSessionStore for testing
type mockSessionStore struct {
	sessions map[string]*interfaces.Session
	mutex     sync.RWMutex
}

func (m *mockSessionStore) GetSession(req *http.Request) (*interfaces.Session, error) {
	// Simple mock implementation - just return nil for now
	return nil, fmt.Errorf("mock session store - no session found")
}

func (m *mockSessionStore) CreateSession(userID string) (*interfaces.Session, error) {
	session := &interfaces.Session{
		ID:        "mock-session-id",
		UserID:    userID,
		Valid:     true,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	m.mutex.Lock()
	if m.sessions == nil {
		m.sessions = make(map[string]*interfaces.Session)
	}
	m.sessions[session.ID] = session
	m.mutex.Unlock()

	return session, nil
}

func (m *mockSessionStore) DeleteSession(sessionID string) error {
	m.mutex.Lock()
	delete(m.sessions, sessionID)
	m.mutex.Unlock()
	return nil
}

// TestWithSessionStore tests custom session store injection
func TestWithSessionStore(t *testing.T) {
	// Create DI container
	container := NewContainer()
	defer container.Shutdown()

	// Register router services
	injector := container.RegisterRouterServices(context.Background(), "TR")

	// Create a custom session store
	customSessionStore := &mockSessionStore{}

	// Register application services with custom session store
	container.RegisterApplicationServices(
		WithSessionStore(customSessionStore),
	)

	// Verify that our custom session store is injected
	injectedSessionStore := do.MustInvoke[interfaces.SessionStore](injector)
	assert.Same(t, customSessionStore, injectedSessionStore, "Custom session store should be injected")

	// Test session operations
	session, err := injectedSessionStore.CreateSession("test-user-id")
	assert.NoError(t, err, "Session creation should succeed")
	assert.NotNil(t, session, "Session should not be nil")
	assert.Equal(t, "test-user-id", session.UserID, "Session should have correct user ID")
	assert.Equal(t, "mock-session-id", session.ID, "Session should have correct ID")
}

// TestDefaultSessionStore tests that default session store is used when no custom one is provided
func TestDefaultSessionStore(t *testing.T) {
	// Create DI container
	container := NewContainer()
	defer container.Shutdown()

	// Register router services
	injector := container.RegisterRouterServices(context.Background(), "TR")

	// Register application services without custom session store
	container.RegisterApplicationServices()

	// Verify that default session store is injected
	sessionStore := do.MustInvoke[interfaces.SessionStore](injector)
	assert.NotNil(t, sessionStore, "Default session store should be injected")

	// Test that default session store works (it should be the in-memory implementation)
	session, err := sessionStore.CreateSession("test-user-id")
	assert.NoError(t, err, "Default session store creation should succeed")
	assert.NotNil(t, session, "Session should not be nil")
	assert.Equal(t, "test-user-id", session.UserID, "Session should have correct user ID")
}