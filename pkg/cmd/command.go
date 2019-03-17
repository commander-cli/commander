package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
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
	}
}

//SetTimeoutMS sets the timeout in milliseconds
func (c *Command) SetTimeoutMS(ms int) {
	if ms == 0 {
		c.Timeout = 1 * time.Minute
	}
	c.Timeout = time.Duration(ms) * time.Millisecond
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

//Execute executes the commande
func (c *Command) Execute() error {
	cmd := exec.Command("sh", "-c", c.Cmd)
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
			log.Println("Command exited with error", c.Cmd, err.Error())
			c.getExitCode(err)
			break
		}
		log.Println("Command exited successfully", c.Cmd)
		c.exitCode = 0
	case <-time.After(c.Timeout):
		log.Println("Command timed out", c.Cmd)
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("Timeout occurred and can not kill process with pid %v", cmd.Process.Pid)
		}
		return fmt.Errorf("Command timed out after %v", c.Timeout)
	}

	// Remove leading and trailing whitespaces
	c.stderr = strings.Trim(errBuff.String(), "\n")
	c.stdout = strings.Trim(outBuff.String(), "\n")
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
