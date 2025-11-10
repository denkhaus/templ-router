package interfaces

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/go-chi/chi/v5"
)

type AssetsService interface {
	SetupRoutesWithRouter(mux chi.Router)
	SetupRoutes(mux *chi.Mux)
}

// ValidationService handles unified validation of routes and configurations
type ValidationService interface {
	ValidateConfiguration(routes []Route, configs map[string]*shared.ConfigFile) error
}

// UserEntity defines the minimal interface that any user implementation must satisfy
type UserEntity interface {
	GetID() string
	GetEmail() string
	GetRoles() []string
}

// AuthValidator defines the authentication hook interface that client applications must implement
// This replaces the concrete AuthService implementation with a flexible hook-based system
type AuthValidator interface {
	// IsAuthenticated checks if the current request is from an authenticated user
	IsAuthenticated(req *http.Request) bool

	// GetCurrentUser returns the authenticated user for the current request
	// Returns error if user cannot be retrieved or is not authenticated
	GetCurrentUser(req *http.Request) (UserEntity, error)

	// HasRole checks if the given user has any of the required roles
	// Returns true if user has at least one of the required roles
	HasRole(user UserEntity, requiredRoles []string) bool
}
