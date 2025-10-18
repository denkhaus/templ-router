package config

import "time"

// Config holds all application configuration
type configImpl struct {
	// Server configuration
	Server ServerConfig `envconfig:"SERVER"`

	// Database configuration
	Database DatabaseConfig `envconfig:"DATABASE"`

	// Authentication configuration
	Auth AuthConfig `envconfig:"AUTH"`

	// Email configuration
	Email EmailConfig `envconfig:"EMAIL"`

	// Security configuration
	Security SecurityConfig `envconfig:"SECURITY"`

	// Logging configuration
	Logging LoggingConfig `envconfig:"LOGGING"`

	// Internationalization configuration
	I18n I18nConfig `envconfig:"I18N"`

	// Layout configuration
	Layout LayoutConfig `envconfig:"LAYOUT"`

	// Template generator configuration
	TemplateGenerator TemplateGeneratorConfig `envconfig:"TEMPLATE_GENERATOR"`

	// Environment configuration
	Environment EnvironmentConfig `envconfig:"ENVIRONMENT"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host            string        `envconfig:"HOST" default:"localhost"`
	Port            int           `envconfig:"PORT" default:"8080"`
	BaseURL         string        `envconfig:"BASE_URL" default:"http://localhost:8080"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"30s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"30s"`
	IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT" default:"120s"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"postgres"`
	Name     string `envconfig:"NAME" default:"router_db"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
}

// AuthConfig holds authentication-related configuration
type AuthConfig struct {
	// Email verification settings
	RequireEmailVerification bool          `envconfig:"REQUIRE_EMAIL_VERIFICATION" default:"true"`
	VerificationTokenExpiry  time.Duration `envconfig:"VERIFICATION_TOKEN_EXPIRY" default:"24h"`

	// Session settings
	SessionCookieName string        `envconfig:"SESSION_COOKIE_NAME" default:"session_id"`
	SessionExpiry     time.Duration `envconfig:"SESSION_EXPIRY" default:"24h"`
	SessionSecure     bool          `envconfig:"SESSION_SECURE" default:"false"`
	SessionHttpOnly   bool          `envconfig:"SESSION_HTTP_ONLY" default:"true"`
	SessionSameSite   string        `envconfig:"SESSION_SAME_SITE" default:"lax"`

	// Password settings
	MinPasswordLength   int  `envconfig:"MIN_PASSWORD_LENGTH" default:"8"`
	RequireStrongPasswd bool `envconfig:"REQUIRE_STRONG_PASSWORD" default:"false"`

	// Default admin user settings
	CreateDefaultAdmin    bool   `envconfig:"CREATE_DEFAULT_ADMIN" default:"true"`
	DefaultAdminEmail     string `envconfig:"DEFAULT_ADMIN_EMAIL" default:"admin@example.com"`
	DefaultAdminPassword  string `envconfig:"DEFAULT_ADMIN_PASSWORD" default:"admin123"`
	DefaultAdminFirstName string `envconfig:"DEFAULT_ADMIN_FIRST_NAME" default:"Default"`
	DefaultAdminLastName  string `envconfig:"DEFAULT_ADMIN_LAST_NAME" default:"Admin"`
}

// EmailConfig holds email-related configuration
type EmailConfig struct {
	// SMTP settings
	SMTPHost     string `envconfig:"SMTP_HOST" default:""`
	SMTPPort     int    `envconfig:"SMTP_PORT" default:"587"`
	SMTPUsername string `envconfig:"SMTP_USERNAME" default:""`
	SMTPPassword string `envconfig:"SMTP_PASSWORD" default:""`
	SMTPUseTLS   bool   `envconfig:"SMTP_USE_TLS" default:"true"`

	// Email settings
	FromEmail    string `envconfig:"FROM_EMAIL" default:"noreply@example.com"`
	FromName     string `envconfig:"FROM_NAME" default:"Router Application"`
	ReplyToEmail string `envconfig:"REPLY_TO_EMAIL" default:""`

	// Development settings
	EnableDummyMode bool `envconfig:"ENABLE_DUMMY_MODE" default:"true"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	// CSRF protection
	CSRFSecret   string `envconfig:"CSRF_SECRET" default:"change-me-in-production"`
	CSRFSecure   bool   `envconfig:"CSRF_SECURE" default:"false"`
	CSRFHttpOnly bool   `envconfig:"CSRF_HTTP_ONLY" default:"true"`
	CSRFSameSite string `envconfig:"CSRF_SAME_SITE" default:"strict"`

	// Rate limiting
	EnableRateLimit   bool `envconfig:"ENABLE_RATE_LIMIT" default:"true"`
	RateLimitRequests int  `envconfig:"RATE_LIMIT_REQUESTS" default:"100"`

	// Security headers
	EnableSecurityHeaders bool `envconfig:"ENABLE_SECURITY_HEADERS" default:"true"`
	EnableHSTS            bool `envconfig:"ENABLE_HSTS" default:"false"`
	HSTSMaxAge            int  `envconfig:"HSTS_MAX_AGE" default:"31536000"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level      string `envconfig:"LEVEL" default:"info"`
	Format     string `envconfig:"FORMAT" default:"json"`
	Output     string `envconfig:"OUTPUT" default:"stdout"`
	EnableFile bool   `envconfig:"ENABLE_FILE" default:"false"`
	FilePath   string `envconfig:"FILE_PATH" default:"logs/router.log"`
}

// I18nConfig holds internationalization configuration
type I18nConfig struct {
	// Supported locales (language codes)
	SupportedLocales []string `envconfig:"SUPPORTED_LOCALES" default:"en,de"`
	DefaultLocale    string   `envconfig:"DEFAULT_LOCALE" default:"en"`
	FallbackLocale   string   `envconfig:"FALLBACK_LOCALE" default:"en"`
}

type EnvironmentConfig struct {
	Kind string `envconfig:"KIND" default:"develop"`
}

// LayoutConfig holds layout system configuration
type LayoutConfig struct {
	// Root directory for templates and layouts
	RootDirectory string `envconfig:"ROOT_DIRECTORY" default:"app"`

	// Assets directory for assets
	AssetsDirectory string `envconfig:"ASSETS_DIRECTORY" default:"assets"`

	// Assets route name used to make assets from AssetsDirectory accessible by the assets service
	AssetsRouteName string `envconfig:"ASSETS_ROUTE_NAME" default:"assets"`

	// Layout file name (without extension)
	LayoutFileName string `envconfig:"LAYOUT_FILE_NAME" default:"layout"`

	// Template file extension
	TemplateExtension string `envconfig:"TEMPLATE_EXTENSION" default:".templ"`

	// YAML metadata file extension
	MetadataExtension string `envconfig:"METADATA_EXTENSION" default:".templ.yaml"`

	// Enable layout inheritance (Next.js style)
	EnableInheritance bool `envconfig:"ENABLE_INHERITANCE" default:"true"`
}

// TemplateGeneratorConfig holds template generator configuration
type TemplateGeneratorConfig struct {
	// Output directory for generated templates
	OutputDir string `envconfig:"OUTPUT_DIR" default:"generated/templates"`

	// Package name for generated templates
	PackageName string `envconfig:"PACKAGE_NAME" default:"templates"`
}
