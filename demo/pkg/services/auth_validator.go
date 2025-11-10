package services

import (
	"net/http"

	"github.com/denkhaus/templ-router/demo/pkg/interfaces"
	sharedInterfaces "github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// demoAuthValidator implements the AuthValidator interface for the demo application
// This shows how to create a custom authentication validator using session and user stores
type demoAuthValidator struct {
	userStore    interfaces.UserStore
	sessionStore interfaces.SessionStore
	logger       *zap.Logger
}

// NewDemoAuthValidator creates a new demo auth validator for DI
func NewDemoAuthValidator(i do.Injector) (sharedInterfaces.AuthValidator, error) {
	userStore := do.MustInvoke[interfaces.UserStore](i)
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &demoAuthValidator{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}, nil
}

// IsAuthenticated checks if the current request is from an authenticated user
func (av *demoAuthValidator) IsAuthenticated(req *http.Request) bool {
	// Get session from request
	session, err := av.sessionStore.GetSession(req)
	if err != nil || !session.Valid {
		return false
	}

	// Verify user exists
	_, err = av.userStore.GetUserByID(session.UserID)
	if err != nil {
		av.logger.Warn("Session references non-existent user",
			zap.String("session_id", session.ID),
			zap.String("user_id", session.UserID),
			zap.Error(err))
		return false
	}

	return true
}

// GetCurrentUser returns the authenticated user for the current request
func (av *demoAuthValidator) GetCurrentUser(req *http.Request) (sharedInterfaces.UserEntity, error) {
	// Get session from request
	session, err := av.sessionStore.GetSession(req)
	if err != nil {
		return nil, err
	}

	if !session.Valid {
		return nil, http.ErrNoCookie // Invalid session
	}

	// Get user from session
	user, err := av.userStore.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// HasRole checks if the given user has any of the required roles
func (av *demoAuthValidator) HasRole(user sharedInterfaces.UserEntity, requiredRoles []string) bool {
	if len(requiredRoles) == 0 {
		return true // No specific roles required
	}

	userRoles := user.GetRoles()
	for _, required := range requiredRoles {
		for _, userRole := range userRoles {
			if userRole == required {
				return true
			}
		}
	}
	return false
}
