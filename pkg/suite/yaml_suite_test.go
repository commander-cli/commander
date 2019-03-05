package suite

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
	got := ParseYAML(yaml)
	tests := got.GetTests()

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
	tests := ParseYAML(yaml).GetTests()


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