package router

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
)

// RouterCore defines the contract for the clean router core
type RouterCore interface {
	Initialize() error
	RegisterRoutes(chiRouter *chi.Mux) error
	GetRoutes() []Route
	GetLayoutTemplates() []LayoutTemplate
	GetErrorTemplates() []ErrorTemplate
	GetMiddlewareSetup() MiddlewareSetup
	GetHandlerBuilder() HandlerBuilder
	GetRouteRegistrar() RouteRegistrar
}

// RouteDiscovery interface for discovering routes, layouts, and error templates
type RouteDiscovery interface {
	DiscoverRoutes(scanPath string) ([]Route, error)
	DiscoverLayouts(scanPath string) ([]LayoutTemplate, error)
	DiscoverErrorTemplates(scanPath string) ([]ErrorTemplate, error)
}

// ConfigLoader interface for loading route configurations
type ConfigLoader interface {
	LoadRouteConfig(templateFile string) (*interfaces.ConfigFile, error)
	LoadConfig(templatePath string) (*interfaces.ConfigFile, error)
	LoadAuthSettings(templatePath string) (*interfaces.AuthSettings, error)
}
