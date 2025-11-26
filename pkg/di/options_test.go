package di

import (
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)

func TestWithTemplateRegistry(t *testing.T) {
	container := NewContainer()
	mockRegistry := &mockTemplateRegistry{}

	option := WithTemplateRegistryFactory(func(i do.Injector) (interfaces.TemplateRegistry, error) {
		return mockRegistry, nil
	})

	// Apply option
	option(container)

	// Verify registry is registered
	retrievedRegistry := do.MustInvoke[interfaces.TemplateRegistry](container.injector)
	if retrievedRegistry == nil {
		t.Error("WithTemplateRegistryFactory option did not register template registry correctly")
	}
}

func TestWithAssetsServiceFactory(t *testing.T) {
	container := NewContainer()

	option := WithAssetsServiceFactory(func(do.Injector) (interfaces.AssetsService, error) {
		return &mockAssetsService{}, nil
	})

	// Apply option
	option(container)

	// Verify assets service is registered
	retrievedAssets := do.MustInvoke[interfaces.AssetsService](container.injector)
	if retrievedAssets == nil {
		t.Error("WithAssetsServiceFactory option did not register assets service correctly")
	}
}

func TestMultipleOptions(t *testing.T) {
	container := NewContainer()
	mockRegistry := &mockTemplateRegistry{}

	// Apply multiple options
	container.RegisterApplicationServices(
		WithTemplateRegistryFactory(func(do.Injector) (interfaces.TemplateRegistry, error) {
			return mockRegistry, nil
		}),
		WithAssetsServiceFactory(func(do.Injector) (interfaces.AssetsService, error) {
			return &mockAssetsService{}, nil
		}),
	)

	// Verify both services are registered
	retrievedRegistry := do.MustInvoke[interfaces.TemplateRegistry](container.injector)
	if retrievedRegistry == nil {
		t.Error("Template registry not registered when using multiple options")
	}

	retrievedAssets := do.MustInvoke[interfaces.AssetsService](container.injector)
	if retrievedAssets == nil {
		t.Error("Assets service not registered when using multiple options")
	}
}

func TestOptionsPattern(t *testing.T) {
	// Test that options are functions that modify the container
	container := NewContainer()

	var optionCalled bool
	testOption := func(c *Container) {
		optionCalled = true
		if c != container {
			t.Error("Option function received wrong container")
		}
	}

	container.RegisterApplicationServices(testOption)

	if !optionCalled {
		t.Error("Option function was not called")
	}
}
