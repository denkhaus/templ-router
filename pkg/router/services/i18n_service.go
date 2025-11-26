package services

import (
	"context"
	"net/http"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// NewI18nService creates a new clean i18n service for DI
func NewI18nService(i do.Injector) (interfaces.I18nService, error) {
	configService := do.MustInvoke[interfaces.ConfigService](i)
	translationStore := do.MustInvoke[TranslationStore](i)
	componentMetadataService := do.MustInvoke[interfaces.ComponentMetadataService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &cleanI18nService{
		configService:            configService,
		translationStore:         translationStore,
		componentMetadataService: componentMetadataService,
		logger:                   logger,
	}, nil
}

// ExtractLocale implements middleware.I18nService
func (cis *cleanI18nService) ExtractLocale(req *http.Request) string {
	// Try to extract from URL path first (e.g., /en/dashboard, /de/admin, /en, /de)
	path := req.URL.Path

	// Handle root locale paths like /en or /de
	if len(path) == 3 && path[0] == '/' {
		locale := path[1:3]
		if cis.isValidLocale(locale) {
			cis.logger.Debug("Extracted locale from root path",
				zap.String("path", path),
				zap.String("locale", locale))
			return locale
		}
	}

	// Handle nested locale paths like /en/dashboard or /de/admin
	if len(path) > 3 && path[0] == '/' && path[3] == '/' {
		locale := path[1:3]
		if cis.isValidLocale(locale) {
			cis.logger.Debug("Extracted locale from nested path",
				zap.String("path", path),
				zap.String("locale", locale))
			return locale
		}
	}

	// Try to extract from Accept-Language header
	if acceptLang := req.Header.Get("Accept-Language"); acceptLang != "" {
		// Simple parsing - take first 2 characters
		if len(acceptLang) >= 2 {
			locale := acceptLang[:2]
			if cis.isValidLocale(locale) {
				return locale
			}
		}
	}

	// Use configured default locale
	return cis.configService.GetDefaultLocale()
}

// CreateContext implements middleware.I18nService - creates i18n context for ANY template
func (cis *cleanI18nService) CreateContext(ctx context.Context, templatePath string) context.Context {
	// Extract locale from context (set by middleware)
	locale, ok := ctx.Value(shared.LocaleKey).(string)
	if !ok {
		locale = cis.configService.GetDefaultLocale()
		cis.logger.Warn("No locale found in context, using default",
			zap.String("default_locale", locale),
			zap.String("template_path", templatePath))
	}

	cis.logger.Debug("Creating i18n context",
		zap.String("locale", locale),
		zap.String("template_path", templatePath))

	// Load translations for this template using TranslationStore
	if err := cis.translationStore.LoadTranslations(templatePath); err != nil {
		cis.logger.Debug("Failed to load translations",
			zap.String("template_path", templatePath),
			zap.String("locale", locale),
			zap.Error(err))
	}

	// Extract translations for current template and locale from TranslationStore
	translations := cis.translationStore.GetTranslationsForTemplate(templatePath, locale)

	// Create i18n data structure that i18n.T() expects
	i18nData := &i18n.I18nData{
		Locale:          locale,
		CurrentTemplate: templatePath,
		Translations:    translations, // Now populated with actual translations!
		FallbackLocale:  cis.configService.GetDefaultLocale(),
		Logger:          cis.logger,
	}

	// Set the context values that i18n.T() expects
	ctx = context.WithValue(ctx, shared.I18nDataKey, i18nData)
	ctx = context.WithValue(ctx, shared.I18nTemplateKey, templatePath)

	return ctx
}

// LoadComponentTranslationsIntoContext loads translations for a specific component into the i18n context
// This method can be called by the template middleware when processing component routes
func (cis *cleanI18nService) LoadComponentTranslationsIntoContext(ctx context.Context, componentName string) context.Context {
	// Get current locale from context
	locale, ok := ctx.Value(shared.LocaleKey).(string)
	if !ok {
		cis.logger.Debug("No locale found in context, skipping component translation loading",
			zap.String("component", componentName))
		return ctx
	}

	// Get the i18n data from context
	i18nData, ok := ctx.Value(shared.I18nDataKey).(*i18n.I18nData)
	if !ok {
		cis.logger.Debug("No i18n data found in context, skipping component translation loading",
			zap.String("component", componentName),
			zap.String("locale", locale))
		return ctx
	}
	if cis.componentMetadataService == nil {
		cis.logger.Debug("ComponentMetadataService not available, skipping component translation loading",
			zap.String("component", componentName),
			zap.String("locale", locale))
		return ctx
	}

	// Load component translations using the ComponentMetadataService
	componentTranslations, err := cis.componentMetadataService.LoadComponentTranslations(componentName, locale)
	if err != nil {
		cis.logger.Debug("Failed to load component translations",
			zap.String("component", componentName),
			zap.String("locale", locale),
			zap.Error(err))
		return ctx
	}

	// Merge component translations into i18n context
	// Component translations take precedence over existing translations
	for key, value := range componentTranslations {
		i18nData.Translations[key] = value
	}

	cis.logger.Debug("Loaded component translations into context",
		zap.String("component", componentName),
		zap.String("locale", locale),
		zap.Int("keys", len(componentTranslations)))

	return ctx
}

// GetSupportedLocales returns supported locales from config
func (cis *cleanI18nService) GetSupportedLocales() []string {
	return cis.configService.GetSupportedLocales()
}

// isValidLocale checks if a locale is supported using config
func (cis *cleanI18nService) isValidLocale(locale string) bool {
	supportedLocales := cis.configService.GetSupportedLocales()
	for _, supported := range supportedLocales {
		if supported == locale {
			return true
		}
	}
	return false
}

// LoadAllTranslations implements interfaces.I18nService by delegating to TranslationStore
func (cis *cleanI18nService) LoadAllTranslations(templatePaths []string) error {
	return cis.translationStore.LoadAllTranslations(templatePaths)
}
