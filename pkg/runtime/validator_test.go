package runtime

import (
	"github.com/SimonBaeumer/commander/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewValidationResult(t *testing.T) {
	mr := matcher.MatcherResult{}
	got := newValidationResult(mr)
	assert.IsType(t, ValidationResult{}, got)
}

func Test_Validate(t *testing.T) {
	test := getExampleTest()

	got := Validate(test)

	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, test, got.TestCase)
}

func Test_ValidateStdoutShouldFail(t *testing.T) {
	test := getExampleTest()
	test.Result = CommandResult{
		Stdout:   "hello\nline2",
		Stderr:   "error",
		ExitCode: 0,
	}

	got := Validate(test)

	assert.False(t, got.ValidationResult.Success)
	assert.Equal(t, "Stdout", got.FailedProperty)
}

func getExampleTest() TestCase {
	test := TestCase{
		Expected: Expected{
			Stdout: ExpectedOut{
				Lines:     map[int]string{0: "hello"},
				LineCount: 1,
				Exactly:   "hello",
				Contains:  []string{"hello"},
			},
			Stderr: ExpectedOut{
				Lines:     map[int]string{0: "error"},
				LineCount: 1,
				Exactly:   "error",
				Contains:  []string{"error"},
			},
			LineCount: 1,
		},
		Result: CommandResult{
			Stdout:   "hello",
			Stderr:   "error",
			ExitCode: 0,
		},
	}
	return test
}
