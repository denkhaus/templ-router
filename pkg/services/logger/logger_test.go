package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Mock ConfigService for logger tests
type mockLoggerConfigService struct {
	logLevel            string
	logFormat           string
	logOutput           string
	fileLoggingEnabled  bool
	logFilePath         string
	isDevelopment       bool
	isProduction        bool
}

func (m *mockLoggerConfigService) GetLogLevel() string                    { return m.logLevel }
func (m *mockLoggerConfigService) GetLogFormat() string                   { return m.logFormat }
func (m *mockLoggerConfigService) GetLogOutput() string                   { return m.logOutput }
func (m *mockLoggerConfigService) IsFileLoggingEnabled() bool             { return m.fileLoggingEnabled }
func (m *mockLoggerConfigService) GetLogFilePath() string                 { return m.logFilePath }
func (m *mockLoggerConfigService) IsDevelopment() bool                    { return m.isDevelopment }
func (m *mockLoggerConfigService) IsProduction() bool                     { return m.isProduction }

// Router configuration methods
func (m *mockLoggerConfigService) GetRouterEnableTrailingSlash() bool     { return true }
func (m *mockLoggerConfigService) GetRouterEnableSlashRedirect() bool     { return true }
func (m *mockLoggerConfigService) GetRouterEnableMethodNotAllowed() bool  { return true }

// Implement remaining interface methods with defaults
func (m *mockLoggerConfigService) GetServerHost() string                     { return "localhost" }
func (m *mockLoggerConfigService) GetServerPort() int                        { return 8080 }
func (m *mockLoggerConfigService) GetServerBaseURL() string                  { return "http://localhost:8080" }
func (m *mockLoggerConfigService) GetServerReadTimeout() time.Duration       { return 30 * time.Second }
func (m *mockLoggerConfigService) GetServerWriteTimeout() time.Duration      { return 30 * time.Second }
func (m *mockLoggerConfigService) GetServerIdleTimeout() time.Duration       { return 60 * time.Second }
func (m *mockLoggerConfigService) GetServerShutdownTimeout() time.Duration   { return 10 * time.Second }
func (m *mockLoggerConfigService) GetSupportedLocales() []string             { return []string{"en"} }
func (m *mockLoggerConfigService) GetDefaultLocale() string                  { return "en" }
func (m *mockLoggerConfigService) GetFallbackLocale() string                 { return "en" }
func (m *mockLoggerConfigService) GetLayoutRootDirectory() string            { return "app" }
func (m *mockLoggerConfigService) GetLayoutFileName() string                 { return "layout.templ" }
func (m *mockLoggerConfigService) GetLayoutAssetsDirectory() string          { return "assets" }
func (m *mockLoggerConfigService) GetLayoutAssetsRouteName() string          { return "/assets/" }
func (m *mockLoggerConfigService) GetTemplateExtension() string              { return ".templ" }
func (m *mockLoggerConfigService) GetMetadataExtension() string              { return ".yaml" }
func (m *mockLoggerConfigService) IsLayoutInheritanceEnabled() bool          { return true }
func (m *mockLoggerConfigService) GetTemplateOutputDir() string              { return "generated" }
func (m *mockLoggerConfigService) GetTemplatePackageName() string            { return "templates" }
func (m *mockLoggerConfigService) GetDatabaseHost() string                   { return "localhost" }
func (m *mockLoggerConfigService) GetDatabasePort() int                      { return 5432 }
func (m *mockLoggerConfigService) GetDatabaseUser() string                   { return "user" }
func (m *mockLoggerConfigService) GetDatabasePassword() string               { return "password" }
func (m *mockLoggerConfigService) GetDatabaseName() string                   { return "testdb" }
func (m *mockLoggerConfigService) GetDatabaseSSLMode() string                { return "disable" }
func (m *mockLoggerConfigService) IsEmailVerificationRequired() bool         { return false }
func (m *mockLoggerConfigService) GetVerificationTokenExpiry() time.Duration { return 24 * time.Hour }
func (m *mockLoggerConfigService) GetSessionCookieName() string              { return "session" }
func (m *mockLoggerConfigService) GetSessionExpiry() time.Duration           { return 24 * time.Hour }
func (m *mockLoggerConfigService) IsSessionSecure() bool                     { return false }
func (m *mockLoggerConfigService) IsSessionHttpOnly() bool                   { return true }
func (m *mockLoggerConfigService) GetSessionSameSite() string                { return "Lax" }
func (m *mockLoggerConfigService) GetMinPasswordLength() int                 { return 8 }
func (m *mockLoggerConfigService) IsStrongPasswordRequired() bool            { return false }
func (m *mockLoggerConfigService) ShouldCreateDefaultAdmin() bool            { return false }
func (m *mockLoggerConfigService) GetDefaultAdminEmail() string              { return "" }
func (m *mockLoggerConfigService) GetDefaultAdminPassword() string           { return "" }
func (m *mockLoggerConfigService) GetDefaultAdminFirstName() string          { return "" }
func (m *mockLoggerConfigService) GetDefaultAdminLastName() string           { return "" }
func (m *mockLoggerConfigService) GetCSRFSecret() string                     { return "secret" }
func (m *mockLoggerConfigService) IsCSRFSecure() bool                        { return false }
func (m *mockLoggerConfigService) IsCSRFHttpOnly() bool                      { return true }
func (m *mockLoggerConfigService) GetCSRFSameSite() string                   { return "Lax" }
func (m *mockLoggerConfigService) IsRateLimitEnabled() bool                  { return false }
func (m *mockLoggerConfigService) GetRateLimitRequests() int                 { return 100 }
func (m *mockLoggerConfigService) AreSecurityHeadersEnabled() bool           { return false }
func (m *mockLoggerConfigService) IsHSTSEnabled() bool                       { return false }
func (m *mockLoggerConfigService) GetHSTSMaxAge() int                        { return 31536000 }
func (m *mockLoggerConfigService) GetSMTPHost() string                       { return "" }
func (m *mockLoggerConfigService) GetSMTPPort() int                          { return 587 }
func (m *mockLoggerConfigService) GetSMTPUsername() string                   { return "" }
func (m *mockLoggerConfigService) GetSMTPPassword() string                   { return "" }
func (m *mockLoggerConfigService) IsSMTPTLSEnabled() bool                    { return true }
func (m *mockLoggerConfigService) GetFromEmail() string                      { return "" }
func (m *mockLoggerConfigService) GetFromName() string                       { return "" }
func (m *mockLoggerConfigService) GetReplyToEmail() string                   { return "" }
func (m *mockLoggerConfigService) IsEmailDummyModeEnabled() bool             { return true }

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected zapcore.Level
	}{
		{"debug", "debug", zapcore.DebugLevel},
		{"info", "info", zapcore.InfoLevel},
		{"warn", "warn", zapcore.WarnLevel},
		{"warning", "warning", zapcore.WarnLevel},
		{"error", "error", zapcore.ErrorLevel},
		{"fatal", "fatal", zapcore.FatalLevel},
		{"panic", "panic", zapcore.PanicLevel},
		{"uppercase", "DEBUG", zapcore.DebugLevel},
		{"mixed case", "WaRn", zapcore.WarnLevel},
		{"invalid", "invalid", zapcore.InfoLevel},
		{"empty", "", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateEncoder(t *testing.T) {
	tests := []struct {
		name   string
		format string
		isJSON bool
	}{
		{"json format", "json", true},
		{"text format", "text", false},
		{"console format", "console", false},
		{"uppercase", "JSON", true},
		{"mixed case", "TeXt", false},
		{"invalid format", "invalid", true}, // defaults to JSON
		{"empty format", "", true},          // defaults to JSON
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := createEncoder(tt.format)
			if encoder == nil {
				t.Fatal("createEncoder() returned nil")
			}

			// Test that encoder can encode a log entry
			entry := zapcore.Entry{
				Level:   zapcore.InfoLevel,
				Time:    time.Now(),
				Message: "test message",
			}

			buf, err := encoder.EncodeEntry(entry, nil)
			if err != nil {
				t.Errorf("Failed to encode entry: %v", err)
			}

			output := buf.String()
			if tt.isJSON {
				// JSON output should contain quotes and braces
				if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
					t.Errorf("Expected JSON format but got: %s", output)
				}
			} else {
				// Console output should be more human-readable
				if !strings.Contains(output, "test message") {
					t.Errorf("Expected message in output but got: %s", output)
				}
			}
		})
	}
}

func TestCreateWriteSyncer(t *testing.T) {
	tests := []struct {
		name           string
		config         *mockLoggerConfigService
		expectError    bool
		expectFileLog  bool
	}{
		{
			name: "stdout only",
			config: &mockLoggerConfigService{
				logOutput:          "stdout",
				fileLoggingEnabled: false,
			},
			expectError:   false,
			expectFileLog: false,
		},
		{
			name: "stderr only",
			config: &mockLoggerConfigService{
				logOutput:          "stderr",
				fileLoggingEnabled: false,
			},
			expectError:   false,
			expectFileLog: false,
		},
		{
			name: "file logging enabled",
			config: &mockLoggerConfigService{
				logOutput:          "stdout",
				fileLoggingEnabled: true,
				logFilePath:        filepath.Join(os.TempDir(), "test_log.log"),
			},
			expectError:   false,
			expectFileLog: true,
		},
		{
			name: "invalid file path",
			config: &mockLoggerConfigService{
				logOutput:          "stdout",
				fileLoggingEnabled: true,
				logFilePath:        "/invalid/path/that/does/not/exist/test.log",
			},
			expectError:   true,
			expectFileLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing test log file
			if tt.config.fileLoggingEnabled && !tt.expectError {
				defer os.Remove(tt.config.logFilePath)
			}

			syncer, err := createWriteSyncer(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if syncer == nil {
				t.Fatal("createWriteSyncer() returned nil syncer")
			}

			// Test that we can write to the syncer
			testMessage := []byte("test log message\n")
			n, err := syncer.Write(testMessage)
			if err != nil {
				t.Errorf("Failed to write to syncer: %v", err)
			}
			if n != len(testMessage) {
				t.Errorf("Expected to write %d bytes, wrote %d", len(testMessage), n)
			}

			// If file logging is enabled, verify file was created
			if tt.expectFileLog {
				if _, err := os.Stat(tt.config.logFilePath); os.IsNotExist(err) {
					t.Error("Expected log file to be created but it doesn't exist")
				}
			}
		})
	}
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name           string
		config         *mockLoggerConfigService
		expectError    bool
		expectCaller   bool
	}{
		{
			name: "development logger",
			config: &mockLoggerConfigService{
				logLevel:           "debug",
				logFormat:          "console",
				logOutput:          "stdout",
				fileLoggingEnabled: false,
				isDevelopment:      true,
				isProduction:       false,
			},
			expectError:  false,
			expectCaller: true,
		},
		{
			name: "production logger",
			config: &mockLoggerConfigService{
				logLevel:           "info",
				logFormat:          "json",
				logOutput:          "stdout",
				fileLoggingEnabled: false,
				isDevelopment:      false,
				isProduction:       true,
			},
			expectError:  false,
			expectCaller: false,
		},
		{
			name: "file logging enabled",
			config: &mockLoggerConfigService{
				logLevel:           "warn",
				logFormat:          "json",
				logOutput:          "stdout",
				fileLoggingEnabled: true,
				logFilePath:        filepath.Join(os.TempDir(), "service_test.log"),
				isDevelopment:      false,
				isProduction:       true,
			},
			expectError:  false,
			expectCaller: false,
		},
		{
			name: "invalid file path",
			config: &mockLoggerConfigService{
				logLevel:           "info",
				logFormat:          "json",
				logOutput:          "stdout",
				fileLoggingEnabled: true,
				logFilePath:        "/invalid/path/service_test.log",
				isDevelopment:      false,
				isProduction:       true,
			},
			expectError:  true,
			expectCaller: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing test log file
			if tt.config.fileLoggingEnabled && !tt.expectError {
				defer os.Remove(tt.config.logFilePath)
			}

			// Create DI container
			injector := do.New()
			defer injector.Shutdown()

			// Register mock config service
			do.ProvideValue[interfaces.ConfigService](injector, tt.config)

			// Test NewService
			logger, err := NewService(injector)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if logger == nil {
				t.Fatal("NewService() returned nil logger")
			}

			// Test that logger works
			logger.Info("test message",
				zap.String("test_key", "test_value"),
				zap.Int("test_number", 42))

			// Test different log levels
			logger.Debug("debug message")
			logger.Warn("warning message")
			logger.Error("error message")

			// Sync to ensure all logs are written
			logger.Sync()

			// If file logging is enabled, verify file was created and has content
			if tt.config.fileLoggingEnabled {
				if _, err := os.Stat(tt.config.logFilePath); os.IsNotExist(err) {
					t.Error("Expected log file to be created but it doesn't exist")
				} else {
					// Check file has content
					content, err := os.ReadFile(tt.config.logFilePath)
					if err != nil {
						t.Errorf("Failed to read log file: %v", err)
					}
					if len(content) == 0 {
						t.Error("Expected log file to have content but it's empty")
					}
				}
			}
		})
	}
}

func TestNewService_InterfaceCompliance(t *testing.T) {
	// Test that NewService returns a proper zap.Logger
	injector := do.New()
	defer injector.Shutdown()

	config := &mockLoggerConfigService{
		logLevel:           "info",
		logFormat:          "json",
		logOutput:          "stdout",
		fileLoggingEnabled: false,
		isDevelopment:      true,
		isProduction:       false,
	}

	do.ProvideValue[interfaces.ConfigService](injector, config)

	logger, err := NewService(injector)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	// Verify it's a proper zap.Logger
	var _ *zap.Logger = logger

	// Test all common logger methods work
	logger.Debug("debug test")
	logger.Info("info test")
	logger.Warn("warn test")
	logger.Error("error test")

	// Test structured logging
	logger.Info("structured test",
		zap.String("string_field", "value"),
		zap.Int("int_field", 123),
		zap.Bool("bool_field", true),
		zap.Duration("duration_field", time.Second),
	)

	// Test with context
	logger.With(zap.String("context", "test")).Info("context test")

	logger.Sync()
}
// Auth redirect routes (only for success cases)
func (m *mockLoggerConfigService) GetSignInSuccessRoute() string  { return "/dashboard" }
func (m *mockLoggerConfigService) GetSignUpSuccessRoute() string  { return "/welcome" }
func (m *mockLoggerConfigService) GetSignOutSuccessRoute() string { return "/" }

// Auth routes
func (m *mockLoggerConfigService) GetSignInRoute() string { return "/login" }
