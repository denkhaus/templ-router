# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure
- File-based routing system
- Internationalization (i18n) support with locale directories
- Authentication and authorization system
- Next.js-style layout inheritance
- Comprehensive configuration service
- Dependency injection container
- Template caching and optimization
- CLI tools for code generation
- Development server with hot reloading

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- CSRF protection implementation
- Rate limiting middleware
- Security headers configuration

## [0.1.0] - 2024-01-XX

### Added
- Initial release of Templ Router
- Core routing functionality
- Basic template system
- Configuration management
- Authentication framework
- Layout system foundation

---

## Release Notes

### Version Numbering
- **Major version** (X.y.z): Breaking changes that require code modifications
- **Minor version** (x.Y.z): New features that are backward compatible
- **Patch version** (x.y.Z): Bug fixes and minor improvements

### Breaking Changes Policy
During the early development phase (0.x.x versions), breaking changes may occur in minor versions. We will clearly document all breaking changes and provide migration guides.

### Support Policy
- **Current version**: Full support with bug fixes and security updates
- **Previous minor version**: Security updates only
- **Older versions**: No support (please upgrade)

### How to Upgrade
1. Check this changelog for breaking changes
2. Update your `go.mod` file
3. Run `go mod tidy`
4. Test your application thoroughly
5. Follow any migration guides provided

### Reporting Issues
If you encounter issues after upgrading, please:
1. Check our [troubleshooting guide](docs/troubleshooting.md)
2. Search existing [GitHub issues](https://github.com/denkhaus/templ-router/issues)
3. Create a new issue with detailed information