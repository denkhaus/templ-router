package middleware

import (
	"context"
	"net/http"
	"net/url"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
)

// routerContext is the concrete implementation of RouterContext
type routerContext struct {
	ctx        context.Context
	request    *http.Request
	chiContext *chi.Context
	
	// Cached parameter maps for efficiency
	urlParams   map[string]string
	queryParams url.Values
}

// NewRouterContext creates a new RouterContext instance
func NewRouterContext(ctx context.Context, req *http.Request) interfaces.RouterContext {
	chiCtx := chi.RouteContext(ctx)
	
	return &routerContext{
		ctx:        ctx,
		request:    req,
		chiContext: chiCtx,
		urlParams:  extractURLParams(chiCtx),
		queryParams: req.URL.Query(),
	}
}

// Context returns the underlying context.Context
func (rc *routerContext) Context() context.Context {
	return rc.ctx
}

// GetURLParam returns a single URL parameter value
func (rc *routerContext) GetURLParam(key string) string {
	if rc.urlParams == nil {
		return ""
	}
	return rc.urlParams[key]
}

// GetAllURLParams returns all URL parameters as a map
func (rc *routerContext) GetAllURLParams() map[string]string {
	if rc.urlParams == nil {
		return make(map[string]string)
	}
	// Return a copy to prevent external modification
	result := make(map[string]string, len(rc.urlParams))
	for k, v := range rc.urlParams {
		result[k] = v
	}
	return result
}

// GetQueryParam returns the first value for a query parameter
func (rc *routerContext) GetQueryParam(key string) string {
	if rc.queryParams == nil {
		return ""
	}
	return rc.queryParams.Get(key)
}

// GetQueryParams returns all values for a query parameter
func (rc *routerContext) GetQueryParams(key string) []string {
	if rc.queryParams == nil {
		return nil
	}
	return rc.queryParams[key]
}

// GetAllQueryParams returns all query parameters
func (rc *routerContext) GetAllQueryParams() url.Values {
	if rc.queryParams == nil {
		return make(url.Values)
	}
	// Return a copy to prevent external modification
	result := make(url.Values, len(rc.queryParams))
	for k, v := range rc.queryParams {
		result[k] = make([]string, len(v))
		copy(result[k], v)
	}
	return result
}

// Request returns the original HTTP request
func (rc *routerContext) Request() *http.Request {
	return rc.request
}

// ChiContext returns the Chi router context
func (rc *routerContext) ChiContext() *chi.Context {
	return rc.chiContext
}

// extractURLParams extracts URL parameters from Chi context
func extractURLParams(chiCtx *chi.Context) map[string]string {
	if chiCtx == nil {
		return make(map[string]string)
	}
	
	params := make(map[string]string)
	for i, key := range chiCtx.URLParams.Keys {
		if i < len(chiCtx.URLParams.Values) {
			value := chiCtx.URLParams.Values[i]
			if value != "" {
				params[key] = value
			}
		}
	}
	return params
}