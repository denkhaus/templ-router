package interfaces

import (
	"net/http"
	"time"

	"github.com/denkhaus/templ-router/pkg/interfaces"
)

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Valid     bool      `json:"valid"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionStore interface for session management
type SessionStore interface {
	GetSession(req *http.Request) (*Session, error)
	CreateSession(userID string) (*Session, error)
	DeleteSession(sessionID string) error

	GetSessionByID(sessionID string) (*Session, error)            // Direct access
	ExtendSession(sessionID string, duration time.Duration) error // Session verlängern
}

// UserStore interface for user management
type UserStore interface {
	GetUserByID(userID string) (interfaces.UserEntity, error)

	// Request-based methods for complete data extraction and validation
	ValidateCredentialsFromRequest(req *http.Request) (interfaces.UserEntity, error)
	CreateUserFromRequest(req *http.Request) (interfaces.UserEntity, error)
}

// AuthHandlers defines the interface for authentication HTTP handlers
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