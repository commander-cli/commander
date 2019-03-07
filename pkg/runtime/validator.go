package runtime

import (
	"github.com/SimonBaeumer/commander/pkg/matcher"
    "strings"
)

type ValidationResult struct {
	Success bool
	Diff    string
}

func NewValidationResult(m matcher.MatcherResult) ValidationResult {
	return ValidationResult{
		Success: m.Success,
		Diff:    m.Diff,
	}
}

func Validate(test TestCase) TestResult {
    equalMatcher := matcher.NewMatcher(matcher.Equal)

    matcherResult := validateExpectedOut(test.Result.Stdout, test.Expected.Stdout)
    if !matcherResult.Success {
        return TestResult{
            ValidationResult: NewValidationResult(matcherResult),
            TestCase:         test,
            FailedProperty:   Stdout,
        }
    }

    matcherResult = validateExpectedOut(test.Result.Stderr, test.Expected.Stderr)
    if !matcherResult.Success {
        return TestResult{
            ValidationResult: NewValidationResult(matcherResult),
            TestCase:         test,
            FailedProperty:   Stderr,
        }
    }

    matcherResult = equalMatcher.Match(test.Result.ExitCode, test.Expected.ExitCode)
    if !matcherResult.Success {
        return TestResult{
            ValidationResult: NewValidationResult(matcherResult),
            TestCase:         test,
            FailedProperty:   ExitCode,
        }
    }

    return TestResult{
        ValidationResult: NewValidationResult(matcherResult),
        TestCase:         test,
    }
}

func validateExpectedOut(got string, expected  ExpectedOut) matcher.MatcherResult {
	var m matcher.Matcher
	var result matcher.MatcherResult

	if expected.Exactly != "" {
		m = matcher.NewMatcher(matcher.Text)
		if result = m.Match(got, expected.Exactly); !result.Success {
			return result
		}
	}

	if len(expected.Contains) > 0 {
		m = matcher.NewMatcher(matcher.Contains)
		for _, c := range expected.Contains {
			if result = m.Match(got, c); !result.Success {
				return result
			}
		}
	}

	if expected.LineCount != 0 {
        m = matcher.NewMatcher(matcher.Equal)
        count := strings.Count(got, "\n") + 1
        if got == "" {
            count = 0
        }

        if result = m.Match(count, expected.LineCount); !result.Success {
            return result
        }
    }

	if len(expected.Lines) > 0 {
	    m = matcher.NewMatcher(matcher.Equal)
	    actualLines := strings.Split(got, "\n")
	    for k, expL := range expected.Lines {
            if result = m.Match(actualLines[k], expL); !result.Success {
                return result
            }
        }
    }

	result.Success = true
	return result
}