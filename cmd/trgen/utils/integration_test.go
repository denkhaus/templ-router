package utils

import (
	"path/filepath"
	"testing"

	"github.com/denkhaus/templ-router/cmd/trgen/types"
)

// TestGetLocalPackageInfoWithRealFiles tests with actual files from the demo project
func TestGetLocalPackageInfoWithRealFiles(t *testing.T) {
	// Get the project root directory
	projectRoot, err := filepath.Abs("../../../")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	tests := []struct {
		name         string
		filePath     string
		moduleName   string
		config       types.Config
		expectedPkg  string
		expectedPath string
	}{
		{
			name:         "Demo root app package",
			filePath:     filepath.Join(projectRoot, "demo/app/page_templ.go"),
			moduleName:   "github.com/denkhaus/templ-router/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "app",
			expectedPath: "github.com/denkhaus/templ-router/demo/app",
		},
		{
			name:         "Demo locale subdirectory",
			filePath:     filepath.Join(projectRoot, "demo/app/locale_/page_templ.go"),
			moduleName:   "github.com/denkhaus/templ-router/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "locale_",
			expectedPath: "github.com/denkhaus/templ-router/demo/app/locale_",
		},
		{
			name:         "Demo admin subdirectory",
			filePath:     filepath.Join(projectRoot, "demo/app/locale_/admin/page_templ.go"),
			moduleName:   "github.com/denkhaus/templ-router/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "admin",
			expectedPath: "github.com/denkhaus/templ-router/demo/app/locale_/admin",
		},
		{
			name:         "Demo dashboard subdirectory",
			filePath:     filepath.Join(projectRoot, "demo/app/locale_/dashboard/page_templ.go"),
			moduleName:   "github.com/denkhaus/templ-router/demo",
			config:       types.Config{ScanPath: "app"},
			expectedPkg:  "dashboard",
			expectedPath: "github.com/denkhaus/templ-router/demo/app/locale_/dashboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualPkg, actualPath := GetLocalPackageInfo(tt.filePath, tt.moduleName, tt.config)

			if actualPkg != tt.expectedPkg {
				t.Errorf("GetLocalPackageInfo() package = %v, want %v", actualPkg, tt.expectedPkg)
			}
			if actualPath != tt.expectedPath {
				t.Errorf("GetLocalPackageInfo() path = %v, want %v", actualPath, tt.expectedPath)
			}

			t.Logf("✅ File: %s", tt.filePath)
			t.Logf("✅ Package: %s", actualPkg)
			t.Logf("✅ Import Path: %s", actualPath)
		})
	}
}