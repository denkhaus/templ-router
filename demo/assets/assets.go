package assets

import (
	"embed"
	"net/http"

	"github.com/denkhaus/templ-router/demo/config"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
)

//go:embed css/* fonts/*
var Assets embed.FS

// createAssetHandler creates the asset handler with proper caching headers
func createAssetHandler(injector do.Injector) http.HandlerFunc {
	cfg := do.MustInvoke[*config.Config](injector)
	isDevelopment := cfg.IsDevelopment()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(Assets))
		}

		fs.ServeHTTP(w, r)
	})
}

func SetupRoutes(mux *chi.Mux, injector do.Injector) {
	assetHandler := createAssetHandler(injector)
	mux.Handle("/assets/*", http.StripPrefix("/assets/", assetHandler))
}

func SetupRoutesForLocale(mux chi.Router, injector do.Injector) {
	assetHandler := createAssetHandler(injector)
	// Handle assets under locale paths (e.g., /en/assets/*)
	mux.Handle("/assets/*", http.StripPrefix("/assets/", assetHandler))
}
