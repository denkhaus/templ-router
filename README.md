# Templ Router

[![Go Version](https://img.shields.io/github/go-mod/go-version/denkhaus/templ-router)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/denkhaus/templ-router)](https://goreportcard.com/report/github.com/denkhaus/templ-router)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/github/actions/workflow/status/denkhaus/templ-router/ci.yml?branch=main)](https://github.com/denkhaus/templ-router/actions)

**A Go library for file-based routing with [templ](https://templ.guide/) templates, dependency injection, and comprehensive middleware support.**

Templ Router is a production-ready library that provides file-based routing, internationalization, authentication, data service integration, validation, caching, and layout inheritance for Go web applications using the templ templating engine and samber/do dependency injection.

## ⚠️ Early Development Warning

**This project is currently in early development stages.** While we strive for stability, please be aware that:

- **API Changes**: The API may change significantly between versions
- **Breaking Changes**: Router methods and interfaces may be modified without prior notice
- **Feature Evolution**: Core features are still being refined and may change
- **Production Use**: Use with caution in production environments

We recommend:
- Pinning to specific versions in your `go.mod`
- Testing thoroughly before upgrading
- Following our [changelog](https://github.com/denkhaus/templ-router/blob/main/CHANGELOG.md) for breaking changes
- Joining our [discussions](https://github.com/denkhaus/templ-router/discussions) for updates

**Stability Target**: We aim for API stability with version 1.0.0

## Features

### 🚀 Core Architecture
- **Dependency Injection**: Built on [samber/do/v2](https://github.com/samber/do) for clean service management
- **Pipeline Architecture**: Composable middleware chain (Template → I18n → Auth)
- **Template Registry**: Generated template registry with automatic route mapping
- **Data Service Integration**: Automatic resolution of named data services for templates

### 🗂️ File-Based Routing
- Routes automatically generated from file structure using `trgen`
- Dynamic parameters: `id_/` (underscore suffix), `locale_/` for internationalization
- Route precedence system for conflict resolution
- Template-to-route mapping with configurable patterns

### 🌍 Internationalization (i18n)
- Multi-language support with `locale_/` directory structure
- YAML-based translations in `.templ.yaml` metadata files
- Context-based translation system (no global `t()` function)
- Automatic locale detection and validation from URLs

### 🔐 Authentication & Authorization
- Three authentication types: `AuthTypePublic`, `AuthTypeUser`, `AuthTypeAdmin`
- Built-in authentication routes: sign in, sign out, sign up
- Session-based authentication with configurable expiry
- Role-based access control with user role validation
- Template-level and route-level auth configuration hierarchy
- Configurable success redirect routes for positive authentication

### 🎨 Layout & Template System
- Layout inheritance with automatic composition
- Error template system with precedence-based resolution
- Template middleware with data service injection
- Configurable template extensions and metadata

### 📊 Data Service Integration
- **Automatic Data Injection**: Templates can declare data service requirements
- **Two Method Patterns**: `GetData()` method or specific `GetDataType()` methods
- **Parameter Injection**: Route parameters automatically passed to data services
- **DI Registration**: Data services registered via `do.ProvideNamed()`

### ⚡ Performance & Validation
- **Cache Service**: Template and route caching for performance optimization
- **Validation Orchestrator**: Comprehensive parameter, route, and template validation
- **Error Handling**: Dedicated error template service with fallback mechanisms
- **File System Abstraction**: Library-agnostic file operations

## Quick Start

### Prerequisites

- Go 1.24 or later (required by go.mod)
- [templ](https://templ.guide/) CLI tool for template compilation
- [trgen](https://github.com/denkhaus/templ-router/cmd/trgen) template generator (install separately)

### Installation

```bash
# Add templ-router to your Go project
go get github.com/denkhaus/templ-router

# Install required tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Generate template registry for your project
trgen --scan-path=app --module-name=github.com/youruser/yourproject
```

**Note:** This is a library, not a standalone application. See the [demo](./demo) directory for a complete example implementation.

## Project Structure

```
# Your project structure (using templ-router as dependency)
your-project/
├── app/                    # Your template directory
│   ├── layout.templ        # Root layout
│   ├── page.templ          # Home page
│   ├── login/
│   │   ├── page.templ      # Login page
│   │   └── page.templ.yaml # Template metadata
│   └── locale_/            # Internationalized routes
│       ├── admin/
│       │   ├── page.templ
│       │   └── page.templ.yaml
│       └── product/
│           └── id_/        # Dynamic parameter (underscore suffix)
│               ├── page.templ
│               └── page.templ.yaml
├── generated/              # Generated by trgen
│   └── templates/
│       └── registry.go     # Template registry
├── pkg/                    # Your application code
├── main.go                 # Your application entry point
└── go.mod                  # Contains: github.com/denkhaus/templ-router v0.x.x

# Note: templ-router is an external dependency in go.mod
# You don't have templ-router/ directory in your project
# Install via: go get github.com/denkhaus/templ-router
```

## File-Based Routing

Routes are automatically generated from your file structure:

```
demo/app/page.templ                    → /
demo/app/login/page.templ              → /login
demo/app/locale_/page.templ            → /en, /de (based on config)
demo/app/locale_/product/id_/page.templ → /en/product/123
```

## Template Metadata System

Each template can have an optional `.templ.yaml` metadata file for configuration:

### Metadata Extraction

Use `metadata.M()` to extract metadata from `.templ.yaml` files:

```go
import "github.com/denkhaus/templ-router/pkg/router/metadata"

templ AdminPage() {
    <h1>{ metadata.M(ctx, "page_title") }</h1>
    <p>{ metadata.M(ctx, "admin_warning") }</p>
    // metadata.M() extracts values from .templ.yaml files
}
```

## Template Configuration

Each template can have an optional `.templ.yaml` configuration file:

```yaml
# demo/app/locale_/admin/page.templ.yaml
auth:
  type: "AdminRequired"
  redirect_url: "/login"

i18n:
  en:
    page_title: "System Administration"
    admin_warning: "Admin Area - Restricted Access"
  de:
    page_title: "Systemadministration"
    admin_warning: "Admin-Bereich - Eingeschränkter Zugang"

dynamic:
  parameters:
    locale:
      validation: "^(en|de)$"
      description: "Language locale code (en or de)"
      supported_values: ["en", "de"]

# Custom metadata accessible via metadata.M()
metadata:
  page_title: "System Administration"
  admin_warning: "Admin Area - Restricted Access"
  custom_field: "Any custom value"
```

### Metadata vs I18n

- **I18n translations:** Use `i18n.T(ctx, "key")` for multi-language content
- **Metadata values:** Use `metadata.M(ctx, "key")` for template-specific configuration
- **Both systems:** Can coexist in the same `.templ.yaml` file

## Authentication

### Authentication Types

Three authentication levels are supported:

```yaml
# Public access (default)
auth:
  type: "Public"

# Requires any authenticated user
auth:
  type: "UserRequired"
  redirect_url: "/login"

# Requires admin privileges
auth:
  type: "AdminRequired"
  redirect_url: "/login"
```

### Authentication Routes

The router provides built-in authentication API endpoints:

```bash
# Authentication API endpoints (automatically registered)
POST /api/auth/signin     # User sign in
POST /api/auth/signout    # User sign out  
POST /api/auth/signup     # User registration
```

### Success Redirect Configuration

Configure where users are redirected after successful authentication:

```bash
# Environment configuration for auth redirects
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/welcome
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/

# Default auth route
TR_AUTH_SIGNIN_ROUTE=/login

# Internationalized success routes with {locale} parameter
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/{locale}/dashboard
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/{locale}/welcome
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/{locale}/

# The {locale} parameter is automatically replaced with the current locale:
# /{locale}/dashboard → /en/dashboard or /de/dashboard
```

### Required Template Routes

You need to create template pages for authentication:

```bash
# Required template files for auth flow
app/login/page.templ      # Sign in page (GET /login)
app/signup/page.templ     # Sign up page (GET /signup)
```

### Auth API Integration

```go
// Auth API routes are automatically registered
authHandlers := do.MustInvoke[interfaces.AuthHandlers](injector)
authHandlers.RegisterRoutes(func(method, path string, handler http.HandlerFunc) {
    mux.Post(path, handler) // Registers POST /api/auth/* routes
})

// Your templates use these API endpoints:
// <form hx-post="/api/auth/signin">
// <form hx-post="/api/auth/signup">  
// <form method="POST" action="/api/auth/signout">
```

## Internationalization

Translation files use locale-specific keys:

```yaml
# demo/app/locale_/dashboard/page.templ.yaml
i18n:
  en:
    page_title: "Dashboard"
    page_subtitle: "Overview of your application metrics"
  de:
    page_title: "Dashboard"
    page_subtitle: "Übersicht Ihrer Anwendungsmetriken"
```

Use in templates with real i18n functions:

```go
import "github.com/denkhaus/templ-router/pkg/router/i18n"

// Real i18n usage with i18n.T() function
templ DashboardPage() {
    <h1>{ i18n.T(ctx, "page_title") }</h1>
    <p>{ i18n.T(ctx, "page_subtitle") }</p>
    <a href={ i18n.LocalizeSafeURL(ctx, "/admin") }>
        { i18n.T(ctx, "nav_admin") }
    </a>
}

// Access current locale and template info
templ DebugPage() {
    <p>Current Locale: { i18n.GetCurrentLocale(ctx) }</p>
    <p>Template: { i18n.GetCurrentTemplate(ctx) }</p>
    <p>Available Keys: { fmt.Sprint(len(i18n.GetAvailableKeys(ctx))) }</p>
}
```

## Development Workflow

### For Library Users

```bash
# 1. Generate templates (required after template changes)
templ generate

# 2. Generate template registry (required after adding/removing templates)
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# 3. Build your application
go build

# 4. Run your application
./your-app

# 5. Development with hot reload (optional)
# Install air: go install github.com/cosmtrek/air@latest
air
```

### For Library Development

The templ-router library itself uses [Mage](https://magefile.org/) for development:

```bash
# Library development (only for templ-router contributors)
mage dev                    # Start demo server
mage test:all               # Run library tests
mage build:all              # Build library for all platforms
```

## Configuration

Environment variables use configurable prefix (set in `RegisterRouterServices(prefix)`):

```bash
# Server Configuration (PREFIX_SECTION_FIELD)
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080
TR_SERVER_BASE_URL=http://localhost:8080
TR_SERVER_READ_TIMEOUT=30s
TR_SERVER_WRITE_TIMEOUT=30s
TR_SERVER_IDLE_TIMEOUT=120s
TR_SERVER_SHUTDOWN_TIMEOUT=30s

# Database Configuration  
TR_DATABASE_HOST=localhost
TR_DATABASE_PORT=5432
TR_DATABASE_DATABASE_USER=postgres
TR_DATABASE_PASSWORD=postgres
TR_DATABASE_NAME=router_db
TR_DATABASE_SSL_MODE=disable

# Authentication & Sessions
TR_AUTH_CREATE_DEFAULT_ADMIN=true
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_NAME=session_id

# Internationalization
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en
TR_I18N_FALLBACK_LOCALE=en

# Layout System
TR_LAYOUT_ROOT_DIRECTORY=app
TR_LAYOUT_ENABLE_INHERITANCE=true
TR_LAYOUT_TEMPLATE_EXTENSION=.templ

# Template Generator
TR_TEMPLATE_GENERATOR_OUTPUT_DIR=generated/templates
TR_TEMPLATE_GENERATOR_PACKAGE_NAME=templates
```

# Email Configuration
TR_EMAIL_SMTP_HOST=
TR_EMAIL_SMTP_PORT=587
TR_EMAIL_SMTP_USERNAME=
TR_EMAIL_SMTP_PASSWORD=
TR_EMAIL_SMTP_USE_TLS=true
TR_EMAIL_FROM_EMAIL=noreply@example.com
TR_EMAIL_FROM_NAME="Router Application"
TR_EMAIL_ENABLE_DUMMY_MODE=true

# Security Configuration
TR_SECURITY_CSRF_SECRET=change-me-in-production
TR_SECURITY_CSRF_SECURE=false
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_SECURITY_RATE_LIMIT_REQUESTS=100
TR_SECURITY_ENABLE_SECURITY_HEADERS=true

# Logging Configuration
TR_LOGGING_LEVEL=info
TR_LOGGING_FORMAT=json
TR_LOGGING_OUTPUT=stdout
TR_LOGGING_ENABLE_FILE=false
TR_LOGGING_FILE_PATH=logs/router.log

# Environment Configuration
TR_ENVIRONMENT_KIND=develop
```

**Configuration Sections:** Server, Database, Auth, Email, Security, Logging, I18n, Layout, TemplateGenerator, Environment

## Data Services

Templates can automatically receive data through the Data Service system:

### Data Service Interface Patterns

**Pattern 1: Simple GetData() Method**
```go
type ProductDataService interface {
    GetData(ctx context.Context, params map[string]string) (*ProductData, error)
}

func (s *productDataServiceImpl) GetData(ctx context.Context, params map[string]string) (*ProductData, error) {
    productID := params["id"] // Route parameters automatically injected
    return &ProductData{
        ID:   productID,
        Name: "Product " + productID,
    }, nil
}
```

**Pattern 2: GetData() + Specific Methods**
```go
type UserDataService interface {
    GetData(ctx context.Context, params map[string]string) (*UserData, error)
    GetUserData(ctx context.Context, params map[string]string) (*UserData, error)
}

// Both methods available - router chooses appropriate one
func (s *userDataServiceImpl) GetUserData(ctx context.Context, params map[string]string) (*UserData, error) {
    userID := params["id"]
    locale := params["locale"] // Multi-language support
    return &UserData{ID: userID, Name: "User " + userID}, nil
}
```

### Data Service Registration

```go
// Register data services with dependency injection
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
```

### Template Integration

Templates automatically receive data when they declare data service requirements:

```go
// Template signature indicates data service requirement
templ UserProfilePage(user *UserData) {
    <h1>{ user.Name }</h1>
    <p>{ user.Email }</p>
}

// Router automatically:
// 1. Detects UserData requirement
// 2. Resolves UserDataService
// 3. Calls GetUserData() or GetData()
// 4. Injects route parameters: params["id"], params["locale"]
// 5. Passes result to template
```

### Parameter Injection

Route parameters are automatically injected into data service methods:

```go
func (s *userDataServiceImpl) GetUserData(ctx context.Context, params map[string]string) (*UserData, error) {
    userID := params["id"]       // From route: /user/123 → id="123"
    locale := params["locale"]   // From route: /en/user/123 → locale="en"
    
    // Use parameters to fetch localized user data
    return s.fetchUser(userID, locale)
}
```

## Dependency Injection

The library uses [samber/do/v2](https://github.com/samber/do) for dependency injection:

```go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/yourproject/generated/templates"
    "github.com/yourproject/pkg/services"
    "github.com/samber/do/v2"
)

func main() {
    // Create DI container
    container := di.NewContainer()
    defer container.Shutdown()

    // Register router services with config prefix
    container.RegisterRouterServices("TR")
    injector := container.GetInjector()

    // Create your services
    templateRegistry, _ := templates.NewRegistry(injector)
    userStore, _ := services.NewDefaultUserStore(injector)

    // Register application services using options pattern
    container.RegisterApplicationServices(
        di.WithTemplateRegistry(templateRegistry),
        di.WithUserStore(userStore),
    )

    // Register named data services for template injection
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
    do.ProvideNamed(injector, "OrderDataService", dataservices.NewOrderDataService)

    // Get router and initialize
    router := container.GetRouter()
    router.Initialize()
    
    // Register routes with your HTTP router (chi, gin, etc.)
    mux := chi.NewRouter()
    router.RegisterRoutes(mux)
    
    // Register auth routes
    authHandlers := do.MustInvoke[interfaces.AuthHandlers](injector)
    authHandlers.RegisterRoutes(func(method, path string, handler http.HandlerFunc) {
        mux.Post(path, handler)
    })
}
```

## Template Generator (trgen)

The `trgen` CLI tool is essential for generating template registries from your file structure:

```bash
# Install generator
go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Generate template registry (required flags)
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Watch mode for development
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch

# Custom output directory
trgen --scan-path=app --module-name=github.com/youruser/yourproject --output-dir=internal/templates
```

**Required Parameters:**
- `--scan-path`: Directory containing your `.templ` files (e.g., `app`, `templates`)
- `--module-name`: Your Go module name from `go.mod`

**Generated Output:**
- Creates `generated/templates/registry.go` with template registry
- Maps file paths to route patterns automatically
- Detects data service requirements from template signatures

## Production Deployment

### Docker

```bash
# Build and run with Docker
mage docker:up

# Or manually
docker build -t templ-router .
docker run -p 8084:8084 templ-router
```

### Binary

```bash
# Build for current platform
mage generator:build

# Build for all platforms
mage build:all
```

## Architecture

The framework follows clean architecture principles:

- **Router Core**: File-based route discovery and registration
- **Middleware**: Authentication, i18n, template rendering
- **Services**: Configuration, caching, validation
- **DI Container**: Dependency management with samber/do
- **Template System**: Templ-based rendering with layout inheritance

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `mage test:all`
5. Submit a pull request

## Development Setup

```bash
git clone https://github.com/denkhaus/templ-router.git
cd templ-router
go mod tidy
mage dev
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- [Documentation](https://github.com/denkhaus/templ-router/wiki)
- [API Reference](https://pkg.go.dev/github.com/denkhaus/templ-router)
- [Issue Tracker](https://github.com/denkhaus/templ-router/issues)
- [Discussions](https://github.com/denkhaus/templ-router/discussions)

---

**Built for the Go and templ community**