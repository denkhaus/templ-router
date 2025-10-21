package interfaces

import (
	"errors"
	"testing"

	"github.com/a-h/templ"
)

// MockTemplateRegistry implements TemplateRegistry for testing
type MockTemplateRegistry struct {
	templates    map[string]interface{}
	routeMapping map[string]string
	shouldError  bool
}

func NewMockTemplateRegistry() *MockTemplateRegistry {
	return &MockTemplateRegistry{
		templates: map[string]interface{}{
			"test-key-1": func() templ.Component { return templ.Raw("test content 1") },
			"test-key-2": func(param string) templ.Component { return templ.Raw("test content 2: " + param) },
			"test-key-3": func(data interface{}) templ.Component { return templ.Raw("test content 3") },
		},
		routeMapping: map[string]string{
			"/":         "test-key-1",
			"/about":    "test-key-2",
			"/contact":  "test-key-3",
		},
		shouldError: false,
	}
}

func (m *MockTemplateRegistry) GetTemplate(key string) (templ.Component, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}
	
	if templateFunc, exists := m.templates[key]; exists {
		switch fn := templateFunc.(type) {
		case func() templ.Component:
			return fn(), nil
		case func(string) templ.Component:
			return fn("test-param"), nil
		case func(interface{}) templ.Component:
			return fn(nil), nil
		default:
			return nil, errors.New("unsupported template function type")
		}
	}
	return nil, errors.New("template not found: " + key)
}

func (m *MockTemplateRegistry) GetTemplateFunction(key string) (func() interface{}, bool) {
	if templateFunc, exists := m.templates[key]; exists {
		return func() interface{} {
			return templateFunc
		}, true
	}
	return nil, false
}

func (m *MockTemplateRegistry) GetAllTemplateKeys() []string {
	keys := make([]string, 0, len(m.templates))
	for key := range m.templates {
		keys = append(keys, key)
	}
	return keys
}

func (m *MockTemplateRegistry) IsAvailable(key string) bool {
	_, exists := m.templates[key]
	return exists
}

func (m *MockTemplateRegistry) GetRouteToTemplateMapping() map[string]string {
	return m.routeMapping
}

func (m *MockTemplateRegistry) GetTemplateByRoute(route string) (templ.Component, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}
	
	if templateKey, exists := m.routeMapping[route]; exists {
		return m.GetTemplate(templateKey)
	}
	return nil, errors.New("no template found for route: " + route)
}

func (m *MockTemplateRegistry) SetShouldError(shouldError bool) {
	m.shouldError = shouldError
}

// RequiresDataService checks if a template requires a data service
func (m *MockTemplateRegistry) RequiresDataService(key string) bool {
	// For testing, return false by default
	return false
}

// GetDataServiceInfo retrieves data service information for a template
func (m *MockTemplateRegistry) GetDataServiceInfo(key string) (DataServiceInfo, bool) {
	// For testing, return empty info
	return DataServiceInfo{}, false
}

// Tests for TemplateRegistry interface
func TestTemplateRegistry_GetTemplate(t *testing.T) {
	registry := NewMockTemplateRegistry()

	tests := []struct {
		name        string
		key         string
		shouldError bool
		expectError bool
	}{
		{"Valid key - simple function", "test-key-1", false, false},
		{"Valid key - parameterized function", "test-key-2", false, false},
		{"Valid key - complex function", "test-key-3", false, false},
		{"Invalid key", "nonexistent-key", false, true},
		{"Error condition", "test-key-1", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry.SetShouldError(tt.shouldError)
			
			component, err := registry.GetTemplate(tt.key)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if component != nil {
					t.Error("Expected nil component on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if component == nil {
					t.Error("Expected component but got nil")
				}
			}
		})
	}
}

func TestTemplateRegistry_GetTemplateFunction(t *testing.T) {
	registry := NewMockTemplateRegistry()

	tests := []struct {
		name   string
		key    string
		exists bool
	}{
		{"Valid key", "test-key-1", true},
		{"Invalid key", "nonexistent-key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, exists := registry.GetTemplateFunction(tt.key)
			
			if exists != tt.exists {
				t.Errorf("Expected exists=%v, got %v", tt.exists, exists)
			}
			
			if tt.exists {
				if fn == nil {
					t.Error("Expected function but got nil")
				}
				// Test that the function returns something
				result := fn()
				if result == nil {
					t.Error("Function should return non-nil value")
				}
			} else {
				if fn != nil {
					t.Error("Expected nil function for non-existent key")
				}
			}
		})
	}
}

func TestTemplateRegistry_GetAllTemplateKeys(t *testing.T) {
	registry := NewMockTemplateRegistry()
	
	keys := registry.GetAllTemplateKeys()
	
	expectedCount := 3
	if len(keys) != expectedCount {
		t.Errorf("Expected %d keys, got %d", expectedCount, len(keys))
	}
	
	expectedKeys := map[string]bool{
		"test-key-1": true,
		"test-key-2": true,
		"test-key-3": true,
	}
	
	for _, key := range keys {
		if !expectedKeys[key] {
			t.Errorf("Unexpected key: %s", key)
		}
	}
}

func TestTemplateRegistry_IsAvailable(t *testing.T) {
	registry := NewMockTemplateRegistry()

	tests := []struct {
		name      string
		key       string
		available bool
	}{
		{"Available key", "test-key-1", true},
		{"Unavailable key", "nonexistent-key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available := registry.IsAvailable(tt.key)
			if available != tt.available {
				t.Errorf("Expected available=%v, got %v", tt.available, available)
			}
		})
	}
}

func TestTemplateRegistry_GetRouteToTemplateMapping(t *testing.T) {
	registry := NewMockTemplateRegistry()
	
	mapping := registry.GetRouteToTemplateMapping()
	
	expectedRoutes := []string{"/", "/about", "/contact"}
	if len(mapping) != len(expectedRoutes) {
		t.Errorf("Expected %d routes, got %d", len(expectedRoutes), len(mapping))
	}
	
	for _, route := range expectedRoutes {
		if _, exists := mapping[route]; !exists {
			t.Errorf("Expected route %s not found in mapping", route)
		}
	}
}

func TestTemplateRegistry_GetTemplateByRoute(t *testing.T) {
	registry := NewMockTemplateRegistry()

	tests := []struct {
		name        string
		route       string
		shouldError bool
		expectError bool
	}{
		{"Valid route", "/", false, false},
		{"Valid route with params", "/about", false, false},
		{"Invalid route", "/nonexistent", false, true},
		{"Error condition", "/", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry.SetShouldError(tt.shouldError)
			
			component, err := registry.GetTemplateByRoute(tt.route)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if component != nil {
					t.Error("Expected nil component on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if component == nil {
					t.Error("Expected component but got nil")
				}
			}
		})
	}
}

// Test interface compliance
func TestTemplateRegistryInterfaceCompliance(t *testing.T) {
	var _ TemplateRegistry = (*MockTemplateRegistry)(nil)
	
	// This test ensures our mock implements the interface correctly
	registry := NewMockTemplateRegistry()
	
	// Test all interface methods are callable
	_ = registry.GetAllTemplateKeys()
	_ = registry.IsAvailable("test")
	_ = registry.GetRouteToTemplateMapping()
	
	_, _ = registry.GetTemplate("test")
	_, _ = registry.GetTemplateFunction("test")
	_, _ = registry.GetTemplateByRoute("/test")
}