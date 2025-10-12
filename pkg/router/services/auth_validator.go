package services

import (
	"fmt"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// AuthValidator handles authentication and authorization validation logic
type AuthValidator struct {
	logger *zap.Logger
	config interfaces.ConfigService
}

// NewAuthValidator creates a new auth validator for DI
func NewAuthValidator(i do.Injector) (*AuthValidator, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	config := do.MustInvoke[interfaces.ConfigService](i)

	return &AuthValidator{
		logger: logger,
		config: config,
	}, nil
}

// ValidateAuthSettings validates authentication and authorization settings for a route
func (av *AuthValidator) ValidateAuthSettings(route *interfaces.Route, config *interfaces.ConfigFile, result *ValidationResult) {
	if config == nil || config.AuthSettings == nil {
		// No auth settings - this is valid for public routes
		av.logger.Debug("No auth settings found for route", zap.String("route", route.Path))
		return
	}

	authSettings := config.AuthSettings

	// Validate authentication requirements
	av.validateAuthenticationSettings(route, authSettings, result)

	// Validate authorization settings
	av.validateAuthorizationSettings(route, authSettings, result)

	// Validate role-based access
	av.validateRoleSettings(route, authSettings, result)

	av.logger.Debug("Auth settings validated",
		zap.String("route", route.Path),
		zap.String("auth_type", authSettings.Type.String()),
		zap.Strings("roles", authSettings.Roles))
}

// validateAuthenticationSettings validates basic authentication requirements
func (av *AuthValidator) validateAuthenticationSettings(route *interfaces.Route, authSettings *interfaces.AuthSettings, result *ValidationResult) {
	// Check auth type validity
	if authSettings.Type < 0 || authSettings.Type > 2 {
		result.Errors = append(result.Errors, ValidationError{
			Type:      "INVALID_AUTH_TYPE",
			Message:   "Invalid authentication type specified",
			RoutePath: route.Path,
			FilePath:  route.TemplateFile,
			Suggestions: []string{
				"Use 'public', 'user', or 'admin' auth type",
			},
		})
	}

	// Validate redirect URL if specified
	if authSettings.RedirectURL != "" {
		av.validateRedirectURL(route, authSettings.RedirectURL, result)
	}
}

// validateAuthorizationSettings validates authorization requirements
func (av *AuthValidator) validateAuthorizationSettings(route *interfaces.Route, authSettings *interfaces.AuthSettings, result *ValidationResult) {
	// Check for admin routes without proper protection
	if authSettings.Type == interfaces.AuthTypeAdmin && len(authSettings.Roles) == 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:      "ADMIN_WITHOUT_ROLES",
			Message:   "Admin route has no specific roles defined",
			RoutePath: route.Path,
			FilePath:  route.TemplateFile,
			Suggestions: []string{
				"Define specific admin roles",
				"Consider role-based access control",
			},
		})
	}
}

// validateRoleSettings validates role-based access control
func (av *AuthValidator) validateRoleSettings(route *interfaces.Route, authSettings *interfaces.AuthSettings, result *ValidationResult) {
	for _, role := range authSettings.Roles {
		if role == "" {
			result.Errors = append(result.Errors, ValidationError{
				Type:      "EMPTY_ROLE",
				Message:   "Empty role specified in required_roles",
				RoutePath: route.Path,
				FilePath:  route.TemplateFile,
			})
			continue
		}

		// Validate role format
		if !av.isValidRoleName(role) {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:      "INVALID_ROLE_FORMAT",
				Message:   fmt.Sprintf("Role name '%s' may not follow naming conventions", role),
				RoutePath: route.Path,
				FilePath:  route.TemplateFile,
				Suggestions: []string{
					"Use lowercase with underscores (e.g., 'admin_user')",
					"Avoid special characters",
				},
			})
		}
	}
}

// validateRedirectURL validates redirect URL configuration
func (av *AuthValidator) validateRedirectURL(route *interfaces.Route, redirectURL string, result *ValidationResult) {
	// Check for circular redirect
	if redirectURL == route.Path {
		result.Errors = append(result.Errors, ValidationError{
			Type:      "CIRCULAR_AUTH_REDIRECT",
			Message:   "Auth failure redirect points to the same route",
			RoutePath: route.Path,
			FilePath:  route.TemplateFile,
		})
	}

	// Validate redirect path format
	if !strings.HasPrefix(redirectURL, "/") {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:      "INVALID_REDIRECT_PATH",
			Message:   fmt.Sprintf("Redirect path should start with '/': %s", redirectURL),
			RoutePath: route.Path,
			FilePath:  route.TemplateFile,
			Suggestions: []string{
				fmt.Sprintf("Change to '/%s'", redirectURL),
			},
		})
	}
}

// isValidRoleName checks if a role name follows naming conventions
func (av *AuthValidator) isValidRoleName(role string) bool {
	// Simple validation: lowercase letters, numbers, underscores
	for _, char := range role {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return len(role) > 0
}
