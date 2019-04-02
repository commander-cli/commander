package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

//Command represents a single command which can be executed
type Command struct {
	Cmd      string
	Env      []string
	Dir      string
	Timeout  time.Duration
	executed bool
	stderr   string
	stdout   string
	exitCode int
}

//NewCommand creates a new command
func NewCommand(cmd string) *Command {
	return &Command{
		Cmd:      cmd,
		Timeout:  1 * time.Minute,
		executed: false,
		Env:      []string{},
	}
}

// AddEnv adds an environment variable to the command
func (c *Command) AddEnv(key string, value string) {
	c.Env = append(c.Env, fmt.Sprintf("%s=%s", key, value))
}

//SetTimeoutMS sets the timeout in milliseconds
func (c *Command) SetTimeoutMS(ms int) {
	if ms == 0 {
		c.Timeout = 1 * time.Minute
		return
	}
	c.Timeout = time.Duration(ms) * time.Millisecond
}

// SetTimeout sets the timeout given a time unit
// Example: SetTimeout("100s") sets the timeout to 100 seconds
func (c *Command) SetTimeout(timeout string) error {
	d, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}

	c.Timeout = d
	return nil
}

//Stdout returns the output to stdout
func (c *Command) Stdout() string {
	c.isExecuted("Stdout")
	return c.stdout
}

//Stderr returns the output to stderr
func (c *Command) Stderr() string {
	c.isExecuted("Stderr")
	return c.stderr
}

//ExitCode returns the exit code of the command
func (c *Command) ExitCode() int {
	c.isExecuted("ExitCode")
	return c.exitCode
}

//Executed returns if the command was already executed
func (c *Command) Executed() bool {
	return c.executed
}

func (c *Command) isExecuted(property string) {
	if !c.executed {
		panic("Can not read " + property + " if command was not executed.")
	}
}

// Execute executes the command and writes the results into it's own instance
// The results can be received with the Stdout(), Stderr() and ExitCode() methods
func (c *Command) Execute() error {
	cmd := createBaseCommand(c)
	cmd.Env = c.Env
	cmd.Dir = c.Dir

	var (
		outBuff bytes.Buffer
		errBuff bytes.Buffer
	)
	cmd.Stdout = &outBuff
	cmd.Stderr = &errBuff

	err := cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			c.getExitCode(err)
			break
		}
		c.exitCode = 0
	case <-time.After(c.Timeout):
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("Timeout occurred and can not kill process with pid %v", cmd.Process.Pid)
		}
		return fmt.Errorf("Command timed out after %v", c.Timeout)
	}

	// Remove leading and trailing whitespaces
	c.stderr = c.removeLineBreaks(errBuff.String())
	c.stdout = c.removeLineBreaks(outBuff.String())
	c.executed = true

	return nil
}

func (c *Command) getExitCode(err error) {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			c.exitCode = status.ExitStatus()
		}
	}
}
