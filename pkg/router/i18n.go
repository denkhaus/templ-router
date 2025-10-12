package router

import (
	"path/filepath"
	"strings"
)

// GenerateI18nIdentifiers creates i18n identifiers based on YAML metadata files associated with each template
// Instead of generating opinionated identifiers, we now load the actual translations from the YAML file
func GenerateI18nIdentifiers(templatePath string) []InternationalizationIdentifier {
	var identifiers []InternationalizationIdentifier
	
	// Extract the path without the app/ prefix and the file extension
	ext := filepath.Ext(templatePath)
	basePath := strings.TrimSuffix(templatePath, ext)
	
	// LIBRARY-AGNOSTIC: Extract path components without hardcoded assumptions
	// This works regardless of the root directory name
	pathParts := strings.Split(basePath, "/")
	var nonEmptySegments []string
	if len(pathParts) > 1 {
		// Skip the first part (root directory) and use the rest
		nonEmptySegments = pathParts[1:]
	} else {
		nonEmptySegments = pathParts
	}
	
	// Generate common identifiers with default values based on the template path
	segments := strings.Split(basePath, "/")
	
	// Skip empty segments (e.g., from leading slash)
	for _, segment := range segments {
		if segment != "" {
			nonEmptySegments = append(nonEmptySegments, segment)
		}
	}
	
	// Generate common identifiers based on the template path
	if len(nonEmptySegments) > 0 {
		// Generate title identifier
		titleKey := strings.Join(append(nonEmptySegments, "title"), ".")
		titleId := InternationalizationIdentifier{
			Key:          titleKey,
			Source:       "opinionated-schema",
			TemplatePath: templatePath,
			DefaultValue: getDefaultTitleValue(nonEmptySegments),
		}
		identifiers = append(identifiers, titleId)
		
		// Generate description identifier
		descriptionKey := strings.Join(append(nonEmptySegments, "description"), ".")
		descriptionId := InternationalizationIdentifier{
			Key:          descriptionKey,
			Source:       "opinionated-schema",
			TemplatePath: templatePath,
			DefaultValue: getDefaultDescriptionValue(nonEmptySegments),
		}
		identifiers = append(identifiers, descriptionId)
		
		// Generate other common identifiers
		// Add more identifiers as needed based on the template type
		if strings.Contains(templatePath, "create") || strings.Contains(templatePath, "new") {
			// Create-specific identifiers
			buttonSubmitKey := strings.Join(append(nonEmptySegments, "submit"), ".")
			buttonSubmitId := InternationalizationIdentifier{
				Key:          buttonSubmitKey,
				Source:       "opinionated-schema",
				TemplatePath: templatePath,
				DefaultValue: "Submit",
			}
			identifiers = append(identifiers, buttonSubmitId)
			
			buttonCancelKey := strings.Join(append(nonEmptySegments, "cancel"), ".")
			buttonCancelId := InternationalizationIdentifier{
				Key:          buttonCancelKey,
				Source:       "opinionated-schema",
				TemplatePath: templatePath,
				DefaultValue: "Cancel",
			}
			identifiers = append(identifiers, buttonCancelId)
		}
	}
	
	// Now try to load translations from the YAML metadata file for this template
	// These will override the default values with actual translations if available
	config, err := ParseYAMLMetadataForTemplate(templatePath)
	if err != nil || config == nil {
		// If there's no YAML file for this template, return the identifiers with default values
		return identifiers
	}
	
	// Replace default values with translations from the YAML file
	for i := range identifiers {
		if translation, exists := config.I18nMappings[identifiers[i].Key]; exists {
			// If there's a direct translation for this key in the YAML file, use it
			identifiers[i].DefaultValue = translation
			identifiers[i].Locales = map[string]string{
				"en": translation, // For simplicity, assuming the translation in the YAML is in English
			}
		}
	}
	
	return identifiers
}

// getDefaultTitleValue generates a default title value from path segments
func getDefaultTitleValue(segments []string) string {
	if len(segments) == 0 {
		return "Home"
	}
	
	// Take the last segment and convert it to a title format
	lastSegment := segments[len(segments)-1]
	
	// If it's a dynamic segment like $id, use a generic name
	if strings.HasPrefix(lastSegment, "$") {
		if lastSegment == "$locale" {
			return "Locale Page"
		}
		return "Dynamic Page"
	}
	
	// Convert snake-case, kebab-case, or camelCase to title format
	title := strings.ReplaceAll(lastSegment, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	
	// Capitalize the first letter of each word
	words := strings.Split(title, " ")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	
	return strings.Join(words, " ")
}

// getDefaultDescriptionValue generates a default description value from path segments
func getDefaultDescriptionValue(segments []string) string {
	if len(segments) == 0 {
		return "Home page description"
	}
	
	// Use the path segments to form a description
	titleValue := getDefaultTitleValue(segments)
	return titleValue + " page"
}

// IntegrateWithCtxI18n provides a way to integrate generated identifiers with the existing ctxi18n system
// This would involve working with the project's existing i18n context system
func IntegrateWithCtxI18n(identifiers []InternationalizationIdentifier) {
	// In a real implementation, this function would integrate the generated i18n identifiers
	// with the existing ctxi18n system in the project
	
	// For now, this serves as a placeholder to indicate where the integration would happen
	// The actual implementation would depend on how ctxi18n is used in the project
	_ = identifiers
}