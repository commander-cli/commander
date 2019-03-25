package output

import (
	"bytes"
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
)

func Test_NewCliOutput(t *testing.T) {
	got := NewCliOutput()
	assert.IsType(t, OutputWriter{}, got)
}

func Test_Start(t *testing.T) {
	var buf bytes.Buffer
	var wg sync.WaitGroup
	results := make(chan runtime.TestResult)

	wg.Add(1)
	go func() {
		defer wg.Done()

		writer := OutputWriter{out: &buf}
		got := writer.Start(results)

		assert.False(t, got)
	}()

	results <- runtime.TestResult{
		TestCase: runtime.TestCase{
			Title:   "Successful test",
			Command: runtime.CommandUnderTest{},
		},
		ValidationResult: runtime.ValidationResult{
			Success: true,
		},
		FailedProperty: "",
	}

	results <- runtime.TestResult{
		TestCase: runtime.TestCase{
			Title:   "Failed test",
			Command: runtime.CommandUnderTest{},
		},
		ValidationResult: runtime.ValidationResult{
			Success: false,
		},
		FailedProperty: "Stdout",
	}

	results <- runtime.TestResult{
		TestCase: runtime.TestCase{
			Title:   "Failed test on stderr",
			Command: runtime.CommandUnderTest{},
		},
		ValidationResult: runtime.ValidationResult{
			Success: false,
		},
		FailedProperty: "Stderr",
	}

	results <- runtime.TestResult{
		TestCase: runtime.TestCase{
			Title: "Invalid command",
			Command: runtime.CommandUnderTest{
				Cmd: "some stupid config",
			},
			Result: runtime.CommandResult{
				Error: fmt.Errorf("Some error message"),
			},
		},
	}

	close(results)
	wg.Wait()

	assert.True(t, true)
	output := buf.String()
	assert.True(t, strings.Contains(output, "✓ Successful test"))
	assert.True(t, strings.Contains(output, "✗ Failed test"))
	assert.True(t, strings.Contains(output, "✗ 'Invalid command' could not be executed"))
	assert.True(t, strings.Contains(output, "✗ 'Failed test on stderr', on property 'Stderr'"))
	assert.True(t, strings.Contains(output, "Some error message"))
}

func Test_SuccessSuite(t *testing.T) {
	var buf bytes.Buffer
	var wg sync.WaitGroup
	results := make(chan runtime.TestResult)

	wg.Add(1)
	go func() {
		defer wg.Done()

		writer := OutputWriter{out: &buf}
		got := writer.Start(results)

		assert.True(t, got)
	}()

	results <- runtime.TestResult{
		TestCase: runtime.TestCase{
			Title:   "Successful test",
			Command: runtime.CommandUnderTest{},
		},
		ValidationResult: runtime.ValidationResult{
			Success: true,
		},
		FailedProperty: "",
	}

	close(results)
	wg.Wait()

	assert.True(t, true)
	assert.True(t, strings.Contains(buf.String(), "✓ Successful test"))
	assert.True(t, strings.Contains(buf.String(), "Duration"))
	assert.True(t, strings.Contains(buf.String(), "Count: 1, Failed: 0"))
}
