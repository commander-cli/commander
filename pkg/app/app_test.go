package app

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"testing"
)

func TestNewAddCommandContextFromCli(t *testing.T) {
	set := flag.NewFlagSet("verbose", 0)
	set.Bool("verbose", true, "")
	set.Bool("no-color", true, "")
	set.Int("concurrent", 5, "")

	context := &cli.Context{}
	ctx := cli.NewContext(nil, set, context)

	r := NewAddContextFromCli(ctx)

	assert.True(t, r.Verbose)
	assert.True(t, r.NoColor)
	assert.Equal(t, 5, r.Concurrent)
}
