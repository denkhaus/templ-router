package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/denkhaus/templ-router/demo/assets"
	"github.com/denkhaus/templ-router/demo/generated/templates"
	"github.com/denkhaus/templ-router/demo/pkg/dataservices"
	"github.com/denkhaus/templ-router/demo/pkg/services"
	"github.com/denkhaus/templ-router/pkg/di"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// main demonstrates the new clean architecture
func main() {
	if err := startupClean(context.Background()); err != nil {
		// Handle startup errors gracefully with structured error handling
		var appErr *shared.AppError
		if errors.As(err, &appErr) {
			// Structured error - log with context and exit gracefully
			fmt.Fprintf(os.Stderr, "Application startup failed: %s\n", appErr.Error())
			if appErr.Context != nil {
				fmt.Fprintf(os.Stderr, "Error context: %+v\n", appErr.Context)
			}
			if appErr.Cause != nil {
				fmt.Fprintf(os.Stderr, "Underlying cause: %v\n", appErr.Cause)
			}
		} else {
			// Generic error - wrap and handle gracefully
			fmt.Fprintf(os.Stderr, "Application startup failed: %v\n", err)
		}
		os.Exit(1)
	}
}

func startupClean(ctx context.Context) error {

	// this is handled by docker compose
	// if err := godotenv.Load(".env"); err != nil {
	// 	return fmt.Errorf("failed to load environment file: %w", err)
	// }

	// Create DI container
	container := di.NewContainer()
	defer container.Shutdown()

	container.RegisterRouterServices("TR")
	injector := container.GetInjector()

	templateRegistry, err := templates.NewRegistry(injector)
	if err != nil {
		return shared.NewServiceError("Failed to create template registry").
			WithCause(err).
			WithContext("component", "template_registry")
	}

	assetsService, err := assets.NewService(injector)
	if err != nil {
		return shared.NewServiceError("Failed to create assets service").
			WithCause(err).
			WithContext("component", "assets_service")
	}

	userStore, err := services.NewDefaultUserStore(injector)
	if err != nil {
		return shared.NewServiceError("Failed to create user store").
			WithCause(err).
			WithContext("component", "user_store")
	}

	// Register application services using options pattern
	container.RegisterApplicationServices(
		di.WithTemplateRegistry(templateRegistry),
		di.WithAssetsService(assetsService),
		di.WithUserStore(userStore),
	)

	// Register DataServices as named dependencies for DataService resolution
	// Use short name without package prefix for cleaner naming
	do.ProvideNamed(container.GetInjector(), "UserDataService", dataservices.NewUserDataService)
	do.ProvideNamed(container.GetInjector(), "ProductDataService", dataservices.NewProductDataService)
	do.ProvideNamed(container.GetInjector(), "OrderDataService", dataservices.NewOrderDataService)
	do.ProvideNamed(container.GetInjector(), "BrokenDataService", dataservices.NewBrokenDataService)
	do.ProvideNamed(container.GetInjector(), "SpecificDataService", dataservices.NewSpecificOnlyDataService)
	do.ProvideNamed(container.GetInjector(), "UserWithIdDataService", dataservices.NewUserWithIdDataService)

	// Get logger from container
	logger := container.GetLogger()
	defer logger.Sync()

	logger.Info("Starting application with clean architecture and dependency injection")

	// Create Chi router
	mux := chi.NewRouter()

	// Get clean router from container
	cleanRouter := container.GetRouter()

	// Initialize the router (discover routes, layouts, etc.)
	if err := cleanRouter.Initialize(); err != nil {
		return shared.NewServiceError("Failed to initialize clean router").
			WithCause(err).
			WithContext("component", "router_initialization")
	}

	// Configure router middleware FIRST (before any routes or other middleware)
	if err := cleanRouter.GetMiddlewareSetup().GetRouterMiddleware().ConfigureRouterMiddleware(mux); err != nil {
		return shared.NewServiceError("Failed to configure router middleware").
			WithCause(err).
			WithContext("component", "router_middleware")
	}

	// Add auth context middleware AFTER router middleware
	authMiddleware, err := middleware.NewAuthContextMiddleware(container.GetInjector())
	if err != nil {
		return shared.NewServiceError("Failed to create auth middleware").
			WithCause(err).
			WithContext("component", "auth_middleware")
	}
	mux.Use(authMiddleware.Middleware)

	// Add API routes
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"status": "healthy",
			"architecture": "clean",
			"dependency_injection": "samber/do",
			"router": "multi-language file-based",
			"i18n": "decentralized",
			"languages": ["en", "de"]
		}`))
	})

	// Register file-based routes
	logger.Info("Registering routes with clean architecture...")
	if err := cleanRouter.RegisterRoutes(mux); err != nil {
		return shared.NewServiceError("Failed to register routes").
			WithCause(err).
			WithContext("component", "route_registration")
	}

	// Register external auth routes (pluggable authentication)
	logger.Info("Registering external auth routes...")
	authHandlers := do.MustInvoke[interfaces.AuthHandlers](injector)
	authHandlers.RegisterRoutes(func(method, path string, handler http.HandlerFunc) {
		switch method {
		case "POST":
			mux.Post(path, handler)
		case "GET":
			mux.Get(path, handler)
		}
		logger.Info("Auth route registered",
			zap.String("method", method),
			zap.String("path", path))
	})

	// Log route information
	logRouteInformation(cleanRouter, logger)

	// Start server
	logger.Info("Starting Clean Architecture Demo Server on 0.0.0.0:8084")
	if err := http.ListenAndServe("0.0.0.0:8084", mux); err != nil {
		return shared.NewServiceError("Failed to start HTTP server").
			WithCause(err).
			WithContext("component", "http_server").
			WithContext("address", "0.0.0.0:8084")
	}

	return nil
}

// logRouteInformation logs information about discovered routes
func logRouteInformation(cleanRouter router.RouterCore, logger *zap.Logger) {
	routes := cleanRouter.GetRoutes()
	layouts := cleanRouter.GetLayoutTemplates()
	errorTemplates := cleanRouter.GetErrorTemplates()

	logger.Info("Route discovery summary",
		zap.Int("routes", len(routes)),
		zap.Int("layouts", len(layouts)),
		zap.Int("error_templates", len(errorTemplates)))

	logger.Info("Available Routes:")
	for _, route := range routes {
		logger.Info("Route registered",
			zap.String("path", route.Path),
			zap.String("template", route.TemplateFile),
			zap.Bool("dynamic", route.IsDynamic),
			zap.Bool("requires_data_service", route.RequiresDataService),
			zap.String("data_service_interface", route.DataServiceInterface))
	}
}
