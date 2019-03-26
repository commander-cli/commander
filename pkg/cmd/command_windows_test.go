package cmd

import (
	"github.com/stretchr/testify/assert"
	"strings"
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
	// This is needed because windows sometimes can not kill the process :(
	containsMsg := strings.Contains(err.Error(), "Timeout occurred and can not kill process with pid") || strings.Contains(err.Error(), "Command timed out after 5ms")
	assert.True(t, containsMsg)
}

func TestCommand_WithValidTimeout(t *testing.T) {
	cmd := NewCommand("timeout 0.01;")
	cmd.SetTimeoutMS(1000)

	err := cmd.Execute()

	assert.Nil(t, err)
}
