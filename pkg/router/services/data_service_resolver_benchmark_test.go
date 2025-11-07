package services

import (
	"context"
	"testing"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/router/test_utils"
	"github.com/samber/do/v2"
)

// BenchmarkDataServiceResolverOriginal benchmarks the original reflection-based resolver
func BenchmarkDataServiceResolverOriginal(b *testing.B) {
	// Setup
	injector := do.New()
	registry := &test_utils.MockTemplateRegistry{}

	resolver := &dataServiceResolverImpl{
		injector:         injector,
		templateRegistry: registry,
	}

	// Register test service
	do.ProvideNamed(injector, "TestDataService", func(_ do.Injector) (any, error) {
		return &TestDataService{}, nil
	})

	// Mock registry to return service info
	test_utils.MockTemplateRegistryWithDataService(registry, "TestDataService")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := resolver.ResolveGenericDataService("TestDataService")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDataServiceResolverOptimized benchmarks the optimized reflection-free resolver
func BenchmarkDataServiceResolverOptimized(b *testing.B) {
	// Setup
	injector := do.New()
	registry := &test_utils.MockTemplateRegistry{}

	resolver := &OptimizedDataServiceResolver{
		injector:         injector,
		templateRegistry: registry,
		serviceRegistry:  make(map[string]interfaces.GenericDataService),
	}

	// Register test service in both DI and type-safe registry
	testService := &TestDataService{}
	do.ProvideNamed(injector, "TestDataService", func(_ do.Injector) (any, error) {
		return testService, nil
	})
	resolver.RegisterDataService("TestDataService", testService)

	// Mock registry to return service info
	test_utils.MockTemplateRegistryWithDataService(registry, "TestDataService")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := resolver.ResolveGenericDataService("TestDataService")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConfigurationValidation benchmarks the centralized validation
func BenchmarkConfigurationValidation(b *testing.B) {
	config := &TestConfig{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host: "localhost",
			Port: 5432,
		},
		Auth: AuthConfig{
			MinPasswordLength:    8,
			CreateDefaultAdmin:   false,
		},
		Email: EmailConfig{
			SMTPPort: 587,
		},
		Security: SecurityConfig{
			RateLimitRequests: 100,
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := config.Validate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test DataService for benchmarking
type TestDataService struct{}

func (s *TestDataService) GetData(routerCtx interfaces.RouterContext) (any, error) {
	return "test-data", nil
}

// Test configuration structures for benchmarking
type TestConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Email    EmailConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host string
	Port int
}

type AuthConfig struct {
	MinPasswordLength  int
	CreateDefaultAdmin bool
}

type EmailConfig struct {
	SMTPPort int
}

type SecurityConfig struct {
	RateLimitRequests int
}

func (c *TestConfig) Validate() error {
	// Use centralized validators
	portValidator := NewPortValidator("")
	positiveValidator := NewPositiveNumberValidator("", 1)

	// Server port
	portValidator.fieldName = "server.port"
	if err := portValidator.ValidatePort(c.Server.Port); err != nil {
		return err
	}

	// Database port
	portValidator.fieldName = "database.port"
	if err := portValidator.ValidatePort(c.Database.Port); err != nil {
		return err
	}

	// Auth
	positiveValidator.fieldName = "auth.min_password_length"
	if err := positiveValidator.ValidatePositiveNumber(c.Auth.MinPasswordLength); err != nil {
		return err
	}

	// Email
	portValidator.fieldName = "email.smtp_port"
	if err := portValidator.ValidatePort(c.Email.SMTPPort); err != nil {
		return err
	}

	// Security
	positiveValidator.fieldName = "security.rate_limit_requests"
	if err := positiveValidator.ValidatePositiveNumber(c.Security.RateLimitRequests); err != nil {
		return err
	}

	return nil
}