# Templ Router

[![Go Version](https://img.shields.io/github/go-mod/go-version/denkhaus/templ-router)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/denkhaus/templ-router)](https://goreportcard.com/report/github.com/denkhaus/templ-router)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/github/actions/workflow/status/denkhaus/templ-router/ci.yml?branch=main)](https://github.com/denkhaus/templ-router/actions)
[![Coverage](https://img.shields.io/codecov/c/github/denkhaus/templ-router)](https://codecov.io/gh/denkhaus/templ-router)
[![GitHub Release](https://img.shields.io/github/v/release/denkhaus/templ-router)](https://github.com/denkhaus/templ-router/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/denkhaus/templ-router.svg)](https://pkg.go.dev/github.com/denkhaus/templ-router)

**A powerful, file-based routing framework for Go applications using [templ](https://templ.guide/) templates with Next.js-inspired features.**

Templ Router brings modern web development patterns to Go, offering file-based routing, internationalization, authentication, layout inheritance, and comprehensive configuration management - all with type safety and excellent performance.

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

## Key Features

### File-Based Routing
- **Zero Configuration**: Routes automatically generated from your file structure
- **Dynamic Parameters**: Support for `[id]`, `[slug]`, and custom parameter patterns
- **Nested Routes**: Unlimited nesting with automatic route hierarchy
- **Route Discovery**: Intelligent scanning and registration of template files

### Advanced Internationalization (i18n)
- **Multi-Language Support**: Built-in locale handling with `[locale]` directory structure
- **Translation Management**: YAML-based translations with fallback support
- **Locale Detection**: Automatic locale detection from URLs and headers
- **Translation Keys**: Structured translation system with dot notation

### Comprehensive Authentication
- **Role-Based Access Control**: Public, User, and Admin authentication levels
- **Route-Level Security**: Per-template authentication configuration
- **Session Management**: Secure session handling with configurable expiry
- **Default Admin Setup**: Automatic admin user creation for development

### Next.js-Style Layout System
- **Layout Inheritance**: Nested layouts with automatic composition
- **Layout Templates**: Reusable layout components across routes
- **Error Boundaries**: Custom error templates with fallback handling
- **Template Composition**: Automatic layout wrapping and content injection

### Powerful Configuration Service
- **Environment-Based Config**: Comprehensive configuration via environment variables
- **Per-Template Settings**: Individual template configuration via YAML metadata
- **Security Settings**: CSRF protection, rate limiting, and security headers
- **Development Tools**: Hot reloading, debug modes, and development helpers

### Dependency Injection Container
- **Service Registration**: Clean dependency management with type safety
- **Data Service Integration**: Automatic data service resolution per route
- **Interface-Based Design**: Modular architecture with clear contracts
- **Lifecycle Management**: Proper service initialization and cleanup

### Performance & Developer Experience
- **Template Caching**: Intelligent caching with automatic invalidation
- **Hot Reloading**: Development server with live template updates
- **Type Safety**: Full Go type safety with templ integration
- **CLI Tools**: Code generation and project scaffolding utilities

## Quick Start

### Installation

```bash
# Clone the starter template
git clone https://github.com/denkhaus/templ-router.git my-app
cd my-app

# Install dependencies
go mod tidy

# Start development server
make dev
```

Your application is now running at [http://localhost:7331](http://localhost:7331)

### Basic Project Structure

```
my-app/
├── app/                          # Template directory
│   ├── layout.templ             # Root layout
│   ├── page.templ               # Home page
│   ├── login/
│   │   ├── page.templ           # Login page
│   │   └── page.templ.yaml      # Page configuration
│   └── [locale]/                # Internationalized routes
│       ├── admin/
│       │   ├── page.templ       # Admin dashboard
│       │   └── page.templ.yaml  # Auth: admin required
│       └── product/
│           └── [id]/
│               ├── page.templ   # Dynamic product page
│               └── page.templ.yaml
├── assets/                      # Static assets
├── pkg/                         # Application logic
└── main.go                      # Application entry point
```

## Core Concepts

### File-Based Routing

Routes are automatically generated from your file structure:

```
app/page.templ                    → /
app/login/page.templ              → /login
app/[locale]/page.templ           → /en, /de, /fr (based on config)
app/[locale]/product/[id]/page.templ → /en/product/123
```

### Template Configuration

Each template can have an optional `.templ.yaml` configuration file:

```yaml
# app/locale_/admin/page.templ.yaml
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
```

### Layout System

Layouts work like Next.js - each directory can have a `layout.templ`:

```go
// app/layout.templ
package app

templ Layout(title string) {
    <!DOCTYPE html>
    <html>
        <head><title>{title}</title></head>
        <body>
            { children... }  // Content injected here
        </body>
    </html>
}
```

### Data Services

Integrate data services seamlessly:

```go
// Define your data service
type ProductService interface {
    GetProduct(id string) (*Product, error)
}

// Use in templates
func ProductPage(productService ProductService, id string) templ.Component {
    product, _ := productService.GetProduct(id)
    return productPageTemplate(product)
}
```

## Internationalization

### Locale Configuration

```bash
# Environment variables
I18N_SUPPORTED_LOCALES=en,de,fr,es
I18N_DEFAULT_LOCALE=en
I18N_FALLBACK_LOCALE=en
```

### Translation Files

```yaml
# app/locale_/dashboard/page.templ.yaml
i18n:
  en:
    page_title: "Dashboard"
    page_subtitle: "Overview of your application metrics and recent activity"
    stats_users: "Total Users"
    recent_activity_title: "Recent Activity"
  de:
    page_title: "Dashboard"
    page_subtitle: "Übersicht Ihrer Anwendungsmetriken und aktuellen Aktivitäten"
    stats_users: "Gesamte Benutzer"
    recent_activity_title: "Letzte Aktivitäten"

auth:
  type: "Public"

metadata:
  title: "Dashboard - Multi-Language Demo"
  theme: "dashboard"
  description: "Application dashboard with metrics"
```

### Using Translations in Templates

```go
templ AdminPage() {
    <h1>{ t("page_title") }</h1>
    <p>{ t("admin_warning") }</p>
}
```

## Authentication & Security

### Route-Level Authentication

```yaml
# app/locale_/admin/page.templ.yaml
auth:
  type: "AdminRequired"    # Public, UserRequired, AdminRequired
  redirect_url: "/login"

# app/login/page.templ.yaml
auth:
  type: "Public"
```

### Security Configuration

```bash
# Environment variables
SECURITY_CSRF_SECRET=your-secret-key
SECURITY_ENABLE_RATE_LIMIT=true
SECURITY_RATE_LIMIT_REQUESTS=100
SECURITY_ENABLE_SECURITY_HEADERS=true
```

### Session Management

```bash
AUTH_SESSION_EXPIRY=24h
AUTH_SESSION_SECURE=true
AUTH_SESSION_HTTP_ONLY=true
AUTH_SESSION_SAME_SITE=strict
```

## Configuration

### Environment Variables

The router supports comprehensive configuration via environment variables:

```bash
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_BASE_URL=http://localhost:8080

# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=router_db

# Authentication
AUTH_CREATE_DEFAULT_ADMIN=true
AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
AUTH_DEFAULT_ADMIN_PASSWORD=admin123
AUTH_MIN_PASSWORD_LENGTH=8

# Internationalization
I18N_SUPPORTED_LOCALES=en,de,fr
I18N_DEFAULT_LOCALE=en

# Layout System
LAYOUT_ROOT_DIRECTORY=app
LAYOUT_ENABLE_INHERITANCE=true
LAYOUT_TEMPLATE_EXTENSION=.templ

# Logging
LOGGING_LEVEL=info
LOGGING_FORMAT=json
LOGGING_OUTPUT=stdout
```

### Template Metadata

Each template can have rich metadata configuration:

```yaml
# app/locale_/product/id_/page.templ.yaml
auth:
  type: "Public"

i18n:
  en:
    page_title: "Product Details"
    product_information: "Product Information"
    dynamic_route_demo: "Dynamic Route Demo"
  de:
    page_title: "Produktdetails"
    product_information: "Produktinformationen"
    dynamic_route_demo: "Dynamische Route Demo"

dynamic:
  parameters:
    locale:
      validation: "^(en|de)$"
      description: "Language locale code (en or de)"
      supported_values: ["en", "de"]
    id:
      validation: "^[a-zA-Z0-9_-]+$"
      description: "Product identifier"
```

## Dependency Injection

### Service Registration

```go
// main.go
func main() {
    container := di.NewContainer()
    
    // Register core services
    container.RegisterRouterServices("TR") // config prefix
    
    // Create your services
    templateRegistry, _ := templates.NewRegistry(container.GetInjector())
    assetsService, _ := assets.NewService(container.GetInjector())
    userStore, _ := services.NewDefaultUserStore(container.GetInjector())
    
    // Register application services
    container.RegisterApplicationServices(
        di.WithTemplateRegistry(templateRegistry),
        di.WithAssetsService(assetsService),
        di.WithUserStore(userStore),
    )
    
    // Register DataServices as named dependencies
    injector := container.GetInjector()
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
    
    router := container.GetRouter()
    router.Initialize()
}
```

### Data Service Integration

```go
// Data services are resolved by name based on template function signature
func ProductPage(productService ProductDataService, id string) templ.Component {
    // ProductDataService is automatically resolved from DI container
    product, _ := productService.GetProduct(id)
    return productTemplate(product)
}
```

## Development Tools

### Available Commands

```bash
# Development
make dev              # Start development server with hot reload
make templ            # Watch and compile templ files
make server           # Run server only

# Building
make build            # Build production binary
make docker           # Build Docker image

# Code Generation
make generate         # Generate template registry
make trgen            # Run template router generator

# Testing
make test             # Run all tests
make test-coverage    # Run tests with coverage
```

### Template Generator

Generate route registries automatically:

```bash
# Install the generator
go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Generate routes
trgen generate --input ./app --output ./generated/routes.go
```

## Production Deployment

### Docker

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/app ./app
COPY --from=builder /app/assets ./assets

EXPOSE 8080
CMD ["./main"]
```

### Environment Configuration

```bash
# Production environment
ENVIRONMENT_KIND=production
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SECURITY_CSRF_SECRET=your-production-secret
LOGGING_LEVEL=warn
LOGGING_FORMAT=json
```

## Advanced Examples

### Custom Data Service

```go
type ProductService struct {
    db *sql.DB
}

func (s *ProductService) GetProduct(id string) (*Product, error) {
    // Database logic
    return product, nil
}

func (s *ProductService) ListProducts() ([]*Product, error) {
    // Database logic
    return products, nil
}

// Register as named service in container
injector := container.GetInjector()
do.ProvideNamed(injector, "ProductDataService", func() *ProductService {
    return &ProductService{db: db}
})
```

### Custom Authentication

```go
type CustomAuthService struct {
    userRepo UserRepository
}

func (s *CustomAuthService) Authenticate(ctx context.Context, token string) (*interfaces.AuthResult, error) {
    // Custom authentication logic
    return &interfaces.AuthResult{
        IsAuthenticated: true,
        User: user,
    }, nil
}
```

### Advanced Routing

```go
// app/api/v1/[version]/[resource]/[id]/page.templ
func APIEndpoint(version, resource, id string) templ.Component {
    // Handle API versioning and resource routing
    return apiTemplate(version, resource, id)
}
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
git clone https://github.com/denkhaus/templ-router.git
cd templ-router
go mod tidy
make dev
```

### Running Tests

```bash
make test              # Run all tests
make test-coverage     # Run with coverage
make test-integration  # Run integration tests
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [templ](https://templ.guide/) - Amazing Go templating engine
- [chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- [Next.js](https://nextjs.org/) - Inspiration for file-based routing and layouts

## Links

- [Documentation](https://github.com/denkhaus/templ-router/wiki)
- [Examples](https://github.com/denkhaus/templ-router/tree/main/examples)
- [API Reference](https://pkg.go.dev/github.com/denkhaus/templ-router)
- [Issue Tracker](https://github.com/denkhaus/templ-router/issues)
- [Discussions](https://github.com/denkhaus/templ-router/discussions)

---

**Built with love for the Go and templ community**