# File-Based Routing System

**Automatic route generation from file structure, similar to Next.js.**

## Overview

Templ Router generates HTTP routes based on your file structure. The `trgen` CLI tool scans template files and creates a registry mapping URLs to templates automatically.

**Core Features:**
- Automatic route generation from file structure
- Dynamic parameters with `_` suffix (`id_/`, `locale_/`)
- Route precedence system for conflict resolution
- Internationalized routes via `locale_/` directory structure
- Template-to-route mapping with configurable patterns

## Configuration

### Environment Prefix

Router services use a configurable prefix (default: `TR`):

```go
// Default prefix
container.RegisterRouterServices("TR")  // Uses TR_* variables

// Custom prefix
container.RegisterRouterServices("MYAPP")  // Uses MYAPP_* variables
```

## Directory Structure

Routes are generated from your template directory structure:

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
│   ├── page.templ            # /{locale}
│   └── user/
│       └── id_/
│           └── page.templ    # /{locale}/user/{id}
└── components/
    ├── header.templ          # Component (no route)
    └── footer.templ          # Component (no route)
```

## Route Generation Rules

### Basic Routes

- `app/page.templ` → `/`
- `app/about.templ` → `/about`
- `app/login/page.templ` → `/login`
- `app/dashboard/settings/page.templ` → `/dashboard/settings`

### Dynamic Parameters

Use `_` suffix for dynamic segments:

- `app/user/id_/page.templ` → `/user/{id}`
- `app/blog/category_/slug_/page.templ` → `/blog/{category}/{slug}`
- `app/api/v1/users/user_id_/page.templ` → `/api/v1/users/{user_id}`

### Internationalization

The `locale_/` directory creates internationalized routes:

- `app/locale_/page.templ` → `/{locale}` (e.g., `/en`, `/de`)
- `app/locale_/dashboard/page.templ` → `/{locale}/dashboard`
- `app/locale_/user/id_/page.templ` → `/{locale}/user/{id}`

### Supported Locales

Configure supported locales:

```bash
# Environment variables
TR_I18N_SUPPORTED_LOCALES=en,de,fr
TR_I18N_DEFAULT_LOCALE=en
TR_I18N_FALLBACK_LOCALE=en
```

## Route Precedence

Routes are prioritized in this order:

1. **Static routes** - `/about`, `/contact`
2. **Dynamic routes** - `/user/{id}`, `/blog/{slug}`
3. **Wildcard routes** - `/{fallback}`
4. **Localized routes** - `/{locale}/page`

### Conflict Resolution

```
app/
├── about.templ           → /about (higher priority)
└── slug_/page.templ      → /{slug} (lower priority)
```

The static route `/about` takes precedence over the dynamic `/{slug}`.

## Template Types

### Page Templates

```go
// app/page.templ
package main

templ Page() {
    <h1>Home Page</h1>
}
```

**Generated Route:**
```go
"/": {
    TemplateFunc: Page,
    Parameters: []string{},
    IsLocalized: false,
}
```

### Layout Templates

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

Layouts are used by pages but don't generate routes themselves.

### Component Templates

```go
// app/components/header.templ
package components

templ Header(title string) {
    <header><h1>{ title }</h1></header>
}
```

Components are reusable and don't generate routes.

## Route Metadata

Templates can have associated `.templ.yaml` files:

```yaml
# app/login/page.templ.yaml
metadata:
  title: "Login Page"
  description: "User authentication"

auth:
  type: "UserRequired"
  redirect_url: "/login"

i18n:
  en:
    page_title: "Login"
    submit_button: "Sign In"
  de:
    page_title: "Anmelden"
    submit_button: "Anmelden"
```

### Metadata Properties

#### Authentication

```yaml
auth:
  type: "Public"           # No authentication required
  type: "UserRequired"     # Any authenticated user
  type: "AdminRequired"    # Admin users only
  redirect_url: "/login"   # Redirect URL for unauthenticated users
  roles: ["admin", "moderator"]  # Specific roles required
```

#### Internationalization

```yaml
i18n:
  en:
    welcome: "Welcome"
    button_text: "Get Started"
  de:
    welcome: "Willkommen"
    button_text: "Loslegen"
  fr:
    welcome: "Bienvenue"
    button_text: "Commencer"
```

#### Route Configuration

```yaml
metadata:
  title: "Page Title"
  description: "Page description"
  keywords: ["keyword1", "keyword2"]
  custom_data: "Custom metadata value"
```

## Dynamic Parameters

### Parameter Types

#### Single Parameters

```
app/user/id_/page.templ
```

**Generated Route:** `/user/{id}`

**Access in Template:**
```go
func (s *service) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")
    return s.getUserByID(userID)
}
```

#### Multiple Parameters

```
app/blog/category_/slug_/page.templ
```

**Generated Route:** `/blog/{category}/{slug}`

**Access in Template:**
```go
func (s *service) GetData(routerCtx interfaces.RouterContext) (*BlogPost, error) {
    category := routerCtx.GetURLParam("category")
    slug := routerCtx.GetURLParam("slug")
    return s.getBlogPost(category, slug)
}
```

#### Optional Parameters

```
app/search/query_/page.templ
app/search/page.templ  # Fallback for no query
```

**Routes:** `/search/{query}` and `/search`

### Parameter Validation

```yaml
# app/user/id_/page.templ.yaml
parameters:
  id:
    type: "int"
    required: true
    min: 1
    pattern: "^[0-9]+$"
```

## Internationalization

### Locale Directory Structure

```
app/
├── locale_/                    # Internationalized routes
│   ├── page.templ            # /{locale}
│   ├── about.templ           # /{locale}/about
│   └── dashboard/
│       └── page.templ        # /{locale}/dashboard
└── page.templ                # / (non-localized fallback)
```

### Locale Detection

1. **URL Path** - `/en/page`, `/de/page`
2. **Query Parameter** - `?locale=en`
3. **Accept-Language Header** - `Accept-Language: en-US,en;q=0.9`
4. **Cookie** - `locale=en`
5. **Default Locale** - Configured fallback

### Translation Files

```yaml
# app/locale_/about/page.templ.yaml
i18n:
  en:
    page_title: "About Us"
    company_description: "We are a technology company..."
  de:
    page_title: "Über Uns"
    company_description: "Wir sind ein Technologieunternehmen..."
  fr:
    page_title: "À Propos"
    company_description: "Nous sommes une entreprise technologique..."
```

### Template Usage

```go
// In templates
templ Page() {
    <h1>{ i18n.T(ctx, "page_title") }</h1>
    <p>{ i18n.T(ctx, "company_description") }</p>
    <a href={ i18n.LocalizeSafeURL(ctx, "/contact") }>Contact</a>
}
```

## Route Generation Process

### trgen Execution

```bash
# Generate routes from file structure
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Watch mode for development
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch
```

### Generated Registry

The generated `registry.go` contains:

```go
// Route mappings
var routeMappings = map[string]*RouteInfo{
    "/": {
        Path:         "/",
        TemplateFunc: Page,
        Parameters:   []string{},
        IsLocalized:  false,
        RequiresAuth: false,
        DataServices: []string{},
    },
    "/login": {
        Path:         "/login",
        TemplateFunc: LoginPage,
        Parameters:   []string{},
        IsLocalized:  false,
        RequiresAuth: false,
        DataServices: []string{},
    },
    "/{locale}/user/{id}": {
        Path:         "/{locale}/user/{id}",
        TemplateFunc: UserProfilePage,
        Parameters:   []string{"locale", "id"},
        IsLocalized:  true,
        RequiresAuth: true,
        DataServices: []string{"UserDataService"},
    },
}

// Registry interface
type Registry interface {
    GetTemplate(path string, params map[string]string) (templ.Component, error)
    GetAllRoutes() []string
    GetRouteInfo(path string) (*RouteInfo, error)
}
```

## Advanced Features

### Route Groups

```go
// Group routes by middleware or functionality
router.Route("/admin", func(r chi.Router) {
    r.Use(AdminMiddleware)
    r.Handle("/dashboard", DashboardPage)
    r.Handle("/users", UsersPage)
})
```

### Custom Route Patterns

```yaml
# Custom route configuration
custom_routes:
  - pattern: "/api/v1/users/{id}"
    template: "api/users/user.templ"
    methods: ["GET", "PUT"]
    middleware: ["auth", "rate-limit"]
```

### Route Aliases

```yaml
# Create aliases for routes
aliases:
  "/home": "/"                # /home redirects to /
  "/profile": "/user/profile" # /profile redirects to /user/profile
```

## Performance Optimization

### Route Caching

Routes are pre-computed at startup for fast matching:

```go
// Pre-computed route tree for O(1) lookups
type RouteTree struct {
    static   map[string]*RouteNode
    dynamic  []*RouteNode
    wildcards []*RouteNode
}
```

### Parameter Parsing

Efficient parameter extraction:

```go
// Fast parameter parsing
func parseParams(route string, path string) map[string]string {
    // Optimized parameter extraction
}
```

## Best Practices

### File Organization

1. **Use descriptive names** - `user-profile.templ` vs `page.templ`
2. **Group related routes** - Use folders for logical grouping
3. **Consistent naming** - Follow naming conventions
4. **Separate concerns** - Pages, layouts, components in appropriate locations

### Route Design

1. **RESTful URLs** - Use standard URL patterns
2. **Clear structure** - Make URLs predictable and intuitive
3. **Avoid deep nesting** - Keep URLs reasonably shallow
4. **Use meaningful parameters** - Clear parameter names

### Internationalization

1. **Plan for i18n** - Use `locale_/` structure from the start
2. **Consistent translations** - Maintain translation quality
3. **Fallback handling** - Provide good fallbacks
4. **URL structure** - Design URLs with internationalization in mind

### Performance

1. **Minimize routes** - Avoid unnecessary route complexity
2. **Use caching** - Leverage built-in caching
3. **Optimize parameters** - Keep parameter validation simple
4. **Monitor performance** - Track route matching performance

## Troubleshooting

### Common Issues

**Route not found:**
- Check file structure matches expected route
- Verify `trgen` has been run
- Check for naming conflicts

**Parameter not working:**
- Ensure `_` suffix is used for dynamic segments
- Check parameter names match in template and metadata
- Verify parameter validation rules

**Internationalization issues:**
- Check supported locales configuration
- Verify translation file structure
- Ensure locale detection is working

### Debug Mode

```bash
# Enable verbose logging
trgen --verbose --scan-path=app --module-name=github.com/youruser/yourproject

# Check generated routes
trgen --scan-path=app --module-name=github.com/youruser/yourproject --dry-run
```

### Route Inspection

```go
// Inspect generated routes
registry := templates.NewRegistry(injector)
routes := registry.GetAllRoutes()
for _, route := range routes {
    info := registry.GetRouteInfo(route)
    fmt.Printf("Route: %s, Template: %s, Auth: %v\n",
        route, info.TemplateFunc, info.RequiresAuth)
}
```

---

**Related Documentation**: [Template Generator](TEMPLATE-GENERATOR.md), [Internationalization](INTERNATIONALIZATION.md), [Authentication](AUTHENTICATION.md)