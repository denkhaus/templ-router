package middleware

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)


// routerMiddleware handles router-level middleware configuration (private implementation)
type routerMiddleware struct {
	configService interfaces.ConfigService
	logger        *zap.Logger
}

// NewRouterMiddleware creates a new router middleware for DI
func NewRouterMiddleware(i do.Injector) (RouterMiddlewareInterface, error) {
	configService := do.MustInvoke[interfaces.ConfigService](i)
	logger := do.MustInvoke[*zap.Logger](i)

	return &routerMiddleware{
		configService: configService,
		logger:        logger,
	}, nil
}

// ConfigureRouterMiddleware configures router-level middleware based on configuration
func (rm *routerMiddleware) ConfigureRouterMiddleware(chiRouter *chi.Mux) error {
	rm.logger.Debug("Configuring router middleware")

	// Configure trailing slash redirection
	if rm.configService.GetRouterEnableTrailingSlash() {
		chiRouter.Use(chimiddleware.RedirectSlashes)
		rm.logger.Info("Enabled trailing slash redirection")
	}

	// Configure slash redirection (clean path)
	if rm.configService.GetRouterEnableSlashRedirect() {
		chiRouter.Use(chimiddleware.CleanPath)
		rm.logger.Info("Enabled slash redirection")
	}

	rm.logger.Debug("Router middleware configuration completed")
	return nil
}