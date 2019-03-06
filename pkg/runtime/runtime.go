package runtime

import (
    "github.com/SimonBaeumer/commander/pkg/cmd"
    "log"
    "sync"
)


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

// CommandResult represents the TestCase and the ValidationResult
type TestResult struct {
    TestCase         TestCase
    ValidationResult ValidationResult
    FailedProperty   string
}

// Start starts the given test suite
func Start(tests []TestCase) <-chan TestResult {
    in := make(chan TestCase)
    out := make(chan TestResult)

     go func(tests []TestCase) {
         defer close(in)
         for _, t := range tests {
             in <- t
         }
     }(tests)

    //TODO: Add more concurrency
    workerCount := 1
    var wg sync.WaitGroup
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for t := range in {
                out <- runTest(t)
            }

        }()
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

func runTest(test TestCase) TestResult {
    // cut = command under test
    cut := cmd.NewCommand(test.Command.Cmd)

    if err := cut.Execute(); err != nil {
        log.Fatal(err)
    }

    // Write test result
    test.Result = CommandResult{
        ExitCode: cut.ExitCode(),
        Stdout:   cut.Stdout(),
        Stderr:   cut.Stderr(),
    }

    validationResult := validateExpectedOut(test.Result.Stdout, test.Expected.Stdout)
    if !validationResult.Success {
        return TestResult{
            ValidationResult: validationResult,
            TestCase:         test,
            FailedProperty:   Stdout,
        }
    }

    validationResult = validateExpectedOut(test.Result.Stderr, test.Expected.Stderr)
    if !validationResult.Success {
        return TestResult{
            ValidationResult: validationResult,
            TestCase:         test,
            FailedProperty: Stderr,
        }
    }

    validator := NewValidator(Equal)
    validationResult = validator.Validate(test.Result.ExitCode, test.Expected.ExitCode)
    if !validationResult.Success {
        return TestResult{
            ValidationResult: validationResult,
            TestCase: test,
            FailedProperty: ExitCode,
        }
    }

    return TestResult{
        ValidationResult: validationResult,
        TestCase:         test,
    }
}

func validateExpectedOut(got string, expected  ExpectedOut) ValidationResult {
    var v Validator
    var result ValidationResult

    if expected.Exactly != ""{
        v = NewValidator(Text)
        if result = v.Validate(got, expected.Exactly); !result.Success {
            return result
        }
    }

    if len(expected.Contains) > 0 {
        v = NewValidator(Contains)
        for _, c := range expected.Contains {
            if result = v.Validate(got, c); !result.Success {
                return result
            }
        }
    }

    result.Success = true
    return result
}