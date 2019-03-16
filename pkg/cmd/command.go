package cmd

import (
    "bytes"
    "os/exec"
    "strings"
    "syscall"
)

//Command represents a single command which can be executed
type Command struct {
    Cmd      string
    Env      []string
    Dir      string
    executed bool
    stderr   string
    stdout   string
    exitCode int
}

//NewCommand creates a new command
func NewCommand(cmd string) *Command {
    return &Command{
        Cmd:      cmd,
        executed: false,
    }
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

//Execute executes the command
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

    c.exitCode = 0
    if err := cmd.Wait(); err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
                c.exitCode = status.ExitStatus()
            }
        }
    }

    // Remove leading and trailing whitespaces
    c.stderr = strings.Trim(errBuff.String(), "\n")
    c.stdout = strings.Trim(outBuff.String(), "\n")
    c.executed = true

    return nil
}
