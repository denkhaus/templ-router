package services

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
)

// dataServiceResolver eliminates reflection through type-safe registration
type dataServiceResolver struct {
	injector         do.Injector
	templateRegistry interfaces.TemplateRegistry
	serviceRegistry  map[string]interfaces.GenericDataService
}

// NewDataServiceResolver creates a new reflection-free data service resolver
func NewDataServiceResolver(i do.Injector) (interfaces.DataServiceResolver, error) {
	templateRegistry := do.MustInvoke[interfaces.TemplateRegistry](i)

	return &dataServiceResolver{
		injector:         i,
		templateRegistry: templateRegistry,
		serviceRegistry:  make(map[string]interfaces.GenericDataService),
	}, nil
}

// RegisterDataService registers a data service without reflection
func (r *dataServiceResolver) RegisterDataService(interfaceType string, service interfaces.GenericDataService) {
	r.serviceRegistry[interfaceType] = service
}

// ResolveDataService resolves a data service by interface type without reflection
func (r *dataServiceResolver) ResolveDataService(interfaceType string) (any, error) {
	// First check the type-safe registry
	if service, exists := r.serviceRegistry[interfaceType]; exists {
		return service, nil
	}

	// Fall back to DI resolution for legacy services
	dataServiceInfo := r.findDataServiceInfo(interfaceType)
	if dataServiceInfo == nil {
		return nil, shared.NewServiceError("data service not found").
			WithDetails("No DataService registered for the specified interface type").
			WithContext("interface_type", interfaceType)
	}

	// Resolve the service from DI using the named dependency
	return r.resolveNamedDataService(dataServiceInfo.InterfaceType)
}

// ResolveGenericDataService resolves a data service as generic interface without reflection
func (r *dataServiceResolver) ResolveGenericDataService(interfaceType string) (interfaces.GenericDataService, error) {
	// First check the type-safe registry
	if service, exists := r.serviceRegistry[interfaceType]; exists {
		return service, nil
	}

	// Fall back to creating a type-safe wrapper
	service, err := r.ResolveDataService(interfaceType)
	if err != nil {
		return nil, err
	}

	// Check if service already implements GenericDataService
	if genericService, ok := service.(interfaces.GenericDataService); ok {
		return genericService, nil
	}

	// Create a type-safe wrapper based on the service type
	dataServiceInfo := r.findDataServiceInfo(interfaceType)
	if dataServiceInfo == nil {
		return nil, shared.NewServiceError("data service info not found").
			WithContext("interface_type", interfaceType)
	}

	// Create a reflection-free wrapper using type assertion
	switch typedService := service.(type) {
	case interfaces.DataService[any]:
		return &typedDataServiceWrapper[any]{service: typedService}, nil
	default:
		// For services that don't implement the generic interface,
		// we still need to use reflection minimally only during registration
		return &legacyDataServiceWrapper{service: service}, nil
	}
}

// typedDataServiceWrapper is a type-safe wrapper for services implementing DataService[T]
type typedDataServiceWrapper[T any] struct {
	service interfaces.DataService[T]
}

// GetData implements GenericDataService without reflection
func (w *typedDataServiceWrapper[T]) GetData(routerCtx interfaces.RouterContext) (any, error) {
	return w.service.GetData(routerCtx)
}

// legacyDataServiceWrapper provides a minimal reflection wrapper for legacy services
type legacyDataServiceWrapper struct {
	service any
}

// GetData implements GenericDataService with minimal reflection (only during wrapper creation)
func (w *legacyDataServiceWrapper) GetData(routerCtx interfaces.RouterContext) (any, error) {
	// This should rarely be used in practice as services should be migrated to the new interface
	return nil, shared.NewServiceError("legacy data service detected").
		WithDetails("Service should be migrated to implement DataService[T] interface").
		WithContext("service_type", "legacy")
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

// HasDataService checks if a data service is registered in DI
func (r *dataServiceResolver) HasDataService(interfaceType string) bool {
	_, err := r.ResolveDataService(interfaceType)
	return err == nil
}
