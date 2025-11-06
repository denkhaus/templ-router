package types

// Config holds the configuration for the template generator
type Config struct {
	ModuleName       string // The Go module name (e.g., "demo", "github.com/user/project")
	ScanPath         string // Absolute path to scan for templates
	OutputDir        string // Absolute path where to generate registry
	PackageName      string // Package name for generated code
	TemplateBasePath string // Base path for template imports relative to module root
}

// TemplateInfo holds information about a discovered template
type TemplateInfo struct {
	FilePath     string // Absolute path to the Go file
	TemplatePath string // Relative template path (e.g., "demo/app/locale_/test/page.templ")
	FunctionName string
	TemplateKey  string
	RoutePattern string
	PackageName  string
	PackageAlias string
	ImportPath   string
	HumanName    string // Human-readable name for documentation

	// Template Classification (for generic discovery)
	TemplateType  string // "layout", "page", "error", or "component"
	ComponentName string // Component name extracted from filename (e.g., "footer" from "footer.templ")

	// YAML Metadata Analysis (determined during scanning)
	YAMLExists    bool   // Whether YAML file exists
	YAMLFile      string // Full path to YAML file if it exists
	YAMLHasI18n   bool   // Whether YAML contains i18n data
	YAMLHasMetadata bool  // Whether YAML contains metadata
	YAMLHasAuth   bool   // Whether YAML contains auth settings

	// Data Service Integration
	RequiresDataService  bool   // true if template has data parameter
	DataServiceInterface string // e.g., "dataservices.UserDataService"
	DataParameterType    string // e.g., "*dataservices.UserData"
}

// ImportInfo represents an import statement with alias
type ImportInfo struct {
	Alias string
	Path  string
}

// TemplateWithAlias combines template info with its package alias
type TemplateWithAlias struct {
	TemplateInfo
	PackageAlias string
}
