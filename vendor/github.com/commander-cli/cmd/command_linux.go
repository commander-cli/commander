package cmd

import (
	"os/exec"
	"syscall"
)

func createBaseCommand(c *Command) *exec.Cmd {
	cmd := exec.Command("/bin/sh", "-c", c.Command)
	return cmd
}

// WithUser allows the command to be run as a different
// user.
//
// Example:
//
//	cred := syscall.Credential{Uid: 1000, Gid: 1000}
//	c := NewCommand("echo hello", cred)
//	c.Execute()
func WithUser(credential syscall.Credential) func(c *Command) {
	return func(c *Command) {
		c.baseCommand.SysProcAttr = &syscall.SysProcAttr{
			Credential: &credential,
		}
	}
}
