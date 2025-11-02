# CRUSH.md - Templ Router Development Guide

## Project Overview

Templ Router is a Go library for file-based routing with [templ](https://templ.guide/) templates, dependency injection, and comprehensive middleware support. This is a **library**, not a standalone application - developers add it as a dependency to their own Go projects.

**Key Technologies:**
- Go 1.24+ with templ templates
- samber/do/v2 for dependency injection
- chi/v5 for HTTP routing
- YAML-based configuration and i18n
- Mage for build automation

## Essential Commands

### Development Workflow
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
mage test:e2eRouting
mage test:e2eI18n
mage test:e2eData
mage test:e2eAuth

# Setup complete testing environment
mage test:devSetup
```

### Library Development
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

## Project Architecture

### Core Components

**Router Core (`pkg/router/clean_router_core.go`)**
- Implements clean architecture with separated concerns
- Manages route discovery, handler building, and middleware setup
- Uses dependency injection throughout

**Dependency Injection (`pkg/di/`)**
- Based on samber/do/v2
- Container-based service management
- Separate registration for router services and application services

**Template Generator (`cmd/trgen/`)**
- CLI tool for generating template registries
- Scans file structure and maps templates to routes
- Configuration-agnostic (no hardcoded paths/names)

**File-Based Routing**
- Routes generated from directory structure
- Dynamic parameters with `_` suffix (e.g., `id_/`, `locale_/`)
- Template metadata via `.templ.yaml` files

### Template System

**Templ Files (`.templ`)**
- Use Go templating syntax
- Import `github.com/denkhaus/templ-router/pkg/router/i18n` for translations
- Package declaration based on directory path

**Metadata Files (`.templ.yaml`)**
- Authentication configuration (`auth.type`)
- Internationalization translations (`i18n`)
- Template metadata (`metadata`)
- Dynamic parameter validation (`dynamic.parameters`)

## Critical Rules from .agent.md

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

## Data Services Pattern

Data services provide template data through dependency injection:

**Interface Pattern:**
```go
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
}
```

**RouterContext Access:**
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

**Registration:**
```go
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
```

## Internationalization (i18n)

**Translation Access:**
```go
// In templates
{ i18n.T(ctx, "translation_key") }
{ i18n.T(ctx, "nested.key.with.dots") }

// Get current locale
{ i18n.GetCurrentLocale(ctx) }

// Localized URLs
<a href={ i18n.LocalizeSafeURL(ctx, "/dashboard") }>Dashboard</a>
```

**Translation File Structure:**
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

## Authentication System

**Auth Types:**
- `Public`: No authentication required
- `UserRequired`: Any authenticated user
- `AdminRequired`: Only admin users

**Configuration:**
```yaml
auth:
  type: "AdminRequired"
  redirect_url: "/login"
```

**Built-in API Routes:**
- `POST /api/auth/signin` - User sign in
- `POST /api/auth/signout` - User sign out  
- `POST /api/auth/signup` - User registration

## Environment Configuration

Use configurable prefix (set in `RegisterRouterServices("TR")`):

```bash
# Server
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080

# Authentication
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard

# Internationalization
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en

# Router Features
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED=true
```

## File Structure Conventions

**Template Directory Structure:**
```
app/
├── layout.templ              # Root layout
├── page.templ               # Home page
├── login/
│   ├── page.templ
│   └── page.templ.yaml      # Metadata
└── locale_/                  # Internationalized routes
    ├── admin/
    │   ├── page.templ
    │   └── page.templ.yaml
    └── product/
        └── id_/             # Dynamic parameter
            ├── page.templ
            └── page.templ.yaml
```

**Generated Files:**
```
generated/
└── templates/
    └── registry.go           # Auto-generated template registry
```

## Code Style and Patterns

**Naming Conventions:**
- Interfaces: `ServiceName` (e.g., `UserDataService`)
- Implementations: `serviceNameImpl` (e.g., `userDataServiceImpl`)
- Methods: `PascalCase` for public, `camelCase` for private
- Files: `snake_case.go`

**DI Pattern:**
```go
// Service creation
func NewUserService(injector do.Injector) interfaces.UserService {
    return &userService{
        config: do.MustInvoke[interfaces.ConfigService](injector),
        logger: do.MustInvoke[*zap.Logger](injector),
    }
}

// Registration
do.Provide(injector, NewUserService)
```

**Error Handling:**
- Use structured errors with context
- Wrap errors with `fmt.Errorf("operation failed: %w", err)`
- Log errors appropriately for the context

## Testing Strategy

**Unit Tests:**
- Standard Go testing with testify
- Mock dependencies for isolation
- Test service interfaces, not implementations

**E2E Tests:**
- Ginkgo/Gomega framework
- Agouti for browser automation
- Test against running Docker service
- Focus on user workflows and integration

**Test Categories:**
- Health checks and basic functionality
- Multi-language routing
- Authentication flows
- Data service integration
- Performance validation

## Development Tools Integration

**Required Tooling:**
```bash
# Install required tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/denkhaus/templ-router/cmd/trgen@latest
go install github.com/magefile/mage@latest
go install github.com/air-verse/air@latest  # For hot reload
```

**VS Code / Editor Support:**
- templ language server for template syntax
- Go extension for Go development
- YAML support for .templ.yaml files

## Common Gotchas

1. **Generated File Overwrites**: Never edit generated files directly
2. **Unicode Characters**: Avoid all non-ASCII in Go source
3. **DI Container**: Only register providers, no business logic
4. **Template Registry**: Must regenerate after adding/removing templates
5. **Route Parameters**: Use `_` suffix for dynamic segments
6. **Context Access**: Always pass context through template functions
7. **Module Names**: trgen requires exact go.mod module name match

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

---

**Remember**: This guide documents the actual observed patterns and commands. Never invent commands or patterns - verify against the codebase and existing configurations.