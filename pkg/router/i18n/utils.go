package i18n

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

func LocalizePath(ctx context.Context, path string) string {
	locale := GetCurrentLocale(ctx)
	return fmt.Sprintf("/%s%s", locale, path)
}

func LocalizeSafeURL(ctx context.Context, path string) templ.SafeURL {
	return templ.URL(LocalizePath(ctx, path))
}

func LocalizeRouteIfRequired(ctx context.Context, path string) string {
	// Check if path contains {locale} placeholder
	if !strings.Contains(path, "{locale}") {
		return path
	}

	// Extract locale from context
	locale := GetCurrentLocale(ctx)
	if locale != "" {
		// Replace {locale} placeholder with actual locale
		localizedPath := strings.ReplaceAll(path, "{locale}", locale)
		return localizedPath
	}

	return path
}
