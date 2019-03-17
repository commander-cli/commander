package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateCliApp(t *testing.T) {
	app := createCliApp()

	assert.Equal(t, "Commander", app.Name)
	assert.Equal(t, "test", app.Commands[0].Name)
}

func Test_TestCommand(t *testing.T) {
	err := testCommand("commander.yaml", "")

	assert.Equal(t, "Error open commander.yaml: no such file or directory", err.Error())
}

func Test_TestCommand_ShouldUseCustomFile(t *testing.T) {
	err := testCommand("my-test.yaml", "")

	assert.Equal(t, "Error open my-test.yaml: no such file or directory", err.Error())
}
