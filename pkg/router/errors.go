package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/metadata"
)

// ErrorTemplate represents an error.templ file for error page presentation
type ErrorTemplate struct {
	// FilePath is the full path to the error.templ file
	FilePath string

	// DirectoryPath is the directory containing this error template
	DirectoryPath string

	// ErrorTypes is a list of error types handled by this template
	ErrorTypes []string

	// ParentErrorTemplate is the path to parent error template if this doesn't override completely
	ParentErrorTemplate string

	// PrecedenceLevel is the level of precedence (closer templates override further ones)
	PrecedenceLevel int

	// ErrorMessages contains mapping of error codes to specific messages
	ErrorMessages map[int]string
}

// FindErrorTemplateForPath finds the appropriate error template for a given path with configurable root directory
func FindErrorTemplateForPath(path string, errorTemplates []ErrorTemplate, rootDirectory, templateExtension string) *ErrorTemplate {
	// If path is empty, use root
	if path == "" {
		path = "/"
	}

	// Convert path to directory format by removing leading slash and replacing internal slashes with file separators
	directoryPath := strings.TrimPrefix(path, "/")
	directoryPath = filepath.FromSlash(directoryPath)

	// Find error templates by going up the directory tree
	templatesByLevel := make(map[int]ErrorTemplate)

	// Check the requested directory and parent directories
	currentDir := directoryPath
	level := 0

	for currentDir != "." && currentDir != "" && currentDir != "/" {
		// Look for error template in current directory
		errorPath := filepath.Join(rootDirectory, currentDir, "error"+templateExtension)

		for _, errorTemplate := range errorTemplates {
			if errorTemplate.FilePath == errorPath {
				templatesByLevel[errorTemplate.PrecedenceLevel] = errorTemplate
				break
			}
		}

		// Move up one directory level
		currentDir = filepath.Dir(currentDir)
		level++
	}

	// Also check the root directory
	rootErrorPath := filepath.Join(rootDirectory, "error"+templateExtension)
	for _, errorTemplate := range errorTemplates {
		if errorTemplate.FilePath == rootErrorPath {
			templatesByLevel[errorTemplate.PrecedenceLevel] = errorTemplate
			break
		}
	}

	// Find the closest error template (highest precedence level number)
	var closestErrorTemplate *ErrorTemplate
	highestLevel := -1
	for level, template := range templatesByLevel {
		if level > highestLevel {
			highestLevel = level
			tempTemplate := template
			closestErrorTemplate = &tempTemplate
		}
	}

	return closestErrorTemplate
}

// ProcessErrorTemplates processes all error templates in the app directory to build error template information
func ProcessErrorTemplates(templates []interfaces.Template) []ErrorTemplate {
	var errorTemplates []ErrorTemplate

	// Filter error templates
	for _, template := range templates {
		if template.Type == "error" {
			// Calculate the precedence level by counting the number of directory segments
			// "app" has 1 segment (level 1), "app/dashboard" has 2 segments (level 2), etc.
			// This matches the test expectations
			segments := strings.Split(filepath.ToSlash(template.DirectoryPath), "/")
			level := len(segments)

			errorTemplate := ErrorTemplate{
				FilePath:      template.FilePath,
				DirectoryPath: template.DirectoryPath,
				// Default error types handled by this template
				ErrorTypes:      []string{"404", "500"},
				PrecedenceLevel: level,
				ErrorMessages:   make(map[int]string),
			}

			// Load error messages from the YAML config file if it exists
			loadErrorMessagesFromConfig(&errorTemplate)

			errorTemplates = append(errorTemplates, errorTemplate)
		}
	}

	return errorTemplates
}

// loadErrorMessagesFromConfig loads error messages from the YAML config file for this error template
func loadErrorMessagesFromConfig(errorTemplate *ErrorTemplate) {
	// Parse the YAML file if it exists
	config, err := metadata.ParseYAMLMetadataForTemplate(errorTemplate.FilePath)
	if err != nil || config == nil {
		// If there's no config, use default messages
		return
	}

	// Extract error messages from the config
	if errorSettings, exists := config.ErrorSettings.(map[interface{}]interface{}); exists {
		if errorTypes, ok := errorSettings["error_types"].(map[interface{}]interface{}); ok {
			for code, message := range errorTypes {
				if codeStr, codeOk := code.(string); codeOk {
					if messageStr, msgOk := message.(string); msgOk {
						var intCode int
						if _, err := fmt.Sscanf(codeStr, "%d", &intCode); err == nil {
							errorTemplate.ErrorMessages[intCode] = messageStr
						}
					}
				}
			}
		}
	} else if errorSettings, exists := config.ErrorSettings.(map[string]interface{}); exists {
		if errorTypes, ok := errorSettings["error_types"].(map[string]interface{}); ok {
			for code, message := range errorTypes {
				if messageStr, ok := message.(string); ok {
					var intCode int
					if _, err := fmt.Sscanf(code, "%d", &intCode); err == nil {
						errorTemplate.ErrorMessages[intCode] = messageStr
					}
				}
			}
		}
	}
}

// CreateErrorHandler creates an HTTP error handler using the appropriate error template
func CreateErrorHandler(errorTemplate *ErrorTemplate, statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use the error template to generate an appropriate error page
		// Include the status code and any error message in the template

		// Set the correct status code in the response
		w.WriteHeader(statusCode)

		var errorMessage string
		if errorTemplate != nil && errorTemplate.ErrorMessages != nil {
			if msg, exists := errorTemplate.ErrorMessages[statusCode]; exists {
				errorMessage = msg
			} else {
				// Use a default error message if no specific message was configured
				switch statusCode {
				case http.StatusNotFound:
					errorMessage = "The requested page could not be found."
				case http.StatusInternalServerError:
					errorMessage = "An internal server error occurred."
				case http.StatusForbidden:
					errorMessage = "Access to this resource is forbidden."
				case http.StatusUnauthorized:
					errorMessage = "You are not authorized to access this resource."
				default:
					errorMessage = fmt.Sprintf("An error occurred (code: %d).", statusCode)
				}
			}
		} else {
			// Use default error message if no specific template was found
			switch statusCode {
			case http.StatusNotFound:
				errorMessage = "The requested page could not be found."
			case http.StatusInternalServerError:
				errorMessage = "An internal server error occurred."
			case http.StatusForbidden:
				errorMessage = "Access to this resource is forbidden."
			case http.StatusUnauthorized:
				errorMessage = "You are not authorized to access this resource."
			default:
				errorMessage = fmt.Sprintf("An error occurred (code: %d).", statusCode)
			}
		}

		// In a real implementation, this would render the actual error.templ file
		// For now, we'll return an HTML response with proper content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		var errorTemplatePath string
		if errorTemplate != nil {
			errorTemplatePath = filepath.Base(errorTemplate.FilePath)
		} else {
			errorTemplatePath = "default error template"
		}

		response := `<!DOCTYPE html>
<html>
<head>
	<title>Error ` + fmt.Sprint(statusCode) + `</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			text-align: center;
			padding: 50px;
			background-color: #f0f0f0;
		}
		.error-container {
			max-width: 600px;
			margin: 0 auto;
			background: white;
			padding: 30px;
			border-radius: 8px;
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		}
		h1 {
			color: #d32f2f;
			font-size: 60px;
			margin: 0;
			line-height: 1;
		}
		h2 {
			color: #555;
			font-size: 24px;
			margin: 20px 0 10px 0;
		}
		p {
			color: #666;
			font-size: 16px;
			line-height: 1.5;
		}
		.error-code {
			display: block;
			font-size: 14px;
			color: #999;
			margin-top: 20px;
		}
	</style>
</head>
<body>
	<div class="error-container">
		<h1>` + fmt.Sprint(statusCode) + `</h1>
		<h2>Oops! Something went wrong.</h2>
		<p>` + errorMessage + `</p>
		<span class="error-code">Error template: ` + errorTemplatePath + `</span>
	</div>
</body>
</html>`

		w.Write([]byte(response))
	}
}
