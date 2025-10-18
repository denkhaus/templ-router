package assets

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

//go:embed css/* fonts/* js/*
var Assets embed.FS

type assetsServiceImpl struct {
	assetRouteName string
	assetsPath     string
	configService  interfaces.ConfigService
}

func NewService(injector do.Injector) (interfaces.AssetsService, error) {
	configService := do.MustInvoke[interfaces.ConfigService](injector)
	logger := do.MustInvoke[*zap.Logger](injector)

	configAssetsDir := configService.GetLayoutAssetsDirectory()
	assetRouteName := configService.GetLayoutAssetsRouteName()

	assetsPath, err := filepath.Abs(configAssetsDir)
	if err != nil {
		return nil, fmt.Errorf("assets directory from config %q can't translate to an absolute path: %w", configAssetsDir, err)
	}

	info, err := os.Stat(assetsPath)
	if os.IsNotExist(err) || !info.IsDir() {
		return nil, fmt.Errorf("assets directory %q does not exist or is no directory", assetsPath)
	}

	logger.Info("assets service successfull initialized",
		zap.String("asset_route_name", assetRouteName),
		zap.String("assets_path", assetsPath),
	)

	return &assetsServiceImpl{
		assetRouteName: assetRouteName,
		configService:  configService,
		assetsPath:     assetsPath,
	}, nil
}

// createAssetHandler creates the asset handler with proper caching headers
func (p *assetsServiceImpl) createAssetHandler() http.HandlerFunc {

	isDevelopment := p.configService.IsDevelopment()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir(p.assetsPath))
		} else {
			fs = http.FileServer(http.FS(Assets))
		}

		fs.ServeHTTP(w, r)
	})
}

func (p *assetsServiceImpl) route() string {
	return fmt.Sprintf("/%s/*", p.assetRouteName)
}

func (p *assetsServiceImpl) prefix() string {
	return fmt.Sprintf("/%s/", p.assetRouteName)
}

func (p *assetsServiceImpl) SetupRoutes(mux *chi.Mux) {
	assetHandler := p.createAssetHandler()
	mux.Handle(p.route(), http.StripPrefix(p.prefix(), assetHandler))
}

func (p *assetsServiceImpl) SetupRoutesWithRouter(mux chi.Router) {
	assetHandler := p.createAssetHandler()
	// Handle assets under locale paths (e.g., /en/assets/*)
	mux.Handle(p.route(), http.StripPrefix(p.prefix(), assetHandler))
}
