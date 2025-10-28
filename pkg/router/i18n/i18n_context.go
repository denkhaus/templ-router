package i18n

import (
	"context"
	"fmt"
	"sync"

	"github.com/denkhaus/templ-router/pkg/shared"
	"go.uber.org/zap"
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

// T translates a key using the current context
func T(ctx context.Context, key string) string {
	data, ok := ctx.Value(shared.I18nDataKey).(*I18nData)
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
	locale, ok := ctx.Value(shared.LocaleKey).(string)
	if !ok {
		// Graceful fallback to default locale
		return "en"
	}
	return locale
}

// GetCurrentTemplate returns the current template from context
func GetCurrentTemplate(ctx context.Context) string {
	template, ok := ctx.Value(shared.I18nTemplateKey).(string)
	if !ok {
		// Graceful fallback
		return "unknown"
	}
	return template
}

// GetAvailableKeys returns all available translation keys for current template
func GetAvailableKeys(ctx context.Context) []string {
	data, ok := ctx.Value(shared.I18nDataKey).(*I18nData)
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
