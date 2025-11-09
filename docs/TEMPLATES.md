# Template System

**Type-safe compiled templates with layout inheritance and metadata integration.**

## Overview

Templ Router uses the [templ](https://templ.guide/) template engine for type-safe, compiled HTML templates with automatic discovery and data service integration.

**Key Features:**
- Type-safe compiled templates
- Layout inheritance system
- Automatic template registry generation
- Metadata-driven configuration
- Data service integration
- Component-based architecture

## Template Types

### Page Templates

Generate HTTP routes (must be named `Page`):

```go
// app/page.templ
package main

templ Page() {
    <h1>Welcome</h1>
    <p>This is the home page.</p>
}
```

### Layout Templates

Used for page structure (must be named `Layout`):

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
    <body>
        { content }
    </body>
    </html>
}
```

### Component Templates

Reusable components (custom names allowed):

```go
// app/components/header.templ
package components

templ Header(title string) {
    <header class="bg-blue-600 text-white p-4">
        <h1>{ title }</h1>
    </header>
}

// app/components/footer.templ
package components

templ Footer() {
    <footer class="bg-gray-800 text-white p-4">
        <p>© 2024 My Company</p>
    </footer>
}
```

## Template Metadata

### Metadata Files

Templates can have associated `.templ.yaml` files:

```yaml
# app/page.templ.yaml
metadata:
  title: "Home Page"
  description: "Welcome to our application"

auth:
  type: "Public"

i18n:
  en:
    welcome: "Welcome"
    page_description: "Welcome to our application"
  de:
    welcome: "Willkommen"
    page_description: "Willkommen bei unserer Anwendung"
```

### Component Metadata

Components can have their own metadata:

```yaml
# app/components/footer.templ.yaml
metadata:
  en:
    company_name: "My Company"
    region: "North America"
  de:
    company_name: "Meine Firma"
    region: "Europa"
  # Global fallback values
  version: "1.0.0"
  author: "Development Team"

i18n:
  en: { copyright: "© 2024 My Company" }
  de: { copyright: "© 2024 Meine Firma" }
```

> **Note**: Metadata can now be localized! See [Internationalization](INTERNATIONALIZATION.md) for complete locale-aware metadata documentation.

## Template Inheritance

### Layout Usage

Layouts are automatically applied by the router:

```go
// Page content - Layout is applied automatically
templ Page() {
    <div>
        <h1>{ i18n.T(ctx, "welcome") }</h1>
        <p>{ i18n.T(ctx, "page_description") }</p>
    </div>
}
```

### Component Integration

Components can be used in pages:

```go
templ Page() {
    @components.Header("Welcome")
    <main>
        <h1>Main Content</h1>
    </main>
    @components.Footer()
}
```

## Data Services

### Service Integration

Templates automatically access required data services:

```go
// UserDataService interface
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
}

// Template with data service
templ Page() {
    userData, _ := userService.GetData(routerCtx)
    <h1>Welcome, { userData.Name }!</h1>
}
```

### RouterContext Access

Access URL parameters and query parameters:

```go
templ Page() {
    // URL parameter: /user/{id}
    userID := routerCtx.GetURLParam("id")

    // Query parameter: ?page=5
    page := routerCtx.GetQueryParam("page")

    <h1>User { userID }, Page { page}</h1>
}
```

## Advanced Features

### Conditional Rendering

```go
templ Page() {
    @if user.IsAuthenticated {
        <div>
            <h1>Welcome back, { user.Name }!</h1>
        @components.UserMenu()
        </div>
    } else {
        <div>
            <h1>Please sign in</h1>
            @components.LoginForm()
        </div>
    }
}
```

### Loops and Iteration

```go
templ Page() {
    <h1>Product List</h1>
    @for _, product := range products {
        <div class="product">
            <h3>{ product.Name }</h3>
            <p>{ product.Description }</p>
            <span>${product.Price}</span>
        </div>
    }
}
```

### CSS Classes and Styling

```go
templ Page() {
    <div class={ "container mx-auto p-4 " + conditionalClass }>
        <h1 class="text-2xl font-bold mb-4">Title</h1>
        <p class="text-gray-600">Content</p>
    </div>
}
```

### Attributes and Properties

```go
templ Page() {
    <div id="main-content" class="content" data-testid="page">
        <h1 title="Page Title">Title</h1>
    </div>
}
```

## Internationalization

### Translation Integration

```go
templ Page() {
    <h1>{ i18n.T(ctx, "page_title") }</h1>
    <p>{ i18n.T(ctx, "welcome_message") }</p>
    <a href={ i18n.LocalizeSafeURL(ctx, "/about") }>
        { i18n.T(ctx, "about_link") }
    </a>
}
```

### Metadata-Based Translations

```yaml
# app/locale_/page.templ.yaml
i18n:
  en:
    page_title: "Welcome"
    welcome_message: "Welcome to our site"
  de:
    page_title: "Willkommen"
    welcome_message: "Willkommen auf unserer Seite"
```

## Error Handling

### Error Templates

```go
// app/error.templ
package main

templ Error() {
    <div class="error-container">
        <h1>{ i18n.T(ctx, "error_title") }</h1>
        <p>{ i18n.T(ctx, "error_message") }</p>
        <a href="/">{ i18n.T(ctx, "go_home") }</a>
    </div>
}
```

### Error Page Configuration

```yaml
# app/error.templ.yaml
metadata:
  title: "Error Page"

i18n:
  en:
    error_title: "Error Occurred"
    error_message: "Something went wrong"
    go_home: "Go Home"
  de:
    error_title: "Fehler Aufgetreten"
    error_message: "Etwas ist schief gelaufen"
    go_home: "Zur Startseite"
```

## Best Practices

### File Organization

1. **Consistent Naming** - Use descriptive file and folder names
2. **Logical Grouping** - Group related templates together
3. **Component Reuse** - Create reusable components for common UI elements
4. **Metadata Integration** - Use metadata for configuration and i18n

### Template Design

1. **Keep Templates Focused** - Single responsibility per template
2. **Use Components** - Break down complex templates into smaller components
3. **Data Service Integration** - Leverage automatic data injection
4. **Consistent Styling** - Use consistent CSS patterns

### Performance

1. **Compiled Templates** - Templates are compiled to Go for performance
2. **Automatic Caching** - Template results are cached automatically
3. **Minimal Overhead** - Avoid complex logic in templates
4. **Efficient Rendering** - Use efficient HTML structure

## Troubleshooting

### Common Issues

**Template not found:**
- Check file path and `.templ` extension
- Verify template compilation with `templ generate`
- Ensure template is in scanned directory

**Layout not applied:**
- Verify Layout template exists and is named `Layout`
- Check that Layout is in the correct location
- Ensure template registry is generated

**Data service not working:**
- Verify service is registered in DI container
- Check service interface implementation
- Ensure template requires the service

**Metadata not loaded:**
- Check `.templ.yaml` file exists and is properly formatted
- Verify metadata syntax is correct
- Ensure template registry is regenerated

### Debug Mode

```bash
# Generate templates with verbose output
templ generate --verbose

# Generate template registry
trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Check generated files
ls generated/templates/
```

### Template Testing

```go
func TestPageTemplate(t *testing.T) {
    // Test template compilation
    page := Page()

    // Test template rendering
    var buf strings.Builder
    page.Render(context.Background(), &buf)

    output := buf.String()
    assert.Contains(t, output, "<h1>")
}
```

---

**Related Documentation**: [File-Based Routing](FILE-BASED-ROUTING.md), [Template Generator](TEMPLATE-GENERATOR.md), [Data Services](DATA-SERVICES.md)