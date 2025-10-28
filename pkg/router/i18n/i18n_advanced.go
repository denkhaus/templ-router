package i18n

import (
	"context"
	"fmt"
	"strings"
)

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
