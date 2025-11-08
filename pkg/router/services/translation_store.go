package services

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/i18n"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

type simpleTranslationStore struct {
	configService      interfaces.ConfigService
	templateRegistry   interfaces.TemplateRegistry
	logger             *zap.Logger
	fallbackLocale     string
	layoutPath         string // Precomputed layout path
	supportedLocales   []string // Preloaded supported locales
	templateExtension  string // Preloaded template extension
	translations       map[string]map[string]map[string]string // [templatePath][locale][key] = value
	loadedPaths        map[string]bool                         // Track which template paths have been loaded
	mu                 sync.RWMutex
}

// NewInMemoryTranslationStore creates a new translation store for DI
func NewInMemoryTranslationStore(i do.Injector) (TranslationStore, error) {
	configService := do.MustInvoke[interfaces.ConfigService](i)
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	logger := do.MustInvoke[*zap.Logger](i)
	
	// Preload frequently used config values for better performance
	templateExtension := configService.GetTemplateExtension()
	layoutPath := filepath.Join(
		configService.GetLayoutRootDirectory(), 
		configService.GetLayoutFileName()+templateExtension,
	)
	
	return &simpleTranslationStore{
		configService:     configService,
		templateRegistry:  templateRegistry,
		logger:            logger,
		fallbackLocale:    configService.GetFallbackLocale(),
		layoutPath:        layoutPath,
		supportedLocales:  configService.GetSupportedLocales(),
		templateExtension: templateExtension,
		translations:      make(map[string]map[string]map[string]string),
		loadedPaths:       make(map[string]bool),
	}, nil
}

func (s *simpleTranslationStore) GetTranslation(locale, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// HIERARCHICAL TRANSLATION RESOLUTION
	// Order: Components > Pages > Layout (higher precedence wins)
	
	var foundTranslation string
	var foundInTemplate string
	foundAny := false
	
	// 1. First pass: Layout translations (lowest precedence)
	if locales, exists := s.translations[s.layoutPath]; exists {
		if localeTranslations, exists := locales[locale]; exists {
			if translation, found := localeTranslations[key]; found {
				foundTranslation = translation
				foundInTemplate = s.layoutPath
				foundAny = true
				s.logger.Debug("Translation found in layout (lowest precedence)",
					zap.String("locale", locale),
					zap.String("key", key),
					zap.String("template_path", s.layoutPath),
					zap.String("translation", translation))
			}
		}
	}
	
	// 2. Second pass: Page translations (medium precedence) - override layout
	for templatePath, locales := range s.translations {
		if templatePath == s.layoutPath || s.isComponentTemplate(templatePath) {
			continue // Skip layout and components in this pass
		}
		
		if localeTranslations, exists := locales[locale]; exists {
			if translation, found := localeTranslations[key]; found {
				foundTranslation = translation
				foundInTemplate = templatePath
				foundAny = true
				s.logger.Debug("Translation found in page (medium precedence)",
					zap.String("locale", locale),
					zap.String("key", key),
					zap.String("template_path", templatePath),
					zap.String("translation", translation))
				// Don't return yet - components might override this
			}
		}
	}
	
	// 3. Third pass: Component translations (highest precedence) - override everything
	for templatePath, locales := range s.translations {
		if !s.isComponentTemplate(templatePath) {
			continue // Only process components
		}
		
		if localeTranslations, exists := locales[locale]; exists {
			if translation, found := localeTranslations[key]; found {
				foundTranslation = translation
				foundInTemplate = templatePath
				foundAny = true
				s.logger.Debug("Translation found in component (highest precedence)",
					zap.String("locale", locale),
					zap.String("key", key),
					zap.String("template_path", templatePath),
					zap.String("translation", translation))
				// Component wins - return immediately
				return translation, true
			}
		}
	}
	
	// Return the best translation found (Page > Layout > nothing)
	if foundAny {
		s.logger.Debug("Translation resolved with hierarchy",
			zap.String("locale", locale),
			zap.String("key", key),
			zap.String("winning_template", foundInTemplate),
			zap.String("translation", foundTranslation))
		return foundTranslation, true
	}

	// Fallback to configured fallback locale if not found in requested locale
	// Apply same hierarchical search for fallback locale
	if locale != s.fallbackLocale {
		return s.getTranslationWithHierarchy(s.fallbackLocale, key)
	}

	s.logger.Debug("Translation not found in any template",
		zap.String("locale", locale),
		zap.String("key", key))
	return "", false
}

// getTranslationWithHierarchy performs hierarchical translation lookup for a specific locale
func (s *simpleTranslationStore) getTranslationWithHierarchy(locale, key string) (string, bool) {
	var foundTranslation string
	var foundInTemplate string
	foundAny := false
	
	// 1. Layout translations (lowest precedence)
	if locales, exists := s.translations[s.layoutPath]; exists {
		if localeTranslations, exists := locales[locale]; exists {
			if translation, found := localeTranslations[key]; found {
				foundTranslation = translation
				foundInTemplate = s.layoutPath
				foundAny = true
			}
		}
	}
	
	// 2. Page translations (medium precedence) - override layout
	for templatePath, locales := range s.translations {
		if templatePath == s.layoutPath || s.isComponentTemplate(templatePath) {
			continue // Skip layout and components
		}
		
		if localeTranslations, exists := locales[locale]; exists {
			if translation, found := localeTranslations[key]; found {
				foundTranslation = translation
				foundInTemplate = templatePath
				foundAny = true
			}
		}
	}
	
	// 3. Component translations (highest precedence) - override everything
	for templatePath, locales := range s.translations {
		if !s.isComponentTemplate(templatePath) {
			continue // Only components
		}
		
		if localeTranslations, exists := locales[locale]; exists {
			if translation, found := localeTranslations[key]; found {
				foundTranslation = translation
				foundInTemplate = templatePath
				foundAny = true
				// Component wins - return immediately
				s.logger.Debug("Translation found in component fallback",
					zap.String("requested_locale", "de"),
					zap.String("fallback_locale", locale),
					zap.String("key", key),
					zap.String("template_path", templatePath),
					zap.String("translation", translation))
				return translation, true
			}
		}
	}
	
	if foundAny {
		s.logger.Debug("Translation found in fallback hierarchy",
			zap.String("requested_locale", "de"),
			zap.String("fallback_locale", locale),
			zap.String("key", key),
			zap.String("winning_template", foundInTemplate),
			zap.String("translation", foundTranslation))
		return foundTranslation, true
	}
	
	return "", false
}

func (s *simpleTranslationStore) GetSupportedLocales() []string {
	return s.supportedLocales
}

func (s *simpleTranslationStore) LoadTranslations(templatePath string) error {
	s.logger.Debug("Loading translations for template", zap.String("template_path", templatePath))

	// Check if already loaded to avoid duplicate work
	s.mu.RLock()
	alreadyLoaded := s.loadedPaths[templatePath]
	s.mu.RUnlock()

	if alreadyLoaded {
		s.logger.Debug("Translations already loaded for template", zap.String("template_path", templatePath))
		return nil
	}

	// LIBRARY-AGNOSTIC: Load layout translations first (if not already loaded)
	if templatePath != s.layoutPath {
		s.mu.RLock()
		layoutLoaded := s.loadedPaths[s.layoutPath]
		s.mu.RUnlock()

		if !layoutLoaded {
			if err := s.loadTranslationsForPath(s.layoutPath); err == nil {
				s.mu.Lock()
				s.loadedPaths[s.layoutPath] = true
				s.mu.Unlock()
			}
		}
	}

	// Load template-specific translations
	err := s.loadTranslationsForPath(templatePath)
	if err == nil {
		// Mark as loaded on success
		s.mu.Lock()
		s.loadedPaths[templatePath] = true
		s.mu.Unlock()
	}
	return err
}

func (s *simpleTranslationStore) loadTranslationsForPath(templatePath string) error {
	s.logger.Debug("Loading translations for path", zap.String("template_path", templatePath))

	// Get YAML path directly from template registry
	yamlPath := s.resolveYamlPath(templatePath)
	if yamlPath == "" {
		s.logger.Debug("No YAML file for template", zap.String("template_path", templatePath))
		return nil // No YAML file exists, not an error
	}

	s.logger.Debug("Using YAML path from registry", 
		zap.String("template_path", templatePath),
		zap.String("yaml_path", yamlPath))

	// Try to load the YAML file
	configFileFound, config, err := i18n.ParseYAMLMetadataExtended(yamlPath, s.logger)
	if err != nil {

		if configFileFound {
			s.logger.Error("Failed to parse config file",
				zap.String("yaml_path", yamlPath),
				zap.Error(err),
			)

			return err
		}

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
	} else {
		// Simple format - use GetI18nMappings() to get flat mappings
		i18nMappings := config.GetI18nMappings()
		if len(i18nMappings) > 0 {
			// Simple format - assume English
			s.logger.Debug("Loading simple format translations as English",
				zap.String("template_path", templatePath),
				zap.Int("keys", len(i18nMappings)))

			if s.translations[templatePath]["en"] == nil {
				s.translations[templatePath]["en"] = make(map[string]string)
			}
			for key, value := range i18nMappings {
				s.translations[templatePath]["en"][key] = value
			}
		}
	}

	s.logger.Info("Successfully loaded translations",
		zap.String("template_path", templatePath),
		zap.String("yaml_path", yamlPath))

	return nil
}

// resolveYamlPath gets the YAML path from template registry metadata
func (s *simpleTranslationStore) resolveYamlPath(templatePath string) string {
	// First, try to get the template key from the route-to-template mapping
	routeMapping := s.templateRegistry.GetRouteToTemplateMapping()
	var templateKey string
	
	// Look for the template key by matching the path
	for _, key := range routeMapping {
		if key == templatePath {
			templateKey = key
			break
		}
	}
	
	// If not found in mapping, try using the path directly as key
	if templateKey == "" {
		templateKey = templatePath
	}
	
	// Get metadata from registry
	metadata, err := s.templateRegistry.GetTemplateMetadata(templateKey)
	if err != nil {
		s.logger.Debug("Template not found in registry, checking alternative keys", 
			zap.String("template_path", templatePath),
			zap.String("attempted_key", templateKey))
		
		// Try alternative key formats
		allMetadata := s.templateRegistry.GetAllTemplateMetadata()
		for _, meta := range allMetadata {
			if meta.TemplatePath == templatePath {
				metadata = meta
				break
			}
		}
		
		if metadata == nil {
			s.logger.Debug("Template metadata not found in registry", zap.String("template_path", templatePath))
			return "" // No YAML file
		}
	}
	
	// Check if YAML exists and return the path
	if metadata.YAMLExists && metadata.YAMLFile != "" {
		s.logger.Debug("Found YAML file in registry", 
			zap.String("template_path", templatePath),
			zap.String("yaml_file", metadata.YAMLFile))
		return metadata.YAMLFile
	}
	
	s.logger.Debug("No YAML file for template according to registry", 
		zap.String("template_path", templatePath))
	return "" // No YAML file exists
}

// isComponentTemplate checks if a template is a component by using the registry
func (s *simpleTranslationStore) isComponentTemplate(templatePath string) bool {
	allMetadata := s.templateRegistry.GetAllTemplateMetadata()
	for _, meta := range allMetadata {
		if meta.TemplatePath == templatePath && meta.Type == interfaces.TemplateTypeComponent {
			return true
		}
	}
	return false
}

// LoadAllTranslations loads translations for multiple template paths in bulk
// HIERARCHICAL SUPPORT: Also automatically discovers and loads all component translations
func (s *simpleTranslationStore) LoadAllTranslations(templatePaths []string) error {
	s.logger.Info("Starting hierarchical translation loading", zap.Int("template_count", len(templatePaths)))

	var errors []error
	loadedCount := 0
	skippedCount := 0

	// 1. First, discover and load ALL component translations (highest precedence)
	componentPaths, err := s.discoverComponentTranslations()
	if err != nil {
		s.logger.Warn("Failed to discover component translations", zap.Error(err))
	} else {
		s.logger.Info("Discovered component translations", zap.Int("component_count", len(componentPaths)))
		
		for _, componentPath := range componentPaths {
			s.mu.RLock()
			alreadyLoaded := s.loadedPaths[componentPath]
			s.mu.RUnlock()

			if alreadyLoaded {
				skippedCount++
				continue
			}

			if err := s.loadTranslationsForPath(componentPath); err != nil {
				s.logger.Warn("Failed to load component translations",
					zap.String("component_path", componentPath),
					zap.Error(err))
				errors = append(errors, err)
			} else {
				s.mu.Lock()
				s.loadedPaths[componentPath] = true
				s.mu.Unlock()
				loadedCount++
				s.logger.Debug("Loaded component translations", zap.String("component_path", componentPath))
			}
		}
	}

	// 2. Load page/template translations (medium precedence)
	for _, templatePath := range templatePaths {
		// Check if already loaded to avoid duplicate work
		s.mu.RLock()
		alreadyLoaded := s.loadedPaths[templatePath]
		s.mu.RUnlock()

		if alreadyLoaded {
			s.logger.Debug("Skipping already loaded template", zap.String("template_path", templatePath))
			skippedCount++
			continue
		}

		// Load translations for this template
		if err := s.loadTranslationsForPath(templatePath); err != nil {
			s.logger.Warn("Failed to load translations for template",
				zap.String("template_path", templatePath),
				zap.Error(err))
			errors = append(errors, err)
		} else {
			// Mark as loaded
			s.mu.Lock()
			s.loadedPaths[templatePath] = true
			s.mu.Unlock()
			loadedCount++
		}
	}

	// 3. Load layout translations (lowest precedence)
	s.mu.RLock()
	layoutLoaded := s.loadedPaths[s.layoutPath]
	s.mu.RUnlock()

	if !layoutLoaded {
		if err := s.loadTranslationsForPath(s.layoutPath); err != nil {
			s.logger.Warn("Failed to load layout translations",
				zap.String("layout_path", s.layoutPath),
				zap.Error(err))
			errors = append(errors, err)
		} else {
			s.mu.Lock()
			s.loadedPaths[s.layoutPath] = true
			s.mu.Unlock()
			loadedCount++
		}
	}

	s.logger.Info("Hierarchical translation loading completed",
		zap.Int("loaded", loadedCount),
		zap.Int("skipped", skippedCount),
		zap.Int("errors", len(errors)))

	// Return aggregated error if any failures occurred
	if len(errors) > 0 {
		return fmt.Errorf("failed to load %d translations: %v", len(errors), errors)
	}

	return nil
}

// discoverComponentTranslations finds all component templates with YAML files using the registry
func (s *simpleTranslationStore) discoverComponentTranslations() ([]string, error) {
	var componentPaths []string
	
	// Use registry to find component templates with YAML files
	allMetadata := s.templateRegistry.GetAllTemplateMetadata()
	for _, metadata := range allMetadata {
		// Check if it's a component with a YAML file
		if metadata.Type == interfaces.TemplateTypeComponent && metadata.YAMLExists {
			componentPaths = append(componentPaths, metadata.TemplatePath)
			s.logger.Debug("Discovered component with YAML", 
				zap.String("template_path", metadata.TemplatePath),
				zap.String("yaml_file", metadata.YAMLFile))
		}
	}
	
	s.logger.Debug("Component discovery completed", zap.Int("found", len(componentPaths)))
	return componentPaths, nil
}

// GetTranslationsForTemplate returns all translations for a specific template and locale
func (s *simpleTranslationStore) GetTranslationsForTemplate(templatePath, locale string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if templateTranslations, exists := s.translations[templatePath][locale]; exists {
		// Return a copy to prevent concurrent modification
		result := make(map[string]string, len(templateTranslations))
		for key, value := range templateTranslations {
			result[key] = value
		}
		return result
	}

	return make(map[string]string)
}
