package runtime

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/cmd"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Constants for defining the various tested properties
const (
	ExitCode  = "ExitCode"
	Stdout    = "Stdout"
	Stderr    = "Stderr"
	LineCount = "LineCount"
)

const WorkerCountMultiplicator = 5

// Result status codes
const (
	//Success status
	Success ResultStatus = iota
	// Failed status
	Failed
	// Skipped status
	Skipped
)

// TestCase represents a test case which will be executed by the runtime
type TestCase struct {
	Title    string
	Command  CommandUnderTest
	Expected Expected
	Result   CommandResult
}

//TestConfig represents the configuration for a test
type TestConfig struct {
	Env      map[string]string
	Dir      string
	Timeout  string
	Retries  int
	Interval string
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
	Contains    []string
	Lines       map[int]string
	Exactly     string
	LineCount   int
	NotContains []string
}

// CommandUnderTest represents the command under test
type CommandUnderTest struct {
	Cmd      string
	Env      map[string]string
	Dir      string
	Timeout  string
	Retries  int
	Interval string
}

// TestResult represents the TestCase and the ValidationResult
type TestResult struct {
	TestCase         TestCase
	ValidationResult ValidationResult
	FailedProperty   string
	Tries            int
}

// Start starts the given test suite and executes all tests
// maxConcurrent configures the amount of go routines which will be started
func Start(tests []TestCase, maxConcurrent int) <-chan TestResult {
	in := make(chan TestCase)
	out := make(chan TestResult)

	go func(tests []TestCase) {
		defer close(in)
		for _, t := range tests {
			in <- t
		}
	}(tests)

	workerCount := maxConcurrent
	if maxConcurrent == 0 {
		workerCount = runtime.NumCPU() * WorkerCountMultiplicator
	}

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(tests chan TestCase) {
			defer wg.Done()
			for t := range tests {
				result := TestResult{}
				for i := 1; i <= t.Command.GetRetries(); i++ {
					result = runTest(t)
					result.Tries = i
					if result.ValidationResult.Success {
						break
					}

					executeRetryInterval(t)
				}

				out <- result
			}
		}(in)
	}

	go func(results chan TestResult) {
		wg.Wait()
		close(results)
	}(out)

	return out
}

func executeRetryInterval(t TestCase) {
	if t.Command.GetRetries() > 1 && t.Command.Interval != "" {
		interval, err := time.ParseDuration(t.Command.Interval)
		if err != nil {
			panic(fmt.Sprintf("'%s' interval error: %s", t.Command.Cmd, err))
		}
		time.Sleep(interval)
	}
}

// runTest executes the current test case
func runTest(test TestCase) TestResult {
	// cut = command under test
	cut := cmd.NewCommand(test.Command.Cmd)
	cut.SetTimeout(test.Command.Timeout)
	cut.Dir = test.Command.Dir
	for k, v := range test.Command.Env {
		cut.AddEnv(k, v)
	}

	if err := cut.Execute(); err != nil {
		log.Println(test.Title, " failed ", err.Error())
		test.Result = CommandResult{
			Error: err,
		}

		return TestResult{
			TestCase: test,
		}
	}

	log.Println("title: '"+test.Title+"'", " Command: ", cut.Cmd)
	log.Println("title: '"+test.Title+"'", " Directory: ", cut.Dir)
	log.Println("title: '"+test.Title+"'", " Env: ", cut.Env)

	// Write test result
	test.Result = CommandResult{
		ExitCode: cut.ExitCode(),
		Stdout:   strings.Replace(cut.Stdout(), "\r\n", "\n", -1),
		Stderr:   strings.Replace(cut.Stderr(), "\r\n", "\n", -1),
	}

	log.Println("title: '"+test.Title+"'", " ExitCode: ", test.Result.ExitCode)
	log.Println("title: '"+test.Title+"'", " Stdout: ", test.Result.Stdout)
	log.Println("title: '"+test.Title+"'", " Stderr: ", test.Result.Stderr)

	return Validate(test)
}

// GetRetries returns the retries of the command
func (c *CommandUnderTest) GetRetries() int {
	if c.Retries == 0 {
		return 1
	}
	return c.Retries
}
