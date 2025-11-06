# API Reference

**Quick reference guide for Templ Router APIs, functions, and configuration.**

## 🚀 Quick Reference

### Configuration Prefix Notice

**Important:** All environment variables in this documentation use the default prefix `TR_`. This prefix is **configurable** when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_SERVER_HOST=localhost`
- Custom: `MYAPP_SERVER_HOST=localhost`
- Multiple apps: `APP1_SERVER_HOST=localhost` and `APP2_SERVER_HOST=localhost`

All environment variable examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix.

### Core Setup

```go
// Main application setup
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/go-chi/chi/v5"
    "github.com/youruser/yourproject/generated/templates"
)

func main() {
    // Create DI container
    container := di.NewContainer()
    container.RegisterRouterServices("TR")  // "TR" is the default prefix

    // Create template registry
    templateRegistry, _ := templates.NewRegistry(container.GetInjector())
    container.RegisterApplicationServices(di.WithTemplateRegistry(templateRegistry))

    // Setup HTTP router
    mux := chi.NewRouter()
    router := container.GetRouter()
    router.Initialize()
    router.RegisterRoutes(mux)

    http.ListenAndServe(":8080", mux)
}
```

### Template Registry Generation

```bash
# Generate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Watch mode
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch

# Environment variables
TRGEN_SCAN_PATH=app
TRGEN_MODULE_NAME=github.com/youruser/yourproject
trgen
```

## 📁 File Structure

```
app/
├── layout.templ                    # Root layout
├── page.templ                      # Home page (/)
├── login/
│   ├── page.templ                  # Login page (/login)
│   └── page.templ.yaml            # Metadata
├── dashboard/
│   └── page.templ                  # Dashboard (/dashboard)
├── locale_/                        # Internationalized routes
│   ├── page.templ                  # /en, /de, /fr
│   └── user/
│       └── id_/
│           └── page.templ          # /en/user/123
└── components/
    ├── navbar.templ                # Component (/components/navbar)
    └── navbar.templ.yaml          # Component metadata
```

## 🛣️ Routing Patterns

### Route Generation

| File Structure | Generated Route | Description |
|---|---|---|
| `app/page.templ` | `/` | Root page |
| `app/login/page.templ` | `/login` | Static route |
| `app/user/id_/page.templ` | `/user/{id}` | Dynamic parameter |
| `app/locale_/page.templ` | `/{locale}` | Internationalized |
| `app/locale_/user/id_/page.templ` | `/{locale}/user/{id}` | Both i18n and dynamic |

### Route Precedence

1. **Specific routes**: `/admin/user/profile`
2. **Dynamic routes**: `/admin/user/{id}`
3. **Localized routes**: `/{locale}/admin/user/{id}`
4. **Fallback routes**: `/`

## 🔐 Authentication Types

### Configuration

```yaml
# In .templ.yaml files
auth:
  type: "Public"                    # No authentication
  type: "UserRequired"              # Any authenticated user
  type: "AdminRequired"             # Admin users only
  redirect_url: "/login"           # Redirect URL for unauthenticated
  roles: ["admin", "super_admin"] # Optional specific roles
```

### API Endpoints

```bash
POST /api/auth/signin     # User login
POST /api/auth/signout    # User logout
POST /api/auth/signup     # User registration
GET  /api/auth/me         # Current user info
```

## 🌍 Internationalization

### Configuration

```yaml
# In .templ.yaml files
i18n:
  en:
    welcome_message: "Welcome!"
    button_submit: "Submit"
  de:
    welcome_message: "Willkommen!"
    button_submit: "Absenden"
```

### Template Functions

```go
// Translation functions
i18n.T(ctx, "key")                        // Simple translation
i18n.T(ctx, "nested.key")                 // Nested translation
i18n.GetCurrentLocale(ctx)                 // Current locale
i18n.LocalizeSafeURL(ctx, "/dashboard")   // Localized URL

// Route information
i18n.GetCurrentRoute(ctx)                  // Full route with locale
i18n.GetCurrentRouteWithoutLocale(ctx)     // Route without locale
i18n.GetAvailableKeys(ctx)                 // Available translation keys
```

### Environment Variables

```bash
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en
TR_I18N_DETECTION_METHOD=url,cookie,header
TR_I18N_URL_PREFIX=true
```

## 📊 Data Services

### Service Interface

```go
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) // Optional
}
```

### RouterContext Methods

```go
// URL parameters (from routes like /{locale}/user/{id})
routerCtx.GetURLParam("locale")     // "en"
routerCtx.GetURLParam("id")         // "123"
routerCtx.GetAllURLParams()         // map[string]string

// Query parameters (from ?page=5&filter=active)
routerCtx.GetQueryParam("page")     // "5"
routerCtx.GetQueryParams("tag")     // ["go", "web"]
routerCtx.GetAllQueryParams()       // map[string][]string

// Context access
routerCtx.Context()                 // context.Context
routerCtx.Request()                 // *http.Request
routerCtx.ChiContext()              // *chi.Context
```

### Service Registration

```go
// Register named data services
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
```

### Template Integration

```go
// Template signature indicates data requirement
templ Page(user *UserData) {
    <h1>{ user.Name }</h1>
}

// Only one data service per template is supported
templ Page(data *CompositeDashboardData) {
    <h1>{ data.User.Name }</h1>
    <p>Total users: { data.SystemStats.TotalUsers }</p>
}
```

## 🎨 Template Functions

### Metadata Functions

```go
// Access metadata from .templ.yaml files
metadata.M(ctx, "page_title")        // Get metadata value
metadata.M(ctx, "theme")             // Get theme setting
```

### Template Structure

```go
// Layout template
templ Layout(title string, content templ.Component) {
    <!DOCTYPE html>
    <html>
    <head><title>{ title }</title></head>
    <body>{ content }</body>
    </html>
}

// Page template
templ HomePage() {
    <div>
        <h1>{ i18n.T(ctx, "welcome_title") }</h1>
        <p>{ metadata.M(ctx, "page_description") }</p>
    </div>
}
```

## ⚙️ Configuration Reference

### Server Configuration

```bash
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080
TR_SERVER_BASE_URL=http://localhost:8080
TR_SERVER_READ_TIMEOUT=30s
TR_SERVER_WRITE_TIMEOUT=30s
```

### Authentication Configuration

```bash
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_NAME=session_id
TR_AUTH_SESSION_COOKIE_SECURE=true
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard
TR_AUTH_CREATE_DEFAULT_ADMIN=false
```

### Internationalization Configuration

```bash
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en
TR_I18N_COOKIE_NAME=locale
TR_I18N_URL_PREFIX=true
```

### Router Configuration

```bash
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED=true
TR_ROUTER_SCAN_PATH=app
```

### Security Configuration

```bash
TR_SECURITY_CSRF_SECRET=your-secret-key
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_SECURITY_RATE_LIMIT_REQUESTS=100
TR_SECURITY_ENABLE_SECURITY_HEADERS=true
```

## 🛠️ CLI Commands

### Mage Commands

```bash
# Development
mage dev                    # Development server
mage dev:setup              # Setup environment

# Building
mage build                  # Build library
mage build:all              # Build all platforms
mage generator:build        # Build trgen

# Testing
mage test                   # Run tests
mage test:all               # All tests
mage test:unit              # Unit tests
mage test:e2e               # E2E tests
mage test:coverage          # Coverage report

# Code quality
mage fmt                    # Format code
mage lint                   # Run linter
mage vet                    # Run go vet

# Dependencies
mage deps:tidy              # Clean dependencies
mage deps:update            # Update dependencies
```

### trgen Commands

```bash
# Basic usage
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# With options
trgen \
  --scan-path=app \
  --module-name=github.com/youruser/yourproject \
  --output-dir=generated/templates \
  --package-name=templates \
  --watch

# Environment variables
TRGEN_SCAN_PATH=app
TRGEN_MODULE_NAME=github.com/youruser/yourproject
TRGEN_WATCH_MODE=true
trgen
```

### templ Commands

```bash
# Generate templates
templ generate

# Watch mode
templ generate --watch

# Specific file
templ generate path/to/template.templ
```

## 🔧 Dependency Injection

### Container Setup

```go
// Create DI container
container := di.NewContainer()
defer container.Shutdown()

// Register router services
container.RegisterRouterServices("TR")

// Register application services
container.RegisterApplicationServices(
    di.WithTemplateRegistry(templateRegistry),
    di.WithUserStore(userStore),
)
```

### Service Registration

```go
// Register named services
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)

// Register singleton services
do.Provide(injector, NewDatabaseConnection)
do.Provide(injector, NewCacheService)
```

### Service Resolution

```go
// Resolve named service
userService := do.MustInvokeNamed[UserDataService](injector, "UserDataService")

// Resolve singleton service
db := do.MustInvoke[*DatabaseConnection](injector)
cache := do.MustInvoke[CacheService](injector)
```

## 🎯 Common Patterns

### Basic Page with Authentication

```go
// app/dashboard/page.templ
package main

templ DashboardPage(user *UserData) {
    <div class="dashboard">
        <h1>Welcome, { user.Name }!</h1>
        <nav>
            <a href={ i18n.LocalizeSafeURL(ctx, "/profile") }>
                { i18n.T(ctx, "profile") }
            </a>
            <form method="POST" action="/api/auth/signout">
                <button type="submit">
                    { i18n.T(ctx, "sign_out") }
                </button>
            </form>
        </nav>
    </div>
}
```

```yaml
# app/dashboard/page.templ.yaml
metadata:
  page_title: "Dashboard"

i18n:
  en:
    welcome: "Welcome"
    profile: "Profile"
    sign_out: "Sign Out"
  de:
    welcome: "Willkommen"
    profile: "Profil"
    sign_out: "Abmelden"

auth:
  type: "UserRequired"
  redirect_url: "/login"

data_services:
  - "UserDataService"
```

### Internationalized Component

```go
// app/components/language-switcher.templ
package components

templ LanguageSwitcher() {
    currentLocale := i18n.GetCurrentLocale(ctx)
    routeWithoutLocale := i18n.GetCurrentRouteWithoutLocale(ctx)

    <div class="language-switcher">
        for _, locale := range []string{"en", "de", "fr"} {
            if locale == currentLocale {
                <span class="current">{ strings.ToUpper(locale) }</span>
            } else {
                <a href={ "/" + locale + routeWithoutLocale }>
                    { strings.ToUpper(locale) }
                </a>
            }
        }
    </div>
}
```

### Data Service Implementation

```go
// pkg/dataservices/user_service.go
package dataservices

type UserData struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserDataService interface {
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
}

type userDataServiceImpl struct {
    userRepo UserRepository
}

func NewUserDataService(injector *do.Injector) (UserDataService, error) {
    return &userDataServiceImpl{
        userRepo: do.MustInvoke[UserRepository](injector),
    }, nil
}

func (s *userDataServiceImpl) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")

    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }

    return &UserData{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

## 🐛 Error Handling

### Service Errors

```go
func (s *userService) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")

    if userID == "" {
        return nil, fmt.Errorf("user ID is required")
    }

    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, fmt.Errorf("user not found: %s", userID)
        }
        return nil, fmt.Errorf("failed to fetch user: %w", err)
    }

    return &UserData{ID: user.ID, Name: user.Name}, nil
}
```

### Template Error Handling

```go
templ UserProfilePage() {
    if user, err := getUserData(ctx); err != nil {
        <div class="error">
            <h2>{ i18n.T(ctx, "error_title") }</h2>
            <p>{ i18n.T(ctx, "user_not_found") }</p>
        </div>
    } else {
        <div class="profile">
            <h1>{ user.Name }</h1>
        </div>
    }
}
```

## 📋 Quick Cheatsheet

### Template File Extensions

| Extension | Purpose |
|---|---|
| `.templ` | Template files |
| `.templ.yaml` | Template metadata |
| `id_/` | Dynamic parameter directory |
| `locale_/` | Internationalization directory |

### Authentication Levels

| Type | Description |
|---|---|
| `Public` | No authentication required |
| `UserRequired` | Any authenticated user |
| `AdminRequired` | Admin users only |

### Common Environment Variables

| Variable | Default | Description |
|---|---|---|
| `TR_SERVER_PORT` | `8080` | Server port |
| `TR_AUTH_SESSION_EXPIRY` | `24h` | Session duration |
| `TR_I18N_DEFAULT_LOCALE` | `en` | Default language |
| `TR_LOGGING_LEVEL` | `info` | Log level |

### Essential Commands

```bash
# Start development
mage dev

# Generate templates
templ generate
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Run tests
mage test:all

# Build project
mage build:all
```

## 🔗 Quick Links

- **[Getting Started](GETTING-STARTED.md)** - Complete setup guide
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Routing system details
- **[Authentication](AUTHENTICATION.md)** - Authentication system
- **[Internationalization](INTERNATIONALIZATION.md)** - Multi-language support
- **[Data Services](DATA-SERVICES.md)** - Data service patterns
- **[Configuration](CONFIGURATION.md)** - Full configuration reference

---

**Need more help? Check the [main documentation](../README.md) or open a [discussion](https://github.com/denkhaus/templ-router/discussions).**