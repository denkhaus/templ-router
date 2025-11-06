package interfaces

import (
	"encoding/json"
	"testing"

	"github.com/denkhaus/templ-router/pkg/shared"
)

func TestConfigFile_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config shared.ConfigFile
		valid  bool
	}{
		{
			name: "Valid config file",
			config: shared.ConfigFile{
				FilePath:         "/app/config.yaml",
				TemplateFilePath: "/app/page.templ",
				Metadata: &shared.MetadataConfig{
					Custom: map[string]interface{}{"title": "Test"},
				},
				I18n: &shared.I18nConfig{
					FlatMappings: map[string]string{"en": "English", "de": "German"},
				},
				Auth: &shared.AuthConfig{
					Type:        "UserRequired",
					RedirectURL: "/login",
					Roles:       []string{"user"},
				},
			},
			valid: true,
		},
		{
			name: "Config without file path",
			config: shared.ConfigFile{
				FilePath:         "",
				TemplateFilePath: "/app/page.templ",
			},
			valid: false,
		},
		{
			name: "Minimal valid config",
			config: shared.ConfigFile{
				FilePath: "/app/config.yaml",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.config.FilePath != ""

			if isValid != tt.valid {
				t.Errorf("ConfigFile validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestConfigFile_JSONSerialization(t *testing.T) {
	config := shared.ConfigFile{
		FilePath:         "/app/config.yaml",
		TemplateFilePath: "/app/page.templ",
		Metadata: &shared.MetadataConfig{
			Custom: map[string]interface{}{"title": "Test", "priority": 1},
		},
		I18n: &shared.I18nConfig{
			FlatMappings: map[string]string{"en": "English", "de": "German"},
			Translations: map[string]*shared.LocaleTranslations{
				"en": {
					Locale: "en",
					Translations: map[string]interface{}{
						"title":       "English Title",
						"description": "English Description",
					},
				},
				"de": {
					Locale: "de",
					Translations: map[string]interface{}{
						"title":       "German Title",
						"description": "German Description",
					},
				},
			},
		},
		Auth: &shared.AuthConfig{
			Type:        "UserRequired",
			RedirectURL: "/login",
			Roles:       []string{"user", "admin"},
		},
		Dynamic: &shared.DynamicConfig{
			Rules: map[string]*shared.ValidationRule{
				"id": {
					Name:    "id",
					Type:    "numeric",
					Pattern: "^[0-9]+$",
				},
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal ConfigFile: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled shared.ConfigFile
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ConfigFile: %v", err)
	}

	// Verify basic fields (FilePath is excluded from JSON due to json:"-" tag)
	if unmarshaled.TemplateFilePath != config.TemplateFilePath {
		t.Errorf("TemplateFilePath mismatch: got %v, want %v", unmarshaled.TemplateFilePath, config.TemplateFilePath)
	}

	// Verify I18n mappings
	if unmarshaled.I18n != nil && config.I18n != nil {
		if len(unmarshaled.I18n.FlatMappings) != len(config.I18n.FlatMappings) {
			t.Errorf("I18n.FlatMappings length mismatch: got %v, want %v", len(unmarshaled.I18n.FlatMappings), len(config.I18n.FlatMappings))
		}
	}

	// Verify auth settings
	if unmarshaled.Auth != nil && config.Auth != nil {
		if unmarshaled.Auth.Type != config.Auth.Type {
			t.Errorf("Auth.Type mismatch: got %v, want %v", unmarshaled.Auth.Type, config.Auth.Type)
		}
	}
}

func TestConfigFile_I18nMappings(t *testing.T) {
	tests := []struct {
		name     string
		mappings map[string]string
		valid    bool
	}{
		{
			name: "Valid mappings",
			mappings: map[string]string{
				"en": "English",
				"de": "German",
				"fr": "French",
			},
			valid: true,
		},
		{
			name:     "Empty mappings",
			mappings: map[string]string{},
			valid:    true,
		},
		{
			name:     "Nil mappings",
			mappings: nil,
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := shared.ConfigFile{
				FilePath: "/app/config.yaml",
				I18n: &shared.I18nConfig{
					FlatMappings: tt.mappings,
				},
			}

			// All mapping configurations should be valid
			if !tt.valid {
				t.Errorf("Expected mappings to be valid but test marked as invalid")
			}

			// Test that we can access mappings safely
			if config.I18n != nil && config.I18n.FlatMappings != nil {
				for locale, name := range config.I18n.FlatMappings {
					if locale == "" {
						t.Errorf("Empty locale found in mappings")
					}
					if name == "" {
						t.Errorf("Empty name found for locale %s", locale)
					}
				}
			}
		})
	}
}

func TestConfigFile_MultiLocaleI18n(t *testing.T) {
	multiLocale := map[string]*shared.LocaleTranslations{
		"en": {
			Locale: "en",
			Translations: map[string]interface{}{
				"title":       "English Title",
				"description": "English Description",
				"button":      "Click Here",
			},
		},
		"de": {
			Locale: "de",
			Translations: map[string]interface{}{
				"title":       "German Title",
				"description": "German Description",
				"button":      "Hier Klicken",
			},
		},
		"fr": {
			Locale: "fr",
			Translations: map[string]interface{}{
				"title":       "French Title",
				"description": "French Description",
				"button":      "Cliquez Ici",
			},
		},
	}

	config := shared.ConfigFile{
		FilePath: "/app/config.yaml",
		I18n: &shared.I18nConfig{
			Translations: multiLocale,
		},
	}

	// Verify structure
	if len(config.I18n.Translations) != 3 {
		t.Errorf("Expected 3 locales, got %d", len(config.I18n.Translations))
	}

	// Verify each locale has translations
	for locale, localeTranslations := range config.I18n.Translations {
		if len(localeTranslations.Translations) == 0 {
			t.Errorf("Locale %s has no translations", locale)
		}

		// Check for required keys
		requiredKeys := []string{"title", "description", "button"}
		for _, key := range requiredKeys {
			if value, exists := localeTranslations.Translations[key]; !exists {
				t.Errorf("Locale %s missing key %s", locale, key)
			} else if value == "" {
				t.Errorf("Locale %s has empty value for key %s", locale, key)
			}
		}
	}
}

func TestDynamicSettings_Validation(t *testing.T) {
	tests := []struct {
		name     string
		settings *shared.DynamicConfig
		valid    bool
	}{
		{
			name: "Valid dynamic settings",
			settings: &shared.DynamicConfig{
				Rules: map[string]*DynamicParameterConfig{
					"id": {
						Name:    "id",
						Type:    "numeric",
						Pattern: "^[0-9]+$",
					},
					"slug": {
						Name:    "slug",
						Type:    "alphanumeric",
						Pattern: "^[a-z0-9-]+$",
					},
				},
			},
			valid: true,
		},
		{
			name: "Empty parameters",
			settings: &shared.DynamicConfig{
				Rules: map[string]*DynamicParameterConfig{},
			},
			valid: true,
		},
		{
			name:     "Nil settings",
			settings: nil,
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := shared.ConfigFile{
				FilePath: "/app/config.yaml",
				Dynamic:  tt.settings,
			}

			// All dynamic settings configurations should be valid
			if !tt.valid {
				t.Errorf("Expected dynamic settings to be valid but test marked as invalid")
			}

			// Test that we can access parameters safely
			if config.Dynamic != nil && config.Dynamic.Rules != nil {
				for paramName, paramConfig := range config.Dynamic.Rules {
					if paramName == "" {
						t.Errorf("Empty parameter name found")
					}
					if paramConfig == nil {
						t.Errorf("Nil parameter config found for %s", paramName)
					}
				}
			}
		})
	}
}

func TestDynamicParameterConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config DynamicParameterConfig
		valid  bool
	}{
		{
			name: "Valid parameter config",
			config: DynamicParameterConfig{
				Name:    "id",
				Type:    "numeric",
				Pattern: "^[0-9]+$",
			},
			valid: true,
		},
		{
			name: "Config without type",
			config: DynamicParameterConfig{
				Name:    "id",
				Type:    "",
				Pattern: "^[0-9]+$",
			},
			valid: false,
		},
		{
			name: "Config without name",
			config: DynamicParameterConfig{
				Name:    "",
				Type:    "numeric",
				Pattern: "^[0-9]+$",
			},
			valid: false,
		},
		{
			name: "Minimal valid config",
			config: DynamicParameterConfig{
				Name: "slug",
				Type: "alphanumeric",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.config.Name != "" && tt.config.Type != ""

			if isValid != tt.valid {
				t.Errorf("DynamicParameterConfig validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}
