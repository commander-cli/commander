package app

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/commander-cli/commander/pkg/output"
	"github.com/commander-cli/commander/pkg/runtime"
	"github.com/commander-cli/commander/pkg/suite"
)

var out output.OutputWriter

// TestCommand executes the test argument
// testPath is the path to the test suite config, it can be a dir or file
// ctx holds the command flags. If directory scanning is enabled with --dir,
// test filtering is not supported
func TestCommand(testPath string, ctx TestCommandContext) error {
	if ctx.Verbose {
		log.SetOutput(os.Stdout)
	}

	out = output.NewCliOutput(!ctx.NoColor)

	if testPath == "" {
		testPath = CommanderFile
	}

	var result runtime.Result
	var err error
	switch {
	case ctx.Dir:
		fmt.Println("Starting test against directory: " + testPath + "...")
		fmt.Println("")
		result, err = testDir(testPath, ctx.Filters)
	case testPath == "-":
		fmt.Println("Starting test from stdin...")
		fmt.Println("")
		result, err = testStdin(ctx.Filters)
	case isURL(testPath):
		fmt.Println("Starting test from " + testPath + "...")
		fmt.Println("")
		result, err = testURL(testPath, ctx.Filters)
	default:
		fmt.Println("Starting test file " + testPath + "...")
		fmt.Println("")
		result, err = testFile(testPath, "", ctx.Filters)
	}

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if !out.PrintSummary(result) && !ctx.Verbose {
		return fmt.Errorf("Test suite failed, use --verbose for more detailed output")
	}

	return nil
}

func testFile(filePath string, fileName string, filters runtime.Filters) (runtime.Result, error) {
	s, err := readFile(filePath, fileName)
	if err != nil {
		return runtime.Result{}, fmt.Errorf("Error " + err.Error())
	}

	return execute(s, filters)
}

func testDir(directory string, filters runtime.Filters) (runtime.Result, error) {
	result := runtime.Result{}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return result, fmt.Errorf("Error: Input is not a directory")
	}

	for _, f := range files {
		if f.IsDir() {
			continue // skip dirs
		}

		p := path.Join(directory, f.Name())
		newResult, err := testFile(p, f.Name(), filters)
		if err != nil {
			return result, err
		}

		result = convergeResults(result, newResult)
	}

	return result, nil
}

func testURL(url string, filters runtime.Filters) (runtime.Result, error) {
	resp, err := http.Get(url)
	if err != nil {
		return runtime.Result{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return runtime.Result{}, err
	}

	s := suite.ParseYAML(body, "")

	return execute(s, filters)
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func convergeResults(result runtime.Result, new runtime.Result) runtime.Result {
	result.TestResults = append(result.TestResults, new.TestResults...)
	result.Failed += new.Failed
	result.Duration += new.Duration

	return result
}

func testStdin(filters runtime.Filters) (runtime.Result, error) {
	f, err := os.Stdin.Stat()
	if err != nil {
		return runtime.Result{}, err
	}

	if f.Mode()&os.ModeCharDevice != 0 {
		return runtime.Result{}, fmt.Errorf("Error: command is intended to work with pipes")
	}

	r := bufio.NewReader(os.Stdin)
	content, err := ioutil.ReadAll(r)
	s := suite.ParseYAML(content, "")

	return execute(s, filters)
}

func execute(s suite.Suite, filters runtime.Filters) (runtime.Result, error) {
	tests := s.GetTests()
	if len(filters) != 0 {
		tests = []runtime.TestCase{}
	}

	// Filter tests if test filters was given
	for _, f := range filters {
		t, err := s.FindTests(f)
		if err != nil {
			return runtime.Result{}, err
		}
		tests = append(tests, t...)
	}

	r := runtime.NewRuntime(out.GetEventHandler(), s.Nodes...)
	result := r.Start(tests)

	return result, nil
}

func readFile(filePath string, fileName string) (suite.Suite, error) {
	s := suite.Suite{}

	f, err := os.Stat(filePath)
	if err != nil {
		return s, fmt.Errorf("open %s: no such file or directory", filePath)
	}

	if f.IsDir() {
		return s, fmt.Errorf("%s: is a directory\nUse --dir to test directories with multiple test files", filePath)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return s, err
	}

	s = suite.ParseYAML(content, fileName)

	return s, nil
}
