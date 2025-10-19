package interfaces

import "net/http"

// AuthHandlers defines the interface for authentication HTTP handlers
// This allows users to provide their own authentication implementation
type AuthHandlers interface {
	// RegisterRoutes registers authentication routes with the provided registration function
	RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc))
	
	// HandleLogin handles user login requests
	HandleLogin(w http.ResponseWriter, r *http.Request)
	
	// HandleSignup handles user signup requests  
	HandleSignup(w http.ResponseWriter, r *http.Request)
	
	// HandleLogout handles user logout requests
	HandleLogout(w http.ResponseWriter, r *http.Request)
}