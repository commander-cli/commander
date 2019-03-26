package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_CommanderFile(t *testing.T) {
	tests := []byte(`
tests:
    echo hello:
        exit-code: 0
`)
	err := ioutil.WriteFile("/tmp/commander_test.yaml", tests, 0755)

	assert.Nil(t, err)

	got := testCommand("/tmp/commander_test.yaml", "")
	assert.Nil(t, got)
}
