package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DockerExecutor_Execute(t *testing.T) {
	d := DockerExecutor{
		Image: "ubuntu:18.04",
	}

	cmd := `echo hello`

	test := TestCase{
		Command: CommandUnderTest{
			Cmd: cmd,
		},
		Expected: Expected{
			Stdout:   ExpectedOut{Exactly: "hello"},
			ExitCode: 0,
		},
	}

	got := d.Execute(test)
	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, "hello", got.TestCase.Result.Stdout)
	assert.Equal(t, 0, got.TestCase.Result.ExitCode)
	assert.Equal(t, "", got.TestCase.Result.Stderr)
}

func Test_DockerExecutor_ExecuteCatchStderr(t *testing.T) {
	d := DockerExecutor{
		Image: "ubuntu:18.04",
	}

	cmd := `cd /invalid`

	test := TestCase{
		Command: CommandUnderTest{
			Cmd: cmd,
		},
		Expected: Expected{
			Stderr:   ExpectedOut{Exactly: "/bin/sh: 1: cd: can't cd to /invalid"},
			ExitCode: 2,
		},
	}

	got := d.Execute(test)
	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, 2, got.TestCase.Result.ExitCode)
	assert.Equal(t, "", got.TestCase.Result.Stdout)
	assert.Equal(t, "/bin/sh: 1: cd: can't cd to /invalid", got.TestCase.Result.Stderr)
}
