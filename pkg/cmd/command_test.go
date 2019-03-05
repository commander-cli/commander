package cmd

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestCommand_NewCommand(t *testing.T) {
    cmd := NewCommand("")
    assert.False(t, cmd.Executed())
}

func TestCommand_Execute(t *testing.T) {
    cmd := NewCommand("/bin/echo hello")

    err := cmd.Execute()

    assert.Nil(t, err)
    assert.True(t, cmd.Executed())
    assert.Equal(t, cmd.Stdout(), "hello")
}

func TestCommand_ExecuteStderr(t *testing.T) {
    cmd := NewCommand(">&2 /bin/echo hello")

    err := cmd.Execute()

    assert.Nil(t, err)
    assert.Equal(t, "hello", cmd.Stderr())
}

func TestCommand_ExitCode(t *testing.T) {
    cmd := NewCommand("exit 120")

    err := cmd.Execute()

    assert.Nil(t, err)
    assert.Equal(t, 120, cmd.ExitCode())
}