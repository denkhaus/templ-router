package main

import (
	"log"
	"os"

	"github.com/denkhaus/templ-router/cmd/template-generator/commands"
	"github.com/urfave/cli/v2"
)

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
			Name:    "scan-path",
			Value:   "app",
			Usage:   "Path to scan for templates",
			EnvVars: []string{"TEMPLATE_SCAN_PATH"},
		},
		&cli.StringFlag{
			Name:    "module-name",
			Value:   "github.com/denkhaus/templ-router",
			Usage:   "Go module name",
			EnvVars: []string{"TEMPLATE_MODULE_NAME"},
		},
	}
}

func main() {
	app := &cli.App{
		Name:   "template-generator",
		Usage:  "Generate templates for templ-router",
		Flags:  appFlags(),
		Action: commands.Run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
