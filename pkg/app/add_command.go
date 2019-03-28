package app

import (
	"github.com/SimonBaeumer/commander/pkg/cmd"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
	"gopkg.in/yaml.v2"
)

func AddCommand(command string, existed []byte) ([]byte, error) {
	conf := suite.YAMLConfig{
		Tests:  make(map[string]suite.YAMLTest),
		Config: suite.YAMLTestConfig{},
	}
	c := cmd.NewCommand(command)

	if err := c.Execute(); err != nil {
		return []byte{}, err
	}

	if len(existed) > 0 {
		err := yaml.UnmarshalStrict(existed, &conf)
		if err != nil {
			panic(err.Error())
		}

		for k, t := range conf.Tests {
			conf.Tests[k] = suite.YAMLTest{
				Title:    t.Title,
				Stdout:   convertExpectedOut(t.Stdout.(runtime.ExpectedOut)),
				Stderr:   convertExpectedOut(t.Stderr.(runtime.ExpectedOut)),
				ExitCode: t.ExitCode,
				Config:   convertConfig(t.Config),
			}
		}
	}

	conf.Tests[command] = suite.YAMLTest{
		Title:    command,
		Stdout:   stringOrNil(c.Stdout()),
		Stderr:   stringOrNil(c.Stderr()),
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
	if config.Dir == "" && len(config.Env) == 0 && config.Timeout == 0 {
		return suite.YAMLTestConfig{}
	}
	return config
}

func convertExpectedOut(out runtime.ExpectedOut) interface{} {
	if len(out.Contains) == 1 && len(out.Lines) == 0 && out.Exactly == "" {
		return out.Contains[0]
	}
	if len(out.Contains) == 0 {
		return nil
	}
	return out
}
