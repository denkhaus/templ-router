# File-Based Routing System

This document explains how to use the new file-based routing system that automatically generates HTTP routes based on your folder structure and *.templ files, similar to Next.js app router.

## Overview

The file-based router automatically discovers `*.templ` files in the `app/` directory and generates corresponding HTTP routes. You can enhance these routes with optional `*.yaml` metadata files that allow customizing paths, auth settings, internationalization, and more.

## Directory Structure

Routes are generated based on the folder structure in the `app/` directory. The basic structure looks like:

```
app/
├── layout.templ          # Root layout (required)
├── page.templ            # Root page (optional)
├── dashboard/
│   ├── layout.templ      # Dashboard layout
│   ├── page.templ        # Dashboard page
│   └── create.templ      # Create page for dashboard
└── user/
    ├── $id/
    │   └── page.templ    # For /user/{id} (dynamic route)
    └── layout.templ      # Layout for user routes
```

### Key File Conventions

- `layout.templ`: Defines the layout structure for a directory and its subdirectories
- `page.templ`: The main content for a directory's route
- `error.templ`: Error page for the directory and subdirectories
- `*.templ`: Component templates for specific routes

## Route Generation

### Static Routes

Static routes are generated based on the directory structure:

- `app/page.templ` → `/`
- `app/dashboard/page.templ` → `/dashboard`
- `app/settings/profile.templ` → `/settings/profile`

### Dynamic Routes

Dynamic routes use the dollar sign (`$`) convention:

- `app/user/$id.templ` → `/user/{id}`
- `app/product/$slug.templ` → `/product/{slug}`
- `app/$locale/dashboard.templ` → `/{locale}/dashboard` (special case for localization)

> **Note**: Only `$locale` is reserved for localization; all other parameters use the `$` prefix (e.g., `$id`, `$slug`).

## Layout Inheritance

Layouts are inherited hierarchically:

1. Layouts in subdirectories override parent layouts
2. If no layout exists in a directory, the system uses the layout from the nearest parent directory
3. The root `app/layout.templ` serves as the fallback layout for the entire application

## YAML Metadata

You can customize routes using optional YAML metadata files named after your template files with a `.yaml` extension:

```
app/
├── dashboard/
│   ├── page.templ
│   ├── page.templ.yaml    # Metadata for page.templ
│   └── layout.templ
```

Example YAML file (`app/dashboard/page.templ.yaml`):

```yaml
route:
  path: /dashboard/custom-path    # Custom route path
  auth:
    required: true
    permissions: ["admin"]

i18n:
  title: "admin.dashboard.title"
  description: "admin.dashboard.description"
  labels:
    submit: "admin.dashboard.submit"
    cancel: "admin.dashboard.cancel"

metadata:
  seo:
    title: "Admin Dashboard"
    description: "Administrative dashboard for managing the application"
```

## Internationalization (i18n)

### Automatic i18n Identifier Generation

The system automatically generates i18n identifiers based on the file structure using an opinionated schema:

- `app/admin/dashboard/create.templ` → `admin.dashboard.create.title`, `admin.dashboard.create.description`, etc.

### Locale-Based Routes

Support dynamic locale switching with routes like:
- `/en/admin/dashboard`
- `/de/admin/dashboard`
- `/fr/admin/dashboard`

Define locale-specific content using the `$locale` parameter:

```
app/
└── $locale/
    └── admin/
        ├── layout.templ
        └── dashboard/
            └── page.templ
```

## Authentication

Authentication can be configured at the directory level or for individual templates:

### Directory-Level Auth

Each directory can define default auth requirements that apply to all templates within that directory.

### Template-Level Auth

YAML metadata can override directory-level auth settings for specific templates.

## Error Handling

Error templates are applied hierarchically:

1. Closest error templates override those in parent directories
2. `app/error.templ` serves as the global fallback error template

## Integration with Manual Routes

Manual routes registered before file-based routes take precedence, allowing you to override specific routes while keeping the rest of the file-based routing system.

Example:

```go
// In main.go
r := chi.NewRouter()

// Register manual routes first (they take precedence)
r.Get("/special", specialHandler)

// Initialize and register file-based routes
fileRouter := router.New("app")
fileRouter.RegisterRoutes(r)
```

## Dependency Injection

The router integrates with the existing DI container through the `di.RouterProvider`:

```go
// Register with DI container
routerProvider := di.NewRouterProvider("app")
container.Register(routerProvider)
```

## Best Practices

1. Use the standard file naming conventions (`layout.templ`, `page.templ`, etc.) for predictable route generation
2. Leverage YAML metadata for route customization instead of hardcoding paths
3. Organize your directory structure logically to take advantage of layout inheritance
4. Use the `$` prefix for dynamic route parameters
5. Remember that `$locale` is reserved for localization
6. Manual routes take precedence, so register them before file-based routes when needed

## Troubleshooting

### Route Not Found
- Verify that the `*.templ` file exists in the correct location
- Confirm that `layout.templ` exists at the app root
- Check that the file name matches the expected route structure

### Layout Not Loading
- Ensure `layout.templ` exists in the appropriate directory
- Check the layout inheritance pattern (child layouts override parent layouts)

### i18n Keys Not Working
- Verify that your i18n keys exist in the appropriate locale files
- Confirm the YAML metadata file is properly formatted
- Ensure the i18n system is properly initialized in your application