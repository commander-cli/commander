package cmd

import (
	"os/exec"
)

func createBaseCommand(c *Command) *exec.Cmd {
	cmd := exec.Command(`C:\windows\system32\cmd.exe`, "/C", c.Command)
	return cmd
}
