package interfaces

import "context"

// DataServiceInfo holds information about required data services
type DataServiceInfo struct {
	InterfaceType string // e.g., "UserDataService" - used for both display and DI resolution
	ParameterType string // e.g., "*dataservices.UserData" (for display)
	// MethodName is always "GetData" - no need to store it
}

// DataService is the generic interface that all data services must implement
// T represents the data type that the service returns
type DataService[T any] interface {
	GetData(ctx context.Context, params map[string]string) (T, error)
}

// DataServiceResolver resolves data services from DI container
type DataServiceResolver interface {
	// ResolveDataService resolves a data service by interface type from DI
	ResolveDataService(interfaceType string) (interface{}, error)
	
	// HasDataService checks if a data service is registered in DI
	HasDataService(interfaceType string) bool
}