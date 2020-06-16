package suite

import (
	"testing"

	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/stretchr/testify/assert"
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
	got := ParseYAML(yaml, "")
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
	tests := ParseYAML(yaml, "").GetTests()

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
	tests := ParseYAML(yaml, "").GetTests()

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
	tests := ParseYAML(yaml, "").GetTests()

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
            not-contains:
                - bonjour
            exactly: exactly hello
            json:
                $.object.attr: jsontest
`)
	tests := ParseYAML(yaml, "").GetTests()

	assert.Equal(t, "hello", tests[0].Expected.Stdout.Contains[0])
	assert.Equal(t, "exactly hello", tests[0].Expected.Stdout.Exactly)
	assert.Equal(t, "bonjour", tests[0].Expected.Stdout.NotContains[0])
	assert.Equal(t, "jsontest", tests[0].Expected.Stdout.JSON["$.object.attr"])
}

func TestYAMLConfig_UnmarshalYAML_ShouldConvertWithoutContains(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
        stderr:
            exactly: exactly stderr
`)
	tests := ParseYAML(yaml, "").GetTests()

	assert.Equal(t, "exactly stderr", tests[0].Expected.Stderr.Exactly)
	assert.IsType(t, runtime.ExpectedOut{}, tests[0].Expected.Stdout)
}

func Test_YAMLConfig_convertToExpectedOut(t *testing.T) {
	in := map[interface{}]interface{}{"exactly": "exactly stderr"}

	y := YAMLSuiteConf{}
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
	_ = ParseYAML(yaml, "")
}

func TestYAMLSuite_GetTestByTitle(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
`)
	test, err := ParseYAML(yaml, "").GetTestByTitle("echo hello")

	assert.Nil(t, err)
	assert.Equal(t, "echo hello", test.Title)
}

func TestYAMLSuite_GetTestByTitleShouldReturnError(t *testing.T) {
	yaml := []byte(`
tests:
    echo hello:
        exit-code: 0
`)
	_, err := ParseYAML(yaml, "").GetTestByTitle("does not exist")

	assert.Equal(t, "could not find test does not exist", err.Error())
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

	got := ParseYAML(yaml, "")
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
    interval: 500ms
    inherit-env: true

tests:
    echo hello:
       exit-code: 0
       config:
           env:
               KEY: local
           dir: /home/test
           timeout: 1s
           retries: 10
           interval: 5s
`)

	got := ParseYAML(yaml, "")

	// Assert global variables
	assert.Equal(t, map[string]string{"KEY": "global", "ANOTHER_KEY": "another_global"}, got.GetGlobalConfig().Env)
	assert.Equal(t, "/home/commander/", got.GetGlobalConfig().Dir)
	assert.Equal(t, "10ms", got.GetGlobalConfig().Timeout)
	assert.Equal(t, 2, got.GetGlobalConfig().Retries)
	assert.Equal(t, "500ms", got.GetGlobalConfig().Interval)
	assert.True(t, got.GetGlobalConfig().InheritEnv)

	// Assert local variables
	assert.Equal(t, map[string]string{"KEY": "local", "ANOTHER_KEY": "another_global"}, got.GetTests()[0].Command.Env)
	assert.Equal(t, "/home/test", got.GetTests()[0].Command.Dir)
	assert.Equal(t, "1s", got.GetTests()[0].Command.Timeout)
	assert.Equal(t, 10, got.GetTests()[0].Command.Retries)
	assert.Equal(t, "5s", got.GetTests()[0].Command.Interval)
	assert.True(t, got.GetTests()[0].Command.InheritEnv)
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

	_ = ParseYAML(yaml, "")
}

func TestYamlSuite_ShouldFailIfArrayIsGivenToExpectedOut(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			assert.Contains(t, r, "Failed to parse Stdout or Stderr with values: [yeah]")
		}
		assert.NotNil(t, r)
	}()

	yaml := []byte(`
tests:
    echo hello:
        stdout: 
          - yeah
`)

	_ = ParseYAML(yaml, "")
}

func Test_YAMLConfig_MarshalYAML(t *testing.T) {
	conf := YAMLSuiteConf{Tests: map[string]YAMLTest{
		"return_string": {
			Stdout: runtime.ExpectedOut{Contains: []string{"stdout string"}},
			Stderr: runtime.ExpectedOut{Contains: []string{"stderr string"}},
		},
		"return_struct": {
			Stdout: runtime.ExpectedOut{
				Contains:  []string{"stdout"},
				LineCount: 10,
			},
			Stderr: runtime.ExpectedOut{
				Contains:  []string{"stderr"},
				LineCount: 10,
			},
		},
		"return_nil": {
			Stdout: runtime.ExpectedOut{},
			Stderr: runtime.ExpectedOut{},
		},
	}}

	out, _ := conf.MarshalYAML()
	r := out.(YAMLSuiteConf)

	assert.Equal(t, "stdout string", r.Tests["return_string"].Stdout)
	assert.Equal(t, "stderr string", r.Tests["return_string"].Stderr)

	assert.Equal(t, conf.Tests["return_struct"].Stdout, r.Tests["return_struct"].Stdout)
	assert.Equal(t, conf.Tests["return_struct"].Stderr, r.Tests["return_struct"].Stderr)

	assert.Nil(t, r.Tests["return_nil"].Stdout)
	assert.Nil(t, r.Tests["return_nil"].Stderr)
}

func Test_convertExpectOut_ReturnNilIfEmpty(t *testing.T) {
	out := runtime.ExpectedOut{
		Contains: []string{""},
	}

	r := convertExpectedOut(out)

	assert.Nil(t, r)
}

func Test_convertExpectedOut_ReturnContainsAsString(t *testing.T) {
	out := runtime.ExpectedOut{
		Contains: []string{"test"},
	}

	r := convertExpectedOut(out)

	assert.Equal(t, "test", r)
}

func Test_convertExpectedOut_ReturnFullStruct(t *testing.T) {
	out := runtime.ExpectedOut{
		Contains:  []string{"hello", "hi"},
		LineCount: 10,
		Exactly:   "test",
	}

	r := convertExpectedOut(out)

	assert.Equal(t, out, r)
}

func TestYAMLSuite_should_parse_ssh(t *testing.T) {
	yaml := []byte(`
nodes:
   ssh-host1:
       type: ssh
       addr: localhost
       user: root
       pass: 12345!
       identity-file: ".ssh/id_rsa"
   docker-host:
       type: docker
       image: ubuntu:18.04
       privileged: true
tests:
   echo hello:
      config:
        nodes:
          - docker-host
          - ssh-host1
      exit-code: 0
`)

	got := ParseYAML(yaml, "")

	assert.Len(t, got.GetNodes(), 2)

	node, err := got.GetNodeByName("ssh-host1")
	assert.Nil(t, err)
	assert.Equal(t, "ssh-host1", node.Name)
	assert.Equal(t, "localhost", node.Addr)
	assert.Equal(t, "root", node.User)
	assert.Equal(t, "12345!", node.Pass)
	assert.Equal(t, "ssh", node.Type)
	assert.Equal(t, ".ssh/id_rsa", node.IdentityFile)

	dockerNode, err := got.GetNodeByName("docker-host")
	assert.Equal(t, "ubuntu:18.04", dockerNode.Image)
	assert.Equal(t, "docker", dockerNode.Type)
	assert.True(t, dockerNode.Privileged)

	assert.Contains(t, got.GetTests()[0].Nodes, "docker-host")
	assert.Contains(t, got.GetTests()[0].Nodes, "ssh-host1")
}
