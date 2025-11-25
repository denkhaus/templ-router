package services

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/demo/pkg/dataservices"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementations for testing

type MockTemplateRegistry struct {
	mock.Mock
}

func (m *MockTemplateRegistry) GetRouteToTemplateMapping() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *MockTemplateRegistry) GetTemplateFunction(uuid string) (func() any, bool) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(func() any), args.Bool(1)
}

func (m *MockTemplateRegistry) RequiresDataService(templateKey string) bool {
	args := m.Called(templateKey)
	return args.Bool(0)
}

func (m *MockTemplateRegistry) GetDataServiceInfo(templateKey string) (interfaces.DataServiceInfo, bool) {
	args := m.Called(templateKey)
	return args.Get(0).(interfaces.DataServiceInfo), args.Bool(1)
}

func (m *MockTemplateRegistry) GetAllTemplateKeys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockTemplateRegistry) GetTemplate(key string) (templ.Component, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(templ.Component), args.Error(1)
}

func (m *MockTemplateRegistry) IsAvailable(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockTemplateRegistry) GetTemplateByRoute(route string) (templ.Component, error) {
	args := m.Called(route)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(templ.Component), args.Error(1)
}

func (m *MockTemplateRegistry) GetTemplateMetadata(key string) (*interfaces.TemplateMetadata, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.TemplateMetadata), args.Error(1)
}

func (m *MockTemplateRegistry) GetTemplateMetadataByRoute(route string) (*interfaces.TemplateMetadata, error) {
	args := m.Called(route)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.TemplateMetadata), args.Error(1)
}

func (m *MockTemplateRegistry) GetAllTemplateMetadata() map[string]*interfaces.TemplateMetadata {
	args := m.Called()
	return args.Get(0).(map[string]*interfaces.TemplateMetadata)
}

func (m *MockTemplateRegistry) FindComponentTemplates() map[string]*interfaces.TemplateMetadata {
	args := m.Called()
	return args.Get(0).(map[string]*interfaces.TemplateMetadata)
}

func (m *MockTemplateRegistry) GetTemplateKeyByComponentName(componentName string) (string, bool) {
	args := m.Called(componentName)
	return args.String(0), args.Bool(1)
}

type MockRouterContext struct {
	mock.Mock
	urlParams  map[string]string
	queryParams url.Values
	request    *http.Request
	chiCtx     *chi.Context
}

func (m *MockRouterContext) Context() context.Context {
	args := m.Called()
	if args.Get(0) == nil {
		return context.Background()
	}
	return args.Get(0).(context.Context)
}

func (m *MockRouterContext) GetURLParam(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockRouterContext) GetAllURLParams() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *MockRouterContext) GetQueryParam(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockRouterContext) GetQueryParams(key string) []string {
	args := m.Called(key)
	return args.Get(0).([]string)
}

func (m *MockRouterContext) GetAllQueryParams() url.Values {
	args := m.Called()
	return args.Get(0).(url.Values)
}

func (m *MockRouterContext) Request() *http.Request {
	args := m.Called()
	if args.Get(0) == nil {
		return m.request
	}
	return args.Get(0).(*http.Request)
}

func (m *MockRouterContext) ChiContext() *chi.Context {
	args := m.Called()
	if args.Get(0) == nil {
		return m.chiCtx
	}
	return args.Get(0).(*chi.Context)
}

// Mock data services for testing

type ValidDataService struct{}

func (s *ValidDataService) GetTestData(ctx interfaces.RouterContext) (*TestData, error) {
	return &TestData{ID: "123", Name: "Test"}, nil
}

func (s *ValidDataService) GetData(ctx interfaces.RouterContext) (any, error) {
	return s.GetTestData(ctx)
}

type InvalidDataService struct{}

func (s *InvalidDataService) DoSomething(ctx interfaces.RouterContext) error {
	return nil
}

type MultipleMethodsDataService struct{}

func (s *MultipleMethodsDataService) GetUserData(ctx interfaces.RouterContext) (*UserData, error) {
	return &UserData{ID: "user123"}, nil
}

func (s *MultipleMethodsDataService) GetProductData(ctx interfaces.RouterContext) (*ProductData, error) {
	return &ProductData{Name: "Product"}, nil
}

type TestData struct {
	ID   string
	Name string
}

type UserData struct {
	ID string
}

type ProductData struct {
	Name string
}

// specificOnlyMockService mimics the real demo SpecificOnlyDataService for testing
type specificOnlyMockService struct{}

func (s *specificOnlyMockService) GetSpecificData(ctx interfaces.RouterContext) (*dataservices.SpecificData, error) {
	locale := ctx.GetURLParam("locale")
	if locale == "" {
		locale = "en"
	}

	return &dataservices.SpecificData{
		ID:          "test-id",
		Title:       "Test Title",
		Description: "Test Description",
		Method:      "GetSpecificData",
		Note:        "This is a mock implementation with locale: " + locale,
	}, nil
}

func TestNewDataServiceResolver(t *testing.T) {
	// Create a mock injector
	injector := do.New()

	// Create and register mock template registry as interface
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Test successful creation
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)
	assert.NotNil(t, resolver)

	// Verify the resolver type
	assert.IsType(t, &dataServiceResolver{}, resolver)
}

func TestDataServiceResolver_ResolveDataService_Success(t *testing.T) {
	// Create a mock injector with registered services
	injector := do.New()
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Register a valid data service
	validService := &ValidDataService{}
	do.ProvideNamedValue(injector, "TestDataService", validService)

	// Create resolver
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)

	// Setup mock registry responses
	templateKeys := []string{"test_template"}
	mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
	mockRegistry.On("RequiresDataService", "test_template").Return(true)
	mockRegistry.On("GetDataServiceInfo", "test_template").Return(interfaces.DataServiceInfo{
		InterfaceType: "TestDataService",
		ParameterType: "*testdata.TestData",
		MethodName:    "GetTestData",
	}, true)

	// Test resolving the service
	service, err := resolver.ResolveDataService("TestDataService")
	require.NoError(t, err)
	assert.NotNil(t, service)

	mockRegistry.AssertExpectations(t)
}

func TestDataServiceResolver_ResolveDataService_NotFound(t *testing.T) {
	// Create a mock injector with registered services
	injector := do.New()
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Create resolver
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)

	// Setup mock registry responses - no matching service
	templateKeys := []string{"other_template"}
	mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
	mockRegistry.On("RequiresDataService", "other_template").Return(true)
	mockRegistry.On("GetDataServiceInfo", "other_template").Return(interfaces.DataServiceInfo{
		InterfaceType: "OtherDataService",
		ParameterType: "*otherdata.OtherData",
		MethodName:    "GetOtherData",
	}, true)

	// Test resolving non-existent service
	service, err := resolver.ResolveDataService("NonExistentService")
	assert.Error(t, err)
	assert.Nil(t, service)

	// Check error type
	var appErr *shared.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, shared.ErrorTypeService, appErr.Type)
	assert.Contains(t, appErr.Message, "data service not found")

	mockRegistry.AssertExpectations(t)
}

func TestDataServiceResolver_ResolveDataService_InvalidService(t *testing.T) {
	// Create a mock injector with registered services
	injector := do.New()
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Register an invalid data service (no Get methods)
	invalidService := &InvalidDataService{}
	do.ProvideNamedValue(injector, "InvalidDataService", invalidService)

	// Create resolver
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)

	// Setup mock registry responses
	templateKeys := []string{"invalid_template"}
	mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
	mockRegistry.On("RequiresDataService", "invalid_template").Return(true)
	mockRegistry.On("GetDataServiceInfo", "invalid_template").Return(interfaces.DataServiceInfo{
		InterfaceType: "InvalidDataService",
		ParameterType: "*invaliddata.InvalidData",
		MethodName:    "GetInvalidData",
	}, true)

	// Test resolving invalid service
	service, err := resolver.ResolveDataService("InvalidDataService")
	assert.Error(t, err)
	assert.Nil(t, service)

	// Check error type
	var appErr *shared.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, shared.ErrorTypeService, appErr.Type)
	assert.Contains(t, appErr.Message, "does not implement required method")

	mockRegistry.AssertExpectations(t)
}

func TestDataServiceResolver_HasDataService(t *testing.T) {
	t.Run("Existing service", func(t *testing.T) {
		// Create a mock injector with registered services
		injector := do.New()
		mockRegistry := &MockTemplateRegistry{}
		do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

		// Register services
		validService := &ValidDataService{}
		do.ProvideNamedValue(injector, "TestDataService", validService)

		// Create resolver
		resolver, err := NewDataServiceResolver(injector)
		require.NoError(t, err)

		// Setup mock registry responses for existing service
		templateKeys := []string{"test_template"}
		mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
		mockRegistry.On("RequiresDataService", "test_template").Return(true)
		mockRegistry.On("GetDataServiceInfo", "test_template").Return(interfaces.DataServiceInfo{
			InterfaceType: "TestDataService",
			ParameterType: "*testdata.TestData",
			MethodName:    "GetTestData",
		}, true)

		// Test existing service
		hasService := resolver.HasDataService("TestDataService")
		assert.True(t, hasService)

		mockRegistry.AssertExpectations(t)
	})

	t.Run("Non-existing service", func(t *testing.T) {
		// Create a mock injector with registered services
		injector := do.New()
		mockRegistry := &MockTemplateRegistry{}
		do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

		// Create resolver
		resolver, err := NewDataServiceResolver(injector)
		require.NoError(t, err)

		// Setup mock registry responses for non-existing service
		templateKeys := []string{"other_template"}
		mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
		mockRegistry.On("RequiresDataService", "other_template").Return(true)
		mockRegistry.On("GetDataServiceInfo", "other_template").Return(interfaces.DataServiceInfo{
			InterfaceType: "OtherDataService",
			ParameterType: "*otherdata.OtherData",
			MethodName:    "GetOtherData",
		}, true)

		// Test non-existing service
		hasService := resolver.HasDataService("NonExistentService")
		assert.False(t, hasService)

		mockRegistry.AssertExpectations(t)
	})
}

func TestSpecificMethodWrapper_GetData_Success(t *testing.T) {
	// Create wrapper with valid service
	service := &ValidDataService{}
	wrapper := &specificMethodWrapper{
		service:    service,
		methodName: "GetTestData",
	}

	// Create mock router context (no expectations since GetTestData doesn't call context methods)
	mockCtx := &MockRouterContext{}

	// Test GetData
	result, err := wrapper.GetData(mockCtx)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify result type
	testData, ok := result.(*TestData)
	require.True(t, ok)
	assert.Equal(t, "123", testData.ID)
	assert.Equal(t, "Test", testData.Name)
}

func TestSpecificMethodWrapper_GetData_MethodNotFound(t *testing.T) {
	// Create wrapper with non-existent method
	service := &ValidDataService{}
	wrapper := &specificMethodWrapper{
		service:    service,
		methodName: "NonExistentMethod",
	}

	// Create mock router context (no expectations needed since method won't be found)
	mockCtx := &MockRouterContext{}

	// Test GetData should fail
	result, err := wrapper.GetData(mockCtx)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Check error type
	var appErr *shared.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, shared.ErrorTypeService, appErr.Type)
	assert.Contains(t, appErr.Message, "method not found")
}

func TestDataServiceResolver_detectSpecificMethod(t *testing.T) {
	// Create a mock injector with registered services
	injector := do.New()
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Create resolver
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)

	// Cast to concrete type to access private methods
	concreteResolver := resolver.(*dataServiceResolver)

	t.Run("Valid service with Get method", func(t *testing.T) {
		service := &ValidDataService{}
		method := concreteResolver.detectSpecificMethod(service)
		assert.Equal(t, "GetTestData", method)
	})

	t.Run("Service with multiple Get methods", func(t *testing.T) {
		service := &MultipleMethodsDataService{}
		method := concreteResolver.detectSpecificMethod(service)
		// Should return one of the Get methods (implementation dependent)
		assert.True(t, method == "GetUserData" || method == "GetProductData")
	})

	t.Run("Invalid service without Get methods", func(t *testing.T) {
		service := &InvalidDataService{}
		method := concreteResolver.detectSpecificMethod(service)
		assert.Empty(t, method)
	})

	t.Run("Real demo service - SpecificOnlyDataService", func(t *testing.T) {
		// Use the predefined mock service that matches the real SpecificOnlyDataService
		mockService := &specificOnlyMockService{}
		method := concreteResolver.detectSpecificMethod(mockService)
		assert.Equal(t, "GetSpecificData", method)
	})
}

func TestDataServiceResolver_findDataServiceInfo(t *testing.T) {
	t.Run("Found service", func(t *testing.T) {
		// Create a mock injector with registered services
		injector := do.New()
		mockRegistry := &MockTemplateRegistry{}
		do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

		// Create resolver
		resolver, err := NewDataServiceResolver(injector)
		require.NoError(t, err)

		// Cast to concrete type to access private methods
		concreteResolver := resolver.(*dataServiceResolver)

		templateKeys := []string{"template1", "template2"}
		mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
		mockRegistry.On("RequiresDataService", "template1").Return(false)
		mockRegistry.On("RequiresDataService", "template2").Return(true)
		mockRegistry.On("GetDataServiceInfo", "template2").Return(interfaces.DataServiceInfo{
			InterfaceType: "TargetService",
			ParameterType: "*targetdata.TargetData",
			MethodName:    "GetTargetData",
		}, true)

		result := concreteResolver.findDataServiceInfo("TargetService")
		assert.NotNil(t, result)
		assert.Equal(t, "TargetService", result.InterfaceType)

		mockRegistry.AssertExpectations(t)
	})

	t.Run("Service not found", func(t *testing.T) {
		// Create a mock injector with registered services
		injector := do.New()
		mockRegistry := &MockTemplateRegistry{}
		do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

		// Create resolver
		resolver, err := NewDataServiceResolver(injector)
		require.NoError(t, err)

		// Cast to concrete type to access private methods
		concreteResolver := resolver.(*dataServiceResolver)

		templateKeys := []string{"template1"}
		mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
		mockRegistry.On("RequiresDataService", "template1").Return(true)
		// This service has a different interface type than what we're looking for
		mockRegistry.On("GetDataServiceInfo", "template1").Return(interfaces.DataServiceInfo{
			InterfaceType: "OtherService", // Different from "TargetService"
			ParameterType: "*otherdata.OtherData",
			MethodName:    "GetOtherData",
		}, true)

		result := concreteResolver.findDataServiceInfo("TargetService") // Looking for "TargetService"
		assert.Nil(t, result) // Should return nil because no service has interfaceType "TargetService"

		mockRegistry.AssertExpectations(t)
	})
}

func TestDataServiceResolver_resolveNamedDataService(t *testing.T) {
	// Create a mock injector with registered services
	injector := do.New()
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Create resolver
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)

	// Cast to concrete type to access private methods
	concreteResolver := resolver.(*dataServiceResolver)

	t.Run("Successfully resolve named service", func(t *testing.T) {
		// Register a service
		testService := &ValidDataService{}
		do.ProvideNamedValue(injector, "TestService", testService)

		result, err := concreteResolver.resolveNamedDataService("TestService")
		require.NoError(t, err)
		assert.Equal(t, testService, result)
	})

	t.Run("Failed to resolve non-existent service", func(t *testing.T) {
		result, err := concreteResolver.resolveNamedDataService("NonExistentService")
		assert.Error(t, err)
		assert.Nil(t, result)

		// Check error type
		var appErr *shared.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, shared.ErrorTypeService, appErr.Type)
		assert.Contains(t, appErr.Message, "failed to resolve DataService from DI container")
	})
}

func TestDataServiceResolver_Integration_GetDataCall(t *testing.T) {
	// Create a mock injector with registered services
	injector := do.New()
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Register a service that returns an error
	errorService := &struct {
		methodCalled bool
	}{}
	do.ProvideNamedValue(injector, "ErrorService", errorService)

	// Create resolver
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)

	// Setup mock registry for error service
	templateKeys := []string{"error_template"}
	mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
	mockRegistry.On("RequiresDataService", "error_template").Return(true)
	mockRegistry.On("GetDataServiceInfo", "error_template").Return(interfaces.DataServiceInfo{
		InterfaceType: "ErrorService",
		ParameterType: "*errordata.ErrorData",
		MethodName:    "GetErrorData",
	}, true)

	// This test demonstrates the integration flow
	// In a real scenario, the service would need proper Get<StructName> methods
	service, err := resolver.ResolveDataService("ErrorService")
	// Should fail because the mock service doesn't have proper Get methods
	assert.Error(t, err)
	assert.Nil(t, service)

	mockRegistry.AssertExpectations(t)
}

func TestDataServiceResolver_Integration_WithDemoSpecificOnlyService(t *testing.T) {
	// Create a mock injector with registered services
	injector := do.New()
	mockRegistry := &MockTemplateRegistry{}
	do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

	// Register the demo-like specific-only service
	specificOnlyService := &specificOnlyMockService{}
	do.ProvideNamedValue(injector, "SpecificOnlyDataService", specificOnlyService)

	// Create resolver
	resolver, err := NewDataServiceResolver(injector)
	require.NoError(t, err)

	// Setup mock registry for the specific-only service
	templateKeys := []string{"specific_template"}
	mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
	mockRegistry.On("RequiresDataService", "specific_template").Return(true)
	mockRegistry.On("GetDataServiceInfo", "specific_template").Return(interfaces.DataServiceInfo{
		InterfaceType: "SpecificOnlyDataService",
		ParameterType: "*dataservices.SpecificData",
		MethodName:    "GetSpecificData",
	}, true)

	// Test resolving the specific-only service
	service, err := resolver.ResolveDataService("SpecificOnlyDataService")
	require.NoError(t, err)
	assert.NotNil(t, service)

	// Verify the service is wrapped with the specific method wrapper
	wrapper, ok := service.(*specificMethodWrapper)
	require.True(t, ok)
	assert.Equal(t, "GetSpecificData", wrapper.methodName)

	// Test calling GetData on the wrapped service
	mockRouterCtx := &MockRouterContext{}
	mockRouterCtx.On("GetURLParam", "locale").Return("en")

	result, err := wrapper.GetData(mockRouterCtx)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the result is the expected SpecificData
	specificData, ok := result.(*dataservices.SpecificData)
	require.True(t, ok)
	assert.Equal(t, "test-id", specificData.ID)
	assert.Equal(t, "Test Title", specificData.Title)
	assert.Equal(t, "GetSpecificData", specificData.Method)

	mockRegistry.AssertExpectations(t)
	mockRouterCtx.AssertExpectations(t)
}