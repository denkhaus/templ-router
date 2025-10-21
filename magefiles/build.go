package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build namespace for build-related commands
type Build mg.Namespace

// TailwindClean builds Tailwind CSS
func (p Build) TailwindClean() error {
	fmt.Println("Building Tailwind CSS...")
	return sh.RunV("npx", "@tailwindcss/cli", "-i", "./demo/assets/css/input.css", "-o", "./demo/assets/css/output.css")
}

// TailwindWatch builds Tailwind CSS in watch mode
func (p Build) TailwindWatch() error {
	fmt.Println("Watching Tailwind CSS...")
	return sh.RunV("npx", "@tailwindcss/cli", "-i", "./demo/assets/css/input.css", "-o", "./demo/assets/css/output.css", "--watch")
}

// TemplWatch runs templ generation in watch mode
func (p Build) TemplWatch() error {
	fmt.Println("Watching Templ files...")
	return sh.RunV("templ", "generate", "--watch", "--proxy=http://localhost:8090", "--open-browser=false")
}

// TemplGenerate generates Templ templates
func (p Build) TemplGenerate() error {
	mg.Deps(Templ.Install)

	fmt.Println("Generating Templ templates...")
	return sh.RunV("templ", "generate")
}

func (p Build) RegistryGenerate() error {

	mg.Deps(p.TemplGenerate, Generator.Install)

	fmt.Println("Generating template registry...")
	// Generate templates
	if err := sh.RunWithV(map[string]string{
		"TEMPLATE_SCAN_PATH":   "app",
		"TEMPLATE_OUTPUT_DIR":  "generated/templates",
		"TEMPLATE_MODULE_NAME": "github.com/denkhaus/templ-router/demo",
	}, "sh", "-c", "cd demo && trgen"); err != nil {
		return err
	}

	// Run go mod tidy to ensure generated packages are recognized
	fmt.Println("Running go mod tidy to register generated packages...")
	return sh.RunV("sh", "-c", "cd demo && go mod tidy")
}

func (Build) RegistryWatch() error {
	fmt.Println("Watching template registry...")

	return sh.RunWithV(map[string]string{
		"TEMPLATE_SCAN_PATH":   "app",
		"TEMPLATE_OUTPUT_DIR":  "generated/templates",
		"TEMPLATE_MODULE_NAME": "github.com/denkhaus/templ-router/demo",
	}, "sh", "-c", "cd demo && trgen --watch")
}
