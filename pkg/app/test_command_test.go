package app

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func Test_TestCommand(t *testing.T) {
	err := TestCommand("commander.yaml", "", AddCommandContext{})

	if runtime.GOOS == "windows" {
		assert.Equal(t, "Error open commander.yaml: The system cannot find the file specified.", err.Error())
	} else {
		assert.Equal(t, "Error open commander.yaml: no such file or directory", err.Error())
	}
}

func Test_TestCommand_ShouldUseCustomFile(t *testing.T) {
	err := TestCommand("my-test.yaml", "", AddCommandContext{})

	if runtime.GOOS == "windows" {
		assert.Equal(t, "Error open my-test.yaml: The system cannot find the file specified.", err.Error())
	} else {
		assert.Equal(t, "Error open my-test.yaml: no such file or directory", err.Error())
	}
}
