package services

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// NewComponentMetadataService creates a new ComponentMetadataService for DI
func NewComponentMetadataService(injector do.Injector) (interfaces.ComponentMetadataService, error) {
	configService := do.MustInvoke[interfaces.ConfigService](injector)
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](injector)
	fileSystemChecker := do.MustInvoke[interfaces.FileSystemChecker](injector)
	logger := do.MustInvoke[*zap.Logger](injector)

	return &componentMetadataService{
		configService:     configService,
		templateRegistry:  templateRegistry,
		fileSystemChecker: fileSystemChecker,
		logger:            logger,
		metadataCache:     make(map[string]*shared.ConfigFile),
		translationCache:  make(map[string]map[string]string),
	}, nil
}

// componentMetadataService implements ComponentMetadataService interface
type componentMetadataService struct {
	configService     interfaces.ConfigService
	templateRegistry  interfaces.TemplateRegistry
	fileSystemChecker interfaces.FileSystemChecker
	logger            *zap.Logger
	metadataCache     map[string]*shared.ConfigFile
	translationCache  map[string]map[string]string
	metadataMutex     sync.RWMutex
	translationMutex  sync.RWMutex
}

// LoadComponentMetadata loads metadata for a specific component by name
func (cms *componentMetadataService) LoadComponentMetadata(componentName string) (*shared.ConfigFile, error) {
	// Check cache first
	if config, found := cms.GetCachedMetadata(componentName); found {
		cms.logger.Debug("Component metadata cache hit", zap.String("component", componentName))
		return config, nil
	}

	// Build component YAML path
	componentYAMLPath := cms.buildComponentYAMLPath(componentName)

	cms.logger.Debug("Loading component metadata",
		zap.String("component", componentName),
		zap.String("yaml_path", componentYAMLPath))

	// Parse YAML metadata
	_, config, err := shared.ParseYAMLMetadata(componentYAMLPath)
	if err != nil {
		cms.logger.Debug("Failed to parse component metadata",
			zap.String("component", componentName),
			zap.String("yaml_path", componentYAMLPath),
			zap.Error(err))
		return nil, fmt.Errorf("failed to load metadata for component '%s': %w", componentName, err)
	}

	// Cache the result
	cms.cacheMetadata(componentName, config)

	return config, nil
}

// GetCachedMetadata returns cached metadata if available
func (cms *componentMetadataService) GetCachedMetadata(componentName string) (*shared.ConfigFile, bool) {
	cms.metadataMutex.RLock()
	defer cms.metadataMutex.RUnlock()

	config, found := cms.metadataCache[componentName]
	return config, found
}

// LoadComponentTranslations loads i18n translations for a specific component and locale
func (cms *componentMetadataService) LoadComponentTranslations(componentName, locale string) (map[string]string, error) {
	// Check cache first
	if translations, found := cms.GetCachedTranslations(componentName, locale); found {
		cms.logger.Debug("Component translation cache hit",
			zap.String("component", componentName),
			zap.String("locale", locale))
		return translations, nil
	}

	// Load component metadata (this will also cache it)
	config, err := cms.LoadComponentMetadata(componentName)
	if err != nil {
		return nil, err
	}

	// Extract translations for the specific locale
	var translations map[string]string

	if config.MultiLocaleI18n != nil {
		if localeTranslations, found := config.MultiLocaleI18n[locale]; found {
			translations = localeTranslations
		} else {
			// Try fallback locale if configured
			if fallbackLocale := cms.getFallbackLocale(); fallbackLocale != "" {
				if fallbackTranslations, found := config.MultiLocaleI18n[fallbackLocale]; found {
					cms.logger.Debug("Using fallback locale for component translations",
						zap.String("component", componentName),
						zap.String("requested_locale", locale),
						zap.String("fallback_locale", fallbackLocale))
					translations = fallbackTranslations
				}
			}
		}
	}

	// Cache the translations (even if empty)
	cms.cacheTranslations(componentName, locale, translations)

	return translations, nil
}

// GetCachedTranslations returns cached translations if available
func (cms *componentMetadataService) GetCachedTranslations(componentName, locale string) (map[string]string, bool) {
	cms.translationMutex.RLock()
	defer cms.translationMutex.RUnlock()

	cacheKey := cms.buildTranslationCacheKey(componentName, locale)
	translations, found := cms.translationCache[cacheKey]
	return translations, found
}

// DetectComponentsFromTemplate parses a template to find component usage
func (cms *componentMetadataService) DetectComponentsFromTemplate(templateContent string) ([]string, error) {
	// Regex to find @components.ComponentName() patterns
	// This matches patterns like @components.Footer(), @components.NavBar(), etc.
	componentRegex := regexp.MustCompile(`@components\.([A-Za-z][A-Za-z0-9_]*)\s*\(`)

	matches := componentRegex.FindAllStringSubmatch(templateContent, -1)

	// Extract unique component names
	components := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			componentName := match[1]
			components[componentName] = true
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(components))
	for component := range components {
		result = append(result, component)
	}

	cms.logger.Debug("Detected components in template",
		zap.Strings("components", result))

	return result, nil
}

// LoadMultipleComponentMetadata loads metadata for multiple components efficiently
func (cms *componentMetadataService) LoadMultipleComponentMetadata(componentNames []string) (map[string]*shared.ConfigFile, error) {
	result := make(map[string]*shared.ConfigFile)

	for _, componentName := range componentNames {
		config, err := cms.LoadComponentMetadata(componentName)
		if err != nil {
			cms.logger.Warn("Failed to load component metadata",
				zap.String("component", componentName),
				zap.Error(err))
			// Continue loading other components, don't fail fast
			continue
		}
		result[componentName] = config
	}

	return result, nil
}

// Helper methods

func (cms *componentMetadataService) buildComponentYAMLPath(componentName string) string {
	// Use template registry to get component metadata directly
	if templateKey, exists := cms.templateRegistry.GetTemplateKeyByComponentName(componentName); exists {
		if metadata, err := cms.templateRegistry.GetTemplateMetadata(templateKey); err == nil {
			if metadata.YAMLExists {
				// Use the YAML file path directly from registry metadata
				yamlPath := metadata.YAMLFile

				cms.logger.Debug("Found component YAML via registry metadata",
					zap.String("component", componentName),
					zap.String("yaml_file", yamlPath),
					zap.Bool("has_i18n", metadata.HasI18n),
					zap.Bool("has_metadata", metadata.HasMetadata))

				return yamlPath
			} else {
				cms.logger.Debug("Component has no YAML file according to registry",
					zap.String("component", componentName))
			}
		}
	}

	cms.logger.Debug("Component not found in registry metadata", zap.String("component", componentName))
	return ""
}


func (cms *componentMetadataService) buildTranslationCacheKey(componentName, locale string) string {
	return fmt.Sprintf("%s:%s", componentName, locale)
}

func (cms *componentMetadataService) cacheMetadata(componentName string, config *shared.ConfigFile) {
	cms.metadataMutex.Lock()
	defer cms.metadataMutex.Unlock()

	cms.metadataCache[componentName] = config
	cms.logger.Debug("Cached component metadata", zap.String("component", componentName))
}

func (cms *componentMetadataService) cacheTranslations(componentName, locale string, translations map[string]string) {
	cms.translationMutex.Lock()
	defer cms.translationMutex.Unlock()

	cacheKey := cms.buildTranslationCacheKey(componentName, locale)
	cms.translationCache[cacheKey] = translations
	cms.logger.Debug("Cached component translations",
		zap.String("component", componentName),
		zap.String("locale", locale))
}

func (cms *componentMetadataService) getFallbackLocale() string {
	// Get fallback locale from config
	return cms.configService.GetFallbackLocale()
}