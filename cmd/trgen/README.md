# trgen (templ-router generator)

A powerful tool for automatically generating template registries from your templ files with proper versioning and cross-platform support.

## Quick Start

### Installation

```bash
# Install globally (recommended)
mage generator:installGlobal

# Or install to GOPATH/bin
mage generator:install

# Development version with race detection
mage generator:dev
```

### Usage

```bash
# Generate templates for current project
trgen

# Custom paths
trgen --scan-path ./app --output-dir ./generated

# Show version
trgen --version
```

## Available Mage Tasks

### 🔧 **Installation & Building**
- `mage generator:install` - Build and install with proper versioning
- `mage generator:installGlobal` - Install globally and verify PATH setup
- `mage generator:build` - Build locally (outputs to `bin/`)
- `mage generator:dev` - Install development version with race detection

### 🧪 **Testing**
- `mage generator:test` - Run all tests (90+ test cases)
- `mage generator:testCoverage` - Generate coverage report

### 🚀 **Release & Distribution**
- `mage generator:release` - Build for multiple platforms
- `mage generator:version` - Show version information
- `mage generator:clean` - Remove build artifacts

## Features

### ✅ **Automatic Versioning**
- Git tag-based versioning (e.g., `v1.2.3`)
- Automatic commit hash inclusion
- Build timestamp tracking
- Development version detection

### ✅ **Cross-Platform Support**
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)
- Automatic platform detection

### ✅ **Comprehensive Testing**
- 90+ test cases covering all components
- Docker vs Local environment testing
- Real-world scenario validation
- Performance and error handling tests

### ✅ **Development Features**
- Race detection in dev builds
- Hot reload support
- Debug symbol inclusion
- Enhanced error reporting

## Example Output

```bash
$ mage generator:install
🔧 Installing template generator...
📦 Building template-generator v1.2.3 (commit: abc1234)
✅ Verifying installation...
🎉 Successfully installed: template-generator version v1.2.3-abc1234
📍 Binary location: /home/user/go/bin/template-generator

$ trgen --version
trgen version v1.2.3-abc1234
Built: 2024-01-15T10:30:45Z
Commit: abc1234
Go: go1.21.0
Platform: linux/amd64
```

## Integration Examples

### Makefile Integration

```makefile
.PHONY: install-generator generate-templates

install-generator:
	@mage generator:install

generate-templates: install-generator
	@trgen --scan-path app --output-dir generated
	@echo "✅ Template registry generated!"
```

### CI/CD Integration

```yaml
# GitHub Actions
- name: Install Template Generator
  run: mage generator:install

- name: Generate Templates  
  run: trgen --scan-path app --output-dir generated

- name: Test Generator
  run: mage generator:test
```

### Docker Integration

```dockerfile
# Multi-stage build
FROM golang:1.21 AS generator-builder
WORKDIR /src
COPY . .
RUN mage generator:build

FROM alpine:latest
COPY --from=generator-builder /src/bin/template-generator /usr/local/bin/
RUN template-generator --version
```

## Development Workflow

```bash
# 1. Install development version
mage generator:dev

# 2. Run tests
mage generator:test

# 3. Generate coverage report
mage generator:testCoverage

# 4. Build for release
mage generator:release

# 5. Clean up
mage generator:clean
```

## Troubleshooting

### PATH Issues
```bash
# Check PATH setup
mage generator:installGlobal

# Manual PATH setup
export PATH="$(go env GOPATH)/bin:$PATH"
```

### Build Issues
```bash
# Clean and rebuild
mage generator:clean
mage generator:install

# Check dependencies
go mod tidy
```

### Version Issues
```bash
# Check current version
mage generator:version

# Force version update
git tag v1.2.3
mage generator:install
```

## Documentation

- **[Installation Guide](../../docs/template-generator-installation.md)** - Comprehensive installation and usage guide
- **[Testing Guide](./utils/README.md)** - Testing framework and coverage details
- **[API Documentation](./types/README.md)** - Type definitions and interfaces

## Support

- **Issues**: [GitHub Issues](https://github.com/denkhaus/templ-router/issues)
- **Tests**: `mage generator:test` (90+ test cases)
- **Coverage**: `mage generator:testCoverage`
- **Version**: `mage generator:version`

## License

MIT License - see [LICENSE](../../LICENSE) for details.