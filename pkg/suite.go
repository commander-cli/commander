package commander

const (
    ExitCode = "ExitCode"
    Stdout = "Stdout"
    Stderr = "Stderr"
)

// TestCase represents a test case which will be executed by the runtime
type TestCase struct {
    Title    string
    Command  string
    Stdout   string
    Stderr   string
    ExitCode int
    Result   TestResult
}

// TestResult holds the result for a specific test
type TestResult struct {
    Success         bool
    //Skipped  bool
    Stdout          string
    Stderr          string
    ExitCode        int
    FailureProperty string
}

// Suite
type Suite interface {
    GetTestCases() []TestCase
}
