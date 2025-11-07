package interfaces

import "time"

// ServerConfigService provides access to server configuration
type ServerConfigService interface {
	GetServerHost() string
	GetServerPort() int
	GetServerBaseURL() string
	GetServerReadTimeout() time.Duration
	GetServerWriteTimeout() time.Duration
	GetServerIdleTimeout() time.Duration
	GetServerShutdownTimeout() time.Duration
}

// I18nConfigService provides access to internationalization configuration
type I18nConfigService interface {
	GetSupportedLocales() []string
	GetDefaultLocale() string
	GetFallbackLocale() string
}

// LayoutConfigService provides access to layout and template configuration
type LayoutConfigService interface {
	GetLayoutRootDirectory() string
	GetLayoutFileName() string
	GetLayoutAssetsDirectory() string
	GetLayoutAssetsRouteName() string
	GetTemplateExtension() string
	GetMetadataExtension() string
	IsLayoutInheritanceEnabled() bool
}

// TemplateGeneratorConfigService provides access to template generator configuration
type TemplateGeneratorConfigService interface {
	GetTemplateOutputDir() string
	GetTemplatePackageName() string
}

// DatabaseConfigService provides access to database configuration
type DatabaseConfigService interface {
	GetDatabaseHost() string
	GetDatabasePort() int
	GetDatabaseUser() string
	GetDatabasePassword() string
	GetDatabaseName() string
	GetDatabaseSSLMode() string
}

// AuthConfigService provides access to authentication configuration
type AuthConfigService interface {
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
}

// AuthRoutesConfigService provides access to authentication route configuration
type AuthRoutesConfigService interface {
	GetSignInRoute() string
	GetSignInSuccessRoute() string
	GetSignUpSuccessRoute() string
	GetSignOutSuccessRoute() string
}

// SecurityConfigService provides access to security configuration
type SecurityConfigService interface {
	GetCSRFSecret() string
	IsCSRFSecure() bool
	IsCSRFHttpOnly() bool
	GetCSRFSameSite() string
	IsRateLimitEnabled() bool
	GetRateLimitRequests() int
	AreSecurityHeadersEnabled() bool
	IsHSTSEnabled() bool
	GetHSTSMaxAge() int
}

// LoggingConfigService provides access to logging configuration
type LoggingConfigService interface {
	GetLogLevel() string
	GetLogFormat() string
	GetLogOutput() string
	IsFileLoggingEnabled() bool
	GetLogFilePath() string
}

// EmailConfigService provides access to email configuration
type EmailConfigService interface {
	GetSMTPHost() string
	GetSMTPPort() int
	GetSMTPUsername() string
	GetSMTPPassword() string
	IsSMTPTLSEnabled() bool
	GetFromEmail() string
	GetFromName() string
	GetReplyToEmail() string
	IsEmailDummyModeEnabled() bool
}

// EnvironmentConfigService provides access to environment configuration
type EnvironmentConfigService interface {
	IsDevelopment() bool
	IsProduction() bool
}

// RouterConfigService provides access to router configuration
type RouterConfigService interface {
	GetRouterEnableTrailingSlash() bool
	GetRouterEnableSlashRedirect() bool
	GetRouterEnableMethodNotAllowed() bool
	GetRouterEnableAuthRoutes() bool
	GetRouterAuthRoutePrefix() string
}

// ConfigService combines all configuration interfaces for backward compatibility
// Note: This interface is maintained for backward compatibility but should be gradually
// replaced by using specific interfaces as needed
type ConfigService interface {
	ServerConfigService
	I18nConfigService
	LayoutConfigService
	TemplateGeneratorConfigService
	DatabaseConfigService
	AuthConfigService
	AuthRoutesConfigService
	SecurityConfigService
	LoggingConfigService
	EmailConfigService
	EnvironmentConfigService
	RouterConfigService
}
