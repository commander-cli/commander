package output

import (
    "fmt"
    "log"
)

type HumanOutput struct {
    buffer []string
}

func (o HumanOutput) BuildHeader() {
    s := "Starting test suite...\n"
    o.add(s)
}

func (o HumanOutput) BuildTestResult(test TestCase) string {
    var s string

    if test.Result.Success {
        s = fmt.Sprintf("✓ %s", test.Title)
    } else {
        s = fmt.Sprintf("✗ %s\n", test.Title)
        for _, p := range test.Result.FailureProperties {
            log.Printf("Printing property result '%s'", p)
            if p == "Stdout" {
                s +=  fmt.Sprintf("Got '%s', expected '%s'\n", test.Result.Stdout, test.Stdout)
            }
            if p == "Stderr" {
                s += fmt.Sprintf("Got %s, expected %s\n", test.Result.Stderr, test.Stderr)
            }
            if p == "ExitCode" {
                s += fmt.Sprintf("Got %d, expected %d\n", test.Result.ExitCode, test.ExitCode)
            }
        }
    }

    o.add(s)
    return s
}

func (o HumanOutput) BuildSuiteResult() {
    s := "Duration: 3.24sec"
    o.add(s)
}

func (o HumanOutput) add(out string) {
    o.buffer = append(o.buffer, out)
}

func (o HumanOutput) GetBuffer() []string {
    return o.buffer
}

func (o HumanOutput) Print() {
    for _, s := range o.GetBuffer() {
        fmt.Println(s)
    }
}
