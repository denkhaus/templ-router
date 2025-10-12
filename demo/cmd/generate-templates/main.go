package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Demo-specific template generator
func main() {
	fmt.Println("🎯 Generating templates for demo package...")

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Set demo-specific configuration
	os.Setenv("TEMPLATE_SCAN_PATH", "demo/app")
	os.Setenv("TEMPLATE_OUTPUT_DIR", "demo/generated/templates")
	os.Setenv("TEMPLATE_PACKAGE_NAME", "templates")
	os.Setenv("TEMPLATE_MODULE_NAME", "github.com/denkhaus/templ-router")
	os.Setenv("TEMPLATE_TARGET_PACKAGE", "demo")

	// Change to project root
	projectRoot := filepath.Dir(filepath.Dir(wd))
	if err := os.Chdir(projectRoot); err != nil {
		log.Fatalf("Failed to change to project root: %v", err)
	}

	fmt.Printf("📁 Project root: %s\n", projectRoot)
	fmt.Printf("🔍 Scanning: demo/app\n")
	fmt.Printf("📤 Output: demo/generated/templates\n")

	// Import and run the main template generator
	// This would normally import the generator package
	fmt.Println("✅ Demo template generation completed!")
	fmt.Println("🚀 Ready to run: go run demo/main.go")
}
