# Template System

**Complete guide to the templ template system and advanced template features in Templ Router.**

## Overview

Templ Router uses the [templ](https://templ.guide/) template engine for type-safe, compiled HTML templates. The system provides automatic template discovery, layout inheritance, metadata integration, and seamless data service injection.

**Key Features:**
- Type-safe compiled templates
- Layout inheritance system
- Automatic template registry generation
- Metadata-driven configuration
- Data service integration
- Component-based architecture
- Internationalization support

## Configuration Prefix Notice

**Important:** Some environment variables in this documentation use the default prefix `TR_`. This prefix is **configurable** when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_TEMPLATE_CACHE_ENABLED=true`
- Custom: `MYAPP_TEMPLATE_CACHE_ENABLED=true`
- Multiple apps: `APP1_TEMPLATE_CACHE_ENABLED=true` and `APP2_TEMPLATE_CACHE_ENABLED=true`

All environment variable examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix.

## Template Basics

### Template File Structure

Templates are Go files with the `.templ` extension:

```
app/
├── layout.templ              # Root layout template
├── page.templ                # Home page template
├── login/
│   └── page.templ           # Login page template
├── dashboard/
│   └── page.templ           # Dashboard page template
└── locale_/                  # Internationalized templates
    └── page.templ           # Localized home page
```

### Basic Template Syntax

```go
// app/page.templ
package main

templ Page() {
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <title>Welcome</title>
    </head>
    <body>
        <h1>Welcome to Templ Router!</h1>
        <p>This is a basic template example.</p>
    </body>
    </html>
}
```

### Template Naming Conventions

**Important:** Template functions must follow specific naming conventions:

```go
// app/page.templ - MUST be named "Page"
package main
templ Page() { ... }

// app/user/id_/page.templ - MUST be named "Page"
package main
templ Page(userID string) { ... }

// app/layout.templ - MUST be named "Layout"
package main
templ Layout(content templ.Component) { ... }

// app/error.templ - MUST be named "Error"
package main
templ Error(errCtx middleware.ErrorContext) { ... }
```

## Layout System

### Layout Templates

Layout templates provide the base HTML structure:

```go
// app/layout.templ
package main

import "github.com/a-h/templ"

// Layout template accepts content as a component
templ Layout(title string, content templ.Component) {
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <title>{ title }</title>
        <script src="https://cdn.tailwindcss.com"></script>
    </head>
    <body class="bg-gray-50">
        <header class="bg-white shadow">
            <nav class="container mx-auto px-4 py-3">
                <div class="flex justify-between items-center">
                    <h1 class="text-xl font-bold">My App</h1>
                    <!-- Navigation content -->
                </div>
            </nav>
        </header>

        <main class="container mx-auto py-6">
            { content }
        </main>

        <footer class="bg-gray-800 text-white py-4 mt-8">
            <div class="container mx-auto text-center">
                <p>&copy; 2024 My Company</p>
            </div>
        </footer>
    </body>
    </html>
}
```

### Page Templates with Layouts

Page templates use layouts for consistent structure:

```go
// app/dashboard/page.templ
package main

import (
    "github.com/denkhaus/templ-router/pkg/router/i18n"
)

templ Page() {
    // Content component that will be wrapped by layout
    DashboardContent()
}

templ DashboardContent() {
    <div class="bg-white rounded-lg shadow p-6">
        <h1 class="text-2xl font-bold mb-4">
            { i18n.T(ctx, "dashboard_title") }
        </h1>
        <p class="text-gray-600">
            { i18n.T(ctx, "dashboard_subtitle") }
        </p>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
            <div class="bg-blue-50 p-4 rounded">
                <h3 class="font-semibold text-blue-900">Users</h3>
                <p class="text-2xl font-bold text-blue-600">1,234</p>
            </div>
            <div class="bg-green-50 p-4 rounded">
                <h3 class="font-semibold text-green-900">Revenue</h3>
                <p class="text-2xl font-bold text-green-600">$12,345</p>
            </div>
            <div class="bg-purple-50 p-4 rounded">
                <h3 class="font-semibold text-purple-900">Orders</h3>
                <p class="text-2xl font-bold text-purple-600">567</p>
            </div>
        </div>
    </div>
}
```

### Layout Inheritance

Layouts follow hierarchical inheritance:

```
app/
├── layout.templ           # Root layout (fallback)
├── dashboard/
│   ├── layout.templ       # Dashboard layout (overrides root)
│   └── page.templ         # Uses dashboard layout
└── user/
    └── page.templ         # Uses root layout
```

## Template Metadata

### YAML Metadata Files

Each template can have a corresponding `.templ.yaml` metadata file:

```yaml
# app/dashboard/page.templ.yaml
metadata:
  page_title: "Dashboard"
  description: "Main dashboard page"
  theme: "dark"
  section: "main"

i18n:
  en:
    dashboard_title: "Dashboard"
    dashboard_subtitle: "Overview of your application"
    users: "Total Users"
    revenue: "Revenue"
    orders: "Orders"
  de:
    dashboard_title: "Dashboard"
    dashboard_subtitle: "Übersicht Ihrer Anwendung"
    users: "Benutzer gesamt"
    revenue: "Umsatz"
    orders: "Bestellungen"

auth:
  type: "UserRequired"
  redirect_url: "/login"

data_services:
  - "DashboardDataService"
```

### Accessing Metadata in Templates

```go
package main

import (
    "github.com/denkhaus/templ-router/pkg/router/metadata"
    "github.com/denkhaus/templ-router/pkg/router/i18n"
)

templ Page() {
    pageTitle := metadata.M(ctx, "page_title")
    theme := metadata.M(ctx, "theme")

    <div class="dashboard" data-theme={ theme }>
        <h1>{ pageTitle }</h1>
        <p>{ i18n.T(ctx, "dashboard_subtitle") }</p>

        <div class="stats">
            <div class="stat">
                <span>{ i18n.T(ctx, "users") }:</span>
                <span class="value">1,234</span>
            </div>
            <div class="stat">
                <span>{ i18n.T(ctx, "revenue") }:</span>
                <span class="value">$12,345</span>
            </div>
        </div>
    </div>
}
```

## Components

### Reusable Components

Create reusable components that can be used across multiple templates:

```go
// app/components/navbar.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ Navbar() {
    <nav class="bg-blue-600 text-white">
        <div class="container mx-auto px-4">
            <div class="flex justify-between items-center h-16">
                <div class="flex items-center">
                    <span class="font-bold text-xl">MyApp</span>
                </div>

                <div class="flex space-x-4">
                    <a href="/" class="hover:bg-blue-700 px-3 py-2 rounded">
                        { i18n.T(ctx, "nav_home") }
                    </a>
                    <a href="/dashboard" class="hover:bg-blue-700 px-3 py-2 rounded">
                        { i18n.T(ctx, "nav_dashboard") }
                    </a>
                    <a href="/profile" class="hover:bg-blue-700 px-3 py-2 rounded">
                        { i18n.T(ctx, "nav_profile") }
                    </a>
                </div>
            </div>
        </div>
    </nav>
}

// app/components/footer.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/metadata"

templ Footer() {
    companyName := metadata.M(ctx, "company_name")
    currentYear := time.Now().Year()

    <footer class="bg-gray-800 text-white py-6">
        <div class="container mx-auto text-center">
            <p>&copy; { currentYear } { companyName }. All rights reserved.</p>
        </div>
    </footer>
}
```

### Component Metadata

#### Self-Contained Components

Components can be **self-contained** with their own metadata and internationalization. This is a powerful feature that makes components truly reusable:

- **No duplication required** - Component translations and metadata are defined once
- **Multiple usage** - Components can be used across different pages without repeating configuration
- **Independence** - Components work standalone with their own complete configuration
- **Consistency** - Same component behaves identically across all pages

```yaml
# app/components/navbar.templ.yaml
metadata:
  component_type: "navigation"
  background_color: "blue"
  show_search: true
  brand_name: "MyApp"

i18n:
  en:
    nav_home: "Home"
    nav_dashboard: "Dashboard"
    nav_profile: "Profile"
    nav_search: "Search"
    nav_logout: "Logout"
  de:
    nav_home: "Startseite"
    nav_dashboard: "Dashboard"
    nav_profile: "Profil"
    nav_search: "Suche"
    nav_logout: "Abmelden"
  fr:
    nav_home: "Accueil"
    nav_dashboard: "Tableau de bord"
    nav_profile: "Profil"
    nav_search: "Rechercher"
    nav_logout: "Déconnexion"
```

**Benefits of Self-Contained Components:**

1. **True Reusability** - Use the same component across multiple pages
2. **Single Source of Truth** - All translations and metadata in one place
3. **Easy Maintenance** - Update component config once, affects all usages
4. **No Duplication** - Avoid repeating translations in every page
5. **Consistent Behavior** - Component works identically everywhere
6. **Independent Testing** - Components can be tested in isolation

**Usage Example:**
```go
// In any page template - no additional configuration needed
<Navbar />
// The Navbar component automatically uses its own translations and metadata
```

Without self-contained components, you would need to repeat all translations in every page's `.templ.yaml` file. With this feature, you define it once in the component and reuse it everywhere.

### Component Routes

Components are accessible via their own routes:

```go
// Component accessible at /components/navbar
// Renders without layout when accessed directly
```

## Data Services Integration

### Templates with Data Services

Templates can automatically receive data from data services:

```go
// app/user/id_/page.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ Page(user *UserData) {
    <div class="user-profile">
        <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center space-x-4">
                <div class="w-16 h-16 bg-gray-300 rounded-full"></div>
                <div>
                    <h1 class="text-2xl font-bold">{ user.Name }</h1>
                    <p class="text-gray-600">{ user.Email }</p>
                </div>
            </div>

            <div class="mt-6 grid grid-cols-2 gap-4">
                <div>
                    <h3 class="font-semibold">User ID</h3>
                    <p class="text-gray-600">{ user.ID }</p>
                </div>
                <div>
                    <h3 class="font-semibold">Email</h3>
                    <p class="text-gray-600">{ user.Email }</p>
                </div>
            </div>

            <div class="mt-6">
                <a href={ i18n.LocalizeSafeURL(ctx, "/user/" + user.ID + "/edit") }
                   class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                    { i18n.T(ctx, "edit_profile") }
                </a>
            </div>
        </div>
    </div>
}
```

### Composite Data Services

For multiple data sources, create composite data services:

```go
// pkg/dataservices/dashboard_service.go
package dataservices

type DashboardData struct {
    UserStats    *UserData    `json:"user_stats"`
    SystemStats  *SystemData  `json:"system_stats"`
}

type DashboardDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*DashboardData, error)
}

func (s *dashboardDataService) GetData(routerCtx interfaces.RouterContext) (*DashboardData, error) {
    // Get data from multiple services
    userStats, _ := s.userService.GetData(routerCtx)
    systemStats, _ := s.systemService.GetData(routerCtx)

    return &DashboardData{
        UserStats:   userStats,
        SystemStats: systemStats,
    }, nil
}

// Template uses composite data
templ Page(data *DashboardData) {
    <div class="dashboard">
        <h1>{ data.UserStats.Name }</h1>
        <p>Total users: { data.SystemStats.TotalUsers }</p>
    </div>
}
```

## Internationalization in Templates

### Translation Functions

Use i18n functions in templates:

```go
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ Page() {
    <div class="content">
        <h1>{ i18n.T(ctx, "welcome_title") }</h1>
        <p>{ i18n.T(ctx, "welcome_message") }</p>

        <!-- Nested translations -->
        <nav>
            <a href="/">{ i18n.T(ctx, "nav.home") }</a>
            <a href="/dashboard">{ i18n.T(ctx, "nav.dashboard") }</a>
        </nav>

        <!-- Conditional content based on locale -->
        if i18n.GetCurrentLocale(ctx) == "de" {
            <p>Dies ist eine deutsche Nachricht</p>
        } else {
            <p>This is an English message</p>
        }
    </div>
}
```

### Localized URLs

Generate URLs with locale prefixes:

```go
templ Navigation() {
    <nav class="navigation">
        <a href={ i18n.LocalizeSafeURL(ctx, "/dashboard") }>
            { i18n.T(ctx, "dashboard") }
        </a>
        <a href={ i18n.LocalizeSafeURL(ctx, "/profile") }>
            { i18n.T(ctx, "profile") }
        </a>

        <!-- Language switcher -->
        <div class="language-switcher">
            currentLocale := i18n.GetCurrentLocale(ctx)
            currentRoute := i18n.GetCurrentRouteWithoutLocale(ctx)

            if currentLocale == "en" {
                <a href={ "/de" + currentRoute }>Deutsch</a>
            } else {
                <a href={ "/en" + currentRoute }>English</a>
            }
        </div>
    </nav>
}
```

## Advanced Template Features

### Conditional Rendering

```go
templ UserProfile(user *UserData) {
    <div class="profile">
        <h1>{ user.Name }</h1>

        <!-- Conditional content -->
        if user.IsAdmin {
            <div class="admin-panel">
                <h2>Admin Panel</h2>
                <button>Manage Users</button>
            </div>
        }

        <!-- Show different content based on user status -->
        if user.IsActive {
            <div class="status active">
                <span class="text-green-600">✓ Active</span>
            </div>
        } else {
            <div class="status inactive">
                <span class="text-red-600">✗ Inactive</span>
            </div>
        }

        <!-- Show optional fields -->
        if user.Department != "" {
            <p>Department: { user.Department }</p>
        }
    </div>
}
```

### Loops and Iteration

```go
templ ProductList(products []Product) {
    <div class="product-list">
        <h1>Products ({ len(products) })</h1>

        <div class="grid grid-cols-3 gap-4">
            for _, product := range products {
                <div class="product-card">
                    <h3>{ product.Name }</h3>
                    <p class="text-gray-600">{ product.Description }</p>
                    <p class="font-bold">${ product.Price }</p>

                    if product.InStock {
                        <button class="bg-green-600 text-white px-4 py-2 rounded">
                            { i18n.T(ctx, "add_to_cart") }
                        </button>
                    } else {
                        <button class="bg-gray-400 text-gray-700 px-4 py-2 rounded" disabled>
                            { i18n.T(ctx, "out_of_stock") }
                        </button>
                    }
                </div>
            }
        </div>
    </div>
}
```

### Template Composition

```go
// Base template
templ BaseCard(title string, content templ.Component) {
    <div class="card">
        <div class="card-header">
            <h2>{ title }</h2>
        </div>
        <div class="card-body">
            { content }
        </div>
    </div>
}

// Using composed template
templ UserCard(user *UserData) {
    BaseCard(user.Name, UserCardDetails(user))
}

templ UserCardDetails(user *UserData) {
    <div>
        <p>Email: { user.Email }</p>
        <p>Joined: { user.CreatedAt.Format("Jan 2, 2006") }</p>
    </div>
}
```

### Template Functions and Helpers

```go
// Helper functions in templates
func formatDate(t time.Time) string {
    return t.Format("January 2, 2006")
}

func formatPrice(price float64) string {
    return fmt.Sprintf("$%.2f", price)
}

templ ProductDetail(product *Product) {
    <div class="product-detail">
        <h1>{ product.Name }</h1>
        <p class="price">{ formatPrice(product.Price) }</p>
        <p class="date">Added: { formatDate(product.CreatedAt) }</p>

        <!-- Custom formatting -->
        <div class="badge">
            if product.IsFeatured {
                <span class="bg-yellow-100 text-yellow-800 px-2 py-1 rounded">
                    { i18n.T(ctx, "featured") }
                </span>
            }
        </div>
    </div>
}
```

## Template Generation

### Command Line Generation

Generate templates from `.templ` files:

```bash
# Generate all templates
templ generate

# Generate specific file
templ generate app/page.templ

# Watch mode for development
templ generate --watch
```

### Template Registry Generation

Generate the template registry for route discovery:

```bash
# Generate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Watch mode for automatic updates
trgen --scan-path=app --module-name=github.com/youruser/yourproject --watch
```

### Generated Registry Structure

The generated registry includes:

```go
// generated/templates/registry.go
package templates

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
    "github.com/samber/do/v2"
)

type Registry struct {
    injector *do.Injector
    templates map[string]interfaces.TemplateFunc
}

func NewRegistry(injector *do.Injector) (*Registry, error) {
    registry := &Registry{
        injector: injector,
        templates: make(map[string]interfaces.TemplateFunc),
    }

    // Templates are auto-registered by trgen
    return registry, nil
}

func (r *Registry) GetTemplate(name string) (interfaces.TemplateFunc, error) {
    if fn, exists := r.templates[name]; exists {
        return fn, nil
    }
    return nil, fmt.Errorf("template not found: %s", name)
}
```

## Error Handling in Templates

### Error Templates

Create custom error templates:

```go
// app/error.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/middleware"

templ Error(errCtx middleware.ErrorContext) {
    <div class="error-page">
        <div class="error-container">
            <h1 class="error-title">Error { errCtx.StatusCode }</h1>
            <p class="error-message">{ errCtx.Message }</p>

            if errCtx.Details != "" {
                <details class="error-details">
                    <summary>Details</summary>
                    <pre>{ errCtx.Details }</pre>
                </details>
            }

            <div class="error-actions">
                <a href="/" class="btn btn-primary">
                    Go Home
                </a>
                <button onclick="history.back()" class="btn btn-secondary">
                    Go Back
                </button>
            </div>
        </div>
    </div>
}
```

### Graceful Degradation

Handle missing data gracefully:

```go
templ UserProfile(user *UserData) {
    <div class="profile">
        if user != nil {
            <h1>{ user.Name }</h1>
            <p>{ user.Email }</p>
        } else {
            <div class="error">
                <h1>User Not Found</h1>
                <p>The requested user could not be found.</p>
            </div>
        }
    </div>
}

// Alternative using templ's error handling
templ UserProfileWithError() {
    if user, err := getUserData(ctx); err != nil {
        <div class="error">
            <h1>Error</h1>
            <p>{ err.Error() }</p>
        </div>
    } else {
        <div class="profile">
            <h1>{ user.Name }</h1>
            <p>{ user.Email }</p>
        </div>
    }
}
```

## Performance Optimization

### Template Caching

Enable template caching in production:

```bash
# Enable template caching
TR_TEMPLATE_CACHE_ENABLED=true
TR_TEMPLATE_CACHE_SIZE=100
TR_TEMPLATE_CACHE_TTL=1h
```

### Component Caching

Cache expensive components:

```go
// Cache expensive computations
templ ExpensiveComponent(data *ComplexData) {
    {{
        // Cache computation results
        cachedResult := computeExpensiveValue(data)
    }}

    <div class="expensive">
        <!-- Use cached result -->
        <p>Result: { cachedResult }</p>
    </div>
}
```

### Minimize Re-renders

Use conditional rendering to avoid unnecessary work:

```go
templ Dashboard(data *DashboardData) {
    // Only render expensive sections if data is available
    if data != nil {
        <div class="dashboard-stats">
            <h1>Dashboard Stats</h1>
            <!-- Expensive rendering -->
        </div>
    }
}
```

## Best Practices

### Template Organization

1. **Keep templates focused** on a single responsibility
2. **Use layouts** for consistent structure
3. **Create reusable components** for common UI elements
4. **Separate concerns** between presentation and business logic
5. **Use descriptive names** for templates and functions

### Performance

1. **Enable template caching** in production
2. **Minimize template complexity** for faster rendering
3. **Use conditional rendering** to avoid unnecessary work
4. **Profile templates** to identify bottlenecks
5. **Optimize data structures** passed to templates

### Maintainability

1. **Document complex templates** with comments
2. **Use consistent formatting** and structure
3. **Extract common patterns** into components
4. **Test templates** with sample data
5. **Keep templates small** and focused

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Create your first template
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Template-based routing
- **[Data Services](DATA-SERVICES.md)** - Data integration with templates
- **[Internationalization](INTERNATIONALIZATION.md)** - Multi-language templates
- **[Middleware](MIDDLEWARE.md)** - Template rendering middleware

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Templ Documentation](https://templ.guide/)** - Official templ documentation