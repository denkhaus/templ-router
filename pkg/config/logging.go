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
func (c *configImpl) LogSummary() {
	fmt.Println("=== Configuration Summary ===")

	// Authentication Configuration
	fmt.Printf("Authentication:\n")
	fmt.Printf("  Sign In Route: %s\n", c.Auth.SignInRoute)
	fmt.Printf("  Sign In Success Route: %s\n", c.Auth.SignInSuccessRoute)
	fmt.Printf("  Sign Up Success Route: %s\n", c.Auth.SignUpSuccessRoute)
	fmt.Printf("  Sign Out Success Route: %s\n", c.Auth.SignOutSuccessRoute)

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
