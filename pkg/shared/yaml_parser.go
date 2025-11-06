package shared

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigFile represents a YAML file containing metadata and settings
type ConfigFile struct {
	// FilePath is the full path to the YAML config file
	FilePath string `yaml:"-" json:"-"`

	// TemplateFilePath is the path to the corresponding *.templ file
	TemplateFilePath string `yaml:"template_file_path,omitempty" json:"template_file_path,omitempty"`

	// Metadata contains custom route configuration
	Metadata *MetadataConfig `yaml:"metadata,omitempty" json:"metadata,omitempty"`

	// I18n contains internationalization configuration
	I18n *I18nConfig `yaml:"i18n,omitempty" json:"i18n,omitempty"`

	// Auth contains authentication settings
	Auth *AuthConfig `yaml:"auth,omitempty" json:"auth,omitempty"`

	// Layout contains layout configuration
	Layout *LayoutConfig `yaml:"layout,omitempty" json:"layout,omitempty"`

	// Error contains error handling configuration
	Error *ErrorConfig `yaml:"error,omitempty" json:"error,omitempty"`

	// Dynamic contains dynamic parameter validation configuration
	Dynamic *DynamicConfig `yaml:"dynamic,omitempty" json:"dynamic,omitempty"`
}

// MetadataConfig contains route metadata with nested structure support
type MetadataConfig struct {
	// Standard metadata fields
	Title       string   `yaml:"title,omitempty" json:"title,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Keywords    []string `yaml:"keywords,omitempty" json:"keywords,omitempty"`
	Author      string   `yaml:"author,omitempty" json:"author,omitempty"`
	Version     string   `yaml:"version,omitempty" json:"version,omitempty"`

	// Custom metadata fields - allows for arbitrary nested structure
	Custom map[string]interface{} `yaml:",inline" json:",inline"`
}

// I18nConfig contains internationalization configuration
type I18nConfig struct {
	// Simple flat mappings (legacy support)
	FlatMappings map[string]string `yaml:"flat,omitempty" json:"flat,omitempty"`

	// Multi-locale translations (locale -> key -> translation)
	Translations map[string]*LocaleTranslations `yaml:"translations,omitempty" json:"translations,omitempty"`
}

// LocaleTranslations contains translations for a specific locale
type LocaleTranslations struct {
	// Locale code (e.g., "en", "de", "fr")
	Locale string `yaml:"-" json:"-"`

	// Nested translation structure
	Translations map[string]interface{} `yaml:"translations,omitempty" json:"translations,omitempty"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	// Type of authentication: "Public", "UserRequired", "AdminRequired"
	Type string `yaml:"type" json:"type"`

	// Redirect URL for unauthenticated users
	RedirectURL string `yaml:"redirect_url,omitempty" json:"redirect_url,omitempty"`

	// Required roles (optional)
	Roles []string `yaml:"roles,omitempty" json:"roles,omitempty"`

	// Additional auth settings
	Settings map[string]interface{} `yaml:",inline" json:",inline"`
}

// LayoutConfig contains layout configuration
type LayoutConfig struct {
	// Layout template to use
	Template string `yaml:"template,omitempty" json:"template,omitempty"`

	// Layout settings
	Settings map[string]interface{} `yaml:",inline" json:",inline"`
}

// ErrorConfig contains error handling configuration
type ErrorConfig struct {
	// Error template to use
	Template string `yaml:"template,omitempty" json:"template,omitempty"`

	// Error handling settings
	Settings map[string]interface{} `yaml:",inline" json:",inline"`
}

// DynamicConfig contains dynamic parameter validation configuration
type DynamicConfig struct {
	// Validation rules
	Rules map[string]*ValidationRule `yaml:"rules,omitempty" json:"rules,omitempty"`

	// Additional dynamic settings
	Settings map[string]interface{} `yaml:"settings,omitempty" json:"settings,omitempty"`
}

// ValidationRule defines a validation rule for dynamic parameters
type ValidationRule struct {
	// Parameter name
	Name string `yaml:"name" json:"name"`

	// Validation type (string, int, email, etc.)
	Type string `yaml:"type" json:"type"`

	// Required flag
	Required bool `yaml:"required,omitempty" json:"required,omitempty"`

	// Validation pattern (regex for strings, min/max for numbers)
	Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty"`

	// Default value
	Default interface{} `yaml:"default,omitempty" json:"default,omitempty"`

	// Additional rule settings
	Settings map[string]interface{} `yaml:",inline" json:",inline"`
}

// ParseYAMLMetadata parses YAML metadata files with validation
func ParseYAMLMetadata(filePath string) (bool, *ConfigFile, error) {
	if filePath == "" {
		return false, nil, fmt.Errorf("file path cannot be empty")
	}

	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read YAML file %s: %w", filePath, err)
	}

	// First pass: decode into raw map for validation
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return true, nil, fmt.Errorf("failed to parse YAML in file %s: %w", filePath, err)
	}

	// Convert map[interface{}]interface{} to map[string]interface{} for JSON compatibility
	converted := convertInterfaceMapToStringMap(rawConfig)
	if convertedMap, ok := converted.(map[string]interface{}); ok {
		rawConfig = convertedMap
	}

	// Validate that only known root keys are used
	if err := validateRootKeys(rawConfig); err != nil {
		return true, nil, fmt.Errorf("invalid YAML structure in file %s: %w", filePath, err)
	}

	// Second pass: decode into our type-safe structure
	var configFile ConfigFile
	if err := yaml.Unmarshal(data, &configFile); err != nil {
		return true, nil, fmt.Errorf("failed to parse YAML into ConfigFile in %s: %w", filePath, err)
	}

	// Set the file path
	configFile.FilePath = filePath

	// Process i18n data to handle structured format
	if i18nData, ok := rawConfig["i18n"]; ok {
		configFile.I18n = processI18nData(i18nData)
	} else {
		// Ensure I18n is initialized
		configFile.I18n = &I18nConfig{
			FlatMappings: make(map[string]string),
			Translations: make(map[string]*LocaleTranslations),
		}
	}

	return true, &configFile, nil
}

// validateRootKeys validates that only known root keys are used in YAML
func validateRootKeys(rawConfig map[string]interface{}) error {
	allowedKeys := map[string]bool{
		"i18n":     true,
		"auth":     true,
		"metadata": true,
		"layout":   true,
		"error":    true,
		"dynamic":  true,
	}

	for key := range rawConfig {
		if !allowedKeys[key] {
			return fmt.Errorf("unknown root key '%s' - allowed keys are: i18n, auth, metadata, layout, error, dynamic", key)
		}
	}

	return nil
}

// processI18nData processes raw i18n data and converts it to our type-safe I18nConfig
func processI18nData(rawI18nData interface{}) *I18nConfig {
	i18nConfig := &I18nConfig{
		FlatMappings: make(map[string]string),
		Translations: make(map[string]*LocaleTranslations),
	}

	// Convert to map[string]interface{} for easier processing
	var i18nMap map[string]interface{}
	switch v := rawI18nData.(type) {
	case map[interface{}]interface{}:
		converted := convertInterfaceMapToStringMap(v)
		if convertedMap, ok := converted.(map[string]interface{}); ok {
			i18nMap = convertedMap
		}
	case map[string]interface{}:
		i18nMap = v
	default:
		return i18nConfig
	}

	// Check for structured format first (flat and translations fields)
	if flatData, hasFlat := i18nMap["flat"]; hasFlat {
		// Structured flat mappings
		if flatMap, ok := flatData.(map[string]interface{}); ok {
			for key, value := range flatMap {
				if strValue, ok := value.(string); ok {
					i18nConfig.FlatMappings[key] = strValue
				}
			}
		}
	}

	if translationsData, hasTranslations := i18nMap["translations"]; hasTranslations {
		// Structured translations
		if translationsMap, ok := translationsData.(map[string]interface{}); ok {
			for localeKey, localeValue := range translationsMap {
				if IsValidLocaleCode(localeKey) {
					if localeTranslations, ok := localeValue.(map[string]interface{}); ok {
						i18nConfig.Translations[localeKey] = &LocaleTranslations{
							Locale:       localeKey,
							Translations: localeTranslations,
						}
					}
				}
			}
		}
		return i18nConfig // Return early if structured format was found
	}

	// If no structured format, fall back to checking if this is a multi-locale configuration
	isMultiLocale := false
	for key := range i18nMap {
		if IsValidLocaleCode(key) {
			isMultiLocale = true
			break
		}
	}

	if isMultiLocale {
		// Multi-locale configuration
		for localeKey, localeValue := range i18nMap {
			if IsValidLocaleCode(localeKey) {
				if localeTranslations, ok := localeValue.(map[string]interface{}); ok {
					i18nConfig.Translations[localeKey] = &LocaleTranslations{
						Locale:       localeKey,
						Translations: localeTranslations,
					}
				}
			}
		}
	} else {
		// Simple flat mappings - check if nested or flat
		isSimpleMapping := true
		for _, value := range i18nMap {
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
			// Flat mappings
			for key, value := range i18nMap {
				if strValue, ok := value.(string); ok {
					i18nConfig.FlatMappings[key] = strValue
				}
			}
		} else {
			// Nested structure - flatten it and store as flat mappings
			flatTranslations := make(map[string]string)
			flattenI18nMapStringKeys(i18nMap, "", flatTranslations)
			i18nConfig.FlatMappings = flatTranslations
		}
	}

	return i18nConfig
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

// Helper functions for the new type-safe ConfigFile

// GetI18nMappings returns flat i18n mappings from the ConfigFile
// For compatibility with existing code that expects map[string]string
func (cf *ConfigFile) GetI18nMappings() map[string]string {
	if cf.I18n == nil {
		return make(map[string]string)
	}
	return cf.I18n.FlatMappings
}

// GetMultiLocaleI18n returns multi-locale i18n translations from the ConfigFile
// For compatibility with existing code that expects map[string]map[string]string
func (cf *ConfigFile) GetMultiLocaleI18n() map[string]map[string]string {
	if cf.I18n == nil {
		return make(map[string]map[string]string)
	}

	result := make(map[string]map[string]string)
	for locale, translations := range cf.I18n.Translations {
		flatTranslations := make(map[string]string)
		flattenI18nMapStringKeys(translations.Translations, "", flatTranslations)
		result[locale] = flatTranslations
	}
	return result
}

// GetRouteMetadata returns the metadata from the ConfigFile
// For compatibility with existing code that expects interface{}
func (cf *ConfigFile) GetRouteMetadata() interface{} {
	if cf.Metadata == nil {
		return nil
	}

	// Convert MetadataConfig to map[string]interface{} for compatibility
	result := make(map[string]interface{})

	// Add standard fields
	if cf.Metadata.Title != "" {
		result["title"] = cf.Metadata.Title
	}
	if cf.Metadata.Description != "" {
		result["description"] = cf.Metadata.Description
	}
	if len(cf.Metadata.Keywords) > 0 {
		result["keywords"] = cf.Metadata.Keywords
	}
	if cf.Metadata.Author != "" {
		result["author"] = cf.Metadata.Author
	}
	if cf.Metadata.Version != "" {
		result["version"] = cf.Metadata.Version
	}

	// Add custom fields
	for key, value := range cf.Metadata.Custom {
		result[key] = value
	}

	return result
}

// GetAuthSettings returns the auth settings from the ConfigFile
// For compatibility with existing code that expects interface{}
func (cf *ConfigFile) GetAuthSettings() interface{} {
	if cf.Auth == nil {
		return nil
	}

	// Convert AuthConfig to map[string]interface{} for compatibility
	result := make(map[string]interface{})

	if cf.Auth.Type != "" {
		result["type"] = cf.Auth.Type
	}
	if cf.Auth.RedirectURL != "" {
		result["redirect_url"] = cf.Auth.RedirectURL
	}
	if len(cf.Auth.Roles) > 0 {
		result["roles"] = cf.Auth.Roles
	}

	// Add custom settings
	for key, value := range cf.Auth.Settings {
		result[key] = value
	}

	return result
}

// GetLayoutSettings returns the layout settings from the ConfigFile
// For compatibility with existing code that expects interface{}
func (cf *ConfigFile) GetLayoutSettings() interface{} {
	if cf.Layout == nil {
		return nil
	}

	result := make(map[string]interface{})

	if cf.Layout.Template != "" {
		result["template"] = cf.Layout.Template
	}

	// Add custom settings
	for key, value := range cf.Layout.Settings {
		result[key] = value
	}

	return result
}

// GetErrorSettings returns the error settings from the ConfigFile
// For compatibility with existing code that expects interface{}
func (cf *ConfigFile) GetErrorSettings() interface{} {
	if cf.Error == nil {
		return nil
	}

	result := make(map[string]interface{})

	if cf.Error.Template != "" {
		result["template"] = cf.Error.Template
	}

	// Add custom settings
	for key, value := range cf.Error.Settings {
		result[key] = value
	}

	return result
}

// GetDynamicSettings returns the dynamic settings from the ConfigFile
// For compatibility with existing code that expects interface{}
func (cf *ConfigFile) GetDynamicSettings() interface{} {
	if cf.Dynamic == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Add validation rules
	if len(cf.Dynamic.Rules) > 0 {
		rules := make(map[string]interface{})
		for name, rule := range cf.Dynamic.Rules {
			ruleMap := make(map[string]interface{})
			ruleMap["name"] = rule.Name
			ruleMap["type"] = rule.Type
			ruleMap["required"] = rule.Required
			if rule.Pattern != "" {
				ruleMap["pattern"] = rule.Pattern
			}
			if rule.Default != nil {
				ruleMap["default"] = rule.Default
			}

			// Add custom rule settings
			for key, value := range rule.Settings {
				ruleMap[key] = value
			}

			rules[name] = ruleMap
		}
		result["rules"] = rules
	}

	// Add custom settings
	for key, value := range cf.Dynamic.Settings {
		result[key] = value
	}

	return result
}

// MergeI18n merges i18n data from another ConfigFile into this one
// The other ConfigFile's i18n data takes precedence (higher priority)
func (cf *ConfigFile) MergeI18n(other *ConfigFile) {
	if other == nil || other.I18n == nil {
		return
	}

	// Ensure this ConfigFile has i18n initialized
	if cf.I18n == nil {
		cf.I18n = &I18nConfig{
			FlatMappings: make(map[string]string),
			Translations: make(map[string]*LocaleTranslations),
		}
	}

	// Merge flat mappings - other takes precedence
	for key, value := range other.I18n.FlatMappings {
		cf.I18n.FlatMappings[key] = value
	}

	// Merge multi-locale translations - other takes precedence
	for locale, translations := range other.I18n.Translations {
		if cf.I18n.Translations[locale] == nil {
			cf.I18n.Translations[locale] = &LocaleTranslations{
				Locale:       locale,
				Translations: make(map[string]interface{}),
			}
		}

		// Merge the translation maps
		for key, value := range translations.Translations {
			cf.I18n.Translations[locale].Translations[key] = value
		}
	}
}

// MergeMetadata merges metadata from another ConfigFile into this one
// The other ConfigFile's metadata takes precedence (higher priority)
func (cf *ConfigFile) MergeMetadata(other *ConfigFile) {
	if other == nil || other.Metadata == nil {
		return
	}

	// Ensure this ConfigFile has metadata initialized
	if cf.Metadata == nil {
		cf.Metadata = &MetadataConfig{
			Custom: make(map[string]interface{}),
		}
	}

	// Merge standard fields - other takes precedence
	if other.Metadata.Title != "" {
		cf.Metadata.Title = other.Metadata.Title
	}
	if other.Metadata.Description != "" {
		cf.Metadata.Description = other.Metadata.Description
	}
	if len(other.Metadata.Keywords) > 0 {
		cf.Metadata.Keywords = other.Metadata.Keywords
	}
	if other.Metadata.Author != "" {
		cf.Metadata.Author = other.Metadata.Author
	}
	if other.Metadata.Version != "" {
		cf.Metadata.Version = other.Metadata.Version
	}

	// Merge custom fields - other takes precedence
	for key, value := range other.Metadata.Custom {
		cf.Metadata.Custom[key] = value
	}
}

// ToInterfacesConfigFile converts this shared.ConfigFile to shared.ConfigFile for backward compatibility
func (cf *ConfigFile) ToInterfacesConfigFile() map[string]interface{} {
	result := make(map[string]interface{})

	result["file_path"] = cf.FilePath
	result["template_file_path"] = cf.TemplateFilePath
	result["route_metadata"] = cf.GetRouteMetadata()
	result["i18n_mappings"] = cf.GetI18nMappings()
	result["multi_locale_i18n"] = cf.GetMultiLocaleI18n()
	result["auth_settings"] = cf.GetAuthSettings()
	result["layout_settings"] = cf.GetLayoutSettings()
	result["error_settings"] = cf.GetErrorSettings()
	result["dynamic_settings"] = cf.GetDynamicSettings()

	return result
}

// FromInterfacesConfigFile creates a shared.ConfigFile from shared.ConfigFile data
func FromInterfacesConfigFile(data map[string]interface{}) *ConfigFile {
	config := &ConfigFile{
		Metadata: &MetadataConfig{
			Custom: make(map[string]interface{}),
		},
		I18n: &I18nConfig{
			FlatMappings: make(map[string]string),
			Translations: make(map[string]*LocaleTranslations),
		},
		Auth: &AuthConfig{
			Settings: make(map[string]interface{}),
		},
		Layout: &LayoutConfig{
			Settings: make(map[string]interface{}),
		},
		Error: &ErrorConfig{
			Settings: make(map[string]interface{}),
		},
		Dynamic: &DynamicConfig{
			Rules:    make(map[string]*ValidationRule),
			Settings: make(map[string]interface{}),
		},
	}

	// Extract file paths
	if filePath, ok := data["file_path"].(string); ok {
		config.FilePath = filePath
	}
	if templatePath, ok := data["template_file_path"].(string); ok {
		config.TemplateFilePath = templatePath
	}

	// Extract route metadata
	if routeMetadata, ok := data["route_metadata"].(map[string]interface{}); ok {
		for key, value := range routeMetadata {
			config.Metadata.Custom[key] = value
		}
	}

	// Extract i18n mappings
	if i18nMappings, ok := data["i18n_mappings"].(map[string]string); ok {
		for key, value := range i18nMappings {
			config.I18n.FlatMappings[key] = value
		}
	}

	// Extract multi-locale i18n
	if multiLocaleI18n, ok := data["multi_locale_i18n"].(map[string]map[string]string); ok {
		for locale, translations := range multiLocaleI18n {
			translationsMap := make(map[string]interface{})
			for key, value := range translations {
				translationsMap[key] = value
			}
			config.I18n.Translations[locale] = &LocaleTranslations{
				Locale:       locale,
				Translations: translationsMap,
			}
		}
	}

	// Extract auth settings
	if authSettings, ok := data["auth_settings"].(map[string]interface{}); ok {
		if authType, ok := authSettings["type"].(string); ok {
			config.Auth.Type = authType
		}
		if redirectURL, ok := authSettings["redirect_url"].(string); ok {
			config.Auth.RedirectURL = redirectURL
		}
		if roles, ok := authSettings["roles"].([]string); ok {
			config.Auth.Roles = roles
		}
		// Copy other settings
		for key, value := range authSettings {
			if key != "type" && key != "redirect_url" && key != "roles" {
				config.Auth.Settings[key] = value
			}
		}
	}

	// Extract layout settings
	if layoutSettings, ok := data["layout_settings"].(map[string]interface{}); ok {
		if template, ok := layoutSettings["template"].(string); ok {
			config.Layout.Template = template
		}
		// Copy other settings
		for key, value := range layoutSettings {
			if key != "template" {
				config.Layout.Settings[key] = value
			}
		}
	}

	// Extract error settings
	if errorSettings, ok := data["error_settings"].(map[string]interface{}); ok {
		if template, ok := errorSettings["template"].(string); ok {
			config.Error.Template = template
		}
		// Copy other settings
		for key, value := range errorSettings {
			if key != "template" {
				config.Error.Settings[key] = value
			}
		}
	}

	// Extract dynamic settings
	if dynamicSettings, ok := data["dynamic_settings"].(map[string]interface{}); ok {
		if rulesData, ok := dynamicSettings["rules"].(map[string]interface{}); ok {
			for ruleName, ruleData := range rulesData {
				if ruleMap, ok := ruleData.(map[string]interface{}); ok {
					rule := &ValidationRule{
						Name:     ruleName,
						Settings: make(map[string]interface{}),
					}
					if ruleType, ok := ruleMap["type"].(string); ok {
						rule.Type = ruleType
					}
					if required, ok := ruleMap["required"].(bool); ok {
						rule.Required = required
					}
					if pattern, ok := ruleMap["pattern"].(string); ok {
						rule.Pattern = pattern
					}
					if defaultValue, ok := ruleMap["default"]; ok {
						rule.Default = defaultValue
					}
					// Copy other rule settings
					for key, value := range ruleMap {
						if key != "name" && key != "type" && key != "required" && key != "pattern" && key != "default" {
							rule.Settings[key] = value
						}
					}
					config.Dynamic.Rules[ruleName] = rule
				}
			}
		}
		// Copy other settings
		for key, value := range dynamicSettings {
			if key != "rules" {
				config.Dynamic.Settings[key] = value
			}
		}
	}

	return config
}
