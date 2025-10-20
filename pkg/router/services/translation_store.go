package services

import (
	"path/filepath"
	"sync"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

type simpleTranslationStore struct {
	configService interfaces.ConfigService
	logger        *zap.Logger
	translations  map[string]map[string]map[string]string // [templatePath][locale][key] = value
	mu            sync.RWMutex
}

// NewInMemoryTranslationStore creates a new translation store for DI
func NewInMemoryTranslationStore(i do.Injector) (TranslationStore, error) {
	configService := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)
	return &simpleTranslationStore{
		configService: configService,
		logger:        logger,
		translations:  make(map[string]map[string]map[string]string),
	}, nil
}

func (s *simpleTranslationStore) GetTranslation(locale, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Search through all templates for this key in the given locale
	for templatePath, locales := range s.translations {
		if localeTranslations, exists := locales[locale]; exists {
			if translation, found := localeTranslations[key]; found {
				s.logger.Debug("Translation found",
					zap.String("locale", locale),
					zap.String("key", key),
					zap.String("template_path", templatePath),
					zap.String("translation", translation))
				return translation, true
			}
		}
	}

	// Fallback to English if not found in requested locale
	if locale != "en" {
		for templatePath, locales := range s.translations {
			if enTranslations, exists := locales["en"]; exists {
				if translation, found := enTranslations[key]; found {
					s.logger.Debug("Translation found in English fallback",
						zap.String("requested_locale", locale),
						zap.String("key", key),
						zap.String("template_path", templatePath),
						zap.String("translation", translation))
					return translation, true
				}
			}
		}
	}

	s.logger.Debug("Translation not found",
		zap.String("locale", locale),
		zap.String("key", key))
	return "", false
}

func (s *simpleTranslationStore) GetSupportedLocales() []string {
	return s.configService.GetSupportedLocales()
}

func (s *simpleTranslationStore) LoadTranslations(templatePath string) error {
	s.logger.Debug("Loading translations for template", zap.String("template_path", templatePath))

	// LIBRARY-AGNOSTIC: Load layout translations first (if not already loaded)
	layoutPath := filepath.Join(s.configService.GetLayoutRootDirectory(), s.configService.GetLayoutFileName()+s.configService.GetTemplateExtension())
	if templatePath != layoutPath {
		s.loadTranslationsForPath(layoutPath)
	}

	// Load template-specific translations
	return s.loadTranslationsForPath(templatePath)
}

func (s *simpleTranslationStore) loadTranslationsForPath(templatePath string) error {
	s.logger.Debug("Loading translations for path", zap.String("template_path", templatePath))

	// Convert template path to YAML path
	yamlPath := templatePath + ".yaml"

	// Try to load the YAML file
	config, err := i18n.ParseYAMLMetadataExtended(yamlPath, s.logger)
	if err != nil {
		s.logger.Debug("No YAML file found or failed to parse",
			zap.String("yaml_path", yamlPath),
			zap.Error(err))
		return nil // Not an error if no YAML file exists
	}

	if config == nil {
		s.logger.Debug("No config loaded", zap.String("yaml_path", yamlPath))
		return nil
	}

	// Load translations into our store
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.translations[templatePath] == nil {
		s.translations[templatePath] = make(map[string]map[string]string)
	}

	// Check if it's multi-locale format
	if config.HasMultiLocaleSupport() {
		s.logger.Debug("Loading multi-locale translations",
			zap.String("template_path", templatePath),
			zap.Int("locales", len(config.MultiLocaleI18n)))

		for locale, translations := range config.MultiLocaleI18n {
			if s.translations[templatePath][locale] == nil {
				s.translations[templatePath][locale] = make(map[string]string)
			}
			for key, value := range translations {
				s.translations[templatePath][locale][key] = value
			}
			s.logger.Debug("Loaded translations for locale",
				zap.String("template_path", templatePath),
				zap.String("locale", locale),
				zap.Int("keys", len(translations)))
		}
	} else if len(config.ConfigFile.I18nMappings) > 0 {
		// Simple format - assume English
		s.logger.Debug("Loading simple format translations as English",
			zap.String("template_path", templatePath),
			zap.Int("keys", len(config.ConfigFile.I18nMappings)))

		if s.translations[templatePath]["en"] == nil {
			s.translations[templatePath]["en"] = make(map[string]string)
		}
		for key, value := range config.ConfigFile.I18nMappings {
			s.translations[templatePath]["en"][key] = value
		}
	}

	s.logger.Info("Successfully loaded translations",
		zap.String("template_path", templatePath),
		zap.String("yaml_path", yamlPath))

	return nil
}
