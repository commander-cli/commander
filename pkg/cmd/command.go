package cmd

import (
    "bytes"
    "os/exec"
    "strings"
    "syscall"
)

type Command struct {
    Cmd      string
    Env      []string
    Dir      string
    executed bool
    stderr   string
    stdout   string
    exitCode int
}

func NewCommand(cmd string) *Command {
    return &Command{
        Cmd: cmd,
        executed: false,
    }
}

func (c *Command) Stdout() string {
    c.evaluateExecuted("Stdout")
    return c.stdout
}

func (c *Command) Stderr() string {
    c.evaluateExecuted("Stderr")
    return c.stderr
}

func (c *Command) ExitCode() int {
    c.evaluateExecuted("ExitCode")
    return c.exitCode
}

func (c *Command) Executed() bool {
    return c.executed
}

func (c *Command) evaluateExecuted(property string) {
    if !c.executed {
        panic("Can not read " + property + " if command was not executed.")
    }
}

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
