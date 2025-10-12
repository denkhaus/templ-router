package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/denkhaus/templ-router/demo/generated/templates"
	"github.com/denkhaus/templ-router/pkg/di"
	"github.com/denkhaus/templ-router/pkg/router"
	"github.com/denkhaus/templ-router/pkg/services/auth"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// main demonstrates the new clean architecture
func main() {
	startupClean()
}

func startupClean() {
	ctx := context.Background()

	// Create DI container
	container := di.NewContainer()
	defer container.Shutdown()

	// Step 2: Create and register the application-specific template registry
	templateRegistry, err := templates.NewTemplateRegistry(container.GetInjector())
	if err != nil {
		log.Fatal("Failed to create template registry:", err)
	}
	container.RegisterTemplateRegistry(templateRegistry)

	// Register all services
	if err := container.RegisterRouterServices(ctx); err != nil {
		log.Fatal("Failed to register services:", err)
	}

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
		logger.Fatal("Failed to initialize clean router", zap.Error(err))
	}

	// Serve static assets
	mux.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

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
		fmt.Printf("FATAL: Failed to register routes - %v\n", err)
		os.Exit(1)
	}

	// Register external auth routes (pluggable authentication)
	logger.Info("Registering external auth routes...")
	authHandlers := do.MustInvoke[auth.AuthHandlersInterface](container.GetInjector())
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
		logger.Fatal("Server failed to start", zap.Error(err))
	}
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

	logger.Info("")
	logger.Info("Clean Architecture Features:")
	logger.Info("  - Separation of Concerns with middleware pattern")
	logger.Info("  - Clean service interfaces without circular dependencies")
	logger.Info("  - Handler pipeline for composable request processing")
	logger.Info("  - Dependency injection with samber/do")
	logger.Info("  - Route discovery abstraction")
	logger.Info("  - Configuration loading abstraction")
	logger.Info("  - Testable and maintainable code structure")
}

// All service implementations have been moved to their respective packages
// and are now managed through dependency injection
