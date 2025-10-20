package interfaces

import "time"

// CENTRAL TYPE DEFINITIONS - Consolidation of all duplicate structs
// This file eliminates the massive struct redundancy identified in code quality analysis

// Route represents a mapping between a URL path and a template file
// CONSOLIDATES: router/models.go:3, interfaces/services.go:54, router/middleware/template_middleware.go:28
type Route struct {
	// Core routing information
	Path         string `json:"path"`
	Handler      string `json:"handler,omitempty"`
	TemplateFile string `json:"template_file"`
	MetadataFile string `json:"metadata_file,omitempty"`

	// Dynamic routing
	IsDynamic  bool `json:"is_dynamic"`
	Precedence int  `json:"precedence,omitempty"`

	// Internationalization
	Locale string `json:"locale,omitempty"`

	// Security
	AuthSettings *AuthSettings `json:"auth_settings,omitempty"`
}

// LayoutTemplate represents a layout template
// CONSOLIDATES: router/models.go:30, interfaces/services.go:65, router/middleware/template_middleware.go:35
type LayoutTemplate struct {
	FilePath      string `json:"file_path"`
	YamlPath      string `json:"yaml_path,omitempty"`
	ComponentName string `json:"component_name,omitempty"`
	Content       string `json:"content,omitempty"`
	LayoutLevel   int    `json:"layout_level,omitempty"`
}

// ErrorTemplate represents an error template
// CONSOLIDATES: Multiple error template definitions
type ErrorTemplate struct {
	FilePath      string `json:"file_path"`
	ComponentName string `json:"component_name,omitempty"`
	Content       string `json:"content,omitempty"`
	ErrorCode     int    `json:"error_code,omitempty"`
}

// AuthSettings contains authentication configuration
// CONSOLIDATES: interfaces/services.go:81, router/auth_types.go:66, router/middleware/auth_middleware.go:27
type AuthSettings struct {
	Type        AuthType `json:"type"`
	RedirectURL string   `json:"redirect_url,omitempty"`
	Roles       []string `json:"roles,omitempty"`
}

// AuthType represents different authentication types
type AuthType int

const (
	AuthTypePublic AuthType = iota
	AuthTypeUser
	AuthTypeAdmin
)

// String returns the string representation of AuthType
func (at AuthType) String() string {
	switch at {
	case AuthTypePublic:
		return "public"
	case AuthTypeUser:
		return "user"
	case AuthTypeAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// AuthResult contains authentication result (generic)
type AuthResult struct {
	IsAuthenticated bool       `json:"is_authenticated"`
	User            UserEntity `json:"user,omitempty"`
	RedirectURL     string     `json:"redirect_url,omitempty"`
	ErrorMessage    string     `json:"error_message,omitempty"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Valid     bool      `json:"valid"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Template represents a *.templ file containing UI components
// CONSOLIDATES: router/models.go:30 and related template definitions
type Template struct {
	// File information
	FilePath      string `json:"file_path"`
	FileName      string `json:"file_name"`
	DirectoryPath string `json:"directory_path"`

	// Template metadata
	Type          string                 `json:"type,omitempty"`
	ComponentName string                 `json:"component_name,omitempty"`
	Content       string                 `json:"content,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
}

// ConfigFile represents a YAML file containing metadata and settings
// CONSOLIDATES: router/models.go:54 and related config definitions
type ConfigFile struct {
	// File paths
	FilePath         string `json:"file_path"`
	TemplateFilePath string `json:"template_file_path,omitempty"`

	// Metadata
	RouteMetadata   interface{}                  `json:"route_metadata,omitempty"`
	I18nMappings    map[string]string            `json:"i18n_mappings,omitempty"`
	MultiLocaleI18n map[string]map[string]string `json:"multi_locale_i18n,omitempty"`

	// Settings
	AuthSettings    *AuthSettings    `json:"auth_settings,omitempty"`
	LayoutSettings  interface{}      `json:"layout_settings,omitempty"`
	ErrorSettings   interface{}      `json:"error_settings,omitempty"`
	DynamicSettings *DynamicSettings `json:"dynamic_settings,omitempty"`
}

// DynamicSettings contains configuration for dynamic route parameters
type DynamicSettings struct {
	Parameters map[string]*DynamicParameterConfig `json:"parameters,omitempty"`
}

// DynamicParameterConfig contains validation configuration for a single parameter
type DynamicParameterConfig struct {
	Validation      string   `json:"validation,omitempty"`
	Description     string   `json:"description,omitempty"`
	SupportedValues []string `json:"supported_values,omitempty"`
}

// InternationalizationIdentifier represents a structured key for translations
type InternationalizationIdentifier struct {
	// Key is the identifier key (e.g., "admin.dashboard.create.title")
	Key string

	// Source is the source of the key ("opinionated-schema" or "yaml-metadata")
	Source string

	// TemplatePath is the path to the template that uses this identifier
	TemplatePath string

	// DefaultValue is the default value if translation is missing
	DefaultValue string

	// Locales contains translations for different locales
	Locales map[string]string
}
