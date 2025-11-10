# Templ Router

[![Go Version](https://img.shields.io/github/go-mod/go-version/denkhaus/templ-router)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/denkhaus/templ-router)](https://goreportcard.com/report/github.com/denkhaus/templ-router)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/github/actions/workflow/status/denkhaus/templ-router/ci.yml?branch=main)](https://github.com/denkhaus/templ-router/actions)

**Go library for file-based routing with [templ](https://templ.guide/) templates, dependency injection, and comprehensive middleware support.**

Templ Router provides automatic route generation from file structure, built-in authentication, internationalization, and data service integration.

> ⚠️ **Early Development**: API may change before v1.0.0. Pin to specific versions in production.

## ✨ Key Features

- **🗂️ File-Based Routing** - Automatic route generation from template files
- **🌍 Internationalization** - Multi-language support with locale-based routing
- **🔐 Authentication** - Built-in session-based authentication and authorization
- **🎨 Self-Contained Components** - Components with their own i18n and metadata
- **📊 Data Services** - Automatic data injection with dependency injection
- **⚡ Performance** - Optimized template caching and routing
- **🔧 Configuration** - Environment-based configuration
- **🛡️ Security** - CSRF protection, rate limiting, and security headers

## 🚀 Quick Start

### 1. Install Dependencies

```bash
go get github.com/denkhaus/templ-router
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/denkhaus/templ-router/cmd/trgen@latest
```

### 2. Create Project Structure

```bash
mkdir your-project && cd your-project
go mod init github.com/youruser/yourproject
mkdir -p app pkg/dataservices generated/templates
```

### 3. Create Templates

```go
// app/layout.templ
package main

templ Layout(title string, content templ.Component) {
    <!DOCTYPE html>
    <html>
    <head><title>{ title }</title></head>
    <body>{ content }</body>
    </html>
}

// app/page.templ
package main

templ Page() {
    <h1>Welcome to Templ Router!</h1>
}
```

### 4. Generate Template Registry

```bash
trgen --scan-path=app --module-name=github.com/youruser/yourproject
```

### 5. Create Main Application

```go
// main.go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/go-chi/chi/v5"
    "github.com/youruser/yourproject/generated/templates"
)

func main() {
    container := di.NewContainer()
    container.RegisterRouterServices("TR")  // Environment variable prefix

    templateRegistry, _ := templates.NewRegistry(container.GetInjector())
    container.RegisterApplicationServices(di.WithTemplateRegistry(templateRegistry))

    mux := chi.NewRouter()
    router := container.GetRouter()
    router.Initialize()
    router.RegisterRoutes(mux)

    http.ListenAndServe(":8080", mux)
}
```

### 6. Run Application

```bash
templ generate  # Generate templ files
go run .
```

Visit `http://localhost:8080` to see your application!

## 📚 Documentation

### Getting Started

- **[Getting Started Guide](docs/GETTING-STARTED.md)** - Complete setup tutorial
- **[Configuration Reference](docs/CONFIGURATION.md)** - All configuration options
- **[Architecture Overview](docs/ARCHITECTURE.md)** - System architecture and design

### Core Features

- **[File-Based Routing](docs/FILE-BASED-ROUTING.md)** - Automatic routing from file structure
- **[Template Generator](docs/TEMPLATE-GENERATOR.md)** - Using the trgen CLI tool
- **[Authentication](docs/AUTHENTICATION.md)** - User authentication and role-based access
- **[Internationalization](docs/INTERNATIONALIZATION.md)** - Multi-language support
- **[Data Services](docs/DATA-SERVICES.md)** - Data injection and service patterns
- **[Dependency Injection](docs/DEPENDENCY-INJECTION.md)** - DI container and service management

### Advanced Topics

- **[Middleware System](docs/MIDDLEWARE.md)** - Custom middleware and pipeline
- **[Template System](docs/TEMPLATES.md)** - Advanced template features
- **[Production Deployment](docs/PRODUCTION-DEPLOYMENT.md)** - Deployment strategies

### Reference

- **[API Reference](docs/REFERENCE.md)** - Quick API reference
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute

## 🏗️ Project Structure

```ini
your-project/
├── app/                        # Your templates
│   ├── layout.templ           # Root layout
│   ├── page.templ             # Home page (/)
│   ├── login/
│   │   └── page.templ         # Login page (/login)
│   └── locale_/               # Internationalized routes
│       ├── page.templ         # /en, /de, /fr
│       └── dashboard/
│           └── page.templ     # /en/dashboard, /de/dashboard
├── pkg/dataservices/          # Your data services
├── generated/templates/       # Generated template registry
└── main.go                    # Application entry point
```

## 🔧 Environment Configuration

Configure with environment variables (prefix configurable, default "TR"):

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

## 🔄 Development Workflow

```bash
# Start development environment (watches templates, CSS, registry, and runs server)
mage dev

# Generate templates from .templ files
mage build:templGenerate

# Generate template registry from file structure
mage build:registryGenerate

# Run tests
mage test:all
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.