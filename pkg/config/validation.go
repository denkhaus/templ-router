package config

import (
	"github.com/denkhaus/templ-router/pkg/shared"
)

// Validate validates the configuration using centralized validators
func (c *configImpl) Validate() error {
	// Initialize reusable validators
	portValidator := NewPortValidator("")
	emailValidator := NewEmailValidator("")
	passwordValidator := NewPasswordValidator("", 0, false)
	positiveNumberValidator := NewPositiveNumberValidator("", 1)
	requiredStringValidator := NewRequiredStringValidator("")

	// Validate server configuration
	portValidator.fieldName = "server.port"
	if err := portValidator.ValidatePort(c.Server.Port); err != nil {
		return err
	}

	// Validate database configuration
	portValidator.fieldName = "database.port"
	if err := portValidator.ValidatePort(c.Database.Port); err != nil {
		return err
	}

	// Validate auth configuration
	positiveNumberValidator.fieldName = "auth.min_password_length"
	if err := positiveNumberValidator.ValidatePositiveNumber(c.Auth.MinPasswordLength); err != nil {
		return err
	}

	// Validate default admin configuration (conditionally)
	if c.Auth.CreateDefaultAdmin {
		emailValidator.fieldName = "auth.default_admin_email"
		if err := emailValidator.ValidateEmail(c.Auth.DefaultAdminEmail); err != nil {
			if appErr, ok := err.(*shared.AppError); ok {
				return appErr.WithContext("create_default_admin", true)
			}
			return err
		}

		passwordValidator.fieldName = "auth.default_admin_password"
		passwordValidator.minLength = c.Auth.MinPasswordLength
		if err := passwordValidator.ValidatePassword(c.Auth.DefaultAdminPassword); err != nil {
			if appErr, ok := err.(*shared.AppError); ok {
				return appErr.WithContext("create_default_admin", true)
			}
			return err
		}

		requiredStringValidator.fieldName = "auth.default_admin_first_name"
		if err := requiredStringValidator.ValidateRequiredString(c.Auth.DefaultAdminFirstName); err != nil {
			if appErr, ok := err.(*shared.AppError); ok {
				return appErr.WithContext("create_default_admin", true)
			}
			return err
		}

		requiredStringValidator.fieldName = "auth.default_admin_last_name"
		if err := requiredStringValidator.ValidateRequiredString(c.Auth.DefaultAdminLastName); err != nil {
			if appErr, ok := err.(*shared.AppError); ok {
				return appErr.WithContext("create_default_admin", true)
			}
			return err
		}
	}

	// Validate email configuration
	portValidator.fieldName = "email.smtp_port"
	if err := portValidator.ValidatePort(c.Email.SMTPPort); err != nil {
		return err
	}

	// Validate security configuration
	positiveNumberValidator.fieldName = "security.rate_limit_requests"
	if err := positiveNumberValidator.ValidatePositiveNumber(c.Security.RateLimitRequests); err != nil {
		return err
	}

	return nil
}
