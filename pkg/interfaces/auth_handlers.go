package interfaces

import "net/http"

// AuthHandlers defines the interface for authentication HTTP handlers
// This allows users to provide their own authentication implementation
type AuthHandlers interface {
	// RegisterRoutes registers authentication routes with the provided registration function
	RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc))

	// HandleSignIn handles user login requests
	HandleSignIn(w http.ResponseWriter, r *http.Request)

	// HandleSignUp handles user signup requests
	HandleSignUp(w http.ResponseWriter, r *http.Request)

	// HandleSignOut handles user logout requests
	HandleSignOut(w http.ResponseWriter, r *http.Request)
}
