package i18n

import (
	"fmt"
	"os"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/metadata"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// ExtendedConfigFile extends ConfigFile to support multi-locale i18n
type ExtendedConfigFile struct {
	*interfaces.ConfigFile
	MultiLocaleI18n map[string]map[string]string `yaml:"i18n"`
}

// ParseYAMLMetadataExtended parses YAML with support for multi-locale i18n
func ParseYAMLMetadataExtended(filePath string, logger *zap.Logger) (bool, *ExtendedConfigFile, error) {
	configFileFound := false

	if filePath == "" {
		return configFileFound, nil, fmt.Errorf("file path cannot be empty")
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	logger.Debug("Parsing extended YAML metadata", zap.String("file_path", filePath))

	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return configFileFound, nil, fmt.Errorf("failed to read YAML file %s: %w", filePath, err)
	}

	// First, try to parse as multi-locale format
	var multiLocaleConfig struct {
		Route   interface{}                  `yaml:"route"`
		I18n    map[string]map[string]string `yaml:"i18n"`
		Auth    interface{}                  `yaml:"auth"`
		Layout  interface{}                  `yaml:"layout"`
		Error   interface{}                  `yaml:"error"`
		Dynamic interface{}                  `yaml:"dynamic"`
	}

	configFileFound = true
	if err := yaml.Unmarshal(data, &multiLocaleConfig); err != nil {
		return configFileFound, nil, fmt.Errorf("failed to parse YAML in file %s: %w", filePath, err)
	}

	// Check if it's multi-locale format
	if len(multiLocaleConfig.I18n) > 0 {
		// Check if the first level contains locale codes
		isMultiLocale := false
		for key := range multiLocaleConfig.I18n {
			if isValidLocaleCode(key) {
				isMultiLocale = true
				break
			}
		}

		if isMultiLocale {
			logger.Debug("Detected multi-locale YAML format",
				zap.String("file_path", filePath),
				zap.Int("locales_count", len(multiLocaleConfig.I18n)))

			// Convert interfaces.AuthSettings back to router.AuthSettings for compatibility
			authSettings := &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}
			// Create a settings parser to handle auth settings
			settingsParser := metadata.NewMetadataSettingsParser()
			if parsedAuth := settingsParser.ParseAuthSettings(multiLocaleConfig.Auth); parsedAuth != nil {
				// Simple conversion - just copy the type
				switch parsedAuth.Type {
				case interfaces.AuthTypePublic:
					authSettings.Type = interfaces.AuthTypePublic
				case interfaces.AuthTypeUser:
					authSettings.Type = interfaces.AuthTypeUser
				case interfaces.AuthTypeAdmin:
					authSettings.Type = interfaces.AuthTypeAdmin
				}
			}

			// Create extended config
			extendedConfig := &ExtendedConfigFile{
				ConfigFile: &interfaces.ConfigFile{
					FilePath:         filePath,
					TemplateFilePath: filePath[:len(filePath)-len(".yaml")] + ".templ",
					RouteMetadata:    multiLocaleConfig.Route,
					AuthSettings:     authSettings,
					LayoutSettings:   multiLocaleConfig.Layout,
					ErrorSettings:    multiLocaleConfig.Error,
					I18nMappings:     make(map[string]string), // Will be populated per locale
				},
				MultiLocaleI18n: multiLocaleConfig.I18n,
			}

			return configFileFound, extendedConfig, nil
		}
	}

	// Fall back to simple format
	logger.Debug("Using simple YAML format", zap.String("file_path", filePath))

	// Create a metadata parser to handle simple format
	metadataParser := &metadata.MetadataParser{}
	simpleConfig, err := metadataParser.ParseYAMLMetadata(filePath)
	if err != nil {
		return configFileFound, nil, err
	}

	// Convert interfaces.ConfigFile back to router.ConfigFile for compatibility
	routerConfig := &interfaces.ConfigFile{
		FilePath:         simpleConfig.FilePath,
		TemplateFilePath: simpleConfig.TemplateFilePath,
		RouteMetadata:    simpleConfig.RouteMetadata,
		I18nMappings:     simpleConfig.I18nMappings,
		MultiLocaleI18n:  simpleConfig.MultiLocaleI18n,
		LayoutSettings:   simpleConfig.LayoutSettings,
		ErrorSettings:    simpleConfig.ErrorSettings,
		// Convert AuthSettings back
		AuthSettings: &interfaces.AuthSettings{Type: interfaces.AuthTypePublic},
	}

	if simpleConfig.AuthSettings != nil {
		switch simpleConfig.AuthSettings.Type {
		case interfaces.AuthTypePublic:
			routerConfig.AuthSettings.Type = interfaces.AuthTypePublic
		case interfaces.AuthTypeUser:
			routerConfig.AuthSettings.Type = interfaces.AuthTypeUser
		case interfaces.AuthTypeAdmin:
			routerConfig.AuthSettings.Type = interfaces.AuthTypeAdmin
		}
	}

	return configFileFound, &ExtendedConfigFile{
		ConfigFile:      routerConfig,
		MultiLocaleI18n: nil,
	}, nil
}

// isValidLocaleCode checks if a string looks like a locale code
func isValidLocaleCode(code string) bool {
	// Simple check for common locale patterns
	validCodes := []string{"en", "de", "fr", "es", "it", "pt", "ru", "ja", "ko", "zh", "ar", "hi"}
	for _, valid := range validCodes {
		if code == valid {
			return true
		}
	}
	return false
}

// HasMultiLocaleSupport checks if this config supports multiple locales
func (ecf *ExtendedConfigFile) HasMultiLocaleSupport() bool {
	return len(ecf.MultiLocaleI18n) > 0
}
