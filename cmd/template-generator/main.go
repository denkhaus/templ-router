package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Get configurable scan path from environment variable or use default
	scanPath := os.Getenv("TEMPLATE_SCAN_PATH")
	if scanPath == "" {
		scanPath = "app" // Default fallback
	}

	// Get configurable module name from environment variable or use default
	moduleName := os.Getenv("TEMPLATE_MODULE_NAME")
	if moduleName == "" {
		moduleName = "github.com/denkhaus/templ-router" // Default fallback
	}

	// Configuration for template scanning
	config := Config{
		ModuleName:  moduleName,
		ScanPath:    scanPath, // Now configurable via TEMPLATE_SCAN_PATH environment variable
		OutputDir:   "generated/templates",
		PackageName: "templates",
	}

	fmt.Println("Starting template generation...")
	fmt.Printf("Scanning path: %s\n", config.ScanPath)
	fmt.Printf("Output directory: %s\n", config.OutputDir)

	// Scan for templates
	templates, validationErrors, err := scanTemplatesWithPackages(config)
	if err != nil {
		log.Fatalf("Failed to scan templates: %v", err)
	}

	// Print validation errors as warnings
	if len(validationErrors) > 0 {
		fmt.Printf("\nValidation warnings:\n")
		for _, warning := range validationErrors {
			fmt.Printf("  WARNING: %s\n", warning)
		}
		fmt.Println()
	}

	fmt.Printf("Found %d templates\n", len(templates))
	for _, template := range templates {
		fmt.Printf("  %s -> %s (%s)\n", template.RoutePattern, template.FunctionName, template.HumanName)
	}

	// Generate registry
	// Validate required parameters
	if config.ModuleName == "" {
		log.Fatalf("TEMPLATE_MODULE_NAME is required and cannot be empty")
	}

	if err := generateRegistry(config, templates); err != nil {
		log.Fatalf("Failed to generate registry: %v", err)
	}

	fmt.Println("Template generation completed successfully!")
}
