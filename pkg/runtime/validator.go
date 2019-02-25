package runtime

import "github.com/SimonBaeumer/commander/pkg"

type ValidationResult struct {
	Success    bool
	Properties []string
}

func Validate(test commander.TestCase) ValidationResult {
	r := &ValidationResult{
		Success:    true,
		Properties: []string{},
	}

	if test.Stdout != "" && (test.Stdout != test.Result.Stdout) {
		r.Properties = append(r.Properties, commander.Stdout)
	}

	if test.Stderr != "" && (test.Stderr != test.Result.Stderr) {
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
