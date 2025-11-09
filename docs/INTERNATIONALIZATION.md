# Internationalization (i18n)

**Multi-language support with locale-based routing and context-based translations.**

## Overview

Templ Router provides comprehensive internationalization with `locale_/` directory structure, YAML-based translations, and context-based translation system.

**Key Features:**
- Multi-language support with `locale_/` directory structure
- YAML-based translations in `.templ.yaml` metadata files
- Context-based translation system (no global `t()` function)
- Automatic locale detection and validation
- URL localization with locale prefixing
- Smart fallback mechanisms

## Configuration

### Environment Variables

Configure i18n behavior (prefix configurable, default `TR`):

```bash
# Language settings
TR_I18N_SUPPORTED_LOCALES=en,de,fr    # Comma-separated list
TR_I18N_DEFAULT_LOCALE=en               # Default fallback
TR_I18N_FALLBACK_LOCALE=en             # Ultimate fallback

# Locale detection
TR_I18N_ENABLE_MIDDLEWARE=true           # Enable i18n middleware
TR_I18N_STRICT_MODE=false              # Allow unknown locales
TR_I18N_REDIRECT_TO_DEFAULT=true        # Redirect unknown locales
```

### Custom Prefix

```go
container.RegisterRouterServices("MYAPP")  // Uses MYAPP_* variables
```

## File Structure

### Locale-Based Routing

Use `locale_/` directory for internationalized routes:

```
app/
├── page.templ                → / (non-localized)
├── about.templ               → /about (non-localized)
└── locale_/                  # Internationalized routes
    ├── page.templ            → /{locale}
    ├── about.templ           → /{locale}/about
    └── user/
        └── profile.templ     → /{locale}/user/profile
```

### Translation Files

Translation keys are defined in `.templ.yaml` files:

```yaml
# app/locale_/about/page.templ.yaml
i18n:
  en:
    page_title: "About Us"
    company_name: "Tech Company"
    description: "Leading technology solutions"
  de:
    page_title: "Über Uns"
    company_name: "Technikunternehmen"
    description: "Führende Technologie-Lösungen"
  fr:
    page_title: "À Propos"
    company_name: "Société Tech"
    description: "Solutions technologiques de pointe"

metadata:
  title: "About"

auth:
  type: "Public"
```

## Template Usage

### Translation Functions

```go
templ Page() {
    <h1>{ i18n.T(ctx, "page_title") }</h1>
    <h2>{ i18n.T(ctx, "company_name") }</h2>
    <p>{ i18n.T(ctx, "description") }</p>

    <a href={ i18n.LocalizeSafeURL(ctx, "/contact") }>
        { i18n.T(ctx, "contact_us") }
    </a>

    <p>Current locale: { i18n.GetCurrentLocale(ctx) }</p>
}
```

### Available Functions

```go
// Get translation for current locale
i18n.T(ctx, "translation_key")

// Localize URL for current locale
i18n.LocalizeSafeURL(ctx, "/path")

// Get current locale
i18n.GetCurrentLocale(ctx)

// Check if translation exists
i18n.HasTranslation(ctx, "translation_key")
```

### Nested Keys

Use dot notation for nested translations:

```yaml
i18n:
  en:
    nav:
      home: "Home"
      about: "About"
      contact: "Contact"
    buttons:
      submit: "Submit"
      cancel: "Cancel"
```

```go
// Usage
i18n.T(ctx, "nav.home")
i18n.T(ctx, "buttons.submit")
```

## Locale Detection

### Detection Order

1. **URL Path** - `/en/page`, `/de/page`
2. **Query Parameter** - `?locale=en`
3. **Accept-Language Header** - Browser preference
4. **Cookie** - `locale=en`
5. **Default Locale** - Configured fallback

### URL Structure

```bash
# Direct locale URLs
https://example.com/en/about
https://example.com/de/über

# Automatic redirect based on browser language
https://example.com/about → redirects to /en/about
```

## Translation Formats

### Multi-Locale Files

```yaml
# app/locale_/page.templ.yaml
i18n:
  en:
    welcome: "Welcome"
  de:
    welcome: "Willkommen"
  fr:
    welcome: "Bienvenue"
```

### Language-Specific Files

```yaml
# app/locale_/en/page.templ.yaml
i18n:
  welcome: "Welcome"
  good_morning: "Good morning"

# app/locale_/de/page.templ.yaml
i18n:
  welcome: "Willkommen"
  good_morning: "Guten Morgen"
```

### Nested Structure

```yaml
# app/locale_/dashboard/page.templ.yaml
i18n:
  en:
    dashboard:
      title: "Dashboard"
      stats:
        users: "Total Users"
        revenue: "Revenue"
    menu:
      overview: "Overview"
      settings: "Settings"
  de:
    dashboard:
      title: "Dashboard"
      stats:
        users: "Gesamte Benutzer"
        revenue: "Umsatz"
    menu:
      overview: "Übersicht"
      settings: "Einstellungen"
```

## Component i18n

### Self-Contained Components

Components can have their own i18n:

```yaml
# app/components/footer.templ.yaml
i18n:
  en:
    copyright: "© 2024 My Company"
    privacy: "Privacy Policy"
    terms: "Terms of Service"
  de:
    copyright: "© 2024 Meine Firma"
    privacy: "Datenschutz"
    terms: "Nutzungsbedingungen"

metadata:
  company_name: "My Company"
  version: "1.0.0"
```

```go
// app/components/footer.templ
package components

templ Footer() {
    <footer class="bg-gray-800 text-white p-4">
        <p>{ i18n.T(ctx, "copyright") }</p>
        <div class="flex space-x-4 text-sm">
            <a href="/privacy">{ i18n.T(ctx, "privacy") }</a>
            <a href="/terms">{ i18n.T(ctx, "terms") }</a>
        </div>
        <p class="text-xs mt-2">{ metadata.M(ctx, "company_name") } v{ metadata.M(ctx, "version") }</p>
    </footer>
}
```

## Locale-Aware Metadata

Metadata can also be localized just like translations. This allows you to have different metadata values for different locales.

### Locale-Specific Metadata Structure

```yaml
# app/components/footer.templ.yaml
i18n:
  en:
    copyright: "© 2024 My Company"
    privacy: "Privacy Policy"
    terms: "Terms of Service"
  de:
    copyright: "© 2024 Meine Firma"
    privacy: "Datenschutz"
    terms: "Nutzungsbedingungen"

metadata:
  en:
    company_name: "My Company"
    company_email: "info@mycompany.com"
    region: "North America"
  de:
    company_name: "Meine Firma"
    company_email: "info@meinefirma.de"
    region: "Europa"
  # Global fallback values
  version: "1.0.0"
  author: "Development Team"
```

### Metadata Resolution Priority

1. **Locale-specific metadata** (highest priority)
2. **Global metadata** (fallback for all locales)
3. **Default fallback** (if no locale match)

### Usage in Templates

```go
// app/components/footer.templ
package components

templ Footer() {
    <footer class="bg-gray-800 text-white p-4">
        <p>{ i18n.T(ctx, "copyright") }</p>
        <div class="flex space-x-4 text-sm">
            <a href="/privacy">{ i18n.T(ctx, "privacy") }</a>
            <a href="/terms">{ i18n.T(ctx, "terms") }</a>
        </div>
        <div class="text-xs mt-2">
            <p>{ metadata.M(ctx, "company_name") }</p>
            <p>{ metadata.M(ctx, "company_email") }</p>
            <p>{ metadata.M(ctx, "region") }</p>
            <p>Version: { metadata.M(ctx, "version") }</p>
            <p>Author: { metadata.M(ctx, "author") }</p>
        </div>
    </footer>
}
```

### How It Works

- **English locale** (`/en/page`): Uses `metadata.en.company_name`, `metadata.en.company_email`, etc.
- **German locale** (`/de/page`): Uses `metadata.de.company_name`, `metadata.de.company_email`, etc.
- **Fallback values**: `metadata.version` and `metadata.author` are available for all locales
- **Unknown locale**: Falls back to global metadata values

### Backward Compatibility

Existing flat metadata structures continue to work:

```yaml
# Existing structure (still supported)
metadata:
  company_name: "My Company"
  version: "1.0.0"
  author: "Development Team"
```

This ensures existing templates continue to work without changes while enabling new locale-aware functionality.

### Component Route Access

Components are accessible via their own routes with i18n:

- `/components/footer` - Footer with English translations
- `/de/components/footer` - Footer with German translations

## Advanced Features

### Locale Validation

Configure strict locale validation:

```yaml
# page.templ.yaml
i18n:
  en:
    title: "Home"
  de:
    title: "Startseite"

validation:
  required_locales: ["en", "de"]
  fallback_locale: "en"
```

### Dynamic Locale Switching

```go
// Language switcher component
templ LanguageSwitcher() {
    <div class="language-switcher">
        { range i18n.GetSupportedLocales(ctx) }
            <a href={ i18n.LocalizeSafeURL(ctx, "/") }
               class={ "locale-link " + i18n.LocaleClass(ctx, .) }>
                { i18n.LocaleDisplayName(ctx, .) }
            </a>
        { endfor }
    </div>
}
```

### Translation Pluralization

```yaml
i18n:
  en:
    items_count:
      one: "1 item"
      other: "{count} items"
    messages:
      zero: "No messages"
      one: "1 message"
      other: "{count} messages"
```

```go
// Pluralization usage
{ i18n.TPlural(ctx, "items_count", itemCount) }
{ i18n.TPlural(ctx, "messages", messageCount) }
```

### Date/Time Formatting

```go
// Date localization
{ i18n.FormatDate(ctx, time.Now(), "2006-01-02") }
{ i18n.FormatTime(ctx, time.Now(), "15:04:05") }
{ i18n.FormatDateTime(ctx, time.Now(), "2006-01-02 15:04:05") }
```

## Best Practices

### File Organization

1. **Consistent Structure** - Use `locale_/` for all internationalized content
2. **Fallback Planning** - Always provide default locale translations
3. **Key Naming** - Use descriptive, consistent translation keys
4. **Component i18n** - Make truly self-contained components

### Translation Management

1. **Centralized Keys** - Use consistent key naming across files
2. **Context Awareness** - Provide context-specific translations
3. **Regular Updates** - Keep translations synchronized with content changes
4. **Quality Assurance** - Review translations for accuracy

### Performance

1. **Lazy Loading** - Load translations only when needed
2. **Caching** - Translation results are cached automatically
3. **Minimal Files** - Keep translation files focused and organized

## Troubleshooting

### Common Issues

**Translation not found:**
- Check key spelling and dot notation
- Verify locale is supported in configuration
- Ensure `.templ.yaml` file exists and is properly formatted

**Locale not detected:**
- Check URL path format: `/{locale}/page`
- Verify supported locales configuration
- Check middleware is enabled

**Component i18n not working:**
- Ensure component has `.templ.yaml` file
- Check component route accessibility
- Verify metadata structure

### Debug Mode

```bash
# Enable verbose i18n logging
TR_I18N_DEBUG=true

# Check supported locales
curl -H "Accept-Language: de,en" https://example.com
```

### Translation Testing

```go
func TestI18n(t *testing.T) {
    ctx := createTestContext("en")

    // Test translation retrieval
    title := i18n.T(ctx, "page_title")
    assert.Equal(t, "Home", title)

    // Test URL localization
    url := i18n.LocalizeSafeURL(ctx, "/about")
    assert.Equal(t, "/en/about", url)
}
```

---

**Related Documentation**: [File-Based Routing](FILE-BASED-ROUTING.md), [Template System](TEMPLATES.md), [Configuration](CONFIGURATION.md)