package metadata

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
)

// MetadataSettingsParser handles parsing of auth and dynamic settings
// Extracted from metadata.go for better separation of concerns
type MetadataSettingsParser struct{}

// NewMetadataSettingsParser creates a new settings parser
func NewMetadataSettingsParser() *MetadataSettingsParser {
	return &MetadataSettingsParser{}
}

// parseAuthSettings parses auth settings from YAML into AuthSettings struct
func (msp *MetadataSettingsParser) ParseAuthSettings(authData interface{}) *interfaces.AuthSettings {
	if authData == nil {
		return &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}
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
			return &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}
		}
	}

	settings := &interfaces.AuthSettings{Type: interfaces.AuthTypePublic}

	// Parse auth type
	if authType, exists := authMap["type"]; exists {
		if authTypeStr, ok := authType.(string); ok {
			switch authTypeStr {
			case "public":
				settings.Type = interfaces.AuthTypePublic
			case "protected", "user":
				settings.Type = interfaces.AuthTypeUser
			case "admin":
				settings.Type = interfaces.AuthTypeAdmin
			default:
				settings.Type = interfaces.AuthTypePublic
			}
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

	return settings
}

// parseDynamicSettings parses dynamic parameter settings from YAML into DynamicSettings struct
func (msp *MetadataSettingsParser) parseDynamicSettings(dynamicData interface{}) *interfaces.DynamicSettings {
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

	// Extract parameters section
	parametersData, exists := dynamicMap["parameters"]
	if !exists {
		return nil
	}

	parametersMap, ok := parametersData.(map[interface{}]interface{})
	if !ok {
		// Try string-keyed map
		if parametersMapStr, ok := parametersData.(map[string]interface{}); ok {
			parametersMap = make(map[interface{}]interface{})
			for k, v := range parametersMapStr {
				parametersMap[k] = v
			}
		} else {
			return nil
		}
	}

	settings := &interfaces.DynamicSettings{
		Parameters: make(map[string]*interfaces.DynamicParameterConfig),
	}

	// Parse each parameter configuration
	for paramKey, paramValue := range parametersMap {
		paramName, ok := paramKey.(string)
		if !ok {
			continue
		}

		paramConfig, ok := paramValue.(map[interface{}]interface{})
		if !ok {
			// Try string-keyed map
			if paramConfigStr, ok := paramValue.(map[string]interface{}); ok {
				paramConfig = make(map[interface{}]interface{})
				for k, v := range paramConfigStr {
					paramConfig[k] = v
				}
			} else {
				continue
			}
		}

		config := &interfaces.DynamicParameterConfig{}

		// Parse validation regex
		if validation, exists := paramConfig["validation"]; exists {
			if validationStr, ok := validation.(string); ok {
				config.Validation = validationStr
			}
		}

		// Parse description
		if description, exists := paramConfig["description"]; exists {
			if descriptionStr, ok := description.(string); ok {
				config.Description = descriptionStr
			}
		}

		// Parse supported values
		if supportedValues, exists := paramConfig["supported_values"]; exists {
			if supportedSlice, ok := supportedValues.([]interface{}); ok {
				config.SupportedValues = make([]string, 0, len(supportedSlice))
				for _, value := range supportedSlice {
					if valueStr, ok := value.(string); ok {
						config.SupportedValues = append(config.SupportedValues, valueStr)
					}
				}
			}
		}

		settings.Parameters[paramName] = config
	}

	return settings
}
