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

#### Request Processing Pipeline

Each request flows through a processing pipeline that handles different concerns:

```ini
Request → Authentication → Internationalization → Template Rendering → Response
```

**What happens in each step:**

1. **Authentication**: Checks if user has access to the requested page
2. **Internationalization**: Detects user language and loads translations
3. **Template Rendering**: Renders the template with data and sends response

This pipeline ensures that templates always receive the correct authentication context and translations for the user's language.

### 🗂️ File-Based Routing

- Routes automatically generated from file structure using `trgen`
- Dynamic parameters e.g. : `id_/` (underscore suffix), `locale_/` for internationalization
- Route precedence system for conflict resolution
- Template-to-route mapping with configurable patterns

### 🌍 Internationalization (i18n)

- **Multi-language support** with `locale_/` directory structure
- **YAML-based translations** in `.templ.yaml` metadata files with nested structure support
- **Context-based translation system** (no global `t()` function)
- **Automatic locale detection** and validation from URLs
- **Flexible i18n formats**: flat, nested, and multi-locale configurations
- **Dot notation support** for deeply nested translation keys

### 🔐 Authentication & Authorization

- Three authentication types: `AuthTypePublic`, `AuthTypeUser`, `AuthTypeAdmin`
- Built-in authentication routes: sign in, sign out, sign up
- Session-based authentication with configurable expiry
- Role-based access control with user role validation
- Hierarchical auth configuration with precedence rules
- Configurable success redirect routes for positive authentication

#### Authentication Configuration

Configure authentication requirements using `.templ.yaml` files alongside your templates:

```yaml
# app/admin/page.templ.yaml
auth:
  type: "AdminRequired"
  redirect_url: "/login"

# app/dashboard/page.templ.yaml
auth:
  type: "UserRequired"
  redirect_url: "/login"

# app/public/page.templ.yaml (or omit auth section)
auth:
  type: "Public"
```

**Authentication Types:**

- `Public`: No authentication required (default)
- `UserRequired`: Any authenticated user can access
- `AdminRequired`: Only admin users can access

**Configuration Priority:**

1. **Template-level** (`.templ.yaml` files) - takes precedence
2. **Default** - public access when no auth specified

#### Built-in Authentication Routes

The router automatically provides authentication API endpoints:

```bash
POST /api/auth/signin      # User sign in
POST /api/auth/signout     # User sign out
POST /api/auth/signup      # User registration
```

These endpoints handle:

- **Session Management**: Automatic session creation and cleanup
- **Redirects**: Configurable success/failure redirects
- **Validation**: Input validation and error handling

#### Session Configuration

Configure session behavior through environment variables:

```bash
# Session settings
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_NAME=session_id

# Redirect routes after successful authentication
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/welcome
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/

# Internationalized redirects (locale parameter replaced automatically)
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/{locale}/dashboard
```

#### Role-Based Access Control

For admin-only pages, the system checks user roles:

```yaml
# app/admin/page.templ.yaml
auth:
  type: "AdminRequired"
  redirect_url: "/login"
  roles: ["admin", "super_admin"]  # Optional: specific roles required
```

**How it works:**

- `UserRequired`: Any authenticated user can access
- `AdminRequired`: Only users with admin privileges
- `roles`: Additional role restrictions (optional)

### 🎨 Layout & Template System

- Layout inheritance with automatic composition
- Error template system with precedence-based resolution
- Template middleware with data service injection
- Configurable template extensions and metadata

### 📊 Data Service Integration

- **Automatic Data Injection**: Templates can declare data service requirements
- **Two Method Patterns**: `GetData()` method or specific `Get<Data Struct Name>()` methods
- **Parameter Injection**: Route parameters automatically passed to data services
- **DI Registration**: Data services registered via `do.ProvideNamed()`

### ⚡ Performance & Validation

- **Cache Service**: Template and route caching for performance optimization
- **Validation Orchestrator**: Comprehensive parameter, route, and template validation
- **Error Handling**: Dedicated error template service with fallback mechanisms
- **File System Abstraction**: Library-agnostic file operations

## Library Usage

**Important:** This is a Go library, not a standalone application. You add it as a dependency to your own Go project.

### For Library Users (Your Project)

#### 1. Add Library Dependency

```bash
# In your Go project directory
go get github.com/denkhaus/templ-router

# Install required tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/denkhaus/templ-router/cmd/trgen@latest
```

#### 2. Create Your Project Structure

```bash
# Your project structure
your-project/
├── go.mod                      # Your module
├── main.go                     # Your application entry point
├── app/                        # Your templates directory
│   ├── layout.templ
│   ├── page.templ
│   └── locale_/
│       └── dashboard/
│           ├── page.templ
│           └── page.templ.yaml
├── pkg/
│   ├── dataservices/          # Your data services
│   └── config/                # Your configuration
└── generated/                 # Generated by trgen
    └── templates/
```

#### 3. Generate Template Registry

```bash
# Navigate to your application directory (where go.mod is located)
cd your-project
trgen --scan-path=app --module-name=github.com/youruser/yourproject
```

#### 4. Integrate Generated Templates

After running `trgen`, you'll have a generated template registry that needs to be integrated:

**What `trgen` generates:**

```bash
# After running trgen, you'll see:
your-project/
└── generated/
    └── templates/
        └── registry.go           # Template registry implementation

```

**Integration in your application:**

```go
// Your main.go
package main

import (
    "context"
    "net/http"

    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/denkhaus/templ-router/pkg/router/middleware"
    "github.com/go-chi/chi/v5"
    "github.com/samber/do/v2"

    // Import your generated template registry
    "github.com/youruser/yourproject/generated/templates"
    // Import your data services
    "github.com/youruser/yourproject/pkg/dataservices"
)

func main() {
    // Create DI container with router services
    container := di.NewContainer()
    defer container.Shutdown()

    // Register router services with config prefix
    container.RegisterRouterServices("TR")

    // Create your template registry
    templateRegistry, err := templates.NewRegistry(container.GetInjector())
    if err != nil {
        panic(err)
    }

    // Register application services using options pattern
    container.RegisterApplicationServices(
        di.WithTemplateRegistry(templateRegistry),
        // Add your other services here
    )

    // Register your data services as named dependencies
    injector := container.GetInjector()
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)

    // Create Chi router
    mux := chi.NewRouter()

    // Add auth context middleware
    authMiddleware, err := middleware.NewAuthContextMiddleware(injector)
    if err != nil {
        panic(err)
    }
    mux.Use(authMiddleware.Middleware)

    // Get clean router and initialize
    cleanRouter := container.GetRouter()
    if err := cleanRouter.Initialize(); err != nil {
        panic(err)
    }

    // Register file-based routes
    if err := cleanRouter.RegisterRoutes(mux); err != nil {
        panic(err)
    }

    // Start server
    http.ListenAndServe(":8080", mux)
}
```

**How it works:**

1. **trgen scans** your `app/` directory for `.templ` files
2. **Generates registry** with all discovered templates and routes
3. **Router uses registry** to map URLs to templates automatically
4. **No manual route registration** needed - everything is automatic

**Important:** Re-run `trgen` whenever you add, remove, or move template files.

### For Demo/Development

If you want to run the included demo or contribute to the library:

```bash
# Clone the repository (for demo/development only)
git clone https://github.com/denkhaus/templ-router.git
cd templ-router

# Run the demo
mage dev
```

**Key Difference:**

- **Library Users**: `go get` the library, implement interfaces, use `trgen` on your own templates
- **Demo/Development**: Clone the repo to run the example or contribute to the library

## Template Generator (trgen)

The `trgen` tool automatically generates template registries from your templ files. It must be run from your project directory.

### Installation Verification

```bash
# Verify installation
trgen --help
trgen --version
```

### Basic Usage

**Important:** `trgen` must be run from your application directory (where your `go.mod` is located).

```bash
# Navigate to your application directory first
cd your-project

# Generate template registry with required flags
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Or use environment variables
TEMPLATE_SCAN_PATH=app TEMPLATE_MODULE_NAME=github.com/youruser/yourproject trgen
```

**Why from application directory?**

- `trgen` needs access to your `go.mod` file
- Generated files are placed relative to your project structure
- Module name must match your `go.mod`

### Required Parameters

- `--scan-path`: Directory containing your `.templ` files (e.g., `app`, `templates`)
- `--module-name`: Your Go module name from `go.mod`

### Advanced Usage

```bash
# Watch mode for development
cd your-project
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch

# Custom watch extensions
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch --watch-extensions=".templ,.yaml"
```

### Environment Variables

All command-line flags have corresponding environment variables:

```bash
export TRGEN_SCAN_PATH=app
export TRGEN_MODULE_NAME=github.com/youruser/yourproject
export TRGEN_WATCH_MODE=true
export TRGEN_WATCH_EXTENSIONS=".templ,.yaml,.yml"

cd your-project
trgen
```

## Project Structure

```ini
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

```sh
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

Each template can have an optional `.templ.yaml` configuration file. Here are real examples from the demo:

## Internationalization (i18n) Configuration Examples

The i18n system supports multiple configuration formats to accommodate different project needs and organizational preferences.

### 1. Multi-Locale Nested Structure (Recommended)

Perfect for complex applications with organized translation hierarchies:

```yaml
# app/locale_/dashboard/page.templ.yaml
i18n:
  en:
    feedback:
      title: "Feedback Dashboard"
      subtitle: "Overview of customer feedback and analytics"
      export: "Export Data"
      refresh: "Refresh Data"
      reviews: "reviews"
      stats:
        total_reviews: "Total Reviews"
        average_rating: "Average Rating"
        productions: "Productions"
        cache_hit_rate: "Cache Hit Rate"
      productions:
        title: "Productions"
        subtitle: "Overview of all productions with review statistics"
      recent:
        title: "Recent Reviews"
        subtitle: "Latest customer feedback and comments"
      actions:
        create_new: "Create New"
        bulk_export: "Bulk Export"
        settings: "Settings"
  de:
    feedback:
      title: "Feedback Dashboard"
      subtitle: "Übersicht über Kundenfeedback und Analysen"
      export: "Daten exportieren"
      refresh: "Daten aktualisieren"
      reviews: "Bewertungen"
      stats:
        total_reviews: "Gesamtbewertungen"
        average_rating: "Durchschnittsbewertung"
        productions: "Produktionen"
        cache_hit_rate: "Cache-Trefferrate"
      productions:
        title: "Produktionen"
        subtitle: "Übersicht aller Produktionen mit Bewertungsstatistiken"
      recent:
        title: "Aktuelle Bewertungen"
        subtitle: "Neuestes Kundenfeedback und Kommentare"
      actions:
        create_new: "Neu erstellen"
        bulk_export: "Massenexport"
        settings: "Einstellungen"

auth:
  type: "UserRequired"
  redirect_url: "/login"
```

**Usage in templates with dot notation:**
```go
templ DashboardPage() {
    <h1>{ i18n.T(ctx, "feedback.title") }</h1>
    <p>{ i18n.T(ctx, "feedback.subtitle") }</p>
    <div class="stats">
        <span>{ i18n.T(ctx, "feedback.stats.total_reviews") }</span>
        <span>{ i18n.T(ctx, "feedback.stats.average_rating") }</span>
    </div>
    <button>{ i18n.T(ctx, "feedback.actions.create_new") }</button>
}
```

### 2. Multi-Locale Flat Structure

Simple key-value pairs for straightforward translations:

```yaml
# app/locale_/admin/page.templ.yaml
i18n:
  en:
    admin_warning: "Admin Area - Restricted Access"
    page_title: "System Administration"
    user_management_title: "User Management"
    user_management_desc: "Manage user accounts, roles, and permissions"
    system_settings_title: "System Settings"
    system_settings_desc: "Configure application settings and preferences"
    btn_save: "Save Changes"
    btn_cancel: "Cancel"

  de:
    admin_warning: "Admin-Bereich - Eingeschränkter Zugang"
    page_title: "Systemadministration"
    user_management_title: "Benutzerverwaltung"
    user_management_desc: "Benutzerkonten, Rollen und Berechtigungen verwalten"
    system_settings_title: "Systemeinstellungen"
    system_settings_desc: "Anwendungseinstellungen und Präferenzen konfigurieren"
    btn_save: "Änderungen speichern"
    btn_cancel: "Abbrechen"

auth:
  type: "AdminRequired"
  redirect_url: "/login"
```

### 3. Single-Locale Nested Structure

For applications that don't need multi-language support but want organized translations:

```yaml
# app/components/navigation/page.templ.yaml
i18n:
  navigation:
    main:
      home: "Home"
      about: "About Us"
      services: "Services"
      contact: "Contact"
    user:
      profile: "My Profile"
      settings: "Account Settings"
      logout: "Sign Out"
    admin:
      dashboard: "Admin Dashboard"
      users: "User Management"
      reports: "System Reports"
  buttons:
    primary:
      submit: "Submit"
      save: "Save"
      continue: "Continue"
    secondary:
      cancel: "Cancel"
      back: "Go Back"
      reset: "Reset Form"
```

**Usage with dot notation:**
```go
templ NavigationComponent() {
    <nav>
        <a href="/">{ i18n.T(ctx, "navigation.main.home") }</a>
        <a href="/about">{ i18n.T(ctx, "navigation.main.about") }</a>
        <a href="/services">{ i18n.T(ctx, "navigation.main.services") }</a>
    </nav>
    <div class="user-menu">
        <a href="/profile">{ i18n.T(ctx, "navigation.user.profile") }</a>
        <button>{ i18n.T(ctx, "navigation.user.logout") }</button>
    </div>
}
```

### 4. Single-Locale Flat Structure

Traditional flat key-value structure:

```yaml
# app/simple/page.templ.yaml
i18n:
  page_title: "Welcome to Our Application"
  page_subtitle: "Get started with our amazing features"
  btn_get_started: "Get Started"
  btn_learn_more: "Learn More"
  feature_1_title: "Fast Performance"
  feature_1_desc: "Lightning-fast response times"
  feature_2_title: "Secure by Design"
  feature_2_desc: "Enterprise-grade security"
```

### 5. Mixed Depth Nested Structure

Combining different nesting levels as needed:

```yaml
# app/locale_/ecommerce/page.templ.yaml
i18n:
  en:
    # Top-level keys
    site_name: "Amazing Store"
    welcome_message: "Welcome to our online store!"
    
    # Nested product information
    products:
      categories:
        electronics: "Electronics"
        clothing: "Clothing"
        books: "Books"
        home_garden: "Home & Garden"
      actions:
        add_to_cart: "Add to Cart"
        buy_now: "Buy Now"
        view_details: "View Details"
        compare: "Compare Products"
      filters:
        price_range: "Price Range"
        brand: "Brand"
        rating: "Customer Rating"
        availability: "Availability"
    
    # Deeply nested checkout process
    checkout:
      steps:
        cart: "Shopping Cart"
        shipping: "Shipping Information"
        payment: "Payment Details"
        confirmation: "Order Confirmation"
      shipping:
        methods:
          standard: "Standard Delivery (5-7 days)"
          express: "Express Delivery (2-3 days)"
          overnight: "Overnight Delivery"
        address:
          street: "Street Address"
          city: "City"
          postal_code: "Postal Code"
          country: "Country"
      payment:
        methods:
          credit_card: "Credit Card"
          paypal: "PayPal"
          bank_transfer: "Bank Transfer"
        security:
          ssl_notice: "Your payment information is secure and encrypted"
          privacy_notice: "We never store your payment details"

  de:
    site_name: "Fantastischer Shop"
    welcome_message: "Willkommen in unserem Online-Shop!"
    
    products:
      categories:
        electronics: "Elektronik"
        clothing: "Kleidung"
        books: "Bücher"
        home_garden: "Haus & Garten"
      actions:
        add_to_cart: "In den Warenkorb"
        buy_now: "Jetzt kaufen"
        view_details: "Details anzeigen"
        compare: "Produkte vergleichen"
      filters:
        price_range: "Preisspanne"
        brand: "Marke"
        rating: "Kundenbewertung"
        availability: "Verfügbarkeit"
    
    checkout:
      steps:
        cart: "Warenkorb"
        shipping: "Versandinformationen"
        payment: "Zahlungsdetails"
        confirmation: "Bestellbestätigung"
      shipping:
        methods:
          standard: "Standardversand (5-7 Tage)"
          express: "Expressversand (2-3 Tage)"
          overnight: "Über-Nacht-Versand"
        address:
          street: "Straße"
          city: "Stadt"
          postal_code: "Postleitzahl"
          country: "Land"
      payment:
        methods:
          credit_card: "Kreditkarte"
          paypal: "PayPal"
          bank_transfer: "Banküberweisung"
        security:
          ssl_notice: "Ihre Zahlungsinformationen sind sicher und verschlüsselt"
          privacy_notice: "Wir speichern niemals Ihre Zahlungsdetails"
```

**Usage in complex templates:**
```go
templ CheckoutPage() {
    <h1>{ i18n.T(ctx, "site_name") }</h1>
    <div class="checkout-steps">
        <span>{ i18n.T(ctx, "checkout.steps.cart") }</span>
        <span>{ i18n.T(ctx, "checkout.steps.shipping") }</span>
        <span>{ i18n.T(ctx, "checkout.steps.payment") }</span>
    </div>
    
    <div class="shipping-methods">
        <label>
            <input type="radio" name="shipping" value="standard"/>
            { i18n.T(ctx, "checkout.shipping.methods.standard") }
        </label>
        <label>
            <input type="radio" name="shipping" value="express"/>
            { i18n.T(ctx, "checkout.shipping.methods.express") }
        </label>
    </div>
    
    <div class="security-notice">
        <p>{ i18n.T(ctx, "checkout.payment.security.ssl_notice") }</p>
    </div>
}
```

### Key Benefits of Nested i18n Structure

1. **Organization**: Group related translations logically
2. **Maintainability**: Easier to find and update related translations
3. **Scalability**: Handle complex applications with hundreds of translation keys
4. **Readability**: Clear hierarchy makes translation files self-documenting
5. **Flexibility**: Mix flat and nested structures as needed
6. **Dot Notation**: Access nested keys with simple `"parent.child.key"` syntax

### Best Practices

- **Use nested structures** for complex applications with many translation keys
- **Group related translations** under common parent keys (e.g., `buttons`, `forms`, `navigation`)
- **Keep nesting levels reasonable** (2-4 levels deep maximum)
- **Use consistent naming conventions** across your translation files
- **Organize by feature or component** rather than by page when possible

### Login Page (Public Access)

```yaml
# app/login/page.templ.yaml
i18n:
  en:
    login_title: "Sign in to your account"
    login_subtitle: "Enter your credentials to access the system"
    username: "Username"
    password: "Password"
    sign_in: "Sign in"
    need_account: "Don't have an account? Sign up"

  de:
    login_title: "Bei Ihrem Konto anmelden"
    login_subtitle: "Geben Sie Ihre Anmeldedaten ein, um auf das System zuzugreifen"
    username: "Benutzername"
    password: "Passwort"
    sign_in: "Anmelden"
    need_account: "Noch kein Konto? Registrieren"

auth:
  type: "Public"
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
// Auth API routes are already build in
// Your templates can use these API endpoints:
// <form hx-post="/api/auth/signin">
// <form hx-post="/api/auth/signup">
// <form method="POST" action="/api/auth/signout">
```

## Internationalization

Translation files support both flat and nested structures with locale-specific keys:

```yaml
# demo/app/locale_/dashboard/page.templ.yaml
i18n:
  en:
    page_title: "Dashboard"
    page_subtitle: "Overview of your application metrics"
    navigation:
      home: "Home"
      settings: "Settings"
      logout: "Sign Out"
    stats:
      users: "Total Users"
      projects: "Active Projects"
      revenue: "Monthly Revenue"
  de:
    page_title: "Dashboard"
    page_subtitle: "Übersicht Ihrer Anwendungsmetriken"
    navigation:
      home: "Startseite"
      settings: "Einstellungen"
      logout: "Abmelden"
    stats:
      users: "Gesamte Benutzer"
      projects: "Aktive Projekte"
      revenue: "Monatlicher Umsatz"
```

**Access nested keys with dot notation:**
```go
templ DashboardPage() {
    <h1>{ i18n.T(ctx, "page_title") }</h1>
    <p>{ i18n.T(ctx, "page_subtitle") }</p>
    <nav>
        <a href="/">{ i18n.T(ctx, "navigation.home") }</a>
        <a href="/settings">{ i18n.T(ctx, "navigation.settings") }</a>
    </nav>
    <div class="stats">
        <span>{ i18n.T(ctx, "stats.users") }: 1,234</span>
        <span>{ i18n.T(ctx, "stats.projects") }: 56</span>
    </div>
}
```

### I18n Helper Functions

The `i18n` package provides several context-based helper functions for templates:

#### Core Translation Function

```go
import "github.com/denkhaus/templ-router/pkg/router/i18n"

// Primary translation function
i18n.T(ctx, "translation_key")
// Returns the translated string for the current locale
// Falls back to "[MISSING_I18N: key]" if translation not found

templ DashboardPage() {
    <h1>{ i18n.T(ctx, "page_title") }</h1>
    <p>{ i18n.T(ctx, "page_subtitle") }</p>
}
```

#### URL Localization

```go
// Automatically adds locale prefix to URLs
i18n.LocalizeSafeURL(ctx, "/dashboard")
// Returns: templ.SafeURL("/en/dashboard") or templ.SafeURL("/de/dashboard")

templ Navigation() {
    <nav>
        <a href={ i18n.LocalizeSafeURL(ctx, "/admin") }>
            { i18n.T(ctx, "nav_admin") }
        </a>
        <a href={ i18n.LocalizeSafeURL(ctx, "/user/profile") }>
            { i18n.T(ctx, "nav_profile") }
        </a>
    </nav>
}
```

#### Context Information Functions

```go
// Get current locale
i18n.GetCurrentLocale(ctx)
// Returns: "en", "de", "fr", etc.

// Get current template path
i18n.GetCurrentTemplate(ctx)
// Returns: "app/locale_/dashboard/page.templ"

// Get all available translation keys for current template
i18n.GetAvailableKeys(ctx)
// Returns: []string{"page_title", "stats_users", "nav_admin", ...}
```

#### Practical Examples

```go
// Language switcher with current locale detection
templ LanguageSwitcher() {
    <div class="language-switcher">
        <span>Current: { i18n.GetCurrentLocale(ctx) }</span>
        if i18n.GetCurrentLocale(ctx) == "en" {
            <a href="/de/dashboard">Switch to Deutsch</a>
        } else {
            <a href="/en/dashboard">Switch to English</a>
        }
    </div>
}

// Conditional content based on locale
templ PriceDisplay(amount float64) {
    <span class="price">
        if i18n.GetCurrentLocale(ctx) == "de" {
            { fmt.Sprintf("%.2f €", amount) }
        } else {
            { fmt.Sprintf("$%.2f", amount) }
        }
    </span>
}

// Debug information panel
templ DebugPanel() {
    <div class="debug-panel">
        <p><strong>Locale:</strong> { i18n.GetCurrentLocale(ctx) }</p>
        <p><strong>Template:</strong> { i18n.GetCurrentTemplate(ctx) }</p>
        <p><strong>Available Keys:</strong> { fmt.Sprint(len(i18n.GetAvailableKeys(ctx))) }</p>
    </div>
}
```

## Development Workflow

### Your Project Development

When developing your own application using templ-router:

```bash
# 1. Generate templates (after changing .templ files)
cd your-project
templ generate

# 2. Generate template registry (after adding/removing templates)
# Must be run from your application directory
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# 3. Build and run your application
go build
./your-app

# 4. Development with hot reload (optional)
# Install air: go install github.com/air-verse/air@latest
air
```

### Contributing to templ-router Library

If you want to contribute to the templ-router library itself:

```bash
# Clone the library repository
git clone https://github.com/denkhaus/templ-router.git
cd templ-router

# Library development commands
mage dev                    # Start demo server
mage test:all               # Run library tests
mage build:all              # Build library for all platforms
mage generator:install      # Install trgen from source
```

## Configuration

Environment variables use configurable prefix (set in `RegisterRouterServices(prefix)`):

```ini
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

# Router Configuration
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED=true
```

# Security Configuration

```ini
TR_SECURITY_CSRF_SECRET=change-me-in-production
TR_SECURITY_CSRF_SECURE=false
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_SECURITY_RATE_LIMIT_REQUESTS=100
TR_SECURITY_ENABLE_SECURITY_HEADERS=true
```

# Logging Configuration

```ini
TR_LOGGING_LEVEL=info
TR_LOGGING_FORMAT=json
TR_LOGGING_OUTPUT=stdout
TR_LOGGING_ENABLE_FILE=false
TR_LOGGING_FILE_PATH=logs/router.log
```

# Environment Configuration

```sh
TR_ENVIRONMENT_KIND=develop
```

**Configuration Sections:** Server, Database, Auth, Email, Security, Logging, I18n, Layout, TemplateGenerator, Router, Environment

### Router Configuration

The router provides several built-in middleware features that can be configured:

| Setting | Default | Description |
|---------|---------|-------------|
| `TR_ROUTER_ENABLE_TRAILING_SLASH` | `true` | Automatically redirects `/path/` to `/path` and vice versa |
| `TR_ROUTER_ENABLE_SLASH_REDIRECT` | `true` | Cleans up double slashes in URLs (e.g., `/path//` → `/path/`) |
| `TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED` | `true` | Enables 405 Method Not Allowed handler for unsupported HTTP methods |

**Examples:**

```bash
# Enable trailing slash redirection (default: true)
TR_ROUTER_ENABLE_TRAILING_SLASH=true
# /dashboard/ → redirects to /dashboard
# /dashboard → redirects to /dashboard/

# Enable slash cleanup (default: true)
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
# /path//to///resource → redirects to /path/to/resource

# Enable method not allowed handler (default: true)
TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED=true
# POST /get-only-route → returns 405 Method Not Allowed
```

**Benefits:**

- **SEO Friendly**: Prevents duplicate content issues from trailing slash variations
- **Clean URLs**: Automatically fixes malformed URLs with double slashes
- **Better UX**: Proper HTTP status codes for unsupported methods
- **Out-of-the-Box**: No manual middleware configuration required

## Data Services

Data services provide data to templates through dependency injection. The router automatically resolves and calls the appropriate data service based on template requirements.

### 🔍 RouterContext - Unified Parameter Access

**NEW**: Data services now use `RouterContext` for unified access to URL parameters, query parameters, and request data:

```go
import "github.com/denkhaus/templ-router/pkg/interfaces"

type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
}

func (s *userDataServiceImpl) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
    // URL Parameters (from route like /{locale}/user/{id})
    locale := routerCtx.GetURLParam("locale")
    userID := routerCtx.GetURLParam("id")

    // Query Parameters (from URL like ?page=5&pageSize=10&filter=active)
    page := routerCtx.GetQueryParam("page")
    pageSize := routerCtx.GetQueryParam("pageSize")
    filter := routerCtx.GetQueryParam("filter")

    // Set defaults for query parameters
    if page == "" {
        page = "1"
    }
    if pageSize == "" {
        pageSize = "10"
    }

    return &UserData{
        ID:       userID,
        Name:     "User " + userID,
        Locale:   locale,
        Page:     page,
        PageSize: pageSize,
        Filter:   filter,
    }, nil
}
```

#### RouterContext Methods

```go
// URL Parameter access (from Chi router path parameters)
routerCtx.GetURLParam("key")           // Single parameter
routerCtx.GetAllURLParams()            // All URL parameters

// Query Parameter access (from URL query string)
routerCtx.GetQueryParam("key")         // First value
routerCtx.GetQueryParams("key")        // All values for key
routerCtx.GetAllQueryParams()          // All query parameters

// Advanced access
routerCtx.Context()                    // context.Context
routerCtx.Request()                    // *http.Request
routerCtx.ChiContext()                 // *chi.Context
```

### Data Service Patterns

**Pattern 1: GetData() Only**

```go
type ProductDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*ProductData, error)
}

func (s *productDataServiceImpl) GetData(routerCtx interfaces.RouterContext) (*ProductData, error) {
    productID := routerCtx.GetURLParam("id") // Route parameters
    category := routerCtx.GetQueryParam("category") // Query parameters

    return &ProductData{
        ID:       productID,
        Name:     "Product " + productID,
        Category: category,
    }, nil
}
```

**Pattern 2: GetData() + Specific Methods**

```go
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
}

// Both methods available - router chooses appropriate one
func (s *userDataServiceImpl) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")
    locale := routerCtx.GetURLParam("locale") // Multi-language support

    return &UserData{ID: userID, Name: "User " + userID, Locale: locale}, nil
}
```

Note: The Method name to access the `UserData` struct must be named `Get<Entity Name>` like `GetUserData`
Both patterns work in parallel where `GetData` is the fallback method.

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

### Query Parameter Demo

Try the live query parameter demo to see RouterContext in action:

```bash
# Start the demo application
cd demo && go run main.go

# Visit these URLs to test query parameters:
# http://localhost:8080/en/query-demo?page=1&pageSize=10&sort=name
# http://localhost:8080/de/query-demo?page=2&pageSize=20&filter=premium&sort=date
```

The demo shows how RouterContext cleanly separates URL parameters from query parameters, preventing conflicts and providing type-safe access.

### Parameter Access Examples

RouterContext provides clean access to all request parameters:

```go
func (s *userDataServiceImpl) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")       // From route: /user/123 → id="123"
    locale := routerCtx.GetURLParam("locale")   // From route: /en/user/123 → locale="en"
    page := routerCtx.GetQueryParam("page")     // From query: ?page=2 → page="2"

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
