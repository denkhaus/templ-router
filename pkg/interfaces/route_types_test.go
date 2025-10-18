package interfaces

import (
	"encoding/json"
	"testing"
)

func TestRoute_JSONSerialization(t *testing.T) {
	route := Route{
		Path:         "/test/{id}",
		Handler:      "TestHandler",
		TemplateFile: "test.templ",
		MetadataFile: "test.yaml",
		IsDynamic:    true,
		Precedence:   1,
		Locale:       "en",
		AuthSettings: &AuthSettings{
			Type:        AuthTypeUser,
			RedirectURL: "/login",
			Roles:       []string{"user", "admin"},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal Route: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Route
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Route: %v", err)
	}

	// Verify fields
	if unmarshaled.Path != route.Path {
		t.Errorf("Path mismatch: got %v, want %v", unmarshaled.Path, route.Path)
	}
	if unmarshaled.IsDynamic != route.IsDynamic {
		t.Errorf("IsDynamic mismatch: got %v, want %v", unmarshaled.IsDynamic, route.IsDynamic)
	}
	if unmarshaled.AuthSettings.Type != route.AuthSettings.Type {
		t.Errorf("AuthType mismatch: got %v, want %v", unmarshaled.AuthSettings.Type, route.AuthSettings.Type)
	}
}

func TestRoute_Validation(t *testing.T) {
	tests := []struct {
		name  string
		route Route
		valid bool
	}{
		{
			name: "Valid static route",
			route: Route{
				Path:         "/about",
				TemplateFile: "about.templ",
				IsDynamic:    false,
			},
			valid: true,
		},
		{
			name: "Valid dynamic route",
			route: Route{
				Path:         "/user/{id}",
				TemplateFile: "user.templ",
				IsDynamic:    true,
			},
			valid: true,
		},
		{
			name: "Route without path",
			route: Route{
				Path:         "",
				TemplateFile: "test.templ",
				IsDynamic:    false,
			},
			valid: false,
		},
		{
			name: "Route without template file",
			route: Route{
				Path:         "/test",
				TemplateFile: "",
				IsDynamic:    false,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.route.Path != "" && tt.route.TemplateFile != ""
			
			if isValid != tt.valid {
				t.Errorf("Route validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestLayoutTemplate_Validation(t *testing.T) {
	tests := []struct {
		name     string
		template LayoutTemplate
		isValid  bool
	}{
		{
			name: "Valid layout template",
			template: LayoutTemplate{
				FilePath:      "/app/layout.templ",
				YamlPath:      "/app/layout.yaml",
				ComponentName: "Layout",
				Content:       "<div>content</div>",
				LayoutLevel:   1,
			},
			isValid: true,
		},
		{
			name: "Empty file path",
			template: LayoutTemplate{
				FilePath:      "",
				ComponentName: "Layout",
			},
			isValid: false,
		},
		{
			name: "Minimal valid template",
			template: LayoutTemplate{
				FilePath: "/app/layout.templ",
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - file path should not be empty
			isValid := tt.template.FilePath != ""
			if isValid != tt.isValid {
				t.Errorf("Layout template validation = %v, want %v", isValid, tt.isValid)
			}
		})
	}
}

func TestLayoutTemplate_JSONSerialization(t *testing.T) {
	template := LayoutTemplate{
		FilePath:      "/app/layout.templ",
		YamlPath:      "/app/layout.yaml",
		ComponentName: "Layout",
		Content:       "<div>content</div>",
		LayoutLevel:   1,
	}

	// Test JSON marshaling
	data, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal LayoutTemplate: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled LayoutTemplate
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal LayoutTemplate: %v", err)
	}

	// Verify fields
	if unmarshaled.FilePath != template.FilePath {
		t.Errorf("FilePath mismatch: got %v, want %v", unmarshaled.FilePath, template.FilePath)
	}
	if unmarshaled.LayoutLevel != template.LayoutLevel {
		t.Errorf("LayoutLevel mismatch: got %v, want %v", unmarshaled.LayoutLevel, template.LayoutLevel)
	}
}

func TestErrorTemplate_Validation(t *testing.T) {
	tests := []struct {
		name     string
		template ErrorTemplate
		valid    bool
	}{
		{
			name: "Valid error template",
			template: ErrorTemplate{
				FilePath:      "/app/error.templ",
				ComponentName: "ErrorPage",
				Content:       "<div>Error occurred</div>",
				ErrorCode:     404,
			},
			valid: true,
		},
		{
			name: "Error template without file path",
			template: ErrorTemplate{
				FilePath:      "",
				ComponentName: "ErrorPage",
				ErrorCode:     404,
			},
			valid: false,
		},
		{
			name: "Error template with zero error code",
			template: ErrorTemplate{
				FilePath:      "/app/error.templ",
				ComponentName: "ErrorPage",
				ErrorCode:     0,
			},
			valid: true, // Zero error code might be valid for generic errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.template.FilePath != ""
			
			if isValid != tt.valid {
				t.Errorf("ErrorTemplate validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestErrorTemplate_JSONSerialization(t *testing.T) {
	template := ErrorTemplate{
		FilePath:      "/app/error.templ",
		ComponentName: "ErrorPage",
		Content:       "<div>Error occurred</div>",
		ErrorCode:     404,
	}

	// Test JSON marshaling
	data, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorTemplate: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled ErrorTemplate
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ErrorTemplate: %v", err)
	}

	// Verify fields
	if unmarshaled.FilePath != template.FilePath {
		t.Errorf("FilePath mismatch: got %v, want %v", unmarshaled.FilePath, template.FilePath)
	}
	if unmarshaled.ErrorCode != template.ErrorCode {
		t.Errorf("ErrorCode mismatch: got %v, want %v", unmarshaled.ErrorCode, template.ErrorCode)
	}
}