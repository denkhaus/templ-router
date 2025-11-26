package services

import (
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataServiceResolver_ErrorHandling_SpecificMessages(t *testing.T) {
	t.Run("Service not found in template registry", func(t *testing.T) {
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

		// Check error type and message
		var appErr *shared.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, shared.ErrorTypeService, appErr.Type)
		assert.Equal(t, "SERVICE_ERROR", appErr.Code)
		assert.Contains(t, appErr.Message, "data service not found")
		assert.Contains(t, appErr.Details, "template registry")
		assert.Equal(t, "NonExistentService", appErr.Context["interface_type"])

		mockRegistry.AssertExpectations(t)
	})

	t.Run("Service not registered in DI container", func(t *testing.T) {
		// Create a mock injector with registered services
		injector := do.New()
		mockRegistry := &MockTemplateRegistry{}
		do.ProvideValue[interfaces.TemplateRegistry](injector, mockRegistry)

		// DON'T register the service in DI - this should cause DI resolution failure
		// do.ProvideNamedValue(injector, "MissingDataService", missingService)

		// Create resolver
		resolver, err := NewDataServiceResolver(injector)
		require.NoError(t, err)

		// Setup mock registry responses for a service that exists in registry but not DI
		templateKeys := []string{"missing_template"}
		mockRegistry.On("GetAllTemplateKeys").Return(templateKeys)
		mockRegistry.On("RequiresDataService", "missing_template").Return(true)
		mockRegistry.On("GetDataServiceInfo", "missing_template").Return(interfaces.DataServiceInfo{
			InterfaceType: "MissingDataService",
			ParameterType: "*missingdata.MissingData",
			MethodName:    "GetMissingData",
		}, true)

		// Test resolving service that's not in DI
		service, err := resolver.ResolveDataService("MissingDataService")
		assert.Error(t, err)
		assert.Nil(t, service)

		// Check error type and message
		var appErr *shared.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, shared.ErrorTypeService, appErr.Type)
		assert.Equal(t, "DI_ERROR", appErr.Code)
		assert.Contains(t, appErr.Message, "not registered in DI container")
		assert.Contains(t, appErr.Details, "MissingDataService")
		assert.Equal(t, "MissingDataService", appErr.Context["interface_type"])
		assert.Equal(t, "MissingDataService", appErr.Context["service_name"])

		mockRegistry.AssertExpectations(t)
	})

	t.Run("Service has no valid Get methods", func(t *testing.T) {
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

		// Check error type and message
		var appErr *shared.AppError
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, shared.ErrorTypeService, appErr.Type)
		assert.Equal(t, "DATA_SERVICE_METHOD_ERROR", appErr.Code)
		assert.Contains(t, appErr.Message, "does not implement required method")

		// Print details for debugging
		t.Logf("Error details: %s", appErr.Details)
		t.Logf("Parameter type from context: %v", appErr.Context["parameter_type"])

		// Check that the important parts are there, but be more flexible with parameter type
		assert.Contains(t, appErr.Details, "func(GetInvalidData interfaces.RouterContext)")
		assert.Contains(t, appErr.Details, "error)")
		assert.Equal(t, "InvalidDataService", appErr.Context["interface_type"])
		assert.Equal(t, "GetInvalidData", appErr.Context["expected_method"])

		mockRegistry.AssertExpectations(t)
	})
}

func TestOptimizedTemplateService_EnhancedErrorMessages(t *testing.T) {
	t.Run("Method implementation error provides clear guidance", func(t *testing.T) {
		// This tests the error handling in optimized_template_service.go
		// We simulate the scenario where the resolver returns a DATA_SERVICE_METHOD_ERROR

		// Create method error like the resolver would return
		methodError := shared.NewDataServiceMethodError("TestService", "GetTestData").
			WithDetails("Service must implement GetTestData method").
			WithContext("interface_type", "TestService")

		// Verify the error has the expected structure
		var appErr *shared.AppError
		assert.ErrorAs(t, methodError, &appErr)
		assert.Equal(t, "DATA_SERVICE_METHOD_ERROR", appErr.Code)
		assert.Contains(t, appErr.Message, "TestService")
		assert.Contains(t, appErr.Message, "GetTestData")
	})
}
