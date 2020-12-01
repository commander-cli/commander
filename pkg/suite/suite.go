package suite

import (
	"fmt"
	"regexp"
)

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

// Suite represents the current tests, nodes and configs.
// It is used by the runtime to execute all tests and is an abstraction for the given source.
// In example it could be possible to add more formats like XML or a custom DSL implementation.
type Suite struct {
	TestCases []TestCase
	Config    GlobalTestConfig
	Nodes     []Node
}

// GetNodes returns all nodes defined in the suite
func (s Suite) GetNodes() []Node {
	return s.Nodes
}

// GetNodeByName returns a node by the given name
func (s Suite) GetNodeByName(name string) (Node, error) {
	for _, n := range s.Nodes {
		if n.Name == name {
			return n, nil
		}
	}
	return Node{}, fmt.Errorf("could not find node with name %s", name)
}

// AddTest pushes a new test to the suite
// if the test was already added it will panic
func (s Suite) AddTest(t TestCase) {
	if _, err := s.GetTestByTitle(t.Title); err == nil {
		panic(fmt.Sprintf("Tests %s was already added to the suite", t.Title))
	}
	s.TestCases = append(s.TestCases, t)
}

// GetTests returns all tests of the test suite
func (s Suite) GetTests() []TestCase {
	return s.TestCases
}

// GetTestByTitle returns a test by title, if the test was not found an error is returned
func (s Suite) GetTestByTitle(title string) (TestCase, error) {
	for _, t := range s.GetTests() {
		if t.Title == title {
			return t, nil
		}
	}

	return TestCase{}, fmt.Errorf("could not find test %s", title)
}

// GetTestByTitle returns a test by title, if the test was not found an error is returned
func (s Suite) FindTests(pattern string) ([]TestCase, error) {
	var r []TestCase
	for _, t := range s.GetTests() {
		matched, err := regexp.Match(pattern, []byte(t.Title))
		if err != nil {
			panic(fmt.Sprintf("Regex error %s: %s", pattern, err.Error()))
		}

		if matched {
			r = append(r, t)
		}
	}

	if len(r) == 0 {
		return []TestCase{}, fmt.Errorf("could not find test with pattern: %s", pattern)
	}

	return r, nil
}

// GetGlobalConfig returns the global configuration which applies to the complete suite
func (s Suite) GetGlobalConfig() GlobalTestConfig {
	return s.Config
}

// GetRetries returns the retries of the command
func (c *CommandUnderTest) GetRetries() int {
	if c.Retries == 0 {
		return 1
	}
	return c.Retries
}
