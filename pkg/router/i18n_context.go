package router

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// I18nContextKey is the context key for i18n data
type I18nContextKey string

const (
	I18nDataKey     I18nContextKey = "router_i18n_data"
	I18nLocaleKey   I18nContextKey = "router_i18n_locale"
	I18nTemplateKey I18nContextKey = "router_i18n_template"
)

// I18nData holds translation data for the current request
type I18nData struct {
	Locale          string
	CurrentTemplate string
	Translations    map[string]string // key -> translation
	FallbackLocale  string
	Logger          *zap.Logger
	mu              sync.RWMutex
}

// I18nRegistry manages all template translations
type I18nRegistry struct {
	// templateTranslations[templatePath][locale][key] = translation
	templateTranslations map[string]map[string]map[string]string
	// requiredKeys[templatePath] = set of required keys for that template
	requiredKeys map[string]map[string]bool
	logger       *zap.Logger
	mu           sync.RWMutex
}

// NewI18nRegistry creates a new i18n registry
func NewI18nRegistry(logger *zap.Logger) *I18nRegistry {
	if logger == nil {
		panic("logger is required for I18nRegistry - no fallback to no-op logger")
	}

	return &I18nRegistry{
		templateTranslations: make(map[string]map[string]map[string]string),
		requiredKeys:         make(map[string]map[string]bool),
		logger:               logger,
	}
}

// RegisterTemplateTranslations registers translations for a template
func (reg *I18nRegistry) RegisterTemplateTranslations(templatePath string, config *ConfigFile) error {
	if templatePath == "" {
		return fmt.Errorf("template path cannot be empty")
	}

	reg.mu.Lock()
	defer reg.mu.Unlock()

	if reg.templateTranslations[templatePath] == nil {
		reg.templateTranslations[templatePath] = make(map[string]map[string]string)
	}

	// Process different i18n formats from YAML
	if config != nil {
		if err := reg.processYAMLTranslations(templatePath, config); err != nil {
			return fmt.Errorf("failed to process YAML translations for %s: %w", templatePath, err)
		}
	}

	reg.logger.Debug("Registered translations for template",
		zap.String("template_path", templatePath),
		zap.Int("locales_count", len(reg.templateTranslations[templatePath])))

	return nil
}

// processYAMLTranslations processes different YAML i18n formats
func (reg *I18nRegistry) processYAMLTranslations(templatePath string, config *ConfigFile) error {
	// Check if we have multi-locale format (preferred)
	if len(config.MultiLocaleI18n) > 0 {
		return reg.processMultiLocaleFormat(templatePath, config)
	}

	// Check if we have simple key-value format (fallback)
	if len(config.I18nMappings) > 0 {
		return reg.processSimpleFormat(templatePath, config)
	}

	return nil
}

// processSimpleFormat processes simple key-value YAML
// Format: i18n: { key: "value", key2: "value2" }
func (reg *I18nRegistry) processSimpleFormat(templatePath string, config *ConfigFile) error {
	// Default to English for simple format
	defaultLocale := "en"

	if reg.templateTranslations[templatePath][defaultLocale] == nil {
		reg.templateTranslations[templatePath][defaultLocale] = make(map[string]string)
	}

	for key, value := range config.I18nMappings {
		reg.templateTranslations[templatePath][defaultLocale][key] = value
	}

	reg.logger.Debug("Processed simple format translations",
		zap.String("template_path", templatePath),
		zap.String("locale", defaultLocale),
		zap.Int("keys_count", len(config.I18nMappings)))

	return nil
}

// RegisterRequiredKey registers a key as required for a template (called during template parsing)
func (reg *I18nRegistry) RegisterRequiredKey(templatePath, key string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if reg.requiredKeys[templatePath] == nil {
		reg.requiredKeys[templatePath] = make(map[string]bool)
	}

	reg.requiredKeys[templatePath][key] = true

	reg.logger.Debug("Registered required translation key",
		zap.String("template_path", templatePath),
		zap.String("key", key))
}

// ValidateAllTranslations validates that all required keys have translations for all supported locales
func (reg *I18nRegistry) ValidateAllTranslations(supportedLocales []string) error {
	reg.mu.RLock()
	defer reg.mu.RUnlock()

	var errors []string

	for templatePath, requiredKeys := range reg.requiredKeys {
		templateTranslations, hasTemplate := reg.templateTranslations[templatePath]
		if !hasTemplate {
			errors = append(errors, fmt.Sprintf("MISSING: Template '%s' has no translations defined", templatePath))
			continue
		}

		for key := range requiredKeys {
			for _, locale := range supportedLocales {
				localeTranslations, hasLocale := templateTranslations[locale]
				if !hasLocale {
					errors = append(errors, fmt.Sprintf("MISSING: Template '%s': missing locale '%s'", templatePath, locale))
					continue
				}

				if _, hasKey := localeTranslations[key]; !hasKey {
					errors = append(errors, fmt.Sprintf("MISSING: Template '%s': missing key '%s' for locale '%s'", templatePath, key, locale))
				}
			}
		}
	}

	if len(errors) > 0 {
		errorMsg := fmt.Sprintf("\n=== Translation Validation Failed ===\n\nMissing translations found:\n%s\n\nFix: Add the missing keys to the corresponding .templ.yaml files\n",
			strings.Join(errors, "\n"))
		return fmt.Errorf("%s", errorMsg)
	}

	reg.logger.Info("All translations validated successfully",
		zap.Int("templates", len(reg.requiredKeys)),
		zap.Int("locales", len(supportedLocales)))

	return nil
}

// processMultiLocaleFormat processes YAML with multiple locales
// Format: i18n: { en: { key: "value" }, de: { key: "wert" } }
func (reg *I18nRegistry) processMultiLocaleFormat(templatePath string, config *ConfigFile) error {
	for locale, translations := range config.MultiLocaleI18n {
		if reg.templateTranslations[templatePath][locale] == nil {
			reg.templateTranslations[templatePath][locale] = make(map[string]string)
		}

		for key, value := range translations {
			reg.templateTranslations[templatePath][locale][key] = value
		}

		reg.logger.Debug("Processed multi-locale format translations",
			zap.String("template_path", templatePath),
			zap.String("locale", locale),
			zap.Int("keys_count", len(translations)))
	}

	return nil
}

// CreateI18nContext creates an i18n context for a request
// func (reg *I18nRegistry) CreateI18nContext(ctx context.Context, locale, templatePath string) context.Context {
// 	reg.mu.RLock()
// 	defer reg.mu.RUnlock()

// 	// Get translations for this template and locale
// 	translations := make(map[string]string)

// 	// LIBRARY-AGNOSTIC: Load layout translations first (as base)
// 	// FAIL FAST: No hardcoded paths allowed - Config injection required
// 	panic("i18n_context.go: Config injection required - hardcoded layout paths forbidden")

// 	// This function needs to be refactored to receive config via DI
// 	// For now, we'll skip layout loading to prevent build errors
// 	layoutPaths := []string{} // Empty to prevent undefined variable error
// 	for _, layoutPath := range layoutPaths {
// 		if templateTranslations, exists := reg.templateTranslations[layoutPath]; exists {
// 			reg.logger.Debug("Found layout translations", zap.String("layout_path", layoutPath))
// 			if localeTranslations, exists := templateTranslations[locale]; exists {
// 				for k, v := range localeTranslations {
// 					translations[k] = v
// 				}
// 			}
// 			// Fallback to English for layout
// 			if locale != "en" {
// 				if enTranslations, exists := templateTranslations["en"]; exists {
// 					for k, v := range enTranslations {
// 						if _, exists := translations[k]; !exists { // Don't override locale-specific
// 							translations[k] = v
// 						}
// 					}
// 				}
// 			}
// 			break // Use first found layout
// 		}
// 	}

// 	// Load template-specific translations (override layout if same key)
// 	if templateTranslations, exists := reg.templateTranslations[templatePath]; exists {
// 		if localeTranslations, exists := templateTranslations[locale]; exists {
// 			for k, v := range localeTranslations {
// 				translations[k] = v // Template-specific overrides layout
// 			}
// 		}

// 		// Fallback to English if locale not found
// 		if locale != "en" {
// 			if enTranslations, exists := templateTranslations["en"]; exists {
// 				for k, v := range enTranslations {
// 					if _, exists := translations[k]; !exists { // Don't override existing
// 						translations[k] = v
// 					}
// 				}
// 			}
// 		}
// 	}

// 	i18nData := &I18nData{
// 		Locale:          locale,
// 		CurrentTemplate: templatePath,
// 		Translations:    translations,
// 		FallbackLocale:  "en",
// 		Logger:          reg.logger,
// 	}

// 	// Add to context
// 	ctx = context.WithValue(ctx, I18nDataKey, i18nData)
// 	ctx = context.WithValue(ctx, I18nLocaleKey, locale)
// 	ctx = context.WithValue(ctx, I18nTemplateKey, templatePath)

// 	return ctx
// }

// T translates a key using the current context
func T(ctx context.Context, key string) string {
	data, ok := ctx.Value(I18nDataKey).(*I18nData)
	if !ok {
		// Graceful fallback when i18n context is missing
		return fmt.Sprintf("[MISSING_I18N_CONTEXT: %s]", key)
	}

	data.mu.RLock()
	defer data.mu.RUnlock()

	if translation, exists := data.Translations[key]; exists {
		data.Logger.Debug("Translation found",
			zap.String("key", key),
			zap.String("locale", data.Locale),
			zap.String("template", data.CurrentTemplate),
			zap.String("translation", translation))
		return translation
	}

	// Graceful fallback for missing translations
	data.Logger.Warn("Translation key not found - using fallback",
		zap.String("key", key),
		zap.String("locale", data.Locale),
		zap.String("template", data.CurrentTemplate),
		zap.String("fallback", fmt.Sprintf("[MISSING: %s]", key)))

	// Return a visible but non-breaking fallback
	return fmt.Sprintf("[MISSING_I18N: %s]", key)
}

// GetCurrentLocale returns the current locale from context
func GetCurrentLocale(ctx context.Context) string {
	locale, ok := ctx.Value(I18nLocaleKey).(string)
	if !ok {
		// Graceful fallback to default locale
		return "en"
	}
	return locale
}

// GetCurrentTemplate returns the current template from context
func GetCurrentTemplate(ctx context.Context) string {
	template, ok := ctx.Value(I18nTemplateKey).(string)
	if !ok {
		// Graceful fallback
		return "unknown"
	}
	return template
}

// GetAvailableKeys returns all available translation keys for current template
func GetAvailableKeys(ctx context.Context) []string {
	data, ok := ctx.Value(I18nDataKey).(*I18nData)
	if !ok {
		return nil
	}

	data.mu.RLock()
	defer data.mu.RUnlock()

	keys := make([]string, 0, len(data.Translations))
	for key := range data.Translations {
		keys = append(keys, key)
	}

	return keys
}
