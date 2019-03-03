package suite

// Constants for defining the various tested properties
const (
	ExitCode = "ExitCode"
	Stdout   = "Stdout"
	Stderr   = "Stderr"
)

// Result status codes
const (
	Success ResultStatus = iota
	Failed
	Skipped
)

// TestCase represents a test case which will be executed by the runtime
type TestCase struct {
	Title    string
	Command  CommandUnderTest
	Expected Expected
	Result   CommandResult
}

// ResultStatus represents the status code of a test result
type ResultStatus int

// CommandResult holds the result for a specific test
type CommandResult struct {
	Status            ResultStatus
	Stdout            string
	Stderr            string
	ExitCode          int
	FailureProperties []string
}

//Expected is the expected output of the command under test
type Expected struct {
	Stdout   ExpectedOut
	Stderr   ExpectedOut
	ExitCode int
}

type ExpectedOut struct {
	Contains    []string
	Line        map[int]string
	Exactly     string
}

// CommandUnderTest represents the command under test
type CommandUnderTest struct {
	Cmd string
	Env []string
	Dir string
}

type Suite struct {
	Tests []TestCase
}

func NewSuite(tests []TestCase) *Suite {
	return &Suite{
		Tests: tests,
	}
}
