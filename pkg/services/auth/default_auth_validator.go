// Package auth provides authentication and authorization services for the templ-router system.
// It includes default implementations, user validation, role-based access control,
// and session management for secure web applications.
package auth

import (
	"fmt"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// defaultAuthValidatorImpl provides a no-op authentication validator
// This is used when no custom AuthValidator is provided by the application
// All authentication checks will fail - applications must provide their own implementation
type defaultAuthValidatorImpl struct {
	logger *zap.Logger
}

// NewDefaultAuthValidator creates a default auth validator that rejects all requests
func NewDefaultAuthValidator(i do.Injector) (interfaces.AuthValidator, error) {
	logger := do.MustInvoke[*zap.Logger](i)

	return &defaultAuthValidatorImpl{
		logger: logger,
	}, nil
}

// IsAuthenticated always returns false - applications must provide their own implementation
func (av *defaultAuthValidatorImpl) IsAuthenticated(_ *http.Request) (bool, error) {
	av.logger.Warn("Default AuthValidator used - authentication will always fail. " +
		"Please provide your own AuthValidator implementation using WithAuthValidatorFactory.")
	return false, nil
}

// GetCurrentUser always returns an error - applications must provide their own implementation
func (av *defaultAuthValidatorImpl) GetCurrentUser(_ *http.Request) (interfaces.UserEntity, error) {
	av.logger.Warn("Default AuthValidator used - cannot get current user. " +
		"Please provide your own AuthValidator implementation using WithAuthValidatorFactory.")
	return nil, fmt.Errorf("no custom AuthValidator provided")
}

// HasRole always returns false - applications must provide their own implementation
func (av *defaultAuthValidatorImpl) HasRole(_ interfaces.UserEntity, _ []string) bool {
	av.logger.Warn("Default AuthValidator used - role checking will always fail. " +
		"Please provide your own AuthValidator implementation using WithAuthValidatorFactory.")
	return false
}
