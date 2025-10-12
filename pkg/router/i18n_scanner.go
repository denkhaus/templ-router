package router

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// ScanTemplateForTranslationKeys scans a template file for router.T() calls
func (reg *I18nRegistry) ScanTemplateForTranslationKeys(templatePath string) error {
	// Convert .templ to _templ.go file
	templGoPath := strings.TrimSuffix(templatePath, ".templ") + "_templ.go"

	// Check if the generated Go file exists
	if _, err := os.Stat(templGoPath); os.IsNotExist(err) {
		reg.logger.Debug("Generated template file not found, skipping scan",
			zap.String("template_path", templatePath),
			zap.String("go_file", templGoPath))
		return nil
	}

	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, templGoPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse generated template file %s: %w", templGoPath, err)
	}

	// Walk the AST to find router.T() calls
	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "router" && sel.Sel.Name == "T" {
					// Found router.T() call, extract the key
					if len(call.Args) >= 2 {
						if lit, ok := call.Args[1].(*ast.BasicLit); ok && lit.Kind == token.STRING {
							// Remove quotes from string literal
							key := strings.Trim(lit.Value, `"`)
							reg.RegisterRequiredKey(templatePath, key)
						}
					}
				}
			}
		}
		return true
	})

	return nil
}

// ScanAllTemplatesForTranslationKeys scans all template files for translation keys
func (reg *I18nRegistry) ScanAllTemplatesForTranslationKeys(templFiles []string) error {
	for _, templFile := range templFiles {
		if err := reg.ScanTemplateForTranslationKeys(templFile); err != nil {
			reg.logger.Warn("Failed to scan template for translation keys",
				zap.String("template", templFile),
				zap.Error(err))
			// Continue with other files instead of failing
		}
	}

	reg.logger.Info("Completed scanning templates for translation keys",
		zap.Int("templates_scanned", len(templFiles)))

	return nil
}

// ScanTemplateFileForKeys scans the original .templ file for translation keys using regex
func (reg *I18nRegistry) ScanTemplateFileForKeys(templatePath string) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}

	// Regex to find router.T(ctx, "key") patterns
	re := regexp.MustCompile(`router\.T\s*\(\s*ctx\s*,\s*"([^"]+)"\s*\)`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		if len(match) > 1 {
			key := match[1]
			reg.RegisterRequiredKey(templatePath, key)
		}
	}

	reg.logger.Debug("Scanned template file for translation keys",
		zap.String("template_path", templatePath),
		zap.Int("keys_found", len(matches)))

	return nil
}

// ScanAllTemplateFilesForKeys scans all .templ files for translation keys using regex
func (reg *I18nRegistry) ScanAllTemplateFilesForKeys(templFiles []string) error {
	for _, templFile := range templFiles {
		if err := reg.ScanTemplateFileForKeys(templFile); err != nil {
			reg.logger.Warn("Failed to scan template file for translation keys",
				zap.String("template", templFile),
				zap.Error(err))
			// Continue with other files instead of failing
		}
	}

	reg.logger.Info("Completed scanning template files for translation keys",
		zap.Int("templates_scanned", len(templFiles)))

	return nil
}
