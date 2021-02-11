package suite

import (
	"fmt"
	"regexp"

	"github.com/commander-cli/commander/pkg/runtime"
)

// Suite represents the current tests, nodes and configs.
// It is used by the runtime to execute all tests and is an abstraction for the given source.
// In example it could be possible to add more formats like XML or a custom DSL implementation.
type Suite struct {
	TestCases []runtime.TestCase
	Config    runtime.GlobalTestConfig
	Nodes     []runtime.Node
}

// NewSuite creates a suite structure from two byte slices,
// suiteContent is the file/suite that is under test
// overwriteConfigContent is an optional slice which overwrites the default configurations
// fileName is the file that is under test
func NewSuite(suiteContent, overwriteConfigContent []byte, fileName string) Suite {
	defaultConfig := ParseYAML(overwriteConfigContent, "default config")
	s := ParseYAML(suiteContent, fileName)

	s.mergeConfigs(defaultConfig.Config, defaultConfig.Nodes)

	return s
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

// FindTests returns a test by the given pattern, if the test was not found an error is returned
func (s Suite) FindTests(pattern string) ([]runtime.TestCase, error) {
	var r []runtime.TestCase
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
		return []runtime.TestCase{}, fmt.Errorf("could not find test with pattern: %s", pattern)
	}

	return r, nil
}

// GetGlobalConfig returns the global configuration which applies to the complete suite
func (s Suite) GetGlobalConfig() runtime.GlobalTestConfig {
	return s.Config
}

// MergeConfigs overwrites a global configuration over an entire suite.
// Config at the lowest level takes precedence
func (s Suite) mergeConfigs(config runtime.GlobalTestConfig, nodes []runtime.Node) {
	s.Config.Env = mergeEnvironmentVariables(s.Config.Env, config.Env)

	if s.Config.Dir == "" {
		s.Config.Dir = config.Dir
	}

	if s.Config.Timeout == "" {
		s.Config.Timeout = config.Timeout
	}

	if s.Config.Retries == 0 {
		s.Config.Retries = config.Retries
	}

	if s.Config.Interval == "" {
		s.Config.Interval = config.Interval
	}

	if !s.Config.InheritEnv {
		s.Config.InheritEnv = config.InheritEnv
	}

	if len(s.Config.Nodes) == 0 {
		s.Config.Nodes = config.Nodes
	}

	// append additional nodes
	s.Nodes = append(s.Nodes, nodes...)

	s.mergeTestConfigs()
}

// mergeConfigs will merge the suites runtime.GlobalTestConfig,
// with each runtime.TestCase in the suite
func (s Suite) mergeTestConfigs() {
	for i := range s.TestCases {

		s.TestCases[i].Command.Env = mergeEnvironmentVariables(s.Config.Env, s.TestCases[i].Command.Env)

		if s.TestCases[i].Command.Dir == "" {
			s.TestCases[i].Command.Dir = s.Config.Dir
		}

		if s.TestCases[i].Command.Timeout == "" {
			s.TestCases[i].Command.Timeout = s.Config.Timeout
		}

		if s.TestCases[i].Command.Retries == 0 {
			s.TestCases[i].Command.Retries = s.Config.Retries
		}

		if s.TestCases[i].Command.Interval == "" {
			s.TestCases[i].Command.Interval = s.Config.Interval
		}

		if !s.TestCases[i].Command.InheritEnv {
			s.TestCases[i].Command.InheritEnv = s.Config.InheritEnv
		}

		if len(s.TestCases[i].Nodes) == 0 {
			s.TestCases[i].Nodes = s.Config.Nodes
		}
	}
}

func mergeEnvironmentVariables(global map[string]string, local map[string]string) map[string]string {
	env := make(map[string]string)
	for k, v := range global {
		env[k] = v
	}
	for k, v := range local {
		env[k] = v
	}
	return env
}
