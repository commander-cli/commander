package suite

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/runtime"
)

// Suite represents the current tests and configs
type Suite struct {
	TestCases []runtime.TestCase
	Config    runtime.TestConfig
	Nodes     []runtime.Node
}

func NewSuite(config runtime.TestConfig, tests ...runtime.TestCase) *Suite {
	return &Suite{
		TestCases: tests,
		Config:    config,
	}
}

func (s Suite) GetNodes() []runtime.Node {
	return s.Nodes
}

func (s Suite) AddTest(t runtime.TestCase) {
	if _, err := s.GetTestByTitle(t.Title); err != nil {
		panic(fmt.Sprintf("Tests %s was already added to the suite", t.Title))
	}
	s.TestCases = append(s.TestCases, t)
}

// GetTests returns all tests of the test suite
func (s Suite) GetTests() []runtime.TestCase {
	return s.TestCases
}

func (s Suite) GetTestByTitle(title string) (runtime.TestCase, error) {
	for _, t := range s.GetTests() {
		if t.Title == title {
			return t, nil
		}
	}
	return runtime.TestCase{}, fmt.Errorf("Could not find test " + title)
}

func (s Suite) GetGlobalConfig() runtime.TestConfig {
	return s.Config
}
