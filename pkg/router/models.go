package router

// Route represents a mapping between a URL path and a template file
type Route struct {
	// Path is the URL path for the route (e.g., "/admin/dashboard", "/user/$id")
	Path string
	
	// Handler is the reference to the templ handler function
	Handler string
	
	// TemplateFile is the path to the corresponding *.templ file
	TemplateFile string
	
	// MetadataFile is the optional path to the corresponding *.yaml metadata file
	MetadataFile string
	
	// AuthSettings contains authentication and authorization settings for this route
	AuthSettings *AuthSettings
	
	// Locale contains locale information for internationalized routes
	Locale string
	
	// IsDynamic indicates whether this is a dynamic route (contains parameters like $id)
	IsDynamic bool
	
	// Precedence is the priority level (manual routes have higher precedence than file-based)
	Precedence int
}

// Template represents a *.templ file containing UI components
type Template struct {
	// FilePath is the full path to the template file in the app directory
	FilePath string
	
	// FileName is the name of the file without path (e.g., "create.templ")
	FileName string
	
	// DirectoryPath is the path to the directory containing this template
	DirectoryPath string
	
	// Type is the type of template (e.g., "layout", "page", "error")
	Type string
	
	// ComponentName is the name of the generated Go component
	ComponentName string
	
	// Content is the content of the template file (for validation)
	Content string
	
	// TemplateParams are parameters expected by the template
	TemplateParams map[string]interface{}
}

// ConfigFile represents a YAML file containing metadata and settings
type ConfigFile struct {
	// FilePath is the full path to the YAML config file
	FilePath string
	
	// TemplateFilePath is the path to the corresponding *.templ file
	TemplateFilePath string
	
	// RouteMetadata contains custom route configuration (HTTP route, auth settings)
	RouteMetadata interface{}
	
	// I18nMappings contains custom i18n identifier mappings
	I18nMappings map[string]string
	
	// MultiLocaleI18n contains multi-locale translations (locale -> key -> translation)
	MultiLocaleI18n map[string]map[string]string
	
	// AuthSettings contains authentication settings that override parent settings
	AuthSettings *AuthSettings
	
	// LayoutSettings contains layout configuration for this template (deprecated - layout inheritance is automatic)
	LayoutSettings interface{}
	
	// ErrorSettings contains error handling configuration
	ErrorSettings interface{}
	
	// DynamicSettings contains dynamic parameter validation configuration
	DynamicSettings *DynamicSettings
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

// DynamicSettings contains configuration for dynamic route parameters
type DynamicSettings struct {
	// Parameters contains validation rules for each parameter
	Parameters map[string]*DynamicParameterConfig
}

// DynamicParameterConfig contains validation configuration for a single parameter
type DynamicParameterConfig struct {
	// Validation is the regex pattern for validating parameter values
	Validation string
	
	// Description is a human-readable description of the parameter
	Description string
	
	// SupportedValues is a list of explicitly allowed values (optional)
	SupportedValues []string
}