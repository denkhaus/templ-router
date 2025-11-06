# Contributing to Templ Router

**Guide to contributing to the Templ Router project.**

Thank you for your interest in contributing to Templ Router! This document provides information on how to contribute to the project, including development setup, coding standards, and submission guidelines.

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+** with modern Go features
- **Git** for version control
- **Make** or **Mage** for build automation
- **Docker** (optional, for testing)

### Development Setup

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/yourusername/templ-router.git
cd templ-router

# 3. Add the original repository as upstream
git remote add upstream https://github.com/denkhaus/templ-router.git

# 4. Install dependencies
go mod tidy

# 5. Install development tools
mage generator:install
mage test:devSetup

# 6. Start development server
mage dev
```

## 📁 Project Structure

Understanding the project structure is essential for contributing:

```
templ-router/
├── cmd/                       # CLI applications
│   └── trgen/                 # Template registry generator
├── demo/                      # Demo application
│   ├── app/                   # Demo templates
│   ├── pkg/                   # Demo packages
│   └── main.go                # Demo entry point
├── docs/                      # Documentation
├── magefiles/                 # Build automation
├── pkg/                       # Library source code
│   ├── di/                    # Dependency injection
│   ├── interfaces/            # Core interfaces
│   ├── router/                # Router implementation
│   │   ├── middleware/        # Middleware components
│   │   ├── pipeline/          # Request pipeline
│   │   └── services/          # Router services
│   ├── services/              # Core services
│   └── shared/                # Shared utilities
├── .github/                   # GitHub configuration
│   ├── workflows/             # CI/CD workflows
│   └── ISSUE_TEMPLATE/        # Issue templates
└── magefile.go               # Main build file
```

## 🛠️ Development Workflow

### 1. Create a Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create a feature branch
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-fix-name
```

### 2. Make Changes

Follow the coding standards and guidelines outlined in this document.

### 3. Test Your Changes

```bash
# Run all tests
mage test:all

# Run specific test categories
mage test:unit               # Unit tests only
mage test:e2e                # E2E tests only
mage test:integration        # Integration tests

# Run tests with coverage
mage test:coverage

# Run linter
mage lint

# Run code formatting
mage fmt
```

### 4. Commit Your Changes

```bash
# Stage changes
git add .

# Commit with conventional commit format
git commit -m "feat: add new authentication middleware"

# Or for bug fixes
git commit -m "fix: resolve template caching issue"
```

### 5. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create a pull request on GitHub
# Provide a clear description and link relevant issues
```

## 📝 Coding Standards

### Go Conventions

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html) guidelines.

### Code Formatting

```bash
# Format code
gofmt -s -w .

# Import organization
goimports -w .

# Run via mage
mage fmt
```

### Naming Conventions

```go
// Packages: lowercase, single words when possible
package router
package middleware

// Interfaces: descriptive names, often ending in -er
type TemplateRenderer interface{}
type RouteRegistrar interface{}

// Structs: PascalCase, descriptive names
type TemplateRegistry struct{}
type AuthenticationMiddleware struct{}

// Methods: PascalCase for exported, camelCase for unexported
func (r *TemplateRegistry) RegisterRoute() {}  // Exported
func (r *TemplateRegistry) validateRoute() {} // Unexported

// Constants: UPPER_SNAKE_CASE
const DEFAULT_PORT = 8080
const SESSION_COOKIE_NAME = "session_id"

// Variables: camelCase
var defaultConfig = &Config{}
```

### Documentation

```go
// Package-level documentation
package router

// Router provides file-based routing functionality for templ templates.
// It automatically generates routes from file structure and supports
// internationalization, authentication, and data service injection.
package router

// Function documentation
// NewRouter creates a new Router instance with the provided configuration.
// It returns an error if the configuration is invalid.
func NewRouter(config *Config) (*Router, error) {
    // Implementation
}

// Method documentation
// RegisterRoutes registers all discovered routes with the provided HTTP router.
// The scanPath parameter specifies the directory containing template files.
func (r *Router) RegisterRoutes(mux http.Handler) error {
    // Implementation
}
```

### Error Handling

```go
// Use explicit error handling
func (s *Service) ProcessData(id string) (*Result, error) {
    if id == "" {
        return nil, fmt.Errorf("data ID cannot be empty")
    }

    data, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find data: %w", err)
    }

    result, err := s.process(data)
    if err != nil {
        return nil, fmt.Errorf("failed to process data: %w", err)
    }

    return result, nil
}

// Wrap errors with context
return nil, fmt.Errorf("failed to register route %s: %w", route.Path, err)
```

## 🧪 Testing Guidelines

### Test Structure

```go
// pkg/router/router_test.go
package router

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRouter_RegisterRoute(t *testing.T) {
    // Arrange
    router := NewRouter(&Config{})
    route := &Route{Path: "/test", Handler: testHandler}

    // Act
    err := router.RegisterRoute(route)

    // Assert
    require.NoError(t, err)
    assert.Contains(t, router.routes, "/test")
}

func TestRouter_RegisterRoute_InvalidPath(t *testing.T) {
    // Arrange
    router := NewRouter(&Config{})
    route := &Route{Path: "", Handler: testHandler}

    // Act
    err := router.RegisterRoute(route)

    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "path cannot be empty")
}
```

### Test Organization

```bash
pkg/
├── router/
│   ├── router.go
│   ├── router_test.go          # Unit tests for router
│   ├── router_integration_test.go  # Integration tests
│   └── router_example_test.go      # Example tests
├── middleware/
│   ├── auth_middleware.go
│   └── auth_middleware_test.go
└── services/
    ├── template_service.go
    └── template_service_test.go
```

### Mock Testing

```go
// Use testify/mock for complex dependencies
type MockTemplateService struct {
    mock.Mock
}

func (m *MockTemplateService) RenderTemplate(name string, data interface{}) (string, error) {
    args := m.Called(name, data)
    return args.String(0), args.Error(1)
}

func TestHandler_RenderTemplate(t *testing.T) {
    // Arrange
    mockService := &MockTemplateService{}
    handler := NewHandler(mockService)

    mockService.On("RenderTemplate", "index", map[string]interface{}{}).
        Return("<html>...</html>", nil)

    // Act
    result, err := handler.RenderTemplate("index", map[string]interface{}{})

    // Assert
    require.NoError(t, err)
    assert.Equal(t, "<html>...</html>", result)
    mockService.AssertExpectations(t)
}
```

## 📋 Commit Guidelines

### Conventional Commits

We use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code formatting, missing semicolons, etc. (no functional changes)
- `refactor`: Code refactoring (no functional changes)
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates, etc.

### Examples

```bash
feat(auth): add JWT token support
fix(router): resolve route precedence issue
docs(i18n): update translation examples
style(format): apply gofmt to all files
refactor(middleware): simplify auth middleware logic
test(data): add unit tests for data service
perf(template): implement template caching
chore(deps): update chi/v5 to v5.0.10
```

### Commit Message Guidelines

- Use the imperative mood ("add" not "added" or "adds")
- Keep the first line under 50 characters
- Wrap the body at 72 characters
- Explain what and why, not how
- Reference relevant issues in the footer

## 🐛 Bug Reports

### Before Creating a Bug Report

1. **Check existing issues** - Search for similar reports
2. **Check if it's fixed** - Try the latest version
3. **Create minimal reproduction** - Simplify the issue to its core

### Bug Report Template

```markdown
## Bug Description
A clear and concise description of the bug.

## To Reproduce
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## Expected Behavior
A clear description of what you expected to happen.

## Actual Behavior
A clear description of what actually happened.

## Environment
- OS: [e.g. macOS 13.0, Ubuntu 20.04]
- Go version: [e.g. 1.21.0]
- Templ Router version: [e.g. v0.3.0]

## Additional Context
Add any other context about the problem here.

## Possible Solution
If you have ideas on how to fix it, include them here.
```

## 💡 Feature Requests

### Feature Request Guidelines

1. **Check existing issues** - Search for similar requests
2. **Consider the scope** - Is this aligned with project goals?
3. **Provide use cases** - Explain why this feature is needed

### Feature Request Template

```markdown
## Feature Description
A clear and concise description of the feature.

## Problem Statement
What problem does this feature solve? What limitations does it address?

## Proposed Solution
Describe the solution you'd like to see implemented.

## Alternatives Considered
Describe any alternative solutions or features you've considered.

## Use Cases
Provide specific examples of how this feature would be used.

## Additional Context
Add any other context, mockups, or examples about the feature request.
```

## 📚 Documentation Contributions

### Documentation Types

- **Code Documentation**: Comments and Go doc strings
- **User Documentation**: Guides, tutorials, and API reference
- **Examples**: Sample code and use cases

### Documentation Guidelines

- **Keep it up to date** - Update docs when code changes
- **Use clear language** - Write for your target audience
- **Include examples** - Show, don't just tell
- **Test examples** - Ensure code examples work

### Documentation Structure

```bash
docs/
├── GETTING-STARTED.md       # New user guide
├── FILE-BASED-ROUTING.md    # Feature documentation
├── AUTHENTICATION.md        # Feature documentation
├── INTERNATIONALIZATION.md   # Feature documentation
├── DATA-SERVICES.md         # Feature documentation
├── CONFIGURATION.md         # Reference documentation
├── CONTRIBUTING.md          # This file
└── REFERENCE.md            # API reference
```

## 🔧 Development Tools

### Mage Commands

```bash
# Development
mage dev                    # Start development server
mage dev:setup              # Setup development environment

# Building
mage build                  # Build library
mage build:all              # Build for all platforms
mage generator:build        # Build trgen tool

# Testing
mage test                   # Run tests
mage test:all               # Run all tests
mage test:unit              # Unit tests only
mage test:e2e               # E2E tests only
mage test:coverage          # Tests with coverage

# Code Quality
mage fmt                    # Format code
mage lint                   # Run linter
mage vet                    # Run go vet

# Dependencies
mage deps:update            # Update dependencies
mage deps:tidy              # Clean dependencies

# Release
mage release:prepare        # Prepare release
mage release:publish        # Publish release
```

### IDE Configuration

#### VS Code

```json
// .vscode/settings.json
{
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.vetOnSave": "workspace",
    "go.buildOnSave": "workspace",
    "go.lintOnSave": "workspace",
    "go.testOnSave": "workspace"
}
```

#### Vim/Neovim

```vim
" vimrc for Go development
set nocompatible
filetype plugin indent on
syntax on

" Go settings
let g:go_fmt_command = "goimports"
let g:go_fmt_autosave = 1
let g:go vet_on_save = 1
let g:go_build_on_save = 0
let g:go_test_on_save = 0
```

## 🚀 Release Process

### Version Management

We follow [Semantic Versioning](https://semver.org/) (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

```bash
# 1. Update version
# Update version in magefiles/version.go
# Update CHANGELOG.md

# 2. Run full test suite
mage test:all

# 3. Build for all platforms
mage build:all

# 4. Test release candidates
mage release:test

# 5. Create release tag
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3

# 6. Publish release
mage release:publish
```

## 🤝 Community Guidelines

### Code of Conduct

We are committed to providing a welcoming and inclusive environment. Please:

- Be respectful and considerate
- Use inclusive language
- Focus on constructive feedback
- Help others learn and grow

### Getting Help

- **GitHub Discussions**: For questions and ideas
- **GitHub Issues**: For bugs and feature requests
- **Documentation**: Check existing docs first

### Communication Channels

- **Issues**: For bug reports and feature requests
- **Discussions**: For general questions and ideas
- **Pull Requests**: For code contributions

## 🏆 Recognition

### Contributors

All contributors are recognized in:

- **README.md**: Contributors section
- **CHANGELOG.md**: Release notes
- **GitHub**: Contributor statistics

### Types of Contributions

We value all types of contributions:

- **Code**: Patches, new features, bug fixes
- **Documentation**: Guides, tutorials, examples
- **Design**: UI/UX improvements, graphics
- **Testing**: Bug reports, test cases
- **Community**: Support, feedback, ideas

## 📖 Learning Resources

### Go Resources

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Blog](https://blog.golang.org/)
- [Go by Example](https://gobyexample.com/)

### Testing Resources

- [Testable Go](https://testablecode.com/posts/go-testing/)
- [Go Testing Tutorial](https://golang.org/doc/tutorial/add-a-test)
- [Testify Documentation](https://github.com/stretchr/testify)

### Project Specific

- [Templ Documentation](https://templ.guide/)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [Samber/Do Documentation](https://github.com/samber/do)

## 🙏 Thank You

Thank you for contributing to Templ Router! Your contributions help make this project better for everyone.

If you have any questions about contributing, please:

- Check existing [Discussions](https://github.com/denkhaus/templ-router/discussions)
- Open a new [Discussion](https://github.com/denkhaus/templ-router/discussions/new)
- Ask in your pull request

---

**Happy coding! 🚀**