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
	fileSystemChecker interfaces.FileSystemChecker
	templateService   interfaces.TemplateService
	logger            *zap.Logger
}

// NewLayoutService creates a new layout service for DI
func NewLayoutService(i do.Injector) (interfaces.LayoutService, error) {
	config := do.MustInvoke[interfaces.ConfigService](i)
	fileSystemChecker := do.MustInvoke[interfaces.FileSystemChecker](i)
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
		configFileFound, layoutConfig, err := shared.ParseYAMLMetadata(layout.YamlPath)
		if err != nil {
			if configFileFound {
				ls.logger.Warn("Failed to load layout metadata",
					zap.String("yaml_path", layout.YamlPath),
					zap.Error(err),
				)
			}
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
					if mergedConfig.GetRouteMetadata() != nil {
						if metadataMap, ok := mergedConfig.GetRouteMetadata().(map[string]interface{}); ok {
							if title, exists := metadataMap["title"]; exists {
								if titleStr, ok := title.(string); ok {
									templateTitle = titleStr
								}
							}
						}
					}
					if layoutConfig.GetRouteMetadata() != nil {
						if metadataMap, ok := layoutConfig.GetRouteMetadata().(map[string]interface{}); ok {
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
						zap.Any("metadata", layoutConfig.GetRouteMetadata()))
				}
			} else {
				// No existing template config, use layout config
				ctx = context.WithValue(ctx, shared.TemplateConfigKey, layoutConfig)
				ls.logger.Info("Added layout metadata to context (no template config)",
					zap.String("yaml_path", layout.YamlPath),
					zap.Any("metadata", layoutConfig.GetRouteMetadata()))
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
		FilePath: layoutConfig.FilePath,
	}

	// Copy layout metadata first
	if layoutConfig.Metadata != nil {
		merged.Metadata = &shared.MetadataConfig{
			Title:       layoutConfig.Metadata.Title,
			Description: layoutConfig.Metadata.Description,
			Keywords:    make([]string, len(layoutConfig.Metadata.Keywords)),
			Author:      layoutConfig.Metadata.Author,
			Version:     layoutConfig.Metadata.Version,
			Custom:      make(map[string]interface{}),
		}
		copy(merged.Metadata.Keywords, layoutConfig.Metadata.Keywords)
		for k, v := range layoutConfig.Metadata.Custom {
			merged.Metadata.Custom[k] = v
		}
	}

	// Copy layout i18n data first
	if layoutConfig.I18n != nil {
		merged.I18n = &shared.I18nConfig{
			FlatMappings: make(map[string]string),
			Translations: make(map[string]*shared.LocaleTranslations),
		}
		// Copy flat mappings
		for k, v := range layoutConfig.I18n.FlatMappings {
			merged.I18n.FlatMappings[k] = v
		}
		// Copy translations
		for locale, translations := range layoutConfig.I18n.Translations {
			merged.I18n.Translations[locale] = &shared.LocaleTranslations{
				Locale:       translations.Locale,
				Translations: make(map[string]interface{}),
			}
			for k, v := range translations.Translations {
				merged.I18n.Translations[locale].Translations[k] = v
			}
		}
	}

	// Copy layout auth settings as default
	if layoutConfig.Auth != nil {
		merged.Auth = &shared.AuthConfig{
			Type:        layoutConfig.Auth.Type,
			RedirectURL: layoutConfig.Auth.RedirectURL,
			Roles:       make([]string, len(layoutConfig.Auth.Roles)),
			Settings:    make(map[string]interface{}),
		}
		copy(merged.Auth.Roles, layoutConfig.Auth.Roles)
		for k, v := range layoutConfig.Auth.Settings {
			merged.Auth.Settings[k] = v
		}
	}

	// Copy layout settings
	if layoutConfig.Layout != nil {
		merged.Layout = &shared.LayoutConfig{
			Template: layoutConfig.Layout.Template,
			Settings: make(map[string]interface{}),
		}
		for k, v := range layoutConfig.Layout.Settings {
			merged.Layout.Settings[k] = v
		}
	}

	// Copy error settings
	if layoutConfig.Error != nil {
		merged.Error = &shared.ErrorConfig{
			Template: layoutConfig.Error.Template,
			Settings: make(map[string]interface{}),
		}
		for k, v := range layoutConfig.Error.Settings {
			merged.Error.Settings[k] = v
		}
	}

	// Copy dynamic settings
	if layoutConfig.Dynamic != nil {
		merged.Dynamic = &shared.DynamicConfig{
			Rules:    make(map[string]*shared.ValidationRule),
			Settings: make(map[string]interface{}),
		}
		for k, v := range layoutConfig.Dynamic.Rules {
			merged.Dynamic.Rules[k] = &shared.ValidationRule{
				Name:     v.Name,
				Type:     v.Type,
				Required: v.Required,
				Pattern:  v.Pattern,
				Default:  v.Default,
				Settings: make(map[string]interface{}),
			}
			for rk, rv := range v.Settings {
				merged.Dynamic.Rules[k].Settings[rk] = rv
			}
		}
		for k, v := range layoutConfig.Dynamic.Settings {
			merged.Dynamic.Settings[k] = v
		}
	}

	// Now merge with template config (template takes precedence)
	merged.MergeMetadata(templateConfig)
	merged.MergeI18n(templateConfig)

	// Override auth settings if template has them
	if templateConfig.Auth != nil {
		merged.Auth = &shared.AuthConfig{
			Type:        templateConfig.Auth.Type,
			RedirectURL: templateConfig.Auth.RedirectURL,
			Roles:       make([]string, len(templateConfig.Auth.Roles)),
			Settings:    make(map[string]interface{}),
		}
		copy(merged.Auth.Roles, templateConfig.Auth.Roles)
		for k, v := range templateConfig.Auth.Settings {
			merged.Auth.Settings[k] = v
		}
	}

	// Override layout settings if template has them
	if templateConfig.Layout != nil {
		merged.Layout = &shared.LayoutConfig{
			Template: templateConfig.Layout.Template,
			Settings: make(map[string]interface{}),
		}
		for k, v := range templateConfig.Layout.Settings {
			merged.Layout.Settings[k] = v
		}
	}

	// Override error settings if template has them
	if templateConfig.Error != nil {
		merged.Error = &shared.ErrorConfig{
			Template: templateConfig.Error.Template,
			Settings: make(map[string]interface{}),
		}
		for k, v := range templateConfig.Error.Settings {
			merged.Error.Settings[k] = v
		}
	}

	// Override dynamic settings if template has them
	if templateConfig.Dynamic != nil {
		merged.Dynamic = &shared.DynamicConfig{
			Rules:    make(map[string]*shared.ValidationRule),
			Settings: make(map[string]interface{}),
		}
		for k, v := range templateConfig.Dynamic.Rules {
			merged.Dynamic.Rules[k] = &shared.ValidationRule{
				Name:     v.Name,
				Type:     v.Type,
				Required: v.Required,
				Pattern:  v.Pattern,
				Default:  v.Default,
				Settings: make(map[string]interface{}),
			}
			for rk, rv := range v.Settings {
				merged.Dynamic.Rules[k].Settings[rk] = rv
			}
		}
		for k, v := range templateConfig.Dynamic.Settings {
			merged.Dynamic.Settings[k] = v
		}
	}

	return merged
}
