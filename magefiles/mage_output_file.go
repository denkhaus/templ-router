// +build ignore

package main

import (
	"context"
	_flag "flag"
	_fmt "fmt"
	_ioutil "io/ioutil"
	_log "log"
	"os"
	"os/signal"
	_filepath "path/filepath"
	_sort "sort"
	"strconv"
	_strings "strings"
	"syscall"
	_tabwriter "text/tabwriter"
	"time"
	
)

func main() {
	// Use local types and functions in order to avoid name conflicts with additional magefiles.
	type arguments struct {
		Verbose       bool          // print out log statements
		List          bool          // print out a list of targets
		Help          bool          // print out help for a specific target
		Timeout       time.Duration // set a timeout to running the targets
		Args          []string      // args contain the non-flag command-line arguments
	}

	parseBool := func(env string) bool {
		val := os.Getenv(env)
		if val == "" {
			return false
		}		
		b, err := strconv.ParseBool(val)
		if err != nil {
			_log.Printf("warning: environment variable %s is not a valid bool value: %v", env, val)
			return false
		}
		return b
	}

	parseDuration := func(env string) time.Duration {
		val := os.Getenv(env)
		if val == "" {
			return 0
		}		
		d, err := time.ParseDuration(val)
		if err != nil {
			_log.Printf("warning: environment variable %s is not a valid duration value: %v", env, val)
			return 0
		}
		return d
	}
	args := arguments{}
	fs := _flag.FlagSet{}
	fs.SetOutput(os.Stdout)

	// default flag set with ExitOnError and auto generated PrintDefaults should be sufficient
	fs.BoolVar(&args.Verbose, "v", parseBool("MAGEFILE_VERBOSE"), "show verbose output when running targets")
	fs.BoolVar(&args.List, "l", parseBool("MAGEFILE_LIST"), "list targets for this binary")
	fs.BoolVar(&args.Help, "h", parseBool("MAGEFILE_HELP"), "print out help for a specific target")
	fs.DurationVar(&args.Timeout, "t", parseDuration("MAGEFILE_TIMEOUT"), "timeout in duration parsable format (e.g. 5m30s)")
	fs.Usage = func() {
		_fmt.Fprintf(os.Stdout, `
%s [options] [target]

Commands:
  -l    list targets in this binary
  -h    show this help

Options:
  -h    show description of a target
  -t <string>
        timeout in duration parsable format (e.g. 5m30s)
  -v    show verbose output when running targets
 `[1:], _filepath.Base(os.Args[0]))
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		// flag will have printed out an error already.
		return
	}
	args.Args = fs.Args()
	if args.Help && len(args.Args) == 0 {
		fs.Usage()
		return
	}
		
	// color is ANSI color type
	type color int

	// If you add/change/remove any items in this constant,
	// you will need to run "stringer -type=color" in this directory again.
	// NOTE: Please keep the list in an alphabetical order.
	const (
		black color = iota
		red
		green
		yellow
		blue
		magenta
		cyan
		white
		brightblack
		brightred
		brightgreen
		brightyellow
		brightblue
		brightmagenta
		brightcyan
		brightwhite
	)

	// AnsiColor are ANSI color codes for supported terminal colors.
	var ansiColor = map[color]string{
		black:         "\u001b[30m",
		red:           "\u001b[31m",
		green:         "\u001b[32m",
		yellow:        "\u001b[33m",
		blue:          "\u001b[34m",
		magenta:       "\u001b[35m",
		cyan:          "\u001b[36m",
		white:         "\u001b[37m",
		brightblack:   "\u001b[30;1m",
		brightred:     "\u001b[31;1m",
		brightgreen:   "\u001b[32;1m",
		brightyellow:  "\u001b[33;1m",
		brightblue:    "\u001b[34;1m",
		brightmagenta: "\u001b[35;1m",
		brightcyan:    "\u001b[36;1m",
		brightwhite:   "\u001b[37;1m",
	}
	
	const _color_name = "blackredgreenyellowbluemagentacyanwhitebrightblackbrightredbrightgreenbrightyellowbrightbluebrightmagentabrightcyanbrightwhite"

	var _color_index = [...]uint8{0, 5, 8, 13, 19, 23, 30, 34, 39, 50, 59, 70, 82, 92, 105, 115, 126}

	colorToLowerString := func (i color) string {
		if i < 0 || i >= color(len(_color_index)-1) {
			return "color(" + strconv.FormatInt(int64(i), 10) + ")"
		}
		return _color_name[_color_index[i]:_color_index[i+1]]
	}

	// ansiColorReset is an ANSI color code to reset the terminal color.
	const ansiColorReset = "\033[0m"

	// defaultTargetAnsiColor is a default ANSI color for colorizing targets.
	// It is set to Cyan as an arbitrary color, because it has a neutral meaning
	var defaultTargetAnsiColor = ansiColor[cyan]

	getAnsiColor := func(color string) (string, bool) {
		colorLower := _strings.ToLower(color)
		for k, v := range ansiColor {
			colorConstLower := colorToLowerString(k)
			if colorConstLower == colorLower {
				return v, true
			}
		}
		return "", false
	}

	// Terminals which  don't support color:
	// 	TERM=vt100
	// 	TERM=cygwin
	// 	TERM=xterm-mono
    var noColorTerms = map[string]bool{
		"vt100":      false,
		"cygwin":     false,
		"xterm-mono": false,
	}

	// terminalSupportsColor checks if the current console supports color output
	//
	// Supported:
	// 	linux, mac, or windows's ConEmu, Cmder, putty, git-bash.exe, pwsh.exe
	// Not supported:
	// 	windows cmd.exe, powerShell.exe
	terminalSupportsColor := func() bool {
		envTerm := os.Getenv("TERM")
		if _, ok := noColorTerms[envTerm]; ok {
			return false
		}
		return true
	}

	// enableColor reports whether the user has requested to enable a color output.
	enableColor := func() bool {
		b, _ := strconv.ParseBool(os.Getenv("MAGEFILE_ENABLE_COLOR"))
		return b
	}

	// targetColor returns the ANSI color which should be used to colorize targets.
	targetColor := func() string {
		s, exists := os.LookupEnv("MAGEFILE_TARGET_COLOR")
		if exists == true {
			if c, ok := getAnsiColor(s); ok == true {
				return c
			}
		}
		return defaultTargetAnsiColor
	}

	// store the color terminal variables, so that the detection isn't repeated for each target
	var enableColorValue = enableColor() && terminalSupportsColor()
	var targetColorValue = targetColor()

	printName := func(str string) string {
		if enableColorValue {
			return _fmt.Sprintf("%s%s%s", targetColorValue, str, ansiColorReset)
		} else {
			return str
		}
	}

	list := func() error {
		
		targets := map[string]string{
			"build:registryGenerate": "",
			"build:registryWatch": "",
			"build:tailwindClean": "builds Tailwind CSS",
			"build:tailwindWatch": "builds Tailwind CSS in watch mode",
			"build:templGenerate": "generates Templ templates",
			"build:templWatch": "runs templ generation in watch mode",
			"clean": "removes build artifacts and temporary files",
			"dev*": "starts the development environment with parallel watch processes",
			"docker:build": "builds the Docker image",
			"docker:clean": "",
			"docker:down": "DevDown stops development Docker services",
			"docker:logs": "shows Docker logs",
			"docker:rebuild": "rebuilds Docker images without cache",
			"docker:startClean": "",
			"docker:up": "DevUp starts development Docker services",
			"generator:build": "builds the template generator without installing",
			"generator:clean": "removes build artifacts",
			"generator:dev": "installs the generator in development mode (rebuilds on changes)",
			"generator:install": "builds and installs the template generator with proper versioning",
			"generator:installGlobal": "Ensure the binary is in PATH",
			"generator:release": "builds the template generator for multiple platforms",
			"generator:test": "runs all tests for the template generator",
			"generator:testCoverage": "runs tests with coverage for the template generator",
			"generator:version": "shows the current version information",
			"templ:install": "",
			"test:all": "runs all tests (unit + E2E)",
			"test:ci": "runs tests in CI mode with coverage",
			"test:checkService": "verifies that the Docker service is running",
			"test:devSetup": "sets up complete development testing environment",
			"test:e2E": "runs end-to-end tests against the Docker service",
			"test:e2EAuth": "runs authentication tests",
			"test:e2EContent": "runs content validation tests",
			"test:e2EData": "runs data service tests",
			"test:e2EI18n": "runs internationalization tests",
			"test:e2EMetadata": "runs metadata and layout tests",
			"test:e2EPerf": "runs performance tests",
			"test:e2ERouting": "runs routing-specific tests",
			"test:e2ESmoke": "runs quick smoke tests",
			"test:e2EWatch": "runs E2E tests in watch mode for development",
			"test:setupE2E": "installs E2E test dependencies",
		}

		keys := make([]string, 0, len(targets))
		for name := range targets {
			keys = append(keys, name)
		}
		_sort.Strings(keys)

		_fmt.Println("Targets:")
		w := _tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		for _, name := range keys {
			_fmt.Fprintf(w, "  %v\t%v\n", printName(name), targets[name])
		}
		err := w.Flush()
			if err == nil {
				_fmt.Println("\n* default target")
			}
		return err
	}

	var ctx context.Context
	ctxCancel := func(){}

	// by deferring in a closure, we let the cancel function get replaced
	// by the getContext function.
	defer func() {
		ctxCancel()
	}()

	getContext := func() (context.Context, func()) {
		if ctx == nil {
			if args.Timeout != 0 {
				ctx, ctxCancel = context.WithTimeout(context.Background(), args.Timeout)
			} else {
				ctx, ctxCancel = context.WithCancel(context.Background())
			}
		}

		return ctx, ctxCancel
	}

	runTarget := func(logger *_log.Logger, fn func(context.Context) error) interface{} {
		var err interface{}
		ctx, cancel := getContext()
		d := make(chan interface{})
		go func() {
			defer func() {
				err := recover()
				d <- err
			}()
			err := fn(ctx)
			d <- err
		}()
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT)
		select {
		case <-sigCh:
			logger.Println("cancelling mage targets, waiting up to 5 seconds for cleanup...")
			cancel()
			cleanupCh := time.After(5 * time.Second)

			select {
			// target exited by itself
			case err = <-d:
				return err
			// cleanup timeout exceeded
			case <-cleanupCh:
				return _fmt.Errorf("cleanup timeout exceeded")
			// second SIGINT received
			case <-sigCh:
				logger.Println("exiting mage")
				return _fmt.Errorf("exit forced")
			}
		case <-ctx.Done():
			cancel()
			e := ctx.Err()
			_fmt.Printf("ctx err: %v\n", e)
			return e
		case err = <-d:
			// we intentionally don't cancel the context here, because
			// the next target will need to run with the same context.
			return err
		}
	}
	// This is necessary in case there aren't any targets, to avoid an unused
	// variable error.
	_ = runTarget

	handleError := func(logger *_log.Logger, err interface{}) {
		if err != nil {
			logger.Printf("Error: %+v\n", err)
			type code interface {
				ExitStatus() int
			}
			if c, ok := err.(code); ok {
				os.Exit(c.ExitStatus())
			}
			os.Exit(1)
		}
	}
	_ = handleError

	// Set MAGEFILE_VERBOSE so mg.Verbose() reflects the flag value.
	if args.Verbose {
		os.Setenv("MAGEFILE_VERBOSE", "1")
	} else {
		os.Setenv("MAGEFILE_VERBOSE", "0")
	}

	_log.SetFlags(0)
	if !args.Verbose {
		_log.SetOutput(_ioutil.Discard)
	}
	logger := _log.New(os.Stderr, "", 0)
	if args.List {
		if err := list(); err != nil {
			_log.Println(err)
			os.Exit(1)
		}
		return
	}

	if args.Help {
		if len(args.Args) < 1 {
			logger.Println("no target specified")
			os.Exit(2)
		}
		switch _strings.ToLower(args.Args[0]) {
			case "build:registrygenerate":
				
				_fmt.Print("Usage:\n\n\tmage build:registrygenerate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "build:registrywatch":
				
				_fmt.Print("Usage:\n\n\tmage build:registrywatch\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "build:tailwindclean":
				_fmt.Println("TailwindClean builds Tailwind CSS")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage build:tailwindclean\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "build:tailwindwatch":
				_fmt.Println("TailwindWatch builds Tailwind CSS in watch mode")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage build:tailwindwatch\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "build:templgenerate":
				_fmt.Println("TemplGenerate generates Templ templates")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage build:templgenerate\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "build:templwatch":
				_fmt.Println("TemplWatch runs templ generation in watch mode")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage build:templwatch\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "clean":
				_fmt.Println("Clean removes build artifacts and temporary files")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage clean\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "dev":
				_fmt.Println("Dev starts the development environment with parallel watch processes")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage dev\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "docker:build":
				_fmt.Println("Build builds the Docker image")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage docker:build\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "docker:clean":
				
				_fmt.Print("Usage:\n\n\tmage docker:clean\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "docker:down":
				_fmt.Println("DevDown stops development Docker services")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage docker:down\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "docker:logs":
				_fmt.Println("Logs shows Docker logs")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage docker:logs\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "docker:rebuild":
				_fmt.Println("Rebuild rebuilds Docker images without cache")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage docker:rebuild\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "docker:startclean":
				
				_fmt.Print("Usage:\n\n\tmage docker:startclean\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "docker:up":
				_fmt.Println("DevUp starts development Docker services")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage docker:up\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:build":
				_fmt.Println("Build builds the template generator without installing")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:build\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:clean":
				_fmt.Println("Clean removes build artifacts")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:clean\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:dev":
				_fmt.Println("Dev installs the generator in development mode (rebuilds on changes)")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:dev\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:install":
				_fmt.Println("Install builds and installs the template generator with proper versioning")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:install\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:installglobal":
				_fmt.Println("Ensure the binary is in PATH")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:installglobal\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:release":
				_fmt.Println("Release builds the template generator for multiple platforms")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:release\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:test":
				_fmt.Println("Test runs all tests for the template generator")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:test\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:testcoverage":
				_fmt.Println("TestCoverage runs tests with coverage for the template generator")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:testcoverage\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "generator:version":
				_fmt.Println("Version shows the current version information")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage generator:version\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "templ:install":
				
				_fmt.Print("Usage:\n\n\tmage templ:install\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:all":
				_fmt.Println("All runs all tests (unit + E2E)")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:all\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:ci":
				_fmt.Println("CI runs tests in CI mode with coverage")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:ci\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:checkservice":
				_fmt.Println("CheckService verifies that the Docker service is running")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:checkservice\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:devsetup":
				_fmt.Println("DevSetup sets up complete development testing environment")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:devsetup\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2e":
				_fmt.Println("E2E runs end-to-end tests against the Docker service")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2e\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2eauth":
				_fmt.Println("E2EAuth runs authentication tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2eauth\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2econtent":
				_fmt.Println("E2EContent runs content validation tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2econtent\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2edata":
				_fmt.Println("E2EData runs data service tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2edata\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2ei18n":
				_fmt.Println("E2EI18n runs internationalization tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2ei18n\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2emetadata":
				_fmt.Println("E2EMetadata runs metadata and layout tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2emetadata\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2eperf":
				_fmt.Println("E2EPerf runs performance tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2eperf\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2erouting":
				_fmt.Println("E2ERouting runs routing-specific tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2erouting\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2esmoke":
				_fmt.Println("E2ESmoke runs quick smoke tests")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2esmoke\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:e2ewatch":
				_fmt.Println("E2EWatch runs E2E tests in watch mode for development")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:e2ewatch\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			case "test:setupe2e":
				_fmt.Println("SetupE2E installs E2E test dependencies")
				_fmt.Println()
				
				_fmt.Print("Usage:\n\n\tmage test:setupe2e\n\n")
				var aliases []string
				if len(aliases) > 0 {
					_fmt.Printf("Aliases: %s\n\n", _strings.Join(aliases, ", "))
				}
				return
			default:
				logger.Printf("Unknown target: %q\n", args.Args[0])
				os.Exit(2)
		}
	}
	if len(args.Args) < 1 {
		ignoreDefault, _ := strconv.ParseBool(os.Getenv("MAGEFILE_IGNOREDEFAULT"))
		if ignoreDefault {
			if err := list(); err != nil {
				logger.Println("Error:", err)
				os.Exit(1)
			}
			return
		}
		
				wrapFn := func(ctx context.Context) error {
					return Dev()
				}
				ret := runTarget(logger, wrapFn)
		handleError(logger, ret)
		return
	}
	for x := 0; x < len(args.Args); {
		target := args.Args[x]
		x++

		// resolve aliases
		switch _strings.ToLower(target) {
		
		}

		switch _strings.ToLower(target) {
		
			case "build:registrygenerate":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:RegistryGenerate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:RegistryGenerate")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.RegistryGenerate()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "build:registrywatch":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:RegistryWatch\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:RegistryWatch")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.RegistryWatch()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "build:tailwindclean":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:TailwindClean\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:TailwindClean")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.TailwindClean()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "build:tailwindwatch":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:TailwindWatch\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:TailwindWatch")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.TailwindWatch()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "build:templgenerate":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:TemplGenerate\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:TemplGenerate")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.TemplGenerate()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "build:templwatch":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Build:TemplWatch\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Build:TemplWatch")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Build{}.TemplWatch()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "clean":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Clean\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Clean")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Clean()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "dev":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Dev\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Dev")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Dev()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "docker:build":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker:Build\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker:Build")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker{}.Build()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "docker:clean":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker:Clean\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker:Clean")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker{}.Clean()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "docker:down":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker:Down\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker:Down")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker{}.Down()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "docker:logs":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker:Logs\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker:Logs")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker{}.Logs()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "docker:rebuild":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker:Rebuild\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker:Rebuild")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker{}.Rebuild()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "docker:startclean":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker:StartClean\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker:StartClean")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker{}.StartClean()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "docker:up":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Docker:Up\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Docker:Up")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Docker{}.Up()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:build":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:Build\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:Build")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.Build()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:clean":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:Clean\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:Clean")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.Clean()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:dev":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:Dev\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:Dev")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.Dev()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:install":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:Install\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:Install")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.Install()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:installglobal":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:InstallGlobal\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:InstallGlobal")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.InstallGlobal()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:release":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:Release\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:Release")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.Release()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:test":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:Test\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:Test")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.Test()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:testcoverage":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:TestCoverage\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:TestCoverage")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.TestCoverage()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "generator:version":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Generator:Version\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Generator:Version")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Generator{}.Version()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "templ:install":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Templ:Install\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Templ:Install")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Templ{}.Install()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:all":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:All\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:All")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.All()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:ci":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:CI\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:CI")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.CI()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:checkservice":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:CheckService\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:CheckService")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.CheckService()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:devsetup":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:DevSetup\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:DevSetup")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.DevSetup()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2e":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2E\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2E")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2E()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2eauth":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2EAuth\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2EAuth")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2EAuth()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2econtent":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2EContent\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2EContent")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2EContent()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2edata":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2EData\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2EData")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2EData()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2ei18n":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2EI18n\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2EI18n")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2EI18n()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2emetadata":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2EMetadata\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2EMetadata")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2EMetadata()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2eperf":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2EPerf\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2EPerf")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2EPerf()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2erouting":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2ERouting\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2ERouting")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2ERouting()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2esmoke":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2ESmoke\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2ESmoke")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2ESmoke()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:e2ewatch":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:E2EWatch\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:E2EWatch")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.E2EWatch()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
			case "test:setupe2e":
				expected := x + 0
				if expected > len(args.Args) {
					// note that expected and args at this point include the arg for the target itself
					// so we subtract 1 here to show the number of args without the target.
					logger.Printf("not enough arguments for target \"Test:SetupE2E\", expected %v, got %v\n", expected-1, len(args.Args)-1)
					os.Exit(2)
				}
				if args.Verbose {
					logger.Println("Running target:", "Test:SetupE2E")
				}
				
				wrapFn := func(ctx context.Context) error {
					return Test{}.SetupE2E()
				}
				ret := runTarget(logger, wrapFn)
				handleError(logger, ret)
		
		default:
			logger.Printf("Unknown target specified: %q\n", target)
			os.Exit(2)
		}
	}
}




