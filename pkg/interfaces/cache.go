package interfaces

// CacheService provides template and route caching functionality
type CacheService interface {
	// Template cache operations
	GetTemplate(key string) (interface{}, bool)
	SetTemplate(key string, value interface{})
	
	// Route cache operations
	GetRoute(key string) (interface{}, bool)
	SetRoute(key string, value interface{})
	
	// Cache management
	ClearAll()
	ClearTemplates()
	ClearRoutes()
	
	// Cache key generation
	BuildTemplateKey(templateFile, locale string, params map[string]string) string
	BuildRouteKey(path string, params map[string]string) string
}