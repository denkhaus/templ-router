package router

import (
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
)

// FindLayoutForTemplate finds the appropriate layout for a given template
func FindLayoutForTemplate(templatePath string, layoutTemplates []interfaces.LayoutTemplate) *interfaces.LayoutTemplate {
	// Extract the directory of the template
	templateDir := filepath.Dir(templatePath)

	// Find all potential layout files by going up the directory tree
	layoutsByLevel := make(map[int]interfaces.LayoutTemplate)

	// Check the same directory as the template first (closest layout overrides)
	currentDir := templateDir
	level := 0

	for currentDir != "." && currentDir != "/" {
		// Look for layout in current directory
		layoutPath := filepath.Join(currentDir, "layout.templ")

		for _, layout := range layoutTemplates {
			if layout.FilePath == layoutPath {
				layoutsByLevel[layout.LayoutLevel] = layout
				break
			}
		}

		// Move up one directory level
		currentDir = filepath.Dir(currentDir)
		level++
	}

	// Find the closest layout (highest level number)
	var closestLayout *interfaces.LayoutTemplate
	maxLevel := -1
	for level, layout := range layoutsByLevel {
		if level > maxLevel {
			maxLevel = level
			tempLayout := layout
			closestLayout = &tempLayout
		}
	}

	return closestLayout
}

// GetLayoutHierarchy returns the layout hierarchy from the root to the specific template
func GetLayoutHierarchy(templatePath string, layoutTemplates []interfaces.LayoutTemplate) []*interfaces.LayoutTemplate {
	var hierarchy []*interfaces.LayoutTemplate

	// Extract the directory of the template
	templateDir := filepath.Dir(templatePath)

	// Go up the directory tree and collect layouts
	currentDir := templateDir
	level := 0

	for currentDir != "." && currentDir != "/" {
		// Look for layout in current directory
		layoutPath := filepath.Join(currentDir, "layout.templ")

		for _, layout := range layoutTemplates {
			if layout.FilePath == layoutPath {
				// Add to the beginning of the hierarchy to maintain top-down order
				tempLayout := layout
				hierarchy = append([]*interfaces.LayoutTemplate{&tempLayout}, hierarchy...)
				break
			}
		}

		// Move up one directory level
		currentDir = filepath.Dir(currentDir)
		level++
	}

	return hierarchy
}

// ProcessLayoutsForApp processes all layout templates in the app directory to build layout information
func ProcessLayoutsForApp(templates []interfaces.Template) []interfaces.LayoutTemplate {
	var layoutTemplates []interfaces.LayoutTemplate

	// Filter layout templates
	for _, template := range templates {
		if template.Type == "layout" {
			// Calculate the layout level by counting the number of directory segments
			// "app" has 1 segment (level 1), "app/dashboard" has 2 segments (level 2), etc.
			// This matches the test expectations
			segments := strings.Split(filepath.ToSlash(template.DirectoryPath), "/")
			level := len(segments)

			layoutTemplate := interfaces.LayoutTemplate{
				FilePath:      template.FilePath,
				DirectoryPath: template.DirectoryPath,
				LayoutLevel:   level,
			}

			layoutTemplates = append(layoutTemplates, layoutTemplate)
		}
	}

	return layoutTemplates
}
