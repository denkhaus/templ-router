package utils

import (
	"testing"
	"github.com/denkhaus/templ-router/cmd/template-generator/types"
)

func TestSanitizePackageName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal package name",
			input:    "admin",
			expected: "admin",
		},
		{
			name:     "Package with hyphen",
			input:    "error-demo",
			expected: "errordemo",
		},
		{
			name:     "Package with multiple hyphens",
			input:    "my-test-package",
			expected: "mytestpackage",
		},
		{
			name:     "Package with dots",
			input:    "package.name",
			expected: "packagename",
		},
		{
			name:     "Package with mixed invalid chars",
			input:    "my-package.test",
			expected: "mypackagetest",
		},
		{
			name:     "Package starting with number",
			input:    "123package",
			expected: "pkg123package",
		},
		{
			name:     "Package with underscores (valid)",
			input:    "my_package",
			expected: "my_package",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "pkg",
		},
		{
			name:     "Only invalid characters",
			input:    "---...",
			expected: "pkg",
		},
		{
			name:     "Mixed case with hyphens",
			input:    "Error-Demo",
			expected: "ErrorDemo",
		},
		{
			name:     "Special characters",
			input:    "package@#$%",
			expected: "package",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := sanitizePackageName(tt.input)
			if actual != tt.expected {
				t.Errorf("sanitizePackageName(%q) = %q, want %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestCreatePackageAlias(t *testing.T) {
	tests := []struct {
		name       string
		packageName string
		importPath string
		config     types.Config
		expected   string
	}{
		{
			name:        "Root package",
			packageName: "app",
			importPath:  "github.com/user/project/app",
			config:      types.Config{ScanPath: "app"},
			expected:    "app",
		},
		{
			name:        "Normal subdirectory",
			packageName: "admin",
			importPath:  "github.com/user/project/app/admin",
			config:      types.Config{ScanPath: "app"},
			expected:    "admin",
		},
		{
			name:        "Directory with hyphen",
			packageName: "errordemo",
			importPath:  "github.com/user/project/app/error-demo",
			config:      types.Config{ScanPath: "app"},
			expected:    "errordemo",
		},
		{
			name:        "Complex path with hyphens",
			packageName: "testpackage",
			importPath:  "github.com/user/project/app/my-test-package",
			config:      types.Config{ScanPath: "app"},
			expected:    "mytestpackage",
		},
		{
			name:        "Path with dots",
			packageName: "packagename",
			importPath:  "github.com/user/project/app/package.name",
			config:      types.Config{ScanPath: "app"},
			expected:    "packagename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CreatePackageAlias(tt.packageName, tt.importPath, tt.config)
			if actual != tt.expected {
				t.Errorf("CreatePackageAlias(%q, %q, %v) = %q, want %q", 
					tt.packageName, tt.importPath, tt.config, actual, tt.expected)
			}
		})
	}
}