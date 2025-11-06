# Data Services

**Complete guide to implementing data services for template data injection in Templ Router.**

## Overview

Data services provide a clean, dependency-injected way to supply data to templates. The system automatically resolves and calls appropriate data services based on template requirements, passing route parameters, query parameters, and request context to each service.

**Key Features:**
- Automatic data service detection and injection
- Two method patterns: `GetData()` fallback and specific `Get<DataStruct>()` methods
- RouterContext interface for unified parameter access
- Support for URL parameters, query parameters, and request context
- Named dependency resolution via DI container
- Automatic parameter injection from routes
- Support for multiple data services per template

## RouterContext Interface

The `RouterContext` interface provides unified access to all request parameters:

```go
// pkg/interfaces/router_context.go
package interfaces

import (
    "context"
    "net/http"
    "github.com/go-chi/chi/v5"
)

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

The simplest pattern using only the generic `GetData()` method:

```go
// pkg/dataservices/product_service.go
package dataservices

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
)

type ProductData struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Category string `json:"category"`
    Price    string `json:"price"`
    InStock  bool   `json:"in_stock"`
}

type ProductDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*ProductData, error)
}

type productDataServiceImpl struct {
    // Dependencies can be injected here
    productRepo ProductRepository
}

func NewProductDataService(injector *do.Injector) (ProductDataService, error) {
    // Inject dependencies
    productRepo := do.MustInvoke[ProductRepository](injector)

    return &productDataServiceImpl{
        productRepo: productRepo,
    }, nil
}

func (s *productDataServiceImpl) GetData(routerCtx interfaces.RouterContext) (*ProductData, error) {
    // Extract route parameters
    productID := routerCtx.GetURLParam("id")

    // Extract query parameters
    category := routerCtx.GetQueryParam("category")
    currency := routerCtx.GetQueryParam("currency")

    // Set defaults
    if currency == "" {
        currency = "USD"
    }

    // Fetch product data
    product, err := s.productRepo.FindByID(productID)
    if err != nil {
        return nil, err
    }

    // Apply filters and transformations
    if category != "" && product.Category != category {
        return nil, fmt.Errorf("product category mismatch")
    }

    return &ProductData{
        ID:       product.ID,
        Name:     product.Name,
        Category: product.Category,
        Price:    s.convertPrice(product.Price, currency),
        InStock:  product.Stock > 0,
    }, nil
}

func (s *productDataServiceImpl) convertPrice(price float64, currency string) string {
    // Currency conversion logic
    return fmt.Sprintf("%.2f %s", price, currency)
}
```

### Pattern 2: GetData() + Specific Methods

Provide both generic and specific methods for flexibility:

```go
// pkg/dataservices/user_service.go
package dataservices

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
)

type UserData struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Locale   string    `json:"locale"`
    JoinDate time.Time `json:"join_date"`
    LastLogin time.Time `json:"last_login"`
}

type UserActivity struct {
    UserID      string `json:"user_id"`
    PageViews   int    `json:"page_views"`
    LastAction  string `json:"last_action"`
    SessionTime string `json:"session_time"`
}

type UserDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
    GetUserActivity(routerCtx interfaces.RouterContext) (*UserActivity, error)
}

type userDataServiceImpl struct {
    userRepo    UserRepository
    activityRepo ActivityRepository
    logger      Logger
}

func NewUserDataService(injector *do.Injector) (UserDataService, error) {
    return &userDataServiceImpl{
        userRepo:    do.MustInvoke[UserRepository](injector),
        activityRepo: do.MustInvoke[ActivityRepository](injector),
        logger:      do.MustInvoke[Logger](injector),
    }, nil
}

func (s *userDataServiceImpl) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    // Extract parameters
    userID := routerCtx.GetURLParam("id")
    locale := routerCtx.GetURLParam("locale")

    // Fetch user data
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        s.logger.Error("Failed to fetch user", "error", err, "userID", userID)
        return nil, err
    }

    // Update last login
    if err := s.userRepo.UpdateLastLogin(userID); err != nil {
        s.logger.Warn("Failed to update last login", "error", err)
    }

    return &UserData{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Locale:    locale,
        JoinDate:  user.CreatedAt,
        LastLogin: time.Now(),
    }, nil
}

func (s *userDataServiceImpl) GetUserActivity(routerCtx interfaces.RouterContext) (*UserActivity, error) {
    userID := routerCtx.GetURLParam("id")

    activity, err := s.activityRepo.GetByUserID(userID)
    if err != nil {
        return nil, err
    }

    return &UserActivity{
        UserID:      userID,
        PageViews:   activity.PageViews,
        LastAction:  activity.LastAction,
        SessionTime: s.formatSessionTime(activity.SessionDuration),
    }, nil
}

func (s *userDataServiceImpl) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
    // Fallback to user data
    return s.GetUserData(routerCtx)
}

func (s *userDataServiceImpl) formatSessionTime(duration time.Duration) string {
    if duration < time.Minute {
        return fmt.Sprintf("%d seconds", int(duration.Seconds()))
    } else if duration < time.Hour {
        return fmt.Sprintf("%d minutes", int(duration.Minutes()))
    } else {
        return fmt.Sprintf("%.1f hours", duration.Hours())
    }
}
```

### Pattern 3: Multiple Data Types

Service that provides multiple related data types:

```go
// pkg/dataservices/dashboard_service.go
package dataservices

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
)

type DashboardStats struct {
    TotalUsers    int     `json:"total_users"`
    ActiveUsers   int     `json:"active_users"`
    TotalRevenue  float64 `json:"total_revenue"`
    RecentOrders  int     `json:"recent_orders"`
    ConversionRate float64 `json:"conversion_rate"`
}

type RecentActivity struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"`
    Description string    `json:"description"`
    Timestamp   time.Time `json:"timestamp"`
    UserID      string    `json:"user_id"`
}

type NotificationData struct {
    Count      int           `json:"count"`
    Items      []Notification `json:"items"`
    LastCheck  time.Time     `json:"last_check"`
}

type DashboardDataService interface {
    GetDashboardStats(routerCtx interfaces.RouterContext) (*DashboardStats, error)
    GetRecentActivity(routerCtx interfaces.RouterContext) ([]RecentActivity, error)
    GetNotifications(routerCtx interfaces.RouterContext) (*NotificationData, error)
}

type dashboardDataServiceImpl struct {
    statsRepo     StatisticsRepository
    activityRepo  ActivityRepository
    notifyRepo    NotificationRepository
    cache         CacheService
}

func NewDashboardDataService(injector *do.Injector) (DashboardDataService, error) {
    return &dashboardDataServiceImpl{
        statsRepo:    do.MustInvoke[StatisticsRepository](injector),
        activityRepo: do.MustInvoke[ActivityRepository](injector),
        notifyRepo:   do.MustInvoke[NotificationRepository](injector),
        cache:        do.MustInvoke[CacheService](injector),
    }, nil
}

func (s *dashboardDataServiceImpl) GetDashboardStats(routerCtx interfaces.RouterContext) (*DashboardStats, error) {
    // Check cache first
    if cached, found := s.cache.Get("dashboard:stats"); found {
        return cached.(*DashboardStats), nil
    }

    // Extract parameters
    locale := routerCtx.GetURLParam("locale")
    timeRange := routerCtx.GetQueryParam("timeRange")
    if timeRange == "" {
        timeRange = "7d"
    }

    // Fetch statistics
    stats, err := s.statsRepo.GetDashboardStats(locale, timeRange)
    if err != nil {
        return nil, err
    }

    dashboardStats := &DashboardStats{
        TotalUsers:     stats.TotalUsers,
        ActiveUsers:    stats.ActiveUsers,
        TotalRevenue:   stats.TotalRevenue,
        RecentOrders:   stats.RecentOrders,
        ConversionRate: stats.ConversionRate,
    }

    // Cache for 5 minutes
    s.cache.Set("dashboard:stats", dashboardStats, 5*time.Minute)

    return dashboardStats, nil
}

func (s *dashboardDataServiceImpl) GetRecentActivity(routerCtx interfaces.RouterContext) ([]RecentActivity, error) {
    // Extract pagination parameters
    page := routerCtx.GetQueryParam("page")
    limit := routerCtx.GetQueryParam("limit")

    // Set defaults
    if page == "" {
        page = "1"
    }
    if limit == "" {
        limit = "10"
    }

    pageNum, _ := strconv.Atoi(page)
    limitNum, _ := strconv.Atoi(limit)

    // Fetch activity
    activities, err := s.activityRepo.GetRecentActivity(pageNum, limitNum)
    if err != nil {
        return nil, err
    }

    // Transform to response format
    result := make([]RecentActivity, len(activities))
    for i, activity := range activities {
        result[i] = RecentActivity{
            ID:          activity.ID,
            Type:        activity.Type,
            Description: activity.Description,
            Timestamp:   activity.CreatedAt,
            UserID:      activity.UserID,
        }
    }

    return result, nil
}

func (s *dashboardDataServiceImpl) GetNotifications(routerCtx interfaces.RouterContext) (*NotificationData, error) {
    userID := routerCtx.GetURLParam("id")

    notifications, err := s.notifyRepo.GetUnreadCount(userID)
    if err != nil {
        return nil, err
    }

    return &NotificationData{
        Count:     notifications.Count,
        Items:     notifications.Items,
        LastCheck: time.Now(),
    }, nil
}
```

## Template Integration

### Single Data Service Template

Template that requires one data service:

```go
// app/user/id_/page.templ
package main

import (
    "fmt"
    "github.com/denkhaus/templ-router/pkg/router/i18n"
)

templ UserProfilePage(user *UserData) {
    <div class="user-profile">
        <div class="profile-header">
            <h1>{ user.Name }</h1>
            <p class="email">{ user.Email }</p>
            <p class="locale">{ i18n.T(ctx, "locale") }: { user.Locale }</p>
        </div>

        <div class="profile-stats">
            <div class="stat">
                <label>{ i18n.T(ctx, "member_since") }</label>
                <span>{ user.JoinDate.Format("Jan 2006") }</span>
            </div>
            <div class="stat">
                <label>{ i18n.T(ctx, "last_login") }</label>
                <span>{ user.LastLogin.Format("Jan 2, 2006 at 3:04 PM") }</span>
            </div>
        </div>

        <div class="profile-actions">
            <a href={ i18n.LocalizeSafeURL(ctx, "/user/" + user.ID + "/edit") }
               class="btn btn-primary">
                { i18n.T(ctx, "edit_profile") }
            </a>
            <a href={ i18n.LocalizeSafeURL(ctx, "/user/" + user.ID + "/settings") }
               class="btn btn-secondary">
                { i18n.T(ctx, "settings") }
            </a>
        </div>
    </div>
}
```

### Template Naming Convention

**Important:** Template functions must follow specific naming conventions:

```go
// app/page.templ - MUST be named "Page"
package main
templ Page() { ... }

// app/user/id_/page.templ - MUST be named "Page"
package main
templ Page(userID string) { ... }

// app/layout.templ - MUST be named "Layout"
package main
templ Layout(content templ.Component) { ... }

// app/error.templ - MUST be named "Error"
package main
templ Error(errCtx middleware.ErrorContext) { ... }
```

## Metadata Precedence Rules

Templ Router uses a hierarchical metadata system with clear precedence rules:

### Precedence Order (Highest to Lowest)

1. **Component Metadata** (highest priority)
2. **Page/Template Metadata** (middle priority)
3. **Layout Metadata** (lowest priority/fallback)

### How Precedence Works

```yaml
# app/layout.templ.yaml (lowest priority)
metadata:
  theme: "light"
  company_name: "Default Company"

i18n:
  en:
    welcome: "Welcome to our app"

# app/dashboard/page.templ.yaml (overrides layout)
metadata:
  theme: "dark"              # Overrides layout theme
  page_title: "Dashboard"     # New metadata

i18n:
  en:
    welcome: "Welcome to dashboard"  # Overrides layout translation

# app/components/navbar.templ.yaml (overrides both layout and page)
metadata:
  theme: "dark"              # Overrides page and layout theme
  navigation_style: "modern" # New metadata only for component
```

### Practical Example

When a component is included in a page that has a layout:

```yaml
# Resulting metadata for navbar component:
theme: "dark"                    # From component (highest priority)
page_title: "Dashboard"         # From page (not relevant to component)
company_name: "Default Company" # From layout (fallback)
navigation_style: "modern"     # From component (highest priority)
```

This system allows components to be truly self-contained while still inheriting sensible defaults from layouts and pages.

## Data Service Limitations

### One Data Service Per Template

**Important:** Templ Router supports **exactly one data service per template**. Templates cannot accept multiple data service parameters.

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

### Composite Data Service Pattern

If you need multiple data sources, create a composite data service:

```go
// Composite data structure
type CompositeDashboardData struct {
    UserStats    *UserData    `json:"user_stats"`
    SystemStats  *SystemData  `json:"system_stats"`
}

// Composite data service
type CompositeDashboardDataService interface {
    GetData(routerCtx interfaces.RouterContext) (*CompositeDashboardData, error)
}

func (s *compositeDataService) GetData(routerCtx interfaces.RouterContext) (*CompositeDashboardData, error) {
    // Get data from multiple services internally
    userData, _ := s.userService.GetData(routerCtx)
    systemData, _ := s.systemService.GetData(routerCtx)

    return &CompositeDashboardData{
        UserStats:   userData,
        SystemStats: systemData,
    }, nil
}

// Template uses composite data
templ Page(data *CompositeDashboardData) {
    <h1>{ data.UserStats.Name }</h1>
    <p>Total users: { data.SystemStats.TotalUsers }</p>
}
```

### Template Metadata with Data Services

Configure data services in template metadata:

```yaml
# app/dashboard/page.templ.yaml
metadata:
  page_title: "Dashboard"
  section: "main"

i18n:
  en:
    dashboard_title: "Dashboard"
    total_users: "Total Users"
    active_users: "Active Users"
    total_revenue: "Total Revenue"
    conversion_rate: "Conversion Rate"
    recent_activity: "Recent Activity"
  de:
    dashboard_title: "Dashboard"
    total_users: "Gesamte Benutzer"
    active_users: "Aktive Benutzer"
    total_revenue: "Gesamteinnahmen"
    conversion_rate: "Konversionsrate"
    recent_activity: "Aktuelle Aktivität"

auth:
  type: "UserRequired"
  redirect_url: "/login"

data_services:
  - "DashboardDataService"
```

## Data Service Registration

### Basic Registration

Register data services with the DI container:

```go
// main.go
package main

import (
    "github.com/denkhaus/templ-router/pkg/di"
    "github.com/denkhaus/templ-router/pkg/services/dataservices"
    "github.com/samber/do/v2"
)

func main() {
    // Create DI container
    container := di.NewContainer()
    injector := container.GetInjector()

    // Register data services as named dependencies
    do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
    do.ProvideNamed(injector, "DashboardDataService", dataservices.NewDashboardDataService)

    // Register application services
    container.RegisterApplicationServices(
        di.WithTemplateRegistry(templateRegistry),
    )
}
```

### Advanced Registration with Dependencies

Register data services with their dependencies:

```go
// pkg/dataservices/registration.go
package dataservices

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
    "github.com/samber/do/v2"
)

func RegisterDataServices(injector *do.Injector) error {
    // Register repositories first
    do.Provide(injector, NewUserRepository)
    do.Provide(injector, NewProductRepository)
    do.Provide(injector, NewStatisticsRepository)
    do.Provide(injector, NewActivityRepository)
    do.Provide(injector, NewNotificationRepository)

    // Register utilities
    do.Provide(injector, NewCacheService)
    do.Provide(injector, NewLogger)

    // Register data services
    do.ProvideNamed(injector, "UserDataService", NewUserDataService)
    do.ProvideNamed(injector, "ProductDataService", NewProductDataService)
    do.ProvideNamed(injector, "DashboardDataService", NewDashboardDataService)

    return nil
}

// Usage in main.go
func main() {
    container := di.NewContainer()
    injector := container.GetInjector()

    // Register all data services
    if err := dataservices.RegisterDataServices(injector); err != nil {
        panic(err)
    }
}
```

## Parameter Access Patterns

### URL Parameters

Access parameters from route patterns like `/{locale}/user/{id}`:

```go
func (s *userDataService) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    // Route: /{locale}/user/{id}/profile
    // URL: /en/user/123/profile

    locale := routerCtx.GetURLParam("locale")  // "en"
    userID := routerCtx.GetURLParam("id")      // "123"

    // Get all URL parameters
    allParams := routerCtx.GetAllURLParams()
    // map[string]string{
    //     "locale": "en",
    //     "id": "123",
    // }

    return s.fetchUser(userID, locale), nil
}
```

### Query Parameters

Access parameters from URL query string like `?page=5&filter=active`:

```go
func (s *productService) GetProducts(routerCtx interfaces.RouterContext) ([]Product, error) {
    // URL: /products?category=electronics&sort=price&order=asc&page=2&limit=10

    category := routerCtx.GetQueryParam("category")     // "electronics"
    sort := routerCtx.GetQueryParam("sort")            // "price"
    order := routerCtx.GetQueryParam("order")          // "asc"
    page := routerCtx.GetQueryParam("page")            // "2"
    limit := routerCtx.GetQueryParam("limit")          // "10"

    // Multiple values for the same parameter
    tags := routerCtx.GetQueryParams("tags")           // ["sale", "featured", "new"]

    // Get all query parameters
    allQueryParams := routerCtx.GetAllQueryParams()
    // map[string][]string{
    //     "category": ["electronics"],
    //     "sort": ["price"],
    //     "order": ["asc"],
    //     "page": ["2"],
    //     "limit": ["10"],
    //     "tags": ["sale", "featured", "new"],
    // }

    // Set defaults
    if page == "" {
        page = "1"
    }
    if limit == "" {
        limit = "20"
    }
    if order == "" {
        order = "desc"
    }

    // Convert and validate
    pageNum, _ := strconv.Atoi(page)
    limitNum, _ := strconv.Atoi(limit)

    return s.fetchProducts(category, sort, order, pageNum, limitNum, tags), nil
}
```

### Advanced Parameter Handling

```go
func (s *searchService) GetSearchResults(routerCtx interfaces.RouterContext) (*SearchResults, error) {
    // Helper method for default values
    query := routerCtx.GetParamWithDefault("q", "")
    page := routerCtx.GetParamWithDefault("page", "1")
    limit := routerCtx.GetParamWithDefault("limit", "20")

    // Check if parameter exists
    hasFilter := routerCtx.HasParam("filter")
    if hasFilter {
        filter := routerCtx.GetQueryParam("filter")
        // Apply filter
    }

    // Complex parameter parsing
    priceRange := routerCtx.GetQueryParam("price_range")
    var minPrice, maxPrice float64
    if priceRange != "" {
        parts := strings.Split(priceRange, "-")
        if len(parts) == 2 {
            minPrice, _ = strconv.ParseFloat(parts[0], 64)
            maxPrice, _ = strconv.ParseFloat(parts[1], 64)
        }
    }

    // Parse date ranges
    dateFrom := routerCtx.GetQueryParam("date_from")
    dateTo := routerCtx.GetQueryParam("date_to")
    var parsedFrom, parsedTo time.Time
    if dateFrom != "" {
        parsedFrom, _ = time.Parse("2006-01-02", dateFrom)
    }
    if dateTo != "" {
        parsedTo, _ = time.Parse("2006-01-02", dateTo)
    }

    // Build search criteria
    criteria := SearchCriteria{
        Query:     query,
        Page:      parseInt(page),
        Limit:     parseInt(limit),
        MinPrice:  minPrice,
        MaxPrice:  maxPrice,
        DateFrom:  parsedFrom,
        DateTo:    parsedTo,
    }

    return s.performSearch(criteria), nil
}

// Helper function
func parseInt(s string) int {
    if i, err := strconv.Atoi(s); err == nil {
        return i
    }
    return 0
}
```

## Error Handling

### Service Error Handling

Implement proper error handling in data services:

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
        s.logger.Error("Failed to fetch user", "error", err, "userID", userID)
        return nil, fmt.Errorf("failed to fetch user data")
    }

    // Check user permissions
    if !s.canUserAccessData(routerCtx, user) {
        return nil, fmt.Errorf("access denied for user: %s", userID)
    }

    return &UserData{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

func (s *userService) canUserAccessData(routerCtx interfaces.RouterContext, user *User) bool {
    // Implement access control logic
    // Check current user's permissions against requested user
    return true
}
```

### Template Error Handling

Handle data service errors gracefully in templates:

```go
templ UserProfilePageWithErrorHandling() {
    // Use templ's error handling
    if user, err := getUserData(ctx); err != nil {
        <div class="error-message">
            <h2>{ i18n.T(ctx, "error_title") }</h2>
            <p>{ i18n.T(ctx, "user_not_found") }</p>
            <a href={ i18n.LocalizeSafeURL(ctx, "/dashboard") }
               class="btn btn-primary">
                { i18n.T(ctx, "back_to_dashboard") }
            </a>
        </div>
    } else {
        <div class="user-profile">
            <h1>{ user.Name }</h1>
            <p>{ user.Email }</p>
        </div>
    }
}
```

## Testing Data Services

### Unit Testing

Test data services with mock RouterContext:

```go
// pkg/dataservices/user_service_test.go
package dataservices

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockRouterContext struct {
    mock.Mock
    urlParams   map[string]string
    queryParams map[string][]string
    request     *http.Request
}

func (m *MockRouterContext) GetURLParam(key string) string {
    args := m.Called(key)
    return args.String(0)
}

func (m *MockRouterContext) GetQueryParam(key string) string {
    args := m.Called(key)
    return args.String(0)
}

func (m *MockRouterContext) Context() context.Context {
    args := m.Called()
    return args.Get(0).(context.Context)
}

func (m *MockRouterContext) Request() *http.Request {
    args := m.Called()
    return args.Get(0).(*http.Request)
}

func TestUserDataService_GetUserData(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := &userDataServiceImpl{userRepo: mockRepo}

    // Create mock RouterContext
    req := httptest.NewRequest("GET", "/en/user/123", nil)
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("locale", "en")
    rctx.URLParams.Add("id", "123")
    req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

    mockCtx := &MockRouterContext{}
    mockCtx.On("GetURLParam", "locale").Return("en")
    mockCtx.On("GetURLParam", "id").Return("123")
    mockCtx.On("Context").Return(req.Context())
    mockCtx.On("Request").Return(req)

    expectedUser := &User{
        ID:   "123",
        Name: "John Doe",
        Email: "john@example.com",
    }
    mockRepo.On("FindByID", "123").Return(expectedUser, nil)

    // Act
    result, err := service.GetUserData(mockCtx)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "123", result.ID)
    assert.Equal(t, "John Doe", result.Name)
    assert.Equal(t, "john@example.com", result.Email)
    assert.Equal(t, "en", result.Locale)

    mockRepo.AssertExpectations(t)
    mockCtx.AssertExpectations(t)
}

func TestUserDataService_GetUserData_UserNotFound(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := &userDataServiceImpl{userRepo: mockRepo}

    mockCtx := &MockRouterContext{}
    mockCtx.On("GetURLParam", "locale").Return("en")
    mockCtx.On("GetURLParam", "id").Return("999")

    mockRepo.On("FindByID", "999").Return(nil, ErrUserNotFound)

    // Act
    result, err := service.GetUserData(mockCtx)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "user not found")

    mockRepo.AssertExpectations(t)
}
```

### Integration Testing

Test data services with real HTTP requests:

```go
func TestDataServiceIntegration(t *testing.T) {
    // Setup test server
    mux := chi.NewRouter()

    // Register routes with data services
    router := setupTestRouter(t)
    if err := router.RegisterRoutes(mux); err != nil {
        t.Fatal(err)
    }

    server := httptest.NewServer(mux)
    defer server.Close()

    // Test request
    resp, err := http.Get(server.URL + "/en/user/123")
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Parse response
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        t.Fatal(err)
    }

    // Assert response contains expected data
    assert.Contains(t, response, "user")
    user := response["user"].(map[string]interface{})
    assert.Equal(t, "123", user["id"])
}
```

## Performance Optimization

### Caching Data Services

Implement caching for frequently accessed data:

```go
type CachedUserDataService struct {
    base    UserDataService
    cache   CacheService
    ttl     time.Duration
}

func NewCachedUserDataService(base UserDataService, cache CacheService) UserDataService {
    return &CachedUserDataService{
        base:  base,
        cache: cache,
        ttl:   5 * time.Minute,
    }
}

func (s *CachedUserDataService) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
    userID := routerCtx.GetURLParam("id")
    cacheKey := fmt.Sprintf("user:%s", userID)

    // Try cache first
    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*UserData), nil
    }

    // Fetch from base service
    data, err := s.base.GetUserData(routerCtx)
    if err != nil {
        return nil, err
    }

    // Cache the result
    s.cache.Set(cacheKey, data, s.ttl)

    return data, nil
}
```

### Concurrent Data Loading

Load multiple data services concurrently:

```go
templ DashboardPageConcurrent() {
    {{
        // Use templ's concurrency features
        user := getUserData(ctx)
        stats := getDashboardStats(ctx)
        activity := getRecentActivity(ctx)
    }}

    <div class="dashboard">
        <h1>Welcome, { user.Name }!</h1>
        <div class="stats">{ stats.TotalUsers } users</div>
        <div class="activity">{ len(activity) } recent activities</div>
    </div>
}
```

## Best Practices

### Service Design

1. **Single Responsibility**: Each service should handle one type of data
2. **Dependency Injection**: Inject all dependencies through constructor
3. **Interface Segregation**: Use small, focused interfaces
4. **Error Handling**: Provide clear error messages and handle edge cases
5. **Validation**: Validate input parameters before processing

### Parameter Handling

1. **Validate Parameters**: Always validate required parameters
2. **Provide Defaults**: Set sensible defaults for optional parameters
3. **Type Safety**: Convert parameters to appropriate types with error handling
4. **Security**: Sanitize and validate user input
5. **Logging**: Log parameter access for debugging and auditing

### Performance

1. **Caching**: Cache frequently accessed data
2. **Database Optimization**: Use efficient queries and indexing
3. **Concurrent Loading**: Load multiple data sources in parallel when possible
4. **Pagination**: Implement pagination for large datasets
5. **Resource Management**: Properly manage database connections and other resources

### Testing

1. **Unit Tests**: Test each service method in isolation
2. **Mock Dependencies**: Use mocks for external dependencies
3. **Integration Tests**: Test service integration with real data
4. **Error Scenarios**: Test error handling and edge cases
5. **Performance Tests**: Test service performance under load

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first data service
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Use data services with routing
- **[Authentication](AUTHENTICATION.md)** - Integrate with authentication
- **[Dependency Injection](DEPENDENCY-INJECTION.md)** - Advanced DI patterns
- **[Configuration](CONFIGURATION.md)** - Configure data service behavior

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Documentation](../README.md)** - Main documentation