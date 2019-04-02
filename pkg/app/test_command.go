package app

import (
	"fmt"
	"github.com/SimonBaeumer/commander/pkg/output"
	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/SimonBaeumer/commander/pkg/suite"
	"io/ioutil"
	"log"
	"os"
)

// TestCommand executes the test argument
// file is the path to the configuration file
// title ist the title of test which should be executed, if empty it will execute all tests
// ctx holds the command flags
func TestCommand(file string, title string, ctx AddCommandContext) error {
	if ctx.Verbose == true {
		log.SetOutput(os.Stdout)
	}

	if file == "" {
		file = CommanderFile
	}

	fmt.Println("Starting test file " + file + "...")
	fmt.Println("")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Error " + err.Error())
	}

	var s suite.Suite
	s = suite.ParseYAML(content)
	tests := s.GetTests()
	// Filter tests if test title was given
	if title != "" {
		test, err := s.GetTestByTitle(title)
		if err != nil {
			return err
		}
		tests = []runtime.TestCase{test}
	}

	results := runtime.Start(tests, ctx.Concurrent)
	out := output.NewCliOutput(!ctx.NoColor)
	if !out.Start(results) {
		return fmt.Errorf("Test suite failed, use --verbose for more detailed output")
	}

	return nil
}
