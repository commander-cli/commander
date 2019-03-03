package runtime

import (
	"github.com/SimonBaeumer/commander/pkg/suite"
	"strings"
)

type ValidationResult struct {
	Success    bool
	Properties []string
}

func Validate(test suite.TestCase) ValidationResult {
	r := &ValidationResult{
		Success:    true,
		Properties: []string{},
	}

	if test.Expected.Stdout.Exactly != "" && !strings.Contains(test.Result.Stdout, test.Expected.Stdout.Exactly) {
		r.Properties = append(r.Properties, suite.Stdout)
	}

	if test.Expected.Stdout.Exactly != "" && !strings.Contains(test.Result.Stderr, test.Expected.Stderr.Exactly) {
		r.Properties = append(r.Properties, suite.Stderr)
	}

	if test.Expected.ExitCode != test.Result.ExitCode {
		r.Properties = append(r.Properties, suite.ExitCode)
	}

	if len(r.Properties) > 0 {
		r.Success = false
	}

	return *r
}
