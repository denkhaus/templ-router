package services

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// OptimizedTemplateService consolidates all template resolution systems
// into a single, performance-optimized service with caching
type OptimizedTemplateService struct {
	logger *zap.Logger

	// Template registry interface for decoupled access
	templateRegistry interfaces.TemplateRegistry

	// Cache service for performance optimization
	cacheService interfaces.CacheService

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
	cacheService := do.MustInvoke[interfaces.CacheService](i)

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
		cacheService:     cacheService,
	}, nil
}

// RenderComponent implements interfaces.TemplateService with optimized resolution
func (ots *OptimizedTemplateService) RenderComponent(route interfaces.Route, routerCtx interfaces.RouterContext, ctx context.Context) (templ.Component, error) {
	routePath := route.Path
	
	// Extract parameters from RouterContext for backward compatibility
	allParams := make(map[string]string)
	// Add URL parameters
	for k, v := range routerCtx.GetAllURLParams() {
		allParams[k] = v
	}
	// Add query parameters with "query_" prefix to avoid conflicts
	for k, values := range routerCtx.GetAllQueryParams() {
		if len(values) > 0 {
			allParams["query_"+k] = values[0]
		}
	}

	ots.logger.Debug("Optimized template service rendering component",
		zap.String("route", routePath),
		zap.String("template_file", route.TemplateFile),
		zap.Any("url_params", routerCtx.GetAllURLParams()),
		zap.Any("query_params", routerCtx.GetAllQueryParams()),
		zap.Any("combined_params", allParams))

	// PERFORMANCE: Check cache first - include parameters in cache key for dynamic templates
	cacheKey := ots.cacheService.BuildTemplateKey(route.TemplateFile, "", allParams)
	if cached, found := ots.cacheService.GetTemplate(cacheKey); found {
		if component, ok := cached.(templ.Component); ok {
			ots.logger.Debug("Template served from cache",
				zap.String("cache_key", cacheKey))
			return component, nil
		}
	}

	// UNIFIED RESOLUTION: Single resolution strategy
	component, err := ots.resolveTemplate(routePath, routerCtx)
	if err != nil {
		return nil, err
	}

	// PERFORMANCE: Cache successful resolution
	ots.cacheService.SetTemplate(cacheKey, component)

	return component, nil
}

// RenderLayoutComponent implements interfaces.TemplateService with layout optimization
func (ots *OptimizedTemplateService) RenderLayoutComponent(layoutPath string, content templ.Component, ctx context.Context) (templ.Component, error) {
	ots.logger.Debug("Optimized layout rendering",
		zap.String("layout_path", layoutPath))

	// PERFORMANCE: Check layout cache
	layoutCacheKey := ots.cacheService.BuildTemplateKey(layoutPath, "", nil)
	if cached, found := ots.cacheService.GetTemplate(layoutCacheKey); found {
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
				ots.cacheService.SetTemplate(layoutCacheKey, layoutFunc)

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
func (ots *OptimizedTemplateService) resolveTemplate(routePath string, routerCtx interfaces.RouterContext) (templ.Component, error) {
	// Extract combined parameters for caching
	allParams := make(map[string]string)
	for k, v := range routerCtx.GetAllURLParams() {
		allParams[k] = v
	}
	for k, values := range routerCtx.GetAllQueryParams() {
		if len(values) > 0 {
			allParams["query_"+k] = values[0]
		}
	}

	// PERFORMANCE: Check route cache first
	routeCacheKey := ots.cacheService.BuildRouteKey(routePath, allParams)
	if templateUUID, found := ots.cacheService.GetRoute(routeCacheKey); found {
		if uuid, ok := templateUUID.(string); ok {
			if templateFunc, exists := ots.templateRegistry.GetTemplateFunction(uuid); exists {
				return ots.executeTemplateFunction(templateFunc, routerCtx, routePath, uuid)
			}
		}
	}

	// Strategy 1: Direct route lookup
	routeMapping := ots.templateRegistry.GetRouteToTemplateMapping()
	if templateUUID, exists := routeMapping[routePath]; exists {
		// PERFORMANCE: Cache successful route mapping
		routeCacheKey := ots.cacheService.BuildRouteKey(routePath, allParams)
		ots.cacheService.SetRoute(routeCacheKey, templateUUID)

		if templateFunc, found := ots.templateRegistry.GetTemplateFunction(templateUUID); found {
			return ots.executeTemplateFunction(templateFunc, routerCtx, routePath, templateUUID)
		}
	}

	// Strategy 2: Dynamic route resolution using RouteConverter
	convertedRoutes := ots.routeConverter.GenerateRouteVariations(routePath)
	for _, convertedRoute := range convertedRoutes {
		if templateUUID, exists := routeMapping[convertedRoute]; exists {
			// PERFORMANCE: Cache successful conversion
			routeCacheKey := ots.cacheService.BuildRouteKey(routePath, allParams)
			ots.cacheService.SetRoute(routeCacheKey, templateUUID)

			if templateFunc, found := ots.templateRegistry.GetTemplateFunction(templateUUID); found {
				ots.logger.Debug("Template resolved via route conversion",
					zap.String("original_route", routePath),
					zap.String("converted_route", convertedRoute),
					zap.String("template_uuid", templateUUID))

				return ots.executeTemplateFunction(templateFunc, routerCtx, routePath, templateUUID)
			}
		}
	}

	ots.logger.Error("Template resolution failed for route",
		zap.String("route", routePath),
		zap.Strings("attempted_conversions", convertedRoutes))

	return nil, middleware.ErrTemplateNotFound
}

// executeTemplateFunction handles different template function signatures
func (ots *OptimizedTemplateService) executeTemplateFunction(templateFunc func() interface{}, routerCtx interfaces.RouterContext, routePath, templateUUID string) (templ.Component, error) {
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
		id := routerCtx.GetURLParam("id")
		if id == "" {
			return nil, shared.NewValidationError("parameter 'id' is required for parameterized template").
				WithContext("route", routePath).
				WithContext("template_uuid", templateUUID)
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
		return ots.executeDataServiceTemplate(result, routerCtx, routePath, templateUUID)
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
		err := shared.NewValidationError("layoutPath cannot be empty").
			WithDetails("invalid template path provided")
		ots.logger.Error("Layout path validation failed", zap.Error(err))
		return ""
	}

	// Library-agnostic conversion: any/path/layout.templ -> /layout
	filename := filepath.Base(layoutPath)
	if !strings.HasSuffix(filename, ".templ") {
		err := shared.NewValidationError("invalid template file extension").
			WithDetails("must have .templ extension").
			WithContext("layout_path", layoutPath).
			WithContext("filename", filename)
		ots.logger.Error("Template file validation failed", zap.Error(err))
		return ""
	}

	routeName := strings.TrimSuffix(filename, ".templ")
	if routeName == "" {
		err := shared.NewValidationError("cannot extract route name from template filename").
			WithDetails("invalid template filename structure").
			WithContext("layout_path", layoutPath).
			WithContext("filename", filename)
		ots.logger.Error("Route name extraction failed", zap.Error(err))
		return ""
	}

	return "/" + routeName
}

// executeDataServiceTemplate handles DataService template execution with optimized method calls
func (ots *OptimizedTemplateService) executeDataServiceTemplate(templateFunc interface{}, routerCtx interfaces.RouterContext, routePath, templateUUID string) (templ.Component, error) {
	// Get DataService info from template registry
	dataServiceInfo, exists := ots.templateRegistry.GetDataServiceInfo(templateUUID)
	if !exists {
		return nil, shared.NewServiceError("template requires data service but no info found").
			WithDetails("DataService information missing from template registry").
			WithContext("template_uuid", templateUUID).
			WithContext("route", routePath)
	}

	ots.logger.Debug("Resolving DataService",
		zap.String("route", routePath),
		zap.String("data_service_interface", dataServiceInfo.InterfaceType))

	// OPTIMIZATION: Use GenericDataService interface (no reflection for DataService calls)
	genericDataService, err := ots.dataResolver.ResolveGenericDataService(dataServiceInfo.InterfaceType)
	if err != nil {
		return nil, shared.NewDependencyInjectionError("failed to resolve generic data service").
			WithDetails("DataService not found or cannot be wrapped as GenericDataService").
			WithCause(err).
			WithContext("interface_type", dataServiceInfo.InterfaceType).
			WithContext("route", routePath).
			WithContext("template_uuid", templateUUID)
	}

	ots.logger.Debug("Using GenericDataService interface (optimized, no reflection for service call)",
		zap.String("interface_type", dataServiceInfo.InterfaceType))

	// Call GetData directly on the generic interface (no reflection!)
	data, err := genericDataService.GetData(routerCtx)
	if err != nil {
		return nil, err
	}

	// Execute template function with reflection (only remaining reflection usage)
	ots.logger.Debug("Executing template function with reflection",
		zap.String("template_uuid", templateUUID))
	
	return ots.executeTemplateWithReflection(templateFunc, data, routePath, templateUUID)
}

// executeTemplateWithReflection executes template function using reflection (fallback)
func (ots *OptimizedTemplateService) executeTemplateWithReflection(templateFunc interface{}, data interface{}, routePath, templateUUID string) (templ.Component, error) {
	// Call template function with data using reflection
	funcValue := reflect.ValueOf(templateFunc)
	
	if funcValue.Kind() != reflect.Func {
		return nil, shared.NewTemplateError("invalid template function type").
			WithDetails("Template must be a function").
			WithContext("actual_type", funcValue.Kind().String()).
			WithContext("route", routePath).
			WithContext("template_uuid", templateUUID)
	}

	// Prepare arguments: data
	templateArgs := []reflect.Value{
		reflect.ValueOf(data),
	}

	// Call the function
	templateResults := funcValue.Call(templateArgs)
	if len(templateResults) != 1 {
		return nil, shared.NewTemplateError("invalid template function signature").
			WithDetails("Template function should return exactly one value").
			WithContext("result_count", len(templateResults)).
			WithContext("route", routePath).
			WithContext("template_uuid", templateUUID)
	}

	// Convert result to templ.Component
	component, ok := templateResults[0].Interface().(templ.Component)
	if !ok {
		return nil, shared.NewTemplateError("invalid template function return type").
			WithDetails("Template function must return templ.Component").
			WithContext("actual_type", fmt.Sprintf("%T", templateResults[0].Interface())).
			WithContext("route", routePath).
			WithContext("template_uuid", templateUUID)
	}

	ots.logger.Debug("Template executed successfully with reflection fallback",
		zap.String("route", routePath),
		zap.String("template_uuid", templateUUID))

	return component, nil
}

// executeDataServiceTemplateWithReflection executes DataService template using reflection (fallback)
func (ots *OptimizedTemplateService) executeDataServiceTemplateWithReflection(templateFunc interface{}, params map[string]string, routePath, templateUUID string, dataServiceInfo interfaces.DataServiceInfo) (templ.Component, error) {
	// Resolve data service from DI
	dataService, err := ots.dataResolver.ResolveDataService(dataServiceInfo.InterfaceType)
	if err != nil {
		return nil, shared.NewDependencyInjectionError("failed to resolve data service").
			WithDetails("DataService not found in DI container").
			WithCause(err).
			WithContext("interface_type", dataServiceInfo.InterfaceType).
			WithContext("route", routePath).
			WithContext("template_uuid", templateUUID)
	}

	// Call specific method based on data type, fallback to GetData
	serviceValue := reflect.ValueOf(dataService)
	methodName := shared.DeriveMethodNameFromDataType(dataServiceInfo.ParameterType)
	getDataMethod := serviceValue.MethodByName(methodName)
	
	// Fallback to GetData if specific method doesn't exist
	if !getDataMethod.IsValid() {
		getDataMethod = serviceValue.MethodByName("GetData")
		methodName = "GetData"
	}
	
	if !getDataMethod.IsValid() {
		return nil, shared.NewServiceError("data service method not found").
			WithDetails("Neither specific method nor GetData method exists on data service").
			WithContext("method_name", methodName).
			WithContext("interface_type", dataServiceInfo.InterfaceType).
			WithContext("route", routePath).
			WithContext("template_uuid", templateUUID)
	}
	
	ots.logger.Debug("Using data service method with reflection fallback",
		zap.String("method_name", methodName),
		zap.String("data_service_interface", dataServiceInfo.InterfaceType))

	// Prepare arguments: ctx, params
	args := []reflect.Value{
		reflect.ValueOf(context.Background()),
		reflect.ValueOf(params),
	}

	// Call the method
	results := getDataMethod.Call(args)
	if len(results) != 2 {
		return nil, shared.NewServiceError("invalid data service method signature").
			WithDetails("GetData method should return (data, error)").
			WithContext("method_name", methodName).
			WithContext("result_count", len(results)).
			WithContext("interface_type", dataServiceInfo.InterfaceType).
			WithContext("route", routePath)
	}

	// Check for error
	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return nil, err
	}

	data := results[0].Interface()

	// Execute template with reflection
	return ots.executeTemplateWithReflection(templateFunc, data, routePath, templateUUID)
}



// ClearCache clears the template cache (useful for development)
func (ots *OptimizedTemplateService) ClearCache() {
	ots.cacheService.ClearAll()
	ots.logger.Info("Template cache cleared")
}
