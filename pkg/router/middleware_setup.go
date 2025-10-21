package router

import (
	"fmt"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// MiddlewareSetup defines the contract for middleware configuration
type MiddlewareSetup interface {
	GetAuthService() interfaces.AuthService
	GetI18nService() interfaces.I18nService
	GetTemplateService() interfaces.TemplateService
	GetLayoutService() interfaces.LayoutService
	GetErrorService() interfaces.ErrorService
	GetAuthMiddleware() middleware.AuthMiddlewareInterface
	GetI18nMiddleware() middleware.I18nMiddlewareInterface
	GetTemplateMiddleware() middleware.TemplateMiddlewareInterface
	GetDataServiceMiddleware() middleware.DataServiceMiddlewareInterface
	ConfigureMiddlewareChain(route interfaces.Route, authSettings interface{}) []interface{}
	ValidateMiddlewareSetup() error
}

// middlewareSetup handles middleware configuration and setup (private implementation)
// SEPARATED FROM: clean_router.go (Separation of Concerns)
type middlewareSetup struct {
	// Clean services (no circular dependencies)
	authService     interfaces.AuthService
	i18nService     interfaces.I18nService
	templateService interfaces.TemplateService
	layoutService   interfaces.LayoutService
	errorService    interfaces.ErrorService

	// Middleware components
	authMiddleware        middleware.AuthMiddlewareInterface
	i18nMiddleware        middleware.I18nMiddlewareInterface
	templateMiddleware    middleware.TemplateMiddlewareInterface
	dataServiceMiddleware middleware.DataServiceMiddlewareInterface

	logger *zap.Logger
}

// NewMiddlewareSetup creates a new middleware setup
func NewMiddlewareSetup(i do.Injector) (MiddlewareSetup, error) {
	// Inject clean services
	authService := do.MustInvoke[interfaces.AuthService](i)
	i18nService := do.MustInvoke[interfaces.I18nService](i)
	templateService := do.MustInvoke[interfaces.TemplateService](i)
	layoutService := do.MustInvoke[interfaces.LayoutService](i)
	errorService := do.MustInvoke[interfaces.ErrorService](i)

	// Inject middleware components
	authMiddleware := do.MustInvoke[middleware.AuthMiddlewareInterface](i)
	i18nMiddleware := do.MustInvoke[middleware.I18nMiddlewareInterface](i)
	templateMiddleware := do.MustInvoke[middleware.TemplateMiddlewareInterface](i)
	dataServiceMiddleware := do.MustInvoke[middleware.DataServiceMiddlewareInterface](i)

	logger := do.MustInvoke[*zap.Logger](i)

	return &middlewareSetup{
		authService:           authService,
		i18nService:           i18nService,
		templateService:       templateService,
		layoutService:         layoutService,
		errorService:          errorService,
		authMiddleware:        authMiddleware,
		i18nMiddleware:        i18nMiddleware,
		templateMiddleware:    templateMiddleware,
		dataServiceMiddleware: dataServiceMiddleware,
		logger:                logger,
	}, nil
}

// GetAuthService returns the auth service
func (ms *middlewareSetup) GetAuthService() interfaces.AuthService {
	return ms.authService
}

// GetI18nService returns the i18n service
func (ms *middlewareSetup) GetI18nService() interfaces.I18nService {
	return ms.i18nService
}

// GetTemplateService returns the template service
func (ms *middlewareSetup) GetTemplateService() interfaces.TemplateService {
	return ms.templateService
}

// GetLayoutService returns the layout service
func (ms *middlewareSetup) GetLayoutService() interfaces.LayoutService {
	return ms.layoutService
}

// GetErrorService returns the error service
func (ms *middlewareSetup) GetErrorService() interfaces.ErrorService {
	return ms.errorService
}

// GetAuthMiddleware returns the auth middleware
func (ms *middlewareSetup) GetAuthMiddleware() middleware.AuthMiddlewareInterface {
	return ms.authMiddleware
}

// GetI18nMiddleware returns the i18n middleware
func (ms *middlewareSetup) GetI18nMiddleware() middleware.I18nMiddlewareInterface {
	return ms.i18nMiddleware
}

// GetTemplateMiddleware returns the template middleware
func (ms *middlewareSetup) GetTemplateMiddleware() middleware.TemplateMiddlewareInterface {
	return ms.templateMiddleware
}

// GetDataServiceMiddleware returns the data service middleware
func (ms *middlewareSetup) GetDataServiceMiddleware() middleware.DataServiceMiddlewareInterface {
	return ms.dataServiceMiddleware
}

// ConfigureMiddlewareChain configures the middleware chain for a specific route
func (ms *middlewareSetup) ConfigureMiddlewareChain(route interfaces.Route, authSettings interface{}) []interface{} {
	var middlewareChain []interface{}

	ms.logger.Debug("Configuring middleware chain",
		zap.String("route", route.Path),
		zap.String("template", route.TemplateFile))

	// Always add i18n middleware for locale extraction
	middlewareChain = append(middlewareChain, ms.i18nMiddleware)

	// Add auth middleware if auth settings are present
	if authSettings != nil {
		ms.logger.Debug("Adding auth middleware to chain",
			zap.String("route", route.Path))
		middlewareChain = append(middlewareChain, ms.authMiddleware)
	}

	// Add data service middleware before template middleware if the route requires data services
	// This ensures data is resolved before template rendering
	if route.RequiresDataService {
		ms.logger.Debug("Adding data service middleware to chain",
			zap.String("route", route.Path),
			zap.String("data_service", route.DataServiceInterface))
		middlewareChain = append(middlewareChain, ms.dataServiceMiddleware)
	}

	// Always add template middleware for rendering
	middlewareChain = append(middlewareChain, ms.templateMiddleware)

	ms.logger.Debug("Middleware chain configured",
		zap.String("route", route.Path),
		zap.Int("middleware_count", len(middlewareChain)))

	return middlewareChain
}

// ValidateMiddlewareSetup validates that all required middleware is properly configured
func (ms *middlewareSetup) ValidateMiddlewareSetup() error {
	ms.logger.Debug("Validating middleware setup")

	// Check that all services are available
	if ms.authService == nil {
		ms.logger.Error("Auth service is nil")
		return fmt.Errorf("auth service not configured")
	}

	if ms.i18nService == nil {
		ms.logger.Error("I18n service is nil")
		return fmt.Errorf("i18n service not configured")
	}

	if ms.templateService == nil {
		ms.logger.Error("Template service is nil")
		return fmt.Errorf("template service not configured")
	}

	if ms.layoutService == nil {
		ms.logger.Error("Layout service is nil")
		return fmt.Errorf("layout service not configured")
	}

	if ms.errorService == nil {
		ms.logger.Error("Error service is nil")
		return fmt.Errorf("error service not configured")
	}

	// Check that all middleware components are available
	if ms.authMiddleware == nil {
		ms.logger.Error("Auth middleware is nil")
		return fmt.Errorf("auth middleware not configured")
	}

	if ms.i18nMiddleware == nil {
		ms.logger.Error("I18n middleware is nil")
		return fmt.Errorf("i18n middleware not configured")
	}

	if ms.templateMiddleware == nil {
		ms.logger.Error("Template middleware is nil")
		return fmt.Errorf("template middleware not configured")
	}

	if ms.dataServiceMiddleware == nil {
		ms.logger.Error("Data service middleware is nil")
		return fmt.Errorf("data service middleware not configured")
	}

	ms.logger.Info("Middleware setup validation successful")
	return nil
}
