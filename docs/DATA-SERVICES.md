# Data Services

**Complete guide to implementing data services for template data injection in Templ Router.**

## Overview

Data services provide clean, dependency-injected data access for templates. The system automatically resolves and calls appropriate data services based on template requirements, passing route parameters, query parameters, and request context.

**Key Features:**
- Automatic data service detection and injection
- Two method patterns: `GetData()` fallback and specific `Get<DataStruct>()` methods
- RouterContext interface for unified parameter access
- Support for URL parameters, query parameters, and request context
- Named dependency resolution via DI container
- One data service per template with composite pattern support

## RouterContext Interface

The `RouterContext` interface provides unified access to all request parameters:

```go
type RouterContext interface {
    // Core context access
    Context() context.Context
    Request() *http.Request
    ChiContext() *chi.Context

    // URL Parameters (from route patterns like /{locale}/user/{id})
    GetURLParam(key string) string
    GetAllURLParams() map[string]string

    // Query Parameters (from URL query string like ?page=5&filter=active)
    GetQueryParam(key string) string
    GetQueryParams(key string) []string
    GetAllQueryParams() map[string][]string

    // Helper methods
    GetParamWithDefault(key, defaultValue string) string
    HasParam(key string) bool
}
```

## Data Service Patterns

### Pattern 1: GetData() Only

Simple pattern using only the generic `GetData()` method:

```go
type ProductDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*ProductData, error)
}

func (s *productDataService) GetData(routerCtx interfaces.RouterContext) (*ProductData, error) {
    productID := routerCtx.GetURLParam("id")
    currency := routerCtx.GetQueryParam("currency")
    if currency == "" {
        currency = "USD"
    }

    product, err := s.productRepo.FindByID(productID)
    if err != nil {
        return nil, err
    }

    return &ProductData{
        ID:       product.ID,
        Name:     product.Name,
        Price:    s.convertPrice(product.Price, currency),
        InStock:  product.Stock > 0,
    }, nil
}
```

### Pattern 2: GetData() + Specific Methods

Provide both generic and specific methods for flexibility:

```go
type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)      // Fallback
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)  // Specific
    GetUserActivity(routerCtx interfaces.RouterContext) (*UserActivity, error)
}

func (s *userDataService) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")
    locale := routerCtx.GetURLParam("locale")

    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }

    return &UserData{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Locale: locale,
    }, nil
}

func (s *userDataService) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
    return s.GetUserData(routerCtx)  // Fallback implementation
}
```

### Pattern 3: Composite Data Services

For multiple data types, create composite services:

```go
type CompositeDashboardData struct {
    UserStats   *UserData     `json:"user_stats"`
    SystemStats *SystemData   `json:"system_stats"`
}

type CompositeDashboardDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*CompositeDashboardData, error)
}

func (s *service) GetData(routerCtx interfaces.RouterContext) (*CompositeDashboardData, error) {
    userData, _ := s.userService.GetData(routerCtx)
    systemData, _ := s.systemService.GetData(routerCtx)

    return &CompositeDashboardData{
        UserStats:   userData,
        SystemStats: systemData,
    }, nil
}
```

## Template Integration

### Template Usage

Templates automatically receive data from registered services:

```go
// app/user/id_/page.templ - MUST be named "Page"
package main

templ Page(user *UserData) {
    <div class="user-profile">
        <div class="profile-header">
            <h1>{ user.Name }</h1>
            <p class="email">{ user.Email }</p>
            <p>{ i18n.T(ctx, "locale") }: { user.Locale }</p>
        </div>

        <div class="profile-actions">
            <a href={ i18n.LocalizeSafeURL(ctx, "/user/" + user.ID + "/edit") }
               class="btn btn-primary">
                { i18n.T(ctx, "edit_profile") }
            </a>
        </div>
    </div>
}
```

**Template Naming Conventions:**
- Pages: `templ Page()` or `templ Page(param string)`
- Layouts: `templ Layout(content templ.Component)`
- Errors: `templ Error(errCtx middleware.ErrorContext)`

## Data Service Limitations

### One Data Service Per Template

**Important:** Templ Router supports **exactly one data service per template**.

```go
// ✅ SUPPORTED - One data service
templ Page(user *UserData) {
    <h1>{ user.Name }</h1>
}

// ❌ NOT SUPPORTED - Multiple data services
templ Page(user *UserData, stats *DashboardStats) {
    // This will NOT work
}
```

### Composite Pattern Solution

Create composite data services for multiple data sources:

```go
type CompositeDashboardData struct {
    UserStats   *UserData   `json:"user_stats"`
    SystemStats *SystemData `json:"system_stats"`
}

templ Page(data *CompositeDashboardData) {
    <h1>{ data.UserStats.Name }</h1>
    <p>Total users: { data.SystemStats.TotalUsers }</p>
}
```

### Template Metadata Configuration

```yaml
# app/dashboard/page.templ.yaml
metadata:
  page_title: "Dashboard"

auth:
  type: "UserRequired"
  redirect_url: "/login"

i18n:
  en:
    dashboard_title: "Dashboard"
    total_users: "Total Users"
  de:
    dashboard_title: "Dashboard"
    total_users: "Gesamte Benutzer"
```

## Data Service Registration

### Basic Registration

Register data services as named dependencies:

```go
// main.go
func main() {
    container := di.NewContainer()
    injector := container.GetInjector()

    // Register data services
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)

    // Register application services
    container.RegisterApplicationServices(
        di.WithTemplateRegistryFactory(func(injector interface{}) (interface{}, error) {
            return templateRegistry, nil
        }),
    )
}
```

### Advanced Registration Pattern

Create registration functions for complex setups:

```go
// pkg/dataservices/registration.go
func RegisterDataServices(injector *do.Injector) error {
    // Register repositories
    do.Provide(injector, NewUserRepository)
    do.Provide(injector, NewProductRepository)

    // Register utilities
    do.Provide(injector, NewCacheService)
    do.Provide(injector, NewLogger)

    // Register data services
    do.ProvideNamed(injector, "UserDataService", NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", NewProductDataService)

    return nil
}
```

## Parameter Access Patterns

### URL Parameters

Access route parameters from patterns like `/{locale}/user/{id}`:

```go
func (s *userDataService) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    // Route: /{locale}/user/{id}/profile
    locale := routerCtx.GetURLParam("locale")  // "en"
    userID := routerCtx.GetURLParam("id")      // "123"

    return s.fetchUser(userID, locale), nil
}
```

### Query Parameters

Access query parameters from URLs like `?page=5&filter=active`:

```go
func (s *productService) GetProducts(routerCtx interfaces.RouterContext) ([]Product, error) {
    category := routerCtx.GetQueryParam("category")     // "electronics"
    sort := routerCtx.GetQueryParam("sort")            // "price"
    page := routerCtx.GetQueryParam("page")            // "2"

    // Multiple values
    tags := routerCtx.GetQueryParams("tags")           // ["sale", "featured"]

    // Set defaults
    if page == "" {
        page = "1"
    }

    return s.fetchProducts(category, sort, page, tags), nil
}
```

### Advanced Parameter Handling

```go
func (s *searchService) GetSearchResults(routerCtx interfaces.RouterContext) (*SearchResults, error) {
    // Helper methods
    query := routerCtx.GetParamWithDefault("q", "")
    page := routerCtx.GetParamWithDefault("page", "1")

    // Check parameter existence
    if routerCtx.HasParam("filter") {
        filter := routerCtx.GetQueryParam("filter")
        // Apply filter
    }

    // Complex parsing
    priceRange := routerCtx.GetQueryParam("price_range")
    if priceRange != "" {
        parts := strings.Split(priceRange, "-")
        // Parse range values
    }

    return s.performSearch(query, parseInt(page)), nil
}
```

## Error Handling

### Service Error Handling

```go
func (s *userService) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")

    // Validate required parameters
    if userID == "" {
        return nil, fmt.Errorf("user ID is required")
    }

    // Fetch user with error handling
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, fmt.Errorf("user not found: %s", userID)
        }
        return nil, fmt.Errorf("failed to fetch user data")
    }

    return &UserData{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

### Template Error Handling

```go
templ UserProfilePageWithErrorHandling() {
    if user, err := getUserData(ctx); err != nil {
        <div class="error-message">
            <h2>{ i18n.T(ctx, "error_title") }</h2>
            <p>{ i18n.T(ctx, "user_not_found") }</p>
            <a href={ i18n.LocalizeSafeURL(ctx, "/dashboard") }>Back</a>
        </div>
    } else {
        <div class="user-profile">
            <h1>{ user.Name }</h1>
            <p>{ user.Email }</p>
        </div>
    }
}
```

## Testing

### Unit Testing

Test data services with mock RouterContext:

```go
func TestUserDataService_GetUserData(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := &userDataServiceImpl{userRepo: mockRepo}

    mockCtx := &MockRouterContext{}
    mockCtx.On("GetURLParam", "id").Return("123")

    expectedUser := &User{ID: "123", Name: "John Doe"}
    mockRepo.On("FindByID", "123").Return(expectedUser, nil)

    result, err := service.GetUserData(mockCtx)
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", result.Name)

    mockRepo.AssertExpectations(t)
}
```

### Integration Testing

Test with real HTTP requests:

```go
func TestDataServiceIntegration(t *testing.T) {
    server := httptest.NewServer(setupTestRouter(t))
    defer server.Close()

    resp, err := http.Get(server.URL + "/en/user/123")
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

## Performance Optimization

### Caching

```go
type CachedUserDataService struct {
    base  UserDataService
    cache CacheService
    ttl   time.Duration
}

func (s *CachedUserDataService) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")
    cacheKey := fmt.Sprintf("user:%s", userID)

    // Try cache first
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*UserData), nil
    }

    // Fetch and cache result
    data, err := s.base.GetUserData(routerCtx)
    if err == nil {
        s.cache.Set(cacheKey, data, s.ttl)
    }

    return data, err
}
```

## Best Practices

### Service Design
- **Single Responsibility**: Each service handles one data type
- **Dependency Injection**: Inject all dependencies through constructor
- **Interface Segregation**: Use focused interfaces
- **Error Handling**: Provide clear error messages

### Parameter Handling
- **Validate Parameters**: Always validate required parameters
- **Provide Defaults**: Set sensible defaults for optional parameters
- **Type Safety**: Convert parameters with proper error handling

### Performance
- **Caching**: Cache frequently accessed data
- **Database Optimization**: Use efficient queries and indexing
- **Resource Management**: Properly manage connections and resources

### Testing
- **Unit Tests**: Test each service method in isolation
- **Mock Dependencies**: Use mocks for external dependencies
- **Integration Tests**: Test with real data
- **Error Scenarios**: Test error handling and edge cases

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first data service
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Use data services with routing
- **[Authentication](AUTHENTICATION.md)** - Integrate with authentication
- **[Dependency Injection](DEPENDENCY-INJECTION.md)** - Advanced DI patterns