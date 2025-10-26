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

const RequestKey ContextType = "request"
const UserContextKey ContextType = "user"
