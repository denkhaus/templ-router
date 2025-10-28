package middleware

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

//go:embed templates/*.html
var i18nTemplates embed.FS

// i18nMiddleware handles internationalization concerns (private implementation)
type i18nMiddleware struct {
	i18nService interfaces.I18nService
	logger      *zap.Logger
}

// LanguageNotSupportedData holds data for the language not supported template
type LanguageNotSupportedData struct {
	Locale             string
	SupportedLanguages []SupportedLanguage
}

// SupportedLanguage represents a supported language with display information
type SupportedLanguage struct {
	Code         string
	Name         string
	ContinueText string
}

// I18nService interface for clean dependency

// NewI18nMiddleware creates a new i18n middleware for DI
func NewI18nMiddleware(i do.Injector) (I18nMiddlewareInterface, error) {
	i18nService := do.MustInvoke[interfaces.I18nService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &i18nMiddleware{
		i18nService: i18nService,
		logger:      logger,
	}, nil
}

// isLocaleInPath checks if the URL path contains a locale segment
func (im *i18nMiddleware) isLocaleInPath(path string) bool {
	// Check if path starts with /xx or /xx/ where xx could be a locale
	parts := strings.Split(path, "/")
	return len(parts) >= 2 && len(parts[1]) == 2
}

// isValidLocale checks if a locale is supported using the translation store
func (im *i18nMiddleware) isValidLocale(locale string) bool {
	supportedLocales := im.i18nService.GetSupportedLocales()
	for _, supported := range supportedLocales {
		if locale == supported {
			return true
		}
	}
	return false
}

// renderLanguageNotSupportedPage renders the language not supported error page using embedded template
func (im *i18nMiddleware) renderLanguageNotSupportedPage(w http.ResponseWriter, locale string) {
	// Parse embedded template
	tmpl, err := template.ParseFS(i18nTemplates, "templates/i18n_language_not_supported.html")
	if err != nil {
		im.logger.Error("Failed to parse language not supported template", zap.Error(err))
		// Fallback to simple error message
		http.Error(w, "Language not supported", http.StatusNotFound)
		return
	}

	// Get supported languages dynamically from i18n service
	supportedLocales := im.i18nService.GetSupportedLocales()
	supportedLanguages := make([]SupportedLanguage, 0, len(supportedLocales))

	// Map locale codes to language names and continue text (ASCII only)
	languageMap := map[string]SupportedLanguage{
		"en": {Code: "en", Name: "English", ContinueText: "Continue in English"},
		"de": {Code: "de", Name: "Deutsch", ContinueText: "Auf Deutsch fortfahren"},
		"fr": {Code: "fr", Name: "Francais", ContinueText: "Continuer en francais"},
		"es": {Code: "es", Name: "Espanol", ContinueText: "Continuar en espanol"},
		"it": {Code: "it", Name: "Italiano", ContinueText: "Continua in italiano"},
		"pt": {Code: "pt", Name: "Portugues", ContinueText: "Continuar em portugues"},
		"nl": {Code: "nl", Name: "Nederlands", ContinueText: "Doorgaan in het Nederlands"},
		"ru": {Code: "ru", Name: "Russian", ContinueText: "Continue in Russian"},
		"zh": {Code: "zh", Name: "Chinese", ContinueText: "Continue in Chinese"},
		"ja": {Code: "ja", Name: "Japanese", ContinueText: "Continue in Japanese"},
	}

	// Build supported languages list from actual supported locales
	for _, localeCode := range supportedLocales {
		if lang, exists := languageMap[localeCode]; exists {
			supportedLanguages = append(supportedLanguages, lang)
		} else {
			// Fallback for unknown locale codes
			supportedLanguages = append(supportedLanguages, SupportedLanguage{
				Code:         localeCode,
				Name:         strings.ToUpper(localeCode),
				ContinueText: "Continue in " + strings.ToUpper(localeCode),
			})
		}
	}

	data := LanguageNotSupportedData{
		Locale:             locale,
		SupportedLanguages: supportedLanguages,
	}

	// Set headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	// Execute template
	if err := tmpl.Execute(w, data); err != nil {
		im.logger.Error("Failed to execute language not supported template",
			zap.Error(err),
			zap.String("locale", locale))
		// Template execution failed, write simple fallback
		w.Write([]byte("Language not supported"))
	}

	im.logger.Debug("Rendered language not supported page",
		zap.String("unsupported_locale", locale),
		zap.Strings("supported_locales", supportedLocales))
}

// Handle processes internationalization for a request
func (im *i18nMiddleware) Handle(next http.Handler, templatePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract locale from request
		locale := im.i18nService.ExtractLocale(r)

		// LOCALE VALIDATION: Check if extracted locale is valid
		if im.isLocaleInPath(r.URL.Path) && !im.isValidLocale(locale) {
			im.logger.Info("Unsupported locale detected",
				zap.String("path", r.URL.Path),
				zap.String("unsupported_locale", locale),
				zap.Strings("supported_locales", im.i18nService.GetSupportedLocales()))

			// Render language not supported page using embedded template
			im.renderLanguageNotSupportedPage(w, locale)
			return
		}

		im.logger.Debug("Processing i18n",
			zap.String("locale", locale),
			zap.String("template_path", templatePath),
			zap.String("path", r.URL.Path))

		// Add locale to context FIRST (before service call)
		ctx := context.WithValue(r.Context(), shared.LocaleKey, locale)
		ctx = context.WithValue(ctx, shared.TemplatePathKey, templatePath)

		// Create i18n context (service will read locale from context)
		ctx = im.i18nService.CreateContext(ctx, templatePath)

		// Continue with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
