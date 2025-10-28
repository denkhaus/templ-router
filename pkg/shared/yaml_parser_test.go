package shared

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseYAMLMetadata_NestedI18n_MultiLocale(t *testing.T) {
	// Create a temporary YAML file with nested multi-locale i18n
	yamlContent := `i18n:
  en:
    feedback:
      title: "Feedback Dashboard"
      subtitle: "Overview of customer feedback and analytics"
      stats:
        total_reviews: "Total Reviews"
        average_rating: "Average Rating"
        cache_hit_rate: "Cache Hit Rate"
      actions:
        export: "Export Data"
        refresh: "Refresh Data"
  de:
    feedback:
      title: "Feedback Dashboard (DE)"
      subtitle: "Übersicht über Kundenfeedback und Analysen"
      stats:
        total_reviews: "Gesamtbewertungen"
        average_rating: "Durchschnittsbewertung"
        cache_hit_rate: "Cache-Trefferrate"
      actions:
        export: "Daten exportieren"
        refresh: "Daten aktualisieren"`

	tmpFile, err := os.CreateTemp("", "test_nested_i18n_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Parse the YAML
	_, config, err := ParseYAMLMetadata(tmpFile.Name())
	require.NoError(t, err)

	// Verify multi-locale i18n is populated with flattened nested keys
	assert.NotEmpty(t, config.MultiLocaleI18n)
	assert.Contains(t, config.MultiLocaleI18n, "en")
	assert.Contains(t, config.MultiLocaleI18n, "de")

	// Test English translations
	enTranslations := config.MultiLocaleI18n["en"]
	assert.Equal(t, "Feedback Dashboard", enTranslations["feedback.title"])
	assert.Equal(t, "Total Reviews", enTranslations["feedback.stats.total_reviews"])
	assert.Equal(t, "Export Data", enTranslations["feedback.actions.export"])

	// Test German translations
	deTranslations := config.MultiLocaleI18n["de"]
	assert.Equal(t, "Feedback Dashboard (DE)", deTranslations["feedback.title"])
	assert.Equal(t, "Gesamtbewertungen", deTranslations["feedback.stats.total_reviews"])
	assert.Equal(t, "Daten exportieren", deTranslations["feedback.actions.export"])

	// Verify I18nMappings is empty for multi-locale configurations
	assert.Empty(t, config.I18nMappings)
}

func TestParseYAMLMetadata_NestedI18n_SimpleStructure(t *testing.T) {
	// Create a temporary YAML file with nested simple i18n (non-multi-locale)
	yamlContent := `i18n:
  feedback:
    title: "Simple Feedback Dashboard"
    subtitle: "Overview of feedback"
    stats:
      total_reviews: "Total Reviews"
      average_rating: "Average Rating"
    actions:
      export: "Export Data"
      refresh: "Refresh Data"`

	tmpFile, err := os.CreateTemp("", "test_simple_nested_i18n_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Parse the YAML
	_, config, err := ParseYAMLMetadata(tmpFile.Name())
	require.NoError(t, err)

	// Verify I18nMappings is populated with flattened nested keys
	assert.NotEmpty(t, config.I18nMappings)
	assert.Equal(t, "Simple Feedback Dashboard", config.I18nMappings["feedback.title"])
	assert.Equal(t, "Total Reviews", config.I18nMappings["feedback.stats.total_reviews"])
	assert.Equal(t, "Export Data", config.I18nMappings["feedback.actions.export"])

	// Verify MultiLocaleI18n is empty for simple structures
	assert.Empty(t, config.MultiLocaleI18n)
}

func TestParseYAMLMetadata_FlatI18n(t *testing.T) {
	// Create a temporary YAML file with flat i18n
	yamlContent := `i18n:
  title: "Dashboard"
  subtitle: "Overview"
  export: "Export"`

	tmpFile, err := os.CreateTemp("", "test_flat_i18n_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Parse the YAML
	_, config, err := ParseYAMLMetadata(tmpFile.Name())
	require.NoError(t, err)

	// Verify I18nMappings is populated
	assert.NotEmpty(t, config.I18nMappings)
	assert.Equal(t, "Dashboard", config.I18nMappings["title"])
	assert.Equal(t, "Overview", config.I18nMappings["subtitle"])
	assert.Equal(t, "Export", config.I18nMappings["export"])

	// Verify MultiLocaleI18n is empty
	assert.Empty(t, config.MultiLocaleI18n)
}

func TestIsValidLocaleCode(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		{"en", true},
		{"de", true},
		{"fr", true},
		{"zh", true},
		{"en-US", true},
		{"de-DE", true},
		{"feedback", false},
		{"dashboard", false},
		{"stats", false},
		{"actions", false},
		{"", false},
		{"xyz", false},
		{"english", false},
	}

	for _, test := range tests {
		t.Run(test.code, func(t *testing.T) {
			result := IsValidLocaleCode(test.code)
			assert.Equal(t, test.expected, result, "IsValidLocaleCode(%q) should return %v", test.code, test.expected)
		})
	}
}
