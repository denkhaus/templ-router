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
		"TRGEN_SCAN_PATH":   "app",
		"TRGEN_OUTPUT_DIR":  "generated/templates",
		"TRGEN_MODULE_NAME": "github.com/denkhaus/templ-router/demo",
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
		"TRGEN_SCAN_PATH":   "app",
		"TRGEN_OUTPUT_DIR":  "generated/templates",
		"TRGEN_MODULE_NAME": "github.com/denkhaus/templ-router/demo",
	}, "sh", "-c", "cd demo && trgen --watch")
}

// All builds all binaries for different platforms
func (Build) All() error {
	fmt.Println("Building all binaries...")

	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	for _, platform := range platforms {
		fmt.Printf("Building for %s/%s...\n", platform.os, platform.arch)

		env := map[string]string{
			"GOOS":   platform.os,
			"GOARCH": platform.arch,
		}

		output := fmt.Sprintf("bin/templ-router-%s-%s", platform.os, platform.arch)
		if platform.os == "windows" {
			output += ".exe"
		}

		if err := sh.RunWithV(env, "go", "build", "-o", output, "./cmd/trgen"); err != nil {
			return fmt.Errorf("failed to build for %s/%s: %w", platform.os, platform.arch, err)
		}
	}

	fmt.Println("All binaries built successfully")
	return nil
}
