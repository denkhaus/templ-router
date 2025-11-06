# File-Based Routing System

**Complete guide to Templ Router's automatic file-based routing system.**

## Overview

Templ Router automatically generates HTTP routes based on your file structure, similar to Next.js app router. The `trgen` CLI tool scans your template files and creates a registry that maps URLs to templates automatically.

**Core Features:**
- Automatic route generation from file structure
- Dynamic parameters with `_` suffix (e.g., `id_/`, `locale_/`)
- Route precedence system for conflict resolution
- Internationalized routes via `locale_/` directory structure
- Template-to-route mapping with configurable patterns

## Configuration Prefix Notice

**Important:** Some environment variables in this documentation use the default prefix `TR_`. This prefix is **configurable** when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_I18N_SUPPORTED_LOCALES=en,de,fr`
- Custom: `MYAPP_I18N_SUPPORTED_LOCALES=en,de,fr`
- Multiple apps: `APP1_I18N_SUPPORTED_LOCALES=en,de` and `APP2_I18N_SUPPORTED_LOCALES=fr,es`

All environment variable examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix.

## Directory Structure

Routes are generated based on the folder structure in your template directory (usually `app/`):

```
app/
├── layout.templ              # Root layout (fallback)
├── page.templ                # Root page (/)
├── login/
│   ├── page.templ            # /login
│   └── page.templ.yaml       # Route metadata
├── dashboard/
│   ├── layout.templ          # Dashboard layout
│   ├── page.templ            # /dashboard
│   └── settings/
│       └── page.templ        # /dashboard/settings
├── locale_/                  # Internationalized routes
│   ├── page.templ            # /en, /de (based on config)
│   ├── dashboard/
│   │   └── page.templ        # /en/dashboard, /de/dashboard
│   └── user/
│       └── id_/              # Dynamic parameter
│           └── page.templ    # /en/user/123, /de/user/123
└── components/
    ├── navbar.templ          # Reusable component
    ├── footer.templ          # Reusable component
    └── navbar.templ.yaml     # Component metadata
```

## Route Generation Patterns

### Static Routes

Static routes are generated directly from directory structure:

- `app/page.templ` → `/`
- `app/login/page.templ` → `/login`
- `app/dashboard/page.templ` → `/dashboard`
- `app/dashboard/settings/page.templ` → `/dashboard/settings`

### Dynamic Routes

Dynamic parameters use the `_` suffix convention:

- `app/user/id_/page.templ` → `/user/123`, `/user/456`
- `app/product/slug_/page.templ` → `/product/my-product`
- `app/locale_/page.templ` → `/en`, `/de` (special case for localization)

### Internationalized Routes

The `locale_/` directory automatically creates localized routes:

- `app/locale_/page.templ` → `/en`, `/de`, `/fr` (based on supported locales)
- `app/locale_/dashboard/page.templ` → `/en/dashboard`, `/de/dashboard`, `/fr/dashboard`
- `app/locale_/user/id_/page.templ` → `/en/user/123`, `/de/user/123`

### Complex Route Examples

```
app/locale_/admin/user/id_/profile/page.templ
```

Generates routes like:
- `/en/admin/user/123/profile`
- `/de/admin/user/123/profile`
- `/fr/admin/user/123/profile`

## Route Precedence

Routes are resolved with the following precedence (highest to lowest):

1. **Specific routes**: `/admin/user/profile` (static)
2. **Dynamic routes**: `/admin/user/{id}` (one parameter)
3. **Localized routes**: `/{locale}/admin/user/{id}` (with locale parameter)
4. **Fallback routes**: `/` (root route)

### Example Precedence

Given these templates:
- `app/admin/page.templ` → `/admin`
- `app/admin/id_/page.templ` → `/admin/{id}`
- `app/locale_/admin/page.templ` → `/{locale}/admin`

URL resolution:
- `/admin` → `app/admin/page.templ`
- `/admin/123` → `app/admin/id_/page.templ`
- `/en/admin` → `app/locale_/admin/page.templ`

## Template Registry Generation

The `trgen` tool generates route mappings by scanning your template files:

```bash
# Navigate to your project directory
cd your-project

# Generate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Watch mode for development
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch
```

### Generated Registry Structure

The generated `registry.go` contains:

```go
// Template registry with route mappings
type Registry struct {
    templates map[string]TemplateFunc
    routes    map[string]RouteInfo
}

// Route information for each template
type RouteInfo struct {
    Path       string
    Template   string
    Parameters []string
    IsLocalized bool
    RequiresAuth bool
}
```

## Layout Inheritance

Layouts are inherited hierarchically through the directory structure:

### Layout Resolution Order

1. **Current directory layout**: `app/dashboard/layout.templ`
2. **Parent directory layout**: `app/layout.templ`
3. **Root fallback layout**: `app/layout.templ`

### Layout Inheritance Examples

```
app/
├── layout.templ              # Root layout
├── page.templ                # Uses root layout
├── dashboard/
│   ├── layout.templ          # Dashboard-specific layout
│   └── page.templ            # Uses dashboard layout
└── user/
    └── id_/
        └── page.templ        # Uses root layout
```

## Component Routes

Components in the `components/` directory are automatically accessible via their own routes:

```
app/components/
├── navbar.templ              # /components/navbar
├── footer.templ              # /components/footer
└── navbar.templ.yaml         # Component metadata
```

### Component Use Cases

- **HTMX partials**: Load components without layout for AJAX requests
- **Reusable components**: Self-contained metadata and i18n for true reusability
- **API endpoints**: Direct component access for dynamic loading

### Self-Contained Components

Components can be **self-contained** with their own `.templ.yaml` files. This is a critical feature for building maintainable applications:

#### Problem Solved
**Without self-contained components**, you would duplicate configuration:
```yaml
# Repeating in EVERY page that uses navbar 😞
# app/home/page.templ.yaml
i18n: { nav_home: "Home", nav_dashboard: "Dashboard" }

# app/profile/page.templ.yaml
i18n: { nav_home: "Home", nav_dashboard: "Dashboard" }

# app/settings/page.templ.yaml
i18n: { nav_home: "Home", nav_dashboard: "Dashboard" }
```

**With self-contained components**, define once, reuse everywhere:
```yaml
# app/components/navbar.templ.yaml - Single source of truth ✅
i18n:
  en: { nav_home: "Home", nav_dashboard: "Dashboard" }
  de: { nav_home: "Startseite", nav_dashboard: "Dashboard" }

# Pages using navbar - NO configuration needed!
```

#### Benefits
1. **No duplication** - Component config defined once
2. **Easy maintenance** - Update in one place, affects all pages
3. **True reusability** - Components work independently
4. **Consistency** - Same behavior across all pages
5. **Clean pages** - Page configs focus on page-specific content

## Template Metadata System

Each template can have an optional `.templ.yaml` metadata file for configuration:

### Metadata File Structure

```yaml
# app/dashboard/page.templ.yaml
metadata:
  page_title: "Dashboard"
  description: "Main dashboard page"
  theme: "dark"

i18n:
  en:
    page_title: "Dashboard"
    welcome_message: "Welcome to your dashboard"
  de:
    page_title: "Dashboard"
    welcome_message: "Willkommen in Ihrem Dashboard"

auth:
  type: "UserRequired"
  redirect_url: "/login"
  roles: ["user", "admin"]  # Optional: specific roles required

data_services:
  - "DashboardDataService"
  - "UserStatsDataService"
```

### Accessing Metadata in Templates

```go
// app/dashboard/page.templ
package main

import (
    "github.com/denkhaus/templ-router/pkg/router/metadata"
    "github.com/denkhaus/templ-router/pkg/router/i18n"
)

templ DashboardPage() {
    pageTitle := metadata.M(ctx, "page_title")
    theme := metadata.M(ctx, "theme")

    <div class="dashboard" data-theme={ theme }>
        <h1>{ i18n.T(ctx, "welcome_message") }</h1>
        <p>Current theme: { theme }</p>
    </div>
}
```

## Route Configuration

### Custom Route Paths

Override auto-generated routes with custom paths:

```yaml
# app/admin/panel/page.templ.yaml
route:
  path: "/admin-control-panel"
  method: "GET"

# Results in route: /admin-control-panel instead of /admin/panel
```

### HTTP Method Restrictions

Restrict routes to specific HTTP methods:

```yaml
# app/api/users/page.templ.yaml
route:
  methods: ["GET", "POST"]

# Route only responds to GET and POST requests
```

## Special Directories

### `locale_/` Directory

Special directory for internationalized content:

- Automatically prepends locale codes to routes
- Uses configured supported locales
- Falls back to default locale for unsupported locales

### `components/` Directory

Special directory for reusable components:

- Components accessible via `/components/*` routes
- Can have their own metadata and i18n
- Render without layout when accessed directly

### Error Templates

Error templates follow the hierarchy:

- `app/error.templ` → Global error template
- `app/admin/error.templ` → Admin section error template
- Falls back to nearest parent error template

## Route Examples

### Basic Static Routes

```
app/
├── page.templ          → /
├── about.templ         → /about
├── contact.templ       → /contact
└── help/
    └── page.templ      → /help
```

### Dynamic Parameters

```
app/
├── user/
│   └── id_/
│       └── page.templ  → /user/123
├── product/
│   └── slug_/
│       └── page.templ  → /product/my-product
└── search/
    └── query_/
        └── page.templ  → /search/javascript
```

### Internationalized Routes

```
app/locale_/dashboard/page.templ
```

With supported locales `en,de,fr`:
- `/en/dashboard`
- `/de/dashboard`
- `/fr/dashboard`

### Complex Combinations

```
app/locale_/user/id_/settings/section_/page.templ
```

Results in:
- `/en/user/123/settings/profile`
- `/de/user/123/settings/profile`
- `/fr/user/123/settings/profile`

## Integration with HTTP Router

### Chi Router Integration

```go
// main.go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/go-chi/chi/v5"
)

func main() {
    // Create DI container and router
    container := di.NewContainer()
    container.RegisterRouterServices("TR")

    router := container.GetRouter()
    router.Initialize()

    // Create Chi router
    mux := chi.NewRouter()

    // Register file-based routes
    if err := router.RegisterRoutes(mux); err != nil {
        panic(err)
    }

    // Manual routes (take precedence)
    mux.Get("/custom", customHandler)

    http.ListenAndServe(":8080", mux)
}
```

## Best Practices

### Directory Organization

1. **Group related pages**: Use directories to group related functionality
2. **Consistent naming**: Use lowercase, hyphenated directory names
3. **Logical hierarchy**: Place nested routes in appropriate subdirectories

### Parameter Naming

1. **Use `_` suffix**: Always use `_` for dynamic parameters
2. **Descriptive names**: Use clear parameter names (`id_`, `slug_`, `uuid_`)
3. **Consistent conventions**: Use the same parameter names across similar routes

### Internationalization

1. **Use `locale_/` directory**: For all internationalized content
2. **Structure translations**: Organize translations by template structure
3. **Fallback content**: Always provide fallback content for missing translations

### Performance

1. **Regenerate registry**: After adding/removing templates
2. **Watch mode**: Use `--watch` flag during development
3. **Cache templates**: Let the router handle template caching

## Troubleshooting

### Route Not Working

**Symptoms**: 404 errors for expected routes

**Solutions**:
```bash
# 1. Regenerate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# 2. Check file structure
ls -la app/your/directory/

# 3. Verify template syntax
templ generate

# 4. Check for naming issues
# Ensure dynamic parameters use _ suffix
# Ensure locale_ is used for internationalization
```

### Layout Not Loading

**Symptoms**: Template renders without expected layout

**Solutions**:
```bash
# 1. Check layout inheritance
# Does a layout exist in the current directory?
# Does a layout exist in a parent directory?

# 2. Verify layout template syntax
templ generate

# 3. Check layout template signature
# Layout should accept (title string, content templ.Component)
```

### Dynamic Parameters Not Working

**Symptoms**: Parameters not being passed to templates

**Solutions**:
```bash
# 1. Check parameter naming
# Use id_ not id
# Use slug_ not slug

# 2. Verify directory structure
# app/user/id_/page.templ  ✓
# app/user/id/page.templ   ✗

# 3. Regenerate registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject
```

### Internationalization Issues

**Symptoms**: Locale not being detected or applied

**Solutions**:
```bash
# 1. Check locale configuration
# Verify TR_I18N_SUPPORTED_LOCALES is set
# Verify TR_I18N_DEFAULT_LOCALE is set

# 2. Check directory structure
# Use locale_ directory for internationalized routes
# app/locale_/page.templ  ✓
# app/en/page.templ       ✗

# 3. Verify translation files
# Check .templ.yaml files for i18n sections
```

## Advanced Features

### Route Groups

Group related routes with shared configuration:

```yaml
# app/admin/page.templ.yaml
metadata:
  section: "admin"
  theme: "admin"

auth:
  type: "AdminRequired"

# This configuration applies to all routes in admin/
```

### Conditional Routing

Use metadata for conditional route behavior:

```yaml
# app/feature-flag/page.templ.yaml
metadata:
  feature_flag: "new_dashboard"
  experimental: true

auth:
  type: "UserRequired"
  roles: ["beta_tester"]
```

### Component Metadata Precedence

Components can override page metadata:

```yaml
# app/components/footer.templ.yaml
metadata:
  company_name: "My Company"
  version: "2.0.0"

# When included in pages, component metadata takes precedence
```

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first project
- **[Template Generator](TEMPLATE-GENERATOR.md)** - Learn about trgen CLI tool
- **[Authentication](AUTHENTICATION.md)** - Add authentication to routes
- **[Internationalization](INTERNATIONALIZATION.md)** - Comprehensive i18n guide
- **[Configuration](CONFIGURATION.md)** - Configure routing behavior

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features