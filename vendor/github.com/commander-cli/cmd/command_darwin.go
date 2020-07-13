package cmd

import (
	"os/exec"
)

func createBaseCommand(c *Command) *exec.Cmd {
	cmd := exec.Command("/bin/sh", "-c", c.Command)
	return cmd
}
