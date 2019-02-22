package runtime

import "github.com/SimonBaeumer/commander/pkg"

type ValidationResult struct {
    Success  bool
    Property string
}

func Validate(test commander.TestCase) ValidationResult {
    r := ValidationResult{}

    result := validateOutput(test.Stdout, test.Result.Stdout)
    if !result {
        r.Success = result
        test.Result.FailureProperty = commander.Stdout
    }

    result = validateOutput(test.Stderr, test.Result.Stderr)
    if !result {
        r.Success = result
        r.Property = commander.Stderr
    }

    result = validateExitCode(test.ExitCode, test.Result.ExitCode)
    if !result {
        r.Success = result
        r.Property = commander.ExitCode
    }

    r.Success = true
    return r
}

func validateOutput(expected string, actual string) bool {
    if expected != actual {
        return false
    }
    return true
}

func validateExitCode(expected int, actual int) bool {
    if expected != actual {
        return false
    }
    return true
}