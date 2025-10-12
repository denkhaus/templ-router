package template

import (
	"context"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

// DiscoveryComponent handles creation of discovery-related template components
type DiscoveryComponent struct{}

// NewDiscoveryComponent creates a new discovery component generator
func NewDiscoveryComponent() *DiscoveryComponent {
	return &DiscoveryComponent{}
}

// CreateDiscoveredTemplateComponent creates a component showing discovered template info
func (dc *DiscoveryComponent) CreateDiscoveredTemplateComponent(templateKey, functionName string, params map[string]string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Discovered Template: %s</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-blue-50">
    <div class="min-h-screen py-8">
        <div class="max-w-4xl mx-auto px-4">
            <div class="bg-white rounded-lg shadow-lg p-8 border-l-4 border-blue-500">
                <h1 class="text-3xl font-bold text-blue-800 mb-4">Template Discovered!</h1>
                <div class="bg-blue-100 border border-blue-200 rounded-lg p-6">
                    <h3 class="text-lg font-semibold text-blue-800 mb-2">AST-Based Template Discovery</h3>
                    <p class="text-blue-700 mb-4">This template was found by scanning the generated Go files!</p>
                    <div class="text-blue-600 text-sm space-y-1">
                        <p><strong>Template Key:</strong> %s</p>
                        <p><strong>Function Name:</strong> %s</p>
                        <p><strong>Discovery Method:</strong> AST Parsing of *_templ.go files</p>
                        <p><strong>Status:</strong> Template-Agnostic Router Working</p>
                    </div>
                </div>

                <div class="mt-6 bg-green-50 border border-green-200 rounded-lg p-4">
                    <h4 class="font-semibold text-green-800 mb-2">Success!</h4>
                    <ul class="text-green-700 text-sm space-y-1">
                        <li>✅ App directory scanned successfully</li>
                        <li>✅ Generated Go files parsed</li>
                        <li>✅ Template functions discovered</li>
                        <li>✅ Router remains template-agnostic</li>
                        <li>✅ No static imports of app package</li>
                    </ul>
                </div>`, templateKey, templateKey, functionName)

		if len(params) > 0 {
			html += `<div class="mt-6 bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                    <h4 class="font-semibold text-yellow-800 mb-2">Parameters:</h4>
                    <div class="text-yellow-600 text-sm">`
			for key, value := range params {
				html += fmt.Sprintf(`<p><strong>%s:</strong> %s</p>`, key, value)
			}
			html += `</div></div>`
		}

		html += `            </div>
        </div>
    </div>
</body>
</html>`

		_, err := w.Write([]byte(html))
		return err
	})
}

// CreateNotFoundComponent creates a component for when a template is not found
func (dc *DiscoveryComponent) CreateNotFoundComponent(templateKey string, availableKeys []string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Template Not Found: %s</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-red-50">
    <div class="min-h-screen py-8">
        <div class="max-w-4xl mx-auto px-4">
            <div class="bg-white rounded-lg shadow-lg p-8 border-l-4 border-red-500">
                <h1 class="text-3xl font-bold text-red-800 mb-4">Template Not Found</h1>
                <div class="bg-red-100 border border-red-200 rounded-lg p-6">
                    <p class="text-red-700 mb-4">The requested template could not be found:</p>
                    <p class="text-red-600 font-mono text-sm"><strong>%s</strong></p>
                </div>`, templateKey, templateKey)

		if len(availableKeys) > 0 {
			html += `<div class="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <h4 class="font-semibold text-blue-800 mb-2">Available Templates:</h4>
                    <div class="text-blue-600 text-sm space-y-1">`
			for _, key := range availableKeys {
				html += fmt.Sprintf(`<p class="font-mono">%s</p>`, key)
			}
			html += `</div></div>`
		}

		html += `            </div>
        </div>
    </div>
</body>
</html>`

		_, err := w.Write([]byte(html))
		return err
	})
}