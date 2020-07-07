package output

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/stretchr/testify/assert"
)

func Test_NewCliOutput(t *testing.T) {
	got := NewCliOutput(true)
	assert.IsType(t, OutputWriter{}, got)
}

func Test_GetEventHandler(t *testing.T) {
	writer := NewCliOutput(true)
	eh := writer.GetEventHandler()
	assert.IsType(t, &runtime.EventHandler{}, eh)
}

func Test_EventHandlerTestFinished(t *testing.T) {
	var buf bytes.Buffer
	writer := NewCliOutput(true)
	writer.out = &buf
	eh := writer.GetEventHandler()

	testResults := createFakeTestResults()

	for _, tr := range testResults {
		eh.TestFinsihed(tr)

	}

	output := buf.String()

	assert.Contains(t, output, "✗ [192.168.0.1] Failed test")
	assert.Contains(t, output, "✓ [docker-host] Successful test")

}

func Test_PrintSummary(t *testing.T) {
	r := runtime.Result{
		Duration:    10,
		Failed:      1,
		TestResults: createFakeTestResults(),
	}

	var buf bytes.Buffer
	writer := NewCliOutput(true)
	writer.out = &buf

	outResult := writer.PrintSummary(r)
	assert.False(t, outResult)

	output := buf.String()
	assert.Contains(t, output, "✗ [192.168.0.1] 'Failed test', on property Stdout")
	assert.NotContains(t, output, "✓ [docker-host] Successful test")
}

func createFakeTestResults() []runtime.TestResult {
	tr := runtime.TestResult{
		TestCase: runtime.TestCase{
			Title:   "Failed test",
			Command: runtime.CommandUnderTest{},
		},
		ValidationResult: runtime.ValidationResult{
			Success: false,
		},
		FailedProperty: "Stdout",
		Node:           "192.168.0.1",
	}

	tr2 := runtime.TestResult{
		TestCase: runtime.TestCase{
			Title:   "Successful test",
			Command: runtime.CommandUnderTest{},
		},
		ValidationResult: runtime.ValidationResult{
			Success: true,
		},
		FailedProperty: "",
		Node:           "docker-host",
	}

	tr3 := runtime.TestResult{
		TestCase: runtime.TestCase{
			Title: "Invalid command",
			Command: runtime.CommandUnderTest{
				Cmd: "some stupid config",
			},
			Result: runtime.CommandResult{
				Error: fmt.Errorf("Some error message"),
			},
		},
		Node:  "local",
		Tries: 2,
	}

	return []runtime.TestResult{tr, tr2, tr3}
}
