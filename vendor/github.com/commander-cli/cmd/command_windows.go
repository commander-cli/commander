package cmd

import (
	"os/exec"
	"syscall"
)

func createBaseCommand(c *Command) *exec.Cmd {
	cmd := exec.Command(`C:\windows\system32\cmd.exe`, "/C", c.Command)
	return cmd
}

// WithUser allows the command to be run as a different
// user.
//
// Example:
//
//	token := syscall.Token(handle)
//	c := NewCommand("echo hello", token)
//	c.Execute()
func WithUser(token syscall.Token) func(c *Command) {
	return func(c *Command) {
		c.baseCommand.SysProcAttr = &syscall.SysProcAttr{
			Token: token,
		}
	}
}
