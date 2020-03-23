package suite

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/runtime"
)

// Suite represents the current tests, nodes and configs.
// It is used by the runtime to execute all tests and is an abstraction for the given source.
// In example it could be possible to add more formats like XML or a custom DSL implementation.
type Suite struct {
	TestCases []runtime.TestCase
	Config    runtime.GlobalTestConfig
	Nodes     []runtime.Node
}

// GetNodes returns all nodes defined in the suite
func (s Suite) GetNodes() []runtime.Node {
	return s.Nodes
}

// GetNodeByName returns a node by the given name
func (s Suite) GetNodeByName(name string) (runtime.Node, error) {
	for _, n := range s.Nodes {
		if n.Name == name {
			return n, nil
		}
	}
	return runtime.Node{}, fmt.Errorf("could not find node with name %s", name)
}

// AddTest pushes a new test to the suite
// if the test was already added it will panic
func (s Suite) AddTest(t runtime.TestCase) {
	if _, err := s.GetTestByTitle(t.Title); err == nil {
		panic(fmt.Sprintf("Tests %s was already added to the suite", t.Title))
	}
	s.TestCases = append(s.TestCases, t)
}

// GetTests returns all tests of the test suite
func (s Suite) GetTests() []runtime.TestCase {
	return s.TestCases
}

// GetTestByTitle returns a test by title, if the test was not found an error is returned
func (s Suite) GetTestByTitle(title string) (runtime.TestCase, error) {
	for _, t := range s.GetTests() {
		if t.Title == title {
			return t, nil
		}
	}
	return runtime.TestCase{}, fmt.Errorf("could not find test %s", title)
}

// GetGlobalConfig returns the global configuration which applies to the complete suite
func (s Suite) GetGlobalConfig() runtime.GlobalTestConfig {
	return s.Config
}
