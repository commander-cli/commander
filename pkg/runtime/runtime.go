package runtime

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SimonBaeumer/cmd"
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

type Filters []string

// NewRuntime creates a new runtime and inits default nodes
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

// Runtime represents the current runtime, please use NewRuntime() instead of creating an instance directly
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
	Error            error
	FileName         string
}

// Start starts the given test suite and executes all tests
func (r *Runtime) Start(tests []TestCase) <-chan TestResult {
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
			// If no node was set use local mode as default
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
		return NewLocalExecutor()
	}

	for _, n := range r.Nodes {
		if n.Name == node {
			switch n.Type {
			case "ssh":
				return NewSSHExecutor(n.Addr, n.User, WithIdentityFile(n.IdentityFile), WithPassword(n.Pass))
			case "local":
				return NewLocalExecutor()
			case "docker":
				log.Println("Use docker executor")
				return DockerExecutor{
					Image:        n.Image,
					Privileged:   n.Privileged,
					ExecUser:     n.DockerExecUser,
					RegistryPass: n.Pass,
					RegistryUser: n.User,
				}
			case "":
				return NewLocalExecutor()
			default:
				log.Fatal(fmt.Sprintf("Node type %s not found for node %s", n.Type, n.Name))
			}
		}
	}

	log.Fatal(fmt.Sprintf("Node %s not found", node))
	return NewLocalExecutor()
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
