package i18n

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParseYAMLMetadataExtended_NestedMultiLocale(t *testing.T) {
	// Create a test YAML file with the exact nested structure from the original example
	yamlContent := `i18n:
  en:
    feedback:
      title: "Feedback Dashboard"
      subtitle: "Overview of customer feedback and analytics"
      export: "Export Data"
      refresh: "Refresh Data"
      reviews: "reviews"
      stats:
        total_reviews: "Total Reviews"
        average_rating: "Average Rating"
        productions: "Productions"
        cache_hit_rate: "Cache Hit Rate"
      productions:
        title: "Productions"
        subtitle: "Overview of all productions with review statistics"
      recent:
        title: "Recent Reviews"
        subtitle: "Latest customer feedback and comments"
  de:
    feedback:
      title: "Feedback Dashboard (DE)"
      subtitle: "Übersicht über Kundenfeedback und Analysen"
      export: "Daten exportieren"
      refresh: "Daten aktualisieren"
      reviews: "Bewertungen"
      stats:
        total_reviews: "Gesamtbewertungen"
        average_rating: "Durchschnittsbewertung"
        productions: "Produktionen"
        cache_hit_rate: "Cache-Trefferrate"
      productions:
        title: "Produktionen"
        subtitle: "Übersicht aller Produktionen mit Bewertungsstatistiken"
      recent:
        title: "Aktuelle Bewertungen"
        subtitle: "Neuestes Kundenfeedback und Kommentare"

auth:
  type: "UserRequired"
  redirect_url: "/login"`

	tmpFile, err := os.CreateTemp("", "test_extended_nested_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Parse the YAML using ParseYAMLMetadataExtended
	logger := zap.NewNop()
	found, config, err := ParseYAMLMetadataExtended(tmpFile.Name(), logger)
	require.NoError(t, err)

	// Verify basic properties
	assert.True(t, found)
	assert.True(t, config.HasMultiLocaleSupport())
	assert.NotNil(t, config.MultiLocaleI18n)
	assert.Contains(t, config.MultiLocaleI18n, "en")
	assert.Contains(t, config.MultiLocaleI18n, "de")

	// Verify I18nMappings is empty for multi-locale configurations
	assert.Empty(t, config.I18nMappings)

	// Test English nested translations
	enTranslations := config.MultiLocaleI18n["en"]
	assert.Equal(t, "Feedback Dashboard", enTranslations["feedback.title"])
	assert.Equal(t, "Overview of customer feedback and analytics", enTranslations["feedback.subtitle"])
	assert.Equal(t, "Total Reviews", enTranslations["feedback.stats.total_reviews"])
	assert.Equal(t, "Average Rating", enTranslations["feedback.stats.average_rating"])
	assert.Equal(t, "Cache Hit Rate", enTranslations["feedback.stats.cache_hit_rate"])
	assert.Equal(t, "Productions", enTranslations["feedback.productions.title"])
	assert.Equal(t, "Recent Reviews", enTranslations["feedback.recent.title"])

	// Test German nested translations
	deTranslations := config.MultiLocaleI18n["de"]
	assert.Equal(t, "Feedback Dashboard (DE)", deTranslations["feedback.title"])
	assert.Equal(t, "Übersicht über Kundenfeedback und Analysen", deTranslations["feedback.subtitle"])
	assert.Equal(t, "Gesamtbewertungen", deTranslations["feedback.stats.total_reviews"])
	assert.Equal(t, "Durchschnittsbewertung", deTranslations["feedback.stats.average_rating"])
	assert.Equal(t, "Cache-Trefferrate", deTranslations["feedback.stats.cache_hit_rate"])
	assert.Equal(t, "Produktionen", deTranslations["feedback.productions.title"])
	assert.Equal(t, "Aktuelle Bewertungen", deTranslations["feedback.recent.title"])

	// Verify auth settings are parsed correctly
	assert.NotNil(t, config.AuthSettings)
	assert.Equal(t, "user", config.AuthSettings.Type.String())
	assert.Equal(t, "/login", config.AuthSettings.RedirectURL)
}

func TestParseYAMLMetadataExtended_SingleLocaleNested(t *testing.T) {
	// Create a test YAML file with single-locale nested structure
	yamlContent := `i18n:
  navigation:
    main:
      home: "Home"
      about: "About Us"
      services: "Services"
    user:
      profile: "My Profile"
      settings: "Account Settings"
      logout: "Sign Out"
  forms:
    buttons:
      submit: "Submit"
      cancel: "Cancel"
      save: "Save Changes"`

	tmpFile, err := os.CreateTemp("", "test_single_nested_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Parse the YAML using ParseYAMLMetadataExtended
	logger := zap.NewNop()
	found, config, err := ParseYAMLMetadataExtended(tmpFile.Name(), logger)
	require.NoError(t, err)

	// Verify basic properties
	assert.True(t, found)
	assert.False(t, config.HasMultiLocaleSupport())
	assert.Empty(t, config.MultiLocaleI18n)

	// Verify I18nMappings contains flattened nested keys
	assert.NotEmpty(t, config.I18nMappings)
	assert.Equal(t, "Home", config.I18nMappings["navigation.main.home"])
	assert.Equal(t, "About Us", config.I18nMappings["navigation.main.about"])
	assert.Equal(t, "My Profile", config.I18nMappings["navigation.user.profile"])
	assert.Equal(t, "Submit", config.I18nMappings["forms.buttons.submit"])
	assert.Equal(t, "Save Changes", config.I18nMappings["forms.buttons.save"])
}

func TestParseYAMLMetadataExtended_FlatStructure(t *testing.T) {
	// Create a test YAML file with flat structure
	yamlContent := `i18n:
  en:
    title: "Dashboard"
    subtitle: "Overview"
    save_btn: "Save"
    cancel_btn: "Cancel"
  de:
    title: "Dashboard"
    subtitle: "Übersicht"
    save_btn: "Speichern"
    cancel_btn: "Abbrechen"`

	tmpFile, err := os.CreateTemp("", "test_flat_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Parse the YAML using ParseYAMLMetadataExtended
	logger := zap.NewNop()
	found, config, err := ParseYAMLMetadataExtended(tmpFile.Name(), logger)
	require.NoError(t, err)

	// Verify basic properties
	assert.True(t, found)
	assert.True(t, config.HasMultiLocaleSupport())
	assert.NotEmpty(t, config.MultiLocaleI18n)

	// Verify flat translations are preserved
	enTranslations := config.MultiLocaleI18n["en"]
	assert.Equal(t, "Dashboard", enTranslations["title"])
	assert.Equal(t, "Overview", enTranslations["subtitle"])
	assert.Equal(t, "Save", enTranslations["save_btn"])

	deTranslations := config.MultiLocaleI18n["de"]
	assert.Equal(t, "Dashboard", deTranslations["title"])
	assert.Equal(t, "Übersicht", deTranslations["subtitle"])
	assert.Equal(t, "Speichern", deTranslations["save_btn"])
}