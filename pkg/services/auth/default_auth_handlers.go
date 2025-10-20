package auth

import (
	"encoding/json"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// authHandlersImpl provides generic authentication API handlers
// Works with any UserEntity implementation through the UserStore interface
type authHandlersImpl struct {
	userStore     interfaces.UserStore
	sessionStore  interfaces.SessionStore
	configService interfaces.ConfigService
	logger        *zap.Logger
}

// NewAuthHandlers creates new generic auth handlers
func NewAuthHandlers(i do.Injector) (interfaces.AuthHandlers, error) {
	userStore := do.MustInvoke[interfaces.UserStore](i)
	configService := do.MustInvoke[interfaces.ConfigService](i)
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &authHandlersImpl{
		userStore:     userStore,
		configService: configService,
		sessionStore:  sessionStore,
		logger:        logger,
	}, nil
}

// RegisterRoutes registers authentication API routes only
func (h *authHandlersImpl) RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc)) {
	registerFunc("POST", "/api/auth/signin", h.HandleSignIn)
	registerFunc("POST", "/api/auth/signup", h.HandleSignUp)
	registerFunc("POST", "/api/auth/signout", h.HandleSignOut)
}

// HandleLogin handles user login API endpoint
// UserStore extracts and validates all relevant data from the request
func (h *authHandlersImpl) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// UserStore handles complete data extraction and validation from request
	user, err := h.userStore.ValidateCredentialsFromRequest(r)
	if err != nil {
		h.logger.Warn("Login failed", zap.Error(err))
		h.respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session, err := h.sessionStore.CreateSession(user.GetID())
	if err != nil {
		h.logger.Error("Failed to create session", zap.Error(err))
		h.respondWithError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.configService.GetSessionCookieName(),
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	h.logger.Info("User logged in successfully",
		zap.String("user_id", user.GetID()),
		zap.String("email", user.GetEmail()))

	// Return JSON response (no redirect - let frontend handle navigation)
	h.respondWithSuccess(w, map[string]interface{}{
		"success": true,
		"user_id": user.GetID(),
		"message": "Login successful",
	})
}

// HandleSignUp handles user registration API endpoint
// UserStore extracts and validates ALL relevant data from the request
func (h *authHandlersImpl) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// UserStore handles complete data extraction, validation, and user creation from request
	user, err := h.userStore.CreateUserFromRequest(r)
	if err != nil {
		h.logger.Warn("Signup failed", zap.Error(err))
		h.respondWithError(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	h.logger.Info("User created successfully",
		zap.String("user_id", user.GetID()),
		zap.String("email", user.GetEmail()))

	// Return JSON response
	h.respondWithSuccess(w, map[string]interface{}{
		"success": true,
		"user_id": user.GetID(),
		"message": "User created successfully",
	})
}

// HandleSignOut handles user logout API endpoint
func (h *authHandlersImpl) HandleSignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookieName := h.configService.GetSessionCookieName()

	// Get session cookie
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		// Delete session
		h.sessionStore.DeleteSession(cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	h.logger.Info("User logged out successfully")

	// Return JSON response
	h.respondWithSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Logout successful",
	})
}

// respondWithError sends an error JSON response
func (h *authHandlersImpl) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// respondWithSuccess sends a success JSON response
func (h *authHandlersImpl) respondWithSuccess(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
