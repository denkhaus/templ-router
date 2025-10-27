package config

import (
	"os"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatabaseDSN(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "default values",
			envVars:  map[string]string{},
			expected: "host=localhost port=5432 user=postgres password=postgres dbname=router_db sslmode=disable",
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"TR_DATABASE_HOST":     "db.example.com",
				"TR_DATABASE_PORT":     "3306",
				"TR_DATABASE_USER":     "myuser",
				"TR_DATABASE_PASSWORD": "mypass",
				"TR_DATABASE_NAME":     "mydb",
				"TR_DATABASE_SSL_MODE": "require",
			},
			expected: "host=db.example.com port=3306 user=myuser password=mypass dbname=mydb sslmode=require",
		},
		{
			name: "with special characters in password",
			envVars: map[string]string{
				"TR_DATABASE_PASSWORD": "p@ssw0rd!#$",
			},
			expected: "host=localhost port=5432 user=postgres password=p@ssw0rd!#$ dbname=router_db sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)
			
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			// Access the internal config to test the helper method
			configSvc := service.(*configService)
			dsn := configSvc.config.GetDatabaseDSN()
			
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestGetServerAddress(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "default values",
			envVars:  map[string]string{},
			expected: "localhost:8080",
		},
		{
			name: "custom host and port",
			envVars: map[string]string{
				"TR_SERVER_HOST": "0.0.0.0",
				"TR_SERVER_PORT": "3000",
			},
			expected: "0.0.0.0:3000",
		},
		{
			name: "IPv6 address",
			envVars: map[string]string{
				"TR_SERVER_HOST": "::1",
				"TR_SERVER_PORT": "8443",
			},
			expected: "::1:8443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)
			
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			// Access the internal config to test the helper method
			configSvc := service.(*configService)
			address := configSvc.config.GetServerAddress()
			
			assert.Equal(t, tt.expected, address)
		})
	}
}

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
				"TR_ENVIRONMENT_KIND":    "production",
				"TR_SERVER_BASE_URL":     "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "production-secret-key",
			},
			isProduction: true,
		},
		{
			name: "development by environment kind override",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":    "develop",
				"TR_SERVER_BASE_URL":     "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "production-secret-key",
			},
			isProduction: false,
		},
		{
			name: "development with localhost base URL",
			envVars: map[string]string{
				"TR_SERVER_BASE_URL":     "http://localhost:8080",
				"TR_SECURITY_CSRF_SECRET": "production-secret-key",
			},
			isProduction: false,
		},
		{
			name: "development with default CSRF secret",
			envVars: map[string]string{
				"TR_SERVER_BASE_URL":     "https://myapp.com",
				"TR_SECURITY_CSRF_SECRET": "change-me-in-production",
			},
			isProduction: false,
		},
		{
			name: "production with staging environment",
			envVars: map[string]string{
				"TR_ENVIRONMENT_KIND":    "staging",
				"TR_SERVER_BASE_URL":     "https://staging.myapp.com",
				"TR_SECURITY_CSRF_SECRET": "staging-secret-key",
			},
			isProduction: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearTestEnv(t)
			
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

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
				"TR_ENVIRONMENT_KIND":    "production",
				"TR_SERVER_BASE_URL":     "https://myapp.com",
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
				os.Setenv(key, value)
			}

			injector := do.New()
			defer injector.Shutdown()

			configFactory := NewConfigService("TR")
			service, err := configFactory(injector)
			require.NoError(t, err)

			// Access the internal config to test the helper method
			configSvc := service.(*configService)
			
			assert.Equal(t, tt.isDevelopment, configSvc.config.IsDevelopment())
		})
	}
}