package auth

import "net/http"

// AuthProvider interface for external authentication functionality (PUBLIC)
// This is a pluggable service that can be injected into any system
type AuthProvider interface {
	ProcessLogin(req LoginRequest) *LoginResult
	ProcessSignup(req SignupRequest) *SignupResult
	ProcessLogout(r *http.Request) error
	SetSessionCookie(w http.ResponseWriter, sessionID string)
	ClearSessionCookie(w http.ResponseWriter)
}

// AuthHandlersInterface defines the contract for authentication HTTP handlers
type AuthHandlersInterface interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleSignup(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc))
}
