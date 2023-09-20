package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestSuiteFile = "/tmp/commander_test.yaml"

func TestRunApp(t *testing.T) {
	tests := []byte(`
tests:
    my title:
        command: echo hello
        exit-code: 0
`)
	err := os.WriteFile(TestSuiteFile, tests, 0o755)
	if err != nil {
		log.Fatal(err)
	}

	got := run([]string{"", "test", TestSuiteFile})
	assert.True(t, got)
}

func Test_AddCommand_ToFile(t *testing.T) {
	got := run([]string{"", "add", "--file", "/tmp/commander.yaml", "echo hello"})

	content, err := os.ReadFile("/tmp/commander.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "tests:\n  echo hello:\n    exit-code: 0\n    stdout: hello\n", string(content))
	assert.True(t, got)
}

func Test_AddCommand_ToStdout(t *testing.T) {
	got := run([]string{"", "add", "--stdout", "--no-file", "echo hello"})

	assert.True(t, got)
}

func Test_AddCommand_ToExistingFile(t *testing.T) {
	existingFile := "/tmp/existing.yaml"
	content := []byte(`
tests:
    echo existing:
        exit-code: 0
`)

	_ = os.WriteFile(existingFile, content, 0o755)

	got := run([]string{"", "add", "--stdout", "--file", existingFile, "echo hello"})

	content, err := os.ReadFile(existingFile)
	assert.Nil(t, err)
	assert.Equal(t, "tests:\n  echo existing:\n    exit-code: 0\n  echo hello:\n    exit-code: 0\n    stdout: hello\n", string(content))
	assert.True(t, got)
}
