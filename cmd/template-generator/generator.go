package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
)

//go:embed templates/*.tmpl
var generatorTemplates embed.FS

// generateRegistry generates the modern interface-based template registry
func generateRegistry(config Config, templates []TemplateInfo) error {
	// Generate interface-based registry only
	if err := generateInterfaceRegistry(config, templates); err != nil {
		return fmt.Errorf("failed to generate interface registry: %w", err)
	}
	
	return nil
}


// processImports processes templates and creates unique imports with aliases
func processImports(templates []TemplateInfo) ([]TemplateWithAlias, []ImportInfo) {
	// Group templates by import path
	importGroups := make(map[string][]TemplateInfo)
	for _, tmpl := range templates {
		importGroups[tmpl.ImportPath] = append(importGroups[tmpl.ImportPath], tmpl)
	}

	// Create unique imports with aliases
	var uniqueImports []ImportInfo
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

		uniqueImports = append(uniqueImports, ImportInfo{
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
	var templatesWithAliases []TemplateWithAlias
	for _, tmpl := range templates {
		// Find the correct alias for this template
		for _, group := range importGroups {
			for _, groupTmpl := range group {
				if groupTmpl.TemplateKey == tmpl.TemplateKey {
					templatesWithAliases = append(templatesWithAliases, TemplateWithAlias{
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
func generateInterfaceRegistry(config Config, templates []TemplateInfo) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Group templates by import path and create aliases
	templatesWithAliases, uniqueImports := processImports(templates)

	// Prepare template data
	data := struct {
		PackageName string
		ModuleName  string
		Templates   []TemplateWithAlias
		Imports     []ImportInfo
	}{
		PackageName: config.PackageName,
		ModuleName:  config.ModuleName,
		Templates:   templatesWithAliases,
		Imports:     uniqueImports,
	}

	// Parse embedded template
	tmpl, err := template.ParseFS(generatorTemplates, "templates/interface_registry.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse interface registry template: %w", err)
	}

	// Create output file
	outputPath := filepath.Join(config.OutputDir, "registry.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create interface registry file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute interface registry template: %w", err)
	}

	return nil
}

