package runtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/commander-cli/commander/pkg/suite"
	"github.com/stretchr/testify/assert"
)

func Test_NewRuntime(t *testing.T) {
	runtime := getRuntime()
	assert.Len(t, runtime.Runner.Nodes, 3)
}

func Test_RuntimeStart(t *testing.T) {
	s := getExampleTestCases()
	r := getRuntime()
	got := r.Start(s)

	assert.IsType(t, Result{}, got)

	count := 0
	for _, r := range got.TestResults {
		assert.Equal(t, "Output hello", r.TestCase.Title)
		assert.True(t, r.ValidationResult.Success)
		count++
	}
	assert.Equal(t, 1, count)
}

func TestRuntime_WithRetries(t *testing.T) {
	s := getExampleTestCases()
	s[0].Command.Retries = 3
	s[0].Command.Cmd = "echo fail"

	r := getRuntime()
	got := r.Start(s)

	var counter = 0
	for _, r := range got.TestResults {
		counter++
		assert.False(t, r.ValidationResult.Success)
		assert.Equal(t, 3, r.Tries)
	}

	assert.Equal(t, 1, counter)
}

func Test_AlphabeticalOrder(t *testing.T) {
	tests := []suite.TestCase{
		{Title: "bbb", Command: suite.CommandUnderTest{Cmd: "exit 0;"}},
		{Title: "aaa"},
		{Title: "111"},
		{Title: "_"},
	}

	got := []string{}
	runtime := NewRuntime(&EventHandler{TestFinished: func(r TestResult) {
		got = append(got, r.TestCase.Title)
	}})

	runtime.Start(tests)

	assert.Equal(t, "111", got[0])
	assert.Equal(t, "_", got[1])
	assert.Equal(t, "aaa", got[2])
	assert.Equal(t, "bbb", got[3])
}

func Test_RuntimeWithRetriesAndInterval(t *testing.T) {
	s := getExampleTestCases()
	s[0].Command.Retries = 3
	s[0].Command.Cmd = "echo fail"
	s[0].Command.Interval = "50ms"

	start := time.Now()
	r := getRuntime()
	got := r.Start(s)

	var counter = 0
	for _, r := range got.TestResults {
		counter++
		assert.False(t, r.ValidationResult.Success)
		assert.Equal(t, 3, r.Tries)
	}
	duration := time.Since(start)

	assert.Equal(t, 1, counter)
	assert.True(t, duration.Seconds() > 0.15, "Retry interval did not work")
}

func Test_RuntimeWithSkip(t *testing.T) {
	s := getExampleTestCases()
	s[0].Skip = true

	r := getRuntime()
	got := r.Start(s)

	assert.Equal(t, 1, got.Skipped)
}

func getRuntime() Runtime {
	eh := EventHandler{
		TestFinished: func(tr TestResult) {
			fmt.Println("I do nothing")
		},
		TestSkipped: func(tr TestResult) {
			fmt.Printf("%s was skipped", tr.TestCase.Title)
		},
	}

	runtime := NewRuntime(&eh, suite.Node{Name: "test"}, suite.Node{Name: "test2"})
	return runtime
}

func getExampleTestCases() []suite.TestCase {
	tests := []suite.TestCase{
		{
			Command: suite.CommandUnderTest{
				Cmd:     "echo hello",
				Timeout: "5s",
			},
			Expected: suite.Expected{
				Stdout: suite.ExpectedOut{
					Exactly: "hello",
				},
				ExitCode: 0,
			},
			Title: "Output hello",
		},
	}
	return tests
}
