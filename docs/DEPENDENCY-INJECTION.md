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

## Configuration Prefix

Router services are configured with environment variables using a configurable prefix:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**

- Default: `TR_SERVER_HOST=localhost`
- Custom: `MYAPP_SERVER_HOST=localhost`

## Quick Start

```go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/denkhaus/templ-router/demo/generated/templates"
    "github.com/denkhaus/templ-router/demo/pkg/dataservices"
    "github.com/denkhaus/templ-router/demo/pkg/services"
    "github.com/samber/do/v2"
)

func main() {
    // Create DI container
    container := di.NewContainer()
    defer container.Shutdown()

    // Register router services with prefix
    container.RegisterRouterServices("TR")

    // Register application services with template registry factory
    container.RegisterApplicationServices(
        di.WithTemplateRegistryFactory(templates.NewRegistry),
        di.WithAuthValidatorFactory(services.NewAuthValidator),
    )

    // Register data services (named)
    injector := container.GetInjector()
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)

    // Register additional services
    do.Provide(injector, services.NewUserStore)
    do.Provide(injector, services.NewSessionStore)

    // Bootstrap router and start server...
}
```

## Service Registration

### Router Services

Built-in services automatically registered:

```go
// Core router services
router := container.GetRouter()
templateRegistry := container.GetTemplateRegistry()
configService := container.GetConfigService()
logger := container.GetLogger()
```

### Application Services

#### Named Services (Data Services)

Used for data services resolved by name:

```go
do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)

// Resolve named service
userDataService := do.MustInvokeNamed[UserDataService](injector, "UserDataService")
```

#### Generic Services

Resolved by type:

```go
// Register
do.Provide(injector, services.NewDatabaseService)
do.Provide(injector, services.NewEmailService)

// Resolve
dbService := do.MustInvoke[*DatabaseService](injector)
emailService := do.MustInvoke[*EmailService](injector)
```

#### Options Pattern

Configure using the options pattern:

```go
container.RegisterApplicationServices(
    di.WithTemplateRegistryFactory(func(injector interface{}) (interface{}, error) {
        return templateRegistry, nil
    }),
    di.WithAuthValidatorFactory(func(i do.Injector) (interfaces.AuthValidator, error) {
        return services.NewAuthValidator(i)
    }),
    di.WithAssetsServiceFactory(func(i do.Injector) (interfaces.AssetsService, error) {
        return assets.NewService(i)
    }),
    di.WithHealthCheck(true, "/api/health"),
)
```

## Service Dependencies

Services automatically receive dependencies through constructor injection:

```go
type OrderService struct {
    db     *DatabaseService
    email  *EmailService
    logger shared.Logger
}

func NewOrderService(injector *do.Injector) (*OrderService, error) {
    return &OrderService{
        db:     do.MustInvoke[*DatabaseService](injector),
        email:  do.MustInvoke[*EmailService](injector),
        logger: do.MustInvoke[shared.Logger](injector),
    }, nil
}
```

## Advanced Patterns

### Conditional Registration

Register services based on configuration:

```go
func registerServices(injector *do.Injector, config *Config) error {
    // Always register core services
    do.Provide(injector, services.NewCoreService)

    // Environment-specific services
    if config.IsDevelopment() {
        do.Provide(injector, services.NewDevToolsService)
        do.Provide(injector, services.NewMockEmailService)
    } else if config.IsProduction() {
        do.Provide(injector, services.NewProductionEmailService)
        do.Provide(injector, services.NewMetricsService)
    }

    return nil
}
```

### Interface-Based Programming

Program against interfaces for better testability:

```go
// Define interface
type UserRepository interface {
    FindByID(id string) (*User, error)
    Create(user *User) error
}

// Implementation
type SQLUserRepository struct {
    db *sql.DB
}

// Register against interface
func NewUserRepository(injector *do.Injector) (UserRepository, error) {
    return &SQLUserRepository{
        db: do.MustInvoke[*sql.DB](injector),
    }, nil
}
```

## Testing

### Mock Services

```go
func TestUserService(t *testing.T) {
    injector := do.New()

    // Register mock service
    do.Provide(injector, func() *MockUserRepository {
        return &MockUserRepository{
            users: map[string]*User{
                "123": {ID: "123", Name: "Test User"},
            },
        }
    })

    // Register service under test
    do.Provide(injector, func(userRepo *MockUserRepository) *UserService {
        return NewUserService(userRepo)
    })

    // Test
    userService := do.MustInvoke[*UserService](injector)
    user, err := userService.GetUser("123")

    require.NoError(t, err)
    assert.Equal(t, "Test User", user.Name)
}
```

### Test Containers

```go
func createTestContainer() *di.Container {
    container := di.NewContainer()
    container.RegisterRouterServices("TR_TEST")

    // Override with mocks
    injector := container.GetInjector()
    do.Provide(injector, services.NewMockDatabaseService)
    do.Provide(injector, services.NewMockEmailService)

    return container
}
```

## Best Practices

1. **Use interfaces** for all service boundaries
2. **Keep constructors simple** - only initialization logic
3. **Group related services** in logical packages
4. **Register services at startup**, not during request handling
5. **Avoid circular dependencies** - use interfaces to break cycles

## Troubleshooting

### Common Issues

**Service Not Found:**

```go
// Ensure service is registered before resolving
do.Provide(injector, NewMyService)
service := do.MustInvoke[*MyService](injector)
```

**Circular Dependency:**

```go
// Use interfaces to break cycles
type ServiceA struct {
    serviceB ServiceBInterface  // Interface, not concrete type
}
```

**Type Mismatch:**

```go
// Types must match exactly
do.Provide(injector, func() *MyService { return &MyService{} })
// Resolve as *MyService, not MyService
```

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Complete setup tutorial
- **[Data Services](DATA-SERVICES.md)** - Use DI with data services
- **[Authentication](AUTHENTICATION.md)** - DI in authentication system
- **[Configuration](CONFIGURATION.md)** - Configure all aspects

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[samber/do Documentation](https://github.com/samber/do)** - DI library documentation