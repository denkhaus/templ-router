package config

import (
	"fmt"
	
	"github.com/denkhaus/templ-router/pkg/shared"
)

// Validate validates the configuration
func (c *configImpl) Validate() error {
	// Validate server configuration
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return shared.NewValidationError("Invalid server port").
			WithDetails(fmt.Sprintf("Port %d is outside valid range 1-65535", c.Server.Port)).
			WithContext("field", "server.port").
			WithContext("value", c.Server.Port).
			WithContext("valid_range", "1-65535")
	}

	// Validate database configuration
	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return shared.NewValidationError("Invalid database port").
			WithDetails(fmt.Sprintf("Port %d is outside valid range 1-65535", c.Database.Port)).
			WithContext("field", "database.port").
			WithContext("value", c.Database.Port).
			WithContext("valid_range", "1-65535")
	}

	// Validate auth configuration
	if c.Auth.MinPasswordLength < 1 {
		return shared.NewValidationError("Invalid minimum password length").
			WithDetails("Minimum password length must be at least 1 character").
			WithContext("field", "auth.min_password_length").
			WithContext("value", c.Auth.MinPasswordLength).
			WithContext("minimum", 1)
	}

	// Validate default admin configuration
	if c.Auth.CreateDefaultAdmin {
		if c.Auth.DefaultAdminEmail == "" {
			return shared.NewConfigurationError("Default admin email is required").
				WithDetails("Email cannot be empty when CreateDefaultAdmin is enabled").
				WithContext("field", "auth.default_admin_email").
				WithContext("create_default_admin", true)
		}
		if c.Auth.DefaultAdminPassword == "" {
			return shared.NewConfigurationError("Default admin password is required").
				WithDetails("Password cannot be empty when CreateDefaultAdmin is enabled").
				WithContext("field", "auth.default_admin_password").
				WithContext("create_default_admin", true)
		}
		if len(c.Auth.DefaultAdminPassword) < c.Auth.MinPasswordLength {
			return shared.NewValidationError("Default admin password too short").
				WithDetails(fmt.Sprintf("Password must be at least %d characters", c.Auth.MinPasswordLength)).
				WithContext("field", "auth.default_admin_password").
				WithContext("current_length", len(c.Auth.DefaultAdminPassword)).
				WithContext("required_length", c.Auth.MinPasswordLength)
		}
		if c.Auth.DefaultAdminFirstName == "" {
			return shared.NewConfigurationError("Default admin first name is required").
				WithDetails("First name cannot be empty when CreateDefaultAdmin is enabled").
				WithContext("field", "auth.default_admin_first_name").
				WithContext("create_default_admin", true)
		}
		if c.Auth.DefaultAdminLastName == "" {
			return shared.NewConfigurationError("Default admin last name is required").
				WithDetails("Last name cannot be empty when CreateDefaultAdmin is enabled").
				WithContext("field", "auth.default_admin_last_name").
				WithContext("create_default_admin", true)
		}
	}

	// Validate email configuration
	if c.Email.SMTPPort < 1 || c.Email.SMTPPort > 65535 {
		return shared.NewValidationError("Invalid SMTP port").
			WithDetails(fmt.Sprintf("Port %d is outside valid range 1-65535", c.Email.SMTPPort)).
			WithContext("field", "email.smtp_port").
			WithContext("value", c.Email.SMTPPort).
			WithContext("valid_range", "1-65535")
	}

	// Validate security configuration
	if c.Security.RateLimitRequests < 1 {
		return shared.NewValidationError("Invalid rate limit configuration").
			WithDetails("Rate limit requests must be at least 1").
			WithContext("field", "security.rate_limit_requests").
			WithContext("value", c.Security.RateLimitRequests).
			WithContext("minimum", 1)
	}

	return nil
}
