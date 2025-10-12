package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/tools/go/packages"
)

// scanTemplatesWithPackages uses the packages API to scan for templates
func scanTemplatesWithPackages(config Config) ([]TemplateInfo, []string, error) {
	var templates []TemplateInfo
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
			if err := validateTemplatePath(filePath, config); err != nil {
				return nil, nil, fmt.Errorf("template validation failed: %w", err)
			}

			fmt.Printf("Processing file: %s\n", filePath)

			fileTemplates, fileValidationErrors, err := extractTemplatesFromFile(file, filePath, pkg, config)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to extract templates from %s: %w", filePath, err)
			}

			templates = append(templates, fileTemplates...)
			allValidationErrors = append(allValidationErrors, fileValidationErrors...)
		}
	}

	return templates, allValidationErrors, nil
}

// extractTemplatesFromFile extracts template functions from a single file using packages API
func extractTemplatesFromFile(file *ast.File, filePath string, pkg *packages.Package, config Config) ([]TemplateInfo, []string, error) {
	var templates []TemplateInfo
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
		if err := validateFunctionNaming(functionName, filePath); err != nil {
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
		fnType, ok := obj.Type().(*types.Signature)
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
		packageName, importPath := getPackageInfo(filePath, config.ModuleName, config)
		
		// Debug: Log the paths being generated
		fmt.Printf("      -> File: %s\n", filePath)
		fmt.Printf("      -> Package: %s, Import: %s\n", packageName, importPath)
		templateKey := uuid.New().String()
		humanName := createHumanName(filePath, functionName)

		// Create route pattern
		routePattern := createRoutePattern(filePath, functionName, config)

		// Create package alias
		packageAlias := createPackageAlias(packageName, importPath, config)

		templateInfo := TemplateInfo{
			FilePath:      filePath,
			FunctionName:  functionName,
			PackageName:   packageName,
			PackageAlias:  packageAlias,
			ImportPath:    importPath,
			TemplateKey:   templateKey,
			RoutePattern:  routePattern,
			HumanName:     humanName,
		}

		templates = append(templates, templateInfo)
		fmt.Printf("      -> %s -> %s\n", functionName, templateKey)
	}

	return templates, validationErrors, nil
}

// validateTemplatePath validates that the template file is in the correct location
func validateTemplatePath(filePath string, config Config) error {
	rootDir := config.ScanPath
	// Basic validation - ensure it's a _templ.go file in the configured root directory
	if !strings.Contains(filePath, "/"+rootDir+"/") && !strings.HasSuffix(filepath.Dir(filePath), "/"+rootDir) {
		return fmt.Errorf("template file %s is not in the %s directory", filePath, rootDir)
	}
	return nil
}

// scanSpecificPackage scans templates for a specific package
func scanSpecificPackage(config Config) ([]TemplateInfo, error) {
	fmt.Printf("Scanning specific package: %s\n", config.PackageName)
	
	// Use the existing scanTemplates logic but filter for specific package
	templates, err := scanAllTemplates(config)
	if err != nil {
		return nil, err
	}
	
	// Filter templates by target package
	var filteredTemplates []TemplateInfo
	for _, template := range templates {
		if template.PackageName == config.PackageName {
			filteredTemplates = append(filteredTemplates, template)
		}
	}
	
	fmt.Printf("Found %d templates in package %s\n", len(filteredTemplates), config.PackageName)
	return filteredTemplates, nil
}

// scanAllTemplates is the original scanTemplates function renamed
func scanAllTemplates(config Config) ([]TemplateInfo, error) {


// scanTemplates is the main entry point for template scanning
	// If TargetPackage is specified, scan only that package
	if config.PackageName != "" {
		return scanSpecificPackage(config)
	}
	
	// Otherwise, scan all templates
	return scanAllTemplates(config)
}
