package metadata

import (
	"os"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
)

// MetadataParser handles YAML metadata parsing operations
// Extracted from metadata.go for better separation of concerns
type MetadataParser struct{}

// NewMetadataParser creates a new metadata parser
func NewMetadataParser() *MetadataParser {
	return &MetadataParser{}
}

// convertSharedToRouterConfig converts shared.ConfigFile to interfaces.ConfigFile
func (mp *MetadataParser) convertSharedToRouterConfig(sharedConfig *shared.ConfigFile) *interfaces.ConfigFile {
	settingsParser := NewMetadataSettingsParser()
	return &interfaces.ConfigFile{
		FilePath:         sharedConfig.FilePath,
		TemplateFilePath: sharedConfig.TemplateFilePath,
		RouteMetadata:    sharedConfig.RouteMetadata,
		I18nMappings:     sharedConfig.I18nMappings,
		MultiLocaleI18n:  sharedConfig.MultiLocaleI18n,
		AuthSettings:     settingsParser.ParseAuthSettings(sharedConfig.AuthSettings),
		LayoutSettings:   sharedConfig.LayoutSettings,
		ErrorSettings:    sharedConfig.ErrorSettings,
		DynamicSettings:  settingsParser.ParseDynamicSettings(sharedConfig.DynamicSettings),
	}
}

// ParseYAMLMetadata parses YAML metadata files to extract route paths, auth settings, i18n mappings, and other configuration
func (mp *MetadataParser) ParseYAMLMetadata(filePath string) (*interfaces.ConfigFile, error) {
	// Use shared parser and convert to router ConfigFile
	sharedConfig, err := shared.ParseYAMLMetadata(filePath)
	if err != nil {
		return nil, err
	}

	// Convert shared.ConfigFile to router.ConfigFile
	return mp.convertSharedToRouterConfig(sharedConfig), nil
}

// ParseYAMLMetadataForTemplate reads and parses the YAML metadata file for a specific template
func (mp *MetadataParser) ParseYAMLMetadataForTemplate(templatePath string) (*interfaces.ConfigFile, error) {
	// Construct the YAML file path based on the template path
	yamlPath := templatePath + ".yaml"

	// Check if the YAML file exists
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		// Return a default config if no YAML file exists
		return &interfaces.ConfigFile{
			FilePath:         yamlPath,
			TemplateFilePath: templatePath,
			AuthSettings:     &interfaces.AuthSettings{Type: interfaces.AuthTypePublic},
		}, nil
	}

	// Parse the YAML file
	return mp.ParseYAMLMetadata(yamlPath)
}

// ParseYAMLMetadataForTemplate is the legacy global function (DEPRECATED)
// Use MetadataParser.ParseYAMLMetadataForTemplate instead
func ParseYAMLMetadataForTemplate(templatePath string) (*interfaces.ConfigFile, error) {
	parser := NewMetadataParser()
	return parser.ParseYAMLMetadataForTemplate(templatePath)
}
