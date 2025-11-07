# Template Generator (trgen)

**CLI tool for generating template registries from file structure.**

## Overview

`trgen` (templ-router generator) automatically scans your templ template files and generates a template registry with route mappings, eliminating manual route registration.

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
go install github.com/denkhaus/templ-router/cmd/trgen@latest
trgen --version
```

### Install from Source

```bash
git clone https://github.com/denkhaus/templ-router.git
cd templ-router
mage generator:install
```

## Usage

### Basic Usage

```bash
# Generate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Watch mode for development
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch
```

### Environment Variables

- `TRGEN_SCAN_PATH`: Directory containing `.templ` files
- `TRGEN_MODULE_NAME`: Go module name from `go.mod`
- `TRGEN_WATCH_MODE`: Enable watch mode (`true`/`false`)
- `TRGEN_WATCH_EXTENSIONS`: File extensions to watch

## File Structure and Routes

`trgen` converts file structure to HTTP routes:

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

### Dynamic Parameters

Use `_` suffix for dynamic segments:

- `id_/` → `{id}`
- `locale_/` → `{locale}`
- `category_/slug_/` → `{category}/{slug}`

### Generated Route Mappings

```go
var routeMappings = map[string]*RouteInfo{
    "/": {
        Path:         "/",
        TemplateFunc: Page,
        Parameters:   []string{},
        IsLocalized:  false,
        RequiresAuth: false,
        DataServices: []string{},
    },
    "/{locale}/user/{id}": {
        Path:         "/{locale}/user/{id}",
        TemplateFunc: Page,
        Parameters:   []string{"locale", "id"},
        IsLocalized:  true,
        RequiresAuth: true,
        DataServices: []string{"UserDataService"},
    },
}
```

## Template Detection

### Template Files

`trgen` scans for `.templ` files and generates Go functions:

```go
// app/page.templ → Page
// app/login/page.templ → Page
// app/locale_/dashboard/page.templ → Page
```

### Component Metadata

`trgen` reads `.templ.yaml` files for metadata:

```yaml
# app/page.templ.yaml
metadata:
  title: "Home Page"
  description: "Welcome to our application"

auth:
  type: "Public"

i18n:
  en:
    welcome: "Welcome"
  de:
    welcome: "Willkommen"
```

### Data Service Detection

`trgen` automatically detects data service requirements:

```go
// In your template
func (s *service) GetData(ctx context.Context, data interface{}) (*UserData, error) {
    // Implementation
}

// Generated registry will automatically inject this service
```

## Configuration

### Command Line Options

```bash
trgen [options]

Options:
  --scan-path string      Directory containing .templ files (default: app)
  --module-name string    Go module name (required)
  --watch                 Enable watch mode
  --output string         Output file (default: generated/templates/registry.go)
  --verbose               Enable verbose logging
  --version               Show version information
  --help                  Show help
```

### Configuration File

Create `.trgen.yaml` in your project root:

```yaml
scan_path: app
module_name: github.com/youruser/yourproject
output: generated/templates/registry.go
watch:
  enabled: true
  extensions: [".templ", ".templ.yaml"]
  ignore_dirs: [".git", "node_modules", "vendor"]
```

## Development Workflow

### Integration with Build Tools

#### Mage (Recommended)

```go
// magefiles/build.go
func BuildRegistry() {
    sh.RunV("trgen", "--scan-path=app", "--module-name=github.com/youruser/yourproject")
}

func Watch() {
    sh.RunV("trgen", "--scan-path=app", "--module-name=github.com/youruser/yourproject", "--watch")
}
```

#### Makefile

```makefile
.PHONY: registry watch
registry:
	trgen --scan-path=app --module-name=github.com/youruser/yourproject

watch:
	trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch
```

### Hot Reload Setup

#### Using Air

```yaml
# .air.toml
[build]
  cmd = "templ generate && trgen --scan-path=app --module-name=github.com/youruser/yourproject && go run ."
  bin = "main"
  include_ext = ["go", "templ"]
```

## Advanced Features

### Internationalization

`trgen` supports locale-based routing:

```
app/
├── locale_/
│   ├── page.templ          → /{locale}
│   └── dashboard/
│       └── page.templ      → /{locale}/dashboard
```

### Component Metadata Inheritance

Components can have their own metadata:

```yaml
# app/components/footer.templ.yaml
i18n:
  en: { footer_copyright: "© 2024 My Company" }
  de: { footer_copyright: "© 2024 Meine Firma" }
```

### Data Service Integration

```go
// Generated registry automatically injects required services
registry, err := templates.NewRegistry(injector)
```

## Troubleshooting

### Common Issues

**Command not found:**
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Module name mismatch:**
```bash
# Use exact module name from go.mod
trgen --module-name=github.com/youruser/yourproject
```

**Scan path issues:**
```bash
# Use relative path from project root
trgen --scan-path=./app
```

### Debug Mode

```bash
trgen --verbose --scan-path=app --module-name=github.com/youruser/yourproject
```

## Examples

### Basic Project Setup

```bash
mkdir my-project && cd my-project
go mod init github.com/youruser/myproject

# Create structure
mkdir -p app generated/templates

# Create template
echo 'package main
templ Page() {
    <h1>Hello World</h1>
}' > app/page.templ

# Generate registry
trgen --scan-path=app --module-name=github.com/youruser/myproject
```

### Watch Mode Development

```bash
# Terminal 1: Watch template changes
trgen --scan-path=app --module-name=github.com/youruser/myproject --watch

# Terminal 2: Run application
templ generate
go run .
```

### Internationalized Project

```
app/
├── layout.templ
├── locale_/
│   ├── page.templ
│   └── about.templ
└── admin/
    └── locale_/
        └── dashboard.templ
```

This generates routes for `/en`, `/de`, `/en/admin/dashboard`, etc.

## Best Practices

1. **Use descriptive file names** - Use descriptive folder structure instead
2. **Organize with folders** - Group related templates
3. **Use dynamic parameters** - `id_/`, `locale_/` for flexible routing
4. **Add metadata** - Use `.templ.yaml` for configuration
5. **Watch mode** - Use `--watch` during development
6. **Version control** - Include generated registry in `.gitignore`

## Integration with IDE

### VS Code

Add to `.vscode/tasks.json`:

```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Generate Registry",
            "type": "shell",
            "command": "trgen",
            "args": ["--scan-path=app", "--module-name=github.com/youruser/yourproject"],
            "group": "build"
        }
    ]
}
```

## API Reference

### Generated Registry Interface

```go
type Registry interface {
    GetTemplate(path string) (templ.Component, error)
    GetAllRoutes() []string
    GetRouteInfo(path string) (*RouteInfo, error)
}
```

### RouteInfo Structure

```go
type RouteInfo struct {
    Path         string
    TemplateFunc templ.Component
    Parameters   []string
    IsLocalized  bool
    RequiresAuth bool
    DataServices []string
}
```

---

**Next Steps**: Read [File-Based Routing](FILE-BASED-ROUTING.md) for detailed routing concepts, or [Configuration](CONFIGURATION.md) for environment setup.