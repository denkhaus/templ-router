package commands

import (
	"fmt"

	"github.com/denkhaus/templ-router/cmd/template-generator/generate"
	"github.com/denkhaus/templ-router/cmd/template-generator/scan"
	"github.com/denkhaus/templ-router/cmd/template-generator/types"

	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	config := newConfig(c)
	runGeneration(config)

	if c.Bool("watch") {
		watchedExtensions := strings.Split(c.String("watch-extensions"), ",")
		return startWatcher(config, watchedExtensions)
	}

	return nil
}

// CRITICAL: This generator MUST be 100% configuration-agnostic!
// These are the ONLY acceptable defaults - generic output paths that work everywhere.
// NEVER add defaults for module names, scan paths, or project-specific values!
func newConfig(c *cli.Context) types.Config {
	return types.Config{
		ModuleName:  c.String("module-name"),  // NO DEFAULT - must be provided
		ScanPath:    c.String("scan-path"),    // NO DEFAULT - must be provided  
		OutputDir:   "generated/templates",    // Generic output path - OK
		PackageName: "templates",              // Generic package name - OK
	}
}

func startWatcher(config types.Config, extensions []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go handleFileEvents(watcher, config, extensions)

	if err := addPathsToWatcher(watcher, config.ScanPath); err != nil {
		return err
	}

	log.Println("Watching for changes in", config.ScanPath)
	<-make(chan struct{})
	return nil
}

func handleFileEvents(watcher *fsnotify.Watcher, config types.Config, extensions []string) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if isWatchedEvent(event) && hasWatchedExtension(event.Name, extensions) {
				log.Println("Modified file:", event.Name)
				runGeneration(config)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func isWatchedEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Write == fsnotify.Write ||
		event.Op&fsnotify.Create == fsnotify.Create ||
		event.Op&fsnotify.Remove == fsnotify.Remove
}

func hasWatchedExtension(fileName string, extensions []string) bool {
	ext := filepath.Ext(fileName)
	for _, watchedExt := range extensions {
		if ext == watchedExt {
			return true
		}
	}
	return false
}

func addPathsToWatcher(watcher *fsnotify.Watcher, scanPath string) error {
	return filepath.Walk(scanPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}

func runGeneration(config types.Config) {
	fmt.Println("Starting template generation...")
	fmt.Printf("Scanning path: %s\n", config.ScanPath)
	fmt.Printf("Output directory: %s\n", config.OutputDir)

	// Scan for templates
	templates, validationErrors, err := scan.ScanTemplatesWithPackages(config)
	if err != nil {
		log.Printf("Failed to scan templates: %v", err)
		return
	}

	// Print validation errors as warnings
	if len(validationErrors) > 0 {
		fmt.Printf("\nValidation warnings:\n")
		for _, warning := range validationErrors {
			fmt.Printf("  WARNING: %s\n", warning)
		}
		fmt.Println()
	}

	fmt.Printf("Found %d templates\n", len(templates))
	for _, template := range templates {
		fmt.Printf("  %s -> %s (%s)\n", template.RoutePattern, template.FunctionName, template.HumanName)
	}

	// Generate registry
	// Validate required parameters
	if config.ModuleName == "" {
		log.Printf("TEMPLATE_MODULE_NAME is required and cannot be empty")
		return
	}

	if err := generate.GenerateRegistry(config, templates); err != nil {
		log.Printf("Failed to generate registry: %v", err)
		return
	}

	fmt.Println("Template generation completed successfully!")
}
