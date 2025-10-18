package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Generator contains tasks for the template generator
type Generator mg.Namespace

// Install builds and installs the template generator with proper versioning
func (p Generator) Install() error {
	fmt.Println("🔧 Installing trgen...")

	// Get version information
	version, err := getVersion()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	gitCommit, err := getGitCommit()
	if err != nil {
		fmt.Printf("Warning: Could not get git commit: %v\n", err)
		gitCommit = "unknown"
	}

	buildTime := time.Now().UTC().Format(time.RFC3339)

	// Build with version information
	ldflags := fmt.Sprintf("-X github.com/denkhaus/templ-router/cmd/trgen/version.Version=%s "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.GitCommit=%s "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.BuildTime=%s",
		version, gitCommit, buildTime)

	fmt.Printf("📦 Building template-generator v%s (commit: %s)\n", version, gitCommit)

	// Install the binary
	err = sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
	}, "go", "install", "-ldflags", ldflags, "./cmd/trgen")

	if err != nil {
		return fmt.Errorf("failed to install template generator: %w", err)
	}

	// Verify installation
	fmt.Println("✅ Verifying installation...")

	// Get the installed binary path
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	binaryName := "trgen"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	binaryPath := filepath.Join(gopath, "bin", binaryName)

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found at %s", binaryPath)
	}

	// Test the binary
	output, err := sh.Output(binaryPath, "--version")
	if err != nil {
		return fmt.Errorf("failed to run installed binary: %w", err)
	}

	fmt.Printf("🎉 Successfully installed: %s\n", strings.TrimSpace(output))
	fmt.Printf("📍 Binary location: %s\n", binaryPath)

	return nil
}

// Build builds the template generator without installing
func (p Generator) Build() error {
	fmt.Println("🔧 Building trgen...")

	// Get version information
	version, err := getVersion()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	gitCommit, err := getGitCommit()
	if err != nil {
		fmt.Printf("Warning: Could not get git commit: %v\n", err)
		gitCommit = "unknown"
	}

	buildTime := time.Now().UTC().Format(time.RFC3339)

	// Build with version information
	ldflags := fmt.Sprintf("-X github.com/denkhaus/templ-router/cmd/trgen/version.Version=%s "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.GitCommit=%s "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.BuildTime=%s",
		version, gitCommit, buildTime)

	fmt.Printf("📦 Building template-generator v%s (commit: %s)\n", version, gitCommit)

	// Create bin directory if it doesn't exist
	if err := os.MkdirAll("bin", 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Build the binary
	binaryName := "trgen"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	outputPath := filepath.Join("bin", binaryName)

	err = sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
	}, "go", "build", "-ldflags", ldflags, "-o", outputPath, "./cmd/trgen")

	if err != nil {
		return fmt.Errorf("failed to build template generator: %w", err)
	}

	// Test the binary
	output, err := sh.Output(outputPath, "--version")
	if err != nil {
		return fmt.Errorf("failed to run built binary: %w", err)
	}

	fmt.Printf("✅ Successfully built: %s\n", strings.TrimSpace(output))
	fmt.Printf("📍 Binary location: %s\n", outputPath)

	return nil
}

// Release builds the template generator for multiple platforms
func (p Generator) Release() error {
	fmt.Println("🚀 Building trgen for release...")

	// Get version information
	version, err := getVersion()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	gitCommit, err := getGitCommit()
	if err != nil {
		fmt.Printf("Warning: Could not get git commit: %v\n", err)
		gitCommit = "unknown"
	}

	buildTime := time.Now().UTC().Format(time.RFC3339)

	// Build with version information
	ldflags := fmt.Sprintf("-X github.com/denkhaus/templ-router/cmd/trgen/version.Version=%s "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.GitCommit=%s "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.BuildTime=%s",
		version, gitCommit, buildTime)

	fmt.Printf("📦 Building trgen v%s (commit: %s) for multiple platforms\n", version, gitCommit)

	// Create release directory
	releaseDir := fmt.Sprintf("release/trgen-v%s", version)
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return fmt.Errorf("failed to create release directory: %w", err)
	}

	// Define target platforms
	platforms := []struct {
		goos   string
		goarch string
		suffix string
	}{
		{"linux", "amd64", "linux-amd64"},
		{"linux", "arm64", "linux-arm64"},
		{"darwin", "amd64", "darwin-amd64"},
		{"darwin", "arm64", "darwin-arm64"},
		{"windows", "amd64", "windows-amd64.exe"},
		{"windows", "arm64", "windows-arm64.exe"},
	}

	// Build for each platform
	for _, platform := range platforms {
		fmt.Printf("🔨 Building for %s/%s...\n", platform.goos, platform.goarch)

		binaryName := fmt.Sprintf("trgen-%s", platform.suffix)
		outputPath := filepath.Join(releaseDir, binaryName)

		env := map[string]string{
			"GOOS":        platform.goos,
			"GOARCH":      platform.goarch,
			"CGO_ENABLED": "0",
		}

		err = sh.RunWith(env, "go", "build", "-ldflags", ldflags, "-o", outputPath, "./cmd/trgen")
		if err != nil {
			return fmt.Errorf("failed to build for %s/%s: %w", platform.goos, platform.goarch, err)
		}

		fmt.Printf("✅ Built %s\n", binaryName)
	}

	fmt.Printf("🎉 Release build complete! Binaries available in: %s\n", releaseDir)

	return nil
}

// Test runs all tests for the template generator
func (p Generator) Test() error {
	fmt.Println("🧪 Running trgen tests...")

	return sh.Run("go", "test", "-v", "./cmd/trgen/...")
}

// TestCoverage runs tests with coverage for the template generator
func (p Generator) TestCoverage() error {
	fmt.Println("🧪 Running trgen tests with coverage...")

	// Create coverage directory
	if err := os.MkdirAll("coverage", 0755); err != nil {
		return fmt.Errorf("failed to create coverage directory: %w", err)
	}

	// Run tests with coverage
	err := sh.Run("go", "test", "-v", "-coverprofile=coverage/generator.out", "./cmd/trgen/...")
	if err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	// Generate HTML coverage report
	err = sh.Run("go", "tool", "cover", "-html=coverage/generator.out", "-o", "coverage/generator.html")
	if err != nil {
		return fmt.Errorf("failed to generate HTML coverage report: %w", err)
	}

	// Show coverage summary
	err = sh.Run("go", "tool", "cover", "-func=coverage/generator.out")
	if err != nil {
		return fmt.Errorf("failed to show coverage summary: %w", err)
	}

	fmt.Println("📊 Coverage report generated: coverage/generator.html")

	return nil
}

// Clean removes build artifacts
func (p Generator) Clean() error {
	fmt.Println("🧹 Cleaning trgen build artifacts...")

	// Remove directories
	dirs := []string{"bin", "release", "coverage"}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("Warning: failed to remove %s: %v\n", dir, err)
		} else {
			fmt.Printf("🗑️  Removed %s\n", dir)
		}
	}

	fmt.Println("✅ Clean complete!")

	return nil
}

// Version shows the current version information
func (p Generator) Version() error {
	version, err := getVersion()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	gitCommit, err := getGitCommit()
	if err != nil {
		gitCommit = "unknown"
	}

	fmt.Printf("trgen Version: v%s\n", version)
	fmt.Printf("Git Commit: %s\n", gitCommit)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	return nil
}

// Helper functions

func getVersion() (string, error) {
	// Try to get version from git tag
	output, err := sh.Output("git", "describe", "--tags", "--abbrev=0")
	if err == nil && strings.HasPrefix(output, "v") {
		return strings.TrimPrefix(strings.TrimSpace(output), "v"), nil
	}

	// Try to get version from git describe
	output, err = sh.Output("git", "describe", "--tags", "--always", "--dirty")
	if err == nil {
		version := strings.TrimSpace(output)
		if strings.HasPrefix(version, "v") {
			return strings.TrimPrefix(version, "v"), nil
		}
		return version, nil
	}

	// Fallback to commit count
	output, err = sh.Output("git", "rev-list", "--count", "HEAD")
	if err == nil {
		return fmt.Sprintf("0.0.%s", strings.TrimSpace(output)), nil
	}

	// Final fallback
	return "dev", nil
}

func getGitCommit() (string, error) {
	output, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Ensure the binary is in PATH
func (p Generator) InstallGlobal() error {
	fmt.Println("🌍 Installing trgen globally...")

	// First build and install
	if err := p.Install(); err != nil {
		return err
	}

	// Check if GOPATH/bin is in PATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	goBin := filepath.Join(gopath, "bin")
	pathEnv := os.Getenv("PATH")

	if !strings.Contains(pathEnv, goBin) {
		fmt.Printf("⚠️  Warning: %s is not in your PATH\n", goBin)
		fmt.Printf("💡 Add this to your shell profile (.bashrc, .zshrc, etc.):\n")
		fmt.Printf("   export PATH=\"%s:$PATH\"\n", goBin)
		fmt.Println("🔄 Then restart your terminal or run: source ~/.bashrc")
	} else {
		fmt.Println("✅ trgen is now available globally!")
		fmt.Println("🚀 Try running: trgen --version")
	}

	return nil
}

// Dev installs the generator in development mode (rebuilds on changes)
func (Generator) Dev() error {
	fmt.Println("🔧 Installing trgen in development mode...")

	// Install with race detection and debug info
	version, err := getVersion()
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	gitCommit, err := getGitCommit()
	if err != nil {
		gitCommit = "dev"
	}

	buildTime := time.Now().UTC().Format(time.RFC3339)

	ldflags := fmt.Sprintf("-X github.com/denkhaus/templ-router/cmd/trgen/version.Version=%s-dev "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.GitCommit=%s "+
		"-X github.com/denkhaus/templ-router/cmd/trgen/version.BuildTime=%s",
		version, gitCommit, buildTime)

	fmt.Printf("📦 Building trgen v%s-dev (commit: %s)\n", version, gitCommit)

	// Install with race detection for development
	err = sh.RunWith(map[string]string{
		"CGO_ENABLED": "1", // Enable for race detection
	}, "go", "install", "-race", "-ldflags", ldflags, "./cmd/trgen")

	if err != nil {
		return fmt.Errorf("failed to install template generator in dev mode: %w", err)
	}

	fmt.Println("✅ Development version installed with race detection!")

	return nil
}
