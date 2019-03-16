package output

import (
    "bytes"
    "github.com/SimonBaeumer/commander/pkg/runtime"
    "github.com/stretchr/testify/assert"
    "strings"
    "sync"
    "testing"
)

func Test_Start(t *testing.T) {
    var buf bytes.Buffer
    var wg sync.WaitGroup
    results := make(chan runtime.TestResult)

    wg.Add(1)
    go func() {
        defer wg.Done()

        writer := OutputWriter{out: &buf}
        got := writer.Start(results)

        assert.False(t, got)
    }()

    results <- runtime.TestResult{
        TestCase: runtime.TestCase{
            Title: "Successful test",
            Command: runtime.CommandUnderTest{},
        },
        ValidationResult: runtime.ValidationResult{
            Success: true,
        },
        FailedProperty: "",
    }

    results <- runtime.TestResult{
        TestCase: runtime.TestCase{
            Title: "Failed test",
            Command: runtime.CommandUnderTest{},
        },
        ValidationResult: runtime.ValidationResult{
            Success: false,
        },
        FailedProperty: "Stdout",
    }

    close(results)
    wg.Wait()
    assert.True(t, true)
    assert.True(t, strings.Contains(buf.String(), "✓ Successful test"))
    assert.True(t, strings.Contains(buf.String(), "✗ Failed test"))
}