package template

import (
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"go.uber.org/zap"
)

// KeyResolver handles template key generation and resolution
type KeyResolver struct {
	logger     *zap.Logger
	config     interfaces.ConfigService
	moduleName string
}

// NewKeyResolver creates a new key resolver
func NewKeyResolver(logger *zap.Logger, config interfaces.ConfigService, moduleName string) *KeyResolver {
	return &KeyResolver{
		logger:     logger,
		config:     config,
		moduleName: moduleName,
	}
}

// CreateTemplateKeyFromPath creates a template key from file path and function name
func (kr *KeyResolver) CreateTemplateKeyFromPath(filePath, functionName string) string {
	// Convert file path to template key
	// app/$locale/dashboard/page_templ.go + DashboardPage -> dashboard.DashboardPage
	// app/page_templ.go + HomePage -> Page

	dir := filepath.Dir(filePath)

	// Remove configurable root directory prefix
	rootDir := kr.config.GetLayoutRootDirectory()
	if strings.HasPrefix(dir, rootDir) {
		dir = strings.TrimPrefix(dir, rootDir)
		dir = strings.TrimPrefix(dir, "/")
	}

	// Handle root directory case
	if dir == "" || dir == "." {
		return functionName
	}

	// Convert directory path to dot notation
	// locale_/dashboard -> locale.dashboard
	parts := strings.Split(dir, "/")
	var cleanParts []string

	for _, part := range parts {
		if part != "" {
			// Convert $locale pattern to locale
			if strings.HasSuffix(part, "_") {
				part = strings.TrimSuffix(part, "_")
			}
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) == 0 {
		return functionName
	}

	// Create hierarchical key
	templateKey := strings.Join(cleanParts, ".") + "." + functionName

	kr.logger.Debug("Created template key",
		zap.String("file_path", filePath),
		zap.String("function_name", functionName),
		zap.String("template_key", templateKey))

	return templateKey
}

// ResolveTemplateKey resolves a template key to its components
func (kr *KeyResolver) ResolveTemplateKey(templateKey string) (string, string) {
	// Split template key into path and function name
	// dashboard.DashboardPage -> dashboard, DashboardPage
	parts := strings.Split(templateKey, ".")
	if len(parts) < 2 {
		return "", templateKey
	}

	functionName := parts[len(parts)-1]
	pathParts := parts[:len(parts)-1]
	path := strings.Join(pathParts, "/")

	return path, functionName
}

// NormalizeTemplateKey normalizes a template key for consistent lookup
func (kr *KeyResolver) NormalizeTemplateKey(templateKey string) string {
	// Convert to lowercase and normalize separators
	normalized := strings.ToLower(templateKey)
	normalized = strings.ReplaceAll(normalized, "/", ".")
	normalized = strings.ReplaceAll(normalized, "_", ".")

	kr.logger.Debug("Normalized template key",
		zap.String("original", templateKey),
		zap.String("normalized", normalized))

	return normalized
}
