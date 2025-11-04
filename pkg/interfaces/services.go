package interfaces

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AssetsService interface {
	SetupRoutesWithRouter(mux chi.Router)
	SetupRoutes(mux *chi.Mux)
}

// ValidationService handles unified validation of routes and configurations
type ValidationService interface {
	ValidateConfiguration(routes []Route, configs map[string]*ConfigFile) error
}

// SessionStore interface for session management (pluggable)
type SessionStore interface {
	GetSession(req *http.Request) (*Session, error)
	CreateSession(userID string) (*Session, error)
	DeleteSession(sessionID string) error
}

// UserEntity defines the minimal interface that any user implementation must satisfy
type UserEntity interface {
	GetID() string
	GetEmail() string
	GetRoles() []string
}

// UserStore interface for user management (pluggable and generic)
type UserStore interface {
	GetUserByID(userID string) (UserEntity, error)
	GetUserByEmail(email string) (UserEntity, error)
	ValidateCredentials(email, password string) (UserEntity, error)
	CreateUser(username, email, password string) (UserEntity, error)
	UserExists(username, email string) (bool, error)

	// Request-based methods for complete data extraction and validation
	ValidateCredentialsFromRequest(req *http.Request) (UserEntity, error)
	CreateUserFromRequest(req *http.Request) (UserEntity, error)
}
