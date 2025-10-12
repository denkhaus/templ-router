package router

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/router/middleware"
)

// DiscoverFiles scans the app directory for *.templ and *.yaml files using FileSystemChecker
func DiscoverFiles(scanPath string, fileSystem middleware.FileSystemChecker) ([]string, []string, error) {
	if scanPath == "" {
		return nil, nil, fmt.Errorf("scan path cannot be empty")
	}

	var templFiles []string
	var yamlFiles []string

	err := fileSystem.WalkDirectory(scanPath, func(path string, isDir bool, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if isDir {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".templ":
			templFiles = append(templFiles, path)
		case ".yaml", ".yml":
			// For YAML files named like *.templ.yaml, check if the corresponding *.templ file exists
			// E.g., test.templ.yaml corresponds to test.templ
			if strings.HasSuffix(path, ".templ.yaml") || strings.HasSuffix(path, ".templ.yml") {
				// Remove the .yaml/.yml extension to get the base name, which should end with .templ
				basePath := strings.TrimSuffix(path, ext) // This removes .yaml or .yml
				// Now check if the basePath (which should be like filename.templ) exists as a file
				if fileSystem.FileExists(basePath) {
					yamlFiles = append(yamlFiles, path)
				}
			} else {
				// For other YAML files, just add them without checking for corresponding template
				yamlFiles = append(yamlFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan directory %s: %w", scanPath, err)
	}

	return templFiles, yamlFiles, nil
}
