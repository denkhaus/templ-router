package interfaces

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

// DataServiceInfo holds information about required data services
type DataServiceInfo struct {
	InterfaceType string // e.g., "UserDataService" - used for both display and DI resolution
	ParameterType string // e.g., "*dataservices.UserData" (for display)
	MethodName    string // e.g., "GetUserData" - specific method name
}

// RouterContext provides unified access to request data for Data Services
// It encapsulates URL parameters, query parameters, HTTP request, and Chi context
type RouterContext interface {
	// Context returns the underlying context.Context
	Context() context.Context

	// URL Parameter access (from Chi router path parameters like /{id})
	GetURLParam(key string) string
	GetAllURLParams() map[string]string

	// Query Parameter access (from URL query string like ?page=5&size=10)
	GetQueryParam(key string) string
	GetQueryParams(key string) []string
	GetAllQueryParams() url.Values

	// Original request access for advanced scenarios
	Request() *http.Request

	// Chi-specific context access
	ChiContext() *chi.Context
}

// DataServiceResolver resolves data services from DI container
type DataServiceResolver interface {
	// ResolveDataService resolves a data service by interface type from DI
	ResolveDataService(interfaceType string) (interface{}, error)

	// HasDataService checks if a data service is registered in DI
	HasDataService(interfaceType string) bool
}
