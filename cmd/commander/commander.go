package main

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/config"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
)

const (
	AppName       = "Commander"
	CommanderFile = "commander.yaml"
)

func main() {
	log.SetOutput(ioutil.Discard)

	log.Println("Starting commander")

	app := cli.NewApp()
	app.Name = AppName

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
			ArgsUsage: "[file]",
			Action: func(c *cli.Context) {
				file := CommanderFile
				if c.Args().First() != "" {
					file = c.Args().First()
				}
				content, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Println("Error " + err.Error())
					os.Exit(1)
				}

				suite := config.ParseYAML(content)
				results := runtime.Start(suite)
				if !start(results) {
					os.Exit(1)
				}
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(results <-chan runtime.TestResult) bool {
	testResults := []runtime.TestResult{}
	success := true

	for r := range results {
		testResults = append(testResults, r)
		if r.ValidationResult.Success {
			fmt.Println("✓ ", r.TestCase.Title)
		}

		if !r.ValidationResult.Success {
			success = false
			fmt.Println("✗ ", r.TestCase.Title)
		}
	}

	return success
}
