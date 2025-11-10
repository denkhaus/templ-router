package config

import (
	"os"
	"path/filepath"
)

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
	rootDir := cs.config.Layout.RootDirectory

	// Normalize the path - convert relative paths to absolute paths
	if !filepath.IsAbs(rootDir) {
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			// Fallback to relative path if we can't get working directory
			return rootDir
		}
		// Convert relative to absolute path
		absPath := filepath.Join(cwd, rootDir)
		return filepath.Clean(absPath)
	}

	// Already absolute, just clean it
	return filepath.Clean(rootDir)
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

// Auth routes
func (cs *configService) GetSignInRoute() string {
	return cs.config.Auth.SignInRoute
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

func (cs *configService) IsDevelopment() bool {
	return cs.config.IsDevelopment()
}

func (cs *configService) IsProduction() bool {
	return cs.config.IsProduction()
}

// Router configuration accessors

// GetRouterEnableTrailingSlash returns whether trailing slash redirection is enabled
func (cs *configService) GetRouterEnableTrailingSlash() bool {
	return cs.config.Router.EnableTrailingSlash
}

// GetRouterEnableSlashRedirect returns whether slash redirection is enabled
func (cs *configService) GetRouterEnableSlashRedirect() bool {
	return cs.config.Router.EnableSlashRedirect
}

// GetRouterEnableMethodNotAllowed returns whether method not allowed handler is enabled
func (cs *configService) GetRouterEnableMethodNotAllowed() bool {
	return cs.config.Router.EnableMethodNotAllowed
}
