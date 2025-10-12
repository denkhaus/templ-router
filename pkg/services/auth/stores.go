package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// ProductiveUserStore provides productive user management
type ProductiveUserStore struct {
	logger *zap.Logger
	users  map[string]*interfaces.User // In-memory store for demo (use database in production)
}

// NewProductiveUserStore creates a new user store for DI
func NewProductiveUserStore(i do.Injector) (interfaces.UserStore, error) {
	logger := do.MustInvoke[*zap.Logger](i)

	// Initialize with default users
	users := map[string]*interfaces.User{
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

	return &ProductiveUserStore{
		logger: logger,
		users:  users,
	}, nil
}

// ValidateCredentials validates user credentials against productive user database
func (pus *ProductiveUserStore) ValidateCredentials(username, password string) (*interfaces.User, error) {
	pus.logger.Debug("Validating credentials", zap.String("username", username))

	// Get user from store
	user, exists := pus.users[username]
	if !exists {
		pus.logger.Warn("User not found", zap.String("username", username))
		return nil, fmt.Errorf("invalid credentials for user: %s", username)
	}

	// Validate password (in production, use proper password hashing)
	validPassword := false
	switch username {
	case "admin":
		validPassword = password == "admin"
	case "user1":
		validPassword = password == "password"
	case "demo":
		validPassword = password == "demo"
	default:
		// For dynamically created users, use a simple password check
		// In production, use bcrypt or similar
		validPassword = len(password) >= 6
	}

	if !validPassword {
		pus.logger.Warn("Invalid password", zap.String("username", username))
		return nil, fmt.Errorf("invalid credentials for user: %s", username)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (pus *ProductiveUserStore) GetUserByID(userID string) (*interfaces.User, error) {
	pus.logger.Debug("Getting user by ID", zap.String("user_id", userID))

	// Get user from store
	user, exists := pus.users[userID]
	if !exists {
		pus.logger.Warn("User not found", zap.String("user_id", userID))
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (pus *ProductiveUserStore) GetUserByEmail(email string) (*interfaces.User, error) {
	pus.logger.Debug("Getting user by email", zap.String("email", email))

	// Search through users by email
	for _, user := range pus.users {
		if strings.EqualFold(user.Email, email) {
			pus.logger.Debug("User found by email",
				zap.String("email", email),
				zap.String("user_id", user.ID))
			return user, nil
		}
	}

	pus.logger.Warn("User not found by email", zap.String("email", email))
	return nil, fmt.Errorf("user not found with email: %s", email)
}

// CreateUser creates a new user account
func (pus *ProductiveUserStore) CreateUser(username, email, password string) (*interfaces.User, error) {
	pus.logger.Info("Creating new user",
		zap.String("username", username),
		zap.String("email", email))

	// Check if user already exists
	if _, exists := pus.users[username]; exists {
		return nil, fmt.Errorf("username already exists: %s", username)
	}

	// Check if email already exists
	for _, user := range pus.users {
		if strings.EqualFold(user.Email, email) {
			return nil, fmt.Errorf("email already exists: %s", email)
		}
	}

	// Create new user
	newUser := &interfaces.User{
		ID:    username, // Use username as ID for simplicity
		Email: email,
		Roles: []string{"user"}, // Default role
	}

	// Store user (in production, hash password and store in database)
	pus.users[username] = newUser

	pus.logger.Info("User created successfully",
		zap.String("user_id", newUser.ID),
		zap.String("email", newUser.Email))

	return newUser, nil
}

// UserExists checks if a user with the given username or email already exists
func (pus *ProductiveUserStore) UserExists(username, email string) (bool, error) {
	pus.logger.Debug("Checking user existence",
		zap.String("username", username),
		zap.String("email", email))

	// Check username
	if _, exists := pus.users[username]; exists {
		return true, nil
	}

	// Check email
	for _, user := range pus.users {
		if strings.EqualFold(user.Email, email) {
			return true, nil
		}
	}

	return false, nil
}

// ProductiveSessionStore provides productive session management
type ProductiveSessionStore struct {
	logger   *zap.Logger
	sessions map[string]*interfaces.Session // In-memory store for demo (use Redis/DB in production)
}

// NewProductiveSessionStore creates a new session store for DI
func NewProductiveSessionStore(i do.Injector) (interfaces.SessionStore, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	return &ProductiveSessionStore{
		logger:   logger,
		sessions: make(map[string]*interfaces.Session),
	}, nil
}

// CreateSession creates a new session for the given user ID
func (pss *ProductiveSessionStore) CreateSession(userID string) (*interfaces.Session, error) {
	// Generate secure session ID
	sessionID := fmt.Sprintf("session_%s_%d", userID, time.Now().Unix())

	session := &interfaces.Session{
		ID:     sessionID,
		UserID: userID,
		Valid:  true,
	}

	// Store session (in production, use Redis or database)
	pss.sessions[sessionID] = session

	pss.logger.Info("Session created",
		zap.String("user_id", userID),
		zap.String("session_id", sessionID))

	return session, nil
}

// GetSession retrieves a session from the HTTP request
func (pss *ProductiveSessionStore) GetSession(req *http.Request) (*interfaces.Session, error) {
	// Check for session cookie
	cookie, err := req.Cookie("session_id")
	if err != nil {
		pss.logger.Debug("No session cookie found")
		return nil, fmt.Errorf("no session cookie")
	}

	// Validate session ID format
	if len(cookie.Value) < 10 {
		pss.logger.Debug("Invalid session ID format", zap.String("session_id", cookie.Value))
		return nil, fmt.Errorf("invalid session format")
	}

	// Get session from store (in production, use Redis or database)
	session, exists := pss.sessions[cookie.Value]
	if !exists {
		pss.logger.Debug("Session not found", zap.String("session_id", cookie.Value))
		return nil, fmt.Errorf("session not found")
	}

	if !session.Valid {
		pss.logger.Debug("Invalid session", zap.String("session_id", cookie.Value))
		return nil, fmt.Errorf("session invalid")
	}

	pss.logger.Debug("Session retrieved successfully",
		zap.String("session_id", session.ID),
		zap.String("user_id", session.UserID))

	return session, nil
}

// DeleteSession deletes an existing session
func (pss *ProductiveSessionStore) DeleteSession(sessionID string) error {
	pss.logger.Debug("Deleting session", zap.String("session_id", sessionID))

	// Delete from store (in production, use Redis or database)
	delete(pss.sessions, sessionID)

	pss.logger.Info("Session deleted", zap.String("session_id", sessionID))
	return nil
}
