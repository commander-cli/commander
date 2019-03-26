package main

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/output"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
)

const (
	appName       = "Commander"
	commanderFile = "commander.yaml"
)

var version string

type CommanderContext struct {
	Verbose    bool
	NoColor    bool
	Concurrent int
}

func NewContextFromCli(c *cli.Context) CommanderContext {
	return CommanderContext{
		Verbose:    c.Bool("verbose"),
		NoColor:    c.Bool("no-color"),
		Concurrent: c.Int("concurrent"),
	}
}

func main() {
	log.SetOutput(ioutil.Discard)

	app := createCliApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func createCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "CLI app testing"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:      "test",
			Usage:     "Execute the test suite",
			ArgsUsage: "[file] [test]",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "concurrent",
					EnvVar: "COMMANDER_CONCURRENT",
					Usage:  "Set the max amount of tests which should run concurrently",
				},
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
			},
			Action: func(c *cli.Context) error {
				return testCommand(c.Args().First(), c.Args().Get(1), NewContextFromCli(c))
			},
		},
	}
	return app
}

func testCommand(file string, title string, ctx CommanderContext) error {
	log.SetOutput(ioutil.Discard)
	if ctx.Verbose == true {
		log.SetOutput(os.Stdout)
	}

	if file == "" {
		file = commanderFile
	}

	fmt.Println("Starting test file " + file + "...")
	fmt.Println("")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Error " + err.Error())
	}

	var s suite.Suite
	s = suite.ParseYAML(content)
	tests := s.GetTests()
	// Filter tests if test title was given
	if title != "" {
		test, err := s.GetTestByTitle(title)
		if err != nil {
			return err
		}
		tests = []runtime.TestCase{test}
	}

	results := runtime.Start(tests, ctx.Concurrent)
	out := output.NewCliOutput(!ctx.NoColor)
	if !out.Start(results) {
		return fmt.Errorf("Test suite failed")
	}

	return nil
}
