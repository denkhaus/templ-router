package router

import (
	"fmt"
	"os"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// ExtendedConfigFile extends ConfigFile to support multi-locale i18n
type ExtendedConfigFile struct {
	*ConfigFile
	MultiLocaleI18n map[string]map[string]string `yaml:"i18n"`
}

// ParseYAMLMetadataExtended parses YAML with support for multi-locale i18n
func ParseYAMLMetadataExtended(filePath string, logger *zap.Logger) (*ExtendedConfigFile, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	if logger == nil {
		logger = zap.NewNop()
	}

	logger.Debug("Parsing extended YAML metadata", zap.String("file_path", filePath))

	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %w", filePath, err)
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

	if err := yaml.Unmarshal(data, &multiLocaleConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in file %s: %w", filePath, err)
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
			authSettings := &AuthSettings{Type: AuthTypePublic}
			// Create a settings parser to handle auth settings
			settingsParser := &MetadataSettingsParser{}
			if parsedAuth := settingsParser.parseAuthSettings(multiLocaleConfig.Auth); parsedAuth != nil {
				// Simple conversion - just copy the type
				switch parsedAuth.Type {
				case interfaces.AuthTypePublic:
					authSettings.Type = AuthTypePublic
				case interfaces.AuthTypeUser:
					authSettings.Type = AuthTypeUserRequired
				case interfaces.AuthTypeAdmin:
					authSettings.Type = AuthTypeAdminRequired
				}
			}

			// Create extended config
			extendedConfig := &ExtendedConfigFile{
				ConfigFile: &ConfigFile{
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

			return extendedConfig, nil
		}
	}

	// Fall back to simple format
	logger.Debug("Using simple YAML format", zap.String("file_path", filePath))

	// Create a metadata parser to handle simple format
	metadataParser := &MetadataParser{}
	simpleConfig, err := metadataParser.ParseYAMLMetadata(filePath)
	if err != nil {
		return nil, err
	}

	// Convert interfaces.ConfigFile back to router.ConfigFile for compatibility
	routerConfig := &ConfigFile{
		FilePath:         simpleConfig.FilePath,
		TemplateFilePath: simpleConfig.TemplateFilePath,
		RouteMetadata:    simpleConfig.RouteMetadata,
		I18nMappings:     simpleConfig.I18nMappings,
		MultiLocaleI18n:  simpleConfig.MultiLocaleI18n,
		LayoutSettings:   simpleConfig.LayoutSettings,
		ErrorSettings:    simpleConfig.ErrorSettings,
		// Convert AuthSettings back
		AuthSettings: &AuthSettings{Type: AuthTypePublic},
	}

	if simpleConfig.AuthSettings != nil {
		switch simpleConfig.AuthSettings.Type {
		case interfaces.AuthTypePublic:
			routerConfig.AuthSettings.Type = AuthTypePublic
		case interfaces.AuthTypeUser:
			routerConfig.AuthSettings.Type = AuthTypeUserRequired
		case interfaces.AuthTypeAdmin:
			routerConfig.AuthSettings.Type = AuthTypeAdminRequired
		}
	}

	return &ExtendedConfigFile{
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

// // GetTranslationsForLocale extracts translations for a specific locale
// func (ecf *ExtendedConfigFile) GetTranslationsForLocale(locale string) map[string]string {
// 	if ecf.MultiLocaleI18n != nil {
// 		if translations, exists := ecf.MultiLocaleI18n[locale]; exists {
// 			return translations
// 		}
// 	}

// 	// Fallback to simple format (assume it's English)
// 	if locale == "en" && len(ecf.ConfigFile.I18nMappings) > 0 {
// 		return ecf.ConfigFile.I18nMappings
// 	}

// 	return make(map[string]string)
// }

// // GetAvailableLocales returns all available locales in this config
// func (ecf *ExtendedConfigFile) GetAvailableLocales() []string {
// 	if ecf.MultiLocaleI18n != nil {
// 		locales := make([]string, 0, len(ecf.MultiLocaleI18n))
// 		for locale := range ecf.MultiLocaleI18n {
// 			locales = append(locales, locale)
// 		}
// 		return locales
// 	}

// 	// Simple format defaults to English
// 	if len(ecf.ConfigFile.I18nMappings) > 0 {
// 		return []string{"en"}
// 	}

// 	return []string{}
// }

// HasMultiLocaleSupport checks if this config supports multiple locales
func (ecf *ExtendedConfigFile) HasMultiLocaleSupport() bool {
	return len(ecf.MultiLocaleI18n) > 0
}
