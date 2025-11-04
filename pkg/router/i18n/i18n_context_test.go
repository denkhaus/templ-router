package i18n

import (
	"context"
	"testing"

	"github.com/denkhaus/templ-router/pkg/shared"
)

func TestGetCurrentRoute(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() context.Context
		expected string
	}{
		{
			name: "Request path in context",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/en/dashboard")
				return ctx
			},
			expected: "/en/dashboard",
		},
		{
			name: "Root path",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/")
				return ctx
			},
			expected: "/",
		},
		{
			name: "Complex path with parameters",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/de/user/123/profile")
				return ctx
			},
			expected: "/de/user/123/profile",
		},
		{
			name: "No request path in context",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expected: "",
		},
		{
			name: "Request path is not a string",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, 123)
				return ctx
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			result := GetCurrentRoute(ctx)
			if result != tt.expected {
				t.Errorf("GetCurrentRoute() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPathMatchesPattern(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		pattern  string
		expected bool
	}{
		{
			name:     "Exact match",
			path:     "/dashboard",
			pattern:  "/dashboard",
			expected: true,
		},
		{
			name:     "Path with locale",
			path:     "/en/dashboard",
			pattern:  "/*/dashboard",
			expected: true,
		},
		{
			name:     "Path with multiple segments",
			path:     "/en/user/123/profile",
			pattern:  "/*/user/*/profile",
			expected: true,
		},
		{
			name:     "Different number of segments",
			path:     "/en/dashboard",
			pattern:  "/*/dashboard/profile",
			expected: false,
		},
		{
			name:     "No match",
			path:     "/login",
			pattern:  "/*/dashboard",
			expected: false,
		},
		{
			name:     "Root path match",
			path:     "/",
			pattern:  "/",
			expected: true,
		},
		{
			name:     "Root path with wildcard",
			path:     "/en",
			pattern:  "/*",
			expected: true,
		},
		{
			name:     "Complex pattern",
			path:     "/de/product/12345/details",
			pattern:  "/*/product/*/details",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PathMatchesPattern(tt.path, tt.pattern)
			if result != tt.expected {
				t.Errorf("PathMatchesPattern(%v, %v) = %v, want %v", tt.path, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestIsLocalizedRoute(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		routeMapping map[string]string
		expected     bool
	}{
		{
			name: "Localized route with locale",
			path: "/en/dashboard",
			routeMapping: map[string]string{
				"/{locale}/dashboard": "template1",
				"/login":             "template2",
			},
			expected: true,
		},
		{
			name: "Non-localized route",
			path: "/login",
			routeMapping: map[string]string{
				"/{locale}/dashboard": "template1",
				"/login":             "template2",
			},
			expected: false,
		},
		{
			name: "Localized route with multiple segments",
			path: "/de/user/123/profile",
			routeMapping: map[string]string{
				"/{locale}/user/{id}/profile": "template1",
				"/login":                      "template2",
			},
			expected: true,
		},
		{
			name: "Route with 2-char segment but not localized",
			path: "/api/v2/users",
			routeMapping: map[string]string{
				"/{locale}/dashboard": "template1",
				"/api/v2/users":       "template2",
			},
			expected: false,
		},
		{
			name: "Empty route mapping",
			path: "/en/dashboard",
			routeMapping: map[string]string{},
			expected:     false,
		},
		{
			name: "Root localized route",
			path: "/fr",
			routeMapping: map[string]string{
				"/{locale}": "template1",
				"/login":   "template2",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLocalizedRoute(tt.path, tt.routeMapping)
			if result != tt.expected {
				t.Errorf("IsLocalizedRoute(%v, routeMapping) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetCurrentRouteWithoutLocale(t *testing.T) {
	tests := []struct {
		name     string
		setupCtx func() context.Context
		expected string
	}{
		{
			name: "Localized route - strip locale",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/en/dashboard")
				ctx = context.WithValue(ctx, shared.RouteMappingKey, map[string]string{
					"/{locale}/dashboard": "template1",
					"/login":             "template2",
				})
				return ctx
			},
			expected: "/dashboard",
		},
		{
			name: "Non-localized route - keep path",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/login")
				ctx = context.WithValue(ctx, shared.RouteMappingKey, map[string]string{
					"/{locale}/dashboard": "template1",
					"/login":             "template2",
				})
				return ctx
			},
			expected: "/login",
		},
		{
			name: "Complex localized route",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/de/user/123/profile")
				ctx = context.WithValue(ctx, shared.RouteMappingKey, map[string]string{
					"/{locale}/user/{id}/profile": "template1",
					"/login":                      "template2",
				})
				return ctx
			},
			expected: "/user/123/profile",
		},
		{
			name: "Route with 2-char segment but not localized",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/api/v2/users")
				ctx = context.WithValue(ctx, shared.RouteMappingKey, map[string]string{
					"/{locale}/dashboard": "template1",
					"/api/v2/users":       "template2",
				})
				return ctx
			},
			expected: "/api/v2/users",
		},
		{
			name: "Root localized route",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/fr")
				ctx = context.WithValue(ctx, shared.RouteMappingKey, map[string]string{
					"/{locale}": "template1",
					"/login":   "template2",
				})
				return ctx
			},
			expected: "/",
		},
		{
			name: "No route mapping - fallback to template check",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/en/dashboard")
				ctx = context.WithValue(ctx, shared.I18nTemplateKey, "app/locale_/dashboard/page.templ")
				return ctx
			},
			expected: "/dashboard",
		},
		{
			name: "No route mapping - no localization active",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RequestPathKey, "/login")
				return ctx
			},
			expected: "/login",
		},
		{
			name: "No request path",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.RouteMappingKey, map[string]string{
					"/{locale}/dashboard": "template1",
				})
				return ctx
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			result := GetCurrentRouteWithoutLocale(ctx)
			if result != tt.expected {
				t.Errorf("GetCurrentRouteWithoutLocale() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCurrentRouteWithoutLocaleFallback(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		setupCtx func() context.Context
		expected string
	}{
		{
			name: "Template present - strip locale",
			path: "/en/dashboard",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.I18nTemplateKey, "app/locale_/dashboard/page.templ")
				return ctx
			},
			expected: "/dashboard",
		},
		{
			name: "No template - keep path",
			path: "/login",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expected: "/login",
		},
		{
			name: "Unknown template - keep path",
			path: "/api/v2/users",
			setupCtx: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, shared.I18nTemplateKey, "unknown")
				return ctx
			},
			expected: "/api/v2/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx()
			result := getCurrentRouteWithoutLocaleFallback(ctx, tt.path)
			if result != tt.expected {
				t.Errorf("getCurrentRouteWithoutLocaleFallback() = %v, want %v", result, tt.expected)
			}
		})
	}
}