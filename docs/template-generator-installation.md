# trgen Installation Guide

trgen (templ-router generator) is a powerful tool for automatically generating template registries from your templ files. This guide covers installation, usage, and development workflows.

## Quick Installation

### Using Mage (Recommended)

```bash
# Install the latest version globally
mage generator:installGlobal

# Or just install to GOPATH/bin
mage generator:install

# For development with race detection
mage generator:dev
```

### Manual Installation

```bash
# Install directly with Go
go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Or build from source
git clone https://github.com/denkhaus/templ-router.git
cd templ-router
mage generator:build
```

## Available Mage Tasks

### Development Tasks

- **`mage generator:build`** - Build the generator locally (outputs to `bin/`)
- **`mage generator:install`** - Build and install to `$GOPATH/bin`
- **`mage generator:installGlobal`** - Install globally and check PATH setup
- **`mage generator:dev`** - Install development version with race detection

### Testing Tasks

- **`mage generator:test`** - Run all template generator tests
- **`mage generator:testCoverage`** - Run tests with coverage report
- **`mage generator:version`** - Show current version information

### Release Tasks

- **`mage generator:release`** - Build for multiple platforms (Linux, macOS, Windows)
- **`mage generator:clean`** - Remove build artifacts

## Version Information

The generator includes comprehensive version information:

```bash
$ trgen --version
trgen version v1.2.3-abc1234

$ mage generator:version
Template Generator Version: v1.2.3
Git Commit: abc1234
Go Version: go1.21.0
Platform: linux/amd64
```

## Build Features

### Automatic Versioning

The build system automatically determines version information:

1. **Git Tags**: Uses the latest git tag (e.g., `v1.2.3`)
2. **Git Describe**: Falls back to `git describe --tags --always --dirty`
3. **Commit Count**: Uses commit count if no tags available
4. **Fallback**: Uses "dev" if git is not available

### Build Information

Each binary includes:
- **Version**: Semantic version or git description
- **Git Commit**: Short commit hash
- **Build Time**: UTC timestamp of build
- **Go Version**: Version of Go used to build
- **Platform**: Target OS and architecture

### Cross-Platform Builds

The release task builds for multiple platforms:

```bash
mage generator:release
```

Creates binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64, arm64)

Output: `release/template-generator-v{version}/`

## Development Workflow

### Local Development

```bash
# Build and test locally
mage generator:build
mage generator:test

# Install development version
mage generator:dev

# Run with coverage
mage generator:testCoverage
```

### Testing

The generator has comprehensive test coverage:

```bash
# Run all tests
mage generator:test

# Generate coverage report
mage generator:testCoverage
# Opens coverage/generator.html in browser
```

### Release Process

```bash
# 1. Tag the release
git tag v1.2.3
git push origin v1.2.3

# 2. Build release binaries
mage generator:release

# 3. Test the release
./release/template-generator-v1.2.3/template-generator-linux-amd64 --version
```

## Usage Examples

### Basic Usage

```bash
# Generate templates for current directory
template-generator

# Specify custom paths
template-generator --scan-path ./templates --output-dir ./generated

# Show help
template-generator --help
```

### Integration with Projects

Add to your project's Makefile:

```makefile
.PHONY: generate-templates
generate-templates:
	@echo "Generating template registry..."
	@template-generator --scan-path app --output-dir generated/templates
	@echo "Template registry generated successfully!"

.PHONY: install-generator
install-generator:
	@echo "Installing template generator..."
	@mage generator:install
```

Or use in CI/CD:

```yaml
# .github/workflows/build.yml
- name: Install Template Generator
  run: mage generator:install

- name: Generate Templates
  run: template-generator --scan-path app --output-dir generated
```

## Troubleshooting

### PATH Issues

If `template-generator` command is not found after installation:

```bash
# Check if GOPATH/bin is in PATH
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✅ GOPATH/bin is in PATH" || echo "❌ GOPATH/bin not in PATH"

# Add to your shell profile (.bashrc, .zshrc, etc.)
export PATH="$(go env GOPATH)/bin:$PATH"

# Or use the global install task
mage generator:installGlobal
```

### Build Issues

```bash
# Clean and rebuild
mage generator:clean
mage generator:build

# Check Go version (requires Go 1.21+)
go version

# Verify dependencies
go mod tidy
```

### Version Issues

```bash
# Check current version
mage generator:version

# Force rebuild with latest version
mage generator:clean
mage generator:install
```

## Advanced Configuration

### Custom Build Flags

For custom builds, you can modify the ldflags in `magefiles/generator.go`:

```go
ldflags := fmt.Sprintf("-X github.com/denkhaus/templ-router/cmd/template-generator/version.Version=%s "+
    "-X github.com/denkhaus/templ-router/cmd/template-generator/version.GitCommit=%s "+
    "-X github.com/denkhaus/templ-router/cmd/template-generator/version.BuildTime=%s",
    version, gitCommit, buildTime)
```

### Development Mode Features

Development builds include:
- Race detection enabled
- Debug symbols included
- "-dev" version suffix
- Enhanced error reporting

```bash
# Install development version
mage generator:dev

# Verify development features
template-generator --version
# Output: template-generator version v1.2.3-dev-abc1234
```

## Support

For issues, feature requests, or contributions:

1. Check existing issues: [GitHub Issues](https://github.com/denkhaus/templ-router/issues)
2. Run tests: `mage generator:test`
3. Check version: `mage generator:version`
4. Provide build information when reporting issues

## License

MIT License - see LICENSE file for details.