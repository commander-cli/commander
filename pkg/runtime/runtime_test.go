package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestRuntime_Start(t *testing.T) {
	s := getExampleTestSuite()
	got := Start(s)

	assert.IsType(t, make(<-chan TestResult), got)

	count := 0
	for r := range got {
		assert.Equal(t, "Output hello", r.TestCase.Title)
		assert.True(t, r.ValidationResult.Success)
		count++
	}
	assert.Equal(t, 1, count)
}

func TestRuntime_WithEnvVariables(t *testing.T) {
	envVar := "$KEY"
	if runtime.GOOS == "windows" {
		envVar = "%KEY%"
	}

	s := TestCase{
		Command: CommandUnderTest{
			Cmd:     fmt.Sprintf("echo %s", envVar),
			Timeout: 50,
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

	got := runTest(s)
	assert.True(t, got.ValidationResult.Success)
}

func Test_runTestShouldReturnError(t *testing.T) {
	test := TestCase{
		Command: CommandUnderTest{
			Cmd: "pwd",
			Dir: "/home/invalid",
		},
	}

	got := runTest(test)

	if runtime.GOOS == "windows" {
		assert.Contains(t, got.TestCase.Result.Error.Error(), "chdir /home/invalid: The system cannot find the path specified.")
	} else {
		assert.Equal(t, "chdir /home/invalid: no such file or directory", got.TestCase.Result.Error.Error())
	}
}

func getExampleTestSuite() []TestCase {
	tests := []TestCase{
		{
			Command: CommandUnderTest{
				Cmd:     "echo hello",
				Timeout: 50,
			},
			Expected: Expected{
				Stdout: ExpectedOut{
					Exactly: "hello",
				},
				ExitCode: 0,
			},
			Title: "Output hello",
		},
	}
	return tests
}
