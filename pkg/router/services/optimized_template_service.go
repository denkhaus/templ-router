package services

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// OptimizedTemplateService consolidates all template resolution systems
// into a single, performance-optimized service with caching
type OptimizedTemplateService struct {
	logger *zap.Logger

	// Template registry interface for decoupled access
	templateRegistry interfaces.TemplateRegistry

	// Performance optimization: Template cache
	templateCache sync.Map // map[string]templ.Component
	routeCache    sync.Map // map[string]string (route -> templateUUID)

	// Route converter for dynamic route handling
	routeConverter RouteConverter
	
	// DataService resolver for DataService templates
	dataResolver interfaces.DataServiceResolver
}

// NewOptimizedTemplateService creates the unified template service
func NewOptimizedTemplateService(i do.Injector) (interfaces.TemplateService, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	dataResolver := do.MustInvoke[interfaces.DataServiceResolver](i)

	// Create route converter for dynamic route handling
	routeConverter, err := NewRouteConverter(i)
	if err != nil {
		return nil, err
	}

	return &OptimizedTemplateService{
		logger:           logger,
		templateRegistry: templateRegistry,
		routeConverter:   routeConverter,
		dataResolver:     dataResolver,
	}, nil
}

// RenderComponent implements interfaces.TemplateService with optimized resolution
func (ots *OptimizedTemplateService) RenderComponent(route interfaces.Route, ctx context.Context, params map[string]string) (templ.Component, error) {
	routePath := route.Path

	ots.logger.Debug("Optimized template service rendering component",
		zap.String("route", routePath),
		zap.String("template_file", route.TemplateFile),
		zap.Any("params", params))

	// PERFORMANCE: Check cache first
	cacheKey := routePath + "|" + route.TemplateFile
	if cached, found := ots.templateCache.Load(cacheKey); found {
		if component, ok := cached.(templ.Component); ok {
			ots.logger.Debug("Template served from cache",
				zap.String("cache_key", cacheKey))
			return component, nil
		}
	}

	// UNIFIED RESOLUTION: Single resolution strategy
	component, err := ots.resolveTemplate(routePath, params)
	if err != nil {
		return nil, err
	}

	// PERFORMANCE: Cache successful resolution
	ots.templateCache.Store(cacheKey, component)

	return component, nil
}

// RenderLayoutComponent implements interfaces.TemplateService with layout optimization
func (ots *OptimizedTemplateService) RenderLayoutComponent(layoutPath string, content templ.Component, ctx context.Context) (templ.Component, error) {
	ots.logger.Debug("Optimized layout rendering",
		zap.String("layout_path", layoutPath))

	// PERFORMANCE: Check layout cache
	layoutCacheKey := "layout:" + layoutPath
	if cached, found := ots.templateCache.Load(layoutCacheKey); found {
		if layoutFunc, ok := cached.(func(templ.Component) templ.Component); ok {
			ots.logger.Debug("Layout function served from cache",
				zap.String("layout_path", layoutPath))
			return layoutFunc(content), nil
		}
	}

	// Convert layout path to route pattern
	layoutRoute := ots.convertLayoutPathToRoute(layoutPath)

	// Look up layout template using template registry
	routeMapping := ots.templateRegistry.GetRouteToTemplateMapping()
	if templateUUID, exists := routeMapping[layoutRoute]; exists {
		if templateFunc, found := ots.templateRegistry.GetTemplateFunction(templateUUID); found {
			layoutFuncResult := templateFunc()

			// Check if it's a layout function
			if layoutFunc, ok := layoutFuncResult.(func(templ.Component) templ.Component); ok {
				ots.logger.Info("Layout function resolved and cached",
					zap.String("layout_route", layoutRoute),
					zap.String("template_uuid", templateUUID))

				// PERFORMANCE: Cache layout function
				ots.templateCache.Store(layoutCacheKey, layoutFunc)

				// Create layout component with content
				return layoutFunc(content), nil
			}
		}
	}

	ots.logger.Warn("Layout rendering failed, returning content without layout",
		zap.String("layout_path", layoutPath))
	return content, nil
}

// resolveTemplate - UNIFIED resolution strategy combining all previous approaches
func (ots *OptimizedTemplateService) resolveTemplate(routePath string, params map[string]string) (templ.Component, error) {
	// PERFORMANCE: Check route cache first
	if templateUUID, found := ots.routeCache.Load(routePath); found {
		if uuid, ok := templateUUID.(string); ok {
			if templateFunc, exists := ots.templateRegistry.GetTemplateFunction(uuid); exists {
				return ots.executeTemplateFunction(templateFunc, params, routePath, uuid)
			}
		}
	}

	// Strategy 1: Direct route lookup
	routeMapping := ots.templateRegistry.GetRouteToTemplateMapping()
	if templateUUID, exists := routeMapping[routePath]; exists {
		// PERFORMANCE: Cache successful route mapping
		ots.routeCache.Store(routePath, templateUUID)

		if templateFunc, found := ots.templateRegistry.GetTemplateFunction(templateUUID); found {
			return ots.executeTemplateFunction(templateFunc, params, routePath, templateUUID)
		}
	}

	// Strategy 2: Dynamic route resolution using RouteConverter
	convertedRoutes := ots.routeConverter.GenerateRouteVariations(routePath)
	for _, convertedRoute := range convertedRoutes {
		if templateUUID, exists := routeMapping[convertedRoute]; exists {
			// PERFORMANCE: Cache successful conversion
			ots.routeCache.Store(routePath, templateUUID)

			if templateFunc, found := ots.templateRegistry.GetTemplateFunction(templateUUID); found {
				ots.logger.Debug("Template resolved via route conversion",
					zap.String("original_route", routePath),
					zap.String("converted_route", convertedRoute),
					zap.String("template_uuid", templateUUID))

				return ots.executeTemplateFunction(templateFunc, params, routePath, templateUUID)
			}
		}
	}

	ots.logger.Error("Template resolution failed for route",
		zap.String("route", routePath),
		zap.Strings("attempted_conversions", convertedRoutes))

	return nil, middleware.ErrTemplateNotFound
}

// executeTemplateFunction handles different template function signatures
func (ots *OptimizedTemplateService) executeTemplateFunction(templateFunc func() interface{}, params map[string]string, routePath, templateUUID string) (templ.Component, error) {
	result := templateFunc()

	// Handle parameterless template functions (most common case)
	if fn, ok := result.(func() templ.Component); ok {
		component := fn()
		ots.logger.Debug("Parameterless template function executed",
			zap.String("route", routePath),
			zap.String("template_uuid", templateUUID))
		return component, nil
	}

	// Handle parameterized templates (e.g., user/product pages)
	if fn, ok := result.(func(string) templ.Component); ok {
		id := params["id"]
		if id == "" {
			panic("OptimizedTemplateService: parameter 'id' is empty - template requires valid parameter")
		}
		component := fn(id)
		ots.logger.Debug("Parameterized template executed",
			zap.String("route", routePath),
			zap.String("template_uuid", templateUUID),
			zap.String("id", id))
		return component, nil
	}

	// Handle DataService templates (e.g., func(*dataservices.UserData) templ.Component)
	if ots.isDataServiceTemplate(result) {
		ots.logger.Debug("DataService template detected - resolving data",
			zap.String("route", routePath),
			zap.String("template_uuid", templateUUID),
			zap.String("result_type", fmt.Sprintf("%T", result)))
		
		// Resolve DataService template directly in TemplateService
		return ots.executeDataServiceTemplate(result, params, routePath, templateUUID)
	}

	// Handle direct components (fallback)
	if component, ok := result.(templ.Component); ok {
		ots.logger.Debug("Direct component executed",
			zap.String("route", routePath),
			zap.String("template_uuid", templateUUID))
		return component, nil
	}

	ots.logger.Error("Unknown template function signature",
		zap.String("route", routePath),
		zap.String("template_uuid", templateUUID),
		zap.String("result_type", fmt.Sprintf("%T", result)))

	return nil, middleware.ErrTemplateNotFound
}

// isDataServiceTemplate checks if the result is a DataService template function
// DataService templates have signature: func(*SomeDataType) templ.Component
func (ots *OptimizedTemplateService) isDataServiceTemplate(result interface{}) bool {
	resultType := reflect.TypeOf(result)
	
	// Must be a function
	if resultType.Kind() != reflect.Func {
		return false
	}
	
	// Must have exactly 1 input parameter and 1 output parameter
	if resultType.NumIn() != 1 || resultType.NumOut() != 1 {
		return false
	}
	
	// Input parameter must be a pointer to a struct (DataService data type)
	inputType := resultType.In(0)
	if inputType.Kind() != reflect.Ptr {
		return false
	}
	
	// The pointer must point to a struct
	if inputType.Elem().Kind() != reflect.Struct {
		return false
	}
	
	// Output must be templ.Component
	outputType := resultType.Out(0)
	// Check if output implements templ.Component interface
	templComponentType := reflect.TypeOf((*templ.Component)(nil)).Elem()
	if !outputType.Implements(templComponentType) {
		return false
	}
	
	return true
}

// convertLayoutPathToRoute converts layout path to route pattern (fail-fast)
func (ots *OptimizedTemplateService) convertLayoutPathToRoute(layoutPath string) string {
	if layoutPath == "" {
		panic("OptimizedTemplateService: layoutPath is empty - invalid template path provided")
	}

	// Library-agnostic conversion: any/path/layout.templ -> /layout
	filename := filepath.Base(layoutPath)
	if !strings.HasSuffix(filename, ".templ") {
		panic(fmt.Sprintf("OptimizedTemplateService: invalid template file '%s' - must have .templ extension", layoutPath))
	}

	routeName := strings.TrimSuffix(filename, ".templ")
	if routeName == "" {
		panic(fmt.Sprintf("OptimizedTemplateService: invalid template filename '%s' - cannot extract route name", layoutPath))
	}

	return "/" + routeName
}

// executeDataServiceTemplate handles DataService template execution
func (ots *OptimizedTemplateService) executeDataServiceTemplate(templateFunc interface{}, params map[string]string, routePath, templateUUID string) (templ.Component, error) {
	// Get DataService info from template registry
	dataServiceInfo, exists := ots.templateRegistry.GetDataServiceInfo(templateUUID)
	if !exists {
		return nil, fmt.Errorf("template %s requires data service but no info found", templateUUID)
	}

	ots.logger.Debug("Resolving DataService",
		zap.String("route", routePath),
		zap.String("data_service_interface", dataServiceInfo.InterfaceType))

	// Resolve data service from DI
	dataService, err := ots.dataResolver.ResolveDataService(dataServiceInfo.InterfaceType)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve data service %s: %w", dataServiceInfo.InterfaceType, err)
	}

	// Call GetData method on the service using reflection
	serviceValue := reflect.ValueOf(dataService)
	getDataMethod := serviceValue.MethodByName("GetData")
	
	if !getDataMethod.IsValid() {
		return nil, fmt.Errorf("GetData method not found on data service")
	}

	// Prepare arguments: ctx (empty for now), params
	args := []reflect.Value{
		reflect.ValueOf(context.Background()),
		reflect.ValueOf(params),
	}

	// Call the method
	results := getDataMethod.Call(args)
	if len(results) != 2 {
		return nil, fmt.Errorf("GetData method should return (data, error)")
	}

	// Check for error
	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return nil, err
	}

	data := results[0].Interface()

	// Call template function with data using reflection
	funcValue := reflect.ValueOf(templateFunc)
	
	if funcValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("template is not a function")
	}

	// Prepare arguments: data
	templateArgs := []reflect.Value{
		reflect.ValueOf(data),
	}

	// Call the function
	templateResults := funcValue.Call(templateArgs)
	if len(templateResults) != 1 {
		return nil, fmt.Errorf("template function should return one value")
	}

	// Convert result to templ.Component
	component, ok := templateResults[0].Interface().(templ.Component)
	if !ok {
		return nil, fmt.Errorf("template function did not return templ.Component")
	}

	ots.logger.Debug("DataService template executed successfully",
		zap.String("route", routePath),
		zap.String("template_uuid", templateUUID))

	return component, nil
}

// ClearCache clears the template cache (useful for development)
func (ots *OptimizedTemplateService) ClearCache() {
	ots.templateCache = sync.Map{}
	ots.routeCache = sync.Map{}
	ots.logger.Info("Template cache cleared")
}
