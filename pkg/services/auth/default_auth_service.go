package auth

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// DefaultAuthService provides a default authentication service implementation
// Users can replace this with their own implementation (OAuth, LDAP, etc.)
type DefaultAuthService struct {
	userStore    interfaces.UserStore
	sessionStore interfaces.SessionStore
	logger       *zap.Logger
}

// NewDefaultAuthService creates a new default auth service for DI
func NewDefaultAuthService(i do.Injector) (interfaces.AuthService, error) {
	userStore := do.MustInvoke[interfaces.UserStore](i)
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &DefaultAuthService{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}, nil
}

// Authenticate implements the AuthService interface
func (s *DefaultAuthService) Authenticate(req *http.Request, requirements *interfaces.AuthSettings) (*interfaces.AuthResult, error) {
	// Check authentication type
	switch requirements.Type {
	case interfaces.AuthTypePublic:
		return &interfaces.AuthResult{
			IsAuthenticated: true,
		}, nil

	case interfaces.AuthTypeUser, interfaces.AuthTypeAdmin:
		return s.authenticateUser(req, requirements)

	default:
		return &interfaces.AuthResult{
			IsAuthenticated: false,
			ErrorMessage:    "Unknown authentication type",
		}, nil
	}
}

// authenticateUser handles user authentication via session
func (s *DefaultAuthService) authenticateUser(req *http.Request, requirements *interfaces.AuthSettings) (*interfaces.AuthResult, error) {
	// Get session from request
	session, err := s.sessionStore.GetSession(req)
	if err != nil || !session.Valid {
		return &interfaces.AuthResult{
			IsAuthenticated: false,
			RedirectURL:     requirements.RedirectURL,
		}, nil
	}

	// Get user from session
	user, err := s.userStore.GetUserByID(session.UserID)
	if err != nil {
		return &interfaces.AuthResult{
			IsAuthenticated: false,
			RedirectURL:     requirements.RedirectURL,
			ErrorMessage:    "User not found",
		}, nil
	}

	// Check role requirements
	if !s.hasRequiredRoles(user, requirements.Roles) {
		return &interfaces.AuthResult{
			IsAuthenticated: true,
			User:            user,
			RedirectURL:     "/unauthorized",
			ErrorMessage:    "Insufficient permissions",
		}, nil
	}

	return &interfaces.AuthResult{
		IsAuthenticated: true,
		User:            user,
	}, nil
}

// HasRequiredPermissions checks if the user has the required permissions
func (s *DefaultAuthService) HasRequiredPermissions(req *http.Request, settings *interfaces.AuthSettings) bool {
	result, err := s.Authenticate(req, settings)
	if err != nil || !result.IsAuthenticated {
		return false
	}

	return s.hasRequiredRoles(result.User, settings.Roles)
}

// hasRequiredRoles checks if user has any of the required roles
func (s *DefaultAuthService) hasRequiredRoles(user interfaces.UserEntity, requiredRoles []string) bool {
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

// Login provides a convenience method for user login
func (s *DefaultAuthService) Login(email, password string) (interfaces.UserEntity, string, error) {
	user, err := s.userStore.ValidateCredentials(email, password)
	if err != nil {
		return nil, "", err
	}

	// Create session
	session, err := s.sessionStore.CreateSession(user.GetID())
	if err != nil {
		return nil, "", err
	}

	s.logger.Info("User logged in successfully",
		zap.String("user_id", user.GetID()),
		zap.String("session_id", session.ID))

	return user, session.ID, nil
}

// Logout provides a convenience method for user logout
func (s *DefaultAuthService) Logout(sessionID string) error {
	err := s.sessionStore.DeleteSession(sessionID)
	if err != nil {
		s.logger.Error("Failed to delete session", zap.Error(err))
		return err
	}

	s.logger.Info("User logged out successfully", zap.String("session_id", sessionID))
	return nil
}
