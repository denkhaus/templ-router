<!--
Sync Impact Report:
- Version change: 1.0.0 → 1.1.0 (added mandatory template naming conventions)
- Modified principles: Added Principle VI (Template Naming Conventions) as mandatory requirement
- Added sections: Template Naming Conventions with detailed examples and consequences
- Removed sections: None
- Templates updated:
  ✅ .specify/templates/plan-template.md - Added constitution compliance gates, updated Go-specific tech stack
  ✅ .specify/templates/tasks-template.md - Updated path conventions for Go library structure
  ⚠️ .specify/templates/spec-template.md - Already aligned with TDD principles (no changes needed)
- Follow-up TODOs: None
- Key architectural clarifications added:
  - Next.js-style routing conventions (Components → Pages → Layouts)
  - Self-contained component metadata with .templ.yaml files
  - Unified dependency injection through di package
  - Component reusability with independent i18n and metadata
  - **NEW**: Mandatory template naming conventions (Page(), Layout(), Error())
  - **NEW**: Dynamic route parameter syntax (id_ not [id])
  - **NEW**: Layout template signature requirements (content templ.Component only)
-->

# Templ Router Constitution

## Core Principles

### I. Library-First Architecture
Templ Router is a library, not a standalone application. Every feature must be designed as a standalone, importable library that developers can integrate into their own Go projects. Libraries must be self-contained, independently testable, well-documented, and serve a clear technical purpose. No organizational-only libraries or framework lock-in.

### II. Type-Safe Templates with Self-Contained Components
All templates use the templ engine for compile-time type safety and Go integration. Templates must be type-safe, compile to Go code, and leverage Go's type system. The system follows Next.js-style conventions: Components can be nested in Pages, Pages are loaded by Layouts, and each component can have its own .templ.yaml metadata file with internationalization and configuration settings that can be used within the component template. No runtime template engines or string-based templates.

### III. Test-First Development (NON-NEGOTIABLE)
TDD is mandatory: Tests must be written before implementation, user-approved, and must fail initially. The Red-Green-Refactor cycle is strictly enforced. All code must have comprehensive unit tests, and critical paths require integration tests. No exceptions.

### IV. Clean Architecture with Unified Dependency Injection
All code follows clean architecture principles with clear separation between infrastructure, application, and domain layers. All services are interconnected through the dependency injection package (di) using samber/do/v2 with named dependencies. The DI container provides unified service management for both router services and application services. No global state, no service locator patterns, and no direct dependencies between layers.

### V. Convention Over Configuration
Routing and configuration follow convention over configuration principles. File-based routing automatically generates routes from directory structure. Configuration uses environment variables with consistent prefixes. No complex configuration files or manual route definitions unless absolutely necessary.

### VI. Template Naming Conventions (MANDATORY)
Templ Router enforces strict naming conventions for template functions. These are **non-negotiable requirements** for all template files:

**Page Templates:**
- MUST be named `Page()`
- Can accept parameters for dynamic routes: `templ Page(id string)`
- Examples: `templ Page()`, `templ Page(userID string)`, `templ Page(locale, id string)`

**Layout Templates:**
- MUST be named `Layout(content templ.Component)`
- MUST accept exactly one parameter: `content templ.Component`
- Title and metadata are handled through the metadata system, not parameters
- Example: `templ Layout(content templ.Component)`

**Error Templates:**
- MUST be named `Error(errCtx middleware.ErrorContext)`
- Used for error page rendering
- Example: `templ Error(errCtx middleware.ErrorContext)`

**Component Templates:**
- Can have custom names (e.g., `Header()`, `Footer()`, `Navbar()`)
- Are reusable and don't generate routes directly
- Can be accessed via `/components/*` routes with their metadata

**Dynamic Route Parameters:**
- Use underscore suffix: `id_`, `locale_`, `slug_`
- NOT bracket syntax: `[id]` or `{id}` - this is incorrect
- Examples: `app/user/id_/page.templ` → `/user/{id}`, `app/locale_/page.templ` → `/{locale}`

**File Structure Examples:**
```
app/
├── layout.templ           # templ Layout(content templ.Component)
├── page.templ             # templ Page()
├── login/
│   └── page.templ         # templ Page()
├── user/
│   └── id_/
│       └── page.templ     # templ Page(userID string)
└── components/
    ├── header.templ       # templ Header()
    └── footer.templ       # templ Footer()
```

**Consequences of Violations:**
- Templates with incorrect names will not be found by the router
- Dynamic routes with incorrect syntax will not work
- Layout templates with wrong signatures will cause build failures
- Code reviews must enforce these conventions without exception

## Technology Standards

### Language and Tooling
- **Go 1.24+** with modern Go features required
- **templ** engine for type-safe templates
- **samber/do/v2** for dependency injection
- **chi/v5** for HTTP routing
- **Ginkgo/Gomega** for E2E testing
- **Mage** for build automation

### Performance and Security
- Template caching is mandatory for production
- All middleware must be configurable via environment variables
- Security headers, CSRF protection, and rate limiting enabled by default
- Session-based authentication with secure defaults
- No Unicode/emojis in Go source code (ASCII only)

## Development Workflow

### Code Quality Standards
- All code must pass golangci-lint with v2.6.1 configuration
- Generated files must never be manually edited
- Component metadata must be self-contained in .templ.yaml files with i18n and configuration
- Components must be truly reusable with independent metadata and internationalization
- All public APIs must have comprehensive documentation

### Testing Requirements
- Unit tests for all service interfaces and implementations
- E2E tests for routing, authentication, i18n, and data services
- Performance tests for template caching and routing
- Integration tests for middleware pipeline and DI container

### Build and Release Process
- Use Mage for all build automation tasks
- Generated files (templates, registry) are built, not committed
- Semantic versioning must be followed for all releases
- All releases must pass full test suite and security scanning

## Governance

This constitution supersedes all other practices and guidelines. Amendments require documentation, approval, and a migration plan. All PRs and code reviews must verify compliance with these principles. Complexity must be justified with clear technical rationale.

For runtime development guidance, refer to CLAUDE.md and the docs/ directory. Use the speckit command system for feature development and planning.

**Version**: 1.1.0 | **Ratified**: 2025-11-07 | **Last Amended**: 2025-11-07
