package auth

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// AuthService provides a default authentication service implementation
// Users can replace this with their own implementation (OAuth, LDAP, etc.)
type AuthService struct {
	userStore     interfaces.UserStore
	sessionStore  interfaces.SessionStore
	configService interfaces.ConfigService
	signInRoute   string
	logger        *zap.Logger
}

// NewAuthService creates a new auth service for DI
func NewAuthService(i do.Injector) (interfaces.AuthService, error) {
	userStore := do.MustInvoke[interfaces.UserStore](i)
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	configService := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	signInRoute := configService.GetSignInRoute()

	return &AuthService{
		signInRoute:   signInRoute,
		userStore:     userStore,
		sessionStore:  sessionStore,
		configService: configService,
		logger:        logger,
	}, nil
}

// Authenticate implements the AuthService interface
func (s *AuthService) Authenticate(req *http.Request, requirements *shared.AuthConfig) (*interfaces.AuthResult, error) {
	// Check authentication type
	switch requirements.Type {
	case "Public":
		return &interfaces.AuthResult{
			IsAuthenticated: true,
		}, nil

	case "UserRequired", "AdminRequired":
		return s.authenticateUser(req, requirements)

	default:
		return &interfaces.AuthResult{
			IsAuthenticated: false,
			ErrorMessage:    "Unknown authentication type",
		}, nil
	}
}

// authenticateUser handles user authentication via session
func (s *AuthService) authenticateUser(req *http.Request, requirements *shared.AuthConfig) (*interfaces.AuthResult, error) {
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
		signInRoute := i18n.LocalizeRouteIfRequired(req.Context(), s.signInRoute)

		return &interfaces.AuthResult{
			IsAuthenticated: true,
			User:            user,
			RedirectURL:     signInRoute,
			ErrorMessage:    "Insufficient permissions",
		}, nil
	}

	return &interfaces.AuthResult{
		IsAuthenticated: true,
		User:            user,
	}, nil
}

// HasRequiredPermissions checks if the user has the required permissions
func (s *AuthService) HasRequiredPermissions(req *http.Request, settings *shared.AuthConfig) bool {
	result, err := s.Authenticate(req, settings)
	if err != nil || !result.IsAuthenticated {
		return false
	}

	return s.hasRequiredRoles(result.User, settings.Roles)
}

// hasRequiredRoles checks if user has any of the required roles
func (s *AuthService) hasRequiredRoles(user interfaces.UserEntity, requiredRoles []string) bool {
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
