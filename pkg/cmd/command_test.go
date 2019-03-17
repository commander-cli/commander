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

func TestCommand_WithEnvVariables(t *testing.T) {
    cmd := NewCommand("echo $TEST")
    cmd.Env = []string{"TEST=hey"}

    _ = cmd.Execute()

    assert.Equal(t, "hey", cmd.Stdout())
}

func TestCommand_WithTimeout(t *testing.T) {
    cmd := NewCommand("sleep 0.005;")
    cmd.SetTimeoutMS(5)

    err := cmd.Execute()

    assert.NotNil(t, err)
    assert.Equal(t, "Command timed out after 5ms", err.Error())
}

func TestCommand_WithValidTimeout(t *testing.T) {
    cmd := NewCommand("sleep 0.01;")
    cmd.SetTimeoutMS(500)

    err := cmd.Execute()

    assert.Nil(t, err)
}
