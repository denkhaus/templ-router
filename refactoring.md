Das System ist unnötig komplex mit unterschiedlichen Behandlungen für Pages, Layouts und Components. Lass uns das komplett verschlanken und vereinheitlichen.

──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
Ready to code?

Here is Claude's plan:
╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
Unified Component System - Eliminate Page/Component Distinction

Core Principle

Everything is a Component - Pages, Layouts, and Components are all templates with the same metadata and i18n loading mechanism.

Root Issues to Fix

1.  Unnecessary Complexity: Different code paths for pages vs components
2.  Redundant Fallbacks: Multiple config types and loading mechanisms
3.  Inconsistent Context: I18n and metadata loaded differently for different template types
4.  Legacy Code: Old functions that are no longer needed

Simplified Architecture

Phase 1: Unified Template Loading

- Single Loading Path: All templates (pages/layouts/components) use the same loading mechanism
- Template Registry Only: Remove redundant config loading, rely solely on registry metadata
- Component Metadata Service: Use this for ALL templates, not just "components"
- Remove Page-Specific Code: Eliminate special handling for "pages" in template middleware

Phase 2: Streamlined I18n System

- Context-First Design: Initialize i18n context ONCE for all requests
- Template-Based Loading: Load translations based on template path from registry, not route type
- Remove Legacy Functions: Delete old i18n loading methods that are no longer needed
- Direct YAML Access: Component metadata service loads directly from YAML files

Phase 3: Unified Metadata System

- Single Context Key: Use only shared.TemplateConfigKey with shared.ConfigFile
- Remove Router Config Types: Eliminate shared.ConfigFile redundancy
- Direct Metadata Loading: All templates get metadata the same way via component metadata service
- No Special Cases: Layouts, pages, components all identical loading

Phase 4: Clean Template Middleware

- Simplify Flow: One path through template middleware for all templates
- Remove Conditional Logic: No more "if component vs if page" branching
- Registry-Driven: All template info comes from registry metadata
- Consistent Context: Same context structure for all template types

Phase 5: Remove Legacy Code

- Delete Old Functions: Remove all obsolete fallback methods
- Clean Interfaces: Simplify interfaces to remove redundant methods
- Consolidate Services: Merge or remove services that are no longer needed
- Update Tests: Remove tests for old legacy behavior

Implementation Details

Registry-First Loading

// All templates use this single path:
templateKey := registry.GetTemplateKey(route.Path)
metadata := registry.GetTemplateMetadata(templateKey)
config := componentMetadataService.LoadMetadata(metadata.TemplatePath)
i18nContext := i18nService.CreateContext(ctx, metadata.TemplatePath)

Unified Context Structure

// Single context for all templates:
ctx = WithTemplateConfig(ctx, config) // metadata.M() access
ctx = WithI18nData(ctx, translations) // i18n.T() access
ctx = WithLocale(ctx, locale) // i18n.GetCurrentLocale()

Template Rendering Flow

1.  Route matched → Get template key from registry
2.  Load template metadata from registry (no file system checks)
3.  Load YAML config via component metadata service (same for all types)
4.  Create i18n context via i18n service (same for all types)
5.  Render template (same for all types)

Benefits

- Dramatic Simplification: 70% reduction in template-related code
- No Special Cases: Consistent behavior across all template types
- Better Performance: No redundant loading or fallback chains
- Easier Maintenance: Single code path to understand and debug
- Clean Architecture: True separation of concerns with single responsibility

Files to Clean Up

- template_middleware.go: Remove all conditional component/page logic
- i18n_service.go: Remove legacy fallback methods
- component_metadata_service.go: Make this the primary metadata loader for ALL templates
- Various interfaces: Remove redundant methods and config types

This creates a truly unified system where every template is treated identically, eliminating the artificial distinction between pages and components.

Double implemented:

1. shared.ConfigFile (Zeile 173 in types.go) - mit typisierten Feldern und Validierung
2. shared.ConfigFile (Zeile 11 in yaml_parser.go) - einfacher struct mit interface{} Feldern

beim Page rendering alle komponenten sammeln und metadata + i18n infos hierarchich mergen: layout->page->component wobei jeweils kindelemente parent elemente überschreiben können, da sie höhere priorität haben.
