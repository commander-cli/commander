package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"testing"
)

// CaptureStandardOutput allows to capture the output which will be written
// to os.Stdout and os.Stderr.
// It returns the captured output and the return value of the called function
func CaptureStandardOutput(f func() interface{}) (string, interface{}) {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	result := f()
	writer.Close()
	return <-out, result
}

func assertEqualWithLineBreak(t *testing.T, expected string, actual string) {
	if runtime.GOOS == "windows" {
		expected = expected + "\r\n"
	} else {
		expected = expected + "\n"
	}

	assert.Equal(t, expected, actual)
}
