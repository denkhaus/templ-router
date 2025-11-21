package interfaces

import (
	"testing"

	"github.com/denkhaus/templ-router/pkg/shared"
)

// TestUser implements UserEntity interface for testing
type TestUser struct {
	ID    string
	Email string
	Roles []string
}

func (u *TestUser) GetID() string {
	return u.ID
}

func (u *TestUser) GetEmail() string {
	return u.Email
}

func (u *TestUser) GetRoles() []string {
	return u.Roles
}

func TestAuthSettings_Validation(t *testing.T) {
	tests := []struct {
		name     string
		settings shared.AuthConfig
		isValid  bool
	}{
		{
			name: "Valid user auth settings",
			settings: shared.AuthConfig{
				Type:        "UserRequired",
				RedirectURL: "/login",
				Roles:       []string{"user"},
			},
			isValid: true,
		},
		{
			name: "Valid admin auth settings",
			settings: shared.AuthConfig{
				Type:        "AdminRequired",
				RedirectURL: "/admin/login",
				Roles:       []string{"admin", "super_admin"},
			},
			isValid: true,
		},
		{
			name: "Public auth with no redirect",
			settings: shared.AuthConfig{
				Type:        "Public",
				RedirectURL: "",
				Roles:       nil,
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation logic
			isValid := true
			if tt.settings.Type == "UserRequired" || tt.settings.Type == "AdminRequired" {
				if tt.settings.RedirectURL == "" {
					isValid = false
				}
			}

			if isValid != tt.isValid {
				t.Errorf("shared.AuthConfig validation = %v, want %v", isValid, tt.isValid)
			}
		})
	}
}

func TestAuthResult_Validation(t *testing.T) {
	tests := []struct {
		name   string
		result AuthResult
		valid  bool
	}{
		{
			name: "Successful authentication",
			result: AuthResult{
				IsAuthenticated: true,
				User: &TestUser{
					ID:    "user123",
					Email: "test@example.com",
					Roles: []string{"user"},
				},
				RedirectURL:  "/dashboard",
				ErrorMessage: "",
			},
			valid: true,
		},
		{
			name: "Failed authentication",
			result: AuthResult{
				IsAuthenticated: false,
				User:            nil,
				RedirectURL:     "/login",
				ErrorMessage:    "Invalid credentials",
			},
			valid: true,
		},
		{
			name: "Authenticated but no user",
			result: AuthResult{
				IsAuthenticated: true,
				User:            nil,
				RedirectURL:     "",
				ErrorMessage:    "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validation: if authenticated, user should be present
			isValid := true
			if tt.result.IsAuthenticated && tt.result.User == nil {
				isValid = false
			}

			if isValid != tt.valid {
				t.Errorf("AuthResult validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestUser_Validation(t *testing.T) {
	user := &TestUser{
		ID:    "user123",
		Email: "test@example.com",
		Roles: []string{"user", "admin"},
	}

	if user.GetID() == "" {
		t.Error("User should have an ID")
	}
	if user.GetEmail() == "" {
		t.Error("User should have an email")
	}
	if len(user.GetRoles()) == 0 {
		t.Error("User should have at least one role")
	}
}

func TestUser_HasRole(t *testing.T) {
	user := &TestUser{
		ID:    "user123",
		Email: "test@example.com",
		Roles: []string{"user", "admin"},
	}

	// Helper function to check if user has role
	hasRole := func(u UserEntity, role string) bool {
		for _, r := range u.GetRoles() {
			if r == role {
				return true
			}
		}
		return false
	}

	if !hasRole(user, "user") {
		t.Error("User should have 'user' role")
	}
	if !hasRole(user, "admin") {
		t.Error("User should have 'admin' role")
	}
	if hasRole(user, "super_admin") {
		t.Error("User should not have 'super_admin' role")
	}
}

func TestAuthSettings_RoleValidation(t *testing.T) {
	settings := shared.AuthConfig{
		Type:        "AdminRequired",
		RedirectURL: "/admin/login",
		Roles:       []string{"admin", "super_admin"},
	}

	// Helper function to check if settings allow role
	allowsRole := func(s shared.AuthConfig, role string) bool {
		for _, r := range s.Roles {
			if r == role {
				return true
			}
		}
		return false
	}

	if !allowsRole(settings, "admin") {
		t.Error("shared.AuthConfig should allow 'admin' role")
	}
	if !allowsRole(settings, "super_admin") {
		t.Error("shared.AuthConfig should allow 'super_admin' role")
	}
	if allowsRole(settings, "user") {
		t.Error("shared.AuthConfig should not allow 'user' role")
	}
}
