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
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// main demonstrates the NEW streamlined bootstrap process
func main() {
	if err := startupStreamlined(context.Background()); err != nil {
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

func startupStreamlined(ctx context.Context) error {
	// 1. Create DI container - this is all the setup needed!
	container := di.NewContainer()
	defer container.Shutdown()

	// 2. Register router services with configuration prefix
	injector := container.RegisterRouterServices(ctx, "TR")

	// 3. Create application services
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

	// 4. Register ALL application services using the streamlined options pattern
	// This replaces ALL the complex setup from the old version!
	container.RegisterApplicationServices(
		di.WithTemplateRegistry(templateRegistry),
		di.WithAssetsService(assetsService),
		di.WithUserStore(userStore),

		// NEW: Health check configuration (auth routes controlled by env vars: TR_ROUTER_ENABLE_AUTH_ROUTES, TR_ROUTER_AUTH_ROUTE_PREFIX)
		di.WithHealthCheck(true, "/api/health"),

		// NEW: Custom routes - add them directly without manual mux manipulation!
		di.WithCustomRoute("GET", "/api/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "operational", "version": "2.0"}`))
		}),

		// NEW: Custom middleware - added to the chain in definition order!
		di.WithCustomMiddleware("request-id", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Request-ID", "custom-request-id")
				next.ServeHTTP(w, r)
			})
		}),

		// NEW: More custom middleware - will execute AFTER request-id middleware
		di.WithCustomMiddleware("timing", func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Timing-Middleware", "executed")
				next.ServeHTTP(w, r)
			})
		}),

		// NEW: Custom Session Store - demonstrates pluggable session management
	// By default, templ-router uses in-memory session store
	// Here we create our own simple in-memory session store to show the principle
	// In production, you might use Redis, database, or other custom implementations
	// Create a factory function that provides our custom session store
		di.WithSessionStoreFactory(func(i do.Injector) (interfaces.SessionStore, error) {
			logger := do.MustInvoke[*zap.Logger](i)
			configService := do.MustInvoke[interfaces.ConfigService](i)
			return services.NewCustomSessionStore(configService, logger)
		}),
	)

	// 5. Register DataServices as named dependencies (unchanged)
	do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
	do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
	do.ProvideNamed(injector, "OrderDataService", dataservices.NewOrderDataService)
	do.ProvideNamed(injector, "BrokenDataService", dataservices.NewBrokenDataService)
	do.ProvideNamed(injector, "SpecificDataService", dataservices.NewSpecificOnlyDataService)
	do.ProvideNamed(injector, "UserWithIdDataService", dataservices.NewUserWithIdDataService)

	// 6. Get logger from container
	logger := container.GetLogger()
	defer logger.Sync()

	logger.Info("Starting application with streamlined bootstrap process")

	// 7. THE MAGIC: Use RouterBootstrap to automatically configure everything!
	// No more: cleanRouter.GetMiddlewareSetup().GetRouterMiddleware().Configure(mux)
	// No more: manual middleware creation and ordering
	// No more: manual route registration
	routerBootstrap := container.GetRouterBootstrap()
	mux, err := routerBootstrap.Bootstrap()
	if err != nil {
		return shared.NewServiceError("Failed to bootstrap router").
			WithCause(err).
			WithContext("component", "router_bootstrap")
	}

	// 8. Log route information
	logRouteInformation(routerBootstrap.GetRouterCore(), logger)

	// 9. Start server - everything is already configured!
	logger.Info("Starting Streamlined Bootstrap Demo Server on 0.0.0.0:8084")
	if err := http.ListenAndServe("0.0.0.0:8084", mux); err != nil {
		return shared.NewServiceError("Failed to start HTTP server").
			WithCause(err).
			WithContext("component", "http_server").
			WithContext("address", "0.0.0.0:8084")
	}

	return nil
}

// logRouteInformation logs information about discovered routes
func logRouteInformation(routerCore interface{}, logger *zap.Logger) {
	// This would need to be adapted based on the actual interface available
	// For now, just log that routes were discovered
	logger.Info("Router bootstrap completed successfully",
		zap.String("architecture", "streamlined"),
		zap.String("bootstrap", "automatic"))
}
