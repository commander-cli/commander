package commander

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
)

// YAMLConfig will be used for unmarshalling the yaml test suite
type YAMLConfig struct {
    Tests map[string]YAMLTest `yaml:"tests"`
}

// YAMLTest represents a test in the yaml test suite
type YAMLTest struct {
    Title    string `yaml:"-"`
    Command  string `yaml:"command"`
    ExitCode int    `yaml:"exit-code"`
    Stdout   string `yaml:"stdout,omitempty"`
    Stderr   string `yaml:"stderr,omitempty"`
}

// YAMLSuite represents the complete test suite and implements the Suite interface
type YAMLSuite struct {
    TestCases []TestCase
}

// GetTestCases returns the test cases of the suite
func (s YAMLSuite) GetTestCases() []TestCase {
    return s.TestCases
}


// ParseYAMLFile takes a file and returns the Suite
func ParseYAMLFile(file string) Suite {
    c, err := ioutil.ReadFile(file)
    if err != nil {
        log.Fatal(err)
    }

    return ParseYAML(c)
}

// ParseYAML parses the Suite from a byte slice
func ParseYAML(content []byte) Suite {
    yamlConfig := YAMLConfig{}

    err := yaml.Unmarshal(content, &yamlConfig)
    if err != nil {
        log.Fatal(err)
    }

    var s Suite
    s = YAMLSuite{TestCases: convertYAMLConfToTestCases(yamlConfig)}
    return s
}

func convertYAMLConfToTestCases(conf YAMLConfig) []TestCase {
    var tests []TestCase
    for _, t := range conf.Tests {
        tests = append(tests, TestCase{
            Title: t.Title,
            Command: t.Command,
            ExitCode: t.ExitCode,
            Stdout: t.Stdout,
            Stderr: t.Stderr,
        })
    }

    return tests
}

// UnmarshalYAML unmarshals the yaml
func (y *YAMLConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var params struct {
        Tests map[string]YAMLTest `yaml:"tests"`
    }

    err := unmarshal(&params)
    if err != nil {
        log.Fatal(err)
    }

    // map key to title property
    y.Tests = make(map[string]YAMLTest)
    for k, v := range params.Tests {
        v.Title = k
        y.Tests[k] = v
    }

    return nil
}