package i18n

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/shared"
	"go.uber.org/zap"
)

// ExtendedConfigFile extends ConfigFile to support multi-locale i18n
type ExtendedConfigFile struct {
	*shared.ConfigFile
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
	multiLocaleI18n := sharedConfig.GetMultiLocaleI18n()
	hasMultiLocale := len(multiLocaleI18n) > 0

	if hasMultiLocale {
		logger.Debug("Detected multi-locale YAML format with nested support",
			zap.String("file_path", filePath),
			zap.Int("locales_count", len(multiLocaleI18n)))

		// Create extended config with multi-locale support using the shared config directly
		extendedConfig := &ExtendedConfigFile{
			ConfigFile:      sharedConfig, // Use the shared config directly since shared.ConfigFile is aliased
			MultiLocaleI18n: multiLocaleI18n,
		}

		return configFileFound, extendedConfig, nil
	}

	// Single-locale configuration (could be nested or flat)
	logger.Debug("Using single-locale YAML format with nested support", zap.String("file_path", filePath))

	// Create extended config for single-locale using the shared config directly
	return configFileFound, &ExtendedConfigFile{
		ConfigFile:      sharedConfig, // Use the shared config directly since shared.ConfigFile is aliased
		MultiLocaleI18n: nil,          // Empty for single-locale
	}, nil
}

// HasMultiLocaleSupport checks if this config supports multiple locales
func (ecf *ExtendedConfigFile) HasMultiLocaleSupport() bool {
	return len(ecf.MultiLocaleI18n) > 0
}
