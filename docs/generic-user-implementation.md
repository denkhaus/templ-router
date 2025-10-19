# Generic User Implementation Guide

Der Templ-Router unterstützt jetzt generische User-Implementationen, sodass Sie Ihre eigene User-Struktur verwenden können.

## UserEntity Interface

Jede User-Implementation muss das `UserEntity` Interface implementieren:

```go
type UserEntity interface {
    GetID() string
    GetEmail() string
    GetRoles() []string
}
```

## Beispiel: Eigene User-Implementation

### 1. Definieren Sie Ihre User-Struktur

```go
package myapp

import "github.com/denkhaus/templ-router/pkg/interfaces"

// MyCustomUser - Ihre eigene User-Implementation
type MyCustomUser struct {
    UserID      string   `json:"user_id"`
    Username    string   `json:"username"`
    EmailAddr   string   `json:"email_address"`
    Permissions []string `json:"permissions"`
    Department  string   `json:"department"`
    IsActive    bool     `json:"is_active"`
}

// Implementierung des UserEntity Interface
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

### 2. Implementieren Sie den UserStore

```go
// MyUserStore implementiert UserStore für MyCustomUser
type MyUserStore struct {
    // Ihre Datenbank-Verbindung, etc.
    db Database
}

func (s *MyUserStore) GetUserByID(userID string) (*MyCustomUser, error) {
    // Ihre Implementation
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
    // Ihre Implementation
    // ...
}

func (s *MyUserStore) ValidateCredentials(email, password string) (*MyCustomUser, error) {
    // Ihre Implementation
    // ...
}

func (s *MyUserStore) CreateUser(username, email, password string) (*MyCustomUser, error) {
    // Ihre Implementation
    // ...
}

func (s *MyUserStore) UserExists(username, email string) (bool, error) {
    // Ihre Implementation
    // ...
}
```

### 3. Verwenden Sie Ihre Implementation

```go
package main

import (
    "github.com/denkhaus/templ-router/pkg/interfaces"
    "myapp"
)

func main() {
    // Erstellen Sie Ihre User-Store Implementation
    userStore := &myapp.MyUserStore{
        db: myapp.NewDatabase(),
    }
    
    // Verwenden Sie den generischen AuthService
    var authService interfaces.AuthService[*myapp.MyCustomUser]
    
    // Konfigurieren Sie Ihren Router mit der benutzerdefinierten User-Implementation
    // ...
}
```

## Vorteile der generischen Implementation

1. **Flexibilität**: Verwenden Sie Ihre eigene User-Struktur mit beliebigen Feldern
2. **Type Safety**: Go's Typsystem stellt sicher, dass alles korrekt typisiert ist
3. **Keine Konvertierungen**: Arbeiten Sie direkt mit Ihren User-Objekten
4. **Erweiterbarkeit**: Fügen Sie beliebige Felder und Methoden zu Ihrer User-Struktur hinzu

## Migration von der Standard-User-Implementation

Wenn Sie bereits die Standard-`User` Struktur verwenden, müssen Sie nichts ändern. Die Standard-Implementation implementiert bereits das `UserEntity` Interface und funktioniert weiterhin:

```go
// Standard-Implementation funktioniert weiterhin
var userStore interfaces.UserStore[*interfaces.User]
var authService interfaces.AuthService[*interfaces.User]
```

## Beispiel: Enterprise User mit zusätzlichen Feldern

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

// Zusätzliche Methoden für Enterprise-Features
func (u *EnterpriseUser) GetEmployeeID() string { return u.EmployeeID }
func (u *EnterpriseUser) GetCostCenter() string { return u.CostCenter }
func (u *EnterpriseUser) IsManager() bool { return contains(u.Roles, "manager") }
```

Diese generische Lösung macht den Templ-Router viel flexibler und ermöglicht es Entwicklern, ihre eigenen User-Modelle zu verwenden, ohne auf die Funktionalität des Routers verzichten zu müssen.