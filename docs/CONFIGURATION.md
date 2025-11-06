# Configuration

**Complete reference for configuring Templ Router applications.**

## Overview

Templ Router uses a comprehensive configuration system with environment variables that support configurable prefixes. All settings are organized into logical sections for server configuration, authentication, internationalization, routing, middleware, and more.

**Key Features:**
- Environment variable-based configuration
- Configurable prefix (default: "TR")
- Logical grouping of related settings
- Default values for all options
- Runtime configuration validation
- Support for multiple deployment environments

## Configuration Prefix

All environment variables use a configurable prefix to avoid conflicts with other applications:

```go
// In your application setup
container.RegisterRouterServices("TR")  // Default prefix
// or
container.RegisterRouterServices("MYAPP")  // Custom prefix
```

### Prefix Examples

```bash
# Default prefix (TR)
TR_SERVER_HOST=localhost
TR_AUTH_SESSION_EXPIRY=24h

# Custom prefix (MYAPP)
MYAPP_SERVER_HOST=localhost
MYAPP_AUTH_SESSION_EXPIRY=24h

# Multiple applications
TR_SERVER_HOST=localhost        # Templ Router
APP_SERVER_HOST=localhost       # Another application
```

## Server Configuration

Basic HTTP server settings:

```bash
# Server Connection
TR_SERVER_HOST=localhost                    # Server bind address
TR_SERVER_PORT=8080                        # Server port
TR_SERVER_BASE_URL=http://localhost:8080   # Public base URL

# Timeouts
TR_SERVER_READ_TIMEOUT=30s                 # Maximum time to read request
TR_SERVER_WRITE_TIMEOUT=30s                # Maximum time to write response
TR_SERVER_IDLE_TIMEOUT=120s                # Maximum idle time for connections
TR_SERVER_SHUTDOWN_TIMEOUT=30s             # Maximum time for graceful shutdown

# Connection Limits
TR_SERVER_MAX_HEADER_BYTES=1048576         # Maximum request header size (1MB)
TR_SERVER_MAX_CONNECTIONS=1000             # Maximum concurrent connections
TR_SERVER_KEEP_ALIVE_TIMEOUT=30s           # Keep-alive timeout
```

### Example Configurations

```bash
# Development Environment
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080
TR_SERVER_BASE_URL=http://localhost:8080
TR_SERVER_READ_TIMEOUT=30s
TR_SERVER_WRITE_TIMEOUT=30s

# Production Environment
TR_SERVER_HOST=0.0.0.0
TR_SERVER_PORT=80
TR_SERVER_BASE_URL=https://example.com
TR_SERVER_READ_TIMEOUT=15s
TR_SERVER_WRITE_TIMEOUT=15s
TR_SERVER_IDLE_TIMEOUT=60s
TR_SERVER_SHUTDOWN_TIMEOUT=10s
```

## Authentication Configuration

Session and authentication settings:

```bash
# Session Management
TR_AUTH_SESSION_EXPIRY=24h                      # Session duration
TR_AUTH_SESSION_COOKIE_NAME=session_id         # Session cookie name
TR_AUTH_SESSION_COOKIE_SECURE=true              # HTTPS-only cookies
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true           # HTTP-only cookies (prevent XSS)
TR_AUTH_SESSION_COOKIE_SAME_SITE=Strict         # CSRF protection
TR_AUTH_SESSION_COOKIE_PATH=/                   # Cookie path
TR_AUTH_SESSION_COOKIE_DOMAIN=""                # Cookie domain (empty = automatic)

# Redirect URLs
TR_AUTH_SIGNIN_ROUTE=/login                     # Custom login route
TR_AUTH_SIGNUP_ROUTE=/signup                    # Custom signup route
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/dashboard         # After successful login
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/welcome           # After successful registration
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/                # After logout

# Internationalized Redirects
TR_AUTH_SIGNIN_SUCCESS_ROUTE=/{locale}/dashboard    # Locale-aware redirect
TR_AUTH_SIGNUP_SUCCESS_ROUTE=/{locale}/welcome      # Locale-aware redirect
TR_AUTH_SIGNOUT_SUCCESS_ROUTE=/{locale}/            # Locale-aware redirect

# User Management
TR_AUTH_CREATE_DEFAULT_ADMIN=false              # Create default admin user
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com    # Default admin email
TR_AUTH_DEFAULT_ADMIN_PASSWORD=admin123         # Default admin password
TR_AUTH_DEFAULT_ADMIN_NAME="Admin User"         # Default admin name

# Password Requirements
TR_AUTH_PASSWORD_MIN_LENGTH=8                   # Minimum password length
TR_AUTH_PASSWORD_REQUIRE_UPPERCASE=true         # Require uppercase letters
TR_AUTH_PASSWORD_REQUIRE_LOWERCASE=true         # Require lowercase letters
TR_AUTH_PASSWORD_REQUIRE_NUMBERS=true           # Require numbers
TR_AUTH_PASSWORD_REQUIRE_SYMBOLS=false          # Require special characters
```

### Authentication Examples

```bash
# Development Configuration
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_SECURE=false
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true
TR_AUTH_CREATE_DEFAULT_ADMIN=true
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
TR_AUTH_DEFAULT_ADMIN_PASSWORD=admin123

# Production Configuration
TR_AUTH_SESSION_EXPIRY=1h
TR_AUTH_SESSION_COOKIE_SECURE=true
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true
TR_AUTH_SESSION_COOKIE_SAME_SITE=Strict
TR_AUTH_CREATE_DEFAULT_ADMIN=false
TR_AUTH_PASSWORD_MIN_LENGTH=12
TR_AUTH_PASSWORD_REQUIRE_SYMBOLS=true
```

## Internationalization Configuration

Multi-language support settings:

```bash
# Supported Languages
TR_I18N_SUPPORTED_LOCALES=en,de,fr,it,es    # Comma-separated locale codes
TR_I18N_DEFAULT_LOCALE=en                   # Default fallback locale
TR_I18N_FALLBACK_LOCALE=en                  # Ultimate fallback if default fails

# Locale Detection Methods
TR_I18N_DETECTION_METHOD=url,cookie,header   # Priority order of detection
TR_I18N_COOKIE_NAME=locale                  # Cookie name for storing locale
TR_I18N_COOKIE_EXPIRY=8760h                 # Cookie expiry (1 year)

# URL Configuration
TR_I18N_URL_PREFIX=true                     # Add locale prefix to URLs
TR_I18N_REDIRECT_ROOT=true                  # Redirect / to default locale
TR_I18N_STRICT_ROUTING=false                # 404 for unsupported locales
TR_I18N_HIDE_DEFAULT_LOCALE=false           # Hide default locale in URLs

# Translation Loading
TR_I18N_CACHE_TRANSLATIONS=true             # Cache translation files
TR_I18N_RELOAD_ON_CHANGE=false              # Auto-reload in development
TR_I18N_MISSING_KEY_BEHAVIOR=log             # How to handle missing keys
```

### I18n Examples

```bash
# Basic English/German Support
TR_I18N_SUPPORTED_LOCALES=en,de
TR_I18N_DEFAULT_LOCALE=en
TR_I18N_DETECTION_METHOD=url
TR_I18N_URL_PREFIX=true

# Full European Support
TR_I18N_SUPPORTED_LOCALES=en-US,en-GB,de-DE,fr-FR,it-IT,es-ES
TR_I18N_DEFAULT_LOCALE=en-US
TR_I18N_DETECTION_METHOD=url,cookie,header
TR_I18N_COOKIE_NAME=user_locale
TR_I18N_REDIRECT_ROOT=true
```

## Router Configuration

File-based routing and middleware settings:

```bash
# URL Normalization
TR_ROUTER_ENABLE_TRAILING_SLASH=true         # Redirect /path/ to /path
TR_ROUTER_ENABLE_SLASH_REDIRECT=true         # Clean double slashes
TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED=true     # 405 handler for wrong methods
TR_ROUTER_ENABLE_AUTO_OPTIONS=true           # Respond to OPTIONS requests
TR_ROUTER_STRICT_SLASH=false                 # Strict slash handling

# Route Discovery
TR_ROUTER_SCAN_PATH=app                      # Template scan path
TR_ROUTER_IGNORE_HIDDEN=true                 # Ignore hidden files/folders
TR_ROUTER_IGNORE_PATTERNS=*.tmp,*.bak        # File patterns to ignore

# Route Generation
TR_ROUTER_GENERATE_ROUTES=true               # Auto-generate routes
TR_ROUTER_CASE_SENSITIVE=true                # Case-sensitive routing
TR_ROUTER_REMOVE_ACCENTS=false               # Remove accents from URLs

# Error Handling
TR_ROUTER_ENABLE_ERROR_ROUTES=true           # Generate error routes
TR_ROUTER_DEFAULT_ERROR_TEMPLATE=error       # Default error template
TR_ROUTER_404_TEMPLATE=404                   # Custom 404 template
TR_ROUTER_500_TEMPLATE=500                   # Custom 500 template
```

### Routing Examples

```bash
# Development Configuration
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
TR_ROUTER_SCAN_PATH=app
TR_ROUTER_CASE_SENSITIVE=false

# Production Configuration
TR_ROUTER_ENABLE_TRAILING_SLASH=true
TR_ROUTER_ENABLE_SLASH_REDIRECT=true
TR_ROUTER_ENABLE_METHOD_NOT_ALLOWED=true
TR_ROUTER_CASE_SENSITIVE=true
TR_ROUTER_STRICT_SLASH=true
```

## Middleware Configuration

Enable/disable various middleware components:

```bash
# Authentication Middleware
TR_AUTH_ENABLE_MIDDLEWARE=true               # Enable auth middleware
TR_AUTH_REQUIRED_HEADER=X-Auth-Required       # Custom auth header
TR_AUTH_SKIP_PATTERNS=/health,/metrics       # Skip auth for these paths

# Internationalization Middleware
TR_I18N_ENABLE_MIDDLEWARE=true                # Enable i18n middleware
TR_I18N_SKIP_PATTERNS=/api,/static           # Skip i18n for these paths
TR_I18N_HEADER_NAME=Accept-Language           # Custom locale header

# Template Middleware
TR_TEMPLATE_ENABLE_MIDDLEWARE=true            # Enable template middleware
TR_TEMPLATE_CACHE_ENABLED=true               # Enable template caching
TR_TEMPLATE_CACHE_SIZE=100                   # Template cache size

# Security Middleware
TR_SECURITY_ENABLE_MIDDLEWARE=true            # Enable security middleware
TR_SECURITY_ENABLE_CORS=true                  # Enable CORS handling
TR_SECURITY_ENABLE_CSRF=true                  # Enable CSRF protection

# Logging Middleware
TR_LOGGING_ENABLE_MIDDLEWARE=true             # Enable logging middleware
TR_LOGGING_REQUEST_ID_HEADER=X-Request-ID     # Request ID header
TR_LOGGING_SKIP_PATHS=/health,/favicon.ico    # Skip logging for these paths
```

### Middleware Examples

```bash
# Full Middleware Stack
TR_AUTH_ENABLE_MIDDLEWARE=true
TR_I18N_ENABLE_MIDDLEWARE=true
TR_TEMPLATE_ENABLE_MIDDLEWARE=true
TR_SECURITY_ENABLE_MIDDLEWARE=true
TR_LOGGING_ENABLE_MIDDLEWARE=true

# Minimal Configuration
TR_AUTH_ENABLE_MIDDLEWARE=false
TR_I18N_ENABLE_MIDDLEWARE=false
TR_TEMPLATE_ENABLE_MIDDLEWARE=true
TR_SECURITY_ENABLE_MIDDLEWARE=false
TR_LOGGING_ENABLE_MIDDLEWARE=true
```

## Template Configuration

Template rendering and caching settings:

```bash
# Template Generation
TR_TEMPLATE_GENERATOR_OUTPUT_DIR=generated/templates  # Output directory
TR_TEMPLATE_GENERATOR_PACKAGE_NAME=templates         # Package name
TR_TEMPLATE_GENERATOR_WATCH_EXTENSIONS=.templ,.yaml  # Watch extensions

# Template Loading
TR_TEMPLATE_ROOT_DIRECTORY=app                       # Template root directory
TR_TEMPLATE_ENABLE_INHERITANCE=true                  # Enable layout inheritance
TR_TEMPLATE_TEMPLATE_EXTENSION=.templ                # Template file extension

# Template Caching
TR_TEMPLATE_CACHE_ENABLED=true                        # Enable template caching
TR_TEMPLATE_CACHE_SIZE=100                            # Cache size (number of templates)
TR_TEMPLATE_CACHE_TTL=1h                              # Cache TTL
TR_TEMPLATE_CACHE_DIR=/tmp/template_cache             # Cache directory

# Template Rendering
TR_TEMPLATE_ENABLE_COMPRESSION=true                   # Enable response compression
TR_TEMPLATE_COMPRESSION_LEVEL=6                       # Compression level (1-9)
TR_TEMPLATE_CONTENT_TYPE=text/html                    # Default content type
```

### Template Examples

```bash
# Development Configuration
TR_TEMPLATE_ROOT_DIRECTORY=app
TR_TEMPLATE_ENABLE_INHERITANCE=true
TR_TEMPLATE_CACHE_ENABLED=false
TR_TEMPLATE_GENERATOR_WATCH_EXTENSIONS=.templ,.yaml,.yml

# Production Configuration
TR_TEMPLATE_ROOT_DIRECTORY=app
TR_TEMPLATE_ENABLE_INHERITANCE=true
TR_TEMPLATE_CACHE_ENABLED=true
TR_TEMPLATE_CACHE_SIZE=500
TR_TEMPLATE_CACHE_TTL=24h
TR_TEMPLATE_ENABLE_COMPRESSION=true
```

## Security Configuration

Security and protection settings:

```bash
# CSRF Protection
TR_SECURITY_CSRF_SECRET=change-me-in-production       # CSRF secret key
TR_SECURITY_CSRF_SECURE=true                          # HTTPS-only CSRF cookies
TR_SECURITY_CSRF_COOKIE_NAME=_csrf                    # CSRF cookie name
TR_SECURITY_CSRF_HEADER_NAME=X-CSRF-Token             # CSRF header name
TR_SECURITY_CSRF_EXPIRY=24h                           # CSRF token expiry

# Rate Limiting
TR_SECURITY_ENABLE_RATE_LIMIT=true                    # Enable rate limiting
TR_SECURITY_RATE_LIMIT_REQUESTS=100                   # Requests per window
TR_SECURITY_RATE_LIMIT_WINDOW=1m                      # Time window
TR_SECURITY_RATE_LIMIT_HEADERS=true                   # Add rate limit headers
TR_SECURITY_RATE_LIMIT_BYPASS=127.0.0.1               # Bypass for this IP

# Security Headers
TR_SECURITY_ENABLE_SECURITY_HEADERS=true              # Enable security headers
TR_SECURITY_STRICT_TRANSPORT_SECURITY=true           # HSTS header
TR_SECURITY_CONTENT_TYPE_NOSNIFF=true                 # X-Content-Type-Options
TR_SECURITY_FRAME_OPTIONS=DENY                        # X-Frame-Options
TR_SECURITY_XSS_PROTECTION=true                       # XSS protection header

# IP Filtering
TR_SECURITY_ALLOWED_IPS=192.168.1.0/24,10.0.0.0/8    # Allowed IP ranges
TR_SECURITY_BLOCKED_IPS=                              # Blocked IP addresses
TR_SECURITY_IP_HEADER=X-Real-IP                      # Custom IP header
```

### Security Examples

```bash
# Development Configuration
TR_SECURITY_CSRF_SECRET=dev-secret
TR_SECURITY_ENABLE_RATE_LIMIT=false
TR_SECURITY_ENABLE_SECURITY_HEADERS=true
TR_SECURITY_STRICT_TRANSPORT_SECURITY=false

# Production Configuration
TR_SECURITY_CSRF_SECRET=$(openssl rand -base64 32)
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_SECURITY_RATE_LIMIT_REQUESTS=60
TR_SECURITY_ENABLE_SECURITY_HEADERS=true
TR_SECURITY_STRICT_TRANSPORT_SECURITY=true
TR_SECURITY_FRAME_OPTIONS=DENY
```

## Logging Configuration

Logging and monitoring settings:

```bash
# General Logging
TR_LOGGING_LEVEL=info                                 # Log level (debug, info, warn, error)
TR_LOGGING_FORMAT=json                                # Log format (json, text)
TR_LOGGING_OUTPUT=stdout                              # Log output (stdout, stderr, file)
TR_LOGGING_ENABLE_FILE=false                          # Enable file logging
TR_LOGGING_FILE_PATH=logs/app.log                    # Log file path
TR_LOGGING_FILE_MAX_SIZE=100MB                        # Max file size
TR_LOGGING_FILE_MAX_BACKUPS=5                         # Max backup files
TR_LOGGING_FILE_MAX_AGE=30d                           # Max file age

# Request Logging
TR_LOGGING_REQUEST_BODY_ENABLED=false                 # Log request bodies
TR_LOGGING_RESPONSE_BODY_ENABLED=false                # Log response bodies
TR_LOGGING_MAX_BODY_SIZE=1KB                          # Max body size to log
TR_LOGGING_QUERY_PARAMS_ENABLED=true                  # Log query parameters
TR_LOGGING_HEADERS_ENABLED=false                      # Log request headers

# Structured Logging
TR_LOGGING_INCLUDE_TIMESTAMP=true                     # Include timestamp
TR_LOGGING_INCLUDE_LEVEL=true                         # Include log level
TR_LOGGING_INCLUDE_CALLER=false                       # Include caller info
TR_LOGGING_INCLUDE_STACK_TRACE=false                  # Include stack trace on errors
TR_LOGGING_FIELDS=service,version                     # Additional log fields
```

### Logging Examples

```bash
# Development Configuration
TR_LOGGING_LEVEL=debug
TR_LOGGING_FORMAT=text
TR_LOGGING_OUTPUT=stdout
TR_LOGGING_REQUEST_BODY_ENABLED=true
TR_LOGGING_INCLUDE_CALLER=true

# Production Configuration
TR_LOGGING_LEVEL=info
TR_LOGGING_FORMAT=json
TR_LOGGING_OUTPUT=stdout
TR_LOGGING_ENABLE_FILE=true
TR_LOGGING_FILE_PATH=/var/log/app/app.log
TR_LOGGING_FILE_MAX_SIZE=500MB
TR_LOGGING_FIELDS=service,version,instance_id
```

## Database Configuration

Database connection and pool settings:

```bash
# Connection Settings
TR_DATABASE_HOST=localhost                          # Database host
TR_DATABASE_PORT=5432                              # Database port
TR_DATABASE_USER=postgres                          # Database user
TR_DATABASE_PASSWORD=postgres                      # Database password
TR_DATABASE_NAME=router_db                        # Database name
TR_DATABASE_SSL_MODE=disable                       # SSL mode
TR_DATABASE_TIMEZONE=UTC                           # Database timezone

# Connection Pool
TR_DATABASE_MAX_OPEN_CONNS=25                      # Max open connections
TR_DATABASE_MAX_IDLE_CONNS=5                       # Max idle connections
TR_DATABASE_CONN_MAX_LIFETIME=1h                   # Connection max lifetime
TR_DATABASE_CONN_MAX_IDLE_TIME=30m                 # Max idle time for connections
TR_DATABASE_HEALTH_CHECK_PERIOD=1m                 # Health check period

# Database Options
TR_DATABASE_ENABLE_LOGGING=false                   # Enable query logging
TR_DATABASE_SLOW_QUERY_THRESHOLD=100ms             # Slow query threshold
TR_DATABASE_CONNECT_TIMEOUT=10s                    # Connection timeout
TR_DATABASE_QUERY_TIMEOUT=30s                      # Query timeout
```

### Database Examples

```bash
# Development Configuration
TR_DATABASE_HOST=localhost
TR_DATABASE_PORT=5432
TR_DATABASE_USER=postgres
TR_DATABASE_PASSWORD=postgres
TR_DATABASE_NAME=router_dev
TR_DATABASE_SSL_MODE=disable

# Production Configuration
TR_DATABASE_HOST=db.example.com
TR_DATABASE_PORT=5432
TR_DATABASE_USER=app_user
TR_DATABASE_PASSWORD=$(vault get db-password)
TR_DATABASE_NAME=router_prod
TR_DATABASE_SSL_MODE=require
TR_DATABASE_MAX_OPEN_CONNS=50
TR_DATABASE_MAX_IDLE_CONNS=10
```

## Environment Configuration

Application environment settings:

```bash
# Environment Type
TR_ENVIRONMENT_KIND=develop                        # Environment (develop, staging, production)
TR_ENVIRONMENT_DEBUG=false                         # Debug mode
TR_ENVIRONMENT_PROFILING=false                     # Enable profiling

# Feature Flags
TR_FEATURE_NEW_UI=false                            # Feature flags
TR_FEATURE_BETA_ACCESS=false
TR_FEATURE_MAINTENANCE_MODE=false

# Performance Settings
TR_ENVIRONMENT_GOMAXPROCS=0                        # Max CPU cores (0 = auto)
TR_ENVIRONMENT_GC_PERCENTAGE=100                   # GC target percentage
TR_ENVIRONMENT_MAX_MEMORY=0                        # Max memory usage (0 = unlimited)

# Monitoring
TR_ENVIRONMENT_METRICS_ENABLED=true                # Enable metrics
TR_ENVIRONMENT_METRICS_PORT=9090                   # Metrics port
TR_ENVIRONMENT_HEALTH_CHECK_ENABLED=true           # Enable health checks
TR_ENVIRONMENT_HEALTH_CHECK_PATH=/health          # Health check path
```

### Environment Examples

```bash
# Development Environment
TR_ENVIRONMENT_KIND=develop
TR_ENVIRONMENT_DEBUG=true
TR_ENVIRONMENT_PROFILING=true
TR_LOGGING_LEVEL=debug
TR_TEMPLATE_CACHE_ENABLED=false

# Production Environment
TR_ENVIRONMENT_KIND=production
TR_ENVIRONMENT_DEBUG=false
TR_ENVIRONMENT_PROFILING=false
TR_LOGGING_LEVEL=info
TR_TEMPLATE_CACHE_ENABLED=true
TR_SECURITY_ENABLE_RATE_LIMIT=true
```

## Configuration Validation

### Required Settings

These settings must be configured:

```bash
# Server Configuration
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080

# Authentication (if enabled)
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_SESSION_COOKIE_NAME=session_id

# Internationalization
TR_I18N_SUPPORTED_LOCALES=en
TR_I18N_DEFAULT_LOCALE=en
```

### Configuration Validation

```go
// Built-in configuration validation
func ValidateConfig() error {
    // Validate required settings
    if config.Server.Host == "" {
        return fmt.Errorf("TR_SERVER_HOST is required")
    }

    // Validate port range
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("TR_SERVER_PORT must be between 1 and 65535")
    }

    // Validate locales
    if !contains(config.I18n.SupportedLocales, config.I18n.DefaultLocale) {
        return fmt.Errorf("TR_I18N_DEFAULT_LOCALE must be in TR_I18N_SUPPORTED_LOCALES")
    }

    return nil
}
```

## Configuration Files

### .env File

Create a `.env` file for local development:

```bash
# .env
# Server Configuration
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080
TR_SERVER_BASE_URL=http://localhost:8080

# Authentication
TR_AUTH_SESSION_EXPIRY=24h
TR_AUTH_CREATE_DEFAULT_ADMIN=true
TR_AUTH_DEFAULT_ADMIN_EMAIL=admin@example.com
TR_AUTH_DEFAULT_ADMIN_PASSWORD=admin123

# Internationalization
TR_I18N_SUPPORTED_LOCALES=en,de
TR_I18N_DEFAULT_LOCALE=en

# Development Settings
TR_ENVIRONMENT_KIND=develop
TR_LOGGING_LEVEL=debug
TR_TEMPLATE_CACHE_ENABLED=false
```

### docker-compose.yml

```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - TR_SERVER_HOST=0.0.0.0
      - TR_SERVER_PORT=8080
      - TR_DATABASE_HOST=postgres
      - TR_DATABASE_USER=postgres
      - TR_DATABASE_PASSWORD=postgres
      - TR_DATABASE_NAME=router_db
      - TR_AUTH_SESSION_EXPIRY=24h
      - TR_I18N_SUPPORTED_LOCALES=en,de
      - TR_I18N_DEFAULT_LOCALE=en
      - TR_LOGGING_LEVEL=info
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=router_db
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### Kubernetes ConfigMap

```yaml
# k8s-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: templ-router-config
data:
  # Server Configuration
  TR_SERVER_HOST: "0.0.0.0"
  TR_SERVER_PORT: "8080"

  # Database Configuration
  TR_DATABASE_HOST: "postgres-service"
  TR_DATABASE_PORT: "5432"
  TR_DATABASE_USER: "app_user"
  TR_DATABASE_NAME: "router_prod"
  TR_DATABASE_SSL_MODE: "require"

  # Authentication
  TR_AUTH_SESSION_EXPIRY: "1h"
  TR_AUTH_SESSION_COOKIE_SECURE: "true"
  TR_AUTH_SESSION_COOKIE_HTTP_ONLY: "true"

  # Internationalization
  TR_I18N_SUPPORTED_LOCALES: "en-US,de-DE,fr-FR"
  TR_I18N_DEFAULT_LOCALE: "en-US"

  # Production Settings
  TR_ENVIRONMENT_KIND: "production"
  TR_LOGGING_LEVEL: "info"
  TR_LOGGING_FORMAT: "json"
  TR_SECURITY_ENABLE_RATE_LIMIT: "true"
  TR_TEMPLATE_CACHE_ENABLED: "true"
```

## Configuration Management

### Custom Configuration

Create custom configuration structs:

```go
// pkg/config/config.go
package config

import (
    "time"
    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    // Server Configuration
    Server ServerConfig `envconfig:"SERVER"`

    // Authentication Configuration
    Auth AuthConfig `envconfig:"AUTH"`

    // Internationalization Configuration
    I18n I18nConfig `envconfig:"I18N"`

    // Custom Application Configuration
    MyApp MyAppConfig `envconfig:"MYAPP"`
}

type ServerConfig struct {
    Host            string        `envconfig:"HOST" default:"localhost"`
    Port            int           `envconfig:"PORT" default:"8080"`
    BaseURL         string        `envconfig:"BASE_URL" default:"http://localhost:8080"`
    ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"30s"`
    WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"30s"`
    IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT" default:"120s"`
    ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

type MyAppConfig struct {
    FeatureFlag     bool   `envconfig:"FEATURE_FLAG" default:"false"`
    ApiKey          string `envconfig:"API_KEY" required:"true"`
    WebhookURL      string `envconfig:"WEBHOOK_URL"`
    CacheTTL        time.Duration `envconfig:"CACHE_TTL" default:"1h"`
}

func Load(prefix string) (*Config, error) {
    var config Config
    err := envconfig.Process(prefix, &config)
    return &config, err
}
```

### Configuration Validation

```go
func (c *Config) Validate() error {
    // Validate server configuration
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("server port must be between 1 and 65535")
    }

    // Validate authentication configuration
    if c.Auth.SessionExpiry < time.Minute {
        return fmt.Errorf("session expiry must be at least 1 minute")
    }

    // Validate internationalization configuration
    if len(c.I18n.SupportedLocales) == 0 {
        return fmt.Errorf("at least one supported locale is required")
    }

    if !contains(c.I18n.SupportedLocales, c.I18n.DefaultLocale) {
        return fmt.Errorf("default locale must be in supported locales")
    }

    return nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

## Environment-Specific Configurations

### Development Environment

```bash
# .env.development
TR_ENVIRONMENT_KIND=develop
TR_SERVER_HOST=localhost
TR_SERVER_PORT=8080
TR_LOGGING_LEVEL=debug
TR_TEMPLATE_CACHE_ENABLED=false
TR_AUTH_CREATE_DEFAULT_ADMIN=true
TR_SECURITY_ENABLE_RATE_LIMIT=false
```

### Staging Environment

```bash
# .env.staging
TR_ENVIRONMENT_KIND=staging
TR_SERVER_HOST=0.0.0.0
TR_SERVER_PORT=80
TR_LOGGING_LEVEL=info
TR_TEMPLATE_CACHE_ENABLED=true
TR_AUTH_CREATE_DEFAULT_ADMIN=false
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_DATABASE_SSL_MODE=require
```

### Production Environment

```bash
# .env.production
TR_ENVIRONMENT_KIND=production
TR_SERVER_HOST=0.0.0.0
TR_SERVER_PORT=443
TR_LOGGING_LEVEL=warn
TR_TEMPLATE_CACHE_ENABLED=true
TR_AUTH_SESSION_EXPIRY=1h
TR_SECURITY_ENABLE_RATE_LIMIT=true
TR_DATABASE_SSL_MODE=require
TR_SECURITY_CSRF_SECRET=$(openssl rand -base64 32)
```

## Best Practices

### Configuration Management

1. **Use Environment Prefixes**: Always use a unique prefix for your application
2. **Provide Defaults**: Set sensible defaults for all configuration options
3. **Validate Configuration**: Validate all required settings on startup
4. **Document Configuration**: Document all configuration options
5. **Use Configuration Files**: Use `.env` files for local development

### Security

1. **Secrets Management**: Never commit secrets to version control
2. **Use Strong Secrets**: Use cryptographically strong random values
3. **Enable Security Features**: Enable all security features in production
4. **Regular Rotation**: Rotate secrets and certificates regularly
5. **Least Privilege**: Use minimal required permissions

### Performance

1. **Optimize Timeouts**: Set appropriate timeouts for your environment
2. **Enable Caching**: Enable template and response caching in production
3. **Connection Pooling**: Use appropriate database connection pool sizes
4. **Resource Limits**: Set appropriate resource limits
5. **Monitor Performance**: Monitor resource usage and performance metrics

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up basic configuration
- **[Authentication](AUTHENTICATION.md)** - Configure authentication settings
- **[Internationalization](INTERNATIONALIZATION.md)** - Configure i18n settings
- **[Production Deployment](PRODUCTION-DEPLOYMENT.md)** - Production configuration
- **[Security](../README.md#security)** - Security configuration

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Documentation](../README.md)** - Main documentation