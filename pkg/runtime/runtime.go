package runtime

import (
    "fmt"
    "github.com/SimonBaeumer/commander/pkg/cmd"
    "github.com/SimonBaeumer/commander/pkg/suite"
    "log"
    "time"
)

type Result struct {
    Success     bool
    Duration    time.Duration
    TestResults []TestResult
}

type TestResult struct {
    TestCase         suite.TestCase
    ValidationResult ValidationResult
}

// Start starts the given test suite
func Start(s suite.Suite) Result {
    testResults := []TestResult{}
    success     := true
    start       := time.Now()

    c := make(chan TestResult)
    for _, t := range s.Tests {
        go runTest(t, c)
    }

    counter := 0
    for r := range c {
        testResults = append(testResults, r)
        if r.ValidationResult.Success {
            fmt.Println("✓ ", r.TestCase.Title)
        }

        if !r.ValidationResult.Success {
            success = false
            fmt.Println("✗ ", r.TestCase.Title)
        }

        counter++
        if counter >= len(s.Tests) {
            close(c)
        }
    }

    return Result{
        Success:     success,
        Duration:    time.Since(start),
        TestResults: testResults,
    }
}

func runTest(test suite.TestCase, results chan<- TestResult) {
    // cut = command under test
    cut := cmd.NewCommand(test.Command.Cmd)

    if err := cut.Execute(); err != nil {
        log.Fatal(err)
    }

    // Write test result
    test.Result = suite.TestResult{
        ExitCode: cut.ExitCode(),
        Stdout:   cut.Stdout(),
        Stderr:   cut.Stderr(),
    }

    validationResult := validateStdout(test)

    result := TestResult{
        ValidationResult: validationResult,
        TestCase:         test,
    }

    // Send to result channel
    results <- result
}

func validateStdout(test suite.TestCase) ValidationResult {
    var v Validator
    var result ValidationResult

    if test.Expected.Stdout.Exactly != ""{
        v = NewValidator(Text)
        if result = v.Validate(test.Result.Stdout, test.Expected.Stdout.Exactly); !result.Success {
            return result
        }
    }

    if len(test.Expected.Stdout.Contains) > 0 {
        v = NewValidator(Contains)
        for _, c := range test.Expected.Stdout.Contains {
            if result = v.Validate(test.Result.Stdout, c); !result.Success {
                return result
            }
        }
    }

    result.Success = true
    return result
}