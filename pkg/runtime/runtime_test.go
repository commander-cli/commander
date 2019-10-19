package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

const SingleConcurrent = 1

func TestRuntime_Start(t *testing.T) {
	s := getExampleTestSuite()
	got := Start(s, SingleConcurrent)

	assert.IsType(t, make(<-chan TestResult), got)

	count := 0
	for r := range got {
		assert.Equal(t, "Output hello", r.TestCase.Title)
		assert.True(t, r.ValidationResult.Success)
		count++
	}
	assert.Equal(t, 1, count)
}

func TestRuntime_WithRetries(t *testing.T) {
	s := getExampleTestSuite()
	s[0].Command.Retries = 3
	s[0].Command.Cmd = "echo fail"

	got := Start(s, 1)

	var counter = 0
	for r := range got {
		counter++
		assert.False(t, r.ValidationResult.Success)
		assert.Equal(t, 3, r.Tries)
	}

	assert.Equal(t, 1, counter)
}

func TestRuntime_WithRetriesAndInterval(t *testing.T) {
	s := getExampleTestSuite()
	s[0].Command.Retries = 3
	s[0].Command.Cmd = "echo fail"
	s[0].Command.Interval = "50ms"

	start := time.Now()
	got := Start(s, 1)

	var counter = 0
	for r := range got {
		counter++
		assert.False(t, r.ValidationResult.Success)
		assert.Equal(t, 3, r.Tries)
	}
	duration := time.Since(start)

	assert.Equal(t, 1, counter)
	assert.True(t, duration.Seconds() > 0.15, "Retry interval did not work")
}

func TestRuntime_WithEnvVariables(t *testing.T) {
	envVar := "$KEY"
	if runtime.GOOS == "windows" {
		envVar = "%KEY%"
	}

	s := TestCase{
		Command: CommandUnderTest{
			Cmd:     fmt.Sprintf("echo %s", envVar),
			Timeout: "50ms",
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

func TestRuntime_WithInvalidDuration(t *testing.T) {
	test := TestCase{
		Command: CommandUnderTest{
			Cmd:     "echo test",
			Timeout: "600lightyears",
		},
	}

	got := runTest(test)

	assert.Equal(t, "time: unknown unit lightyears in duration 600lightyears", got.TestCase.Result.Error.Error())
}

func getExampleTestSuite() []TestCase {
	tests := []TestCase{
		{
			Command: CommandUnderTest{
				Cmd:     "echo hello",
				Timeout: "50ms",
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
