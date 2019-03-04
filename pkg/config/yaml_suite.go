package config

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"gopkg.in/yaml.v2"
	"log"
)

// YAMLConfig will be used for unmarshalling the yaml test suite
type YAMLConfig struct {
	Tests map[string]YAMLTest `yaml:"tests"`
}

// YAMLTest represents a test in the yaml test suite
type YAMLTest struct {
	Title    string      `yaml:"-"`
	Command  string      `yaml:"command"`
	ExitCode int         `yaml:"exit-code"`
	Stdout   interface{} `yaml:"stdout,omitempty"`
	Stderr   interface{} `yaml:"stderr,omitempty"`
}

// ParseYAML parses the Suite from a byte slice
func ParseYAML(content []byte) []runtime.TestCase {
	yamlConfig := YAMLConfig{}

	err := yaml.Unmarshal(content, &yamlConfig)
	if err != nil {
		log.Fatal(err)
	}

	return convertYAMLConfToTestCases(yamlConfig)
}

//Convert YAMlConfig to runtime TestCases
func convertYAMLConfToTestCases(conf YAMLConfig) []runtime.TestCase {
	var tests []runtime.TestCase
	for _, t := range conf.Tests {
		tests = append(tests, runtime.TestCase{
			Title:    t.Title,
			Command:  runtime.CommandUnderTest{
				Cmd: t.Command,
			},
			Expected: runtime.Expected{
				ExitCode: t.ExitCode,
				Stdout: t.Stdout.(runtime.ExpectedOut),
				Stderr: t.Stderr.(runtime.ExpectedOut),
			},
		})
	}

	return tests
}

func toString(s interface{}) string {
	return fmt.Sprintf("%s", s)
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
		test := YAMLTest{
			Title: k,
			Command: v.Command,
			ExitCode: v.ExitCode,
			Stdout: y.convertToExpectedOut(v.Stdout),
			Stderr: y.convertToExpectedOut(v.Stderr),
		}

		// Set key as command, if command property was empty
		if v.Command == "" {
			test.Command = k
		}

		y.Tests[k] = test
	}

	return nil
}

//Converts given value to an ExpectedOut. Especially used for Stdout and Stderr.
func (y *YAMLConfig) convertToExpectedOut(value interface{}) runtime.ExpectedOut {
    exp := runtime.ExpectedOut{}

    switch value.(type) {
    case string:
        exp.Exactly = toString(value)
        break
    case map[interface{}]interface{}:
    	if contains := value.(map[interface{}]interface{})["contains"]; contains != nil {
		    values := contains.([]interface{})
		    for _, v := range values {
			    exp.Contains = append(exp.Contains, toString(v))
		    }
	    }

        if exactly := value.(map[interface{}]interface{})["exactly"]; exactly != nil {
        	exact := toString(exactly)
            exp.Exactly = exact
        }
        break
    }

    return exp
}
