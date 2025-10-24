package cache

import (
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewCacheService(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	// Provide logger
	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	service, err := NewCacheService(injector)
	require.NoError(t, err)
	require.NotNil(t, service)
}

func TestCacheService_TemplateOperations(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	service, err := NewCacheService(injector)
	require.NoError(t, err)

	// Test cache miss
	value, exists := service.GetTemplate("nonexistent")
	assert.False(t, exists)
	assert.Nil(t, value)

	// Test cache set and hit
	testValue := "test template content"
	service.SetTemplate("test-key", testValue)

	value, exists = service.GetTemplate("test-key")
	assert.True(t, exists)
	assert.Equal(t, testValue, value)

	// Test cache clear
	service.ClearTemplates()
	value, exists = service.GetTemplate("test-key")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestCacheService_RouteOperations(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	service, err := NewCacheService(injector)
	require.NoError(t, err)

	// Test cache miss
	value, exists := service.GetRoute("nonexistent")
	assert.False(t, exists)
	assert.Nil(t, value)

	// Test cache set and hit
	testRoute := interfaces.Route{Path: "/test", TemplateFile: "test.templ"}
	service.SetRoute("test-route", testRoute)

	value, exists = service.GetRoute("test-route")
	assert.True(t, exists)
	assert.Equal(t, testRoute, value)

	// Test cache clear
	service.ClearRoutes()
	value, exists = service.GetRoute("test-route")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestCacheService_ClearAll(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	service, err := NewCacheService(injector)
	require.NoError(t, err)

	// Add items to both caches
	service.SetTemplate("template-key", "template-value")
	service.SetRoute("route-key", "route-value")

	// Verify items exist
	_, exists := service.GetTemplate("template-key")
	assert.True(t, exists)
	_, exists = service.GetRoute("route-key")
	assert.True(t, exists)

	// Clear all
	service.ClearAll()

	// Verify items are gone
	_, exists = service.GetTemplate("template-key")
	assert.False(t, exists)
	_, exists = service.GetRoute("route-key")
	assert.False(t, exists)
}

func TestCacheService_BuildTemplateKey(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	service, err := NewCacheService(injector)
	require.NoError(t, err)

	tests := []struct {
		name         string
		templateFile string
		locale       string
		params       map[string]string
		expected     string
	}{
		{
			name:         "no parameters",
			templateFile: "user/profile.templ",
			locale:       "en",
			params:       nil,
			expected:     "template:user/profile.templ:en",
		},
		{
			name:         "single parameter",
			templateFile: "user/profile.templ",
			locale:       "en",
			params:       map[string]string{"id": "123"},
			expected:     "template:user/profile.templ:en:id=123",
		},
		{
			name:         "multiple parameters sorted",
			templateFile: "user/profile.templ",
			locale:       "de",
			params:       map[string]string{"id": "123", "action": "edit"},
			expected:     "template:user/profile.templ:de:action=edit&id=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.BuildTemplateKey(tt.templateFile, tt.locale, tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCacheService_BuildRouteKey(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	service, err := NewCacheService(injector)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		params   map[string]string
		expected string
	}{
		{
			name:     "no parameters",
			path:     "/user/profile",
			params:   nil,
			expected: "route:/user/profile",
		},
		{
			name:     "single parameter",
			path:     "/user/profile",
			params:   map[string]string{"id": "123"},
			expected: "route:/user/profile:id=123",
		},
		{
			name:     "multiple parameters sorted",
			path:     "/user/profile",
			params:   map[string]string{"id": "123", "tab": "settings"},
			expected: "route:/user/profile:id=123&tab=settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.BuildRouteKey(tt.path, tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCacheService_KeyConsistency(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	service, err := NewCacheService(injector)
	require.NoError(t, err)

	// Test that parameter order doesn't affect key generation
	params1 := map[string]string{"a": "1", "b": "2", "c": "3"}
	params2 := map[string]string{"c": "3", "a": "1", "b": "2"}

	key1 := service.BuildTemplateKey("test.templ", "en", params1)
	key2 := service.BuildTemplateKey("test.templ", "en", params2)

	assert.Equal(t, key1, key2, "Keys should be identical regardless of parameter order")

	routeKey1 := service.BuildRouteKey("/test", params1)
	routeKey2 := service.BuildRouteKey("/test", params2)

	assert.Equal(t, routeKey1, routeKey2, "Route keys should be identical regardless of parameter order")
}