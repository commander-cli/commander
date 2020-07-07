package runtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewRuntime(t *testing.T) {
	runtime := getRuntime()
	assert.Len(t, runtime.Runner.Nodes, 3)
}

func Test_RuntimeStart(t *testing.T) {
	s := getExampleTestSuite()
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
	s := getExampleTestSuite()
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

func Test_RuntimeWithRetriesAndInterval(t *testing.T) {
	s := getExampleTestSuite()
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

func getRuntime() Runtime {
	eh := EventHandler{
		TestFinsihed: func(tr TestResult) {
			fmt.Println("I do nothing")
		},
	}

	runtime := NewRuntime(&eh, Node{Name: "test"}, Node{Name: "test2"})
	return runtime
}

func getExampleTestSuite() []TestCase {
	tests := []TestCase{
		{
			Command: CommandUnderTest{
				Cmd:     "echo hello",
				Timeout: "5s",
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
