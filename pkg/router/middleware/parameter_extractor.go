package middleware

import (
	"net/http"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// ConfigurableParameterExtractor is a library-agnostic parameter extractor
// It uses configuration to determine how to extract parameters from URLs
type ConfigurableParameterExtractor struct {
	config        ParameterExtractionConfig
	logger        *zap.Logger
	configService interfaces.ConfigService
}

// ParameterExtractionConfig defines how parameters should be extracted from URLs
type ParameterExtractionConfig struct {
	// DynamicSegmentPatterns defines patterns for dynamic route segments
	// Key: pattern (e.g., "$id", "$slug"), Value: parameter name
	DynamicSegmentPatterns map[string]string

	// RoutePatterns defines specific route patterns and their parameter extraction rules
	// Key: route pattern (e.g., "/{locale}/{type}/{id}"), Value: extraction rules
	RoutePatterns map[string]RouteExtractionRule
}

// RouteExtractionRule defines how to extract parameters from a specific route pattern
type RouteExtractionRule struct {
	// SegmentMappings maps URL segment positions to parameter names
	// Key: segment index, Value: parameter name
	SegmentMappings map[int]string

	// RequiredSegments defines the minimum number of segments required
	RequiredSegments int

	// DynamicSegmentIndicators defines which segments are dynamic
	// Key: segment index, Value: true if dynamic
	DynamicSegmentIndicators map[int]bool
}

// NewConfigurableParameterExtractor creates a new parameter extractor for DI
func NewConfigurableParameterExtractor(i do.Injector) (ParameterExtractor, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	configService := do.MustInvoke[interfaces.ConfigService](i)

	// Get supported locales from config (not used in constructor, but available for validation)

	// Default configuration - can be overridden by config service
	defaultConfig := ParameterExtractionConfig{
		DynamicSegmentPatterns: map[string]string{
			"$id":     "id",
			"$slug":   "slug",
			"$locale": "locale",
		},
		RoutePatterns: map[string]RouteExtractionRule{
			// Generic pattern for /{locale}/{type}/{id} routes
			"generic_dynamic": {
				SegmentMappings: map[int]string{
					0: "locale",
					2: "id",
				},
				RequiredSegments: 3,
				DynamicSegmentIndicators: map[int]bool{
					2: true, // Third segment (index 2) is dynamic
				},
			},
		},
	}

	return &ConfigurableParameterExtractor{
		config:        defaultConfig,
		logger:        logger,
		configService: configService,
	}, nil
}

// ExtractParameters extracts parameters from URL path using configurable rules
func (cpe *ConfigurableParameterExtractor) ExtractParameters(urlPath string, route interfaces.Route) map[string]string {
	params := make(map[string]string)

	// Clean and split the URL path
	cleanPath := strings.TrimPrefix(urlPath, "/")
	if cleanPath == "" {
		return params
	}

	pathSegments := strings.Split(cleanPath, "/")

	cpe.logger.Debug("Extracting parameters from URL",
		zap.String("url_path", urlPath),
		zap.Strings("segments", pathSegments),
		zap.String("route_path", route.Path),
		zap.Bool("is_dynamic", route.IsDynamic))

	// Apply generic dynamic route pattern if route is marked as dynamic
	if route.IsDynamic {
		if rule, exists := cpe.config.RoutePatterns["generic_dynamic"]; exists {
			cpe.applyExtractionRule(pathSegments, rule, params)
		}
	}

	// Apply route-specific patterns based on route path
	for pattern, rule := range cpe.config.RoutePatterns {
		if pattern != "generic_dynamic" && cpe.matchesPattern(route.Path, pattern) {
			cpe.applyExtractionRule(pathSegments, rule, params)
		}
	}

	cpe.logger.Debug("Parameter extraction completed",
		zap.Any("extracted_params", params),
		zap.String("url_path", urlPath))

	return params
}

// ExtractParametersFromRequest extracts parameters from HTTP request using Chi's URL parameters
func (cpe *ConfigurableParameterExtractor) ExtractParametersFromRequest(r *http.Request, route interfaces.Route) map[string]string {
	params := make(map[string]string)

	// Extract all Chi URL parameters generically
	rctx := chi.RouteContext(r.Context())
	if rctx != nil {
		for i, key := range rctx.URLParams.Keys {
			if i < len(rctx.URLParams.Values) {
				value := rctx.URLParams.Values[i]
				if value != "" {
					params[key] = value
					cpe.logger.Debug("Extracted URL parameter from Chi",
						zap.String("key", key),
						zap.String("value", value),
						zap.String("url_path", r.URL.Path))
				}
			}
		}
	}

	// Extract locale from URL path (first segment)
	cleanPath := strings.TrimPrefix(r.URL.Path, "/")
	if cleanPath != "" {
		pathSegments := strings.Split(cleanPath, "/")
		if len(pathSegments) > 0 {
			firstSegment := pathSegments[0]
			// Check if first segment is a valid locale
			if cpe.isValidLocale(firstSegment) {
				params["locale"] = firstSegment
				cpe.logger.Debug("Extracted locale parameter",
					zap.String("locale", firstSegment),
					zap.String("url_path", r.URL.Path))
			}
		}
	}

	cpe.logger.Debug("Parameter extraction from request completed",
		zap.Any("extracted_params", params),
		zap.String("url_path", r.URL.Path),
		zap.String("route_path", route.Path))

	return params
}

// isValidLocale checks if a string is a valid locale code
func (cpe *ConfigurableParameterExtractor) isValidLocale(locale string) bool {
	// ConfigService must be properly injected - no fallbacks
	if cpe.configService == nil {
		cpe.logger.Error("ConfigService is nil - this is a DI configuration error")
		return false
	}

	// Get supported locales from config service
	supportedLocales := cpe.configService.GetSupportedLocales()
	for _, valid := range supportedLocales {
		if locale == valid {
			return true
		}
	}
	return false
}

// applyExtractionRule applies a specific extraction rule to path segments
func (cpe *ConfigurableParameterExtractor) applyExtractionRule(segments []string, rule RouteExtractionRule, params map[string]string) {
	if len(segments) < rule.RequiredSegments {
		cpe.logger.Debug("Insufficient segments for extraction rule",
			zap.Int("available", len(segments)),
			zap.Int("required", rule.RequiredSegments))
		return
	}

	for segmentIndex, paramName := range rule.SegmentMappings {
		if segmentIndex < len(segments) {
			// Check if this segment should be dynamic
			if isDynamic, exists := rule.DynamicSegmentIndicators[segmentIndex]; exists && isDynamic {
				params[paramName] = segments[segmentIndex]
				cpe.logger.Debug("Extracted dynamic parameter",
					zap.String("param_name", paramName),
					zap.String("param_value", segments[segmentIndex]),
					zap.Int("segment_index", segmentIndex))
			} else if !exists {
				// If not specified in DynamicSegmentIndicators, extract anyway
				params[paramName] = segments[segmentIndex]
				cpe.logger.Debug("Extracted parameter",
					zap.String("param_name", paramName),
					zap.String("param_value", segments[segmentIndex]),
					zap.Int("segment_index", segmentIndex))
			}
		}
	}
}

// matchesPattern checks if a route path matches a specific pattern
func (cpe *ConfigurableParameterExtractor) matchesPattern(routePath, pattern string) bool {
	// Simple pattern matching - can be enhanced with regex or more sophisticated matching
	return strings.Contains(routePath, pattern)
}
