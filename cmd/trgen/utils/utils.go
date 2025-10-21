// Package utils provides utility functions for the template generator.
//
// CRITICAL REQUIREMENT: All functions in this package must be file/directory agnostic.
// They should not hardcode any project-specific paths, module names, or base routes.
// The generator must work with any project structure and be completely generic.
package utils

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/cmd/trgen/types"
)

// toSnakeCase converts CamelCase to snake_case for URL-friendly route names
func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		if r >= 'A' && r <= 'Z' {
			result.WriteRune(r - 'A' + 'a')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// sanitizePackageName converts directory names to valid Go package identifiers
// Removes hyphens, dots, and other invalid characters
func sanitizePackageName(name string) string {
	// Replace hyphens and dots with empty string
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, ".", "")
	
	// Remove any other non-alphanumeric characters except underscores
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	
	sanitized := result.String()
	
	// Ensure it doesn't start with a number
	if len(sanitized) > 0 && sanitized[0] >= '0' && sanitized[0] <= '9' {
		sanitized = "pkg" + sanitized
	}
	
	// Ensure it's not empty
	if sanitized == "" {
		sanitized = "pkg"
	}
	
	return sanitized
}

// createRoutePattern creates a route pattern from a file path and function name
// IMPORTANT: This function must be file/directory agnostic and not hardcode any base paths.
// Routes should be generated purely based on the file structure without project-specific assumptions.
func CreateRoutePattern(filePath, functionName string, config types.Config) string {
	dir := filepath.Dir(filePath)
	rootDir := config.ScanPath // Use configurable root directory

	// Convert to absolute paths for proper comparison
	absScanPath, err := filepath.Abs(rootDir)
	if err != nil {
		// Fallback to original logic
		dir = filepath.ToSlash(dir)
		dirParts := strings.Split(dir, "/")
		var scanPathIndex = -1
		for i := len(dirParts) - 1; i >= 0; i-- {
			if dirParts[i] == rootDir {
				scanPathIndex = i
				break
			}
		}
		
		if scanPathIndex == -1 {
			// Keep original dir
		} else if scanPathIndex == len(dirParts)-1 {
			dir = ""
		} else {
			relativeParts := dirParts[scanPathIndex+1:]
			dir = strings.Join(relativeParts, "/")
		}
	} else {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			// Fallback to original logic
			dir = filepath.ToSlash(dir)
			dirParts := strings.Split(dir, "/")
			var scanPathIndex = -1
			for i := len(dirParts) - 1; i >= 0; i-- {
				if dirParts[i] == rootDir {
					scanPathIndex = i
					break
				}
			}
			
			if scanPathIndex == -1 {
				// Keep original dir
			} else if scanPathIndex == len(dirParts)-1 {
				dir = ""
			} else {
				relativeParts := dirParts[scanPathIndex+1:]
				dir = strings.Join(relativeParts, "/")
			}
		} else {
			// Calculate relative path from scan path to template directory
			relativeDir, err := filepath.Rel(absScanPath, absDir)
			if err != nil || strings.HasPrefix(relativeDir, "..") {
				// Fallback to original logic
				dir = filepath.ToSlash(dir)
				dirParts := strings.Split(dir, "/")
				var scanPathIndex = -1
				for i := len(dirParts) - 1; i >= 0; i-- {
					if dirParts[i] == rootDir {
						scanPathIndex = i
						break
					}
				}
				
				if scanPathIndex == -1 {
					// Keep original dir
				} else if scanPathIndex == len(dirParts)-1 {
					dir = ""
				} else {
					relativeParts := dirParts[scanPathIndex+1:]
					dir = strings.Join(relativeParts, "/")
				}
			} else {
				// Use the calculated relative path
				dir = filepath.ToSlash(relativeDir)
				if dir == "." {
					dir = ""
				}
			}
		}
	}

	if dir == "" || dir == "." || dir == rootDir {
		// Root level templates
		if functionName == "Page" {
			return "/"
		} else if functionName == "Layout" {
			return "/layout"
		} else if functionName == "Error" {
			return "/error"
		} else {
			// Other templates get their own routes based on function name
			return "/" + ToSnakeCase(functionName)
		}
	}

	// Convert identifier_ back to {identifier} for URL patterns
	// and apply snake_case conversion for better URL conventions
	parts := strings.Split(dir, "/")
	var cleanParts []string
	for _, part := range parts {
		if part == "" || part == rootDir {
			continue // Skip empty parts and root directory prefix
		}
		if strings.HasSuffix(part, "_") {
			paramName := strings.TrimSuffix(part, "_")
			cleanParts = append(cleanParts, "{"+paramName+"}")
		} else {
			// Convert CamelCase/PascalCase to snake_case for URL-friendly routes
			cleanParts = append(cleanParts, ToSnakeCase(part))
		}
	}

	var routePath string
	if len(cleanParts) == 0 {
		routePath = ""
	} else {
		routePath = "/" + strings.Join(cleanParts, "/")
	}

	// Add function-specific suffix for different template types
	if functionName == "Page" {
		return routePath
	} else if functionName == "Layout" {
		return routePath + "/layout"
	} else if functionName == "Error" {
		return routePath + "/error"
	} else {
		// For other templates (components), use the function name
		return routePath + "/" + ToSnakeCase(functionName)
	}
}

// getLocalPackageInfo handles local packages generically
func GetLocalPackageInfo(filePath, moduleName string, config types.Config) (string, string) {
	rootDir := config.ScanPath

	// Parse the file to get package declaration
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return rootDir, moduleName + "/" + rootDir
	}

	packageName := node.Name.Name
	dir := filepath.Dir(filePath)
	dir = filepath.ToSlash(dir)

	// Find the scan path in the directory path
	pathParts := strings.Split(dir, "/")
	var scanPathIndex = -1
	for i := len(pathParts) - 1; i >= 0; i-- {
		if pathParts[i] == rootDir {
			scanPathIndex = i
			break
		}
	}
	
	var importPath string
	
	if scanPathIndex == -1 {
		// Scan path not found in directory - use base import path
		importPath = moduleName + "/" + rootDir
	} else if scanPathIndex == len(pathParts)-1 {
		// File is directly in the scan path directory
		importPath = moduleName + "/" + rootDir
	} else {
		// File is in a subdirectory of the scan path
		subParts := pathParts[scanPathIndex+1:]
		importPath = moduleName + "/" + rootDir + "/" + strings.Join(subParts, "/")
		
		// Update package name to be the last directory in the path
		rawPackageName := pathParts[len(pathParts)-1]
		packageName = sanitizePackageName(rawPackageName)
	}
	
	// Clean up any "./" in the path
	importPath = strings.ReplaceAll(importPath, "/./", "/")

	return packageName, importPath
}

// getPackageInfo extracts package name and import path from a Go file
func GetPackageInfo(filePath, moduleName string, config types.Config) (string, string) {
	// Use generic local package info for all modules
	return GetLocalPackageInfo(filePath, moduleName, config)
}

// createHumanName creates a human-readable name for documentation
func CreateHumanName(filePath, functionName string) string {
	dir := filepath.Dir(filePath)

	// Extract relative path from app/
	if strings.Contains(dir, "/app/") {
		parts := strings.Split(dir, "/app/")
		if len(parts) > 1 {
			dir = parts[1]
		} else {
			dir = ""
		}
	} else if strings.HasSuffix(dir, "/app") {
		dir = ""
	}

	if dir == "" || dir == "." {
		// Root level: just function name
		return functionName
	}

	// Remove dynamic parameters (identifier_) for human name
	parts := strings.Split(dir, "/")
	var cleanParts []string
	for _, part := range parts {
		if !strings.HasSuffix(part, "_") {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) == 0 {
		// Only dynamic parts - create generic name
		return functionName
	}

	// Create human name like "dashboard.Page" or "admin.Layout"
	packageName := cleanParts[len(cleanParts)-1]
	return packageName + "." + functionName
}

// createPackageAlias creates a package alias for imports
func CreatePackageAlias(packageName, importPath string, config types.Config) string {
	rootDir := config.ScanPath // Use configurable root directory
	// For root directory package, no alias needed
	if packageName == rootDir {
		return packageName
	}

	// For other packages, create alias to avoid conflicts
	parts := strings.Split(importPath, "/")
	if len(parts) > 1 {
		// Use the last part of the import path as alias and sanitize it
		rawAlias := parts[len(parts)-1]
		return sanitizePackageName(rawAlias)
	}

	return sanitizePackageName(packageName)
}
