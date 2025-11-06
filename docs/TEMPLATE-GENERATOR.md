# Template Generator (trgen)

**Complete guide to the trgen CLI tool for generating template registries.**

## Overview

`trgen` (templ-router generator) is a CLI tool that automatically scans your templ template files and generates a template registry with route mappings. It eliminates the need for manual route registration by analyzing your file structure and creating the necessary mappings.

**Key Features:**
- Automatic route discovery from file structure
- Dynamic parameter detection (`id_`, `locale_`, etc.)
- Template registry generation for dependency injection
- Watch mode for development
- Support for internationalized routes
- Data service requirement detection

## Installation

### Quick Install

```bash
# Install the latest release
go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Verify installation
trgen --version
trgen --help
```

### Install from Source

```bash
# Clone the repository
git clone https://github.com/denkhaus/templ-router.git
cd templ-router

# Install from source
mage generator:install

# Or build manually
go install ./cmd/trgen
```

### Using Mage (Recommended)

```bash
# Install the latest version globally
mage generator:installGlobal

# Or just install to GOPATH/bin
mage generator:install

# For development with race detection
mage generator:dev

# Build locally (outputs to bin/)
mage generator:build
```

### Available Mage Tasks

#### Development Tasks
- **`mage generator:build`** - Build the generator locally (outputs to `bin/`)
- **`mage generator:install`** - Build and install to `$GOPATH/bin`
- **`mage generator:installGlobal`** - Install globally and check PATH setup
- **`mage generator:dev`** - Install development version with race detection

#### Testing Tasks
- **`mage generator:test`** - Run all template generator tests
- **`mage generator:testCoverage`** - Run tests with coverage report
- **`mage generator:version`** - Show current version information

#### Release Tasks
- **`mage generator:release`** - Build for multiple platforms (Linux, macOS, Windows)
- **`mage generator:clean`** - Remove build artifacts

### Development Installation

```bash
# Install development version with race detection
mage generator:dev

# Build locally (outputs to bin/)
mage generator:build
```

### Version Information

The generator includes comprehensive version information:

```bash
$ trgen --version
trgen version v1.2.3-abc1234

$ mage generator:version
Template Generator Version: v1.2.3
Git Commit: abc1234
Go Version: go1.21.0
Platform: linux/amd64
```

## Basic Usage

### Required Parameters

`trgen` requires two parameters:

- `--scan-path`: Directory containing your `.templ` files
- `--module-name`: Your Go module name from `go.mod`

### Generate Template Registry

```bash
# Navigate to your project directory (where go.mod is located)
cd your-project

# Generate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Example output
# ✅ Scanning templates in: app
# ✅ Found 15 templates
# ✅ Generated registry: generated/templates/registry.go
# ✅ Processing complete!
```

### Using Environment Variables

You can use environment variables instead of command-line flags:

```bash
export TRGEN_SCAN_PATH=app
export TRGEN_MODULE_NAME=github.com/youruser/yourproject

# Generate with environment variables
trgen

# Or one-liner
TRGEN_SCAN_PATH=app TRGEN_MODULE_NAME=github.com/youruser/yourproject trgen
```

## Command Line Options

### Core Options

```bash
trgen [flags]

--scan-path string          # Directory containing .templ files (required)
--module-name string        # Go module name from go.mod (required)
--output-dir string         # Output directory for generated files (default: generated/templates)
--package-name string       # Package name for generated registry (default: templates)
--watch                     # Enable watch mode for development
--watch-extensions strings  # File extensions to watch in watch mode (default: [.templ,.yaml,.yml])
--verbose                   # Enable verbose output
--help                      # Show help information
--version                   # Show version information
```

### Example Commands

```bash
# Basic usage
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Custom output directory
trgen --scan-path=app --module-name=github.com/youruser/yourproject --output-dir=internal/registry

# Custom package name
trgen --scan-path=app --module-name=github.com/youruser/yourproject --package-name=mytemplates

# Verbose output
trgen --scan-path=app --module-name=github.com/youruser/yourproject --verbose
```

## Watch Mode

### Development Watch Mode

```bash
# Start watch mode for development
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch

# Output
# 🔄 Starting watch mode...
# 📁 Watching: app
# 🔍 Extensions: .templ,.yaml,.yml
# ⏳ Waiting for file changes...
```

### Custom Watch Extensions

```bash
# Watch specific file types
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch --watch-extensions=".templ,.yaml"

# Watch additional files like CSS or JS
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch --watch-extensions=".templ,.yaml,.css"
```

### Watch Mode Features

- **Automatic regeneration**: Registry updates when files change
- **Debouncing**: Prevents excessive regeneration during rapid changes
- **Error handling**: Shows generation errors without stopping watch mode
- **File filtering**: Only watches specified file extensions

## Generated Output

### Registry Structure

`trgen` generates a `registry.go` file with the following structure:

```go
// generated/templates/registry.go
package templates

import (
    "context"
    "github.com/denkhaus/templ-router/pkg/interfaces"
)

// Template registry with route mappings
type Registry struct {
    injector *do.Injector
}

// Route information for each template
type RouteInfo struct {
    Path           string
    TemplateFunc   interfaces.TemplateFunc
    Parameters     []string
    IsLocalized    bool
    RequiresAuth   bool
    DataServices   []string
}

// NewRegistry creates a new template registry
func NewRegistry(injector *do.Injector) (*Registry, error) {
    registry := &Registry{
        injector: injector,
    }
    return registry, nil
}

// GetRouteInfo returns route information for a template
func (r *Registry) GetRouteInfo(templateName string) (*RouteInfo, error) {
    // Generated route mappings
}

// GetAllRoutes returns all available routes
func (r *Registry) GetAllRoutes() map[string]*RouteInfo {
    // Generated route mappings
}
```

### Route Mapping Examples

Given this file structure:

```
app/
├── page.templ                    → /
├── login/
│   └── page.templ              → /login
├── dashboard/
│   └── page.templ              → /dashboard
└── locale_/
    ├── page.templ              → /{locale}
    └── user/
        └── id_/
            └── page.templ      → /{locale}/user/{id}
```

`trgen` generates route mappings like:

```go
// Generated route mappings
var routeMappings = map[string]*RouteInfo{
    "/": {
        Path:         "/",
        TemplateFunc: HomePage,
        Parameters:   []string{},
        IsLocalized:  false,
        RequiresAuth: false,
        DataServices: []string{},
    },
    "/login": {
        Path:         "/login",
        TemplateFunc: LoginPage,
        Parameters:   []string{},
        IsLocalized:  false,
        RequiresAuth: false,
        DataServices: []string{},
    },
    "/{locale}/user/{id}": {
        Path:         "/{locale}/user/{id}",
        TemplateFunc: UserProfilePage,
        Parameters:   []string{"locale", "id"},
        IsLocalized:  true,
        RequiresAuth: true,
        DataServices: []string{"UserDataService"},
    },
}
```

## Template Detection

### Template File Discovery

`trgen` automatically discovers `.templ` files in your scan path:

```bash
# Scans for all .templ files recursively
trgen --scan-path=app --module-name=github.com/youruser/yourproject --verbose

# Output shows discovered templates
# 🔍 Discovered template: app/page.templ
# 🔍 Discovered template: app/login/page.templ
# 🔍 Discovered template: app/dashboard/page.templ
# 🔍 Discovered template: app/locale_/user/id_/page.templ
```

### Dynamic Parameter Detection

`trgen` automatically detects dynamic parameters from directory names:

```bash
# Dynamic parameters end with underscore _
app/user/id_/page.templ        → Parameter: id
app/product/slug_/page.templ   → Parameter: slug
app/docs/version_/page.templ   → Parameter: version

# Special locale parameter
app/locale_/page.templ         → Parameter: locale (special handling)
```

### Data Service Detection

`trgen` analyzes template signatures to detect data service requirements:

```go
// Template with data service requirement
templ UserProfilePage(user *UserData) {
    // trgen detects UserDataService requirement
}

// Template with multiple data services
templ DashboardPage(user *UserData, stats *DashboardStats) {
    // trgen detects UserDataService and DashboardStatsService requirements
}
```

## Metadata Processing

### YAML Metadata Files

`trgen` processes `.templ.yaml` files alongside templates:

```yaml
# app/dashboard/page.templ.yaml
metadata:
  page_title: "Dashboard"
  theme: "dark"

i18n:
  en:
    page_title: "Dashboard"
    welcome_message: "Welcome to your dashboard"
  de:
    page_title: "Dashboard"
    welcome_message: "Willkommen in Ihrem Dashboard"

auth:
  type: "UserRequired"
  redirect_url: "/login"

data_services:
  - "DashboardDataService"
  - "UserStatsDataService"
```

### Metadata Integration

Generated registry includes metadata information:

```go
// Generated route info includes metadata
"/dashboard": {
    Path:         "/dashboard",
    TemplateFunc: DashboardPage,
    Parameters:   []string{},
    IsLocalized:  false,
    RequiresAuth: true,        // From YAML metadata
    DataServices: []string{   // From YAML metadata
        "DashboardDataService",
        "UserStatsDataService",
    },
}
```

## Project Integration

### Integration in Main Application

```go
// main.go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/youruser/yourproject/generated/templates"
    "github.com/samber/do/v2"
)

func main() {
    // Create DI container
    container := di.NewContainer()
    injector := container.GetInjector()

    // Create template registry from generated code
    templateRegistry, err := templates.NewRegistry(injector)
    if err != nil {
        panic(err)
    }

    // Register with application
    container.RegisterApplicationServices(
        di.WithTemplateRegistry(templateRegistry),
    )
}
```

### Build Process Integration

#### Makefile Integration

```makefile
.PHONY: generate-templates
generate-templates:
	@echo "Generating template registry..."
	@trgen --scan-path=app --module-name=github.com/youruser/yourproject
	@echo "✅ Template registry generated!"

.PHONY: generate-templates-watch
generate-templates-watch:
	@echo "Starting template registry watch mode..."
	@trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch

.PHONY: build
build: generate-templates
	@echo "Building application..."
	@go build -o bin/your-app

.PHONY: dev
dev: generate-templates-watch
	@echo "Starting development server..."
	@air -c .air.toml
```

#### Mage Integration

```go
// magefiles/main.go
package main

var Default = Dev

func GenerateTemplates() error {
    return sh.RunV("trgen",
        "--scan-path=app",
        "--module-name=github.com/youruser/yourproject")
}

func GenerateTemplatesWatch() error {
    return sh.RunV("trgen",
        "--scan-path=app",
        "--module-name=github.com/youruser/yourproject",
        "--watch")
}

func Build() error {
    if err := GenerateTemplates(); err != nil {
        return err
    }
    return sh.RunV("go", "build")
}

func Dev() error {
    // Run template generation and watch in parallel
    return RunParallel(
        GenerateTemplatesWatch,
        func() error { return sh.RunV("air") },
    )
}
```

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/build.yml
name: Build and Test

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21

    - name: Install trgen
      run: go install github.com/denkhaus/templ-router/cmd/trgen@latest

    - name: Generate templates
      run: trgen --scan-path=app --module-name=github.com/${{ github.repository }}

    - name: Build
      run: go build

    - name: Test
      run: go test ./...
```

### Docker Integration

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

# Install trgen
RUN go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Copy source and generate templates
WORKDIR /app
COPY . .
RUN trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Build application
RUN go build -o main .

# Final image
FROM alpine:latest
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
```

## Advanced Usage

### Multiple Template Directories

```bash
# Generate from multiple directories (run multiple times)
trgen --scan-path=app --module-name=github.com/youruser/yourproject --output-dir=generated/app
trgen --scan-path=components --module-name=github.com/youruser/yourproject --output-dir=generated/components
```

### Custom Package Configuration

```bash
# Generate with custom package and directory
trgen \
  --scan-path=app \
  --module-name=github.com/youruser/yourproject \
  --output-dir=internal/registry \
  --package-name=registry
```

### Integration with Tools

#### Air Hot Reload

```toml
# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "templ", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

#### Pre-commit Hooks

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: trgen
        name: Generate template registry
        entry: trgen
        language: system
        args: ["--scan-path=app", "--module-name=github.com/youruser/yourproject"]
        pass_filenames: false
        always_run: true
```

## Troubleshooting

### Common Issues

#### Command Not Found

```bash
# Check if trgen is installed
which trgen

# If not found, install it
go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Check GOPATH/bin is in your PATH
echo $PATH | grep -q "$(go env GOPATH)/bin"
```

#### Module Name Mismatch

```bash
# Error: module name doesn't match go.mod
trgen --scan-path=app --module-name=wrong/module/name

# Solution: Use exact module name from go.mod
head -n 1 go.mod
# Output: module github.com/youruser/yourproject

trgen --scan-path=app --module-name=github.com/youruser/yourproject
```

#### Scan Path Issues

```bash
# Error: scan path not found
trgen --scan-path=wrong/path --module-name=github.com/youruser/yourproject

# Solution: Use correct relative path
ls -la app/
trgen --scan-path=app --module-name=github.com/youruser/yourproject
```

#### Generated Files Not Found

```bash
# Check if generated files exist
ls -la generated/templates/

# Verify output directory
trgen --scan-path=app --module-name=github.com/youruser/yourproject --verbose

# Custom output directory
trgen --scan-path=app --module-name=github.com/youruser/yourproject --output-dir=internal/registry
```

### Debug Mode

```bash
# Use verbose output for debugging
trgen --scan-path=app --module-name=github.com/youruser/yourproject --verbose

# Output includes detailed information
# 🔍 Starting template generation...
# 📁 Scan path: app
# 📦 Module name: github.com/youruser/yourproject
# 📂 Output directory: generated/templates
# 🏷️  Package name: templates
# 🔍 Scanning for .templ files...
# 🔍 Found: app/page.templ
# 🔍 Processing: app/page.templ
# ✅ Route: /
# ✅ Generated 1 routes
# ✅ Registry written to: generated/templates/registry.go
```

### Version Information

```bash
# Check trgen version
trgen --version
# Output: trgen version v1.2.3-abc1234

# Check for updates
go install github.com/denkhaus/templ-router/cmd/trgen@latest
```

## Performance Tips

### Large Projects

```bash
# Use watch mode for large projects
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch

# Limit file extensions to watch
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch --watch-extensions=".templ"

# Use verbose output sparingly
trgen --scan-path=app --module-name=github.com/youruser/yourproject  # Without --verbose
```

### Optimization

```bash
# Regenerate only when templates change
# Use watch mode or integrate with build tools

# Exclude unnecessary files
# Don't put large non-template files in scan path

# Use appropriate output directory structure
trgen --scan-path=app --module-name=github.com/youruser/yourproject --output-dir=internal/registry
```

## Development Workflow

### Local Development

```bash
# Build and test locally
mage generator:build
mage generator:test

# Install development version
mage generator:dev

# Run with coverage
mage generator:testCoverage
```

### Testing

The generator has comprehensive test coverage:

```bash
# Run all tests
mage generator:test

# Generate coverage report
mage generator:testCoverage
# Opens coverage/generator.html in browser
```

### Release Process

```bash
# 1. Tag the release
git tag v1.2.3
git push origin v1.2.3

# 2. Build release binaries
mage generator:release

# 3. Test the release
./release/template-generator-v1.2.3/template-generator-linux-amd64 --version
```

## Troubleshooting

### PATH Issues

If `trgen` command is not found after installation:

```bash
# Check if GOPATH/bin is in PATH
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✅ GOPATH/bin is in PATH" || echo "❌ GOPATH/bin not in PATH"

# Add to your shell profile (.bashrc, .zshrc, etc.)
export PATH="$(go env GOPATH)/bin:$PATH"

# Or use the global install task
mage generator:installGlobal
```

### Build Issues

```bash
# Clean and rebuild
mage generator:clean
mage generator:build

# Check Go version (requires Go 1.21+)
go version

# Verify dependencies
go mod tidy
```

### Version Issues

```bash
# Check current version
mage generator:version

# Force rebuild with latest version
mage generator:clean
mage generator:install
```

### Common Problems

**Error: "no Go files in directory"**
- Ensure you're in the correct project directory with `go.mod`
- Check that `--scan-path` points to your template directory

**Error: "module name not found"**
- Verify `--module-name` matches your go.mod file exactly
- Use the format: `github.com/youruser/yourproject`

**Permission denied**
```bash
# Fix Go bin directory permissions
chmod +x $(go env GOPATH)/bin/trgen

# Or install globally
mage generator:installGlobal
```

## Advanced Features

### Cross-Platform Builds

The release task builds for multiple platforms:

```bash
mage generator:release
```

Creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

Output: `release/template-generator-v{version}/`

### Automatic Versioning

The build system automatically determines version information:

1. **Git Tags**: Uses the latest git tag (e.g., `v1.2.3`)
2. **Git Describe**: Falls back to `git describe --tags --always --dirty`
3. **Commit Count**: Uses commit count if no tags available
4. **Fallback**: Uses "dev" if git is not available

## Integration Examples

### Makefile Integration

```makefile
.PHONY: generate-templates
generate-templates:
	@echo "Generating template registry..."
	@trgen --scan-path app --module-name=github.com/youruser/yourproject
	@echo "Template registry generated successfully!"

.PHONY: install-generator
install-generator:
	@echo "Installing template generator..."
	@mage generator:install

.PHONY: watch-templates
watch-templates:
	@echo "Watching templates for changes..."
	@trgen --scan-path app --module-name=github.com/youruser/yourproject --watch
```

### CI/CD Integration

```yaml
# .github/workflows/build.yml
- name: Install Template Generator
  run: mage generator:install

- name: Generate Templates
  run: trgen --scan-path app --module-name=github.com/youruser/yourproject

- name: Verify Generated Files
  run: git diff --exit-code generated/
```

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first project
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Understand routing patterns
- **[Configuration](CONFIGURATION.md)** - Configure trgen behavior
- **[Data Services](DATA-SERVICES.md)** - Use data services with templates

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Documentation](../README.md)** - Main documentation