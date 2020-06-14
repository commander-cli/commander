package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCliOutput(t *testing.T) {
	got := NewCliOutput(true)
	assert.IsType(t, OutputWriter{}, got)
}
