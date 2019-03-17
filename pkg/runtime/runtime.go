package runtime

import (
	"github.com/SimonBaeumer/commander/pkg/cmd"
	"log"
	"sync"
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
	Success ResultStatus = iota
	Failed
	Skipped
)

// TestCase represents a test case which will be executed by the runtime
type TestCase struct {
	Title    string
	Command  CommandUnderTest
	Expected Expected
	Result   CommandResult
}

//Config
type TestConfig struct {
	Env []string
	Dir string
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

type ExpectedOut struct {
	Contains  []string
	Lines     map[int]string
	Exactly   string
	LineCount int
}

// CommandUnderTest represents the command under test
type CommandUnderTest struct {
	Cmd string
	Env []string
	Dir string
}

// CommandResult represents the TestCase and the ValidationResult
type TestResult struct {
	TestCase         TestCase
	ValidationResult ValidationResult
	FailedProperty   string
}

// Start starts the given test suite
func Start(tests []TestCase) <-chan TestResult {
	in := make(chan TestCase)
	out := make(chan TestResult)

	go func(tests []TestCase) {
		defer close(in)
		for _, t := range tests {
			in <- t
		}
	}(tests)

	//TODO: Add more concurrency
	workerCount := 1
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(tests chan TestCase) {
			defer wg.Done()
			for t := range tests {
				out <- runTest(t)
			}
		}(in)
	}

	go func(results chan TestResult) {
		wg.Wait()
		close(results)
	}(out)

	return out
}

func runTest(test TestCase) TestResult {
	// cut = command under test
	cut := cmd.NewCommand(test.Command.Cmd)
	cut.Env = test.Command.Env
	cut.Dir = test.Command.Dir

	if err := cut.Execute(); err != nil {
		log.Println("Command failed ", err.Error())
		test.Result = CommandResult{
			Error: err,
		}

		return TestResult{
			TestCase: test,
		}
	}

	log.Println("Executed command ", test.Command.Cmd)

	// Write test result
	test.Result = CommandResult{
		ExitCode: cut.ExitCode(),
		Stdout:   cut.Stdout(),
		Stderr:   cut.Stderr(),
	}

	return Validate(test)
}
