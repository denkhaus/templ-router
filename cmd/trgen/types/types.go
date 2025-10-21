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
	FilePath     string
	FunctionName string
	TemplateKey  string
	RoutePattern string
	PackageName  string
	PackageAlias string
	ImportPath   string
	HumanName    string // Human-readable name for documentation
	
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
