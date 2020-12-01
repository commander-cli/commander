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

func Execute(eh *EventHandler, suites []suite.Suite, filters Filters) (Result, error) {
	var result Result

	for _, s := range suites {
		tests := s.GetTests()
		if len(filters) != 0 {
			tests = []suite.TestCase{}
		}

		// Filter tests if test filters was given
		for _, f := range filters {
			t, err := s.FindTests(f)
			if err != nil {
				return Result{}, err
			}
			tests = append(tests, t...)
		}

		r := NewRuntime(eh, s.Nodes...)
		newResult := r.Start(tests)

		result = convergeResults(result, newResult)

	}
	return result, nil
}

func convergeResults(result Result, new Result) Result {
	result.TestResults = append(result.TestResults, new.TestResults...)
	result.Failed += new.Failed
	result.Duration += new.Duration

	return result
}

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
