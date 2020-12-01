package runtime

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/commander-cli/commander/pkg/suite"
	"github.com/stretchr/testify/assert"
)

func TestRuntime_WithEnvVariables(t *testing.T) {
	envVar := "$KEY"
	if runtime.GOOS == "windows" {
		envVar = "%KEY%"
	}

	s := suite.TestCase{
		Command: suite.CommandUnderTest{
			Cmd:     fmt.Sprintf("echo %s", envVar),
			Timeout: "2s",
			Env:     map[string]string{"KEY": "value"},
		},
		Expected: suite.Expected{
			Stdout: suite.ExpectedOut{
				Contains: []string{"value"},
			},
			ExitCode: 0,
		},
		Title: "Output env variable",
	}

	e := LocalExecutor{}
	got := e.Execute(s)
	assert.True(t, got.ValidationResult.Success)
}

func Test_runTestShouldReturnError(t *testing.T) {
	test := suite.TestCase{
		Command: suite.CommandUnderTest{
			Cmd: "pwd",
			Dir: "/home/invalid",
		},
	}

	e := LocalExecutor{}
	got := e.Execute(test)

	if runtime.GOOS == "windows" {
		assert.Contains(t, got.TestCase.Result.Error.Error(), "chdir /home/invalid")
	} else {
		assert.Equal(t, "chdir /home/invalid: no such file or directory", got.TestCase.Result.Error.Error())
	}
}

func TestRuntime_WithInvalidDuration(t *testing.T) {
	test := suite.TestCase{
		Command: suite.CommandUnderTest{
			Cmd:     "echo test",
			Timeout: "600lightyears",
		},
	}

	e := LocalExecutor{}
	got := e.Execute(test)

	assert.Equal(t, "time: unknown unit lightyears in duration 600lightyears", got.TestCase.Result.Error.Error())
}
