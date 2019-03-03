package suite

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func Test_NewSuite(t *testing.T) {
    tests := []TestCase{}
    suite := NewSuite(tests)

    assert.False(t, suite.Executed)
}
