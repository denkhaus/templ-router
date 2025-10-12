package auth

import (
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// authProviderImpl implements AuthProvider (PRIVATE)
type authProviderImpl struct {
	userStore    interfaces.UserStore
	sessionStore interfaces.SessionStore
	logger       *zap.Logger
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string
	Password string
}

// SignupRequest represents a signup request
type SignupRequest struct {
	Username string
	Email    string
	Password string
}

// LoginResult represents the result of a login attempt
type LoginResult struct {
	Success     bool
	User        *interfaces.User
	SessionID   string
	RedirectURL string
	Error       string
}

// SignupResult represents the result of a signup attempt
type SignupResult struct {
	Success     bool
	User        *interfaces.User
	SessionID   string
	RedirectURL string
	Error       string
}

// NewAuthProvider creates a new auth provider for DI
func NewAuthProvider(i do.Injector) (AuthProvider, error) {
	userStore := do.MustInvoke[interfaces.UserStore](i)
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &authProviderImpl{
		userStore:    userStore,
		sessionStore: sessionStore,
		logger:       logger,
	}, nil
}

// ProcessLogin handles the complete login flow
func (ap *authProviderImpl) ProcessLogin(req LoginRequest) *LoginResult {
	ap.logger.Info("Processing login request", zap.String("username", req.Username))

	// Validate credentials
	user, err := ap.userStore.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		ap.logger.Warn("Login failed - invalid credentials",
			zap.String("username", req.Username),
			zap.Error(err))
		return &LoginResult{
			Success: false,
			Error:   "Invalid credentials",
		}
	}

	// Create session
	session, err := ap.sessionStore.CreateSession(user.ID)
	if err != nil {
		ap.logger.Error("Failed to create session",
			zap.String("user_id", user.ID),
			zap.Error(err))
		return &LoginResult{
			Success: false,
			Error:   "Session creation failed",
		}
	}

	// Use interfaces.User directly (no conversion needed)
	interfacesUser := &interfaces.User{
		ID:    user.ID,
		Email: user.Email,
		Roles: user.Roles,
	}

	// Determine redirect URL based on user role
	redirectURL := ap.getRedirectURL(interfacesUser)

	ap.logger.Info("Login successful",
		zap.String("user_id", user.ID),
		zap.String("session_id", session.ID))

	return &LoginResult{
		Success:     true,
		User:        interfacesUser,
		SessionID:   session.ID,
		RedirectURL: redirectURL,
	}
}

// ProcessSignup handles the complete signup flow
func (ap *authProviderImpl) ProcessSignup(req SignupRequest) *SignupResult {
	ap.logger.Info("Processing signup request", zap.String("username", req.Username))

	// Check if user already exists
	exists, err := ap.userStore.UserExists(req.Username, req.Email)
	if err != nil {
		ap.logger.Error("Failed to check user existence",
			zap.String("username", req.Username),
			zap.Error(err))
		return &SignupResult{
			Success: false,
			Error:   "Internal server error",
		}
	}

	if exists {
		ap.logger.Warn("Signup failed - user already exists",
			zap.String("username", req.Username),
			zap.String("email", req.Email))
		return &SignupResult{
			Success: false,
			Error:   "Username or email already exists",
		}
	}

	// Create new user
	user, err := ap.userStore.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		ap.logger.Error("Failed to create user",
			zap.String("username", req.Username),
			zap.Error(err))
		return &SignupResult{
			Success: false,
			Error:   "Failed to create user account",
		}
	}

	// Create session for auto-login
	session, err := ap.sessionStore.CreateSession(user.ID)
	if err != nil {
		ap.logger.Error("Failed to create session after signup",
			zap.String("user_id", user.ID),
			zap.Error(err))
		return &SignupResult{
			Success: false,
			Error:   "Account created but login failed",
		}
	}

	// Use interfaces.User directly (no conversion needed)
	interfacesUser := &interfaces.User{
		ID:    user.ID,
		Email: user.Email,
		Roles: user.Roles,
	}

	// Determine redirect URL
	redirectURL := ap.getRedirectURL(interfacesUser)

	ap.logger.Info("Signup successful",
		zap.String("user_id", user.ID),
		zap.String("session_id", session.ID))

	return &SignupResult{
		Success:     true,
		User:        interfacesUser,
		SessionID:   session.ID,
		RedirectURL: redirectURL,
	}
}

// ProcessLogout handles the complete logout flow
func (ap *authProviderImpl) ProcessLogout(r *http.Request) error {
	ap.logger.Info("Processing logout request")

	// Get session from request
	session, err := ap.sessionStore.GetSession(r)
	if err != nil {
		ap.logger.Debug("No session found during logout", zap.Error(err))
		return nil // Not an error if no session exists
	}

	// Delete session
	if err := ap.sessionStore.DeleteSession(session.ID); err != nil {
		ap.logger.Warn("Failed to delete session during logout",
			zap.String("session_id", session.ID),
			zap.Error(err))
		return err
	}

	ap.logger.Info("Logout successful", zap.String("session_id", session.ID))
	return nil
}

// SetSessionCookie sets the session cookie in the HTTP response
func (ap *authProviderImpl) SetSessionCookie(w http.ResponseWriter, sessionID string) {
	ap.logger.Debug("Setting session cookie", zap.String("session_id", sessionID))

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		MaxAge:   3600,  // 1 hour
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

// ClearSessionCookie clears the session cookie from the HTTP response
func (ap *authProviderImpl) ClearSessionCookie(w http.ResponseWriter) {
	ap.logger.Debug("Clearing session cookie")

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1, // Delete immediately
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

// getRedirectURL determines the post-login redirect URL based on user roles
func (ap *authProviderImpl) getRedirectURL(user *interfaces.User) string {
	// Check if user has admin role
	for _, role := range user.Roles {
		if role == "admin" {
			return "/en/admin"
		}
	}

	// Default to dashboard for regular users
	return "/en/dashboard"
}
