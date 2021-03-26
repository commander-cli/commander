package output

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOutput(t *testing.T) {
	var o Output
	o, err := NewOutput(CLI, false)
	assert.NoError(t, err)
	assert.Implements(t, (*Output)(nil), o)
}
