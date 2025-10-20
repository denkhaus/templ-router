package metadata

// MetadataI18nExtractor handles i18n extraction from YAML metadata
// Extracted from metadata.go for better separation of concerns
type MetadataI18nExtractor struct{}

// NewMetadataI18nExtractor creates a new i18n extractor
func NewMetadataI18nExtractor() *MetadataI18nExtractor {
	return &MetadataI18nExtractor{}
}
