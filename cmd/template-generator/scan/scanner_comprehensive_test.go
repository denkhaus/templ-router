package scan

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
	"golang.org/x/tools/go/packages"
)

// TestExtractTemplatesFromFile tests the core template extraction logic
func TestExtractTemplatesFromFile(t *testing.T) {
	tests := []struct {
		name            string
		goFileContent   string
		expectedFuncs   []string
		expectedSkipped []string
		description     string
	}{
		{
			name: "Simple template functions",
			goFileContent: `package templates

import "github.com/a-h/templ"

func HomePage() templ.Component {
	return templ.Raw("Home")
}

func AboutPage() templ.Component {
	return templ.Raw("About")
}`,
			expectedFuncs:   []string{"HomePage", "AboutPage"},
			expectedSkipped: []string{},
			description:     "Should extract simple template functions",
		},
		{
			name: "Functions with parameters",
			goFileContent: `package templates

import "github.com/a-h/templ"

func UserProfile(userID string, isAdmin bool) templ.Component {
	return templ.Raw("Profile")
}

func ProductCard(product Product) templ.Component {
	return templ.Raw("Product")
}

func SimpleCard() templ.Component {
	return templ.Raw("Card")
}`,
			expectedFuncs:   []string{"SimpleCard"},
			expectedSkipped: []string{"UserProfile", "ProductCard"},
			description:     "Should skip functions with parameters except special ones",
		},
		{
			name: "Special routing functions",
			goFileContent: `package templates

import "github.com/a-h/templ"

func Page(data interface{}) templ.Component {
	return templ.Raw("Page")
}

func Layout(content templ.Component) templ.Component {
	return templ.Raw("Layout")
}

func Error(code int, message string) templ.Component {
	return templ.Raw("Error")
}`,
			expectedFuncs:   []string{"Page", "Layout", "Error"},
			expectedSkipped: []string{},
			description:     "Should include special routing functions even with parameters",
		},
		{
			name: "Mixed function types",
			goFileContent: `package templates

import "github.com/a-h/templ"

// Regular template function
func Header() templ.Component {
	return templ.Raw("Header")
}

// Method (should be skipped)
func (t *Template) Render() templ.Component {
	return templ.Raw("Render")
}

// Function with parameters (should be skipped)
func Button(text string) templ.Component {
	return templ.Raw("Button")
}

// Non-template function (should be skipped)
func HelperFunction() string {
	return "helper"
}`,
			expectedFuncs:   []string{"Header"},
			expectedSkipped: []string{"Render", "Button", "HelperFunction"},
			description:     "Should handle mixed function types correctly",
		},
		{
			name: "Invalid naming conventions",
			goFileContent: `package templates

import "github.com/a-h/templ"

func validFunction() templ.Component {
	return templ.Raw("Valid")
}

func invalid_function() templ.Component {
	return templ.Raw("Invalid")
}

func ValidFunction() templ.Component {
	return templ.Raw("Valid")
}`,
			expectedFuncs:   []string{"ValidFunction"},
			expectedSkipped: []string{"validFunction", "invalid_function"},
			description:     "Should enforce naming conventions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test_templ.go")
			
			err := os.WriteFile(tempFile, []byte(tt.goFileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			// Parse the file
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tempFile, nil, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse file: %v", err)
			}

			// Create mock package info
			config := types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    tempDir,
				PackageName: "templates",
			}

			// Load package for type information
			cfg := &packages.Config{
				Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
				Dir:  tempDir,
			}

			pkgs, err := packages.Load(cfg, ".")
			if err != nil {
				t.Fatalf("Failed to load package: %v", err)
			}

			if len(pkgs) == 0 {
				t.Fatal("No packages loaded")
			}

			pkg := pkgs[0]

			// Extract templates
			templates, validationErrors, err := ExtractTemplatesFromFile(file, tempFile, pkg, config)
			if err != nil {
				t.Fatalf("ExtractTemplatesFromFile failed: %v", err)
			}

			// Check expected functions are found
			foundFuncs := make(map[string]bool)
			for _, template := range templates {
				foundFuncs[template.FunctionName] = true
			}

			for _, expectedFunc := range tt.expectedFuncs {
				if !foundFuncs[expectedFunc] {
					t.Errorf("Expected function %s not found in templates", expectedFunc)
				}
			}

			// Check unexpected functions are not found
			for _, skippedFunc := range tt.expectedSkipped {
				if foundFuncs[skippedFunc] {
					t.Errorf("Function %s should have been skipped but was found", skippedFunc)
				}
			}

			t.Logf("✅ %s", tt.description)
			t.Logf("   Found %d templates", len(templates))
			t.Logf("   Validation errors: %d", len(validationErrors))
		})
	}
}

// TestScanTemplatesWithPackages_ErrorHandling tests error scenarios
func TestScanTemplatesWithPackages_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		config      types.Config
		setupFiles  map[string]string
		expectError bool
		description string
	}{
		{
			name: "Non-existent directory",
			config: types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    "/nonexistent/path",
				PackageName: "templates",
			},
			setupFiles:  nil,
			expectError: true,
			description: "Should handle non-existent scan directory",
		},
		{
			name: "Empty directory",
			config: types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    "", // Will be set to temp dir
				PackageName: "templates",
			},
			setupFiles:  map[string]string{}, // Empty directory
			expectError: false,
			description: "Should handle empty directory gracefully",
		},
		{
			name: "Directory with syntax errors",
			config: types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    "",
				PackageName: "templates",
			},
			setupFiles: map[string]string{
				"broken_templ.go": `package templates
// This file has syntax errors
func BrokenFunction( {
	return "broken"
}`,
			},
			expectError: false, // Should skip files with errors
			description: "Should handle files with syntax errors",
		},
		{
			name: "Valid templates",
			config: types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    "",
				PackageName: "templates",
			},
			setupFiles: map[string]string{
				"page_templ.go": `package templates

import "github.com/a-h/templ"

func HomePage() templ.Component {
	return templ.Raw("Home")
}`,
				"about_templ.go": `package templates

import "github.com/a-h/templ"

func AboutPage() templ.Component {
	return templ.Raw("About")
}`,
			},
			expectError: false,
			description: "Should process valid template files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tempDir string
			if tt.setupFiles != nil {
				tempDir = t.TempDir()
				
				// Create files
				for filename, content := range tt.setupFiles {
					fullPath := filepath.Join(tempDir, filename)
					err := os.WriteFile(fullPath, []byte(content), 0644)
					if err != nil {
						t.Fatalf("Failed to create file %s: %v", filename, err)
					}
				}
				
				tt.config.ScanPath = tempDir
			}

			// Scan templates
			templates, validationErrors, err := ScanTemplatesWithPackages(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else {
					t.Logf("✅ %s - Got expected error: %v", tt.description, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			t.Logf("✅ %s", tt.description)
			t.Logf("   Found %d templates", len(templates))
			t.Logf("   Validation errors: %d", len(validationErrors))

			// For valid templates test, verify we found the expected templates
			if len(tt.setupFiles) > 0 && strings.Contains(tt.name, "Valid templates") {
				if len(templates) == 0 {
					t.Error("Expected to find templates but found none")
				}
			}
		})
	}
}

// TestScanSpecificPackage tests package-specific scanning
func TestScanSpecificPackage(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create files in different packages
	files := map[string]string{
		"admin/users_templ.go": `package admin

import "github.com/a-h/templ"

func UsersPage() templ.Component {
	return templ.Raw("Users")
}`,
		"admin/settings_templ.go": `package admin

import "github.com/a-h/templ"

func SettingsPage() templ.Component {
	return templ.Raw("Settings")
}`,
		"public/home_templ.go": `package public

import "github.com/a-h/templ"

func HomePage() templ.Component {
	return templ.Raw("Home")
}`,
		"api/docs_templ.go": `package api

import "github.com/a-h/templ"

func DocsPage() templ.Component {
	return templ.Raw("Docs")
}`,
	}
	
	// Create files
	for filename, content := range files {
		fullPath := filepath.Join(tempDir, filename)
		dir := filepath.Dir(fullPath)
		
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}

	tests := []struct {
		name            string
		targetPackage   string
		expectedCount   int
		expectedFuncs   []string
		description     string
	}{
		{
			name:            "Admin package",
			targetPackage:   "admin",
			expectedCount:   2,
			expectedFuncs:   []string{"UsersPage", "SettingsPage"},
			description:     "Should find only admin package templates",
		},
		{
			name:            "Public package",
			targetPackage:   "public",
			expectedCount:   1,
			expectedFuncs:   []string{"HomePage"},
			description:     "Should find only public package templates",
		},
		{
			name:            "API package",
			targetPackage:   "api",
			expectedCount:   1,
			expectedFuncs:   []string{"DocsPage"},
			description:     "Should find only API package templates",
		},
		{
			name:            "Non-existent package",
			targetPackage:   "nonexistent",
			expectedCount:   0,
			expectedFuncs:   []string{},
			description:     "Should find no templates for non-existent package",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.Config{
				ModuleName:  "github.com/test/project",
				ScanPath:    tempDir,
				PackageName: tt.targetPackage,
			}

			templates, err := ScanSpecificPackage(config)
			if err != nil {
				t.Fatalf("ScanSpecificPackage failed: %v", err)
			}

			if len(templates) != tt.expectedCount {
				t.Errorf("Expected %d templates, got %d", tt.expectedCount, len(templates))
			}

			// Check expected functions
			foundFuncs := make(map[string]bool)
			for _, template := range templates {
				foundFuncs[template.FunctionName] = true
			}

			for _, expectedFunc := range tt.expectedFuncs {
				if !foundFuncs[expectedFunc] {
					t.Errorf("Expected function %s not found", expectedFunc)
				}
			}

			// Verify all templates belong to the target package
			for _, template := range templates {
				if template.PackageName != tt.targetPackage {
					t.Errorf("Found template from wrong package: %s (expected %s)", template.PackageName, tt.targetPackage)
				}
			}

			t.Logf("✅ %s", tt.description)
			t.Logf("   Found %d templates in package %s", len(templates), tt.targetPackage)
		})
	}
}

// TestTemplateInfoGeneration tests the template info structure generation
func TestTemplateInfoGeneration(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_templ.go")
	
	goContent := `package templates

import "github.com/a-h/templ"

func HomePage() templ.Component {
	return templ.Raw("Home")
}

func UserProfile(userID string) templ.Component {
	return templ.Raw("Profile")
}

func Page(data interface{}) templ.Component {
	return templ.Raw("Page")
}`

	err := os.WriteFile(tempFile, []byte(goContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	config := types.Config{
		ModuleName:  "github.com/company/app",
		ScanPath:    tempDir,
		PackageName: "templates",
	}

	templates, _, err := ScanTemplatesWithPackages(config)
	if err != nil {
		t.Fatalf("ScanTemplatesWithPackages failed: %v", err)
	}

	// Should find HomePage and Page (UserProfile skipped due to parameters)
	expectedFunctions := map[string]bool{
		"HomePage": false,
		"Page":     false,
	}

	for _, template := range templates {
		if _, expected := expectedFunctions[template.FunctionName]; expected {
			expectedFunctions[template.FunctionName] = true

			// Validate template info structure
			if template.FilePath == "" {
				t.Errorf("Template %s has empty FilePath", template.FunctionName)
			}
			if template.PackageName == "" {
				t.Errorf("Template %s has empty PackageName", template.FunctionName)
			}
			if template.ImportPath == "" {
				t.Errorf("Template %s has empty ImportPath", template.FunctionName)
			}
			if template.TemplateKey == "" {
				t.Errorf("Template %s has empty TemplateKey", template.FunctionName)
			}
			if template.HumanName == "" {
				t.Errorf("Template %s has empty HumanName", template.FunctionName)
			}

			// Validate UUID format for template key
			if len(template.TemplateKey) != 36 {
				t.Errorf("Template %s has invalid UUID format: %s", template.FunctionName, template.TemplateKey)
			}

			t.Logf("Template %s:", template.FunctionName)
			t.Logf("  FilePath: %s", template.FilePath)
			t.Logf("  PackageName: %s", template.PackageName)
			t.Logf("  ImportPath: %s", template.ImportPath)
			t.Logf("  TemplateKey: %s", template.TemplateKey)
			t.Logf("  HumanName: %s", template.HumanName)
		}
	}

	// Check all expected functions were found
	for funcName, found := range expectedFunctions {
		if !found {
			t.Errorf("Expected function %s not found in templates", funcName)
		}
	}

	t.Logf("✅ Template info generation working correctly")
	t.Logf("   Generated %d valid template infos", len(templates))
}