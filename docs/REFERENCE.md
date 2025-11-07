# API Reference

**Quick reference for Templ Router APIs, functions, and configuration.**

## Quick Setup

```go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/go-chi/chi/v5"
    "github.com/youruser/yourproject/generated/templates"
)

func main() {
    container := di.NewContainer()
    container.RegisterRouterServices("TR")  // Configurable prefix

    templateRegistry, _ := templates.NewRegistry(container.GetInjector())
    container.RegisterApplicationServices(di.WithTemplateRegistry(templateRegistry))

    mux := chi.NewRouter()
    router := container.GetRouter()
    router.Initialize()
    router.RegisterRoutes(mux)

    http.ListenAndServe(":8080", mux)
}
```

## Configuration

### Environment Variables

Use configurable prefix (default `TR`):

```bash
# Server
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8084
TR_SERVER_BASE_URL=http://localhost:8084

# Authentication
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard

# Internationalization
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en

# Router
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
```

### Custom Prefix

```go
container.RegisterRouterServices("MYAPP")  // Uses MYAPP_* variables
```

## Core APIs

### DI Container

```go
container := di.NewContainer()
container.RegisterRouterServices("TR")
container.RegisterApplicationServices(di.WithTemplateRegistry(registry))

router := container.GetRouter()
injector := container.GetInjector()
```

### Template Registry

```go
type Registry interface {
    GetTemplate(path string, params map[string]string) (templ.Component, error)
    GetAllRoutes() []string
    GetRouteInfo(path string) (*RouteInfo, error)
}

registry := templates.NewRegistry(injector)
```

### Router

```go
router := container.GetRouter()
router.Initialize()
router.RegisterRoutes(mux)

// Custom middleware
router.UseMiddleware(customMiddleware)
```

## Template System

### Template Types

```go
// Page template (generates route)
templ Page() {
    <h1>Home</h1>
}

// Layout template (no route)
templ Layout(content templ.Component) {
    <!DOCTYPE html>
    <html>
    <head>
        { title := metadata.M(ctx, "title") }
        if title != "" {
            <title>{ title }</title>
        } else {
            <title>{ i18n.T(ctx, "site_title") }</title>
        }
    </head>
    <body>{ content }</body>
    </html>
}

// Error component (no route)
templ Error() {
    <div class="error">
        <h1>Error Occurred</h1>
        <p>Something went wrong.</p>
    </div>
}

// Component template (no route)
templ Header(title string) {
    <header><h1>{ title }</h1></header>
}
```

### Template Metadata

```yaml
# page.templ.yaml
metadata:
  title: "Page Title"
  description: "Page description"

auth:
  type: "UserRequired"  # Public, UserRequired, AdminRequired
  redirect_url: "/login"

i18n:
  en:
    welcome: "Welcome"
  de:
    welcome: "Willkommen"
```

## File-Based Routing

### Route Generation

```
app/
├── page.templ                → /
├── about.templ               → /about
├── login/
│   └── page.templ           → /login
├── user/
│   └── id_/page.templ       → /user/{id}
└── locale_/
    ├── page.templ           → /{locale}
    └── dashboard/
        └── page.templ       → /{locale}/dashboard
```

### Dynamic Parameters

Use `_` suffix for dynamic segments:

```go
// Access in data service
func (s *service) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
    id := routerCtx.GetURLParam("id")
    locale := routerCtx.GetURLParam("locale")
    page := routerCtx.GetQueryParam("page")
    return s.getUser(id, locale, page)
}
```

## Data Services

### Service Interface

```go
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
}
```

### Service Registration

```go
// Register named service
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)

// Service implementation
func (s *userService) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")
    return s.repository.GetUser(userID)
}
```

### RouterContext

```go
type RouterContext interface {
    GetURLParam(key string) string
    GetQueryParam(key string) string
    GetContext() context.Context
    GetRequest() *http.Request
    GetResponse() http.ResponseWriter
}
```

## Internationalization

### Template Usage

```go
templ Page() {
    <h1>{ i18n.T(ctx, "welcome_title") }</h1>
    <p>{ i18n.T(ctx, "welcome_message") }</p>
    <a href={ i18n.LocalizeSafeURL(ctx, "/dashboard") }>Dashboard</a>
    <p>Current locale: { i18n.GetCurrentLocale(ctx) }</p>
}
```

### Translation Files

```yaml
# locale_/page.templ.yaml
i18n:
  en:
    welcome_title: "Welcome"
    welcome_message: "Welcome to our application"
  de:
    welcome_title: "Willkommen"
    welcome_message: "Willkommen in unserer Anwendung"
```

### Locale Detection

1. URL path (`/en/page`, `/de/page`)
2. Query parameter (`?locale=en`)
3. Accept-Language header
4. Cookie (`locale=en`)
5. Default locale

## Authentication

### Auth Types

```yaml
auth:
  type: "Public"           # No authentication
  type: "UserRequired"     # Any authenticated user
  type: "AdminRequired"    # Admin users only
  redirect_url: "/login"   # Redirect for unauthenticated users
  roles: ["admin", "moderator"]  # Specific roles
```

### Auth Endpoints

- `POST /api/auth/signin` - User sign in
- `POST /api/auth/signout` - User sign out
- `POST /api/auth/signup` - User registration

### Session Management

```go
// Session configuration
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_NAME=session_id
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard
```

## Middleware

### Built-in Middleware

```bash
# Enable/disable middleware
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_AUTH_ENABLE_MIDDLEWARE=true
TR_I18N_ENABLE_MIDDLEWARE=true
TR_TEMPLATE_ENABLE_MIDDLEWARE=true
```

### Custom Middleware

```go
func CustomMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Custom logic
        next.ServeHTTP(w, r)
    })
}

router.UseMiddleware(CustomMiddleware)
```

## CLI Tools

### trgen (Template Generator)

```bash
# Generate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Watch mode
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch

# Environment variables
TRGEN_SCAN_PATH=app
TRGEN_MODULE_NAME=github.com/youruser/yourproject
TRGEN_WATCH_MODE=true
```

### templ (Template Compiler)

```bash
# Generate Go files from .templ files
templ generate

# Watch mode
templ generate --watch

# Clean generated files
templ clean
```

## Error Handling

### Error Pages

```yaml
# error.templ.yaml
metadata:
  title: "Error Page"

i18n:
  en:
    error_title: "Error Occurred"
    error_message: "Something went wrong"
  de:
    error_title: "Fehler Aufgetreten"
    error_message: "Etwas ist schief gelaufen"
```

### Custom Error Handling

```go
router.SetErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
    // Custom error handling
    log.Printf("Error: %v", err)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
})
```

## Testing

### Unit Testing

```go
func TestDataService(t *testing.T) {
    service := NewUserService(mockRepo)

    mockCtx := &mockRouterContext{
        urlParams: map[string]string{"id": "123"},
    }

    data, err := service.GetData(mockCtx)
    assert.NoError(t, err)
    assert.NotNil(t, data)
}
```

### Integration Testing

```go
func TestRouter(t *testing.T) {
    router := setupTestRouter()

    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

## Performance

### Template Caching

```go
// Templates are compiled to Go code for performance
// Registry is cached for fast lookups
// Route matching is optimized for O(1) complexity
```

### Configuration Optimization

```bash
# Production settings
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
TR_TEMPLATE_ENABLE_MIDDLEWARE=true
```

## Security

### Built-in Security

- CSRF protection
- Rate limiting
- Security headers
- Input validation
- Secure session management

### Security Configuration

```bash
# Security settings
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_SECURE=true
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true
```

## Deployment

### Docker Configuration

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN templ generate && trgen --scan-path=app --module-name=github.com/youruser/yourproject
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/generated ./generated
EXPOSE 8080
CMD ["./main"]
```

### Environment Configuration

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    environment:
      - TR_SERVER_HOST=0.0.0.0
      - TR_SERVER_PORT=8080
      - TR_I18N_SUPPORTED_LOCALES=en,de,fr
    ports:
      - "8080:8080"
```

## Common Patterns

### Layout Inheritance

```go
// base layout
templ Layout(content templ.Component) {
    <!DOCTYPE html>
    <html>
    <head>
        { title := metadata.M(ctx, "title") }
        if title != "" {
            <title>{ title }</title>
        } else {
            <title>{ i18n.T(ctx, "site_title") }</title>
        }
    </head>
    <body>
        { Header() }
        <main>{ content }</main>
        { Footer() }
    </body>
    </html>
}

// page content component (not a page template)
templ AboutContent() {
    <h1>About Us</h1>
    <p>Company information...</p>
}

// Note: Pages use Page() function name, layout inheritance is automatic
// The Layout() template is used by the router automatically
```

### Component Reuse

```go
// Reusable component
templ Alert(message string, alertType string) {
    <div class={ "alert alert-" + alertType }>
        { message }
    </div>
}

// Using component in a page
templ Page() {
    Alert("Welcome!", "success")
    Alert("Error occurred", "error")
}
```

### Data Service Pattern

```go
// Service interface
type BlogService interface {
    GetData(routerCtx interfaces.RouterContext) (*BlogData, error)
}

// Implementation
func (s *blogService) GetData(routerCtx interfaces.RouterContext) (*BlogData, error) {
    category := routerCtx.GetURLParam("category")
    page := routerCtx.GetQueryParam("page")
    return s.repository.GetPosts(category, page)
}

// Template usage
templ Page() {
    data, _ := blogService.GetData(routerCtx)
    for _, post := range data.Posts {
        <article>{ post.Title }</article>
    }
}
```

---

**Related Documentation**: [Getting Started](GETTING-STARTED.md), [Architecture](ARCHITECTURE.md), [Configuration](CONFIGURATION.md)