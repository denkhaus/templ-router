package config

import "fmt"

// Validate validates the configuration
func (c *configImpl) Validate() error {
	// Validate server configuration
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// Validate database configuration
	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	// Validate auth configuration
	if c.Auth.MinPasswordLength < 1 {
		return fmt.Errorf("minimum password length must be at least 1")
	}

	// Validate default admin configuration
	if c.Auth.CreateDefaultAdmin {
		if c.Auth.DefaultAdminEmail == "" {
			return fmt.Errorf("default admin email cannot be empty when CreateDefaultAdmin is enabled")
		}
		if c.Auth.DefaultAdminPassword == "" {
			return fmt.Errorf("default admin password cannot be empty when CreateDefaultAdmin is enabled")
		}
		if len(c.Auth.DefaultAdminPassword) < c.Auth.MinPasswordLength {
			return fmt.Errorf("default admin password must meet minimum password length requirement (%d characters)", c.Auth.MinPasswordLength)
		}
		if c.Auth.DefaultAdminFirstName == "" {
			return fmt.Errorf("default admin first name cannot be empty when CreateDefaultAdmin is enabled")
		}
		if c.Auth.DefaultAdminLastName == "" {
			return fmt.Errorf("default admin last name cannot be empty when CreateDefaultAdmin is enabled")
		}
	}

	// Validate email configuration
	if c.Email.SMTPPort < 1 || c.Email.SMTPPort > 65535 {
		return fmt.Errorf("invalid SMTP port: %d", c.Email.SMTPPort)
	}

	// Validate security configuration
	if c.Security.RateLimitRequests < 1 {
		return fmt.Errorf("rate limit requests must be at least 1")
	}

	return nil
}
