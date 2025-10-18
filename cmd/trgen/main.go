package main

import (
	"fmt"
	"log"
	"os"

	"github.com/denkhaus/templ-router/cmd/trgen/commands"
	"github.com/denkhaus/templ-router/cmd/trgen/version"
	"github.com/urfave/cli/v2"
)

// CRITICAL: This generator MUST be 100% configuration-agnostic!
// It's a LIBRARY for thousands of developers, NOT a local tool.
// NEVER hardcode project names, paths, or module names.
// EVERY project has different structures and names.
// NO DEFAULTS for project-specific values!
func appFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "watch",
			Usage:   "Watch for file changes and regenerate templates",
			EnvVars: []string{"TEMPLATE_WATCH_MODE"},
		},
		&cli.StringFlag{
			Name:    "watch-extensions",
			Value:   ".templ,.yaml,.yml",
			Usage:   "Comma-separated list of file extensions to watch",
			EnvVars: []string{"TEMPLATE_WATCH_EXTENSIONS"},
		},
		&cli.StringFlag{
			Name:     "scan-path",
			Usage:    "Path to scan for templates (required)",
			EnvVars:  []string{"TEMPLATE_SCAN_PATH"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "module-name", 
			Usage:    "Go module name (required)",
			EnvVars:  []string{"TEMPLATE_MODULE_NAME"},
			Required: true,
		},
	}
}

func main() {
	buildInfo := version.GetBuildInfo()
	
	app := &cli.App{
		Name:    "trgen",
		Usage:   "templ-router generator - Generate templates for templ-router",
		Version: buildInfo.Short(),
		Flags:   appFlags(),
		Action: func(c *cli.Context) error {
			// Always show version at start of generation
			fmt.Printf("trgen %s\n", buildInfo.String())
			fmt.Println()
			return commands.Run(c)
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version information",
				Action: func(c *cli.Context) error {
					fmt.Println(buildInfo.String())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
