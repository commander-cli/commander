package commander

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var ExampleYaml = []byte(`
tests:
    "should print something":
        command: echo something
        exit-code: 0
        stdout: something
    `)

func Test_ParseYamlFile(t *testing.T) {
	got := ParseYAMLFile("./_fixtures/commander.yaml")

	assert.Implements(t, (*Suite)(nil), got)
}

func Test_ParseYaml(t *testing.T) {
	got := ParseYAML(ExampleYaml)

	assert.Implements(t, (*Suite)(nil), got)
	assert.Len(t, got.GetTestCases(), 1)

	test := got.GetTestCases()[0]
	assert.Equal(t, "should print something", test.Title)
	assert.Equal(t, "echo something", test.Command)
	assert.Equal(t, "something", test.Stdout)
	assert.Equal(t, 0, test.ExitCode)
}

func Test_ParseMutlipleYamlTests(t *testing.T) {
	var multiExampleYaml = []byte(`
tests:
    "should print something":
        command: echo something
        exit-code: 0
    "should print hello":
        command: echo hello
        exit-code: 0
    `)

	got := ParseYAML(multiExampleYaml)

	assert.Len(t, got.GetTestCases(), 2)
}
