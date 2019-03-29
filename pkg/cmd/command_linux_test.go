package cmd

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCommand_ExecuteStderr(t *testing.T) {
	cmd := NewCommand(">&2 echo hello")

	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "hello", cmd.Stderr())
}

func TestCommand_WithTimeout(t *testing.T) {
	cmd := NewCommand("sleep 0.01;")
	cmd.SetTimeoutMS(1)

	err := cmd.Execute()

	assert.NotNil(t, err)
	// Sadly a process can not be killed every time :(
	containsMsg := strings.Contains(err.Error(), "Timeout occurred and can not kill process with pid") || strings.Contains(err.Error(), "Command timed out after 1ms")
	assert.True(t, containsMsg)
}

func TestCommand_WithValidTimeout(t *testing.T) {
	cmd := NewCommand("sleep 0.01;")
	cmd.SetTimeoutMS(500)

	err := cmd.Execute()

	assert.Nil(t, err)
}
