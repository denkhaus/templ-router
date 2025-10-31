package generate

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/denkhaus/templ-router/cmd/trgen/types"
	"github.com/denkhaus/templ-router/cmd/trgen/version"
)

//go:embed templates/*.tmpl
var generatorTemplates embed.FS

// generateRegistry generates the modern interface-based template registry
func GenerateRegistry(config types.Config, templates []types.TemplateInfo) error {
	// Generate interface-based registry only
	if err := generateInterfaceRegistry(config, templates); err != nil {
		return fmt.Errorf("failed to generate interface registry: %w", err)
	}

	return nil
}

// processImports processes templates and creates unique imports with aliases
func processImports(templates []types.TemplateInfo) ([]types.TemplateWithAlias, []types.ImportInfo) {
	// Group templates by import path
	importGroups := make(map[string][]types.TemplateInfo)
	for _, tmpl := range templates {
		importGroups[tmpl.ImportPath] = append(importGroups[tmpl.ImportPath], tmpl)
	}

	// Create unique imports with aliases
	var uniqueImports []types.ImportInfo
	aliasCounter := make(map[string]int)

	for importPath, groupTemplates := range importGroups {
		packageName := groupTemplates[0].PackageName
		alias := packageName

		// Handle alias conflicts
		if count, exists := aliasCounter[alias]; exists {
			aliasCounter[alias] = count + 1
			alias = fmt.Sprintf("%s%d", packageName, count+1)
		} else {
			aliasCounter[alias] = 1
		}

		uniqueImports = append(uniqueImports, types.ImportInfo{
			Alias: alias,
			Path:  importPath,
		})

		// Update all templates in this group with the final alias
		for i := range groupTemplates {
			importGroups[importPath][i].PackageAlias = alias
		}
	}

	// Sort imports for consistent output
	sort.Slice(uniqueImports, func(i, j int) bool {
		return uniqueImports[i].Path < uniqueImports[j].Path
	})

	// Flatten templates with updated aliases
	var templatesWithAliases []types.TemplateWithAlias
	for _, tmpl := range templates {
		// Find the correct alias for this template
		for _, group := range importGroups {
			for _, groupTmpl := range group {
				if groupTmpl.TemplateKey == tmpl.TemplateKey {
					templatesWithAliases = append(templatesWithAliases, types.TemplateWithAlias{
						TemplateInfo: tmpl,
						PackageAlias: groupTmpl.PackageAlias,
					})
					break
				}
			}
		}
	}

	return templatesWithAliases, uniqueImports
}

// generateInterfaceRegistry generates the interface-based template registry
func generateInterfaceRegistry(config types.Config, templates []types.TemplateInfo) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Clean up existing registry file before regeneration
	outputPath := filepath.Join(config.OutputDir, "registry.go")
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return fmt.Errorf("failed to remove existing registry file: %w", err)
		}
	}

	// Group templates by import path and create aliases
	templatesWithAliases, uniqueImports := processImports(templates)

	// Prepare template data
	buildInfo := version.GetBuildInfo()
	data := struct {
		PackageName      string
		ModuleName       string
		Templates        []types.TemplateWithAlias
		Imports          []types.ImportInfo
		GeneratorVersion string
		GeneratedAt      string
	}{
		PackageName:      config.PackageName,
		ModuleName:       config.ModuleName,
		Templates:        templatesWithAliases,
		Imports:          uniqueImports,
		GeneratorVersion: buildInfo.Short(),
		GeneratedAt:      time.Now().Format("2006-01-02 15:04:05 MST"),
	}

	// Parse embedded template
	tmpl, err := template.ParseFS(generatorTemplates, "templates/interface_registry.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse interface registry template: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create interface registry file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", cerr)
		}
	}()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute interface registry template: %w", err)
	}

	return nil
}
