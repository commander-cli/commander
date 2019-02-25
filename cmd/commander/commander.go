package main

import (
	"github.com/SimonBaeumer/commander/pkg"
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
				log.Println("Starting test suite")
				file := c.Args().First()
				if file == "" {
					file = CommanderFile
				}

				suite := commander.ParseYAMLFile(file)
				runtime.Start(suite)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
