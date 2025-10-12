package middleware

import (
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// ErrorTemplateResolver interface for error template resolution (PUBLIC)
type ErrorTemplateResolver interface {
	FindErrorTemplateForPath(path string) *interfaces.ErrorTemplate
}

// errorTemplateResolverImpl implements ErrorTemplateResolver (PRIVATE)
type errorTemplateResolverImpl struct {
	configService     interfaces.ConfigService
	fileSystemChecker FileSystemChecker
	logger            *zap.Logger
}

// NewErrorTemplateResolver creates a new error template resolver (RETURNS INTERFACE)
func NewErrorTemplateResolver(i do.Injector) (ErrorTemplateResolver, error) {
	configService := do.MustInvoke[interfaces.ConfigService](i)
	fileSystemChecker := do.MustInvoke[FileSystemChecker](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &errorTemplateResolverImpl{
		configService:     configService,
		fileSystemChecker: fileSystemChecker,
		logger:            logger,
	}, nil
}

// FindErrorTemplateForPath finds the most specific error template for a given path
// This implements the hierarchical error template resolution logic
func (etr *errorTemplateResolverImpl) FindErrorTemplateForPath(path string) *interfaces.ErrorTemplate {
	etr.logger.Debug("Finding error template for path", zap.String("path", path))

	// Get template root directory from config
	templateRoot := etr.configService.GetLayoutRootDirectory()

	// Generate candidate paths for error templates (most specific to least specific)
	candidatePaths := etr.generateErrorTemplateCandidates(path, templateRoot)

	// Try each candidate path
	for _, candidatePath := range candidatePaths {
		if etr.fileSystemChecker.FileExists(candidatePath) {
			etr.logger.Debug("Found error template",
				zap.String("path", path),
				zap.String("template", candidatePath))

			return &interfaces.ErrorTemplate{
				FilePath:      candidatePath,
				ComponentName: etr.generateComponentName(candidatePath),
				ErrorCode:     etr.extractErrorCodeFromPath(candidatePath),
			}
		}
	}

	etr.logger.Debug("No specific error template found", zap.String("path", path))
	return nil
}

// generateErrorTemplateCandidates generates candidate paths for error templates
// in order of specificity (most specific first)
func (etr *errorTemplateResolverImpl) generateErrorTemplateCandidates(path, templateRoot string) []string {
	var candidates []string

	// Clean and normalize the path
	cleanPath := strings.Trim(path, "/")
	pathSegments := strings.Split(cleanPath, "/")

	// Generate candidates from most specific to least specific
	// Example: /admin/users/123 -> [admin/users/error.templ, admin/error.templ, error.templ]

	for i := len(pathSegments); i >= 0; i-- {
		var candidatePath string

		if i == 0 {
			// Root error template
			candidatePath = filepath.Join(templateRoot, "error.templ")
		} else {
			// Directory-specific error template
			dirPath := strings.Join(pathSegments[:i], "/")
			candidatePath = filepath.Join(templateRoot, dirPath, "error.templ")
		}

		candidates = append(candidates, candidatePath)
	}

	etr.logger.Debug("Generated error template candidates",
		zap.String("path", path),
		zap.Strings("candidates", candidates))

	return candidates
}

// generateComponentName generates a component name from the template file path
func (etr *errorTemplateResolverImpl) generateComponentName(templatePath string) string {
	// Convert file path to component name
	// Example: app/admin/error.templ -> AdminError

	dir := filepath.Dir(templatePath)
	base := filepath.Base(dir)

	if base == "." || base == "/" {
		return "Error"
	}

	// Capitalize first letter and add "Error" suffix
	componentName := strings.Title(base) + "Error"
	return componentName
}

// extractErrorCodeFromPath extracts error code from template path if present
func (etr *errorTemplateResolverImpl) extractErrorCodeFromPath(templatePath string) int {
	// Check if the path contains specific error codes
	// Example: app/404/error.templ -> 404

	dir := filepath.Dir(templatePath)
	segments := strings.Split(dir, "/")

	for _, segment := range segments {
		switch segment {
		case "404":
			return 404
		case "500":
			return 500
		case "403":
			return 403
		case "401":
			return 401
		}
	}

	// Default to 500 for generic error templates
	return 500
}
