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

	"github.com/denkhaus/templ-router/cmd/template-generator/types"
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

// createRoutePattern creates a route pattern from a file path and function name
// IMPORTANT: This function must be file/directory agnostic and not hardcode any base paths.
// Routes should be generated purely based on the file structure without project-specific assumptions.
func CreateRoutePattern(filePath, functionName string, config types.Config) string {
	dir := filepath.Dir(filePath)
	rootDir := config.ScanPath // Use configurable root directory

	// Extract relative path from scan directory

	// Extract relative path from configurable root directory
	// GENERIC APPROACH: Work with the actual working directory and scan path
	// to determine the relative path regardless of project structure
	
	// Convert to forward slashes for consistent handling
	dir = filepath.ToSlash(dir)
	
	// The key insight: we need to find where the scan path ends in the file path
	// and extract everything after that point, regardless of the absolute path structure
	
	// Split the directory path into parts
	dirParts := strings.Split(dir, "/")
	
	// Find the rightmost occurrence of the scan path in the directory parts
	// This handles cases like: /any/path/structure/scanPath/sub/dirs
	var scanPathIndex = -1
	for i := len(dirParts) - 1; i >= 0; i-- {
		if dirParts[i] == rootDir {
			scanPathIndex = i
			break
		}
	}
	
	
	if scanPathIndex == -1 {
		// Scan path not found in directory - this shouldn't happen in normal operation
		// Keep original dir for now, but this might indicate a configuration issue
	} else if scanPathIndex == len(dirParts)-1 {
		// The scan path is the last part - we're in the root scan directory
		dir = ""
	} else {
		// Extract everything after the scan path
		relativeParts := dirParts[scanPathIndex+1:]
		dir = strings.Join(relativeParts, "/")
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

	// CRITICAL: We need to extract the full relative path from the module root,
	// not just from the scan path. The import path must be relative to where
	// the go.mod file is located.
	
	// WORKING DIRECTORY AWARE APPROACH: 
	// The generator is executed from a working directory, and scan-path is relative to that.
	// We need to determine the working directory's position relative to the module root.
	
	pathParts := strings.Split(dir, "/")
	
	// Find the scan path in the directory
	var scanPathIndex = -1
	for i := len(pathParts) - 1; i >= 0; i-- {
		if pathParts[i] == rootDir {
			scanPathIndex = i
			break
		}
	}
	
	var importPath string
	if scanPathIndex == -1 {
		// Fallback if scan path not found
		importPath = moduleName + "/" + rootDir
	} else {
		// Extract the working directory from the file path
		// For /path/to/project/demo/app/file.go, if scanPath is "app",
		// then working directory is "demo" and full path should be "demo/app"
		
		var workingDirPath string
		if scanPathIndex > 0 {
			// Look at the directory immediately before the scan path
			// This represents the working directory where the generator was executed
			workingDir := pathParts[scanPathIndex-1]
			
			// Only include the working directory if it's not an absolute path prefix
			// (skip things like "", "app", "usr", etc. that are clearly not project dirs)
			if workingDir != "" && workingDir != "app" && workingDir != "usr" && 
			   workingDir != "home" && workingDir != "tmp" && len(workingDir) > 1 {
				workingDirPath = workingDir + "/"
			}
		}
		
		if scanPathIndex == len(pathParts)-1 {
			// We're in the root scan directory
			importPath = moduleName + "/" + workingDirPath + rootDir
		} else {
			// We're in a subdirectory of the scan path
			subParts := pathParts[scanPathIndex+1:]
			importPath = moduleName + "/" + workingDirPath + rootDir + "/" + strings.Join(subParts, "/")
		}
	}

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
		// Use the last part of the import path as alias
		return parts[len(parts)-1]
	}

	return packageName
}
