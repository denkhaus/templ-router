# Authentication & Authorization

**Complete guide to implementing authentication and authorization in Templ Router.**

## Overview

Templ Router provides a comprehensive authentication system with session-based authentication, role-based access control, and built-in API endpoints for user management. The authentication system integrates seamlessly with the file-based routing system through `.templ.yaml` metadata files.

**Key Features:**
- Three authentication types: `Public`, `UserRequired`, `AdminRequired`
- Session-based authentication with configurable expiry
- Built-in authentication API endpoints
- Role-based access control
- Hierarchical authentication configuration
- Configurable redirect routes
- Support for internationalized authentication flows

## Configuration Prefix

**Important:** All environment variables in this documentation use the default prefix `TR_`. This prefix is **configurable** when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_AUTH_SESSION_EXPIRY=24h`
- Custom: `MYAPP_AUTH_SESSION_EXPIRY=24h`
- Multiple apps: `APP1_AUTH_SESSION_EXPIRY=24h` and `APP2_AUTH_SESSION_EXPIRY=24h`

All examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix.

## Authentication Types

### 1. Public Access

Default access level - no authentication required:

```yaml
# app/public/page.templ.yaml
auth:
  type: "Public"
```

Or simply omit the auth section:

```yaml
# app/public/page.templ.yaml
metadata:
  page_title: "Public Page"

i18n:
  en:
    page_title: "Welcome!"
  de:
    page_title: "Willkommen!"
```

### 2. User Required

Any authenticated user can access:

```yaml
# app/dashboard/page.templ.yaml
auth:
  type: "UserRequired"
  redirect_url: "/login"

metadata:
  page_title: "Dashboard"

i18n:
  en:
    page_title: "Dashboard"
  de:
    page_title: "Dashboard"
```

### 3. Admin Required

Only users with admin privileges can access:

```yaml
# app/admin/page.templ.yaml
auth:
  type: "AdminRequired"
  redirect_url: "/login"
  roles: ["admin", "super_admin"]  # Optional: specific roles required

metadata:
  page_title: "Admin Panel"

i18n:
  en:
    page_title: "Admin Panel"
  de:
    page_title: "Admin-Panel"
```

## Configuration

### Environment Variables

Configure authentication behavior through environment variables with configurable prefix:

```bash
# Session Configuration
TR_AUTH_SESSION_EXPIRY=24h                # Session duration
TR_AUTH_SESSION_COOKIE_NAME=session_id     # Session cookie name
TR_AUTH_SESSION_COOKIE_SECURE=true         # HTTPS-only cookies
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true      # HTTP-only cookies

# Redirect Configuration
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard    # After successful login
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/welcome      # After successful registration
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/           # After logout
TR_AUTH_SIGNIN_ROUTE=/login               # Custom login route

# User Management
TR_AUTH_CREATE_DEFAULT_ADMIN=true         # Create default admin user
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
TR_AUTH_DEFAULT_ADMIN_PASSWORD=admin123   # Only for development
```

### Internationalized Redirects

Support for internationalized success routes with locale parameter:

```bash
# Locale-aware redirects
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/{locale}/dashboard
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/{locale}/welcome
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/{locale}/

# The {locale} parameter is automatically replaced with the current locale:
# /{locale}/dashboard → /en/dashboard or /de/dashboard
```

## Built-in Authentication API Endpoints

Templ Router automatically provides authentication API endpoints:

### Available Endpoints

```bash
POST /api/auth/signin      # User sign in
POST /api/auth/signout     # User sign out
POST /api/auth/signup      # User registration
GET  /api/auth/me          # Get current user info
```

### Sign In Endpoint

```bash
POST /api/auth/signin
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "userpassword"
}
```

**Response:**
```json
{
  "success": true,
  "redirect_url": "/dashboard",
  "user": {
    "id": "123",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user"
  }
}
```

### Sign Up Endpoint

```bash
POST /api/auth/signup
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "Jane Doe"
}
```

**Response:**
```json
{
  "success": true,
  "redirect_url": "/welcome",
  "user": {
    "id": "456",
    "email": "newuser@example.com",
    "name": "Jane Doe",
    "role": "user"
  }
}
```

### Sign Out Endpoint

```bash
POST /api/auth/signout
```

**Response:**
```json
{
  "success": true,
  "redirect_url": "/"
}
```

### Current User Endpoint

```bash
GET /api/auth/me
Authorization: Bearer <session_cookie>
```

**Response:**
```json
{
  "user": {
    "id": "123",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "last_login": "2024-01-15T10:30:00Z"
  }
}
```

## Template Implementation

### Login Page

```go
// app/login/page.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ LoginPage() {
    <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div class="max-w-md w-full space-y-8">
            <div>
                <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                    { i18n.T(ctx, "login_title") }
                </h2>
                <p class="mt-2 text-center text-sm text-gray-600">
                    { i18n.T(ctx, "login_subtitle") }
                </p>
            </div>
            <form class="mt-8 space-y-6"
                  hx-post="/api/auth/signin"
                  hx-target="#form-error"
                  hx-swap="innerHTML">
                <input type="hidden" name="remember" value="true"/>
                <div class="rounded-md shadow-sm -space-y-px">
                    <div>
                        <label for="email-address" class="sr-only">
                            { i18n.T(ctx, "email") }
                        </label>
                        <input id="email-address"
                               name="email"
                               type="email"
                               autocomplete="email"
                               required
                               class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                               placeholder={ i18n.T(ctx, "email") }/>
                    </div>
                    <div>
                        <label for="password" class="sr-only">
                            { i18n.T(ctx, "password") }
                        </label>
                        <input id="password"
                               name="password"
                               type="password"
                               autocomplete="current-password"
                               required
                               class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                               placeholder={ i18n.T(ctx, "password") }/>
                    </div>
                </div>

                <div id="form-error" class="text-red-600 text-sm"></div>

                <div>
                    <button type="submit"
                            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                        { i18n.T(ctx, "sign_in") }
                    </button>
                </div>

                <div class="text-center">
                    <span class="text-sm text-gray-600">
                        { i18n.T(ctx, "need_account") }
                    </span>
                    <a href="/signup" class="font-medium text-indigo-600 hover:text-indigo-500">
                        { i18n.T(ctx, "sign_up") }
                    </a>
                </div>
            </form>
        </div>
    </div>
}
```

### Login Page Metadata

```yaml
# app/login/page.templ.yaml
metadata:
  page_title: "Sign In"
  description: "Sign in to your account"

i18n:
  en:
    login_title: "Sign in to your account"
    login_subtitle: "Enter your credentials to access the system"
    email: "Email address"
    password: "Password"
    sign_in: "Sign in"
    need_account: "Don't have an account?"
    sign_up: "Sign up"
  de:
    login_title: "Bei Ihrem Konto anmelden"
    login_subtitle: "Geben Sie Ihre Anmeldedaten ein, um auf das System zuzugreifen"
    email: "E-Mail-Adresse"
    password: "Passwort"
    sign_in: "Anmelden"
    need_account: "Noch kein Konto?"
    sign_up: "Registrieren"

auth:
  type: "Public"
```

### Sign Up Page

```go
// app/signup/page.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ SignUpPage() {
    <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div class="max-w-md w-full space-y-8">
            <div>
                <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                    { i18n.T(ctx, "signup_title") }
                </h2>
                <p class="mt-2 text-center text-sm text-gray-600">
                    { i18n.T(ctx, "signup_subtitle") }
                </p>
            </div>
            <form class="mt-8 space-y-6"
                  hx-post="/api/auth/signup"
                  hx-target="#form-error"
                  hx-swap="innerHTML">
                <div class="space-y-4">
                    <div>
                        <label for="name" class="block text-sm font-medium text-gray-700">
                            { i18n.T(ctx, "name") }
                        </label>
                        <input id="name"
                               name="name"
                               type="text"
                               required
                               class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                               placeholder={ i18n.T(ctx, "name") }/>
                    </div>
                    <div>
                        <label for="email-address" class="block text-sm font-medium text-gray-700">
                            { i18n.T(ctx, "email") }
                        </label>
                        <input id="email-address"
                               name="email"
                               type="email"
                               autocomplete="email"
                               required
                               class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                               placeholder={ i18n.T(ctx, "email") }/>
                    </div>
                    <div>
                        <label for="password" class="block text-sm font-medium text-gray-700">
                            { i18n.T(ctx, "password") }
                        </label>
                        <input id="password"
                               name="password"
                               type="password"
                               required
                               class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                               placeholder={ i18n.T(ctx, "password") }/>
                    </div>
                </div>

                <div id="form-error" class="text-red-600 text-sm"></div>

                <div>
                    <button type="submit"
                            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                        { i18n.T(ctx, "create_account") }
                    </button>
                </div>

                <div class="text-center">
                    <span class="text-sm text-gray-600">
                        { i18n.T(ctx, "have_account") }
                    </span>
                    <a href="/login" class="font-medium text-indigo-600 hover:text-indigo-500">
                        { i18n.T(ctx, "sign_in") }
                    </a>
                </div>
            </form>
        </div>
    </div>
}
```

### Protected Page Example

```go
// app/dashboard/page.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ DashboardPage(user *UserData) {
    <div class="min-h-screen bg-gray-50">
        <header class="bg-white shadow">
            <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                <div class="flex justify-between items-center">
                    <h1 class="text-3xl font-bold text-gray-900">
                        { i18n.T(ctx, "welcome") }, { user.Name }!
                    </h1>
                    <div class="flex space-x-4">
                        <span class="text-sm text-gray-600">
                            { user.Email }
                        </span>
                        <form method="POST" action="/api/auth/signout" class="inline">
                            <button type="submit"
                                    class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium">
                                { i18n.T(ctx, "sign_out") }
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </header>

        <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
            <div class="px-4 py-6 sm:px-0">
                <div class="border-4 border-dashed border-gray-200 rounded-lg h-96 flex items-center justify-center">
                    <div class="text-center">
                        <h2 class="text-2xl font-semibold text-gray-900">
                            { i18n.T(ctx, "dashboard_title") }
                        </h2>
                        <p class="mt-2 text-gray-600">
                            { i18n.T(ctx, "dashboard_subtitle") }
                        </p>
                    </div>
                </div>
            </div>
        </main>
    </div>
}
```

### Dashboard Metadata

```yaml
# app/dashboard/page.templ.yaml
metadata:
  page_title: "Dashboard"
  description: "User dashboard"

i18n:
  en:
    welcome: "Welcome"
    dashboard_title: "Dashboard"
    dashboard_subtitle: "Your personalized dashboard"
    sign_out: "Sign Out"
  de:
    welcome: "Willkommen"
    dashboard_title: "Dashboard"
    dashboard_subtitle: "Ihr persönliches Dashboard"
    sign_out: "Abmelden"

auth:
  type: "UserRequired"
  redirect_url: "/login"

data_services:
  - "UserDataService"
```

## Role-Based Access Control

### Admin-Only Page

```go
// app/admin/page.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ AdminPage(user *UserData) {
    <div class="min-h-screen bg-gray-50">
        <header class="bg-white shadow">
            <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                <h1 class="text-3xl font-bold text-gray-900">
                    { i18n.T(ctx, "admin_panel") }
                </h1>
            </div>
        </header>

        <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
            <div class="bg-white shadow rounded-lg p-6">
                <h2 class="text-xl font-semibold text-gray-900 mb-4">
                    { i18n.T(ctx, "user_management") }
                </h2>
                <p class="text-gray-600">
                    { i18n.T(ctx, "admin_welcome", user.Name) }
                </p>
            </div>
        </main>
    </div>
}
```

### Admin Page Metadata

```yaml
# app/admin/page.templ.yaml
metadata:
  page_title: "Admin Panel"
  description: "Administrative interface"

i18n:
  en:
    admin_panel: "Admin Panel"
    user_management: "User Management"
    admin_welcome: "Welcome to the admin panel, %s"
  de:
    admin_panel: "Admin-Panel"
    user_management: "Benutzerverwaltung"
    admin_welcome: "Willkommen im Admin-Panel, %s"

auth:
  type: "AdminRequired"
  redirect_url: "/login"
  roles: ["admin", "super_admin"]  # Optional: specific roles required

data_services:
  - "UserDataService"
```

## Authentication Middleware Integration

### Custom Middleware

```go
// pkg/middleware/auth.go
package middleware

import (
    "net/http"
    "github.com/denkhaus/templ-router/pkg/interfaces"
)

func CustomAuthMiddleware(authService interfaces.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Custom authentication logic
            user, err := authService.GetCurrentUser(r.Context())
            if err != nil {
                // Handle authentication error
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }

            // Add user to context
            ctx := context.WithValue(r.Context(), "user", user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Integration in Main Application

```go
// main.go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/denkhaus/templ-router/pkg/router/middleware"
    "github.com/go-chi/chi/v5"
)

func main() {
    // Create DI container
    container := di.NewContainer()
    container.RegisterRouterServices("TR")

    // Create router
    mux := chi.NewRouter()

    // Add authentication middleware
    authMiddleware, err := middleware.NewAuthContextMiddleware(container.GetInjector())
    if err != nil {
        panic(err)
    }
    mux.Use(authMiddleware.Middleware)

    // Register file-based routes
    router := container.GetRouter()
    if err := router.Initialize(); err != nil {
        panic(err)
    }

    if err := router.RegisterRoutes(mux); err != nil {
        panic(err)
    }

    // Register authentication API routes
    authHandlers := do.MustInvoke[interfaces.AuthHandlers](container.GetInjector())
    authHandlers.RegisterRoutes(func(method, path string, handler http.HandlerFunc) {
        mux.Post(path, handler)
        mux.Get(path, handler)
    })

    http.ListenAndServe(":8080", mux)
}
```

## Security Features

### Session Security

```yaml
# Production configuration
TR_AUTH_SESSION_COOKIE_SECURE=true      # HTTPS-only cookies
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true   # Prevent XSS attacks
TR_AUTH_SESSION_COOKIE_SAME_SITE=Strict # Prevent CSRF attacks
TR_AUTH_SESSION_EXPIRY=1h               # Shorter session for production
```

### CSRF Protection

```yaml
# Enable CSRF protection
TR_SECURITY_CSRF_SECRET=your-secret-key
TR_SECURITY_ENABLE_CSRF=true
```

### Rate Limiting

```yaml
# Rate limiting configuration
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_SECURITY_RATE_LIMIT_REQUESTS=10
TR_SECURITY_RATE_LIMIT_WINDOW=1m
```

## User Management

### Default Admin User

Create a default admin user automatically:

```yaml
# Development configuration
TR_AUTH_CREATE_DEFAULT_ADMIN=true
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
TR_AUTH_DEFAULT_ADMIN_PASSWORD=admin123
TR_AUTH_DEFAULT_ADMIN_NAME="Admin User"
```

### User Roles

Supported user roles:

- `user`: Regular authenticated user
- `admin`: Administrative user
- `super_admin`: Super administrative user

### Custom User Data

Extend user model with custom fields:

```go
// pkg/models/user.go
package models

type User struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Name     string `json:"name"`
    Role     string `json:"role"`

    // Custom fields
    Avatar   string `json:"avatar"`
    Bio      string `json:"bio"`
    Settings map[string]interface{} `json:"settings"`
}
```

## Internationalization

### Internationalized Login Flow

Create internationalized authentication pages:

```bash
app/
├── login/
│   ├── page.templ
│   └── page.templ.yaml
└── locale_/
    └── login/
        ├── page.templ
        └── page.templ.yaml
```

### Internationalized Redirects

Configure locale-aware redirect URLs:

```yaml
# app/login/page.templ.yaml
i18n:
  en:
    login_title: "Sign In"
    login_subtitle: "Enter your credentials"
  de:
    login_title: "Anmelden"
    login_subtitle: "Geben Sie Ihre Anmeldedaten ein"

auth:
  type: "Public"
  redirect_url: "/login"  # Will be localized to /{locale}/login
```

## Testing Authentication

### Test Authenticated Routes

```go
// pkg/tests/auth_test.go
package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAuthenticatedRoute(t *testing.T) {
    // Create test request with authentication
    req := httptest.NewRequest("GET", "/dashboard", nil)

    // Add authentication cookie/session
    req.AddCookie(&http.Cookie{
        Name:  "session_id",
        Value: "test-session-token",
    })

    w := httptest.NewRecorder()

    // Test the route
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}

func TestUnauthorizedRoute(t *testing.T) {
    // Create test request without authentication
    req := httptest.NewRequest("GET", "/dashboard", nil)
    w := httptest.NewRecorder()

    // Test the route
    router.ServeHTTP(w, req)

    // Should redirect to login
    if w.Code != http.StatusSeeOther {
        t.Errorf("Expected redirect to login, got %d", w.Code)
    }
}
```

## Troubleshooting

### Common Issues

#### Authentication Not Working

```bash
# Check authentication configuration
env | grep TR_AUTH

# Verify middleware is registered
# Check that authMiddleware.Middleware is added to your router

# Check session configuration
# Verify TR_AUTH_SESSION_* variables are set correctly
```

#### Redirect Loops

```bash
# Check redirect URLs
# Ensure public pages don't require authentication
# Verify redirect_url in auth configuration

# Check route precedence
# Ensure login route is accessible to public users
```

#### Session Issues

```bash
# Check session cookie configuration
# Verify TR_AUTH_SESSION_COOKIE_* variables

# Check session storage
# Ensure session store is properly configured
```

#### Permission Denied

```bash
# Check user roles
# Verify user has required role for protected routes

# Check role configuration
# Ensure roles are properly set in auth configuration
```

### Debug Mode

```bash
# Enable debug logging
TR_LOGGING_LEVEL=debug

# Check authentication flow
# Look for authentication-related log messages
```

## Best Practices

### Security

1. **Use HTTPS in production**: Always enable secure cookies
2. **Set appropriate session expiry**: Balance security and user experience
3. **Implement rate limiting**: Prevent brute force attacks
4. **Validate user input**: Sanitize all user inputs
5. **Use CSRF protection**: Prevent cross-site request forgery

### User Experience

1. **Clear error messages**: Provide helpful authentication errors
2. **Remember me functionality**: Allow users to stay logged in
3. **Password strength requirements**: Enforce strong passwords
4. **Account recovery**: Implement password reset functionality
5. **Multi-factor authentication**: Consider 2FA for sensitive applications

### Performance

1. **Session caching**: Cache session data for better performance
2. **Database optimization**: Optimize user queries
3. **Cookie size**: Keep session cookies small
4. **Connection pooling**: Use database connection pooling

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first project
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Understand routing with authentication
- **[Internationalization](INTERNATIONALIZATION.md)** - Add i18n to auth flows
- **[Configuration](CONFIGURATION.md)** - Configure authentication behavior
- **[Data Services](DATA-SERVICES.md)** - Use user data in templates

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Documentation](../README.md)** - Main documentation