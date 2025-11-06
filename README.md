# Templ Router

[![Go Version](https://img.shields.io/github/go-mod/go-version/denkhaus/templ-router)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/denkhaus/templ-router)](https://goreportcard.com/report/github.com/denkhaus/templ-router)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/github/actions/workflow/status/denkhaus/templ-router/ci.yml?branch=main)](https://github.com/denkhaus/templ-router/actions)

**A production-ready Go library for file-based routing with [templ](https://templ.guide/) templates, dependency injection, and comprehensive middleware support.**

Templ Router provides automatic route generation from file structure, built-in authentication, internationalization, data service integration, and a clean architecture that follows Go best practices.

## ⚠️ Early Development Warning

This project is currently in early development. The API may change before v1.0.0. We recommend pinning to specific versions in production.

## ✨ Key Features

- **🗂️ File-Based Routing** - Automatic route generation from your template files
- **🌍 Internationalization** - Multi-language support with locale-based routing
- **🔐 Authentication** - Built-in session-based authentication and authorization
- **🎨 Template System** - Layout inheritance and **self-contained components** with their own i18n
- **📊 Data Services** - Automatic data injection with dependency injection
- **⚡ Performance** - Optimized template caching and routing
- **🔧 Configuration** - Comprehensive environment-based configuration
- **🛡️ Security** - CSRF protection, rate limiting, and security headers

## ✨ Unique Features

### Self-Contained Components

Components can have their own metadata and internationalization files, making them **truly reusable**:

```yaml
# app/components/footer.templ.yaml
i18n:
  en: { footer_copyright: "© 2024 My Company" }
  de: { footer_copyright: "© 2024 Meine Firma" }

metadata:
  company_name: "My Company"
  company_email: "info@company.com"
```

**Benefits:**
- **No duplication** - Define component translations once, use everywhere
- **True reusability** - Components work independently across pages
- **Easy maintenance** - Update component config in one place
- **Consistent behavior** - Same component works identically everywhere

## 🚀 Quick Start

### 1. Install Dependencies

```bash
# Install templ-router
go get github.com/denkhaus/templ-router

# Install required tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/denkhaus/templ-router/cmd/trgen@latest
```

### 2. Create Your Project

```bash
mkdir your-project && cd your-project
go mod init github.com/youruser/yourproject

# Create basic structure
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
    container.RegisterRouterServices("TR")  // "TR" is the default prefix for env vars (configurable)

    templateRegistry, _ := templates.NewRegistry(container.GetInjector())
    container.RegisterApplicationServices(di.WithTemplateRegistry(templateRegistry))

    mux := chi.NewRouter()
    router := container.GetRouter()
    router.Initialize()
    router.RegisterRoutes(mux)

    http.ListenAndServe(":8080", mux)
}
```

### 6. Run Your Application

```bash
templ generate  # Generate templ files
go run .
```

Visit `http://localhost:8080` to see your application!

## 📚 Documentation

### Getting Started

- **[Getting Started Guide](docs/GETTING-STARTED.md)** - Complete setup and integration tutorial
- **[Configuration Reference](docs/CONFIGURATION.md)** - All configuration options
- **[Architecture Overview](docs/ARCHITECTURE.md)** - System architecture and design

### Core Features

- **[File-Based Routing](docs/FILE-BASED-ROUTING.md)** - Automatic routing from file structure
- **[Template Generator](docs/TEMPLATE-GENERATOR.md)** - Using the trgen CLI tool
- **[Authentication & Authorization](docs/AUTHENTICATION.md)** - User authentication and role-based access
- **[Internationalization](docs/INTERNATIONALIZATION.md)** - Multi-language support
- **[Data Services](docs/DATA-SERVICES.md)** - Data injection and service patterns
- **[Dependency Injection](docs/DEPENDENCY-INJECTION.md)** - DI container and service management

### Advanced Topics

- **[Middleware System](docs/MIDDLEWARE.md)** - Custom middleware and pipeline
- **[Template System](docs/TEMPLATES.md)** - Advanced template features
- **[Production Deployment](docs/PRODUCTION-DEPLOYMENT.md)** - Deployment strategies and best practices

### Reference

- **[API Reference](docs/REFERENCE.md)** - Quick API reference and cheat sheet
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute to the project

## 🏗️ Project Structure

```
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
├── pkg/
│   └── dataservices/          # Your data services
├── generated/
│   └── templates/
│       └── registry.go        # Generated template registry
├── main.go                    # Your application entry point
└── go.mod                     # Your module dependencies
```

## 🎯 Use Cases

### Perfect For

- **Content Management Systems** - Blog, documentation sites, marketing pages
- **SaaS Applications** - Multi-tenant applications with authentication
- **E-commerce Platforms** - Product catalogs with internationalization
- **Corporate Websites** - Multi-language corporate websites
- **Admin Panels** - Administrative interfaces with role-based access
- **API Backends** - RESTful APIs with templated documentation

### Key Benefits

- **🚀 Rapid Development** - Get started in minutes with automatic routing
- **🌍 Global Ready** - Built-in internationalization support
- **🔒 Secure by Default** - Authentication, CSRF protection, and security headers
- **📈 Scalable** - Clean architecture with dependency injection
- **🎨 Developer Experience** - Hot reload, comprehensive tooling
- **🛡️ Production Ready** - Battle-tested with comprehensive error handling

## 🛠️ Development Commands

```bash
# Development workflow
mage dev                    # Start development server with hot reload

# Template generation
templ generate              # Generate templ files
trgen --scan-path=app --module-name=github.com/youruser/yourproject
trgen --watch              # Watch mode for automatic regeneration

# Testing
mage test:all              # Run all tests
mage test:e2e              # Run E2E tests
mage test:devSetup         # Setup test environment

# Building
mage build:all             # Build for all platforms
mage generator:install     # Install trgen tool
```

## 🔧 Installation

### As a Library

```bash
# Add to your project
go get github.com/denkhaus/templ-router
```

### CLI Tools

```bash
# Install templ (template engine)
go install github.com/a-h/templ/cmd/templ@latest

# Install trgen (template registry generator)
go install github.com/denkhaus/templ-router/cmd/trgen@latest
```

### Using the Demo

```bash
# Clone and run the demo
git clone https://github.com/denkhaus/templ-router.git
cd templ-router/demo
go run main.go
```

## 📊 Features Overview

### File-Based Routing System

Routes are automatically generated from your file structure:

```
app/page.templ                    → /
app/login/page.templ              → /login
app/locale_/dashboard/page.templ  → /en/dashboard, /de/dashboard
app/user/id_/page.templ           → /user/123, /user/456
```

### Authentication & Authorization

Built-in authentication with three levels:

- `Public` - No authentication required
- `UserRequired` - Any authenticated user
- `AdminRequired` - Admin users only

### Internationalization

Comprehensive i18n support with:

- Locale-based routing (`/en/dashboard`, `/de/dashboard`)
- YAML-based translations
- Context-based translation functions
- Automatic locale detection

### Data Services

Clean data injection pattern:

```go
type UserDataService interface {
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
}

// Automatic injection into templates (one per template)
templ Page(user *UserData) {
    <h1>{ user.Name }</h1>
}
```

### Configuration

Comprehensive environment-based configuration with **configurable prefix**:

```bash
# Default prefix (TR) - configurable in application setup
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080
TR_AUTH_SESSION_EXPIRY=24h
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_LOGGING_LEVEL=info

# Custom prefix examples:
MYAPP_SERVER_HOST=localhost      # Using "MYAPP" prefix
BLOG_SERVER_HOST=localhost       # Using "BLOG" prefix
API_SERVER_HOST=localhost        # Using "API" prefix
```

**Note:** The prefix `TR_` is the default but fully configurable via `container.RegisterRouterServices("YOUR_PREFIX")`.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

### Development Setup

```bash
git clone https://github.com/denkhaus/templ-router.git
cd templ-router
go mod tidy
mage dev                    # Start development server
mage test:all               # Run tests
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **[Documentation](docs/)** - Complete documentation
- **[API Reference](docs/REFERENCE.md)** - Quick API reference
- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[pkg.go.dev](https://pkg.go.dev/github.com/denkhaus/templ-router)** - Go package documentation

---

**Built with ❤️ for the Go and templ community**