package main

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// toSnakeCase converts CamelCase to snake_case for URL-friendly route names
func toSnakeCase(str string) string {
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
func createRoutePattern(filePath, functionName string, config Config) string {
	dir := filepath.Dir(filePath)
	rootDir := config.ScanPath // Use configurable root directory

	// Extract relative path from configurable root directory
	if strings.Contains(dir, "/"+rootDir+"/") {
		parts := strings.Split(dir, "/"+rootDir+"/")
		if len(parts) > 1 {
			dir = parts[1]
		} else {
			dir = ""
		}
	} else if strings.HasSuffix(dir, "/"+rootDir) {
		dir = ""
	}

	if dir == "" || dir == "." || dir == rootDir {
		// Root level templates - use generic pattern based on function name
		if functionName == "Page" {
			return "/"
		} else if functionName == "Layout" {
			return "/layout"
		} else if functionName == "Error" {
			return "/error"
		} else {
			// Other templates get their own routes based on function name
			return "/" + toSnakeCase(functionName)
		}
	}

	// Convert identifier_ back to $identifier for URL patterns
	// and apply snake_case conversion for better URL conventions
	parts := strings.Split(dir, "/")
	var cleanParts []string
	for _, part := range parts {
		if part == "" || part == rootDir {
			continue // Skip empty parts and root directory prefix
		}
		if strings.HasSuffix(part, "_") {
			paramName := strings.TrimSuffix(part, "_")
			cleanParts = append(cleanParts, "$"+paramName)
		} else {
			// Convert CamelCase/PascalCase to snake_case for URL-friendly routes
			cleanParts = append(cleanParts, toSnakeCase(part))
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
		return routePath + "/" + toSnakeCase(functionName)
	}
}

// getLocalPackageInfo handles local packages generically
func getLocalPackageInfo(filePath, moduleName string, config Config) (string, string) {
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

	// Generic path handling - find the scan path and extract relative path from there
	// Split the path by "/" and find the last occurrence of rootDir
	pathParts := strings.Split(dir, "/")
	var foundIndex = -1
	
	// Find the last occurrence of the rootDir in the path
	for i := len(pathParts) - 1; i >= 0; i-- {
		if pathParts[i] == rootDir {
			foundIndex = i
			break
		}
	}
	
	if foundIndex != -1 {
		// Extract everything after the last occurrence of rootDir
		if foundIndex == len(pathParts)-1 {
			// We're in the root scan directory
			dir = rootDir
		} else {
			// We're in a subdirectory - join the remaining parts
			subParts := pathParts[foundIndex+1:]
			dir = rootDir + "/" + strings.Join(subParts, "/")
		}
	} else {
		// Fallback - assume we're in root
		dir = rootDir
	}

	// Create import path with module name prefix
	importPath := moduleName + "/" + dir
	
	return packageName, importPath
}

// getPackageInfo extracts package name and import path from a Go file
func getPackageInfo(filePath, moduleName string, config Config) (string, string) {
	// Use generic local package info for all modules
	return getLocalPackageInfo(filePath, moduleName, config)
}

// createHumanName creates a human-readable name for documentation
func createHumanName(filePath, functionName string) string {
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
func createPackageAlias(packageName, importPath string, config Config) string {
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