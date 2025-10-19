package interfaces

import (
	"testing"
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
		settings AuthSettings
		isValid  bool
	}{
		{
			name: "Valid user auth settings",
			settings: AuthSettings{
				Type:        AuthTypeUser,
				RedirectURL: "/login",
				Roles:       []string{"user"},
			},
			isValid: true,
		},
		{
			name: "Valid admin auth settings",
			settings: AuthSettings{
				Type:        AuthTypeAdmin,
				RedirectURL: "/admin/login",
				Roles:       []string{"admin", "super_admin"},
			},
			isValid: true,
		},
		{
			name: "Public auth with no redirect",
			settings: AuthSettings{
				Type:        AuthTypePublic,
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
			if tt.settings.Type == AuthTypeUser || tt.settings.Type == AuthTypeAdmin {
				if tt.settings.RedirectURL == "" {
					isValid = false
				}
			}

			if isValid != tt.isValid {
				t.Errorf("AuthSettings validation = %v, want %v", isValid, tt.isValid)
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

func TestSession_Validation(t *testing.T) {
	session := Session{
		ID:     "session123",
		UserID: "user123",
		Valid:  true,
	}

	if session.ID == "" {
		t.Error("Session should have an ID")
	}
	if session.UserID == "" {
		t.Error("Session should have a UserID")
	}
	if !session.Valid {
		t.Error("Session should be valid in this test")
	}
}

func TestSession_InvalidSession(t *testing.T) {
	session := Session{
		ID:     "session123",
		UserID: "user123",
		Valid:  false,
	}

	if session.Valid {
		t.Error("Session should be invalid in this test")
	}
}

func TestAuthSettings_RoleValidation(t *testing.T) {
	settings := AuthSettings{
		Type:        AuthTypeAdmin,
		RedirectURL: "/admin/login",
		Roles:       []string{"admin", "super_admin"},
	}

	// Helper function to check if settings allow role
	allowsRole := func(s AuthSettings, role string) bool {
		for _, r := range s.Roles {
			if r == role {
				return true
			}
		}
		return false
	}

	if !allowsRole(settings, "admin") {
		t.Error("AuthSettings should allow 'admin' role")
	}
	if !allowsRole(settings, "super_admin") {
		t.Error("AuthSettings should allow 'super_admin' role")
	}
	if allowsRole(settings, "user") {
		t.Error("AuthSettings should not allow 'user' role")
	}
}