package config

import (
	"fmt"
	"strings"
)

// maskSensitive masks sensitive data for logging
func maskSensitive(value string) string {
	if value == "" {
		return "<empty>"
	}
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

// LogSummary logs a summary of the configuration with sensitive data masked
func (c *Config) LogSummary() {
	fmt.Println("=== Configuration Summary ===")
	
	// Server Configuration
	fmt.Printf("Server:\n")
	fmt.Printf("  Host: %s\n", c.Server.Host)
	fmt.Printf("  Port: %d\n", c.Server.Port)
	fmt.Printf("  Base URL: %s\n", c.Server.BaseURL)
	fmt.Printf("  Read Timeout: %s\n", c.Server.ReadTimeout)
	fmt.Printf("  Write Timeout: %s\n", c.Server.WriteTimeout)
	fmt.Printf("  Idle Timeout: %s\n", c.Server.IdleTimeout)
	fmt.Printf("  Shutdown Timeout: %s\n", c.Server.ShutdownTimeout)
	
	// Database Configuration
	fmt.Printf("Database:\n")
	fmt.Printf("  Host: %s\n", c.Database.Host)
	fmt.Printf("  Port: %d\n", c.Database.Port)
	fmt.Printf("  User: %s\n", c.Database.User)
	fmt.Printf("  Password: %s\n", maskSensitive(c.Database.Password))
	fmt.Printf("  Name: %s\n", c.Database.Name)
	fmt.Printf("  SSL Mode: %s\n", c.Database.SSLMode)
	
	// Authentication Configuration
	fmt.Printf("Authentication:\n")
	fmt.Printf("  Require Email Verification: %t\n", c.Auth.RequireEmailVerification)
	fmt.Printf("  Verification Token Expiry: %s\n", c.Auth.VerificationTokenExpiry)
	fmt.Printf("  Session Cookie Name: %s\n", c.Auth.SessionCookieName)
	fmt.Printf("  Session Expiry: %s\n", c.Auth.SessionExpiry)
	fmt.Printf("  Session Secure: %t\n", c.Auth.SessionSecure)
	fmt.Printf("  Session HTTP Only: %t\n", c.Auth.SessionHttpOnly)
	fmt.Printf("  Session Same Site: %s\n", c.Auth.SessionSameSite)
	fmt.Printf("  Min Password Length: %d\n", c.Auth.MinPasswordLength)
	fmt.Printf("  Require Strong Password: %t\n", c.Auth.RequireStrongPasswd)
	fmt.Printf("  Create Default Admin: %t\n", c.Auth.CreateDefaultAdmin)
	if c.Auth.CreateDefaultAdmin {
		fmt.Printf("  Default Admin Email: %s\n", c.Auth.DefaultAdminEmail)
		fmt.Printf("  Default Admin Password: %s\n", maskSensitive(c.Auth.DefaultAdminPassword))
		fmt.Printf("  Default Admin First Name: %s\n", c.Auth.DefaultAdminFirstName)
		fmt.Printf("  Default Admin Last Name: %s\n", c.Auth.DefaultAdminLastName)
	}
	
	// Email Configuration
	fmt.Printf("Email:\n")
	fmt.Printf("  SMTP Host: %s\n", func() string {
		if c.Email.SMTPHost == "" {
			return "<not configured>"
		}
		return c.Email.SMTPHost
	}())
	fmt.Printf("  SMTP Port: %d\n", c.Email.SMTPPort)
	fmt.Printf("  SMTP Username: %s\n", func() string {
		if c.Email.SMTPUsername == "" {
			return "<not configured>"
		}
		return maskSensitive(c.Email.SMTPUsername)
	}())
	fmt.Printf("  SMTP Password: %s\n", func() string {
		if c.Email.SMTPPassword == "" {
			return "<not configured>"
		}
		return maskSensitive(c.Email.SMTPPassword)
	}())
	fmt.Printf("  SMTP Use TLS: %t\n", c.Email.SMTPUseTLS)
	fmt.Printf("  From Email: %s\n", c.Email.FromEmail)
	fmt.Printf("  From Name: %s\n", c.Email.FromName)
	fmt.Printf("  Reply To Email: %s\n", func() string {
		if c.Email.ReplyToEmail == "" {
			return "<not set>"
		}
		return c.Email.ReplyToEmail
	}())
	fmt.Printf("  Enable Dummy Mode: %t\n", c.Email.EnableDummyMode)
	
	// Security Configuration
	fmt.Printf("Security:\n")
	fmt.Printf("  CSRF Secret: %s\n", maskSensitive(c.Security.CSRFSecret))
	fmt.Printf("  CSRF Secure: %t\n", c.Security.CSRFSecure)
	fmt.Printf("  CSRF HTTP Only: %t\n", c.Security.CSRFHttpOnly)
	fmt.Printf("  CSRF Same Site: %s\n", c.Security.CSRFSameSite)
	fmt.Printf("  Enable Rate Limit: %t\n", c.Security.EnableRateLimit)
	fmt.Printf("  Rate Limit Requests: %d\n", c.Security.RateLimitRequests)
	fmt.Printf("  Enable Security Headers: %t\n", c.Security.EnableSecurityHeaders)
	fmt.Printf("  Enable HSTS: %t\n", c.Security.EnableHSTS)
	fmt.Printf("  HSTS Max Age: %d\n", c.Security.HSTSMaxAge)
	
	// Logging Configuration
	fmt.Printf("Logging:\n")
	fmt.Printf("  Level: %s\n", c.Logging.Level)
	fmt.Printf("  Format: %s\n", c.Logging.Format)
	fmt.Printf("  Output: %s\n", c.Logging.Output)
	fmt.Printf("  Enable File: %t\n", c.Logging.EnableFile)
	if c.Logging.EnableFile {
		fmt.Printf("  File Path: %s\n", c.Logging.FilePath)
	}
	
	// Environment Detection
	fmt.Printf("Environment:\n")
	fmt.Printf("  Production Mode: %t\n", c.IsProduction())
	fmt.Printf("  Development Mode: %t\n", c.IsDevelopment())
	
	fmt.Println("=============================")
}