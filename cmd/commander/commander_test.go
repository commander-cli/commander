package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateCliApp(t *testing.T) {
	app := createCliApp()

	assert.Equal(t, "Commander", app.Name)
	assert.Equal(t, "test", app.Commands[0].Name)
}
