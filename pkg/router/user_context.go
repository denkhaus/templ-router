package router

import (
	"context"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
)

// User context keys - using shared package to avoid import cycles
// Deprecated: Use shared.UserContextKey instead

// GetCurrentUser retrieves the current user from context
// This is a generic utility function that works with any UserEntity implementation
func GetCurrentUser[T interfaces.UserEntity](ctx context.Context) T {
	var zero T
	if user, ok := ctx.Value(shared.UserContextKey).(T); ok {
		return user
	}
	return zero
}

// HasUser checks if a user is present in the context
func HasUser(ctx context.Context) bool {
	return ctx.Value(shared.UserContextKey) != nil
}

// GetUserID retrieves the current user's ID from context
func GetUserID[T interfaces.UserEntity](ctx context.Context) string {
	user := GetCurrentUser[T](ctx)
	var zero T
	if any(user) != any(zero) {
		return user.GetID()
	}
	return ""
}

// GetUserEmail retrieves the current user's email from context
func GetUserEmail[T interfaces.UserEntity](ctx context.Context) string {
	user := GetCurrentUser[T](ctx)
	var zero T
	if any(user) != any(zero) {
		return user.GetEmail()
	}
	return ""
}

// GetUserRoles retrieves the current user's roles from context
func GetUserRoles[T interfaces.UserEntity](ctx context.Context) []string {
	user := GetCurrentUser[T](ctx)
	var zero T
	if any(user) != any(zero) {
		return user.GetRoles()
	}
	return nil
}

// UserHasRole checks if the current user has a specific role
func UserHasRole[T interfaces.UserEntity](ctx context.Context, role string) bool {
	roles := GetUserRoles[T](ctx)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// SetUserInContext sets a user in the context
// This is set in auth middleware
func SetUserInContext[T interfaces.UserEntity](ctx context.Context, user T) context.Context {
	return context.WithValue(ctx, shared.UserContextKey, user)
}
