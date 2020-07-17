package app

import (
	"github.com/urfave/cli"
)

const (
	//AppName defines the app name
	AppName = "Commander"
	//CommanderFile holds the default config file which is loaded
	CommanderFile = "commander.yaml"
)

//TestCommandContext holds all flags for the add command
type TestCommandContext struct {
	Verbose    bool
	NoColor    bool
	Dir        bool
	Concurrent int
	Filters    []string
}

//NewTestContextFromCli is a constructor which creates the context
func NewTestContextFromCli(c *cli.Context) TestCommandContext {
	return TestCommandContext{
		Verbose:    c.Bool("verbose"),
		NoColor:    c.Bool("no-color"),
		Dir:        c.Bool("dir"),
		Concurrent: c.Int("concurrent"),
		Filters:    c.StringSlice("filter"),
	}
}
