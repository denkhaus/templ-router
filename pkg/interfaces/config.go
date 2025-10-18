package interfaces

import "time"

// ConfigService provides access to application configuration
type ConfigService interface {
	// Server configuration
	GetServerHost() string
	GetServerPort() int
	GetServerBaseURL() string
	GetServerReadTimeout() time.Duration
	GetServerWriteTimeout() time.Duration
	GetServerIdleTimeout() time.Duration
	GetServerShutdownTimeout() time.Duration

	// I18n configuration
	GetSupportedLocales() []string
	GetDefaultLocale() string
	GetFallbackLocale() string

	// Layout configuration
	GetLayoutRootDirectory() string
	GetLayoutFileName() string
	GetLayoutAssetsDirectory() string
	GetLayoutAssetsRouteName() string
	GetTemplateExtension() string
	GetMetadataExtension() string
	IsLayoutInheritanceEnabled() bool

	// Template generator configuration
	GetTemplateOutputDir() string
	GetTemplatePackageName() string

	// Database configuration
	GetDatabaseHost() string
	GetDatabasePort() int
	GetDatabaseUser() string
	GetDatabasePassword() string
	GetDatabaseName() string
	GetDatabaseSSLMode() string

	// Auth configuration
	IsEmailVerificationRequired() bool
	GetVerificationTokenExpiry() time.Duration
	GetSessionCookieName() string
	GetSessionExpiry() time.Duration
	IsSessionSecure() bool
	IsSessionHttpOnly() bool
	GetSessionSameSite() string
	GetMinPasswordLength() int
	IsStrongPasswordRequired() bool
	ShouldCreateDefaultAdmin() bool
	GetDefaultAdminEmail() string
	GetDefaultAdminPassword() string
	GetDefaultAdminFirstName() string
	GetDefaultAdminLastName() string

	// Security configuration
	GetCSRFSecret() string
	IsCSRFSecure() bool
	IsCSRFHttpOnly() bool
	GetCSRFSameSite() string
	IsRateLimitEnabled() bool
	GetRateLimitRequests() int
	AreSecurityHeadersEnabled() bool
	IsHSTSEnabled() bool
	GetHSTSMaxAge() int

	// Logging configuration
	GetLogLevel() string
	GetLogFormat() string
	GetLogOutput() string
	IsFileLoggingEnabled() bool
	GetLogFilePath() string

	// Email configuration
	GetSMTPHost() string
	GetSMTPPort() int
	GetSMTPUsername() string
	GetSMTPPassword() string
	IsSMTPTLSEnabled() bool
	GetFromEmail() string
	GetFromName() string
	GetReplyToEmail() string
	IsEmailDummyModeEnabled() bool

	// Environment configuration
	IsDevelopment() bool
	IsProduction() bool
}
