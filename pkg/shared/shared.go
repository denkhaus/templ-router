package shared

type ContextType string

const (
	UserContextKey        ContextType = "user"
	LocaleKey             ContextType = "locale"
	TemplateConfigKey     ContextType = "template_config"
	TemplatePathKey       ContextType = "template_path"
	I18nDataKey           ContextType = "router_i18n_data"
	I18nTemplateKey       ContextType = "router_i18n_template"
	RequestPathKey        ContextType = "request_path"
	RouteMappingKey       ContextType = "route_mapping"
	ComponentsMetadataKey ContextType = "components_metadata"
)
