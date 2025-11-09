package interfaces

import (
	"net/http"
	"time"

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

// SessionStore interface for session management (pluggable)
type SessionStore interface {
	GetSession(req *http.Request) (*Session, error)
	CreateSession(userID string) (*Session, error)
	DeleteSession(sessionID string) error

	GetSessionByID(sessionID string) (*Session, error)            // Direct access
	ExtendSession(sessionID string, duration time.Duration) error // Session verlängern
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

	// Request-based methods for complete data extraction and validation
	ValidateCredentialsFromRequest(req *http.Request) (UserEntity, error)
	CreateUserFromRequest(req *http.Request) (UserEntity, error)
}
