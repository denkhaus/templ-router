package services

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// NewAuthService creates a new clean auth service for DI
func NewAuthService(i do.Injector) (interfaces.AuthService, error) {
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	userStore := do.MustInvoke[interfaces.UserStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &CleanAuthService{
		sessionStore: sessionStore,
		userStore:    userStore,
		logger:       logger,
	}, nil
}

// Authenticate implements interfaces.AuthService
func (cas *CleanAuthService) Authenticate(req *http.Request, requirements *interfaces.AuthSettings) (*interfaces.AuthResult, error) {
	if requirements == nil || requirements.Type == interfaces.AuthTypePublic {
		return &interfaces.AuthResult{IsAuthenticated: true}, nil
	}

	// Get session from request
	session, err := cas.sessionStore.GetSession(req)
	if err != nil {
		cas.logger.Debug("No valid session found", zap.Error(err))
		return &interfaces.AuthResult{
			IsAuthenticated: false,
			RedirectURL:     requirements.RedirectURL,
			ErrorMessage:    "Authentication required",
		}, nil
	}

	if !session.Valid {
		cas.logger.Debug("Invalid session", zap.String("session_id", session.ID))
		return &interfaces.AuthResult{
			IsAuthenticated: false,
			RedirectURL:     requirements.RedirectURL,
			ErrorMessage:    "Session expired",
		}, nil
	}

	// Get user from session
	user, err := cas.userStore.GetUserByID(session.UserID)
	if err != nil {
		cas.logger.Error("Failed to get user from session",
			zap.String("user_id", session.UserID),
			zap.Error(err))
		return &interfaces.AuthResult{
			IsAuthenticated: false,
			RedirectURL:     requirements.RedirectURL,
			ErrorMessage:    "User not found",
		}, nil
	}

	// Use interfaces.User directly (no conversion needed)
	interfacesUser := &interfaces.User{
		ID:    user.ID,
		Email: user.Email,
		Roles: user.Roles,
	}

	return &interfaces.AuthResult{
		IsAuthenticated: true,
		User:            interfacesUser,
	}, nil
}

// HasRequiredPermissions implements interfaces.AuthService
func (cas *CleanAuthService) HasRequiredPermissions(req *http.Request, settings *interfaces.AuthSettings) bool {
	if settings == nil || settings.Type == interfaces.AuthTypePublic {
		return true
	}

	session, err := cas.sessionStore.GetSession(req)
	if err != nil || !session.Valid {
		return false
	}

	user, err := cas.userStore.GetUserByID(session.UserID)
	if err != nil {
		return false
	}

	return cas.userHasRequiredRoles(user, settings)
}

// userHasRequiredRoles checks if user has required roles
func (cas *CleanAuthService) userHasRequiredRoles(user *interfaces.User, settings *interfaces.AuthSettings) bool {
	switch settings.Type {
	case interfaces.AuthTypePublic:
		return true
	case interfaces.AuthTypeUser:
		return len(user.Roles) > 0 // Any authenticated user
	case interfaces.AuthTypeAdmin:
		return cas.userHasRole(user, "admin")
	default:
		// Check specific roles if provided
		if len(settings.Roles) == 0 {
			return true // No specific roles required
		}
		for _, requiredRole := range settings.Roles {
			if cas.userHasRole(user, requiredRole) {
				return true
			}
		}
		return false
	}
}

// userHasRole checks if user has a specific role
func (cas *CleanAuthService) userHasRole(user *interfaces.User, role string) bool {
	for _, userRole := range user.Roles {
		if userRole == role {
			return true
		}
	}
	return false
}
