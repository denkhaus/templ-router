package services

import (
	"encoding/json"
	"net/http"

	"github.com/denkhaus/templ-router/demo/pkg/interfaces"
	sharedInterfaces "github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// demoAuthHandlers provides authentication API handlers for the demo application
// This shows how to implement authentication endpoints using custom stores
type demoAuthHandlers struct {
	userStore     interfaces.UserStore
	sessionStore  interfaces.SessionStore
	configService sharedInterfaces.ConfigService
	logger        *zap.Logger
}

// NewDemoAuthHandlers creates new demo auth handlers
func NewDemoAuthHandlers(i do.Injector) (interfaces.AuthHandlers, error) {
	userStore := do.MustInvoke[interfaces.UserStore](i)
	configService := do.MustInvoke[sharedInterfaces.ConfigService](i)
	sessionStore := do.MustInvoke[interfaces.SessionStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &demoAuthHandlers{
		userStore:     userStore,
		configService: configService,
		sessionStore:  sessionStore,
		logger:        logger,
	}, nil
}

// RegisterRoutes registers authentication routes
func (h *demoAuthHandlers) RegisterRoutes(registerFunc func(method, path string, handler http.HandlerFunc)) {
	registerFunc("POST", "/api/auth/signin", h.HandleSignIn)
	registerFunc("POST", "/api/auth/signup", h.HandleSignUp)
	registerFunc("POST", "/api/auth/signout", h.HandleSignOut)
}

// HandleSignIn handles user login API endpoint
func (h *demoAuthHandlers) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, r, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// UserStore handles complete data extraction and validation from request
	user, err := h.userStore.ValidateCredentialsFromRequest(r)
	if err != nil {
		h.logger.Warn("Login failed", zap.Error(err))
		h.respondWithError(w, r, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session, err := h.sessionStore.CreateSession(user.GetID())
	if err != nil {
		h.logger.Error("Failed to create session", zap.Error(err))
		h.respondWithError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	h.logger.Info("User logged in successfully",
		zap.String("user_id", user.GetID()),
		zap.String("email", user.GetEmail()))

	// Redirect to success route on successful login
	successRoute := h.configService.GetSignInSuccessRoute()
	if successRoute != "" {
		successRoute := i18n.LocalizeRouteIfRequired(r.Context(), successRoute)

		// Check if this is an HTMX request
		if h.isHTMXRequest(r) {
			w.Header().Set("HX-Redirect", successRoute)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Redirect(w, r, successRoute, http.StatusSeeOther)
		return
	}

	// Fallback to JSON response if no redirect route configured
	h.respondWithSuccess(w, map[string]interface{}{
		"success": true,
		"user_id": user.GetID(),
		"message": "Login successful",
	})
}

// HandleSignUp handles user registration API endpoint
func (h *demoAuthHandlers) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, r, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// UserStore handles complete data extraction, validation, and user creation from request
	user, err := h.userStore.CreateUserFromRequest(r)
	if err != nil {
		h.logger.Warn("Signup failed", zap.Error(err))
		h.respondWithError(w, r, "Failed to create user", http.StatusBadRequest)
		return
	}

	h.logger.Info("User created successfully",
		zap.String("user_id", user.GetID()),
		zap.String("email", user.GetEmail()))

	// Redirect to success route on successful signup
	successRoute := h.configService.GetSignUpSuccessRoute()
	if successRoute != "" {
		successRoute := i18n.LocalizeRouteIfRequired(r.Context(), successRoute)

		// Check if this is an HTMX request
		if h.isHTMXRequest(r) {
			w.Header().Set("HX-Redirect", successRoute)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Redirect(w, r, successRoute, http.StatusSeeOther)
		return
	}

	// Fallback to JSON response if no redirect route configured
	h.respondWithSuccess(w, map[string]interface{}{
		"success": true,
		"user_id": user.GetID(),
		"message": "User created successfully",
	})
}

// HandleSignOut handles user logout API endpoint
func (h *demoAuthHandlers) HandleSignOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, r, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionCookieName := "session_id"

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

	// Redirect to success route on successful logout
	successRoute := h.configService.GetSignOutSuccessRoute()
	if successRoute != "" {
		successRoute := i18n.LocalizeRouteIfRequired(r.Context(), successRoute)
		http.Redirect(w, r, successRoute, http.StatusSeeOther)
		return
	}

	// Fallback to JSON response if no redirect route configured
	h.respondWithSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Logout successful",
	})
}

// respondWithError sends an error response (HTML for HTMX, JSON for API)
func (h *demoAuthHandlers) respondWithError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	// Check if this is an HTMX request
	if h.isHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(statusCode)
		w.Write([]byte(`<span class="font-medium">` + message + `</span>`))
		return
	}

	// Default JSON response for API requests
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// isHTMXRequest checks if the request is from HTMX
func (h *demoAuthHandlers) isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// respondWithSuccess sends a success JSON response
func (h *demoAuthHandlers) respondWithSuccess(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
