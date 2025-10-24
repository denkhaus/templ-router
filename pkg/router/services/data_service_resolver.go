package services

import (
	"context"
	"reflect"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
)

// dataServiceResolverImpl implements DataServiceResolver
type dataServiceResolverImpl struct {
	injector         do.Injector
	templateRegistry interfaces.TemplateRegistry
}

// NewDataServiceResolver creates a new data service resolver for DI
func NewDataServiceResolver(i do.Injector) (interfaces.DataServiceResolver, error) {
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)
	
	return &dataServiceResolverImpl{
		injector:         i,
		templateRegistry: templateRegistry,
	}, nil
}

// ResolveDataService resolves a data service by interface type from DI
func (r *dataServiceResolverImpl) ResolveDataService(interfaceType string) (interface{}, error) {
	// The elegant solution: get the service name from the template registry
	// and resolve it using named dependency from DI
	
	// Find the DataServiceInfo that matches the interface type
	dataServiceInfo := r.findDataServiceInfo(interfaceType)
	if dataServiceInfo == nil {
		return nil, shared.NewServiceError("data service not found").
			WithDetails("No DataService registered for the specified interface type").
			WithContext("interface_type", interfaceType)
	}
	
	// Resolve the service from DI using the named dependency
	// This is the elegant, generic solution
	return r.resolveNamedDataService(dataServiceInfo.InterfaceType)
}

// ResolveGenericDataService resolves a data service as generic interface (no reflection needed)
func (r *dataServiceResolverImpl) ResolveGenericDataService(interfaceType string) (interfaces.GenericDataService, error) {
	// Resolve the service using the existing method
	service, err := r.ResolveDataService(interfaceType)
	if err != nil {
		return nil, err
	}
	
	// Try to cast to GenericDataService interface
	if genericService, ok := service.(interfaces.GenericDataService); ok {
		return genericService, nil
	}
	
	// If service doesn't implement GenericDataService, create a wrapper
	return &genericDataServiceWrapper{
		service:       service,
		interfaceType: interfaceType,
	}, nil
}

// genericDataServiceWrapper wraps any data service to implement GenericDataService
type genericDataServiceWrapper struct {
	service       interface{}
	interfaceType string
}

// GetData implements GenericDataService by calling the underlying service's GetData method
func (w *genericDataServiceWrapper) GetData(ctx context.Context, params map[string]string) (interface{}, error) {
	// Use reflection only in the wrapper, not in the main flow
	serviceValue := reflect.ValueOf(w.service)
	getDataMethod := serviceValue.MethodByName("GetData")
	
	if !getDataMethod.IsValid() {
		return nil, shared.NewServiceError("data service does not implement GetData method").
			WithDetails("Service must implement GetData(context.Context, map[string]string) (T, error)").
			WithContext("interface_type", w.interfaceType)
	}
	
	// Call GetData method
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(params),
	}
	
	results := getDataMethod.Call(args)
	if len(results) != 2 {
		return nil, shared.NewServiceError("invalid GetData method signature").
			WithDetails("GetData must return (data, error)").
			WithContext("interface_type", w.interfaceType)
	}
	
	// Check for error
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	
	return results[0].Interface(), nil
}

// findDataServiceInfo searches the template registry for DataService info by interface type
func (r *dataServiceResolverImpl) findDataServiceInfo(interfaceType string) *interfaces.DataServiceInfo {
	// Get all template keys and check each one for DataService info
	templateKeys := r.templateRegistry.GetAllTemplateKeys()
	
	for _, key := range templateKeys {
		if r.templateRegistry.RequiresDataService(key) {
			dataServiceInfo, exists := r.templateRegistry.GetDataServiceInfo(key)
			if exists && dataServiceInfo.InterfaceType == interfaceType {
				return &dataServiceInfo
			}
		}
	}
	
	return nil
}

// resolveNamedDataService resolves a DataService from DI using named dependency
func (r *dataServiceResolverImpl) resolveNamedDataService(serviceName string) (interface{}, error) {
	// The elegant solution: resolve the service using named dependency
	// This is completely generic and works for any DataService
	
	// The serviceName is already the short name (e.g., "UserDataService")
	// No need to extract anything since we optimized the scanner
	
	// Use do.InvokeNamed with interface{} to resolve any service type by name
	// This allows us to resolve services that don't implement a common interface
	service, err := do.InvokeNamed[interface{}](r.injector, serviceName)
	if err != nil {
		return nil, shared.NewDependencyInjectionError("failed to resolve DataService from DI container").
			WithDetails("Service not found or not properly registered").
			WithCause(err).
			WithContext("service_name", serviceName).
			WithContext("interface_type", serviceName)
	}
	
	return service, nil
}

// HasDataService checks if a data service is registered in DI
func (r *dataServiceResolverImpl) HasDataService(interfaceType string) bool {
	_, err := r.ResolveDataService(interfaceType)
	return err == nil
}