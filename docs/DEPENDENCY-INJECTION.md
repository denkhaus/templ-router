# Dependency Injection

**Complete guide to the dependency injection system in Templ Router.**

## Overview

Templ Router uses [samber/do/v2](https://github.com/samber/do) for dependency injection, providing a clean, type-safe, and flexible way to manage services and their dependencies. The DI system separates router services from application services, allowing for clear separation of concerns and easy testing.

**Key Features:**
- Two-tier DI system (Router Services + Application Services)
- Type-safe service resolution
- Named service registration for data services
- Generic service registration for application services
- Clean separation of concerns
- Easy testing with mock implementations

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

All environment variable examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix.

## DI Container Architecture

The DI system consists of two main layers:

### 1. Router Services Layer
Built-in services that power the router functionality:
- Template Service
- Authentication Service
- Internationalization Service
- Configuration Service
- Route Discovery Service

### 2. Application Services Layer
User-defined services specific to your application:
- Data Services
- Business Logic Services
- External API Services
- Database Services

## Container Setup

### Basic Container Creation

```go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/samber/do/v2"
)

func main() {
    // Create DI container
    container := di.NewContainer()
    defer container.Shutdown()

    // Register router services with configurable prefix
    container.RegisterRouterServices("TR")  // "TR" is the prefix for env vars

    // Get injector for registering application services
    injector := container.GetInjector()

    // Your application setup continues...
}
```

### Container Lifecycle

```go
// Container creation
container := di.NewContainer()

// Register services (see sections below)
registerServices(container)

// Use container
injector := container.GetInjector()
service := do.MustInvoke[*MyService](injector)

// Cleanup when done
defer container.Shutdown()
```

## Router Services Registration

Router services are automatically registered with a configurable prefix:

```go
// Register router services with custom prefix
container.RegisterRouterServices("MYAPP")

// Environment variables will now use MYAPP_ prefix:
// MYAPP_SERVER_PORT, MYAPP_AUTH_SESSION_EXPIRY, etc.
```

### Available Router Services

Router services provide these interfaces:

```go
// Core router services
router := container.GetRouter()                    // interfaces.Router
templateRegistry := container.GetTemplateRegistry()   // interfaces.TemplateRegistry
authService := container.GetAuthService()            // interfaces.AuthService
i18nService := container.GetI18nService()            // interfaces.I18nService
configService := container.GetConfigService()        // interfaces.ConfigService

// Utility services
logger := container.GetLogger()                      // shared.Logger
cache := container.GetCache()                        // shared.CacheService
```

## Application Services Registration

### Named Service Registration

Named services are used for data services and other services that need to be resolved by name:

```go
import (
    "github.com/denkhaus/templ-router/pkg/dataservices"
    "github.com/samber/do/v2"
)

func registerApplicationServices(injector *do.Injector) error {
    // Register data services (named)
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
    do.ProvideNamed(injector, "OrderDataService", dataservices.NewOrderDataService)

    // Register other named services
    do.ProvideNamed(injector, "EmailService", services.NewEmailService)
    do.ProvideNamed(injector, "PaymentService", services.NewPaymentService)

    return nil
}
```

### Generic Service Registration

Generic services are resolved by their type:

```go
import (
    "github.com/yourproject/pkg/services"
    "github.com/samber/do/v2"
)

func registerApplicationServices(injector *do.Injector) error {
    // Register services by type
    do.Provide(injector, services.NewDatabaseService)
    do.Provide(injector, services.NewCacheService)
    do.Provide(injector, services.NewEmailService)

    return nil
}
```

### Application Services Options Pattern

Use the options pattern to register application services:

```go
func main() {
    container := di.NewContainer()
    container.RegisterRouterServices("TR")

    // Create application services
    templateRegistry, _ := templates.NewRegistry(container.GetInjector())
    userStore, _ := services.NewUserStore(container.GetInjector())

    // Register using options pattern
    container.RegisterApplicationServices(
        di.WithTemplateRegistry(templateRegistry),
        di.WithUserStore(userStore),
        di.WithLogger(customLogger),
    )
}
```

## Service Resolution

### Named Service Resolution

```go
// Resolve named service
userDataService, err := di.ResolveNamed[UserDataService](injector, "UserDataService")
if err != nil {
    return fmt.Errorf("failed to resolve UserDataService: %w", err)
}

// Or using do.MustInvoke
userDataService := do.MustInvokeNamed[UserDataService](injector, "UserDataService")
```

### Generic Service Resolution

```go
// Resolve service by type
dbService := do.MustInvoke[*DatabaseService](injector)
cacheService := do.MustInvoke[*CacheService](injector)
logger := do.MustInvoke[*Logger](injector)
```

### Router Services Access

```go
// Access built-in router services
router := container.GetRouter()
authService := container.GetAuthService()
i18nService := container.GetI18nService()
configService := container.GetConfigService()
```

## Service Dependencies

### Automatic Dependency Injection

Services automatically receive their dependencies through constructor injection:

```go
package services

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
    "github.com/samber/do/v2"
)

type EmailService struct {
    logger shared.Logger
    config *shared.ConfigService
}

func NewEmailService(injector *do.Injector) (*EmailService, error) {
    return &EmailService{
        logger: do.MustInvoke[shared.Logger](injector),
        config: do.MustInvoke[*shared.ConfigService](injector),
    }, nil
}

type OrderService struct {
    db      *DatabaseService
    email   *EmailService
    logger  shared.Logger
}

func NewOrderService(injector *do.Injector) (*OrderService, error) {
    return &OrderService{
        db:     do.MustInvoke[*DatabaseService](injector),
        email:  do.MustInvoke[*EmailService](injector),
        logger: do.MustInvoke[shared.Logger](injector),
    }, nil
}
```

### Circular Dependencies

The DI container handles circular dependencies automatically:

```go
// Circular dependency is resolved by the DI container
type ServiceA struct {
    serviceB *ServiceB
}

type ServiceB struct {
    serviceA *ServiceA
}

func NewServiceA(injector *do.Injector) (*ServiceA, error) {
    return &ServiceA{
        serviceB: do.MustInvoke[*ServiceB](injector),
    }, nil
}

func NewServiceB(injector *do.Injector) (*ServiceB, error) {
    return &ServiceB{
        serviceA: do.MustInvoke[*ServiceA](injector),
    }, nil
}
```

## Service Lifetime Management

### Singleton Services

By default, all services are singletons:

```go
// All services are singletons by default
do.Provide(injector, NewDatabaseService)

// Both calls return the same instance
db1 := do.MustInvoke[*DatabaseService](injector)
db2 := do.MustInvoke[*DatabaseService](injector)
// db1 == db2
```

### Custom Service Factories

You can provide custom factories for complex initialization:

```go
func NewDatabaseServiceWithConfig(config *DatabaseConfig) *DatabaseService {
    return &DatabaseService{
        host: config.Host,
        port: config.Port,
        // Custom initialization
    }
}

// Register with factory
do.ProvideValue(injector, NewDatabaseServiceWithConfig(dbConfig))
```

## Advanced Patterns

### Service Groups

Organize related services into groups:

```go
// Service group registration
func registerUserServices(injector *do.Injector) {
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "UserRepository", repositories.NewUserRepository)
    do.ProvideNamed(injector, "UserValidator", validators.NewUserValidator)
}

func registerProductServices(injector *do.Injector) {
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
    do.ProvideNamed(injector, "ProductRepository", repositories.NewProductRepository)
    do.ProvideNamed(injector, "ProductCache", caches.NewProductCache)
}

// Register all service groups
registerUserServices(injector)
registerProductServices(injector)
```

### Conditional Service Registration

Register services based on configuration:

```go
func registerServices(injector *do.Injector, config *Config) error {
    // Always register core services
    do.Provide(injector, services.NewCoreService)

    // Register development services in development mode
    if config.IsDevelopment() {
        do.Provide(injector, services.NewDevToolsService)
        do.Provide(injector, services.NewMockEmailService)
    }

    // Register production services in production mode
    if config.IsProduction() {
        do.Provide(injector, services.NewProductionEmailService)
        do.Provide(injector, services.NewMetricsService)
    }

    // Register services based on feature flags
    if config.Features.EnableCaching {
        do.Provide(injector, services.NewRedisCache)
    }

    return nil
}
```

### Service Interfaces

Program against interfaces for better testability:

```go
// Define interfaces
type UserRepository interface {
    FindByID(id string) (*User, error)
    Create(user *User) error
}

type EmailSender interface {
    SendEmail(to, subject, body string) error
}

// Implementations
type SQLUserRepository struct {
    db *sql.DB
}

type SMTPEmailSender struct {
    client *smtp.Client
}

// Register implementations against interfaces
func NewUserRepository(injector *do.Injector) (UserRepository, error) {
    return &SQLUserRepository{
        db: do.MustInvoke[*sql.DB](injector),
    }, nil
}

func NewEmailSender(injector *do.Injector) (EmailSender, error) {
    return &SMTPEmailSender{
        client: do.MustInvoke[*smtp.Client](injector),
    }, nil
}
```

## Testing with DI

### Mock Service Registration

```go
func TestUserService(t *testing.T) {
    // Create test container
    injector := do.New()

    // Register mock services
    do.Provide(injector, func() *MockUserRepository {
        return &MockUserRepository{
            users: map[string]*User{
                "123": {ID: "123", Name: "Test User"},
            },
        }
    })

    do.Provide(injector, func(userRepo *MockUserRepository) *UserService {
        return NewUserService(userRepo)
    })

    // Test with mocks
    userService := do.MustInvoke[*UserService](injector)
    user, err := userService.GetUser("123")

    require.NoError(t, err)
    assert.Equal(t, "Test User", user.Name)
}
```

### Test Containers

Create test-specific containers:

```go
func createTestContainer() *di.Container {
    container := di.NewContainer()
    container.RegisterRouterServices("TR_TEST")

    // Override services for testing
    injector := container.GetInjector()
    do.Provide(injector, services.NewMockDatabaseService)
    do.Provide(injector, services.NewMockEmailService)

    return container
}

func TestOrderService(t *testing.T) {
    container := createTestContainer()
    defer container.Shutdown()

    orderService := do.MustInvoke[*OrderService](container.GetInjector())
    // Test with mocked dependencies
}
```

### Integration Testing

```go
func TestFullWorkflow(t *testing.T) {
    // Setup container with all real services
    container := di.NewContainer()
    container.RegisterRouterServices("TR_TEST")

    // Register test data
    injector := container.GetInjector()
    do.ProvideNamed(injector, "UserDataService", dataservices.NewTestDataService)

    // Test the complete workflow
    userService := do.MustInvokeNamed[UserDataService](injector, "UserDataService")
    result := userService.ProcessUserWorkflow("123")

    assert.NotNil(t, result)
    assert.Equal(t, "processed", result.Status)
}
```

## Best Practices

### Service Organization

1. **Group related services** in separate files/packages
2. **Use interfaces** for all service boundaries
3. **Keep constructors simple** and focused on initialization
4. **Register services at application startup**
5. **Avoid business logic in constructors**

### Error Handling

```go
func NewService(injector *do.Injector) (*Service, error) {
    // Validate dependencies
    db, err := do.Invoke[*DatabaseService](injector)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve database service: %w", err)
    }

    // Validate configuration
    config, err := do.Invoke[*ConfigService](injector)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve config service: %w", err)
    }

    if !config.IsValid() {
        return nil, fmt.Errorf("invalid configuration")
    }

    return &Service{
        db:    db,
        config: config,
    }, nil
}
```

### Configuration

```go
// Configure services based on environment
func registerServices(injector *do.Injector, config *Config) error {
    switch config.Environment {
    case "development":
        do.Provide(injector, services.NewDevelopmentEmailService)
        do.Provide(injector, services.NewMockPaymentService)
    case "staging":
        do.Provide(injector, services.NewSMTPEmailService)
        do.Provide(injector, services.NewTestPaymentService)
    case "production":
        do.Provide(injector, services.NewSMTPEmailService)
        do.Provide(injector, services.NewProductionPaymentService)
    default:
        return fmt.Errorf("unknown environment: %s", config.Environment)
    }

    return nil
}
```

### Naming Conventions

```go
// Service interfaces
type UserService interface{...}
type EmailService interface{...}
type PaymentService interface{...}

// Service implementations
type userService struct{...}
type emailService struct{...}
type paymentService struct{...}

// Constructor functions
func NewUserService(...)*userService{...}
func NewEmailService(...)*emailService{...}
func NewPaymentService(...)*paymentService{...}
```

## Troubleshooting

### Common DI Issues

#### Service Not Found

```go
// Error: service not found
service := do.MustInvoke[*MyService](injector)

// Solution: Ensure service is registered
do.Provide(injector, NewMyService)
```

#### Circular Dependency

```go
// Error: circular dependency detected
// Solution: Use interfaces or break the circular dependency
type ServiceA struct {
    serviceB ServiceBInterface  // Use interface instead of concrete type
}
```

#### Type Mismatch

```go
// Error: type mismatch
// Solution: Ensure types match exactly
do.Provide(injector, func() *MyService { return &MyService{} })
// Must resolve as *MyService, not MyService
```

### Debugging DI

```go
// List all registered services
services := do.GetHealth(injector)
for _, service := range services {
    fmt.Printf("Service: %s\n", service.Name)
}

// Check if a service is registered
if do.HealthCheck[injector, *MyService]() {
    fmt.Println("MyService is registered")
} else {
    fmt.Println("MyService is not registered")
}
```

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first DI container
- **[Data Services](DATA-SERVICES.md)** - Use DI with data services
- **[Authentication](AUTHENTICATION.md)** - DI in authentication system
- **[Configuration](CONFIGURATION.md)** - Configure DI behavior
- **[Testing](../CONTRIBUTING.md#testing)** - Testing with DI

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[samber/do Documentation](https://github.com/samber/do)** - DI library documentation