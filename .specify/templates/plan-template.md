# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.24+
**Primary Dependencies**: templ, samber/do/v2, chi/v5, Ginkgo/Gomega
**Storage**: [if applicable, e.g., PostgreSQL, files, or N/A]
**Testing**: Go testing + Ginkgo/Gomega for E2E, testify for unit tests
**Target Platform**: Linux server, Docker containers
**Project Type**: Go library (determines source structure)
**Performance Goals**: Template caching, routing performance, or NEEDS CLARIFICATION
**Constraints**: <200ms p95 response time, <100MB memory, ASCII-only source code
**Scale/Scope**: Designed for developer integration, not end-user scale

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Required Compliance Gates
- [ ] **Library-First**: Feature designed as importable library, not standalone application
- [ ] **Type-Safe Templates**: Uses templ engine, Next.js-style component hierarchy, self-contained .templ.yaml metadata
- [ ] **Test-First**: TDD approach with Red-Green-Refactor cycle planned
- [ ] **Clean Architecture**: Clear layer separation with unified DI through di package
- [ ] **Convention Over Configuration**: File-based routing and environment variable configuration
- [ ] **Technology Standards**: Go 1.24+, Mage build automation, golangci-lint compliance
- [ ] **Performance & Security**: Template caching, middleware configuration, ASCII-only source code

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
