package output

import (
	"errors"
	"fmt"
	"github.com/commander-cli/commander/pkg/runtime"
)

const (
	TAP = "tap"
	CLI = "cli"
)

type Output interface {
	GetEventHandler() *runtime.EventHandler
	PrintSummary(result runtime.Result)
}

// NewOutput creates a new output
func NewOutput(format string, color bool) (Output, error) {
	switch format {
	case TAP:
		return NewTAPOutputWriter(), nil
	case CLI:
		return NewCliOutput(color), nil
	}
	return nil, errors.New(fmt.Sprintf("Invalid format type %s", format))
}
