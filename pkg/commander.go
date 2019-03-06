package commander

import (
    "fmt"
    "time"

    "github.com/SimonBaeumer/commander/pkg/runtime"
    "github.com/logrusorgru/aurora"
)

var au aurora.Aurora

func Start(results <-chan runtime.TestResult) bool {
    au = aurora.NewAurora(true)
    failed := 0
    testResults := []runtime.TestResult{}
    start := time.Now()

    for r := range results {
        testResults = append(testResults, r)
        if r.ValidationResult.Success {
            fmt.Println("✓ " + r.TestCase.Title)
        }

        if !r.ValidationResult.Success {
            failed++
            fmt.Println(au.Red("✗ " + r.TestCase.Title))
        }
    }

    duration := time.Since(start)

    if failed > 0 {
        printFailures(testResults)
    }

    fmt.Println("")
    fmt.Println(fmt.Sprintf("Duration: %.3fs", duration.Seconds()))
    summary := fmt.Sprintf("Count: %d, Failed: %d", len(testResults), failed)
    if failed > 0 {
        fmt.Println(au.Red(summary))
    } else {
        fmt.Println(au.Green(summary))
    }

    return failed == 0
}

func printFailures(results []runtime.TestResult) {
    fmt.Println("")
    fmt.Println(au.Bold("Results"))
    fmt.Println(au.Bold(""))

    for _, r := range results {
        if !r.ValidationResult.Success {
            fmt.Println(au.Bold(au.Red("✗ " + r.TestCase.Title + ", on property '" + r.FailedProperty + "'")))
            fmt.Println(r.ValidationResult.Diff)
        }
    }

}
