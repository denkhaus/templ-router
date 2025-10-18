package validate

import (
	"strings"
	"testing"

	"github.com/denkhaus/templ-router/cmd/trgen/types"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      types.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			config: types.Config{
				ScanPath:    "app",
				OutputDir:   "generated/templates",
				ModuleName:  "github.com/user/project",
				PackageName: "templates",
			},
			expectError: false,
		},
		{
			name: "Empty scan path",
			config: types.Config{
				ScanPath:    "",
				OutputDir:   "generated/templates",
				ModuleName:  "github.com/user/project",
				PackageName: "templates",
			},
			expectError: true,
			errorMsg:    "scan path cannot be empty",
		},
		{
			name: "Empty output directory",
			config: types.Config{
				ScanPath:    "app",
				OutputDir:   "",
				ModuleName:  "github.com/user/project",
				PackageName: "templates",
			},
			expectError: true,
			errorMsg:    "output directory cannot be empty",
		},
		{
			name: "Empty module name",
			config: types.Config{
				ScanPath:    "app",
				OutputDir:   "generated/templates",
				ModuleName:  "",
				PackageName: "templates",
			},
			expectError: true,
			errorMsg:    "module name cannot be empty",
		},
		{
			name: "Empty package name",
			config: types.Config{
				ScanPath:    "app",
				OutputDir:   "generated/templates",
				ModuleName:  "github.com/user/project",
				PackageName: "",
			},
			expectError: true,
			errorMsg:    "package name cannot be empty",
		},
		{
			name: "Invalid package name with hyphen",
			config: types.Config{
				ScanPath:    "app",
				OutputDir:   "generated/templates",
				ModuleName:  "github.com/user/project",
				PackageName: "my-templates",
			},
			expectError: true,
			errorMsg:    "package name contains invalid characters",
		},
		{
			name: "Invalid package name starting with number",
			config: types.Config{
				ScanPath:    "app",
				OutputDir:   "generated/templates",
				ModuleName:  "github.com/user/project",
				PackageName: "123templates",
			},
			expectError: true,
			errorMsg:    "package name cannot start with a number",
		},
		{
			name: "Invalid module name format",
			config: types.Config{
				ScanPath:    "app",
				OutputDir:   "generated/templates",
				ModuleName:  "invalid module name",
				PackageName: "templates",
			},
			expectError: true,
			errorMsg:    "module name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateTemplates(t *testing.T) {
	tests := []struct {
		name        string
		templates   []types.TemplateInfo
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid templates",
			templates: []types.TemplateInfo{
				{
					FunctionName: "Page",
					PackageName:  "app",
					ImportPath:   "github.com/user/project/app",
					PackageAlias: "app",
					RoutePattern: "/",
					TemplateKey:  "key-1",
				},
				{
					FunctionName: "Layout",
					PackageName:  "app",
					ImportPath:   "github.com/user/project/app",
					PackageAlias: "app",
					RoutePattern: "/layout",
					TemplateKey:  "key-2",
				},
			},
			expectError: false,
		},
		{
			name:        "Empty templates",
			templates:   []types.TemplateInfo{},
			expectError: false, // Empty templates should be allowed
		},
		{
			name: "Duplicate template keys",
			templates: []types.TemplateInfo{
				{
					FunctionName: "Page",
					PackageName:  "app",
					ImportPath:   "github.com/user/project/app",
					PackageAlias: "app",
					RoutePattern: "/",
					TemplateKey:  "duplicate-key",
				},
				{
					FunctionName: "Layout",
					PackageName:  "app",
					ImportPath:   "github.com/user/project/app",
					PackageAlias: "app",
					RoutePattern: "/layout",
					TemplateKey:  "duplicate-key",
				},
			},
			expectError: true,
			errorMsg:    "duplicate template key",
		},
		{
			name: "Duplicate route patterns",
			templates: []types.TemplateInfo{
				{
					FunctionName: "Page",
					PackageName:  "app",
					ImportPath:   "github.com/user/project/app",
					PackageAlias: "app",
					RoutePattern: "/same-route",
					TemplateKey:  "key-1",
				},
				{
					FunctionName: "Layout",
					PackageName:  "app",
					ImportPath:   "github.com/user/project/app",
					PackageAlias: "app",
					RoutePattern: "/same-route",
					TemplateKey:  "key-2",
				},
			},
			expectError: true,
			errorMsg:    "duplicate route pattern",
		},
		{
			name: "Invalid package alias",
			templates: []types.TemplateInfo{
				{
					FunctionName: "Page",
					PackageName:  "errordemo",
					ImportPath:   "github.com/user/project/app/error-demo",
					PackageAlias: "error-demo", // Invalid Go identifier
					RoutePattern: "/error-demo",
					TemplateKey:  "key-1",
				},
			},
			expectError: true,
			errorMsg:    "invalid package alias",
		},
		{
			name: "Empty required fields",
			templates: []types.TemplateInfo{
				{
					FunctionName: "",
					PackageName:  "app",
					ImportPath:   "github.com/user/project/app",
					PackageAlias: "app",
					RoutePattern: "/",
					TemplateKey:  "key-1",
				},
			},
			expectError: true,
			errorMsg:    "function name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplates(tt.templates)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}