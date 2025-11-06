package metadata

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
)

// MetadataSettingsParser handles parsing of auth and dynamic settings
// Extracted from metadata.go for better separation of concerns
type MetadataSettingsParser struct{}

// NewMetadataSettingsParser creates a new settings parser
func NewMetadataSettingsParser() *MetadataSettingsParser {
	return &MetadataSettingsParser{}
}

// ParseAuthSettings parses auth settings from YAML into AuthSettings struct
func (msp *MetadataSettingsParser) ParseAuthSettings(authData interface{}) *interfaces.AuthSettings {
	if authData == nil {
		return &interfaces.AuthSettings{Type: "Public"}
	}

	authMap, ok := authData.(map[interface{}]interface{})
	if !ok {
		// Try string-keyed map
		if authMapStr, ok := authData.(map[string]interface{}); ok {
			authMap = make(map[interface{}]interface{})
			for k, v := range authMapStr {
				authMap[k] = v
			}
		} else {
			return &interfaces.AuthSettings{Type: "Public"}
		}
	}

	settings := &interfaces.AuthSettings{
		Type:     "Public",
		Settings: make(map[string]interface{}),
	}

	// Parse auth type - now using string type instead of AuthType enum
	if authType, exists := authMap["type"]; exists {
		if authTypeStr, ok := authType.(string); ok {
			// Normalize auth type string
			switch authTypeStr {
			case "Public", "public":
				settings.Type = "Public"
			case "UserRequired", "User", "protected", "user":
				settings.Type = "UserRequired"
			case "AdminRequired", "Admin", "admin":
				settings.Type = "AdminRequired"
			default:
				settings.Type = "Public"
			}
		}
	}

	// Parse redirect URL
	if redirectURL, exists := authMap["redirect_url"]; exists {
		if redirectURLStr, ok := redirectURL.(string); ok {
			settings.RedirectURL = redirectURLStr
		}
	}

	// Parse required roles (AuthSettings uses Roles, not RequiredPermissions)
	if permissions, exists := authMap["permissions"]; exists {
		if permissionsSlice, ok := permissions.([]interface{}); ok {
			settings.Roles = make([]string, 0, len(permissionsSlice))
			for _, perm := range permissionsSlice {
				if permStr, ok := perm.(string); ok {
					settings.Roles = append(settings.Roles, permStr)
				}
			}
		}
	}

	// Parse roles field directly
	if roles, exists := authMap["roles"]; exists {
		if rolesSlice, ok := roles.([]interface{}); ok {
			if settings.Roles == nil {
				settings.Roles = make([]string, 0, len(rolesSlice))
			}
			for _, role := range rolesSlice {
				if roleStr, ok := role.(string); ok {
					settings.Roles = append(settings.Roles, roleStr)
				}
			}
		}
	}

	// Copy any additional settings
	for key, value := range authMap {
		if keyStr, ok := key.(string); ok {
			switch keyStr {
			case "type", "redirect_url", "permissions", "roles":
				// These are handled above
				continue
			default:
				// Copy additional settings
				settings.Settings[keyStr] = value
			}
		}
	}

	return settings
}

// ParseDynamicSettings parses dynamic parameter settings from YAML into DynamicSettings struct
func (msp *MetadataSettingsParser) ParseDynamicSettings(dynamicData interface{}) *interfaces.DynamicSettings {
	if dynamicData == nil {
		return nil
	}

	dynamicMap, ok := dynamicData.(map[interface{}]interface{})
	if !ok {
		// Try string-keyed map
		if dynamicMapStr, ok := dynamicData.(map[string]interface{}); ok {
			dynamicMap = make(map[interface{}]interface{})
			for k, v := range dynamicMapStr {
				dynamicMap[k] = v
			}
		} else {
			return nil
		}
	}

	// Create DynamicSettings with the new shared.DynamicConfig structure
	settings := &interfaces.DynamicSettings{
		Rules:    make(map[string]*interfaces.DynamicParameterConfig),
		Settings: make(map[string]interface{}),
	}

	// Parse validation rules directly from the dynamic config
	for paramKey, paramValue := range dynamicMap {
		paramName, ok := paramKey.(string)
		if !ok || paramName == "parameters" {
			continue // Skip "parameters" key and non-string keys
		}

		// Try to parse as validation rule
		if paramConfig, ok := paramValue.(map[interface{}]interface{}); ok {
			config := &interfaces.DynamicParameterConfig{
				Name:     paramName,
				Settings: make(map[string]interface{}),
			}

			// Parse validation type
			if validationType, exists := paramConfig["type"]; exists {
				if typeStr, ok := validationType.(string); ok {
					config.Type = typeStr
				}
			}

			// Parse required flag
			if required, exists := paramConfig["required"]; exists {
				if requiredBool, ok := required.(bool); ok {
					config.Required = requiredBool
				}
			}

			// Parse pattern (validation regex)
			if pattern, exists := paramConfig["pattern"]; exists {
				if patternStr, ok := pattern.(string); ok {
					config.Pattern = patternStr
				}
			}

			// Parse default value
			if defaultValue, exists := paramConfig["default"]; exists {
				config.Default = defaultValue
			}

			// Copy additional settings
			for key, value := range paramConfig {
				if keyStr, ok := key.(string); ok {
					switch keyStr {
					case "name", "type", "required", "pattern", "default":
						// These are handled above
						continue
					default:
						config.Settings[keyStr] = value
					}
				}
			}

			settings.Rules[paramName] = config
		} else {
			// Store as general setting
			if paramNameStr, ok := paramKey.(string); ok {
				settings.Settings[paramNameStr] = paramValue
			}
		}
	}

	// Handle legacy "parameters" section if present
	if parametersData, exists := dynamicMap["parameters"]; exists {
		if parametersMap, ok := parametersData.(map[interface{}]interface{}); ok {
			for paramKey, paramValue := range parametersMap {
				paramName, ok := paramKey.(string)
				if !ok {
					continue
				}

				// Try to parse as validation rule
				if paramConfig, ok := paramValue.(map[interface{}]interface{}); ok {
					config := &interfaces.DynamicParameterConfig{
						Name:     paramName,
						Settings: make(map[string]interface{}),
					}

					// Parse validation type (default to "string" if not specified)
					config.Type = "string"
					if validationType, exists := paramConfig["type"]; exists {
						if typeStr, ok := validationType.(string); ok {
							config.Type = typeStr
						}
					}

					// Parse required flag
					if required, exists := paramConfig["required"]; exists {
						if requiredBool, ok := required.(bool); ok {
							config.Required = requiredBool
						}
					}

					// Parse validation pattern (for backward compatibility)
					if validation, exists := paramConfig["validation"]; exists {
						if validationStr, ok := validation.(string); ok {
							config.Pattern = validationStr
						}
					}

					// Parse pattern (new format)
					if pattern, exists := paramConfig["pattern"]; exists {
						if patternStr, ok := pattern.(string); ok {
							config.Pattern = patternStr
						}
					}

					// Parse default value
					if defaultValue, exists := paramConfig["default"]; exists {
						config.Default = defaultValue
					}

					// Copy additional settings
					for key, value := range paramConfig {
						if keyStr, ok := key.(string); ok {
							switch keyStr {
							case "name", "type", "required", "validation", "pattern", "default":
								// These are handled above
								continue
							default:
								config.Settings[keyStr] = value
							}
						}
					}

					settings.Rules[paramName] = config
				}
			}
		}
	}

	return settings
}

// ParseConfigFileFromShared parses a shared.ConfigFile and returns auth and dynamic settings
// This helper function provides a bridge between the new shared.ConfigFile and the existing parsing logic
func (msp *MetadataSettingsParser) ParseConfigFileFromShared(configFile *shared.ConfigFile) (*interfaces.AuthSettings, *interfaces.DynamicSettings) {
	if configFile == nil {
		return msp.ParseAuthSettings(nil), msp.ParseDynamicSettings(nil)
	}

	// Parse auth settings from the shared config
	authSettings := msp.ParseAuthSettings(configFile.GetAuthSettings())

	// Parse dynamic settings from the shared config
	dynamicSettings := msp.ParseDynamicSettings(configFile.GetDynamicSettings())

	return authSettings, dynamicSettings
}