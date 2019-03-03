package runtime

import (
    "github.com/SimonBaeumer/commander/pkg/cmd"
    "github.com/SimonBaeumer/commander/pkg/suite"
    "log"
    "sync"
)

// CommandResult represents the TestCase and the ValidationResult
type TestResult struct {
    TestCase         suite.TestCase
    ValidationResult ValidationResult
}

// Start starts the given test suite
func Start(s suite.Suite) <-chan TestResult {
    in := make(chan suite.TestCase)
    out := make(chan TestResult)

     go func(tests []suite.TestCase) {
         defer close(in)
         for _, t := range tests {
             in <- t
         }
     }(s.Tests)

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

func runTest(test suite.TestCase) TestResult {
    // cut = command under test
    cut := cmd.NewCommand(test.Command.Cmd)

    if err := cut.Execute(); err != nil {
        log.Fatal(err)
    }

    // Write test result
    test.Result = suite.CommandResult{
        ExitCode: cut.ExitCode(),
        Stdout:   cut.Stdout(),
        Stderr:   cut.Stderr(),
    }

    validationResult := validateStdout(test)

    return TestResult{
        ValidationResult: validationResult,
        TestCase:         test,
    }
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