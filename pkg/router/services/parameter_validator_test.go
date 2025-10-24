package services

import (
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParameterValidator_ValidateParameterInheritance(t *testing.T) {
	// Setup
	injector := do.New()
	defer injector.Shutdown()

	// Provide logger and config
	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockConfigService{}, nil
	})

	validator, err := NewParameterValidator(injector)
	require.NoError(t, err)

	tests := []struct {
		name           string
		route          *interfaces.Route
		config         *interfaces.ConfigFile
		hierarchyMap   map[string][]string
		configs        map[string]*interfaces.ConfigFile
		expectedErrors int
		expectedWarnings int
		expectedWarningTypes []string
	}{
		{
			name: "no inheritance - no conflicts",
			route: &interfaces.Route{
				Path:         "/user/profile",
				TemplateFile: "user/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[0-9]+$",
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{},
			configs:      map[string]*interfaces.ConfigFile{},
			expectedErrors: 0,
			expectedWarnings: 0,
		},
		{
			name: "compatible parameter inheritance",
			route: &interfaces.Route{
				Path:         "/user/$id/profile",
				TemplateFile: "user/id_/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[0-9]+$",
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/profile": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								Validation: "^[0-9]+$",
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 0,
		},
		{
			name: "conflicting validation regex",
			route: &interfaces.Route{
				Path:         "/user/$id/profile",
				TemplateFile: "user/id_/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[a-z]+$", // Different from parent
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/profile": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								Validation: "^[0-9]+$", // Different from child
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"PARAMETER_VALIDATION_CONFLICT"},
		},
		{
			name: "conflicting supported values",
			route: &interfaces.Route{
				Path:         "/user/$role/dashboard",
				TemplateFile: "user/role_/dashboard/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"role": {
							SupportedValues: []string{"admin", "user", "guest", "superuser"}, // Contains values not in parent
							Description: "User role",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$role/dashboard": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"role": {
								SupportedValues: []string{"admin", "user", "guest"}, // Subset of child
								Description: "User role",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"PARAMETER_VALUES_CONFLICT"},
		},
		{
			name: "child adds validation (good practice)",
			route: &interfaces.Route{
				Path:         "/user/$id/settings",
				TemplateFile: "user/id_/settings/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						"id": {
							Validation: "^[0-9]+$", // Child adds validation
							Description: "User ID",
						},
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/settings": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								// No validation in parent
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"PARAMETER_INHERITANCE_INFO"},
		},
		{
			name: "missing inherited parameter config",
			route: &interfaces.Route{
				Path:         "/user/$id/profile",
				TemplateFile: "user/id_/profile/page.templ",
			},
			config: &interfaces.ConfigFile{
				DynamicSettings: &interfaces.DynamicSettings{
					Parameters: map[string]*interfaces.DynamicParameterConfig{
						// Missing "id" parameter config
					},
				},
			},
			hierarchyMap: map[string][]string{
				"/user/$id/profile": {"/user"},
			},
			configs: map[string]*interfaces.ConfigFile{
				"user/page.templ": {
					DynamicSettings: &interfaces.DynamicSettings{
						Parameters: map[string]*interfaces.DynamicParameterConfig{
							"id": {
								Validation: "^[0-9]+$",
								Description: "User ID",
							},
						},
					},
				},
			},
			expectedErrors: 0,
			expectedWarnings: 1,
			expectedWarningTypes: []string{"INHERITED_PARAMETER_MISSING_CONFIG"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &ValidationResult{
				Errors:   make([]ValidationError, 0),
				Warnings: make([]ValidationWarning, 0),
			}

			validator.ValidateParameterInheritance(tt.route, tt.config, tt.hierarchyMap, tt.configs, result)

			assert.Equal(t, tt.expectedErrors, len(result.Errors), "Expected %d errors, got %d", tt.expectedErrors, len(result.Errors))
			assert.Equal(t, tt.expectedWarnings, len(result.Warnings), "Expected %d warnings, got %d", tt.expectedWarnings, len(result.Warnings))

			// Check warning types
			if len(tt.expectedWarningTypes) > 0 {
				actualWarningTypes := make([]string, len(result.Warnings))
				for i, warning := range result.Warnings {
					actualWarningTypes[i] = warning.Type
				}
				assert.ElementsMatch(t, tt.expectedWarningTypes, actualWarningTypes, "Warning types don't match")
			}
		})
	}
}

func TestParameterValidator_IsSubsetOfValues(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockConfigService{}, nil
	})

	validator, err := NewParameterValidator(injector)
	require.NoError(t, err)

	tests := []struct {
		name         string
		childValues  []string
		parentValues []string
		expected     bool
	}{
		{
			name:         "empty child values",
			childValues:  []string{},
			parentValues: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:         "child is subset of parent",
			childValues:  []string{"a", "b"},
			parentValues: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:         "child equals parent",
			childValues:  []string{"a", "b", "c"},
			parentValues: []string{"a", "b", "c"},
			expected:     true,
		},
		{
			name:         "child has extra values",
			childValues:  []string{"a", "b", "c", "d"},
			parentValues: []string{"a", "b", "c"},
			expected:     false,
		},
		{
			name:         "completely different values",
			childValues:  []string{"x", "y", "z"},
			parentValues: []string{"a", "b", "c"},
			expected:     false,
		},
		{
			name:         "partial overlap",
			childValues:  []string{"a", "x"},
			parentValues: []string{"a", "b", "c"},
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isSubsetOfValues(tt.childValues, tt.parentValues)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParameterValidator_IsConfigForRoute(t *testing.T) {
	injector := do.New()
	defer injector.Shutdown()

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})
	do.Provide(injector, func(i do.Injector) (interfaces.ConfigService, error) {
		return &mockConfigService{}, nil
	})

	validator, err := NewParameterValidator(injector)
	require.NoError(t, err)

	tests := []struct {
		name         string
		templateFile string
		routePath    string
		expected     bool
	}{
		{
			name:         "exact match",
			templateFile: "user/page.templ",
			routePath:    "/user",
			expected:     true,
		},
		{
			name:         "nested route",
			templateFile: "user/profile/page.templ",
			routePath:    "/user/profile",
			expected:     true,
		},
		{
			name:         "dynamic parameter",
			templateFile: "user/id_/page.templ",
			routePath:    "/user/$id",
			expected:     true, // "user" segment matches
		},
		{
			name:         "no match",
			templateFile: "admin/page.templ",
			routePath:    "/user",
			expected:     false,
		},
		{
			name:         "partial match",
			templateFile: "user/settings/page.templ",
			routePath:    "/user/profile",
			expected:     false, // "profile" not in template file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isConfigForRoute(tt.templateFile, tt.routePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

