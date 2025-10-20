package middleware

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

const (
	DefaultSigninRoute = "/login"
)

// authMiddleware handles authentication concerns separately (private implementation)
type authMiddleware struct {
	authService   interfaces.AuthService
	configService interfaces.ConfigService
	logger        *zap.Logger
}

// AuthService interface for clean dependency

// Import central auth types to eliminate redundancy
// AuthResult, AuthSettings, AuthType, User are now imported from interfaces package

// AuthType methods moved to interfaces package

// NewAuthMiddleware creates a new auth middleware for DI
func NewAuthMiddleware(i do.Injector) (AuthMiddlewareInterface, error) {
	authService := do.MustInvoke[interfaces.AuthService](i)
	configService := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &authMiddleware{
		authService:   authService,
		configService: configService,
		logger:        logger,
	}, nil
}

// Handle processes authentication for a request
func (am *authMiddleware) Handle(next http.Handler, requirements *interfaces.AuthSettings) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public routes
		if requirements == nil || requirements.Type == interfaces.AuthTypePublic {
			next.ServeHTTP(w, r)
			return
		}

		// Authenticate the request
		authResult, err := am.authService.Authenticate(r, requirements)
		if err != nil {
			am.logger.Error("Authentication error",
				zap.String("path", r.URL.Path),
				zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Handle authentication failure
		if !authResult.IsAuthenticated {
			am.handleAuthFailure(w, r, authResult, requirements)
			return
		}

		// Check permissions
		if !am.authService.HasRequiredPermissions(r, requirements) {
			am.handlePermissionFailure(w, r, requirements)
			return
		}

		// Authentication successful, continue to next handler
		next.ServeHTTP(w, r)
	})
}

// handleAuthFailure handles authentication failures
func (am *authMiddleware) handleAuthFailure(w http.ResponseWriter, r *http.Request, authResult *interfaces.AuthResult, requirements *interfaces.AuthSettings) {
	am.logger.Info("Authentication required but user not authenticated",
		zap.String("path", r.URL.Path),
		zap.String("auth_type", requirements.Type.String()))

	// Determine redirect URL
	var redirectURL string
	if authResult.RedirectURL != "" {
		redirectURL = authResult.RedirectURL
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
func (am *authMiddleware) handlePermissionFailure(w http.ResponseWriter, r *http.Request, requirements *interfaces.AuthSettings) {
	am.logger.Warn("User lacks required permissions",
		zap.String("path", r.URL.Path),
		zap.String("required_auth_type", requirements.Type.String()))

	if requirements.RedirectURL != "" {
		http.Redirect(w, r, requirements.RedirectURL, http.StatusFound)
	} else {
		am.logger.Warn("Auth-required page has no redirect_url configured",
			zap.String("path", r.URL.Path),
			zap.String("auth_type", requirements.Type.String()))
		http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
	}
}
