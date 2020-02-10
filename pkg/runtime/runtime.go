package runtime

import (
	"fmt"
	"github.com/SimonBaeumer/cmd"
	"log"
	"os"
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

func NewRuntime(nodes ...Node) Runtime {
	local := Node{
		Name: "local",
		Type: "local",
		Addr: "localhost",
	}

	nodes = append(nodes, local)
	return Runtime{
		Nodes: nodes,
	}
}

type Runtime struct {
	Nodes []Node
}

// TestCase represents a test case which will be executed by the runtime
type TestCase struct {
	Title    string
	Command  CommandUnderTest
	Expected Expected
	Result   CommandResult
	Nodes    []string
}

type Node struct {
	Name         string
	Type         string
	User         string
	Pass         string
	Addr         string
	Image        string
	IdentityFile string
}

//TestConfig represents the configuration for a test
type TestConfig struct {
	Env        map[string]string
	Dir        string
	Timeout    string
	Retries    int
	Interval   string
	InheritEnv bool
	SSH        SSHExecutor
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
	Executors  []Executor
}

// TestResult represents the TestCase and the ValidationResult
type TestResult struct {
	TestCase         TestCase
	ValidationResult ValidationResult
	FailedProperty   string
	Tries            int
	Node             string
}

// Start starts the given test suite and executes all tests
// maxConcurrent configures the amount of go routines which will be started
func (r *Runtime) Start(tests []TestCase, maxConcurrent int) <-chan TestResult {
	in := make(chan TestCase)
	out := make(chan TestResult)

	go func(tests []TestCase) {
		defer close(in)
		for _, t := range tests {
			in <- t
		}
	}(tests)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(tests chan TestCase) {
		defer wg.Done()

		for t := range tests {
			// If no node was set use local mode
			if len(t.Nodes) == 0 {
				t.Nodes = []string{"local"}
			}

			for _, n := range t.Nodes {
				result := TestResult{}
				for i := 1; i <= t.Command.GetRetries(); i++ {

					e := r.getExecutor(n)
					result = e.Execute(t)
					result.Node = n
					result.Tries = i

					if result.ValidationResult.Success {
						break
					}

					executeRetryInterval(t)
				}
				out <- result
			}

		}
	}(in)

	go func(results chan TestResult) {
		wg.Wait()
		close(results)
	}(out)

	return out
}

func (r *Runtime) getExecutor(node string) Executor {
	if len(r.Nodes) == 0 {
		return LocalExecutor{}
	}

	for _, n := range r.Nodes {
		if n.Name == node {
			switch n.Type {
			case "ssh":
				return SSHExecutor{
					Password:     n.Pass,
					IdentityFile: n.IdentityFile,
					User:         n.User,
					Host:         n.Addr,
				}
			case "local":
				return LocalExecutor{}
			case "":
				return LocalExecutor{}
			default:
				log.Fatal(fmt.Sprintf("Node type %s not found for node %s", n.Type, n.Name))
			}
		}
	}

	log.Fatal(fmt.Sprintf("Node %s not found", node))
	return LocalExecutor{}
}

func execTest(t TestCase) TestResult {
	result := TestResult{}
	for i := 1; i <= t.Command.GetRetries(); i++ {

		e := LocalExecutor{}
		result = e.Execute(t)

		result.Tries = i
		if result.ValidationResult.Success {
			break
		}

		executeRetryInterval(t)
	}
	return result
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

func createEnvVarsOption(test TestCase) func(c *cmd.Command) {
	return func(c *cmd.Command) {
		// Add all env variables from parent process
		if test.Command.InheritEnv {
			for _, v := range os.Environ() {
				split := strings.Split(v, "=")
				c.AddEnv(split[0], split[1])
			}
		}

		// Add custom env variables
		for k, v := range test.Command.Env {
			c.AddEnv(k, v)
		}
	}
}

func createTimeoutOption(timeout string) (func(c *cmd.Command), error) {
	timeoutOpt := cmd.WithoutTimeout
	if timeout != "" {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return func(c *cmd.Command) {}, err
		}
		timeoutOpt = cmd.WithTimeout(d)
	}

	return timeoutOpt, nil
}

// GetRetries returns the retries of the command
func (c *CommandUnderTest) GetRetries() int {
	if c.Retries == 0 {
		return 1
	}
	return c.Retries
}
