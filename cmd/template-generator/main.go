package main

import (
	"fmt"
	"log"
	"os"

	"github.com/denkhaus/templ-router/cmd/template-generator/commands"
	"github.com/denkhaus/templ-router/cmd/template-generator/version"
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
	buildInfo := version.GetBuildInfo()
	
	app := &cli.App{
		Name:    "template-generator",
		Usage:   "Generate templates for templ-router",
		Version: buildInfo.Short(),
		Flags:   appFlags(),
		Action: func(c *cli.Context) error {
			// Always show version at start of generation
			fmt.Printf("Template Generator %s\n", buildInfo.String())
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
