package cmd

import (
	"os/exec"
	"strings"
)

func createBaseCommand(c *Command) *exec.Cmd {
	cmd := exec.Command("sh", "-c", c.Cmd)
	return cmd
}

func (c *Command) removeLineBreaks(text string) string {
	return strings.Trim(text, "\n")
}
