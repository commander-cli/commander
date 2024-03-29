package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/urfave/cli"

	"github.com/commander-cli/commander/v2/pkg/app"
)

var version string

func main() {
	run(os.Args)
}

func run(args []string) bool {
	log.SetOutput(io.Discard)

	cliapp := createCliApp()

	if err := cliapp.Run(args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return true
}

func createCliApp() *cli.App {
	cliapp := cli.NewApp()
	cliapp.Name = app.AppName
	cliapp.Usage = "CLI app testing"
	cliapp.Version = version

	cliapp.Commands = []cli.Command{
		createTestCommand(),
		createAddCommand(),
	}
	return cliapp
}

func createTestCommand() cli.Command {
	return cli.Command{
		Name:  "test",
		Usage: "Execute cli app tests",
		UsageText: `Execute cli app tests

By default it will use the commander.yaml from your current directory.
Tests are always executed in alphabetical order.

Examples:

Directory test:
commander test --dir /your/dir/

Stdin test:
cat commander.yaml | commander test -

HTTP test:
commander test https://your-url/commander_test.yaml

Filtering tests:
commander test commander.yaml --filter="my test"

Multiple filters:
commander test commander.yaml --filter=filter1 --filter=filter2

Regex filters:
commander test commander.yaml --filter="^filter1$"
`,
		ArgsUsage: "[file] [--filter]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "no-color",
				EnvVar: "COMMANDER_NO_COLOR",
				Usage:  "Activate or deactivate colored output",
			},
			cli.BoolFlag{
				Name:   "verbose",
				Usage:  "More output for debugging",
				EnvVar: "COMMANDER_VERBOSE",
			},
			cli.BoolFlag{
				Name:  "dir",
				Usage: "Execute all test files in a directory sorted by file name, this is not recursive - e.g. /path/to/test_files/",
			},
			cli.StringFlag{
				Name:  "workdir",
				Usage: "Set the working directory of commander's execution",
			},
			cli.StringFlag{
				Name:  "config",
				Usage: "Default config file to be applied to all test files. Lowest level config values take precednce",
			},
			cli.StringSliceFlag{
				Name:  "filter",
				Usage: `Filter tests by a given regex pattern. Tests are filtered by its title.`,
			},
		},
		Action: func(c *cli.Context) error {
			return app.TestCommand(c.Args().First(), app.NewTestContextFromCli(c))
		},
	}
}

func createAddCommand() cli.Command {
	return cli.Command{
		Name:      "add",
		Usage:     "Automatically add a test to your test suite",
		ArgsUsage: "[command]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "stdout",
				Usage: "Output test file to stdout",
			},
			cli.BoolFlag{
				Name:  "no-file",
				Usage: "Don't create a commander.yaml",
			},
			cli.StringFlag{
				Name:  "file",
				Usage: "Write to another file, default is commander.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			file := ""
			var existedContent []byte

			if !c.Bool("no-file") {
				dir, _ := os.Getwd()
				file = path.Join(dir, app.CommanderFile)
				if c.String("file") != "" {
					file = c.String("file")
				}
				existedContent, _ = os.ReadFile(file)
			}

			content, err := app.AddCommand(strings.Join(c.Args(), " "), existedContent)
			if err != nil {
				return err
			}

			if c.Bool("stdout") {
				fmt.Println(string(content))
			}
			if !c.Bool("no-file") {
				fmt.Println("written to", file)
				err := os.WriteFile(file, content, 0o755)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
