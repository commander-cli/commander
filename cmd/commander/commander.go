package main

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	AppName       = "Commander"
	CommanderFile = "commander.yaml"
)

func main() {
	tests := []runtime.TestCase{
		{
			Command: runtime.CommandUnderTest{
				Cmd: "echo hello",
			},
			Expected: runtime.Expected{
				Stdout: runtime.ExpectedOut{
					Exactly: "hello",
				},
				ExitCode: 0,
			},
			Title: "Output hello",
		},
	}

	results := runtime.Start(tests)
	r := start(results)

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
