package i18n

import (
	"context"
	"testing"

	"github.com/denkhaus/templ-router/pkg/shared"
)

func TestLocalizeRouteIfRequired(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		locale   string
		expected string
	}{
		{
			name:     "Path with locale placeholder",
			path:     "/{locale}/dashboard",
			locale:   "en",
			expected: "/en/dashboard",
		},
		{
			name:     "Path with locale placeholder - German",
			path:     "/{locale}/benutzer/{id}",
			locale:   "de",
			expected: "/de/benutzer/{id}",
		},
		{
			name:     "Path without locale placeholder",
			path:     "/login",
			locale:   "en",
			expected: "/login",
		},
		{
			name:     "Path without locale placeholder - with params",
			path:     "/api/users/{id}",
			locale:   "fr",
			expected: "/api/users/{id}",
		},
		{
			name:     "Root path with locale",
			path:     "/{locale}",
			locale:   "es",
			expected: "/es",
		},
		{
			name:     "Multiple locale placeholders",
			path:     "/{locale}/admin/{locale}/settings",
			locale:   "it",
			expected: "/it/admin/it/settings",
		},
		{
			name:     "Empty locale defaults to en",
			path:     "/{locale}/profile",
			locale:   "",
			expected: "/en/profile",
		},
		{
			name:     "Complex nested path",
			path:     "/{locale}/category/{category}/product/{id}",
			locale:   "ja",
			expected: "/ja/category/{category}/product/{id}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create context with locale
			ctx := context.Background()
			if tt.locale != "" {
				ctx = context.WithValue(ctx, shared.LocaleKey, tt.locale)
			}

			result := LocalizeRouteIfRequired(ctx, tt.path)
			if result != tt.expected {
				t.Errorf("LocalizeRouteIfRequired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLocalizeRouteIfRequired_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		locale   string
		expected string
	}{
		{
			name:     "Empty path",
			path:     "",
			locale:   "en",
			expected: "",
		},
		{
			name:     "Path with only locale placeholder",
			path:     "{locale}",
			locale:   "en",
			expected: "en",
		},
		{
			name:     "Path with locale in middle",
			path:     "/api/{locale}/v1/users",
			locale:   "en",
			expected: "/api/en/v1/users",
		},
		{
			name:     "Path with locale at end",
			path:     "/switch-language/{locale}",
			locale:   "fr",
			expected: "/switch-language/fr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create context with locale
			ctx := context.Background()
			if tt.locale != "" {
				ctx = context.WithValue(ctx, shared.LocaleKey, tt.locale)
			}

			result := LocalizeRouteIfRequired(ctx, tt.path)
			if result != tt.expected {
				t.Errorf("LocalizeRouteIfRequired() = %v, want %v", result, tt.expected)
			}
		})
	}
}
