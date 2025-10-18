package di

import (
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)

func TestWithTemplateRegistry(t *testing.T) {
	container := NewContainer()
	mockRegistry := &mockTemplateRegistry{}
	
	option := WithTemplateRegistry(mockRegistry)
	
	// Apply option
	option(container)
	
	// Verify registry is registered
	retrievedRegistry := do.MustInvoke[interfaces.TemplateRegistry](container.injector)
	if retrievedRegistry == nil {
		t.Error("WithTemplateRegistry option did not register template registry correctly")
	}
}

func TestWithAssetsService(t *testing.T) {
	container := NewContainer()
	mockAssets := &mockAssetsService{}
	
	option := WithAssetsService(mockAssets)
	
	// Apply option
	option(container)
	
	// Verify assets service is registered
	retrievedAssets := do.MustInvoke[interfaces.AssetsService](container.injector)
	if retrievedAssets == nil {
		t.Error("WithAssetsService option did not register assets service correctly")
	}
}

func TestMultipleOptions(t *testing.T) {
	container := NewContainer()
	mockRegistry := &mockTemplateRegistry{}
	mockAssets := &mockAssetsService{}
	
	// Apply multiple options
	container.RegisterApplicationServices(
		WithTemplateRegistry(mockRegistry),
		WithAssetsService(mockAssets),
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
