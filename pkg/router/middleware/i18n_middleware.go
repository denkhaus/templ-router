package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// i18nMiddleware handles internationalization concerns (private implementation)
type i18nMiddleware struct {
	i18nService interfaces.I18nService
	logger      *zap.Logger
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

			// Return a proper "Language not supported" error page
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)

			// TODO: put this intemplate and make it accessible via embed.FS
			errorHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Language Not Supported - 404</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; background: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .error-code { font-size: 72px; font-weight: bold; color: #e74c3c; margin-bottom: 20px; }
        .btn { display: inline-block; padding: 12px 24px; background: #3498db; color: white; text-decoration: none; border-radius: 5px; margin: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="error-code">404</div>
        <h1>Language Not Supported</h1>
        <p>The language "<strong>` + locale + `</strong>" is not supported.</p>
        <p><strong>Supported Languages:</strong> English (en) | Deutsch (de)</p>
        <div>
            <a href="/en" class="btn">Continue in English</a>
            <a href="/de" class="btn">Auf Deutsch fortfahren</a>
            <a href="/" class="btn">Language Selection</a>
        </div>
    </div>
</body>
</html>`

			w.Write([]byte(errorHTML))
			return
		}

		im.logger.Debug("Processing i18n",
			zap.String("locale", locale),
			zap.String("template_path", templatePath),
			zap.String("path", r.URL.Path))

		// Create i18n context
		ctx := im.i18nService.CreateContext(r.Context(), locale, templatePath)

		// Add locale to context for easy access
		ctx = context.WithValue(ctx, "locale", locale)
		ctx = context.WithValue(ctx, "template_path", templatePath)

		// Continue with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetLocaleFromContext extracts locale from context
func GetLocaleFromContext(ctx context.Context) string {
	if locale, ok := ctx.Value("locale").(string); ok {
		return locale
	}
	return "en" // default fallback
}

// GetTemplatePathFromContext extracts template path from context
func GetTemplatePathFromContext(ctx context.Context) string {
	if templatePath, ok := ctx.Value("template_path").(string); ok {
		return templatePath
	}
	return ""
}
