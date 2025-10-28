package i18n

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/metadata"
	"github.com/denkhaus/templ-router/pkg/shared"
	"go.uber.org/zap"
)

// ExtendedConfigFile extends ConfigFile to support multi-locale i18n
type ExtendedConfigFile struct {
	*interfaces.ConfigFile
	MultiLocaleI18n map[string]map[string]string `yaml:"i18n"`
}

// ParseYAMLMetadataExtended parses YAML with support for multi-locale i18n and nested structures
func ParseYAMLMetadataExtended(filePath string, logger *zap.Logger) (bool, *ExtendedConfigFile, error) {

	if filePath == "" {
		return false, nil, fmt.Errorf("file path cannot be empty")
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	logger.Debug("Parsing extended YAML metadata with nested i18n support", zap.String("file_path", filePath))

	// Use the enhanced shared parser that supports nested structures
	configFileFound, sharedConfig, err := shared.ParseYAMLMetadata(filePath)
	if err != nil {
		return configFileFound, nil, fmt.Errorf("failed to parse YAML file %s: %w", filePath, err)
	}

	// Check if this is a multi-locale configuration
	hasMultiLocale := len(sharedConfig.MultiLocaleI18n) > 0

	if hasMultiLocale {
		logger.Debug("Detected multi-locale YAML format with nested support",
			zap.String("file_path", filePath),
			zap.Int("locales_count", len(sharedConfig.MultiLocaleI18n)))

		// Convert shared.AuthSettings to interfaces.AuthSettings
		authSettings := &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}
		if sharedConfig.AuthSettings != nil {
			// Create a settings parser to handle auth settings conversion
			settingsParser := metadata.NewMetadataSettingsParser()
			if parsedAuth := settingsParser.ParseAuthSettings(sharedConfig.AuthSettings); parsedAuth != nil {
				authSettings = parsedAuth
			}
		}

		// Convert DynamicSettings if present
		var dynamicSettings *interfaces.DynamicSettings
		if sharedConfig.DynamicSettings != nil {
			settingsParser := metadata.NewMetadataSettingsParser()
			dynamicSettings = settingsParser.ParseDynamicSettings(sharedConfig.DynamicSettings)
		}

		// Create extended config with multi-locale support
		extendedConfig := &ExtendedConfigFile{
			ConfigFile: &interfaces.ConfigFile{
				FilePath:         sharedConfig.FilePath,
				TemplateFilePath: sharedConfig.TemplateFilePath,
				RouteMetadata:    sharedConfig.RouteMetadata,
				AuthSettings:     authSettings,
				LayoutSettings:   sharedConfig.LayoutSettings,
				ErrorSettings:    sharedConfig.ErrorSettings,
				DynamicSettings:  dynamicSettings,
				I18nMappings:     make(map[string]string), // Empty for multi-locale
				MultiLocaleI18n:  sharedConfig.MultiLocaleI18n,
			},
			MultiLocaleI18n: sharedConfig.MultiLocaleI18n,
		}

		return configFileFound, extendedConfig, nil
	}

	// Single-locale configuration (could be nested or flat)
	logger.Debug("Using single-locale YAML format with nested support", zap.String("file_path", filePath))

	// Convert shared.AuthSettings to interfaces.AuthSettings
	authSettings := &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}
	if sharedConfig.AuthSettings != nil {
		settingsParser := metadata.NewMetadataSettingsParser()
		if parsedAuth := settingsParser.ParseAuthSettings(sharedConfig.AuthSettings); parsedAuth != nil {
			authSettings = parsedAuth
		}
	}

	// Convert DynamicSettings if present
	var dynamicSettings *interfaces.DynamicSettings
	if sharedConfig.DynamicSettings != nil {
		settingsParser := metadata.NewMetadataSettingsParser()
		dynamicSettings = settingsParser.ParseDynamicSettings(sharedConfig.DynamicSettings)
	}

	// Create extended config for single-locale
	routerConfig := &interfaces.ConfigFile{
		FilePath:         sharedConfig.FilePath,
		TemplateFilePath: sharedConfig.TemplateFilePath,
		RouteMetadata:    sharedConfig.RouteMetadata,
		I18nMappings:     sharedConfig.I18nMappings, // Contains flattened nested keys
		MultiLocaleI18n:  sharedConfig.MultiLocaleI18n,
		LayoutSettings:   sharedConfig.LayoutSettings,
		ErrorSettings:    sharedConfig.ErrorSettings,
		DynamicSettings:  dynamicSettings,
		AuthSettings:     authSettings,
	}

	return configFileFound, &ExtendedConfigFile{
		ConfigFile:      routerConfig,
		MultiLocaleI18n: nil, // Empty for single-locale
	}, nil
}

// HasMultiLocaleSupport checks if this config supports multiple locales
func (ecf *ExtendedConfigFile) HasMultiLocaleSupport() bool {
	return len(ecf.MultiLocaleI18n) > 0
}
