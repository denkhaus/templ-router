# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Templ Router** is a production-ready Go library for file-based routing with templ templates, dependency injection, and comprehensive middleware support. This is a **library**, not a standalone application - developers add it as a dependency to their own Go projects.

**Current Status**: Early development (API may change, aiming for stability with v1.0.0)

## Core Technologies

- **Go 1.24+** with modern Go features
- **templ** template engine (https://templ.guide/)
- **samber/do/v2** for dependency injection
- **chi/v5** for HTTP routing
- **YAML** for configuration and internationalization
- **Mage** for build automation
- **Ginkgo/Gomega** for E2E testing
- **Tailwind CSS** for styling (in demo)

## Key Architecture Components

### 1. **File-Based Routing System**
- Routes automatically generated from file structure using `trgen` CLI tool
- Dynamic parameters with `_` suffix (e.g., `id_/`, `locale_/`)
- Route precedence system for conflict resolution
- Template-to-route mapping with configurable patterns
- Support for internationalized routes via `locale_/` directory structure

### 2. **Clean Architecture with Dependency Injection**
- **DI Container**: Built on samber/do/v2 with separate registration for:
  - Router services (configurable via prefix, default "TR")
  - Application services (options pattern)
- **Service Management**: Named dependencies for data service resolution
- **Clean Separation**: Infrastructure, application, and domain layers separated

### 3. **Template System**
- **Templ Files**: Go-based templates with type safety
- **Metadata System**: `.templ.yaml` files for:
  - Authentication configuration
  - Internationalization translations
  - Template metadata
  - Dynamic parameter validation
- **Component Metadata**: Self-contained component configurations:
  - Components can have their own `.templ.yaml` files (e.g., `footer.templ.yaml`)
  - Component metadata overrides page metadata when accessed via `/components/*` routes
  - Support for component-specific i18n, auth settings, and metadata
  - Enables truly reusable components with built-in internationalization
- **Layout Inheritance**: Automatic template composition
- **Error Templates**: Fallback mechanisms with precedence resolution

### 4. **Internationalization (i18n)**
- **Multi-language support** with YAML-based translations
- **Context-based system** (no global `t()` function)
- **Nested key structure** with dot notation support
- **Multiple formats**: flat, nested, and multi-locale configurations
- **Smart locale detection** and routing

### 5. **Authentication & Authorization**
- **Three auth types**: `Public`, `UserRequired`, `AdminRequired`
- **Session-based authentication** with configurable expiry
- **Built-in API endpoints**: `/api/auth/signin`, `/api/auth/signout`, `/api/auth/signup`
- **Role-based access control** with user validation
- **Configurable redirect routes** for authentication flows

### 6. **Data Service Integration**
- **Automatic injection** based on template requirements
- **RouterContext interface** for unified parameter access:
  - URL parameters from route paths
  - Query parameters from URL strings
  - Request context access
- **Two method patterns**: `GetData()` fallback and specific `Get<DataStruct>()` methods
- **Named dependency resolution** via DI container

## Project Structure

```
templ-router/
├── cmd/trgen/                    # CLI template generator
├── demo/                        # Demo application
│   ├── app/                     # Template files
│   ├── generated/templates/     # Generated registry
│   ├── pkg/dataservices/        # Example data services
│   └── main.go                  # Demo application entry point
├── magefiles/                   # Build automation (Mage)
│   ├── main.go                  # Default targets and dev environment
│   ├── build.go                 # Build tasks (Templ, CSS, registry)
│   ├── test.go                  # Testing tasks (E2E, unit tests)
│   ├── generator.go            # Generator-specific tasks
│   └── docker.go                # Docker operations
├── pkg/                         # Library source code
│   ├── di/                      # Dependency injection container
│   ├── router/                 # Router core and middleware
│   ├── interfaces/             # Service interfaces
│   ├── config/                 # Configuration management
│   ├── services/               # Core services
│   └── shared/                 # Shared utilities
├── docs/                       # Documentation
├── .goreleaser.yml             # Release configuration
├── CHANGELOG.md                # Version history
├── README.md                  # Comprehensive documentation
├── CRUSH.md                   # Development guide
└── Dockerfile.release         # Production Docker image
```

## Build & Development Commands

### Essential Development Workflow
```bash
# Start full development environment (watches templates, CSS, registry, and runs server)
mage dev

# Generate templ templates from .templ files
mage build:templGenerate

# Generate template registry from file structure (after adding/removing templates)
mage build:registryGenerate

# Build Tailwind CSS
mage build:tailwindClean

# Clean all build artifacts
mage clean
```

### Testing Commands
```bash
# Run all tests (unit + E2E)
mage test:all

# Run E2E tests against running service
mage test:e2e

# Run E2E tests in watch mode
mage test:e2eWatch

# Quick smoke tests
mage test:e2eSmoke

# Specific test categories
mage test:e2eRouting     # Multi-language routing tests
mage test:e2eI18n        # Internationalization tests
mage test:e2eData        # Data service tests
mage test:e2eAuth        # Authentication tests
mage test:e2ePerf        # Performance tests

# Setup complete testing environment
mage test:devSetup
```

### Generator Development
```bash
# Install trgen tool from source
mage generator:install

# Run generator tests
mage generator:test

# Build trgen binary
mage generator:build

# Build for all platforms
mage build:all
```

### Docker Operations
```bash
# Start demo service
mage docker:up

# Stop service
mage docker:down

# View logs
mage docker:logs
```

## Template Generator (trgen)

**Purpose**: CLI tool that generates template registries from file structure

**Installation**:
```bash
go install github.com/denkhaus/templ-router/cmd/trgen@latest
```

**Usage** (must be run from application directory):
```bash
trgen --scan-path=app --module-name=github.com/youruser/yourproject
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch
```

**Environment Variables**:
- `TRGEN_SCAN_PATH`: Directory containing `.templ` files
- `TRGEN_MODULE_NAME`: Go module name from `go.mod`
- `TRGEN_WATCH_MODE`: Enable watch mode
- `TRGEN_WATCH_EXTENSIONS`: File extensions to watch

## Configuration System

**Environment Variables** with configurable prefix (default "TR"):

### Server Configuration
```bash
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8084
TR_SERVER_BASE_URL=http://localhost:8084
TR_SERVER_READ_TIMEOUT=30s
TR_SERVER_WRITE_TIMEOUT=30s
TR_SERVER_IDLE_TIMEOUT=120s
TR_SERVER_SHUTDOWN_TIMEOUT=30s
```

### Authentication Configuration
```bash
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_NAME=session_id
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/welcome
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/
```

### Internationalization Configuration
```bash
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en
TR_I18N_FALLBACK_LOCALE=en
```

### Router Configuration
```bash
TR_ROUTER_ENABLE_TRAILING_SLASH=true      # Redirect /path/ to /path
TR_ROUTER_ENABLE_SLASH_REDIRECT=true       # Clean double slashes
TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED=true   # 405 handler
```

### Middleware Configuration
All middleware is controlled via environment variables:
```bash
# Router middleware (URL normalization)
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true

# Authentication middleware
TR_AUTH_ENABLE_MIDDLEWARE=true             # Enable/disable auth middleware

# Internationalization middleware
TR_I18N_ENABLE_MIDDLEWARE=true              # Enable/disable i18n middleware

# Template middleware
TR_TEMPLATE_ENABLE_MIDDLEWARE=true          # Enable/disable template middleware

# Authentication routes
TR_ROUTER_ENABLE_AUTH_ROUTES=true           # Enable/disable authentication routes
TR_ROUTER_AUTH_ROUTE_PREFIX=/api            # Authentication route prefix
```

## Data Service Pattern

### Interface Definition
```go
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) // Optional specific method
}
```

### RouterContext Usage
```go
func (s *service) GetData(routerCtx interfaces.RouterContext) (*Data, error) {
    // URL parameters from routes like /{locale}/user/{id}
    locale := routerCtx.GetURLParam("locale")
    userID := routerCtx.GetURLParam("id")

    // Query parameters from ?page=5&filter=active
    page := routerCtx.GetQueryParam("page")
    filter := routerCtx.GetQueryParam("filter")

    return &Data{...}, nil
}
```

### Service Registration
```go
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
```

## Component Metadata Pattern

### Component YAML Structure
Components can have their own metadata files for self-contained configuration:

```yaml
# app/components/footer.templ.yaml
i18n:
  en:
    footer_copyright: "© 2024 My Company. All rights reserved."
    footer_privacy: "Privacy Policy"
    footer_contact: "Contact Us"
  de:
    footer_copyright: "© 2024 Meine Firma. Alle Rechte vorbehalten."
    footer_privacy: "Datenschutz"
    footer_contact: "Kontakt"

metadata:
  company_name: "My Company"
  company_email: "info@company.com"
  company_address: "123 Main Street, City"
  version: "1.0.0"

auth:
  type: "Public"
```

### Component Template Usage
Components can access their own metadata using the same functions as pages:

```go
// app/components/footer.templ
package components

import (
    "github.com/denkhaus/templ-router/pkg/router/i18n"
    "github.com/denkhaus/templ-router/pkg/router/metadata"
)

templ Footer() {
    companyEmail := metadata.M(ctx, "company_email")
    companyName := metadata.M(ctx, "company_name")

    <footer class="bg-gray-800 text-white p-4 mt-8">
        <div class="container mx-auto text-center">
            <p>{ i18n.T(ctx, "footer_copyright") }</p>
            <div class="flex justify-center space-x-4 mt-2 text-sm">
                <a href="/privacy" class="hover:underline">{ i18n.T(ctx, "footer_privacy") }</a>
                <a href={ "mailto:" + companyEmail } class="hover:underline">{ i18n.T(ctx, "footer_contact") }</a>
            </div>
            <div class="mt-2 text-xs text-gray-400">
                { companyName } - v{ metadata.M(ctx, "version") }
            </div>
        </div>
    </footer>
}
```

### Metadata Precedence Chain
The system follows a clear precedence order for metadata resolution:

1. **Component metadata** (highest priority) - when accessing `/components/footer`
2. **Page metadata** (middle priority) - when component is nested in a page
3. **Layout metadata** (lowest priority) - fallback for inherited settings

### Component Routes
Components are automatically accessible via their own routes:

- `/components/footer` - Renders the Footer component with its metadata
- `/components/navbar` - Renders the Navbar component with its metadata
- `/components/language-switcher` - Renders the LanguageSwitcher component

This enables:
- **HTMX partials** - Load components without layout for AJAX requests
- **Reusable components** - Self-contained metadata and i18n
- **API endpoints** - Direct component access for dynamic loading

## Internationalization System

### Translation File Structure
```yaml
# app/locale_/dashboard/page.templ.yaml
i18n:
  en:
    page_title: "Dashboard"
    stats:
      users: "Total Users"
      projects: "Active Projects"
  de:
    page_title: "Dashboard"
    stats:
      users: "Gesamte Benutzer"
      projects: "Aktive Projekte"

auth:
  type: "UserRequired"
  redirect_url: "/login"
```

### Template Usage
```go
// In templ files
{ i18n.T(ctx, "translation_key") }
{ i18n.T(ctx, "nested.key.with.dots") }
{ i18n.GetCurrentLocale(ctx) }
<a href={ i18n.LocalizeSafeURL(ctx, "/dashboard") }>Dashboard</a>
```

## Authentication System

### Auth Types
- `Public`: No authentication required (default)
- `UserRequired`: Any authenticated user can access
- `AdminRequired`: Only admin users can access

### Configuration
```yaml
# .templ.yaml file
auth:
  type: "AdminRequired"
  redirect_url: "/login"
  roles: ["admin", "super_admin"]  # Optional: specific roles required
```

### Built-in API Routes
- `POST /api/auth/signin` - User sign in
- `POST /api/auth/signout` - User sign out
- `POST /api/auth/signup` - User registration

## Critical Development Rules

### ⚠️ ABSOLUTELY CRITICAL - NEVER IGNORE

1. **NO Unicode/Emojis in Code**
   - NEVER use Unicode characters like 🔧 in Go code
   - ONLY use ASCII characters in all code files
   - REASON: Destroys files and build processes

2. **NEVER Manually Edit Generated Files**
   - NEVER edit `generated/templates/registry.go`
   - NEVER edit `*_templ.go` files
   - These are overwritten with every build
   - CORRECT: Modify generator templates in `cmd/trgen/`

3. **Template Generator Paths**
   - Generator is in `cmd/trgen/`
   - Install: `go install github.com/denkhaus/templ-router/cmd/trgen@main`
   - Execution: `trgen --scan-path=app --module-name=github.com/youruser/yourproject`

4. **NO Production Code Replacement**
   - This is a PRODUCTION PROJECT, not test example
   - Adapt real implementations for DI, don't replace them

5. **DI Container Rules**
   - ONLY provider registrations in the container
   - NO logic in the container

6. **Always View Code Before Editing**
   - NEVER find_and_replace without seeing current code
   - ALWAYS use view/expand_code_chunks first

7. **Use KNOT CLI for Work**
   - Use `knot` to structure work and track tasks
   - Create issues for complex tasks
   - NEVER mechanically process tasks without doing actual work

## Testing Strategy

### Unit Tests
- Standard Go testing with testify
- Mock dependencies for isolation
- Test service interfaces, not implementations

### E2E Tests
- Ginkgo/Gomega framework
- Agouti for browser automation
- Test against running Docker service
- Test categories:
  - Health checks and basic functionality
  - Multi-language routing
  - Authentication flows
  - Data service integration
  - Performance validation

## Performance Considerations

- Template caching is handled automatically
- Route discovery is optimized for startup
- Data services are resolved per request
- Use appropriate timeouts for external calls
- Monitor memory usage with large template sets

## Security Notes

- CSRF protection configurable via environment
- Rate limiting available for API endpoints
- Security headers enabled by default
- Input validation through parameter validators
- Session-based authentication with secure defaults

## Common Gotchas

1. **Generated File Overwrites**: Never edit generated files directly
2. **Unicode Characters**: Avoid all non-ASCII in Go source
3. **DI Container**: Only register providers, no business logic
4. **Template Registry**: Must regenerate after adding/removing templates
5. **Route Parameters**: Use `_` suffix for dynamic segments
6. **Context Access**: Always pass context through template functions
7. **Module Names**: trgen requires exact go.mod module name match

## CI/CD Pipeline

**GitHub Actions Workflow** (.github/workflows/ci.yml):
- Multi-version Go testing (1.24)
- Mage-based test execution
- GolangCI linting
- Security scanning (Gosec, govulncheck)
- Code coverage reporting
- Platform-agnostic building

## Release Management

**GoReleaser Configuration** (.goreleaser.yml):
- Multi-platform binary builds
- Docker image publishing (GHCR)
- Homebrew tap integration
- Automated changelog generation
- Structured release notes

This comprehensive architecture provides a solid foundation for building modern web applications with file-based routing, type-safe templating, and comprehensive middleware support. The library emphasizes clean architecture principles, dependency injection, and developer experience while maintaining production-ready standards.