package router

import (
	"fmt"
	"strings"
)

// AuthType represents the authentication requirement for a route
type AuthType string

const (
	// AuthTypePublic allows access without authentication
	AuthTypePublic AuthType = "Public"
	
	// AuthTypeUserRequired requires any authenticated user
	AuthTypeUserRequired AuthType = "UserRequired"
	
	// AuthTypeAdminRequired requires admin privileges
	AuthTypeAdminRequired AuthType = "AdminRequired"
)

// String returns the string representation of AuthType
func (at AuthType) String() string {
	return string(at)
}

// IsValid checks if the AuthType is valid
func (at AuthType) IsValid() bool {
	switch at {
	case AuthTypePublic, AuthTypeUserRequired, AuthTypeAdminRequired:
		return true
	default:
		return false
	}
}

// RequiresAuthentication returns true if the auth type requires authentication
func (at AuthType) RequiresAuthentication() bool {
	return at != AuthTypePublic
}

// RequiresAdmin returns true if the auth type requires admin privileges
func (at AuthType) RequiresAdmin() bool {
	return at == AuthTypeAdminRequired
}

// ParseAuthType parses a string into an AuthType
func ParseAuthType(s string) (AuthType, error) {
	s = strings.TrimSpace(s)
	
	// Only accept the new format - no legacy support
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

// AuthSettings represents authentication settings for a route
type AuthSettings struct {
	Type        AuthType  `yaml:"type"`
	RedirectURL string    `yaml:"redirect_url,omitempty"`
}

