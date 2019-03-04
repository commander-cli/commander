package config

import (
    "github.com/SimonBaeumer/commander/pkg/runtime"
    "github.com/stretchr/testify/assert"
	"testing"
)

func TestYAMLConfig_UnmarshalYAML(t *testing.T) {
	yaml := []byte(`
tests:
    it should print hello:
        command: echo hello
        exit-code: 0
        stdout: hello
        stderr: anything
`)
	tests := ParseYAML(yaml)

	assert.Len(t, tests, 1)
	assert.Equal(t, "echo hello", tests[0].Command.Cmd)
	assert.Equal(t, 0, tests[0].Expected.ExitCode)
	assert.Equal(t, "it should print hello", tests[0].Title)
	assert.Equal(t, "hello", tests[0].Expected.Stdout.Exactly)
	assert.Equal(t, "anything", tests[0].Expected.Stderr.Exactly)
}

func TestYAMLConfig_UnmarshalYAML_ShouldUseTitleAsCommand(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
        stdout: hello
        stderr: anything
`)
	tests := ParseYAML(yaml)

	assert.Equal(t, "echo hello", tests[0].Command.Cmd)
	assert.Equal(t, "echo hello", tests[0].Title)
}

func TestYAMLConfig_UnmarshalYAML_ShouldConvertStdoutToExpectedOut(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
        stdout:
            contains:
                - hello
                - another hello
            exactly: exactly hello
`)
	tests := ParseYAML(yaml)

	assert.Equal(t, "hello", tests[0].Expected.Stdout.Contains[0])
	assert.Equal(t, "exactly hello", tests[0].Expected.Stdout.Exactly)
}

func TestYAMLConfig_UnmarshalYAML_ShouldConvertWithoutContains(t *testing.T) {
    yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
        stderr:
            exactly: exactly stderr
`)
    tests := ParseYAML(yaml)

    assert.Equal(t, "exactly stderr", tests[0].Expected.Stderr.Exactly)
    assert.IsType(t, runtime.ExpectedOut{}, tests[0].Expected.Stdout)
}

func Test_YAMLConfig_convertToExpectedOut(t *testing.T) {
    in := map[interface{}]interface{}{"exactly": "exactly stderr"}

    y := YAMLConfig{}
    got := y.convertToExpectedOut(in)

    assert.IsType(t, runtime.ExpectedOut{}, got)
    assert.Equal(t, "exactly stderr", got.Exactly)
}