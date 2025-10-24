package config

import (
	"os"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid default configuration",
			envVars:     map[string]string{},
			expectError: false,
		},
		{
			name: "valid custom configuration",
			envVars: map[string]string{
				"TR_SERVER_PORT":                "3000",
				"TR_DATABASE_PORT":              "5432",
				"TR_AUTH_MIN_PASSWORD_LENGTH":   "10",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD": "validpassword123",
				"TR_EMAIL_SMTP_PORT":            "587",
				"TR_SECURITY_RATE_LIMIT_REQUESTS": "50",
			},
			expectError: false,
		},
		// Server port validation tests
		{
			name: "invalid server port - too low",
			envVars: map[string]string{
				"TR_SERVER_PORT": "0",
			},
			expectError: true,
			errorMsg:    "invalid server port: 0",
		},
		{
			name: "invalid server port - too high",
			envVars: map[string]string{
				"TR_SERVER_PORT": "65536",
			},
			expectError: true,
			errorMsg:    "invalid server port: 65536",
		},
		{
			name: "valid server port - minimum",
			envVars: map[string]string{
				"TR_SERVER_PORT": "1",
			},
			expectError: false,
		},
		{
			name: "valid server port - maximum",
			envVars: map[string]string{
				"TR_SERVER_PORT": "65535",
			},
			expectError: false,
		},
		// Database port validation tests
		{
			name: "invalid database port - too low",
			envVars: map[string]string{
				"TR_DATABASE_PORT": "0",
			},
			expectError: true,
			errorMsg:    "invalid database port: 0",
		},
		{
			name: "invalid database port - too high",
			envVars: map[string]string{
				"TR_DATABASE_PORT": "65536",
			},
			expectError: true,
			errorMsg:    "invalid database port: 65536",
		},
		{
			name: "valid database port - minimum",
			envVars: map[string]string{
				"TR_DATABASE_PORT": "1",
			},
			expectError: false,
		},
		{
			name: "valid database port - maximum",
			envVars: map[string]string{
				"TR_DATABASE_PORT": "65535",
			},
			expectError: false,
		},
		// Auth password length validation tests
		{
			name: "invalid min password length - zero",
			envVars: map[string]string{
				"TR_AUTH_MIN_PASSWORD_LENGTH": "0",
			},
			expectError: true,
			errorMsg:    "minimum password length must be at least 1",
		},
		{
			name: "invalid min password length - negative",
			envVars: map[string]string{
				"TR_AUTH_MIN_PASSWORD_LENGTH": "-1",
			},
			expectError: true,
			errorMsg:    "minimum password length must be at least 1",
		},
		{
			name: "valid min password length - minimum",
			envVars: map[string]string{
				"TR_AUTH_MIN_PASSWORD_LENGTH": "1",
			},
			expectError: false,
		},
		// Default admin validation tests
		{
			name: "default admin enabled with empty email",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "true",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":  "",
			},
			expectError: true,
			errorMsg:    "default admin email cannot be empty when CreateDefaultAdmin is enabled",
		},
		{
			name: "default admin enabled with empty password",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN":   "true",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":    "admin@test.com",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD": "",
			},
			expectError: true,
			errorMsg:    "default admin password cannot be empty when CreateDefaultAdmin is enabled",
		},
		{
			name: "default admin password too short",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN":    "true",
				"TR_AUTH_MIN_PASSWORD_LENGTH":     "10",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":     "admin@test.com",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD":  "short",
				"TR_AUTH_DEFAULT_ADMIN_FIRST_NAME": "Admin",
				"TR_AUTH_DEFAULT_ADMIN_LAST_NAME":  "User",
			},
			expectError: true,
			errorMsg:    "default admin password must meet minimum password length requirement (10 characters)",
		},
		{
			name: "default admin enabled with empty first name",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN":     "true",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":      "admin@test.com",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD":   "validpassword",
				"TR_AUTH_DEFAULT_ADMIN_FIRST_NAME": "",
				"TR_AUTH_DEFAULT_ADMIN_LAST_NAME":  "User",
			},
			expectError: true,
			errorMsg:    "default admin first name cannot be empty when CreateDefaultAdmin is enabled",
		},
		{
			name: "default admin enabled with empty last name",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN":     "true",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":      "admin@test.com",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD":   "validpassword",
				"TR_AUTH_DEFAULT_ADMIN_FIRST_NAME": "Admin",
				"TR_AUTH_DEFAULT_ADMIN_LAST_NAME":  "",
			},
			expectError: true,
			errorMsg:    "default admin last name cannot be empty when CreateDefaultAdmin is enabled",
		},
		{
			name: "valid default admin configuration",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN":     "true",
				"TR_AUTH_MIN_PASSWORD_LENGTH":      "8",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":      "admin@test.com",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD":   "validpassword123",
				"TR_AUTH_DEFAULT_ADMIN_FIRST_NAME": "Admin",
				"TR_AUTH_DEFAULT_ADMIN_LAST_NAME":  "User",
			},
			expectError: false,
		},
		{
			name: "default admin disabled - validation skipped",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "false",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":  "",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD": "",
			},
			expectError: false,
		},
		// Email SMTP port validation tests
		{
			name: "invalid SMTP port - too low",
			envVars: map[string]string{
				"TR_EMAIL_SMTP_PORT": "0",
			},
			expectError: true,
			errorMsg:    "invalid SMTP port: 0",
		},
		{
			name: "invalid SMTP port - too high",
			envVars: map[string]string{
				"TR_EMAIL_SMTP_PORT": "65536",
			},
			expectError: true,
			errorMsg:    "invalid SMTP port: 65536",
		},
		{
			name: "valid SMTP port - minimum",
			envVars: map[string]string{
				"TR_EMAIL_SMTP_PORT": "1",
			},
			expectError: false,
		},
		{
			name: "valid SMTP port - maximum",
			envVars: map[string]string{
				"TR_EMAIL_SMTP_PORT": "65535",
			},
			expectError: false,
		},
		// Security rate limit validation tests
		{
			name: "invalid rate limit - zero",
			envVars: map[string]string{
				"TR_SECURITY_RATE_LIMIT_REQUESTS": "0",
			},
			expectError: true,
			errorMsg:    "rate limit requests must be at least 1",
		},
		{
			name: "invalid rate limit - negative",
			envVars: map[string]string{
				"TR_SECURITY_RATE_LIMIT_REQUESTS": "-1",
			},
			expectError: true,
			errorMsg:    "rate limit requests must be at least 1",
		},
		{
			name: "valid rate limit - minimum",
			envVars: map[string]string{
				"TR_SECURITY_RATE_LIMIT_REQUESTS": "1",
			},
			expectError: false,
		},
		// Multiple validation errors (should return first error)
		{
			name: "multiple validation errors",
			envVars: map[string]string{
				"TR_SERVER_PORT":                   "0",
				"TR_DATABASE_PORT":                 "0",
				"TR_AUTH_MIN_PASSWORD_LENGTH":      "0",
			},
			expectError: true,
			errorMsg:    "invalid server port: 0", // Should return first validation error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)
			
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			_, err := configFactory(injector)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDirectValidation(t *testing.T) {
	// Test the Validate method directly on configImpl
	tests := []struct {
		name        string
		config      configImpl
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			config: configImpl{
				Server: ServerConfig{
					Port: 8080,
				},
				Database: DatabaseConfig{
					Port: 5432,
				},
				Auth: AuthConfig{
					MinPasswordLength:   8,
					CreateDefaultAdmin:  false,
				},
				Email: EmailConfig{
					SMTPPort: 587,
				},
				Security: SecurityConfig{
					RateLimitRequests: 100,
				},
			},
			expectError: false,
		},
		{
			name: "invalid server port",
			config: configImpl{
				Server: ServerConfig{
					Port: 0,
				},
				Database: DatabaseConfig{
					Port: 5432,
				},
				Auth: AuthConfig{
					MinPasswordLength: 8,
				},
				Email: EmailConfig{
					SMTPPort: 587,
				},
				Security: SecurityConfig{
					RateLimitRequests: 100,
				},
			},
			expectError: true,
			errorMsg:    "invalid server port: 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}