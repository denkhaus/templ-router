# Architecture Overview

**Complete guide to the architecture and design principles of Templ Router.**

## Overview

Templ Router follows clean architecture principles with a focus on separation of concerns, dependency injection, and modular design. The system is built as a library that developers can integrate into their own Go applications, providing file-based routing, template management, and comprehensive middleware support.

**Key Architectural Principles:**
- **Clean Architecture** - Clear separation between infrastructure, application, and domain layers
- **Dependency Injection** - Type-safe service management with samber/do/v2
- **File-Based Routing** - Convention over configuration for route generation
- **Template System** - Type-safe templating with metadata management
- **Middleware Pipeline** - Configurable request processing chain
- **Service Orientation** - Modular services with clear interfaces

## Configuration Prefix Notice

**Important:** Router services are configured with environment variables using a configurable prefix. The default prefix is `TR_`, but this is fully configurable when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_SERVER_HOST=localhost`
- Custom: `MYAPP_SERVER_HOST=localhost`
- Multiple apps: `APP1_SERVER_HOST=localhost` and `APP2_SERVER_HOST=localhost`

All environment variable examples in this documentation use the default `TR_` prefix, but you can replace `TR` with your custom prefix. This design allows multiple Templ Router applications to run in the same environment without conflicts.

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Application Layer                           │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   HTTP Routes   │  │  Auth Endpoints │  │  API Endpoints  │ │
│  │                 │  │                 │  │                 │ │
│  │  • Chi Router   │  │  • Signin       │  │  • REST APIs    │ │
│  │  • Middleware   │  │  • Signout      │  │  • Custom APIs  │ │
│  │  • Handlers     │  │  • Signup       │  │  • Data APIs    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                    Middleware Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Auth Mw       │  │   I18n Mw       │  │  Template Mw    │ │
│  │                 │  │                 │  │                 │ │
│  │  • Session Mgmt │  │  • Locale       │  │  • Template     │ │
│  │  • Route Guard  │  │  • Translation  │  │  • Rendering    │ │
│  │  • Context      │  │  • URL Gen      │  │  • Metadata     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                    Service Layer                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │  Router Service │  │  Auth Service   │  │ Template Svc    │ │
│  │                 │  │                 │  │                 │ │
│  │  • Route Disc   │  │  • Session      │  │  • Template     │ │
│  │  • Template Reg │  │  • User Mgmt    │  │  • Registry     │ │
│  │  • Handler      │  │  • Validation   │  │  • Rendering    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                 Infrastructure Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   DI Container  │  │   Config Svc    │  │   Storage Svc   │ │
│  │                 │  │                 │  │                 │ │
│  │  • Service Reg  │  │  • Env Vars     │  │  • File System  │ │
│  │  • Lifecycle    │  │  • Validation   │  │  • Cache        │ │
│  │  • Resolution   │  │  • Defaults     │  │  • Metadata     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### Component Interaction

```
Request → Middleware → Router → Template → Response
    │           │           │           │
    │           │           │           ├─► Template Service
    │           │           │           ├─► Data Services (DI)
    │           │           │           └─► Registry (trgen)
    │           │           ├─► Route Discovery
    │           │           ├─► Handler Builder
    │           │           └─► URL Parameters
    │           ├─► Auth Context
    │           ├─► Locale Detection
    │           └─► Error Handling
    └─► HTTP Router (Chi)
```

## Core Components

### 1. Router Core

**Location**: `pkg/router/`

The router core is the heart of the system, responsible for route discovery, template management, and request handling.

```go
type Router struct {
    // Dependencies injected via DI
    templateRegistry interfaces.TemplateRegistry
    authService      interfaces.AuthService
    configService    interfaces.ConfigService
    i18nService      interfaces.I18nService

    // Internal state
    routeMap        map[string]Route
    templateCache   map[string]templ.Component
    middlewareChain []func(http.Handler) http.Handler
}
```

**Key Responsibilities:**
- **Route Discovery**: Automatic route generation from file structure
- **Template Registry**: Template registration and management
- **Handler Building**: Request handler creation and execution
- **URL Parameter Parsing**: Dynamic parameter extraction
- **Error Handling**: Template resolution and error fallbacks

### 2. Template System

**Location**: `pkg/router/services/`

The template system manages template rendering, metadata, and component composition.

```go
type TemplateService struct {
    registry        interfaces.TemplateRegistry
    configService   interfaces.ConfigService
    metadataLoader  *MetadataLoader
    i18nService     interfaces.I18nService
    cache           map[string]CachedTemplate
}

type TemplateMetadata struct {
    Auth       AuthConfig       `yaml:"auth"`
    I18n       I18nConfig       `yaml:"i18n"`
    Metadata   map[string]interface{} `yaml:"metadata"`
    Validation ValidationConfig `yaml:"validation"`
}
```

**Template Processing Pipeline:**

```
Template File → Metadata Parsing → Template Compilation → Component Creation
       │                 │                    │                   │
       ├─► .templ Files  ├─► .templ.yaml      ├─► templ Generate  ├─► Registry
       └─► Components    └─► Component YAML  └─► Type Safety     └─► Runtime
```

### 3. Dependency Injection Container

**Location**: `pkg/di/`

Two-tier DI system separating router services from application services.

```go
type Container struct {
    injector     *do.Injector
    routerServices   []do.Provider
    appServices      []do.Provider
}

type ServiceRegistry struct {
    // Router Services (built-in)
    Router           interfaces.Router
    TemplateRegistry interfaces.TemplateRegistry
    AuthService      interfaces.AuthService
    ConfigService    interfaces.ConfigService

    // Application Services (user-defined)
    DataServices     map[string]interface{}
    BusinessServices []interface{}
}
```

**Service Lifecycle:**

1. **Container Creation**: `di.NewContainer()`
2. **Router Service Registration**: `container.RegisterRouterServices("TR")`
3. **Application Service Registration**: `container.RegisterApplicationServices(...)`
4. **Service Resolution**: `do.MustInvoke[*Service](injector)`
5. **Cleanup**: `container.Shutdown()`

### 4. Middleware Pipeline

**Location**: `pkg/router/middleware/`

Configurable middleware pipeline for request processing.

```go
type MiddlewareChain struct {
    auth      *AuthMiddleware
    i18n      *I18nMiddleware
    template  *TemplateMiddleware
    error     *ErrorHandlingMiddleware
    custom    []http.Middleware
}

type Middleware interface {
    Middleware(next http.Handler) http.Handler
}
```

**Middleware Execution Order:**

```
Request → Auth → I18n → Template → Custom → Handler → Response
    │       │      │         │        │        │
    │       │      │         │        │        ├─► Template Rendering
    │       │      │         │        │        ├─► Data Service Resolution
    │       │      │         │        │        └─► Response Writing
    │       │      │         │        └─► Custom Processing
    │       │      │         └─► Template Context
    │       │      └─► Locale/Translation
    │       └─► Authentication/Authorization
    └─► HTTP Router (Chi)
```

## Data Flow

### Request Processing Flow

```
1. HTTP Request Received
   ↓
2. Chi Router Matching
   ↓
3. Middleware Pipeline
   ├─► Authentication Middleware
   │   • Session validation
   │   • User context injection
   │   • Route guard checks
   │   ↓
   ├─► I18n Middleware
   │   • Locale detection
   │   • Translation loading
   │   • URL localization
   │   ↓
   ├─► Template Middleware
   │   • Template discovery
   │   • Metadata loading
   │   • Data service resolution
   │   ↓
   └─► Custom Middleware
       • Custom processing
       • Error handling
       • Response headers
   ↓
4. Handler Execution
   ├─► Template rendering
   ├─► Data service calls
   ├─► Context injection
   ↓
5. Response Generation
   ├─► Template output
   ├─► Error handling
   └─► HTTP status codes
```

### Template Rendering Flow

```
Template Request
   ↓
Template Discovery
   ├─► Route mapping
   ├─► File lookup
   └─► Template registry
   ↓
Metadata Loading
   ├─► Component metadata (highest priority)
   ├─► Page metadata (middle priority)
   └─► Layout metadata (lowest priority)
   ↓
Data Service Resolution
   ├─► Service discovery
   ├─► Parameter injection
   └─► Data fetching
   ↓
Template Rendering
   ├─► Component composition
   ├─► Layout inheritance
   ├─► I18n translation
   └─► Context data injection
   ↓
Response Output
```

## Configuration Architecture

### Environment-Based Configuration

The system uses environment variables with configurable prefixes for configuration:

```go
type Config struct {
    Server    ServerConfig    `envPrefix:"TR_SERVER_"`
    Auth      AuthConfig      `envPrefix:"TR_AUTH_"`
    I18n      I18nConfig      `envPrefix:"TR_I18N_"`
    Router    RouterConfig    `envPrefix:"TR_ROUTER_"`
    Template  TemplateConfig  `envPrefix:"TR_TEMPLATE_"`
    Logging   LoggingConfig   `envPrefix:"TR_LOGGING_"`
}
```

### Configuration Hierarchy

```
1. Environment Variables (highest priority)
   └─► Runtime configuration
   └─► Environment-specific settings
   └─► Secrets and sensitive data

2. Template Metadata Files (middle priority)
   ├─► Component metadata: component.templ.yaml
   ├─► Page metadata: page.templ.yaml
   └─► Layout metadata: layout.templ.yaml

3. Default Values (lowest priority)
   └─► Built-in defaults
   └─► Fallback configuration
   └─► Safe initial values
```

## Security Architecture

### Authentication & Authorization

```go
type AuthService struct {
    sessionStore   SessionStore
    userValidator  UserValidator
    configService  interfaces.ConfigService
}

type AuthContext struct {
    UserID     string
    UserName   string
    Roles      []string
    SessionID  string
    ExpiresAt  time.Time
}

type AuthLevel string
const (
    AuthPublic        AuthLevel = "Public"
    AuthUserRequired  AuthLevel = "UserRequired"
    AuthAdminRequired AuthLevel = "AdminRequired"
)
```

### Security Layers

1. **Session Management**
   - Secure cookie configuration
   - Session expiration handling
   - Session store interface

2. **Route Guards**
   - Template-based auth configuration
   - Role-based access control
   - Redirect handling

3. **Input Validation**
   - Parameter validation
   - XSS prevention
   - CSRF protection

4. **Security Headers**
   - Content Security Policy
   - X-Frame-Options
   - X-Content-Type-Options

## Performance Architecture

### Caching Strategy

```go
type CacheManager struct {
    templateCache  map[string]templ.Component
    metadataCache  map[string]TemplateMetadata
    i18nCache      map[string]map[string]string
    configCache    map[string]interface{}
}

type CacheConfig struct {
    TemplateEnabled   bool          `env:"TEMPLATE_CACHE_ENABLED"`
    TemplateTTL       time.Duration `env:"TEMPLATE_CACHE_TTL"`
    MetadataEnabled   bool          `env:"METADATA_CACHE_ENABLED"`
    I18nEnabled       bool          `env:"I18N_CACHE_ENABLED"`
}
```

### Performance Optimizations

1. **Template Caching**
   - Pre-compiled templates
   - Metadata caching
   - Component caching

2. **Route Optimization**
   - Efficient route matching
   - Parameter extraction
   - Handler reuse

3. **I18n Optimization**
   - Translation caching
   - Locale detection caching
   - Lazy loading

4. **Data Service Optimization**
   - Connection pooling
   - Query caching
   - Batch operations

## Extensibility Architecture

### Plugin Points

1. **Custom Middleware**
   ```go
   type CustomMiddleware struct {
       config interface{}
       logger Logger
   }

   func (m *CustomMiddleware) Middleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           // Custom processing
           next.ServeHTTP(w, r)
       })
   }
   ```

2. **Data Services**
   ```go
   type CustomDataService interface {
       GetData(routerCtx interfaces.RouterContext) (*CustomData, error)
   }

   do.ProvideNamed(injector, "CustomDataService", dataservices.NewCustomDataService)
   ```

3. **Template Components**
   ```go
   // Custom component with metadata
   templ CustomComponent(data *CustomData) {
       <div class="custom-component">
           { data.Content }
       </div>
   }
   ```

### Extension Patterns

1. **Service Registration**
   - Named services for data services
   - Generic services for business logic
   - Interface-based programming

2. **Middleware Composition**
   - Chainable middleware
   - Conditional registration
   - Custom error handling

3. **Template Composition**
   - Component reusability
   - Layout inheritance
   - Metadata overriding

## Testing Architecture

### Test Organization

```
tests/
├── unit/                    # Unit tests
│   ├── services/           # Service tests
│   ├── middleware/         # Middleware tests
│   └── utils/              # Utility tests
├── integration/            # Integration tests
│   ├── api/                # API endpoint tests
│   ├── auth/               # Authentication tests
│   └── i18n/               # Internationalization tests
├── e2e/                    # End-to-end tests
│   ├── routing/            # Routing tests
│   ├── authentication/     # Auth flow tests
│   ├── internationalization/ # I18n tests
│   └── performance/        # Performance tests
└── fixtures/               # Test data
    ├── templates/          # Test templates
    ├── configs/            # Test configs
    └── data/               # Test data
```

### Testing Patterns

1. **Unit Testing**
   - Interface mocking
   - Service isolation
   - Dependency injection

2. **Integration Testing**
   - Service composition
   - Middleware chains
   - Database integration

3. **E2E Testing**
   - Browser automation
   - Full request/response cycles
   - Multi-language testing

## Deployment Architecture

### Container Strategy

```dockerfile
# Multi-stage build
FROM golang:1.24-alpine AS builder
# Build application
FROM alpine:latest AS runtime
# Runtime image
COPY --from=builder /app /app
```

### Deployment Models

1. **Single Instance**
   - Simple deployment
   - Local development
   - Small applications

2. **Load-Balanced**
   - Multiple instances
   - Health checks
   - Session affinity

3. **Container Orchestration**
   - Kubernetes deployment
   - Service discovery
   - Configuration management

## Development Workflow

### Architecture Considerations

1. **Local Development**
   - Hot reload support
   - Development server
   - Debugging tools

2. **Template Development**
   - Live reloading
   - Component isolation
   - Metadata validation

3. **Service Development**
   - Interface design
   - Dependency injection
   - Testing patterns

### Build Pipeline

```
1. Template Generation (templ generate)
   ↓
2. Registry Generation (trgen)
   ↓
3. Go Build (go build)
   ↓
4. Testing (mage test:all)
   ↓
5. Packaging (docker build)
   ↓
6. Deployment
```

## Architectural Decisions

### Key Design Choices

1. **File-Based Routing**
   - **Rationale**: Convention over configuration, reduced boilerplate
   - **Trade-off**: Less flexible than programmatic routing
   - **Impact**: Simpler development, clearer project structure

2. **Dependency Injection**
   - **Rationale**: Testability, loose coupling, modularity
   - **Trade-off**: Additional complexity in service setup
   - **Impact**: Better testability, cleaner architecture

3. **Type-Safe Templates**
   - **Rationale**: Compile-time safety, better IDE support
   - **Trade-off**: Additional build step
   - **Impact**: Fewer runtime errors, better developer experience

4. **Metadata System**
   - **Rationale**: Configuration close to templates, clear precedence
   - **Trade-off**: Additional file management
   - **Impact**: Self-documenting templates, better organization

5. **Middleware Pipeline**
   - **Rationale**: Composable request processing
   - **Trade-off**: Performance overhead
   - **Impact**: Flexible request handling, cleaner code

### Future Considerations

1. **Performance Optimization**
   - Template pre-compilation
   - Route caching
   - Connection pooling

2. **Scalability Enhancements**
   - Horizontal scaling support
   - Distributed caching
   - Load balancing strategies

3. **Feature Additions**
   - WebSocket support
   - GraphQL integration
   - Advanced caching strategies

## Best Practices

### Architectural Guidelines

1. **Separation of Concerns**
   - Clear layer boundaries
   - Interface-based programming
   - Dependency injection

2. **Configuration Management**
   - Environment-based configuration
   - Sensitive data protection
   - Validation and defaults

3. **Error Handling**
   - Graceful degradation
   - Proper HTTP status codes
   - Comprehensive logging

4. **Security First**
   - Input validation
   - Secure defaults
   - Regular security updates

5. **Performance Awareness**
   - Efficient resource usage
   - Proper caching strategies
   - Monitoring and metrics

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first application
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Understanding the routing system
- **[Data Services](DATA-SERVICES.md)** - Building data-driven applications
- **[Middleware System](MIDDLEWARE.md)** - Custom middleware development
- **[Production Deployment](PRODUCTION-DEPLOYMENT.md)** - Deployment strategies

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute to the project