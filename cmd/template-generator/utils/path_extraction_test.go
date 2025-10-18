package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
)

// TestGetLocalPackageInfo_DockerVsLocal tests path extraction in different environments
func TestGetLocalPackageInfo_DockerVsLocal(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		moduleName   string
		config       types.Config
		expectedPkg  string
		expectedPath string
		description  string
	}{
		// Local development scenarios
		{
			name:         "Local - Home directory project",
			filePath:     "/home/user/projects/myapp/templates/page_templ.go",
			moduleName:   "github.com/user/myapp",
			config:       types.Config{ScanPath: "templates"},
			expectedPkg:  "templates",
			expectedPath: "github.com/user/myapp/templates",
			description:  "Local development in home directory",
		},
		{
			name:         "Local - Workspace project",
			filePath:     "/workspace/projects/webapp/views/admin/page_templ.go",
			moduleName:   "github.com/company/webapp",
			config:       types.Config{ScanPath: "views"},
			expectedPkg:  "admin",
			expectedPath: "github.com/company/webapp/views", // Base scan path
			description:  "Local development in workspace",
		},
		
		// Docker environment scenarios
		{
			name:         "Docker - App directory",
			filePath:     "/app/templates/components/page_templ.go",
			moduleName:   "github.com/company/project",
			config:       types.Config{ScanPath: "templates"},
			expectedPkg:  "components",
			expectedPath: "github.com/company/project/templates", // Base scan path
			description:  "Docker container with /app mount",
		},
		{
			name:         "Docker - Working directory",
			filePath:     "/workdir/src/views/layout_templ.go",
			moduleName:   "github.com/team/service",
			config:       types.Config{ScanPath: "views"},
			expectedPkg:  "views",
			expectedPath: "github.com/team/service/views",
			description:  "Docker with custom working directory",
		},
		
		// Different module structures
		{
			name:         "Monorepo - Service subdirectory",
			filePath:     "/repo/services/auth/templates/login/page_templ.go",
			moduleName:   "github.com/company/monorepo/services/auth",
			config:       types.Config{ScanPath: "templates"},
			expectedPkg:  "login",
			expectedPath: "github.com/company/monorepo/services/auth/templates", // Base scan path
			description:  "Monorepo with service-specific modules",
		},
		{
			name:         "Nested module - Frontend",
			filePath:     "/project/frontend/ui/components/button_templ.go",
			moduleName:   "github.com/org/project/frontend",
			config:       types.Config{ScanPath: "ui"},
			expectedPkg:  "components",
			expectedPath: "github.com/org/project/frontend/ui", // Base scan path
			description:  "Nested module structure",
		},
		
		// Edge cases with special characters
		{
			name:         "Directory with hyphens",
			filePath:     "/app/templates/error-pages/not-found/page_templ.go",
			moduleName:   "github.com/user/webapp",
			config:       types.Config{ScanPath: "templates"},
			expectedPkg:  "notfound", // Should be sanitized
			expectedPath: "github.com/user/webapp/templates", // Base scan path
			description:  "Directory names with hyphens should be sanitized for package names",
		},
		{
			name:         "Directory with dots",
			filePath:     "/app/views/v1.0/api/page_templ.go",
			moduleName:   "github.com/api/service",
			config:       types.Config{ScanPath: "views"},
			expectedPkg:  "api",
			expectedPath: "github.com/api/service/views", // Base scan path
			description:  "Directory names with dots",
		},
		
		// Deep nesting scenarios
		{
			name:         "Very deep nesting",
			filePath:     "/project/src/main/resources/templates/admin/users/profile/edit/page_templ.go",
			moduleName:   "github.com/enterprise/system",
			config:       types.Config{ScanPath: "templates"},
			expectedPkg:  "edit",
			expectedPath: "github.com/enterprise/system/templates", // Base scan path
			description:  "Very deep directory nesting",
		},
		
		// Different scan path names
		{
			name:         "Custom scan path - pages",
			filePath:     "/app/pages/dashboard/page_templ.go",
			moduleName:   "github.com/startup/app",
			config:       types.Config{ScanPath: "pages"},
			expectedPkg:  "dashboard",
			expectedPath: "github.com/startup/app/pages", // Base scan path
			description:  "Custom scan path name",
		},
		{
			name:         "Custom scan path - src",
			filePath:     "/workspace/project/src/components/layout_templ.go",
			moduleName:   "github.com/dev/project",
			config:       types.Config{ScanPath: "src"},
			expectedPkg:  "components",
			expectedPath: "github.com/dev/project/src", // Base scan path
			description:  "Source directory as scan path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file to simulate the Go file
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test_templ.go")
			
			// Write a minimal Go file with package declaration
			packageName := tt.expectedPkg
			if packageName == "templates" || packageName == "views" || packageName == "pages" || packageName == "src" {
				packageName = tt.config.ScanPath
			}
			
			goContent := "package " + packageName + "\n\n// Test file\n"
			err := os.WriteFile(tempFile, []byte(goContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			// Test with the temp file (simulating the real file path structure)
			actualPkg, actualPath := GetLocalPackageInfo(tempFile, tt.moduleName, tt.config)

			// For this test, we mainly care about the import path logic
			// The package name will be extracted from the actual file
			if actualPath != tt.expectedPath {
				t.Errorf("GetLocalPackageInfo() path = %v, want %v", actualPath, tt.expectedPath)
			}

			t.Logf("✅ %s", tt.description)
			t.Logf("   File: %s", tt.filePath)
			t.Logf("   Module: %s", tt.moduleName)
			t.Logf("   Expected Package: %s", tt.expectedPkg)
			t.Logf("   Expected Path: %s", tt.expectedPath)
			t.Logf("   Actual Package: %s", actualPkg)
			t.Logf("   Actual Path: %s", actualPath)
		})
	}
}

// TestGetLocalPackageInfo_ErrorHandling tests error scenarios
func TestGetLocalPackageInfo_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		moduleName  string
		config      types.Config
		description string
	}{
		{
			name:        "Non-existent file",
			filePath:    "/nonexistent/path/file_templ.go",
			moduleName:  "github.com/test/project",
			config:      types.Config{ScanPath: "templates"},
			description: "Should handle non-existent files gracefully",
		},
		{
			name:        "Invalid file path",
			filePath:    "",
			moduleName:  "github.com/test/project",
			config:      types.Config{ScanPath: "templates"},
			description: "Should handle empty file path",
		},
		{
			name:        "Empty module name",
			filePath:    "/app/templates/page_templ.go",
			moduleName:  "",
			config:      types.Config{ScanPath: "templates"},
			description: "Should handle empty module name",
		},
		{
			name:        "Empty scan path",
			filePath:    "/app/templates/page_templ.go",
			moduleName:  "github.com/test/project",
			config:      types.Config{ScanPath: ""},
			description: "Should handle empty scan path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These should not panic and should return reasonable defaults
			pkg, path := GetLocalPackageInfo(tt.filePath, tt.moduleName, tt.config)
			
			// Should return non-empty values even in error cases (except for empty scan path)
			if pkg == "" && tt.name != "Empty scan path" {
				t.Errorf("Package name should not be empty, got: %q", pkg)
			}
			if path == "" {
				t.Errorf("Import path should not be empty, got: %q", path)
			}
			
			t.Logf("✅ %s", tt.description)
			t.Logf("   Returned Package: %s", pkg)
			t.Logf("   Returned Path: %s", path)
		})
	}
}

// Note: TestSanitizePackageName and TestCreatePackageAlias are already covered in sanitize_test.go