package app

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	commanderRuntime "github.com/commander-cli/commander/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func Test_TestCommand_Verbose(t *testing.T) {
	out := captureOutput(func() {
		TestCommand("commander.yaml", TestCommandContext{Verbose: true})
		log.Println("test test test")
	})

	assert.Contains(t, out, "test test test")
}

func Test_TestCommand_DefaultFile(t *testing.T) {
	err := TestCommand("", TestCommandContext{Verbose: true})
	assert.Contains(t, err.Error(), "commander.yaml")
}

func Test_TestCommand(t *testing.T) {
	err := TestCommand("commander.yaml", TestCommandContext{})

	if runtime.GOOS == "windows" {
		assert.Contains(t, err.Error(), "Error open commander.yaml:")
	} else {
		assert.Equal(t, "Error open commander.yaml: no such file or directory", err.Error())
	}
}

func Test_TestCommand_ShouldUseCustomFile(t *testing.T) {
	err := TestCommand("my-test.yaml", TestCommandContext{})

	if runtime.GOOS == "windows" {
		assert.Contains(t, err.Error(), "Error open my-test.yaml:")
	} else {
		assert.Equal(t, "Error open my-test.yaml: no such file or directory", err.Error())
	}
}

func Test_TestCommand_File_WithDir(t *testing.T) {
	err := TestCommand("../../examples", TestCommandContext{})

	if runtime.GOOS == "windows" {
		assert.Contains(t, err.Error(), "is a directory")
	} else {
		assert.Equal(t, "Error ../../examples: is a directory\nUse --dir to test directories with multiple test files", err.Error())
	}
}

func Test_TestCommand_Dir(t *testing.T) {
	err := TestCommand("../../examples", TestCommandContext{Dir: true})

	if runtime.GOOS == "windows" {
		assert.Contains(t, err.Error(), "Test suite failed, use --verbose for more detailed output")
	} else {
		assert.Equal(t, "Test suite failed, use --verbose for more detailed output", err.Error())
	}
}

func Test_TestCommand_Dir_Err(t *testing.T) {
	err := TestCommand("http://foo.com/bar", TestCommandContext{Dir: true})

	if runtime.GOOS == "windows" {
		assert.Contains(t, err.Error(), "Error: Input is not a directory")
	} else {
		assert.Equal(t, "Error: Input is not a directory", err.Error())
	}
}

func Test_TestCommand_Http(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Get("/bar").
		Reply(200).
		BodyString(`
tests:
  echo hello:
    exit-code: 0
`)

	out := captureOutput(func() {
		TestCommand("http://foo.com/bar", TestCommandContext{Verbose: true})
	})

	assert.Contains(t, out, "✓ [local] echo hello")
}

func Test_TestCommand_Http_Err(t *testing.T) {
	err := TestCommand("http://error/not/a/url", TestCommandContext{Dir: false})

	assert.NotNil(t, err)
}

func Test_TestCommand_StdIn(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) //clean up

	content := []byte(`
tests:
   echo hello:
     exit-code: 0
`)

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// set stdin to tempfile
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	os.Stdin = tmpfile

	out := captureOutput(func() {
		TestCommand("-", TestCommandContext{Verbose: false})
	})

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	assert.Contains(t, out, "✓ [local] echo hello")
}

func Test_TestCommand_StdIn_Err(t *testing.T) {
	// set stdin to nil
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	os.Stdin = nil

	err := TestCommand("-", TestCommandContext{Verbose: false})

	assert.NotNil(t, err)
}

func Test_ConvergeResults(t *testing.T) {
	duration, _ := time.ParseDuration("5s")

	result1 := commanderRuntime.Result{
		TestResults: []commanderRuntime.TestResult{},
		Duration:    duration,
		Failed:      1,
	}

	result2 := commanderRuntime.Result{
		TestResults: []commanderRuntime.TestResult{},
		Duration:    duration,
		Failed:      0,
	}

	actual := convergeResults(result1, result2)

	expectDur, _ := time.ParseDuration("10s")
	assert.Equal(t, expectDur, actual.Duration)
	assert.Equal(t, 1, actual.Failed)
}

func captureOutput(f func()) string {
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
	f()
	writer.Close()
	return <-out
}
