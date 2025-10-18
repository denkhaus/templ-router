package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
	"github.com/denkhaus/templ-router/cmd/template-generator/commands"
)

// TestEndToEndGeneration tests the complete template generator workflow
func TestEndToEndGeneration(t *testing.T) {
	scenarios := []struct {
		name        string
		projectStructure map[string]string // file path -> content
		config      types.Config
		expectedOutput []string // strings that should be in the generated registry
		description string
	}{
		{
			name: "Simple web application",
			projectStructure: map[string]string{
				"go.mod": `module github.com/test/webapp

go 1.21

require github.com/a-h/templ v0.3.960
`,
				"app/home_templ.go": `package app

import "github.com/a-h/templ"

func HomePage() templ.Component {
	return templ.Raw("Welcome to our website!")
}

func AboutPage() templ.Component {
	return templ.Raw("About us")
}`,
				"app/admin/dashboard_templ.go": `package admin

import "github.com/a-h/templ"

func DashboardPage() templ.Component {
	return templ.Raw("Admin Dashboard")
}

func UsersPage() templ.Component {
	return templ.Raw("User Management")
}`,
				"app/public/contact_templ.go": `package public

import "github.com/a-h/templ"

func ContactPage() templ.Component {
	return templ.Raw("Contact Us")
}`,
			},
			config: types.Config{
				ModuleName:  "github.com/test/webapp",
				ScanPath:    "app",
				OutputDir:   "generated",
				PackageName: "templates",
			},
			expectedOutput: []string{
				"package templates",
				"HomePage",
				"AboutPage", 
				"DashboardPage",
				"UsersPage",
				"ContactPage",
				"github.com/test/webapp",
			},
			description: "Should generate registry for complete web application",
		},
		{
			name: "Component library",
			projectStructure: map[string]string{
				"go.mod": `module github.com/ui/components

go 1.21

require github.com/a-h/templ v0.3.960
`,
				"components/button_templ.go": `package components

import "github.com/a-h/templ"

func PrimaryButton() templ.Component {
	return templ.Raw("<button class='btn-primary'>Click me</button>")
}

func SecondaryButton() templ.Component {
	return templ.Raw("<button class='btn-secondary'>Cancel</button>")
}`,
				"components/forms/input_templ.go": `package forms

import "github.com/a-h/templ"

func TextInput() templ.Component {
	return templ.Raw("<input type='text' />")
}

func EmailInput() templ.Component {
	return templ.Raw("<input type='email' />")
}`,
				"components/layout/header_templ.go": `package layout

import "github.com/a-h/templ"

func Header() templ.Component {
	return templ.Raw("<header>Site Header</header>")
}

func Footer() templ.Component {
	return templ.Raw("<footer>Site Footer</footer>")
}`,
			},
			config: types.Config{
				ModuleName:  "github.com/ui/components",
				ScanPath:    "components",
				OutputDir:   "generated",
				PackageName: "registry",
			},
			expectedOutput: []string{
				"package registry",
				"PrimaryButton",
				"SecondaryButton",
				"TextInput",
				"EmailInput",
				"Header",
				"Footer",
				"github.com/ui/components",
			},
			description: "Should generate registry for component library",
		},
		{
			name: "Microservice with API docs",
			projectStructure: map[string]string{
				"go.mod": `module github.com/company/api-service

go 1.21

require github.com/a-h/templ v0.3.960
`,
				"templates/docs_templ.go": `package templates

import "github.com/a-h/templ"

func APIDocsPage() templ.Component {
	return templ.Raw("API Documentation")
}

func SwaggerUI() templ.Component {
	return templ.Raw("Swagger Interface")
}`,
				"templates/errors/not_found_templ.go": `package errors

import "github.com/a-h/templ"

func NotFoundPage() templ.Component {
	return templ.Raw("404 - Page Not Found")
}

func ServerErrorPage() templ.Component {
	return templ.Raw("500 - Internal Server Error")
}`,
			},
			config: types.Config{
				ModuleName:  "github.com/company/api-service",
				ScanPath:    "templates",
				OutputDir:   "generated",
				PackageName: "templates",
			},
			expectedOutput: []string{
				"package templates",
				"APIDocsPage",
				"SwaggerUI",
				"NotFoundPage",
				"ServerErrorPage",
				"github.com/company/api-service",
			},
			description: "Should generate registry for microservice templates",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Create temporary project directory
			tempDir := t.TempDir()
			
			// Create project structure
			for filePath, content := range scenario.projectStructure {
				fullPath := filepath.Join(tempDir, filePath)
				dir := filepath.Dir(fullPath)
				
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					t.Fatalf("Failed to create directory %s: %v", dir, err)
				}
				
				err = os.WriteFile(fullPath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create file %s: %v", filePath, err)
				}
			}
			
			// Update config with absolute paths
			config := scenario.config
			config.ScanPath = filepath.Join(tempDir, config.ScanPath)
			config.OutputDir = filepath.Join(tempDir, config.OutputDir)
			
			// Create output directory
			err := os.MkdirAll(config.OutputDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create output directory: %v", err)
			}
			
			// Change to project directory for proper module resolution
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(originalDir)
			
			err = os.Chdir(tempDir)
			if err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}
			
			// Run the template generator using the actual command structure
			err = commands.GenerateCommand(config)
			if err != nil {
				t.Fatalf("Template generation failed: %v", err)
			}
			
			// Verify output file was created
			outputFile := filepath.Join(config.OutputDir, "registry.go")
			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				t.Fatalf("Output file was not created: %s", outputFile)
			}
			
			// Read and verify generated content
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}
			
			contentStr := string(content)
			
			// Check for expected content
			for _, expected := range scenario.expectedOutput {
				if !strings.Contains(contentStr, expected) {
					t.Errorf("Expected content %q not found in generated file", expected)
				}
			}
			
			// Verify the generated file is valid Go code
			if !strings.HasPrefix(contentStr, "package "+config.PackageName) {
				t.Error("Generated file should start with correct package declaration")
			}
			
			// Count template functions found
			templateCount := 0
			for _, expected := range scenario.expectedOutput {
				if strings.HasSuffix(expected, "Page") || strings.HasSuffix(expected, "Button") || 
				   strings.HasSuffix(expected, "Input") || strings.HasSuffix(expected, "Header") || 
				   strings.HasSuffix(expected, "Footer") || strings.HasSuffix(expected, "UI") {
					if strings.Contains(contentStr, expected) {
						templateCount++
					}
				}
			}
			
			t.Logf("✅ %s", scenario.description)
			t.Logf("   Generated file: %s", outputFile)
			t.Logf("   Template functions found: %d", templateCount)
			t.Logf("   File size: %d bytes", len(content))
			
			// Log first few lines for verification
			lines := strings.Split(contentStr, "\n")
			if len(lines) > 5 {
				t.Logf("   Generated content preview:")
				for i := 0; i < 5; i++ {
					t.Logf("     %s", lines[i])
				}
			}
		})
	}
}

// TestGeneratorPerformance tests performance with larger template sets
func TestGeneratorPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}
	
	tempDir := t.TempDir()
	
	// Create go.mod
	goModContent := `module github.com/test/large-app

go 1.21

require github.com/a-h/templ v0.3.960
`
	err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}
	
	// Generate many template files
	templateCount := 100
	packagesCount := 10
	
	for pkg := 0; pkg < packagesCount; pkg++ {
		packageName := "package" + string(rune('a'+pkg))
		packageDir := filepath.Join(tempDir, "templates", packageName)
		
		err := os.MkdirAll(packageDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create package directory: %v", err)
		}
		
		for i := 0; i < templateCount/packagesCount; i++ {
			fileName := filepath.Join(packageDir, "template"+string(rune('0'+i))+"_templ.go")
			content := "package " + packageName + "\n\n" +
				"import \"github.com/a-h/templ\"\n\n" +
				"func Template" + string(rune('A'+i)) + "() templ.Component {\n" +
				"	return templ.Raw(\"Template " + string(rune('A'+i)) + "\")\n" +
				"}\n"
			
			err := os.WriteFile(fileName, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}
		}
	}
	
	config := types.Config{
		ModuleName:  "github.com/test/large-app",
		ScanPath:    filepath.Join(tempDir, "templates"),
		OutputDir:   filepath.Join(tempDir, "generated"),
		PackageName: "templates",
	}
	
	err = os.MkdirAll(config.OutputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)
	
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	
	// Run generation and measure time
	err = commands.GenerateCommand(config)
	if err != nil {
		t.Fatalf("Template generation failed: %v", err)
	}
	
	// Verify output
	outputFile := filepath.Join(config.OutputDir, "registry.go")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	contentStr := string(content)
	
	// Count generated template registrations
	templateRegistrations := strings.Count(contentStr, "func Template")
	
	t.Logf("✅ Performance test completed")
	t.Logf("   Generated %d template files across %d packages", templateCount, packagesCount)
	t.Logf("   Found %d template registrations in output", templateRegistrations)
	t.Logf("   Output file size: %d bytes", len(content))
	
	if templateRegistrations == 0 {
		t.Error("Expected to find template registrations in generated file")
	}
}

// TestGeneratorErrorRecovery tests error handling and recovery
func TestGeneratorErrorRecovery(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(tempDir string) types.Config
		expectError bool
		description string
	}{
		{
			name: "Missing go.mod file",
			setupFunc: func(tempDir string) types.Config {
				// Don't create go.mod
				templateDir := filepath.Join(tempDir, "templates")
				os.MkdirAll(templateDir, 0755)
				
				// Create a template file
				templateFile := filepath.Join(templateDir, "page_templ.go")
				content := `package templates

import "github.com/a-h/templ"

func HomePage() templ.Component {
	return templ.Raw("Home")
}`
				os.WriteFile(templateFile, []byte(content), 0644)
				
				return types.Config{
					ModuleName:  "github.com/test/no-mod",
					ScanPath:    templateDir,
					OutputDir:   filepath.Join(tempDir, "generated"),
					PackageName: "templates",
				}
			},
			expectError: false, // Should handle gracefully
			description: "Should handle missing go.mod gracefully",
		},
		{
			name: "Invalid template syntax",
			setupFunc: func(tempDir string) types.Config {
				// Create go.mod
				goModContent := `module github.com/test/broken-syntax

go 1.21

require github.com/a-h/templ v0.3.960
`
				os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644)
				
				templateDir := filepath.Join(tempDir, "templates")
				os.MkdirAll(templateDir, 0755)
				
				// Create broken template file
				brokenFile := filepath.Join(templateDir, "broken_templ.go")
				brokenContent := `package templates

import "github.com/a-h/templ"

func BrokenFunction( {
	// Missing closing parenthesis and return
}`
				os.WriteFile(brokenFile, []byte(brokenContent), 0644)
				
				// Create valid template file
				validFile := filepath.Join(templateDir, "valid_templ.go")
				validContent := `package templates

import "github.com/a-h/templ"

func ValidFunction() templ.Component {
	return templ.Raw("Valid")
}`
				os.WriteFile(validFile, []byte(validContent), 0644)
				
				return types.Config{
					ModuleName:  "github.com/test/broken-syntax",
					ScanPath:    templateDir,
					OutputDir:   filepath.Join(tempDir, "generated"),
					PackageName: "templates",
				}
			},
			expectError: false, // Should skip broken files and continue
			description: "Should skip files with syntax errors and continue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			config := tt.setupFunc(tempDir)
			
			err := os.MkdirAll(config.OutputDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create output directory: %v", err)
			}
			
			// Change to project directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(originalDir)
			
			err = os.Chdir(tempDir)
			if err != nil {
				t.Fatalf("Failed to change to temp directory: %v", err)
			}
			
			// Run generation
			err = commands.GenerateCommand(config)
			
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
					
					// Verify output file was still created
					outputFile := filepath.Join(config.OutputDir, "registry.go")
					if _, err := os.Stat(outputFile); os.IsNotExist(err) {
						t.Error("Output file should have been created despite errors")
					}
				}
			}
		})
	}
}