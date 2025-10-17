package utils

import (
	"testing"

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
)

func TestGetLocalPackageInfo(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		moduleName   string
		config       types.Config
		expectedPkg  string
		expectedPath string
	}{
		{
			name:         "Local development - demo root",
			filePath:     "/home/user/project/demo/app/page_templ.go",
			moduleName:   "github.com/user/project/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "app",
			expectedPath: "github.com/user/project/demo/app",
		},
		{
			name:         "Local development - demo subdirectory",
			filePath:     "/home/user/project/demo/app/locale_/admin/page_templ.go",
			moduleName:   "github.com/user/project/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "admin",
			expectedPath: "github.com/user/project/demo/app/locale_/admin",
		},
		{
			name:         "Docker environment - demo root",
			filePath:     "/app/demo/app/page_templ.go",
			moduleName:   "github.com/user/project/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "app",
			expectedPath: "github.com/user/project/demo/app",
		},
		{
			name:         "Docker environment - demo subdirectory",
			filePath:     "/app/demo/app/locale_/admin/page_templ.go",
			moduleName:   "github.com/user/project/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "admin",
			expectedPath: "github.com/user/project/demo/app/locale_/admin",
		},
		{
			name:         "Different module structure - root project",
			filePath:     "/home/user/myproject/templates/page_templ.go",
			moduleName:   "github.com/user/myproject",
			config:       types.Config{ScanPath: "templates"},
			expectedPkg:  "templates",
			expectedPath: "github.com/user/myproject/templates",
		},
		{
			name:         "Different module structure - subdirectory",
			filePath:     "/home/user/myproject/templates/components/button_templ.go",
			moduleName:   "github.com/user/myproject",
			config:       types.Config{ScanPath: "templates"},
			expectedPkg:  "components",
			expectedPath: "github.com/user/myproject/templates/components",
		},
		{
			name:         "Complex nested structure",
			filePath:     "/workspace/projects/webapp/frontend/views/admin/users/page_templ.go",
			moduleName:   "github.com/company/webapp/frontend",
			config:       types.Config{ScanPath: "views"},
			expectedPkg:  "users",
			expectedPath: "github.com/company/webapp/frontend/views/admin/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary test file with package declaration
			testContent := "package " + tt.expectedPkg + "\n"
			testFile := createTempFile(t, testContent)
			defer removeTempFile(t, testFile)

			// Use the test file path but with the expected directory structure
			actualPkg, actualPath := GetLocalPackageInfo(tt.filePath, tt.moduleName, tt.config)

			if actualPkg != tt.expectedPkg {
				t.Errorf("GetLocalPackageInfo() package = %v, want %v", actualPkg, tt.expectedPkg)
			}
			if actualPath != tt.expectedPath {
				t.Errorf("GetLocalPackageInfo() path = %v, want %v", actualPath, tt.expectedPath)
			}
		})
	}
}

func TestCreateRoutePattern(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		functionName string
		config       types.Config
		expected     string
	}{
		{
			name:         "Root page",
			filePath:     "/demo/app/page_templ.go",
			functionName: "Page",
			config:       types.Config{ScanPath: "app"},
			expected:     "/",
		},
		{
			name:         "Locale page",
			filePath:     "/demo/app/locale_/page_templ.go",
			functionName: "Page",
			config:       types.Config{ScanPath: "app"},
			expected:     "/{locale}",
		},
		{
			name:         "Admin page",
			filePath:     "/demo/app/locale_/admin/page_templ.go",
			functionName: "Page",
			config:       types.Config{ScanPath: "app"},
			expected:     "/{locale}/admin",
		},
		{
			name:         "Dynamic product page",
			filePath:     "/demo/app/locale_/product/id_/page_templ.go",
			functionName: "Page",
			config:       types.Config{ScanPath: "app"},
			expected:     "/{locale}/product/{id}",
		},
		{
			name:         "Error template",
			filePath:     "/demo/app/locale_/dashboard/error_templ.go",
			functionName: "Error",
			config:       types.Config{ScanPath: "app"},
			expected:     "/{locale}/dashboard/error",
		},
		{
			name:         "Layout template",
			filePath:     "/demo/app/layout_templ.go",
			functionName: "Layout",
			config:       types.Config{ScanPath: "app"},
			expected:     "/layout",
		},
		{
			name:         "Component template",
			filePath:     "/demo/app/layout_templ.go",
			functionName: "Navbar",
			config:       types.Config{ScanPath: "app"},
			expected:     "/navbar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CreateRoutePattern(tt.filePath, tt.functionName, tt.config)
			if actual != tt.expected {
				t.Errorf("CreateRoutePattern() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

// Helper functions for testing
func createTempFile(t *testing.T, content string) string {
	// For testing, we'll mock the file parsing
	// In a real implementation, you'd create actual temp files
	return "/tmp/test_file.go"
}

func removeTempFile(t *testing.T, path string) {
	// Cleanup temp file
}

// Mock the file parsing for testing
func init() {
	// We need to override the file parsing logic for tests
	// This is a simplified approach for testing
}