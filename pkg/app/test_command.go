package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/SimonBaeumer/commander/pkg/output"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
)

// TestCommand executes the test argument
// testPath is the path to the test suite config, it can be a dir or file
// ctx holds the command flags. If directory scanning is enabled with --dir it is
// not supported to filter tests, therefore testFilterTitle is an empty string
func TestCommand(testPath string, ctx AddCommandContext) error {
	if ctx.Verbose == true {
		log.SetOutput(os.Stdout)
	}

	if testPath == "" {
		testPath = CommanderFile
	}

	var results <-chan runtime.TestResult
	var err error
	if ctx.Dir {
		fmt.Println("Starting test against directory: " + testPath + "...")
		fmt.Println("")
		results, err = testDir(testPath, ctx.Filters)
	} else {
		fmt.Println("Starting test file " + testPath + "...")
		fmt.Println("")
		results, err = testFile(testPath, ctx.Filters)
	}

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	out := output.NewCliOutput(!ctx.NoColor)
	if !out.Start(results) {
		return fmt.Errorf("Test suite failed, use --verbose for more detailed output")
	}

	return nil
}

func testDir(directory string, filters runtime.Filters) (<-chan runtime.TestResult, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	results := make(chan runtime.TestResult)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, f := range files {
			// Skip reading dirs for now. Should we also check valid file types?
			if f.IsDir() {
				continue
			}

			fileResults, err := testFile(path.Join(directory, f.Name()), filters)
			if err != nil {
				panic(fmt.Sprintf("%s: %s", f.Name(), err))
			}

			for r := range fileResults {
				r.FileName = f.Name()
				results <- r
			}
		}
	}()

	go func(ch chan runtime.TestResult) {
		wg.Wait()
		close(results)
	}(results)

	return results, nil
}

func testFile(filePath string, filters runtime.Filters) (<-chan runtime.TestResult, error) {
	content, err := readFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error " + err.Error())
	}

	var s suite.Suite
	s = suite.ParseYAML(content)

	tests := s.GetTests()
	if len(filters) != 0 {
		tests = []runtime.TestCase{}
	}

	// Filter tests if test filters was given
	for _, f := range filters {
		t, err := s.GetTestByTitle(f)
		if err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}

	r := runtime.NewRuntime(s.Nodes...)
	results := r.Start(tests)

	return results, nil
}

func readFile(filePath string) ([]byte, error) {
	f, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("open %s: no such file or directory", filePath)
	}

	if f.IsDir() {
		return nil, fmt.Errorf("%s: is a directory\nUse --dir to test directories with multiple test files", filePath)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}
