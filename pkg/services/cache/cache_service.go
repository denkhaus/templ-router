package cache

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// cacheService implements the CacheService interface
type cacheService struct {
	templateCache sync.Map
	routeCache    sync.Map
	logger        *zap.Logger
}

// NewCacheService creates a new cache service instance
func NewCacheService(i do.Injector) (interfaces.CacheService, error) {
	logger := do.MustInvoke[*zap.Logger](i)

	service := &cacheService{
		templateCache: sync.Map{},
		routeCache:    sync.Map{},
		logger:        logger,
	}

	logger.Info("Cache service initialized")
	return service, nil
}

// GetTemplate retrieves a template from cache
func (cs *cacheService) GetTemplate(key string) (interface{}, bool) {
	value, exists := cs.templateCache.Load(key)
	if exists {
		cs.logger.Debug("Template cache hit", zap.String("key", key))
	} else {
		cs.logger.Debug("Template cache miss", zap.String("key", key))
	}
	return value, exists
}

// SetTemplate stores a template in cache
func (cs *cacheService) SetTemplate(key string, value interface{}) {
	cs.templateCache.Store(key, value)
	cs.logger.Debug("Template cached", zap.String("key", key))
}

// GetRoute retrieves a route from cache
func (cs *cacheService) GetRoute(key string) (interface{}, bool) {
	value, exists := cs.routeCache.Load(key)
	if exists {
		cs.logger.Debug("Route cache hit", zap.String("key", key))
	} else {
		cs.logger.Debug("Route cache miss", zap.String("key", key))
	}
	return value, exists
}

// SetRoute stores a route in cache
func (cs *cacheService) SetRoute(key string, value interface{}) {
	cs.routeCache.Store(key, value)
	cs.logger.Debug("Route cached", zap.String("key", key))
}

// ClearAll clears both template and route caches
func (cs *cacheService) ClearAll() {
	cs.templateCache.Range(func(key, value interface{}) bool {
		cs.templateCache.Delete(key)
		return true
	})
	cs.routeCache.Range(func(key, value interface{}) bool {
		cs.routeCache.Delete(key)
		return true
	})
	cs.logger.Info("All caches cleared")
}

// ClearTemplates clears only the template cache
func (cs *cacheService) ClearTemplates() {
	cs.templateCache.Range(func(key, value interface{}) bool {
		cs.templateCache.Delete(key)
		return true
	})
	cs.logger.Info("Template cache cleared")
}

// ClearRoutes clears only the route cache
func (cs *cacheService) ClearRoutes() {
	cs.routeCache.Range(func(key, value interface{}) bool {
		cs.routeCache.Delete(key)
		return true
	})
	cs.logger.Info("Route cache cleared")
}

// BuildTemplateKey generates a cache key for templates
func (cs *cacheService) BuildTemplateKey(templateFile, locale string, params map[string]string) string {
	if len(params) == 0 {
		return fmt.Sprintf("template:%s:%s", templateFile, locale)
	}

	// Sort parameters for consistent key generation
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramParts []string
	for _, k := range keys {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s", k, params[k]))
	}

	return fmt.Sprintf("template:%s:%s:%s", templateFile, locale, strings.Join(paramParts, "&"))
}

// BuildRouteKey generates a cache key for routes
func (cs *cacheService) BuildRouteKey(path string, params map[string]string) string {
	if len(params) == 0 {
		return fmt.Sprintf("route:%s", path)
	}

	// Sort parameters for consistent key generation
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramParts []string
	for _, k := range keys {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s", k, params[k]))
	}

	return fmt.Sprintf("route:%s:%s", path, strings.Join(paramParts, "&"))
}