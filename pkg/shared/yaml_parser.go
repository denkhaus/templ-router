package shared

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigFile represents a YAML file containing metadata and settings
type ConfigFile struct {
	// FilePath is the full path to the YAML config file
	FilePath string
	
	// TemplateFilePath is the path to the corresponding *.templ file
	TemplateFilePath string
	
	// RouteMetadata contains custom route configuration (metadata from YAML)
	RouteMetadata interface{}
	
	// I18nMappings contains custom i18n identifier mappings
	I18nMappings map[string]string
	
	// MultiLocaleI18n contains multi-locale translations (locale -> key -> translation)
	MultiLocaleI18n map[string]map[string]string
	
	// AuthSettings contains authentication settings
	AuthSettings interface{}
	
	// LayoutSettings contains layout configuration
	LayoutSettings interface{}
	
	// ErrorSettings contains error handling configuration
	ErrorSettings interface{}
	
	// DynamicSettings contains dynamic parameter validation configuration
	DynamicSettings interface{}
}

// ParseYAMLMetadata parses YAML metadata files with validation
func ParseYAMLMetadata(filePath string) (*ConfigFile, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %w", filePath, err)
	}

	// Create a struct to decode the YAML into
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in file %s: %w", filePath, err)
	}
	
	// Convert map[interface{}]interface{} to map[string]interface{} for JSON compatibility
	converted := convertInterfaceMapToStringMap(rawConfig)
	if convertedMap, ok := converted.(map[string]interface{}); ok {
		rawConfig = convertedMap
	}

	// Validate that only known root keys are used
	if err := validateRootKeys(rawConfig); err != nil {
		return nil, fmt.Errorf("invalid YAML structure in file %s: %w", filePath, err)
	}

	// Create a ConfigFile struct with the parsed data
	configFile := &ConfigFile{
		FilePath:         filePath,
		RouteMetadata:    rawConfig["metadata"],
		I18nMappings:     extractI18nMappings(rawConfig),
		MultiLocaleI18n:  extractMultiLocaleI18n(rawConfig),
		AuthSettings:     rawConfig["auth"],
		LayoutSettings:   rawConfig["layout"],
		ErrorSettings:    rawConfig["error"],
		DynamicSettings:  rawConfig["dynamic"],
	}

	return configFile, nil
}

// validateRootKeys validates that only known root keys are used in YAML
func validateRootKeys(rawConfig map[string]interface{}) error {
	allowedKeys := map[string]bool{
		"i18n":     true,
		"auth":     true,
		"metadata": true,
		"layout":   true, // deprecated but still allowed
		"error":    true,
		"dynamic":  true,
		"route":    true, // legacy field name
	}
	
	for key := range rawConfig {
		if !allowedKeys[key] {
			return fmt.Errorf("unknown root key '%s' - allowed keys are: i18n, auth, metadata, layout, error, dynamic", key)
		}
	}
	
	return nil
}

// extractI18nMappings extracts i18n mappings from the raw config
func extractI18nMappings(rawConfig map[string]interface{}) map[string]string {
	i18nMappings := make(map[string]string)

	if i18nData, ok := rawConfig["i18n"].(map[interface{}]interface{}); ok {
		for key, value := range i18nData {
			if strKey, ok := key.(string); ok {
				if strValue, ok := value.(string); ok {
					i18nMappings[strKey] = strValue
				}
			}
		}
	} else if i18nData, ok := rawConfig["i18n"].(map[string]interface{}); ok {
		for key, value := range i18nData {
			if strValue, ok := value.(string); ok {
				i18nMappings[key] = strValue
			}
		}
	}

	return i18nMappings
}

// extractMultiLocaleI18n extracts multi-locale i18n mappings from the raw config
func extractMultiLocaleI18n(rawConfig map[string]interface{}) map[string]map[string]string {
	multiLocaleI18n := make(map[string]map[string]string)

	if i18nData, ok := rawConfig["i18n"].(map[interface{}]interface{}); ok {
		for localeKey, localeValue := range i18nData {
			if localeStr, ok := localeKey.(string); ok {
				if localeTranslations, ok := localeValue.(map[interface{}]interface{}); ok {
					translations := make(map[string]string)
					for key, value := range localeTranslations {
						if keyStr, ok := key.(string); ok {
							if valueStr, ok := value.(string); ok {
								translations[keyStr] = valueStr
							}
						}
					}
					multiLocaleI18n[localeStr] = translations
				}
			}
		}
	} else if i18nData, ok := rawConfig["i18n"].(map[string]interface{}); ok {
		for localeStr, localeValue := range i18nData {
			if localeTranslations, ok := localeValue.(map[string]interface{}); ok {
				translations := make(map[string]string)
				for key, value := range localeTranslations {
					if valueStr, ok := value.(string); ok {
						translations[key] = valueStr
					}
				}
				multiLocaleI18n[localeStr] = translations
			}
		}
	}

	return multiLocaleI18n
}

// convertInterfaceMapToStringMap recursively converts map[interface{}]interface{} to map[string]interface{}
func convertInterfaceMapToStringMap(input interface{}) interface{} {
	switch v := input.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if strKey, ok := key.(string); ok {
				result[strKey] = convertInterfaceMapToStringMap(value)
			}
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = convertInterfaceMapToStringMap(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = convertInterfaceMapToStringMap(item)
		}
		return result
	default:
		return v
	}
}