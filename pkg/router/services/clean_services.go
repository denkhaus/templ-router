package services

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"go.uber.org/zap"
)

// CleanAuthService provides authentication without dependencies on router internals
type CleanAuthService struct {
	sessionStore interfaces.SessionStore
	userStore    interfaces.UserStore
	logger       *zap.Logger
}

// cleanI18nService provides internationalization without router dependencies
type cleanI18nService struct {
	configService    interfaces.ConfigService
	translationStore TranslationStore
	logger           *zap.Logger
}

// TranslationStore interface for translation management
type TranslationStore interface {
	GetTranslation(locale, key string) (string, bool)
	GetSupportedLocales() []string
	LoadTranslations(templatePath string) error
	LoadAllTranslations(templatePaths []string) error
}

// ContextAwareLayoutComponent ensures layout renders with correct context
type ContextAwareLayoutComponent struct {
	layoutComponent     templ.Component
	contextWithMetadata context.Context
	logger              *zap.Logger
}

// Render renders the layout component with the metadata context
func (calc *ContextAwareLayoutComponent) Render(ctx context.Context, w io.Writer) error {
	calc.logger.Debug("Rendering layout with metadata context")
	// Use the context that contains template_config instead of the passed context
	return calc.layoutComponent.Render(calc.contextWithMetadata, w)
}
