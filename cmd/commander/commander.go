package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"github.com/urfave/cli"
	"github.com/SimonBaeumer/commander/pkg"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
)

const (
	AppName       = "Commander"
	CommanderFile = "commander.yaml"
)

var version string

func main() {
	log.SetOutput(ioutil.Discard)

	log.Println("Starting commander")

	app := cli.NewApp()
	app.Name = AppName
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
			Action: func(c *cli.Context) {
				file := CommanderFile
				if c.Args().First() != "" {
					file = c.Args().First()
				}
				fmt.Println("Starting test file " + file + "...")
				fmt.Println("")

				content, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Println("Error " + err.Error())
					os.Exit(1)
				}

				var s suite.Suite
				s = suite.ParseYAML(content)

				tests := s.GetTests()
				// Filter tests if test title was given
				if title := c.Args().Get(1); title != "" {
					test, err := s.GetTestByTitle(title)
					if err != nil {
						log.Fatal(err.Error())
						os.Exit(1)
					}
					tests = []runtime.TestCase{test}
				}

				results := runtime.Start(tests)
				if !commander.Start(results) {
					os.Exit(1)
				}
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
