# Middleware System

**Complete guide to the middleware system and request processing pipeline in Templ Router.**

## Overview

Templ Router provides a comprehensive middleware system that processes requests through a configurable pipeline. The middleware handles authentication, internationalization, template rendering, error handling, and more in a clean, composable architecture.

**Key Features:**
- Configurable middleware pipeline
- Built-in authentication, i18n, and template middleware
- Custom middleware support
- Hierarchical metadata merging
- Error handling and recovery
- Request context management

## Configuration Prefix Notice

**Important:** Some environment variables in this documentation use the default prefix `TR_`. This prefix is **configurable** when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_AUTH_ENABLE_MIDDLEWARE=true`
- Custom: `MYAPP_AUTH_ENABLE_MIDDLEWARE=true`
- Multiple apps: `APP1_AUTH_ENABLE_MIDDLEWARE=true` and `APP2_AUTH_ENABLE_MIDDLEWARE=true`

All environment variable examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix.

## Request Processing Pipeline

Each request flows through a series of middleware layers:

```
Request → Authentication → I18n → Template Rendering → Response
```

### Pipeline Stages

1. **Authentication Middleware**
   - Validates user authentication status
   - Enforces route-based access control
   - Sets user context

2. **Internationalization Middleware**
   - Detects and validates locale
   - Loads appropriate translations
   - Sets locale context

3. **Template Middleware**
   - Loads template metadata
   - Merges hierarchical configurations
   - Renders templates with data

4. **Error Handling Middleware**
   - Catches and handles errors
   - Renders appropriate error pages
   - Provides graceful degradation

## Built-in Middleware

### Authentication Middleware

Handles user authentication and authorization:

```go
// Enable authentication middleware
authMiddleware, err := middleware.NewAuthContextMiddleware(injector)
if err != nil {
    return err
}
mux.Use(authMiddleware.Middleware)
```

**Features:**
- Session validation
- Route-based access control
- User context injection
- Redirect handling for unauthenticated users

### Internationalization Middleware

Manages locale detection and translation loading:

```go
// Enable i18n middleware
i18nMiddleware, err := middleware.NewI18nMiddleware(injector)
if err != nil {
    return err
}
mux.Use(i18nMiddleware.Middleware)
```

**Features:**
- Automatic locale detection from URL/headers/cookies
- Translation file loading
- Locale context injection
- Fallback locale handling

### Template Middleware

Handles template rendering and metadata management:

```go
// Enable template middleware
templateMiddleware, err := middleware.NewTemplateMiddleware(injector)
if err != nil {
    return err
}
mux.Use(templateMiddleware.Middleware)
```

**Features:**
- Template metadata loading
- Hierarchical configuration merging
- Data service resolution
- Template rendering with context

## Custom Middleware

### Creating Custom Middleware

```go
package middleware

import (
    "context"
    "net/http"
)

// CustomMiddleware adds custom functionality
type CustomMiddleware struct {
    config *CustomConfig
    logger Logger
}

func NewCustomMiddleware(config *CustomConfig, logger Logger) *CustomMiddleware {
    return &CustomMiddleware{
        config: config,
        logger: logger,
    }
}

func (m *CustomMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing logic
        startTime := time.Now()

        // Add custom headers
        w.Header().Set("X-Custom-Header", "value")

        // Log request
        m.logger.Info("Processing request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.String("remote", r.RemoteAddr),
        )

        // Call next middleware in chain
        next.ServeHTTP(w, r)

        // Post-processing logic
        duration := time.Since(startTime)
        m.logger.Info("Request completed",
            zap.Duration("duration", duration),
            zap.Int("status", w.(*responseWriter).Status()),
        )
    })
}
```

### Request/Response Wrapper

For capturing response details:

```go
type responseWriter struct {
    http.ResponseWriter
    status int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
    return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### Custom Middleware Registration

```go
func main() {
    // Create custom middleware
    customMiddleware := middleware.NewCustomMiddleware(config, logger)

    // Register middleware with router
    mux := chi.NewRouter()
    mux.Use(customMiddleware.Middleware)

    // Register other middleware
    authMiddleware, _ := middleware.NewAuthContextMiddleware(injector)
    mux.Use(authMiddleware.Middleware)

    // Register routes
    router.RegisterRoutes(mux)
}
```

## Middleware Configuration

### Environment-Based Configuration

Configure middleware through environment variables:

```bash
# Authentication middleware
TR_AUTH_ENABLE_MIDDLEWARE=true
TR_AUTH_SKIP_PATTERNS=/health,/metrics

# Internationalization middleware
TR_I18N_ENABLE_MIDDLEWARE=true
TR_I18N_SKIP_PATTERNS=/api,/static

# Template middleware
TR_TEMPLATE_ENABLE_MIDDLEWARE=true
TR_TEMPLATE_CACHE_ENABLED=true

# Security middleware
TR_SECURITY_ENABLE_MIDDLEWARE=true
TR_SECURITY_ENABLE_CSRF=true
```

### Conditional Middleware

Enable/disable middleware based on configuration:

```go
func setupMiddleware(mux *chi.Mux, injector *do.Injector) error {
    config := do.MustInvoke[*ConfigService](injector)

    // Authentication middleware
    if config.Auth.Enabled {
        authMiddleware, err := middleware.NewAuthContextMiddleware(injector)
        if err != nil {
            return err
        }
        mux.Use(authMiddleware.Middleware)
    }

    // I18n middleware
    if config.I18n.Enabled {
        i18nMiddleware, err := middleware.NewI18nMiddleware(injector)
        if err != nil {
            return err
        }
        mux.Use(i18nMiddleware.Middleware)
    }

    // Security middleware
    if config.Security.Enabled {
        securityMiddleware, err := middleware.NewSecurityMiddleware(injector)
        if err != nil {
            return err
        }
        mux.Use(securityMiddleware.Middleware)
    }

    return nil
}
```

## Middleware Ordering

The order of middleware registration is important:

```go
func main() {
    mux := chi.NewRouter()

    // 1. Security middleware (first)
    securityMiddleware, _ := middleware.NewSecurityMiddleware(injector)
    mux.Use(securityMiddleware.Middleware)

    // 2. Logging middleware
    loggingMiddleware, _ := middleware.NewLoggingMiddleware(injector)
    mux.Use(loggingMiddleware.Middleware)

    // 3. Authentication middleware
    authMiddleware, _ := middleware.NewAuthContextMiddleware(injector)
    mux.Use(authMiddleware.Middleware)

    // 4. Internationalization middleware
    i18nMiddleware, _ := middleware.NewI18nMiddleware(injector)
    mux.Use(i18nMiddleware.Middleware)

    // 5. Template middleware (last)
    templateMiddleware, _ := middleware.NewTemplateMiddleware(injector)
    mux.Use(templateMiddleware.Middleware)

    // Register routes after all middleware
    router.RegisterRoutes(mux)
}
```

## Advanced Middleware Patterns

### Context Enhancement

Add custom data to request context:

```go
type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserAgentKey  contextKey = "user_agent"
)

func ContextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Generate request ID
        requestID := uuid.New().String()

        // Create enhanced context
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        ctx = context.WithValue(ctx, UserAgentKey, r.UserAgent())

        // Continue with enhanced context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Helper functions to access context data
func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(RequestIDKey).(string); ok {
        return id
    }
    return ""
}
```

### Request Metrics

Collect request metrics:

```go
type MetricsMiddleware struct {
    requestCount    prometheus.Counter
    requestDuration prometheus.Histogram
    activeRequests  prometheus.Gauge
}

func (m *MetricsMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Increment active requests
        m.activeRequests.Inc()
        defer m.activeRequests.Dec()

        // Record start time
        start := time.Now()

        // Increment request count
        m.requestCount.Inc()

        // Create response writer to capture status
        rw := NewResponseWriter(w)

        // Process request
        next.ServeHTTP(rw, r)

        // Record metrics
        duration := time.Since(start)
        m.requestDuration.Observe(duration.Seconds())

        // Record status code
        labels := prometheus.Labels{
            "method": r.Method,
            "path":   r.URL.Path,
            "status": fmt.Sprintf("%d", rw.status),
        }
        m.requestCount.With(labels).Inc()
    })
}
```

### Rate Limiting

Implement rate limiting per user/IP:

```go
type RateLimitMiddleware struct {
    limiter *rate.Limiter
    requests map[string]*rate.Limiter
    mutex    sync.RWMutex
}

func (m *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get client identifier (IP or user ID)
        clientID := m.getClientIdentifier(r)

        // Get or create limiter for client
        limiter := m.getLimiter(clientID)

        // Check rate limit
        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        // Process request
        next.ServeHTTP(w, r)
    })
}

func (m *RateLimitMiddleware) getClientIdentifier(r *http.Request) string {
    // Try to get user ID from context
    if userID := getUserID(r.Context()); userID != "" {
        return userID
    }

    // Fall back to IP address
    return r.RemoteAddr
}

func (m *RateLimitMiddleware) getLimiter(clientID string) *rate.Limiter {
    m.mutex.RLock()
    limiter, exists := m.requests[clientID]
    m.mutex.RUnlock()

    if !exists {
        m.mutex.Lock()
        defer m.mutex.Unlock()

        // Double-check after acquiring write lock
        if limiter, exists := m.requests[clientID]; exists {
            return limiter
        }

        // Create new limiter
        limiter = rate.NewLimiter(rate.Limit(10), 20) // 10 requests per second, burst of 20
        m.requests[clientID] = limiter
    }

    return limiter
}
```

## Error Handling Middleware

### Centralized Error Handling

```go
type ErrorHandlingMiddleware struct {
    logger Logger
}

func (m *ErrorHandlingMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                m.logger.Error("Panic recovered",
                    zap.Any("error", err),
                    zap.String("request", r.URL.String()),
                    zap.Stack("stack"),
                )

                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

### Custom Error Pages

```go
type ErrorPageMiddleware struct {
    errorTemplate *template.Template
    logger       Logger
}

func (m *ErrorPageMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                m.logger.Error("Application error",
                    zap.Any("error", err),
                    zap.String("path", r.URL.Path),
                )

                m.renderErrorPage(w, http.StatusInternalServerError, err)
            }
        }()

        next.ServeHTTP(w, r)
    })
}

func (m *ErrorPageMiddleware) renderErrorPage(w http.ResponseWriter, statusCode int, err error) {
    w.WriteHeader(statusCode)

    data := map[string]interface{}{
        "StatusCode": statusCode,
        "Error":      err.Error(),
        "Timestamp": time.Now(),
    }

    if err := m.errorTemplate.Execute(w, data); err != nil {
        m.logger.Error("Failed to render error page",
            zap.Error(err),
            zap.Int("status_code", statusCode),
        )

        // Fallback to simple error message
        fmt.Fprintf(w, "Error %d: %s", statusCode, err.Error())
    }
}
```

## Testing Middleware

### Unit Testing Middleware

```go
func TestCustomMiddleware(t *testing.T) {
    // Create test handler
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Create middleware
    middleware := NewCustomMiddleware(config, logger)
    handler := middleware.Middleware(testHandler)

    // Create test request
    req := httptest.NewRequest("GET", "/test", nil)
    rr := httptest.NewRecorder()

    // Execute middleware
    handler.ServeHTTP(rr, req)

    // Assertions
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "OK", rr.Body.String())
    assert.Equal(t, "value", rr.Header().Get("X-Custom-Header"))
}
```

### Integration Testing Middleware

```go
func TestMiddlewareChain(t *testing.T) {
    // Create middleware chain
    mux := chi.NewRouter()
    mux.Use(loggingMiddleware.Middleware)
    mux.Use(authMiddleware.Middleware)
    mux.Use(templateMiddleware.Middleware)

    // Register test route
    mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Test OK"))
    })

    // Create test server
    server := httptest.NewServer(mux)
    defer server.Close()

    // Test authenticated request
    req := httptest.NewRequest("GET", server.URL+"/test", nil)
    req.Header.Set("Authorization", "Bearer valid-token")

    resp, err := http.DefaultClient.Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

## Best Practices

### Middleware Design

1. **Keep middleware focused** on a single responsibility
2. **Make middleware composable** and reusable
3. **Handle errors gracefully** and provide meaningful responses
4. **Log important events** for debugging and monitoring
5. **Test middleware independently** of the application

### Performance Considerations

1. **Avoid expensive operations** in middleware hot paths
2. **Use connection pooling** for external service calls
3. **Implement timeouts** for external dependencies
4. **Cache frequently accessed data**
5. **Monitor middleware performance**

### Security

1. **Validate input** and sanitize data
2. **Implement rate limiting** to prevent abuse
3. **Use HTTPS** and secure cookies
4. **Set appropriate security headers**
5. **Log security events** for monitoring

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up basic middleware
- **[Authentication](AUTHENTICATION.md)** - Authentication middleware
- **[Internationalization](INTERNATIONALIZATION.md)** - I18n middleware
- **[Dependency Injection](DEPENDENCY-INJECTION.md)** - DI with middleware
- **[Configuration](CONFIGURATION.md)** - Configure middleware

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Chi Router Documentation](https://github.com/go-chi/chi)** - HTTP router documentation