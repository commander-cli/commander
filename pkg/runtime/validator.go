package runtime

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/matcher"
	"log"
	"strings"
)

// ValidationResult will be returned after the validation was executed
type ValidationResult struct {
	Success bool
	Diff    string
}

func newValidationResult(m matcher.MatcherResult) ValidationResult {
	return ValidationResult{
		Success: m.Success,
		Diff:    m.Diff,
	}
}

// Validate validates the test results with the expected values
// The test should hold the result and expected to validate the result
func Validate(test TestCase) TestResult {
	equalMatcher := matcher.NewMatcher(matcher.Equal)

	log.Println("title: '"+test.Title+"'", " Stdout-Expected: ", test.Expected.Stdout)
	matcherResult := validateExpectedOut(test.Result.Stdout, test.Expected.Stdout)
	log.Println("title: '"+test.Title+"'", " Stdout-Result: ", matcherResult.Success)
	if !matcherResult.Success {
		return TestResult{
			ValidationResult: newValidationResult(matcherResult),
			TestCase:         test,
			FailedProperty:   Stdout,
		}
	}

	log.Println("title: '"+test.Title+"'", " Stderr-Expected: ", test.Expected.Stderr)
	matcherResult = validateExpectedOut(test.Result.Stderr, test.Expected.Stderr)
	log.Println("title: '"+test.Title+"'", " Stderr-Result: ", matcherResult.Success)
	if !matcherResult.Success {
		return TestResult{
			ValidationResult: newValidationResult(matcherResult),
			TestCase:         test,
			FailedProperty:   Stderr,
		}
	}

	log.Println("title: '"+test.Title+"'", " Exit-Expected: ", test.Expected.ExitCode)
	matcherResult = equalMatcher.Match(test.Result.ExitCode, test.Expected.ExitCode)
	log.Println("title: '"+test.Title+"'", " Exit-Result: ", matcherResult.Success)
	if !matcherResult.Success {
		return TestResult{
			ValidationResult: newValidationResult(matcherResult),
			TestCase:         test,
			FailedProperty:   ExitCode,
		}
	}

	return TestResult{
		ValidationResult: newValidationResult(matcherResult),
		TestCase:         test,
	}
}

func validateExpectedOut(got string, expected ExpectedOut) matcher.MatcherResult {
	var m matcher.Matcher
	var result matcher.MatcherResult
	result.Success = true

	if expected.Exactly != "" {
		m = matcher.NewMatcher(matcher.Text)
		if result = m.Match(got, expected.Exactly); !result.Success {
			return result
		}
	}

	m = matcher.NewMatcher(matcher.Contains)
	for _, c := range expected.Contains {
		if result = m.Match(got, c); !result.Success {
			return result
		}
	}

	if expected.LineCount != 0 {
		result = validateExpectedLineCount(got, expected)
		if !result.Success {
			return result
		}
	}

	if len(expected.Lines) > 0 {
		result = validateExpectedLines(got, expected)
		if !result.Success {
			return result
		}
	}

	m = matcher.NewMatcher(matcher.NotContains)
	for _, c := range expected.NotContains {
		if result = m.Match(got, c); !result.Success {
			return result
		}
	}

	m = matcher.NewMatcher(matcher.JSON)
	for i, v := range expected.JSON {
		if result = m.Match(got, map[string]string{i: v}); !result.Success {
			return result
		}
	}

	m = matcher.NewMatcher(matcher.XML)
	for i, v := range expected.XML {
		if result = m.Match(got, map[string]string{i: v}); !result.Success {
			return result
		}
	}

	return result
}

func validateExpectedLineCount(got string, expected ExpectedOut) matcher.MatcherResult {
	m := matcher.NewMatcher(matcher.Equal)
	count := strings.Count(got, getLineBreak()) + 1

	if got == "" {
		count = 0
	}

	return m.Match(count, expected.LineCount)
}

func validateExpectedLines(got string, expected ExpectedOut) matcher.MatcherResult {
	m := matcher.NewMatcher(matcher.Equal)
	actualLines := strings.Split(got, getLineBreak())
	result := matcher.MatcherResult{Success: true}

	for k, expL := range expected.Lines {
		if (k-1 > len(actualLines)) || (k-1 < 0) {
			panic(fmt.Sprintf("Invalid line number given %d", k))
		}

		if result = m.Match(actualLines[k-1], expL); !result.Success {
			return result
		}
	}

	return result
}

func getLineBreak() string {
	return "\n"
}
