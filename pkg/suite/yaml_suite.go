package suite

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"gopkg.in/yaml.v2"
	"strings"
)

const defaultTimeoutMs = 5000

// YAMLConfig will be used for unmarshalling the yaml test suite
type YAMLConfig struct {
	Tests  map[string]YAMLTest `yaml:"tests"`
	Config YAMLTestConfig      `yaml:"config"`
}

type YAMLTestConfig struct {
	Env     []string `yaml:"env"`
	Dir     string   `yaml:"dir"`
	Timeout int      `yaml:"timeout"`
}

// YAMLTest represents a test in the yaml test suite
type YAMLTest struct {
	Title    string         `yaml:"-"`
	Command  string         `yaml:"command"`
	ExitCode int            `yaml:"exit-code"`
	Stdout   interface{}    `yaml:"stdout,omitempty"`
	Stderr   interface{}    `yaml:"stderr,omitempty"`
	Config   YAMLTestConfig `yaml:"config,omitempty"`
}

//YAMLSuite represents a test suite which was configured in yaml
type YAMLSuite struct {
	TestCases []runtime.TestCase
	Config    runtime.TestConfig
}

// GetTests returns all tests of the test suite
func (s YAMLSuite) GetTests() []runtime.TestCase {
	return s.TestCases
}

//GetTestByTitle returns the first test it finds for the given title
func (s YAMLSuite) GetTestByTitle(title string) (runtime.TestCase, error) {
	for _, t := range s.GetTests() {
		if t.Title == title {
			return t, nil
		}
	}
	return runtime.TestCase{}, fmt.Errorf("Could not find test " + title)
}

func (s YAMLSuite) GetGlobalConfig() runtime.TestConfig {
	return s.Config
}

// ParseYAML parses the Suite from a yaml byte slice
func ParseYAML(content []byte) Suite {
	yamlConfig := YAMLConfig{}

	err := yaml.UnmarshalStrict(content, &yamlConfig)
	if err != nil {
		panic(err.Error())
	}

	return YAMLSuite{
		TestCases: convertYAMLConfToTestCases(yamlConfig),
		Config: runtime.TestConfig{
			Env:     yamlConfig.Config.Env,
			Dir:     yamlConfig.Config.Dir,
			Timeout: yamlConfig.Config.Timeout,
		},
	}
}

//Convert YAMlConfig to runtime TestCases
func convertYAMLConfToTestCases(conf YAMLConfig) []runtime.TestCase {
	var tests []runtime.TestCase
	for _, t := range conf.Tests {
		tests = append(tests, runtime.TestCase{
			Title: t.Title,
			Command: runtime.CommandUnderTest{
				Cmd:     t.Command,
				Env:     t.Config.Env,
				Dir:     t.Config.Dir,
				Timeout: t.Config.Timeout,
			},
			Expected: runtime.Expected{
				ExitCode: t.ExitCode,
				Stdout:   t.Stdout.(runtime.ExpectedOut),
				Stderr:   t.Stderr.(runtime.ExpectedOut),
			},
		})
	}

	return tests
}

// Convert variable to string and remove trailing blank lines
func toString(s interface{}) string {
	return strings.Trim(fmt.Sprintf("%s", s), "\n")
}

// UnmarshalYAML unmarshals the yaml
func (y *YAMLConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params struct {
		Tests  map[string]YAMLTest `yaml:"tests"`
		Config YAMLTestConfig      `yaml:"config"`
	}

	err := unmarshal(&params)
	if err != nil {
		return err
	}

	// map key to title property
	y.Tests = make(map[string]YAMLTest)
	for k, v := range params.Tests {
		test := YAMLTest{
			Title:    k,
			Command:  v.Command,
			ExitCode: v.ExitCode,
			Stdout:   y.convertToExpectedOut(v.Stdout),
			Stderr:   y.convertToExpectedOut(v.Stderr),
			Config:   y.mergeConfigs(v.Config, params.Config),
		}

		// Set key as command, if command property was empty
		if v.Command == "" {
			test.Command = k
		}

		y.Tests[k] = test
	}

	//Parse global configuration
	y.Config = YAMLTestConfig{
		Env:     params.Config.Env,
		Dir:     params.Config.Dir,
		Timeout: params.Config.Timeout,
	}

	return nil
}

//Converts given value to an ExpectedOut. Especially used for Stdout and Stderr.
func (y *YAMLConfig) convertToExpectedOut(value interface{}) runtime.ExpectedOut {
	exp := runtime.ExpectedOut{}

	switch value.(type) {
	//If only a string was passed it is assigned to exactly automatically
	case string:
		exp.Contains = []string{toString(value)}
		break

	//If there is nested map set the properties will be assigned to the contains
	case map[interface{}]interface{}:
		v := value.(map[interface{}]interface{})
		// Check if keys are parsable
		// TODO: Could be refactored into a registry maybe which holds all parsers
		for k, _ := range v {
			switch k {
			case
				"contains",
				"exactly",
				"line-count",
				"lines":
				break
			default:
				panic(fmt.Sprintf("Key %s is not allowed.", k))
			}
		}

		//Parse contains key
		if contains := v["contains"]; contains != nil {
			values := contains.([]interface{})
			for _, v := range values {
				exp.Contains = append(exp.Contains, toString(v))
			}
		}

		//Parse exactly key
		if exactly := v["exactly"]; exactly != nil {
			exp.Exactly = toString(exactly)
		}

		//Parse line-count key
		if lc := v["line-count"]; lc != nil {
			exp.LineCount = lc.(int)
		}

		// Parse lines
		if l := v["lines"]; l != nil {
			exp.Lines = make(map[int]string)
			for k, v := range l.(map[interface{}]interface{}) {
				exp.Lines[k.(int)] = toString(v)
			}
		}
		break
	}

	return exp
}

func (y *YAMLConfig) mergeConfigs(local YAMLTestConfig, global YAMLTestConfig) YAMLTestConfig {
	conf := global
	if len(local.Env) > 0 {
		conf.Env = local.Env
	}

	if local.Dir != "" {
		conf.Dir = local.Dir
	}

	if local.Timeout != 0 {
		conf.Timeout = local.Timeout
	}
	if conf.Timeout == 0 {
		conf.Timeout = defaultTimeoutMs
	}

	return conf
}
