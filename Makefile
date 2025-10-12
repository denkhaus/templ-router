.PHONY: generate generate_templ generate_registry

# Generate templ files
generate_templ:
	@cd demo && templ generate

# Generate template registry
generate_registry:
	@cd demo && TEMPLATE_SCAN_PATH=app TEMPLATE_OUTPUT_DIR=generated/templates TEMPLATE_MODULE_NAME=github.com/denkhaus/templ-router/demo go run ../cmd/template-generator/*.go

# Generate all
generate: generate_templ generate_registry
