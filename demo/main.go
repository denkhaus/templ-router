package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/denkhaus/templ-router/demo/assets"
	"github.com/denkhaus/templ-router/demo/generated/templates"
	"github.com/denkhaus/templ-router/demo/pkg/services"
	"github.com/denkhaus/templ-router/pkg/di"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// main demonstrates the new clean architecture
func main() {
	startupClean(context.Background())
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
		return fmt.Errorf("failed to create template registry: %w", err)
	}

	assetsService, err := assets.NewService(injector)
	if err != nil {
		return fmt.Errorf("failed to create assets service: %w", err)
	}

	userStore, err := services.NewDefaultUserStore(injector)
	if err != nil {
		return fmt.Errorf("failed to create userStore: %w", err)
	}

	// Register application services using options pattern
	container.RegisterApplicationServices(
		di.WithTemplateRegistry(templateRegistry),
		di.WithAssetsService(assetsService),
		di.WithUserStore(userStore),
	)

	// Get logger from container
	logger := container.GetLogger()
	defer logger.Sync()

	logger.Info("Starting application with clean architecture and dependency injection")

	// Create Chi router
	mux := chi.NewRouter()

	// Add auth context middleware
	authMiddleware, err := middleware.NewAuthContextMiddleware(container.GetInjector())
	if err != nil {
		return fmt.Errorf("failed to create auth middleware: %w", err)
	}
	mux.Use(authMiddleware.Middleware)

	// Get clean router from container
	cleanRouter := container.GetRouter()

	// Initialize the router (discover routes, layouts, etc.)
	if err := cleanRouter.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize clean router: %w", err)
	}

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
		return fmt.Errorf("failed to register routes: %w", err)
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
	logger.Info("Starting Clean Architecture Demo Server on :8084")
	if err := http.ListenAndServe(":8084", mux); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
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
			zap.Bool("dynamic", route.IsDynamic))
	}
}
