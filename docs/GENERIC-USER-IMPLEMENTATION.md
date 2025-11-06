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

### 2. Implement the UserStore

```go
// MyUserStore implements UserStore for MyCustomUser
type MyUserStore struct {
    // Your database connection, etc.
    db Database
}

func (s *MyUserStore) GetUserByID(userID string) (*MyCustomUser, error) {
    // Your implementation
    user, err := s.db.FindUserByID(userID)
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

func (s *MyUserStore) GetUserByEmail(email string) (*MyCustomUser, error) {
    // Your implementation
    // ...
}

func (s *MyUserStore) ValidateCredentials(email, password string) (*MyCustomUser, error) {
    // Your implementation
    // ...
}

func (s *MyUserStore) CreateUser(username, email, password string) (*MyCustomUser, error) {
    // Your implementation
    // ...
}

func (s *MyUserStore) UserExists(username, email string) (bool, error) {
    // Your implementation
    // ...
}
```

### 3. Use Your Implementation

```go
package main

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
    "myapp"
)

func main() {
    // Create your user store implementation
    userStore := &myapp.MyUserStore{
        db: myapp.NewDatabase(),
    }

    // Use the generic AuthService
    var authService interfaces.AuthService[*myapp.MyCustomUser]

    // Configure your router with the custom user implementation
    // ...
}
```

## Benefits of Generic Implementation

1. **Flexibility**: Use your own user structure with any fields you need
2. **Type Safety**: Go's type system ensures everything is properly typed
3. **No Conversions**: Work directly with your user objects
4. **Extensibility**: Add any fields and methods to your user structure

## Migration from Default User Implementation

If you're already using the default `User` structure, you don't need to change anything. The default implementation already implements the `UserEntity` interface and continues to work:

```go
// Default implementation continues to work
var userStore interfaces.UserStore[*interfaces.User]
var authService interfaces.AuthService[*interfaces.User]
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