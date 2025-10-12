package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// authHandlers provides HTTP handlers for authentication routes (private implementation)
type authHandlers struct {
	authProvider AuthProvider // Use interface, not pointer to interface
	logger       *zap.Logger
}

// NewAuthHandlers creates new auth handlers for DI
func NewAuthHandlers(i do.Injector) (AuthHandlersInterface, error) {
	authProvider := do.MustInvoke[AuthProvider](i)
	logger := do.MustInvoke[*zap.Logger](i)
	
	return &authHandlers{
		authProvider: authProvider,
		logger:       logger,
	}, nil
}

// HandleLogin processes login form submissions
func (ah *authHandlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ah.logger.Info("Processing login request")
	
	// Parse form data
	if err := r.ParseForm(); err != nil {
		ah.logger.Error("Failed to parse login form", zap.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	
	// Extract credentials
	loginReq := LoginRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	
	// Process login through auth provider
	result := ah.authProvider.ProcessLogin(loginReq)
	
	if !result.Success {
		ah.logger.Warn("Login failed", zap.String("error", result.Error))
		ah.redirectWithError(w, r, "/login", result.Error)
		return
	}
	
	// Set session cookie
	ah.authProvider.SetSessionCookie(w, result.SessionID)
	
	// Redirect to appropriate page
	http.Redirect(w, r, result.RedirectURL, http.StatusFound)
}

// HandleSignup processes signup form submissions
func (ah *authHandlers) HandleSignup(w http.ResponseWriter, r *http.Request) {
	ah.logger.Info("Processing signup request")
	
	// Parse form data
	if err := r.ParseForm(); err != nil {
		ah.logger.Error("Failed to parse signup form", zap.Error(err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	
	// Extract signup data
	signupReq := SignupRequest{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	
	// Validate signup data
	if err := ah.validateSignupRequest(signupReq); err != nil {
		ah.logger.Warn("Signup validation failed", zap.Error(err))
		ah.redirectWithError(w, r, "/signup", err.Error())
		return
	}
	
	// Process signup through auth provider
	result := ah.authProvider.ProcessSignup(signupReq)
	
	if !result.Success {
		ah.logger.Warn("Signup failed", zap.String("error", result.Error))
		ah.redirectWithError(w, r, "/signup", result.Error)
		return
	}
	
	// Set session cookie for auto-login after signup
	ah.authProvider.SetSessionCookie(w, result.SessionID)
	
	// Redirect to dashboard or welcome page
	redirectURL := result.RedirectURL
	if redirectURL == "" {
		redirectURL = "/en/dashboard" // Default redirect after signup
	}
	
	ah.logger.Info("Signup successful", 
		zap.String("username", signupReq.Username),
		zap.String("redirect_url", redirectURL))
	
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// HandleLogout processes logout requests
func (ah *authHandlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ah.logger.Info("Processing logout request")
	
	// Process logout through auth provider
	if err := ah.authProvider.ProcessLogout(r); err != nil {
		ah.logger.Error("Logout failed", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	// Clear session cookie
	ah.authProvider.ClearSessionCookie(w)
	
	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusFound)
}

// redirectWithError redirects with an error parameter
func (ah *authHandlers) redirectWithError(w http.ResponseWriter, r *http.Request, url, errorMsg string) {
	redirectURL := url + "?error=" + errorMsg
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// validateSignupRequest validates signup form data
func (ah *authHandlers) validateSignupRequest(req SignupRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !ah.isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

// isValidEmail performs basic email validation
func (ah *authHandlers) isValidEmail(email string) bool {
	// Basic email validation (in production, use proper regex or library)
	return len(email) > 5 && 
		   strings.Contains(email, "@") && 
		   strings.Contains(email, ".") &&
		   !strings.HasPrefix(email, "@") &&
		   !strings.HasSuffix(email, "@")
}

// RegisterRoutes registers auth routes with a router (pluggable interface)
func (ah *authHandlers) RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc)) {
	ah.logger.Info("Registering auth routes")
	
	registerFunc("POST", "/login", ah.HandleLogin)
	registerFunc("POST", "/signup", ah.HandleSignup)
	registerFunc("POST", "/logout", ah.HandleLogout)
	
	ah.logger.Info("Auth routes registered successfully")
}