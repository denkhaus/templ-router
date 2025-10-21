package services

import (
	"fmt"
	
	"github.com/denkhaus/templ-router/pkg/interfaces"
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
		return nil, fmt.Errorf("no DataService found for interface type: %s", interfaceType)
	}
	
	// Resolve the service from DI using the named dependency
	// This is the elegant, generic solution
	return r.resolveNamedDataService(dataServiceInfo.InterfaceType)
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
		return nil, fmt.Errorf("failed to resolve DataService '%s': %w", serviceName, err)
	}
	
	return service, nil
}

// HasDataService checks if a data service is registered in DI
func (r *dataServiceResolverImpl) HasDataService(interfaceType string) bool {
	_, err := r.ResolveDataService(interfaceType)
	return err == nil
}