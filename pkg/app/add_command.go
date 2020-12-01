package app

import (
	"strings"

	"github.com/commander-cli/cmd"
	"github.com/commander-cli/commander/pkg/suite"
	"gopkg.in/yaml.v2"
)

// AddCommand executes the add command
// command is the command which should be added to the test suite
// existed holds the existing yaml content
func AddCommand(command string, existed []byte) ([]byte, error) {
	conf := suite.YAMLSuiteConf{
		Tests:  make(map[string]suite.YAMLTest),
		Config: suite.YAMLTestConfigConf{},
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
				Stdout:   t.Stdout.(suite.ExpectedOut),
				Stderr:   t.Stderr.(suite.ExpectedOut),
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

	stdout := strings.TrimSpace(c.Stdout())
	stderr := strings.TrimSpace(c.Stderr())
	conf.Tests[command] = suite.YAMLTest{
		Title:    command,
		Stdout:   suite.ExpectedOut{Contains: []string{stdout}},
		Stderr:   suite.ExpectedOut{Contains: []string{stderr}},
		ExitCode: c.ExitCode(),
	}

	out, err := yaml.Marshal(conf)
	if err != nil {
		return []byte{}, err
	}

	return out, nil
}

func convertConfig(config suite.YAMLTestConfigConf) suite.YAMLTestConfigConf {
	if config.Dir == "" && len(config.Env) == 0 && config.Timeout == "" {
		return suite.YAMLTestConfigConf{}
	}
	return config
}
