package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/SimonBaeumer/commander/pkg/output"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
)

// TestCommand executes the test argument
// testPath is the path to the configuration object, an object can be a dir or file
// titleFilterTitle is the title of test which should be executed, if empty it will execute all tests
// ctx holds the command flags
// when --dir is enabled testFilterPath must be a zero value
func TestCommand(testPath string, testFilterTitle string, ctx AddCommandContext) error {
	if ctx.Verbose == true {
		log.SetOutput(os.Stdout)
	}

	if testPath == "" {
		testPath = CommanderFile
	}

	var results <-chan runtime.TestResult
	var err error
	if ctx.Dir {
		if testFilterTitle != "" {
			return fmt.Errorf("Test may not be filtered when --dir is enabled")
		}
		fmt.Println("Starting test against directory: " + testPath + "...")
		fmt.Println("")
		results, err = testDir(testPath)
	} else {
		fmt.Println("Starting test file " + testPath + "...")
		fmt.Println("")
		results, err = testFile(testPath, testFilterTitle)
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

func testDir(directory string) (<-chan runtime.TestResult, error) {
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

			fileResults, err := testFile(directory+"/"+f.Name(), "")
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

func testFile(filePath string, title string) (<-chan runtime.TestResult, error) {
	content, err := readFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error " + err.Error())
	}

	var s suite.Suite
	s = suite.ParseYAML(content)
	tests := s.GetTests()
	// Filter tests if test title was given
	if title != "" {
		test, err := s.GetTestByTitle(title)
		if err != nil {
			return nil, err
		}
		tests = []runtime.TestCase{test}
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
