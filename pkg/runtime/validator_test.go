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
}

func Test_ValidateExitCodeSuccess(t *testing.T) {
	test := commander.TestCase{}
	test.ExitCode = SuccessCode
	test.Result = commander.TestResult{
		ExitCode: SuccessCode,
	}

	got := Validate(test)

	assert.True(t, got.Success)
}

func Test_ValidateStdoutFail(t *testing.T) {
	test := commander.TestCase{}
	test.Stdout = "foo"
	test.Result = commander.TestResult{
		Stdout: "bar",
	}

	got := Validate(test)

	assert.False(t, got.Success)
}

func Test_ValidateStdoutSuccess(t *testing.T) {
	test := commander.TestCase{}
	test.Stdout = "foo"
	test.Result = commander.TestResult{
		Stdout: "foo",
	}

	got := Validate(test)

	assert.True(t, got.Success)
}

func Test_ValidateStderrFail(t *testing.T) {
	test := commander.TestCase{Stdout: "foo"}
	test.Result = commander.TestResult{
		Stdout: "bar",
	}

	got := Validate(test)

	assert.False(t, got.Success)
}

func Test_ValidateStderrSuccess(t *testing.T) {
	test := commander.TestCase{Stdout: "foo"}
	test.Result = commander.TestResult{
		Stdout: "foo",
	}

	got := Validate(test)

	assert.True(t, got.Success)
}

func Test_Validate(t *testing.T) {
	test := commander.TestCase{
		Stdout:   "foo",
		Stderr:   "bar",
		ExitCode: 0,
	}
	test.Result = commander.TestResult{
		Stdout:   "foo",
		Stderr:   "bar",
		ExitCode: 0,
	}

	got := Validate(test)

	assert.True(t, got.Success)
	assert.Len(t, got.Properties, 0)
}

func Test_ValidateFail(t *testing.T) {
	test := commander.TestCase{
		Stdout:   "fail",
		Stderr:   "fail",
		ExitCode: 1,
	}
	test.Result = commander.TestResult{
		Stdout:   "foo",
		Stderr:   "bar",
		ExitCode: 0,
	}

	got := Validate(test)

	assert.False(t, got.Success)
	assert.Equal(t, []string{"Stdout", "Stderr", "ExitCode"}, got.Properties)
}
