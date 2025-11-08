package services

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
)

// dataServiceResolver resolves data services with specific method detection
type dataServiceResolver struct {
	injector         do.Injector
	templateRegistry interfaces.TemplateRegistry
}

// NewDataServiceResolver creates a new data service resolver for specific methods
func NewDataServiceResolver(i do.Injector) (interfaces.DataServiceResolver, error) {
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)

	return &dataServiceResolver{
		injector:         i,
		templateRegistry: templateRegistry,
	}, nil
}

// ResolveDataService resolves a data service by interface type with specific method detection
func (r *dataServiceResolver) ResolveDataService(interfaceType string) (any, error) {
	// Get service info from template registry
	dataServiceInfo := r.findDataServiceInfo(interfaceType)
	if dataServiceInfo == nil {
		return nil, shared.NewServiceError("data service not found").
			WithDetails("No DataService registered for the specified interface type").
			WithContext("interface_type", interfaceType)
	}

	// Resolve the service from DI using the named dependency
	service, err := r.resolveNamedDataService(dataServiceInfo.InterfaceType)
	if err != nil {
		return nil, err
	}

	// Detect specific method - this must exist
	specificMethodName := r.detectSpecificMethod(service)
	if specificMethodName == "" {
		return nil, shared.NewServiceError("data service does not implement a supported method").
			WithDetails("Service must implement a specific Get<StructName>() method").
			WithContext("interface_type", interfaceType).
			WithContext("service_type", fmt.Sprintf("%T", service))
	}

	// Return service wrapped for specific method calling
	return &specificMethodWrapper{service: service, methodName: specificMethodName}, nil
}

// HasDataService checks if a data service is registered in DI
func (r *dataServiceResolver) HasDataService(interfaceType string) bool {
	_, err := r.ResolveDataService(interfaceType)
	return err == nil
}

// specificMethodWrapper wraps services that implement Get<StructName>() methods
type specificMethodWrapper struct {
	service    any
	methodName string
}

// GetData calls the specific Get<StructName>() method
func (w *specificMethodWrapper) GetData(routerCtx interfaces.RouterContext) (any, error) {
	serviceValue := reflect.ValueOf(w.service)
	method := serviceValue.MethodByName(w.methodName)
	if !method.IsValid() {
		return nil, shared.NewServiceError("specific method not found").
			WithDetails("Method disappeared after detection").
			WithContext("method_name", w.methodName)
	}

	// Call the method with routerCtx as argument
	results := method.Call([]reflect.Value{reflect.ValueOf(routerCtx)})
	if len(results) != 2 {
		return nil, shared.NewServiceError("method signature invalid").
			WithDetails("Method must return (data, error)").
			WithContext("method_name", w.methodName).
			WithContext("return_count", len(results))
	}

	// Extract the data and error from the results
	data := results[0].Interface()
	err, _ := results[1].Interface().(error)

	return data, err
}

// detectSpecificMethod detects Get<StructName>() methods in a service
func (r *dataServiceResolver) detectSpecificMethod(service any) string {
	serviceType := reflect.TypeOf(service)

	// Look for methods that start with "Get" and return (*Struct, error)
	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)

		// Skip if method doesn't start with "Get"
		if !strings.HasPrefix(method.Name, "Get") {
			continue
		}

		// Check method signature: func (r) GetXxx(routerCtx RouterContext) (*Xxx, error)
		if method.Type.NumIn() != 2 || method.Type.NumOut() != 2 {
			continue
		}

		// Check first parameter is RouterContext
		routerCtxType := method.Type.In(1)
		if !routerCtxType.Implements(reflect.TypeOf((*interfaces.RouterContext)(nil)).Elem()) {
			continue
		}

		// Check return values are (*Struct, error)
		if !method.Type.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			continue
		}

		// Check first return is a pointer to a struct
		returnType := method.Type.Out(0)
		if returnType.Kind() == reflect.Ptr && returnType.Elem().Kind() == reflect.Struct {
			// This is a valid Get<StructName>() method
			return method.Name
		}
	}

	return ""
}

// findDataServiceInfo searches the template registry for DataService info by interface type
func (r *dataServiceResolver) findDataServiceInfo(interfaceType string) *interfaces.DataServiceInfo {
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
func (r *dataServiceResolver) resolveNamedDataService(serviceName string) (any, error) {
	service, err := do.InvokeNamed[any](r.injector, serviceName)
	if err != nil {
		return nil, shared.NewDependencyInjectionError("failed to resolve DataService from DI container").
			WithDetails("Service not found or not properly registered").
			WithCause(err).
			WithContext("service_name", serviceName).
			WithContext("interface_type", serviceName)
	}

	return service, nil
}