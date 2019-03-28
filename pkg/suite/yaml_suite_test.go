package suite

import (
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

const ExpectedLineCount = 10

func TestYAMLConfig_UnmarshalYAML(t *testing.T) {
	yaml := []byte(`
tests:
    it should print hello:
        command: echo hello
        exit-code: 0
        stdout: hello
        stderr: anything
`)
	got := ParseYAML(yaml)
	tests := got.GetTests()

	assert.Len(t, tests, 1)
	assert.Equal(t, "echo hello", tests[0].Command.Cmd)
	assert.Equal(t, 0, tests[0].Expected.ExitCode)
	assert.Equal(t, "it should print hello", tests[0].Title)
	assert.Equal(t, "hello", tests[0].Expected.Stdout.Contains[0])
	assert.Equal(t, "anything", tests[0].Expected.Stderr.Contains[0])
}

func TestYAMLConfig_UnmarshalYAML_ShouldUseTitleAsCommand(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
        stdout: hello
        stderr: anything
`)
	tests := ParseYAML(yaml).GetTests()

	assert.Equal(t, "echo hello", tests[0].Command.Cmd)
	assert.Equal(t, "echo hello", tests[0].Title)
}

func TestYAMLConfig_UnmarshalYAML_ShouldParseLineCount(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
        stdout:
            line-count: 10
`)
	tests := ParseYAML(yaml).GetTests()

	assert.Equal(t, ExpectedLineCount, tests[0].Expected.Stdout.LineCount)
}

func TestYAMLConfig_UnmarshalYAML_ShouldParseLines(t *testing.T) {
	yaml := []byte(`
tests:
    printf "line1\nline2\nline3\nline4":
        exit-code: 0
        stdout:
            lines:
                0: line1
                1: line2
                3: line4
`)
	tests := ParseYAML(yaml).GetTests()

	assert.Equal(t, "line1", tests[0].Expected.Stdout.Lines[0])
	assert.Equal(t, "line2", tests[0].Expected.Stdout.Lines[1])
	assert.Equal(t, "line4", tests[0].Expected.Stdout.Lines[3])
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
	tests := ParseYAML(yaml).GetTests()

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
	tests := ParseYAML(yaml).GetTests()

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

func TestYAMLConfig_UnmarshalYAML_ShouldPanicIfKeyDoesNotExist(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Unknown keys should not be parsed, the program should panic")
		}
	}()

	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
        stderr:
            typo: exactly stderr
`)
	_ = ParseYAML(yaml)
}

func TestYAMLSuite_GetTestByTitle(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
`)
	test, err := ParseYAML(yaml).GetTestByTitle("echo hello")

	assert.Nil(t, err)
	assert.Equal(t, "echo hello", test.Title)
}

func TestYAMLSuite_GetTestByTitleShouldReturnError(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
`)
	_, err := ParseYAML(yaml).GetTestByTitle("does not exist")

	assert.Equal(t, "Could not find test does not exist", err.Error())
}

func TestYAMLSuite_ShouldParseGlobalConfig(t *testing.T) {
	yaml := []byte(`
config:
    env:
        KEY: value
    dir: /home/commander/
tests:
    echo hello:
       exit-code: 0
`)

	got := ParseYAML(yaml)
	assert.Equal(t, map[string]string{"KEY": "value"}, got.GetGlobalConfig().Env)
	assert.Equal(t, map[string]string{"KEY": "value"}, got.GetTests()[0].Command.Env)
	assert.Equal(t, "/home/commander/", got.GetTests()[0].Command.Dir)
	assert.Equal(t, "/home/commander/", got.GetGlobalConfig().Dir)
}

func TestYAMLSuite_ShouldPreferLocalTestConfigs(t *testing.T) {
	yaml := []byte(`
config:
    env:
        KEY: global
        ANOTHER_KEY: another_global
    dir: /home/commander/
    timeout: 10ms
    retries: 2

tests:
    echo hello:
       exit-code: 0
       config:
           env:
               KEY: local
           dir: /home/test
           timeout: 1s
           retries: 10
`)

	got := ParseYAML(yaml)
	assert.Equal(t, map[string]string{"KEY": "global", "ANOTHER_KEY": "another_global"}, got.GetGlobalConfig().Env)
	assert.Equal(t, "/home/commander/", got.GetGlobalConfig().Dir)
	assert.Equal(t, "10ms", got.GetGlobalConfig().Timeout)
	assert.Equal(t, 2, got.GetGlobalConfig().Retries)

	assert.Equal(t, map[string]string{"KEY": "local", "ANOTHER_KEY": "another_global"}, got.GetTests()[0].Command.Env)
	assert.Equal(t, "/home/test", got.GetTests()[0].Command.Dir)
	assert.Equal(t, "1s", got.GetTests()[0].Command.Timeout)
	assert.Equal(t, 10, got.GetTests()[0].Command.Retries)
}

func TestYAMLSuite_ShouldThrowAnErrorIfFieldIsNotRegistered(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			assert.Contains(t, r, "field stdot not found in type suite.YAMLTest")
		}
		assert.NotNil(t, r)
	}()

	yaml := []byte(`
tests:
    echo hello:
        stdot: yeah
`)

	_ = ParseYAML(yaml)
}
