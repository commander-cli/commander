package runtime

import (
    "fmt"
    "github.com/SimonBaeumer/commander/pkg/cmd"
    "github.com/SimonBaeumer/commander/pkg/suite"
    "log"
    "os"
)

// Start starts the given test suite
func Start(s suite.Suite) {
    s.Start()

    c := make(chan suite.TestCase)

    for _, t := range s.Tests {
        go runTest(t, c)
    }

    counter := 0
    s.Success = true
    for r := range c {
        if r.Result.Status == suite.Failed {
            fmt.Println("Failed test " + r.Title)
            s.Success = false
        } else {
            fmt.Println("Success test " + r.Title)
        }

        counter++
        if counter >= len(s.Tests) {
            close(c)
        }
    }

    s.Finish()

    if !s.Success {
        log.Println("Suite failed, set error code to 1")
        os.Exit(1)
    }
}

func runTest(test suite.TestCase, results chan<- suite.TestCase) {
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

    result := Validate(test)
    if result.Success {
        test.Result.Status = suite.Success
    } else {
        test.Result.Status = suite.Failed
    }

    // Send to result channel
    results <- test
}
