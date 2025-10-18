package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// Mock AuthProvider for testing handlers
type mockAuthProviderForHandlers struct {
	loginResult  *LoginResult
	signupResult *SignupResult
	logoutError  error
	sessionID    string
}

func (m *mockAuthProviderForHandlers) ProcessLogin(req LoginRequest) *LoginResult {
	if m.loginResult != nil {
		return m.loginResult
	}
	return &LoginResult{
		Success:     true,
		User:        &interfaces.User{ID: "user123", Email: req.Username},
		SessionID:   m.sessionID,
		RedirectURL: "/dashboard",
	}
}

func (m *mockAuthProviderForHandlers) ProcessSignup(req SignupRequest) *SignupResult {
	if m.signupResult != nil {
		return m.signupResult
	}
	return &SignupResult{
		Success:     true,
		User:        &interfaces.User{ID: "user123", Email: req.Email},
		SessionID:   m.sessionID,
		RedirectURL: "/en/dashboard",
	}
}

func (m *mockAuthProviderForHandlers) ProcessLogout(r *http.Request) error {
	return m.logoutError
}

func (m *mockAuthProviderForHandlers) SetSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionID,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func (m *mockAuthProviderForHandlers) ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func createAuthHandlersTestContainer() do.Injector {
	injector := do.New()

	do.Provide(injector, func(i do.Injector) (AuthProvider, error) {
		return &mockAuthProviderForHandlers{sessionID: "test-session-123"}, nil
	})

	do.Provide(injector, func(i do.Injector) (*zap.Logger, error) {
		return zap.NewNop(), nil
	})

	return injector
}

func TestNewAuthHandlers(t *testing.T) {
	injector := createAuthHandlersTestContainer()
	defer injector.Shutdown()

	handlers, err := NewAuthHandlers(injector)
	if err != nil {
		t.Fatalf("NewAuthHandlers() error = %v", err)
	}

	if handlers == nil {
		t.Fatal("NewAuthHandlers() returned nil")
	}

	// Verify interface compliance
	var _ AuthHandlersInterface = handlers
}

func TestAuthHandlers_HandleLogin(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		mockProvider   *mockAuthProviderForHandlers
		expectStatus   int
		expectLocation string
		expectCookie   bool
	}{
		{
			name: "successful login",
			formData: url.Values{
				"username": []string{"testuser"},
				"password": []string{"password123"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				sessionID: "session123",
				loginResult: &LoginResult{
					Success:     true,
					User:        &interfaces.User{ID: "user123", Email: "testuser"},
					SessionID:   "session123",
					RedirectURL: "/dashboard",
				},
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/dashboard",
			expectCookie:   true,
		},
		{
			name: "failed login",
			formData: url.Values{
				"username": []string{"testuser"},
				"password": []string{"wrongpassword"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				loginResult: &LoginResult{
					Success: false,
					Error:   "Invalid credentials",
				},
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/login?error=Invalid credentials",
			expectCookie:   false,
		},
		{
			name: "missing username",
			formData: url.Values{
				"password": []string{"password123"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				sessionID: "session123",
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/dashboard",
			expectCookie:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injector := do.New()
			defer injector.Shutdown()

			do.ProvideValue[AuthProvider](injector, tt.mockProvider)
			do.ProvideValue[*zap.Logger](injector, zap.NewNop())

			handlers, err := NewAuthHandlers(injector)
			if err != nil {
				t.Fatalf("NewAuthHandlers() error = %v", err)
			}

			// Create request
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			// Call handler
			handlers.HandleLogin(w, req)

			// Check status code
			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}

			// Check redirect location
			location := w.Header().Get("Location")
			if location != tt.expectLocation {
				t.Errorf("Expected location %q, got %q", tt.expectLocation, location)
			}

			// Check session cookie
			cookies := w.Result().Cookies()
			hasCookie := false
			for _, cookie := range cookies {
				if cookie.Name == "session" && cookie.Value != "" {
					hasCookie = true
					break
				}
			}

			if hasCookie != tt.expectCookie {
				t.Errorf("Expected cookie presence %v, got %v", tt.expectCookie, hasCookie)
			}
		})
	}
}

func TestAuthHandlers_HandleSignup(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		mockProvider   *mockAuthProviderForHandlers
		expectStatus   int
		expectLocation string
		expectCookie   bool
	}{
		{
			name: "successful signup",
			formData: url.Values{
				"username": []string{"newuser"},
				"email":    []string{"newuser@example.com"},
				"password": []string{"password123"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				sessionID: "session123",
				signupResult: &SignupResult{
					Success:     true,
					User:        &interfaces.User{ID: "user123", Email: "newuser@example.com"},
					SessionID:   "session123",
					RedirectURL: "/welcome",
				},
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/welcome",
			expectCookie:   true,
		},
		{
			name: "failed signup",
			formData: url.Values{
				"username": []string{"existinguser"},
				"email":    []string{"existing@example.com"},
				"password": []string{"password123"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				signupResult: &SignupResult{
					Success: false,
					Error:   "Username already exists",
				},
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/signup?error=Username already exists",
			expectCookie:   false,
		},
		{
			name: "validation error - short username",
			formData: url.Values{
				"username": []string{"ab"},
				"email":    []string{"test@example.com"},
				"password": []string{"password123"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				sessionID: "session123",
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/signup?error=username must be at least 3 characters",
			expectCookie:   false,
		},
		{
			name: "validation error - invalid email",
			formData: url.Values{
				"username": []string{"testuser"},
				"email":    []string{"invalid-email"},
				"password": []string{"password123"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				sessionID: "session123",
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/signup?error=invalid email format",
			expectCookie:   false,
		},
		{
			name: "validation error - short password",
			formData: url.Values{
				"username": []string{"testuser"},
				"email":    []string{"test@example.com"},
				"password": []string{"123"},
			},
			mockProvider: &mockAuthProviderForHandlers{
				sessionID: "session123",
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/signup?error=password must be at least 6 characters",
			expectCookie:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injector := do.New()
			defer injector.Shutdown()

			do.ProvideValue[AuthProvider](injector, tt.mockProvider)
			do.ProvideValue[*zap.Logger](injector, zap.NewNop())

			handlers, err := NewAuthHandlers(injector)
			if err != nil {
				t.Fatalf("NewAuthHandlers() error = %v", err)
			}

			// Create request
			req := httptest.NewRequest("POST", "/signup", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			// Call handler
			handlers.HandleSignup(w, req)

			// Check status code
			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}

			// Check redirect location
			location := w.Header().Get("Location")
			if location != tt.expectLocation {
				t.Errorf("Expected location %q, got %q", tt.expectLocation, location)
			}

			// Check session cookie
			cookies := w.Result().Cookies()
			hasCookie := false
			for _, cookie := range cookies {
				if cookie.Name == "session" && cookie.Value != "" {
					hasCookie = true
					break
				}
			}

			if hasCookie != tt.expectCookie {
				t.Errorf("Expected cookie presence %v, got %v", tt.expectCookie, hasCookie)
			}
		})
	}
}

func TestAuthHandlers_HandleLogout(t *testing.T) {
	tests := []struct {
		name           string
		mockProvider   *mockAuthProviderForHandlers
		expectStatus   int
		expectLocation string
		expectCookie   bool
	}{
		{
			name: "successful logout",
			mockProvider: &mockAuthProviderForHandlers{
				logoutError: nil,
			},
			expectStatus:   http.StatusFound,
			expectLocation: "/login",
			expectCookie:   true, // Cookie should be cleared (set with MaxAge: -1)
		},
		{
			name: "logout with error",
			mockProvider: &mockAuthProviderForHandlers{
				logoutError: http.ErrNoCookie,
			},
			expectStatus:   http.StatusInternalServerError,
			expectLocation: "",
			expectCookie:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injector := do.New()
			defer injector.Shutdown()

			do.ProvideValue[AuthProvider](injector, tt.mockProvider)
			do.ProvideValue[*zap.Logger](injector, zap.NewNop())

			handlers, err := NewAuthHandlers(injector)
			if err != nil {
				t.Fatalf("NewAuthHandlers() error = %v", err)
			}

			// Create request
			req := httptest.NewRequest("POST", "/logout", nil)
			w := httptest.NewRecorder()

			// Call handler
			handlers.HandleLogout(w, req)

			// Check status code
			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}

			// Check redirect location
			location := w.Header().Get("Location")
			if location != tt.expectLocation {
				t.Errorf("Expected location %q, got %q", tt.expectLocation, location)
			}

			// Check session cookie clearing
			if tt.expectCookie {
				cookies := w.Result().Cookies()
				hasClearedCookie := false
				for _, cookie := range cookies {
					if cookie.Name == "session" && cookie.MaxAge == -1 {
						hasClearedCookie = true
						break
					}
				}

				if !hasClearedCookie {
					t.Error("Expected session cookie to be cleared")
				}
			}
		})
	}
}

func TestAuthHandlers_RegisterRoutes(t *testing.T) {
	injector := createAuthHandlersTestContainer()
	defer injector.Shutdown()

	handlers, err := NewAuthHandlers(injector)
	if err != nil {
		t.Fatalf("NewAuthHandlers() error = %v", err)
	}

	// Track registered routes
	registeredRoutes := make(map[string]string)

	registerFunc := func(method, path string, handler http.HandlerFunc) {
		registeredRoutes[method+" "+path] = method + " " + path
		if handler == nil {
			t.Errorf("Handler for %s %s is nil", method, path)
		}
	}

	// Call RegisterRoutes
	handlers.RegisterRoutes(registerFunc)

	// Verify expected routes were registered
	expectedRoutes := []string{
		"POST /login",
		"POST /signup",
		"POST /logout",
	}

	for _, expectedRoute := range expectedRoutes {
		if _, exists := registeredRoutes[expectedRoute]; !exists {
			t.Errorf("Expected route %s was not registered", expectedRoute)
		}
	}

	if len(registeredRoutes) != len(expectedRoutes) {
		t.Errorf("Expected %d routes, got %d", len(expectedRoutes), len(registeredRoutes))
	}
}

func TestAuthHandlers_ValidationHelpers(t *testing.T) {
	injector := createAuthHandlersTestContainer()
	defer injector.Shutdown()

	handlers, err := NewAuthHandlers(injector)
	if err != nil {
		t.Fatalf("NewAuthHandlers() error = %v", err)
	}

	// Access the private implementation to test validation methods
	authHandlersImpl := handlers.(*authHandlers)

	// Test email validation
	emailTests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user@domain.co.uk", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
		{"a@b.c", false}, // Too short for our validation
	}

	for _, tt := range emailTests {
		t.Run("email_"+tt.email, func(t *testing.T) {
			result := authHandlersImpl.isValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}

	// Test signup validation
	signupTests := []struct {
		name    string
		request SignupRequest
		valid   bool
	}{
		{
			name: "valid signup",
			request: SignupRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "short username",
			request: SignupRequest{
				Username: "ab",
				Email:    "test@example.com",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "invalid email",
			request: SignupRequest{
				Username: "testuser",
				Email:    "invalid",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "short password",
			request: SignupRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "123",
			},
			valid: false,
		},
	}

	for _, tt := range signupTests {
		t.Run(tt.name, func(t *testing.T) {
			err := authHandlersImpl.validateSignupRequest(tt.request)
			isValid := err == nil

			if isValid != tt.valid {
				t.Errorf("validateSignupRequest() valid = %v, want %v, error = %v", isValid, tt.valid, err)
			}
		})
	}
}