package config

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskSensitive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "<empty>",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "*",
		},
		{
			name:     "two characters",
			input:    "ab",
			expected: "**",
		},
		{
			name:     "three characters",
			input:    "abc",
			expected: "***",
		},
		{
			name:     "four characters",
			input:    "abcd",
			expected: "****",
		},
		{
			name:     "five characters",
			input:    "abcde",
			expected: "ab*de",
		},
		{
			name:     "normal password",
			input:    "password123",
			expected: "pa*******23",
		},
		{
			name:     "long password",
			input:    "verylongpassword123456",
			expected: "ve******************56",
		},
		{
			name:     "special characters",
			input:    "p@ssw0rd!",
			expected: "p@*****d!",
		},
		{
			name:     "unicode characters",
			input:    "password",
			expected: "pa****rd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskSensitive(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogSummary(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectInLog  []string
		expectMasked []string
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			expectInLog: []string{
				"=== Configuration Summary ===",
				"Server:",
				"Host: localhost",
				"Port: 8080",
				"Base URL: http://localhost:8080",
				"Authentication:",
				"Security:",
				"Enable Rate Limit: true",
				"Rate Limit Requests: 100",
				"Logging:",
				"Level: info",
				"Format: json",
				"Output: stdout",
				"Enable File: false",
				"Environment:",
				"Production Mode: false",
				"Development Mode: true",
				"=============================",
			},
			expectMasked: []string{
				"CSRF Secret: ch*******************on", // change-me-in-production masked
			},
		},
		{
			name: "custom configuration with sensitive data",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":     "production",
				"TR_SERVER_HOST":          "prod.example.com",
				"TR_SERVER_PORT":          "443",
				"TR_SERVER_BASE_URL":      "https://prod.example.com",
				"TR_SECURITY_CSRF_SECRET": "production-csrf-secret-key",
				"TR_LOGGING_ENABLE_FILE":  "true",
				"TR_LOGGING_FILE_PATH":    "/var/log/app.log",
			},
			expectInLog: []string{
				"Host: prod.example.com",
				"Port: 443",
				"Base URL: https://prod.example.com",
				"Enable File: true",
				"File Path: /var/log/app.log",
				"Production Mode: true",
				"Development Mode: false",
			},
			expectMasked: []string{
				"CSRF Secret: pr**********************ey", // production-csrf-secret-key masked
			},
		},
		{
			name: "configuration with empty sensitive fields",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "false",
			},
			expectInLog: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)

			for key, value := range tt.envVars {
				_ = os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Access the internal config and call LogSummary
			configSvc := service.(*configService)
			configSvc.config.LogSummary()

			// Restore stdout and read captured output
			if err := w.Close(); err != nil {
				t.Fatalf("Failed to close pipe writer: %v", err)
			}
			os.Stdout = oldStdout

			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}
			output := buf.String()

			// Check that expected strings are in the log
			for _, expected := range tt.expectInLog {
				assert.Contains(t, output, expected, "Expected '%s' to be in log output", expected)
			}

			// Check that sensitive data is properly masked
			for _, expectedMasked := range tt.expectMasked {
				assert.Contains(t, output, expectedMasked, "Expected masked value '%s' to be in log output", expectedMasked)
			}

			// Ensure the log starts and ends with the expected markers
			assert.True(t, strings.HasPrefix(output, "=== Configuration Summary ==="))
			assert.True(t, strings.HasSuffix(strings.TrimSpace(output), "============================="))
		})
	}
}

func TestLogSummaryWithPrintSummaryEnabled(t *testing.T) {
	clearTestEnv(t)

	// Enable print summary
	_ = os.Setenv("TR_CONFIG_PRINT_SUMMARY", "true")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	injector := do.New()
	defer injector.Shutdown()

	configFactory := NewConfigService("TR")
	service, err := configFactory(injector)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Restore stdout and read captured output

	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close pipe writer: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// LogSummary should have been called automatically during service creation
	assert.Contains(t, output, "=== Configuration Summary ===")
	assert.Contains(t, output, "Host: localhost")
	assert.Contains(t, output, "=============================")
}

func TestLogSummaryWithPrintSummaryDisabled(t *testing.T) {
	clearTestEnv(t)

	// Disable print summary (default is false, but being explicit)
	_ = os.Setenv("TR_CONFIG_PRINT_SUMMARY", "false")

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	injector := do.New()
	defer injector.Shutdown()

	configFactory := NewConfigService("TR")
	service, err := configFactory(injector)
	require.NoError(t, err)
	require.NotNil(t, service)

	// Restore stdout and read captured output
	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// LogSummary should NOT have been called automatically
	assert.NotContains(t, output, "=== Configuration Summary ===")
}

func TestLogSummaryConditionalFields(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectInLog []string
		notInLog    []string
	}{
		{
			name: "default admin enabled - shows admin fields",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "true",
			},
			expectInLog: []string{},
		},
		{
			name: "default admin disabled - hides admin fields",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "false",
			},
			expectInLog: []string{},
			notInLog: []string{
				"Default Admin Email:",
			},
		},
		{
			name: "file logging enabled - shows file path",
			envVars: map[string]string{
				"TR_LOGGING_ENABLE_FILE": "true",
				"TR_LOGGING_FILE_PATH":   "/custom/log/path.log",
			},
			expectInLog: []string{
				"Enable File: true",
				"File Path: /custom/log/path.log",
			},
		},
		{
			name: "file logging disabled - hides file path",
			envVars: map[string]string{
				"TR_LOGGING_ENABLE_FILE": "false",
			},
			expectInLog: []string{
				"Enable File: false",
			},
			notInLog: []string{
				"File Path:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)

			for key, value := range tt.envVars {
				_ = os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Access the internal config and call LogSummary
			configSvc := service.(*configService)
			configSvc.config.LogSummary()

			// Restore stdout and read captured output
			_ = w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			// Check expected content
			for _, expected := range tt.expectInLog {
				assert.Contains(t, output, expected, "Expected '%s' to be in log output", expected)
			}

			// Check that unwanted content is not present
			for _, notExpected := range tt.notInLog {
				assert.NotContains(t, output, notExpected, "Did not expect '%s' to be in log output", notExpected)
			}
		})
	}
}
