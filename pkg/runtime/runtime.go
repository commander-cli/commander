package runtime

import (
	"time"

	"github.com/SimonBaeumer/commander/pkg/output"
)

// Constants for defining the various tested properties
const (
	ExitCode  = "ExitCode"
	Stdout    = "Stdout"
	Stderr    = "Stderr"
	LineCount = "LineCount"
)

// Result status codes
const (
	//Success status
	Success ResultStatus = iota
	// Failed status
	Failed
	// Skipped status
	Skipped
)

// NewRuntime creates a new runtime and inits default nodes
func NewRuntime(out *output.OutputWriter, nodes ...Node) Runtime {
	local := Node{
		Name: "local",
		Type: "local",
		Addr: "localhost",
	}

	nodes = append(nodes, local)
	runner := Runner{
		Nodes: nodes,
	}

	return Runtime{
		Runner: &runner,
		Output: out,
	}
}

// Runtime represents the current runtime, please use NewRuntime() instead of creating an instance directly
type Runtime struct {
	Runner *Runner
	Output *output.OutputWriter
}

// TestCase represents a test case which will be executed by the runtime
type TestCase struct {
	Title    string
	Command  CommandUnderTest
	Expected Expected
	Result   CommandResult
	Nodes    []string
}

//GlobalTestConfig represents the configuration for a test
type GlobalTestConfig struct {
	Env        map[string]string
	Dir        string
	Timeout    string
	Retries    int
	Interval   string
	InheritEnv bool
	Nodes      []string
}

// ResultStatus represents the status code of a test result
type ResultStatus int

// CommandResult holds the result for a specific test
type CommandResult struct {
	Status            ResultStatus
	Stdout            string
	Stderr            string
	ExitCode          int
	FailureProperties []string
	Error             error
}

//Expected is the expected output of the command under test
type Expected struct {
	Stdout    ExpectedOut
	Stderr    ExpectedOut
	LineCount int
	ExitCode  int
}

//ExpectedOut represents the assertions on stdout and stderr
type ExpectedOut struct {
	Contains    []string          `yaml:"contains,omitempty"`
	Lines       map[int]string    `yaml:"lines,omitempty"`
	Exactly     string            `yaml:"exactly,omitempty"`
	LineCount   int               `yaml:"line-count,omitempty"`
	NotContains []string          `yaml:"not-contains,omitempty"`
	JSON        map[string]string `yaml:"json,omitempty"`
	XML         map[string]string `yaml:"xml,omitempty"`
}

// CommandUnderTest represents the command under test
type CommandUnderTest struct {
	Cmd        string
	InheritEnv bool
	Env        map[string]string
	Dir        string
	Timeout    string
	Retries    int
	Interval   string
}

// TestResult represents the TestCase and the ValidationResult
type TestResult struct {
	TestCase         TestCase
	ValidationResult ValidationResult
	FailedProperty   string
	Tries            int
	Node             string
	FileName         string
}

// Start starts the given test suite and executes all tests
func (r *Runtime) Start(tests []TestCase) output.Result {
	result := output.Result{}
	testCh := r.Runner.Execute(tests)
	start := time.Now()
	for tr := range testCh {
		tr := convertTestResult(tr)
		r.Output.PrintResult(tr)

		if !tr.Success {
			result.Failed++
		}

		result.TestResults = append(result.TestResults, tr)
	}
	result.Duration = time.Since(start)

	return result
}

// convert runtime.TestResult to output.TestResult
func convertTestResult(tr TestResult) output.TestResult {
	testResult := output.TestResult{
		FileName:       "", //TODO: Get filename from TestßCase
		Title:          tr.TestCase.Title,
		Node:           tr.Node,
		Tries:          tr.Tries,
		Success:        tr.ValidationResult.Success,
		FailedProperty: tr.FailedProperty,
		Diff:           tr.ValidationResult.Diff,
		Error:          tr.TestCase.Result.Error,
	}

	return testResult
}

// GetRetries returns the retries of the command
func (c *CommandUnderTest) GetRetries() int {
	if c.Retries == 0 {
		return 1
	}
	return c.Retries
}
