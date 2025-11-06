# Internationalization (i18n)

**Complete guide to implementing multi-language support in Templ Router applications.**

## Overview

Templ Router provides a comprehensive internationalization system that supports multiple languages, locale-based routing, context-based translations, and flexible translation file formats. The i18n system integrates seamlessly with the file-based routing system through `.templ.yaml` metadata files.

**Key Features:**
- Multi-language support with `locale_/` directory structure
- YAML-based translations in `.templ.yaml` metadata files
- Context-based translation system (no global `t()` function)
- Automatic locale detection and validation from URLs
- Flexible i18n formats: flat, nested, and multi-locale configurations
- Dot notation support for deeply nested translation keys
- URL localization with automatic locale prefixing
- Smart locale detection with fallback mechanisms

## Configuration Prefix Notice

**Important:** All environment variables in this documentation use the default prefix `TR_`. This prefix is **configurable** when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_I18N_SUPPORTED_LOCALES=en,de,fr`
- Custom: `MYAPP_I18N_SUPPORTED_LOCALES=en,de,fr`
- Multiple apps: `APP1_I18N_SUPPORTED_LOCALES=en,de` and `APP2_I18N_SUPPORTED_LOCALES=fr,es`

All examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix in all environment variable names.

## Configuration

### Environment Variables

Configure internationalization behavior through environment variables:

```bash
# Supported Languages
TR_I18N_SUPPORTED_LOCALES=en,de,fr,it,es    # Comma-separated list
TR_I18N_DEFAULT_LOCALE=en                   # Default fallback language
TR_I18N_FALLBACK_LOCALE=en                  # Ultimate fallback

# Locale Detection
TR_I18N_DETECTION_METHOD=url                # url, header, cookie
TR_I18N_COOKIE_NAME=locale                  # Cookie name for locale storage
TR_I18N_COOKIE_EXPIRY=8760h                 # Cookie expiry (1 year)

# Routing Behavior
TR_I18N_URL_PREFIX=true                     # Add locale prefix to URLs
TR_I18N_REDIRECT_ROOT=true                  # Redirect / to default locale
TR_I18N_STRICT_ROUTING=false                # 404 for unsupported locales
```

### Supported Locales

Configure the languages your application supports:

```bash
# Basic English and German
TR_I18N_SUPPORTED_LOCALES=en,de
TR_I18N_DEFAULT_LOCALE=en

# Multiple European languages
TR_I18N_SUPPORTED_LOCALES=en,de,fr,it,es,nl
TR_I18N_DEFAULT_LOCALE=en

# Languages with regions
TR_I18N_SUPPORTED_LOCALES=en-US,en-GB,de-DE,fr-FR
TR_I18N_DEFAULT_LOCALE=en-US
```

## File Structure

### Internationalized Directory Structure

Create internationalized routes using the `locale_/` directory:

```
app/
├── page.templ                    # / (root page)
├── about.templ                   # /about (non-localized)
├── locale_/                      # Internationalized routes
│   ├── page.templ               # /en, /de, /fr (based on config)
│   ├── about/
│   │   └── page.templ           # /en/about, /de/about, /fr/about
│   ├── dashboard/
│   │   ├── page.templ           # /en/dashboard, /de/dashboard
│   │   └── settings/
│   │       └── page.templ       # /en/dashboard/settings
│   └── user/
│       └── id_/
│           └── page.templ       # /en/user/123, /de/user/123
└── components/
    ├── navbar.templ             # /components/navbar (shared)
    └── footer.templ             # /components/footer (shared)
```

### Route Generation Examples

Given supported locales `en,de,fr`:

```
app/locale_/page.templ                    → /en, /de, /fr
app/locale_/about/page.templ              → /en/about, /de/about, /fr/about
app/locale_/user/id_/page.templ           → /en/user/123, /de/user/123, /fr/user/123
app/locale_/admin/settings/page.templ     → /en/admin/settings, /de/admin/settings, /fr/admin/settings
```

## Translation File Formats

Templ Router supports multiple translation file formats to accommodate different project needs:

### 1. Multi-Locale Nested Structure (Recommended)

Perfect for complex applications with organized translation hierarchies:

```yaml
# app/locale_/dashboard/page.templ.yaml
i18n:
  en:
    feedback:
      title: "Feedback Dashboard"
      subtitle: "Overview of customer feedback and analytics"
      export: "Export Data"
      refresh: "Refresh Data"
      reviews: "reviews"
      stats:
        total_reviews: "Total Reviews"
        average_rating: "Average Rating"
        productions: "Productions"
        cache_hit_rate: "Cache Hit Rate"
      productions:
        title: "Productions"
        subtitle: "Overview of all productions with review statistics"
      recent:
        title: "Recent Reviews"
        subtitle: "Latest customer feedback and comments"
      actions:
        create_new: "Create New"
        bulk_export: "Bulk Export"
        settings: "Settings"
  de:
    feedback:
      title: "Feedback Dashboard"
      subtitle: "Übersicht über Kundenfeedback und Analysen"
      export: "Daten exportieren"
      refresh: "Daten aktualisieren"
      reviews: "Bewertungen"
      stats:
        total_reviews: "Gesamtbewertungen"
        average_rating: "Durchschnittsbewertung"
        productions: "Produktionen"
        cache_hit_rate: "Cache-Trefferrate"
      productions:
        title: "Produktionen"
        subtitle: "Übersicht aller Produktionen mit Bewertungsstatistiken"
      recent:
        title: "Aktuelle Bewertungen"
        subtitle: "Neuestes Kundenfeedback und Kommentare"
      actions:
        create_new: "Neu erstellen"
        bulk_export: "Massenexport"
        settings: "Einstellungen"

metadata:
  page_title: "Dashboard"
  section: "analytics"

auth:
  type: "UserRequired"
  redirect_url: "/login"
```

**Usage in templates with dot notation:**
```go
templ DashboardPage() {
    <h1>{ i18n.T(ctx, "feedback.title") }</h1>
    <p>{ i18n.T(ctx, "feedback.subtitle") }</p>
    <div class="stats">
        <span>{ i18n.T(ctx, "feedback.stats.total_reviews") }</span>
        <span>{ i18n.T(ctx, "feedback.stats.average_rating") }</span>
    </div>
    <button>{ i18n.T(ctx, "feedback.actions.create_new") }</button>
}
```

### 2. Multi-Locale Flat Structure

Simple key-value pairs for straightforward translations:

```yaml
# app/locale_/admin/page.templ.yaml
i18n:
  en:
    admin_warning: "Admin Area - Restricted Access"
    page_title: "System Administration"
    user_management_title: "User Management"
    user_management_desc: "Manage user accounts, roles, and permissions"
    system_settings_title: "System Settings"
    system_settings_desc: "Configure application settings and preferences"
    btn_save: "Save Changes"
    btn_cancel: "Cancel"

  de:
    admin_warning: "Admin-Bereich - Eingeschränkter Zugang"
    page_title: "Systemadministration"
    user_management_title: "Benutzerverwaltung"
    user_management_desc: "Benutzerkonten, Rollen und Berechtigungen verwalten"
    system_settings_title: "Systemeinstellungen"
    system_settings_desc: "Anwendungseinstellungen und Präferenzen konfigurieren"
    btn_save: "Änderungen speichern"
    btn_cancel: "Abbrechen"

auth:
  type: "AdminRequired"
  redirect_url: "/login"
```

### 3. Single-Locale Nested Structure

For applications that don't need multi-language support but want organized translations:

```yaml
# app/components/navigation/page.templ.yaml
i18n:
  navigation:
    main:
      home: "Home"
      about: "About Us"
      services: "Services"
      contact: "Contact"
    user:
      profile: "My Profile"
      settings: "Account Settings"
      logout: "Sign Out"
    admin:
      dashboard: "Admin Dashboard"
      users: "User Management"
      reports: "System Reports"
  buttons:
    primary:
      submit: "Submit"
      save: "Save"
      continue: "Continue"
    secondary:
      cancel: "Cancel"
      back: "Go Back"
      reset: "Reset Form"

metadata:
  component_type: "navigation"
  theme: "default"
```

### 4. Mixed Depth Nested Structure

Combining different nesting levels as needed:

```yaml
# app/locale_/ecommerce/page.templ.yaml
i18n:
  en:
    # Top-level keys
    site_name: "Amazing Store"
    welcome_message: "Welcome to our online store!"

    # Nested product information
    products:
      categories:
        electronics: "Electronics"
        clothing: "Clothing"
        books: "Books"
        home_garden: "Home & Garden"
      actions:
        add_to_cart: "Add to Cart"
        buy_now: "Buy Now"
        view_details: "View Details"
        compare: "Compare Products"

    # Deeply nested checkout process
    checkout:
      steps:
        cart: "Shopping Cart"
        shipping: "Shipping Information"
        payment: "Payment Details"
        confirmation: "Order Confirmation"
      shipping:
        methods:
          standard: "Standard Delivery (5-7 days)"
          express: "Express Delivery (2-3 days)"
          overnight: "Overnight Delivery"
      payment:
        methods:
          credit_card: "Credit Card"
          paypal: "PayPal"
          bank_transfer: "Bank Transfer"

  de:
    site_name: "Fantastischer Shop"
    welcome_message: "Willkommen in unserem Online-Shop!"

    products:
      categories:
        electronics: "Elektronik"
        clothing: "Kleidung"
        books: "Bücher"
        home_garden: "Haus & Garten"
      actions:
        add_to_cart: "In den Warenkorb"
        buy_now: "Jetzt kaufen"
        view_details: "Details anzeigen"
        compare: "Produkte vergleichen"

    checkout:
      steps:
        cart: "Warenkorb"
        shipping: "Versandinformationen"
        payment: "Zahlungsdetails"
        confirmation: "Bestellbestätigung"
      shipping:
        methods:
          standard: "Standardversand (5-7 Tage)"
          express: "Expressversand (2-3 Tage)"
          overnight: "Über-Nacht-Versand"
      payment:
        methods:
          credit_card: "Kreditkarte"
          paypal: "PayPal"
          bank_transfer: "Banküberweisung"
```

## Template Integration

### Basic Translation Usage

```go
// app/locale_/dashboard/page.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ DashboardPage() {
    <div class="dashboard">
        <h1>{ i18n.T(ctx, "page_title") }</h1>
        <p>{ i18n.T(ctx, "page_subtitle") }</p>

        <nav class="dashboard-nav">
            <a href={ i18n.LocalizeSafeURL(ctx, "/dashboard") }>
                { i18n.T(ctx, "nav_dashboard") }
            </a>
            <a href={ i18n.LocalizeSafeURL(ctx, "/profile") }>
                { i18n.T(ctx, "nav_profile") }
            </a>
            <a href={ i18n.LocalizeSafeURL(ctx, "/settings") }>
                { i18n.T(ctx, "nav_settings") }
            </a>
        </nav>

        <div class="stats">
            <div class="stat-item">
                <span class="stat-label">{ i18n.T(ctx, "stats_users") }</span>
                <span class="stat-value">1,234</span>
            </div>
            <div class="stat-item">
                <span class="stat-label">{ i18n.T(ctx, "stats_projects") }</span>
                <span class="stat-value">56</span>
            </div>
            <div class="stat-item">
                <span class="stat-label">{ i18n.T(ctx, "stats_revenue") }</span>
                <span class="stat-value">$12,345</span>
            </div>
        </div>

        <div class="actions">
            <button class="btn btn-primary">
                { i18n.T(ctx, "btn_create_new") }
            </button>
            <button class="btn btn-secondary">
                { i18n.T(ctx, "btn_export") }
            </button>
        </div>
    </div>
}
```

### Language Switcher Component

```go
// app/components/language-switcher.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ LanguageSwitcher() {
    currentLocale := i18n.GetCurrentLocale(ctx)
    currentRoute := i18n.GetCurrentRouteWithoutLocale(ctx)

    <div class="language-switcher">
        <div class="flex items-center space-x-2 bg-blue-700 rounded px-3 py-1">
            <span class="text-sm">🌐</span>
            <div class="flex space-x-1">
                if currentLocale == "en" {
                    <span class="bg-white text-blue-600 px-2 py-1 rounded text-sm font-semibold">EN</span>
                    <a href={ "/de" + currentRoute } class="text-blue-200 hover:text-white px-2 py-1 rounded text-sm">DE</a>
                    <a href={ "/fr" + currentRoute } class="text-blue-200 hover:text-white px-2 py-1 rounded text-sm">FR</a>
                } else if currentLocale == "de" {
                    <a href={ "/en" + currentRoute } class="text-blue-200 hover:text-white px-2 py-1 rounded text-sm">EN</a>
                    <span class="bg-white text-blue-600 px-2 py-1 rounded text-sm font-semibold">DE</span>
                    <a href={ "/fr" + currentRoute } class="text-blue-200 hover:text-white px-2 py-1 rounded text-sm">FR</a>
                } else {
                    <a href={ "/en" + currentRoute } class="text-blue-200 hover:text-white px-2 py-1 rounded text-sm">EN</a>
                    <a href={ "/de" + currentRoute } class="text-blue-200 hover:text-white px-2 py-1 rounded text-sm">DE</a>
                    <span class="bg-white text-blue-600 px-2 py-1 rounded text-sm font-semibold">FR</span>
                }
            </div>
        </div>
    </div>
}
```

### Breadcrumb Navigation

```go
// app/components/breadcrumb.templ
package main

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ Breadcrumb() {
    currentLocale := i18n.GetCurrentLocale(ctx)
    currentRoute := i18n.GetCurrentRouteWithoutLocale(ctx)

    <nav class="breadcrumb">
        <a href={ "/" + currentLocale } class="breadcrumb-link">
            { i18n.T(ctx, "nav_home") }
        </a>

        if currentRoute != "/" {
            <span class="breadcrumb-separator">/</span>
            <span class="breadcrumb-current">
                { strings.TrimPrefix(currentRoute, "/") }
            </span>
        }
    </nav>
}
```

### Conditional Content Based on Locale

```go
// app/components/currency-display.templ
package main

templ PriceDisplay(amount float64) {
    currentLocale := i18n.GetCurrentLocale(ctx)

    <div class="price-display">
        <span class="amount">
            if currentLocale == "de" {
                { fmt.Sprintf("%.2f €", amount) }
            } else if currentLocale == "fr" {
                { fmt.Sprintf("%.2f €", amount) }
            } else {
                { fmt.Sprintf("$%.2f", amount) }
            }
        </span>
        <span class="currency-info">
            { i18n.T(ctx, "price_includes_vat") }
        </span>
    </div>
}
```

## I18n Helper Functions

### Core Translation Function

```go
import "github.com/denkhaus/templ-router/pkg/router/i18n"

// Primary translation function
i18n.T(ctx, "translation_key")
// Returns the translated string for the current locale
// Falls back to "[MISSING_I18N: key]" if translation not found

templ ExamplePage() {
    <h1>{ i18n.T(ctx, "page_title") }</h1>
    <p>{ i18n.T(ctx, "page_subtitle") }</p>
    <p>{ i18n.T(ctx, "nested.key.with.dots") }</p>
}
```

### URL Localization

```go
// Automatically adds locale prefix to URLs
i18n.LocalizeSafeURL(ctx, "/dashboard")
// Returns: templ.SafeURL("/en/dashboard") or templ.SafeURL("/de/dashboard")

templ Navigation() {
    <nav>
        <a href={ i18n.LocalizeSafeURL(ctx, "/admin") }>
            { i18n.T(ctx, "nav_admin") }
        </a>
        <a href={ i18n.LocalizeSafeURL(ctx, "/user/profile") }>
            { i18n.T(ctx, "nav_profile") }
        </a>
    </nav>
}
```

### Context Information Functions

```go
// Get current locale
i18n.GetCurrentLocale(ctx)
// Returns: "en", "de", "fr", etc.

// Get current template path
i18n.GetCurrentTemplate(ctx)
// Returns: "app/locale_/dashboard/page.templ"

// Get all available translation keys for current template
i18n.GetAvailableKeys(ctx)
// Returns: []string{"page_title", "stats_users", "nav_admin", ...}
```

### Route Convenience Functions

```go
// Get current route path (full path including locale)
i18n.GetCurrentRoute(ctx)
// Returns: "/en/dashboard", "/de/user/123/profile", "/login", etc.

// Get current route path with locale stripped (only for actually localized routes)
i18n.GetCurrentRouteWithoutLocale(ctx)
// Returns: "/dashboard" from "/en/dashboard"
// Returns: "/user/123/profile" from "/de/user/123/profile"
// Returns: "/login" from "/login" (unchanged - not localized)
```

### Debug Information

```go
// Debug information panel
templ DebugPanel() {
    <div class="debug-panel">
        <p><strong>Locale:</strong> { i18n.GetCurrentLocale(ctx) }</p>
        <p><strong>Template:</strong> { i18n.GetCurrentTemplate(ctx) }</p>
        <p><strong>Available Keys:</strong> { fmt.Sprint(len(i18n.GetAvailableKeys(ctx))) }</p>
        <p><strong>Current Route:</strong> { i18n.GetCurrentRoute(ctx) }</p>
        <p><strong>Route Without Locale:</strong> { i18n.GetCurrentRouteWithoutLocale(ctx) }</p>
    </div>
}
```

## Advanced Features

### Language Switching with Route Preservation

```go
// Enhanced language switcher that preserves current route
templ LanguageSwitcherAdvanced() {
    currentLocale := i18n.GetCurrentLocale(ctx)
    routeWithoutLocale := i18n.GetCurrentRouteWithoutLocale(ctx)

    <div class="language-switcher">
        <div class="flex space-x-2">
            // Generate language switch links for all supported locales
            for _, locale := range []string{"en", "de", "fr"} {
                if locale == currentLocale {
                    <span class="px-3 py-1 bg-blue-600 text-white rounded text-sm font-semibold">
                        { strings.ToUpper(locale) }
                    </span>
                } else {
                    <a href={ "/" + locale + routeWithoutLocale }
                       class="px-3 py-1 bg-gray-200 hover:bg-gray-300 rounded text-sm">
                        { strings.ToUpper(locale) }
                    </a>
                }
            }
        </div>
    </div>
}
```

### Fallback Translation System

```go
// Custom fallback handling
templ SafeTranslation(key string, fallback string) {
    translation := i18n.T(ctx, key)
    if strings.HasPrefix(translation, "[MISSING_I18N:") {
        <span>{ fallback }</span>
    } else {
        <span>{ translation }</span>
    }
}

templ ExamplePage() {
    <h1>{ SafeTranslation("page_title", "Default Title") }</h1>
    <p>{ SafeTranslation("page_description", "Default description") }</p>
}
```

### Translation Key Validation

```go
// Validate that all required translations exist
templ ValidateTranslations() {
    availableKeys := i18n.GetAvailableKeys(ctx)
    requiredKeys := []string{"page_title", "page_subtitle", "btn_submit", "btn_cancel"}

    for _, key := range requiredKeys {
        if !contains(availableKeys, key) {
            <div class="error">
                Missing translation key: { key }
            </div>
        }
    }
}
```

## Component Internationalization

### Self-Contained Components with i18n

Components can be **self-contained** with their own internationalized metadata. This is a crucial feature for building truly reusable components:

#### Why Self-Contained Components Matter

**Without self-contained components**, you would need to repeat translations in every page:
```yaml
# Repeating translations in EVERY page that uses the footer 😞
# app/home/page.templ.yaml
i18n:
  en: { footer_copyright: "...", footer_privacy: "..." }
  de: { footer_copyright: "...", footer_privacy: "..." }

# app/about/page.templ.yaml
i18n:
  en: { footer_copyright: "...", footer_privacy: "..." }
  de: { footer_copyright: "...", footer_privacy: "..." }

# app/contact/page.templ.yaml
i18n:
  en: { footer_copyright: "...", footer_privacy: "..." }
  de: { footer_copyright: "...", footer_privacy: "..." }
```

**With self-contained components**, define once, use everywhere:
```yaml
# app/components/footer.templ.yaml
i18n:
  en:
    footer_copyright: "© 2024 My Company. All rights reserved."
    footer_privacy: "Privacy Policy"
    footer_contact: "Contact Us"
    footer_sitemap: "Sitemap"
    footer_links: "Quick Links"
    footer_follow_us: "Follow Us"
  de:
    footer_copyright: "© 2024 Meine Firma. Alle Rechte vorbehalten."
    footer_privacy: "Datenschutz"
    footer_contact: "Kontakt"
    footer_sitemap: "Sitemap"
    footer_links: "Quick Links"
    footer_follow_us: "Folgen Sie uns"
  fr:
    footer_copyright: "© 2024 Mon Entreprise. Tous droits réservés."
    footer_privacy: "Politique de confidentialité"
    footer_contact: "Contactez-nous"
    footer_sitemap: "Plan du site"
    footer_links: "Liens rapides"
    footer_follow_us: "Suivez-nous"

metadata:
  company_name: "My Company"
  company_email: "info@company.com"
  company_address: "123 Main Street, City"
  version: "1.0.0"
  social_links:
    twitter: "https://twitter.com/mycompany"
    linkedin: "https://linkedin.com/company/mycompany"
    github: "https://github.com/mycompany"
```

**Benefits:**
1. **No duplication** - Define translations once, not in every page
2. **Easy maintenance** - Update component translations in one place
3. **True reusability** - Components work independently across pages
4. **Consistency** - Same component behaves identically everywhere
5. **Single source of truth** - All component config in one file
6. **Independent testing** - Test components in isolation

#### Usage in Templates

```go
// In any page template - just use the component
<Footer />
// The Footer automatically uses its own translations and metadata
// No need to add anything to the page's .templ.yaml file
```

### Component Template with i18n

```go
// app/components/footer.templ
package components

import (
    "github.com/denkhaus/templ-router/pkg/router/i18n"
    "github.com/denkhaus/templ-router/pkg/router/metadata"
)

templ Footer() {
    companyEmail := metadata.M(ctx, "company_email")
    companyName := metadata.M(ctx, "company_name")
    currentYear := time.Now().Year()

    <footer class="bg-gray-800 text-white p-4 mt-8">
        <div class="container mx-auto">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <!-- Company Info -->
                <div>
                    <h3 class="text-lg font-semibold mb-2">{ companyName }</h3>
                    <p class="text-sm text-gray-300">
                        { metadata.M(ctx, "company_address") }
                    </p>
                    <p class="text-sm text-gray-300">
                        <a href={ "mailto:" + companyEmail } class="hover:text-white">
                            { companyEmail }
                        </a>
                    </p>
                </div>

                <!-- Links -->
                <div>
                    <h3 class="text-lg font-semibold mb-2">
                        { i18n.T(ctx, "footer_links") }
                    </h3>
                    <ul class="text-sm space-y-1">
                        <li>
                            <a href="/privacy" class="text-gray-300 hover:text-white">
                                { i18n.T(ctx, "footer_privacy") }
                            </a>
                        </li>
                        <li>
                            <a href="/contact" class="text-gray-300 hover:text-white">
                                { i18n.T(ctx, "footer_contact") }
                            </a>
                        </li>
                        <li>
                            <a href="/sitemap" class="text-gray-300 hover:text-white">
                                { i18n.T(ctx, "footer_sitemap") }
                            </a>
                        </li>
                    </ul>
                </div>

                <!-- Copyright -->
                <div>
                    <p class="text-sm text-gray-300">
                        { i18n.T(ctx, "footer_copyright", fmt.Sprintf("%d", currentYear)) }
                    </p>
                    <p class="text-xs text-gray-400 mt-2">
                        Version { metadata.M(ctx, "version") }
                    </p>
                </div>
            </div>
        </div>
    </footer>
}
```

## Locale Detection

### URL-Based Detection (Default)

Locale is detected from the URL path:

```bash
# URL patterns
/en/dashboard        → locale = "en"
/de/benutzerprofil   → locale = "de"
/fr/connexion        → locale = "fr"
```

### Header-Based Detection

Detect locale from Accept-Language header:

```yaml
# Configuration
TR_I18N_DETECTION_METHOD=header
TR_I18N_DEFAULT_LOCALE=en
```

```bash
# Request headers
Accept-Language: de-DE,de;q=0.9,en;q=0.8  → locale = "de"
Accept-Language: fr-FR,fr;q=0.9,en;q=0.8  → locale = "fr"
```

### Cookie-Based Detection

Store and detect locale from cookies:

```yaml
# Configuration
TR_I18N_DETECTION_METHOD=cookie
TR_I18N_COOKIE_NAME=locale
TR_I18N_COOKIE_EXPIRY=8760h  # 1 year
```

### Hybrid Detection

Combine multiple detection methods:

```yaml
# Priority: URL > Cookie > Header > Default
TR_I18N_DETECTION_METHOD=url,cookie,header
TR_I18N_DEFAULT_LOCALE=en
```

## SEO Optimization

### hreflang Tags

Generate proper hreflang tags for search engines:

```go
// app/components/hreflang.templ
package components

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ HreflangTags() {
    currentRoute := i18n.GetCurrentRouteWithoutLocale(ctx)
    supportedLocales := []string{"en", "de", "fr"}

    for _, locale := range supportedLocales {
        <link rel="alternate" hreflang={ locale }
              href={ "https://example.com/" + locale + currentRoute } />
    }

    // x-default for international users
    <link rel="alternate" hreflang="x-default"
          href={ "https://example.com/en" + currentRoute } />
}
```

### Localized Meta Tags

```go
// app/components/meta-tags.templ
package components

import "github.com/denkhaus/templ-router/pkg/router/i18n"

templ MetaTags() {
    <title>{ i18n.T(ctx, "meta_title") }</title>
    <meta name="description" content={ i18n.T(ctx, "meta_description") } />
    <meta name="keywords" content={ i18n.T(ctx, "meta_keywords") } />

    // Open Graph tags
    <meta property="og:title" content={ i18n.T(ctx, "og_title") } />
    <meta property="og:description" content={ i18n.T(ctx, "og_description") } />

    // Twitter Card tags
    <meta name="twitter:title" content={ i18n.T(ctx, "twitter_title") } />
    <meta name="twitter:description" content={ i18n.T(ctx, "twitter_description") } />
}
```

## Testing Internationalization

### Test Translation Keys

```go
// pkg/tests/i18n_test.go
package tests

import (
    "context"
    "testing"
    "github.com/denkhaus/templ-router/pkg/router/i18n"
)

func TestTranslationKeys(t *testing.T) {
    // Create test context with locale
    ctx := i18n.WithLocale(context.Background(), "en")

    // Test basic translation
    if got := i18n.T(ctx, "page_title"); got != "Dashboard" {
        t.Errorf("Expected 'Dashboard', got '%s'", got)
    }

    // Test nested translation
    if got := i18n.T(ctx, "nav.main.dashboard"); got != "Dashboard" {
        t.Errorf("Expected 'Dashboard', got '%s'", got)
    }

    // Test missing translation fallback
    if got := i18n.T(ctx, "nonexistent_key"); !strings.Contains(got, "[MISSING_I18N:") {
        t.Errorf("Expected missing translation indicator, got '%s'", got)
    }
}
```

### Test URL Localization

```go
func TestURLLocalization(t *testing.T) {
    tests := []struct {
        locale string
        path   string
        want   string
    }{
        {"en", "/dashboard", "/en/dashboard"},
        {"de", "/dashboard", "/de/dashboard"},
        {"fr", "/user/profile", "/fr/user/profile"},
    }

    for _, tt := range tests {
        ctx := i18n.WithLocale(context.Background(), tt.locale)
        got := i18n.LocalizeSafeURL(ctx, tt.path)

        if string(got) != tt.want {
            t.Errorf("LocalizeSafeURL(%s, %s) = %s, want %s",
                tt.locale, tt.path, got, tt.want)
        }
    }
}
```

### Test Locale Detection

```go
func TestLocaleDetection(t *testing.T) {
    // Test URL-based detection
    req := httptest.NewRequest("GET", "/de/dashboard", nil)
    locale := i18n.DetectLocaleFromURL(req, []string{"en", "de", "fr"})

    if locale != "de" {
        t.Errorf("Expected locale 'de', got '%s'", locale)
    }

    // Test header-based detection
    req = httptest.NewRequest("GET", "/dashboard", nil)
    req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en;q=0.8")
    locale = i18n.DetectLocaleFromHeader(req, []string{"en", "de", "fr"})

    if locale != "fr" {
        t.Errorf("Expected locale 'fr', got '%s'", locale)
    }
}
```

## Best Practices

### Translation Organization

1. **Use nested structures** for complex applications with many translation keys
2. **Group related translations** under common parent keys
3. **Keep nesting levels reasonable** (2-4 levels deep maximum)
4. **Use consistent naming conventions** across translation files
5. **Organize by feature or component** rather than by page

### Key Naming Conventions

```yaml
# Good: Consistent and descriptive
i18n:
  navigation:
    main:
      home: "Home"
      dashboard: "Dashboard"
    user:
      profile: "Profile"
      settings: "Settings"
  buttons:
    primary:
      submit: "Submit"
      save: "Save"
    secondary:
      cancel: "Cancel"

# Avoid: Inconsistent or unclear naming
i18n:
  home_button: "Home"
  main_nav_dashboard: "Dashboard"
  user_profile_link: "Profile"
  submit_btn: "Submit"
```

### Performance Optimization

1. **Lazy load translations**: Load translations only when needed
2. **Cache translations**: Cache parsed translation files
3. **Minimize key lookups**: Cache frequently used translations
4. **Use locale detection hierarchy**: Most specific to least specific

### Maintenance

1. **Track missing translations**: Monitor for missing translation keys
2. **Validate translation completeness**: Ensure all locales have required keys
3. **Regular translation reviews**: Keep translations updated and accurate
4. **Use translation management tools**: Consider external translation services

## Troubleshooting

### Common Issues

#### Translations Not Loading

```bash
# Check locale configuration
env | grep TR_I18N

# Verify translation file format
# Ensure YAML is properly formatted
# Check file naming: page.templ.yaml

# Check template directory structure
# Ensure locale_/ directory is used correctly
```

#### Wrong Locale Detected

```bash
# Check detection method
TR_I18N_DETECTION_METHOD=url,cookie,header

# Verify supported locales
TR_I18N_SUPPORTED_LOCALES=en,de,fr

# Check default locale
TR_I18N_DEFAULT_LOCALE=en
```

#### Missing Translation Keys

```bash
# Enable debug mode
TR_LOGGING_LEVEL=debug

# Check for missing keys in logs
# Validate YAML syntax
# Ensure keys are nested correctly
```

#### URL Generation Issues

```bash
# Check URL prefix configuration
TR_I18N_URL_PREFIX=true

# Verify route generation
# Test trgen output for localized routes

# Check route precedence
# Ensure localized routes are generated correctly
```

### Debug Mode

```bash
# Enable debug logging
TR_LOGGING_LEVEL=debug

# Check i18n-specific logs
# Look for locale detection messages
# Verify translation loading
```

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Set up your first internationalized project
- **[File-Based Routing](FILE-BASED-ROUTING.md)** - Understand localized routing
- **[Authentication](AUTHENTICATION.md)** - Add i18n to authentication flows
- **[Configuration](CONFIGURATION.md)** - Configure i18n behavior
- **[Templates](TEMPLATES.md)** - Advanced template internationalization

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Documentation](../README.md)** - Main documentation