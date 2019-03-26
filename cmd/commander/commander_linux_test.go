package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

	got := testCommand(TestSuiteFile, "")
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

	got := testCommand(TestSuiteFile, "")
	assert.Equal(t, "Test suite failed", got.Error())

}
