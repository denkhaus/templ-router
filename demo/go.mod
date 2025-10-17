module github.com/denkhaus/templ-router/demo

go 1.25.0

require (
	github.com/Oudwins/tailwind-merge-go v0.2.1
	github.com/a-h/templ v0.3.960
	github.com/denkhaus/templ-router v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.2.3
	github.com/samber/do/v2 v2.0.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/samber/go-type-to-string v1.8.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/denkhaus/templ-router => ../
