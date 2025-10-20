package config

import "time"

// Server configuration methods
func (cs *configService) GetServerHost() string {
	return cs.config.Server.Host
}

func (cs *configService) GetServerPort() int {
	return cs.config.Server.Port
}

func (cs *configService) GetServerBaseURL() string {
	return cs.config.Server.BaseURL
}

func (cs *configService) GetServerReadTimeout() time.Duration {
	return cs.config.Server.ReadTimeout
}

func (cs *configService) GetServerWriteTimeout() time.Duration {
	return cs.config.Server.WriteTimeout
}

func (cs *configService) GetServerIdleTimeout() time.Duration {
	return cs.config.Server.IdleTimeout
}

func (cs *configService) GetServerShutdownTimeout() time.Duration {
	return cs.config.Server.ShutdownTimeout
}

// I18n configuration methods
func (cs *configService) GetSupportedLocales() []string {
	return cs.config.I18n.SupportedLocales
}

func (cs *configService) GetDefaultLocale() string {
	return cs.config.I18n.DefaultLocale
}

func (cs *configService) GetFallbackLocale() string {
	return cs.config.I18n.FallbackLocale
}

// Layout configuration methods
func (cs *configService) GetLayoutRootDirectory() string {
	return cs.config.Layout.RootDirectory
}

func (cs *configService) GetLayoutAssetsDirectory() string {
	return cs.config.Layout.AssetsDirectory
}

func (cs *configService) GetLayoutAssetsRouteName() string {
	return cs.config.Layout.AssetsRouteName
}

func (cs *configService) GetLayoutFileName() string {
	return cs.config.Layout.LayoutFileName
}

func (cs *configService) GetTemplateExtension() string {
	return cs.config.Layout.TemplateExtension
}

func (cs *configService) GetMetadataExtension() string {
	return cs.config.Layout.MetadataExtension
}

func (cs *configService) IsLayoutInheritanceEnabled() bool {
	return cs.config.Layout.EnableInheritance
}

func (cs *configService) GetTemplateOutputDir() string {
	return cs.config.TemplateGenerator.OutputDir
}

func (cs *configService) GetTemplatePackageName() string {
	return cs.config.TemplateGenerator.PackageName
}

// Database configuration methods
func (cs *configService) GetDatabaseHost() string {
	return cs.config.Database.Host
}

func (cs *configService) GetDatabasePort() int {
	return cs.config.Database.Port
}

func (cs *configService) GetDatabaseUser() string {
	return cs.config.Database.User
}

func (cs *configService) GetDatabasePassword() string {
	return cs.config.Database.Password
}

func (cs *configService) GetDatabaseName() string {
	return cs.config.Database.Name
}

func (cs *configService) GetDatabaseSSLMode() string {
	return cs.config.Database.SSLMode
}

// Auth configuration methods
func (cs *configService) IsEmailVerificationRequired() bool {
	return cs.config.Auth.RequireEmailVerification
}

func (cs *configService) GetVerificationTokenExpiry() time.Duration {
	return cs.config.Auth.VerificationTokenExpiry
}

func (cs *configService) GetSessionCookieName() string {
	return cs.config.Auth.SessionCookieName
}

func (cs *configService) GetSessionExpiry() time.Duration {
	return cs.config.Auth.SessionExpiry
}

func (cs *configService) IsSessionSecure() bool {
	return cs.config.Auth.SessionSecure
}

func (cs *configService) IsSessionHttpOnly() bool {
	return cs.config.Auth.SessionHttpOnly
}

func (cs *configService) GetSessionSameSite() string {
	return cs.config.Auth.SessionSameSite
}

func (cs *configService) GetMinPasswordLength() int {
	return cs.config.Auth.MinPasswordLength
}

func (cs *configService) IsStrongPasswordRequired() bool {
	return cs.config.Auth.RequireStrongPasswd
}

func (cs *configService) ShouldCreateDefaultAdmin() bool {
	return cs.config.Auth.CreateDefaultAdmin
}

func (cs *configService) GetDefaultAdminEmail() string {
	return cs.config.Auth.DefaultAdminEmail
}

func (cs *configService) GetDefaultAdminPassword() string {
	return cs.config.Auth.DefaultAdminPassword
}

func (cs *configService) GetDefaultAdminFirstName() string {
	return cs.config.Auth.DefaultAdminFirstName
}

func (cs *configService) GetDefaultAdminLastName() string {
	return cs.config.Auth.DefaultAdminLastName
}

// Auth redirect routes (only for success cases)
func (cs *configService) GetSignInSuccessRoute() string {
	return cs.config.Auth.SignInSuccessRoute
}

func (cs *configService) GetSignUpSuccessRoute() string {
	return cs.config.Auth.SignUpSuccessRoute
}

func (cs *configService) GetSignOutSuccessRoute() string {
	return cs.config.Auth.SignOutSuccessRoute
}

// Security configuration methods
func (cs *configService) GetCSRFSecret() string {
	return cs.config.Security.CSRFSecret
}

func (cs *configService) IsCSRFSecure() bool {
	return cs.config.Security.CSRFSecure
}

func (cs *configService) IsCSRFHttpOnly() bool {
	return cs.config.Security.CSRFHttpOnly
}

func (cs *configService) GetCSRFSameSite() string {
	return cs.config.Security.CSRFSameSite
}

func (cs *configService) IsRateLimitEnabled() bool {
	return cs.config.Security.EnableRateLimit
}

func (cs *configService) GetRateLimitRequests() int {
	return cs.config.Security.RateLimitRequests
}

func (cs *configService) AreSecurityHeadersEnabled() bool {
	return cs.config.Security.EnableSecurityHeaders
}

func (cs *configService) IsHSTSEnabled() bool {
	return cs.config.Security.EnableHSTS
}

func (cs *configService) GetHSTSMaxAge() int {
	return cs.config.Security.HSTSMaxAge
}

// Logging configuration methods
func (cs *configService) GetLogLevel() string {
	return cs.config.Logging.Level
}

func (cs *configService) GetLogFormat() string {
	return cs.config.Logging.Format
}

func (cs *configService) GetLogOutput() string {
	return cs.config.Logging.Output
}

func (cs *configService) IsFileLoggingEnabled() bool {
	return cs.config.Logging.EnableFile
}

func (cs *configService) GetLogFilePath() string {
	return cs.config.Logging.FilePath
}

// Email configuration methods
func (cs *configService) GetSMTPHost() string {
	return cs.config.Email.SMTPHost
}

func (cs *configService) GetSMTPPort() int {
	return cs.config.Email.SMTPPort
}

func (cs *configService) GetSMTPUsername() string {
	return cs.config.Email.SMTPUsername
}

func (cs *configService) GetSMTPPassword() string {
	return cs.config.Email.SMTPPassword
}

func (cs *configService) IsSMTPTLSEnabled() bool {
	return cs.config.Email.SMTPUseTLS
}

func (cs *configService) GetFromEmail() string {
	return cs.config.Email.FromEmail
}

func (cs *configService) GetFromName() string {
	return cs.config.Email.FromName
}

func (cs *configService) GetReplyToEmail() string {
	return cs.config.Email.ReplyToEmail
}

func (cs *configService) IsEmailDummyModeEnabled() bool {
	return cs.config.Email.EnableDummyMode
}

func (cs *configService) IsDevelopment() bool {
	return cs.config.IsDevelopment()
}

func (cs *configService) IsProduction() bool {
	return cs.config.IsProduction()
}
