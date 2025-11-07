# Authentication & Authorization

**Complete guide to implementing authentication and authorization in Templ Router.**

## Overview

Templ Router provides session-based authentication with role-based access control and built-in API endpoints. The system integrates seamlessly with file-based routing through `.templ.yaml` metadata files.

**Key Features:**
- Three authentication types: `Public`, `UserRequired`, `AdminRequired`
- Session-based authentication with configurable expiry
- Built-in authentication API endpoints
- Role-based access control
- Configurable redirect routes
- Internationalized authentication flows

## Configuration Prefix

All environment variables use a configurable prefix (default `TR_`):

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Uses MYAPP_ prefix
```

## Authentication Types

### 1. Public Access

No authentication required (default):

```yaml
auth:
  type: "Public"
```

Or omit the auth section entirely.

### 2. User Required

Any authenticated user can access:

```yaml
auth:
  type: "UserRequired"
  redirect_url: "/login"
```

### 3. Admin Required

Only admin users can access:

```yaml
auth:
  type: "AdminRequired"
  redirect_url: "/login"
  roles: ["admin", "super_admin"]  # Optional: specific roles
```

## Configuration

### Environment Variables

```bash
# Session Configuration
TR_AUTH_SESSION_EXPIRY=24h                # Session duration
TR_AUTH_SESSION_COOKIE_NAME=session_id     # Cookie name
TR_AUTH_SESSION_COOKIE_SECURE=true         # HTTPS-only cookies
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true      # HTTP-only cookies

# Redirect Configuration
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard    # After login
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/welcome      # After registration
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/           # After logout

# Default Admin User (development)
TR_AUTH_CREATE_DEFAULT_ADMIN=true
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
TR_AUTH_DEFAULT_ADMIN_PASSWORD=admin123
```

### Internationalized Redirects

Support for locale-aware redirects:

```bash
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/{locale}/dashboard
# → /en/dashboard or /de/dashboard
```

## Built-in API Endpoints

### Available Endpoints

```bash
POST /api/auth/signin      # User sign in
POST /api/auth/signout     # User sign out
POST /api/auth/signup      # User registration
GET  /api/auth/me          # Get current user info
```

### Sign In

```bash
POST /api/auth/signin
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

### Sign Up

```bash
POST /api/auth/signup
{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "Jane Doe"
}
```

### Current User

```bash
GET /api/auth/me
```

**Response:**
```json
{
  "user": {
    "id": "123",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user"
  }
}
```

## Template Implementation

### Login Page

```go
// app/login/page.templ - MUST be named "Page"
templ Page() {
    <div class="min-h-screen flex items-center justify-center bg-gray-50">
        <div class="max-w-md w-full space-y-8">
            <h2 class="text-center text-3xl font-extrabold">
                { i18n.T(ctx, "login_title") }
            </h2>
            <form hx-post="/api/auth/signin" hx-target="#form-error">
                <div class="space-y-4">
                    <input name="email" type="email"
                           placeholder={ i18n.T(ctx, "email") } required/>
                    <input name="password" type="password"
                           placeholder={ i18n.T(ctx, "password") } required/>
                </div>
                <div id="form-error" class="text-red-600 text-sm"></div>
                <button type="submit" class="w-full bg-indigo-600 text-white py-2 rounded">
                    { i18n.T(ctx, "sign_in") }
                </button>
            </form>
        </div>
    </div>
}
```

### Login Page Metadata

```yaml
# app/login/page.templ.yaml
auth:
  type: "Public"

i18n:
  en:
    login_title: "Sign in to your account"
    email: "Email address"
    password: "Password"
    sign_in: "Sign in"
  de:
    login_title: "Bei Ihrem Konto anmelden"
    email: "E-Mail-Adresse"
    password: "Passwort"
    sign_in: "Anmelden"
```

### Protected Page Example

```go
// app/dashboard/page.templ - MUST be named "Page"
templ Page(user *UserData) {
    <div class="min-h-screen bg-gray-50">
        <header class="bg-white shadow">
            <div class="max-w-7xl mx-auto py-6 px-4">
                <div class="flex justify-between items-center">
                    <h1 class="text-3xl font-bold">
                        { i18n.T(ctx, "welcome") }, { user.Name }!
                    </h1>
                    <form method="POST" action="/api/auth/signout" class="inline">
                        <button type="submit" class="bg-red-600 text-white px-4 py-2 rounded">
                            { i18n.T(ctx, "sign_out") }
                        </button>
                    </form>
                </div>
            </div>
        </header>
        <main class="max-w-7xl mx-auto py-6">
            <h2>{ i18n.T(ctx, "dashboard_title") }</h2>
        </main>
    </div>
}
```

### Dashboard Metadata

```yaml
# app/dashboard/page.templ.yaml
auth:
  type: "UserRequired"
  redirect_url: "/login"

i18n:
  en:
    welcome: "Welcome"
    dashboard_title: "Dashboard"
    sign_out: "Sign Out"
  de:
    welcome: "Willkommen"
    dashboard_title: "Dashboard"
    sign_out: "Abmelden"
```

### Admin Page

```go
// app/admin/page.templ
templ Page(user *UserData) {
    <div class="min-h-screen bg-gray-50">
        <header class="bg-white shadow">
            <h1 class="text-3xl font-bold">
                { i18n.T(ctx, "admin_panel") }
            </h1>
        </header>
        <main class="max-w-7xl mx-auto py-6">
            <h2>{ i18n.T(ctx, "user_management") }</h2>
        </main>
    </div>
}
```

### Admin Page Metadata

```yaml
# app/admin/page.templ.yaml
auth:
  type: "AdminRequired"
  redirect_url: "/login"
  roles: ["admin", "super_admin"]

i18n:
  en:
    admin_panel: "Admin Panel"
    user_management: "User Management"
  de:
    admin_panel: "Admin-Panel"
    user_management: "Benutzerverwaltung"
```

## Integration

### Main Application Setup

```go
// main.go
func main() {
    container := di.NewContainer()
    container.RegisterRouterServices("TR")

    mux := chi.NewRouter()

    // Add authentication middleware
    authMiddleware, err := middleware.NewAuthContextMiddleware(container.GetInjector())
    if err != nil {
        panic(err)
    }
    mux.Use(authMiddleware.Middleware)

    // Register routes
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
TR_SECURITY_CSRF_SECRET=your-secret-key
TR_SECURITY_ENABLE_CSRF=true
```

### Rate Limiting

```yaml
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_SECURITY_RATE_LIMIT_REQUESTS=10
TR_SECURITY_RATE_LIMIT_WINDOW=1m
```

## User Management

### Default Admin User

```yaml
# Development configuration
TR_AUTH_CREATE_DEFAULT_ADMIN=true
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
TR_AUTH_DEFAULT_ADMIN_PASSWORD=admin123
```

### User Roles

- `user`: Regular authenticated user
- `admin`: Administrative user
- `super_admin`: Super administrative user

## Testing

### Test Authentication

```go
func TestAuthenticatedRoute(t *testing.T) {
    req := httptest.NewRequest("GET", "/dashboard", nil)
    req.AddCookie(&http.Cookie{
        Name:  "session_id",
        Value: "test-session-token",
    })

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Troubleshooting

### Common Issues

**Authentication not working:**
```bash
env | grep TR_AUTH
# Verify middleware registration
# Check session configuration
```

**Redirect loops:**
```bash
# Ensure public pages don't require authentication
# Check redirect_url configuration
# Verify login route accessibility
```

**Session issues:**
```bash
# Check session cookie configuration
# Verify session store setup
```

**Permission denied:**
```bash
# Check user roles
# Verify role configuration
```

## Best Practices

### Security
- Use HTTPS in production
- Set appropriate session expiry
- Implement rate limiting
- Use CSRF protection
- Validate user input

### User Experience
- Clear error messages
- Password strength requirements
- Account recovery functionality
- Consider multi-factor authentication

### Performance
- Session caching
- Database optimization
- Small session cookies
- Connection pooling

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first project
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Routing with authentication
- **[Internationalization](INTERNATIONALIZATION.md)** - Add i18n to auth flows
- **[Data Services](DATA-SERVICES.md)** - Use user data in templates