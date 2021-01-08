package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestRuntime_WithEnvVariables(t *testing.T) {
	envVar := "$KEY"
	if runtime.GOOS == "windows" {
		envVar = "%KEY%"
	}

	s := TestCase{
		Command: CommandUnderTest{
			Cmd:     fmt.Sprintf("echo %s", envVar),
			Timeout: "2s",
			Env:     map[string]string{"KEY": "value"},
		},
		Expected: Expected{
			Stdout: ExpectedOut{
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
	test := TestCase{
		Command: CommandUnderTest{
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
	test := TestCase{
		Command: CommandUnderTest{
			Cmd:     "echo test",
			Timeout: "600lightyears",
		},
	}

	e := LocalExecutor{}
	got := e.Execute(test)

	assert.Equal(t, `time: unknown unit "lightyears" in duration "600lightyears"`, got.TestCase.Result.Error.Error())
}
