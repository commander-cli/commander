package runtime

import (
	"log"
	"sort"
	"time"

	"github.com/commander-cli/commander/pkg/suite"
)

// Constants for defining the various tested properties
const (
	ExitCode  = "ExitCode"
	Stdout    = "Stdout"
	Stderr    = "Stderr"
	LineCount = "LineCount"
)

// Filters represent runtime filters
type Filters []string

// NewRuntime creates a new runtime and inits default nodes
func NewRuntime(eh *EventHandler, nodes ...suite.Node) Runtime {
	local := suite.Node{
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

// TestResult represents the TestCase and the ValidationResult
type TestResult struct {
	TestCase         suite.TestCase
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
func (r *Runtime) Start(tests []suite.TestCase) Result {
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
