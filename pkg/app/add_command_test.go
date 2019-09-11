package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AddCommand(t *testing.T) {
	content, err := AddCommand("echo hello", []byte{})

	assert.Nil(t, err)
	assert.Equal(t, "tests:\n  echo hello:\n    exit-code: 0\n    stdout: hello\n", string(content))
}

func Test_AddCommand_AddToExisting(t *testing.T) {
	existing := []byte(`
tests:
    echo exists:
        exit-code: 0
`)

	content, err := AddCommand("echo hello", existing)

	assert.Nil(t, err)
	assert.Equal(t, "tests:\n  echo exists:\n    exit-code: 0\n  echo hello:\n    exit-code: 0\n    stdout: hello\n", string(content))
}

func Test_AddCommand_AddToExistingWithComplexStdoutStructure(t *testing.T) {
	existing := []byte(`
tests:
  exists:
    command: echo exists
    stdout:
      contains:
        - exists
      not-contains:
        - byebye
    exit-code: 0
`)

	content, err := AddCommand("echo hello", existing)

	expected := []byte(`tests:
  echo hello:
    exit-code: 0
    stdout: hello
  exists:
    command: echo exists
    exit-code: 0
    stdout:
      contains:
      - exists
      not-contains:
      - byebye
`)

	assert.Nil(t, err)
	assert.Equal(t, string(expected), string(content))
}
