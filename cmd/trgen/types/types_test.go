package types

import (
	"testing"
)

func TestTemplateInfo(t *testing.T) {
	// Test TemplateInfo struct creation and field access
	template := TemplateInfo{
		FunctionName:  "Page",
		PackageName:   "app",
		ImportPath:    "github.com/user/project/app",
		PackageAlias:  "app",
		RoutePattern:  "/",
		TemplateKey:   "test-key",
		FilePath:      "/path/to/file.go",
		HumanName:     "Page",
	}

	if template.FunctionName != "Page" {
		t.Errorf("Expected FunctionName to be 'Page', got %s", template.FunctionName)
	}

	if template.PackageName != "app" {
		t.Errorf("Expected PackageName to be 'app', got %s", template.PackageName)
	}

	if template.ImportPath != "github.com/user/project/app" {
		t.Errorf("Expected ImportPath to be 'github.com/user/project/app', got %s", template.ImportPath)
	}

	if template.PackageAlias != "app" {
		t.Errorf("Expected PackageAlias to be 'app', got %s", template.PackageAlias)
	}

	if template.RoutePattern != "/" {
		t.Errorf("Expected RoutePattern to be '/', got %s", template.RoutePattern)
	}

	if template.TemplateKey != "test-key" {
		t.Errorf("Expected TemplateKey to be 'test-key', got %s", template.TemplateKey)
	}

	if template.FilePath != "/path/to/file.go" {
		t.Errorf("Expected FilePath to be '/path/to/file.go', got %s", template.FilePath)
	}

	if template.HumanName != "Page" {
		t.Errorf("Expected HumanName to be 'Page', got %s", template.HumanName)
	}
}

func TestConfig(t *testing.T) {
	// Test Config struct creation and field access
	config := Config{
		ScanPath:    "app",
		OutputDir:   "generated/templates",
		ModuleName:  "github.com/user/project",
		PackageName: "templates",
	}

	if config.ScanPath != "app" {
		t.Errorf("Expected ScanPath to be 'app', got %s", config.ScanPath)
	}

	if config.OutputDir != "generated/templates" {
		t.Errorf("Expected OutputDir to be 'generated/templates', got %s", config.OutputDir)
	}

	if config.ModuleName != "github.com/user/project" {
		t.Errorf("Expected ModuleName to be 'github.com/user/project', got %s", config.ModuleName)
	}

	if config.PackageName != "templates" {
		t.Errorf("Expected PackageName to be 'templates', got %s", config.PackageName)
	}
}

func TestTemplateInfoZeroValues(t *testing.T) {
	// Test zero values of TemplateInfo
	var template TemplateInfo

	if template.FunctionName != "" {
		t.Errorf("Expected zero value FunctionName to be empty, got %s", template.FunctionName)
	}

	if template.PackageName != "" {
		t.Errorf("Expected zero value PackageName to be empty, got %s", template.PackageName)
	}

	if template.ImportPath != "" {
		t.Errorf("Expected zero value ImportPath to be empty, got %s", template.ImportPath)
	}

	if template.PackageAlias != "" {
		t.Errorf("Expected zero value PackageAlias to be empty, got %s", template.PackageAlias)
	}

	if template.RoutePattern != "" {
		t.Errorf("Expected zero value RoutePattern to be empty, got %s", template.RoutePattern)
	}

	if template.TemplateKey != "" {
		t.Errorf("Expected zero value TemplateKey to be empty, got %s", template.TemplateKey)
	}

	if template.FilePath != "" {
		t.Errorf("Expected zero value FilePath to be empty, got %s", template.FilePath)
	}

	if template.HumanName != "" {
		t.Errorf("Expected zero value HumanName to be empty, got %s", template.HumanName)
	}
}

func TestConfigZeroValues(t *testing.T) {
	// Test zero values of Config
	var config Config

	if config.ScanPath != "" {
		t.Errorf("Expected zero value ScanPath to be empty, got %s", config.ScanPath)
	}

	if config.OutputDir != "" {
		t.Errorf("Expected zero value OutputDir to be empty, got %s", config.OutputDir)
	}

	if config.ModuleName != "" {
		t.Errorf("Expected zero value ModuleName to be empty, got %s", config.ModuleName)
	}

	if config.PackageName != "" {
		t.Errorf("Expected zero value PackageName to be empty, got %s", config.PackageName)
	}
}

func TestTemplateInfoSlice(t *testing.T) {
	// Test working with slices of TemplateInfo
	templates := []TemplateInfo{
		{
			FunctionName: "Page",
			PackageName:  "app",
			RoutePattern: "/",
		},
		{
			FunctionName: "Layout",
			PackageName:  "app",
			RoutePattern: "/layout",
		},
		{
			FunctionName: "Page",
			PackageName:  "admin",
			RoutePattern: "/admin",
		},
	}

	if len(templates) != 3 {
		t.Errorf("Expected 3 templates, got %d", len(templates))
	}

	// Test accessing slice elements
	if templates[0].FunctionName != "Page" {
		t.Errorf("Expected first template FunctionName to be 'Page', got %s", templates[0].FunctionName)
	}

	if templates[1].RoutePattern != "/layout" {
		t.Errorf("Expected second template RoutePattern to be '/layout', got %s", templates[1].RoutePattern)
	}

	if templates[2].PackageName != "admin" {
		t.Errorf("Expected third template PackageName to be 'admin', got %s", templates[2].PackageName)
	}
}