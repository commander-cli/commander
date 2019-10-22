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

func Test_ValidateStderrShouldFail(t *testing.T) {
	test := getExampleTest()
	test.Expected.Stdout = ExpectedOut{}
	test.Result = CommandResult{
		Stderr:   "is not in message",
		ExitCode: 0,
	}

	got := Validate(test)

	assert.False(t, got.ValidationResult.Success)
	assert.Equal(t, "Stderr", got.FailedProperty)
}

func Test_ValidateExitCodeShouldFail(t *testing.T) {
	test := getExampleTest()
	test.Expected.Stdout = ExpectedOut{}
	test.Expected.Stderr = ExpectedOut{}
	test.Result = CommandResult{
		ExitCode: 1,
	}

	got := Validate(test)

	assert.False(t, got.ValidationResult.Success)
	assert.Equal(t, "ExitCode", got.FailedProperty)
}

func Test_ValidateExpectedOut_Contains_Fails(t *testing.T) {
	value := `test`

	got := validateExpectedOut(value, ExpectedOut{Contains: []string{"not-exists"}})

	diff := `
Expected

test

to contain

not-exists
`

	assert.False(t, got.Success)
	assert.Equal(t, diff, got.Diff)
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

func Test_ValidateExpectedOut_MatchLines_ExpectedLineDoesNotExists(t *testing.T) {
	value := `test`

	got := validateExpectedOut(value, ExpectedOut{Lines: map[int]string{2: "my"}})

	assert.False(t, got.Success)
	diff := `Line number 2 does not exists in result: 

test`
	assert.Equal(t, diff, got.Diff)
}

func Test_ValidateExpectedOut_MatchLines_Fails(t *testing.T) {
	value := `test
line 2
line 3`

	got := validateExpectedOut(value, ExpectedOut{Lines: map[int]string{2: "line 3"}})

	assert.False(t, got.Success)
	diff := `--- Got
+++ Expected
@@ -1 +1 @@
-line 2
+line 3
`
	assert.Equal(t, diff, got.Diff)
}

func Test_ValidateExpectedOut_LineCount_Fails(t *testing.T) {
	value := ``

	got := validateExpectedOut(value, ExpectedOut{LineCount: 2})

	assert.False(t, got.Success)
	diff := `--- Got
+++ Expected
@@ -1 +1 @@
-0
+2
`
	assert.Equal(t, diff, got.Diff)
}

func Test_ValidateExpectedOut_NotContains_Fails(t *testing.T) {
	value := `my string contains`

	got := validateExpectedOut(value, ExpectedOut{NotContains: []string{"contains"}})

	diff := `
Expected

my string contains

to not contain

contains
`
	assert.False(t, got.Success)
	assert.Equal(t, diff, got.Diff)
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
