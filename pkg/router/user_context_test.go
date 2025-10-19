package router

import (
	"context"
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
)

// testUserEntity implements UserEntity for testing
type testUserEntity struct {
	id    string
	email string
	roles []string
}

func (u *testUserEntity) GetID() string    { return u.id }
func (u *testUserEntity) GetEmail() string { return u.email }
func (u *testUserEntity) GetRoles() []string { return u.roles }

func TestGetCurrentUser(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected interfaces.UserEntity
	}{
		{
			name:     "No user in context",
			ctx:      context.Background(),
			expected: nil,
		},
		{
			name: "User in context",
			ctx: context.WithValue(context.Background(), UserContextKey, &testUserEntity{
				id:    "user123",
				email: "test@example.com",
				roles: []string{"user"},
			}),
			expected: &testUserEntity{
				id:    "user123",
				email: "test@example.com",
				roles: []string{"user"},
			},
		},
		{
			name:     "Wrong type in context",
			ctx:      context.WithValue(context.Background(), UserContextKey, "not-a-user"),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCurrentUser(tt.ctx)
			
			if tt.expected == nil {
				if result != nil {
					t.Errorf("getCurrentUser() = %v, want nil", result)
				}
				return
			}
			
			if result == nil {
				t.Errorf("getCurrentUser() = nil, want %v", tt.expected)
				return
			}
			
			if result.GetID() != tt.expected.GetID() {
				t.Errorf("getCurrentUser().GetID() = %v, want %v", result.GetID(), tt.expected.GetID())
			}
			
			if result.GetEmail() != tt.expected.GetEmail() {
				t.Errorf("getCurrentUser().GetEmail() = %v, want %v", result.GetEmail(), tt.expected.GetEmail())
			}
		})
	}
}

func TestSetUserInContext(t *testing.T) {
	user := &testUserEntity{
		id:    "user123",
		email: "test@example.com",
		roles: []string{"user"},
	}
	
	ctx := setUserInContext(context.Background(), user)
	result := GetCurrentUser(ctx)
	
	if result == nil {
		t.Fatal("setUserInContext() did not set user in context")
	}
	
	if result.GetID() != user.GetID() {
		t.Errorf("setUserInContext() user ID = %v, want %v", result.GetID(), user.GetID())
	}
}

func TestHasUser(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected bool
	}{
		{
			name:     "No user in context",
			ctx:      context.Background(),
			expected: false,
		},
		{
			name: "User in context",
			ctx: context.WithValue(context.Background(), UserContextKey, &testUserEntity{
				id: "user123",
			}),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasUser(tt.ctx)
			if result != tt.expected {
				t.Errorf("hasUser() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "No user in context",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name: "User in context",
			ctx: context.WithValue(context.Background(), UserContextKey, &testUserEntity{
				id: "user123",
			}),
			expected: "user123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getUserID(tt.ctx)
			if result != tt.expected {
				t.Errorf("getUserID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetUserEmail(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "No user in context",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name: "User in context",
			ctx: context.WithValue(context.Background(), UserContextKey, &testUserEntity{
				email: "test@example.com",
			}),
			expected: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getUserEmail(tt.ctx)
			if result != tt.expected {
				t.Errorf("getUserEmail() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetUserRoles(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected []string
	}{
		{
			name:     "No user in context",
			ctx:      context.Background(),
			expected: nil,
		},
		{
			name: "User in context",
			ctx: context.WithValue(context.Background(), UserContextKey, &testUserEntity{
				roles: []string{"user", "admin"},
			}),
			expected: []string{"user", "admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getUserRoles(tt.ctx)
			
			if tt.expected == nil {
				if result != nil {
					t.Errorf("getUserRoles() = %v, want nil", result)
				}
				return
			}
			
			if len(result) != len(tt.expected) {
				t.Errorf("getUserRoles() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			
			for i, role := range result {
				if role != tt.expected[i] {
					t.Errorf("getUserRoles()[%d] = %v, want %v", i, role, tt.expected[i])
				}
			}
		})
	}
}

func TestHasRole(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		role     string
		expected bool
	}{
		{
			name:     "No user in context",
			ctx:      context.Background(),
			role:     "admin",
			expected: false,
		},
		{
			name: "User has role",
			ctx: context.WithValue(context.Background(), UserContextKey, &testUserEntity{
				roles: []string{"user", "admin"},
			}),
			role:     "admin",
			expected: true,
		},
		{
			name: "User does not have role",
			ctx: context.WithValue(context.Background(), UserContextKey, &testUserEntity{
				roles: []string{"user"},
			}),
			role:     "admin",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasRole(tt.ctx, tt.role)
			if result != tt.expected {
				t.Errorf("hasRole() = %v, want %v", result, tt.expected)
			}
		})
	}
}