package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DockerExecutor_Execute(t *testing.T) {
	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
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
	assert.Nil(t, got.TestCase.Result.Error)
}

func Test_DockerExecutor_ExecuteCatchStderr(t *testing.T) {
	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
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
	assert.Nil(t, got.TestCase.Result.Error)
}

func Test_DockerExecutor_Execute_Dir(t *testing.T) {
	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
	}

	test := TestCase{
		Command: CommandUnderTest{
			Cmd: "pwd",
			Dir: "/tmp",
		},
		Expected: Expected{
			Stdout:   ExpectedOut{Exactly: "/tmp"},
			ExitCode: 0,
		},
	}

	got := d.Execute(test)
	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, "/tmp", got.TestCase.Result.Stdout)
	assert.Nil(t, got.TestCase.Result.Error)
}

func Test_DockerExecutor_Execute_Env(t *testing.T) {
	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
	}

	test := TestCase{
		Command: CommandUnderTest{
			Cmd: "echo $ENV_KEY",
			Env: map[string]string{
				"ENV_KEY": "env-value",
			},
		},
		Expected: Expected{
			Stdout:   ExpectedOut{Exactly: "env-value"},
			ExitCode: 0,
		},
	}

	got := d.Execute(test)
	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, "env-value", got.TestCase.Result.Stdout)
	assert.Nil(t, got.TestCase.Result.Error)
}
