package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// DemoUser represents a concrete user implementation for demo purposes
type DemoUser struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// Implement UserEntity interface
func (u *DemoUser) GetID() string      { return u.ID }
func (u *DemoUser) GetEmail() string   { return u.Email }
func (u *DemoUser) GetRoles() []string { return u.Roles }

// DefaultUserStore provides a default in-memory user store implementation
// Users can replace this with their own database-backed implementation
type DefaultUserStore struct {
	logger *zap.Logger
	users  map[string]*DemoUser // In-memory store for development
}

// NewDefaultUserStore creates a new default user store for DI
func NewDefaultUserStore(i do.Injector) (interfaces.UserStore, error) {
	logger := do.MustInvoke[*zap.Logger](i)

	// Initialize with default demo users
	users := map[string]*DemoUser{
		"admin": {
			ID:    "admin",
			Email: "admin@example.com",
			Roles: []string{"admin", "user"},
		},
		"user1": {
			ID:    "user1",
			Email: "user1@example.com",
			Roles: []string{"user"},
		},
		"demo": {
			ID:    "demo",
			Email: "demo@example.com",
			Roles: []string{"user"},
		},
	}

	return &DefaultUserStore{
		logger: logger,
		users:  users,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *DefaultUserStore) GetUserByID(userID string) (interfaces.UserEntity, error) {
	user, exists := s.users[userID]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", userID)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *DefaultUserStore) GetUserByEmail(email string) (interfaces.UserEntity, error) {
	for _, user := range s.users {
		if user.GetEmail() == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with email %s not found", email)
}

// ValidateCredentials validates user credentials (simple demo implementation)
func (s *DefaultUserStore) ValidateCredentials(email, password string) (interfaces.UserEntity, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Simple demo validation - in production use bcrypt
	// For demo: admin@example.com / admin123, user1@example.com / user123, etc.
	expectedPassword := "demo123"
	switch email {
	case "admin@example.com":
		expectedPassword = "admin123"
	case "user1@example.com":
		expectedPassword = "user123"
	}

	if password != expectedPassword {
		return nil, errors.New("invalid credentials")
	}

	s.logger.Info("User authenticated successfully",
		zap.String("user_id", user.GetID()),
		zap.String("email", user.GetEmail()))

	return user, nil
}

// CreateUser creates a new user
func (s *DefaultUserStore) CreateUser(username, email, password string) (interfaces.UserEntity, error) {
	// Check if user already exists
	if exists, _ := s.UserExists(username, email); exists {
		return nil, fmt.Errorf("user with username %s or email %s already exists", username, email)
	}

	// Generate new ID
	userID := fmt.Sprintf("user_%d", len(s.users)+1)

	user := &DemoUser{
		ID:    userID,
		Email: email,
		Roles: []string{"user"}, // Default role
	}

	s.users[userID] = user

	s.logger.Info("User created successfully",
		zap.String("user_id", user.GetID()),
		zap.String("email", email))

	return user, nil
}

// UserExists checks if a user exists by username or email
func (s *DefaultUserStore) UserExists(username, email string) (bool, error) {
	for _, user := range s.users {
		if user.GetEmail() == email {
			return true, nil
		}
	}
	return false, nil
}

// ValidateCredentialsFromRequest extracts and validates credentials from HTTP request
// This method handles all data extraction and validation logic
func (s *DefaultUserStore) ValidateCredentialsFromRequest(req *http.Request) (interfaces.UserEntity, error) {
	// Extract credentials from request
	email := req.FormValue("email")
	password := req.FormValue("password")

	// Validate required fields
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	// Additional validation can be added here
	// - Email format validation
	// - Password strength validation
	// - Rate limiting checks
	// - CAPTCHA validation
	// etc.

	// Use existing validation method
	user, err := s.ValidateCredentials(email, password)
	if err != nil {
		return nil, err
	}

	s.logger.Info("User credentials validated from request",
		zap.String("user_id", user.GetID()),
		zap.String("email", user.GetEmail()),
		zap.String("remote_addr", req.RemoteAddr))

	return user, nil
}

// CreateUserFromRequest extracts and creates user from HTTP request
// This method handles all data extraction, validation, and user creation logic
func (s *DefaultUserStore) CreateUserFromRequest(req *http.Request) (interfaces.UserEntity, error) {
	// Extract user data from request
	username := req.FormValue("username")
	email := req.FormValue("email")
	password := req.FormValue("password")

	// Demo: Extract additional fields that could be in the form
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	phone := req.FormValue("phone")

	// Validate required fields
	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email and password are required")
	}

	// Additional validation can be added here
	// - Email format validation
	// - Username uniqueness validation
	// - Password strength validation
	// - Phone number format validation
	// - Terms acceptance validation
	// - Age verification
	// etc.

	// Check if user already exists
	if exists, _ := s.UserExists(username, email); exists {
		return nil, fmt.Errorf("user with username %s or email %s already exists", username, email)
	}

	// Generate new ID
	userID := fmt.Sprintf("user_%d", len(s.users)+1)

	// Create user with all extracted data
	user := &DemoUser{
		ID:    userID,
		Email: email,
		Roles: []string{"user"}, // Default role
		// Additional fields could be stored in a custom user struct:
		// Firstname: firstname,
		// Lastname:  lastname,
		// Phone:     phone,
	}

	s.users[userID] = user

	s.logger.Info("User created from request",
		zap.String("user_id", user.GetID()),
		zap.String("username", username),
		zap.String("email", email),
		zap.String("firstname", firstname),
		zap.String("lastname", lastname),
		zap.String("phone", phone),
		zap.String("remote_addr", req.RemoteAddr))

	return user, nil
}
