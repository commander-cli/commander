package runtime

import (
	"github.com/SimonBaeumer/commander/pkg"
	"strings"
)

type ValidationResult struct {
	Success    bool
	Properties []string
}

func Validate(test commander.TestCase) ValidationResult {
	r := &ValidationResult{
		Success:    true,
		Properties: []string{},
	}

	if test.Stdout != "" && !strings.Contains(test.Result.Stdout, test.Stdout) {
		r.Properties = append(r.Properties, commander.Stdout)
	}

	if test.Stderr != "" && !strings.Contains(test.Result.Stderr, test.Stderr) {
		r.Properties = append(r.Properties, commander.Stderr)
	}

	if test.ExitCode != test.Result.ExitCode {
		r.Properties = append(r.Properties, commander.ExitCode)
	}

	if len(r.Properties) > 0 {
		r.Success = false
	}

	return *r
}
