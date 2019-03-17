package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"github.com/urfave/cli"
	"github.com/SimonBaeumer/commander/pkg/output"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
)

const (
	appName       = "Commander"
	commanderFile = "commander.yaml"
)

var version string

func main() {
	log.SetOutput(ioutil.Discard)

	log.Println("Starting commander")

	app := createCliApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "CLI app testing"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "More output for debugging",
			EnvVar: "COMMANDER_VERBOSE",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "test",
			Usage:     "Execute the test suite",
			ArgsUsage: "[file] [test]",
			Action: func(c *cli.Context) error {
				return testCommand(c.Args().First(), c.Args().Get(1))
			},
		},
	}
	return app
}

func testCommand(file string, title string) error {
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

	results := runtime.Start(tests)
	out := output.NewCliOutput()
	if !out.Start(results) {
		return fmt.Errorf("Test suite failed")
	}

	return nil
}
