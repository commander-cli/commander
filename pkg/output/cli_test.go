package output

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/stretchr/testify/assert"
)

func Test_NewCliOutput(t *testing.T) {
	got := NewCliOutput(true, true)
	assert.IsType(t, OutputWriter{}, got)
}

func Test_Start(t *testing.T) {
	var buf bytes.Buffer
	var wg sync.WaitGroup
	results := make(chan runtime.TestResult)

	wg.Add(1)
	go func() {
		defer wg.Done()

		writer := OutputWriter{out: &buf, order: true}
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
		Node:           "docker-host",
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
		Node:           "192.168.0.1",
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
		Node:           "ssh-host1",
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
		Node: "local",
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
		Node:  "local",
		Tries: 2,
	}

	results <- runtime.TestResult{
		FileName:  "MySweetFile",
		FileError: fmt.Errorf("Some file error message"),
	}

	close(results)
	wg.Wait()

	assert.True(t, true)
	output := buf.String()
	assert.Contains(t, output, "✓ [docker-host] Successful test")
	assert.Contains(t, output, "✗ [192.168.0.1] Failed test")
	assert.Contains(t, output, "✗ [local] 'Invalid command' could not be executed with error message")
	assert.Contains(t, output, "✗ [ssh-host1] 'Failed test on stderr', on property 'Stderr'")
	assert.Contains(t, output, "Some error message")
	assert.Contains(t, output, "Some file error message")
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
		Node:           "local",
	}

	close(results)
	wg.Wait()

	assert.True(t, true)
	assert.True(t, strings.Contains(buf.String(), "✓ [local] Successful test"))
	assert.True(t, strings.Contains(buf.String(), "Duration"))
	assert.True(t, strings.Contains(buf.String(), "Count: 1, Failed: 0"))
}
