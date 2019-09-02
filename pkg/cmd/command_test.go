package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestCommand_NewCommand(t *testing.T) {
	cmd := NewCommand("")
	assert.False(t, cmd.Executed())
}

func TestCommand_Execute(t *testing.T) {
	cmd := NewCommand("echo hello")

	err := cmd.Execute()

	assert.Nil(t, err)
	assert.True(t, cmd.Executed())
	assert.Equal(t, cmd.Stdout(), "hello")
}

func TestCommand_ExitCode(t *testing.T) {
	cmd := NewCommand("exit 120")

	err := cmd.Execute()

	assert.Nil(t, err)
	assert.Equal(t, 120, cmd.ExitCode())
}

func TestCommand_WithEnvVariables(t *testing.T) {
	envVar := "$TEST"
	if runtime.GOOS == "windows" {
		envVar = "%TEST%"
	}
	cmd := NewCommand(fmt.Sprintf("echo %s", envVar))
	cmd.Env = []string{"TEST=hey"}

	_ = cmd.Execute()

	assert.Equal(t, cmd.Stdout(), "hey")
}

func TestCommand_Executed(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			assert.Contains(t, r, "Can not read Stdout if command was not executed")
		}
		assert.NotNil(t, r)
	}()

	c := NewCommand("echo will not be executed")
	_ = c.Stdout()
}

func TestCommand_AddEnv(t *testing.T) {
	c := NewCommand("echo test")
	c.AddEnv("key", "value")
	assert.Equal(t, []string{"key=value"}, c.Env)
}

func TestCommand_AddEnvWithShellVariable(t *testing.T) {
	const TestEnvKey = "COMMANDER_TEST_SOME_KEY"
	os.Setenv(TestEnvKey, "test from shell")
	defer os.Unsetenv(TestEnvKey)

	c := NewCommand("echo $SOME_KEY")
	c.AddEnv("SOME_KEY", fmt.Sprintf("${%s}", TestEnvKey))

	err := c.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "test from shell", c.Stdout())
}

func TestCommand_AddMultipleEnvWithShellVariable(t *testing.T) {
	const TestEnvKeyPlanet = "COMMANDER_TEST_PLANET"
	const TestEnvKeyName = "COMMANDER_TEST_NAME"
	os.Setenv(TestEnvKeyPlanet, "world")
	os.Setenv(TestEnvKeyName, "Simon")
	defer func() {
		os.Unsetenv(TestEnvKeyPlanet)
		os.Unsetenv(TestEnvKeyName)
	}()

	c := NewCommand("echo $SOME_KEY")
	envValue := fmt.Sprintf("Hello ${%s}, I am ${%s}", TestEnvKeyPlanet, TestEnvKeyName)
	c.AddEnv("SOME_KEY", envValue)

	err := c.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "Hello world, I am Simon", c.Stdout())
}

func TestCommand_SetTimeoutMS_DefaultTimeout(t *testing.T) {
	c := NewCommand("echo test")
	c.SetTimeoutMS(0)
	assert.Equal(t, (1 * time.Minute), c.Timeout)
}

func TestCommand_SetTimeoutMS(t *testing.T) {
	c := NewCommand("echo test")
	c.SetTimeoutMS(100)
	assert.Equal(t, 100*time.Millisecond, c.Timeout)
}

func TestCommand_SetTimeout(t *testing.T) {
	c := NewCommand("echo test")
	_ = c.SetTimeout("100s")
	duration, _ := time.ParseDuration("100s")
	assert.Equal(t, duration, c.Timeout)
}

func TestCommand_SetInvalidTimeout(t *testing.T) {
	c := NewCommand("echo test")
	err := c.SetTimeout("1")
	assert.Equal(t, "time: missing unit in duration 1", err.Error())
}
