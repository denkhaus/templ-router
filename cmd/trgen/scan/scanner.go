package scan

import (
	"fmt"
	"go/ast"
	gotypes "go/types"
	"strings"

	"github.com/denkhaus/templ-router/cmd/trgen/types"
	"github.com/denkhaus/templ-router/cmd/trgen/utils"
	"github.com/denkhaus/templ-router/cmd/trgen/validate"

	"github.com/google/uuid"
	"golang.org/x/tools/go/packages"
)

// isDataServiceType checks if a given type implements the DataService interface
func isDataServiceType(paramType gotypes.Type, pkg *packages.Package) (bool, string, string) {
	// Get the underlying type, handling pointers
	actualType := paramType
	if ptr, ok := paramType.(*gotypes.Pointer); ok {
		actualType = ptr.Elem()
	}
	
	// Get the type string for the parameter
	paramTypeStr := paramType.String()
	
	// Try to find a corresponding service interface
	// Look for types that could be data types (ending with "Data" or similar patterns)
	typeStr := actualType.String()
	
	// Extract package and type name
	var packagePath, typeName string
	if named, ok := actualType.(*gotypes.Named); ok {
		obj := named.Obj()
		if obj != nil && obj.Pkg() != nil {
			packagePath = obj.Pkg().Path()
			typeName = obj.Name()
		}
	}
	
	// If we can't extract proper type info, fall back to string analysis
	if packagePath == "" || typeName == "" {
		// Parse from type string like "*github.com/path/dataservices.UserData"
		if strings.Contains(typeStr, "/") {
			parts := strings.Split(typeStr, "/")
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				if strings.Contains(lastPart, ".") {
					dotParts := strings.Split(lastPart, ".")
					if len(dotParts) == 2 {
						packageName := dotParts[0]
						typeName = dotParts[1]
						// Reconstruct a reasonable package path
						packagePath = strings.Join(parts[:len(parts)-1], "/") + "/" + packageName
					}
				}
			}
		}
	}
	
	// Check if this looks like a data type that would have a corresponding service
	if typeName != "" {
		var serviceInterfaceName string
		var servicePackage string
		
		// Common patterns for data service types
		if strings.HasSuffix(typeName, "Data") {
			// Convert "UserData" to "UserDataService"
			serviceInterfaceName = strings.TrimSuffix(typeName, "Data") + "DataService"
		} else if strings.HasSuffix(typeName, "Model") {
			// Convert "UserModel" to "UserDataService"  
			serviceInterfaceName = strings.TrimSuffix(typeName, "Model") + "DataService"
		} else if strings.HasSuffix(typeName, "Entity") {
			// Convert "UserEntity" to "UserDataService"
			serviceInterfaceName = strings.TrimSuffix(typeName, "Entity") + "DataService"
		} else {
			// For other types, assume they might have a corresponding service
			serviceInterfaceName = typeName + "DataService"
		}
		
		// Determine service package - look for common service package patterns
		if strings.Contains(packagePath, "/") {
			pathParts := strings.Split(packagePath, "/")
			for _, part := range pathParts {
				if strings.Contains(part, "dataservices") || strings.Contains(part, "data") {
					servicePackage = part
					break
				}
			}
			// If no specific service package found, use the same package as the data type
			if servicePackage == "" {
				servicePackage = pathParts[len(pathParts)-1]
			}
		}
		
		// Construct the full service interface name
		fullServiceInterface := servicePackage + "." + serviceInterfaceName
		
		// For now, we assume it's a data service if it matches our naming patterns
		// In a more sophisticated implementation, we could try to resolve and check
		// if the service actually implements the DataService interface
		isDataService := strings.HasSuffix(typeName, "Data") || 
						strings.HasSuffix(typeName, "Model") || 
						strings.HasSuffix(typeName, "Entity") ||
						strings.Contains(packagePath, "dataservices") ||
						strings.Contains(packagePath, "data")
		
		return isDataService, fullServiceInterface, paramTypeStr
	}
	
	return false, "", paramTypeStr
}

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

		// Analyze parameters for data service integration
		var requiresDataService bool
		var dataServiceInterface, dataParameterType string

		if hasParams && (functionName == "Page" || functionName == "Layout" || functionName == "Error") {
			// Check if first parameter is a data service type
			if fnType.Params().Len() > 0 {
				firstParam := fnType.Params().At(0)
				paramType := firstParam.Type()
				
				// Use generic function to check if parameter implements DataService interface
				isDataService, serviceInterface, paramTypeStr := isDataServiceType(paramType, pkg)
				
				if isDataService {
					requiresDataService = true
					dataParameterType = paramTypeStr
					dataServiceInterface = serviceInterface
					
					fmt.Printf("      -> %s (has data service: %s with GetData method)\n", functionName, dataServiceInterface)
				} else {
					fmt.Printf("      -> %s (has params, including for routing)\n", functionName)
				}
			}
		} else if hasParams && functionName != "Page" && functionName != "Layout" && functionName != "Error" {
			fmt.Printf("      -> %s (has params, skipping)\n", functionName)
			continue
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
			
			// Data Service Integration
			RequiresDataService:  requiresDataService,
			DataServiceInterface: dataServiceInterface,
			DataParameterType:    dataParameterType,
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
