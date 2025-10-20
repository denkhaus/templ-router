package router

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
