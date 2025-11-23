# Architecture Overview

**Clean architecture principles with dependency injection and file-based routing.**

## Overview

Templ Router follows clean architecture with clear separation of concerns, dependency injection, and modular design. Built as a Go library for integration into applications.

**Key Principles:**
- **Clean Architecture** - Separation of infrastructure, application, and domain layers
- **Dependency Injection** - Type-safe service management with samber/do/v2
- **File-Based Routing** - Convention over configuration for route generation
- **Template System** - Type-safe templating with metadata management
- **Middleware Pipeline** - Configurable request processing chain

## Configuration

### Environment Variables

Router services use a configurable prefix (default: `TR`):

```go
// Default prefix
container.RegisterRouterServices("TR")  // Uses TR_* variables

// Custom prefix
container.RegisterRouterServices("MYAPP")  // Uses MYAPP_* variables
```

**Examples:**
- Default: `TR_SERVER_HOST=localhost`
- Custom: `MYAPP_SERVER_HOST=localhost`

## System Architecture

### High-Level Structure

```
┌─────────────────────────────────────────┐
│              Application Layer           │
├─────────────────────────────────────────┤
│  HTTP Routes  │ Auth Endpoints │ APIs   │
├─────────────────────────────────────────┤
│            Router Layer                 │
├─────────────────────────────────────────┤
│  Template Engine │ Middleware │ DI     │
├─────────────────────────────────────────┤
│           Infrastructure Layer          │
└─────────────────────────────────────────┘
```

### Core Components

#### 1. Router Layer
- **HTTP Router**: Chi-based routing with middleware support
- **Template Registry**: Auto-generated template mappings
- **Route Handler**: Request routing and template rendering
- **Middleware Pipeline**: Configurable request processing

#### 2. Template System
- **Template Engine**: Go-based templating with type safety
- **Metadata System**: YAML-based configuration
- **Component System**: Reusable template components
- **Layout Inheritance**: Template composition and inheritance

#### 3. Dependency Injection
- **DI Container**: Service registration and resolution
- **Service Management**: Type-safe service injection
- **Named Dependencies**: Service identification and resolution
- **Lifecycle Management**: Service initialization and cleanup

#### 4. Configuration System
- **Environment Variables**: Runtime configuration
- **Structured Config**: YAML-based configuration files
- **Validation**: Configuration validation and defaults
- **Hot Reload**: Runtime configuration updates

## File-Based Routing

### Route Generation

Templ Router generates routes from file structure:

```
app/
├── layout.templ           → Base layout
├── page.templ             → /
├── login/
│   └── page.templ         → /login
└── locale_/
    ├── page.templ         → /{locale}
    └── dashboard/
        └── page.templ     → /{locale}/dashboard
```

### Dynamic Parameters

Use `_` suffix for dynamic segments:

- `id_/` → `{id}`
- `locale_/` → `{locale}`
- `category_/slug_/` → `{category}/{slug}`

### Route Precedence

1. **Static routes** - `/about`
2. **Dynamic routes** - `/user/{id}`
3. **Wildcard routes** - `/{fallback}`
4. **Localized routes** - `/{locale}/page`

## Template System

### Template Types

#### Page Templates
```go
// app/page.templ
package main

templ Page() {
    <h1>Welcome</h1>
}
```

#### Layout Templates
```go
// app/layout.templ
package main

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
```

#### Component Templates
```go
// app/components/header.templ
package components

templ Header(title string) {
    <header><h1>{ title }</h1></header>
}
```

### Metadata System

Templates can have associated `.templ.yaml` files:

```yaml
# app/page.templ.yaml
metadata:
  title: "Home Page"
  description: "Welcome page"

auth:
  type: "Public"

i18n:
  en:
    welcome: "Welcome"
  de:
    welcome: "Willkommen"
```

### Component Metadata

Components can have their own metadata:

```yaml
# app/components/footer.templ.yaml
i18n:
  en: { footer_copyright: "© 2024 My Company" }
  de: { footer_copyright: "© 2024 Meine Firma" }

metadata:
  company_name: "My Company"
  version: "1.0.0"
```

## Dependency Injection

### DI Container

```go
// Initialize container
container := di.NewContainer()

// Register router services
container.RegisterRouterServices("TR")

// Register application services
container.RegisterApplicationServices(
    di.WithTemplateRegistryFactory(func(injector interface{}) (interface{}, error) {
        return registry, nil
    }),
)

// Get router instance
router := container.GetRouter()
```

### Service Registration

```go
// Data service registration
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)

// Custom service registration
container.RegisterService("customService", customService)
```

### Service Resolution

```go
// Get service by name
userService := injector.MustInvokeNamed("UserDataService").(UserDataService)

// Get router
router := container.GetRouter()
```

## Middleware System

### Built-in Middleware

1. **URL Normalization** - Trailing slashes and path cleaning
2. **Authentication** - Session-based auth and role checking
3. **Internationalization** - Locale detection and routing
4. **Template Rendering** - Template selection and rendering
5. **Error Handling** - Error page rendering and logging

### Custom Middleware

```go
// Custom middleware
func CustomMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Custom logic
        next.ServeHTTP(w, r)
    })
}

// Register middleware
router.UseMiddleware(CustomMiddleware)
```

### Middleware Configuration

```bash
# Enable/disable middleware via environment
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_AUTH_ENABLE_MIDDLEWARE=true
TR_I18N_ENABLE_MIDDLEWARE=true
TR_TEMPLATE_ENABLE_MIDDLEWARE=true
```

## Data Services

### Service Interface

```go
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
}
```

### RouterContext

```go
// Access URL parameters
locale := routerCtx.GetURLParam("locale")
userID := routerCtx.GetURLParam("id")

// Access query parameters
page := routerCtx.GetQueryParam("page")
filter := routerCtx.GetQueryParam("filter")
```

### Service Implementation

```go
func (s *userService) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
    // Get parameters
    userID := routerCtx.GetURLParam("id")

    // Fetch data
    return s.repository.GetUser(userID)
}
```

## Internationalization

### Locale Detection

1. **URL Path** - `/en/page`, `/de/page`
2. **Query Parameter** - `?locale=en`
3. **Header** - `Accept-Language` header
4. **Cookie** - `locale` cookie
5. **Default** - Configured default locale

### Translation Files

```yaml
# app/locale_/page.templ.yaml
i18n:
  en:
    page_title: "Home"
    welcome_message: "Welcome to our site"
  de:
    page_title: "Startseite"
    welcome_message: "Willkommen auf unserer Seite"
```

### Template Usage

```go
// In templates
{ i18n.T(ctx, "welcome_message") }
{ i18n.LocalizeSafeURL(ctx, "/dashboard") }
{ i18n.GetCurrentLocale(ctx) }
```

## Authentication System

### Auth Types

- **Public** - No authentication required
- **UserRequired** - Any authenticated user
- **AdminRequired** - Admin users only

### Session Management

```go
// Session configuration
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_NAME=session_id
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard
```

### Auth Endpoints

- `POST /api/auth/signin` - User sign in
- `POST /api/auth/signout` - User sign out
- `POST /api/auth/signup` - User registration

## Configuration Management

### Environment Variables

```bash
# Server configuration
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8084
TR_SERVER_BASE_URL=http://localhost:8084

# Authentication
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard

# Internationalization
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en

# Router settings
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
```

### Configuration Files

```yaml
# config.yaml
server:
  host: localhost
  port: 8084
  base_url: http://localhost:8084

auth:
  session_expiry: 24h
  signin_success_route: /dashboard

i18n:
  supported_locales: [en, de, fr]
  default_locale: en
```

## Performance Considerations

### Template Caching

- **Compile-time Generation** - Templates compiled to Go code
- **Runtime Caching** - Template instances cached
- **Memory Management** - Efficient template lifecycle

### Route Optimization

- **Route Precomputation** - Routes generated at startup
- **Fast Matching** - Efficient route matching algorithm
- **Parameter Parsing** - Optimized parameter extraction

### Database Connections

- **Connection Pooling** - Efficient database connection management
- **Query Optimization** - Optimized data service queries
- **Caching Strategy** - Result caching for frequently accessed data

## Security Features

### Built-in Security

- **CSRF Protection** - Cross-site request forgery prevention
- **Rate Limiting** - Request rate limiting
- **Security Headers** - Security-focused HTTP headers
- **Input Validation** - Parameter validation and sanitization

### Session Security

- **Secure Cookies** - HTTP-only, secure cookies
- **Session Expiration** - Configurable session lifetime
- **Session Invalidation** - Secure session termination

## Testing Architecture

### Unit Testing

```go
func TestUserDataService(t *testing.T) {
    service := NewUserService(mockRepo)
    data, err := service.GetData(mockRouterCtx)
    assert.NoError(t, err)
    assert.NotNil(t, data)
}
```

### Integration Testing

```go
func TestRouterIntegration(t *testing.T) {
    router := setupTestRouter()
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

### E2E Testing

```go
func TestUserFlow(t *testing.T) {
    // Start test server
    server := startTestServer()
    defer server.Close()

    // Test user login flow
    client := &http.Client{}
    // ... test implementation
}
```

## Deployment Architecture

### Production Deployment

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    environment:
      - TR_SERVER_HOST=0.0.0.0
      - TR_SERVER_PORT=8080
    ports:
      - "8080:8080"
```

### Scaling Considerations

- **Horizontal Scaling** - Multiple application instances
- **Load Balancing** - Distribute requests across instances
- **Database Scaling** - Read replicas and connection pooling
- **Caching Layer** - Redis or similar for session storage

## Best Practices

### Code Organization

1. **Clear Layer Separation** - Maintain clean architecture boundaries
2. **Interface-Driven Design** - Program to interfaces, not implementations
3. **Dependency Injection** - Use DI container for service management
4. **Configuration Management** - Externalize configuration

### Performance Optimization

1. **Template Optimization** - Efficient template design
2. **Database Optimization** - Query optimization and indexing
3. **Caching Strategy** - Appropriate caching layers
4. **Resource Management** - Efficient resource utilization

### Security Best Practices

1. **Input Validation** - Validate all user inputs
2. **Secure Defaults** - Secure by default configuration
3. **Regular Updates** - Keep dependencies updated
4. **Security Testing** - Regular security audits

---

**Related Documentation**: [File-Based Routing](FILE-BASED-ROUTING.md), [Data Services](DATA-SERVICES.md), [Configuration](CONFIGURATION.md)