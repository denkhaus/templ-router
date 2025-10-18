package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
)

// ValidateTemplatePath validates that the template file is in the correct location
func ValidateTemplatePath(filePath string, config types.Config) error {
	rootDir := config.ScanPath
	// Basic validation - ensure it's a _templ.go file in the configured root directory
	if !strings.Contains(filePath, "/"+rootDir+"/") && !strings.HasSuffix(filepath.Dir(filePath), "/"+rootDir) {
		return fmt.Errorf("template file %s is not in the %s directory", filePath, rootDir)
	}
	return nil
}

// ValidateFunctionNaming validates that template functions follow naming conventions
func ValidateFunctionNaming(functionName, filePath string) error {
	// Extract the template file name without extension
	fileName := filepath.Base(filePath)
	templateName := strings.TrimSuffix(fileName, "_templ.go")

	// Generic validation rules that work for any app structure:

	// Rule 1: Reserved function names must be in correctly named files
	switch functionName {
	case "Page":
		// Page functions must be in page.templ files
		if templateName != "page" {
			return fmt.Errorf("function 'Page' found in '%s.templ' but should only be in 'page.templ'", templateName)
		}
		return nil
	case "Layout":
		// Layout functions must be in layout.templ files
		if templateName != "layout" {
			return fmt.Errorf("function 'Layout' found in '%s.templ' but should only be in 'layout.templ'", templateName)
		}
		return nil
	case "Error":
		// Error functions must be in error.templ files
		if templateName != "error" {
			return fmt.Errorf("function 'Error' found in '%s.templ' but should only be in 'error.templ'", templateName)
		}
		return nil
	}

	// Rule 2: Components cannot have "Page" suffix (reserved for actual pages)
	if strings.HasSuffix(functionName, "Page") && functionName != "Page" {
		return fmt.Errorf("function '%s' has 'Page' suffix but is not in 'page.templ' - components cannot use 'Page' suffix", functionName)
	}

	// Rule 3: Check for parameter conflicts in YAML metadata
	if err := ValidateYamlParameters(filePath); err != nil {
		return err
	}

	// Rule 2: Specific template files must have correct function names
	switch templateName {
	case "page":
		// page.templ files must have Page function
		if functionName != "Page" {
			return fmt.Errorf("function '%s' found in 'page.templ' but should be 'Page'", functionName)
		}
	case "layout":
		// layout.templ files can have Layout function plus other components
		// No strict validation - layouts are flexible
	case "error":
		// error.templ files must have Error function
		if functionName != "Error" {
			return fmt.Errorf("function '%s' found in 'error.templ' but should be 'Error'", functionName)
		}
	default:
		// All other files are considered component files and can have any function names
		// This makes the generator agnostic to specific app structures
		return nil
	}

	return nil
}

// ValidateYamlParameters checks for parameter conflicts in YAML metadata files
func ValidateYamlParameters(filePath string) error {
	// Get the corresponding YAML file path
	yamlPath := strings.TrimSuffix(filePath, "_templ.go") + ".templ.yaml"

	// Check if YAML file exists
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		// No YAML file, no validation needed
		return nil
	}

	// Read YAML file
	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		// Can't read YAML, skip validation
		return nil
	}

	// Check for parameter conflicts based on directory structure
	dir := filepath.Dir(filePath)

	// Check if we're in a locale_ subdirectory and defining locale parameter
	if strings.Contains(dir, "/locale_/") && strings.Contains(string(yamlContent), "locale:") {
		// Check if it's in the dynamic parameters section
		if strings.Contains(string(yamlContent), "dynamic:") &&
			strings.Contains(string(yamlContent), "parameters:") &&
			strings.Contains(string(yamlContent), "locale:") {
			return fmt.Errorf("YAML parameter conflict: 'locale' parameter defined in '%s' but locale is already inherited from parent 'locale_/' directory", yamlPath)
		}
	}

	// Add more parameter conflict checks here for other inherited parameters
	// For example, if we had user_/ directories, we'd check for user parameter conflicts

	return nil
}

// ValidateConfig validates the configuration
func ValidateConfig(config types.Config) error {
	if config.ScanPath == "" {
		return fmt.Errorf("scan path cannot be empty")
	}
	
	if config.OutputDir == "" {
		return fmt.Errorf("output directory cannot be empty")
	}
	
	if config.ModuleName == "" {
		return fmt.Errorf("module name cannot be empty")
	}
	
	if config.PackageName == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	
	// Validate package name format
	if strings.Contains(config.PackageName, "-") {
		return fmt.Errorf("package name contains invalid characters (hyphens not allowed)")
	}
	
	if len(config.PackageName) > 0 && config.PackageName[0] >= '0' && config.PackageName[0] <= '9' {
		return fmt.Errorf("package name cannot start with a number")
	}
	
	// Validate module name format
	if strings.Contains(config.ModuleName, " ") {
		return fmt.Errorf("module name contains invalid characters (spaces not allowed)")
	}
	
	return nil
}

// ValidateTemplates validates a slice of template info
func ValidateTemplates(templates []types.TemplateInfo) error {
	if len(templates) == 0 {
		return nil // Empty templates are allowed
	}
	
	// Check for duplicate template keys
	templateKeys := make(map[string]bool)
	for _, tmpl := range templates {
		if tmpl.FunctionName == "" {
			return fmt.Errorf("function name cannot be empty")
		}
		
		if tmpl.TemplateKey == "" {
			return fmt.Errorf("template key cannot be empty")
		}
		
		if templateKeys[tmpl.TemplateKey] {
			return fmt.Errorf("duplicate template key: %s", tmpl.TemplateKey)
		}
		templateKeys[tmpl.TemplateKey] = true
	}
	
	// Check for duplicate route patterns
	routePatterns := make(map[string]bool)
	for _, tmpl := range templates {
		if routePatterns[tmpl.RoutePattern] {
			return fmt.Errorf("duplicate route pattern: %s", tmpl.RoutePattern)
		}
		routePatterns[tmpl.RoutePattern] = true
	}
	
	// Check for invalid package aliases
	for _, tmpl := range templates {
		if strings.Contains(tmpl.PackageAlias, "-") {
			return fmt.Errorf("invalid package alias: %s (contains hyphens)", tmpl.PackageAlias)
		}
	}
	
	return nil
}
