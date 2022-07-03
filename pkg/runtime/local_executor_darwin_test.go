package runtime

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntime_WithInheritFromShell(t *testing.T) {
	os.Setenv("TEST_COMMANDER", "test")
	defer os.Unsetenv("TEST_COMMANDER")

	test := TestCase{
		Command: CommandUnderTest{
			Cmd:        "echo $TEST_COMMANDER",
			InheritEnv: true,
		},
	}

	e := LocalExecutor{}
	got, err := e.Execute(test)
	assert.NoError(t, err)

	assert.Equal(t, "test", got.TestCase.Result.Stdout)
}

func TestRuntime_WithInheritFromShell_Overwrite(t *testing.T) {
	os.Setenv("TEST_COMMANDER", "test")
	os.Setenv("ANOTHER_ENV", "from-parent")
	defer func() {
		os.Unsetenv("TEST_COMMANDER")
		os.Unsetenv("ANOTHER_ENV")
	}()

	test := TestCase{
		Command: CommandUnderTest{
			Cmd:        "echo $TEST_COMMANDER $ANOTHER_ENV",
			InheritEnv: true,
			Env:        map[string]string{"TEST_COMMANDER": "overwrite"},
		},
	}

	e := LocalExecutor{}
	got, err := e.Execute(test)
	assert.NoError(t, err)

	assert.Equal(t, "overwrite from-parent", got.TestCase.Result.Stdout)
}
