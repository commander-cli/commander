package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommand_ExecuteStderr(t *testing.T) {
	cmd := NewCommand("echo hello 1>&2")

	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "hello ", cmd.Stderr())
}

func TestCommand_WithTimeout(t *testing.T) {
	cmd := NewCommand("timeout 0.005;")
	cmd.SetTimeoutMS(5)

	err := cmd.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, "Command timed out after 5ms", err.Error())
}

func TestCommand_WithValidTimeout(t *testing.T) {
	cmd := NewCommand("timeout 0.01;")
	cmd.SetTimeoutMS(500)

	err := cmd.Execute()

	assert.Nil(t, err)
}
