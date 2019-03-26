package main

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func Test_CreateCliApp(t *testing.T) {
	app := createCliApp()

	assert.Equal(t, "Commander", app.Name)
	assert.Equal(t, "test", app.Commands[0].Name)
}

func Test_TestCommand(t *testing.T) {
	err := testCommand("commander.yaml", "", CommanderContext{})

	if runtime.GOOS == "windows" {
		assert.Equal(t, "Error open commander.yaml: The system cannot find the file specified.", err.Error())
	} else {
		assert.Equal(t, "Error open commander.yaml: no such file or directory", err.Error())
	}
}

func Test_TestCommand_ShouldUseCustomFile(t *testing.T) {
	err := testCommand("my-test.yaml", "", CommanderContext{})

	if runtime.GOOS == "windows" {
		assert.Equal(t, "Error open my-test.yaml: The system cannot find the file specified.", err.Error())
	} else {
		assert.Equal(t, "Error open my-test.yaml: no such file or directory", err.Error())
	}
}
