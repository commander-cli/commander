package runtime

import (
	"log"
	"sort"
	"time"
)

// Constants for defining the various tested properties
const (
	ExitCode  = "ExitCode"
	Stdout    = "Stdout"
	Stderr    = "Stderr"
	LineCount = "LineCount"
)

type Filters []string

// NewRuntime creates a new runtime and inits default nodes
func NewRuntime(eh *EventHandler, nodes ...Node) Runtime {
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
		Runner:       &runner,
		EventHandler: eh,
	}
}

// Runtime represents the current runtime, please use NewRuntime() instead of creating an instance directly
type Runtime struct {
	Runner       *Runner
	EventHandler *EventHandler
}

// EventHandler is a configurable event system that handles events such as test completion
type EventHandler struct {
	TestFinished func(TestResult)
	TestSkipped  func(TestResult)
}

// TestCase represents a test case which will be executed by the runtime
type TestCase struct {
	Title    string
	Command  CommandUnderTest
	Expected Expected
	Result   CommandResult
	Nodes    []string
	FileName string
	Skip     bool
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
	File        string            `yaml:"file,omitempty"`
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
	Skipped          bool
}

// Result respresents the aggregation of all TestResults/summary of a runtime
type Result struct {
	TestResults []TestResult
	Duration    time.Duration
	Failed      int
	Skipped     int
}

// Start starts the given test suite and executes all tests
func (r *Runtime) Start(tests []TestCase) Result {
	// Sort tests alphabetically to preserve a reproducible execution order
	sort.SliceStable(tests, func(i, j int) bool {
		return tests[i].Title < tests[j].Title
	})

	result := Result{}
	testCh := r.Runner.Run(tests)
	start := time.Now()
	for tr := range testCh {
		if tr.Skipped {
			result.Skipped++

			log.Println("title: '"+tr.TestCase.Title+"'", " was skipped")
			log.Println("title: '"+tr.TestCase.Title+"'", " Command: ", tr.TestCase.Command.Cmd)
			log.Println("title: '"+tr.TestCase.Title+"'", " Directory: ", tr.TestCase.Command.Dir)
			log.Println("title: '"+tr.TestCase.Title+"'", " Env: ", tr.TestCase.Command.Env)

			r.EventHandler.TestSkipped(tr)
			result.TestResults = append(result.TestResults, tr)
			continue
		}

		if !tr.ValidationResult.Success {
			result.Failed++
		}

		r.EventHandler.TestFinished(tr)
		result.TestResults = append(result.TestResults, tr)
	}
	result.Duration = time.Since(start)

	return result
}
