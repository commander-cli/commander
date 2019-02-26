package runtime

import (
	"github.com/SimonBaeumer/commander/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	SuccessCode = 0
	ErrorCode   = 1
)

func Test_ValidateExitCodeFail(t *testing.T) {
	test := commander.TestCase{}
	test.ExitCode = SuccessCode
	test.Result = commander.TestResult{
		ExitCode: ErrorCode,
	}

	got := Validate(test)

	assert.False(t, got.Success)
    assert.Len(t, got.Properties, 1)
}

func Test_ValidateExitCodeSuccess(t *testing.T) {
	test := commander.TestCase{}
	test.ExitCode = SuccessCode
	test.Result = commander.TestResult{
		ExitCode: SuccessCode,
	}

	got := Validate(test)

	assert.True(t, got.Success)
    assert.Len(t, got.Properties, 0)
}

func Test_ValidateStdoutFail(t *testing.T) {
	test := commander.TestCase{}
	test.Stdout = "foo"
	test.Result = commander.TestResult{
		Stdout: "bar",
	}

	got := Validate(test)

	assert.False(t, got.Success)
    assert.Len(t, got.Properties, 1)
}

func Test_ValidateStdoutSuccess(t *testing.T) {
	test := commander.TestCase{}
	test.Stdout = "foo"
	test.Result = commander.TestResult{
		Stdout: "foo",
	}

	got := Validate(test)

	assert.True(t, got.Success)
    assert.Len(t, got.Properties, 0)
}

func Test_ValidateStderrFail(t *testing.T) {
	test := commander.TestCase{Stdout: "foo"}
	test.Result = commander.TestResult{
		Stdout: "bar",
	}

	got := Validate(test)

	assert.False(t, got.Success)
    assert.Len(t, got.Properties, 1)
}

func Test_ValidateStderrSuccess(t *testing.T) {
	test := commander.TestCase{Stdout: "foo"}
	test.Result = commander.TestResult{
		Stdout: "foo",
	}

	got := Validate(test)

	assert.True(t, got.Success)
    assert.Len(t, got.Properties, 0)
}

func Test_Validate(t *testing.T) {
	test := commander.TestCase{
		Stdout:   "foo",
		Stderr:   "bar",
		ExitCode: SuccessCode,
	}
	test.Result = commander.TestResult{
		Stdout:   "foo",
		Stderr:   "bar",
		ExitCode: SuccessCode,
	}

	got := Validate(test)

	assert.True(t, got.Success)
	assert.Len(t, got.Properties, 0)
}

func Test_ValidateFail(t *testing.T) {
	test := commander.TestCase{
		Stdout:   "fail",
		Stderr:   "fail",
		ExitCode: ErrorCode,
	}
	test.Result = commander.TestResult{
		Stdout:   "foo",
		Stderr:   "bar",
		ExitCode: SuccessCode,
	}

	got := Validate(test)

	assert.False(t, got.Success)
	assert.Equal(t, []string{"Stdout", "Stderr", "ExitCode"}, got.Properties)
}

func Test_ValidateContains(t *testing.T) {
	test := commander.TestCase{
		Stdout:   `
✓ it should assert stderr
`,
		Stderr:   `
this
`,
		ExitCode: 0,
	}
	test.Result = commander.TestResult{
		Stdout:   `
✓ it should exit with error code
✓ it should assert stderr
✓ it should assert stdout
✓ it should fail'
`,
		Stderr:   `
this
is
my
stderr
and
more`,
		ExitCode: SuccessCode,
	}

	got := Validate(test)

	assert.True(t, got.Success)
	assert.Len(t, got.Properties, 0)
}

func Test_ValidateWithEmptyStdoutAndStderr(t *testing.T) {
    test := commander.TestCase{ExitCode: SuccessCode}
    test.Result = commander.TestResult{
        Stdout: "out",
        Stderr: "err",
        ExitCode: SuccessCode,
    }

    got := Validate(test)

    assert.True(t, got.Success)
    assert.Len(t, got.Properties, 0)
}