package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddCommand_ToFile(t *testing.T) {
	tmpFile := os.Getenv("Temp") + "\\commander.yaml"

	got := run([]string{"", "add", "--file", tmpFile, "echo hello"})

	content, err := os.ReadFile(tmpFile)
	assert.Nil(t, err)
	assert.Equal(t, "tests:\n  echo hello:\n    exit-code: 0\n    stdout: hello\n", string(content))
	assert.True(t, got)
}
