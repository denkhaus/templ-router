package services

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/a-h/templ"
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/middleware"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Use unique mock names to avoid conflicts with existing mocks
type mockOTSTemplateRegistry struct {
	mock.Mock
}

func (m *mockOTSTemplateRegistry) GetRouteToTemplateMapping() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *mockOTSTemplateRegistry) GetTemplateFunction(uuid string) (func() interface{}, bool) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(func() interface{}), args.Bool(1)
}

func (m *mockOTSTemplateRegistry) RequiresDataService(templateKey string) bool {
	args := m.Called(templateKey)
	return args.Bool(0)
}

func (m *mockOTSTemplateRegistry) GetDataServiceInfo(templateKey string) (interfaces.DataServiceInfo, bool) {
	args := m.Called(templateKey)
	return args.Get(0).(interfaces.DataServiceInfo), args.Bool(1)
}

func (m *mockOTSTemplateRegistry) GetAllTemplateKeys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *mockOTSTemplateRegistry) GetTemplate(key string) (templ.Component, error) {
	args := m.Called(key)
	return args.Get(0).(templ.Component), args.Error(1)
}

func (m *mockOTSTemplateRegistry) IsAvailable(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *mockOTSTemplateRegistry) GetTemplateByRoute(route string) (templ.Component, error) {
	args := m.Called(route)
	return args.Get(0).(templ.Component), args.Error(1)
}

type mockOTSCacheService struct {
	mock.Mock
}

func (m *mockOTSCacheService) BuildTemplateKey(templateFile, locale string, params map[string]string) string {
	args := m.Called(templateFile, locale, params)
	return args.String(0)
}

func (m *mockOTSCacheService) BuildRouteKey(routePath string, params map[string]string) string {
	args := m.Called(routePath, params)
	return args.String(0)
}

func (m *mockOTSCacheService) GetTemplate(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *mockOTSCacheService) SetTemplate(key string, template interface{}) {
	m.Called(key, template)
}

func (m *mockOTSCacheService) GetRoute(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *mockOTSCacheService) SetRoute(key string, templateUUID interface{}) {
	m.Called(key, templateUUID)
}

func (m *mockOTSCacheService) ClearAll() {
	m.Called()
}

func (m *mockOTSCacheService) ClearTemplates() {
	m.Called()
}

func (m *mockOTSCacheService) ClearRoutes() {
	m.Called()
}

type mockOTSDataServiceResolver struct {
	mock.Mock
}

func (m *mockOTSDataServiceResolver) ResolveDataService(interfaceType string) (interface{}, error) {
	args := m.Called(interfaceType)
	return args.Get(0), args.Error(1)
}

func (m *mockOTSDataServiceResolver) ResolveGenericDataService(interfaceType string) (interfaces.GenericDataService, error) {
	args := m.Called(interfaceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(interfaces.GenericDataService), args.Error(1)
}

func (m *mockOTSDataServiceResolver) HasDataService(interfaceType string) bool {
	args := m.Called(interfaceType)
	return args.Bool(0)
}

type mockOTSRouteConverter struct {
	mock.Mock
}

func (m *mockOTSRouteConverter) GenerateRouteVariations(routePath string) []string {
	args := m.Called(routePath)
	return args.Get(0).([]string)
}

func (m *mockOTSRouteConverter) ConvertLayoutPathToRoute(layoutPath string) string {
	args := m.Called(layoutPath)
	return args.String(0)
}

func (m *mockOTSRouteConverter) GenerateTemplateKey(templateFile string) string {
	args := m.Called(templateFile)
	return args.String(0)
}

// Mock component for testing
type mockOTSComponent struct{}

func (m mockOTSComponent) Render(ctx context.Context, w io.Writer) error {
	return nil
}

// Mock RouterContext for testing
type mockOTSRouterContext struct {
	mock.Mock
	ctx         context.Context
	urlParams   map[string]string
	queryParams url.Values
	request     *http.Request
	chiCtx      *chi.Context
}

func (m *mockOTSRouterContext) Context() context.Context {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(context.Context)
	}
	return m.ctx
}

func (m *mockOTSRouterContext) GetURLParam(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *mockOTSRouterContext) GetAllURLParams() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *mockOTSRouterContext) GetQueryParam(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *mockOTSRouterContext) GetQueryParams(key string) []string {
	args := m.Called(key)
	return args.Get(0).([]string)
}

func (m *mockOTSRouterContext) GetAllQueryParams() url.Values {
	args := m.Called()
	return args.Get(0).(url.Values)
}

func (m *mockOTSRouterContext) Request() *http.Request {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*http.Request)
	}
	return m.request
}

func (m *mockOTSRouterContext) ChiContext() *chi.Context {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*chi.Context)
	}
	return m.chiCtx
}

// Test helper to create OptimizedTemplateService with mocks
func createTestOTS(t *testing.T) (*OptimizedTemplateService, *mockOTSTemplateRegistry, *mockOTSCacheService, *mockOTSDataServiceResolver, *mockOTSRouteConverter) {
	logger := zap.NewNop()
	
	mockRegistry := &mockOTSTemplateRegistry{}
	mockCache := &mockOTSCacheService{}
	mockDataResolver := &mockOTSDataServiceResolver{}
	mockConverter := &mockOTSRouteConverter{}

	service := &OptimizedTemplateService{
		logger:           logger,
		templateRegistry: mockRegistry,
		cacheService:     mockCache,
		dataResolver:     mockDataResolver,
		routeConverter:   mockConverter,
	}

	return service, mockRegistry, mockCache, mockDataResolver, mockConverter
}

func TestOptimizedTemplateService_RenderComponent_CacheHit(t *testing.T) {
	service, _, mockCache, _, _ := createTestOTS(t)

	route := interfaces.Route{
		Path:         "/test",
		TemplateFile: "test.templ",
	}
	params := map[string]string{"id": "123"}
	ctx := context.Background()

	// Create mock RouterContext
	mockRouterCtx := &mockOTSRouterContext{
		ctx:       ctx,
		urlParams: params,
	}
	mockRouterCtx.On("GetAllURLParams").Return(params)
	mockRouterCtx.On("GetAllQueryParams").Return(url.Values{})

	// Mock cache hit
	cacheKey := "test-cache-key"
	mockComponent := mockOTSComponent{}
	mockCache.On("BuildTemplateKey", "test.templ", "", params).Return(cacheKey)
	mockCache.On("GetTemplate", cacheKey).Return(mockComponent, true)

	// Execute
	result, err := service.RenderComponent(route, mockRouterCtx, ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockComponent, result)
	mockCache.AssertExpectations(t)
	mockRouterCtx.AssertExpectations(t)
}

func TestOptimizedTemplateService_RenderComponent_DirectRouteMatch(t *testing.T) {
	service, mockRegistry, mockCache, _, _ := createTestOTS(t)

	route := interfaces.Route{
		Path:         "/test",
		TemplateFile: "test.templ",
	}
	params := map[string]string{}
	ctx := context.Background()

	// Create mock RouterContext
	mockRouterCtx := &mockOTSRouterContext{
		ctx:       ctx,
		urlParams: params,
	}
	mockRouterCtx.On("GetAllURLParams").Return(params)
	mockRouterCtx.On("GetAllQueryParams").Return(url.Values{})

	// Mock cache miss
	cacheKey := "test-cache-key"
	routeCacheKey := "route-cache-key"
	templateUUID := "template-123"
	
	mockCache.On("BuildTemplateKey", "test.templ", "", params).Return(cacheKey)
	mockCache.On("GetTemplate", cacheKey).Return(nil, false)
	mockCache.On("BuildRouteKey", "/test", params).Return(routeCacheKey)
	mockCache.On("GetRoute", routeCacheKey).Return(nil, false)

	// Mock direct route mapping
	routeMapping := map[string]string{"/test": templateUUID}
	mockRegistry.On("GetRouteToTemplateMapping").Return(routeMapping)
	
	// Mock template function that returns a parameterless function
	templateFunc := func() interface{} {
		return func() templ.Component {
			return mockOTSComponent{}
		}
	}
	mockRegistry.On("GetTemplateFunction", templateUUID).Return(templateFunc, true)
	
	// Mock cache operations
	mockCache.On("SetRoute", routeCacheKey, templateUUID).Return()
	mockCache.On("SetTemplate", cacheKey, mock.Anything).Return()

	// Execute
	result, err := service.RenderComponent(route, mockRouterCtx, ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRegistry.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockRouterCtx.AssertExpectations(t)
}

func TestOptimizedTemplateService_RenderComponent_ParameterizedTemplate_MissingID(t *testing.T) {
	service, mockRegistry, mockCache, _, _ := createTestOTS(t)

	route := interfaces.Route{
		Path:         "/user/123",
		TemplateFile: "user.templ",
	}
	params := map[string]string{} // Missing ID parameter
	ctx := context.Background()

	// Create mock RouterContext
	mockRouterCtx := &mockOTSRouterContext{
		ctx:       ctx,
		urlParams: params,
	}
	mockRouterCtx.On("GetAllURLParams").Return(params)
	mockRouterCtx.On("GetAllQueryParams").Return(url.Values{})
	mockRouterCtx.On("GetURLParam", "id").Return("") // Missing ID parameter

	// Mock cache miss
	cacheKey := "test-cache-key"
	routeCacheKey := "route-cache-key"
	templateUUID := "template-123"
	
	mockCache.On("BuildTemplateKey", "user.templ", "", params).Return(cacheKey)
	mockCache.On("GetTemplate", cacheKey).Return(nil, false)
	mockCache.On("BuildRouteKey", "/user/123", params).Return(routeCacheKey)
	mockCache.On("GetRoute", routeCacheKey).Return(nil, false)

	// Mock direct route mapping
	routeMapping := map[string]string{"/user/123": templateUUID}
	mockRegistry.On("GetRouteToTemplateMapping").Return(routeMapping)
	
	// Mock template function that returns a parameterized function
	templateFunc := func() interface{} {
		return func(id string) templ.Component {
			return mockOTSComponent{}
		}
	}
	mockRegistry.On("GetTemplateFunction", templateUUID).Return(templateFunc, true)

	// Mock the cache operations that happen before the validation error
	mockCache.On("SetRoute", routeCacheKey, templateUUID).Return()

	// Execute
	result, err := service.RenderComponent(route, mockRouterCtx, ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	// Check that it's a validation error
	var appErr *shared.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, shared.ErrorTypeValidation, appErr.Type)
	assert.Contains(t, appErr.Message, "parameter 'id' is required")
	
	mockRegistry.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockRouterCtx.AssertExpectations(t)
}

func TestOptimizedTemplateService_RenderComponent_RouteNotFound(t *testing.T) {
	service, mockRegistry, mockCache, _, mockConverter := createTestOTS(t)

	route := interfaces.Route{
		Path:         "/nonexistent",
		TemplateFile: "nonexistent.templ",
	}
	params := map[string]string{}
	ctx := context.Background()

	// Create mock RouterContext
	mockRouterCtx := &mockOTSRouterContext{
		ctx:       ctx,
		urlParams: params,
	}
	mockRouterCtx.On("GetAllURLParams").Return(params)
	mockRouterCtx.On("GetAllQueryParams").Return(url.Values{})

	// Mock cache miss
	cacheKey := "test-cache-key"
	routeCacheKey := "route-cache-key"
	
	mockCache.On("BuildTemplateKey", "nonexistent.templ", "", params).Return(cacheKey)
	mockCache.On("GetTemplate", cacheKey).Return(nil, false)
	mockCache.On("BuildRouteKey", "/nonexistent", params).Return(routeCacheKey)
	mockCache.On("GetRoute", routeCacheKey).Return(nil, false)

	// Mock empty route mapping
	routeMapping := map[string]string{}
	mockRegistry.On("GetRouteToTemplateMapping").Return(routeMapping)
	
	// Mock route converter returning no variations
	mockConverter.On("GenerateRouteVariations", "/nonexistent").Return([]string{})

	// Execute
	result, err := service.RenderComponent(route, mockRouterCtx, ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, middleware.ErrTemplateNotFound, err)
	
	mockRegistry.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockConverter.AssertExpectations(t)
	mockRouterCtx.AssertExpectations(t)
}

func TestOptimizedTemplateService_convertLayoutPathToRoute(t *testing.T) {
	service, _, _, _, _ := createTestOTS(t)

	tests := []struct {
		name         string
		layoutPath   string
		expected     string
		expectEmpty  bool
	}{
		{
			name:       "valid layout path",
			layoutPath: "/app/layout.templ",
			expected:   "/layout",
		},
		{
			name:        "empty layout path",
			layoutPath:  "",
			expectEmpty: true,
		},
		{
			name:        "invalid extension",
			layoutPath:  "/app/layout.html",
			expectEmpty: true,
		},
		{
			name:        "no filename",
			layoutPath:  "/app/.templ",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.convertLayoutPathToRoute(tt.layoutPath)
			
			if tt.expectEmpty {
				assert.Empty(t, result)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestOptimizedTemplateService_isDataServiceTemplate(t *testing.T) {
	service, _, _, _, _ := createTestOTS(t)

	// Mock data type for testing
	type UserData struct {
		ID   string
		Name string
	}

	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name: "valid DataService template function",
			input: func(*UserData) templ.Component {
				return mockOTSComponent{}
			},
			expected: true,
		},
		{
			name: "parameterless function",
			input: func() templ.Component {
				return mockOTSComponent{}
			},
			expected: false,
		},
		{
			name: "string parameter function",
			input: func(string) templ.Component {
				return mockOTSComponent{}
			},
			expected: false,
		},
		{
			name:     "not a function",
			input:    "not a function",
			expected: false,
		},
		{
			name: "function with wrong return type",
			input: func(*UserData) string {
				return "not a component"
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isDataServiceTemplate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOptimizedTemplateService_ClearCache(t *testing.T) {
	service, _, mockCache, _, _ := createTestOTS(t)

	mockCache.On("ClearAll").Return()

	// Execute
	service.ClearCache()

	// Assert
	mockCache.AssertExpectations(t)
}