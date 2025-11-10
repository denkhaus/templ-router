package config

// Config holds all application configuration
type configImpl struct {
	// Authentication configuration
	Auth AuthConfig `envconfig:"AUTH"`

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

	// Config Configuration
	Config ConfigConfig `envconfig:"CONFIG"`

	// Router configuration
	Router RouterConfig `envconfig:"ROUTER"`
}

// RouterConfig holds router-related configuration
type RouterConfig struct {
	// Enable automatic trailing slash redirection
	// When true, /path/ redirects to /path and vice versa
	EnableTrailingSlash bool `envconfig:"ENABLE_TRAILING_SLASH" default:"true"`

	// Enable automatic slash redirection
	// When true, /path// redirects to /path/
	EnableSlashRedirect bool `envconfig:"ENABLE_SLASH_REDIRECT" default:"true"`

	// Enable method not allowed handler
	EnableMethodNotAllowed bool `envconfig:"ENABLE_METHOD_NOT_ALLOWED" default:"true"`

	// Authentication routes configuration
	EnableAuthRoutes bool   `envconfig:"ENABLE_AUTH_ROUTES" default:"true"`
	AuthRoutePrefix  string `envconfig:"AUTH_ROUTE_PREFIX" default:"/api"`
}

type ConfigConfig struct {
	PrintSummary bool `envconfig:"PRINT_SUMMARY" default:"false"`
}

// AuthConfig holds authentication-related configuration for router
// Only redirect routes remain for router-level authentication failure handling
type AuthConfig struct {
	// Auth routes - kept for router-level authentication failure handling
	SignInRoute string `envconfig:"SIGNIN_ROUTE" default:"/login"`

	// Auth redirect routes (only for success cases) - kept for router-level handling
	SignInSuccessRoute  string `envconfig:"SIGNIN_SUCCESS_ROUTE" default:"/"`
	SignUpSuccessRoute  string `envconfig:"SIGNUP_SUCCESS_ROUTE" default:"/"`
	SignOutSuccessRoute string `envconfig:"SIGNOUT_SUCCESS_ROUTE" default:"/"`
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
