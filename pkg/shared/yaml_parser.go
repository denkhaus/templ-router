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
// Supports nested structures by flattening them with dot notation
// For multi-locale configurations, this returns empty map as translations are in MultiLocaleI18n
func extractI18nMappings(rawConfig map[string]interface{}) map[string]string {
	i18nMappings := make(map[string]string)

	if i18nData, ok := rawConfig["i18n"].(map[interface{}]interface{}); ok {
		// Check if this is a multi-locale configuration
		isMultiLocale := false
		for key := range i18nData {
			if strKey, ok := key.(string); ok {
				if IsValidLocaleCode(strKey) {
					isMultiLocale = true
					break
				}
			}
		}
		
		if isMultiLocale {
			// Multi-locale configuration - return empty map
			// Translations will be handled by extractMultiLocaleI18n
			return i18nMappings
		}
		
		// Check if this is a simple key-value mapping (not nested)
		isSimpleMapping := true
		for _, value := range i18nData {
			if _, isMap := value.(map[interface{}]interface{}); isMap {
				isSimpleMapping = false
				break
			}
			if _, isMap := value.(map[string]interface{}); isMap {
				isSimpleMapping = false
				break
			}
		}
		
		if isSimpleMapping {
			// Simple flat mapping
			for key, value := range i18nData {
				if strKey, ok := key.(string); ok {
					if strValue, ok := value.(string); ok {
						i18nMappings[strKey] = strValue
					}
				}
			}
		} else {
			// Nested structure - flatten it
			flattenI18nMap(i18nData, "", i18nMappings)
		}
	} else if i18nData, ok := rawConfig["i18n"].(map[string]interface{}); ok {
		// Check if this is a multi-locale configuration
		isMultiLocale := false
		for key := range i18nData {
			if IsValidLocaleCode(key) {
				isMultiLocale = true
				break
			}
		}
		
		if isMultiLocale {
			// Multi-locale configuration - return empty map
			// Translations will be handled by extractMultiLocaleI18n
			return i18nMappings
		}
		
		// Check if this is a simple key-value mapping (not nested)
		isSimpleMapping := true
		for _, value := range i18nData {
			if _, isMap := value.(map[string]interface{}); isMap {
				isSimpleMapping = false
				break
			}
			if _, isMap := value.(map[interface{}]interface{}); isMap {
				isSimpleMapping = false
				break
			}
		}
		
		if isSimpleMapping {
			// Simple flat mapping
			for key, value := range i18nData {
				if strValue, ok := value.(string); ok {
					i18nMappings[key] = strValue
				}
			}
		} else {
			// Nested structure - flatten it
			flattenI18nMapStringKeys(i18nData, "", i18nMappings)
		}
	}

	return i18nMappings
}

// IsValidLocaleCode checks if a string looks like a locale code
func IsValidLocaleCode(code string) bool {
	// Check for common 2-letter ISO 639-1 language codes
	validCodes := []string{
		"en", "de", "fr", "es", "it", "pt", "ru", "ja", "ko", "zh", "ar", "hi", 
		"nl", "sv", "da", "no", "fi", "pl", "tr", "he", "th", "vi", "cs", "hu",
		"ro", "bg", "hr", "sk", "sl", "et", "lv", "lt", "mt", "ga", "cy", "eu",
		"ca", "gl", "is", "fo", "kl", "se", "fi", "et", "lv", "lt", "be", "uk",
		"mk", "sq", "sr", "bs", "me", "xk",
	}
	
	// Only accept exact matches for known locale codes
	// This prevents words like "feedback" from being treated as locales
	for _, valid := range validCodes {
		if code == valid {
			return true
		}
	}
	
	// Also check for common locale patterns like "en-US", "de-DE", etc.
	if len(code) == 5 && code[2] == '-' {
		langCode := code[:2]
		for _, valid := range validCodes {
			if langCode == valid {
				return true
			}
		}
	}
	
	return false
}

// extractMultiLocaleI18n extracts multi-locale i18n mappings from the raw config
// Supports nested structures by flattening them with dot notation
// Only processes keys that are valid locale codes
func extractMultiLocaleI18n(rawConfig map[string]interface{}) map[string]map[string]string {
	multiLocaleI18n := make(map[string]map[string]string)

	if i18nData, ok := rawConfig["i18n"].(map[interface{}]interface{}); ok {
		for localeKey, localeValue := range i18nData {
			if localeStr, ok := localeKey.(string); ok {
				// Only process if this is a valid locale code
				if IsValidLocaleCode(localeStr) {
					if localeTranslations, ok := localeValue.(map[interface{}]interface{}); ok {
						translations := make(map[string]string)
						flattenI18nMap(localeTranslations, "", translations)
						multiLocaleI18n[localeStr] = translations
					}
				}
			}
		}
	} else if i18nData, ok := rawConfig["i18n"].(map[string]interface{}); ok {
		for localeStr, localeValue := range i18nData {
			// Only process if this is a valid locale code
			if IsValidLocaleCode(localeStr) {
				if localeTranslations, ok := localeValue.(map[string]interface{}); ok {
					translations := make(map[string]string)
					flattenI18nMapStringKeys(localeTranslations, "", translations)
					multiLocaleI18n[localeStr] = translations
				}
			}
		}
	}

	return multiLocaleI18n
}

// flattenI18nMap recursively flattens nested i18n structures with interface{} keys
// Example: {"feedback": {"title": "Dashboard"}} becomes {"feedback.title": "Dashboard"}
func flattenI18nMap(data map[interface{}]interface{}, prefix string, result map[string]string) {
	for key, value := range data {
		keyStr, ok := key.(string)
		if !ok {
			continue
		}
		
		currentKey := keyStr
		if prefix != "" {
			currentKey = prefix + "." + keyStr
		}
		
		switch v := value.(type) {
		case string:
			// Direct string value
			result[currentKey] = v
		case map[interface{}]interface{}:
			// Nested map - recurse
			flattenI18nMap(v, currentKey, result)
		case map[string]interface{}:
			// Convert and recurse
			converted := make(map[interface{}]interface{})
			for k, val := range v {
				converted[k] = val
			}
			flattenI18nMap(converted, currentKey, result)
		}
	}
}

// flattenI18nMapStringKeys recursively flattens nested i18n structures with string keys
// Example: {"feedback": {"title": "Dashboard"}} becomes {"feedback.title": "Dashboard"}
func flattenI18nMapStringKeys(data map[string]interface{}, prefix string, result map[string]string) {
	for key, value := range data {
		currentKey := key
		if prefix != "" {
			currentKey = prefix + "." + key
		}
		
		switch v := value.(type) {
		case string:
			// Direct string value
			result[currentKey] = v
		case map[string]interface{}:
			// Nested map - recurse
			flattenI18nMapStringKeys(v, currentKey, result)
		case map[interface{}]interface{}:
			// Convert and recurse
			converted := make(map[string]interface{})
			for k, val := range v {
				if strKey, ok := k.(string); ok {
					converted[strKey] = val
				}
			}
			flattenI18nMapStringKeys(converted, currentKey, result)
		}
	}
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