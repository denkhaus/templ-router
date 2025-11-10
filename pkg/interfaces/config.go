package interfaces

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

// AuthRoutesConfigService provides access to authentication route configuration
// Only redirect routes remain for router-level authentication failure handling
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
}

// ConfigService combines all configuration interfaces for backward compatibility
// Note: This interface is maintained for backward compatibility but should be gradually
// replaced by using specific interfaces as needed
type ConfigService interface {
	I18nConfigService
	LayoutConfigService
	TemplateGeneratorConfigService
	AuthRoutesConfigService
	SecurityConfigService
	LoggingConfigService
	EnvironmentConfigService
	RouterConfigService
}
