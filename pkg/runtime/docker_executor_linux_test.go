package runtime

import (
	"log"
	"os"
	"testing"

	"github.com/commander-cli/commander/pkg/suite"
	"github.com/stretchr/testify/assert"
)

func isEnabled() bool {
	v := os.Getenv("COMMANDER_TEST_ALL")
	if v != "1" {
		log.Println("Skip test, set env COMMANDER_TEST_ALL to 1")
		return false
	}
	return true
}

func Test_DockerExecutor_Execute(t *testing.T) {
	if !isEnabled() {
		return
	}

	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
	}

	cmd := `echo hello`

	test := suite.TestCase{
		Command: suite.CommandUnderTest{
			Cmd: cmd,
		},
		Expected: suite.Expected{
			Stdout:   suite.ExpectedOut{Exactly: "hello"},
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
	if !isEnabled() {
		return
	}

	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
	}

	cmd := `cd /invalid`

	test := suite.TestCase{
		Command: suite.CommandUnderTest{
			Cmd: cmd,
		},
		Expected: suite.Expected{
			Stderr:   suite.ExpectedOut{Exactly: "/bin/sh: 1: cd: can't cd to /invalid"},
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
	if !isEnabled() {
		return
	}

	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
	}

	test := suite.TestCase{
		Command: suite.CommandUnderTest{
			Cmd: "pwd",
			Dir: "/tmp",
		},
		Expected: suite.Expected{
			Stdout:   suite.ExpectedOut{Exactly: "/tmp"},
			ExitCode: 0,
		},
	}

	got := d.Execute(test)
	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, "/tmp", got.TestCase.Result.Stdout)
	assert.Nil(t, got.TestCase.Result.Error)
}

func Test_DockerExecutor_Execute_Env(t *testing.T) {
	if !isEnabled() {
		return
	}

	d := DockerExecutor{
		Image: "docker.io/library/ubuntu:18.04",
	}

	test := suite.TestCase{
		Command: suite.CommandUnderTest{
			Cmd: "echo $ENV_KEY",
			Env: map[string]string{
				"ENV_KEY": "env-value",
			},
		},
		Expected: suite.Expected{
			Stdout:   suite.ExpectedOut{Exactly: "env-value"},
			ExitCode: 0,
		},
	}

	got := d.Execute(test)
	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, "env-value", got.TestCase.Result.Stdout)
	assert.Nil(t, got.TestCase.Result.Error)
}
