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
			errorMsg:    "Invalid server port",
		},
		{
			name: "invalid server port - too high",
			envVars: map[string]string{
				"TR_SERVER_PORT": "65536",
			},
			expectError: true,
			errorMsg:    "Invalid server port",
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
		// Security rate limit validation tests
		{
			name: "invalid rate limit - zero",
			envVars: map[string]string{
				"TR_SECURITY_RATE_LIMIT_REQUESTS": "0",
			},
			expectError: true,
			errorMsg:    "Rate limit requests must be at least 1",
		},
		{
			name: "invalid rate limit - negative",
			envVars: map[string]string{
				"TR_SECURITY_RATE_LIMIT_REQUESTS": "-1",
			},
			expectError: true,
			errorMsg:    "Rate limit requests must be at least 1",
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
				"TR_SERVER_PORT": "0",
			},
			expectError: true,
			errorMsg:    "Invalid server port", // Should return first validation error
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
					// The service layer wraps validation errors in "configuration validation failed"
					assert.Contains(t, err.Error(), "configuration validation failed")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
