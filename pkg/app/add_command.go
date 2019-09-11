package app

import (
	"github.com/SimonBaeumer/commander/pkg/cmd"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
	"gopkg.in/yaml.v2"
)

// AddCommand executes the add command
// command is the command which should be added to the test suite
// existed holds the existing yaml content
func AddCommand(command string, existed []byte) ([]byte, error) {
	conf := suite.YAMLConfig{
		Tests:  make(map[string]suite.YAMLTest),
		Config: suite.YAMLTestConfig{},
	}
	c := cmd.NewCommand(command)

	if err := c.Execute(); err != nil {
		return []byte{}, err
	}

	//If a suite existed before adding the new command it is need to parse it and re-add it
	if len(existed) > 0 {
		err := yaml.UnmarshalStrict(existed, &conf)
		if err != nil {
			panic(err.Error())
		}

		for k, t := range conf.Tests {
			test := suite.YAMLTest{
				Title:    t.Title,
				Stdout:   t.Stdout.(runtime.ExpectedOut),
				Stderr:   t.Stderr.(runtime.ExpectedOut),
				ExitCode: t.ExitCode,
				Config:   convertConfig(t.Config),
			}

			//If title and command are not equal add the command property to the struct
			if t.Title != t.Command {
				test.Command = t.Command
			}

			conf.Tests[k] = test
		}
	}

	conf.Tests[command] = suite.YAMLTest{
		Title:    command,
		Stdout:   runtime.ExpectedOut{Contains: []string{c.Stdout()}},
		Stderr:   runtime.ExpectedOut{Contains: []string{c.Stderr()}},
		ExitCode: c.ExitCode(),
	}

	out, err := yaml.Marshal(conf)
	if err != nil {
		return []byte{}, err
	}

	return out, nil
}

func stringOrNil(str string) interface{} {
	if str == "" {
		return nil
	}
	return str
}

func convertConfig(config suite.YAMLTestConfig) suite.YAMLTestConfig {
	if config.Dir == "" && len(config.Env) == 0 && config.Timeout == "" {
		return suite.YAMLTestConfig{}
	}
	return config
}
