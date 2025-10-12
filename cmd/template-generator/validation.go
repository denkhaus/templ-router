package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// validateFunctionNaming validates that template functions follow naming conventions
func validateFunctionNaming(functionName, filePath string) error {
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
	if err := validateYamlParameters(filePath); err != nil {
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

// validateYamlParameters checks for parameter conflicts in YAML metadata files
func validateYamlParameters(filePath string) error {
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