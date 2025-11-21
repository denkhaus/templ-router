package middleware

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

const (
	DefaultSigninRoute = "/login"
)

// authMiddleware handles authentication concerns using AuthValidator hooks
type authMiddleware struct {
	authValidator interfaces.AuthValidator
	configService interfaces.ConfigService
	logger        *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware for DI using AuthValidator
func NewAuthMiddleware(i do.Injector) (interfaces.AuthMiddlewareInterface, error) {
	authValidator := do.MustInvoke[interfaces.AuthValidator](i)
	configService := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &authMiddleware{
		authValidator: authValidator,
		configService: configService,
		logger:        logger,
	}, nil
}

// Handle processes authentication for a request
func (am *authMiddleware) Handle(next http.Handler, requirements *shared.AuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public routes
		if requirements == nil || requirements.Type == "Public" {
			next.ServeHTTP(w, r)
			return
		}

		// Check if user is authenticated
		authenticated, err := am.authValidator.IsAuthenticated(r)
		if err != nil {
			am.logger.Error("Failed to get authentication status",
				zap.String("path", r.URL.Path),
				zap.Error(err))
			am.handleAuthFailure(w, r, requirements)
			return
		}

		if !authenticated {
			am.handleAuthFailure(w, r, requirements)
			return
		}

		// Get current user
		user, err := am.authValidator.GetCurrentUser(r)
		if err != nil {
			am.logger.Error("Failed to get current user",
				zap.String("path", r.URL.Path),
				zap.Error(err))
			am.handleAuthFailure(w, r, requirements)
			return
		}

		// Check role-based permissions
		if !am.authValidator.HasRole(user, requirements.Roles) {
			am.handlePermissionFailure(w, r, requirements)
			return
		}

		// Authentication successful, continue to next handler
		next.ServeHTTP(w, r)
	})
}

// handleAuthFailure handles authentication failures
func (am *authMiddleware) handleAuthFailure(w http.ResponseWriter, r *http.Request, requirements *shared.AuthConfig) {
	am.logger.Info("Authentication required but user not authenticated",
		zap.String("path", r.URL.Path),
		zap.String("auth_type", requirements.Type))

	// Determine redirect URL
	var redirectURL string
	if requirements.RedirectURL != "" {
		redirectURL = requirements.RedirectURL
	} else {
		// Default to signin route from config
		redirectURL = DefaultSigninRoute
		if signinRoute := am.configService.GetSignInRoute(); signinRoute != "" {
			redirectURL = i18n.LocalizeRouteIfRequired(r.Context(), signinRoute)
		}
	}

	// Add return URL parameter so user can be redirected back after login
	if redirectURL != "" {
		if r.URL.RawQuery != "" {
			redirectURL += "?return_to=" + r.URL.Path + "?" + r.URL.RawQuery
		} else {
			redirectURL += "?return_to=" + r.URL.Path
		}

		am.logger.Info("Redirecting unauthenticated user to signin",
			zap.String("original_path", r.URL.Path),
			zap.String("redirect_url", redirectURL))

		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		am.logger.Warn("No signin route configured, falling back to error response",
			zap.String("path", r.URL.Path))
		http.Error(w, "Authentication required", http.StatusUnauthorized)
	}
}

// handlePermissionFailure handles permission failures
func (am *authMiddleware) handlePermissionFailure(w http.ResponseWriter, r *http.Request, requirements *shared.AuthConfig) {
	am.logger.Warn("User lacks required permissions",
		zap.String("path", r.URL.Path),
		zap.String("required_auth_type", requirements.Type))

	if requirements.RedirectURL != "" {
		http.Redirect(w, r, requirements.RedirectURL, http.StatusFound)
	} else {
		am.logger.Warn("Auth-required page has no redirect_url configured",
			zap.String("path", r.URL.Path),
			zap.String("auth_type", requirements.Type))
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
	}
}
