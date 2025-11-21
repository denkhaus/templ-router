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
	demoservices "github.com/denkhaus/templ-router/demo/pkg/services"
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

	// 4. Register ALL application services using the streamlined options pattern
	// This replaces ALL the complex setup from the old version!
	container.RegisterApplicationServices(
		di.WithTemplateRegistry(templateRegistry),

		// Assets Service Factory - demonstrates pluggable asset management
		// This creates the assets service using the DI injector for dependencies
		di.WithAssetsServiceFactory(func(i do.Injector) (interfaces.AssetsService, error) {
			return assets.NewService(i)
		}),

		di.WithHealthCheck(true, "/api/health"),

		// NEW: Custom routes - add them directly without manual mux manipulation!
		di.WithCustomRoute("GET", "/api/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "operational", "version": "2.0"}`))
		}),

		// Additional custom routes to test multiple routes support
		di.WithCustomRoute("GET", "/api/health/detailed", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "healthy", "checks": ["database", "cache", "external_api"]}`))
		}),

		di.WithCustomRoute("POST", "/api/webhook", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "webhook received"}`))
		}),

		di.WithCustomRoute("GET", "/api/metrics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"requests": 1234, "errors": 0, "uptime": "99.9%"}`))
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

		// NEW: Custom Auth Validator - demonstrates hook-based authentication
		// This replaces the old AuthService with the new AuthValidator interface
		// The router middleware will use this to check authentication
		di.WithAuthValidatorFactory(func(i do.Injector) (interfaces.AuthValidator, error) {
			return demoservices.NewDemoAuthValidator(i)
		}),

	)

	// 5. Register Demo Authentication Services as client-side dependencies
	// These are now handled by the client application, not the router framework
	do.Provide(injector, demoservices.NewDefaultUserStore)
	do.Provide(injector, demoservices.NewDemoSessionStore)

	// Create auth handlers and register routes manually (client-side)
	authHandlers, err := demoservices.NewDemoAuthHandlers(injector)
	if err != nil {
		return shared.NewServiceError("Failed to create auth handlers").
			WithCause(err).
			WithContext("component", "auth_handlers")
	}

	// 6. Register DataServices as named dependencies (unchanged)
	do.ProvideNamed(injector, "UserDataService", dataservices.NewUserDataService)
	do.ProvideNamed(injector, "ProductDataService", dataservices.NewProductDataService)
	do.ProvideNamed(injector, "OrderDataService", dataservices.NewOrderDataService)
	do.ProvideNamed(injector, "BrokenDataService", dataservices.NewBrokenDataService)
	do.ProvideNamed(injector, "SpecificDataService", dataservices.NewSpecificOnlyDataService)
	do.ProvideNamed(injector, "UserWithIdDataService", dataservices.NewUserWithIdDataService)

	// 7. Get logger from container
	logger := container.GetLogger()
	defer logger.Sync()

	logger.Info("Starting application with streamlined bootstrap process")

	routerBootstrap := container.GetRouterBootstrap()
	mux, err := routerBootstrap.Bootstrap()
	if err != nil {
		return shared.NewServiceError("Failed to bootstrap router").
			WithCause(err).
			WithContext("component", "router_bootstrap")
	}

	// 8. Register client-side auth routes manually
	// Since auth is now client-side, we register the routes directly to the mux
	authHandlers.RegisterRoutes(func(method, path string, handler http.HandlerFunc) {
		switch method {
		case "GET":
			mux.Get(path, handler)
		case "POST":
			mux.Post(path, handler)
		case "PUT":
			mux.Put(path, handler)
		case "DELETE":
			mux.Delete(path, handler)
		case "PATCH":
			mux.Patch(path, handler)
		default:
			logger.Warn("Unsupported HTTP method for auth route",
				zap.String("method", method),
				zap.String("path", path))
		}
	})

	// 9. Log route information
	logRouteInformation(routerBootstrap.GetRouterCore(), logger)

	// 10. Start server - everything is already configured!
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
	// TODO: This would need to be adapted based on the actual interface available
	// For now, just log that routes were discovered
	logger.Info("Router bootstrap completed successfully",
		zap.String("architecture", "streamlined"),
		zap.String("bootstrap", "automatic"))
}
