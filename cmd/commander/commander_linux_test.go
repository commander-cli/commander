package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

const TestSuiteFile = "/tmp/commander_test.yaml"

func Test_CommanderFile(t *testing.T) {
	tests := []byte(`
tests:
    echo hello:
        exit-code: 0
`)
	err := ioutil.WriteFile(TestSuiteFile, tests, 0755)

	assert.Nil(t, err)

	got := testCommand(TestSuiteFile, "", CommanderContext{})
	assert.Nil(t, got)
}

func Test_FailingSuite(t *testing.T) {
	tests := []byte(`
tests:
    echo hello:
        exit-code: 1
`)
	err := ioutil.WriteFile(TestSuiteFile, tests, 0755)

	assert.Nil(t, err)

	got := testCommand(TestSuiteFile, "", CommanderContext{})
	assert.Equal(t, "Test suite failed, use --verbose for more detailed output", got.Error())

}

func Test_WithTitle(t *testing.T) {
	tests := []byte(`
tests:
    my title:
        command: echo hello
        exit-code: 0
    another:
        command: echo another
        exit-code: 1
`)
	err := ioutil.WriteFile(TestSuiteFile, tests, 0755)

	assert.Nil(t, err)

	got := testCommand(TestSuiteFile, "my title", CommanderContext{})
	assert.Nil(t, got)
}

func TestRunApp(t *testing.T) {
	tests := []byte(`
tests:
    my title:
        command: echo hello
        exit-code: 0
`)
	err := ioutil.WriteFile(TestSuiteFile, tests, 0755)
	if err != nil {
		log.Fatal(err)
	}

	got := run([]string{"", "test", TestSuiteFile})
	assert.True(t, got)
}
