package i18n

import (
	"context"
	"fmt"
	"strings"
)

// I18nAdvanced provides advanced i18n functionality
// type I18nAdvanced struct {
// 	registry *I18nRegistry
// 	logger   *zap.Logger
// }

// // NewI18nAdvanced creates a new advanced i18n manager
// func NewI18nAdvanced(registry *I18nRegistry, logger *zap.Logger) *I18nAdvanced {
// 	return &I18nAdvanced{
// 		registry: registry,
// 		logger:   logger,
// 	}
// }

// TWithParams translates a key with parameter substitution
func TWithParams(ctx context.Context, key string, params map[string]string) string {
	translation := T(ctx, key)

	// Replace parameters in the format {{param}}
	for paramKey, paramValue := range params {
		placeholder := fmt.Sprintf("{{%s}}", paramKey)
		translation = strings.ReplaceAll(translation, placeholder, paramValue)
	}

	return translation
}

// TPlural handles pluralization based on count
// func TPlural(ctx context.Context, key string, count int) string {
// 	locale := GetCurrentLocale(ctx)

// 	// Get plural form based on locale rules
// 	pluralForm := getPluralForm(locale, count)
// 	pluralKey := fmt.Sprintf("%s.%s", key, pluralForm)

// 	// Try to get plural-specific translation
// 	data, ok := ctx.Value(I18nDataKey).(*I18nData)
// 	if !ok {
// 		return fmt.Sprintf("[MISSING_CONTEXT_PLURAL: %s]", key)
// 	}

// 	data.mu.RLock()
// 	defer data.mu.RUnlock()

// 	// Try plural form first
// 	if translation, exists := data.Translations[pluralKey]; exists {
// 		return strings.ReplaceAll(translation, "{{count}}", fmt.Sprint(count))
// 	}

// 	// Fallback to base key
// 	if translation, exists := data.Translations[key]; exists {
// 		return strings.ReplaceAll(translation, "{{count}}", fmt.Sprint(count))
// 	}

// 	// Graceful fallback for missing plural translations
// 	return fmt.Sprintf("[MISSING_PLURAL: %s (%d)]", key, count)
// }

// TDate formats a date according to locale
// func TDate(ctx context.Context, t time.Time, format string) string {
// 	locale := GetCurrentLocale(ctx)

// 	switch locale {
// 	case "de":
// 		return formatDateGerman(t, format)
// 	case "en":
// 		return formatDateEnglish(t, format)
// 	default:
// 		return formatDateEnglish(t, format)
// 	}
// }

// // TCurrency formats currency according to locale
// func TCurrency(ctx context.Context, amount float64) string {
// 	locale := GetCurrentLocale(ctx)

// 	switch locale {
// 	case "de":
// 		return fmt.Sprintf("%.2f €", amount)
// 	case "en":
// 		return fmt.Sprintf("$%.2f", amount)
// 	default:
// 		return fmt.Sprintf("$%.2f", amount)
// 	}
// }

// // TNumber formats numbers according to locale
// func TNumber(ctx context.Context, number float64) string {
// 	locale := GetCurrentLocale(ctx)

// 	switch locale {
// 	case "de":
// 		// German uses comma as decimal separator and dot as thousands separator
// 		return strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%.2f", number), ".", ","), ",", ".")
// 	case "en":
// 		// English uses dot as decimal separator and comma as thousands separator
// 		return fmt.Sprintf("%.2f", number)
// 	default:
// 		return fmt.Sprintf("%.2f", number)
// 	}
// }

// // ValidateTranslations checks if all required translations exist
// func (ia *I18nAdvanced) ValidateTranslations(templatePath string, requiredKeys []string) error {
// 	ia.logger.Debug("Validating translations",
// 		zap.String("template_path", templatePath),
// 		zap.Strings("required_keys", requiredKeys))

// 	// Get available translations for this template
// 	ia.registry.mu.RLock()
// 	templateTranslations, exists := ia.registry.templateTranslations[templatePath]
// 	ia.registry.mu.RUnlock()

// 	if !exists {
// 		return fmt.Errorf("no translations found for template %s", templatePath)
// 	}

// 	// Check each required key for each locale
// 	missingKeys := make(map[string][]string)

// 	for locale, translations := range templateTranslations {
// 		for _, key := range requiredKeys {
// 			if _, exists := translations[key]; !exists {
// 				missingKeys[locale] = append(missingKeys[locale], key)
// 			}
// 		}
// 	}

// 	if len(missingKeys) > 0 {
// 		return fmt.Errorf("missing translations for template %s: %v", templatePath, missingKeys)
// 	}

// 	ia.logger.Debug("All translations validated successfully",
// 		zap.String("template_path", templatePath))

// 	return nil
// }

// // GetTranslationCoverage returns translation coverage statistics
// func (ia *I18nAdvanced) GetTranslationCoverage() map[string]interface{} {
// 	ia.registry.mu.RLock()
// 	defer ia.registry.mu.RUnlock()

// 	coverage := make(map[string]interface{})
// 	totalTemplates := len(ia.registry.templateTranslations)

// 	localeStats := make(map[string]int)
// 	templateStats := make(map[string]map[string]int)

// 	for templatePath, localeTranslations := range ia.registry.templateTranslations {
// 		templateStats[templatePath] = make(map[string]int)

// 		for locale, translations := range localeTranslations {
// 			keyCount := len(translations)
// 			localeStats[locale] += keyCount
// 			templateStats[templatePath][locale] = keyCount
// 		}
// 	}

// 	coverage["total_templates"] = totalTemplates
// 	coverage["locale_stats"] = localeStats
// 	coverage["template_stats"] = templateStats

// 	return coverage
// }

// // Helper functions for pluralization
// func getPluralForm(locale string, count int) string {
// 	switch locale {
// 	case "en":
// 		if count == 1 {
// 			return "one"
// 		}
// 		return "other"
// 	case "de":
// 		if count == 1 {
// 			return "one"
// 		}
// 		return "other"
// 	default:
// 		if count == 1 {
// 			return "one"
// 		}
// 		return "other"
// 	}
// }

// // Helper functions for date formatting
// func formatDateGerman(t time.Time, format string) string {
// 	switch format {
// 	case "short":
// 		return t.Format("02.01.2006")
// 	case "long":
// 		return t.Format("2. January 2006")
// 	case "datetime":
// 		return t.Format("02.01.2006 15:04")
// 	default:
// 		return t.Format("02.01.2006")
// 	}
// }

// func formatDateEnglish(t time.Time, format string) string {
// 	switch format {
// 	case "short":
// 		return t.Format("01/02/2006")
// 	case "long":
// 		return t.Format("January 2, 2006")
// 	case "datetime":
// 		return t.Format("01/02/2006 3:04 PM")
// 	default:
// 		return t.Format("01/02/2006")
// 	}
// }

// // I18nMiddleware creates middleware for automatic i18n context setup
// func I18nMiddleware(registry *I18nRegistry) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Extract locale from URL or headers
// 			locale := extractLocaleFromRequest(r)

// 			// Determine template path from request
// 			templatePath := determineTemplateFromRequest(r)

// 			// Create i18n context
// 			ctx := registry.CreateI18nContext(r.Context(), locale, templatePath)

// 			// Continue with updated context
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// Helper to determine template from request path
// func determineTemplateFromRequest(r *http.Request) string {
// 	path := strings.TrimPrefix(r.URL.Path, "/")

// 	// Handle locale-prefixed paths
// 	parts := strings.Split(path, "/")
// 	if len(parts) > 0 && isValidLocaleCode(parts[0]) {
// 		// Remove locale from path
// 		path = strings.Join(parts[1:], "/")
// 	}

// 	// Convert path to template path
// 	if path == "" {
// 		return "app/page.templ"
// 	}

// 	templatePath := fmt.Sprintf("app/%s/page.templ", strings.ReplaceAll(path, "/", "/"))
// 	return templatePath
// }
