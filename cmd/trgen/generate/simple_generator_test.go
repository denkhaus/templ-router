package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/denkhaus/templ-router/cmd/trgen/types"
)

// TestGenerateRegistry_Basic tests the basic registry generation
func TestGenerateRegistry_Basic(t *testing.T) {
	tempDir := t.TempDir()
	
	config := types.Config{
		ModuleName:  "github.com/test/project",
		ScanPath:    tempDir,
		OutputDir:   tempDir,
		PackageName: "templates",
	}
	
	// Create some sample template info
	templates := []types.TemplateInfo{
		{
			FilePath:     "/app/templates/home_templ.go",
			FunctionName: "HomePage",
			TemplateKey:  "home-page",
			PackageName:  "templates",
			PackageAlias: "templates",
			ImportPath:   "github.com/test/project/templates",
			HumanName:    "Home Page",
		},
		{
			FilePath:     "/app/templates/about_templ.go",
			FunctionName: "AboutPage",
			TemplateKey:  "about-page",
			PackageName:  "templates",
			PackageAlias: "templates",
			ImportPath:   "github.com/test/project/templates",
			HumanName:    "About Page",
		},
	}
	
	// Generate registry
	err := GenerateRegistry(config, templates)
	if err != nil {
		t.Fatalf("GenerateRegistry failed: %v", err)
	}
	
	// Check output file exists
	outputFile := filepath.Join(config.OutputDir, "registry.go")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}
	
	// Read and verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	contentStr := string(content)
	
	// Check for expected content
	expectedContent := []string{
		"package templates",
		"HomePage",
		"AboutPage",
		"home-page",
		"about-page",
		"github.com/test/project/templates",
	}
	
	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content %q not found in generated file", expected)
		}
	}
	
	t.Logf("✅ Basic registry generation working")
	t.Logf("   Generated file: %s", outputFile)
	t.Logf("   Templates: %d", len(templates))
}

// TestGenerateRegistry_EmptyTemplates tests generation with no templates
func TestGenerateRegistry_EmptyTemplates(t *testing.T) {
	tempDir := t.TempDir()
	
	config := types.Config{
		ModuleName:  "github.com/test/project",
		ScanPath:    tempDir,
		OutputDir:   tempDir,
		PackageName: "templates",
	}
	
	// Generate registry with empty templates
	err := GenerateRegistry(config, []types.TemplateInfo{})
	if err != nil {
		t.Fatalf("GenerateRegistry failed: %v", err)
	}
	
	// Check output file exists
	outputFile := filepath.Join(config.OutputDir, "registry.go")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %s", outputFile)
	}
	
	// Read and verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	contentStr := string(content)
	
	// Should still have package declaration and basic structure
	if !strings.Contains(contentStr, "package templates") {
		t.Error("Expected package declaration in generated file")
	}
	
	t.Logf("✅ Empty templates generation working")
	t.Logf("   Generated file: %s", outputFile)
}

// TestGenerateRegistry_NestedPackages tests generation with nested packages
func TestGenerateRegistry_NestedPackages(t *testing.T) {
	tempDir := t.TempDir()
	
	config := types.Config{
		ModuleName:  "github.com/company/app",
		ScanPath:    tempDir,
		OutputDir:   tempDir,
		PackageName: "templates",
	}
	
	// Create templates from different packages
	templates := []types.TemplateInfo{
		{
			FilePath:     "/app/admin/users_templ.go",
			FunctionName: "UsersPage",
			TemplateKey:  "admin-users",
			PackageName:  "admin",
			PackageAlias: "admin",
			ImportPath:   "github.com/company/app/admin",
			HumanName:    "Admin Users",
		},
		{
			FilePath:     "/app/public/home_templ.go",
			FunctionName: "HomePage",
			TemplateKey:  "public-home",
			PackageName:  "public",
			PackageAlias: "public",
			ImportPath:   "github.com/company/app/public",
			HumanName:    "Public Home",
		},
		{
			FilePath:     "/app/api/docs_templ.go",
			FunctionName: "DocsPage",
			TemplateKey:  "api-docs",
			PackageName:  "api",
			PackageAlias: "api",
			ImportPath:   "github.com/company/app/api",
			HumanName:    "API Documentation",
		},
	}
	
	// Generate registry
	err := GenerateRegistry(config, templates)
	if err != nil {
		t.Fatalf("GenerateRegistry failed: %v", err)
	}
	
	// Check output file exists
	outputFile := filepath.Join(config.OutputDir, "registry.go")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	contentStr := string(content)
	
	// Check for all packages and functions
	expectedContent := []string{
		"package templates",
		"UsersPage",
		"HomePage", 
		"DocsPage",
		"admin-users",
		"public-home",
		"api-docs",
		"github.com/company/app/admin",
		"github.com/company/app/public",
		"github.com/company/app/api",
	}
	
	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected content %q not found in generated file", expected)
		}
	}
	
	t.Logf("✅ Nested packages generation working")
	t.Logf("   Generated file: %s", outputFile)
	t.Logf("   Packages: admin, public, api")
}

// TestGenerateRegistry_ConfigValidation tests configuration validation
func TestGenerateRegistry_ConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      types.Config
		expectError bool
		description string
	}{
		{
			name: "Valid config",
			config: types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    "/tmp",
				OutputDir:   "/tmp",
				PackageName: "templates",
			},
			expectError: false,
			description: "Should accept valid configuration",
		},
		{
			name: "Empty module name",
			config: types.Config{
				ModuleName:  "",
				ScanPath:    "/tmp",
				OutputDir:   "/tmp",
				PackageName: "templates",
			},
			expectError: false, // Actually doesn't error, just uses empty module
			description: "Should handle empty module name gracefully",
		},
		{
			name: "Empty package name",
			config: types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    "/tmp",
				OutputDir:   "/tmp",
				PackageName: "",
			},
			expectError: false, // Actually doesn't error, just uses empty package
			description: "Should handle empty package name gracefully",
		},
		{
			name: "Empty output directory",
			config: types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    "/tmp",
				OutputDir:   "",
				PackageName: "templates",
			},
			expectError: true,
			description: "Should reject empty output directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateRegistry(tt.config, []types.TemplateInfo{})
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else {
					t.Logf("✅ %s - Got expected error: %v", tt.description, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					t.Logf("✅ %s", tt.description)
				}
			}
		})
	}
}