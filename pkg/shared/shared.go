package shared

import "strings"

// DeriveMethodNameFromDataType converts "*dataservices.UserData" to "GetUserData"
// This is a shared helper function used by both optimized_template_service.go and data_service_resolver.go
func DeriveMethodNameFromDataType(parameterType string) string {
	// Extract type name from "*github.com/path/dataservices.UserData"
	parts := strings.Split(parameterType, ".")
	if len(parts) > 0 {
		typeName := parts[len(parts)-1]
		// Remove pointer prefix if present
		typeName = strings.TrimPrefix(typeName, "*")
		return "Get" + typeName
	}
	return "GetData" // fallback
}

type ContextType string

const (
	UserContextKey    ContextType = "user"
	LocaleKey         ContextType = "locale"
	TemplateConfigKey ContextType = "template_config"
	TemplatePathKey   ContextType = "template_path"
	I18nDataKey       ContextType = "router_i18n_data"
	I18nTemplateKey   ContextType = "router_i18n_template"
)
