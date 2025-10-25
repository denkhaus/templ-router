# Templ Router Demo E2E Tests

Comprehensive End-to-End tests for the Templ Router Demo application using Agouti (Selenium WebDriver for Go).

## 🎯 Purpose

This test suite ensures that all Templ Router features work correctly after changes to the router library. Tests run on the **host system** against the **Docker Compose service**.

## 🏗️ Architecture

```
Host System (Tests)  ←→  Docker Container (Demo App)
     ↓                        ↓
  Agouti/Chrome         localhost:8084
```

## 🚀 Quick Start

### 1. Start Docker Service
```bash
cd ../docker
docker-compose -f docker-compose.dev.yml up
```

### 2. Run Tests
```bash
cd demo/tests
make install
make test
```

## 📋 Test Coverage

### ✅ **Core Router Features**
- [x] Multi-language routing (`/en/*`, `/de/*`)
- [x] Dynamic routes (`/user/{id}`, `/product/{id}`)
- [x] Static routes (`/login`, `/signup`, `/admin`)
- [x] Language switching via URL
- [x] Template discovery and rendering

### ✅ **I18n Features**
- [x] Locale-specific content rendering
- [x] Currency formatting (€ vs $)
- [x] Date formatting
- [x] Translation key resolution
- [x] Language switcher functionality

### ✅ **Navigation & UX**
- [x] Navbar navigation
- [x] Language selection landing page
- [x] Responsive design
- [x] Error handling (404, invalid locales)

### ✅ **Performance & Reliability**
- [x] Page load times
- [x] Service health checks
- [x] Mobile viewport testing

## 🛠️ Available Commands

```bash
# Basic usage
make install       # Install dependencies
make test          # Run all tests
make test-verbose  # Verbose output
make test-watch    # Watch mode for development

# Specific test categories
make test-routing  # Multi-language & dynamic routing
make test-i18n     # Internationalization features
make test-navigation # Navigation & UX

# Service management
make check-service # Verify Docker service is running
make docker-up     # Start Docker service
make docker-down   # Stop Docker service
make docker-logs   # View service logs

# Development
make smoke         # Quick smoke tests
make report        # Generate JUnit report
make clean         # Clean test artifacts
```

## 🔧 Configuration

### Environment Variables
```bash
# Test configuration
AGOUTI_TIMEOUT=10s
AGOUTI_HEADLESS=true
DEMO_BASE_URL=http://localhost:8084

# Chrome options
CHROME_ARGS="--headless,--no-sandbox,--disable-dev-shm-usage"
```

### Test Timeouts
- **Page Load**: 5 seconds
- **Element Wait**: 10 seconds
- **Service Ready**: 30 seconds

## 📊 Test Structure

```
tests/
├── e2e_test.go          # Main test suite
├── go.mod               # Test dependencies
├── Makefile             # Test automation
└── README.md            # This file
```

### Test Categories

1. **Health Check**: Service availability
2. **Language Selection**: Landing page functionality
3. **Multi-Language Routing**: EN/DE route handling
4. **Dynamic Routes**: Parameter-based routing
5. **Language Switching**: Navbar language switcher
6. **Navigation**: Inter-page navigation
7. **Authentication**: Login/signup pages
8. **Error Handling**: 404 and invalid routes
9. **Responsive Design**: Mobile viewport
10. **Performance**: Load time validation

## 🐛 Debugging

### View Test Output
```bash
make test-verbose
```

### Run Specific Tests
```bash
# Focus on routing tests
ginkgo run --focus="Multi-Language Routing"

# Focus on specific test
ginkgo run --focus="should load English dashboard"
```

### Check Service Logs
```bash
make docker-logs
```

### Manual Service Check
```bash
curl http://localhost:8084/api/health
```

## 🔄 CI/CD Integration

### GitHub Actions Example
```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      - name: Start Demo Service
        run: |
          cd docker
          docker-compose -f docker-compose.dev.yml up -d
      
      - name: Run E2E Tests
        run: |
          cd demo/tests
          make ci
      
      - name: Upload Test Results
        uses: actions/upload-artifact@v4
        with:
          name: test-results
          path: demo/tests/test-results.xml
```

## 📈 Extending Tests

### Add New Route Test
```go
It("should handle new feature route", func() {
    Expect(page.Navigate(BaseURL + "/en/new-feature")).To(Succeed())
    Eventually(page.HTML, PageLoadTimeout).Should(ContainSubstring("New Feature"))
})
```

### Add New Language Test
```go
It("should support French locale", func() {
    Expect(page.Navigate(BaseURL + "/fr/dashboard")).To(Succeed())
    Eventually(page.HTML, PageLoadTimeout).Should(ContainSubstring("Tableau de bord"))
})
```

## 🚨 Troubleshooting

### Common Issues

1. **Service Not Running**
   ```bash
   make check-service
   # If fails: make docker-up
   ```

2. **Chrome Not Found**
   ```bash
   # Install Chrome/Chromium
   sudo apt-get install chromium-browser  # Ubuntu
   brew install --cask google-chrome      # macOS
   ```

3. **Port Conflicts**
   ```bash
   # Check what's using port 8084
   lsof -i :8084
   ```

4. **Test Timeouts**
   - Increase timeouts in `e2e_test.go`
   - Check Docker service performance
   - Verify network connectivity

### Debug Mode
```bash
# Run with Chrome visible (non-headless)
AGOUTI_HEADLESS=false make test-verbose
```

## 📝 Best Practices

1. **Always check service health** before running tests
2. **Use descriptive test names** that explain the feature being tested
3. **Group related tests** in contexts
4. **Set appropriate timeouts** for different operations
5. **Clean up resources** after each test
6. **Test both positive and negative scenarios**
7. **Include performance validations** where relevant

This test suite ensures your Templ Router changes don't break existing functionality! 🚀