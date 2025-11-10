package utils

import (
	"testing"

	"github.com/denkhaus/templ-router/cmd/trgen/types"
)

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
		{
			name:         "Directory with hyphen - error-demo",
			filePath:     "/demo/app/locale_/dashboard/error-demo/page_templ.go",
			functionName: "Page",
			config:       types.Config{ScanPath: "app"},
			expected:     "/{locale}/dashboard/error-demo",
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
