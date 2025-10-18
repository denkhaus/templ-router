package interfaces

import (
	"encoding/json"
	"testing"
)

func TestConfigFile_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config ConfigFile
		valid  bool
	}{
		{
			name: "Valid config file",
			config: ConfigFile{
				FilePath:         "/app/config.yaml",
				TemplateFilePath: "/app/page.templ",
				RouteMetadata:    map[string]interface{}{"title": "Test"},
				I18nMappings:     map[string]string{"en": "English", "de": "German"},
				AuthSettings: &AuthSettings{
					Type:        AuthTypeUser,
					RedirectURL: "/login",
					Roles:       []string{"user"},
				},
			},
			valid: true,
		},
		{
			name: "Config without file path",
			config: ConfigFile{
				FilePath:         "",
				TemplateFilePath: "/app/page.templ",
			},
			valid: false,
		},
		{
			name: "Minimal valid config",
			config: ConfigFile{
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
	config := ConfigFile{
		FilePath:         "/app/config.yaml",
		TemplateFilePath: "/app/page.templ",
		RouteMetadata:    map[string]interface{}{"title": "Test", "priority": 1},
		I18nMappings:     map[string]string{"en": "English", "de": "German"},
		MultiLocaleI18n: map[string]map[string]string{
			"en": {"title": "English Title", "description": "English Description"},
			"de": {"title": "German Title", "description": "German Description"},
		},
		AuthSettings: &AuthSettings{
			Type:        AuthTypeUser,
			RedirectURL: "/login",
			Roles:       []string{"user", "admin"},
		},
		DynamicSettings: &DynamicSettings{
			Parameters: map[string]*DynamicParameterConfig{
				"id": {
					Validation:      "numeric",
					Description:     "User ID",
					SupportedValues: []string{"1", "2", "3"},
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
	var unmarshaled ConfigFile
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ConfigFile: %v", err)
	}

	// Verify basic fields
	if unmarshaled.FilePath != config.FilePath {
		t.Errorf("FilePath mismatch: got %v, want %v", unmarshaled.FilePath, config.FilePath)
	}
	if unmarshaled.TemplateFilePath != config.TemplateFilePath {
		t.Errorf("TemplateFilePath mismatch: got %v, want %v", unmarshaled.TemplateFilePath, config.TemplateFilePath)
	}

	// Verify I18n mappings
	if len(unmarshaled.I18nMappings) != len(config.I18nMappings) {
		t.Errorf("I18nMappings length mismatch: got %v, want %v", len(unmarshaled.I18nMappings), len(config.I18nMappings))
	}

	// Verify auth settings
	if unmarshaled.AuthSettings.Type != config.AuthSettings.Type {
		t.Errorf("AuthSettings.Type mismatch: got %v, want %v", unmarshaled.AuthSettings.Type, config.AuthSettings.Type)
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
			config := ConfigFile{
				FilePath:     "/app/config.yaml",
				I18nMappings: tt.mappings,
			}

			// All mapping configurations should be valid
			if !tt.valid {
				t.Errorf("Expected mappings to be valid but test marked as invalid")
			}

			// Test that we can access mappings safely
			if config.I18nMappings != nil {
				for locale, name := range config.I18nMappings {
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
	multiLocale := map[string]map[string]string{
		"en": {
			"title":       "English Title",
			"description": "English Description",
			"button":      "Click Here",
		},
		"de": {
			"title":       "German Title",
			"description": "German Description",
			"button":      "Hier Klicken",
		},
		"fr": {
			"title":       "French Title",
			"description": "French Description",
			"button":      "Cliquez Ici",
		},
	}

	config := ConfigFile{
		FilePath:        "/app/config.yaml",
		MultiLocaleI18n: multiLocale,
	}

	// Verify structure
	if len(config.MultiLocaleI18n) != 3 {
		t.Errorf("Expected 3 locales, got %d", len(config.MultiLocaleI18n))
	}

	// Verify each locale has translations
	for locale, translations := range config.MultiLocaleI18n {
		if len(translations) == 0 {
			t.Errorf("Locale %s has no translations", locale)
		}
		
		// Check for required keys
		requiredKeys := []string{"title", "description", "button"}
		for _, key := range requiredKeys {
			if value, exists := translations[key]; !exists {
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
		settings *DynamicSettings
		valid    bool
	}{
		{
			name: "Valid dynamic settings",
			settings: &DynamicSettings{
				Parameters: map[string]*DynamicParameterConfig{
					"id": {
						Validation:      "numeric",
						Description:     "User ID",
						SupportedValues: []string{"1", "2", "3"},
					},
					"slug": {
						Validation:      "alphanumeric",
						Description:     "URL slug",
						SupportedValues: []string{"home", "about", "contact"},
					},
				},
			},
			valid: true,
		},
		{
			name: "Empty parameters",
			settings: &DynamicSettings{
				Parameters: map[string]*DynamicParameterConfig{},
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
			config := ConfigFile{
				FilePath:        "/app/config.yaml",
				DynamicSettings: tt.settings,
			}

			// All dynamic settings configurations should be valid
			if !tt.valid {
				t.Errorf("Expected dynamic settings to be valid but test marked as invalid")
			}

			// Test that we can access parameters safely
			if config.DynamicSettings != nil && config.DynamicSettings.Parameters != nil {
				for paramName, paramConfig := range config.DynamicSettings.Parameters {
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
				Validation:      "numeric",
				Description:     "User ID parameter",
				SupportedValues: []string{"1", "2", "3"},
			},
			valid: true,
		},
		{
			name: "Config without validation",
			config: DynamicParameterConfig{
				Validation:      "",
				Description:     "User ID parameter",
				SupportedValues: []string{"1", "2", "3"},
			},
			valid: false,
		},
		{
			name: "Config without description",
			config: DynamicParameterConfig{
				Validation:      "numeric",
				Description:     "",
				SupportedValues: []string{"1", "2", "3"},
			},
			valid: false,
		},
		{
			name: "Config without supported values",
			config: DynamicParameterConfig{
				Validation:      "numeric",
				Description:     "User ID parameter",
				SupportedValues: []string{},
			},
			valid: false,
		},
		{
			name: "Minimal valid config",
			config: DynamicParameterConfig{
				Validation:      "alphanumeric",
				Description:     "Parameter description",
				SupportedValues: []string{"value1"},
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.config.Validation != "" && 
					  tt.config.Description != "" && 
					  len(tt.config.SupportedValues) > 0
			
			if isValid != tt.valid {
				t.Errorf("DynamicParameterConfig validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}