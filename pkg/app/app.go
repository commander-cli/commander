package app

import (
	"github.com/urfave/cli"
	"strings"
)

const (
	//AppName defines the app name
	AppName = "Commander"
	//CommanderFile holds the default config file which is loaded
	CommanderFile = "commander.yaml"
)

//AddCommandContext holds all flags for the add command
type AddCommandContext struct {
	Verbose    bool
	NoColor    bool
	Dir        bool
	Concurrent int
	Filters    []string
}

//NewAddContextFromCli is a constructor which creates the context
func NewAddContextFromCli(c *cli.Context) AddCommandContext {
	filters := strings.Split(c.String("filter"), ",")
	if filters[0] == "" {
		filters = make([]string, 0)
	}

	return AddCommandContext{
		Verbose:    c.Bool("verbose"),
		NoColor:    c.Bool("no-color"),
		Dir:        c.Bool("dir"),
		Concurrent: c.Int("concurrent"),
		Filters:    filters,
	}
}
