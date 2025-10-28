package middleware

import (
	"context"
	"io"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// layoutServiceImpl implements LayoutService
type layoutServiceImpl struct {
	config            interfaces.ConfigService
	fileSystemChecker FileSystemChecker
	templateService   interfaces.TemplateService
	logger            *zap.Logger
}

// NewLayoutService creates a new layout service for DI
func NewLayoutService(i do.Injector) (interfaces.LayoutService, error) {
	config := do.MustInvoke[interfaces.ConfigService](i)
	fileSystemChecker := do.MustInvoke[FileSystemChecker](i)
	templateService := do.MustInvoke[interfaces.TemplateService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	return &layoutServiceImpl{
		config:            config,
		fileSystemChecker: fileSystemChecker,
		templateService:   templateService,
		logger:            logger,
	}, nil
}

// FindLayoutForTemplate finds the appropriate layout for a template using Next.js-style layout inheritance
func (ls *layoutServiceImpl) FindLayoutForTemplate(templatePath string) *interfaces.LayoutTemplate {
	ls.logger.Debug("Finding layout for template", zap.String("template_path", templatePath))

	// Get layout configuration
	rootDir := ls.config.GetLayoutRootDirectory()
	layoutFileName := ls.config.GetLayoutFileName() + ls.config.GetTemplateExtension()
	metadataExtension := ls.config.GetMetadataExtension()

	// Start from the template's directory and walk up the directory tree
	dir := filepath.Dir(templatePath)

	layoutLevel := 0

	for {
		layoutPath := filepath.Join(dir, layoutFileName)
		ls.logger.Debug("Checking for layout",
			zap.String("layout_path", layoutPath),
			zap.String("current_dir", dir),
			zap.String("root_dir", rootDir))

		// Check if layout file actually exists at this level (library-agnostic)
		if ls.fileSystemChecker.FileExists(layoutPath) {
			// Build metadata path correctly: layout.templ -> layout.templ.yaml
			metadataPath := strings.TrimSuffix(layoutPath, ls.config.GetTemplateExtension()) + metadataExtension

			ls.logger.Info("Found layout file",
				zap.String("layout_path", layoutPath),
				zap.String("metadata_path", metadataPath),
				zap.Int("layout_level", layoutLevel))

			return &interfaces.LayoutTemplate{
				FilePath:    layoutPath,
				YamlPath:    metadataPath,
				LayoutLevel: layoutLevel,
			}
		}

		// Move to parent directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir || parentDir == "." || parentDir == "" {
			// Reached filesystem root without finding layout
			ls.logger.Debug("Reached filesystem root, no layout found",
				zap.String("template_path", templatePath),
				zap.String("searched_up_to", dir))
			return nil
		}

		dir = parentDir
		layoutLevel++

		// Prevent infinite loops - reasonable depth limit
		if layoutLevel > 10 {
			ls.logger.Warn("Layout search depth limit reached",
				zap.String("template_path", templatePath),
				zap.Int("max_depth", layoutLevel))
			return nil
		}
	}
}

// WrapInLayout wraps a component in a layout
func (ls *layoutServiceImpl) WrapInLayout(component templ.Component, layout *interfaces.LayoutTemplate, ctx context.Context) templ.Component {
	ls.logger.Debug("Wrapping component in layout", zap.String("layout_path", layout.FilePath))

	// Load layout metadata and merge with existing template metadata
	if layout.YamlPath != "" {
		layoutConfig, err := shared.ParseYAMLMetadata(layout.YamlPath)
		if err != nil {
			ls.logger.Warn("Failed to load layout metadata",
				zap.String("yaml_path", layout.YamlPath),
				zap.Error(err))
		} else {
			// CRITICAL FIX: Template metadata should override layout metadata
			// Get existing template config from context
			if existingConfig := ctx.Value(shared.TemplateConfigKey); existingConfig != nil {
				if templateConfig, ok := existingConfig.(*shared.ConfigFile); ok {
					// Merge configs: template metadata takes precedence over layout metadata
					mergedConfig := mergeConfigs(layoutConfig, templateConfig)
					ctx = context.WithValue(ctx, shared.TemplateConfigKey, mergedConfig)

					// Safe access to metadata
					templateTitle := ""
					layoutTitle := ""
					if mergedConfig.RouteMetadata != nil {
						if metadataMap, ok := mergedConfig.RouteMetadata.(map[string]interface{}); ok {
							if title, exists := metadataMap["title"]; exists {
								if titleStr, ok := title.(string); ok {
									templateTitle = titleStr
								}
							}
						}
					}
					if layoutConfig.RouteMetadata != nil {
						if metadataMap, ok := layoutConfig.RouteMetadata.(map[string]interface{}); ok {
							if title, exists := metadataMap["title"]; exists {
								if titleStr, ok := title.(string); ok {
									layoutTitle = titleStr
								}
							}
						}
					}

					ls.logger.Info("Merged template and layout metadata (template takes precedence)",
						zap.String("layout_yaml", layout.YamlPath),
						zap.String("template_title", templateTitle),
						zap.String("layout_title", layoutTitle))
				} else {
					// Fallback: use layout config if template config is invalid
					ctx = context.WithValue(ctx, shared.TemplateConfigKey, layoutConfig)
					ls.logger.Info("Added layout metadata to context (fallback)",
						zap.String("yaml_path", layout.YamlPath),
						zap.Any("metadata", layoutConfig.RouteMetadata))
				}
			} else {
				// No existing template config, use layout config
				ctx = context.WithValue(ctx, shared.TemplateConfigKey, layoutConfig)
				ls.logger.Info("Added layout metadata to context (no template config)",
					zap.String("yaml_path", layout.YamlPath),
					zap.Any("metadata", layoutConfig.RouteMetadata))
			}
		}
	}

	// Create a wrapped component that includes the layout context and template service
	return &LayoutWrappedComponent{
		innerComponent:  component,
		layoutContext:   ctx,
		layoutPath:      layout.FilePath,
		templateService: ls.templateService,
		logger:          ls.logger,
	}
}

// LayoutWrappedComponent wraps a component with layout context
type LayoutWrappedComponent struct {
	innerComponent  templ.Component
	layoutContext   context.Context
	layoutPath      string
	templateService interfaces.TemplateService
	logger          *zap.Logger
}

// Render renders the component with the layout context
func (lwc *LayoutWrappedComponent) Render(ctx context.Context, w io.Writer) error {
	lwc.logger.Debug("Rendering component with layout context", zap.String("layout_path", lwc.layoutPath))

	// LIBRARY-AGNOSTIC: Render any layout template using template service
	if lwc.layoutPath != "" {
		lwc.logger.Info("Rendering layout via template service", zap.String("layout_path", lwc.layoutPath))

		// Use template service to render layout with content (library-agnostic)
		layoutComponent, err := lwc.templateService.RenderLayoutComponent(lwc.layoutPath, lwc.innerComponent, lwc.layoutContext)
		if err != nil {
			lwc.logger.Warn("Failed to render layout",
				zap.String("layout_path", lwc.layoutPath),
				zap.Error(err))
		} else {
			lwc.logger.Info("Successfully rendered layout component", zap.String("layout_path", lwc.layoutPath))
			return layoutComponent.Render(lwc.layoutContext, w)
		}
	}

	// Use the layout context instead of the original context
	// This ensures router.M() has access to template_config
	return lwc.innerComponent.Render(lwc.layoutContext, w)
}

// mergeConfigs merges layout and template configs with template taking precedence
func mergeConfigs(layoutConfig, templateConfig *shared.ConfigFile) *shared.ConfigFile {
	// Start with layout config as base
	merged := &shared.ConfigFile{
		RouteMetadata:   layoutConfig.RouteMetadata,
		MultiLocaleI18n: make(map[string]map[string]string),
		AuthSettings:    layoutConfig.AuthSettings, // Use layout auth settings as default
	}

	// Copy layout i18n data first
	if layoutConfig.MultiLocaleI18n != nil {
		for locale, translations := range layoutConfig.MultiLocaleI18n {
			merged.MultiLocaleI18n[locale] = make(map[string]string)
			for key, value := range translations {
				merged.MultiLocaleI18n[locale][key] = value
			}
		}
	}

	// Override with template metadata (template takes precedence)
	if templateConfig.RouteMetadata != nil {
		merged.RouteMetadata = templateConfig.RouteMetadata
	}

	// Override with template i18n data (template takes precedence)
	if templateConfig.MultiLocaleI18n != nil {
		for locale, translations := range templateConfig.MultiLocaleI18n {
			if merged.MultiLocaleI18n[locale] == nil {
				merged.MultiLocaleI18n[locale] = make(map[string]string)
			}
			for key, value := range translations {
				merged.MultiLocaleI18n[locale][key] = value
			}
		}
	}

	// Use template auth settings if available
	if templateConfig.AuthSettings != nil {
		merged.AuthSettings = templateConfig.AuthSettings
	}

	return merged
}
