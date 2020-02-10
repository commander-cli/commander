package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const SingleConcurrent = 1

func TestRuntime_Start(t *testing.T) {
	s := getExampleTestSuite()
	r := Runtime{}
	got := r.Start(s, SingleConcurrent)

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

	r := Runtime{}
	got := r.Start(s, 1)

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
	r := Runtime{}
	got := r.Start(s, 0)

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
