package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test namespace for testing-related commands
type Test mg.Namespace

// E2E runs end-to-end tests against the Docker service
func (p Test) E2E() error {
	mg.Deps(p.CheckService)

	fmt.Println("🧪 Running E2E tests against Docker service...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --randomize-all --randomize-suites --race --trace")
}

// E2EWatch runs E2E tests in watch mode for development
func (p Test) E2EWatch() error {
	mg.Deps(p.CheckService)

	fmt.Println("🔄 Running E2E tests in watch mode...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo watch --randomize-all --randomize-suites --race --trace")
}

// E2ESmoke runs quick smoke tests
func (p Test) E2ESmoke() error {
	mg.Deps(p.CheckService)

	fmt.Println("💨 Running E2E smoke tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Health Check|Language Selection' --randomize-all --race --trace")
}

// E2ERouting runs routing-specific tests
func (p Test) E2ERouting() error {
	mg.Deps(p.CheckService)

	fmt.Println("🛣️  Running routing tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Multi-Language Routing|Dynamic Routes' --randomize-all --race --trace")
}

// E2EI18n runs internationalization tests
func (p Test) E2EI18n() error {
	mg.Deps(p.CheckService)

	fmt.Println("🌐 Running i18n tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Language' --randomize-all --race --trace")
}

// E2EPerf runs performance tests
func (p Test) E2EPerf() error {
	mg.Deps(p.CheckService)

	fmt.Println("⚡ Running performance tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Performance' --randomize-all --race --trace")
}

// E2EAuth runs authentication tests
func (p Test) E2EAuth() error {
	mg.Deps(p.CheckService)

	fmt.Println("🔐 Running authentication tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Authentication' --randomize-all --race --trace")
}

// E2EData runs data service tests
func (p Test) E2EData() error {
	mg.Deps(p.CheckService)

	fmt.Println("💾 Running data service tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Data Service|TestProduct|TestSpecific' --randomize-all --race --trace")
}

// E2EDataServiceI18n runs DataService + i18n integration tests
func (p Test) E2EDataServiceI18n() error {
	mg.Deps(p.CheckService)

	fmt.Println("🌐💾 Running DataService + i18n integration tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='TestProduct.*i18n|TestSpecific.*i18n' --randomize-all --race --trace")
}

// E2EUserWithIdData runs UserWithId DataService tests
func (p Test) E2EUserWithIdData() error {
	mg.Deps(p.CheckService)

	fmt.Println("👤💾 Running UserWithId DataService tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='UserWithId DataService Integration' --randomize-all --race --trace")
}

// E2EContent runs content validation tests
func (p Test) E2EContent() error {
	mg.Deps(p.CheckService)

	fmt.Println("📄 Running content validation tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Content Validation' --randomize-all --race --trace")
}

// E2EMetadata runs metadata and layout tests
func (p Test) E2EMetadata() error {
	mg.Deps(p.CheckService)

	fmt.Println("🏗️ Running metadata and layout tests...")

	return sh.RunWithV(map[string]string{
		"TEST_BASE_URL": "http://localhost:8084",
	}, "sh", "-c", "cd demo/tests && ginkgo run --focus='Metadata and Layout' --randomize-all --race --trace")
}

// SetupE2E installs E2E test dependencies
func (p Test) SetupE2E() error {
	fmt.Println("📦 Installing E2E test dependencies...")

	// Install ginkgo CLI
	if err := sh.Run("go", "install", "github.com/onsi/ginkgo/v2/ginkgo@latest"); err != nil {
		return fmt.Errorf("failed to install ginkgo: %w", err)
	}

	// Download test dependencies
	if err := sh.RunV("sh", "-c", "cd demo/tests && go mod download"); err != nil {
		return fmt.Errorf("failed to download test dependencies: %w", err)
	}

	fmt.Println("✅ E2E test dependencies installed")
	return nil
}

// CheckService verifies that the Docker service is running
func (p Test) CheckService() error {
	fmt.Println("🔍 Checking if Docker service is running...")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:8084/api/health")
	if err != nil {
		return fmt.Errorf("❌ Docker service not running on localhost:8084. Start with: mage docker:up")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("❌ Service unhealthy (status: %d). Check with: mage docker:logs", resp.StatusCode)
	}

	fmt.Println("✅ Docker service is running and healthy")
	return nil
}

// All runs all tests (unit + E2E)
func (p Test) All() error {
	fmt.Println("🧪 Running all tests...")

	// Run generator tests
	if err := (Generator{}).Test(); err != nil {
		return fmt.Errorf("generator tests failed: %w", err)
	}

	// Run E2E tests
	if err := p.E2E(); err != nil {
		return fmt.Errorf("E2E tests failed: %w", err)
	}

	fmt.Println("✅ All tests passed!")
	return nil
}

// CI runs tests in CI mode with coverage
func (p Test) CI() error {
	fmt.Println("🤖 Running tests in CI mode...")

	// Run generator tests with coverage
	if err := (Generator{}).TestCoverage(); err != nil {
		return fmt.Errorf("generator tests failed: %w", err)
	}

	// Run E2E tests
	if err := p.E2E(); err != nil {
		return fmt.Errorf("E2E tests failed: %w", err)
	}

	fmt.Println("✅ CI tests completed!")
	return nil
}

// DevSetup sets up complete development testing environment
func (p Test) DevSetup() error {
	fmt.Println("🚀 Setting up development testing environment...")

	// Install E2E dependencies
	if err := p.SetupE2E(); err != nil {
		return err
	}

	// Start Docker service
	if err := (Docker{}).Up(); err != nil {
		return fmt.Errorf("failed to start Docker service: %w", err)
	}

	// Wait for service to be ready
	fmt.Println("⏳ Waiting for service to be ready...")
	time.Sleep(10 * time.Second)

	// Check service
	if err := p.CheckService(); err != nil {
		return err
	}

	// Run smoke test
	if err := p.E2ESmoke(); err != nil {
		return fmt.Errorf("smoke test failed: %w", err)
	}

	fmt.Println("✅ Development testing environment ready!")
	fmt.Println("💡 Available commands:")
	fmt.Println("   mage test:e2e               - Run all E2E tests")
	fmt.Println("   mage test:e2eWatch          - Watch mode for development")
	fmt.Println("   mage test:e2eSmoke          - Quick smoke tests")
	fmt.Println("   mage test:e2eRouting        - Routing tests")
	fmt.Println("   mage test:e2eI18n           - i18n tests")
	fmt.Println("   mage test:e2eData           - DataService tests")
	fmt.Println("   mage test:e2eDataServiceI18n - DataService + i18n integration tests")
	fmt.Println("   mage test:e2eUserWithIdData - UserWithId DataService tests")

	return nil
}
