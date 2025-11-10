package config

import (
	"os"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsProduction(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		isProduction bool
	}{
		{
			name:         "default development",
			envVars:      map[string]string{},
			isProduction: false,
		},
		{
			name: "production by non-localhost base URL and custom CSRF secret",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":     "production",
				"TR_SERVER_BASE_URL":      "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "production-secret-key",
			},
			isProduction: true,
		},
		{
			name: "development by environment kind override",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":     "develop",
				"TR_SERVER_BASE_URL":      "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "production-secret-key",
			},
			isProduction: false,
		},
		{
			name: "development with localhost base URL",
			envVars: map[string]string{
				"TR_SERVER_BASE_URL":      "http://localhost:8080",
				"TR_SECURITY_CSRF_SECRET": "production-secret-key",
			},
			isProduction: false,
		},
		{
			name: "development with default CSRF secret",
			envVars: map[string]string{
				"TR_SERVER_BASE_URL":      "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "change-me-in-production",
			},
			isProduction: false,
		},
		{
			name: "production with staging environment",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":     "staging",
				"TR_SERVER_BASE_URL":      "https://staging.myapp.com",
				"TR_SECURITY_CSRF_SECRET": "staging-secret-key",
			},
			isProduction: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)

			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set environment variable %s: %v", key, err)
				}
			}

			injector := do.New()
			defer func() {
				_ = injector.Shutdown()
			}()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			// Access the internal config to test the helper method
			configSvc := service.(*configService)

			assert.Equal(t, tt.isProduction, configSvc.config.IsProduction())
			assert.Equal(t, !tt.isProduction, configSvc.config.IsDevelopment())
		})
	}
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		isDevelopment bool
	}{
		{
			name:          "default development",
			envVars:       map[string]string{},
			isDevelopment: true,
		},
		{
			name: "production environment",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":     "production",
				"TR_SERVER_BASE_URL":      "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "production-secret-key",
			},
			isDevelopment: false,
		},
		{
			name: "explicit development environment",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND": "develop",
			},
			isDevelopment: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)

			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set environment variable %s: %v", key, err)
				}
			}

			injector := do.New()
			defer func() {
				_ = injector.Shutdown()
			}()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			// Access the internal config to test the helper method
			configSvc := service.(*configService)

			assert.Equal(t, tt.isDevelopment, configSvc.config.IsDevelopment())
		})
	}
}
