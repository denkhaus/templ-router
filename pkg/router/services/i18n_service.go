package services

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// NewI18nService creates a new clean i18n service for DI
func NewI18nService(i do.Injector) (interfaces.I18nService, error) {
	configService := do.MustInvoke[interfaces.ConfigService](i)
	translationStore := do.MustInvoke[TranslationStore](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &cleanI18nService{
		configService:    configService,
		translationStore: translationStore,
		logger:           logger,
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

// CreateContext implements middleware.I18nService
func (cis *cleanI18nService) CreateContext(ctx context.Context, locale, templatePath string) context.Context {
	cis.logger.Debug("Creating i18n context",
		zap.String("locale", locale),
		zap.String("template_path", templatePath))

	// Load translations for this template
	if err := cis.translationStore.LoadTranslations(templatePath); err != nil {
		cis.logger.Warn("Failed to load translations",
			zap.String("template_path", templatePath),
			zap.String("locale", locale),
			zap.Error(err))
	}

	// Create i18n data structure that router.T() expects
	i18nData := &router.I18nData{
		Locale:          locale,
		CurrentTemplate: templatePath,
		Translations:    make(map[string]string),
		FallbackLocale:  "en",
		Logger:          cis.logger,
	}

	// Load all translations for this template and locale into the context
	if store, ok := cis.translationStore.(*simpleTranslationStore); ok {
		store.mu.RLock()

		// LIBRARY-AGNOSTIC: Load layout translations first (as base)
		layoutPath := filepath.Join(cis.configService.GetLayoutRootDirectory(), cis.configService.GetLayoutFileName()+cis.configService.GetTemplateExtension())
		if layoutTranslations, exists := store.translations[layoutPath]; exists {
			if localeTranslations, exists := layoutTranslations[locale]; exists {
				for key, value := range localeTranslations {
					i18nData.Translations[key] = value
				}
				cis.logger.Debug("Loaded layout translations into context",
					zap.String("locale", locale),
					zap.String("layout_path", layoutPath),
					zap.Int("keys", len(localeTranslations)))
			} else if locale != "en" && layoutTranslations["en"] != nil {
				// Fallback to English for layout
				for key, value := range layoutTranslations["en"] {
					i18nData.Translations[key] = value
				}
				cis.logger.Debug("Loaded English layout fallback translations into context",
					zap.String("requested_locale", locale),
					zap.String("layout_path", layoutPath),
					zap.Int("keys", len(layoutTranslations["en"])))
			}
		}

		// Load template-specific translations (override layout if same key)
		if templateTranslations, exists := store.translations[templatePath]; exists {
			if localeTranslations, exists := templateTranslations[locale]; exists {
				for key, value := range localeTranslations {
					i18nData.Translations[key] = value // Template-specific overrides layout
				}
				cis.logger.Debug("Loaded template translations into context",
					zap.String("locale", locale),
					zap.String("template_path", templatePath),
					zap.Int("keys", len(localeTranslations)))
			} else if locale != "en" && templateTranslations["en"] != nil {
				// Fallback to English for template
				for key, value := range templateTranslations["en"] {
					if _, exists := i18nData.Translations[key]; !exists { // Don't override existing
						i18nData.Translations[key] = value
					}
				}
				cis.logger.Debug("Loaded English template fallback translations into context",
					zap.String("requested_locale", locale),
					zap.String("template_path", templatePath),
					zap.Int("keys", len(templateTranslations["en"])))
			}
		}

		store.mu.RUnlock()
	}

	// Set the context values that router.T() expects
	ctx = context.WithValue(ctx, router.I18nDataKey, i18nData)
	ctx = context.WithValue(ctx, router.I18nLocaleKey, locale)
	ctx = context.WithValue(ctx, router.I18nTemplateKey, templatePath)

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
