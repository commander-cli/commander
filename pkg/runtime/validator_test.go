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

func Test_ValidateExpectedOut_MatchLines(t *testing.T) {
	value := `my
multi
line
output`

	got := validateExpectedOut(value, ExpectedOut{Lines: map[int]string{1: "my", 3: "line"}})

	assert.True(t, got.Success)
	assert.Empty(t, got.Diff)
}

func Test_ValidateExpectedOut_PanicIfLineDoesNotExist_TooHigh(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			assert.Equal(t, "Invalid line number given 99", r)
		}
		assert.NotNil(t, r)
	}()

	value := `my
multi
line
output`

	_ = validateExpectedOut(value, ExpectedOut{Lines: map[int]string{99: "my"}})
}

func Test_ValidateExpectedOut_PanicIfLineDoesNotExist(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			assert.Equal(t, "Invalid line number given 0", r)
		}
		assert.NotNil(t, r)
	}()

	value := `my`
	_ = validateExpectedOut(value, ExpectedOut{Lines: map[int]string{0: "my"}})
}

func Test_ValidateExpectedOut_ValidateJSON(t *testing.T) {
	json := `
{
  "object": {
    "attr": "test"
  }
}
`
	r := validateExpectedOut(json, ExpectedOut{JSON: map[string]string{"object.attr": "test"}})
	assert.True(t, r.Success)

	diff := `Expected json path "object.attr" with result

no

to be equal to

test`
	r = validateExpectedOut(json, ExpectedOut{JSON: map[string]string{"object.attr": "no"}})
	assert.False(t, r.Success)
	assert.Equal(t, diff, r.Diff)
}

func Test_ValidateExpectedOut_ValidateXML(t *testing.T) {
	xml := `<book>
  <author>J. R. R. Tolkien</author>
</book>`

	r := validateExpectedOut(xml, ExpectedOut{XML: map[string]string{"/book//author": "J. R. R. Tolkien"}})
	assert.True(t, r.Success)
	assert.Equal(t, "", r.Diff)

	diff := `Expected xml path "/book//author" with result

Joanne K. Rowling

to be equal to

J. R. R. Tolkien`
	r = validateExpectedOut(xml, ExpectedOut{XML: map[string]string{"/book//author": "Joanne K. Rowling"}})
	assert.False(t, r.Success)
	assert.Equal(t, diff, r.Diff)

	r = validateExpectedOut(xml, ExpectedOut{XML: map[string]string{"/book//title": "J. R. R. Tolkien"}})
	assert.False(t, r.Success)
	assert.Equal(t, `Query "/book//title" did not match a path`, r.Diff)
}

func getExampleTest() TestCase {
	test := TestCase{
		Expected: Expected{
			Stdout: ExpectedOut{
				Lines:       map[int]string{1: "hello"},
				LineCount:   1,
				Exactly:     "hello",
				Contains:    []string{"hello"},
				NotContains: []string{"not-exist"},
			},
			Stderr: ExpectedOut{
				Lines:       map[int]string{1: "error"},
				LineCount:   1,
				Exactly:     "error",
				Contains:    []string{"error"},
				NotContains: []string{"not-exist"},
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
