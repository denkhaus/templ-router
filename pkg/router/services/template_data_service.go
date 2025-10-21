package services

import (
	"context"
	"fmt"
	"reflect"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)


// templateDataServiceImpl implements TemplateDataService
type templateDataServiceImpl struct {
	templateRegistry interfaces.TemplateRegistry
	dataResolver     interfaces.DataServiceResolver
}

// NewTemplateDataService creates a new template data service for DI
func NewTemplateDataService(i do.Injector) (interfaces.TemplateDataService, error) {
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	dataResolver := do.MustInvoke[interfaces.DataServiceResolver](i)
	
	return &templateDataServiceImpl{
		templateRegistry: templateRegistry,
		dataResolver:     dataResolver,
	}, nil
}

// ResolveTemplateWithData resolves template data and renders the template
func (s *templateDataServiceImpl) ResolveTemplateWithData(
	ctx context.Context,
	templateKey string,
	routeParams map[string]string,
) (templ.Component, error) {
	// Check if template requires data service
	if !s.templateRegistry.RequiresDataService(templateKey) {
		// No data service required, use normal template resolution
		return s.templateRegistry.GetTemplate(templateKey)
	}

	// Get data service info
	dataServiceInfo, exists := s.templateRegistry.GetDataServiceInfo(templateKey)
	if !exists {
		return nil, fmt.Errorf("template %s requires data service but no info found", templateKey)
	}

	// Resolve data service from DI
	dataService, err := s.dataResolver.ResolveDataService(dataServiceInfo.InterfaceType)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve data service %s: %w", dataServiceInfo.InterfaceType, err)
	}

	// Call GetData method on the service
	data, err := s.callGetDataMethod(ctx, dataService, routeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from service: %w", err)
	}

	// Get template function and call it with data
	templateFunc, exists := s.templateRegistry.GetTemplateFunction(templateKey)
	if !exists {
		return nil, fmt.Errorf("template function not found for key: %s", templateKey)
	}

	// Call template function with data
	component, err := s.callTemplateWithData(templateFunc(), data)
	if err != nil {
		return nil, fmt.Errorf("failed to call template with data: %w", err)
	}

	return component, nil
}

// RequiresData checks if a template requires data service
func (s *templateDataServiceImpl) RequiresData(templateKey string) bool {
	return s.templateRegistry.RequiresDataService(templateKey)
}

// callGetDataMethod calls the GetData method on the data service using reflection
func (s *templateDataServiceImpl) callGetDataMethod(
	ctx context.Context,
	dataService interface{},
	params map[string]string,
) (interface{}, error) {
	// Use reflection to call GetData method
	serviceValue := reflect.ValueOf(dataService)
	getDataMethod := serviceValue.MethodByName("GetData")
	
	if !getDataMethod.IsValid() {
		return nil, fmt.Errorf("GetData method not found on data service")
	}

	// Prepare arguments: ctx, params
	args := []reflect.Value{
		reflect.ValueOf(ctx),
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

	return results[0].Interface(), nil
}

// callTemplateWithData calls the template function with data using reflection
func (s *templateDataServiceImpl) callTemplateWithData(
	templateFunc interface{},
	data interface{},
) (templ.Component, error) {
	// Use reflection to call template function with data
	funcValue := reflect.ValueOf(templateFunc)
	
	if funcValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("template is not a function")
	}

	// Prepare arguments: data
	args := []reflect.Value{
		reflect.ValueOf(data),
	}

	// Call the function
	results := funcValue.Call(args)
	if len(results) != 1 {
		return nil, fmt.Errorf("template function should return one value")
	}

	// Convert result to templ.Component
	component, ok := results[0].Interface().(templ.Component)
	if !ok {
		return nil, fmt.Errorf("template function did not return templ.Component")
	}

	return component, nil
}