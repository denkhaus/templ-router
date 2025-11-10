package config

import (
	"bytes"
	"os"
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
