package main

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	AppName       = "Commander"
	CommanderFile = "commander.yaml"
)

func main() {
	tests := []suite.TestCase{
		{
			Command: suite.CommandUnderTest{
				Cmd: "echo hello",
			},
			Expected: suite.Expected{
				Stdout: suite.ExpectedOut{
					Exactly: "hello",
				},
				ExitCode: 0,
			},
			Title: "Output hello",
		},
	}

	s := suite.NewSuite(tests)
	r := start(*s)

	if !r {
		os.Exit(1)
	}
	os.Exit(0)

	os.Args = []string{"./commander", "test", "my.yml"}
	//log.SetOutput(ioutil.Discard)

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
//				suite := commander.ParseYAMLFile(file)
//				runtime.Start(&suite)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(s suite.Suite) bool {
	testResults := []runtime.TestResult{}
	results := runtime.Start(s)
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
