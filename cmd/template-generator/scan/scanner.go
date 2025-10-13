package scan

import (
	"fmt"
	"go/ast"
	gotypes "go/types"
	"strings"

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
	"github.com/denkhaus/templ-router/cmd/template-generator/utils"
	"github.com/denkhaus/templ-router/cmd/template-generator/validate"

	"github.com/google/uuid"
	"golang.org/x/tools/go/packages"
)

// ScanTemplatesWithPackages uses the packages API to scan for templates
func ScanTemplatesWithPackages(config types.Config) ([]types.TemplateInfo, []string, error) {
	var templates []types.TemplateInfo
	var allValidationErrors []string

	// Load packages
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:  config.ScanPath,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load packages: %w", err)
	}

	// Process each package
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			fmt.Printf("Package %s has errors, skipping\n", pkg.PkgPath)
			continue
		}

		// Process each file in the package
		for _, file := range pkg.Syntax {
			// Get the file path
			filePath := pkg.Fset.Position(file.Pos()).Filename

			// Only process _templ.go files
			if !strings.HasSuffix(filePath, "_templ.go") {
				continue
			}

			// Validate template path
			if err := validate.ValidateTemplatePath(filePath, config); err != nil {
				return nil, nil, fmt.Errorf("template validation failed: %w", err)
			}

			fmt.Printf("Processing file: %s\n", filePath)

			fileTemplates, fileValidationErrors, err := ExtractTemplatesFromFile(file, filePath, pkg, config)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to extract templates from %s: %w", filePath, err)
			}

			templates = append(templates, fileTemplates...)
			allValidationErrors = append(allValidationErrors, fileValidationErrors...)
		}
	}

	return templates, allValidationErrors, nil
}

// ExtractTemplatesFromFile extracts template functions from a single file using packages API
func ExtractTemplatesFromFile(file *ast.File, filePath string, pkg *packages.Package, config types.Config) ([]types.TemplateInfo, []string, error) {
	var templates []types.TemplateInfo
	var validationErrors []string

	// Walk through all declarations in the file
	for _, decl := range file.Decls {
		// Only process function declarations
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil {
			continue
		}

		// Skip methods (functions with receivers)
		if fn.Recv != nil {
			continue
		}

		// Extract function name
		functionName := fn.Name.Name

		fmt.Printf("    Found function: %s\n", functionName)

		// Validate naming convention
		if err := validate.ValidateFunctionNaming(functionName, filePath); err != nil {
			fmt.Printf("      -> NAMING VIOLATION: %s\n", err.Error())
			fmt.Printf("      -> Skipping due to naming convention violation\n")
			validationErrors = append(validationErrors, err.Error())
			continue
		}

		// Get function type information from packages
		obj := pkg.TypesInfo.Defs[fn.Name]
		if obj == nil {
			fmt.Printf("      -> No type info found, skipping\n")
			continue
		}

		// Check if it's a function
		fnType, ok := obj.Type().(*gotypes.Signature)
		if !ok {
			fmt.Printf("      -> Not a function type, skipping\n")
			continue
		}

		// Check if function has parameters
		hasParams := fnType.Params() != nil && fnType.Params().Len() > 0

		// Skip templates with parameters, except for Page, Layout and Error which are needed for routing
		if hasParams && functionName != "Page" && functionName != "Layout" && functionName != "Error" {
			fmt.Printf("      -> %s (has params, skipping)\n", functionName)
			continue
		}

		// Special handling for parametrized templates (Layout, Error)
		if hasParams {
			fmt.Printf("      -> %s (has params, including for routing)\n", functionName)
		}

		// Create template key using UUID
		packageName, importPath := utils.GetPackageInfo(filePath, config.ModuleName, config)

		// Debug: Log the paths being generated
		fmt.Printf("      -> File: %s\n", filePath)
		fmt.Printf("      -> Package: %s, Import: %s\n", packageName, importPath)
		templateKey := uuid.New().String()
		humanName := utils.CreateHumanName(filePath, functionName)

		// Create route pattern
		routePattern := utils.CreateRoutePattern(filePath, functionName, config)

		// Create package alias
		packageAlias := utils.CreatePackageAlias(packageName, importPath, config)

		templateInfo := types.TemplateInfo{
			FilePath:     filePath,
			FunctionName: functionName,
			PackageName:  packageName,
			PackageAlias: packageAlias,
			ImportPath:   importPath,
			TemplateKey:  templateKey,
			RoutePattern: routePattern,
			HumanName:    humanName,
		}

		templates = append(templates, templateInfo)
		fmt.Printf("      -> %s -> %s\n", functionName, templateKey)
	}

	return templates, validationErrors, nil
}

// ScanSpecificPackage scans templates for a specific package
func ScanSpecificPackage(config types.Config) ([]types.TemplateInfo, error) {
	fmt.Printf("Scanning specific package: %s\n", config.PackageName)

	// Use the existing scanTemplates logic but filter for specific package
	templates, err := scanAllTemplates(config)
	if err != nil {
		return nil, err
	}

	// Filter templates by target package
	var filteredTemplates []types.TemplateInfo
	for _, template := range templates {
		if template.PackageName == config.PackageName {
			filteredTemplates = append(filteredTemplates, template)
		}
	}

	fmt.Printf("Found %d templates in package %s\n", len(filteredTemplates), config.PackageName)
	return filteredTemplates, nil
}

// scanAllTemplates is the original scanTemplates function renamed
func scanAllTemplates(config types.Config) ([]types.TemplateInfo, error) {

	// scanTemplates is the main entry point for template scanning
	// If TargetPackage is specified, scan only that package
	if config.PackageName != "" {
		return ScanSpecificPackage(config)
	}

	// Otherwise, scan all templates
	return scanAllTemplates(config)
}
