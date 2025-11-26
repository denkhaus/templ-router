package shared

import (
	"fmt"
	"strings"
)

type ContextType string

const (
	UserContextKey        ContextType = "user"
	LocaleKey             ContextType = "locale"
	TemplateConfigKey     ContextType = "template_config"
	TemplatePathKey       ContextType = "template_path"
	I18nDataKey           ContextType = "router_i18n_data"
	I18nTemplateKey       ContextType = "router_i18n_template"
	RequestPathKey        ContextType = "request_path"
	RouteMappingKey       ContextType = "route_mapping"
	ComponentsMetadataKey ContextType = "components_metadata"
)

// AuthType represents different authentication types
type AuthType string

const (
	AuthTypePublic AuthType = "Public"
	AuthTypeUser   AuthType = "UserRequired"
	AuthTypeAdmin  AuthType = "AdminRequired"
)

// String returns the string representation of AuthType
func (at AuthType) String() string {
	return string(at)
}

// IsValid checks if the AuthType is valid
func (at AuthType) IsValid() bool {
	switch at {
	case AuthTypePublic, AuthTypeUser, AuthTypeAdmin:
		return true
	default:
		return false
	}
}

// ParseAuthType parses a string into an AuthType
func ParseAuthType(s string) (AuthType, error) {
	s = strings.TrimSpace(s)

	authType := AuthType(s)
	if authType.IsValid() {
		return authType, nil
	}

	// Default to Public for empty strings
	if s == "" {
		return AuthTypePublic, nil
	}

	return "", fmt.Errorf("invalid auth type: %s. Valid types are: Public, UserRequired, AdminRequired", s)
}
