package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
)

func TestGenerateRegistry(t *testing.T) {
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "generated", "templates")
	
	// Create test templates
	templates := []types.TemplateInfo{
		{
			FunctionName:  "Page",
			PackageName:   "app",
			ImportPath:    "github.com/test/project/app",
			PackageAlias:  "app",
			RoutePattern:  "/",
			TemplateKey:   "test-key-1",
			FilePath:      "/test/app/page_templ.go",
			HumanName:     "Page",
		},
		{
			FunctionName:  "Layout",
			PackageName:   "app",
			ImportPath:    "github.com/test/project/app",
			PackageAlias:  "app",
			RoutePattern:  "/layout",
			TemplateKey:   "test-key-2",
			FilePath:      "/test/app/layout_templ.go",
			HumanName:     "Layout",
		},
		{
			FunctionName:  "Page",
			PackageName:   "errordemo",
			ImportPath:    "github.com/test/project/app/error-demo",
			PackageAlias:  "errordemo",
			RoutePattern:  "/error-demo",
			TemplateKey:   "test-key-3",
			FilePath:      "/test/app/error-demo/page_templ.go",
			HumanName:     "error-demo.Page",
		},
	}

	config := types.Config{
		ScanPath:      "app",
		OutputDir:     outputDir,
		ModuleName:    "github.com/test/project",
		PackageName:   "templates",
	}

	err := GenerateRegistry(config, templates)
	if err != nil {
		t.Fatalf("GenerateRegistry failed: %v", err)
	}

	// Check if registry file was created
	registryPath := filepath.Join(outputDir, "registry.go")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		t.Fatalf("Registry file was not created: %s", registryPath)
	}

	// Read and verify registry content
	content, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("Failed to read registry file: %v", err)
	}

	contentStr := string(content)

	// Check package declaration
	if !strings.Contains(contentStr, "package templates") {
		t.Error("Registry should have correct package declaration")
	}

	// Check imports
	expectedImports := []string{
		`app "github.com/test/project/app"`,
		`errordemo "github.com/test/project/app/error-demo"`,
	}
	for _, imp := range expectedImports {
		if !strings.Contains(contentStr, imp) {
			t.Errorf("Registry should contain import: %s", imp)
		}
	}

	// Check template mappings
	expectedMappings := []string{
		`"test-key-1": app.Page,`,
		`"test-key-2": app.Layout,`,
		`"test-key-3": errordemo.Page,`,
	}
	for _, mapping := range expectedMappings {
		if !strings.Contains(contentStr, mapping) {
			t.Errorf("Registry should contain template mapping: %s", mapping)
		}
	}

	// Check route mappings
	expectedRoutes := []string{
		`"/": "test-key-1",`,
		`"/layout": "test-key-2",`,
		`"/error-demo": "test-key-3",`,
	}
	for _, route := range expectedRoutes {
		if !strings.Contains(contentStr, route) {
			t.Errorf("Registry should contain route mapping: %s", route)
		}
	}

	// Verify no invalid identifiers
	if strings.Contains(contentStr, "error-demo \"") {
		t.Error("Registry should not contain invalid Go identifiers like 'error-demo'")
	}
}

func TestGenerateRegistryEmptyTemplates(t *testing.T) {
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "generated", "templates")
	
	config := types.Config{
		ScanPath:      "app",
		OutputDir:     outputDir,
		ModuleName:    "github.com/test/project",
		PackageName:   "templates",
	}

	err := GenerateRegistry(config, []types.TemplateInfo{})
	if err != nil {
		t.Fatalf("GenerateRegistry failed with empty templates: %v", err)
	}

	// Check if registry file was created
	registryPath := filepath.Join(outputDir, "registry.go")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		t.Fatalf("Registry file was not created: %s", registryPath)
	}

	// Read and verify registry content
	content, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("Failed to read registry file: %v", err)
	}

	contentStr := string(content)

	// Should still have basic structure
	if !strings.Contains(contentStr, "package templates") {
		t.Error("Registry should have correct package declaration")
	}

	if !strings.Contains(contentStr, "func NewTemplateRegistry") {
		t.Error("Registry should have NewTemplateRegistry function")
	}
}

func TestGenerateRegistryInvalidOutputDir(t *testing.T) {
	// Try to write to a non-existent directory without creating it
	invalidDir := "/nonexistent/path/that/should/not/exist"
	
	templates := []types.TemplateInfo{
		{
			FunctionName: "Page",
			PackageName:  "app",
			ImportPath:   "github.com/test/project/app",
			PackageAlias: "app",
			RoutePattern: "/",
			TemplateKey:  "test-key-1",
		},
	}

	config := types.Config{
		ScanPath:    "app",
		OutputDir:   invalidDir,
		ModuleName:  "github.com/test/project",
		PackageName: "templates",
	}

	err := GenerateRegistry(config, templates)
	if err == nil {
		t.Error("GenerateRegistry should fail with invalid output directory")
	}
}