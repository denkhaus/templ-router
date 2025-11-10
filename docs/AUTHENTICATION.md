# Authentication & Authorization

**Complete guide to implementing hook-based authentication and authorization in Templ Router.**

## Overview

Templ Router provides a **hook-based authentication system** where the router handles route protection and access control, while client applications implement their own authentication logic through interfaces. This approach ensures the library remains generic and works with any authentication system (OAuth, JWT, LDAP, custom session stores, etc.).

**Key Features:**
- Hook-based authentication via `AuthValidator` interface
- Three authentication types: `Public`, `UserRequired`, `AdminRequired`
- Route protection middleware that calls client-provided authentication hooks
- Role-based access control
- Configurable redirect routes for authentication failures
- Internationalized authentication flows
- Complete separation of authentication logic from router core

## Hook-Based Authentication Architecture

### Router Responsibilities
- **Route Protection**: Middleware that checks authentication based on template metadata
- **Access Control**: Validates user roles against route requirements
- **Failure Handling**: Redirects unauthenticated users to login routes
- **Configuration**: Manages redirect URLs for authentication failures

### Client Application Responsibilities
- **Authentication Logic**: Implement `AuthValidator` interface with your auth system
- **Session Management**: Handle cookies, tokens, or other session mechanisms
- **User Stores**: Implement user lookup and validation (database, LDAP, OAuth, etc.)
- **API Endpoints**: Provide signin, signup, signout endpoints
- **User Management**: Handle user creation, validation, and role management

### AuthValidator Interface

```go
type AuthValidator interface {
    // IsAuthenticated checks if the current request is from an authenticated user
    IsAuthenticated(req *http.Request) bool

    // GetCurrentUser returns the authenticated user for the current request
    GetCurrentUser(req *http.Request) (UserEntity, error)

    // HasRole checks if the given user has any of the required roles
    HasRole(user UserEntity, requiredRoles []string) bool
}
```

### Implementation Example

```go
// Your custom authentication validator
type MyAuthValidator struct {
    userStore    MyUserStore
    sessionStore MySessionStore
    jwtService   MyJWTService
}

func (av *MyAuthValidator) IsAuthenticated(req *http.Request) bool {
    // Check JWT token, session cookie, or other auth mechanism
    token := av.extractToken(req)
    return av.jwtService.ValidateToken(token)
}

func (av *MyAuthValidator) GetCurrentUser(req *http.Request) (UserEntity, error) {
    // Extract user from token, session, or other mechanism
    token := av.extractToken(req)
    userID, err := av.jwtService.GetUserIDFromToken(token)
    if err != nil {
        return nil, err
    }
    return av.userStore.GetUserByID(userID)
}

func (av *MyAuthValidator) HasRole(user UserEntity, requiredRoles []string) bool {
    // Check if user has required roles
    userRoles := user.GetRoles()
    for _, required := range requiredRoles {
        for _, userRole := range userRoles {
            if userRole == required {
                return true
            }
        }
    }
    return false
}
```

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

**Router-Level Configuration** (handled by templ-router):

```bash
# Redirect Routes for authentication failures
TR_AUTH_SIGNIN_ROUTE=/login                # Where to redirect unauthenticated users
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard    # Where to redirect after successful login
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/welcome      # Where to redirect after successful registration
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/           # Where to redirect after successful logout

# Route Configuration
TR_ROUTER_ENABLE_AUTH_ROUTES=true          # Enable/disable auth routes
TR_ROUTER_AUTH_ROUTE_PREFIX=/api           # Prefix for auth API endpoints
```

**Client Application Configuration** (handled by your application):

```bash
# Session Management - implement in your application
# Example: MYAPP_SESSION_EXPIRY=24h
# Example: MYAPP_SESSION_COOKIE_NAME=session_id
# Example: MYAPP_JWT_SECRET=your-secret-key

# User Management - implement in your application
# Example: MYAPP_DEFAULT_ADMIN_EMAIL=admin@example.com
# Example: MYAPP_PASSWORD_MIN_LENGTH=8

# Email Configuration - implement in your application
# Example: MYAPP_SMTP_HOST=smtp.example.com
# Example: MYAPP_SMTP_PORT=587
```

**Note**: Most authentication configuration has been moved to client applications. Only redirect routes remain in the router to handle authentication failures.

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

### Main Application Setup with Hook-Based Authentication

```go
// main.go
func main() {
    container := di.NewContainer()

    // Register router services
    injector := container.RegisterRouterServices(context.Background(), "TR")

    // Register your custom authentication components
    container.RegisterApplicationServices(
        // Your custom AuthValidator implementation
        di.WithAuthValidatorFactory(func(i do.Injector) (interfaces.AuthValidator, error) {
            return NewMyAuthValidator(i) // Your implementation
        }),

        // Your custom AuthHandlers for API endpoints
        di.WithAuthHandlersFactory(func(i do.Injector) (interfaces.AuthHandlers, error) {
            return NewMyAuthHandlers(i) // Your implementation
        }),

        // Your custom SessionStore
        di.WithSessionStoreFactory(func(i do.Injector) (interfaces.SessionStore, error) {
            return NewMySessionStore(i) // Your implementation
        }),

        // Your custom UserStore
        di.WithUserStoreFactory(func(i do.Injector) (interfaces.UserStore, error) {
            return NewMyUserStore(i) // Your implementation
        }),
    )

    // Bootstrap router - auth middleware automatically uses your AuthValidator
    routerBootstrap := container.GetRouterBootstrap()
    mux, err := routerBootstrap.Bootstrap()
    if err != nil {
        panic(err)
    }

    http.ListenAndServe(":8080", mux)
}
```

### Custom AuthValidator Example

```go
// auth_validator.go
type MyAuthValidator struct {
    userStore    MyUserStore
    sessionStore MySessionStore
    jwtSecret    string
}

func NewMyAuthValidator(i do.Injector) (interfaces.AuthValidator, error) {
    userStore := do.MustInvoke[interfaces.UserStore](i)
    sessionStore := do.MustInvoke[interfaces.SessionStore](i)

    return &MyAuthValidator{
        userStore:    userStore.(MyUserStore),
        sessionStore: sessionStore.(MySessionStore),
        jwtSecret:    os.Getenv("JWT_SECRET"),
    }, nil
}

func (av *MyAuthValidator) IsAuthenticated(req *http.Request) bool {
    // Extract JWT from Authorization header or cookie
    token := av.extractJWTFromRequest(req)
    if token == "" {
        return false
    }

    // Validate JWT token
    claims, err := av.validateJWT(token)
    if err != nil {
        return false
    }

    // Check if user still exists and is active
    _, err = av.userStore.GetUserByID(claims.UserID)
    return err == nil
}

func (av *MyAuthValidator) GetCurrentUser(req *http.Request) (interfaces.UserEntity, error) {
    token := av.extractJWTFromRequest(req)
    if token == "" {
        return nil, fmt.Errorf("no authentication token found")
    }

    claims, err := av.validateJWT(token)
    if err != nil {
        return nil, err
    }

    return av.userStore.GetUserByID(claims.UserID)
}

func (av *MyAuthValidator) HasRole(user interfaces.UserEntity, requiredRoles []string) bool {
    userRoles := user.GetRoles()
    for _, required := range requiredRoles {
        for _, userRole := range userRoles {
            if userRole == required {
                return true
            }
        }
    }
    return false
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