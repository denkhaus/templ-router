# Generic User Implementation Guide

Templ Router now supports generic user implementations, allowing you to use your own user structure.

## UserEntity Interface

Every user implementation must implement the `UserEntity` interface:

```go
type UserEntity interface {
    GetID() string
    GetEmail() string
    GetRoles() []string
}
```

## Example: Custom User Implementation

### 1. Define Your User Structure

```go
package myapp

import "github.com/denkhaus/templ-router/pkg/interfaces"

// MyCustomUser - Your custom user implementation
type MyCustomUser struct {
    UserID      string   `json:"user_id"`
    Username    string   `json:"username"`
    EmailAddr   string   `json:"email_address"`
    Permissions []string `json:"permissions"`
    Department  string   `json:"department"`
    IsActive    bool     `json:"is_active"`
}

// UserEntity interface implementation
func (u *MyCustomUser) GetID() string {
    return u.UserID
}

func (u *MyCustomUser) GetEmail() string {
    return u.EmailAddr
}

func (u *MyCustomUser) GetRoles() []string {
    return u.Permissions
}
```

### 2. Implement User Repository or Store

You can implement user data access using any pattern that fits your architecture:

```go
// MyUserRepository implements user data access for MyCustomUser
type MyUserRepository struct {
    // Your database connection, etc.
    db Database
}

func (r *MyUserRepository) FindByID(userID string) (*MyCustomUser, error) {
    // Your implementation
    user, err := r.db.FindUserByID(userID)
    if err != nil {
        return nil, err
    }

    return &MyCustomUser{
        UserID:      user.ID,
        Username:    user.Name,
        EmailAddr:   user.Email,
        Permissions: user.Roles,
        Department:  user.Dept,
        IsActive:    user.Active,
    }, nil
}

func (r *MyUserRepository) FindByEmail(email string) (*MyCustomUser, error) {
    // Your implementation for finding users by email
    // ...
}

func (r *MyUserRepository) ValidateCredentials(email, password string) (*MyCustomUser, error) {
    // Your implementation for credential validation
    // ...
}

// Add other methods as needed for your authentication system
```

### 3. Use Your Implementation in AuthValidator

```go
package myapp

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
)

// MyAuthValidator uses your custom user implementation
type MyAuthValidator struct {
    userRepository *MyUserRepository
    // Add other dependencies like JWT service, etc.
}

func NewMyAuthValidator(injector *do.Injector) (interfaces.AuthValidator, error) {
    // Register your user repository
    userRepo := &MyUserRepository{
        db: NewDatabase(),
    }

    return &MyAuthValidator{
        userRepository: userRepo,
    }, nil
}

func (av *MyAuthValidator) IsAuthenticated(req *http.Request) bool {
    // Your authentication logic using MyCustomUser
    // For example, validate JWT token or check session
    return true
}

func (av *MyAuthValidator) GetCurrentUser(req *http.Request) (interfaces.UserEntity, error) {
    // Get user ID from token/session
    userID := "extract-user-id-from-request"

    // Use your user repository to get MyCustomUser
    return av.userRepository.FindByID(userID)
}

func (av *MyAuthValidator) HasRole(user interfaces.UserEntity, requiredRoles []string) bool {
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

## Benefits of Generic Implementation

1. **Flexibility**: Use your own user structure with any fields you need
2. **Type Safety**: Go's type system ensures everything is properly typed
3. **No Conversions**: Work directly with your user objects
4. **Extensibility**: Add any fields and methods to your user structure

## Integration with Templ Router

To use your custom user implementation with Templ Router:

```go
func main() {
    container := di.NewContainer()
    container.RegisterRouterServices("TR")

    // Register your AuthValidator that uses your custom user
    container.RegisterApplicationServices(
        di.WithAuthValidatorFactory(func(i do.Injector) (interfaces.AuthValidator, error) {
            return myapp.NewMyAuthValidator(i)
        }),
    )

    // Register your user repository as a dependency
    do.Provide(i, myapp.NewUserRepository)

    // Continue with bootstrap...
}
```

## Example: Enterprise User with Additional Fields

```go
type EnterpriseUser struct {
    ID           string    `json:"id"`
    Email        string    `json:"email"`
    Roles        []string  `json:"roles"`
    EmployeeID   string    `json:"employee_id"`
    CostCenter   string    `json:"cost_center"`
    Manager      string    `json:"manager"`
    LastLogin    time.Time `json:"last_login"`
    Preferences  map[string]interface{} `json:"preferences"`
}

// UserEntity Interface Implementation
func (u *EnterpriseUser) GetID() string    { return u.ID }
func (u *EnterpriseUser) GetEmail() string { return u.Email }
func (u *EnterpriseUser) GetRoles() []string { return u.Roles }

// Additional methods for enterprise features
func (u *EnterpriseUser) GetEmployeeID() string { return u.EmployeeID }
func (u *EnterpriseUser) GetCostCenter() string { return u.CostCenter }
func (u *EnterpriseUser) IsManager() bool { return contains(u.Roles, "manager") }
```

This generic solution makes Templ Router much more flexible and enables developers to use their own user models without sacrificing router functionality.