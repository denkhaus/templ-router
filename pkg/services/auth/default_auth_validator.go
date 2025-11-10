package auth

import (
	"fmt"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// DefaultAuthValidator provides a no-op authentication validator
// This is used when no custom AuthValidator is provided by the application
// All authentication checks will fail - applications must provide their own implementation
type DefaultAuthValidator struct {
	logger *zap.Logger
}

// NewDefaultAuthValidator creates a default auth validator that rejects all requests
func NewDefaultAuthValidator(i do.Injector) (interfaces.AuthValidator, error) {
	logger := do.MustInvoke[*zap.Logger](i)

	return &DefaultAuthValidator{
		logger: logger,
	}, nil
}

// IsAuthenticated always returns false - applications must provide their own implementation
func (av *DefaultAuthValidator) IsAuthenticated(req *http.Request) bool {
	av.logger.Warn("Default AuthValidator used - authentication will always fail. " +
		"Please provide your own AuthValidator implementation using WithAuthValidatorFactory.")
	return false
}

// GetCurrentUser always returns an error - applications must provide their own implementation
func (av *DefaultAuthValidator) GetCurrentUser(req *http.Request) (interfaces.UserEntity, error) {
	av.logger.Warn("Default AuthValidator used - cannot get current user. " +
		"Please provide your own AuthValidator implementation using WithAuthValidatorFactory.")
	return nil, fmt.Errorf("no custom AuthValidator provided")
}

// HasRole always returns false - applications must provide their own implementation
func (av *DefaultAuthValidator) HasRole(user interfaces.UserEntity, requiredRoles []string) bool {
	av.logger.Warn("Default AuthValidator used - role checking will always fail. " +
		"Please provide your own AuthValidator implementation using WithAuthValidatorFactory.")
	return false
}