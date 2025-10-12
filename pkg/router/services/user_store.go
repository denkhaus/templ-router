package services

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

type simpleUserStore struct {
	logger *zap.Logger
}

// NewInMemoryUserStore creates a new user store for DI
func NewInMemoryUserStore(i do.Injector) (interfaces.UserStore, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	return &simpleUserStore{logger: logger}, nil
}

func (s *simpleUserStore) GetUserByID(userID string) (*interfaces.User, error) {
	s.logger.Debug("Getting user by ID", zap.String("user_id", userID))

	// Produktive User-Datenbank-Simulation
	switch userID {
	case "admin":
		return &interfaces.User{
			ID:    "admin",
			Email: "admin@example.com",
			Roles: []string{"admin", "user"},
		}, nil
	case "user1":
		return &interfaces.User{
			ID:    "user1",
			Email: "user1@example.com",
			Roles: []string{"user"},
		}, nil
	default:
		s.logger.Warn("User not found", zap.String("user_id", userID))
		return nil, fmt.Errorf("user not found: %s", userID)
	}
}

func (s *simpleUserStore) GetUserByEmail(email string) (*interfaces.User, error) {
	return s.GetUserByID("user1")
}

func (s *simpleUserStore) CreateUser(username, email, password string) (*interfaces.User, error) {
	s.logger.Info("CreateUser not implemented in simple store", zap.String("username", username))
	return nil, fmt.Errorf("CreateUser not implemented in simple store")
}

func (s *simpleUserStore) UserExists(username, email string) (bool, error) {
	s.logger.Debug("UserExists not implemented in simple store", zap.String("username", username))
	return false, nil
}

func (s *simpleUserStore) ValidateCredentials(email, password string) (*interfaces.User, error) {
	s.logger.Debug("Validating credentials", zap.String("email", email))

	// Produktive Credential-Validierung
	switch {
	case email == "admin" && password == "admin":
		return s.GetUserByID("admin")
	case email == "user1" && password == "password":
		return s.GetUserByID("user1")
	default:
		s.logger.Warn("Invalid credentials", zap.String("email", email))
		return nil, fmt.Errorf("invalid credentials")
	}
}
