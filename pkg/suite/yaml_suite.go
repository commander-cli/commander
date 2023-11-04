package suite

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/commander-cli/commander/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
)

// YAMLSuiteConf will be used for unmarshalling the yaml test suite
type YAMLSuiteConf struct {
	Tests  map[string]YAMLTest     `yaml:"tests"`
	Config YAMLTestConfigConf      `yaml:"config,omitempty"`
	Nodes  map[string]YAMLNodeConf `yaml:"nodes,omitempty"`
}

// YAMLTestConfigConf is a struct to represent the test config
type YAMLTestConfigConf struct {
	InheritEnv bool              `yaml:"inherit-env,omitempty"`
	Env        map[string]string `yaml:"env,omitempty"`
	Dir        string            `yaml:"dir,omitempty"`
	Timeout    string            `yaml:"timeout,omitempty"`
	Retries    int               `yaml:"retries,omitempty"`
	Interval   string            `yaml:"interval,omitempty"`
	Nodes      []string          `yaml:"nodes,omitempty"`
}

type YAMLNodeConf struct {
	Name           string `yaml:"-"`
	Type           string `yaml:"type"`
	User           string `yaml:"user,omitempty"`
	Pass           string `yaml:"pass,omitempty"`
	Addr           string `yaml:"addr,omitempty"`
	Image          string `yaml:"image,omitempty"`
	IdentityFile   string `yaml:"identity-file,omitempty"`
	Privileged     bool   `yaml:"privileged,omitempty"`
	DockerExecUser string `yaml:"docker-exec-user,omitempty"`
}

// YAMLTest represents a test in the yaml test suite
type YAMLTest struct {
	Title    string             `yaml:"-"`
	Command  string             `yaml:"command,omitempty"`
	ExitCode int                `yaml:"exit-code"`
	Stdout   interface{}        `yaml:"stdout,omitempty"`
	Stderr   interface{}        `yaml:"stderr,omitempty"`
	Config   YAMLTestConfigConf `yaml:"config,omitempty"`
	Skip     bool               `yaml:"skip,omitempty"`
}

// ParseYAML parses the Suite from a yaml byte slice
func ParseYAML(content []byte, fileName string) Suite {
	yamlConfig := YAMLSuiteConf{}

	err := yaml.UnmarshalStrict(content, &yamlConfig)
	if err != nil {
		panic(err.Error())
	}

	tests := convertYAMLSuiteConfToTestCases(yamlConfig, fileName)

	return Suite{
		TestCases: tests,
		Config: runtime.GlobalTestConfig{
			InheritEnv: yamlConfig.Config.InheritEnv,
			Env:        yamlConfig.Config.Env,
			Dir:        yamlConfig.Config.Dir,
			Timeout:    yamlConfig.Config.Timeout,
			Retries:    yamlConfig.Config.Retries,
			Interval:   yamlConfig.Config.Interval,
			Nodes:      yamlConfig.Config.Nodes,
		},
		Nodes: convertNodes(yamlConfig.Nodes),
	}
}

func convertNodes(nodeConfs map[string]YAMLNodeConf) []runtime.Node {
	var nodes []runtime.Node
	for _, v := range nodeConfs {
		node := runtime.Node{
			Pass:           v.Pass,
			Type:           v.Type,
			User:           v.User,
			Addr:           v.Addr,
			Name:           v.Name,
			Image:          v.Image,
			IdentityFile:   v.IdentityFile,
			Privileged:     v.Privileged,
			DockerExecUser: v.DockerExecUser,
		}

		node.ExpandEnv()
		nodes = append(nodes, node)
	}
	return nodes
}

// Convert YAMLSuiteConf to runtime TestCases
func convertYAMLSuiteConfToTestCases(conf YAMLSuiteConf, fileName string) []runtime.TestCase {
	var tests []runtime.TestCase
	for _, t := range conf.Tests {
		tests = append(tests, runtime.TestCase{
			Title: t.Title,
			Command: runtime.CommandUnderTest{
				Cmd:        t.Command,
				InheritEnv: t.Config.InheritEnv,
				Env:        t.Config.Env,
				Dir:        t.Config.Dir,
				Timeout:    t.Config.Timeout,
				Retries:    t.Config.Retries,
				Interval:   t.Config.Interval,
			},
			Expected: runtime.Expected{
				ExitCode: t.ExitCode,
				Stdout:   t.Stdout.(runtime.ExpectedOut),
				Stderr:   t.Stderr.(runtime.ExpectedOut),
			},
			Nodes:    t.Config.Nodes,
			FileName: fileName,
			Skip:     t.Skip,
		})
	}

	return tests
}

// Convert variable to string and remove trailing blank lines
func toString(s interface{}) string {
	return strings.Trim(fmt.Sprintf("%s", s), "\n")
}

// UnmarshalYAML unmarshals the yaml
func (y *YAMLSuiteConf) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var params struct {
		Tests  map[string]YAMLTest     `yaml:"tests"`
		Config YAMLTestConfigConf      `yaml:"config"`
		Nodes  map[string]YAMLNodeConf `yaml:"nodes"`
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
			Config:   v.Config,
			Skip:     v.Skip,
		}

		// Set key as command, if command property was empty
		if v.Command == "" {
			test.Command = k
		}

		y.Tests[k] = test
	}

	y.Nodes = make(map[string]YAMLNodeConf)
	for k, v := range params.Nodes {
		node := YAMLNodeConf{
			Name:           k,
			Addr:           v.Addr,
			User:           v.User,
			Type:           v.Type,
			Pass:           v.Pass,
			IdentityFile:   v.IdentityFile,
			Image:          v.Image,
			Privileged:     v.Privileged,
			DockerExecUser: v.DockerExecUser,
		}

		y.Nodes[k] = node
	}

	//Parse global configuration
	y.Config = YAMLTestConfigConf{
		InheritEnv: params.Config.InheritEnv,
		Env:        params.Config.Env,
		Dir:        params.Config.Dir,
		Timeout:    params.Config.Timeout,
		Retries:    params.Config.Retries,
		Interval:   params.Config.Interval,
		Nodes:      params.Config.Nodes,
	}

	return nil
}

// Converts given value to an ExpectedOut. Especially used for Stdout and Stderr.
func (y *YAMLSuiteConf) convertToExpectedOut(value interface{}) runtime.ExpectedOut {
	exp := runtime.ExpectedOut{
		JSON: make(map[string]string),
	}

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
		for k := range v {
			switch k {
			case
				"contains",
				"exactly",
				"line-count",
				"lines",
				"json",
				"xml",
				"file",
				"not-contains":
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

		// Parse file key
		if file := v["file"]; file != nil {
			exp.File = toString(file)
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

		if notContains := v["not-contains"]; notContains != nil {
			values := notContains.([]interface{})
			for _, v := range values {
				exp.NotContains = append(exp.NotContains, toString(v))
			}
		}

		if json := v["json"]; json != nil {
			values := json.(map[interface{}]interface{})
			for k, v := range values {
				exp.JSON[k.(string)] = v.(string)
			}
		}
		break

	case nil:
		break
	default:
		panic(fmt.Sprintf("Failed to parse Stdout or Stderr with values: %v", value))
	}

	return exp
}

// MarshalYAML adds custom logic to the struct to yaml conversion
func (y YAMLSuiteConf) MarshalYAML() (interface{}, error) {
	//Detect which values of the stdout/stderr assertions should be filled.
	//If all values are empty except Contains it will convert it to a single string
	//to match the easiest test suite definitions
	for k, t := range y.Tests {
		t.Stdout = convertExpectedOut(t.Stdout.(runtime.ExpectedOut))
		if reflect.ValueOf(t.Stdout).Kind() == reflect.Struct {
			t.Stdout = t.Stdout.(runtime.ExpectedOut)
		}

		t.Stderr = convertExpectedOut(t.Stderr.(runtime.ExpectedOut))
		if reflect.ValueOf(t.Stderr).Kind() == reflect.Struct {
			t.Stderr = t.Stderr.(runtime.ExpectedOut)
		}

		y.Tests[k] = t
	}

	return y, nil
}

func (y *YAMLSuiteConf) mergeNodes(nodes map[string]YAMLNodeConf, globalNodes map[string]YAMLNodeConf) map[string]YAMLNodeConf {
	return nodes
}

func convertExpectedOut(out runtime.ExpectedOut) interface{} {
	//If the property contains consists of only one element it will be set without the struct structure
	if isContainsASingleNonEmptyString(out) && propertiesAreEmpty(out) {
		return out.Contains[0]
	}

	//If the contains property only has one empty string element it should not be displayed
	//in the marshaled yaml file
	if len(out.Contains) == 1 && out.Contains[0] == "" {
		out.Contains = nil
	}

	if len(out.Contains) == 0 && propertiesAreEmpty(out) {
		return nil
	}
	return out
}

func propertiesAreEmpty(out runtime.ExpectedOut) bool {
	return out.Lines == nil &&
		out.Exactly == "" &&
		out.LineCount == 0 &&
		out.NotContains == nil
}

func isContainsASingleNonEmptyString(out runtime.ExpectedOut) bool {
	return len(out.Contains) == 1 &&
		out.Contains[0] != ""
}
