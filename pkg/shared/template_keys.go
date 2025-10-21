package shared

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"
)

// TemplateKeyGenerator provides centralized template key generation
type TemplateKeyGenerator struct{}

// NewTemplateKeyGenerator creates a new template key generator
func NewTemplateKeyGenerator() *TemplateKeyGenerator {
	return &TemplateKeyGenerator{}
}

// GenerateTemplateKey generates a consistent template key from template file path
// This ensures all services use the same key for the same template
func (tkg *TemplateKeyGenerator) GenerateTemplateKey(templatePath string) string {
	// Normalize the path to ensure consistency
	normalizedPath := filepath.Clean(templatePath)
	normalizedPath = strings.ReplaceAll(normalizedPath, "\\", "/")
	
	// Generate a deterministic hash from the normalized path
	hash := md5.Sum([]byte(normalizedPath))
	return fmt.Sprintf("%x", hash)[:32] // Use first 32 characters for shorter keys
}

// GenerateRouteKey generates a consistent route key from route pattern
func (tkg *TemplateKeyGenerator) GenerateRouteKey(routePattern string) string {
	// Normalize the route pattern
	normalizedRoute := strings.TrimSpace(routePattern)
	if !strings.HasPrefix(normalizedRoute, "/") {
		normalizedRoute = "/" + normalizedRoute
	}
	
	// Generate a deterministic hash
	hash := md5.Sum([]byte(normalizedRoute))
	return fmt.Sprintf("route_%x", hash)[:32]
}

// ExtractRouteFromTemplatePath extracts route pattern from template file path
// e.g., "demo/app/locale_/test/page.templ" -> "/{locale}/test"
func (tkg *TemplateKeyGenerator) ExtractRouteFromTemplatePath(templatePath string) string {
	// Remove file extension
	pathWithoutExt := strings.TrimSuffix(templatePath, ".templ")
	
	// Remove base path and page.templ
	pathWithoutExt = strings.TrimSuffix(pathWithoutExt, "/page")
	
	// Extract the route part
	parts := strings.Split(pathWithoutExt, "/")
	var routeParts []string
	
	// Skip until we find "app" directory
	appFound := false
	for _, part := range parts {
		if part == "app" {
			appFound = true
			continue
		}
		if !appFound {
			continue
		}
		
		// Convert special directory names to route parameters
		if strings.HasSuffix(part, "_") {
			// locale_ -> {locale}
			paramName := strings.TrimSuffix(part, "_")
			routeParts = append(routeParts, "{"+paramName+"}")
		} else {
			routeParts = append(routeParts, part)
		}
	}
	
	if len(routeParts) == 0 {
		return "/"
	}
	
	return "/" + strings.Join(routeParts, "/")
}

// Global instance for consistent key generation across the application
var DefaultKeyGenerator = NewTemplateKeyGenerator()

// Convenience functions for global access
func GenerateTemplateKey(templatePath string) string {
	return DefaultKeyGenerator.GenerateTemplateKey(templatePath)
}

func GenerateRouteKey(routePattern string) string {
	return DefaultKeyGenerator.GenerateRouteKey(routePattern)
}

func ExtractRouteFromTemplatePath(templatePath string) string {
	return DefaultKeyGenerator.ExtractRouteFromTemplatePath(templatePath)
}