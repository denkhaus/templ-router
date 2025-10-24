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
				"Database:",
				"Host: localhost",
				"Port: 5432",
				"User: postgres",
				"Name: router_db",
				"SSL Mode: disable",
				"Authentication:",
				"Require Email Verification: true",
				"Min Password Length: 8",
				"Create Default Admin: true",
				"Default Admin Email: admin@example.com",
				"Email:",
				"SMTP Host: <not configured>",
				"SMTP Port: 587",
				"SMTP Username: <not configured>",
				"SMTP Password: <not configured>",
				"From Email: noreply@example.com",
				"From Name: Router Application",
				"Reply To Email: <not set>",
				"Enable Dummy Mode: true",
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
				"Password: po****es", // postgres password masked
				"CSRF Secret: ch*******************on", // change-me-in-production masked
				"Default Admin Password: ad****23", // admin123 masked
			},
		},
		{
			name: "custom configuration with sensitive data",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":          "production",
				"TR_SERVER_HOST":               "prod.example.com",
				"TR_SERVER_PORT":               "443",
				"TR_SERVER_BASE_URL":           "https://prod.example.com",
				"TR_DATABASE_HOST":             "db.prod.com",
				"TR_DATABASE_PASSWORD":         "supersecretdbpass",
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "true",
				"TR_AUTH_DEFAULT_ADMIN_EMAIL":  "admin@prod.com",
				"TR_AUTH_DEFAULT_ADMIN_PASSWORD": "verysecureadminpass",
				"TR_SECURITY_CSRF_SECRET":      "production-csrf-secret-key",
				"TR_EMAIL_SMTP_HOST":           "smtp.prod.com",
				"TR_EMAIL_SMTP_USERNAME":       "smtp@prod.com",
				"TR_EMAIL_SMTP_PASSWORD":       "smtpsecretpass",
				"TR_EMAIL_REPLY_TO_EMAIL":      "support@prod.com",
				"TR_LOGGING_ENABLE_FILE":       "true",
				"TR_LOGGING_FILE_PATH":         "/var/log/app.log",
			},
			expectInLog: []string{
				"Host: prod.example.com",
				"Port: 443",
				"Base URL: https://prod.example.com",
				"Host: db.prod.com",
				"Default Admin Email: admin@prod.com",
				"SMTP Host: smtp.prod.com",
				"Reply To Email: support@prod.com",
				"Enable File: true",
				"File Path: /var/log/app.log",
				"Production Mode: true",
				"Development Mode: false",
			},
			expectMasked: []string{
				"Password: su*************ss", // supersecretdbpass masked
				"Default Admin Password: ve***************ss", // verysecureadminpass masked
				"CSRF Secret: pr**********************ey", // production-csrf-secret-key masked
				"SMTP Username: sm*********om", // smtp@prod.com masked
				"SMTP Password: sm**********ss", // smtpsecretpass masked
			},
		},
		{
			name: "configuration with empty sensitive fields",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "false",
				"TR_EMAIL_SMTP_HOST":           "",
				"TR_EMAIL_SMTP_USERNAME":       "",
				"TR_EMAIL_SMTP_PASSWORD":       "",
				"TR_EMAIL_REPLY_TO_EMAIL":      "",
			},
			expectInLog: []string{
				"Create Default Admin: false",
				"SMTP Host: <not configured>",
				"SMTP Username: <not configured>",
				"SMTP Password: <not configured>",
				"Reply To Email: <not set>",
			},
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
			w.Close()
			os.Stdout = oldStdout
			
			var buf bytes.Buffer
			buf.ReadFrom(r)
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
	os.Setenv("TR_CONFIG_PRINT_SUMMARY", "true")

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
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// LogSummary should have been called automatically during service creation
	assert.Contains(t, output, "=== Configuration Summary ===")
	assert.Contains(t, output, "Host: localhost")
	assert.Contains(t, output, "=============================")
}

func TestLogSummaryWithPrintSummaryDisabled(t *testing.T) {
	clearTestEnv(t)
	
	// Disable print summary (default is false, but being explicit)
	os.Setenv("TR_CONFIG_PRINT_SUMMARY", "false")

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
	w.Close()
	os.Stdout = oldStdout
	
	var buf bytes.Buffer
	buf.ReadFrom(r)
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
			expectInLog: []string{
				"Create Default Admin: true",
				"Default Admin Email: admin@example.com",
				"Default Admin Password: ad****23",
				"Default Admin First Name: Default",
				"Default Admin Last Name: Admin",
			},
		},
		{
			name: "default admin disabled - hides admin fields",
			envVars: map[string]string{
				"TR_AUTH_CREATE_DEFAULT_ADMIN": "false",
			},
			expectInLog: []string{
				"Create Default Admin: false",
			},
			notInLog: []string{
				"Default Admin Email:",
				"Default Admin Password:",
				"Default Admin First Name:",
				"Default Admin Last Name:",
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
				os.Setenv(key, value)
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
			w.Close()
			os.Stdout = oldStdout
			
			var buf bytes.Buffer
			buf.ReadFrom(r)
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