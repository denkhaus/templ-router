package router

import (
	"context"

	"github.com/denkhaus/templ-router/pkg/interfaces"
)

// User context keys
const (
	UserContextKey = "user"
)

// GetCurrentUser retrieves the current user from context
// This is a generic utility function that works with any UserEntity implementation
func GetCurrentUser(ctx context.Context) interfaces.UserEntity {
	if user, ok := ctx.Value(UserContextKey).(interfaces.UserEntity); ok {
		return user
	}
	return nil
}

// getCurrentUser is an internal alias for backward compatibility
func getCurrentUser(ctx context.Context) interfaces.UserEntity {
	return GetCurrentUser(ctx)
}

// setUserInContext sets a user in the context
// This is a helper function for middleware to set the authenticated user
func setUserInContext(ctx context.Context, user interfaces.UserEntity) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// HasUser checks if a user is present in the context
func HasUser(ctx context.Context) bool {
	return GetCurrentUser(ctx) != nil
}

// hasUser is an internal alias for backward compatibility
func hasUser(ctx context.Context) bool {
	return HasUser(ctx)
}

// getUserID retrieves the current user's ID from context
func getUserID(ctx context.Context) string {
	if user := getCurrentUser(ctx); user != nil {
		return user.GetID()
	}
	return ""
}

// getUserEmail retrieves the current user's email from context
func getUserEmail(ctx context.Context) string {
	if user := getCurrentUser(ctx); user != nil {
		return user.GetEmail()
	}
	return ""
}

// getUserRoles retrieves the current user's roles from context
func getUserRoles(ctx context.Context) []string {
	if user := getCurrentUser(ctx); user != nil {
		return user.GetRoles()
	}
	return nil
}

// hasRole checks if the current user has a specific role
func hasRole(ctx context.Context, role string) bool {
	roles := getUserRoles(ctx)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}